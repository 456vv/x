//go:build darwin

package v0_2

import (
	"io/fs"
	"os"
	"syscall"
	"time"

	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"
)

func goFileInfoToDescriptorStat(info fs.FileInfo) DescriptorStat {
	var stat DescriptorStat
	stat.Type = goModeToDescriptorType(info.Mode())
	stat.Size = Filesize(info.Size())
	stat.DataModificationTimestamp = witgo.Some(datetimeFromTime(info.ModTime()))

	if sys, ok := info.Sys().(*syscall.Stat_t); ok {
		stat.LinkCount = uint64(sys.Nlink)
		stat.DataAccessTimestamp = witgo.Some(timeToDatetime(sys.Atimespec))
		stat.StatusChangeTimestamp = witgo.Some(timeToDatetime(sys.Ctimespec))
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
	return time.Unix(stat.Atimespec.Sec, int64(stat.Atimespec.Nsec)), nil
}

func adviseFile(_ *os.File, _ Filesize, _ Filesize, _ Advice) error {
	// Darwin 没有 posix_fadvise/Fadvise，保持 WASI Unsupported
	return errAdviseUnsupported
}

func syncDataFile(f *os.File) error {
	if f == nil {
		return os.ErrInvalid
	}
	// 退回 fsync，满足 sync-data 的落盘语义。
	return f.Sync()
}

func timeToDatetime(ts syscall.Timespec) Datetime {
	return datetimeFromUnix(int64(ts.Sec), int64(ts.Nsec))
}
