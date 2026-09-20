package v0_2

import (
	"context"

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
	return uint16(resp.Response.StatusCode)
}

// Headers 实现了 [method]incoming-response.headers。
func (i *incomingResponseImpl) Headers(_ context.Context, this IncomingResponse) Fields {
	resp, ok := i.hm.Responses.Get(this)
	if !ok || resp == nil || resp.Response == nil {
		panic("invalid incoming-respone handle")
	}

	// WASI 规定 headers 是 child resource，多次调用必须返回同一句柄。
	resp.EnsureHeaderHandle(func() {
		if resp.HeaderHandle != 0 {
			return
		}
		// http.Header.Clone 保留规范键，fields.Get 用 ToLower 会 miss
		cloned := cloneFieldsLower(resp.Response.Header)
		handle := i.hm.Fields.Add(cloned)
		i.hm.MarkFieldsImmutable(handle)
		resp.HeaderHandle = handle
	})
	return resp.HeaderHandle
}

// Consume 实现了 [method]incoming-response.consume。
// 返回一个 incoming-body，用于读取响应体。
func (i *incomingResponseImpl) Consume(_ context.Context, this IncomingResponse) witgo.Result[IncomingBody, witgo.Unit] {
	resp, ok := i.hm.Responses.Get(this)
	if !ok || resp == nil || resp.Response == nil {
		return witgo.Err[IncomingBody, witgo.Unit](witgo.Unit{})
	}
	if !resp.Consumed.CompareAndSwap(false, true) {
		// 原先误写为 Err[OutgoingBody, Unit]，与方法返回类型 Result[IncomingBody, Unit] 不一致
		return witgo.Err[IncomingBody, witgo.Unit](witgo.Unit{})
	}

	body := &manager_http.IncomingBody{
		Stream: resp.Response.Body,
		GetTrailers: func() (trailers manager_http.Fields) {
			if resp.Response == nil {
				return nil
			}
			return cloneFieldsLower(resp.Response.Trailer)
		},
	}
	resp.BodyHandle = i.hm.IncomingBodies.Add(body)
	return witgo.Ok[IncomingBody, witgo.Unit](resp.BodyHandle)
}
