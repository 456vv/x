package v0_2

import (
	"context"
	"time"
)

type wallClockImpl struct{}

func newWallClockImpl() *wallClockImpl {
	return &wallClockImpl{}
}

// Now returns the current wall-clock time.
func (i *wallClockImpl) Now(_ context.Context) Datetime {
	now := time.Now()
	sec := now.Unix()
	if sec < 0 {
		// 时钟被拨到 1970 前时 uint64(负数) 会绕成超大秒数。
		sec = 0
	}
	return Datetime{
		Seconds:     uint64(sec),
		Nanoseconds: uint32(now.Nanosecond()),
	}
}

// Resolution returns the resolution of the wall-clock.
func (i *wallClockImpl) Resolution(_ context.Context) Datetime {
	// Go's time resolution is typically 1 nanosecond on most systems.
	return Datetime{
		Seconds:     0,
		Nanoseconds: 1,
	}
}
