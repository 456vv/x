package v0_2

import (
	"context"
	"sync"

	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"
)

type insecureSeedImpl struct {
	once sync.Once
	src  *insecureImpl
}

func newInsecureSeedImpl() *insecureSeedImpl {
	return &insecureSeedImpl{}
}

// InsecureSeed 实现了 insecure-seed 函数。
// 它返回一个 128 位的值（作为两个 u64），用于初始化哈希表等。
// 这里我们简单地重用 insecure 的实现来生成种子。
func (i *insecureSeedImpl) InsecureSeed(ctx context.Context) witgo.Tuple[uint64, uint64] {
	// 复用同一 PRNG，避免每次调用新建 Source；并保证并发安全
	i.once.Do(func() {
		i.src = newInsecureImpl()
	})
	return witgo.Tuple[uint64, uint64]{
		F0: i.src.GetInsecureRandomU64(ctx),
		F1: i.src.GetInsecureRandomU64(ctx),
	}
}
