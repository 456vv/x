//go:build unix

package v0_2

import (
	"context"
	"encoding/binary"
	"errors"
	"hash"
	"io/fs"
	"os"
	"syscall"
	"time"

	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"

	"golang.org/x/sys/unix"
)

func (i *typesImpl) GetFlags(ctx context.Context, this Descriptor) witgo.Result[DescriptorFlags, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil || d.File == nil {
		return witgo.Err[DescriptorFlags, ErrorCode](ErrorCodeBadDescriptor)
	}

	raw, err := d.File.SyscallConn()
	if err != nil {
		return witgo.Err[DescriptorFlags, ErrorCode](mapOsError(err))
	}
	var sysErr error
	var flags int
	ctrlErr := raw.Control(func(fd uintptr) {
		flags, sysErr = unix.FcntlInt(fd, unix.F_GETFL, 0)
	})
	if ctrlErr != nil {
		return witgo.Err[DescriptorFlags, ErrorCode](mapOsError(ctrlErr))
	}

	if sysErr != nil {
		return witgo.Err[DescriptorFlags, ErrorCode](mapOsError(sysErr))
	}

	var wasiFlags DescriptorFlags
	accmode := flags & unix.O_ACCMODE
	switch accmode {
	case unix.O_RDWR:
		wasiFlags.Read = true
		wasiFlags.Write = true
	case unix.O_WRONLY:
		wasiFlags.Write = true
	default:
		wasiFlags.Read = true
	}

	if flags&unix.O_DSYNC != 0 {
		wasiFlags.DataIntegritySync = true
	}
	if flags&unix.O_SYNC != 0 {
		wasiFlags.FileIntegritySync = true
		wasiFlags.RequestedWriteSync = true
	}

	return witgo.Ok[DescriptorFlags, ErrorCode](wasiFlags)
}

func timeToDatetime(ts syscall.Timespec) Datetime {
	return Datetime{
		Seconds:     uint64(ts.Sec),
		Nanoseconds: uint32(ts.Nsec),
	}
}

func goModeToDescriptorType(mode fs.FileMode) DescriptorType {
	switch {
	case mode.IsRegular():
		return DescriptorTypeRegularFile
	case mode.IsDir():
		return DescriptorTypeDirectory
	case mode&fs.ModeSymlink != 0:
		return DescriptorTypeSymbolicLink
	case mode&fs.ModeDevice != 0:
		if mode&fs.ModeCharDevice != 0 {
			return DescriptorTypeCharacterDevice
		}
		return DescriptorTypeBlockDevice
	case mode&fs.ModeNamedPipe != 0:
		return DescriptorTypeFifo
	case mode&fs.ModeSocket != 0:
		return DescriptorTypeSocket
	default:
		return DescriptorTypeUnknown
	}
}

func mapOsError(err error) ErrorCode {
	if err == nil {
		return 0
	}
	if errors.Is(err, fs.ErrPermission) {
		return ErrorCodeAccess
	}
	if errors.Is(err, fs.ErrExist) {
		return ErrorCodeExist
	}
	if errors.Is(err, fs.ErrNotExist) {
		return ErrorCodeNoEntry
	}
	if errors.Is(err, fs.ErrInvalid) {
		return ErrorCodeInvalid
	}
	var errno syscall.Errno
	if errors.As(err, &errno) {
		switch errno {
		case unix.EACCES:
			return ErrorCodeAccess
		case unix.EAGAIN:
			return ErrorCodeWouldBlock
		case unix.EBADF:
			return ErrorCodeBadDescriptor
		case unix.EBUSY:
			return ErrorCodeBusy
		case unix.EEXIST:
			return ErrorCodeExist
		case unix.EFBIG:
			return ErrorCodeFileTooLarge
		case unix.EINTR:
			return ErrorCodeInterrupted
		case unix.EINVAL:
			return ErrorCodeInvalid
		case unix.EIO:
			return ErrorCodeIo
		case unix.EISDIR:
			return ErrorCodeIsDirectory
		case unix.ELOOP:
			return ErrorCodeLoop
		case unix.EMLINK:
			return ErrorCodeTooManyLinks
		case unix.ENAMETOOLONG:
			return ErrorCodeNameTooLong
		case unix.ENODEV:
			return ErrorCodeNoDevice
		case unix.ENOENT:
			return ErrorCodeNoEntry
		case unix.ENOLCK:
			return ErrorCodeNoLock
		case unix.ENOMEM:
			return ErrorCodeInsufficientMemory
		case unix.ENOSPC:
			return ErrorCodeInsufficientSpace
		case unix.ENOTDIR:
			return ErrorCodeNotDirectory
		case unix.ENOTEMPTY:
			return ErrorCodeNotEmpty
		case unix.ENOTSUP:
			return ErrorCodeUnsupported
		case unix.ENXIO:
			return ErrorCodeNoSuchDevice
		case unix.EOVERFLOW:
			return ErrorCodeOverflow
		case unix.EPERM:
			return ErrorCodeNotPermitted
		case unix.EPIPE:
			return ErrorCodePipe
		case unix.EROFS:
			return ErrorCodeReadOnly
		case unix.ESPIPE:
			return ErrorCodeInvalidSeek
		case unix.ETXTBSY:
			return ErrorCodeTextFileBusy
		case unix.EXDEV:
			return ErrorCodeCrossDevice
		}
	}
	return ErrorCodeUnsupported
}

// extraOpenFlags 把 WASI path-flags / open-flags 映射到 Unix open(2) 标志。
func extraOpenFlags(pathFlags PathFlags, openFlags OpenFlags) int {
	var f int
	if !pathFlags.SymlinkFollow {
		f |= unix.O_NOFOLLOW // 原先 OpenAt 忽略 symlink-follow=false，符号链接会被跟随
	}
	if openFlags.Directory {
		f |= unix.O_DIRECTORY // open-flags.directory 要求目标必须是目录
	}
	return f
}

func lutimes(path string, atime, mtime time.Time) error {
	ts := []unix.Timespec{
		unix.NsecToTimespec(atime.UnixNano()),
		unix.NsecToTimespec(mtime.UnixNano()),
	}
	// SetTimesAt(symlink-follow=false) 必须 utimensat(..., AT_SYMLINK_NOFOLLOW)
	return unix.UtimesNanoAt(unix.AT_FDCWD, path, ts, unix.AT_SYMLINK_NOFOLLOW)
}

// chtimesFile 对已打开的 fd 设置 atime/mtime。
// *os.File 没有 Chtimes；不能用 f.Fd()（会把文件改成阻塞模式）。
func chtimesFile(f *os.File, atime, mtime time.Time) error {
	if f == nil {
		return os.ErrInvalid
	}
	raw, err := f.SyscallConn()
	if err != nil {
		return err
	}
	tv := []unix.Timeval{
		unix.NsecToTimeval(atime.UnixNano()),
		unix.NsecToTimeval(mtime.UnixNano()),
	}
	var sysErr error
	ctrlErr := raw.Control(func(fd uintptr) {
		sysErr = unix.Futimes(int(fd), tv)
	})
	if ctrlErr != nil {
		return ctrlErr
	}
	return sysErr
}

func writePlatformFileIdentity(h hash.Hash, info fs.FileInfo) {
	if sys, ok := info.Sys().(*syscall.Stat_t); ok {
		var b [16]byte
		binary.LittleEndian.PutUint64(b[0:8], uint64(sys.Ino))
		binary.LittleEndian.PutUint64(b[8:16], uint64(sys.Dev))
		_, _ = h.Write(b[:])
	}
}
