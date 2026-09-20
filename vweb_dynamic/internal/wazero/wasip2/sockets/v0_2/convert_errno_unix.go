//go:build unix

package v0_2

import (
	"syscall"

	"golang.org/x/sys/unix"
)

func mapNetErrno(errno syscall.Errno) (ErrorCode, bool) {
	if errno == unix.EAGAIN || errno == unix.EWOULDBLOCK {
		return ErrorCodeWouldBlock, true
	}
	switch errno {
	case unix.EACCES, unix.EPERM:
		return ErrorCodeAccessDenied, true
	case unix.EADDRINUSE:
		return ErrorCodeAddressInUse, true
	case unix.EADDRNOTAVAIL:
		return ErrorCodeAddressNotBindable, true
	case unix.EAFNOSUPPORT:
		return ErrorCodeNotSupported, true
	case unix.EALREADY:
		return ErrorCodeConcurrencyConflict, true
	case unix.ECONNABORTED:
		return ErrorCodeConnectionAborted, true
	case unix.ECONNREFUSED:
		return ErrorCodeConnectionRefused, true
	case unix.ECONNRESET:
		return ErrorCodeConnectionReset, true
	case unix.EINPROGRESS:
		return ErrorCodeConcurrencyConflict, true
	case unix.EINVAL:
		return ErrorCodeInvalidArgument, true
	case unix.EISCONN:
		return ErrorCodeInvalidState, true
	case unix.ENETUNREACH:
		return ErrorCodeRemoteUnreachable, true
	case unix.ENFILE, unix.EMFILE:
		return ErrorCodeNewSocketLimit, true
	case unix.ENOTCONN:
		return ErrorCodeInvalidState, true
	case unix.EOPNOTSUPP:
		return ErrorCodeNotSupported, true
	case unix.ETIMEDOUT:
		return ErrorCodeTimeout, true
	}
	return 0, false
}
