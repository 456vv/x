package v0_2

import (
	"context"
	"io"
	"strconv"
	"strings"

	manager_http "github.com/456vv/x/vweb_dynamic/internal/wazero/manager/http"
	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"
)

// guest 发起请求 host 处理
type outgoingRequestImpl struct {
	hm *manager_http.HTTPManager
}

func newOutgoingRequestImpl(hm *manager_http.HTTPManager) *outgoingRequestImpl {
	return &outgoingRequestImpl{hm: hm}
}

func (i *outgoingRequestImpl) Constructor(fields Fields) OutgoingRequest {
	// 消耗掉资源
	header, _ := i.hm.Fields.Pop(fields)
	req := &manager_http.OutgoingRequest{
		Headers: header,
		Method:  "GET", // WASI outgoing-request 默认 method 为 GET；空字符串会让 Method() 与 net/http 行为偏离规范。
	}
	return i.hm.OutgoingRequests.Add(req)
}

func (i *outgoingRequestImpl) Drop(_ context.Context, handle OutgoingRequest) {
	i.hm.OutgoingRequests.Remove(handle)
}

// 返回当前请求对应的输出体（outgoing-body）资源。
// 仅首次调用成功，最多获取一次；后续调用返回错误。
func (i *outgoingRequestImpl) Body(_ context.Context, this OutgoingRequest) witgo.Result[OutgoingBody, witgo.Unit] {
	req, ok := i.hm.OutgoingRequests.Get(this)
	if !ok || req == nil { // Get 成功但资源为 nil 时解引用会 panic
		return witgo.Err[OutgoingBody, witgo.Unit](witgo.Unit{})
	}

	// NewOutgoingBody 已 Bodies.Add，必须在锁内写入 BodyHandle，否则 drop 漏回收。
	handle, ok := req.InstallOutgoingBody(func() (uint32, io.Reader, *io.PipeWriter) {
		var contentLength *uint64
		i.hm.LockFields()
		cl := headerValuesGet(req.Headers, "Content-Length")
		i.hm.UnlockFields()
		if len(cl) > 0 {
			if val, err := strconv.ParseUint(cl, 10, 64); err == nil {
				contentLength = &val
			}
		}
		h, body, pw := i.hm.NewOutgoingBody(contentLength, func(trailers manager_http.Fields) error {
			req.StoreTrailers(trailers)
			return nil
		})
		return h, body, pw
	})
	if !ok {
		return witgo.Err[OutgoingBody, witgo.Unit](witgo.Unit{})
	}
	return witgo.Ok[OutgoingBody, witgo.Unit](handle)
}

func (i *outgoingRequestImpl) Method(_ context.Context, this OutgoingRequest) Method {
	req, ok := i.hm.OutgoingRequests.Get(this)
	if !ok || req == nil {
		return Method{Other: witgo.String("")}
	}
	method, _, _, _ := req.LoadRequestLine()
	return toWasiMethod(method)
}

func (i *outgoingRequestImpl) PathWithQuery(_ context.Context, this OutgoingRequest) witgo.Option[string] {
	req, ok := i.hm.OutgoingRequests.Get(this)
	if !ok || req == nil {
		return witgo.None[string]()
	}
	_, path, _, _ := req.LoadRequestLine()
	if path == "" {
		return witgo.None[string]()
	}
	return witgo.Some(path)
}

func (i *outgoingRequestImpl) Scheme(_ context.Context, this OutgoingRequest) witgo.Option[Scheme] {
	req, ok := i.hm.OutgoingRequests.Get(this)
	if !ok || req == nil {
		return witgo.None[Scheme]()
	}
	_, _, scheme, _ := req.LoadRequestLine()
	if scheme == nil {
		return witgo.None[Scheme]()
	}
	return witgo.Some(toWasiScheme(*scheme))
}

func (i *outgoingRequestImpl) Authority(_ context.Context, this OutgoingRequest) witgo.Option[string] {
	req, ok := i.hm.OutgoingRequests.Get(this)
	if !ok || req == nil {
		return witgo.None[string]()
	}
	_, _, _, authority := req.LoadRequestLine()
	if authority == nil {
		return witgo.None[string]()
	}
	return witgo.Some(*authority)
}

func (i *outgoingRequestImpl) Headers(_ context.Context, this OutgoingRequest) Headers {
	req, ok := i.hm.OutgoingRequests.Get(this)
	if !ok || req == nil {
		return 0
	}
	// HeadersHandle 与 drop 并发时无锁读写是 data race；
	// 必须走 Store/Load。Once 保证只 Add 一次；drop 的 TakeHeadersHandleForDrop
	// 会先 Once.Do(空)，等本函数写完句柄再取走，避免漏回收。
	req.EnsureHeadersHandle(func() {
		if req.LoadHeadersHandle() != 0 {
			return
		}
		if req.Headers == nil {
			req.Headers = make(manager_http.Fields)
		}
		req.StoreHeadersHandle(i.hm.Fields.Add(manager_http.Fields(req.Headers)))
	})
	return req.LoadHeadersHandle()
}

// validHTTPToken 按 RFC 9110 tchar 校验 method。
// set-method 原先原样保存，空串、空格和分隔符会进请求行。
func validHTTPToken(s string) bool {
	if s == "" || len(s) > 64 {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9':
		case c == '!' || c == '#' || c == '$' || c == '%' || c == '&' || c == '\'' || c == '*' ||
			c == '+' || c == '-' || c == '.' || c == '^' || c == '_' || c == '`' || c == '|' || c == '~':
		default:
			return false
		}
	}
	return true
}

func validPort(port string) bool {
	if port == "" || len(port) > 5 {
		return false
	}
	n := 0
	for i := 0; i < len(port); i++ {
		c := port[i]
		if c < '0' || c > '9' {
			return false
		}
		n = n*10 + int(c-'0')
		if n > 65535 {
			return false
		}
	}
	return true
}

// validRequestAuthority 拒绝 userinfo 和非法 host。
// authority 中的 @ 会让 net/http 把请求发到用户指定的另一台主机。
// 不用 url.Parse 一刀切，避免带 zone 的 IPv6 字面量被误拒。
func validRequestAuthority(authority string) bool {
	if authority == "" || len(authority) > 255 || strings.ContainsAny(authority, "\r\n\x00/\\@ \t") {
		return false
	}
	if strings.HasPrefix(authority, "[") {
		end := strings.LastIndex(authority, "]")
		if end <= 1 {
			return false
		}
		portPart := authority[end+1:]
		if portPart == "" {
			return true
		}
		return strings.HasPrefix(portPart, ":") && validPort(portPart[1:])
	}
	if strings.Count(authority, ":") > 1 {
		return false
	}
	host := authority
	if i := strings.LastIndex(host, ":"); i >= 0 {
		if !validPort(host[i+1:]) {
			return false
		}
		host = host[:i]
	}
	return host != ""
}

// validRequestPath 校验 path-with-query。空字符串表示清除。
// 不带前导 / 的值会被拼进 URL；反斜杠和 CTL 会破坏请求行。
func validRequestPath(path string) bool {
	if path == "" {
		return true
	}
	if len(path) > 8192 || strings.ContainsAny(path, "\r\n\x00\\") {
		return false
	}
	return path == "*" || strings.HasPrefix(path, "/")
}

// validRequestScheme 只校验 scheme 语法，不在这里限制 http/https。
func validRequestScheme(scheme string) bool {
	if scheme == "" || len(scheme) > 32 {
		return false
	}
	for i := 0; i < len(scheme); i++ {
		c := scheme[i]
		if i == 0 {
			if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')) {
				return false
			}
			continue
		}
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '+' || c == '-' || c == '.') {
			return false
		}
	}
	return true
}

func (i *outgoingRequestImpl) SetMethod(_ context.Context, this OutgoingRequest, method Method) witgo.UnitResult {
	req, ok := i.hm.OutgoingRequests.Get(this)
	if !ok || req == nil {
		return witgo.UintErr()
	}
	m := fromWasiMethod(method)
	if !validHTTPToken(m) {
		return witgo.UintErr()
	}
	req.StoreMethod(m)
	return witgo.UintOk()
}

func (i *outgoingRequestImpl) SetPathWithQuery(_ context.Context, this OutgoingRequest, path witgo.Option[string]) witgo.UnitResult {
	req, ok := i.hm.OutgoingRequests.Get(this)
	if !ok || req == nil {
		return witgo.UintErr()
	}
	if path.Some != nil {
		if !validRequestPath(*path.Some) {
			return witgo.UintErr()
		}
		req.StorePath(*path.Some)
	} else {
		req.StorePath("")
	}
	return witgo.UintOk()
}

func (i *outgoingRequestImpl) SetScheme(_ context.Context, this OutgoingRequest, scheme witgo.Option[Scheme]) witgo.UnitResult {
	req, ok := i.hm.OutgoingRequests.Get(this)
	if !ok || req == nil {
		return witgo.UintErr()
	}
	if scheme.Some != nil {
		s := fromWasiScheme(*scheme.Some)
		if s == nil || !validRequestScheme(*s) {
			return witgo.UintErr()
		}
		req.StoreScheme(s)
	} else {
		req.StoreScheme(nil)
	}
	return witgo.UintOk()
}

func (i *outgoingRequestImpl) SetAuthority(_ context.Context, this OutgoingRequest, authority witgo.Option[string]) witgo.UnitResult {
	req, ok := i.hm.OutgoingRequests.Get(this)
	if !ok || req == nil {
		return witgo.UintErr()
	}
	if authority.Some != nil && !validRequestAuthority(*authority.Some) {
		return witgo.UintErr()
	}
	req.StoreAuthority(authority.Some)
	return witgo.UintOk()
}
