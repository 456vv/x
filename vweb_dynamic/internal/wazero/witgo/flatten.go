package witgo

import (
	"context"
	"fmt"
	"math"
	"reflect"
)

// ============================================================
// 扁平化：Go 值 -> uint64 参数流
// ============================================================

// flattenParam 将 Go 值 val 扁平化为 uint64 参数列表，附加到 flatParams 切片。
// 支持基本类型、结构体、数组、切片、字符串、flags、variant。
//
// 参数:
//   - ctx:        上下文（用于字符串/切片的内存分配）
//   - val:        要扁平化的 Go 值
//   - flatParams: 输出的 uint64 参数列表（会被追加）
//
// 返回:
//   - error: 扁平化失败时的错误
//
// 示例:
//
//	var params []uint64
//	h.flattenParam(ctx, reflect.ValueOf(42), &params)
//	// params = [42]
func (h *Host) flattenParam(ctx context.Context, val reflect.Value, flatParams *[]uint64) error {
	if h == nil || h.allocator == nil || h.module == nil {
		return fmt.Errorf("invalid host")
	}
	if flatParams == nil {
		return fmt.Errorf("flatParams is nil")
	}
	// 解引用指针
	for val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return fmt.Errorf("cannot flatten nil pointer of type %v", val.Type())
		}
		val = val.Elem()
	}
	if !val.IsValid() {
		return fmt.Errorf("invalid reflect value encountered during flattening")
	}

	typ := val.Type()
	if isVariant(typ) {
		return h.flattenVariant(ctx, val, flatParams)
	}
	if isFlags(typ) {
		return h.flattenFlags(val, flatParams)
	}

	switch val.Kind() {
	case reflect.String:
		return h.flattenString(ctx, val, flatParams)
	case reflect.Slice:
		return h.flattenSlice(ctx, val, flatParams)
	case reflect.Struct:
		return h.flattenStruct(ctx, val, flatParams)
	case reflect.Array:
		return h.flattenArray(ctx, val, flatParams)
	case reflect.Bool:
		if val.Bool() {
			*flatParams = append(*flatParams, 1)
		} else {
			*flatParams = append(*flatParams, 0)
		}
		return nil
	case reflect.Float32:
		*flatParams = append(*flatParams, uint64(math.Float32bits(float32(val.Float()))))
		return nil
	case reflect.Float64:
		*flatParams = append(*flatParams, math.Float64bits(val.Float()))
		return nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		*flatParams = append(*flatParams, val.Uint())
		return nil
	case reflect.Int8, reflect.Int16, reflect.Int32:
		*flatParams = append(*flatParams, uint64(uint32(int32(val.Int()))))
		return nil
	case reflect.Int64:
		*flatParams = append(*flatParams, uint64(val.Int()))
		return nil
	case reflect.Int:
		if typ.Size() == 4 {
			*flatParams = append(*flatParams, uint64(uint32(int32(val.Int()))))
		} else {
			*flatParams = append(*flatParams, uint64(val.Int()))
		}
		return nil
	case reflect.Uint:
		if typ.Size() == 4 {
			*flatParams = append(*flatParams, uint64(uint32(val.Uint())))
		} else {
			*flatParams = append(*flatParams, val.Uint())
		}
		return nil
	default:
		return fmt.Errorf("unsupported parameter kind for flattening: %v", val.Kind())
	}
}

// flattenStruct 递归扁平化结构体的每个导出字段。
// 字段按声明顺序依次追加到 flatParams。
func (h *Host) flattenStruct(ctx context.Context, val reflect.Value, flatParams *[]uint64) error {
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		sf := typ.Field(i)
		if !isExportedField(sf) {
			continue
		}
		if err := h.flattenParam(ctx, val.Field(i), flatParams); err != nil {
			return fmt.Errorf("failed to flatten field %s: %w", sf.Name, err)
		}
	}
	return nil
}

// flattenString 将字符串写入 guest 内存，并将 (ptr, len) 添加到 flatParams。
func (h *Host) flattenString(ctx context.Context, val reflect.Value, flatParams *[]uint64) error {
	if !val.IsValid() || val.IsZero() {
		*flatParams = append(*flatParams, 0, 0) // 空字符串表示为 (0, 0)
		return nil
	}
	s := val.String()
	if uint64(len(s)) > math.MaxUint32 {
		return fmt.Errorf("string length exceeds 32-bit limit: %d", len(s))
	}
	strLen := uint32(len(s))
	contentPtr, err := h.allocator.Allocate(ctx, strLen, 1)
	if err != nil {
		return err
	}
	if !h.module.Memory().Write(contentPtr, []byte(s)) {
		_ = h.allocator.Free(ctx, contentPtr, strLen, 1)
		return fmt.Errorf("failed to write string content for param")
	}
	*flatParams = append(*flatParams, uint64(contentPtr), uint64(strLen))
	return nil
}

// flattenSlice 将切片提升到 guest 内存，添加 (ptr, len) 到 flatParams。
func (h *Host) flattenSlice(ctx context.Context, val reflect.Value, flatParams *[]uint64) error {
	if !val.IsValid() || val.IsZero() {
		*flatParams = append(*flatParams, 0, 0)
		return nil
	}
	liftedStructPtr, err := Lift(ctx, h, val)
	if err != nil {
		return fmt.Errorf("failed to lift slice for param: %w", err)
	}
	mem := h.module.Memory()
	contentPtr, ok := mem.ReadUint32Le(liftedStructPtr)
	if !ok {
		return fmt.Errorf("failed to read content pointer for slice param at ptr %d", liftedStructPtr)
	}
	contentLen, ok := mem.ReadUint32Le(liftedStructPtr + 4)
	if !ok {
		return fmt.Errorf("failed to read content length for slice param at ptr %d", liftedStructPtr+4)
	}
	_ = h.allocator.Free(ctx, liftedStructPtr, 8, 4) // 只释放临时 header，内容保留在作用域
	*flatParams = append(*flatParams, uint64(contentPtr), uint64(contentLen))
	return nil
}

// flattenArray 逐一扁平化数组的每个元素。
func (h *Host) flattenArray(ctx context.Context, val reflect.Value, flatParams *[]uint64) error {
	for i := 0; i < val.Len(); i++ {
		if err := h.flattenParam(ctx, val.Index(i), flatParams); err != nil {
			return fmt.Errorf("failed to flatten array element %d: %w", i, err)
		}
	}
	return nil
}

// flattenFlags 将 flags 结构体编码为单个位图，添加到 flatParams。
// 每个导出布尔字段对应一位，按字段顺序排列。
// Canonical ABI 规定 1–32 个 flag 扁平化为 i32，33–64 才是 i64。
// 原先一律追加 uint64，与 flattenType 返回的 [uint32] 签名不一致，导致核心函数参数类型不匹配。
// 现在根据 layout.Size 决定写入宽度，保持与 flattenType 一致。
func (h *Host) flattenFlags(val reflect.Value, flatParams *[]uint64) error {
	layout, err := GetOrCalculateLayout(val.Type())
	if err != nil {
		return err
	}
	var bits uint64
	bitIdx := 0
	for i := 0; i < val.NumField(); i++ {
		sf := val.Type().Field(i)
		if !isExportedField(sf) {
			continue
		}
		f := val.Field(i)
		if f.Kind() == reflect.Bool && f.Bool() {
			if bitIdx >= 64 {
				return fmt.Errorf("flags bit index %d exceeds 64", bitIdx)
			}
			bits |= 1 << uint(bitIdx)
		}
		if f.Kind() == reflect.Bool {
			bitIdx++
		}
	}
	// 按布局宽度追加，保证与 flattenType 返回的核心类型一致
	if layout.Size <= 4 {
		*flatParams = append(*flatParams, uint64(uint32(bits)))
	} else {
		*flatParams = append(*flatParams, bits)
	}
	return nil
}

// flattenVariant 将 variant 编码为 (discriminant, payload...) 并填充到 flatParams，
// 载荷对齐到所有 case 中最长的扁平化长度。
func (h *Host) flattenVariant(ctx context.Context, val reflect.Value, flatParams *[]uint64) error {
	typ := val.Type()
	fields := variantFields(typ)
	maxPayloadLen := 0

	// 只需各 case 扁平化后的槽数，用 flattenType 做纯类型计算。
	// 原先 flattenParam(Zero) 会走分配器（字符串/切片 Allocate），探宽度不该碰 guest 内存。
	for _, field := range fields {
		fieldType := field.Type
		if fieldType.Kind() == reflect.Pointer {
			fieldType = fieldType.Elem()
		}
		shape, err := flattenType(fieldType)
		if err != nil {
			return fmt.Errorf("could not get shape for variant case %s: %w", field.Name, err)
		}
		if len(shape) > maxPayloadLen {
			maxPayloadLen = len(shape)
		}
	}

	if val.IsZero() {
		// Result 的 case(0) 是 Ok。零值两边都是 nil，不能编码成 Ok(零值)。
		if isResult(typ) {
			return fmt.Errorf("invalid result: no case set for %v", typ)
		}
		*flatParams = append(*flatParams, make([]uint64, 1+maxPayloadLen)...)
		return nil
	}

	for i, field := range fields {
		fieldVal := val.FieldByIndex(field.Index)
		if fieldVal.IsZero() {
			continue
		}
		activeDiscriminant := uint64(witCaseIndex(field.Tag.Get("wit"), i))
		activePayloadField := fieldVal
		if activePayloadField.Kind() == reflect.Pointer {
			if activePayloadField.IsNil() {
				*flatParams = append(*flatParams, activeDiscriminant)
				*flatParams = append(*flatParams, make([]uint64, maxPayloadLen)...)
				return nil
			}
			activePayloadField = activePayloadField.Elem()
		}
		*flatParams = append(*flatParams, activeDiscriminant)
		startLen := len(*flatParams)
		if err := h.flattenParam(ctx, activePayloadField, flatParams); err != nil {
			return err
		}
		pad := maxPayloadLen - (len(*flatParams) - startLen)
		if pad > 0 {
			*flatParams = append(*flatParams, make([]uint64, pad)...)
		}
		return nil
	}
	if isOption(typ) {
		*flatParams = append(*flatParams, make([]uint64, 1+maxPayloadLen)...)
		return nil
	}
	return fmt.Errorf("invalid variant: no case set for %v", typ)
}

// flattenType 返回类型 typ 扁平化后的 reflect.Type 列表。
// 例如 string -> [uint32, uint32] (ptr, len)。
func flattenType(typ reflect.Type) ([]reflect.Type, error) {
	return flattenTypeStack(typ, make(map[reflect.Type]struct{}))
}

// flattenTypeStack 内部实现，支持递归检测。
func flattenTypeStack(typ reflect.Type, stack map[reflect.Type]struct{}) ([]reflect.Type, error) {
	if typ == nil {
		return nil, fmt.Errorf("cannot flatten nil type")
	}
	if _, cycling := stack[typ]; cycling {
		return nil, fmt.Errorf("recursive type not supported for flattening: %v", typ)
	}
	stack[typ] = struct{}{}
	defer delete(stack, typ)

	if isVariant(typ) {
		if typ.Kind() != reflect.Struct {
			return nil, fmt.Errorf("variant type %v must be a struct", typ)
		}
		var casePayloads [][]reflect.Type
		for _, field := range variantFields(typ) {
			fieldType := field.Type
			if fieldType.Kind() == reflect.Pointer {
				fieldType = fieldType.Elem()
			}
			shape, err := flattenTypeStack(fieldType, stack)
			if err != nil {
				return nil, err
			}
			casePayloads = append(casePayloads, shape)
		}
		payload := maxFlat(casePayloads...)
		result := make([]reflect.Type, 1+len(payload))
		result[0] = typeOf[uint32]()
		copy(result[1:], payload)
		for i := 1; i < len(result); i++ {
			if result[i] == nil {
				result[i] = typeOf[uint32]()
			}
		}
		return result, nil
	}

	if isFlags(typ) {
		layout, err := GetOrCalculateLayout(typ)
		if err != nil {
			return nil, err
		}
		// Canonical ABI 规定 1–32 个 flag 扁平化为 i32，33–64 才是 i64。
		// 原先一律 uint64，和 wit-bindgen 导出的 i32 参数对不齐。
		// layout.Size 已按位数取 1/2/4/8，与内存布局一致。
		if layout.Size <= 4 {
			return []reflect.Type{typeOf[uint32]()}, nil
		}
		return []reflect.Type{typeOf[uint64]()}, nil
	}

	switch typ.Kind() {
	case reflect.Bool:
		return []reflect.Type{typeOf[uint32]()}, nil
	case reflect.Int8, reflect.Uint8:
		return []reflect.Type{typeOf[uint32]()}, nil
	case reflect.Int16, reflect.Uint16:
		return []reflect.Type{typeOf[uint32]()}, nil
	case reflect.Int32, reflect.Uint32, reflect.Float32:
		return []reflect.Type{typeOf[uint32]()}, nil
	case reflect.Int64, reflect.Uint64, reflect.Float64:
		return []reflect.Type{typeOf[uint64]()}, nil
	case reflect.Int, reflect.Uint:
		if typ.Size() == 4 {
			return []reflect.Type{typeOf[uint32]()}, nil
		}
		return []reflect.Type{typeOf[uint64]()}, nil
	case reflect.String, reflect.Slice:
		return []reflect.Type{typeOf[uint32](), typeOf[uint32]()}, nil // {ptr, len}
	case reflect.Struct:
		var result []reflect.Type
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			if !isExportedField(field) {
				continue
			}
			fields, err := flattenTypeStack(field.Type, stack)
			if err != nil {
				return nil, err
			}
			result = append(result, fields...)
		}
		return result, nil
	case reflect.Array:
		elemTypes, err := flattenTypeStack(typ.Elem(), stack)
		if err != nil {
			return nil, err
		}
		result := make([]reflect.Type, typ.Len()*len(elemTypes))
		for i := 0; i < typ.Len(); i++ {
			copy(result[i*len(elemTypes):(i+1)*len(elemTypes)], elemTypes)
		}
		return result, nil
	case reflect.Pointer:
		return flattenTypeStack(typ.Elem(), stack)
	default:
		return nil, fmt.Errorf("unsupported type for flattening: %v", typ)
	}
}
