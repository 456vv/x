package v0_2

import (
	"context"
	"reflect"

	"github.com/456vv/x/vweb_dynamic/internal/wazero/manager/io"
)

// pollImpl 结构体持有 wasi:io/poll 的具体实现逻辑。
type pollImpl struct {
	pm *io.PollManager
}

func newPollImpl(pm *io.PollManager) *pollImpl {
	return &pollImpl{pm: pm}
}

// DropPollable 是 pollable 资源的析构函数。
func (i *pollImpl) DropPollable(_ context.Context, handle Pollable) {
	i.pm.Remove(handle)
}

// Ready 实现 [method]pollable.ready 方法。
func (i *pollImpl) Ready(_ context.Context, this Pollable) bool {
	p, ok := i.pm.Get(this)
	if !ok {
		return true // 无效句柄被认为是“就绪”的，以便调用者可以发现错误。
	}
	return p.IsReady()
}

// Block 实现 [method]pollable.block 方法。
func (i *pollImpl) Block(ctx context.Context, this Pollable) {
	p, ok := i.pm.Get(this)
	if !ok {
		return
	}
	select {
	case <-p.Channel():
	case <-ctx.Done(): // 忽略 ctx 时 guest 取消后 host 永久阻塞
	}
}

func (i *pollImpl) Poll(ctx context.Context, handles []Pollable) []uint32 {
	if len(handles) == 0 {
		return nil
	}
	for {
		if ctx != nil {
			if err := ctx.Err(); err != nil {
				return nil
			}
		}
		var readyIndexes []uint32
		for j, handle := range handles {
			if i.Ready(ctx, handle) {
				readyIndexes = append(readyIndexes, uint32(j))
			}
		}
		if len(readyIndexes) > 0 {
			return readyIndexes
		}

		cases := make([]reflect.SelectCase, 0, len(handles)+1)
		for _, handle := range handles {
			var ch <-chan struct{}
			if p, ok := i.pm.Get(handle); ok && p != nil {
				ch = p.Channel()
			} else {
				ready := make(chan struct{})
				close(ready)
				ch = ready
			}
			cases = append(cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)})
		}
		if ctx != nil {
			if done := ctx.Done(); done != nil {
				cases = append(cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(done)})
			}
		}
		chosen, _, _ := reflect.Select(cases)
		if chosen >= len(handles) {
			return nil
		}
	}
}
