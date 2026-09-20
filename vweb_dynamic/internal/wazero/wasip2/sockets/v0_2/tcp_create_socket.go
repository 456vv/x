package v0_2

import (
	"github.com/456vv/x/vweb_dynamic/internal/wazero/wasip2"
)

type tcpCreateSocketImpl struct {
	host *wasip2.Host
}

func newTCPCreateSocketImpl(h *wasip2.Host) *tcpCreateSocketImpl {
	return &tcpCreateSocketImpl{host: h}
}
