package v0_2

import (
	"context"
	"math/rand"
	"sync"
	"time"
)

// insecureImpl 实现了 wasi:random/insecure 接口。
// 它为每个实例包含一个独立的、正确初始化的随机数生成器。
type insecureImpl struct {
	mu sync.Mutex // rand.Rand 非并发安全的，因此需要加锁以确保线程安全。
	r  *rand.Rand
}

// newInsecureImpl 创建一个新的 insecureImpl 实例。
// 它使用 rand.NewSource 来创建一个新的、非共享的随机数源，
// 这是当前推荐的最佳实践。
func newInsecureImpl() *insecureImpl {
	return &insecureImpl{
		r: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GetInsecureRandomBytes 实现了 get-insecure-random-bytes 函数。
func (i *insecureImpl) GetInsecureRandomBytes(_ context.Context, length uint64) []byte {
	if length == 0 {
		return []byte{}
	}
	const maxChunk = 16 << 20
	// 补 io import；32 位 make 溢出；超大 length 可 OOM
	if length > uint64(^uint(0)>>1) || length > maxChunk {
		length = maxChunk
	}
	buf := make([]byte, length)
	i.mu.Lock()
	_, err := i.r.Read(buf)
	i.mu.Unlock()
	if err != nil {
		panic(err)
	}
	return buf
}

// GetInsecureRandomU64 实现了 get-insecure-random-u64 函数。
func (i *insecureImpl) GetInsecureRandomU64(_ context.Context) uint64 {
	i.mu.Lock()
	v := i.r.Uint64()
	i.mu.Unlock()
	return v
}
