//go:build linux

package v0_2

import (
	"io/fs"
	"os"
	"syscall"
	"time"

	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"

	"golang.org/x/sys/unix"
)

func goFileInfoToDescriptorStat(info fs.FileInfo) DescriptorStat {
	var stat DescriptorStat
	stat.Type = goModeToDescriptorType(info.Mode())
	stat.Size = Filesize(info.Size())
	modTime := info.ModTime()
	stat.DataModificationTimestamp = witgo.Some(Datetime{
		Seconds:     uint64(modTime.Unix()),
		Nanoseconds: uint32(modTime.Nanosecond()),
	})

	if sys, ok := info.Sys().(*syscall.Stat_t); ok {
		stat.LinkCount = uint64(sys.Nlink)
		stat.DataAccessTimestamp = witgo.Some(timeToDatetime(sys.Atim))
		stat.StatusChangeTimestamp = witgo.Some(timeToDatetime(sys.Ctim))
	} else {
		stat.LinkCount = 1
		stat.DataAccessTimestamp = witgo.None[Datetime]()
		stat.StatusChangeTimestamp = witgo.None[Datetime]()
	}
	return stat
}

func GetATime(info os.FileInfo) (time.Time, error) {
	// 从FileInfo中获取底层的syscall.Stat_t
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return time.Time{}, syscall.EINVAL
	}

	// 转换为time.Time
	return time.Unix(stat.Atim.Sec, int64(stat.Atim.Nsec)), nil
}

func adviseFile(f *os.File, offset Filesize, length Filesize, advice Advice) error {
	if f == nil {
		return os.ErrInvalid
	}
	var posix int
	switch advice {
	case AdviceNormal:
		posix = unix.FADV_NORMAL
	case AdviceSequential:
		posix = unix.FADV_SEQUENTIAL
	case AdviceRandom:
		posix = unix.FADV_RANDOM
	case AdviceWillNeed:
		posix = unix.FADV_WILLNEED
	case AdviceDontNeed:
		posix = unix.FADV_DONTNEED
	case AdviceNoReuse:
		posix = unix.FADV_NOREUSE
	default:
		return unix.EINVAL
	}
	raw, err := f.SyscallConn()
	if err != nil {
		return err
	}
	var sysErr error
	ctrlErr := raw.Control(func(fd uintptr) {
		sysErr = unix.Fadvise(int(fd), int64(offset), int64(length), posix)
	})
	if ctrlErr != nil {
		return ctrlErr
	}
	return sysErr
}
