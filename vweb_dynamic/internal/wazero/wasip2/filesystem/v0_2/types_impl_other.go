//go:build !unix && !windows

package v0_2

import (
	"context"
	"errors"
	"hash"
	"io/fs"
	"os"
	"time"

	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"
)

func (i *typesImpl) GetFlags(ctx context.Context, this Descriptor) witgo.Result[DescriptorFlags, ErrorCode] {
	return witgo.Err[DescriptorFlags, ErrorCode](ErrorCodeUnsupported)
}

func goFileInfoToDescriptorStat(info fs.FileInfo) DescriptorStat {
	var stat DescriptorStat
	stat.Type = goModeToDescriptorType(info.Mode())
	stat.Size = Filesize(info.Size())
	stat.DataModificationTimestamp = witgo.Some(datetimeFromTime(info.ModTime()))

	// For non-Unix platforms, we provide best-effort information.
	stat.DataAccessTimestamp = witgo.None[Datetime]()
	stat.StatusChangeTimestamp = witgo.None[Datetime]()
	// Link count might not be available.
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
	if errors.Is(err, os.ErrClosed) {
		return ErrorCodeBadDescriptor
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

	return ErrorCodeUnsupported
}

func GetATime(info os.FileInfo) (time.Time, error) {
	return time.Time{}, errors.New("GetATime not supported on this platform")
}

func extraOpenFlags(_ PathFlags, _ OpenFlags) int {
	return 0 // 非 unix/windows 无对应 open 标志
}

func adviseFile(_ *os.File, _ Filesize, _ Filesize, _ Advice) error {
	return errAdviseUnsupported
}

func lutimes(path string, atime, mtime time.Time) error {
	// 非 unix 无稳定的 AT_SYMLINK_NOFOLLOW；尽力而为走 os.Chtimes。
	return os.Chtimes(path, atime, mtime)
}

// chtimesFile 在无 futimes/SetFileTime 的平台上回退到路径 API。
// 注意：*os.File 没有 Chtimes 方法。
func chtimesFile(f *os.File, atime, mtime time.Time) error {
	if f == nil {
		return os.ErrInvalid
	}
	return os.Chtimes(f.Name(), atime, mtime)
}

func writePlatformFileIdentity(_ hash.Hash, _ fs.FileInfo) {}
func syncDataFile(f *os.File) error {
	if f == nil {
		return os.ErrInvalid
	}
	// 退回 fsync，满足 sync-data 的落盘语义。
	return f.Sync()
}
