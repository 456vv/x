package v0_2

import (
	"context"
	"errors"
	"io"
	"math"
	"time"

	manager_io "github.com/456vv/x/vweb_dynamic/internal/wazero/manager/io"
	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"
)

type streamsImpl struct {
	sm *manager_io.StreamManager
	em *manager_io.ErrorManager
	pm *manager_io.PollManager
}

func newStreamsImpl(sm *manager_io.StreamManager, em *manager_io.ErrorManager, pm *manager_io.PollManager) *streamsImpl {
	return &streamsImpl{sm: sm, em: em, pm: pm}
}

// subscribeToStream 是 SubscribeToInputStream 和 SubscribeToOutputStream 的通用实现。
func (i *streamsImpl) subscribeToStream(this uint32) Pollable {
	s, ok := i.sm.Get(this)
	if !ok || s == nil {
		// 无效的 stream 句柄，返回一个立即就绪的 pollable。
		return i.pm.Add(manager_io.NewReadyPollable())
	}

	// 如果 stream 的创建者提供了 OnSubscribe 回调，则调用它。
	if s.OnSubscribe != nil {
		pollable := s.OnSubscribe()
		if pollable != nil {
			return i.pm.Add(pollable)
		}
	}

	// 否则，回退到默认行为：为通用阻塞流创建一个立即就绪的 pollable。
	return i.pm.Add(manager_io.NewReadyPollable())
}

func (i *streamsImpl) DropInputStream(_ context.Context, handle InputStream) {
	i.sm.Remove(handle)
}

func (i *streamsImpl) DropOutputStream(_ context.Context, handle OutputStream) {
	i.sm.Remove(handle)
}

func capReadLen(maxLen uint64) int {
	const maxChunk = 8 << 20
	// WIT 允许短读；32 位上 make([]byte, uint64) 会溢出 panic，超大分配会 OOM
	if maxLen == 0 {
		return 0
	}
	if maxLen > uint64(^uint(0)>>1) || maxLen > maxChunk {
		return maxChunk
	}
	return int(maxLen)
}

func minU64(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

func (i *streamsImpl) Read(_ context.Context, this InputStream, maxLen uint64) witgo.Result[[]byte, StreamError] {
	s, ok := i.sm.Get(this)
	if !ok || s == nil || s.Reader == nil {
		return witgo.Err[[]byte, StreamError](StreamError{Closed: &witgo.Unit{}})
	}
	ncap := capReadLen(maxLen)
	if ncap == 0 {
		return witgo.Ok[[]byte, StreamError]([]byte{})
	}
	buf := make([]byte, ncap)
	n, err := s.Reader.Read(buf)
	if n == 0 && err != nil {
		if err == io.EOF {
			return witgo.Err[[]byte, StreamError](StreamError{Closed: &witgo.Unit{}})
		}
		errHandle := i.em.Add(err)
		return witgo.Err[[]byte, StreamError](StreamError{LastOperationFailed: &errHandle})
	}
	return witgo.Ok[[]byte, StreamError](buf[:n])
}

func (i *streamsImpl) BlockingRead(ctx context.Context, this InputStream, maxLen uint64) witgo.Result[[]byte, StreamError] {
	s, ok := i.sm.Get(this)
	if !ok || s == nil || s.Reader == nil {
		return witgo.Err[[]byte, StreamError](StreamError{Closed: &witgo.Unit{}})
	}

	if maxLen == 0 {
		return witgo.Ok[[]byte, StreamError]([]byte{})
	}

	for {
		if err := waitSubscribe(ctx, s); err != nil {
			errHandle := i.em.Add(err)
			return witgo.Err[[]byte, StreamError](StreamError{LastOperationFailed: &errHandle})
		}

		ncap := capReadLen(maxLen)
		if ncap == 0 {
			return witgo.Ok[[]byte, StreamError]([]byte{})
		}
		buf := make([]byte, ncap)
		n, err := s.Reader.Read(buf)
		if n > 0 {
			return witgo.Ok[[]byte, StreamError](buf[:n])
		}
		if err != nil {
			if err == io.EOF {
				return witgo.Err[[]byte, StreamError](StreamError{Closed: &witgo.Unit{}})
			}
			errHandle := i.em.Add(err)
			return witgo.Err[[]byte, StreamError](StreamError{LastOperationFailed: &errHandle})
		}
		// AsyncReadWrapper 无数据返回 (0,nil)；blocking-read 必须等到至少 1 字节或关闭
	}
}

// Skip (非阻塞) 尝试跳过最多 maxLen 字节并立即返回。
func (i *streamsImpl) Skip(ctx context.Context, this InputStream, maxLen uint64) witgo.Result[uint64, StreamError] {
	s, ok := i.sm.Get(this)
	if !ok || s == nil || s.Reader == nil {
		return witgo.Err[uint64, StreamError](StreamError{Closed: &witgo.Unit{}})
	}

	// 优先使用 Seeker 实现高效跳转。
	if s.Seeker != nil {
		currentPos, err := s.Seeker.Seek(0, io.SeekCurrent)
		if err == nil {
			skip := minU64(maxLen, uint64(math.MaxInt64))
			newPos, err := s.Seeker.Seek(int64(skip), io.SeekCurrent)
			if err == nil {
				return witgo.Ok[uint64, StreamError](uint64(newPos - currentPos))
			}
		}
	}

	// 回退到读取和丢弃方法，并指定为非阻塞模式。
	return i.skipByReading(ctx, s, maxLen, false) // blocking = false
}

// BlockingSkip (阻塞) 会跳过 maxLen 字节，并在必要时等待数据。
func (i *streamsImpl) BlockingSkip(ctx context.Context, this InputStream, maxLen uint64) witgo.Result[uint64, StreamError] {
	s, ok := i.sm.Get(this)
	if !ok || s == nil || s.Reader == nil {
		return witgo.Err[uint64, StreamError](StreamError{Closed: &witgo.Unit{}})
	}

	// 优先使用 Seeker 实现高效跳转。
	if s.Seeker != nil {
		currentPos, err := s.Seeker.Seek(0, io.SeekCurrent)
		if err == nil {
			skip := minU64(maxLen, uint64(math.MaxInt64))
			newPos, err := s.Seeker.Seek(int64(skip), io.SeekCurrent)
			if err == nil {
				return witgo.Ok[uint64, StreamError](uint64(newPos - currentPos))
			}
		}
	}

	// 回退到读取和丢弃方法，并指定为阻塞模式。
	return i.skipByReading(ctx, s, maxLen, true) // blocking = true
}

// skipByReading 是跳过字节的核心实现，支持阻塞和非阻塞两种模式。
func (i *streamsImpl) skipByReading(ctx context.Context, s *manager_io.Stream, maxLen uint64, blocking bool) witgo.Result[uint64, StreamError] {
	var totalSkipped uint64
	buf := make([]byte, 32*1024)

	for totalSkipped < maxLen {
		readSize := uint64(len(buf))
		if remaining := maxLen - totalSkipped; remaining < readSize {
			readSize = remaining
		}

		n, err := s.Reader.Read(buf[:readSize])
		if n > 0 {
			totalSkipped += uint64(n)
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			errHandle := i.em.Add(err)
			return witgo.Err[uint64, StreamError](StreamError{LastOperationFailed: &errHandle})
		}

		// 当 n == 0 时，根据 blocking 参数决定行为
		if n == 0 {
			if blocking {
				if err := waitSubscribe(ctx, s); err != nil {
					errHandle := i.em.Add(err)
					return witgo.Err[uint64, StreamError](StreamError{LastOperationFailed: &errHandle})
				}
			} else {
				// 非阻塞模式：立即停止
				break
			}
		}
	}

	return witgo.Ok[uint64, StreamError](totalSkipped)
}

func (i *streamsImpl) SubscribeToInputStream(_ context.Context, this InputStream) Pollable {
	return i.subscribeToStream(this)
}

func (i *streamsImpl) CheckWrite(_ context.Context, this OutputStream) witgo.Result[uint64, StreamError] {
	s, ok := i.sm.Get(this)
	if !ok || s == nil || s.Writer == nil {
		return witgo.Err[uint64, StreamError](StreamError{Closed: &witgo.Unit{}})
	}
	if s.CheckWriter != nil {
		return witgo.Ok[uint64, StreamError](s.CheckWriter.CheckWrite())
	}
	return witgo.Ok[uint64, StreamError](4096)
}

func (i *streamsImpl) Write(_ context.Context, this OutputStream, contents []byte) witgo.Result[witgo.Unit, StreamError] {
	s, ok := i.sm.Get(this)
	if !ok || s == nil || s.Writer == nil {
		return witgo.Err[witgo.Unit, StreamError](StreamError{Closed: &witgo.Unit{}})
	}
	if len(contents) == 0 {
		return witgo.Ok[witgo.Unit, StreamError](witgo.Unit{})
	}
	if s.CheckWriter != nil {
		// WASI 规定 write 长度不得超过 check-write 许可
		if uint64(len(contents)) > s.CheckWriter.CheckWrite() {
			errHandle := i.em.Add(errors.New("write exceeds check-write permit"))
			return witgo.Err[witgo.Unit, StreamError](StreamError{LastOperationFailed: &errHandle})
		}
	}
	// io.Writer 允许短写且 err==nil；原先忽略 n 会丢数据。
	off := 0
	for off < len(contents) {
		n, err := s.Writer.Write(contents[off:])
		if n > 0 {
			off += n
		}
		if err != nil {
			errHandle := i.em.Add(err)
			return witgo.Err[witgo.Unit, StreamError](StreamError{LastOperationFailed: &errHandle})
		}
		if n == 0 {
			errHandle := i.em.Add(io.ErrShortWrite)
			return witgo.Err[witgo.Unit, StreamError](StreamError{LastOperationFailed: &errHandle})
		}
	}
	return witgo.Ok[witgo.Unit, StreamError](witgo.Unit{})
}

func (i *streamsImpl) BlockingWriteAndFlush(ctx context.Context, this OutputStream, contents []byte) witgo.Result[witgo.Unit, StreamError] {
	s, ok := i.sm.Get(this)
	if !ok || s == nil || s.Writer == nil {
		return witgo.Err[witgo.Unit, StreamError](StreamError{Closed: &witgo.Unit{}})
	}

	writeSize := uint64(4096)
	for len(contents) > 0 {
		if s.CheckWriter != nil {
			writeSize = s.CheckWriter.CheckWrite()
		}
		if writeSize == 0 {
			if err := waitSubscribe(ctx, s); err != nil {
				errHandle := i.em.Add(err)
				return witgo.Err[witgo.Unit, StreamError](StreamError{LastOperationFailed: &errHandle})
			}
			continue
		}

		chunk := uint64(len(contents))
		if chunk > writeSize {
			chunk = writeSize
		}
		n, err := s.Writer.Write(contents[:chunk])
		if err != nil {
			errHandle := i.em.Add(err)
			return witgo.Err[witgo.Unit, StreamError](StreamError{LastOperationFailed: &errHandle})
		}
		// 短写时推进已写部分，避免死循环或丢数据
		if n == 0 {
			if err := waitSubscribe(ctx, s); err != nil {
				errHandle := i.em.Add(err)
				return witgo.Err[witgo.Unit, StreamError](StreamError{LastOperationFailed: &errHandle})
			}
			continue
		}
		contents = contents[n:]
	}
	return i.BlockingFlush(ctx, this)
}

func (i *streamsImpl) Flush(_ context.Context, this OutputStream) witgo.Result[witgo.Unit, StreamError] {
	s, ok := i.sm.Get(this)
	if !ok || s == nil || s.Writer == nil {
		return witgo.Err[witgo.Unit, StreamError](StreamError{Closed: &witgo.Unit{}})
	}

	if s.Flusher != nil {
		if err := s.Flusher.Flush(); err != nil {
			errHandle := i.em.Add(err)
			return witgo.Err[witgo.Unit, StreamError](StreamError{LastOperationFailed: &errHandle})
		}
	}

	return witgo.Ok[witgo.Unit, StreamError](witgo.Unit{})
}

func (i *streamsImpl) BlockingFlush(ctx context.Context, this OutputStream) witgo.Result[witgo.Unit, StreamError] {
	s, ok := i.sm.Get(this)
	if !ok || s == nil || s.Writer == nil {
		return witgo.Err[witgo.Unit, StreamError](StreamError{Closed: &witgo.Unit{}})
	}
	// 原先直接调非阻塞 Flush，AsyncWriteWrapper 缓冲未排空就返回 Ok。
	if s.BlockingFlusher != nil {
		if err := s.BlockingFlusher.BlockingFlush(); err != nil {
			errHandle := i.em.Add(err)
			return witgo.Err[witgo.Unit, StreamError](StreamError{LastOperationFailed: &errHandle})
		}
		return witgo.Ok[witgo.Unit, StreamError](witgo.Unit{})
	}
	return i.Flush(ctx, this)
}

func (i *streamsImpl) SubscribeToOutputStream(_ context.Context, this OutputStream) Pollable {
	return i.subscribeToStream(this)
}

// Splice 将最多 maxLen 字节从 src 拷到 this。
// 原先用 io.CopyN。AsyncReadWrapper 在无数据时返回 (0, nil)，
// CopyN 会立刻重试，形成忙等。改为显式等待 src/dst 的 subscribe。
func (i *streamsImpl) Splice(ctx context.Context, this OutputStream, src InputStream, maxLen uint64) witgo.Result[uint64, StreamError] {
	dst, ok := i.sm.Get(this)
	if !ok || dst == nil || dst.Writer == nil {
		return witgo.Err[uint64, StreamError](StreamError{Closed: &witgo.Unit{}})
	}

	var srcStream *manager_io.Stream
	var srcReader io.Reader
	if src == 0 {
		// WriteZeroes 走这里：无穷零字节，不会 (0, nil)
		srcReader = zeroReader{}
	} else {
		var okSrc bool
		srcStream, okSrc = i.sm.Get(src)
		if !okSrc || srcStream.Reader == nil {
			return witgo.Err[uint64, StreamError](StreamError{Closed: &witgo.Unit{}})
		}
		srcReader = srcStream.Reader
	}

	var totalWritten uint64
	buf := make([]byte, 32*1024)

	for totalWritten < maxLen {
		// 1. 目标可写空间
		writePermit := uint64(4096)
		if dst.CheckWriter != nil {
			writePermit = dst.CheckWriter.CheckWrite()
		}
		if writePermit == 0 {
			if err := waitSubscribe(ctx, dst); err != nil {
				errHandle := i.em.Add(err)
				return witgo.Err[uint64, StreamError](StreamError{LastOperationFailed: &errHandle})
			}
			continue
		}

		remaining := maxLen - totalWritten
		chunk := writePermit
		if remaining < chunk {
			chunk = remaining
		}
		if chunk > uint64(len(buf)) {
			chunk = uint64(len(buf))
		}

		// 2. 从源读取；n==0 && err==nil 表示暂时没数据，必须阻塞等待，不能忙等
		n, err := srcReader.Read(buf[:chunk])
		if n == 0 && err == nil {
			if srcStream != nil {
				if werr := waitSubscribe(ctx, srcStream); werr != nil {
					errHandle := i.em.Add(werr)
					return witgo.Err[uint64, StreamError](StreamError{LastOperationFailed: &errHandle})
				}
			} else {
				// zeroReader 不应走到这里
				time.Sleep(20 * time.Millisecond)
			}
			continue
		}
		if n == 0 && err != nil {
			if err == io.EOF {
				break
			}
			errHandle := i.em.Add(err)
			return witgo.Err[uint64, StreamError](StreamError{LastOperationFailed: &errHandle})
		}

		// 3. 写入已读到的数据（可能短写，短写同样要等可写）
		off := 0
		for off < n {
			wn, werr := dst.Writer.Write(buf[off:n])
			if wn > 0 {
				off += wn
				totalWritten += uint64(wn)
			}
			if werr != nil {
				errHandle := i.em.Add(werr)
				return witgo.Err[uint64, StreamError](StreamError{LastOperationFailed: &errHandle})
			}
			if wn == 0 {
				if waitErr := waitSubscribe(ctx, dst); waitErr != nil {
					errHandle := i.em.Add(waitErr)
					return witgo.Err[uint64, StreamError](StreamError{LastOperationFailed: &errHandle})
				}
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			errHandle := i.em.Add(err)
			return witgo.Err[uint64, StreamError](StreamError{LastOperationFailed: &errHandle})
		}
	}
	return witgo.Ok[uint64, StreamError](totalWritten)
}

func waitPollable(ctx context.Context, p manager_io.IPollable) error {
	if p == nil {
		time.Sleep(20 * time.Millisecond)
		return nil
	}
	for !p.IsReady() {
		ch := p.Channel()
		if ctx == nil {
			<-ch
			continue
		}
		select {
		case <-ch:
		case <-ctx.Done():
			return ctx.Err() // 阻塞 skip/splice/write 忽略 ctx 会在 guest 取消后挂死 host
		}
	}
	return nil
}

// waitSubscribe 阻塞到流就绪；无 OnSubscribe 时短睡眠兜底。
func waitSubscribe(ctx context.Context, s *manager_io.Stream) error {
	if s != nil && s.OnSubscribe != nil {
		if p := s.OnSubscribe(); p != nil {
			return waitPollable(ctx, p)
		}
	}
	if ctx != nil {
		if err := ctx.Err(); err != nil {
			return err
		}
	}
	time.Sleep(20 * time.Millisecond)
	return nil
}

// BlockingSplice 与 Splice 行为一致（写入路径已按 pollable 阻塞）。
func (i *streamsImpl) BlockingSplice(ctx context.Context, this OutputStream, src InputStream, maxLen uint64) witgo.Result[uint64, StreamError] {
	return i.Splice(ctx, this, src, maxLen)
}

// WriteZeroes 从 Splice 继承了新的阻塞行为。
func (i *streamsImpl) WriteZeroes(ctx context.Context, this OutputStream, len uint64) witgo.Result[witgo.Unit, StreamError] {
	// 调用阻塞式的 Splice，并使用一个虚拟的零字节流 (src=0) 作为源。
	if spliceResult := i.Splice(ctx, this, 0, len); spliceResult.Err != nil {
		return witgo.Err[witgo.Unit, StreamError](*spliceResult.Err)
	}
	return witgo.Ok[witgo.Unit, StreamError](witgo.Unit{})
}

// BlockingWriteZeroesAndFlush
func (i *streamsImpl) BlockingWriteZeroesAndFlush(ctx context.Context, this OutputStream, len uint64) witgo.Result[witgo.Unit, StreamError] {
	if writeResult := i.WriteZeroes(ctx, this, len); writeResult.Err != nil {
		return writeResult
	}
	return i.BlockingFlush(ctx, this)
}

type zeroReader struct{}

func (z zeroReader) Read(p []byte) (n int, err error) {
	// clear(p) 需要 Go 1.21；手写清零兼容更旧编译器
	for i := range p {
		p[i] = 0
	}
	return len(p), nil
}
