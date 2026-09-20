package v0_2

import (
	"context"
	"math"
	"time"

	manager_http "github.com/456vv/x/vweb_dynamic/internal/wazero/manager/http"
	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"
)

type requestOptionsImpl struct {
	hm *manager_http.HTTPManager
}

func newRequestOptionsImpl(hm *manager_http.HTTPManager) *requestOptionsImpl {
	return &requestOptionsImpl{hm: hm}
}

func (i *requestOptionsImpl) Constructor() RequestOptions {
	return i.hm.Options.Add(&manager_http.RequestOptions{})
}

func (i *requestOptionsImpl) Drop(_ context.Context, handle RequestOptions) {
	i.hm.Options.Remove(handle)
}

// durationFromOption 把 WASI Duration(ns) 转成 *time.Duration。
// uint64 ns → int64 溢出后会变成负超时。
func durationFromOption(duration witgo.Option[Duration]) *time.Duration {
	if duration.Some == nil {
		return nil
	}
	ns := uint64(*duration.Some)
	var d time.Duration
	if ns > uint64(math.MaxInt64) {
		d = time.Duration(math.MaxInt64)
	} else {
		d = time.Duration(ns)
	}
	return &d
}

// --- Getters ---
func (i *requestOptionsImpl) ConnectTimeout(_ context.Context, this RequestOptions) witgo.Option[Duration] {
	opts, ok := i.hm.Options.Get(this)
	if !ok || opts.ConnectTimeout == nil {
		return witgo.None[Duration]()
	}
	return witgo.Some(Duration(*opts.ConnectTimeout))
}

func (i *requestOptionsImpl) FirstByteTimeout(_ context.Context, this RequestOptions) witgo.Option[Duration] {
	opts, ok := i.hm.Options.Get(this)
	if !ok || opts.FirstByteTimeout == nil {
		return witgo.None[Duration]()
	}
	return witgo.Some(Duration(*opts.FirstByteTimeout))
}

func (i *requestOptionsImpl) BetweenBytesTimeout(_ context.Context, this RequestOptions) witgo.Option[Duration] {
	opts, ok := i.hm.Options.Get(this)
	if !ok || opts.BetweenBytesTimeout == nil {
		return witgo.None[Duration]()
	}
	return witgo.Some(Duration(*opts.BetweenBytesTimeout))
}

// --- Setters ---
func (i *requestOptionsImpl) SetConnectTimeout(_ context.Context, this RequestOptions, duration witgo.Option[Duration]) witgo.UnitResult {
	opts, ok := i.hm.Options.Get(this)
	if !ok {
		return witgo.UintErr()
	}
	opts.ConnectTimeout = durationFromOption(duration)
	return witgo.UintOk()
}

func (i *requestOptionsImpl) SetFirstByteTimeout(_ context.Context, this RequestOptions, duration witgo.Option[Duration]) witgo.UnitResult {
	opts, ok := i.hm.Options.Get(this)
	if !ok {
		return witgo.UintErr()
	}
	opts.FirstByteTimeout = durationFromOption(duration)
	return witgo.UintOk()
}

func (i *requestOptionsImpl) SetBetweenBytesTimeout(_ context.Context, this RequestOptions, duration witgo.Option[Duration]) witgo.UnitResult {
	opts, ok := i.hm.Options.Get(this)
	if !ok {
		return witgo.UintErr()
	}
	opts.BetweenBytesTimeout = durationFromOption(duration)
	return witgo.UintOk()
}
