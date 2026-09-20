package v0_2

import (
	"math"
	"time"

	io_v0_2 "github.com/456vv/x/vweb_dynamic/internal/wazero/wasip2/io/v0_2"
)

// Pollable 从 wasi:io 导入
type Pollable = io_v0_2.Pollable

// --- monotonic-clock types ---
type (
	Instant  = uint64
	Duration uint64
)

func (d Duration) ToTime() time.Time {
	return time.Unix(0, int64(d.ToDuration()))
}

func (d Duration) ToDuration() time.Duration {
	// uint64 纳秒转 int64 溢出后变成负超时
	if uint64(d) > uint64(math.MaxInt64) {
		return time.Duration(math.MaxInt64)
	}
	return time.Duration(d)
}

// --- wall-clock types ---
type Datetime struct {
	Seconds     uint64
	Nanoseconds uint32
}
