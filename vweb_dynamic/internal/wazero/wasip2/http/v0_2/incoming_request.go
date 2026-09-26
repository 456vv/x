package v0_2

import (
	"io"

	manager_http "github.com/456vv/x/vweb_dynamic/internal/wazero/manager/http"
	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"
)

type incomingRequestImpl struct {
	hm *manager_http.HTTPManager
}

func newIncomingRequestImpl(hm *manager_http.HTTPManager) *incomingRequestImpl {
	return &incomingRequestImpl{hm: hm}
}

func (i *incomingRequestImpl) Drop(this IncomingRequest) {
	i.hm.IncomingRequests.Remove(this)
}

func (i *incomingRequestImpl) Method(this IncomingRequest) Method {
	req, ok := i.hm.IncomingRequests.Get(this)
	if !ok || req == nil { // Get 成功但资源为 nil 时解引用会 panic
		return Method{Other: new(string)}
	}
	return toWasiMethod(req.Method)
}

func (i *incomingRequestImpl) PathWithQuery(this IncomingRequest) witgo.Option[string] {
	req, ok := i.hm.IncomingRequests.Get(this)
	if !ok || req == nil {
		return witgo.None[string]()
	}
	path := req.Path
	if req.Query != "" {
		path += "?" + req.Query
	}
	return witgo.Some(path)
}

func (i *incomingRequestImpl) Scheme(this IncomingRequest) witgo.Option[Scheme] {
	req, ok := i.hm.IncomingRequests.Get(this)
	if !ok || req == nil {
		return witgo.None[Scheme]()
	}
	if req.Scheme == nil {
		return witgo.None[Scheme]()
	}
	return witgo.Some(toWasiScheme(*req.Scheme))
}

func (i *incomingRequestImpl) Authority(this IncomingRequest) witgo.Option[string] {
	req, ok := i.hm.IncomingRequests.Get(this)
	if !ok || req == nil {
		return witgo.None[string]()
	}
	if req.Authority == nil {
		return witgo.None[string]()
	}
	return witgo.Some(*req.Authority)
}

func (i *incomingRequestImpl) Headers(this IncomingRequest) Headers {
	req, ok := i.hm.IncomingRequests.Get(this)
	if !ok || req == nil {
		return 0
	}
	return req.LoadHeadersHandle()
}

func (i *incomingRequestImpl) Consume(this IncomingRequest) witgo.Result[IncomingBody, witgo.Unit] {
	req, ok := i.hm.IncomingRequests.Get(this)
	if !ok || req == nil {
		return witgo.Err[IncomingBody, witgo.Unit](witgo.Unit{})
	}
	// CAS/Add/BodyHandle 必须在 manager 锁内完成，避免与 resource-drop 交错。
	handle, ok := req.InstallIncomingBody(func(rc io.ReadCloser) uint32 {
		body := &manager_http.IncomingBody{
			Stream: rc,
			GetTrailers: func() manager_http.Fields {
				if req.Request == nil {
					return nil
				}
				// net/http Trailer 是规范键，fields.Get 用小写查找会 miss
				return cloneFieldsLower(req.Request.Trailer)
			},
		}
		return i.hm.IncomingBodies.Add(body)
	})
	if !ok {
		return witgo.Err[IncomingBody, witgo.Unit](witgo.Unit{})
	}
	return witgo.Ok[IncomingBody, witgo.Unit](handle)
}
