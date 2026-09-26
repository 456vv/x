package http

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"reflect"
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

	Method      string
	Path        string
	Query       string
	Scheme      *string
	Authority   *string
	Headers     uint32
	HeadersData Fields     // 与 Fields 表里同一块 map，供 drop 时 RemoveIf
	hdrMu       sync.Mutex // Headers 句柄与 Headers() 并发

	Body io.ReadCloser

	// BodyHandle 用于 Guest 端消费 Body
	BodyHandle uint32

	// 与 IncomingResponse 对齐，保证 consume 至多一次（并发安全）
	Consumed atomic.Bool

	bodyMu   sync.Mutex // consume 写 BodyHandle 与 drop 关 Body 必须同一把锁
	dropping bool       // drop 已取走 Body 后禁止 consume 再 Add
}

func (r *IncomingRequest) TakeHeadersForDrop() (handle uint32, data Fields) {
	if r == nil {
		return 0, nil
	}
	r.hdrMu.Lock()
	handle = r.Headers
	data = r.HeadersData
	r.Headers = 0
	r.hdrMu.Unlock()
	return handle, data
}

func (r *IncomingRequest) LoadHeadersHandle() uint32 {
	if r == nil {
		return 0
	}
	r.hdrMu.Lock()
	h := r.Headers
	r.hdrMu.Unlock()
	return h
}

// InstallIncomingBody 在锁内 CAS + 登记 incoming-body，并把 Body 所有权交出去。
// wasip2 与 manager 不同包，不能直接锁未导出的 bodyMu。
func (r *IncomingRequest) InstallIncomingBody(add func(io.ReadCloser) uint32) (uint32, bool) {
	if r == nil || add == nil {
		return 0, false
	}
	r.bodyMu.Lock()
	defer r.bodyMu.Unlock()
	if r.dropping || !r.Consumed.CompareAndSwap(false, true) {
		return 0, false
	}
	h := add(r.Body)
	r.BodyHandle = h
	r.Body = nil // 所有权已在 IncomingBody；drop 只认 BodyHandle
	return h, true
}

// ReleaseForDrop 在析构时取出尚未移交的 Body 或已登记的句柄（互斥）。
func (r *IncomingRequest) ReleaseForDrop() (bodyHandle uint32, body io.ReadCloser) {
	if r == nil {
		return 0, nil
	}
	r.bodyMu.Lock()
	r.dropping = true
	r.Consumed.Store(true) // drop 已开始则并发 consume 必须失败，不能对已关 Body 再 Add
	bodyHandle = r.BodyHandle
	r.BodyHandle = 0
	body = r.Body
	r.Body = nil
	r.bodyMu.Unlock()
	return bodyHandle, body
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
	headersOnce   sync.Once  // 并发 headers() 会重复 Add；Once 必须留在本包，外包不能碰未导出字段
	hdrMu         sync.Mutex // HeadersHandle 与 headers()/drop 并发

	Body io.Reader

	// 这里只是引用资源，不属于OutgoingRequest生命周期管理
	// BodyWriter 用于在 Host 端写入 Guest 提供的数据
	BodyWriter *io.PipeWriter
	BodyHandle uint32 // 指向 outgoing-body 资源的句柄

	// 消耗标记
	Consumed atomic.Bool

	bodyMu   sync.Mutex // body() 里 NewOutgoingBody 已 Add，写入 BodyHandle 前 drop 会漏回收
	dropping bool

	metaMu sync.RWMutex

	pendingTrailers Fields
}

func (o *OutgoingRequest) TakeHeadersHandleForDrop() (uint32, Fields) {
	if o == nil {
		return 0, nil
	}
	// 等正在进行的 headers() 把句柄写入；或抢先让后续 headers() 变成 no-op。
	o.headersOnce.Do(func() {})
	o.hdrMu.Lock()
	h := o.HeadersHandle
	hdr := o.Headers
	o.HeadersHandle = 0
	o.hdrMu.Unlock()
	return h, hdr
}

func (o *OutgoingRequest) StoreHeadersHandle(h uint32) {
	if o == nil {
		return
	}
	o.hdrMu.Lock()
	o.HeadersHandle = h
	o.hdrMu.Unlock()
}

func (o *OutgoingRequest) LoadHeadersHandle() uint32 {
	if o == nil {
		return 0
	}
	o.hdrMu.Lock()
	h := o.HeadersHandle
	o.hdrMu.Unlock()
	return h
}

// StoreTrailers 保存 finish 时的 trailer。Request 已经存在时同时写进即将发送的请求。
// finish 既可能早于 handle，也可能晚于 handle；只处理其中一条路径都会丢 trailer。
func (o *OutgoingRequest) StoreTrailers(t Fields) {
	if o == nil || len(t) == 0 {
		return
	}
	o.metaMu.Lock()
	defer o.metaMu.Unlock()
	if o.pendingTrailers == nil {
		o.pendingTrailers = make(Fields, len(t))
	}
	for k, vv := range t {
		cp := make([]string, len(vv))
		copy(cp, vv)
		o.pendingTrailers[k] = append(o.pendingTrailers[k], cp...)
	}
	if o.Request == nil {
		return
	}
	if o.Request.Trailer == nil {
		o.Request.Trailer = make(Fields)
	}
	for k, vv := range t {
		cp := make([]string, len(vv))
		copy(cp, vv)
		o.Request.Trailer[k] = append(o.Request.Trailer[k], cp...)
		o.Request.Header.Add("Trailer", k)
	}
}

// BindRequest 发布 http.Request，并把此前 finish 保存的 trailer 挂上去。
func (o *OutgoingRequest) BindRequest(r *http.Request) {
	if o == nil || r == nil {
		return
	}
	o.metaMu.Lock()
	defer o.metaMu.Unlock()
	o.Request = r
	if len(o.pendingTrailers) == 0 {
		return
	}
	if r.Trailer == nil {
		r.Trailer = make(Fields, len(o.pendingTrailers))
	}
	for k, vv := range o.pendingTrailers {
		cp := make([]string, len(vv))
		copy(cp, vv)
		r.Trailer[k] = cp
		r.Header.Add("Trailer", k)
	}
}

// LoadRequestLine 复制请求行。修改原因：set-* 与 buildGoRequest 会并发读写这些字段。
func (o *OutgoingRequest) LoadRequestLine() (method, path string, scheme, authority *string) {
	if o == nil {
		return "", "", nil, nil
	}
	o.metaMu.RLock()
	defer o.metaMu.RUnlock()
	return o.Method, o.Path, cloneStringPtr(o.Scheme), cloneStringPtr(o.Authority)
}

func (o *OutgoingRequest) StoreMethod(method string) {
	if o == nil {
		return
	}
	o.metaMu.Lock()
	o.Method = method
	o.metaMu.Unlock()
}

func (o *OutgoingRequest) StorePath(path string) {
	if o == nil {
		return
	}
	o.metaMu.Lock()
	o.Path = path
	o.metaMu.Unlock()
}

func (o *OutgoingRequest) StoreScheme(scheme *string) {
	if o == nil {
		return
	}
	o.metaMu.Lock()
	o.Scheme = cloneStringPtr(scheme)
	o.metaMu.Unlock()
}

func (o *OutgoingRequest) StoreAuthority(authority *string) {
	if o == nil {
		return
	}
	o.metaMu.Lock()
	o.Authority = cloneStringPtr(authority)
	o.metaMu.Unlock()
}

func (o *OutgoingRequest) Close() error {
	if o == nil {
		return nil
	}
	// 无锁读 Body 与 InstallOutgoingBody 并发是 data race。
	o.bodyMu.Lock()
	body := o.Body
	o.Body = nil
	o.bodyMu.Unlock()
	if closer, ok := body.(io.Closer); ok {
		return closer.Close()
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

// InstallOutgoingBody 在锁内创建并登记 outgoing-body。
func (o *OutgoingRequest) InstallOutgoingBody(create func() (handle uint32, body io.Reader, pw *io.PipeWriter)) (uint32, bool) {
	if o == nil || create == nil {
		return 0, false
	}
	o.bodyMu.Lock()
	defer o.bodyMu.Unlock()
	if o.dropping || !o.Consumed.CompareAndSwap(false, true) {
		return 0, false
	}
	h, body, pw := create()
	o.BodyHandle = h
	o.Body = body
	o.BodyWriter = pw
	return h, true
}

// TakeBodyHandleForDrop 取出 outgoing-body 句柄；drop 期间禁止再 body()。
func (o *OutgoingRequest) TakeBodyHandleForDrop() uint32 {
	if o == nil {
		return 0
	}
	o.bodyMu.Lock()
	o.dropping = true
	o.Consumed.Store(true)
	h := o.BodyHandle
	o.BodyHandle = 0
	o.bodyMu.Unlock()
	return h
}

// IncomingResponse 代表一个已到达的、由 Host 接收的 HTTP 响应。
type IncomingResponse struct {
	Response *http.Response

	StatusCode int
	Headers    Fields

	// WASI incoming-response.headers 每次返回同一不可变快照句柄
	HeaderHandle uint32
	HeaderSnap   Fields     // headers() 实际 Add 进表的那份 clone，不是 Headers 原件
	headersOnce  sync.Once  // 并发 headers() 会重复 Clone+Add
	hdrMu        sync.Mutex // HeadersHandle 与 headers()/drop 并发

	Body       *IncomingBody
	BodyHandle uint32 // 指向 incoming-body 的句柄
	// 消耗标记
	Consumed atomic.Bool

	bodyMu   sync.Mutex // 与 IncomingRequest 相同，consume/drop 争用 Response.Body
	dropping bool
}

func (r *IncomingResponse) TakeHeaderHandleForDrop() (uint32, Fields) {
	if r == nil {
		return 0, nil
	}
	r.headersOnce.Do(func() {})
	r.hdrMu.Lock()
	h := r.HeaderHandle
	snap := r.HeaderSnap
	r.HeaderHandle = 0
	r.hdrMu.Unlock()
	return h, snap
}

func (r *IncomingResponse) StoreHeaderHandle(h uint32, snap Fields) {
	if r == nil {
		return
	}
	r.hdrMu.Lock()
	r.HeaderHandle = h
	r.HeaderSnap = snap
	r.hdrMu.Unlock()
}

func (r *IncomingResponse) LoadHeaderHandle() uint32 {
	if r == nil {
		return 0
	}
	r.hdrMu.Lock()
	h := r.HeaderHandle
	r.hdrMu.Unlock()
	return h
}

// EnsureHeaderHandle 保证 incoming-response.headers 只创建一次快照句柄。
func (r *IncomingResponse) EnsureHeaderHandle(init func()) {
	if r == nil || init == nil {
		return
	}
	r.headersOnce.Do(init)
}

// InstallIncomingBody 在锁内移交 http.Response.Body，避免与 drop 双关。
func (r *IncomingResponse) InstallIncomingBody(add func(io.ReadCloser) uint32) (uint32, bool) {
	if r == nil || add == nil {
		return 0, false
	}
	r.bodyMu.Lock()
	defer r.bodyMu.Unlock()
	if r.dropping || !r.Consumed.CompareAndSwap(false, true) {
		return 0, false
	}
	var rc io.ReadCloser
	if r.Response != nil {
		rc = r.Response.Body
	}
	h := add(rc)
	r.BodyHandle = h
	if r.Response != nil {
		// 必须先 Add 再摘 Body。先改 NoBody 再 Add，create panic 会丢失未关闭 Body。
		r.Response.Body = http.NoBody
	}
	return h, true
}

// ReleaseForDrop 在析构时取出 incoming-body 句柄或尚未 consume 的 Body。
func (r *IncomingResponse) ReleaseForDrop() (bodyHandle uint32, body io.ReadCloser) {
	if r == nil {
		return 0, nil
	}
	r.bodyMu.Lock()
	r.dropping = true
	r.Consumed.Store(true)
	bodyHandle = r.BodyHandle
	r.BodyHandle = 0
	if r.Response != nil {
		body = r.Response.Body
		r.Response.Body = http.NoBody
	}
	r.bodyMu.Unlock()
	if body == http.NoBody {
		body = nil
	}
	return bodyHandle, body
}

// OutgoingResponse 代表一个由 Guest 构建的出站 HTTP 响应。
type OutgoingResponse struct {
	Response http.ResponseWriter

	StatusCode int
	Headers    Fields

	// 与 OutgoingRequest 相同，缓存 headers() 子资源
	HeadersHandle uint32
	headersOnce   sync.Once  // 并发 headers() 会重复 Add
	hdrMu         sync.Mutex // HeadersHandle 与 headers()/drop 并发

	Body io.Reader

	// 这里只是引用资源，不属于OutgoingResponse生命周期管理
	BodyWriter *io.PipeWriter
	BodyHandle uint32

	// 消耗标记
	Consumed atomic.Bool

	bodyMu   sync.Mutex // 与 OutgoingRequest.body 相同的 BodyHandle 窗口
	dropping bool

	metaMu sync.RWMutex

	pendingTrailers Fields
}

func (o *OutgoingResponse) TakeHeadersHandleForDrop() (uint32, Fields) {
	if o == nil {
		return 0, nil
	}
	// 修改原因：等正在进行的 headers() 把句柄写入；或抢先让后续 headers() 变成 no-op。
	o.headersOnce.Do(func() {})
	o.hdrMu.Lock()
	h := o.HeadersHandle
	hdr := o.Headers
	o.HeadersHandle = 0
	o.hdrMu.Unlock()
	return h, hdr
}

func (o *OutgoingResponse) StoreHeadersHandle(h uint32) {
	if o == nil {
		return
	}
	o.hdrMu.Lock()
	o.HeadersHandle = h
	o.hdrMu.Unlock()
}

func (o *OutgoingResponse) LoadHeadersHandle() uint32 {
	if o == nil {
		return 0
	}
	o.hdrMu.Lock()
	h := o.HeadersHandle
	o.hdrMu.Unlock()
	return h
}

// StoreTrailers 在写响应头之前暂存 trailer。net/http 要求 Trailer 前缀先于 WriteHeader。
func (o *OutgoingResponse) StoreTrailers(t Fields) {
	if o == nil || len(t) == 0 {
		return
	}
	o.metaMu.Lock()
	defer o.metaMu.Unlock()
	if o.pendingTrailers == nil {
		o.pendingTrailers = make(Fields, len(t))
	}
	for k, vv := range t {
		cp := make([]string, len(vv))
		copy(cp, vv)
		o.pendingTrailers[k] = append(o.pendingTrailers[k], cp...)
	}
}

// TakeTrailers 取出并清空 trailer，避免重复写入。
func (o *OutgoingResponse) TakeTrailers() Fields {
	if o == nil {
		return nil
	}
	o.metaMu.Lock()
	defer o.metaMu.Unlock()
	t := o.pendingTrailers
	o.pendingTrailers = nil
	return t
}

func (o *OutgoingResponse) LoadStatus() int {
	if o == nil {
		return 0
	}
	o.metaMu.RLock()
	defer o.metaMu.RUnlock()
	return o.StatusCode
}

func (o *OutgoingResponse) StoreStatus(code int) {
	if o == nil {
		return
	}
	o.metaMu.Lock()
	o.StatusCode = code
	o.metaMu.Unlock()
}

func (o *OutgoingResponse) Close() error {
	if o == nil {
		return nil
	}
	// 与 OutgoingRequest.Close 相同，必须和 body() 共用 bodyMu。
	o.bodyMu.Lock()
	body := o.Body
	o.Body = nil
	o.bodyMu.Unlock()
	if closer, ok := body.(io.Closer); ok {
		return closer.Close()
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

// InstallOutgoingBody 在锁内创建并登记 outgoing-body。
func (o *OutgoingResponse) InstallOutgoingBody(create func() (handle uint32, body io.Reader, pw *io.PipeWriter)) (uint32, bool) {
	if o == nil || create == nil {
		return 0, false
	}
	o.bodyMu.Lock()
	defer o.bodyMu.Unlock()
	if o.dropping || !o.Consumed.CompareAndSwap(false, true) {
		return 0, false
	}
	h, body, pw := create()
	o.BodyHandle = h
	o.Body = body
	o.BodyWriter = pw
	return h, true
}

// TakeBodyHandleForDrop 取出 outgoing-body 句柄。
func (o *OutgoingResponse) TakeBodyHandleForDrop() uint32 {
	if o == nil {
		return 0
	}
	o.bodyMu.Lock()
	o.dropping = true
	o.Consumed.Store(true)
	h := o.BodyHandle
	o.BodyHandle = 0
	o.bodyMu.Unlock()
	return h
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
	StreamHandle uint32             // 指向 input-stream 的句柄
	streamRes    *manager_io.Stream // finish/drop 不能按复用后的句柄 Get/Remove 别人的流
	Stream       io.Reader

	// 可选方法
	GetTrailers func() (trailers Fields)

	// 消耗标记
	Consumed atomic.Bool

	// stream() 写 StreamHandle 与 drop/finish 读句柄无同步，
	// 既是数据竞争，也会在句柄尚未写入时漏掉 Streams.Remove。
	streamMu sync.Mutex
	dropping bool
}

// InstallStream 在锁内完成一次性 stream()，并登记 input-stream 句柄。
func (o *IncomingBody) InstallStream(create func(r io.Reader) (uint32, *manager_io.Stream)) (uint32, bool) {
	if o == nil || create == nil {
		return 0, false
	}
	o.streamMu.Lock()
	defer o.streamMu.Unlock()
	if o.dropping || o.StreamHandle != 0 || !o.Consumed.CompareAndSwap(false, true) {
		return 0, false
	}
	r := o.Stream
	if r == nil {
		r = bytes.NewReader(nil)
	}
	h, st := create(r)
	o.StreamHandle = h
	o.streamRes = st
	return h, true
}

// TakeStreamHandleForDrop 取出流句柄。drop/finish 之后 stream() 必须失败。
func (o *IncomingBody) TakeStreamHandleForDrop() (uint32, *manager_io.Stream) {
	if o == nil {
		return 0, nil
	}
	o.streamMu.Lock()
	o.dropping = true
	o.Consumed.Store(true)
	h := o.StreamHandle
	st := o.streamRes
	o.StreamHandle = 0
	o.streamRes = nil
	o.streamMu.Unlock()
	return h, st
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

	// stream() 写 StreamHandle 与 drop/finish 读句柄无同步，
	// 既是数据竞争，也会在句柄尚未写入时漏掉 Streams.Remove。
	streamMu  sync.Mutex
	dropping  bool
	streamRes *manager_io.Stream // finish/drop 不能按复用后的句柄 Get/Remove 别人的流
}

// InstallOutputStream 在锁内完成一次性 write()。
func (o *OutgoingBody) InstallOutputStream(create func(w *io.PipeWriter) (uint32, *manager_io.Stream)) (uint32, bool) {
	if o == nil || create == nil {
		return 0, false
	}
	o.streamMu.Lock()
	defer o.streamMu.Unlock()
	if o.dropping || o.OutputStreamHandle != 0 || !o.Consumed.CompareAndSwap(false, true) {
		return 0, false
	}
	h, st := create(o.BodyWriter)
	o.OutputStreamHandle = h
	o.streamRes = st
	return h, true
}

// TakeOutputStreamHandleForDrop 取出 output-stream 句柄。
func (o *OutgoingBody) TakeOutputStreamHandleForDrop() (uint32, *manager_io.Stream) {
	if o == nil {
		return 0, nil
	}
	o.streamMu.Lock()
	o.dropping = true
	o.Consumed.Store(true)
	h := o.OutputStreamHandle
	st := o.streamRes
	o.OutputStreamHandle = 0
	o.streamRes = nil
	o.streamMu.Unlock()
	return h, st
}

func (o *OutgoingBody) Close() error {
	if o == nil {
		return nil
	}
	if o.BodyWriter != nil {
		// resource-drop 且未 finish 时，WASI 要求把 body 视为不完整。
		// PipeWriter.Close() 是干净 EOF，http.Client / io.Copy 会把截断当发送成功。
		// Finish() 走 Pop，不进析构，自行 Close() 表示正常结束。
		return o.BodyWriter.CloseWithError(io.ErrUnexpectedEOF)
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
	mu                  sync.Mutex
}

func cloneDurationPtr(d *time.Duration) *time.Duration {
	if d == nil {
		return nil
	}
	v := *d
	return &v
}

func cloneStringPtr(s *string) *string {
	if s == nil {
		return nil
	}
	v := *s
	return &v
}

// CopyTimeouts 返回超时副本。修改原因：setter 替换 *time.Duration 与 getClient 读取是数据竞争。
func (o *RequestOptions) CopyTimeouts() (connect, first, between *time.Duration) {
	if o == nil {
		return nil, nil, nil
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	return cloneDurationPtr(o.ConnectTimeout), cloneDurationPtr(o.FirstByteTimeout), cloneDurationPtr(o.BetweenBytesTimeout)
}

func (o *RequestOptions) StoreConnectTimeout(d *time.Duration) {
	if o == nil {
		return
	}
	o.mu.Lock()
	o.ConnectTimeout = cloneDurationPtr(d)
	o.mu.Unlock()
}

func (o *RequestOptions) StoreFirstByteTimeout(d *time.Duration) {
	if o == nil {
		return
	}
	o.mu.Lock()
	o.FirstByteTimeout = cloneDurationPtr(d)
	o.mu.Unlock()
}

func (o *RequestOptions) StoreBetweenBytesTimeout(d *time.Duration) {
	if o == nil {
		return
	}
	o.mu.Lock()
	o.BetweenBytesTimeout = cloneDurationPtr(d)
	o.mu.Unlock()
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
		if h, hdr := resource.TakeHeadersHandleForDrop(); h != 0 {
			hm.Fields.RemoveIf(h, func(cur Fields) bool { return SameFields(cur, hdr) })
		}
		if h := resource.TakeBodyHandleForDrop(); h != 0 {
			pw := resource.BodyWriter
			hm.Bodies.RemoveIf(h, func(cur *OutgoingBody) bool {
				return cur != nil && pw != nil && cur.BodyWriter == pw
			})
		}
		resource.Close()
	})
	hm.Futures = witgo.NewResourceManager[*FutureIncomingResponse](func(resource *FutureIncomingResponse) {
		if resource == nil {
			return
		}
		if resource.Cancel != nil {
			resource.Cancel()
			resource.Cancel = nil // Clear/Remove 二次路径不要重复观察 Cancel
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
			resource.Result.Response.Body = http.NoBody
		}
	})
	hm.Responses = witgo.NewResourceManager[*IncomingResponse](func(resource *IncomingResponse) {
		if resource == nil {
			return
		}
		if h, snap := resource.TakeHeaderHandleForDrop(); h != 0 {
			hm.Fields.RemoveIfWith(h, func(cur Fields) bool {
				return SameFields(cur, snap)
			}, func(Fields) { hm.UnmarkFieldsImmutable(h) })
		}
		bh, body := resource.ReleaseForDrop()
		if bh != 0 {
			hm.IncomingBodies.Remove(bh)
		} else if body != nil {
			body.Close()
		}
	})
	hm.IncomingRequests = witgo.NewResourceManager[*IncomingRequest](func(resource *IncomingRequest) {
		if resource == nil {
			return
		}
		if h, data := resource.TakeHeadersForDrop(); h != 0 {
			hm.Fields.RemoveIfWith(h, func(cur Fields) bool {
				return SameFields(cur, data)
			}, func(Fields) { hm.UnmarkFieldsImmutable(h) })
		}
		bh, body := resource.ReleaseForDrop()
		if bh != 0 {
			hm.IncomingBodies.Remove(bh)
		} else if body != nil {
			body.Close()
		}
	})
	hm.Bodies = witgo.NewResourceManager[*OutgoingBody](func(resource *OutgoingBody) {
		if resource == nil {
			return
		}
		h, st := resource.TakeOutputStreamHandleForDrop()
		// 先 CloseWithError(pw) 再 Flush，会把未写出的缓冲写成 broken pipe。
		// 先停 output-stream（DontCloseWriter 时会 BlockingFlush），再标记 body 不完整。
		if h != 0 && sm != nil {
			sm.RemoveIf(h, func(cur *manager_io.Stream) bool { return cur != nil && cur == st })
		}
		resource.Close()
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
		if h, hdr := resource.TakeHeadersHandleForDrop(); h != 0 {
			hm.Fields.RemoveIf(h, func(cur Fields) bool { return SameFields(cur, hdr) })
		}
		if h := resource.TakeBodyHandleForDrop(); h != 0 {
			pw := resource.BodyWriter
			hm.Bodies.RemoveIf(h, func(cur *OutgoingBody) bool {
				return cur != nil && pw != nil && cur.BodyWriter == pw
			})
		}
		resource.Close()
	})
	hm.IncomingBodies = witgo.NewResourceManager[*IncomingBody](func(resource *IncomingBody) {
		if resource == nil {
			return
		}
		h, st := resource.TakeStreamHandleForDrop()
		// 先关底层 Body 再 Remove 流，后台 Read 更容易 use-after-close。
		if h != 0 && sm != nil {
			sm.RemoveIf(h, func(cur *manager_io.Stream) bool { return cur != nil && cur == st })
		}
		resource.Close()
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
	// Host.Close 后不再持有 Client 引用，避免 map 只增不减导致 Transport 常驻
	hm.httpClients = make(map[*http.Client]struct{})
	hm.httpClientsMu.Unlock()
	for _, c := range clients {
		if c != nil {
			c.CloseIdleConnections()
		}
	}
}

// SameFields 用 map 头指针比较 http.Header 身份，供 RemoveIf 使用。
// 句柄 uint32 会复用；不能只凭编号 Remove child fields。
func SameFields(a, b Fields) bool {
	if a == nil || b == nil {
		return false
	}
	return reflect.ValueOf(a).Pointer() == reflect.ValueOf(b).Pointer()
}
