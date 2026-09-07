package witgo

import (
	"context"
	"fmt"
	"reflect"

	"github.com/tetratelabs/wazero/api"
)

// ============================================================
// Lower：Guest 内存 -> Go 值
// ============================================================

// Lower 从 Guest 内存指针 ptr 处读取数据，填充到 Go 值 goVal。
// 自动选择合适的 lowerer 并执行读取。
//
// 参数:
//   - ctx:    上下文
//   - h:      Host 实例
//   - ptr:    Guest 内存中的读取起始指针
//   - goVal:  要填充的 Go 值（必须是可设置的 reflect.Value）
//
// 返回:
//   - error: 读取失败时的错误
//
// 示例:
//
//	var result MyStruct
//	err := Lower(ctx, host, ptr, reflect.ValueOf(&result).Elem())
func Lower(ctx context.Context, h *Host, ptr uint32, goVal reflect.Value) error {
	if !goVal.IsValid() {
		return fmt.Errorf("cannot lower into invalid value")
	}
	l, err := getOrGenerateLowerer(goVal.Type())
	if err != nil {
		return fmt.Errorf("failed to get lowerer for type %v: %w", goVal.Type(), err)
	}
	return l.lower(ctx, h, ptr, goVal)
}

// read 是实际读取的核心函数。
// 递归处理指针解引用，分发给各类型的读取逻辑。
func read(ctx context.Context, mem api.Memory, ptr uint32, val reflect.Value, layout *TypeLayout) error {
	if mem == nil {
		return fmt.Errorf("memory is nil")
	}
	// 解引用指针（空指针会自动分配）
	for val.Kind() == reflect.Pointer {
		if val.IsNil() {
			if !val.CanSet() {
				return fmt.Errorf("cannot allocate nil pointer of type %v", val.Type())
			}
			val.Set(reflect.New(val.Type().Elem()))
		}
		val = val.Elem()
	}
	if !val.CanSet() {
		return fmt.Errorf("value of type %v is not settable", val.Type())
	}

	typ := val.Type()
	if isVariant(typ) {
		return readVariant(ctx, mem, ptr, val, layout)
	}
	if isFlags(typ) {
		return readFlags(mem, ptr, val, layout)
	}

	switch val.Kind() {
	case reflect.Bool:
		b, ok := mem.ReadByte(ptr)
		if !ok {
			return fmt.Errorf("memory read failed for bool at ptr %d", ptr)
		}
		val.SetBool(b != 0)
		return nil
	case reflect.Int8:
		b, ok := mem.ReadByte(ptr)
		if !ok {
			return fmt.Errorf("memory read failed for int8 at ptr %d", ptr)
		}
		val.SetInt(int64(int8(b)))
		return nil
	case reflect.Uint8:
		b, ok := mem.ReadByte(ptr)
		if !ok {
			return fmt.Errorf("memory read failed for uint8 at ptr %d", ptr)
		}
		val.SetUint(uint64(b))
		return nil
	case reflect.Int16:
		s, ok := mem.ReadUint16Le(ptr)
		if !ok {
			return fmt.Errorf("memory read failed for int16 at ptr %d", ptr)
		}
		val.SetInt(int64(int16(s)))
		return nil
	case reflect.Uint16:
		s, ok := mem.ReadUint16Le(ptr)
		if !ok {
			return fmt.Errorf("memory read failed for uint16 at ptr %d", ptr)
		}
		val.SetUint(uint64(s))
		return nil
	case reflect.Int32:
		i, ok := mem.ReadUint32Le(ptr)
		if !ok {
			return fmt.Errorf("memory read failed for int32 at ptr %d", ptr)
		}
		val.SetInt(int64(int32(i)))
		return nil
	case reflect.Uint32:
		i, ok := mem.ReadUint32Le(ptr)
		if !ok {
			return fmt.Errorf("memory read failed for uint32 at ptr %d", ptr)
		}
		val.SetUint(uint64(i))
		return nil
	case reflect.Int64:
		i, ok := mem.ReadUint64Le(ptr)
		if !ok {
			return fmt.Errorf("memory read failed for int64 at ptr %d", ptr)
		}
		val.SetInt(int64(i))
		return nil
	case reflect.Uint64:
		i, ok := mem.ReadUint64Le(ptr)
		if !ok {
			return fmt.Errorf("memory read failed for uint64 at ptr %d", ptr)
		}
		val.SetUint(i)
		return nil
	case reflect.Int:
		if layout != nil && layout.Size == 4 {
			i, ok := mem.ReadUint32Le(ptr)
			if !ok {
				return fmt.Errorf("memory read failed for int at ptr %d", ptr)
			}
			val.SetInt(int64(int32(i)))
		} else {
			i, ok := mem.ReadUint64Le(ptr)
			if !ok {
				return fmt.Errorf("memory read failed for int at ptr %d", ptr)
			}
			val.SetInt(int64(i))
		}
		return nil
	case reflect.Uint:
		if layout != nil && layout.Size == 4 {
			i, ok := mem.ReadUint32Le(ptr)
			if !ok {
				return fmt.Errorf("memory read failed for uint at ptr %d", ptr)
			}
			val.SetUint(uint64(i))
		} else {
			i, ok := mem.ReadUint64Le(ptr)
			if !ok {
				return fmt.Errorf("memory read failed for uint at ptr %d", ptr)
			}
			val.SetUint(i)
		}
		return nil
	case reflect.Float32:
		f, ok := mem.ReadFloat32Le(ptr)
		if !ok {
			return fmt.Errorf("memory read failed for float32 at ptr %d", ptr)
		}
		val.SetFloat(float64(f))
		return nil
	case reflect.Float64:
		f, ok := mem.ReadFloat64Le(ptr)
		if !ok {
			return fmt.Errorf("memory read failed for float64 at ptr %d", ptr)
		}
		val.SetFloat(f)
		return nil
	case reflect.String:
		s, err := lowerString(mem, ptr)
		if err != nil {
			return err
		}
		val.SetString(s)
		return nil
	case reflect.Slice:
		return lowerSlice(ctx, mem, ptr, val)
	case reflect.Struct:
		return lowerStruct(ctx, mem, ptr, val, layout)
	case reflect.Array:
		return lowerArray(ctx, mem, ptr, val)
	default:
		return fmt.Errorf("unsupported type for lowering: %v", val.Kind())
	}
}

// readVariant 读取 variant，根据判别值选择 case 并填充。
func readVariant(ctx context.Context, mem api.Memory, ptr uint32, val reflect.Value, layout *TypeLayout) error {
	typ := val.Type()
	if layout == nil || !layout.IsSum {
		var err error
		layout, err = GetOrCalculateLayout(typ)
		if err != nil {
			return err
		}
	}
	discSize := layout.DiscSize
	if discSize == 0 {
		discSize = 1
	}

	// 读取判别值
	var disc uint32
	var ok bool
	switch discSize {
	case 1:
		var b byte
		b, ok = mem.ReadByte(ptr)
		disc = uint32(b)
	case 2:
		var s uint16
		s, ok = mem.ReadUint16Le(ptr)
		disc = uint32(s)
	case 4:
		disc, ok = mem.ReadUint32Le(ptr)
	default:
		return fmt.Errorf("unsupported variant discriminant size: %d", discSize)
	}
	if !ok {
		return fmt.Errorf("failed to read variant discriminant at ptr %d", ptr)
	}

	// 根据判别值找到对应的 case
	fields := variantFields(typ)
	var active *reflect.StructField
	for i := range fields {
		if uint32(witCaseIndex(fields[i].Tag.Get("wit"), i)) == disc {
			f := fields[i]
			active = &f
			break
		}
	}
	if active == nil {
		return fmt.Errorf("invalid variant discriminant %d for type %v", disc, typ)
	}

	// 重置所有字段
	for i := range fields {
		fv := val.FieldByIndex(fields[i].Index)
		if !fv.CanSet() || fields[i].Name == active.Name {
			continue
		}
		fv.Set(reflect.Zero(fv.Type()))
	}

	// 读取载荷
	payloadField := val.FieldByIndex(active.Index)
	if payloadField.Kind() == reflect.Pointer {
		if payloadField.IsNil() {
			if !payloadField.CanSet() {
				return fmt.Errorf("cannot set payload field %s", active.Name)
			}
			payloadField.Set(reflect.New(payloadField.Type().Elem()))
		}
		payloadField = payloadField.Elem()
	}
	cl, err := GetOrCalculateLayout(payloadField.Type())
	if err != nil {
		return err
	}
	p, err := addU32(ptr, layout.PayloadOffset)
	if err != nil {
		return err
	}
	return read(ctx, mem, p, payloadField, cl)
}

// readFlags 读取 flags 位图并设置结构体的布尔字段。
func readFlags(mem api.Memory, ptr uint32, val reflect.Value, layout *TypeLayout) error {
	if layout == nil {
		var err error
		layout, err = GetOrCalculateLayout(val.Type())
		if err != nil {
			return err
		}
	}
	var bits uint64
	var ok bool
	switch layout.Size {
	case 1:
		var b byte
		b, ok = mem.ReadByte(ptr)
		bits = uint64(b)
	case 2:
		var s uint16
		s, ok = mem.ReadUint16Le(ptr)
		bits = uint64(s)
	case 4:
		var i uint32
		i, ok = mem.ReadUint32Le(ptr)
		bits = uint64(i)
	case 8:
		bits, ok = mem.ReadUint64Le(ptr)
	default:
		return fmt.Errorf("unsupported flags size: %d", layout.Size)
	}
	if !ok {
		return fmt.Errorf("memory read failed for flags at ptr %d", ptr)
	}
	// 按位设置布尔字段
	bitIdx := 0
	for i := 0; i < val.NumField(); i++ {
		sf := val.Type().Field(i)
		if !isExportedField(sf) {
			continue
		}
		f := val.Field(i)
		if f.Kind() != reflect.Bool || !f.CanSet() {
			continue
		}
		if bitIdx >= 64 {
			f.SetBool(false)
		} else {
			f.SetBool(bits&(1<<uint(bitIdx)) != 0)
		}
		bitIdx++
	}
	return nil
}

// lowerString 从 guest 内存读取字符串（header + content）。
// header 格式: {content_ptr(uint32), length(uint32)}
func lowerString(mem api.Memory, ptr uint32) (string, error) {
	if mem == nil {
		return "", fmt.Errorf("memory is nil")
	}
	contentPtr, ok1 := mem.ReadUint32Le(ptr)
	contentLen, ok2 := mem.ReadUint32Le(ptr + 4)
	if !ok1 || !ok2 {
		return "", fmt.Errorf("failed to read string header at ptr %d", ptr)
	}
	if contentLen == 0 {
		return "", nil
	}
	return LowerStringFromParts(mem, contentPtr, contentLen)
}

// LowerStringFromParts 从线性内存的 (ptr,len) 读取字符串。
func LowerStringFromParts(mem api.Memory, ptr, length uint32) (str string, err error) {
	if length == 0 {
		return "", nil
	}
	if mem == nil {
		return "", fmt.Errorf("memory is nil")
	}
	content, ok := mem.Read(ptr, length)
	if !ok {
		return "", fmt.Errorf("failed to read string content at ptr %d with length %d", ptr, length)
	}
	return string(content), nil
}

// lowerSlice 读取切片 header，然后调用 lowerSlice2 填充 Go 切片。
func lowerSlice(ctx context.Context, mem api.Memory, ptr uint32, val reflect.Value) error {
	if mem == nil {
		return fmt.Errorf("memory is nil")
	}
	contentPtr, ok1 := mem.ReadUint32Le(ptr)
	contentLen, ok2 := mem.ReadUint32Le(ptr + 4)
	if !ok1 || !ok2 {
		return fmt.Errorf("failed to read slice header at ptr %d", ptr)
	}
	return lowerSlice2(ctx, mem, contentPtr, contentLen, val)
}

// lowerSlice2 从内容指针和长度填充 Go 切片。
func lowerSlice2(ctx context.Context, mem api.Memory, contentPtr uint32, contentLen uint32, val reflect.Value) error {
	if !val.CanSet() {
		return fmt.Errorf("slice value is not settable")
	}
	elemType := val.Type().Elem()
	if contentLen == 0 {
		val.Set(reflect.MakeSlice(val.Type(), 0, 0))
		return nil
	}
	if !sliceLenFitsInt(contentLen) {
		return fmt.Errorf("slice length %d exceeds platform int range", contentLen)
	}

	if isExactByteSlice(elemType) {
		content, err := LowerSliceFromParts(mem, contentPtr, contentLen)
		if err != nil {
			return err
		}
		val.Set(reflect.ValueOf(content).Convert(val.Type()))
		return nil
	}

	elemLayout, err := GetOrCalculateLayout(elemType)
	if err != nil {
		return fmt.Errorf("failed to get element layout for slice: %w", err)
	}
	stride := align(elemLayout.Size, elemLayout.Alignment)
	if stride == 0 {
		stride = elemLayout.Size
	}

	n := int(contentLen)
	newSlice := reflect.MakeSlice(val.Type(), n, n)
	for i := 0; i < n; i++ {
		elemPtr, err := addU32(contentPtr, uint32(i)*stride)
		if err != nil {
			return err
		}
		if err := read(ctx, mem, elemPtr, newSlice.Index(i), elemLayout); err != nil {
			return fmt.Errorf("failed to read slice element %d: %w", i, err)
		}
	}
	val.Set(newSlice)
	return nil
}

// LowerSliceFromParts 从线性内存读取 []byte 并返回独立副本。
func LowerSliceFromParts(mem api.Memory, ptr, length uint32) ([]byte, error) {
	if length == 0 {
		return []byte{}, nil
	}
	if mem == nil {
		return nil, fmt.Errorf("memory is nil")
	}
	content, ok := mem.Read(ptr, length)
	if !ok {
		return nil, fmt.Errorf("failed to read slice content at ptr %d with length %d", ptr, length)
	}
	result := make([]byte, length)
	copy(result, content)
	return result, nil
}

// lowerArray 读取数组元素。
func lowerArray(ctx context.Context, mem api.Memory, ptr uint32, val reflect.Value) error {
	elemLayout, err := GetOrCalculateLayout(val.Type().Elem())
	if err != nil {
		return fmt.Errorf("failed to get array element layout: %w", err)
	}
	currentOffset := ptr
	for i := 0; i < val.Len(); i++ {
		elemPtr := align(currentOffset, elemLayout.Alignment)
		if err := read(ctx, mem, elemPtr, val.Index(i), elemLayout); err != nil {
			return fmt.Errorf("failed to read array element %d: %w", i, err)
		}
		currentOffset, err = addU32(elemPtr, elemLayout.Size)
		if err != nil {
			return err
		}
	}
	return nil
}

// lowerStruct 读取结构体字段。
func lowerStruct(ctx context.Context, mem api.Memory, ptr uint32, val reflect.Value, layout *TypeLayout) error {
	if layout == nil {
		var err error
		layout, err = GetOrCalculateLayout(val.Type())
		if err != nil {
			return err
		}
	}
	for _, fieldLayout := range layout.Fields {
		fieldVal := val.FieldByIndex(fieldLayout.StructField.Index)
		if !fieldVal.IsValid() || !fieldVal.CanSet() {
			continue
		}
		fieldPtr, err := addU32(ptr, fieldLayout.Offset)
		if err != nil {
			return err
		}
		if err := read(ctx, mem, fieldPtr, fieldVal, fieldLayout.Layout); err != nil {
			return fmt.Errorf("failed to read struct field %s: %w", fieldLayout.StructField.Name, err)
		}
	}
	return nil
}
