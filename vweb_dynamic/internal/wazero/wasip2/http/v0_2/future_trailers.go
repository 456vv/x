package v0_2

import (
	"context"

	manager_http "github.com/456vv/x/vweb_dynamic/internal/wazero/manager/http"
	manager_io "github.com/456vv/x/vweb_dynamic/internal/wazero/manager/io"
	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"
)

type futureTrailersImpl struct {
	hm *manager_http.HTTPManager
}

func newFutureTrailersImpl(hm *manager_http.HTTPManager) *futureTrailersImpl {
	return &futureTrailersImpl{hm: hm}
}

func (i *futureTrailersImpl) Drop(this FutureTrailers) {
	i.hm.FutureTrailers.Remove(this)
}

func (i *futureTrailersImpl) Subscribe(this FutureTrailers) Pollable {
	future, ok := i.hm.FutureTrailers.Get(this)
	if !ok || future == nil || future.Pollable == nil {
		// 对于无效句柄，返回一个立即就绪的 pollable
		return i.hm.Poll.Add(manager_io.NewReadyPollable())
	}

	// 与 future-incoming-response 相同，不能把内部 Pollable 直接交给 PollManager 析构。
	return i.hm.Poll.Add(manager_io.NewLevelPollable(future.Pollable.IsReady, future.Pollable))
}

func (i *futureTrailersImpl) Get(ctx context.Context, this FutureTrailers) witgo.Option[witgo.Result[witgo.Result[witgo.Option[Trailers], ErrorCode], witgo.Unit]] {
	_ = ctx // 与 future-incoming-response.get 相同：非阻塞，不因 ctx 丢就绪结果
	future, ok := i.hm.FutureTrailers.Get(this)
	if !ok || future == nil {
		return witgo.None[witgo.Result[witgo.Result[witgo.Option[Trailers], ErrorCode], witgo.Unit]]()
	}

	if future.Pollable == nil || !future.Pollable.IsReady() {
		return witgo.None[witgo.Result[witgo.Result[witgo.Option[Trailers], ErrorCode], witgo.Unit]]()
	}

	if !future.Consumed.CompareAndSwap(false, true) {
		// 与 future-incoming-response 对齐，重复 get 为 Some(Err(unit)) 而非 None
		outer := witgo.Err[witgo.Result[witgo.Option[Trailers], ErrorCode], witgo.Unit](witgo.Unit{})
		return witgo.Some(outer)
	}

	res := future.LoadResult()
	if res.Err != nil {
		// 读取 body 过程中发生错误
		errorCode := ErrorCode{InternalError: witgo.SomePtr(res.Err.Error())}
		return witgo.Some(witgo.Ok[witgo.Result[witgo.Option[Trailers], ErrorCode], witgo.Unit](
			witgo.Err[witgo.Option[Trailers], ErrorCode](errorCode),
		))
	}

	var trailers witgo.Option[Trailers]
	if res.Trailers != nil {
		handle := i.hm.Fields.Add(res.Trailers)
		// future-trailers.get 返回的 option<trailers> 同样是不可变 fields
		i.hm.MarkFieldsImmutable(handle)
		trailers = witgo.Some(handle)
	} else {
		trailers = witgo.None[Trailers]()
	}

	return witgo.Some(witgo.Ok[witgo.Result[witgo.Option[Trailers], ErrorCode], witgo.Unit](
		witgo.Ok[witgo.Option[Trailers], ErrorCode](trailers),
	))
}
