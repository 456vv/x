package v0_2

import (
	"context"
	"errors"
	"io"
	"math"
	"net"
	"os"
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

func (i *streamsImpl) subscribeToStream(this uint32) Pollable {
	s, ok := i.sm.Get(this)
	if !ok || s == nil {
		return i.pm.Add(manager_io.NewReadyPollable())
	}

	if s.OnSubscribe != nil {
		pollable := s.OnSubscribe()
		if pollable != nil {
			return i.pm.Add(pollable)
		}
	}

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
	if maxLen == 0 {
		return 0
	}
	if maxLen > uint64(^uint(0)>>1) || maxLen > maxChunk {
		return maxChunk
	}
	return int(maxLen)
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
		if isStreamClosedErr(err) {
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
			if isStreamClosedErr(err) {
				return witgo.Err[[]byte, StreamError](StreamError{Closed: &witgo.Unit{}})
			}
			errHandle := i.em.Add(err)
			return witgo.Err[[]byte, StreamError](StreamError{LastOperationFailed: &errHandle})
		}
	}
}

func (i *streamsImpl) Skip(ctx context.Context, this InputStream, maxLen uint64) witgo.Result[uint64, StreamError] {
	s, ok := i.sm.Get(this)
	if !ok || s == nil || s.Reader == nil {
		return witgo.Err[uint64, StreamError](StreamError{Closed: &witgo.Unit{}})
	}
	return i.skipByReading(ctx, s, maxLen, false)
}

func (i *streamsImpl) BlockingSkip(ctx context.Context, this InputStream, maxLen uint64) witgo.Result[uint64, StreamError] {
	s, ok := i.sm.Get(this)
	if !ok || s == nil || s.Reader == nil {
		return witgo.Err[uint64, StreamError](StreamError{Closed: &witgo.Unit{}})
	}
	return i.skipByReading(ctx, s, maxLen, true)
}

// sectionRemaining 在 Seeker 实现了 Size()（如 io.SectionReader）时给出剩余字节。
// read-via-stream 用 SectionReader；skip 若逐块读大偏移会浪费 CPU。
// 不用任意 Seeker.Seek：*os.File 作为 Seeker 会移动共享 fd，干扰同 fd 上的 ReadAt 流。
func sectionRemaining(s *manager_io.Stream) (remain int64, ok bool) {
	if s == nil || s.Seeker == nil {
		return 0, false
	}
	type sizer interface{ Size() int64 }
	sz, ok := s.Seeker.(sizer)
	if !ok {
		return 0, false
	}
	cur, err := s.Seeker.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, false
	}
	size := sz.Size()
	// 不知道真实文件长度时用 MaxInt64 占位。把它当成 section 末尾，skip 会报告跳过了不存在的字节。
	if size < 0 || size >= math.MaxInt64/2 {
		return 0, false
	}
	remain = size - cur
	if remain < 0 {
		remain = 0
	}
	return remain, true
}

func (i *streamsImpl) skipByReading(ctx context.Context, s *manager_io.Stream, maxLen uint64, blocking bool) witgo.Result[uint64, StreamError] {
	if maxLen == 0 {
		return witgo.Ok[uint64, StreamError](0)
	}

	if remain, ok := sectionRemaining(s); ok {
		if remain == 0 {
			return witgo.Err[uint64, StreamError](StreamError{Closed: &witgo.Unit{}})
		}
		n := remain
		if maxLen <= uint64(math.MaxInt64) && int64(maxLen) < n {
			n = int64(maxLen)
		}
		if _, err := s.Seeker.Seek(n, io.SeekCurrent); err != nil {
			if isStreamClosedErr(err) {
				return witgo.Err[uint64, StreamError](StreamError{Closed: &witgo.Unit{}})
			}
			errHandle := i.em.Add(err)
			return witgo.Err[uint64, StreamError](StreamError{LastOperationFailed: &errHandle})
		}
		return witgo.Ok[uint64, StreamError](uint64(n))
	}

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
			if isStreamClosedErr(err) {
				if totalSkipped == 0 {
					return witgo.Err[uint64, StreamError](StreamError{Closed: &witgo.Unit{}})
				}
				break
			}
			errHandle := i.em.Add(err)
			return witgo.Err[uint64, StreamError](StreamError{LastOperationFailed: &errHandle})
		}

		if n == 0 {
			if blocking {
				if err := waitSubscribe(ctx, s); err != nil {
					errHandle := i.em.Add(err)
					return witgo.Err[uint64, StreamError](StreamError{LastOperationFailed: &errHandle})
				}
			} else {
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
		avail := s.CheckWriter.CheckWrite()
		if avail == 0 {
			if se, ok := s.CheckWriter.(manager_io.StreamErrorer); ok {
				if err := se.StreamErr(); err != nil {
					if isStreamClosedErr(err) {
						return witgo.Err[uint64, StreamError](StreamError{Closed: &witgo.Unit{}})
					}
					errHandle := i.em.Add(err)
					return witgo.Err[uint64, StreamError](StreamError{LastOperationFailed: &errHandle})
				}
			}
		}
		return witgo.Ok[uint64, StreamError](avail)
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
		if uint64(len(contents)) > s.CheckWriter.CheckWrite() {
			errHandle := i.em.Add(errors.New("write exceeds check-write permit"))
			return witgo.Err[witgo.Unit, StreamError](StreamError{LastOperationFailed: &errHandle})
		}
	}
	off := 0
	for off < len(contents) {
		n, err := s.Writer.Write(contents[off:])
		if n > 0 {
			off += n
		}
		if err != nil {
			if isStreamClosedErr(err) {
				return witgo.Err[witgo.Unit, StreamError](StreamError{Closed: &witgo.Unit{}})
			}
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
			if se, ok := s.CheckWriter.(manager_io.StreamErrorer); ok {
				if err := se.StreamErr(); err != nil {
					if isStreamClosedErr(err) {
						return witgo.Err[witgo.Unit, StreamError](StreamError{Closed: &witgo.Unit{}})
					}
					errHandle := i.em.Add(err)
					return witgo.Err[witgo.Unit, StreamError](StreamError{LastOperationFailed: &errHandle})
				}
			}
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
			if isStreamClosedErr(err) {
				return witgo.Err[witgo.Unit, StreamError](StreamError{Closed: &witgo.Unit{}})
			}
			errHandle := i.em.Add(err)
			return witgo.Err[witgo.Unit, StreamError](StreamError{LastOperationFailed: &errHandle})
		}
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

func isStreamClosedErr(err error) bool {
	// 关闭文件是 os.ErrClosed，关闭管道是 ErrClosedPipe，关闭 TCP 是 net.ErrClosed。
	// 只认 io.EOF 会让 guest 把正常结束当成失败并不停重试。
	return err != nil && (errors.Is(err, io.EOF) || errors.Is(err, io.ErrClosedPipe) || errors.Is(err, net.ErrClosed) || errors.Is(err, os.ErrClosed))
}

func (i *streamsImpl) BlockingFlush(ctx context.Context, this OutputStream) witgo.Result[witgo.Unit, StreamError] {
	s, ok := i.sm.Get(this)
	if !ok || s == nil || s.Writer == nil {
		return witgo.Err[witgo.Unit, StreamError](StreamError{Closed: &witgo.Unit{}})
	}
	if s.BlockingFlusher != nil {
		var err error
		type ctxFlusher interface {
			BlockingFlushContext(context.Context) error
		}
		if cf, ok := s.BlockingFlusher.(ctxFlusher); ok {
			err = cf.BlockingFlushContext(ctx)
		} else {
			err = s.BlockingFlusher.BlockingFlush()
		}
		if err != nil {
			if isStreamClosedErr(err) {
				return witgo.Err[witgo.Unit, StreamError](StreamError{Closed: &witgo.Unit{}})
			}
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
// 修改原因：WASI splice 是非阻塞的；原先 waitSubscribe 等于 blocking-splice，
// 会在 wasm 线程上挂死，且与 subscribe 轮询冲突。
func (i *streamsImpl) Splice(ctx context.Context, this OutputStream, src InputStream, maxLen uint64) witgo.Result[uint64, StreamError] {
	return i.splice(ctx, this, src, maxLen, false)
}

// BlockingSplice 在源无数据或目标不可写时等待 pollable。
func (i *streamsImpl) BlockingSplice(ctx context.Context, this OutputStream, src InputStream, maxLen uint64) witgo.Result[uint64, StreamError] {
	return i.splice(ctx, this, src, maxLen, true)
}

func (i *streamsImpl) splice(ctx context.Context, this OutputStream, src InputStream, maxLen uint64, blocking bool) witgo.Result[uint64, StreamError] {
	dst, ok := i.sm.Get(this)
	if !ok || dst == nil || dst.Writer == nil {
		return witgo.Err[uint64, StreamError](StreamError{Closed: &witgo.Unit{}})
	}

	srcStream, okSrc := i.sm.Get(src)
	// 修改原因：src==0 原是 WriteZeroes 内部 hack；句柄 0 非法，不能当成无限零流
	if !okSrc || srcStream == nil || srcStream.Reader == nil {
		return witgo.Err[uint64, StreamError](StreamError{Closed: &witgo.Unit{}})
	}
	srcReader := srcStream.Reader

	var totalWritten uint64
	buf := make([]byte, 32*1024)

	for totalWritten < maxLen {
		writePermit := uint64(4096)
		if dst.CheckWriter != nil {
			writePermit = dst.CheckWriter.CheckWrite()
		}
		if writePermit == 0 {
			if se, ok := dst.CheckWriter.(manager_io.StreamErrorer); ok {
				if err := se.StreamErr(); err != nil {
					if isStreamClosedErr(err) {
						return witgo.Err[uint64, StreamError](StreamError{Closed: &witgo.Unit{}})
					}
					errHandle := i.em.Add(err)
					return witgo.Err[uint64, StreamError](StreamError{LastOperationFailed: &errHandle})
				}
			}
			if !blocking {
				break
			}
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

		n, err := srcReader.Read(buf[:chunk])
		if n == 0 && err == nil {
			if !blocking {
				break
			}
			if werr := waitSubscribe(ctx, srcStream); werr != nil {
				errHandle := i.em.Add(werr)
				return witgo.Err[uint64, StreamError](StreamError{LastOperationFailed: &errHandle})
			}
			continue
		}
		if n == 0 && err != nil {
			if errors.Is(err, io.EOF) || isStreamClosedErr(err) {
				if totalWritten == 0 {
					return witgo.Err[uint64, StreamError](StreamError{Closed: &witgo.Unit{}})
				}
				break
			}
			errHandle := i.em.Add(err)
			return witgo.Err[uint64, StreamError](StreamError{LastOperationFailed: &errHandle})
		}

		off := 0
		for off < n {
			wn, werr := dst.Writer.Write(buf[off:n])
			if wn > 0 {
				off += wn
				totalWritten += uint64(wn)
			}
			if werr != nil {
				if isStreamClosedErr(werr) {
					return witgo.Err[uint64, StreamError](StreamError{Closed: &witgo.Unit{}})
				}
				errHandle := i.em.Add(werr)
				return witgo.Err[uint64, StreamError](StreamError{LastOperationFailed: &errHandle})
			}
			if wn == 0 {
				// 修改原因：已从 src 读出的字节不能丢。即使 splice 非阻塞，
				// 也必须把本轮已读数据写完；否则数据从源抽走却未送达。
				if waitErr := waitSubscribe(ctx, dst); waitErr != nil {
					errHandle := i.em.Add(waitErr)
					return witgo.Err[uint64, StreamError](StreamError{LastOperationFailed: &errHandle})
				}
			}
		}

		if errors.Is(err, io.EOF) {
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
			return ctx.Err()
		}
	}
	return nil
}

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

// zeroChunk 全零且只读，并发 WriteZeroes 可共享，避免每次 make。
var zeroChunk = make([]byte, 32*1024)

// WriteZeroes 与 write 相同：不得超过 check-write，不得阻塞等待。
func (i *streamsImpl) WriteZeroes(ctx context.Context, this OutputStream, length uint64) witgo.Result[witgo.Unit, StreamError] {
	s, ok := i.sm.Get(this)
	if !ok || s == nil || s.Writer == nil {
		return witgo.Err[witgo.Unit, StreamError](StreamError{Closed: &witgo.Unit{}})
	}
	if length == 0 {
		return witgo.Ok[witgo.Unit, StreamError](witgo.Unit{})
	}
	if s.CheckWriter != nil {
		permit := s.CheckWriter.CheckWrite()
		if length > permit {
			if permit == 0 {
				if se, ok := s.CheckWriter.(manager_io.StreamErrorer); ok {
					if err := se.StreamErr(); err != nil {
						if isStreamClosedErr(err) {
							return witgo.Err[witgo.Unit, StreamError](StreamError{Closed: &witgo.Unit{}})
						}
						errHandle := i.em.Add(err)
						return witgo.Err[witgo.Unit, StreamError](StreamError{LastOperationFailed: &errHandle})
					}
				}
			}
			// 修改原因：write-zeroes 不返回写入计数；permit=0 且 len>0 不能 Ok（guest 会当成已写完）
			errHandle := i.em.Add(errors.New("write-zeroes exceeds check-write permit"))
			return witgo.Err[witgo.Unit, StreamError](StreamError{LastOperationFailed: &errHandle})
		}
	}
	remaining := length
	for remaining > 0 {
		n := len(zeroChunk)
		if remaining < uint64(n) {
			n = int(remaining)
		}
		res := i.Write(ctx, this, zeroChunk[:n])
		if res.Err != nil {
			return witgo.Err[witgo.Unit, StreamError](*res.Err)
		}
		remaining -= uint64(n)
	}
	return witgo.Ok[witgo.Unit, StreamError](witgo.Unit{})
}

// BlockingWriteZeroesAndFlush 写满 len 个零字节并阻塞 flush。
// 修改原因：不能再委托非阻塞 WriteZeroes，否则 check-write 不足时会失败而非等待。
func (i *streamsImpl) BlockingWriteZeroesAndFlush(ctx context.Context, this OutputStream, length uint64) witgo.Result[witgo.Unit, StreamError] {
	s, ok := i.sm.Get(this)
	if !ok || s == nil || s.Writer == nil {
		return witgo.Err[witgo.Unit, StreamError](StreamError{Closed: &witgo.Unit{}})
	}

	writeSize := uint64(4096)
	remaining := length
	for remaining > 0 {
		if s.CheckWriter != nil {
			writeSize = s.CheckWriter.CheckWrite()
		}
		if writeSize == 0 {
			if se, ok := s.CheckWriter.(manager_io.StreamErrorer); ok {
				if err := se.StreamErr(); err != nil {
					if isStreamClosedErr(err) {
						return witgo.Err[witgo.Unit, StreamError](StreamError{Closed: &witgo.Unit{}})
					}
					errHandle := i.em.Add(err)
					return witgo.Err[witgo.Unit, StreamError](StreamError{LastOperationFailed: &errHandle})
				}
			}
			if err := waitSubscribe(ctx, s); err != nil {
				errHandle := i.em.Add(err)
				return witgo.Err[witgo.Unit, StreamError](StreamError{LastOperationFailed: &errHandle})
			}
			continue
		}

		chunk := remaining
		if chunk > writeSize {
			chunk = writeSize
		}
		if chunk > uint64(len(zeroChunk)) {
			chunk = uint64(len(zeroChunk))
		}
		n, err := s.Writer.Write(zeroChunk[:chunk])
		if err != nil {
			if isStreamClosedErr(err) {
				return witgo.Err[witgo.Unit, StreamError](StreamError{Closed: &witgo.Unit{}})
			}
			errHandle := i.em.Add(err)
			return witgo.Err[witgo.Unit, StreamError](StreamError{LastOperationFailed: &errHandle})
		}
		if n == 0 {
			if err := waitSubscribe(ctx, s); err != nil {
				errHandle := i.em.Add(err)
				return witgo.Err[witgo.Unit, StreamError](StreamError{LastOperationFailed: &errHandle})
			}
			continue
		}
		remaining -= uint64(n)
	}
	return i.BlockingFlush(ctx, this)
}
