package tls

import (
	"context"
	"crypto/tls"
	"sync"
	"sync/atomic"

	manager_io "github.com/456vv/x/vweb_dynamic/internal/wazero/manager/io"
	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"
)

// ClientHandshake 代表一个客户端 TLS 握手操作。
type ClientHandshake struct {
	ServerName string
	Input      manager_io.Stream // The underlying input stream (e.g., from a TCP socket)
	Output     manager_io.Stream // The underlying output stream (e.g., from a TCP socket)
}

func (c *ClientHandshake) Close() error {
	if c == nil {
		return nil
	}
	return manager_io.NewMultiCloser(c.Input.Closer, c.Output.Closer).Close()
}

// ClientConnection 代表一个已建立的 TLS 连接。
type ClientConnection struct {
	Conn *tls.Conn
}

func (c *ClientConnection) Close() error {
	if c != nil && c.Conn != nil {
		return c.Conn.Close()
	}
	return nil
}

// FutureClientStreams 代表一个尚未完成的 TLS 握手，最终会产生加密流。
type FutureClientStreams struct {
	Pollable *manager_io.ChannelPollable
	Result   Result
	Consumed atomic.Bool
	resultMu sync.Mutex // 握手 goroutine 写 Result 与 get/drop 并发
	// Cancel 在 drop future 时取消 HandshakeContext，避免握手 goroutine 泄漏。
	Cancel context.CancelFunc
}

func (c *FutureClientStreams) Close() error {
	if c == nil {
		return nil
	}
	conn := c.TakeTlsConn()
	if conn != nil {
		return conn.Close()
	}
	return nil
}

// StoreResult 在锁内写入握手结果。
func (c *FutureClientStreams) StoreResult(res Result) {
	if c == nil {
		return
	}
	c.resultMu.Lock()
	c.Result = res
	c.resultMu.Unlock()
}

// LoadResult 在锁内拷贝握手结果。
func (c *FutureClientStreams) LoadResult() Result {
	if c == nil {
		return Result{}
	}
	c.resultMu.Lock()
	defer c.resultMu.Unlock()
	return c.Result
}

// TakeTlsConn 转移 tls.Conn 所有权，避免 drop future 与 client-connection 双关。
func (c *FutureClientStreams) TakeTlsConn() *tls.Conn {
	if c == nil {
		return nil
	}
	c.resultMu.Lock()
	defer c.resultMu.Unlock()
	conn := c.Result.TlsConn
	c.Result.TlsConn = nil
	return conn
}

// Result 是一个内部类型，用于在 goroutine 之间传递 TLS 握手的结果。
type Result struct {
	TlsConn *tls.Conn
	Err     error // 或一个 Go 的 error
}

// TLSManager 是所有 TLS 相关资源的总管理器。
type TLSManager struct {
	ClientHandshakes    *witgo.ResourceManager[*ClientHandshake]
	ClientConnections   *witgo.ResourceManager[*ClientConnection]
	FutureClientStreams *witgo.ResourceManager[*FutureClientStreams]
}

func NewTLSManager() *TLSManager {
	return &TLSManager{
		ClientHandshakes: witgo.NewResourceManager[*ClientHandshake](func(resource *ClientHandshake) {
			if resource != nil {
				resource.Close()
			}
		}),
		ClientConnections: witgo.NewResourceManager[*ClientConnection](func(resource *ClientConnection) {
			if resource != nil {
				resource.Close()
			}
		}),
		FutureClientStreams: witgo.NewResourceManager[*FutureClientStreams](func(resource *FutureClientStreams) {
			if resource == nil {
				return
			}
			// drop 时取消握手，否则 Handshake 可能一直阻塞在底层 Read。
			if resource.Cancel != nil {
				resource.Cancel()
				resource.Cancel = nil // 与 HTTP future 对齐，避免二次析构重复观察 Cancel
			}
			if resource.Pollable != nil {
				resource.Pollable.Block()
			}
			// 已 Consumed 时 TlsConn 已交给 ClientConnection，不能再关
			if !resource.Consumed.CompareAndSwap(false, true) {
				return
			}
			_ = resource.Close()
		}),
	}
}
