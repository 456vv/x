//go:build windows

package v0_2

import (
	"syscall"

	"golang.org/x/sys/windows"
)

func mapNetErrno(errno syscall.Errno) (ErrorCode, bool) {
	switch errno {
	case windows.WSAEACCES, windows.ERROR_ACCESS_DENIED:
		return ErrorCodeAccessDenied, true
	case windows.WSAEADDRINUSE:
		return ErrorCodeAddressInUse, true
	case windows.WSAEADDRNOTAVAIL:
		return ErrorCodeAddressNotBindable, true
	case windows.WSAEAFNOSUPPORT:
		return ErrorCodeNotSupported, true
	case windows.WSAEALREADY:
		return ErrorCodeConcurrencyConflict, true
	case windows.WSAECONNABORTED:
		return ErrorCodeConnectionAborted, true
	case windows.WSAECONNREFUSED:
		return ErrorCodeConnectionRefused, true
	case windows.WSAECONNRESET:
		return ErrorCodeConnectionReset, true
	case windows.WSAEINPROGRESS:
		return ErrorCodeConcurrencyConflict, true
	case windows.WSAEINVAL:
		return ErrorCodeInvalidArgument, true
	case windows.WSAEISCONN:
		return ErrorCodeInvalidState, true
	case windows.WSAENETUNREACH:
		return ErrorCodeRemoteUnreachable, true
	case windows.WSAEMFILE:
		return ErrorCodeNewSocketLimit, true
	case windows.WSAENOTCONN:
		return ErrorCodeInvalidState, true
	case windows.WSAEOPNOTSUPP:
		return ErrorCodeNotSupported, true
	case windows.WSAETIMEDOUT:
		return ErrorCodeTimeout, true
	case windows.WSAEWOULDBLOCK:
		return ErrorCodeWouldBlock, true
	}
	return 0, false
}
