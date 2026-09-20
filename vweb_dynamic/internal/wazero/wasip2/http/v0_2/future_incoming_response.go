package v0_2

import (
	"context"

	manager_http "github.com/456vv/x/vweb_dynamic/internal/wazero/manager/http"
	manager_io "github.com/456vv/x/vweb_dynamic/internal/wazero/manager/io"
	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"
)

type futureIncomingResponseImpl struct {
	hm *manager_http.HTTPManager
}

func newFutureIncomingResponseImpl(hm *manager_http.HTTPManager) *futureIncomingResponseImpl {
	return &futureIncomingResponseImpl{hm: hm}
}

// Drop 是 future-incoming-response 资源的析构函数。
func (i *futureIncomingResponseImpl) Drop(_ context.Context, handle FutureIncomingResponse) {
	i.hm.Futures.Remove(handle)
}

// Subscribe 实现了 [method]future-incoming-response.subscribe。
func (i *futureIncomingResponseImpl) Subscribe(_ context.Context, this FutureIncomingResponse) Pollable {
	future, ok := i.hm.Futures.Get(this)
	if !ok {
		// 对于无效句柄，返回一个立即就绪的 pollable
		return i.hm.Poll.Add(manager_io.NewReadyPollable())
	}

	return i.hm.Poll.Add(future.Pollable)
}

// Get implements [method]future-incoming-response.get.
// It returns the response at most once.
func (i *futureIncomingResponseImpl) Get(
	ctx context.Context,
	this FutureIncomingResponse,
) witgo.Option[witgo.Result[witgo.Result[IncomingResponse, ErrorCode], witgo.Unit]] {
	future, ok := i.hm.Futures.Get(this)
	if !ok {
		// Invalid handle, return None. The WIT doesn't specify an error here.
		return witgo.None[witgo.Result[witgo.Result[IncomingResponse, ErrorCode], witgo.Unit]]()
	}

	select {
	case <-future.Pollable.Channel():
	case <-ctx.Done():
		return witgo.None[witgo.Result[witgo.Result[IncomingResponse, ErrorCode], witgo.Unit]]()
	}

	if !future.Consumed.CompareAndSwap(false, true) {
		// It was already consumed. Return Some(Err()).
		outerResult := witgo.Err[witgo.Result[IncomingResponse, ErrorCode], witgo.Unit](witgo.Unit{})
		return witgo.Some(outerResult)
	}

	var innerResult witgo.Result[IncomingResponse, ErrorCode]
	res := future.LoadResult() // 与 executeRequest 无锁写 Result 数据竞争
	if res.Err != nil {
		innerResult = witgo.Err[IncomingResponse, ErrorCode](mapGoErrToWasiHttpErr(res.Err))
	} else if res.Response == nil {
		// 在 Client.Do 异常路径可能 Err 与 Response 皆空，解引用会 panic
		innerResult = witgo.Err[IncomingResponse, ErrorCode](ErrorCode{InternalError: witgo.SomePtr("empty http response")})
	} else {
		responseHandle := i.hm.Responses.Add(&manager_http.IncomingResponse{
			Response:   res.Response,
			StatusCode: res.Response.StatusCode,
			// 与 fields.Get 的小写键约定对齐，且不与 net/http 共享底层 map。
			Headers: cloneFieldsLower(res.Response.Header),
		})
		innerResult = witgo.Ok[IncomingResponse, ErrorCode](responseHandle)
	}

	// Wrap the inner result in Ok() to signify a successful 'get' operation.
	outerResult := witgo.Ok[witgo.Result[IncomingResponse, ErrorCode], witgo.Unit](innerResult)
	return witgo.Some(outerResult)
}
