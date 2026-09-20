package v0_2

import (
	manager_http "github.com/456vv/x/vweb_dynamic/internal/wazero/manager/http"
	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"
)

type responseOutparamImpl struct {
	hm *manager_http.HTTPManager
}

func newResponseOutparamImpl(hm *manager_http.HTTPManager) *responseOutparamImpl {
	return &responseOutparamImpl{hm: hm}
}

func (i *responseOutparamImpl) Drop(this ResponseOutparam) {
	i.hm.ResponseOutparams.Remove(this)
}

func (i *responseOutparamImpl) Set(param ResponseOutparam, response witgo.Result[OutgoingResponse, ErrorCode]) {
	p, ok := i.hm.ResponseOutparams.Pop(param)
	if !ok || p == nil || p.ResultChan == nil {
		return
	}
	var v any
	if response.Err != nil {
		v = *response.Err
	} else if response.Ok != nil {
		v = *response.Ok
	} else {
		return
	}
	// Drop 可能已 close(ResultChan)；重复 set 时缓冲已满。send 关闭 channel 会 panic。
	defer func() { _ = recover() }()
	select {
	case p.ResultChan <- v:
	default:
	}
}
