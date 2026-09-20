//go:build unix

package v0_2

import (
	"errors"
	"syscall"

	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"

	"golang.org/x/sys/unix"
)

func mapHTTPSysErrno(err error) (ErrorCode, bool) {
	switch {
	case errors.Is(err, syscall.ETIMEDOUT), errors.Is(err, unix.ETIMEDOUT):
		return ErrorCode{ConnectionTimeout: &witgo.Unit{}}, true
	case errors.Is(err, syscall.ECONNREFUSED), errors.Is(err, unix.ECONNREFUSED):
		return ErrorCode{ConnectionRefused: &witgo.Unit{}}, true
	case errors.Is(err, syscall.ECONNRESET), errors.Is(err, unix.ECONNRESET):
		return ErrorCode{ConnectionTerminated: &witgo.Unit{}}, true
	}
	return ErrorCode{}, false
}
