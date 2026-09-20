package io

import (
	"errors"
	"io"
	"sync"
)

// MultiCloser 用于合并多个io.Closer，统一管理关闭操作
type MultiCloser struct {
	closers []io.Closer
	once    sync.Once // Close 幂等，避免重复关闭
}

// NewMultiCloser 创建一个新的MultiCloser，自动过滤nil值
func NewMultiCloser(closers ...io.Closer) *MultiCloser {
	mc := &MultiCloser{
		closers: make([]io.Closer, 0, len(closers)),
	}

	// 过滤掉nil的closer，避免后续关闭时panic
	for _, c := range closers {
		if c != nil {
			mc.closers = append(mc.closers, c)
		}
	}

	return mc
}

// Close 关闭所有非nil的io.Closer，返回合并后的错误
func (m *MultiCloser) Close() error {
	var errList []error
	m.once.Do(func() {
		for _, c := range m.closers {
			if c != nil {
				if err := c.Close(); err != nil {
					errList = append(errList, err)
				}
			}
		}
	})
	return errors.Join(errList...)
}
