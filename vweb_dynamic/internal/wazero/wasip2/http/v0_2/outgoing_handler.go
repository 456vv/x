package v0_2

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	gohttp "net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	manager_http "github.com/456vv/x/vweb_dynamic/internal/wazero/manager/http"
	manager_io "github.com/456vv/x/vweb_dynamic/internal/wazero/manager/io"
	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"

	lru "github.com/hashicorp/golang-lru/v2"
)

type timeoutConfig struct {
	connect      time.Duration
	firstByte    time.Duration
	betweenBytes time.Duration
}

// outgoingHandlerImpl 封装了 wasi:http/outgoing-handler 的所有操作。
type outgoingHandlerImpl struct {
	hm              *manager_http.HTTPManager
	clientCache     *lru.Cache[timeoutConfig, *gohttp.Client]
	clientMu        sync.Mutex
	fallbackClients map[timeoutConfig]*gohttp.Client // 修改原因：lru 创建失败时仍按超时配置复用 Client，避免每次 Handle 泄漏 Transport
}

func newOutgoingHandlerImpl(hm *manager_http.HTTPManager) *outgoingHandlerImpl {
	impl := &outgoingHandlerImpl{hm: hm}
	clientCache, err := lru.NewWithEvict[timeoutConfig, *gohttp.Client](32, func(_ timeoutConfig, c *gohttp.Client) {
		if c != nil {
			hm.UntrackClient(c)
			c.CloseIdleConnections()
		}
	})
	if err == nil {
		impl.clientCache = clientCache
	} else {
		impl.fallbackClients = make(map[timeoutConfig]*gohttp.Client)
	}
	return impl
}

func (i *outgoingHandlerImpl) getClient(opts *manager_http.RequestOptions) *gohttp.Client {
	var cfg timeoutConfig
	if opts != nil {
		c, f, b := opts.CopyTimeouts()
		if c != nil {
			cfg.connect = *c
		}
		if f != nil {
			cfg.firstByte = *f
		}
		if b != nil {
			cfg.betweenBytes = *b
		}
	}

	i.clientMu.Lock()
	defer i.clientMu.Unlock()

	if i.clientCache != nil {
		if client, ok := i.clientCache.Get(cfg); ok {
			return client
		}
	} else if i.fallbackClients != nil {
		if client, ok := i.fallbackClients[cfg]; ok {
			return client
		}
	}

	transport, ok := gohttp.DefaultTransport.(*gohttp.Transport)
	if ok {
		transport = transport.Clone()
	} else {
		transport = &gohttp.Transport{Proxy: gohttp.ProxyFromEnvironment}
	}

	transport.Proxy = gohttp.ProxyFromEnvironment
	if transport.TLSClientConfig == nil {
		transport.TLSClientConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	} else {
		transport.TLSClientConfig = transport.TLSClientConfig.Clone()
		transport.TLSClientConfig.InsecureSkipVerify = false
		if transport.TLSClientConfig.MinVersion == 0 {
			transport.TLSClientConfig.MinVersion = tls.VersionTLS12
		}
	}
	transport.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
		c, err := (&net.Dialer{Timeout: cfg.connect, KeepAlive: 30 * time.Second}).DialContext(ctx, network, address)
		if err != nil {
			return nil, err
		}
		if cfg.betweenBytes > 0 {
			c = &betweenBytesConn{Conn: c, idle: cfg.betweenBytes}
		}
		return c, nil
	}
	transport.ResponseHeaderTimeout = cfg.firstByte
	client := &gohttp.Client{
		Transport:     transport,
		Timeout:       0,
		Jar:           nil,
		CheckRedirect: func(req *gohttp.Request, via []*gohttp.Request) error { return gohttp.ErrUseLastResponse },
	}

	if i.clientCache != nil {
		i.clientCache.Add(cfg, client)
	} else {
		if i.fallbackClients == nil {
			i.fallbackClients = make(map[timeoutConfig]*gohttp.Client)
		}
		// 修改原因：fallback map 无 LRU；超过 32 淘汰一条，避免 timeout 组合把 Transport 撑爆
		if len(i.fallbackClients) >= 32 {
			for k, c := range i.fallbackClients {
				delete(i.fallbackClients, k)
				i.hm.UntrackClient(c)
				if c != nil {
					c.CloseIdleConnections()
				}
				break
			}
		}
		i.fallbackClients[cfg] = client
	}
	i.hm.TrackClient(client)
	return client
}

// Handle 实现了 outgoing-handler.handle 接口。
// 这是执行 HTTP 请求的核心。
// 调用后会消耗掉 request 和 options（WIT 按值移动，Ok/Err 都不再归 guest）。
func (i *outgoingHandlerImpl) Handle(
	request OutgoingRequest,
	options witgo.Option[RequestOptions],
) witgo.Result[FutureIncomingResponse, ErrorCode] {
	req, ok := i.hm.OutgoingRequests.Pop(request)
	if !ok {
		if options.IsSome() && options.Some != nil {
			i.hm.Options.Remove(*options.Some)
		}
		return witgo.Err[FutureIncomingResponse, ErrorCode](ErrorCode{InternalError: witgo.SomePtr("invalid request handle")})
	}

	var opts *manager_http.RequestOptions
	if options.IsSome() && options.Some != nil {
		opts, _ = i.hm.Options.Pop(*options.Some)
	}

	goReq, err := i.buildGoRequest(req)
	if err != nil {
		if h, hdr := req.TakeHeadersHandleForDrop(); h != 0 {
			i.hm.Fields.RemoveIf(h, func(cur manager_http.Fields) bool {
				return manager_http.SameFields(cur, hdr)
			})
		}
		// Do 没跑，Body 没有消费者。不 Remove 会泄漏 Pipe 和 outgoing-body。
		if req.BodyHandle != 0 {
			pw := req.BodyWriter
			i.hm.Bodies.RemoveIf(req.BodyHandle, func(cur *manager_http.OutgoingBody) bool {
				return cur != nil && pw != nil && cur.BodyWriter == pw
			})
			req.BodyHandle = 0
		}
		req.Close()
		return witgo.Err[FutureIncomingResponse, ErrorCode](ErrorCode{InternalError: witgo.SomePtr(err.Error())})
	}

	// —— 成功路径 —— headers 已拷进 goReq，按身份回收 Fields。
	// Body 仍由 guest 的 finish()/resource-drop 负责。
	if h, hdr := req.TakeHeadersHandleForDrop(); h != 0 {
		i.hm.Fields.RemoveIf(h, func(cur manager_http.Fields) bool {
			return manager_http.SameFields(cur, hdr)
		})
	}
	// 禁止 i.hm.Bodies.Remove(req.BodyHandle)：
	// goReq.Body 就是 body() 创建的 PipeReader；executeRequest → Client.Do 还在读。
	// Remove 会走 OutgoingBody 析构，CloseWithError(ErrUnexpectedEOF)，把未写完的请求截断。
	// outgoing-body 仍由 guest 的 finish()/resource-drop 负责回收。

	client := i.getClient(opts)

	ctx, cancel := context.WithCancel(context.Background())
	future := &manager_http.FutureIncomingResponse{
		Pollable: manager_io.NewPollable(nil),
		Cancel:   cancel, // drop future 时取消 in-flight Client.Do
	}
	goReq = goReq.WithContext(ctx)

	// 启动一个新的 goroutine 来异步执行 HTTP 请求。
	go i.executeRequest(client, goReq, future)

	// 立即返回 future 句柄，不阻塞。
	futureHandle := i.hm.Futures.Add(future)
	return witgo.Ok[FutureIncomingResponse, ErrorCode](futureHandle)
}

// executeRequest 在一个单独的 goroutine 中运行。
func (i *outgoingHandlerImpl) executeRequest(client *gohttp.Client, goReq *gohttp.Request, future *manager_http.FutureIncomingResponse) {
	defer future.Pollable.SetReady()
	defer func() {
		if rec := recover(); rec != nil {
			// Do/回调 panic 时仍须 StoreResult，否则 future.get/drop 会死等
			future.StoreResult(manager_http.Result{Err: fmt.Errorf("http roundtrip panic: %v", rec)})
		}
	}()

	// drop future 只 Cancel context；Client.Do 在部分路径仍等待 Body EOF，
	// 与 Futures destructor 里 Pollable.Block() 死锁。取消时关掉 PipeReader 唤醒 Do。
	reqDone := make(chan struct{})
	defer close(reqDone)
	if goReq.Body != nil {
		body := goReq.Body
		go func() {
			select {
			case <-goReq.Context().Done():
				_ = body.Close()
			case <-reqDone:
			}
		}()
	}

	resp, err := client.Do(goReq)
	if err != nil {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close() // Do 在部分错误路径仍返回非 nil Response
		}
		future.StoreResult(manager_http.Result{Err: err})
		return
	}
	future.StoreResult(manager_http.Result{Response: resp})
}

// buildGoRequest 是一个辅助函数，用于将 wasi-http 请求转换为 Go 的 http.Request。
func (i *outgoingHandlerImpl) buildGoRequest(req *manager_http.OutgoingRequest) (*gohttp.Request, error) {
	method, path, schemePtr, authority := req.LoadRequestLine()
	if authority == nil || !validRequestAuthority(*authority) {
		return nil, fmt.Errorf("request authority is invalid")
	}
	if method == "" {
		method = gohttp.MethodGet
	}
	if !validHTTPToken(method) {
		return nil, fmt.Errorf("invalid method")
	}

	scheme := "https"
	if schemePtr != nil {
		switch strings.ToLower(strings.TrimSpace(*schemePtr)) {
		case "", "https":
			scheme = "https"
		case "http":
			scheme = "http"
		default:
			return nil, fmt.Errorf("unsupported scheme %q", *schemePtr)
		}
	}

	if path == "" {
		path = "/"
	} else if !validRequestPath(path) {
		// 不再把缺少前导 / 的字符串偷偷改成路径，避免和 set-path 的失败语义不一致。
		return nil, fmt.Errorf("invalid path")
	}
	rawURL := fmt.Sprintf("%s://%s%s", scheme, *authority, path)

	// 创建 Go 的 http.Request。req.Body 是一个 io.PipeReader，
	// 当 Guest 向 outgoing-body 写入数据时，这里就能读到。
	goReq, err := gohttp.NewRequest(method, rawURL, req.Body)
	if err != nil {
		return nil, err
	}

	i.hm.LockFields()
	hdrs := cloneFieldsLower(req.Headers)
	i.hm.UnlockFields()
	for k, vv := range hdrs {
		for _, v := range vv {
			goReq.Header.Add(k, v)
		}
	}
	goReq.Trailer = make(gohttp.Header)
	if v := goReq.Header.Get("Content-Length"); v != "" {
		if n, perr := strconv.ParseInt(v, 10, 64); perr == nil && n >= 0 {
			//  NewRequest 对 PipeReader 把 ContentLength 设为 -1，会改成 chunked，与 finish 的长度校验不一致。
			goReq.ContentLength = n
		}
	}

	req.BindRequest(goReq)
	return goReq, nil
}

type betweenBytesConn struct {
	net.Conn
	idle time.Duration
}

func (c *betweenBytesConn) Read(b []byte) (int, error) {
	if c.idle > 0 {
		_ = c.Conn.SetReadDeadline(time.Now().Add(c.idle))
	}
	return c.Conn.Read(b)
}
