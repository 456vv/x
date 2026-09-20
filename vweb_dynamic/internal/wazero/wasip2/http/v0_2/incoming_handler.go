package v0_2

import (
	"context"

	manager_http "github.com/456vv/x/vweb_dynamic/internal/wazero/manager/http"
	"github.com/456vv/x/vweb_dynamic/internal/wazero/wasip2"

	"github.com/tetratelabs/wazero"
)

// --- wasi:http/incoming-handler@0.2.0 implementation ---

type incomingHandler struct {
	hm *manager_http.HTTPManager
}

func NewIncomingHandler(hm *manager_http.HTTPManager) wasip2.Implementation {
	return &incomingHandler{hm: hm}
}

func (i *incomingHandler) Name() string       { return "wasi:http/incoming-handler" }
func (i *incomingHandler) Versions() []string {
	return []string{"0.2.0", "0.2.1", "0.2.2", "0.2.3", "0.2.4", "0.2.5", "0.2.6", "0.2.7"}
}

func (i *incomingHandler) Instantiate(_ context.Context, h *wasip2.Host, builder wazero.HostModuleBuilder) error {
	// handle 由 guest 导出，host 不导出
	return nil
}
