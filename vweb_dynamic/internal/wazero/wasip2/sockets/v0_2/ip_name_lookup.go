package v0_2

import (
	"context"
	"net"

	manager_io "github.com/456vv/x/vweb_dynamic/internal/wazero/manager/io"
	"github.com/456vv/x/vweb_dynamic/internal/wazero/manager/sockets"
	"github.com/456vv/x/vweb_dynamic/internal/wazero/wasip2"
	wasip2_io "github.com/456vv/x/vweb_dynamic/internal/wazero/wasip2/io/v0_2"
	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"
)

type ipNameLookupImpl struct {
	host *wasip2.Host
}

func newIPNameLookupImpl(h *wasip2.Host) *ipNameLookupImpl {
	return &ipNameLookupImpl{host: h}
}

func (i *ipNameLookupImpl) ResolveAddresses(ctx context.Context, network Network, name string) witgo.Result[ResolveAddressStream, ErrorCode] {
	if name == "" {
		return witgo.Err[ResolveAddressStream, ErrorCode](ErrorCodeInvalidArgument)
	}
	state := &sockets.ResolveAddressStreamState{
		Done: make(chan struct{}),
	}
	// 不能用 host 调用的 ctx，函数返回后常被取消，DNS 会立刻失败；
	// 用可 drop 取消的独立 context。SetCancel 跨包登记未导出的 cancel。
	lookupCtx, cancel := context.WithCancel(context.Background())
	state.SetCancel(cancel)
	handle := i.host.ResolveAddressStreamManager().Add(state)

	if ip := net.ParseIP(name); ip != nil {
		cancel() //  不走 LookupIPAddr 必须立刻 cancel，避免 context/goroutine 泄漏
		state.Addresses = []net.IP{ip}
		close(state.Done)
		return witgo.Ok[ResolveAddressStream, ErrorCode](handle)
	}

	go func() {
		defer close(state.Done)
		addrs, err := net.DefaultResolver.LookupIPAddr(lookupCtx, name)
		state.Mu.Lock()
		defer state.Mu.Unlock()
		if err != nil {
			state.Error = err
			return
		}
		state.Addresses = make([]net.IP, len(addrs))
		for i, ipAddr := range addrs {
			state.Addresses[i] = ipAddr.IP
		}
	}()

	return witgo.Ok[ResolveAddressStream, ErrorCode](handle)
}

func (i *ipNameLookupImpl) DropResolveAddressStream(_ context.Context, handle ResolveAddressStream) {
	i.host.ResolveAddressStreamManager().Remove(handle)
}

func (i *ipNameLookupImpl) ResolveNextAddress(_ context.Context, this ResolveAddressStream) witgo.Result[witgo.Option[IPAddress], ErrorCode] {
	state, ok := i.host.ResolveAddressStreamManager().Get(this)
	if !ok {
		return witgo.Err[witgo.Option[IPAddress], ErrorCode](ErrorCodeInvalidArgument)
	}

	select {
	case <-state.Done:
	default:
		return witgo.Err[witgo.Option[IPAddress], ErrorCode](ErrorCodeWouldBlock)
	}

	state.Mu.Lock()
	defer state.Mu.Unlock()
	if state.Error != nil {
		return witgo.Err[witgo.Option[IPAddress], ErrorCode](mapDnsError(state.Error))
	}

	// 转换失败不再递归（可栈溢出）；持锁推进 Index
	for state.Index < len(state.Addresses) {
		addr := state.Addresses[state.Index]
		state.Index++
		wasiAddr, err := toIPAddress(addr)
		if err != nil {
			continue
		}
		return witgo.Ok[witgo.Option[IPAddress], ErrorCode](witgo.Some(wasiAddr))
	}
	return witgo.Ok[witgo.Option[IPAddress], ErrorCode](witgo.None[IPAddress]())
}

func (i *ipNameLookupImpl) Subscribe(_ context.Context, this ResolveAddressStream) wasip2_io.Pollable {
	state, ok := i.host.ResolveAddressStreamManager().Get(this)
	if !ok {
		// 如果句柄无效，返回一个立即就绪的 pollable
		return i.host.PollManager().Add(manager_io.NewReadyPollable())
	}

	p := manager_io.NewPollableByChan(state.Done, nil)
	handle := i.host.PollManager().Add(p)
	return handle
}
