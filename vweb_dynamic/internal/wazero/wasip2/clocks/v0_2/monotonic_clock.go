package v0_2

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/456vv/x/vweb_dynamic/internal/wazero/manager/io"
)

// Monotonic clock zero value, initialized at package load time.
var programStart = time.Now()

type monotonicClockImpl struct {
	pm *io.PollManager
}

func newMonotonicClockImpl(pm *io.PollManager) *monotonicClockImpl {
	return &monotonicClockImpl{pm: pm}
}

// Now returns the current time from the monotonic clock in nanoseconds.
func (i *monotonicClockImpl) Now(_ context.Context) Instant {
	return Instant(time.Since(programStart).Nanoseconds())
}

// Resolution returns the resolution of the monotonic clock.
func (i *monotonicClockImpl) Resolution(_ context.Context) Duration {
	// Go's time resolution is 1 nanosecond.
	return 1
}

func (i *monotonicClockImpl) subscribeTimer(d time.Duration) Pollable {
	if d <= 0 {
		return i.pm.Add(io.ReadyPollable)
	}
	timer := time.NewTimer(d)
	done := make(chan struct{})
	var once sync.Once
	p := io.NewPollable(func() {
		once.Do(func() { close(done) })
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
	})
	handle := i.pm.Add(p)
	go func() {
		select {
		case <-timer.C:
			p.SetReady()
		case <-done:
			// pollable 被 drop 且 Stop 成功时，原 goroutine 会永久阻塞在 timer.C
		}
	}()
	return handle
}

func durationFromWasi(ns uint64) time.Duration {
	if ns > uint64(math.MaxInt64) {
		return time.Duration(math.MaxInt64)
	}
	return time.Duration(ns)
}

func (i *monotonicClockImpl) SubscribeInstant(_ context.Context, when Instant) Pollable {
	now := i.Now(context.Background())
	if when <= now {
		return i.pm.Add(io.ReadyPollable)
	}
	return i.subscribeTimer(durationFromWasi(uint64(when - now)))
}

func (i *monotonicClockImpl) SubscribeDuration(_ context.Context, when Duration) Pollable {
	return i.subscribeTimer(durationFromWasi(uint64(when)))
}
