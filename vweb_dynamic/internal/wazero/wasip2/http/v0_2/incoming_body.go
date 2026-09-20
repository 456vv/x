package v0_2

import (
	"bytes"
	"context"
	"io"

	manager_http "github.com/456vv/x/vweb_dynamic/internal/wazero/manager/http"
	manager_io "github.com/456vv/x/vweb_dynamic/internal/wazero/manager/io"
	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"
)

type incomingBodyImpl struct {
	hm *manager_http.HTTPManager
}

func newIncomingBodyImpl(hm *manager_http.HTTPManager) *incomingBodyImpl {
	return &incomingBodyImpl{hm: hm}
}

// Drop 是析构函数
func (i *incomingBodyImpl) Drop(_ context.Context, handle IncomingBody) {
	i.hm.IncomingBodies.Remove(handle)
}

// Stream 实现了 [method]incoming-body.stream。
// 将响应体转换为 wasi:io 的 input-stream。
func (i *incomingBodyImpl) Stream(_ context.Context, this IncomingBody) witgo.Result[InputStream, witgo.Unit] {
	body, ok := i.hm.IncomingBodies.Get(this)
	if !ok || body == nil {
		return witgo.Err[InputStream, witgo.Unit](witgo.Unit{})
	}

	if !body.Consumed.CompareAndSwap(false, true) {
		// 原先误写为 Err[OutgoingBody, Unit]，与方法返回类型 Result[InputStream, Unit] 不一致
		return witgo.Err[InputStream, witgo.Unit](witgo.Unit{})
	}

	// Stream 和 IncomingBody 生命周期绑定，这里不 Close 底层 Reader
	reader := body.Stream
	if reader == nil {
		// 空 body 给出立即 EOF 的 Reader，避免后台异步读 goroutine 对 nil 解引用
		reader = bytes.NewReader(nil)
	}
	stream := manager_io.NewAsyncStreamForReader(reader, manager_io.DontCloseReader())
	body.StreamHandle = i.hm.Streams.Add(stream)
	return witgo.Ok[InputStream, witgo.Unit](body.StreamHandle)
}

// Finish 是一个静态方法，消费 incoming-body 并返回 future-trailers。
func (i *incomingBodyImpl) Finish(_ context.Context, this IncomingBody) FutureTrailers {
	// 1. 获取 incoming-body 实例。
	body, ok := i.hm.IncomingBodies.Pop(this)
	if !ok || body == nil {
		// 根据 WIT 规范，无效句柄应该触发陷阱。
		// 在 wazero 的 Go host function 中，panic 会被转换为 trap。
		panic("invalid incoming-body handle")
	}

	// 必须先 Close 底层 body 才能打断 DontCloseReader 的 Read，再 Remove stream，否则死锁
	if body.StreamHandle != 0 {
		body.Close()
		i.hm.Streams.Remove(body.StreamHandle)
		body.StreamHandle = 0
	} else if body.Stream != nil {
		// net/http Trailer 要在 body 读到 EOF 后才填充
		io.Copy(io.Discard, body.Stream)
	}
	body.Close()

	if body.Consumed.CompareAndSwap(false, true) {
		// 根据 WIT，流仍存活时应 trap；此处为兼容不报错
	}

	future := &manager_http.FutureTrailers{
		Pollable: manager_io.ReadyPollable,
	}
	if body.GetTrailers != nil {
		future.StoreResult(manager_http.ResultTrailers{Trailers: body.GetTrailers()})
	}
	return i.hm.FutureTrailers.Add(future)
}
