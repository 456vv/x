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
	// 同一 wazero Module 不是并发安全的；只在 Call 期间互斥，不要覆盖后续等待。
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

	req, reqHandle, err := s.createIncomingRequest(ctx, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create incoming-request: %v", err), http.StatusInternalServerError)
		return
	}

	respChan := make(chan any, 1)
	out := &manager_http.ResponseOutparam{ResultChan: respChan}
	outparamHandle := hm.ResponseOutparams.Add(out)

	// WASI 规范路径是先 response-outparam.set 再写 body。
	// io.Pipe 无缓冲，必须在 handle() 返回前就开始 Copy，否则 guest 写体会阻塞 wasm 线程。
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
		// 不能和 ctx.Done() 一起 select。两者同时就绪时可能丢掉已经 set 的响应，主流程再误报 503。
		// 一直等到 set 投递，或 Remove/Drop 关闭 channel。
		respResult, ok := <-respChan
		if ok {
			startWrite(respResult)
		}
	}()

	// Call panic 时原先没有 Unlock，后续请求会全部堵在 guestMu 上。
	// 锁只包住 Call，不能盖住后面的 setDone/copyDone，否则响应写完之前别的请求进不了 wasm。
	callErr := func() (callErr error) {
		s.guestMu.Lock()
		defer s.guestMu.Unlock()
		defer func() {
			if rec := recover(); rec != nil {
				callErr = fmt.Errorf("guest handle panic: %v", rec)
			}
		}()
		_, callErr = s.handleFunc.Call(ctx, uint64(reqHandle), uint64(outparamHandle))
		return callErr
	}()
	if callErr != nil {
		// trap 时 guest 可能已经 drop；句柄号会被复用。只移除仍是本次请求的那一项。
		hm.IncomingRequests.RemoveIf(reqHandle, func(cur *manager_http.IncomingRequest) bool {
			return cur == req
		})
	}

	// 无条件 Remove(outparamHandle) 会在句柄复用后关掉另一个请求的 ResultChan。
	hm.ResponseOutparams.RemoveIf(outparamHandle, func(cur *manager_http.ResponseOutparam) bool {
		return cur == out
	})
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

// createIncomingRequest 将 http.Request 转换为 wasi:http/types.incoming-request 资源。
// 返回对象指针，供 trap 路径按身份 RemoveIf，避免句柄复用误删。
func (s *Server) createIncomingRequest(ctx context.Context, r *http.Request) (*manager_http.IncomingRequest, v0_2.IncomingRequest, error) {
	_ = ctx
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
	path := r.URL.EscapedPath()
	if path == "" {
		// 星号等请求没有转义路径时退回解码路径，避免得到空 path。
		path = r.URL.Path
	}
	var authority *string
	if r.Host != "" {
		host := r.Host
		authority = &host
	}
	req := &manager_http.IncomingRequest{
		Request:     r,
		Method:      r.Method,
		Path:        path,
		Query:       r.URL.RawQuery,
		Scheme:      &scheme,
		Authority:   authority,
		Headers:     headersHandle,
		HeadersData: headers, // drop 时用 map 身份 RemoveIf，避免句柄复用误删
		Body:        r.Body,
	}
	return req, hm.IncomingRequests.Add(req), nil
}

// writeOutgoingResponse 将 guest 返回的 OutgoingResponse 写入 http.ResponseWriter。
func (s *Server) writeOutgoingResponse(ctx context.Context, w http.ResponseWriter, respHandle v0_2.OutgoingResponse) {
	hm := s.wasiHost.HTTPManager()
	resp, ok := hm.OutgoingResponses.Pop(respHandle)
	if !ok {
		http.Error(w, "internal error: invalid outgoing-response handle", http.StatusInternalServerError)
		return
	}
	defer func() {
		// Pop 父资源不会清 child；guest 已 drop/finish 后编号会复用。
		if h, hdr := resp.TakeHeadersHandleForDrop(); h != 0 {
			hm.Fields.RemoveIf(h, func(cur manager_http.Fields) bool {
				return manager_http.SameFields(cur, hdr)
			})
		}

		if h := resp.TakeBodyHandleForDrop(); h != 0 {
			pw := resp.BodyWriter
			hm.Bodies.RemoveIf(h, func(cur *manager_http.OutgoingBody) bool {
				return cur != nil && pw != nil && cur.BodyWriter == pw
			})
		}
	}()
	resp.Response = w
	if resp.Headers != nil {
		hm.LockFields()
		hdrs := make(manager_http.Fields, len(resp.Headers))
		for k, vv := range resp.Headers {
			cp := make([]string, len(vv))
			copy(cp, vv)
			hdrs[k] = cp
		}
		hm.UnlockFields()
		for k, vv := range hdrs {
			for _, v := range vv {
				w.Header().Add(k, v)
			}
		}
	}

	var bodyBuf []byte
	if resp.Body != nil {
		errCh := make(chan error, 1)
		go func() {
			var copyErr error
			bodyBuf, copyErr = io.ReadAll(resp.Body)
			errCh <- copyErr
		}()
		var copyErr error
		aborted := false
		select {
		case copyErr = <-errCh:
		case <-ctx.Done():
			aborted = true
			if c, ok := resp.Body.(io.Closer); ok {
				c.Close()
			}
			copyErr = <-errCh
		}
		if c, ok := resp.Body.(io.Closer); ok {
			c.Close()
		}
		// 必须先读完 body（finish 在 EOF 时写入 trailer），再 TakeTrailers + WriteHeader。
		// 上一轮流式 Copy 会在 trailer 到达前 WriteHeader，trailer 丢失。
		if aborted || copyErr != nil {
			http.Error(w, "failed to read guest response body", http.StatusBadGateway)
			return
		}
	}

	for k, vv := range resp.TakeTrailers() {
		for _, v := range vv {
			w.Header().Add(http.TrailerPrefix+k, v)
		}
	}
	status := resp.LoadStatus()
	if status < 100 || status > 599 {
		status = http.StatusOK
	}
	w.WriteHeader(status)
	if len(bodyBuf) > 0 {
		w.Write(bodyBuf)
	}
}
