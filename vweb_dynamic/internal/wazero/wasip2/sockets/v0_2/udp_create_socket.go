package v0_2

import (
	"github.com/456vv/x/vweb_dynamic/internal/wazero/wasip2"
)

type udpCreateSocketImpl struct {
	host *wasip2.Host
}

func newUDPCreateSocketImpl(h *wasip2.Host) *udpCreateSocketImpl {
	return &udpCreateSocketImpl{host: h}
}
