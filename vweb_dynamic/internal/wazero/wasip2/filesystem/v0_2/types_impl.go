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

// resolveSandboxedPath 保持原签名，默认跟随已存在的符号链接。
// 未改调用点的旧代码仍能编译。
func resolveSandboxedPath(baseDir, guestPath string) (string, ErrorCode) {
	return resolveSandboxedPathEx(baseDir, guestPath, true)
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

func (i *typesImpl) DropDescriptor(_ context.Context, handle Descriptor) {
	i.host.FilesystemManager().Remove(handle)
}

func (i *typesImpl) ReadViaStream(_ context.Context, this Descriptor, offset Filesize) witgo.Result[InputStream, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil || d.File == nil {
		return witgo.Err[InputStream, ErrorCode](ErrorCodeBadDescriptor)
	}
	off, code := filesizeToInt64(offset)
	if code != 0 {
		return witgo.Err[InputStream, ErrorCode](code)
	}
	// 本方法无 length 参数；单次读上限由 wasi:io streams.read 的 capReadLen 保证
	// MaxInt64-off 避免 SectionReader 内部 offset+limit 溢出
	limit := int64(math.MaxInt64)
	if off > 0 {
		limit = math.MaxInt64 - off
	}
	reader := io.NewSectionReader(d.File, off, limit)
	stream := &manager_io.Stream{Reader: reader, Seeker: reader}
	handle := i.host.StreamManager().Add(stream)
	return witgo.Ok[InputStream, ErrorCode](handle)
}

func (i *typesImpl) WriteViaStream(_ context.Context, this Descriptor, offset Filesize) witgo.Result[OutputStream, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil || d.File == nil {
		return witgo.Err[OutputStream, ErrorCode](ErrorCodeBadDescriptor)
	}
	off, code := filesizeToInt64(offset)
	if code != 0 {
		return witgo.Err[OutputStream, ErrorCode](code)
	}
	writer := &sectionWriter{w: d.File, offset: off}
	stream := &manager_io.Stream{Writer: writer}
	handle := i.host.StreamManager().Add(stream)
	return witgo.Ok[OutputStream, ErrorCode](handle)
}

func (i *typesImpl) AppendViaStream(_ context.Context, this Descriptor) witgo.Result[OutputStream, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil || d.File == nil {
		return witgo.Err[OutputStream, ErrorCode](ErrorCodeBadDescriptor)
	}
	stream := &manager_io.Stream{Writer: d.File, Flusher: &OsFileFlusher{w: d.File}}
	handle := i.host.StreamManager().Add(stream)
	return witgo.Ok[OutputStream, ErrorCode](handle)
}

func (i *typesImpl) Advise(_ context.Context, this Descriptor, offset Filesize, length Filesize, advice Advice) witgo.Result[witgo.Unit, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil || d.File == nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeBadDescriptor)
	}
	if _, code := filesizeToInt64(offset); code != 0 {
		return witgo.Err[witgo.Unit, ErrorCode](code)
	}
	if _, code := filesizeToInt64(length); code != 0 {
		return witgo.Err[witgo.Unit, ErrorCode](code)
	}
	// 原先一律 Unsupported；Linux 可用 Fadvise，其它平台仍返回 Unsupported
	if err := adviseFile(d.File, offset, length, advice); err != nil {
		if errors.Is(err, errAdviseUnsupported) {
			return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeUnsupported)
		}
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
}

func (i *typesImpl) SyncData(ctx context.Context, this Descriptor) witgo.Result[witgo.Unit, ErrorCode] {
	return i.Sync(ctx, this)
}

func (i *typesImpl) GetType(ctx context.Context, this Descriptor) witgo.Result[DescriptorType, ErrorCode] {
	statResult := i.Stat(ctx, this)
	if statResult.Err != nil {
		return witgo.Err[DescriptorType, ErrorCode](*statResult.Err)
	}
	return witgo.Ok[DescriptorType, ErrorCode](statResult.Ok.Type)
}

func (i *typesImpl) SetSize(ctx context.Context, this Descriptor, size Filesize) witgo.Result[witgo.Unit, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil || d.File == nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeBadDescriptor)
	}
	sz, code := filesizeToInt64(size)
	if code != 0 {
		return witgo.Err[witgo.Unit, ErrorCode](code)
	}
	err := d.File.Truncate(sz)
	if err != nil {
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
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

func (i *typesImpl) SetTimes(_ context.Context, this Descriptor, data_access_timestamp NewTimestamp, data_modification_timestamp NewTimestamp) witgo.Result[witgo.Unit, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil || d.File == nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeBadDescriptor)
	}

	path := d.File.Name()
	atime, mtime, err := applyNewTimestamps(path, PathFlags{SymlinkFollow: true}, data_access_timestamp, data_modification_timestamp)
	if err != nil {
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
	}

	//  chtimesFile 走 fd，unlink 后仍有效。os.Chtimes 走路径，unlink 后会失败。
	// 通过平台 chtimesFile 对已打开的 fd 改 atime/mtime（Unix futimes / Windows SetFileTime）。
	if err := chtimesFile(d.File, atime, mtime); err != nil {
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
}

func (i *typesImpl) Read(ctx context.Context, this Descriptor, length Filesize, offset Filesize) witgo.Result[witgo.Tuple[[]byte, bool], ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil || d.File == nil {
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
	n, err := d.File.ReadAt(buf, off)
	endOfFile := err == io.EOF
	if err != nil && err != io.EOF {
		return witgo.Err[witgo.Tuple[[]byte, bool], ErrorCode](mapOsError(err))
	}
	// 即使请求 length 更大，也只返回本轮读到的字节（WASI 允许短读）
	result := witgo.Tuple[[]byte, bool]{F0: buf[:n], F1: endOfFile}
	return witgo.Ok[witgo.Tuple[[]byte, bool], ErrorCode](result)
}

func (i *typesImpl) Write(ctx context.Context, this Descriptor, buffer []byte, offset Filesize) witgo.Result[Filesize, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil || d.File == nil {
		return witgo.Err[Filesize, ErrorCode](ErrorCodeBadDescriptor)
	}
	off, code := filesizeToInt64(offset)
	if code != 0 {
		return witgo.Err[Filesize, ErrorCode](code)
	}
	n, err := d.File.WriteAt(buffer, off)
	if err != nil {
		return witgo.Err[Filesize, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[Filesize, ErrorCode](Filesize(n))
}

func (i *typesImpl) ReadDirectory(ctx context.Context, this Descriptor) witgo.Result[DirectoryEntryStream, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil || d.File == nil {
		return witgo.Err[DirectoryEntryStream, ErrorCode](ErrorCodeBadDescriptor)
	}

	entries, err := d.File.ReadDir(-1) // 用 fd 列举，避免 unlink/换挂后 Name() 指到错误路径
	if err != nil {
		return witgo.Err[DirectoryEntryStream, ErrorCode](mapOsError(err))
	}

	// WASI 规定 directory-entry-stream 不得返回 "." / ".."
	filtered := make([]fs.DirEntry, 0, len(entries))
	for _, e := range entries {
		name := e.Name()
		if name == "." || name == ".." {
			continue
		}
		filtered = append(filtered, e)
	}

	streamState := &filesystem.DirectoryEntryStreamState{
		Entries: filtered,
		Index:   0,
	}
	handle := i.host.DirectoryEntryStreamManager().Add(streamState)
	return witgo.Ok[DirectoryEntryStream, ErrorCode](handle)
}

func (i *typesImpl) Sync(ctx context.Context, this Descriptor) witgo.Result[witgo.Unit, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil || d.File == nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeBadDescriptor)
	}
	err := d.File.Sync()
	if err != nil {
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
}

func (i *typesImpl) CreateDirectoryAt(ctx context.Context, this Descriptor, path string) witgo.Result[witgo.Unit, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil || d.File == nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeBadDescriptor)
	}

	fullPath, code := resolveSandboxedPathEx(d.File.Name(), path, false)
	if code != 0 {
		return witgo.Err[witgo.Unit, ErrorCode](code)
	}

	err := os.Mkdir(fullPath, 0o755)
	if err != nil {
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
}

func (i *typesImpl) Stat(ctx context.Context, this Descriptor) witgo.Result[DescriptorStat, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil || d.File == nil {
		return witgo.Err[DescriptorStat, ErrorCode](ErrorCodeBadDescriptor)
	}
	info, err := d.File.Stat()
	if err != nil {
		return witgo.Err[DescriptorStat, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[DescriptorStat, ErrorCode](goFileInfoToDescriptorStat(info))
}

func (i *typesImpl) StatAt(ctx context.Context, this Descriptor, pathFlags PathFlags, path string) witgo.Result[DescriptorStat, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil || d.File == nil {
		return witgo.Err[DescriptorStat, ErrorCode](ErrorCodeBadDescriptor)
	}

	fullPath, code := resolveSandboxedPathEx(d.File.Name(), path, pathFlags.SymlinkFollow)
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
	if !ok || d == nil || d.File == nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeBadDescriptor)
	}

	fullPath, code := resolveSandboxedPathEx(d.File.Name(), path, pathFlags.SymlinkFollow)
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
	if !ok || oldDir == nil || oldDir.File == nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeBadDescriptor)
	}
	newDir, ok := i.host.FilesystemManager().Get(newDescriptor)
	if !ok || newDir == nil || newDir.File == nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeBadDescriptor)
	}

	oldFullPath, code := resolveSandboxedPathEx(oldDir.File.Name(), oldPath, oldPathFlags.SymlinkFollow)
	if code != 0 {
		return witgo.Err[witgo.Unit, ErrorCode](code)
	}
	newFullPath, code := resolveSandboxedPathEx(newDir.File.Name(), newPath, false)
	if code != 0 {
		return witgo.Err[witgo.Unit, ErrorCode](code)
	}

	err := os.Link(oldFullPath, newFullPath)
	if err != nil {
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
}

func (i *typesImpl) OpenAt(ctx context.Context, this Descriptor, pathFlags PathFlags, path string, openFlags OpenFlags, flags DescriptorFlags) witgo.Result[Descriptor, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil || d.File == nil {
		return witgo.Err[Descriptor, ErrorCode](ErrorCodeBadDescriptor)
	}

	// open 是否跟随符号链接由 path-flags.symlink-follow 决定，不能一律 follow
	fullPath, code := resolveSandboxedPathEx(d.File.Name(), path, pathFlags.SymlinkFollow)
	if code != 0 {
		return witgo.Err[Descriptor, ErrorCode](code)
	}

	// WASI create+directory 表示“若不存在则创建目录”；
	// os.OpenFile(..., O_CREATE) 会建成普通文件，不是目录。
	if openFlags.Create && openFlags.Directory {
		if err := os.Mkdir(fullPath, 0o755); err != nil {
			if os.IsExist(err) {
				// create+exclusive+directory 在已存在时应 Exist，不能忽略 IsExist。
				if openFlags.Exclusive {
					return witgo.Err[Descriptor, ErrorCode](ErrorCodeExist)
				}
			} else {
				return witgo.Err[Descriptor, ErrorCode](mapOsError(err))
			}
		}
	}

	// Windows/other 没有 O_NOFOLLOW；follow=false 时若最终分量是 symlink 应返回 loop
	if !pathFlags.SymlinkFollow {
		if info, err := os.Lstat(fullPath); err == nil && info.Mode()&os.ModeSymlink != 0 {
			return witgo.Err[Descriptor, ErrorCode](ErrorCodeLoop)
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
		// 目录已由 Mkdir 创建；再带 O_CREATE 在部分平台对目录不合法
		osFlags |= os.O_CREATE
	}
	if openFlags.Exclusive {
		osFlags |= os.O_EXCL
	}
	if openFlags.Truncate {
		osFlags |= os.O_TRUNC
	}
	// 补全 path-flags.symlink-follow=false 与 open-flags.directory（unix 上为 O_NOFOLLOW/O_DIRECTORY）
	osFlags |= extraOpenFlags(pathFlags, openFlags)

	file, err := os.OpenFile(fullPath, osFlags, 0o644)
	if err != nil {
		return witgo.Err[Descriptor, ErrorCode](mapOsError(err))
	}

	if openFlags.Directory {
		// 无 O_DIRECTORY 的平台（Windows/other）打开后必须校验是目录
		info, statErr := file.Stat()
		if statErr != nil {
			_ = file.Close()
			return witgo.Err[Descriptor, ErrorCode](mapOsError(statErr))
		}
		if !info.IsDir() {
			_ = file.Close()
			return witgo.Err[Descriptor, ErrorCode](ErrorCodeNotDirectory)
		}
	}

	newDesc := &filesystem.Descriptor{
		File: file,
		Path: path,
	}
	handle := i.host.FilesystemManager().Add(newDesc)
	return witgo.Ok[Descriptor, ErrorCode](handle)
}

func (i *typesImpl) ReadlinkAt(ctx context.Context, this Descriptor, path string) witgo.Result[string, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil || d.File == nil {
		return witgo.Err[string, ErrorCode](ErrorCodeBadDescriptor)
	}
	fullPath, code := resolveSandboxedPathEx(d.File.Name(), path, false)
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
	if !ok || d == nil || d.File == nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeBadDescriptor)
	}
	fullPath, code := resolveSandboxedPathEx(d.File.Name(), path, false)
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
	err := os.Remove(fullPath)
	if err != nil {
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
}

func (i *typesImpl) RenameAt(ctx context.Context, this Descriptor, oldPath string, newDescriptor Descriptor, newPath string) witgo.Result[witgo.Unit, ErrorCode] {
	oldDir, ok := i.host.FilesystemManager().Get(this)
	if !ok || oldDir == nil || oldDir.File == nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeBadDescriptor)
	}
	newDir, ok := i.host.FilesystemManager().Get(newDescriptor)
	if !ok || newDir == nil || newDir.File == nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeBadDescriptor)
	}

	oldFullPath, code := resolveSandboxedPathEx(oldDir.File.Name(), oldPath, false)
	if code != 0 {
		return witgo.Err[witgo.Unit, ErrorCode](code)
	}
	newFullPath, code := resolveSandboxedPathEx(newDir.File.Name(), newPath, false)
	if code != 0 {
		return witgo.Err[witgo.Unit, ErrorCode](code)
	}

	err := os.Rename(oldFullPath, newFullPath)
	if err != nil {
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
}

func (i *typesImpl) SymlinkAt(ctx context.Context, this Descriptor, oldPath string, newPath string) witgo.Result[witgo.Unit, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil || d.File == nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeBadDescriptor)
	}
	// newPath 是要创建的链接名，不得跟随；oldPath 是链接内容，按 WASI 原样写入
	newFullPath, code := resolveSandboxedPathEx(d.File.Name(), newPath, false)
	if code != 0 {
		return witgo.Err[witgo.Unit, ErrorCode](code)
	}

	err := os.Symlink(oldPath, newFullPath)
	if err != nil {
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
}

func (i *typesImpl) UnlinkFileAt(ctx context.Context, this Descriptor, path string) witgo.Result[witgo.Unit, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil || d.File == nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeBadDescriptor)
	}
	fullPath, code := resolveSandboxedPathEx(d.File.Name(), path, false)
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
	err := os.Remove(fullPath)
	if err != nil {
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
}

func (i *typesImpl) IsSameObject(ctx context.Context, this Descriptor, other Descriptor) bool {
	d1, ok1 := i.host.FilesystemManager().Get(this)
	d2, ok2 := i.host.FilesystemManager().Get(other)
	if !ok1 || !ok2 || d1 == nil || d2 == nil || d1.File == nil || d2.File == nil {
		return false
	}

	info1, err1 := d1.File.Stat()
	info2, err2 := d2.File.Stat()
	if err1 != nil || err2 != nil {
		return false
	}

	return os.SameFile(info1, info2)
}

func (i *typesImpl) MetadataHash(ctx context.Context, this Descriptor) witgo.Result[MetadataHashValue, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil || d.File == nil {
		return witgo.Err[MetadataHashValue, ErrorCode](ErrorCodeBadDescriptor)
	}
	info, err := d.File.Stat()
	if err != nil {
		return witgo.Err[MetadataHashValue, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[MetadataHashValue, ErrorCode](metadataHashFromInfo(info, d.File.Name()))
}

func (i *typesImpl) MetadataHashAt(ctx context.Context, this Descriptor, pathFlags PathFlags, path string) witgo.Result[MetadataHashValue, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil || d.File == nil {
		return witgo.Err[MetadataHashValue, ErrorCode](ErrorCodeBadDescriptor)
	}
	fullPath, code := resolveSandboxedPathEx(d.File.Name(), path, pathFlags.SymlinkFollow)
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

	info, err := entry.Info()
	if err != nil {
		return witgo.Err[witgo.Option[DirectoryEntry], ErrorCode](mapOsError(err))
	}

	dirEntry := DirectoryEntry{
		Type: goModeToDescriptorType(info.Mode()),
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

// --- Helper: sectionWriter for WriteViaStream ---
type sectionWriter struct {
	mu     sync.Mutex // WriteViaStream 可能并发写，offset 会乱
	w      *os.File
	offset int64
}

func (s *sectionWriter) Write(p []byte) (n int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	n, err = s.w.WriteAt(p, s.offset)
	s.offset += int64(n)
	return
}

type OsFileFlusher struct {
	w interface{ Sync() error }
}

func (f *OsFileFlusher) Flush() error {
	return f.w.Sync()
}
