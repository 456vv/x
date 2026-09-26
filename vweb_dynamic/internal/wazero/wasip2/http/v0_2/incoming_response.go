package v0_2

import (
	"context"
	"io"

	manager_http "github.com/456vv/x/vweb_dynamic/internal/wazero/manager/http"
	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"
)

type incomingResponseImpl struct {
	hm *manager_http.HTTPManager
}

func newIncomingResponseImpl(hm *manager_http.HTTPManager) *incomingResponseImpl {
	return &incomingResponseImpl{hm: hm}
}

// Drop 是析构函数
func (i *incomingResponseImpl) Drop(_ context.Context, handle IncomingResponse) {
	i.hm.Responses.Remove(handle)
}

// Status 实现了 [method]incoming-response.status。
func (i *incomingResponseImpl) Status(_ context.Context, this IncomingResponse) uint16 {
	resp, ok := i.hm.Responses.Get(this)
	if !ok || resp == nil || resp.Response == nil {
		panic("invalid incoming-respone handle")
	}
	// consume 会写 Response.Body；与读 StatusCode 并发是 data race。创建时已拷到 StatusCode。
	return uint16(resp.StatusCode)
}

// Headers 实现了 [method]incoming-response.headers。
func (i *incomingResponseImpl) Headers(_ context.Context, this IncomingResponse) Fields {
	resp, ok := i.hm.Responses.Get(this)
	if !ok || resp == nil || resp.Response == nil {
		panic("invalid incoming-respone handle")
	}

	// WASI 规定 headers 是 child resource，多次调用必须返回同一句柄。
	resp.EnsureHeaderHandle(func() {
		if resp.LoadHeaderHandle() != 0 {
			return
		}
		src := resp.Headers
		if src == nil && resp.Response != nil {
			src = cloneFieldsLower(resp.Response.Header)
		} else {
			src = cloneFieldsLower(src)
		}
		handle := i.hm.Fields.Add(src)
		i.hm.MarkFieldsImmutable(handle)
		// 表里是 clone，不能用 resp.Headers 做 RemoveIf 身份。
		resp.StoreHeaderHandle(handle, src)
	})
	return resp.LoadHeaderHandle()
}

// Consume 实现了 [method]incoming-response.consume。
// 返回一个 incoming-body，用于读取响应体。
func (i *incomingResponseImpl) Consume(_ context.Context, this IncomingResponse) witgo.Result[IncomingBody, witgo.Unit] {
	resp, ok := i.hm.Responses.Get(this)
	if !ok || resp == nil || resp.Response == nil {
		return witgo.Err[IncomingBody, witgo.Unit](witgo.Unit{})
	}

	// 与 incoming-request.consume 相同，锁内移交 Body。
	handle, ok := resp.InstallIncomingBody(func(rc io.ReadCloser) uint32 {
		body := &manager_http.IncomingBody{
			Stream: rc,
			GetTrailers: func() (trailers manager_http.Fields) {
				if resp.Response == nil {
					return nil
				}
				return cloneFieldsLower(resp.Response.Trailer)
			},
		}
		return i.hm.IncomingBodies.Add(body)
	})
	if !ok {
		return witgo.Err[IncomingBody, witgo.Unit](witgo.Unit{})
	}
	return witgo.Ok[IncomingBody, witgo.Unit](handle)
}
