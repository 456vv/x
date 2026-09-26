package wasip2

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/456vv/x/vweb_dynamic/internal/wazero/manager/filesystem"
	"github.com/456vv/x/vweb_dynamic/internal/wazero/manager/http"
	"github.com/456vv/x/vweb_dynamic/internal/wazero/manager/io"
	"github.com/456vv/x/vweb_dynamic/internal/wazero/manager/sockets"
	"github.com/456vv/x/vweb_dynamic/internal/wazero/manager/tls"

	"github.com/tetratelabs/wazero"
)

// Implementation 是所有 WASI 模块必须实现的接口。
type Implementation interface {
	// Name 返回模块的名称，例如 "wasi:io"。
	Name() string
	// Versions 返回此实现兼容的 WIT 版本列表，例如 ["0.2.0", "0.2.1"]。
	Versions() []string
	// Instantiate 将模块的函数导出到 wazero 运行时。
	Instantiate(context.Context, *Host, wazero.HostModuleBuilder) error
}

// Host 是所有 WASI 实现的容器。
type Host struct {
	streamManager *io.StreamManager
	errorManager  *io.ErrorManager
	pollManager   *io.PollManager
	httpManager   *http.HTTPManager
	tlsManager    *tls.TLSManager

	// filesystem 管理器
	filesystemManager           *filesystem.Manager
	directoryEntryStreamManager *filesystem.DirectoryEntryStreamManager

	//  Sockets 管理器
	networkManager              *sockets.NetworkManager
	tcpSocketManager            *sockets.TCPSocketManager
	udpSocketManager            *sockets.UDPSocketManager
	resolveAddressStreamManager *sockets.ResolveAddressStreamManager
	// 未来可以在这里添加 httpManager 等其他状态管理器

	implementations []Implementation
	implMu          sync.Mutex // AddImplementation 与 Instantiate 并发会踩切片
}

// ModuleOption 是用于配置 Host 的选项函数。
type ModuleOption func(*Host)

// NewHost 创建一个新的 Host 实例，并应用所有提供的模块选项。
func NewHost(opts ...ModuleOption) *Host {
	streamManager, pollManager, errorManager := io.NewManager()
	h := &Host{
		streamManager: streamManager,
		errorManager:  errorManager,
		pollManager:   pollManager,
		httpManager:   http.NewHTTPManager(streamManager, pollManager),
		tlsManager:    tls.NewTLSManager(),

		filesystemManager:           filesystem.NewManager(),
		directoryEntryStreamManager: filesystem.NewDirectoryEntryStreamManager(),

		networkManager:              sockets.NewNetworkManager(),
		tcpSocketManager:            sockets.NewTCPSocketManager(),
		udpSocketManager:            sockets.NewUDPSocketManager(),
		resolveAddressStreamManager: sockets.NewResolveAddressStreamManager(),
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

func (h *Host) AddImplementation(impl Implementation) {
	if h == nil || impl == nil {
		return
	}
	h.implMu.Lock()
	h.implementations = append(h.implementations, impl)
	h.implMu.Unlock()
}

// Instantiate 将所有已配置的模块实例化到 wazero 运行时。
func (h *Host) Instantiate(ctx context.Context, r wazero.Runtime) error {
	if r == nil {
		return fmt.Errorf("wazero runtime is nil")
	}
	h.implMu.Lock()
	impls := append([]Implementation(nil), h.implementations...)
	h.implMu.Unlock()
	for _, impl := range impls {
		for _, version := range impl.Versions() {
			moduleName := impl.Name() + "@" + version
			builder := r.NewHostModuleBuilder(moduleName)
			if err := impl.Instantiate(ctx, h, builder); err != nil {
				// 带上模块名，便于定位实例化失败原因
				return fmt.Errorf("instantiate %s: %w", moduleName, err)
			}

			if _, err := builder.Instantiate(ctx); err != nil {
				return fmt.Errorf("builder.Instantiate %s: %w", moduleName, err)
			}
		}
	}
	return nil
}

func (h *Host) StreamManager() *io.StreamManager {
	return h.streamManager
}

func (h *Host) ErrorManager() *io.ErrorManager {
	return h.errorManager
}

func (h *Host) PollManager() *io.PollManager {
	return h.pollManager
}

func (h *Host) HTTPManager() *http.HTTPManager {
	return h.httpManager
}

// FilesystemManager 返回文件系统管理器。
func (h *Host) FilesystemManager() *filesystem.Manager {
	return h.filesystemManager
}

// AddPreopen 登记一个预打开目录，供 wasi:filesystem/preopens.get-directories 返回。
// guestPath 是 guest 看到的路径（常见 "/" 或 "."），不是宿主绝对路径。
// 注意：调用方原先 fsm.Add(&Descriptor{File, Path}) 不会置 IsPreopen，
// ForPreopen 修正后 get-directories 会漏掉全部预打开项。
func (h *Host) AddPreopen(file *os.File, guestPath string) uint32 {
	if h == nil || h.filesystemManager == nil || file == nil {
		return 0
	}
	return h.filesystemManager.Add(filesystem.NewPreopenDescriptor(file, guestPath))
}

// DirectoryEntryStreamManager 返回目录条目流管理器。
func (h *Host) DirectoryEntryStreamManager() *filesystem.DirectoryEntryStreamManager {
	return h.directoryEntryStreamManager
}

func (h *Host) NetworkManager() *sockets.NetworkManager {
	return h.networkManager
}

func (h *Host) TCPSocketManager() *sockets.TCPSocketManager {
	return h.tcpSocketManager
}

func (h *Host) UDPSocketManager() *sockets.UDPSocketManager {
	return h.udpSocketManager
}

func (h *Host) ResolveAddressStreamManager() *sockets.ResolveAddressStreamManager {
	return h.resolveAddressStreamManager
}

func (h *Host) TLSManager() *tls.TLSManager {
	return h.tlsManager
}

func (h *Host) Close() {
	if h == nil {
		return
	}
	// 必须先关 HTTP Transport 空闲连接，再 Clear 资源；
	// 各 manager 判空后再 Clear，避免测试/部分构造失败时对 nil 解引用 panic。
	if h.httpManager != nil {
		h.httpManager.CloseIdleConnections()
		h.httpManager.ClearImmutableFields()
		// 先取消并等待 in-flight Client.Do，再关 pipe body。
		// finish() 已 Pop 的 request 不在 OutgoingRequests 里，Bodies.Clear 先关 Writer
		// 会把未发完的请求截断；Futures.Clear 的 Cancel+Block 才能按取消路径结束 Do。
		if h.httpManager.Futures != nil {
			h.httpManager.Futures.Clear()
		}
		if h.httpManager.FutureTrailers != nil {
			h.httpManager.FutureTrailers.Clear()
		}
		if h.httpManager.IncomingBodies != nil {
			h.httpManager.IncomingBodies.Clear()
		}
		if h.httpManager.Bodies != nil {
			h.httpManager.Bodies.Clear()
		}
		if h.httpManager.Responses != nil {
			h.httpManager.Responses.Clear()
		}
		if h.httpManager.IncomingRequests != nil {
			h.httpManager.IncomingRequests.Clear()
		}
		if h.httpManager.OutgoingRequests != nil {
			h.httpManager.OutgoingRequests.Clear()
		}
		if h.httpManager.OutgoingResponses != nil {
			h.httpManager.OutgoingResponses.Clear()
		}
		if h.httpManager.ResponseOutparams != nil {
			h.httpManager.ResponseOutparams.Clear()
		}
		if h.httpManager.Options != nil {
			h.httpManager.Options.Clear()
		}
		if h.httpManager.Fields != nil {
			h.httpManager.Fields.Clear()
		}
	}
	if h.tlsManager != nil {
		// handshake.finish 已 Pop ClientHandshake，真正需要先取消的是 future。
		// 先 Clear handshakes 对 in-flight 握手无影响，但已 Consumed 的连接应先于 future 关闭，
		// 避免 future 与 connection 双关。未 Consumed 的 TlsConn 仍在 future 里。
		if h.tlsManager.FutureClientStreams != nil {
			h.tlsManager.FutureClientStreams.Clear()
		}
		if h.tlsManager.ClientConnections != nil {
			h.tlsManager.ClientConnections.Clear()
		}
		if h.tlsManager.ClientHandshakes != nil {
			h.tlsManager.ClientHandshakes.Clear()
		}
	}

	// 先停 stream 后台 Read/Write，再关 TCP/UDP，避免与 Conn.Close 并发。
	if h.streamManager != nil {
		h.streamManager.Clear()
	}

	if h.tcpSocketManager != nil {
		h.tcpSocketManager.Clear()
	}
	if h.udpSocketManager != nil {
		h.udpSocketManager.Clear()
	}
	if h.resolveAddressStreamManager != nil {
		h.resolveAddressStreamManager.Clear()
	}
	if h.networkManager != nil {
		h.networkManager.Clear()
	}
	if h.filesystemManager != nil {
		h.filesystemManager.Clear()
	}
	if h.directoryEntryStreamManager != nil {
		h.directoryEntryStreamManager.Clear()
	}

	if h.pollManager != nil {
		h.pollManager.Clear()
	}
	if h.errorManager != nil {
		h.errorManager.Clear()
	}
}
