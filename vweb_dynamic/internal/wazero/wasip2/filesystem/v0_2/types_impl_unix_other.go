//go:build unix && !linux && !darwin

package v0_2

import "os"

func adviseFile(_ *os.File, _ Filesize, _ Filesize, _ Advice) error {
	// 非 Linux 的 unix 默认无 Fadvise 绑定，避免链接失败
	return errAdviseUnsupported
}

func syncDataFile(f *os.File) error {
	if f == nil {
		return os.ErrInvalid
	}
	// 退回 fsync，满足 sync-data 的落盘语义。
	return f.Sync()
}
