package io

import (
	"bytes"
	"context"
	"io"
	"sync"
	"sync/atomic"
	"time"
)

// 默认的读取缓冲区大小，用于后台读取操作。
const defaultBufferSize = 8192

// 后台读 goroutine 每轮若新建 8KiB 切片，在高频流上制造 GC 压力；Pool 复用读缓冲。
var readBufPool = sync.Pool{
	New: func() any {
		b := make([]byte, defaultBufferSize)
		return &b
	},
}

// --- 优化后的异步读取封装器 ---

// AsyncReadWrapperOption 是用于配置 AsyncReadWrapper 的函数类型。
type AsyncReadWrapperOption func(*AsyncReadWrapper)

// DontCloseReader 是一个选项，用于阻止在关闭 wrapper 时关闭底层的 reader。
func DontCloseReader() AsyncReadWrapperOption {
	return func(arw *AsyncReadWrapper) {
		arw.closeUnderlying = false
	}
}

// WithReadMaxBufferSize 限制后台读缓冲，避免慢消费者导致无界内存。
func WithReadMaxBufferSize(size int) AsyncReadWrapperOption {
	return func(arw *AsyncReadWrapper) {
		if size > 0 {
			arw.maxBufferSize = size
		}
	}
}

// DontCloseReader 时必须能打断底层阻塞 Read（TCP/TLS），否则 Close 死等 exited。
type deadlineReadSetter interface {
	SetReadDeadline(time.Time) error
}

type deadlineWriteSetter interface {
	SetWriteDeadline(time.Time) error
}

type AsyncReadWrapper struct {
	reader          io.Reader
	buffer          *bytes.Buffer
	mutex           sync.Mutex
	cond            *sync.Cond // 缓冲满时让 run 等待
	ready           *ChannelPollable
	done            chan struct{}
	exited          chan struct{} // Close 等待 run 退出，避免关 reader 后仍 Read
	err             error
	once            sync.Once
	closeUnderlying bool
	maxBufferSize   int
}

// NewAsyncReadWrapper 创建并启动一个新的异步读取封装器。
func NewAsyncReadWrapper(r io.Reader, opts ...AsyncReadWrapperOption) *AsyncReadWrapper {
	if r == nil {
		// nil Reader 会让后台 goroutine 对 Read 解引用 panic
		r = bytes.NewReader(nil)
	}
	wrapper := &AsyncReadWrapper{
		reader:          r,
		buffer:          &bytes.Buffer{},
		ready:           NewPollable(nil),
		done:            make(chan struct{}),
		exited:          make(chan struct{}),
		closeUnderlying: true,    // 默认在 Close 时关闭底层 reader。
		maxBufferSize:   1 << 20, // 默认 1MiB 上限，防止 OOM
	}
	for _, opt := range opts {
		opt(wrapper)
	}
	wrapper.cond = sync.NewCond(&wrapper.mutex)
	go wrapper.run()
	return wrapper
}

func (arw *AsyncReadWrapper) closed() bool {
	select {
	case <-arw.done:
		return true
	default:
		return false
	}
}

func (arw *AsyncReadWrapper) readableLocked() bool {
	return arw.buffer.Len() > 0 || arw.err != nil || arw.closed()
}

func (arw *AsyncReadWrapper) readable() bool {
	arw.mutex.Lock()
	defer arw.mutex.Unlock()
	return arw.readableLocked()
}

func (arw *AsyncReadWrapper) run() {
	defer func() {
		close(arw.exited)
		arw.mutex.Lock()
		arw.ready.SetReady()
		if arw.cond != nil {
			arw.cond.Broadcast()
		}
		arw.mutex.Unlock()
	}()

	bufp := readBufPool.Get().(*[]byte)
	readBuf := *bufp
	// goroutine 退出时归还缓冲，避免 Pool 条目泄漏
	defer readBufPool.Put(bufp)

	for {
		arw.mutex.Lock()
		for arw.buffer.Len() >= arw.maxBufferSize && arw.err == nil && !arw.closed() {
			arw.cond.Wait()
		}
		if arw.closed() {
			arw.mutex.Unlock()
			return
		}
		arw.mutex.Unlock()

		n, readErr := arw.reader.Read(readBuf)

		// io.Reader 允许 len(p)>0 时返回 (0, nil) 表示“无进展”；原先会忙等占满 CPU。
		if n == 0 && readErr == nil {
			if arw.closed() {
				return
			}

			t := time.NewTimer(time.Millisecond)
			select {
			case <-arw.done:
				if !t.Stop() {
					select {
					case <-t.C:
					default:
					}
				}
				return
			case <-t.C:
			}
			continue
		}

		arw.mutex.Lock()
		// Close 之后丢弃迟到的 Read，避免向已结束的流继续塞数据。
		if arw.closed() {
			if readErr != nil && arw.err == nil {
				arw.err = readErr
			}
			arw.mutex.Unlock()
			return
		}
		wasEmpty := arw.buffer.Len() == 0
		if n > 0 {
			arw.buffer.Write(readBuf[:n])
		}
		if readErr != nil && arw.err == nil {
			arw.err = readErr
		}
		if (wasEmpty && arw.buffer.Len() > 0) || readErr != nil {
			arw.ready.SetReady()
		}
		if arw.cond != nil {
			arw.cond.Broadcast()
		}
		arw.mutex.Unlock()

		// 如果遇到任何错误，就停止读取。
		if readErr != nil {
			return
		}
	}
}

// Read 从内部缓冲区非阻塞地读取数据。
// 如果缓冲区为空，它会返回 (0, nil)。
// 调用者应该使用 subscribe() 来等待数据变为可用。
func (arw *AsyncReadWrapper) Read(p []byte) (n int, err error) {
	arw.mutex.Lock()
	defer arw.mutex.Unlock()

	// 1. 检查缓冲区是否有数据。
	if arw.buffer.Len() > 0 {
		n, _ = arw.buffer.Read(p)
		if arw.cond != nil {
			arw.cond.Signal()
		}
		if arw.buffer.Len() == 0 && arw.err == nil {
			// 排空后 Reset 的是唤醒通道，IsReady 改由 LevelPollable 看缓冲，避免 poll 忙等。
			arw.ready.Reset()
		}

		return n, nil
	}

	// 2. 缓冲区为空，检查是否有已记录的错误（如 EOF）。
	if arw.err != nil {
		return 0, arw.err
	}

	// Close 后若未设 err，guest 会永远看到“暂无数据”。
	if arw.closed() {
		arw.err = io.EOF
		return 0, io.EOF
	}

	// 3. 缓冲区为空且无错误，表示需要等待后台 goroutine 读取更多数据。
	return 0, nil
}

// subscribe 返回一个 pollable 对象，当有数据可读或发生错误时，该对象会变为就绪状态。
func (arw *AsyncReadWrapper) subscribe() IPollable {
	arw.mutex.Lock()
	if arw.readableLocked() {
		arw.ready.SetReady()
	} else {
		arw.ready.Reset()
		// 持同一 mutex 复检，避免 Reset 与 run.SetReady 交错丢失就绪。
		if arw.readableLocked() {
			arw.ready.SetReady()
		}
	}
	arw.mutex.Unlock()
	// 每次 subscribe 返回电平触发包装，IsReady 看缓冲而不是通道是否关过。
	return NewLevelPollable(arw.readable, arw.ready)
}

// Close 停止后台 goroutine 并根据配置决定是否关闭底层资源。
func (arw *AsyncReadWrapper) Close() error {
	var closeErr error
	arw.once.Do(func() {
		arw.mutex.Lock()
		// 必须先标结束并 close(done)，run 在 Read 返回后看到 closed 不再入队。
		// 原先先关 reader 再 close(done)，迟到的成功 Read 会往已关闭流塞数据。
		if arw.err == nil {
			arw.err = io.EOF
		}
		arw.buffer.Reset() // drop 时丢弃未读缓冲，尽快释放内存
		select {
		case <-arw.done:
		default:
			close(arw.done)
		}
		arw.ready.SetReady()
		if arw.cond != nil {
			arw.cond.Broadcast()
		}
		arw.mutex.Unlock()

		var clearer deadlineReadSetter
		if arw.closeUnderlying {
			if c, ok := arw.reader.(io.Closer); ok {
				closeErr = c.Close()
			}
		} else if sd, ok := arw.reader.(deadlineReadSetter); ok {
			// DontCloseReader 时 Conn.Read 不会因 close(done) 返回；用过去 deadline 唤醒 run
			sd.SetReadDeadline(time.Now())
			clearer = sd
		}
		// 关闭底层 reader 后等待 run 退出，避免 use-after-close
		<-arw.exited
		if clearer != nil {
			// TCP 仍由 socket 持有，必须清 deadline，否则后续读立刻超时
			clearer.SetReadDeadline(time.Time{})
		}
	})
	return closeErr
}

// NewAsyncStreamForReader 是一个便捷的辅助函数，
// 将一个阻塞的 io.Reader 转换为完全支持异步 subscribe 的 *Stream。
func NewAsyncStreamForReader(r io.Reader, opts ...AsyncReadWrapperOption) *Stream {
	wrapper := NewAsyncReadWrapper(r, opts...)
	return &Stream{
		Reader: wrapper,
		Closer: wrapper,
		OnSubscribe: func() IPollable {
			return wrapper.subscribe()
		},
	}
}

// --- Asynchronous Writer ---

// AsyncWriteWrapperOption 是用于配置 AsyncWriteWrapper 的函数类型。
type AsyncWriteWrapperOption func(*AsyncWriteWrapper)

// DontCloseWriter 是一个选项，用于阻止在关闭 wrapper 时关闭底层的 writer。
func DontCloseWriter() AsyncWriteWrapperOption {
	return func(aww *AsyncWriteWrapper) {
		aww.closeUnderlying = false
	}
}

// 记录写入数量
func WriterWritten(bytesWritten *atomic.Uint64) AsyncWriteWrapperOption {
	return func(aww *AsyncWriteWrapper) {
		aww.bytesWritten = bytesWritten
	}
}

func WithMaxBufferSize(size int) AsyncWriteWrapperOption {
	return func(aww *AsyncWriteWrapper) {
		if size > 0 {
			aww.maxBufferSize = size
		}
	}
}

// closedWriter 避免 NewAsyncWriteWrapper(nil) 时后台 goroutine 对 Writer 解引用 panic。
type closedWriter struct{}

func (closedWriter) Write([]byte) (int, error) { return 0, io.ErrClosedPipe }

// AsyncWriteWrapper 将一个阻塞的 io.Writer 封装成一个非阻塞的 writer，
// 带有内部缓冲区，并通过 IPollable 接口提供空间可用性通知。
type AsyncWriteWrapper struct {
	writer          io.Writer
	buffer          *bytes.Buffer
	mutex           sync.Mutex
	cond            *sync.Cond
	ready           *ChannelPollable
	done            chan struct{}
	exited          chan struct{} // Close 等待后台写 goroutine 退出
	err             error
	maxBufferSize   int
	once            sync.Once
	closeUnderlying bool
	bytesWritten    *atomic.Uint64
	closed          atomic.Bool // 避免 Close 后 Write 继续写入
	scratch         []byte      // 复用写出临时缓冲，避免每次 run 都 make
	inflight        int         // 已离开 buffer、尚未写完的字节数
}

// NewAsyncWriteWrapper 创建并启动一个新的异步写入封装器。
func NewAsyncWriteWrapper(w io.Writer, opts ...AsyncWriteWrapperOption) *AsyncWriteWrapper {
	if w == nil {
		// nil Writer 会让 run() 中 Write 解引用 panic
		w = closedWriter{}
	}
	wrapper := &AsyncWriteWrapper{
		writer:          w,
		buffer:          &bytes.Buffer{},
		ready:           NewPollable(nil),
		done:            make(chan struct{}),
		exited:          make(chan struct{}),
		maxBufferSize:   defaultBufferSize,
		closeUnderlying: true,
	}
	for _, opt := range opts {
		opt(wrapper)
	}
	wrapper.cond = sync.NewCond(&wrapper.mutex)
	wrapper.ready.SetReady() // 一开始缓冲区是空的，所以是就绪状态
	go wrapper.run()
	return wrapper
}

func (aww *AsyncWriteWrapper) pendingLocked() int {
	return aww.buffer.Len() + aww.inflight
}

func (aww *AsyncWriteWrapper) writableLocked() bool {
	// inflight 也占用背压额度，否则后台正在 Write 时 buffer 为空，check-write 会虚报空间。
	return aww.err != nil || aww.closed.Load() || aww.pendingLocked() < aww.maxBufferSize
}

func (aww *AsyncWriteWrapper) writable() bool {
	aww.mutex.Lock()
	defer aww.mutex.Unlock()
	return aww.writableLocked()
}

func (aww *AsyncWriteWrapper) run() {
	defer close(aww.exited)
	for {
		aww.mutex.Lock()
		for aww.buffer.Len() == 0 {
			select {
			case <-aww.done:
				aww.mutex.Unlock()
				return
			default:
				aww.cond.Wait()
				select {
				case <-aww.done:
					if aww.buffer.Len() == 0 {
						aww.mutex.Unlock()
						return
					}
				default:
				}
			}
		}

		nbuf := aww.buffer.Len()
		if cap(aww.scratch) < nbuf {
			aww.scratch = make([]byte, nbuf)
		}
		tempBuf := aww.scratch[:nbuf]
		aww.buffer.Read(tempBuf)
		// 字节已经不在 buffer 里，但仍未交给底层 writer。
		// Flush 只看 buffer.Len()==0 会提前成功，HTTP body 被截断。
		aww.inflight = nbuf
		aww.mutex.Unlock()

		n, err := aww.writer.Write(tempBuf)

		if n == 0 && err == nil {
			aww.mutex.Lock()
			aww.inflight = 0
			// 底层 Write 允许 (0,nil) 表示无进展；必须把已从 buffer 取出的字节完整回灌，
			// 且放在当前 buffer 之前，避免乱序。原先若只 Reset 会丢数据。
			rest := append([]byte(nil), aww.buffer.Bytes()...)
			aww.buffer.Reset()
			aww.buffer.Write(tempBuf)
			aww.buffer.Write(rest)
			stop := aww.closed.Load() || aww.err != nil
			select {
			case <-aww.done:
				stop = true
			default:
			}
			if stop && aww.err == nil {
				aww.err = io.ErrClosedPipe
			}
			aww.ready.SetReady()
			aww.cond.Broadcast()
			aww.mutex.Unlock()
			if stop {
				return
			}
			t := time.NewTimer(time.Millisecond)
			select {
			case <-aww.done:
				if !t.Stop() {
					select {
					case <-t.C:
					default:
					}
				}
				// 退避期间 Close 已关闭 done。直接 return 会留下缓冲且不再 Broadcast，BlockingFlush 永久 Wait。
				aww.mutex.Lock()
				if aww.err == nil && aww.pendingLocked() > 0 {
					aww.err = io.ErrClosedPipe
				}
				aww.ready.SetReady()
				aww.cond.Broadcast()
				aww.mutex.Unlock()
				return
			case <-t.C:
			}
			continue
		}

		aww.mutex.Lock()
		aww.inflight = 0
		if n > 0 && aww.bytesWritten != nil {
			aww.bytesWritten.Add(uint64(n))
		}
		if err != nil {
			aww.err = err
		}
		if n < len(tempBuf) && n >= 0 {
			remain := append([]byte(nil), tempBuf[n:]...)
			rest := append([]byte(nil), aww.buffer.Bytes()...)
			aww.buffer.Reset()
			aww.buffer.Write(remain)
			aww.buffer.Write(rest)
		}
		if aww.writableLocked() {
			aww.ready.SetReady()
		}
		aww.cond.Broadcast()
		aww.mutex.Unlock()

		if err != nil {
			return
		}
	}
}

func (aww *AsyncWriteWrapper) Write(p []byte) (n int, err error) {
	aww.mutex.Lock()
	defer aww.mutex.Unlock()

	if aww.closed.Load() {
		return 0, io.ErrClosedPipe
	}
	if aww.err != nil {
		return 0, aww.err
	}
	available := aww.maxBufferSize - aww.pendingLocked()
	if available <= 0 {
		return 0, nil
	}
	if len(p) > available {
		n, _ = aww.buffer.Write(p[:available])
	} else {
		n, _ = aww.buffer.Write(p)
	}
	if n > 0 {
		aww.cond.Signal()
	}
	if aww.pendingLocked() >= aww.maxBufferSize {
		aww.ready.Reset()
	}
	return n, nil
}

func (aww *AsyncWriteWrapper) BlockingFlush() error {
	aww.mutex.Lock()
	defer aww.mutex.Unlock()

	// inflight>0 表示 tempBuf 里还有未完成的写，buffer 暂时为空不能算刷新完成。
	for aww.pendingLocked() > 0 && aww.err == nil {
		select {
		case <-aww.done:
			if aww.err == nil && aww.pendingLocked() > 0 {
				aww.err = io.ErrClosedPipe
			}
			return aww.err
		default:
		}
		aww.cond.Signal()
		aww.cond.Wait()
	}
	return aww.err
}

// Flush 触发一次非阻塞的刷新。
func (aww *AsyncWriteWrapper) Flush() error {
	aww.mutex.Lock()
	defer aww.mutex.Unlock()
	if aww.err != nil {
		return aww.err
	}
	if aww.buffer.Len() > 0 {
		aww.cond.Signal()
	}
	return nil
}

func (aww *AsyncWriteWrapper) CheckWrite() uint64 {
	aww.mutex.Lock()
	defer aww.mutex.Unlock()
	if aww.err != nil || aww.closed.Load() {
		return 0
	}
	avail := aww.maxBufferSize - aww.pendingLocked()
	if avail < 0 {
		return 0
	}
	return uint64(avail)
}

// StreamErr 把“关闭/失败”与“缓冲已满”区分开。
// check-write 在关闭后若只返回 0，而 subscribe 因 closed 仍就绪，guest 会空转。
func (aww *AsyncWriteWrapper) StreamErr() error {
	aww.mutex.Lock()
	defer aww.mutex.Unlock()
	if aww.err != nil {
		return aww.err
	}
	if aww.closed.Load() {
		return io.ErrClosedPipe
	}
	return nil
}

func (aww *AsyncWriteWrapper) subscribe() IPollable {
	aww.mutex.Lock()
	if aww.writableLocked() {
		aww.ready.SetReady()
	} else {
		// 缓冲满时若不 Reset，初始 SetReady 会让 guest 误以为可写
		aww.ready.Reset()
		// Reset 与 run 排空并发时可能丢就绪；持锁复检。
		if aww.writableLocked() {
			aww.ready.SetReady()
		}
	}
	aww.mutex.Unlock()
	return NewLevelPollable(aww.writable, aww.ready)
}

func (aww *AsyncWriteWrapper) Close() error {
	// 必须先在锁内标记 closed，禁止 Flush 与 closeInternal 之间的并发 Write 再入队
	aww.mutex.Lock()
	aww.closed.Store(true)
	aww.mutex.Unlock()

	var flushErr error
	if !aww.closeUnderlying {
		// HTTP pipe（DontCloseWriter）必须先排空，读端在 Copy；
		// 会关底层 writer 时若先 Flush，对端不读则 Write 卡住，Flush 与 Close 死锁。
		flushErr = aww.BlockingFlush()
	}
	closeErr := aww.closeInternal()
	if flushErr != nil {
		return flushErr
	}
	return closeErr
}

func (aww *AsyncWriteWrapper) closeInternal() error {
	var closeErr error
	aww.once.Do(func() {
		aww.closed.Store(true)
		aww.mutex.Lock()
		select {
		case <-aww.done:
		default:
			close(aww.done)
		}
		if aww.cond != nil {
			aww.cond.Broadcast()
		}
		aww.ready.SetReady()
		aww.mutex.Unlock()

		// 必须先关底层 writer 才能打断阻塞 Write。若先 <-exited 再 Close，对端不读时会死锁。
		// 关完后再等 run 退出，避免 Close 返回后仍并发 Write。
		var clearer deadlineWriteSetter
		if aww.closeUnderlying {
			if c, ok := aww.writer.(io.Closer); ok {
				closeErr = c.Close()
			}
		} else if sd, ok := aww.writer.(deadlineWriteSetter); ok {
			// DontCloseWriter 时阻塞 Write 不会因 close(done) 返回
			_ = sd.SetWriteDeadline(time.Now())
			clearer = sd
		}
		<-aww.exited
		if clearer != nil {
			_ = clearer.SetWriteDeadline(time.Time{})
		}
	})
	return closeErr
}

// BlockingFlushContext 在 ctx 取消时返回。
// 不能和 BlockingFlush 互相调用，context.Background 的 Done 为 nil，会无限递归。
// 取消标志必须在同一把锁下设置并 Broadcast。只 Broadcast 一次且发生在 Wait 之前会丢唤醒；
// 成功路径不使用定时轮询。返回前等待辅助 goroutine 退出，避免 cancelled 数据竞争。
func (aww *AsyncWriteWrapper) BlockingFlushContext(ctx context.Context) error {
	if aww == nil {
		return io.ErrClosedPipe
	}
	if ctx == nil || ctx.Done() == nil {
		return aww.BlockingFlush()
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	stop := make(chan struct{})
	exited := make(chan struct{})
	var cancelled bool
	go func() {
		defer close(exited)
		select {
		case <-ctx.Done():
			aww.mutex.Lock()
			cancelled = true
			if aww.cond != nil {
				aww.cond.Broadcast()
			}
			aww.mutex.Unlock()
		case <-stop:
		}
	}()

	aww.mutex.Lock()
	var flushErr error
	for aww.pendingLocked() > 0 && aww.err == nil && !cancelled {
		select {
		case <-aww.done:
			if aww.err == nil && aww.pendingLocked() > 0 {
				aww.err = io.ErrClosedPipe
			}
			flushErr = aww.err
			aww.mutex.Unlock()
			close(stop)
			<-exited
			return flushErr
		default:
		}
		if aww.cond == nil {
			flushErr = aww.err
			break
		}
		aww.cond.Signal()
		aww.cond.Wait()
	}
	if flushErr == nil {
		if aww.pendingLocked() == 0 || aww.err != nil {
			flushErr = aww.err
		} else if cancelled {
			flushErr = ctx.Err()
			if flushErr == nil {
				flushErr = context.Canceled
			}
		}
	}
	aww.mutex.Unlock()
	close(stop)
	<-exited
	return flushErr
}

// NewAsyncStreamForWriter 是一个便捷辅助函数，
// 将一个阻塞的 io.Writer 转换为完全支持异步 subscribe 的 *Stream。
func NewAsyncStreamForWriter(w io.Writer, opts ...AsyncWriteWrapperOption) *Stream {
	wrapper := NewAsyncWriteWrapper(w, opts...)
	return &Stream{
		Writer:          wrapper,
		Closer:          wrapper,
		Flusher:         wrapper,
		BlockingFlusher: wrapper, // blocking-flush 需走到 BlockingFlush 而不是非阻塞 Flush
		CheckWriter:     wrapper,
		OnSubscribe: func() IPollable {
			return wrapper.subscribe()
		},
	}
}
