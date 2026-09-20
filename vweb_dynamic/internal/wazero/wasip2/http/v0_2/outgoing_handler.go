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
	hm          *manager_http.HTTPManager
	clientCache *lru.Cache[timeoutConfig, *gohttp.Client]
	clientMu    sync.Mutex // hashicorp/lru 默认非并发安全
}

func newOutgoingHandlerImpl(hm *manager_http.HTTPManager) *outgoingHandlerImpl {
	impl := &outgoingHandlerImpl{hm: hm}
	clientCache, _ := lru.NewWithEvict[timeoutConfig, *gohttp.Client](32, func(_ timeoutConfig, c *gohttp.Client) {
		// 如果 LRU 淘汰不 Untrack 会导致 Transport 只增不减
		if c != nil {
			hm.UntrackClient(c)
			c.CloseIdleConnections()
		}
	})
	impl.clientCache = clientCache
	return impl
}

func (i *outgoingHandlerImpl) getClient(opts *manager_http.RequestOptions) *gohttp.Client {
	var cfg timeoutConfig
	if opts != nil {
		if opts.ConnectTimeout != nil {
			cfg.connect = *opts.ConnectTimeout
		}
		if opts.FirstByteTimeout != nil {
			cfg.firstByte = *opts.FirstByteTimeout
		}
		if opts.BetweenBytesTimeout != nil {
			cfg.betweenBytes = *opts.BetweenBytesTimeout
		}
	}

	i.clientMu.Lock()
	defer i.clientMu.Unlock()

	if client, ok := i.clientCache.Get(cfg); ok {
		return client
	}

	transport, ok := gohttp.DefaultTransport.(*gohttp.Transport)
	if ok {
		transport = transport.Clone()
	} else {
		transport = &gohttp.Transport{Proxy: gohttp.ProxyFromEnvironment}
	}

	// 保留原始代码中的代理和 TLS 设置
	// p, _ := url.Parse("http://192.168.3.121:8888")
	// transport.Proxy = gohttp.ProxyURL(p)
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

	i.clientCache.Add(cfg, client)
	i.hm.TrackClient(client) // Host.Close 需要 CloseIdleConnections
	return client
}

// Handle 实现了 outgoing-handler.handle 接口。
// 这是执行 HTTP 请求的核心。
// 调用后会消耗掉request和options
func (i *outgoingHandlerImpl) Handle(
	request OutgoingRequest,
	options witgo.Option[RequestOptions], // options 是可选的
) witgo.Result[FutureIncomingResponse, ErrorCode] {
	var opts *manager_http.RequestOptions
	var optHandle uint32
	if options.IsSome() && options.Some != nil {
		optHandle = *options.Some
		opts, _ = i.hm.Options.Pop(optHandle)
	}

	req, ok := i.hm.OutgoingRequests.Pop(request)
	if !ok {
		// 已 Pop 的 options 必须塞回，否则 Handle 失败泄漏
		if optHandle != 0 && opts != nil {
			i.hm.Options.Set(optHandle, opts)
		}
		return witgo.Err[FutureIncomingResponse, ErrorCode](ErrorCode{InternalError: witgo.SomePtr("invalid request handle")})
	}

	goReq, err := i.buildGoRequest(req)
	if err != nil {
		// Pop 不再调用 destructor，失败时必须回收 headers/body 子句柄
		if req.HeadersHandle != 0 {
			i.hm.Fields.Remove(req.HeadersHandle)
			req.HeadersHandle = 0
		}
		if req.BodyHandle != 0 {
			i.hm.Bodies.Remove(req.BodyHandle)
			req.BodyHandle = 0
		}
		req.Close()
		return witgo.Err[FutureIncomingResponse, ErrorCode](ErrorCode{InternalError: witgo.SomePtr(err.Error())})
	}

	// Pop 不跑 destructor，headers() 分配的 Fields 子句柄会一直留在 Fields 表里。
	// 头部已拷进 goReq.Header，可以立刻 Remove。
	// 不能 Bodies.Remove(BodyHandle)：Client.Do 仍通过 req.Body（PipeReader）读 guest 写入的 body。
	if req.HeadersHandle != 0 {
		i.hm.Fields.Remove(req.HeadersHandle)
		req.HeadersHandle = 0
	}

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
	// 构造 URL
	scheme := "https"
	if req.Scheme != nil && *req.Scheme == "http" {
		scheme = "http"
	}
	if req.Authority == nil {
		return nil, fmt.Errorf("request authority cannot be empty")
	}
	path := req.Path
	if path == "" {
		path = "/"
	} else if !strings.HasPrefix(path, "/") {
		// 缺 leading slash 时会拼成 https://hostfoo 而不是 https://host/foo
		path = "/" + path
	}
	url := fmt.Sprintf("%s://%s%s", scheme, *req.Authority, path)

	// 创建 Go 的 http.Request。req.Body 是一个 io.PipeReader，
	// 当 Guest 向 outgoing-body 写入数据时，这里就能读到。
	goReq, err := gohttp.NewRequest(req.Method, url, req.Body)
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

	req.Request = goReq
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
