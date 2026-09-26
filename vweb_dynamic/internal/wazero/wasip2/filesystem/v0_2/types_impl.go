package v0_2

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"io"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/456vv/x/vweb_dynamic/internal/wazero/manager/filesystem"
	manager_io "github.com/456vv/x/vweb_dynamic/internal/wazero/manager/io"
	"github.com/456vv/x/vweb_dynamic/internal/wazero/wasip2"
	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"
)

type typesImpl struct {
	host *wasip2.Host
}

func newTypesImpl(h *wasip2.Host) *typesImpl {
	return &typesImpl{host: h}
}

/// TODO:
/// 1. 完善沙盒机制（resolveSandboxedPath 为第一阶段防护）
/// 2. 完善poll机制(可能没有必要)

// errAdviseUnsupported 表示当前平台没有 posix_fadvise/Fadvise。
var errAdviseUnsupported = errors.New("file advise is not supported on this platform")

// capReadLen 限制单次 read 的分配大小。
// WIT 允许短读；32 位上 make([]byte, uint64) 会溢出 panic，超大 length 会 OOM。
func capReadLen(maxLen uint64) int {
	const maxChunk = 8 << 20 // 8MiB
	if maxLen == 0 {
		return 0
	}
	if maxLen > uint64(^uint(0)>>1) || maxLen > maxChunk {
		return maxChunk
	}
	return int(maxLen)
}

// filesizeToInt64 把 WASI filesize 转给 Go 的 int64 API。
// uint64→int64 溢出后变成负偏移，ReadAt/Truncate/SectionReader 行为未定义。
func filesizeToInt64(v Filesize) (int64, ErrorCode) {
	if uint64(v) > uint64(math.MaxInt64) {
		return 0, ErrorCodeOverflow
	}
	return int64(v), 0
}

// lexicallyContained 判断 target 是否落在 base 目录树内（含 base 自身）。
func lexicallyContained(base, target string) bool {
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return false
	}
	rel = filepath.Clean(rel)
	sep := string(os.PathSeparator)
	return rel != ".." && !strings.HasPrefix(rel, ".."+sep)
}

func inSandbox(realBase, candidate string) bool {
	c := filepath.Clean(candidate)
	b := filepath.Clean(realBase)
	return c == b || lexicallyContained(b, c)
}

// splitGuestComponents 按 filepath 规则拆相对路径分量。
func splitGuestComponents(cleanGuest string) []string {
	if cleanGuest == "" || cleanGuest == "." {
		return nil
	}
	var comps []string
	rest := cleanGuest
	for rest != "" && rest != "." {
		base := filepath.Base(rest)
		dir := filepath.Dir(rest)
		if base == string(os.PathSeparator) || base == "/" || base == `\` {
			break
		}
		comps = append([]string{base}, comps...)
		if dir == rest {
			break
		}
		rest = dir
	}
	return comps
}

// resolveSandboxedPathEx 将 guest 相对路径解析到基目录下并拒绝逃逸。
// follow 只作用于最后一截；中间 symlink 必须解析并检查仍在 realBase 内。
func resolveSandboxedPathEx(baseDir, guestPath string, follow bool) (string, ErrorCode) {
	if guestPath == "" {
		return "", ErrorCodeInvalid
	}
	// 路径嵌入 NUL 时，C/syscall 会截断，可能绕过 ".." 检查
	if strings.IndexByte(guestPath, 0) >= 0 {
		return "", ErrorCodeInvalid
	}
	// 拒绝绝对路径与 Windows 盘符路径
	if filepath.IsAbs(guestPath) || (len(guestPath) >= 2 && guestPath[1] == ':') {
		return "", ErrorCodeNotPermitted
	}
	cleanGuest := filepath.Clean(guestPath)
	if cleanGuest == ".." || strings.HasPrefix(cleanGuest, ".."+string(os.PathSeparator)) {
		return "", ErrorCodeNotPermitted
	}

	absBase, err := filepath.Abs(baseDir)
	if err != nil {
		return "", mapOsError(err)
	}
	realBase, err := filepath.EvalSymlinks(absBase)
	if err != nil {
		realBase = absBase
	}
	realBase = filepath.Clean(realBase)

	comps := splitGuestComponents(cleanGuest)
	if len(comps) == 0 {
		return realBase, 0
	}

	current := realBase
	for i, comp := range comps {
		if comp == ".." || comp == "." || comp == "" {
			return "", ErrorCodeNotPermitted
		}
		nextAbs, err := filepath.Abs(filepath.Join(current, comp))
		if err != nil {
			return "", mapOsError(err)
		}
		last := i == len(comps)-1
		info, err := os.Lstat(nextAbs)
		if err != nil {
			if os.IsNotExist(err) {
				if !inSandbox(realBase, nextAbs) {
					return "", ErrorCodeNotPermitted
				}
				if last {
					return nextAbs, 0
				}
				return "", mapOsError(err)
			}
			return "", mapOsError(err)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			if !follow && last {
				// 最后一截不跟随：把 symlink 本身交给 open(O_NOFOLLOW)/Lstat
				if !inSandbox(realBase, nextAbs) {
					return "", ErrorCodeNotPermitted
				}
				return nextAbs, 0
			}
			// 仅检查最终词法路径拦不住 “dir/symlink -> /etc” 再拼 “passwd”
			resolved, rerr := filepath.EvalSymlinks(nextAbs)
			if rerr != nil {
				return "", mapOsError(rerr)
			}
			resolvedAbs, aerr := filepath.Abs(resolved)
			if aerr != nil {
				return "", mapOsError(aerr)
			}
			if !inSandbox(realBase, resolvedAbs) {
				return "", ErrorCodeNotPermitted
			}
			current = resolvedAbs
			continue
		}
		if !inSandbox(realBase, nextAbs) {
			return "", ErrorCodeNotPermitted
		}
		current = nextAbs
	}
	return current, 0
}

// unixTimeFromDatetime 避免 uint64 秒/纳秒转 int64 溢出成 1970 前的负时间。
func unixTimeFromDatetime(dt Datetime) time.Time {
	sec := dt.Seconds
	if sec > uint64(math.MaxInt64) {
		sec = uint64(math.MaxInt64)
	}
	ns := int64(dt.Nanoseconds)
	if ns < 0 {
		ns = 0
	}
	if ns >= 1e9 {
		ns = 1e9 - 1
	}
	return time.Unix(int64(sec), ns)
}

func datetimeFromUnix(sec, nsec int64) Datetime {
	// 负的 Timespec 直接转 uint64 会变成 2^64 附近，guest 看到的是几百万年。
	if nsec < 0 {
		borrow := (-nsec + 1_000_000_000 - 1) / 1_000_000_000
		sec -= borrow
		nsec += borrow * 1_000_000_000
	}
	if nsec >= 1_000_000_000 {
		sec += nsec / 1_000_000_000
		nsec %= 1_000_000_000
	}
	if sec < 0 {
		return Datetime{}
	}
	return Datetime{Seconds: uint64(sec), Nanoseconds: uint32(nsec)}
}

func datetimeFromTime(t time.Time) Datetime {
	return datetimeFromUnix(t.Unix(), int64(t.Nanosecond()))
}

func chtimesPath(path string, atime, mtime time.Time, follow bool) error {
	if follow {
		return os.Chtimes(path, atime, mtime)
	}
	// os.Chtimes 会跟随符号链接；WASI symlink-follow=false 应对链接本身。
	return lutimes(path, atime, mtime)
}

// applyNewTimestamps 把 WASI NewTimestamp 落到 atime/mtime。
// 三个 case 都为空时原先 atime/mtime 为零值 epoch，会把时间戳打成 1970。
func applyNewTimestamps(path string, pathFlags PathFlags, access, modification NewTimestamp) (atime, mtime time.Time, err error) {
	needCurrent := access.Now == nil && access.Timestamp == nil ||
		modification.Now == nil && modification.Timestamp == nil
	var info fs.FileInfo
	if needCurrent {
		if pathFlags.SymlinkFollow {
			info, err = os.Stat(path)
		} else {
			info, err = os.Lstat(path)
		}
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		mtime = info.ModTime()
		atime = mtime
		if t, aerr := GetATime(info); aerr == nil {
			atime = t
		}
	}
	now := time.Now()
	switch {
	case access.Timestamp != nil:
		atime = unixTimeFromDatetime(*access.Timestamp)
	case access.Now != nil:
		atime = now
	}
	switch {
	case modification.Timestamp != nil:
		mtime = unixTimeFromDatetime(*modification.Timestamp)
	case modification.Now != nil:
		mtime = now
	}
	return atime, mtime, nil
}

// metadataHashFromInfo 实现 wasi:filesystem metadata-hash（实现定义的 128bit）。
// 原先一律 Unsupported，guest 无法做缓存失效。
func metadataHashFromInfo(info fs.FileInfo, path string) MetadataHashValue {
	h := sha256.New()
	_, _ = io.WriteString(h, path)
	var tmp [16]byte
	binary.LittleEndian.PutUint64(tmp[0:8], uint64(info.Size()))
	binary.LittleEndian.PutUint64(tmp[8:16], uint64(info.ModTime().UnixNano()))
	_, _ = h.Write(tmp[:])
	var mode [4]byte
	binary.LittleEndian.PutUint32(mode[:], uint32(info.Mode()))
	_, _ = h.Write(mode[:])
	writePlatformFileIdentity(h, info)
	sum := h.Sum(nil)
	return MetadataHashValue{
		Lower: binary.LittleEndian.Uint64(sum[0:8]),
		Upper: binary.LittleEndian.Uint64(sum[8:16]),
	}
}

func (i *typesImpl) DropDescriptor(_ context.Context, handle Descriptor) {
	i.host.FilesystemManager().Remove(handle)
}

func (i *typesImpl) GetType(ctx context.Context, this Descriptor) witgo.Result[DescriptorType, ErrorCode] {
	statResult := i.Stat(ctx, this)
	if statResult.Err != nil {
		return witgo.Err[DescriptorType, ErrorCode](*statResult.Err)
	}
	return witgo.Ok[DescriptorType, ErrorCode](statResult.Ok.Type)
}

func (i *typesImpl) DropDirectoryEntryStream(_ context.Context, handle DirectoryEntryStream) {
	i.host.DirectoryEntryStreamManager().Remove(handle)
}

func (i *typesImpl) ReadDirectoryEntry(_ context.Context, this DirectoryEntryStream) witgo.Result[witgo.Option[DirectoryEntry], ErrorCode] {
	stream, ok := i.host.DirectoryEntryStreamManager().Get(this)
	if !ok || stream == nil {
		return witgo.Err[witgo.Option[DirectoryEntry], ErrorCode](ErrorCodeBadDescriptor)
	}

	// 通过 Next() 持锁递增 Index，避免并发读同一流时重复/越界
	entry, ok := stream.Next()
	if !ok {
		return witgo.Ok[witgo.Option[DirectoryEntry], ErrorCode](witgo.None[DirectoryEntry]())
	}

	// entry.Info() 对断裂 symlink / 并发 unlink 会失败并中断整个 directory-entry-stream。
	// DirEntry.Type() 来自 d_type，失败时再回退 Info，最后 Unknown，保证还能返回 Name。
	mode := entry.Type()
	if mode == 0 {
		if info, err := entry.Info(); err == nil {
			mode = info.Mode()
		}
	}
	dirEntry := DirectoryEntry{
		Type: goModeToDescriptorType(mode),
		Name: entry.Name(),
	}

	return witgo.Ok[witgo.Option[DirectoryEntry], ErrorCode](witgo.Some(dirEntry))
}

func (i *typesImpl) FilesystemErrorCode(ctx context.Context, err WasiError) witgo.Option[ErrorCode] {
	// 这个函数用于将 wasi:io/error 向 wasi:filesystem/error-code "向下转型"
	// 我们需要一种方法来存储原始的 os/syscall 错误
	if e, ok := i.host.ErrorManager().Get(err); ok && e != nil {
		// 检查 e 是否是我们可以映射的错误类型
		if code := mapOsError(e); code != ErrorCodeUnsupported {
			return witgo.Some(code)
		}
	}
	return witgo.None[ErrorCode]()
}

func (i *typesImpl) ReadViaStream(_ context.Context, this Descriptor, offset Filesize) witgo.Result[InputStream, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil {
		return witgo.Err[InputStream, ErrorCode](ErrorCodeBadDescriptor)
	}
	off, code := filesizeToInt64(offset)
	if code != 0 {
		return witgo.Err[InputStream, ErrorCode](code)
	}
	limit := int64(math.MaxInt64)
	if off > 0 {
		limit = math.MaxInt64 - off
	}
	// 常规文件用真实剩余长度，skip 才能在 EOF 停下。管道和设备没有稳定长度，保持占位值。
	if st, statErr := d.Stat(); statErr == nil && st.Mode().IsRegular() {
		sz := st.Size()
		if sz <= off {
			limit = 0
		} else {
			limit = sz - off
		}
	}
	dup, err := d.DupFile()
	if errors.Is(err, os.ErrClosed) {
		return witgo.Err[InputStream, ErrorCode](ErrorCodeBadDescriptor)
	}
	if err == nil {
		// 流拥有 dup 出的 fd，ReadAt 不共享目录偏移。descriptor 先 drop 只关原 fd。
		reader := io.NewSectionReader(dup, off, limit)
		stream := &manager_io.Stream{Reader: reader, Seeker: reader, Closer: dup}
		return witgo.Ok[InputStream, ErrorCode](i.host.StreamManager().Add(stream))
	}
	// 无 dup 的平台不能直接失败，也不能把共享 *os.File 交给 SectionReader。
	sec := &descriptorSection{d: d, base: off, limit: limit}
	stream := &manager_io.Stream{Reader: sec, Seeker: sec}
	return witgo.Ok[InputStream, ErrorCode](i.host.StreamManager().Add(stream))
}

func (i *typesImpl) WriteViaStream(_ context.Context, this Descriptor, offset Filesize) witgo.Result[OutputStream, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil {
		return witgo.Err[OutputStream, ErrorCode](ErrorCodeBadDescriptor)
	}
	off, code := filesizeToInt64(offset)
	if code != 0 {
		return witgo.Err[OutputStream, ErrorCode](code)
	}
	writer := &sectionWriter{d: d, offset: off}
	return witgo.Ok[OutputStream, ErrorCode](i.host.StreamManager().Add(&manager_io.Stream{Writer: writer}))
}

func (i *typesImpl) AppendViaStream(_ context.Context, this Descriptor) witgo.Result[OutputStream, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil {
		return witgo.Err[OutputStream, ErrorCode](ErrorCodeBadDescriptor)
	}
	stream := &manager_io.Stream{
		Writer:  &appendWriter{d: d},
		Flusher: &OsFileFlusher{w: descriptorSyncer{d: d}},
	}
	return witgo.Ok[OutputStream, ErrorCode](i.host.StreamManager().Add(stream))
}

func (i *typesImpl) Advise(_ context.Context, this Descriptor, offset Filesize, length Filesize, advice Advice) witgo.Result[witgo.Unit, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeBadDescriptor)
	}
	if _, code := filesizeToInt64(offset); code != 0 {
		return witgo.Err[witgo.Unit, ErrorCode](code)
	}
	if _, code := filesizeToInt64(length); code != 0 {
		return witgo.Err[witgo.Unit, ErrorCode](code)
	}
	// 修改原因：advise 内部拿 SyscallConn，必须与 Close 互斥。
	err := d.DoR(func(f *os.File) error { return adviseFile(f, offset, length, advice) })
	if err != nil {
		if errors.Is(err, errAdviseUnsupported) {
			return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeUnsupported)
		}
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
}

func (i *typesImpl) SyncData(ctx context.Context, this Descriptor) witgo.Result[witgo.Unit, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeBadDescriptor)
	}
	if err := d.DoR(func(f *os.File) error { return syncDataFile(f) }); err != nil {
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
}

func (i *typesImpl) SetSize(ctx context.Context, this Descriptor, size Filesize) witgo.Result[witgo.Unit, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeBadDescriptor)
	}
	sz, code := filesizeToInt64(size)
	if code != 0 {
		return witgo.Err[witgo.Unit, ErrorCode](code)
	}
	if err := d.Truncate(sz); err != nil {
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
}

func (i *typesImpl) SetTimes(_ context.Context, this Descriptor, data_access_timestamp NewTimestamp, data_modification_timestamp NewTimestamp) witgo.Result[witgo.Unit, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeBadDescriptor)
	}
	// 修改原因：取路径和 futimes 必须在同一次锁内完成，避免中间 Close。
	err := d.DoR(func(f *os.File) error {
		atime, mtime, aerr := applyNewTimestamps(f.Name(), PathFlags{SymlinkFollow: true}, data_access_timestamp, data_modification_timestamp)
		if aerr != nil {
			return aerr
		}
		return chtimesFile(f, atime, mtime)
	})
	if err != nil {
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
}

func (i *typesImpl) Read(ctx context.Context, this Descriptor, length Filesize, offset Filesize) witgo.Result[witgo.Tuple[[]byte, bool], ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil {
		return witgo.Err[witgo.Tuple[[]byte, bool], ErrorCode](ErrorCodeBadDescriptor)
	}
	off, code := filesizeToInt64(offset)
	if code != 0 {
		return witgo.Err[witgo.Tuple[[]byte, bool], ErrorCode](code)
	}
	ncap := capReadLen(uint64(length))
	if ncap == 0 {
		return witgo.Ok[witgo.Tuple[[]byte, bool], ErrorCode](witgo.Tuple[[]byte, bool]{F0: []byte{}, F1: false})
	}
	buf := make([]byte, ncap)
	n, err := d.ReadAt(buf, off)
	if err != nil && !errors.Is(err, io.EOF) {
		return witgo.Err[witgo.Tuple[[]byte, bool], ErrorCode](mapOsError(err))
	}
	return witgo.Ok[witgo.Tuple[[]byte, bool], ErrorCode](witgo.Tuple[[]byte, bool]{F0: buf[:n], F1: errors.Is(err, io.EOF)})
}

func (i *typesImpl) Write(ctx context.Context, this Descriptor, buffer []byte, offset Filesize) witgo.Result[Filesize, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil {
		return witgo.Err[Filesize, ErrorCode](ErrorCodeBadDescriptor)
	}
	off, code := filesizeToInt64(offset)
	if code != 0 {
		return witgo.Err[Filesize, ErrorCode](code)
	}
	// WriteAt 允许短写。原先把已写入前缀和错误一起丢弃，guest 会重写整段导致重复。
	written := 0
	for written < len(buffer) {
		if off > math.MaxInt64-int64(written) {
			if written > 0 {
				return witgo.Ok[Filesize, ErrorCode](Filesize(written))
			}
			return witgo.Err[Filesize, ErrorCode](ErrorCodeOverflow)
		}
		n, err := d.WriteAt(buffer[written:], off+int64(written))
		if n > 0 {
			written += n
		}
		if err != nil {
			if written > 0 {
				return witgo.Ok[Filesize, ErrorCode](Filesize(written))
			}
			return witgo.Err[Filesize, ErrorCode](mapOsError(err))
		}
		if n == 0 {
			break
		}
	}
	return witgo.Ok[Filesize, ErrorCode](Filesize(written))
}

func (i *typesImpl) ReadDirectory(ctx context.Context, this Descriptor) witgo.Result[DirectoryEntryStream, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil {
		return witgo.Err[DirectoryEntryStream, ErrorCode](ErrorCodeBadDescriptor)
	}
	entries, err := d.ReadDir(-1)
	if err != nil {
		return witgo.Err[DirectoryEntryStream, ErrorCode](mapOsError(err))
	}
	filtered := make([]fs.DirEntry, 0, len(entries))
	for _, e := range entries {
		if name := e.Name(); name == "." || name == ".." {
			continue
		}
		filtered = append(filtered, e)
	}
	handle := i.host.DirectoryEntryStreamManager().Add(&filesystem.DirectoryEntryStreamState{Entries: filtered})
	return witgo.Ok[DirectoryEntryStream, ErrorCode](handle)
}

func (i *typesImpl) Sync(ctx context.Context, this Descriptor) witgo.Result[witgo.Unit, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeBadDescriptor)
	}
	if err := d.Sync(); err != nil {
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
}

func (i *typesImpl) CreateDirectoryAt(ctx context.Context, this Descriptor, path string) witgo.Result[witgo.Unit, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeBadDescriptor)
	}
	base, code := descriptorBase(d)
	if code != 0 {
		return witgo.Err[witgo.Unit, ErrorCode](code)
	}
	fullPath, code := resolveSandboxedPathEx(base, path, false)
	if code != 0 {
		return witgo.Err[witgo.Unit, ErrorCode](code)
	}
	if err := os.Mkdir(fullPath, 0o755); err != nil {
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
}

func (i *typesImpl) Stat(ctx context.Context, this Descriptor) witgo.Result[DescriptorStat, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil {
		return witgo.Err[DescriptorStat, ErrorCode](ErrorCodeBadDescriptor)
	}
	info, err := d.Stat()
	if err != nil {
		return witgo.Err[DescriptorStat, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[DescriptorStat, ErrorCode](goFileInfoToDescriptorStat(info))
}

func (i *typesImpl) StatAt(ctx context.Context, this Descriptor, pathFlags PathFlags, path string) witgo.Result[DescriptorStat, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil {
		return witgo.Err[DescriptorStat, ErrorCode](ErrorCodeBadDescriptor)
	}
	base, code := descriptorBase(d)
	if code != 0 {
		return witgo.Err[DescriptorStat, ErrorCode](code)
	}
	fullPath, code := resolveSandboxedPathEx(base, path, pathFlags.SymlinkFollow)
	if code != 0 {
		return witgo.Err[DescriptorStat, ErrorCode](code)
	}
	var info fs.FileInfo
	var err error
	if pathFlags.SymlinkFollow {
		info, err = os.Stat(fullPath)
	} else {
		info, err = os.Lstat(fullPath)
	}
	if err != nil {
		return witgo.Err[DescriptorStat, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[DescriptorStat, ErrorCode](goFileInfoToDescriptorStat(info))
}

func (i *typesImpl) SetTimesAt(ctx context.Context, this Descriptor, pathFlags PathFlags, path string, data_access_timestamp NewTimestamp, data_modification_timestamp NewTimestamp) witgo.Result[witgo.Unit, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeBadDescriptor)
	}
	base, code := descriptorBase(d)
	if code != 0 {
		return witgo.Err[witgo.Unit, ErrorCode](code)
	}
	fullPath, code := resolveSandboxedPathEx(base, path, pathFlags.SymlinkFollow)
	if code != 0 {
		return witgo.Err[witgo.Unit, ErrorCode](code)
	}
	atime, mtime, err := applyNewTimestamps(fullPath, pathFlags, data_access_timestamp, data_modification_timestamp)
	if err != nil {
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
	}
	if err := chtimesPath(fullPath, atime, mtime, pathFlags.SymlinkFollow); err != nil {
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
}

func (i *typesImpl) LinkAt(ctx context.Context, this Descriptor, oldPathFlags PathFlags, oldPath string, newDescriptor Descriptor, newPath string) witgo.Result[witgo.Unit, ErrorCode] {
	oldDir, ok := i.host.FilesystemManager().Get(this)
	if !ok || oldDir == nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeBadDescriptor)
	}
	newDir, ok := i.host.FilesystemManager().Get(newDescriptor)
	if !ok || newDir == nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeBadDescriptor)
	}
	oldBase, code := descriptorBase(oldDir)
	if code != 0 {
		return witgo.Err[witgo.Unit, ErrorCode](code)
	}
	newBase, code := descriptorBase(newDir)
	if code != 0 {
		return witgo.Err[witgo.Unit, ErrorCode](code)
	}
	oldFullPath, code := resolveSandboxedPathEx(oldBase, oldPath, oldPathFlags.SymlinkFollow)
	if code != 0 {
		return witgo.Err[witgo.Unit, ErrorCode](code)
	}
	newFullPath, code := resolveSandboxedPathEx(newBase, newPath, false)
	if code != 0 {
		return witgo.Err[witgo.Unit, ErrorCode](code)
	}
	if err := os.Link(oldFullPath, newFullPath); err != nil {
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
}

func (i *typesImpl) OpenAt(ctx context.Context, this Descriptor, pathFlags PathFlags, path string, openFlags OpenFlags, flags DescriptorFlags) witgo.Result[Descriptor, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil {
		return witgo.Err[Descriptor, ErrorCode](ErrorCodeBadDescriptor)
	}

	// WASI open-at 要求 flags 至少包含 read 或 write。两者都没有时不能打开。
	if !flags.Read && !flags.Write {
		return witgo.Err[Descriptor, ErrorCode](ErrorCodeInvalid)
	}

	base, code := descriptorBase(d)
	if code != 0 {
		return witgo.Err[Descriptor, ErrorCode](code)
	}
	fullPath, code := resolveSandboxedPathEx(base, path, pathFlags.SymlinkFollow)
	if code != 0 {
		return witgo.Err[Descriptor, ErrorCode](code)
	}
	createdDir := false
	// Mkdir 成功后，后面的 OpenFile/Stat 失败必须删掉刚建的空目录。
	fail := func(c ErrorCode) witgo.Result[Descriptor, ErrorCode] {
		if createdDir {
			os.Remove(fullPath)
		}
		return witgo.Err[Descriptor, ErrorCode](c)
	}
	if openFlags.Create && openFlags.Directory {
		if err := os.Mkdir(fullPath, 0o755); err != nil {
			if os.IsExist(err) {
				if openFlags.Exclusive {
					return witgo.Err[Descriptor, ErrorCode](ErrorCodeExist)
				}
			} else {
				return witgo.Err[Descriptor, ErrorCode](mapOsError(err))
			}
		} else {
			createdDir = true
		}
	}
	if !pathFlags.SymlinkFollow {
		if info, err := os.Lstat(fullPath); err == nil && info.Mode()&os.ModeSymlink != 0 {
			return fail(ErrorCodeLoop)
		}
	}
	var osFlags int
	if flags.Read && flags.Write {
		osFlags |= os.O_RDWR
	} else if flags.Write {
		osFlags |= os.O_WRONLY
	} else {
		osFlags |= os.O_RDONLY
	}
	if openFlags.Create && !openFlags.Directory {
		osFlags |= os.O_CREATE
	}
	if openFlags.Exclusive && !createdDir {
		osFlags |= os.O_EXCL
	}
	if openFlags.Truncate {
		osFlags |= os.O_TRUNC
	}
	osFlags |= extraOpenFlags(pathFlags, openFlags)
	file, err := os.OpenFile(fullPath, osFlags, 0o644)
	if err != nil {
		return fail(mapOsError(err))
	}
	if openFlags.Directory {
		info, statErr := file.Stat()
		if statErr != nil {
			_ = file.Close()
			return fail(mapOsError(statErr))
		}
		if !info.IsDir() {
			_ = file.Close()
			return fail(ErrorCodeNotDirectory)
		}
	}
	handle := i.host.FilesystemManager().Add(&filesystem.Descriptor{File: file, Path: path})
	return witgo.Ok[Descriptor, ErrorCode](handle)
}

func (i *typesImpl) ReadlinkAt(ctx context.Context, this Descriptor, path string) witgo.Result[string, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil {
		return witgo.Err[string, ErrorCode](ErrorCodeBadDescriptor)
	}
	base, code := descriptorBase(d)
	if code != 0 {
		return witgo.Err[string, ErrorCode](code)
	}
	fullPath, code := resolveSandboxedPathEx(base, path, false)
	if code != 0 {
		return witgo.Err[string, ErrorCode](code)
	}
	target, err := os.Readlink(fullPath)
	if err != nil {
		return witgo.Err[string, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[string, ErrorCode](target)
}

func (i *typesImpl) RemoveDirectoryAt(ctx context.Context, this Descriptor, path string) witgo.Result[witgo.Unit, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeBadDescriptor)
	}
	base, code := descriptorBase(d)
	if code != 0 {
		return witgo.Err[witgo.Unit, ErrorCode](code)
	}
	fullPath, code := resolveSandboxedPathEx(base, path, false)
	if code != 0 {
		return witgo.Err[witgo.Unit, ErrorCode](code)
	}
	info, statErr := os.Lstat(fullPath)
	if statErr != nil {
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(statErr))
	}
	if !info.IsDir() {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeNotDirectory)
	}
	if err := os.Remove(fullPath); err != nil {
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
}

func (i *typesImpl) RenameAt(ctx context.Context, this Descriptor, oldPath string, newDescriptor Descriptor, newPath string) witgo.Result[witgo.Unit, ErrorCode] {
	oldDir, ok := i.host.FilesystemManager().Get(this)
	if !ok || oldDir == nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeBadDescriptor)
	}
	newDir, ok := i.host.FilesystemManager().Get(newDescriptor)
	if !ok || newDir == nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeBadDescriptor)
	}
	oldBase, code := descriptorBase(oldDir)
	if code != 0 {
		return witgo.Err[witgo.Unit, ErrorCode](code)
	}
	newBase, code := descriptorBase(newDir)
	if code != 0 {
		return witgo.Err[witgo.Unit, ErrorCode](code)
	}
	oldFullPath, code := resolveSandboxedPathEx(oldBase, oldPath, false)
	if code != 0 {
		return witgo.Err[witgo.Unit, ErrorCode](code)
	}
	newFullPath, code := resolveSandboxedPathEx(newBase, newPath, false)
	if code != 0 {
		return witgo.Err[witgo.Unit, ErrorCode](code)
	}
	if err := os.Rename(oldFullPath, newFullPath); err != nil {
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
}

func (i *typesImpl) SymlinkAt(ctx context.Context, this Descriptor, oldPath string, newPath string) witgo.Result[witgo.Unit, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeBadDescriptor)
	}
	base, code := descriptorBase(d)
	if code != 0 {
		return witgo.Err[witgo.Unit, ErrorCode](code)
	}
	newFullPath, code := resolveSandboxedPathEx(base, newPath, false)
	if code != 0 {
		return witgo.Err[witgo.Unit, ErrorCode](code)
	}
	if err := os.Symlink(oldPath, newFullPath); err != nil {
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
}

func (i *typesImpl) UnlinkFileAt(ctx context.Context, this Descriptor, path string) witgo.Result[witgo.Unit, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeBadDescriptor)
	}
	base, code := descriptorBase(d)
	if code != 0 {
		return witgo.Err[witgo.Unit, ErrorCode](code)
	}
	fullPath, code := resolveSandboxedPathEx(base, path, false)
	if code != 0 {
		return witgo.Err[witgo.Unit, ErrorCode](code)
	}
	info, statErr := os.Lstat(fullPath)
	if statErr != nil {
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(statErr))
	}
	if info.IsDir() {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeIsDirectory)
	}
	if err := os.Remove(fullPath); err != nil {
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
}

func (i *typesImpl) IsSameObject(ctx context.Context, this Descriptor, other Descriptor) bool {
	d1, ok1 := i.host.FilesystemManager().Get(this)
	d2, ok2 := i.host.FilesystemManager().Get(other)
	if !ok1 || !ok2 || d1 == nil || d2 == nil {
		return false
	}
	info1, err1 := d1.Stat()
	info2, err2 := d2.Stat()
	if err1 != nil || err2 != nil {
		return false
	}
	return os.SameFile(info1, info2)
}

func (i *typesImpl) MetadataHash(ctx context.Context, this Descriptor) witgo.Result[MetadataHashValue, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil {
		return witgo.Err[MetadataHashValue, ErrorCode](ErrorCodeBadDescriptor)
	}
	var info fs.FileInfo
	var name string
	// 修改原因：Stat 和 Name 分两次加锁，中间 Close 会让哈希用到已关闭文件的路径。
	err := d.DoR(func(f *os.File) error {
		var statErr error
		info, statErr = f.Stat()
		if statErr != nil {
			return statErr
		}
		name = f.Name()
		return nil
	})
	if err != nil {
		return witgo.Err[MetadataHashValue, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[MetadataHashValue, ErrorCode](metadataHashFromInfo(info, name))
}

func (i *typesImpl) MetadataHashAt(ctx context.Context, this Descriptor, pathFlags PathFlags, path string) witgo.Result[MetadataHashValue, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil {
		return witgo.Err[MetadataHashValue, ErrorCode](ErrorCodeBadDescriptor)
	}
	base, code := descriptorBase(d)
	if code != 0 {
		return witgo.Err[MetadataHashValue, ErrorCode](code)
	}
	fullPath, code := resolveSandboxedPathEx(base, path, pathFlags.SymlinkFollow)
	if code != 0 {
		return witgo.Err[MetadataHashValue, ErrorCode](code)
	}
	var info fs.FileInfo
	var err error
	if pathFlags.SymlinkFollow {
		info, err = os.Stat(fullPath)
	} else {
		info, err = os.Lstat(fullPath)
	}
	if err != nil {
		return witgo.Err[MetadataHashValue, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[MetadataHashValue, ErrorCode](metadataHashFromInfo(info, fullPath))
}

type OsFileFlusher struct {
	w interface{ Sync() error }
}

func (f *OsFileFlusher) Flush() error {
	return f.w.Sync()
}

type appendWriter struct {
	d *filesystem.Descriptor
}

func (w *appendWriter) Write(p []byte) (int, error) {
	if w == nil || w.d == nil {
		return 0, os.ErrInvalid
	}
	return w.d.AppendWrite(p)
}

type sectionWriter struct {
	mu     sync.Mutex // WriteViaStream 可能并发写，offset 会乱
	d      *filesystem.Descriptor
	offset int64
}

func (s *sectionWriter) Write(p []byte) (n int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.d == nil {
		return 0, os.ErrInvalid
	}
	// 走 Descriptor.WriteAt，和 set-size/append 共用同一把锁。
	n, err = s.d.WriteAt(p, s.offset)
	if n > 0 {
		if s.offset > math.MaxInt64-int64(n) {
			s.offset = math.MaxInt64
		} else {
			s.offset += int64(n)
		}
	}
	return
}

// descriptorBase 在锁内取目录路径。已关闭的描述符不再做路径运算
func descriptorBase(d *filesystem.Descriptor) (string, ErrorCode) {
	if d == nil {
		return "", ErrorCodeBadDescriptor
	}
	var name string
	var isDir bool
	err := d.DoR(func(f *os.File) error {
		info, statErr := f.Stat()
		if statErr != nil {
			return statErr
		}
		isDir = info.IsDir()
		name = f.Name()
		return nil
	})
	if err != nil {
		return "", ErrorCodeBadDescriptor
	}
	// WASI *-at 要求基描述符是目录。只拿 Name 时，普通文件路径会被当成目录拼接。
	// 不能改用 Descriptor.Path：OpenAt 存的是 guest 相对路径，不是宿主绝对路径。
	if !isDir {
		return "", ErrorCodeNotDirectory
	}
	return name, 0
}

// descriptorSyncer 让 flush 走 Descriptor.Sync。
// 保存 *os.File 后 Sync 会与 Close 并发。
type descriptorSyncer struct{ d *filesystem.Descriptor }

func (s descriptorSyncer) Sync() error {
	if s.d == nil {
		return os.ErrInvalid
	}
	return s.d.Sync()
}

// descriptorSection 在无法 dup 的平台按偏移读。每次 ReadAt 都与 Close 互斥。
type descriptorSection struct {
	d     *filesystem.Descriptor
	base  int64
	off   int64
	limit int64
	mu    sync.Mutex
}

func (s *descriptorSection) Size() int64 { return s.limit }

func (s *descriptorSection) Seek(offset int64, whence int) (int64, error) {
	if s == nil {
		return 0, os.ErrInvalid
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	var n int64
	switch whence {
	case io.SeekStart:
		n = offset
	case io.SeekCurrent:
		n = s.off + offset
	case io.SeekEnd:
		n = s.limit + offset
	default:
		return 0, os.ErrInvalid
	}
	if n < 0 {
		return 0, os.ErrInvalid
	}
	s.off = n
	return n, nil
}

func (s *descriptorSection) Read(p []byte) (int, error) {
	if s == nil || s.d == nil {
		return 0, os.ErrInvalid
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.off >= s.limit {
		return 0, io.EOF
	}
	remain := s.limit - s.off
	if int64(len(p)) > remain {
		p = p[:remain]
	}
	// base+off 溢出后变成负偏移，ReadAt 会读到文件头或直接报错。
	if s.base < 0 || s.off > math.MaxInt64-s.base {
		return 0, os.ErrInvalid
	}
	n, err := s.d.ReadAt(p, s.base+s.off)
	s.off += int64(n)
	if n == 0 && errors.Is(err, os.ErrClosed) {
		return 0, io.EOF
	}
	return n, err
}
