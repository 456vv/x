package witgo

import (
	"context"
	"fmt"
	"math"
	"reflect"

	"github.com/tetratelabs/wazero/api"
)

// unflattenParam 从参数流 ps 中读取值，还原为 targetType 类型的 Go 值。
func (h *Host) unflattenParam(ctx context.Context, mem api.Memory, ps *paramStream, targetType reflect.Type) (reflect.Value, error) {
	if targetType == nil {
		return reflect.Value{}, fmt.Errorf("cannot unflatten nil type")
	}
	if targetType.Kind() == reflect.Pointer {
		elem, err := h.unflattenParam(ctx, mem, ps, targetType.Elem())
		if err != nil {
			return reflect.Value{}, err
		}
		ptr := reflect.New(targetType.Elem())
		if elem.IsValid() {
			ptr.Elem().Set(elem)
		}
		return ptr, nil
	}

	if isVariant(targetType) {
		return h.unflattenVariant(ctx, mem, ps, targetType)
	}
	if isFlags(targetType) {
		return h.unflattenFlags(ps, targetType)
	}

	outVal := reflect.New(targetType).Elem()
	switch targetType.Kind() {
	case reflect.String:
		return h.unflattenString(mem, ps, targetType)
	case reflect.Slice:
		return h.unflattenSlice(ctx, mem, ps, targetType)
	case reflect.Struct:
		return h.unflattenStruct(ctx, mem, ps, targetType)
	case reflect.Array:
		return h.unflattenArray(ctx, mem, ps, targetType)
	case reflect.Bool:
		p, ok := ps.Next()
		if !ok {
			return reflect.Value{}, fmt.Errorf("not enough params on stack for bool")
		}
		outVal.SetBool(p != 0)
	case reflect.Int8:
		p, ok := ps.Next()
		if !ok {
			return reflect.Value{}, fmt.Errorf("not enough params on stack for int8")
		}
		outVal.SetInt(int64(int8(p)))
	case reflect.Int16:
		p, ok := ps.Next()
		if !ok {
			return reflect.Value{}, fmt.Errorf("not enough params on stack for int16")
		}
		outVal.SetInt(int64(int16(p)))
	case reflect.Int32:
		p, ok := ps.Next()
		if !ok {
			return reflect.Value{}, fmt.Errorf("not enough params on stack for int32")
		}
		outVal.SetInt(int64(int32(p)))
	case reflect.Int64:
		p, ok := ps.Next()
		if !ok {
			return reflect.Value{}, fmt.Errorf("not enough params on stack for int64")
		}
		outVal.SetInt(int64(p))
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		p, ok := ps.Next()
		if !ok {
			return reflect.Value{}, fmt.Errorf("not enough params on stack for %v", targetType)
		}
		bits := targetType.Bits()
		if bits < 64 {
			p &= (1 << bits) - 1 // 截断到类型位宽
		}
		outVal.SetUint(p)
	case reflect.Int:
		p, ok := ps.Next()
		if !ok {
			return reflect.Value{}, fmt.Errorf("not enough params on stack for int")
		}
		if targetType.Size() == 4 {
			outVal.SetInt(int64(int32(p)))
		} else {
			outVal.SetInt(int64(p))
		}
	case reflect.Float32:
		p, ok := ps.Next()
		if !ok {
			return reflect.Value{}, fmt.Errorf("not enough params on stack for float32")
		}
		outVal.SetFloat(float64(math.Float32frombits(uint32(p))))
	case reflect.Float64:
		p, ok := ps.Next()
		if !ok {
			return reflect.Value{}, fmt.Errorf("not enough params on stack for float64")
		}
		outVal.SetFloat(math.Float64frombits(p))
	default:
		return reflect.Value{}, fmt.Errorf("unsupported type for unflattening: %v", targetType.Kind())
	}
	return outVal, nil
}

// unflattenStruct 逐个字段反扁平化。
func (h *Host) unflattenStruct(ctx context.Context, mem api.Memory, ps *paramStream, targetType reflect.Type) (reflect.Value, error) {
	outVal := reflect.New(targetType).Elem()
	for i := 0; i < outVal.NumField(); i++ {
		sf := targetType.Field(i)
		if !isExportedField(sf) {
			continue
		}
		fieldVal, err := h.unflattenParam(ctx, mem, ps, sf.Type)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("failed to unflatten field %s: %w", sf.Name, err)
		}
		fv := outVal.Field(i)
		if fv.CanSet() {
			fv.Set(fieldVal)
		}
	}
	return outVal, nil
}

// unflattenString 从参数流读取 (ptr,len) 并返回字符串。
func (h *Host) unflattenString(mem api.Memory, ps *paramStream, targetType reflect.Type) (reflect.Value, error) {
	outVal := reflect.New(targetType).Elem()
	ptr, ok1 := ps.Next()
	length, ok2 := ps.Next()
	if !ok1 || !ok2 {
		return reflect.Value{}, fmt.Errorf("not enough params on stack for string")
	}
	s, err := LowerStringFromParts(mem, uint32(ptr), uint32(length))
	if err != nil {
		return reflect.Value{}, err
	}
	outVal.SetString(s)
	return outVal, nil
}

// unflattenSlice 从参数流读取 (ptr,len) 并返回切片。
func (h *Host) unflattenSlice(ctx context.Context, mem api.Memory, ps *paramStream, targetType reflect.Type) (reflect.Value, error) {
	outVal := reflect.New(targetType).Elem()
	ptr, ok1 := ps.Next()
	length, ok2 := ps.Next()
	if !ok1 || !ok2 {
		return reflect.Value{}, fmt.Errorf("not enough params on stack for slice")
	}
	err := lowerSlice2(ctx, mem, uint32(ptr), uint32(length), outVal)
	return outVal, err
}

// unflattenArray 从参数流逐个元素读取数组。
func (h *Host) unflattenArray(ctx context.Context, mem api.Memory, ps *paramStream, targetType reflect.Type) (reflect.Value, error) {
	outVal := reflect.New(targetType).Elem()
	for i := 0; i < outVal.Len(); i++ {
		elemVal, err := h.unflattenParam(ctx, mem, ps, outVal.Index(i).Type())
		if err != nil {
			return reflect.Value{}, fmt.Errorf("failed to unflatten array element %d: %w", i, err)
		}
		outVal.Index(i).Set(elemVal)
	}
	return outVal, nil
}

// unflattenFlags 从参数流读取一个 uint64 并解码为 flags 结构体。
func (h *Host) unflattenFlags(ps *paramStream, targetType reflect.Type) (reflect.Value, error) {
	outVal := reflect.New(targetType).Elem()
	p, ok := ps.Next()
	if !ok {
		return reflect.Value{}, fmt.Errorf("not enough params on stack for flags %v", targetType)
	}
	// 按位解码布尔字段
	bitIdx := 0
	for i := 0; i < outVal.NumField(); i++ {
		sf := targetType.Field(i)
		if !isExportedField(sf) {
			continue
		}
		f := outVal.Field(i)
		if f.Kind() != reflect.Bool || !f.CanSet() {
			continue
		}
		if bitIdx < 64 && p&(1<<uint(bitIdx)) != 0 {
			f.SetBool(true)
		} else {
			f.SetBool(false)
		}
		bitIdx++
	}
	return outVal, nil
}

// unflattenVariant 从参数流读取判别值，然后读取对应载荷。
// 需要消费 maxPayloadLen 个参数（包括 padding）。
func (h *Host) unflattenVariant(ctx context.Context, mem api.Memory, ps *paramStream, targetType reflect.Type) (reflect.Value, error) {
	outVal := reflect.New(targetType).Elem()
	disc, ok := ps.Next()
	if !ok {
		return reflect.Value{}, fmt.Errorf("not enough params on stack for variant discriminant")
	}

	fields := variantFields(targetType)
	var active *reflect.StructField
	for i := range fields {
		if uint64(witCaseIndex(fields[i].Tag.Get("wit"), i)) == disc {
			f := fields[i]
			active = &f
			break
		}
	}
	if active == nil {
		return reflect.Value{}, fmt.Errorf("invalid variant discriminant %d for type %v", disc, targetType)
	}

	// 原先用 flattenParam(Zero) 探测各 case 扁平长度。
	// flattenParam 会走分配器（字符串/切片等），对零值也可能在 guest 里分配；
	// 这里只需要“形状有多长”，用 flattenType 按类型计算即可，无副作用。
	maxPayloadLen := 0
	for _, field := range fields {
		fieldType := field.Type
		if fieldType.Kind() == reflect.Pointer {
			fieldType = fieldType.Elem()
		}
		shape, err := flattenType(fieldType)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("could not get shape for variant case %s: %w", field.Name, err)
		}
		if len(shape) > maxPayloadLen {
			maxPayloadLen = len(shape)
		}
	}

	// 读取载荷
	payloadField := outVal.FieldByIndex(active.Index)
	wantType := payloadField.Type()
	if payloadField.Kind() == reflect.Pointer {
		if payloadField.IsNil() && payloadField.CanSet() {
			payloadField.Set(reflect.New(payloadField.Type().Elem()))
		}
		payloadField = payloadField.Elem()
		wantType = payloadField.Type()
	}

	startPos := 0
	if ps != nil {
		startPos = ps.pos
	}
	val, err := h.unflattenParam(ctx, mem, ps, wantType)
	if err != nil {
		return reflect.Value{}, err
	}
	if payloadField.CanSet() && val.IsValid() {
		payloadField.Set(val)
	}
	// 消费 padding（对齐到 maxPayloadLen）
	consumed := 0
	if ps != nil {
		consumed = ps.pos - startPos
	}
	for j := consumed; j < maxPayloadLen; j++ {
		ps.Next()
	}
	return outVal, nil
}
