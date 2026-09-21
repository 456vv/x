package wasi_http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"

	manager_http "github.com/456vv/x/vweb_dynamic/internal/wazero/manager/http"
	"github.com/456vv/x/vweb_dynamic/internal/wazero/wasip2"
	v0_2 "github.com/456vv/x/vweb_dynamic/internal/wazero/wasip2/http/v0_2"
	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"

	"github.com/tetratelabs/wazero/api"
)

// Server 实现了 http.Handler，将传入的 HTTP 请求转发给 wasi-http guest 模块处理。
type Server struct {
	guest      api.Module
	wasiHost   *wasip2.Host
	witHost    *witgo.Host
	handleFunc api.Function
	// 同一 wazero Module 不是并发安全的；只在 Call 期间互斥，不要覆盖 setDone/copy
	guestMu sync.Mutex
}

// NewServer 创建一个新的 wasi-http 服务器实例。
// guest 模块必须导出一个 `wasi:http/incoming-handler.handle` 函数。
func NewServer(guest api.Module, wasiHost *wasip2.Host) (*Server, error) {
	handleFunc := guest.ExportedFunction("wasi:http/incoming-handler#handle")
	if handleFunc == nil {
		return nil, fmt.Errorf("guest module must export wasi:http/incoming-handler#handle function")
	}

	witHost, err := witgo.NewHost(guest)
	if err != nil {
		return nil, fmt.Errorf("failed to create wit-go host: %w", err)
	}

	return &Server{
		guest:      guest,
		wasiHost:   wasiHost,
		witHost:    witHost,
		handleFunc: handleFunc,
	}, nil
}

// ServeHTTP 是 http.Handler 接口的实现。
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	hm := s.wasiHost.HTTPManager()

	reqHandle, err := s.createIncomingRequest(ctx, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create incoming-request: %v", err), http.StatusInternalServerError)
		return
	}

	respChan := make(chan any, 1)
	outparamHandle := hm.ResponseOutparams.Add(&manager_http.ResponseOutparam{ResultChan: respChan})

	// WASI 规范路径是先 response-outparam.set 再写 body。
	// io.Pipe 无缓冲，必须在 handle() 返回前就开始 Copy，否则 guest 写体会阻塞 wasm 线程。
	// 先写后 set 仍可能死锁，那是 guest 违反规范。
	var started atomic.Bool
	copyDone := make(chan struct{})
	startWrite := func(v any) {
		if !started.CompareAndSwap(false, true) {
			return
		}
		go func() {
			defer close(copyDone)
			switch result := v.(type) {
			case v0_2.OutgoingResponse:
				s.writeOutgoingResponse(ctx, w, result)
			case v0_2.ErrorCode:
				http.Error(w, fmt.Sprintf("guest returned an error code: %+v", result), http.StatusInternalServerError)
			default:
				http.Error(w, "internal error: unknown type from response channel", http.StatusInternalServerError)
			}
		}()
	}

	setDone := make(chan struct{})
	go func() {
		defer close(setDone)
		select {
		case respResult, ok := <-respChan:
			if ok {
				startWrite(respResult)
			}
		case <-ctx.Done():
			// 不在这里 http.Error。若随后已经 set，由 startWrite 写响应。
		}
	}()

	s.guestMu.Lock()
	_, callErr := s.handleFunc.Call(ctx, uint64(reqHandle), uint64(outparamHandle))
	s.guestMu.Unlock()
	if callErr != nil {
		// handle 成功后 incoming-request 所有权在 guest，不能 Remove；
		// 仅 trap 时 guest 未 drop，才由宿主回收，避免拆掉仍在读的 body。
		hm.IncomingRequests.Remove(reqHandle)
	}

	// 未 set 时关闭 ResultChan，让上面的接收结束
	hm.ResponseOutparams.Remove(outparamHandle)
	<-setDone

	if !started.Load() {
		switch {
		case ctx.Err() != nil:
			http.Error(w, "context cancelled", http.StatusServiceUnavailable)
		case callErr != nil:
			http.Error(w, fmt.Sprintf("guest handle function trapped: %v", callErr), http.StatusInternalServerError)
		default:
			http.Error(w, "guest did not set a response", http.StatusInternalServerError)
		}
		return
	}
	<-copyDone
}

// createIncomingRequest 将 http.Request 转换为 wasi:http/types.incoming-request 资源
func (s *Server) createIncomingRequest(ctx context.Context, r *http.Request) (v0_2.IncomingRequest, error) {
	hm := s.wasiHost.HTTPManager()

	headers := make(manager_http.Fields, len(r.Header))
	for k, v := range r.Header {
		// 直接赋值会与 net/http 共享底层 []string
		cp := make([]string, len(v))
		copy(cp, v)
		headers[strings.ToLower(k)] = cp
	}

	headersHandle := hm.Fields.Add(headers)
	hm.MarkFieldsImmutable(headersHandle) // incoming-request.headers 按规范不可变

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	req := &manager_http.IncomingRequest{
		Request:   r,
		Method:    r.Method,
		Path:      r.URL.Path,
		Query:     r.URL.RawQuery,
		Scheme:    &scheme,
		Authority: &r.Host,
		Headers:   headersHandle,
		Body:      r.Body,
	}
	return hm.IncomingRequests.Add(req), nil
}

// writeOutgoingResponse 将 guest 返回的 OutgoingResponse 写入 http.ResponseWriter
func (s *Server) writeOutgoingResponse(ctx context.Context, w http.ResponseWriter, respHandle v0_2.OutgoingResponse) {
	hm := s.wasiHost.HTTPManager()
	resp, ok := hm.OutgoingResponses.Pop(respHandle)
	if !ok {
		http.Error(w, "internal error: invalid outgoing-response handle", http.StatusInternalServerError)
		return
	}
	// Pop 不跑 destructor，必须回收 headers/body 子句柄，否则长驻 server 泄漏
	defer func() {
		if resp.HeadersHandle != 0 {
			hm.Fields.Remove(resp.HeadersHandle)
			resp.HeadersHandle = 0
		}
		if resp.BodyHandle != 0 {
			hm.Bodies.Remove(resp.BodyHandle)
			resp.BodyHandle = 0
		}
	}()
	resp.Response = w
	if resp.Headers != nil {
		for k, vv := range resp.Headers {
			for _, v := range vv {
				w.Header().Add(k, v) // 规范化
			}
		}
	}
	status := resp.StatusCode
	if status < 100 || status > 999 {
		status = http.StatusOK
	}
	w.WriteHeader(status)
	if resp.Body != nil {
		// 客户端取消时 io.Copy 会一直等到 guest finish；关 Body 打断 Pipe
		errCh := make(chan error, 1)
		go func() {
			_, copyErr := io.Copy(w, resp.Body)
			errCh <- copyErr
		}()
		select {
		case <-errCh:
		case <-ctx.Done():
			if c, ok := resp.Body.(io.Closer); ok {
				_ = c.Close()
			}
			<-errCh
		}
		if c, ok := resp.Body.(io.Closer); ok {
			_ = c.Close()
		}
	}
}

// Close 释放与 guest 模块绑定的 witgo.Host 缓存，避免模块销毁后分配器泄漏。
func (s *Server) Close() {
	if s == nil || s.guest == nil {
		return
	}
	witgo.ReleaseHost(s.guest)
}
