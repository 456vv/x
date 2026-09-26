package io

import (
	"sync"

	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"
)

type IPollable interface {
	// IsReady 以非阻塞方式检查 Pollable 是否就绪。
	IsReady() bool
	// Block 阻塞直到 Pollable 就绪。
	Block()
	// Channel 返回内部的 channel，用于 reflect.Select。
	Channel() <-chan struct{}
	// Close 用于释放与 Pollable 关联的资源，例如取消底层的定时器。
	Close()
}

// ChannelPollable 是 IPollable 接口的一个具体实现，它使用 channel 来进行阻塞。
// 这个实现是线程安全的。
type ChannelPollable struct {
	mu        sync.Mutex
	readyChan chan struct{}
	cancel    func() // 用于 Close()
	closed    bool   // Close 幂等，避免重复调用 cancel
	sticky    bool   // 单例 ReadyPollable 禁止 Reset/Close 破坏共享引用
}

// NewPollable 创建一个新的 channelPollable 实例。
func NewPollable(cancel func()) *ChannelPollable {
	return &ChannelPollable{
		readyChan: make(chan struct{}),
		cancel:    cancel,
	}
}

func NewPollableByChan(c chan struct{}, cancel func()) *ChannelPollable {
	// 外部传入 nil channel 时用空 channel，避免 select 永久阻塞在 nil 上语义不清
	if c == nil {
		c = make(chan struct{})
	}
	// 不得把调用方的 channel 当作 readyChan。Close→SetReady 会 close 它：
	// 1) DNS ResolveAddressStreamState.Done 在 lookup 结束时还会 close → 双重 close panic
	// 2) TCP ConnectDone 被 socket/析构/多个 subscribe 共享 → drop 某一个 pollable 会让 finish-connect 误判完成
	stop := make(chan struct{})
	var once sync.Once
	p := NewPollable(func() {
		once.Do(func() { close(stop) })
		if cancel != nil {
			cancel()
		}
	})
	select {
	case <-c:
		p.SetReady()
	default:
		go func() {
			select {
			case <-c:
				p.SetReady()
			case <-stop:
			}
		}()
	}
	return p
}

var ReadyPollable = NewReadyPollable()

// NewReadyPollable 创建一个已经处于“就绪”状态的 ChannelPollable。
func NewReadyPollable() *ChannelPollable {
	ch := make(chan struct{})
	close(ch)
	return &ChannelPollable{
		readyChan: ch,
		sticky:    true,
	}
}

// IsReady 以非阻塞方式检查 Pollable 是否就绪。
func (p *ChannelPollable) IsReady() bool {
	if p == nil {
		return true // 无效对象视为就绪，便于调用方发现后续错误
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	// 必须在锁内观察 channel，避免与 Reset 交错误判
	select {
	case <-p.readyChan:
		return true
	default:
		return false
	}
}

// Block 阻塞直到 Pollable 就绪。
func (p *ChannelPollable) Block() {
	if p == nil {
		return
	}
	// Reset 会换新 channel；必须在锁内拿到“当前视图”并确认仍未就绪后再等，
	// 否则只等到已关闭的旧 channel 后误判就绪，或与 SetReady 交错丢失唤醒。
	for {
		p.mu.Lock()
		if p.sticky {
			p.mu.Unlock()
			return
		}
		select {
		case <-p.readyChan:
			p.mu.Unlock()
			return
		default:
			ch := p.readyChan
			p.mu.Unlock()
			<-ch
		}
	}
}

// SetReady 将 Pollable 状态设置为就绪。这个操作是幂等的。
func (p *ChannelPollable) SetReady() {
	if p == nil {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.sticky {
		return
	}
	select {
	case <-p.readyChan:
		return
	default:
		close(p.readyChan)
	}
}

// Reset 将 Pollable 重置为“未就绪”状态，使其可以被再次使用。
func (p *ChannelPollable) Reset() {
	if p == nil {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.sticky {
		return // 全局 ReadyPollable 必须永远就绪
	}
	select {
	case <-p.readyChan:
		p.readyChan = make(chan struct{})
	default:
	}
}

// Channel 返回当前的内部 channel。
// 警告：返回的 channel 可能会在 Reset 调用后失效。
func (p *ChannelPollable) Channel() <-chan struct{} {
	if p == nil {
		return ReadyPollable.Channel()
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.readyChan
}

// Close 调用与此 pollable 关联的取消函数（如果存在）。
func (p *ChannelPollable) Close() {
	if p == nil {
		return
	}
	p.mu.Lock()
	// ReadyPollable 被 PollManager 多次 Add/Remove 时不得把全局单例标成 closed。
	if p.sticky || p.closed {
		p.mu.Unlock()
		return
	}
	p.closed = true
	cancel := p.cancel
	p.cancel = nil
	p.mu.Unlock()

	// 在锁外调用 cancel，避免 cancel 回调里再碰 pollable 造成死锁
	if cancel != nil {
		cancel()
	}
	// drop/Host.Close 必须唤醒 Block 与 poll.select；定时器 Stop 后原 channel 可能永不关闭
	p.SetReady()
}

// LevelPollable 按条件判断就绪（电平触发）。
// 粘滞 ChannelPollable 在数据被读空后仍 IsReady，或 Channel 已关但条件为假，
// 会让 wasi:io/poll.poll 忙等。IsReady 必须看实际条件；未就绪时 Channel 必须是未关闭通道。
type LevelPollable struct {
	check func() bool
	wake  *ChannelPollable
}

// NewLevelPollable 用 check 做电平就绪，用 wake 做阻塞唤醒（可 Reset）。
func NewLevelPollable(check func() bool, wake *ChannelPollable) *LevelPollable {
	if wake == nil {
		wake = NewPollable(nil)
	}
	return &LevelPollable{check: check, wake: wake}
}

func (p *LevelPollable) IsReady() bool {
	if p == nil {
		return true
	}
	if p.check != nil && p.check() {
		return true
	}
	return false
}

func (p *LevelPollable) Channel() <-chan struct{} {
	if p == nil {
		return ReadyPollable.Channel()
	}
	if p.IsReady() {
		return ReadyPollable.Channel()
	}
	var ch <-chan struct{}
	if p.wake != nil {
		ch = p.wake.Channel()
	} else {
		ch = ReadyPollable.Channel()
	}
	// 取 wake 与数据到达之间可能已就绪；再读一次避免丢失唤醒后仍等旧通道。
	if p.IsReady() {
		return ReadyPollable.Channel()
	}
	return ch
}

func (p *LevelPollable) Block() {
	if p == nil {
		return
	}
	for !p.IsReady() {
		<-p.Channel()
	}
}

func (p *LevelPollable) Close() {
	// wake 由 stream 共享，不能在 drop 某一个 subscribe 句柄时 Reset/Close 掉其它 waiter。
}

// PollManager 是用于管理所有 Pollable 资源的管理器。
type PollManager = witgo.ResourceManager[IPollable]

// NewPollManager 创建一个新的 Poll 管理器。
func NewPollManager() *PollManager {
	return witgo.NewResourceManager[IPollable](func(resource IPollable) {
		if resource != nil {
			resource.Close()
		}
	})
}
