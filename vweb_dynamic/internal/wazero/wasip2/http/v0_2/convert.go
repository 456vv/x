package v0_2

import (
	"errors"
	"net"
	"net/http"
	"os"
	"strings"

	manager_http "github.com/456vv/x/vweb_dynamic/internal/wazero/manager/http"
	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"
)

// fromWasiMethod 将 WIT 的 Method variant 转换为 Go 的 HTTP 方法字符串。
func fromWasiMethod(method Method) string {
	switch {
	case method.Get != nil:
		return http.MethodGet
	case method.Head != nil:
		return http.MethodHead
	case method.Post != nil:
		return http.MethodPost
	case method.Put != nil:
		return http.MethodPut
	case method.Delete != nil:
		return http.MethodDelete
	case method.Connect != nil:
		return http.MethodConnect
	case method.Options != nil:
		return http.MethodOptions
	case method.Trace != nil:
		return http.MethodTrace
	case method.Patch != nil:
		return http.MethodPatch
	case method.Other != nil:
		return *method.Other
	default:
		return "" // 或者返回一个错误
	}
}

// toWasiMethod 将 Go 的 HTTP 方法字符串转换为 WIT 的 Method variant。
func toWasiMethod(method string) Method {
	switch strings.ToUpper(method) {
	case http.MethodGet:
		return Method{Get: &witgo.Unit{}}
	case http.MethodHead:
		return Method{Head: &witgo.Unit{}}
	case http.MethodPost:
		return Method{Post: &witgo.Unit{}}
	case http.MethodPut:
		return Method{Put: &witgo.Unit{}}
	case http.MethodDelete:
		return Method{Delete: &witgo.Unit{}}
	case http.MethodConnect:
		return Method{Connect: &witgo.Unit{}}
	case http.MethodOptions:
		return Method{Options: &witgo.Unit{}}
	case http.MethodTrace:
		return Method{Trace: &witgo.Unit{}}
	case http.MethodPatch:
		return Method{Patch: &witgo.Unit{}}
	default:
		// 对于不在标准列表中的方法，使用 Other case。
		return Method{Other: &method}
	}
}

// fromWasiScheme 将 WIT 的 Scheme variant 转换为 Go 的 URL scheme 字符串指针。
func fromWasiScheme(scheme Scheme) *string {
	var s string
	switch {
	case scheme.HTTP != nil:
		s = "http"
	case scheme.HTTPS != nil:
		s = "https"
	case scheme.Other != nil:
		s = *scheme.Other
	default:
		return nil
	}
	return &s
}

// toWasiScheme 将 Go 的 URL scheme 字符串转换为 WIT 的 Scheme variant。
func toWasiScheme(scheme string) Scheme {
	switch strings.ToLower(scheme) {
	case "http":
		return Scheme{HTTP: &witgo.Unit{}}
	case "https":
		return Scheme{HTTPS: &witgo.Unit{}}
	default:
		return Scheme{Other: &scheme}
	}
}

// mapGoErerToWasiHTTPErr 将 Go 的 net/http 和 net 错误映射到 wasi:http 的 ErrorCode。
func mapGoErerToWasiHTTPErr(err error) ErrorCode {
	if err == nil {
		return ErrorCode{}
	}

	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		if dnsErr.IsTimeout {
			return ErrorCode{DNSTimeout: &witgo.Unit{}}
		}
		if dnsErr.IsNotFound {
			return ErrorCode{DestinationNotFound: &witgo.Unit{}}
		}
		return ErrorCode{DNSError: &DNSErrorPayload{
			Rcode:    witgo.Some(dnsErr.Err),
			InfoCode: witgo.None[uint16](),
		}}
	}

	for i := 0; i < 8 && err != nil; i++ {
		var opErr *net.OpError
		if !errors.As(err, &opErr) {
			break
		}
		if opErr.Timeout() {
			return ErrorCode{ConnectionTimeout: &witgo.Unit{}}
		}
		// Err 为空时再递归会得到空 ErrorCode；Err 指向自己会爆栈。
		if opErr.Err == nil || opErr.Err == error(opErr) {
			break
		}
		err = opErr.Err
	}

	var nerr net.Error
	if errors.As(err, &nerr) && nerr.Timeout() {
		return ErrorCode{ConnectionTimeout: &witgo.Unit{}}
	}

	var syscallErr *os.SyscallError
	if errors.As(err, &syscallErr) {
		if code, ok := mapHTTPSysErrno(syscallErr.Err); ok {
			return code
		}
	}
	if code, ok := mapHTTPSysErrno(err); ok {
		return code
	}

	errMsg := err.Error()
	return ErrorCode{InternalError: witgo.SomePtr(errMsg)}
}

func headerValuesGet(h manager_http.Fields, name string) string {
	if h == nil {
		return ""
	}
	if v := h[strings.ToLower(name)]; len(v) > 0 {
		return v[0]
	}
	// fields 用 ToLower 存键，http.Header.Get 走 CanonicalMIMEHeaderKey，Content-Length 会读空
	return h.Get(name)
}

// cloneFieldsLower 把 net/http 的规范键转成 WASI fields 使用的小写键。
// http.Header.Clone 保留 "Content-Type"，fields.Get 用 ToLower 查找会 miss。
func cloneFieldsLower(h http.Header) manager_http.Fields {
	if h == nil {
		return nil
	}
	out := make(manager_http.Fields, len(h))
	for k, v := range h {
		nk := strings.ToLower(k)
		cp := make([]string, len(v))
		copy(cp, v)
		out[nk] = append(out[nk], cp...)
	}
	return out
}
