//go:build windows

package v0_2

import (
	"errors"
	"syscall"

	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"

	"golang.org/x/sys/windows"
)

func mapHTTPSysErrno(err error) (ErrorCode, bool) {
	switch {
	case errors.Is(err, syscall.ETIMEDOUT), errors.Is(err, windows.WSAETIMEDOUT):
		return ErrorCode{ConnectionTimeout: &witgo.Unit{}}, true
	case errors.Is(err, syscall.ECONNREFUSED), errors.Is(err, windows.WSAECONNREFUSED):
		return ErrorCode{ConnectionRefused: &witgo.Unit{}}, true
	case errors.Is(err, syscall.ECONNRESET), errors.Is(err, windows.WSAECONNRESET):
		return ErrorCode{ConnectionTerminated: &witgo.Unit{}}, true
	}
	return ErrorCode{}, false
}
