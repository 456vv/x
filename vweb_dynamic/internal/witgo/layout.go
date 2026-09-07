package witgo

import (
	"fmt"
	"reflect"
	"sync"
)

// ============================================================
// 类型布局计算：描述 WIT 类型在内存中的排布
// ============================================================

// TypeLayout 描述 WIT 类型在内存中的布局。
// 包含大小、对齐、字段偏移等关键信息，用于正确的内存读写。
type TypeLayout struct {
	Size          uint32        // 类型总大小（字节）
	Alignment     uint32        // 对齐要求（字节，2 的幂）
	Fields        []FieldLayout // 用于结构体：各字段的布局
	DiscSize      uint32        // 用于 variant：判别值大小（1/2/4 字节）
	PayloadOffset uint32        // 用于 variant：载荷相对于结构体起始的偏移
	IsSum         bool          // 是否为 variant/sum 类型
}

// FieldLayout 描述结构体中一个字段的布局。
type FieldLayout struct {
	StructField reflect.StructField // Go 结构体字段信息
	Offset      uint32              // 字段相对于结构体起始的偏移
	Layout      *TypeLayout         // 字段自身的布局
}

var layoutCache = sync.Map{} // 全局缓存，避免重复计算布局

// GetOrCalculateLayout 获取或计算类型 typ 的布局，结果缓存在全局缓存中。
// 布局计算是昂贵的（尤其对嵌套结构体），缓存可显著提升性能。
//
// 参数:
//   - typ: Go 类型
//
// 返回:
//   - *TypeLayout: 计算得到的布局
//   - error: 类型不支持时的错误
//
// 示例:
//
//	layout, err := GetOrCalculateLayout(reflect.TypeOf(MyStruct{}))
func GetOrCalculateLayout(typ reflect.Type) (*TypeLayout, error) {
	if typ == nil {
		return nil, fmt.Errorf("cannot calculate layout for nil type")
	}
	if layout, ok := layoutCache.Load(typ); ok {
		return layout.(*TypeLayout), nil
	}
	layout, err := calculateLayout(typ)
	if err != nil {
		return nil, err
	}
	actual, _ := layoutCache.LoadOrStore(typ, layout)
	return actual.(*TypeLayout), nil
}

// calculateLayout 内部计算布局。
// 根据类型种类分发到不同的计算函数。
func calculateLayout(typ reflect.Type) (*TypeLayout, error) {
	// 先检查特殊类型
	if isVariant(typ) {
		return calculateSumLayout(typ, getVariantCaseTypes(typ))
	}
	if isFlags(typ) {
		numFlags := 0
		for i := 0; i < typ.NumField(); i++ {
			f := typ.Field(i)
			if isExportedField(f) && f.Type.Kind() == reflect.Bool {
				numFlags++
			}
		}
		// flags 编码为位图，大小取决于位数
		switch {
		case numFlags <= 8:
			return &TypeLayout{Size: 1, Alignment: 1}, nil
		case numFlags <= 16:
			return &TypeLayout{Size: 2, Alignment: 2}, nil
		case numFlags <= 32:
			return &TypeLayout{Size: 4, Alignment: 4}, nil
		case numFlags <= 64:
			return &TypeLayout{Size: 8, Alignment: 8}, nil
		default:
			return nil, fmt.Errorf("flags with >64 members not supported (got %d)", numFlags)
		}
	}

	switch typ.Kind() {
	case reflect.Uint8, reflect.Int8, reflect.Bool:
		return &TypeLayout{Size: 1, Alignment: 1}, nil
	case reflect.Uint16, reflect.Int16:
		return &TypeLayout{Size: 2, Alignment: 2}, nil
	case reflect.Uint32, reflect.Int32, reflect.Float32:
		return &TypeLayout{Size: 4, Alignment: 4}, nil
	case reflect.Uint64, reflect.Int64, reflect.Float64:
		return &TypeLayout{Size: 8, Alignment: 8}, nil
	case reflect.Int, reflect.Uint:
		size := uint32(typ.Size())
		if size != 4 && size != 8 {
			return nil, fmt.Errorf("unsupported int/uint size %d for layout", size)
		}
		return &TypeLayout{Size: size, Alignment: size}, nil
	case reflect.String, reflect.Slice:
		return &TypeLayout{Size: 8, Alignment: 4}, nil // {ptr, len} 各 4 字节
	case reflect.Struct:
		return calculateStructLayout(typ)
	case reflect.Array:
		return calculateArrayLayout(typ)
	case reflect.Pointer:
		return GetOrCalculateLayout(typ.Elem())
	default:
		return nil, fmt.Errorf("unsupported type for layout calculation: %v", typ)
	}
}

// calculateSumLayout 计算 variant/option/result 的布局（判别值 + 最大载荷）。
// variant 在内存中排列为: [discriminator][padding][payload]
// 总大小为 (payloadOffset + maxPayloadSize) 向上对齐到 maxPayloadAlignment。
//
// 参数:
//   - typ:  variant 类型
//   - cases: 各 case 的类型列表
//
// 返回:
//   - *TypeLayout: variant 的内存布局
//   - error: 计算失败时的错误
func calculateSumLayout(typ reflect.Type, cases []reflect.Type) (*TypeLayout, error) {
	numCases := len(cases)
	if numCases == 0 {
		numCases = typ.NumField()
	}

	// 根据 case 数量确定判别值大小
	var discSize uint32
	switch {
	case numCases <= 256:
		discSize = 1
	case numCases <= 65536:
		discSize = 2
	default:
		discSize = 4
	}

	// 找到最大载荷大小和对齐
	var maxCaseSize, maxCaseAlign uint32 = 0, 1
	for _, caseType := range cases {
		if caseType.Kind() == reflect.Pointer {
			caseType = caseType.Elem()
		}
		caseLayout, err := GetOrCalculateLayout(caseType)
		if err != nil {
			return nil, err
		}
		if caseLayout.Size > maxCaseSize {
			maxCaseSize = caseLayout.Size
		}
		if caseLayout.Alignment > maxCaseAlign {
			maxCaseAlign = caseLayout.Alignment
		}
	}

	// payload 紧跟 discriminator，需要对齐
	payloadOffset := align(discSize, maxCaseAlign)
	totalSize := payloadOffset + maxCaseSize

	return &TypeLayout{
		Size:          align(totalSize, maxCaseAlign), // 整体对齐
		Alignment:     maxCaseAlign,
		DiscSize:      discSize,
		PayloadOffset: payloadOffset,
		IsSum:         true,
	}, nil
}

func getVariantCaseTypes(typ reflect.Type) []reflect.Type {
	var types []reflect.Type
	for i := 0; i < typ.NumField(); i++ {
		if _, ok := typ.Field(i).Tag.Lookup("wit"); ok {
			types = append(types, typ.Field(i).Type)
		}
	}
	return types
}

// calculateStructLayout 计算结构体布局，按字段顺序排列并考虑对齐。
// 每个字段的偏移 = align(当前偏移, 字段对齐)，然后偏移 += 字段大小。
//
// 参数:
//   - typ: 结构体类型
//
// 返回:
//   - *TypeLayout: 结构体的内存布局
//   - error: 计算失败时的错误
func calculateStructLayout(typ reflect.Type) (*TypeLayout, error) {
	var fields []FieldLayout
	var currentOffset uint32 = 0
	var maxAlignment uint32 = 1

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !isExportedField(field) {
			continue
		}
		fieldLayout, err := GetOrCalculateLayout(field.Type)
		if err != nil {
			return nil, fmt.Errorf("field %s: %w", field.Name, err)
		}

		// 对齐当前偏移到字段对齐要求
		currentOffset = align(currentOffset, fieldLayout.Alignment)
		fields = append(fields, FieldLayout{
			StructField: field,
			Offset:      currentOffset,
			Layout:      fieldLayout,
		})

		currentOffset += fieldLayout.Size
		if fieldLayout.Alignment > maxAlignment {
			maxAlignment = fieldLayout.Alignment
		}
	}

	return &TypeLayout{
		Size:      align(currentOffset, maxAlignment), // 整体大小对齐
		Alignment: maxAlignment,
		Fields:    fields,
	}, nil
}

// calculateArrayLayout 计算数组（固定长度）的布局，连续存储元素。
// 数组大小 = len * 元素大小（考虑元素间对齐）。
//
// 参数:
//   - typ: 数组类型
//
// 返回:
//   - *TypeLayout: 数组的内存布局
//   - error: 计算失败时的错误
func calculateArrayLayout(typ reflect.Type) (*TypeLayout, error) {
	if typ.Len() == 0 {
		return &TypeLayout{Size: 0, Alignment: 1}, nil
	}
	elemLayout, err := GetOrCalculateLayout(typ.Elem())
	if err != nil {
		return nil, err
	}
	var currentOffset uint32 = 0
	for i := 0; i < typ.Len(); i++ {
		currentOffset = align(currentOffset, elemLayout.Alignment)
		currentOffset += elemLayout.Size
	}
	return &TypeLayout{
		Size:      align(currentOffset, elemLayout.Alignment),
		Alignment: elemLayout.Alignment,
	}, nil
}
