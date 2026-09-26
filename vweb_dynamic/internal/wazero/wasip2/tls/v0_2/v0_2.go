package v0_2

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	manager_io "github.com/456vv/x/vweb_dynamic/internal/wazero/manager/io"
	manager_tls "github.com/456vv/x/vweb_dynamic/internal/wazero/manager/tls"
	"github.com/456vv/x/vweb_dynamic/internal/wazero/wasip2"
	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"

	"github.com/tetratelabs/wazero"
)

// timeoutErr 实现 net.Error，供 crypto/tls 把 deadline 当作可重试超时。
type timeoutErr struct{}

func (timeoutErr) Error() string   { return "tls stream deadline exceeded" }
func (timeoutErr) Timeout() bool   { return true }
func (timeoutErr) Temporary() bool { return true }

// streamConn 把 WASI 异步 Stream 适配成阻塞 net.Conn。
// HandshakeContext 通过 SetDeadline(过去) 打断 Read/Write；只调用一次 OnSubscribe。
type streamConn struct {
	in     *manager_io.Stream
	out    *manager_io.Stream
	closer io.Closer

	mu        sync.Mutex
	readDead  time.Time
	writeDead time.Time
	closed    bool
	dlCh      chan struct{} // Close/SetDeadline 换新 channel，唤醒所有等待者
}

func (c *streamConn) wakeLocked() {
	if c.dlCh == nil {
		c.dlCh = make(chan struct{})
		return
	}
	select {
	case <-c.dlCh:
	default:
		close(c.dlCh)
	}
	c.dlCh = make(chan struct{})
}

func (c *streamConn) errIfClosedOrDeadline(isWrite bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return net.ErrClosed
	}
	dl := c.readDead
	if isWrite {
		dl = c.writeDead
	}
	if !dl.IsZero() && !dl.After(time.Now()) {
		return timeoutErr{}
	}
	return nil
}

func (c *streamConn) waitReady(s *manager_io.Stream, isWrite bool) error {
	if err := c.errIfClosedOrDeadline(isWrite); err != nil {
		return err
	}
	var p manager_io.IPollable
	if s != nil && s.OnSubscribe != nil {
		p = s.OnSubscribe() // 只能订一次，二次 subscribe 会 Reset 刚就绪的 pollable
	}
	if p != nil && p.IsReady() {
		return c.errIfClosedOrDeadline(isWrite)
	}

	c.mu.Lock()
	dlCh := c.dlCh
	dl := c.readDead
	if isWrite {
		dl = c.writeDead
	}
	c.mu.Unlock()

	var timeout <-chan time.Time
	if !dl.IsZero() {
		d := time.Until(dl)
		if d <= 0 {
			return timeoutErr{}
		}
		t := time.NewTimer(d)
		defer t.Stop()
		timeout = t.C
	}

	var ready <-chan struct{}
	if p != nil {
		ready = p.Channel()
	}

	if ready == nil && timeout == nil {
		// time.After 在热路径泄漏 timer
		t := time.NewTimer(time.Millisecond)
		defer t.Stop()
		select {
		case <-dlCh:
			return c.errIfClosedOrDeadline(isWrite)
		case <-t.C:
			return c.errIfClosedOrDeadline(isWrite)
		}
	}

	select {
	case <-ready:
	case <-dlCh:
	case <-timeout:
		return timeoutErr{}
	}

	return c.errIfClosedOrDeadline(isWrite)
}

func (c *streamConn) Read(b []byte) (int, error) {
	if c.in == nil || c.in.Reader == nil {
		return 0, io.EOF
	}
	for {
		if err := c.errIfClosedOrDeadline(false); err != nil {
			return 0, err
		}
		n, err := c.in.Reader.Read(b)
		if n > 0 || err != nil {
			return n, err
		}
		// (0,nil) 是异步流“暂无数据”，必须阻塞，否则 tls 握手忙等占满 CPU
		if err := c.waitReady(c.in, false); err != nil {
			return 0, err
		}
	}
}

func (c *streamConn) Write(b []byte) (int, error) {
	if c.out == nil || c.out.Writer == nil {
		return 0, io.ErrClosedPipe
	}
	total := 0
	for total < len(b) {
		if err := c.errIfClosedOrDeadline(true); err != nil {
			return total, err
		}
		if c.out.CheckWriter != nil && c.out.CheckWriter.CheckWrite() == 0 {
			if err := c.waitReady(c.out, true); err != nil {
				return total, err
			}
			continue
		}
		n, err := c.out.Writer.Write(b[total:])
		total += n
		if err != nil {
			return total, err
		}
		if n == 0 {
			if err := c.waitReady(c.out, true); err != nil {
				return total, err
			}
		}
	}
	return total, nil
}

func (c *streamConn) Close() error {
	c.mu.Lock()
	c.closed = true
	c.wakeLocked()
	c.mu.Unlock()
	if c.closer != nil {
		return c.closer.Close()
	}
	return nil
}

type dummyNetAddr struct{}

func (dummyNetAddr) Network() string { return "tcp" }
func (dummyNetAddr) String() string  { return "0.0.0.0:0" }

func (c *streamConn) LocalAddr() net.Addr {
	// crypto/tls 部分路径会调 addr.String()，返回 nil 会 panic
	return dummyNetAddr{}
}
func (c *streamConn) RemoteAddr() net.Addr { return dummyNetAddr{} }

func (c *streamConn) SetDeadline(t time.Time) error {
	c.mu.Lock()
	c.readDead = t
	c.writeDead = t
	c.wakeLocked() // HandshakeContext 取消时用过去时间打断阻塞 Read
	c.mu.Unlock()
	return nil
}

func (c *streamConn) SetReadDeadline(t time.Time) error {
	c.mu.Lock()
	c.readDead = t
	c.wakeLocked()
	c.mu.Unlock()
	return nil
}

func (c *streamConn) SetWriteDeadline(t time.Time) error {
	c.mu.Lock()
	c.writeDead = t
	c.wakeLocked()
	c.mu.Unlock()
	return nil
}

type tlsTypes struct {
	tm *manager_tls.TLSManager
	sm *manager_io.StreamManager
	em *manager_io.ErrorManager
}

func NewTypes(tm *manager_tls.TLSManager, sm *manager_io.StreamManager, em *manager_io.ErrorManager) wasip2.Implementation {
	return &tlsTypes{tm: tm, sm: sm, em: em}
}

func (i *tlsTypes) Name() string       { return "wasi:tls/types" }
func (i *tlsTypes) Versions() []string { return []string{"0.2.0-draft"} }

func (i *tlsTypes) Instantiate(_ context.Context, h *wasip2.Host, builder wazero.HostModuleBuilder) error {
	exporter := witgo.NewExporter(builder)

	tm := h.TLSManager()
	sm := h.StreamManager()
	em := h.ErrorManager()

	// --- client-handshake ---
	exporter.Export("[constructor]client-handshake", func(serverName string, inputStream InputStream, outputStream OutputStream) ClientHandshake {
		inStream, inOk := sm.Pop(inputStream)
		outStream, outOk := sm.Pop(outputStream)
		if !inOk || !outOk {
			// 只成功 Pop 其中一个时若不 Close，底层 TCP 流会泄漏且不再出现在 StreamManager 里
			if inOk && inStream != nil && inStream.Closer != nil {
				_ = inStream.Closer.Close()
			}
			if outOk && outStream != nil && outStream.Closer != nil {
				_ = outStream.Closer.Close()
			}
			panic("invalid input or output stream for TLS handshake")
		}

		handshake := &manager_tls.ClientHandshake{
			ServerName: serverName,
			Input:      *inStream,
			Output:     *outStream,
		}
		return tm.ClientHandshakes.Add(handshake)
	})
	exporter.Export("[resource-drop]client-handshake", tm.ClientHandshakes.Remove)

	exporter.Export("[static]client-handshake.finish", func(this ClientHandshake) FutureClientStreams {
		// 根据 WIT 规范，finish 会消费掉 client-handshake 句柄。
		handshake, ok := tm.ClientHandshakes.Pop(this)
		if !ok {
			panic("invalid client-handshake handle")
		}

		ctx, cancel := context.WithCancel(context.Background())
		future := &manager_tls.FutureClientStreams{
			Pollable: manager_io.NewPollable(nil),
			Cancel:   cancel, // drop 必须能打断 Handshake；不能用 host 调用 ctx（返回后常被取消）
		}
		futureHandle := tm.FutureClientStreams.Add(future)
		go func() {
			defer future.Pollable.SetReady()
			var tlsConn *tls.Conn
			defer func() {
				if rec := recover(); rec != nil {
					// Handshake/streamConn panic 时必须唤醒 pollable 并记录错误，避免 future 泄漏
					if tlsConn != nil {
						_ = tlsConn.Close()
					} else {
						_ = handshake.Close()
					}
					future.StoreResult(manager_tls.Result{Err: fmt.Errorf("tls handshake panic: %v", rec)})
				}
			}()
			underlyingConn := &streamConn{
				in:     &handshake.Input,
				out:    &handshake.Output,
				closer: handshake,
				dlCh:   make(chan struct{}),
			}
			tlsConn = tls.Client(underlyingConn, &tls.Config{
				ServerName: handshake.ServerName,
				MinVersion: tls.VersionTLS12,
			})
			if err := tlsConn.HandshakeContext(ctx); err != nil {
				_ = tlsConn.Close()
				future.StoreResult(manager_tls.Result{Err: err})
				return
			}
			future.StoreResult(manager_tls.Result{TlsConn: tlsConn})
		}()

		return futureHandle
	})

	// --- client-connection ---
	exporter.Export("[resource-drop]client-connection", tm.ClientConnections.Remove)

	exporter.Export("[method]client-connection.close-output", func(this ClientConnection) {
		conn, ok := tm.ClientConnections.Get(this)
		// 无效句柄或 Conn 尚未填充时解引用会 panic
		if !ok || conn == nil || conn.Conn == nil {
			return
		}
		_ = conn.Conn.CloseWrite()
	})

	// --- future-client-streams ---
	exporter.Export("[resource-drop]future-client-streams", tm.FutureClientStreams.Remove)
	exporter.Export("[method]future-client-streams.subscribe", func(this FutureClientStreams) Pollable {
		future, ok := tm.FutureClientStreams.Get(this)
		if !ok || future == nil || future.Pollable == nil {
			return h.PollManager().Add(manager_io.NewReadyPollable())
		}
		// 直接 Add 内部 Pollable 时 drop subscribe 句柄会 SetReady，get 误判握手完成。
		return h.PollManager().Add(manager_io.NewLevelPollable(future.Pollable.IsReady, future.Pollable))
	})

	exporter.Export("[method]future-client-streams.get", func(ctx context.Context, this FutureClientStreams) witgo.Option[witgo.Result[witgo.Result[witgo.Tuple3[ClientConnection, InputStream, OutputStream], WasiError], witgo.Unit]] {
		none := witgo.None[witgo.Result[witgo.Result[witgo.Tuple3[ClientConnection, InputStream, OutputStream], WasiError], witgo.Unit]]()
		_ = ctx // WASI get 非阻塞；不在 ctx 上等待，以免丢掉已就绪的握手结果

		// 原先先 Pop 再等 Channel/ctx。ctx 取消时 future 已从管理器摘掉，
		// 握手结果丢失且 tls.Conn 无人关闭。Get 保留句柄，允许取消后再次 get。
		future, ok := tm.FutureClientStreams.Get(this)
		if !ok || future == nil {
			return none
		}

		if future.Pollable == nil || !future.Pollable.IsReady() {
			return none
		}

		// 规范为至多消费一次；重复 get 返回 Some(Err(unit))
		if !future.Consumed.CompareAndSwap(false, true) {
			return witgo.Some(witgo.Err[witgo.Result[witgo.Tuple3[ClientConnection, InputStream, OutputStream], WasiError], witgo.Unit](witgo.Unit{}))
		}

		res := future.LoadResult()
		if res.Err != nil {
			errHandle := em.Add(res.Err)
			return witgo.Some(witgo.Ok[witgo.Result[witgo.Tuple3[ClientConnection, InputStream, OutputStream], WasiError], witgo.Unit](
				witgo.Err[witgo.Tuple3[ClientConnection, InputStream, OutputStream], WasiError](errHandle),
			))
		}

		tlsConn := future.TakeTlsConn()
		if tlsConn == nil {
			errHandle := em.Add(io.ErrUnexpectedEOF)
			return witgo.Some(witgo.Ok[witgo.Result[witgo.Tuple3[ClientConnection, InputStream, OutputStream], WasiError], witgo.Unit](
				witgo.Err[witgo.Tuple3[ClientConnection, InputStream, OutputStream], WasiError](errHandle),
			))
		}

		// 先把 tls.Conn 放进 ClientConnections。若先 Add 流再 Add conn，
		// Host.Close 插在中间会得到 DontClose 的流且 tls.Conn 不在任何 manager 里。
		conn := &manager_tls.ClientConnection{Conn: tlsConn}
		connHandle := tm.ClientConnections.Add(conn)

		inStreamEncrypted := manager_io.NewAsyncStreamForReader(tlsConn, manager_io.DontCloseReader())
		outStreamEncrypted := manager_io.NewAsyncStreamForWriter(tlsConn, manager_io.DontCloseWriter())
		inStreamHandle := sm.Add(inStreamEncrypted)
		outStreamHandle := sm.Add(outStreamEncrypted)

		tuple := witgo.Tuple3[ClientConnection, InputStream, OutputStream]{
			F0: connHandle,
			F1: inStreamHandle,
			F2: outStreamHandle,
		}

		return witgo.Some(witgo.Ok[witgo.Result[witgo.Tuple3[ClientConnection, InputStream, OutputStream], WasiError], witgo.Unit](
			witgo.Ok[witgo.Tuple3[ClientConnection, InputStream, OutputStream], WasiError](tuple),
		))
	})

	return nil
}
