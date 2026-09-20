package v0_2

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"io"
)

type randomImpl struct{}

func newRandomImpl() *randomImpl {
	return &randomImpl{}
}

func (i *randomImpl) GetRandomBytes(_ context.Context, length uint64) []byte {
	if length == 0 {
		return []byte{}
	}
	const maxChunk = 16 << 20
	// 补 io import；32 位 make 溢出；超大 length 可 OOM
	if length > uint64(^uint(0)>>1) || length > maxChunk {
		length = maxChunk
	}
	buf := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		panic(err)
	}
	return buf
}

func (i *randomImpl) GetRandomU64(_ context.Context) uint64 {
	var buf [8]byte
	if _, err := io.ReadFull(rand.Reader, buf[:]); err != nil {
		panic(err)
	}
	return binary.LittleEndian.Uint64(buf[:])
}
