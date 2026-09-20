package witgo

import (
	"fmt"
	"reflect"
)

// ============================================================
// 参数流与工具函数
// ============================================================

// paramStream 是一个简单的参数流，用于顺序消费 uint64 参数。
// 模拟 Wasm 函数的 flat 参数栈，按顺序取出值。
type paramStream struct {
	params []uint64 // 扁平化后的参数列表
	pos    int      // 当前读取位置
}

// Next 返回下一个参数，若已耗尽则返回 (0, false)。
//
// 返回:
//   - uint64: 下一个参数值
//   - bool:   是否成功读取
func (s *paramStream) Next() (uint64, bool) {
	if s == nil || s.pos >= len(s.params) {
		return 0, false
	}
	p := s.params[s.pos]
	s.pos++
	return p, true
}

// maxType 返回两个 reflect.Type 中"较大"的类型，用于统一 variant 载荷的扁平类型。
// 优先级: float64 > float32 > int64 > int32 > uint32
func maxType(a, b reflect.Type) reflect.Type {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a == b {
		return a
	}
	aKind, bKind := a.Kind(), b.Kind()
	if aKind == reflect.Float64 || bKind == reflect.Float64 {
		return typeOf[float64]()
	}
	if aKind == reflect.Float32 || bKind == reflect.Float32 {
		return typeOf[float32]()
	}
	aIs64 := aKind == reflect.Int64 || aKind == reflect.Uint64 || (aKind == reflect.Int && a.Size() == 8) || (aKind == reflect.Uint && a.Size() == 8)
	bIs64 := bKind == reflect.Int64 || bKind == reflect.Uint64 || (bKind == reflect.Int && b.Size() == 8) || (bKind == reflect.Uint && b.Size() == 8)
	if aIs64 || bIs64 {
		return typeOf[int64]()
	}
	if aKind == reflect.Int32 || bKind == reflect.Int32 || aKind == reflect.Int || bKind == reflect.Int {
		return typeOf[int32]()
	}
	return typeOf[uint32]()
}

// maxFlat 接受多个类型列表，返回每个位置的最大类型列表（用于对齐 variant 载荷）。
// 例如: maxFlat([u32, u64], [u32, u32]) = [u32, u64]
func maxFlat(lists ...[]reflect.Type) []reflect.Type {
	maxLength := 0
	for _, l := range lists {
		if len(l) > maxLength {
			maxLength = len(l)
		}
	}
	if maxLength == 0 {
		return nil
	}
	result := make([]reflect.Type, maxLength)
	for _, l := range lists {
		for i, t := range l {
			result[i] = maxType(result[i], t)
		}
	}
	return result
}

// unflattenState 已废弃，保留仅为兼容（实际未使用）。
// 旧版使用此结构体管理反扁平化状态，现已被 paramStream 替代。
type unflattenState struct {
	flat   []uint64
	offset int
}

func (u *unflattenState) unflatten() (uint64, error) {
	if u == nil || u.offset >= len(u.flat) {
		return 0, fmt.Errorf("not enough values to unflatten")
	}
	val := u.flat[u.offset]
	u.offset++
	return val, nil
}
