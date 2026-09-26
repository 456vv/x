package v0_2

import (
	"context"
	"reflect"
	"time"

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

func (i *pollImpl) Ready(_ context.Context, this Pollable) bool {
	p, ok := i.pm.Get(this)
	// 句柄在表里但资源是 nil 时，直接 IsReady 会 panic。
	if !ok || p == nil {
		return true
	}
	return p.IsReady()
}

func (i *pollImpl) Block(ctx context.Context, this Pollable) {
	p, ok := i.pm.Get(this)
	if !ok || p == nil {
		return
	}
	if ctx == nil {
		p.Block()
		return
	}
	for !p.IsReady() {
		select {
		case <-p.Channel():
		case <-ctx.Done():
			return
		}
	}
}

func (i *pollImpl) Poll(ctx context.Context, handles []Pollable) []uint32 {
	if len(handles) == 0 {
		return nil
	}
	// 修改原因：reflect.Select 最多 65536 个 case。截断等待集后必须定时重扫全部 IsReady，
	// 否则被裁掉的 pollable 永远无法把 poll 唤醒。
	const maxWait = 65533
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

		waitN := len(handles)
		needTimer := false
		if waitN > maxWait {
			waitN = maxWait
			needTimer = true
		}

		ncase := waitN
		hasCtx := ctx != nil && ctx.Done() != nil
		if hasCtx {
			ncase++
		}
		if needTimer {
			ncase++
		}
		cases := make([]reflect.SelectCase, 0, ncase)
		for _, handle := range handles[:waitN] {
			var ch <-chan struct{}
			if p, ok := i.pm.Get(handle); ok && p != nil {
				ch = p.Channel()
			}
			if ch == nil {
				ready := make(chan struct{})
				close(ready)
				ch = ready
			}
			cases = append(cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)})
		}
		ctxIdx := -1
		if hasCtx {
			ctxIdx = len(cases)
			cases = append(cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ctx.Done())})
		}
		var timer *time.Timer
		if needTimer {
			timer = time.NewTimer(50 * time.Millisecond)
			cases = append(cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(timer.C)})
		}
		chosen, _, _ := reflect.Select(cases)
		if timer != nil && !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
		if chosen == ctxIdx {
			return nil
		}
		// pollable 或 timer：回到循环做全量 IsReady
	}
}
