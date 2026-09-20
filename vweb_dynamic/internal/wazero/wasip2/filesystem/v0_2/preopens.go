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
	// Note: In a real implementation, this would iterate over pre-opened
	// directories configured by the host environment. For this example,
	// we assume the manager contains only pre-opens.
	i.fsm.Range(func(handle uint32, desc *filesystem.Descriptor) bool {
		if desc == nil {
			return true
		}
		// get-directories 只应返回目录 preopen；文件描述符混入会让 guest 对文件调 read-directory。
		if desc.File != nil {
			st, err := desc.File.Stat()
			if err != nil || !st.IsDir() {
				return true
			}
		}
		results = append(results, witgo.Tuple[Descriptor, string]{
			F0: handle,
			F1: desc.Path,
		})
		return true
	})
	return results
}
