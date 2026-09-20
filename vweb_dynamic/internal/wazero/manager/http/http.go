package http

import (
	"context"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	manager_io "github.com/456vv/x/vweb_dynamic/internal/wazero/manager/io"
	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"
)

// Fields 代表 HTTP 头部或尾部。
type Fields = http.Header

// IncomingRequest 代表一个由 Host 接收的、传递给 Guest 的 HTTP 请求的内部表示。
type IncomingRequest struct {
	Request *http.Request

	Method    string
	Path      string
	Query     string
	Scheme    *string
	Authority *string
	Headers   uint32

	Body io.ReadCloser

	// BodyHandle 用于 Guest 端消费 Body
	BodyHandle uint32

	// 与 IncomingResponse 对齐，保证 consume 至多一次（并发安全）
	Consumed atomic.Bool
}

// OutgoingRequest 代表一个由 guest 构建的出站 HTTP 请求。
type OutgoingRequest struct {
	// 在handle调用后设置
	Request *http.Request

	Method    string
	Scheme    *string
	Authority *string
	Path      string
	Headers   Fields

	// headers() 必须每次返回同一 child handle，且不能每次 Add 共享 map
	HeadersHandle uint32
	headersOnce   sync.Once // 并发 headers() 会重复 Add；Once 必须留在本包，外包不能碰未导出字段

	Body io.Reader

	// 这里只是引用资源，不属于OutgoingRequest生命周期管理
	// BodyWriter 用于在 Host 端写入 Guest 提供的数据
	BodyWriter *io.PipeWriter
	BodyHandle uint32 // 指向 outgoing-body 资源的句柄

	// 消耗标记
	Consumed atomic.Bool
}

func (o *OutgoingRequest) Close() error {
	if o == nil {
		return nil
	}
	if o.Body != nil {
		if closer, ok := o.Body.(io.Closer); ok {
			return closer.Close()
		}
	}
	return nil
}

// EnsureHeadersHandle 保证 headers() 子句柄只分配一次（跨包调用 Once）。
func (o *OutgoingRequest) EnsureHeadersHandle(init func()) {
	if o == nil || init == nil {
		return
	}
	o.headersOnce.Do(init)
}

// IncomingResponse 代表一个已到达的、由 Host 接收的 HTTP 响应。
type IncomingResponse struct {
	Response *http.Response

	StatusCode int
	Headers    Fields

	// WASI incoming-response.headers 每次返回同一不可变快照句柄
	HeaderHandle uint32
	headersOnce  sync.Once // 并发 headers() 会重复 Clone+Add

	Body       *IncomingBody
	BodyHandle uint32 // 指向 incoming-body 的句柄
	// 消耗标记
	Consumed atomic.Bool
}

// EnsureHeaderHandle 保证 incoming-response.headers 只创建一次快照句柄。
func (r *IncomingResponse) EnsureHeaderHandle(init func()) {
	if r == nil || init == nil {
		return
	}
	r.headersOnce.Do(init)
}

// OutgoingResponse 代表一个由 Guest 构建的出站 HTTP 响应。
type OutgoingResponse struct {
	Response http.ResponseWriter

	StatusCode int
	Headers    Fields

	// 与 OutgoingRequest 相同，缓存 headers() 子资源
	HeadersHandle uint32
	headersOnce   sync.Once // 并发 headers() 会重复 Add

	Body io.Reader

	// 这里只是引用资源，不属于OutgoingResponse生命周期管理
	BodyWriter *io.PipeWriter
	BodyHandle uint32

	// 消耗标记
	Consumed atomic.Bool
}

func (o *OutgoingResponse) Close() error {
	if o == nil {
		return nil
	}
	if o.Body != nil {
		if closer, ok := o.Body.(io.Closer); ok {
			return closer.Close()
		}
	}
	return nil
}

// EnsureHeadersHandle 保证 outgoing-response.headers 子句柄只分配一次。
func (o *OutgoingResponse) EnsureHeadersHandle(init func()) {
	if o == nil || init == nil {
		return
	}
	o.headersOnce.Do(init)
}

// ResponseOutparam 是一个一次性的句柄，用于让 Guest 设置对 IncomingRequest 的响应。
type ResponseOutparam struct {
	// 当 Guest 调用 response-outparam.set 时，结果会通过这个 channel 发送。
	// Result 包含一个 OutgoingResponse 的句柄或一个 ErrorCode。
	ResultChan chan<- any
}

// IncomingBody 代表一个入站的 HTTP Body。
type IncomingBody struct {
	// 因为go http 的限制Body和Stream 生命周期统一管理
	StreamHandle uint32 // 指向 input-stream 的句柄
	Stream       io.Reader

	// 可选方法
	GetTrailers func() (trailers Fields)

	// 消耗标记
	Consumed atomic.Bool
}

func (o *IncomingBody) Close() error {
	if o == nil {
		return nil
	}
	if o.Stream != nil {
		if closer, ok := o.Stream.(io.Closer); ok {
			return closer.Close()
		}
	}
	return nil
}

// OutgoingBody 代表一个出站的 HTTP Body。
type OutgoingBody struct {
	// 因为go http 的限制Body和Stream 生命周期统一管理
	OutputStreamHandle uint32
	BodyWriter         *io.PipeWriter

	// 可选方法
	SetTrailers func(trailers Fields) error

	ContentLength *uint64
	BytesWritten  atomic.Uint64

	// 消耗标记
	Consumed atomic.Bool
}

func (o *OutgoingBody) Close() error {
	if o == nil {
		return nil
	}
	if o.BodyWriter != nil {
		// 正常关闭用 Close，避免把 EOF 当作错误路径误导读端
		return o.BodyWriter.Close()
	}
	return nil
}

// FutureTrailers 代表一个尚未到达的 HTTP Trailers。
type FutureTrailers struct {
	Pollable *manager_io.ChannelPollable
	Result   ResultTrailers
	Consumed atomic.Bool
	resultMu sync.Mutex // Finish 写 Result 与 get 并发，无锁是 data race
}

func (f *FutureTrailers) StoreResult(res ResultTrailers) {
	if f == nil {
		return
	}
	f.resultMu.Lock()
	f.Result = res
	f.resultMu.Unlock()
}

func (f *FutureTrailers) LoadResult() ResultTrailers {
	if f == nil {
		return ResultTrailers{}
	}
	f.resultMu.Lock()
	defer f.resultMu.Unlock()
	return f.Result
}

type ResultTrailers struct {
	Trailers Fields
	Err      error
}

// RequestOptions 存储了 wasi:http/types.request-options 的状态。
type RequestOptions struct {
	ConnectTimeout      *time.Duration
	FirstByteTimeout    *time.Duration
	BetweenBytesTimeout *time.Duration
}

// FutureIncomingResponse 代表一个尚未到达的 HTTP 响应。
// ResultChan 是实现异步的核心。
type FutureIncomingResponse struct {
	Pollable *manager_io.ChannelPollable
	Consumed atomic.Bool
	Result   Result
	// Cancel 在 drop future 时取消 in-flight 的 Client.Do。
	// 否则 drop 后仍占用连接/goroutine
	Cancel context.CancelFunc

	resultMu sync.Mutex // executeRequest 写入 Result 与 get/drop 并发，无锁会泄漏或 double-close Body
}

// Result 是一个内部类型，用于在 goroutine 之间传递 HTTP 请求的结果。
type Result struct {
	Response *http.Response
	Err      error // 或一个 Go 的 error
}

// StoreResult 在锁内写入 Result。
// executeRequest 与 get/drop 并发写 Result 是 data race。
func (f *FutureIncomingResponse) StoreResult(res Result) {
	if f == nil {
		return
	}
	f.resultMu.Lock()
	f.Result = res
	f.resultMu.Unlock()
}

// LoadResult 在锁内拷贝 Result。
func (f *FutureIncomingResponse) LoadResult() Result {
	if f == nil {
		return Result{}
	}
	f.resultMu.Lock()
	defer f.resultMu.Unlock()
	return f.Result
}

// FieldsManager 使用通用 ResourceManager 来管理 Fields 资源。
type FieldsManager = witgo.ResourceManager[Fields]

func NewFieldsManager() *FieldsManager {
	return witgo.NewResourceManager[Fields](nil)
}

// HTTPManager 是所有 HTTP 相关资源的总管理器。
type HTTPManager struct {
	Fields  *FieldsManager
	Streams *manager_io.StreamManager
	Poll    *manager_io.PollManager

	Options          *witgo.ResourceManager[*RequestOptions]
	OutgoingRequests *witgo.ResourceManager[*OutgoingRequest]
	Futures          *witgo.ResourceManager[*FutureIncomingResponse]
	Responses        *witgo.ResourceManager[*IncomingResponse]
	Bodies           *witgo.ResourceManager[*OutgoingBody]
	FutureTrailers   *witgo.ResourceManager[*FutureTrailers]

	IncomingRequests  *witgo.ResourceManager[*IncomingRequest]
	ResponseOutparams *witgo.ResourceManager[*ResponseOutparam]
	OutgoingResponses *witgo.ResourceManager[*OutgoingResponse]
	IncomingBodies    *witgo.ResourceManager[*IncomingBody]

	// Fields 仍是 http.Header 别名，不能改导出类型；用旁路表标记 WASI 不可变 headers/trailers
	immutableFields sync.Map

	// outgoing-handler 的 http.Client 在 Host.Close 时要关空闲连接；用 map 去重，避免 LRU 淘汰后切片重复追加
	httpClientsMu sync.Mutex
	httpClients   map[*http.Client]struct{}

	// http.Header 非并发安全；fields 方法与 outgoing-handler 分属不同包，必须把锁放在 manager
	fieldsMu sync.Mutex
}

func NewHTTPManager(sm *manager_io.StreamManager, poll *manager_io.PollManager) *HTTPManager {
	// 先建 hm 再注册 destructor，才能在 drop 时回收 headers/body 子句柄
	hm := &HTTPManager{
		Fields:      NewFieldsManager(),
		Streams:     sm,
		Poll:        poll,
		httpClients: make(map[*http.Client]struct{}),
	}

	hm.Options = witgo.NewResourceManager[*RequestOptions](nil)
	hm.OutgoingRequests = witgo.NewResourceManager[*OutgoingRequest](func(resource *OutgoingRequest) {
		if resource == nil {
			return
		}
		if resource.HeadersHandle != 0 {
			hm.Fields.Remove(resource.HeadersHandle)
			resource.HeadersHandle = 0
		}
		if resource.BodyHandle != 0 {
			hm.Bodies.Remove(resource.BodyHandle)
			resource.BodyHandle = 0
		}
		_ = resource.Close()
	})
	hm.Futures = witgo.NewResourceManager[*FutureIncomingResponse](func(resource *FutureIncomingResponse) {
		if resource == nil {
			return
		}
		if resource.Cancel != nil {
			resource.Cancel()
		}
		// 取消后必须等到 Do 结束再关 Body，否则与 executeRequest 赋值 Result 竞态
		if resource.Pollable != nil {
			resource.Pollable.Block()
		}
		// 用 CAS 与 get() 争用 Body 所有权，避免 double-close
		if !resource.Consumed.CompareAndSwap(false, true) {
			return
		}
		resource.resultMu.Lock()
		defer resource.resultMu.Unlock()
		if resource.Result.Response != nil && resource.Result.Response.Body != nil {
			resource.Result.Response.Body.Close()
		}
	})
	hm.Responses = witgo.NewResourceManager[*IncomingResponse](func(resource *IncomingResponse) {
		if resource == nil {
			return
		}
		if resource.HeaderHandle != 0 {
			hm.UnmarkFieldsImmutable(resource.HeaderHandle)
			hm.Fields.Remove(resource.HeaderHandle)
			resource.HeaderHandle = 0
		}
		if resource.BodyHandle != 0 {
			hm.IncomingBodies.Remove(resource.BodyHandle)
			resource.BodyHandle = 0
		} else if !resource.Consumed.Load() && resource.Response != nil && resource.Response.Body != nil {
			_ = resource.Response.Body.Close()
		}
	})
	hm.IncomingRequests = witgo.NewResourceManager[*IncomingRequest](func(resource *IncomingRequest) {
		if resource == nil {
			return
		}
		if resource.Headers != 0 {
			hm.UnmarkFieldsImmutable(resource.Headers)
			hm.Fields.Remove(resource.Headers)
		}
		if resource.BodyHandle != 0 {
			hm.IncomingBodies.Remove(resource.BodyHandle)
			resource.BodyHandle = 0
		} else if resource.Body != nil && !resource.Consumed.Load() {
			_ = resource.Body.Close()
		}
	})
	hm.Bodies = witgo.NewResourceManager[*OutgoingBody](func(resource *OutgoingBody) {
		if resource == nil {
			return
		}
		_ = resource.Close()
		if resource.OutputStreamHandle != 0 && sm != nil {
			sm.Remove(resource.OutputStreamHandle)
		}
	})
	hm.FutureTrailers = witgo.NewResourceManager[*FutureTrailers](nil)
	hm.ResponseOutparams = witgo.NewResourceManager[*ResponseOutparam](func(resource *ResponseOutparam) {
		if resource != nil && resource.ResultChan != nil {
			close(resource.ResultChan)
		}
	})
	hm.OutgoingResponses = witgo.NewResourceManager[*OutgoingResponse](func(resource *OutgoingResponse) {
		if resource == nil {
			return
		}
		if resource.HeadersHandle != 0 {
			hm.Fields.Remove(resource.HeadersHandle)
			resource.HeadersHandle = 0
		}
		if resource.BodyHandle != 0 {
			hm.Bodies.Remove(resource.BodyHandle)
			resource.BodyHandle = 0
		}
		_ = resource.Close()
	})
	hm.IncomingBodies = witgo.NewResourceManager[*IncomingBody](func(resource *IncomingBody) {
		if resource == nil {
			return
		}
		_ = resource.Close()
		// NOTE: 为了防止忘记关闭，这里的生命周期和Stream绑定
		if resource.StreamHandle != 0 && sm != nil {
			sm.Remove(resource.StreamHandle)
		}
	})
	return hm
}

func (hm *HTTPManager) NewOutgoingBody(contentLength *uint64, setTrailers func(trailers Fields) error) (bodyHandle uint32, bodyReader *io.PipeReader, bodyWriter *io.PipeWriter) {
	pr, pw := io.Pipe()

	body := &OutgoingBody{
		BodyWriter:    pw,
		SetTrailers:   setTrailers,
		ContentLength: contentLength,
	}

	bodyHandle = hm.Bodies.Add(body)
	return bodyHandle, pr, pw
}

// MarkFieldsImmutable 将 fields 句柄标为不可变（incoming headers/trailers）。
func (hm *HTTPManager) MarkFieldsImmutable(handle uint32) {
	if hm == nil || handle == 0 {
		return
	}
	hm.immutableFields.Store(handle, struct{}{})
}

// IsFieldsImmutable 报告该 fields 是否禁止 set/append/delete。
func (hm *HTTPManager) IsFieldsImmutable(handle uint32) bool {
	if hm == nil {
		return false
	}
	_, ok := hm.immutableFields.Load(handle)
	return ok
}

// UnmarkFieldsImmutable 在 drop 时清理标记，避免 map 无限增长。
func (hm *HTTPManager) UnmarkFieldsImmutable(handle uint32) {
	if hm == nil {
		return
	}
	hm.immutableFields.Delete(handle)
}

// ClearImmutableFields 在 Host.Close 时清空旁路表。
func (hm *HTTPManager) ClearImmutableFields() {
	if hm == nil {
		return
	}
	hm.immutableFields.Range(func(k, _ any) bool {
		hm.immutableFields.Delete(k)
		return true
	})
}

// TrackClient 登记 outgoing-handler 创建的 http.Client，供 CloseIdleConnections 去重关闭。
func (hm *HTTPManager) TrackClient(c *http.Client) {
	if hm == nil || c == nil {
		return
	}
	hm.httpClientsMu.Lock()
	if hm.httpClients == nil {
		hm.httpClients = make(map[*http.Client]struct{})
	}
	hm.httpClients[c] = struct{}{}
	hm.httpClientsMu.Unlock()
}

// LockFields / UnlockFields 保护 Fields 内 http.Header 的 map 内容。
func (hm *HTTPManager) LockFields() {
	if hm != nil {
		hm.fieldsMu.Lock()
	}
}

func (hm *HTTPManager) UnlockFields() {
	if hm != nil {
		hm.fieldsMu.Unlock()
	}
}

// UntrackClient 在 LRU 淘汰时去掉引用，避免 httpClients 只增不减。
func (hm *HTTPManager) UntrackClient(c *http.Client) {
	if hm == nil || c == nil {
		return
	}
	hm.httpClientsMu.Lock()
	delete(hm.httpClients, c)
	hm.httpClientsMu.Unlock()
}

// CloseIdleConnections 关闭缓存 Client 的空闲连接（Host.Close 调用）。
func (hm *HTTPManager) CloseIdleConnections() {
	if hm == nil {
		return
	}
	hm.httpClientsMu.Lock()
	clients := make([]*http.Client, 0, len(hm.httpClients))
	for c := range hm.httpClients {
		clients = append(clients, c)
	}
	hm.httpClientsMu.Unlock()
	for _, c := range clients {
		if c != nil {
			c.CloseIdleConnections()
		}
	}
}
