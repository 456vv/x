package witgo

import (
	"context"
	"fmt"
	"math"
	"reflect"

	"github.com/tetratelabs/wazero/api"
)

var errPtrOverflow = fmt.Errorf("guest pointer arithmetic overflow")

// ============================================================
// Lift：Go 值 -> Guest 内存
// ============================================================

// liftByWrite 根据类型布局分配内存并写入 val，返回分配的内存指针。
// 这是 most lifter 的实际实现，负责分配内存并调用 write 写入数据。
//
// 参数:
//   - ctx:    上下文
//   - h:      Host 实例
//   - val:    要提升的 Go 值
//   - layout: 类型布局（为 nil 时自动计算）
//
// 返回:
//   - uint32: 分配的 Guest 内存指针
//   - error:  分配或写入失败时的错误
func liftByWrite(ctx context.Context, h *Host, val reflect.Value, layout *TypeLayout) (uint32, error) {
	if h == nil || h.allocator == nil || h.module == nil {
		return 0, fmt.Errorf("invalid host")
	}
	mem := h.module.Memory()
	if mem == nil {
		return 0, fmt.Errorf("module memory is nil")
	}
	if layout == nil {
		var err error
		layout, err = GetOrCalculateLayout(val.Type())
		if err != nil {
			return 0, err
		}
	}
	// 分配 Guest 内存
	ptr, err := h.allocator.Allocate(ctx, layout.Size, layout.Alignment)
	if err != nil {
		return 0, err
	}
	// 写入数据
	if err := write(ctx, mem, h.allocator, val, ptr, layout); err != nil {
		// 写入失败时回滚分配
		_ = h.allocator.Free(ctx, ptr, layout.Size, layout.Alignment)
		return 0, err
	}
	return ptr, nil
}

// Lift 将 Go 值 val 提升到 Guest 内存，并返回指向该内存的指针。
// 自动选择合适的 lifter 并执行提升。
//
// 参数:
//   - ctx: 上下文
//   - h:   Host 实例
//   - val: 要提升的 Go 值
//
// 返回:
//   - uint32: Guest 内存指针
//   - error:  提升失败时的错误
//
// 示例:
//
//	ptr, err := Lift(ctx, host, reflect.ValueOf(myStruct))
func Lift(ctx context.Context, h *Host, val reflect.Value) (uint32, error) {
	if !val.IsValid() {
		return 0, fmt.Errorf("cannot lift invalid value")
	}
	layout, err := GetOrCalculateLayout(val.Type())
	if err != nil {
		return 0, err
	}
	l, err := getOrGenerateLifter(val.Type())
	if err != nil {
		return 0, fmt.Errorf("failed to get lifter for type %v: %w", val.Type(), err)
	}
	return l.lift(ctx, h, val, layout)
}

// LiftToPtr 将 Go 值 val 写入已经分配好的 Guest 内存指针 ptr 处。
// 与 Lift 不同，此函数不分配新内存，而是写入已有指针。
// 用于 retptr 模式（返回值写入调用方提供的缓冲区）。
//
// 参数:
//   - mem:  Guest 内存接口
//   - alloc: 分配器（用于字符串/切片等需要额外分配的情况）
//   - val:   要写入的 Go 值
//   - ptr:   目标 Guest 内存指针
//
// 返回:
//   - error: 写入失败时的错误
func LiftToPtr(ctx context.Context, mem api.Memory, alloc *GuestAllocator, val reflect.Value, ptr uint32) error {
	if mem == nil || alloc == nil {
		return fmt.Errorf("memory/allocator cannot be nil")
	}
	if !val.IsValid() {
		return fmt.Errorf("cannot lift invalid value")
	}
	layout, err := GetOrCalculateLayout(val.Type())
	if err != nil {
		return err
	}
	return write(ctx, mem, alloc, val, ptr, layout)
}

// write 是实际写入的核心函数，处理各种类型。
// 递归处理指针解引用，分发给各类型的写入逻辑。
//
// 参数:
//   - ctx:    上下文
//   - mem:    Guest 内存
//   - alloc:  分配器
//   - val:    要写入的 Go 值
//   - ptr:    目标 Guest 内存指针
//   - layout: 类型布局
//
// 返回:
//   - error: 写入失败时的错误
func write(ctx context.Context, mem api.Memory, alloc *GuestAllocator, val reflect.Value, ptr uint32, layout *TypeLayout) error {
	if mem == nil {
		return fmt.Errorf("memory is nil")
	}
	// 解引用指针
	for val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return fmt.Errorf("cannot write nil pointer")
		}
		val = val.Elem()
	}

	typ := val.Type()
	if isVariant(typ) {
		return writeVariant(ctx, mem, alloc, val, ptr, layout)
	}
	if isFlags(typ) {
		return writeFlags(mem, val, ptr, layout)
	}

	switch val.Kind() {
	case reflect.Bool:
		return checkWrite(mem.WriteByte(ptr, boolToByte(val.Bool())), "bool", ptr)
	case reflect.Int8:
		return checkWrite(mem.WriteByte(ptr, byte(val.Int())), "int8", ptr)
	case reflect.Uint8:
		return checkWrite(mem.WriteByte(ptr, byte(val.Uint())), "uint8", ptr)
	case reflect.Int16:
		return checkWrite(mem.WriteUint16Le(ptr, uint16(val.Int())), "int16", ptr)
	case reflect.Uint16:
		return checkWrite(mem.WriteUint16Le(ptr, uint16(val.Uint())), "uint16", ptr)
	case reflect.Int32:
		return checkWrite(mem.WriteUint32Le(ptr, uint32(val.Int())), "int32", ptr)
	case reflect.Uint32:
		return checkWrite(mem.WriteUint32Le(ptr, uint32(val.Uint())), "uint32", ptr)
	case reflect.Int64:
		return checkWrite(mem.WriteUint64Le(ptr, uint64(val.Int())), "int64", ptr)
	case reflect.Uint64:
		return checkWrite(mem.WriteUint64Le(ptr, val.Uint()), "uint64", ptr)
	case reflect.Int:
		if layout != nil && layout.Size == 4 {
			return checkWrite(mem.WriteUint32Le(ptr, uint32(int32(val.Int()))), "int(32)", ptr)
		}
		return checkWrite(mem.WriteUint64Le(ptr, uint64(val.Int())), "int(64)", ptr)
	case reflect.Uint:
		if layout != nil && layout.Size == 4 {
			return checkWrite(mem.WriteUint32Le(ptr, uint32(val.Uint())), "uint(32)", ptr)
		}
		return checkWrite(mem.WriteUint64Le(ptr, val.Uint()), "uint(64)", ptr)
	case reflect.Float32:
		return checkWrite(mem.WriteFloat32Le(ptr, float32(val.Float())), "float32", ptr)
	case reflect.Float64:
		return checkWrite(mem.WriteFloat64Le(ptr, val.Float()), "float64", ptr)
	case reflect.String:
		return liftString(ctx, mem, alloc, val.String(), ptr)
	case reflect.Slice:
		return liftSlice(ctx, mem, alloc, val, ptr)
	case reflect.Struct:
		return liftStruct(ctx, mem, alloc, val, ptr, layout)
	case reflect.Array:
		return liftArray(ctx, mem, alloc, val, ptr)
	default:
		return fmt.Errorf("unsupported type for lifting: %v", val.Kind())
	}
}

func checkWrite(ok bool, kind string, ptr uint32) error {
	if !ok {
		return fmt.Errorf("memory write failed for %s at ptr %d", kind, ptr)
	}
	return nil
}

// writeVariant 写入 variant，先写判别值，再写载荷。
// 判别值大小根据 case 数量确定（1/2/4 字节），载荷从 PayloadOffset 开始写入。
func writeVariant(ctx context.Context, mem api.Memory, alloc *GuestAllocator, val reflect.Value, ptr uint32, layout *TypeLayout) error {
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
	fields := variantFields(typ)

	// 找到激活的 case
	active := -1
	for i, f := range fields {
		if !val.FieldByIndex(f.Index).IsZero() {
			active = i
			break
		}
	}
	if active < 0 {
		if isOption(typ) && len(fields) > 0 {
			active = 0 // None case
		} else {
			return fmt.Errorf("invalid variant: no case set for type %v", typ)
		}
	}

	// 写入判别值
	disc := uint32(witCaseIndex(fields[active].Tag.Get("wit"), active))
	switch discSize {
	case 1:
		if !mem.WriteByte(ptr, byte(disc)) {
			return fmt.Errorf("failed to write variant discriminant (u8)")
		}
	case 2:
		if !mem.WriteUint16Le(ptr, uint16(disc)) {
			return fmt.Errorf("failed to write variant discriminant (u16)")
		}
	case 4:
		if !mem.WriteUint32Le(ptr, disc) {
			return fmt.Errorf("failed to write variant discriminant (u32)")
		}
	default:
		return fmt.Errorf("unsupported variant discriminant size: %d", discSize)
	}

	// 写入载荷
	payloadField := val.FieldByIndex(fields[active].Index)
	if payloadField.Kind() == reflect.Pointer {
		if payloadField.IsNil() {
			return nil
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
	return write(ctx, mem, alloc, payloadField, p, cl)
}

// writeFlags 将 flags 的布尔字段编码为位图写入内存。
// 每个布尔字段对应一位，按字段顺序排列。
func writeFlags(mem api.Memory, val reflect.Value, ptr uint32, layout *TypeLayout) error {
	if layout == nil {
		var err error
		layout, err = GetOrCalculateLayout(val.Type())
		if err != nil {
			return err
		}
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
	switch layout.Size {
	case 1:
		return check(mem.WriteByte(ptr, byte(bits)))
	case 2:
		return check(mem.WriteUint16Le(ptr, uint16(bits)))
	case 4:
		return check(mem.WriteUint32Le(ptr, uint32(bits)))
	case 8:
		return check(mem.WriteUint64Le(ptr, bits))
	default:
		return fmt.Errorf("unsupported flags size: %d", layout.Size)
	}
}

// liftString 将字符串写入 guest 内存（header + content）。
// header 在 ptr 处: {content_ptr(uint32), length(uint32)}
// content 在 content_ptr 处: 原始字节
func liftString(ctx context.Context, mem api.Memory, alloc *GuestAllocator, s string, ptr uint32) error {
	if mem == nil || alloc == nil {
		return fmt.Errorf("memory/allocator cannot be nil")
	}
	if uint64(len(s)) > math.MaxUint32 {
		return fmt.Errorf("string length exceeds 32-bit limit: %d", len(s))
	}
	length := uint32(len(s))
	var contentPtr uint32
	if length != 0 {
		var err error
		contentPtr, err = alloc.Allocate(ctx, length, 1)
		if err != nil {
			return fmt.Errorf("failed to allocate string content (len=%d): %w", length, err)
		}
		if !mem.Write(contentPtr, []byte(s)) {
			_ = alloc.Free(ctx, contentPtr, length, 1)
			return fmt.Errorf("failed to write string content at ptr %d (len=%d)", contentPtr, length)
		}
	}
	if !mem.WriteUint32Le(ptr, contentPtr) || !mem.WriteUint32Le(ptr+4, length) {
		if contentPtr != 0 {
			alloc.Free(ctx, contentPtr, length, 1) // 未赋值 Free 触发 ineffectual-assignment 检查
		}
		return fmt.Errorf("failed to write string header back at ptr %d", ptr)
	}
	return nil
}

// liftSlice 将切片写入 guest 内存（header + 元素内容）。
func liftSlice(ctx context.Context, mem api.Memory, alloc *GuestAllocator, val reflect.Value, ptr uint32) error {
	if mem == nil || alloc == nil {
		return fmt.Errorf("memory/allocator cannot be nil")
	}
	sliceLen := val.Len()
	elemType := val.Type().Elem()
	if sliceLen == 0 {
		if !mem.WriteUint32Le(ptr, 0) || !mem.WriteUint32Le(ptr+4, 0) {
			return fmt.Errorf("failed to write empty slice header at ptr %d", ptr)
		}
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
	if stride == 0 {
		return fmt.Errorf("invalid zero stride for slice element %v", elemType)
	}
	if uint64(sliceLen) > uint64(math.MaxUint32)/uint64(stride) {
		return fmt.Errorf("slice content size overflow: len=%d stride=%d", sliceLen, stride)
	}

	contentSize := uint32(sliceLen-1)*stride + elemLayout.Size
	contentPtr, err := alloc.Allocate(ctx, contentSize, elemLayout.Alignment)
	if err != nil {
		return fmt.Errorf("failed to allocate slice content (len=%d, size=%d): %w", sliceLen, contentSize, err)
	}

	if isExactByteSlice(elemType) {
		if !mem.Write(contentPtr, val.Bytes()) {
			alloc.Free(ctx, contentPtr, contentSize, elemLayout.Alignment)
			return fmt.Errorf("failed to write []byte slice content at ptr %d", contentPtr)
		}
	} else {
		for i := 0; i < sliceLen; i++ {
			elemPtr, err := addU32(contentPtr, uint32(i)*stride)
			if err != nil {
				alloc.Free(ctx, contentPtr, contentSize, elemLayout.Alignment)
				return err
			}
			if err := write(ctx, mem, alloc, val.Index(i), elemPtr, elemLayout); err != nil {
				alloc.Free(ctx, contentPtr, contentSize, elemLayout.Alignment)
				return fmt.Errorf("failed to write slice element %d: %w", i, err)
			}
		}
	}

	if !mem.WriteUint32Le(ptr, contentPtr) || !mem.WriteUint32Le(ptr+4, uint32(sliceLen)) {
		alloc.Free(ctx, contentPtr, contentSize, elemLayout.Alignment) // header 写失败泄漏元素缓冲
		return fmt.Errorf("failed to write slice header back at ptr %d", ptr)
	}
	return nil
}

// liftArray 将数组连续写入内存。
func liftArray(ctx context.Context, mem api.Memory, alloc *GuestAllocator, val reflect.Value, ptr uint32) error {
	elemLayout, err := GetOrCalculateLayout(val.Type().Elem())
	if err != nil {
		return fmt.Errorf("failed to get array element layout: %w", err)
	}
	currentOffset := ptr
	for i := 0; i < val.Len(); i++ {
		elemPtr := align(currentOffset, elemLayout.Alignment)
		if err := write(ctx, mem, alloc, val.Index(i), elemPtr, elemLayout); err != nil {
			return fmt.Errorf("failed to write array element %d: %w", i, err)
		}
		currentOffset, err = addU32(elemPtr, elemLayout.Size)
		if err != nil {
			return err
		}
	}
	return nil
}

// liftStruct 将结构体字段按布局依次写入。
func liftStruct(ctx context.Context, mem api.Memory, alloc *GuestAllocator, val reflect.Value, ptr uint32, layout *TypeLayout) error {
	if layout == nil {
		var err error
		layout, err = GetOrCalculateLayout(val.Type())
		if err != nil {
			return err
		}
	}
	for _, fieldLayout := range layout.Fields {
		fieldVal := val.FieldByIndex(fieldLayout.StructField.Index)
		if !fieldVal.IsValid() {
			continue
		}
		fieldPtr, err := addU32(ptr, fieldLayout.Offset)
		if err != nil {
			return err
		}
		if err := write(ctx, mem, alloc, fieldVal, fieldPtr, fieldLayout.Layout); err != nil {
			return fmt.Errorf("failed to write struct field %s: %w", fieldLayout.StructField.Name, err)
		}
	}
	return nil
}

func boolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}

func check(ok bool) error {
	if !ok {
		return fmt.Errorf("memory access failed")
	}
	return nil
}
