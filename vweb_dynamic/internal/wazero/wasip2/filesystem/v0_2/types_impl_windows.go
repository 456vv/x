//go:build windows

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

	"golang.org/x/sys/windows"
)

func (i *typesImpl) GetFlags(ctx context.Context, this Descriptor) witgo.Result[DescriptorFlags, ErrorCode] {
	d, ok := i.host.FilesystemManager().Get(this)
	if !ok || d == nil || d.File == nil {
		return witgo.Err[DescriptorFlags, ErrorCode](ErrorCodeBadDescriptor)
	}

	info, err := d.File.Stat()
	if err != nil {
		return witgo.Err[DescriptorFlags, ErrorCode](mapOsError(err))
	}

	var wasiFlags DescriptorFlags
	wasiFlags.Read = true
	if info.Mode()&0o200 == 0 {
		wasiFlags.Write = false
	} else {
		wasiFlags.Write = true
	}
	wasiFlags.FileIntegritySync = false
	wasiFlags.DataIntegritySync = false
	wasiFlags.RequestedWriteSync = false

	return witgo.Ok[DescriptorFlags, ErrorCode](wasiFlags)
}

func timeToDatetime(ft syscall.Filetime) Datetime {
	t := time.Unix(0, ft.Nanoseconds()).UTC()
	return Datetime{
		Seconds:     uint64(t.Unix()),
		Nanoseconds: uint32(t.Nanosecond()),
	}
}

func goFileInfoToDescriptorStat(info fs.FileInfo) DescriptorStat {
	var stat DescriptorStat
	stat.Type = goModeToDescriptorType(info.Mode())
	stat.Size = Filesize(info.Size())
	modTime := info.ModTime()
	stat.DataModificationTimestamp = witgo.Some(Datetime{
		Seconds:     uint64(modTime.Unix()),
		Nanoseconds: uint32(modTime.Nanosecond()),
	})

	if sys, ok := info.Sys().(*syscall.Win32FileAttributeData); ok {
		stat.DataAccessTimestamp = witgo.Some(timeToDatetime(sys.LastAccessTime))
		stat.StatusChangeTimestamp = witgo.Some(timeToDatetime(sys.LastWriteTime))
	} else {
		stat.DataAccessTimestamp = witgo.None[Datetime]()
		stat.StatusChangeTimestamp = witgo.None[Datetime]()
	}
	stat.LinkCount = 1
	return stat
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
		// 原 Windows 分支把块设备也映射成 CharacterDevice
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
		case windows.ERROR_ACCESS_DENIED:
			return ErrorCodeAccess
		case windows.ERROR_ALREADY_EXISTS, windows.ERROR_FILE_EXISTS:
			return ErrorCodeExist
		case windows.ERROR_FILE_NOT_FOUND, windows.ERROR_PATH_NOT_FOUND:
			return ErrorCodeNoEntry
		case windows.ERROR_DIR_NOT_EMPTY:
			return ErrorCodeNotEmpty
		case windows.ERROR_INVALID_HANDLE:
			return ErrorCodeBadDescriptor
		case windows.ERROR_INVALID_PARAMETER:
			return ErrorCodeInvalid
		case windows.ERROR_SHARING_VIOLATION:
			return ErrorCodeBusy
		case windows.ERROR_NOT_SUPPORTED:
			return ErrorCodeUnsupported
		case windows.ERROR_DISK_FULL:
			return ErrorCodeInsufficientSpace
		case windows.ERROR_BROKEN_PIPE:
			return ErrorCodePipe
		case windows.ERROR_NOT_A_REPARSE_POINT:
			return ErrorCodeInvalid
		case windows.ERROR_DIRECTORY:
			return ErrorCodeNotDirectory
		}
	}
	return ErrorCodeUnsupported
}

func GetATime(info os.FileInfo) (time.Time, error) {
	data, ok := info.Sys().(*syscall.Win32FileAttributeData)
	if !ok {
		return time.Time{}, syscall.EINVAL
	}
	return time.Unix(0, data.LastAccessTime.Nanoseconds()), nil
}

func extraOpenFlags(_ PathFlags, _ OpenFlags) int {
	// Windows 的 os.OpenFile 没有等价的 O_NOFOLLOW/O_DIRECTORY；
	// symlink-follow 由 Lstat/Open 语义处理，directory 由 OpenAt 打开后 Stat 校验。
	return 0
}

func adviseFile(_ *os.File, _ Filesize, _ Filesize, _ Advice) error {
	return errAdviseUnsupported // Windows 无 posix_fadvise
}

func lutimes(path string, atime, mtime time.Time) error {
	// 非 unix 无稳定的 AT_SYMLINK_NOFOLLOW；尽力而为走 os.Chtimes。
	return os.Chtimes(path, atime, mtime)
}

// chtimesFile 对已打开的 HANDLE 调用 SetFileTime。
// *os.File 没有 Chtimes；路径 os.Chtimes 在文件已 unlink 时会失败。
func chtimesFile(f *os.File, atime, mtime time.Time) error {
	if f == nil {
		return os.ErrInvalid
	}
	a := windows.NsecToFiletime(atime.UnixNano())
	m := windows.NsecToFiletime(mtime.UnixNano())
	raw, err := f.SyscallConn()
	if err != nil {
		return err
	}
	var sysErr error
	ctrlErr := raw.Control(func(fd uintptr) {
		sysErr = windows.SetFileTime(windows.Handle(fd), nil, &a, &m)
	})
	if ctrlErr != nil {
		return ctrlErr
	}
	return sysErr
}

func writePlatformFileIdentity(h hash.Hash, info fs.FileInfo) {
	//  Win32FileAttributeData 不含 inode；写入属性位与 unix 签名对齐即可。
	if sys, ok := info.Sys().(*syscall.Win32FileAttributeData); ok {
		var b [4]byte
		binary.LittleEndian.PutUint32(b[:], sys.FileAttributes)
		_, _ = h.Write(b[:])
	}
}
