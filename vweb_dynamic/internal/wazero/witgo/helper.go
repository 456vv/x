package witgo

import (
	"math"
	"reflect"
	"strconv"
	"strings"
)

// ============================================================
// 工具函数：内存对齐、类型检测等
// ============================================================

// align 将 ptr 按 alignment 向上对齐，处理非 2 次幂，并防止 uint32 溢出。
// 例如: align(3, 4) = 4, align(5, 8) = 8, align(8, 4) = 8
//
// 参数:
//   - ptr:        原始指针值
//   - alignment:  对齐要求
//
// 返回:
//   - uint32: 对齐后的指针值
func align(ptr, alignment uint32) uint32 {
	if alignment == 0 || alignment == 1 {
		return ptr
	}
	if alignment&(alignment-1) == 0 {
		// 2 的幂：使用位运算快速对齐
		add := alignment - 1
		if ptr > math.MaxUint32-add {
			return math.MaxUint32
		}
		return (ptr + add) &^ add
	}
	// 非 2 的幂：使用取模
	rem := ptr % alignment
	if rem == 0 {
		return ptr
	}
	add := alignment - rem
	if math.MaxUint32-ptr < add {
		return math.MaxUint32
	}
	return ptr + add
}

// maxU32 返回两个 uint32 中的较大值。
func maxU32(a, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}

// addU32 安全地计算 base+off，若溢出返回错误。
//
// 参数:
//   - base: 基地址
//   - off:  偏移量
//
// 返回:
//   - uint32: 结果地址
//   - error:  溢出时的错误
func addU32(base, off uint32) (uint32, error) {
	if off > math.MaxUint32-base {
		return 0, errPtrOverflow
	}
	return base + off, nil
}

// isVariant 检查结构体是否包含带 "wit" tag 的字段，用于识别 variant/option/result。
// variant 类型在 WIT 中表示可区分的联合类型（类似 Rust 的 enum）。
//
// 参数:
//   - typ: Go 类型
//
// 返回:
//   - bool: 是否为 variant 类型
func isVariant(typ reflect.Type) bool {
	if typ == nil || typ.Kind() != reflect.Struct {
		return false
	}
	for i := 0; i < typ.NumField(); i++ {
		if _, ok := typ.Field(i).Tag.Lookup("wit"); ok {
			return true
		}
	}
	return false
}

// variantFields 返回结构体中所有带有 "wit" tag 的字段。
// 这些字段代表 variant 的各个 case。
//
// 参数:
//   - typ: variant 结构体类型
//
// 返回:
//   - []reflect.StructField: 带 wit tag 的字段列表
func variantFields(typ reflect.Type) []reflect.StructField {
	if typ == nil || typ.Kind() != reflect.Struct {
		return nil
	}
	fields := make([]reflect.StructField, 0, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if _, ok := f.Tag.Lookup("wit"); ok {
			fields = append(fields, f)
		}
	}
	return fields
}

// witCaseIndex 从 tag 中解析 `wit:"case(N)"`，返回整数值；解析失败则返回 fallback。
//
// 参数:
//   - tag:     wit tag 字符串
//   - fallback: 解析失败时的默认值
//
// 返回:
//   - int: case 索引值
func witCaseIndex(tag string, fallback int) int {
	if tag == "" {
		return fallback
	}
	const prefix = "case("
	i := strings.Index(tag, prefix)
	if i < 0 {
		return fallback
	}
	rest := tag[i+len(prefix):]
	j := strings.IndexByte(rest, ')')
	if j <= 0 {
		return fallback
	}
	n, err := strconv.Atoi(rest[:j])
	if err != nil || n < 0 {
		return fallback
	}
	return n
}

// Flagger 是标志类型的标记接口。
// 实现此接口的结构体被视为 flags 类型。
type Flagger interface {
	IsFlags()
}

var flaggerType = reflect.TypeFor[Flagger]()

// isFlags 检查类型是否实现了 Flagger 接口（即标志类型）。
// flags 类型被编码为单个整数位图。
//
// 参数:
//   - typ: Go 类型
//
// 返回:
//   - bool: 是否为 flags 类型
func isFlags(typ reflect.Type) bool {
	if typ == nil || typ.Kind() != reflect.Struct {
		return false
	}
	ptrType := reflect.PointerTo(typ)
	return typ.Implements(flaggerType) || ptrType.Implements(flaggerType)
}

// Optioner 是 Option 类型的标记接口。
type Optioner interface {
	IsOption()
}

var optionerType = reflect.TypeFor[Optioner]()

// isOption 检查类型是否为 Option 类型。
func isOption(typ reflect.Type) bool {
	if typ == nil || typ.Kind() != reflect.Struct {
		return false
	}
	ptrType := reflect.PointerTo(typ)
	return typ.Implements(optionerType) || ptrType.Implements(optionerType)
}

// Resulter 是 Result 类型的标记接口。
type Resulter interface {
	IsResult()
}

var resulterType = reflect.TypeFor[Resulter]()

// isResult 检查类型是否为 Result 类型。
func isResult(typ reflect.Type) bool {
	if typ == nil || typ.Kind() != reflect.Struct {
		return false
	}
	ptrType := reflect.PointerTo(typ)
	return typ.Implements(resulterType) || ptrType.Implements(resulterType)
}

// isExportedField 判断结构体字段是否导出（首字母大写）。
// Go 的 reflect 包中，未导出字段的 PkgPath 非空。
//
// 参数:
//   - f: 结构体字段信息
//
// 返回:
//   - bool: 是否导出
func isExportedField(f reflect.StructField) bool {
	return f.PkgPath == ""
}

// sliceLenFitsInt 检查 uint32 长度能否安全转换为 int（平台相关）。
// 在 32 位平台上，int 最大为 2^31-1，需要检查溢出。
//
// 参数:
//   - n: uint32 长度
//
// 返回:
//   - bool: 是否可以安全转换
func sliceLenFitsInt(n uint32) bool {
	max := uint64(^uint(0) >> 1)
	return uint64(n) <= max
}

// isExactByteSlice 判断类型是否为未命名的 []uint8（即 []byte）。
// []byte 有特殊的内存布局优化路径。
//
// 参数:
//   - elemType: 元素类型
//
// 返回:
//   - bool: 是否为 []byte 的元素类型
func isExactByteSlice(elemType reflect.Type) bool {
	return elemType != nil && elemType.Kind() == reflect.Uint8 && elemType.Name() == "" && elemType.PkgPath() == ""
}

// isScalarKind 判断 kind 是否为标量类型（布尔、整数、浮点）。
// 标量类型可以直接映射到 Wasm 的 i32/i64/f32/f64。
//
// 参数:
//   - k: reflect.Kind
//
// 返回:
//   - bool: 是否为标量类型
func isScalarKind(k reflect.Kind) bool {
	switch k {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

// typeOf 兼容 Go 1.18+；reflect.TypeFor 需要 1.22。
// var zero T 对 interface 是无类型 nil，TypeOf 会 panic；必须用 (*T)(nil).Elem()
func typeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}
