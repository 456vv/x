//go:build !unix && !windows

package v0_2

func mapHTTPSysErrno(err error) (ErrorCode, bool) {
	return ErrorCode{}, false
}
