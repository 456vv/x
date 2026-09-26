package v0_2

import (
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
	// CAS、创建异步流、写 StreamHandle 必须原子完成。
	handle, ok := body.InstallStream(func(r io.Reader) (uint32, *manager_io.Stream) {
		stream := manager_io.NewAsyncStreamForReader(r, manager_io.DontCloseReader())
		return i.hm.Streams.Add(stream), stream
	})
	if !ok {
		return witgo.Err[InputStream, witgo.Unit](witgo.Unit{})
	}
	return witgo.Ok[InputStream, witgo.Unit](handle)
}

// drainIncomingStream 阻塞读到 EOF，让 net/http 填好 Trailer。
// stream() 使用 DontCloseReader 的异步包装，Read 非阻塞且独占底层 Body；
// 原先只在未调用 stream() 时 Copy，走 stream() 后直接 Close，Trailer 经常为空。
func drainIncomingStream(s *manager_io.Stream) {
	if s == nil || s.Reader == nil {
		return
	}
	buf := make([]byte, 32*1024)
	for {
		n, err := s.Reader.Read(buf)
		if n == 0 && err == nil {
			if s.OnSubscribe == nil {
				return
			}
			p := s.OnSubscribe()
			if p == nil {
				return
			}
			p.Block()
			continue
		}
		if err != nil {
			return
		}
	}
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

	// 不能无锁读 StreamHandle；未调用 stream() 时仍排空原始 Body。
	if h, st := body.TakeStreamHandleForDrop(); h != 0 {
		// 句柄复用后 Get(h) 会 drain 并关掉另一个 guest 的 input-stream。
		if st != nil {
			if cur, ok := i.hm.Streams.Get(h); ok && cur == st {
				drainIncomingStream(st)
			}
			body.Close()
			i.hm.Streams.RemoveIf(h, func(cur *manager_io.Stream) bool { return cur == st })
		} else {
			body.Close()
		}
	} else if body.Stream != nil {
		io.Copy(io.Discard, body.Stream)
		body.Close()
	} else {
		body.Close()
	}

	future := &manager_http.FutureTrailers{
		Pollable: manager_io.ReadyPollable,
	}
	if body.GetTrailers != nil {
		future.StoreResult(manager_http.ResultTrailers{Trailers: body.GetTrailers()})
	}
	return i.hm.FutureTrailers.Add(future)
}
