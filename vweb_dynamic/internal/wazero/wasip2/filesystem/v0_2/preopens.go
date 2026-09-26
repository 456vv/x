package v0_2

import (
	"context"

	"github.com/456vv/x/vweb_dynamic/internal/wazero/manager/filesystem"
	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"
)

type preopensImpl struct {
	fsm *filesystem.Manager
}

func newPreopensImpl(fsm *filesystem.Manager) *preopensImpl {
	return &preopensImpl{fsm: fsm}
}

// GetDirectories returns the list of pre-opened directories.
func (i *preopensImpl) GetDirectories(_ context.Context) []witgo.Tuple[Descriptor, string] {
	var results []witgo.Tuple[Descriptor, string]
	i.fsm.Range(func(handle uint32, desc *filesystem.Descriptor) bool {
		// 无锁 desc.File.Stat 与 Close 是数据竞争；已关闭项由 ForPreopen 排除。
		if desc == nil || !desc.ForPreopen() {
			return true
		}
		results = append(results, witgo.Tuple[Descriptor, string]{
			F0: handle,
			F1: desc.Path,
		})
		return true
	})

	return results
}
