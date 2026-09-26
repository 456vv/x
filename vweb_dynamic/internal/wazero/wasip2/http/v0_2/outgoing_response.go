package v0_2

import (
	"io"
	gohttp "net/http"
	"strconv"

	manager_http "github.com/456vv/x/vweb_dynamic/internal/wazero/manager/http"
	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"
)

type outgoingResponseImpl struct {
	hm *manager_http.HTTPManager
}

func newOutgoingResponseImpl(hm *manager_http.HTTPManager) *outgoingResponseImpl {
	return &outgoingResponseImpl{hm: hm}
}

func (i *outgoingResponseImpl) Constructor(headers Headers) OutgoingResponse {
	header, _ := i.hm.Fields.Pop(headers)
	resp := &manager_http.OutgoingResponse{
		StatusCode: gohttp.StatusOK,
		Headers:    header,
	}

	return i.hm.OutgoingResponses.Add(resp)
}

func (i *outgoingResponseImpl) Drop(this OutgoingResponse) {
	i.hm.OutgoingResponses.Remove(this)
}

func (i *outgoingResponseImpl) StatusCode(this OutgoingResponse) StatusCode {
	resp, ok := i.hm.OutgoingResponses.Get(this)
	if !ok || resp == nil {
		return 0
	}
	return StatusCode(resp.LoadStatus())
}

func (i *outgoingResponseImpl) SetStatusCode(this OutgoingResponse, statusCode StatusCode) witgo.Result[witgo.Unit, witgo.Unit] {
	resp, ok := i.hm.OutgoingResponses.Get(this)
	if !ok || resp == nil {
		return witgo.Err[witgo.Unit, witgo.Unit](witgo.Unit{})
	}
	// WASI status-code 只有 100–599 有效，其余必须返回错误而不是留给 server 改写。
	if statusCode < 100 || statusCode > 599 {
		return witgo.Err[witgo.Unit, witgo.Unit](witgo.Unit{})
	}
	resp.StoreStatus(int(statusCode))
	return witgo.Ok[witgo.Unit, witgo.Unit](witgo.Unit{})
}

func (i *outgoingResponseImpl) Headers(this OutgoingResponse) Headers {
	resp, ok := i.hm.OutgoingResponses.Get(this)
	if !ok || resp == nil {
		return 0
	}

	// WASI headers 是 child resource，多次调用必须返回同一句柄。
	// 并发 headers() 若各自 Add 会泄漏 Fields，且 guest 拿到的 handle 互不相等。
	// EnsureHeadersHandle 在 manager/http 包内跑 sync.Once，外包不能直接碰未导出的 headersOnce。
	resp.EnsureHeadersHandle(func() {
		if resp.LoadHeadersHandle() != 0 {
			return
		}
		if resp.Headers == nil {
			resp.Headers = make(manager_http.Fields)
		}
		resp.StoreHeadersHandle(i.hm.Fields.Add(manager_http.Fields(resp.Headers)))
	})
	return resp.LoadHeadersHandle()
}

func (i *outgoingResponseImpl) Body(this OutgoingResponse) witgo.Result[OutgoingBody, witgo.Unit] {
	resp, ok := i.hm.OutgoingResponses.Get(this)
	if !ok || resp == nil {
		return witgo.Err[OutgoingBody, witgo.Unit](witgo.Unit{})
	}

	handle, ok := resp.InstallOutgoingBody(func() (uint32, io.Reader, *io.PipeWriter) {
		var contentLength *uint64
		i.hm.LockFields()
		cl := headerValuesGet(resp.Headers, "Content-Length")
		i.hm.UnlockFields()
		if len(cl) > 0 {
			if val, err := strconv.ParseUint(cl, 10, 64); err == nil {
				contentLength = &val
			}
		}
		h, body, pw := i.hm.NewOutgoingBody(contentLength, func(trailers manager_http.Fields) error {
			resp.StoreTrailers(trailers)
			return nil
		})
		return h, body, pw
	})
	if !ok {
		return witgo.Err[OutgoingBody, witgo.Unit](witgo.Unit{})
	}
	return witgo.Ok[OutgoingBody, witgo.Unit](handle)
}
