package sockets

import (
	"context"
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"

	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"
)

// ErrInvalidSocketState 表示套接字处于不允许该操作的状态。
// wasip2 与 manager 不同包，不能直接碰未导出锁，用错误值区分 InvalidState。
var ErrInvalidSocketState = errors.New("socket invalid state")

// IPAddressFamily 是一个版本无关的 IP 地址族枚举。
type IPAddressFamily uint8

const (
	IPAddressFamilyIPV4 IPAddressFamily = iota
	IPAddressFamilyIPV6
)

type (
	IPv4Address = [4]byte
	IPv6Address = [8]uint16
)

type IPAddress struct {
	IPV4 *IPv4Address `wit:"case(0)"`
	IPV6 *IPv6Address `wit:"case(1)"`
}

type IPv4SocketAddress struct {
	Port    uint16
	Address IPv4Address
}

type IPv6SocketAddress struct {
	Port     uint16
	FlowInfo uint32
	Address  IPv6Address
	ScopeID  uint32
}

type IPSocketAddress struct {
	IPV4 *IPv4SocketAddress `wit:"case(0)"`
	IPV6 *IPv6SocketAddress `wit:"case(1)"`
}

// IncomingDatagram 代表一个接收到的 UDP 数据报。
type IncomingDatagram struct {
	Data          []byte
	RemoteAddress IPSocketAddress
}

// OutgoingDatagram 代表一个待发送的 UDP 数据报。
type OutgoingDatagram struct {
	Data          []byte
	RemoteAddress witgo.Option[IPSocketAddress]
}

// Network represents the capability to access the network.
type Network struct{}

// ConnectResult 用于在 goroutine 之间传递异步连接的结果。
type ConnectResult struct {
	Conn *net.TCPConn
	Err  error
}

// TCPSocket 代表一个 TCP 套接字资源。
type TCPSocket struct {
	Fd       int
	Listener *net.TCPListener
	Conn     *net.TCPConn
	Family   IPAddressFamily
	State    TCPState

	// ConnectResult 用于异步连接操作（保留字段，避免破坏已有引用）。
	ConnectResult chan ConnectResult

	// Subscribe 与 FinishConnect 不能争抢同一条缓冲 channel
	ConnectDone chan struct{}

	stateAtomic atomic.Uint32
	closeFd     func() error

	connectMu     sync.Mutex
	connectStored ConnectResult
	connectReady  atomic.Bool

	// drop 时取消 in-flight Dial/Poll；无原始 fd 时 DialTCP 无法被 Conn.Close 打断
	connectCtx    context.Context
	connectCancel context.CancelFunc
}

// TCPState represents the state of a TCP socket as defined in the WIT world.
type TCPState uint8

const (
	TCPStateUnbound TCPState = iota
	TCPStateBound
	TCPStateListening
	TCPStateConnecting
	TCPStateConnected
	TCPStateClosed
)

// HasFd 报告是否仍持有需由 closeFd 关闭的原始套接字。
// 注意：Fd==0 不是可靠哨兵（fd 0 可能是被回收的 stdin）。
func (s *TCPSocket) HasFd() bool { return s != nil && s.closeFd != nil }

// DetachFd 在所有权已交给 net.Conn/Listener 后清除 closeFd，不再关原 fd。
func (s *TCPSocket) DetachFd() {
	if s == nil {
		return
	}
	s.closeFd = nil
	s.Fd = -1
}

// GetState 返回当前 TCP 状态（优先原子字段）。
func (s *TCPSocket) GetState() TCPState {
	if s == nil {
		return TCPStateClosed
	}
	// Unbound==0 时也以 atomic 为准，避免无锁读 State 与 SetState 竞态。
	return TCPState(s.stateAtomic.Load())
}

// SetState 更新 TCP 状态。
func (s *TCPSocket) SetState(st TCPState) {
	if s == nil {
		return
	}
	s.State = st
	s.stateAtomic.Store(uint32(st))
}

// ConnectContext 返回当前 connect 的可取消 context；无则 Background。
// wasip2 与 manager 不同包，不能直接读未导出字段。
func (s *TCPSocket) ConnectContext() context.Context {
	if s == nil {
		return context.Background()
	}
	s.connectMu.Lock()
	defer s.connectMu.Unlock()
	if s.connectCtx == nil {
		return context.Background()
	}
	return s.connectCtx
}

// DoBind 在 Unbound 下执行 bind，成功则进入 Bound。持锁避免与 connect/drop 并发。
func (s *TCPSocket) DoBind(fn func() error) error {
	if s == nil {
		return ErrInvalidSocketState
	}
	s.connectMu.Lock()
	defer s.connectMu.Unlock()
	if s.GetState() != TCPStateUnbound {
		return ErrInvalidSocketState
	}
	if fn == nil {
		return ErrInvalidSocketState
	}
	if err := fn(); err != nil {
		return err
	}
	s.SetState(TCPStateBound)
	return nil
}

// DoListen 在 Bound 下执行 listen，成功则进入 Listening。
func (s *TCPSocket) DoListen(fn func() error) error {
	if s == nil {
		return ErrInvalidSocketState
	}
	s.connectMu.Lock()
	defer s.connectMu.Unlock()
	if s.GetState() != TCPStateBound {
		return ErrInvalidSocketState
	}
	if fn == nil {
		return ErrInvalidSocketState
	}
	if err := fn(); err != nil {
		return err
	}
	s.SetState(TCPStateListening)
	return nil
}

// StoreConnectResult 一次性保存 connect 结果并唤醒等待方。
func (s *TCPSocket) StoreConnectResult(res ConnectResult) {
	if s == nil {
		// 无接收方时仍可能带着已拨通的 Conn，必须关掉
		if res.Conn != nil {
			_ = res.Conn.Close()
		}
		return
	}
	s.connectMu.Lock()
	defer s.connectMu.Unlock()
	if s.connectReady.Load() {
		// 二次结果不能覆盖已保存的 Conn，否则旧/新连接泄漏或被误关
		if res.Conn != nil && res.Conn != s.connectStored.Conn {
			_ = res.Conn.Close()
		}
		return
	}
	if s.GetState() != TCPStateConnecting {
		// drop 已 Closed 时 finish 不会再 Claim，Conn 必须当场关
		if res.Conn != nil {
			_ = res.Conn.Close()
			res.Conn = nil
		}
	}
	s.connectStored = res
	s.connectReady.Store(true)
	if s.ConnectDone != nil {
		select {
		case <-s.ConnectDone:
		default:
			close(s.ConnectDone)
		}
	}
	if s.ConnectResult != nil {
		select {
		case s.ConnectResult <- res:
		default:
		}
	}
}

// TryConnectResult 非阻塞读取已完成的 connect 结果（可重复读，不消费）。
func (s *TCPSocket) TryConnectResult() (ConnectResult, bool) {
	if s == nil || !s.connectReady.Load() {
		return ConnectResult{}, false
	}
	s.connectMu.Lock()
	defer s.connectMu.Unlock()
	return s.connectStored, true
}

// BeginConnect 将 Unbound/Bound 转为 Connecting（并发安全）。
// 两个 start-connect 同时看到 Unbound 会双拨号。
func (s *TCPSocket) BeginConnect() (ok bool, conflict bool) {
	if s == nil {
		return false, false
	}
	s.connectMu.Lock()
	defer s.connectMu.Unlock()
	st := s.GetState()
	if st == TCPStateConnecting {
		return false, true
	}
	if st != TCPStateUnbound && st != TCPStateBound {
		return false, false
	}
	// 先建 ConnectDone 再标 Connecting，避免 subscribe 看到 nil channel 误报就绪。
	if s.ConnectDone == nil {
		s.ConnectDone = make(chan struct{})
	}
	if s.ConnectResult == nil {
		s.ConnectResult = make(chan ConnectResult, 1)
	}
	// 每次 connect 独立可取消 context，drop 才能打断 Dial/Poll
	if s.connectCancel != nil {
		s.connectCancel()
	}
	s.connectCtx, s.connectCancel = context.WithCancel(context.Background())
	s.SetState(TCPStateConnecting)
	return true, false
}

// ClaimConnectSuccess 仅允许一次 Connecting→Connected，并写入 Conn。
// 并发 finish-connect 会重复建 input/output stream。
func (s *TCPSocket) ClaimConnectSuccess(conn *net.TCPConn) bool {
	if s == nil {
		return false
	}
	s.connectMu.Lock()
	defer s.connectMu.Unlock()
	if s.GetState() != TCPStateConnecting {
		return false
	}
	s.Conn = conn
	// 所有权转到 sock.Conn 后清空 stored，避免析构再关一次
	if s.connectStored.Conn == conn {
		s.connectStored.Conn = nil
	}
	s.SetState(TCPStateConnected)
	return true
}

// ClaimConnectFailure 仅允许一次 Connecting→Closed。
func (s *TCPSocket) ClaimConnectFailure() bool {
	if s == nil {
		return false
	}
	s.connectMu.Lock()
	defer s.connectMu.Unlock()
	if s.GetState() != TCPStateConnecting {
		return false
	}
	s.SetState(TCPStateClosed)
	return true
}

// UDPSocket represents a UDP socket resource.
type UDPSocket struct {
	Fd     int
	Conn   *net.UDPConn
	Family IPAddressFamily

	Reader *AsyncUDPReader
	Writer *AsyncUDPWriter

	closeFd func() error
	mu      sync.Mutex // Stream/Receive/Drop 与析构并发改 Reader/Writer/Conn
}

func (s *UDPSocket) HasFd() bool { return s != nil && s.closeFd != nil }

func (s *UDPSocket) DetachFd() {
	if s == nil {
		return
	}
	s.closeFd = nil
	s.Fd = -1
}

// DoBind 在尚未拥有 Conn 时执行 bind/listenUDP。
func (s *UDPSocket) DoBind(fn func() error) error {
	if s == nil {
		return ErrInvalidSocketState
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Conn != nil {
		return ErrInvalidSocketState
	}
	if fn == nil {
		return ErrInvalidSocketState
	}
	return fn()
}

func (s *UDPSocket) GetReader() *AsyncUDPReader {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Reader
}

func (s *UDPSocket) GetWriter() *AsyncUDPWriter {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Writer
}

func (s *UDPSocket) TakeReader() *AsyncUDPReader {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	r := s.Reader
	s.Reader = nil
	s.mu.Unlock()
	return r
}

func (s *UDPSocket) TakeWriter() *AsyncUDPWriter {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	w := s.Writer
	s.Writer = nil
	s.mu.Unlock()
	return w
}

// ReplaceDatagramStreams 在同一把锁下 connect + 停旧流 + 开新流，避免双 ReadFromUDP。
func (s *UDPSocket) ReplaceDatagramStreams(connect func() error) error {
	if s == nil {
		return ErrInvalidSocketState
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Conn == nil {
		return ErrInvalidSocketState
	}
	if connect != nil {
		if err := connect(); err != nil {
			return err
		}
	}
	if s.Reader != nil {
		s.Reader.Close()
		s.Reader.WaitExit()
		s.Reader = nil
	}
	if s.Writer != nil {
		s.Writer.Close()
		s.Writer.WaitExit()
		s.Writer = nil
	}
	if s.Conn != nil {
		_ = s.Conn.SetReadDeadline(time.Time{})
		_ = s.Conn.SetWriteDeadline(time.Time{})
	}
	s.Reader = NewAsyncUDPReader(s.Conn)
	s.Writer = NewAsyncUDPWriter(s.Conn)
	return nil
}

// ResolveAddressStreamState 保存了域名解析操作的状态。
type ResolveAddressStreamState struct {
	// 存储解析出的 IP 地址列表。
	Addresses []net.IP
	// 当前已返回给 Guest 的地址索引。
	Index int
	// 在异步解析过程中可能发生的错误。
	Error error
	// 一个 channel，当后台解析任务完成时，它会被关闭。
	Done chan struct{}
	// DNS goroutine 与 resolve-next-address 并发
	Mu sync.Mutex
	// drop 时取消 LookupIPAddr，避免 goroutine 泄漏
	cancel context.CancelFunc
}

// SetCancel 登记 DNS 查找的取消函数（wasip2 与 manager 不同包，cancel 未导出）。
func (s *ResolveAddressStreamState) SetCancel(c context.CancelFunc) {
	if s == nil {
		return
	}
	s.cancel = c
}

// --- Resource Managers ---

type (
	NetworkManager              = witgo.ResourceManager[*Network]
	TCPSocketManager            = witgo.ResourceManager[*TCPSocket]
	UDPSocketManager            = witgo.ResourceManager[*UDPSocket]
	ResolveAddressStreamManager = witgo.ResourceManager[*ResolveAddressStreamState]
)

func NewNetworkManager() *NetworkManager {
	return witgo.NewResourceManager[*Network](nil)
}

func NewTCPSocketManager() *TCPSocketManager {
	return witgo.NewResourceManager[*TCPSocket](func(resource *TCPSocket) {
		if resource == nil {
			return
		}
		// 与 bind/connect 同一把锁交接 Conn/Listener/fd，避免关 fd 后仍 Connect。
		resource.connectMu.Lock()
		resource.SetState(TCPStateClosed)
		if resource.connectCancel != nil {
			// 打断 in-flight DialContext/Poll，否则 drop 后 dial goroutine 泄漏
			resource.connectCancel()
			resource.connectCancel = nil
		}
		conn := resource.Conn
		resource.Conn = nil
		stored := resource.connectStored.Conn
		resource.connectStored.Conn = nil
		ln := resource.Listener
		resource.Listener = nil
		cfd := resource.closeFd
		resource.closeFd = nil
		resource.Fd = -1
		done := resource.ConnectDone
		resource.connectMu.Unlock()
		if conn != nil {
			_ = conn.Close()
		}
		// finish-connect 尚未 Claim 时连接只在 connectStored 里
		if stored != nil && stored != conn {
			_ = stored.Close()
		}
		if ln != nil {
			_ = ln.Close()
		}
		if cfd != nil {
			_ = cfd()
		}
		if done != nil {
			// connect 被取消且尚未 Store 时，subscribe 会永久堵在 ConnectDone
			select {
			case <-done:
			default:
				close(done)
			}
		}
	})
}

func NewUDPSocketManager() *UDPSocketManager {
	return witgo.NewResourceManager[*UDPSocket](func(resource *UDPSocket) {
		// 统一释放 reader/writer/conn
		if resource == nil {
			return
		}
		resource.mu.Lock()
		r := resource.Reader
		w := resource.Writer
		c := resource.Conn
		resource.Reader = nil
		resource.Writer = nil
		resource.mu.Unlock()

		if r != nil {
			r.Close()
		}
		if w != nil {
			w.Close()
		}
		if c != nil {
			_ = c.Close()
		}

		if r != nil {
			r.waitExit()
		}
		if w != nil {
			w.waitExit()
		}

		resource.mu.Lock()
		cfd := resource.closeFd
		resource.closeFd = nil
		resource.Fd = -1
		resource.mu.Unlock()
		if cfd != nil {
			_ = cfd()
		}
	})
}

func NewResolveAddressStreamManager() *ResolveAddressStreamManager {
	return witgo.NewResourceManager[*ResolveAddressStreamState](func(resource *ResolveAddressStreamState) {
		if resource != nil && resource.cancel != nil {
			resource.cancel()
		}
	})
}
