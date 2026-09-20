//go:build !unix && !windows

package v0_2

import "syscall"

func mapNetErrno(errno syscall.Errno) (ErrorCode, bool) {
	return 0, false
}
