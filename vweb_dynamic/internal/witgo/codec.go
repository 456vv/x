package witgo

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

// ============================================================
// Lifter / Lowerer 接口：Go <-> Guest 内存的数据转换
// ============================================================
//
// 整体架构说明：
//   - "提升 (lift)"：把 Go 值编码(序列化)到 Guest(Wasm) 内存，返回其在 Guest 侧的指针。
//   - "降低 (lower)"：从 Guest 内存按指针读取数据，解码(反序列化)成 Go 值。
//   - 一个 reflect.Type 对应一套 codec(lifter+lowerer)，由本文件底部的工厂函数按类型生成，
//     结果缓存在 lifterCache / lowererCache 中，避免重复生成带来的反射开销。

// lifter 接口：将 Go 值 "提升"(写入)到 Guest 内存。
// 由实现方负责把 Go 数据结构按 TypeLayout 描述的内存布局序列化到 Guest 内存，
// 一般流程是：先向 Guest 分配内存 -> 按布局写入各字段/元素 -> 返回起始指针。
//
// 方法参数:
//   - ctx   : context.Context，用于取消与分配作用域控制(贯穿到 allocator)。
//   - h     : *Host，宿主实例，提供 allocator(分配器)与 module(Wasm 模块)引用。
//   - val   : reflect.Value，要被提升的 Go 值；需是 val.CanInterface() 可读的值。
//   - layout: *TypeLayout，该类型的对齐/偏移等内存布局描述；nil 时部分实现会自行计算。
//
// 返回值:
//   - uint32: 提升后数据在 Guest 内存中的起始地址(0 表示空指针/失败占位)。
//   - error : 失败时非 nil，例如：unsupported type、分配失败、nil 指针等。
//
// 调用示例:
//
//	lf, _ := getOrGenerateLifter(reflect.TypeOf("hello"))   // 获取 string 的 lifter
//	ptr, err := lf.lift(ctx, host, reflect.ValueOf("hello"), nil)
//	// 成功: ptr = 0x10420 (guest 内存中该字符串所在地址), err = nil
//	// 失败: ptr = 0,           err = <具体错误信息>
type lifter interface {
	lift(ctx context.Context, h *Host, val reflect.Value, layout *TypeLayout) (uint32, error)
}

// lowerer 接口：从 Guest 内存 "降低"(读取)数据到 Go 值。
// 与 lifter 相反：根据 Guest 内存指针与布局信息，把二进制内存反序列化填充进 outVal。
//
// 方法参数:
//   - ctx   : context.Context，贯穿到底层 read 操作。
//   - h     : *Host，宿主实例，经 h.module.Memory() 取 Guest 线性内存。
//   - ptr   : uint32，Guest 内存中的读取起始地址(0 通常代表空引用)。
//   - outVal: reflect.Value，接收数据的 Go 值；**必须可设置(CanSet)**，
//     否则实现内部会返回 "outVal is not settable" 类错误。
//
// 返回值:
//   - error: 成功为 nil；失败非 nil，例如读取越界、类型不匹配、outVal 不可写等。
//
// 调用示例:
//
//	lr, _ := getOrGenerateLowerer(reflect.TypeOf(""))
//	out := new(string)                                   // *string 可解引用设置
//	err := lr.lower(ctx, host, 0x10420, reflect.ValueOf(out).Elem())
//	// 成功: err = nil 且 *out == "hello"(读回的内容)
//	// 失败: err = <具体错误信息>
type lowerer interface {
	lower(ctx context.Context, h *Host, ptr uint32, outVal reflect.Value) error
}

// 全局 codec 缓存：以 reflect.Type 为 key，避免对同一类型反复生成 lifter/lowerer。
// 用 sync.Map 保证并发安全(Go <-> Guest 编解码常发生在多个 goroutine 调用 Wasm 时)。
var (
	// lifterCache 缓存 reflect.Type -> lifter，按类型缓存"提升器"。
	lifterCache sync.Map
	// lowererCache 缓存 reflect.Type -> lowerer，按类型缓存"降低器"。
	lowererCache sync.Map
)

// getOrGenerateLifter 获取或生成类型 typ 的 lifter，并附带循环引用检测。
// 对外统一入口：先查全局缓存，命中直接返回；未命中则递归生成，成功后写回缓存。
//
// 参数:
//   - typ: 需要提升的 Go 类型(reflect.Type)，不能为 nil。
//
// 返回值:
//   - lifter: 该类型对应的提升器；已生成过一次后每次调用返回同一缓存实例。
//   - error : 生成失败时非 nil，例如 typ 为 nil 或不支持的类型。
//
// 调用示例:
//
//	lf, err := getOrGenerateLifter(reflect.TypeOf([]int32{1, 2})) // 切片类型
//	// 成功: lf = *witgo.sliceLifter, err = nil
//	// 之后每次相同类型调用都会命中缓存，直接返回，不再重新生成。
func getOrGenerateLifter(typ reflect.Type) (lifter, error) {
	return getOrGenerateLifterStack(typ, make(map[reflect.Type]struct{}))
}

// getOrGenerateLifterStack getOrGenerateLifter 的内部实现。
// 多带一个 stack(当前展开路径上的类型集合)用于检测循环引用，
// 防止 A->B->A 这样的自引用/互引用结构体在生成时无限递归。
//
// 参数:
//   - typ  : 需要生成 lifter 的类型。
//   - stack: 记录"正在生成中"的类型集合；命中即代表出现递归，报错。
//
// 返回值:
//   - lifter: 生成(或从缓存命中)的提升器。
//   - error : nil 类型、缓存条目损坏、循环类型、不支持类型都会返回错误。
func getOrGenerateLifterStack(typ reflect.Type, stack map[reflect.Type]struct{}) (lifter, error) {
	if typ == nil {
		return nil, fmt.Errorf("cannot generate lifter for nil type")
	}
	// 1. 先查缓存：命中则直接复用，避免重复做昂贵的类型递归展开。
	if l, ok := lifterCache.Load(typ); ok {
		lf, ok := l.(lifter)
		if !ok {
			// 理论不可能发生：缓存里存的本就该是 lifter。防御性校验。
			return nil, fmt.Errorf("invalid lifter cache entry for %v", typ)
		}
		return lf, nil
	}
	// 2. 循环检测：若该类型已在本条展开路径上，说明存在递归类型，暂不支持。
	if _, cycling := stack[typ]; cycling {
		return nil, fmt.Errorf("recursive type not supported for lifter generation: %v", typ)
	}
	// 3. 把当前类型压栈，开始生成；defer 负责出栈，
	//    这样同一类型可在不同分支(不同路径)合法出现，只在单一路径内判环。
	stack[typ] = struct{}{}
	defer delete(stack, typ)

	l, err := generateLifter(typ, stack)
	if err != nil {
		return nil, err
	}
	// 4. 写入缓存。若并发下别的 goroutine 已抢先写入，则采用先写入者，
	//    保证同一类型全局只有一份 lifter(幂等)。
	actual, _ := lifterCache.LoadOrStore(typ, l)
	if lf, ok := actual.(lifter); ok {
		return lf, nil
	}
	return l, nil
}

// getOrGenerateLowerer 获取或生成类型 typ 的 lowerer。
// 逻辑与 getOrGenerateLifter 完全对称，只是缓存和生成目标换成 lowerer。
//
// 参数:
//   - typ: 需要从 Guest 内存读取的 Go 类型(reflect.Type)。
//
// 返回值:
//   - lowerer: 该类型的降低器(全局缓存单例)。
//   - error  : 生成失败时非 nil。
//
// 调用示例:
//
//	lr, err := getOrGenerateLowerer(reflect.TypeOf(MyStruct{}))
//	// 成功: lr = *witgo.structLowerer, err = nil
func getOrGenerateLowerer(typ reflect.Type) (lowerer, error) {
	return getOrGenerateLowererStack(typ, make(map[reflect.Type]struct{}))
}

// getOrGenerateLowererStack getOrGenerateLowerer 的内部实现(带循环引用检测)。
// 流程与 getOrGenerateLifterStack 相同：查缓存 -> 判环 -> 生成 -> 写缓存。
func getOrGenerateLowererStack(typ reflect.Type, stack map[reflect.Type]struct{}) (lowerer, error) {
	if typ == nil {
		return nil, fmt.Errorf("cannot generate lowerer for nil type")
	}
	// 查缓存命中直接返回
	if l, ok := lowererCache.Load(typ); ok {
		lr, ok := l.(lowerer)
		if !ok {
			return nil, fmt.Errorf("invalid lowerer cache entry for %v", typ)
		}
		return lr, nil
	}
	// 循环引用检测
	if _, cycling := stack[typ]; cycling {
		return nil, fmt.Errorf("recursive type not supported for lowerer generation: %v", typ)
	}
	stack[typ] = struct{}{}
	defer delete(stack, typ)

	l, err := generateLowerer(typ, stack)
	if err != nil {
		return nil, err
	}
	// LoadOrStore 保证并发下同一类型只有一份实例
	actual, _ := lowererCache.LoadOrStore(typ, l)
	if lr, ok := actual.(lowerer); ok {
		return lr, nil
	}
	return l, nil
}

// ============================================================
// Lifter / Lowerer 生成器(类型分发工厂)
// ============================================================

// generateLifter 依据类型的 Kind 及特殊标记(variant/flags)，分派生成对应的 lifter 实现。
// 这是"类型 -> 具体 codec"的分发中心：
//   - 先判断是否是 variant(option/result 联合体) 或 flags(位图) 这类 wit 语义特殊类型；
//   - 再按 reflect.Kind 判断基本类型 / string / slice / struct / array / pointer；
//   - 其余一概报 unsupported，做到生成期即失败(fail-fast)。
//
// 参数:
//   - typ  : Go 类型。
//   - stack: 循环引用检测栈(会透传给递归的构造函数)。
//
// 返回值:
//   - lifter: 与该类型匹配的具体 lifter 实例。
//   - error : 类型不受支持时非 nil。
//
// 调用示例:
//
//	lf, err := generateLifter(reflect.TypeOf(uint8(7)), emptyStack)
//	// 成功: lf = &witgo.primitiveLifter{}
//	lf, err = generateLifter(reflect.TypeOf(map[string]int{}), emptyStack)
//	// 失败: err = "unsupported type for lifter generation: map[string]int"
func generateLifter(typ reflect.Type, stack map[reflect.Type]struct{}) (lifter, error) {
	switch {
	case isVariant(typ): // wit 的 variant / option / result 语义(见 variantFields 等辅助函数)
		return newVariantLifter(typ, stack)
	case isFlags(typ): // wit 的 flags(布尔集合按位编码)语义
		return newFlagsLifter(), nil
	}

	switch typ.Kind() {
	case reflect.Bool, reflect.Int8, reflect.Uint8, reflect.Int16, reflect.Uint16,
		reflect.Int32, reflect.Uint32, reflect.Int64, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Int, reflect.Uint:
		return newPrimitiveLifter(), nil // 定长标量：无需分配，直接按字节写
	case reflect.String:
		return &stringLifter{}, nil // 字符串：guest 侧表示为 {ptr,len}，需分配正文内存
	case reflect.Slice:
		return newSliceLifter(typ, stack) // 切片：guest 侧为 {ptr,len,cap} header + 连续元素
	case reflect.Struct:
		return newStructLifter(typ, stack) // 结构体：逐字段按其偏移写入
	case reflect.Array:
		return newArrayLifter(typ, stack) // 定长数组：元素连续存放，无需 header
	case reflect.Pointer:
		return newPointerLifter(typ, stack) // 指针：解引用后提升被指向的值
	default:
		return nil, fmt.Errorf("unsupported type for lifter generation: %v", typ)
	}
}

// generateLowerer 依据类型生成对应的 lowerer 实现，与 generateLifter 对称。
// 多了一步：先 GetOrCalculateLayout(typ) 求取该类型的布局(lower 侧常需要布局
// 才能定位字段读取，且 struct lowerer 会把 layout 缓存下来复用)。
//
// 参数:
//   - typ  : Go 类型。
//   - stack: 循环引用检测栈。
//
// 返回值:
//   - lowerer: 匹配的降低器实例。
//   - error  : 布局计算失败或类型不受支持时非 nil。
//
// 调用示例:
//
//	lr, err := generateLowerer(reflect.TypeOf(uint8(0)), emptyStack)
//	// 成功: lr = &witgo.primitiveLowerer{}
func generateLowerer(typ reflect.Type, stack map[reflect.Type]struct{}) (lowerer, error) {
	layout, err := GetOrCalculateLayout(typ)
	if err != nil {
		return nil, err
	}

	switch {
	case isVariant(typ):
		return newVariantLowerer(typ, stack)
	case isFlags(typ):
		return newFlagsLowerer(), nil
	}

	switch typ.Kind() {
	case reflect.Bool, reflect.Int8, reflect.Uint8, reflect.Int16, reflect.Uint16,
		reflect.Int32, reflect.Uint32, reflect.Int64, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Int, reflect.Uint:
		return newPrimitiveLowerer(), nil
	case reflect.String:
		return &stringLowerer{}, nil
	case reflect.Slice:
		return newSliceLowerer(typ, stack)
	case reflect.Struct:
		return newStructLowerer(typ, layout, stack)
	case reflect.Array:
		return newArrayLowerer(typ, stack)
	case reflect.Pointer:
		return newPointerLowerer(typ, stack)
	default:
		return nil, fmt.Errorf("unsupported type for lowerer generation: %v", typ)
	}
}

// ============================================================
// 各类型的 Lifter 实现(Go 值 -> Guest 内存)
// ============================================================

// primitiveLifter 基本类型提升器：处理 bool / 各种宽度的有符号无符号整数 / 浮点。
// 这些是定长标量，无需额外分配，直接把值按字节写入由 liftByWrite 计算好的目标地址即可。
type primitiveLifter struct{}

// newPrimitiveLifter 构造一个基本类型提升器(无内部状态，返回空结构体)。
func newPrimitiveLifter() *primitiveLifter { return &primitiveLifter{} }

// lift 实现 lifter 接口：委托 liftByWrite 完成"写入 Guest 内存"。
//
// 参数示例: 把 uint8 值 7 提升进 Guest 内存
//   - val = reflect.ValueOf(uint8(7))
//
// 返回示例:
//   - (ptr=0x10000, err=nil)：值已写入 0x10000 处 1 字节
//   - (ptr=0,      err=非nil)：分配/写入失败
func (l *primitiveLifter) lift(ctx context.Context, h *Host, val reflect.Value, layout *TypeLayout) (uint32, error) {
	return liftByWrite(ctx, h, val, layout)
}

// stringLifter 字符串提升器。
// wit 里字符串在 Guest 内存中呈现为 {data_ptr: i32, len: i32} 这样一对(指针+长度)，
// 正文需要另行分配连续内存并拷贝；具体写入逻辑同样收拢在 liftByWrite 中。
type stringLifter struct{}

// lift 实现 lifter 接口。
//
// 参数示例: 提升字符串 "hi"
//   - val = reflect.ValueOf("hi")
//
// 返回示例:
//   - (ptr=0x20010, err=nil)：0x20010 处是 {ptr,len} 记录，正文被拷贝到另一块内存
//   - (ptr=0,       err=非nil)：分配/拷贝失败
func (l *stringLifter) lift(ctx context.Context, h *Host, val reflect.Value, layout *TypeLayout) (uint32, error) {
	return liftByWrite(ctx, h, val, layout)
}

// structLifter 结构体提升器。
// 构造时会为每个导出字段预生成子 lifter——主要目的不是运行期使用(真正的字段写入
// 仍由 liftByWrite 统一完成)，而是**提前校验**该结构体所有字段类型都能生成 codec，
// 并借此触发循环引用检测，让"不支持/递归类型"的错误在构造期(而非运行期)暴露。
type structLifter struct{}

// newStructLifter 构造结构体提升器，并预校验全部导出字段。
//
// 参数:
//   - typ  : 结构体类型。
//   - stack: 循环引用检测栈。
//
// 返回值:
//   - *structLifter: 结构体提升器。
//   - error        : 任一导出字段类型无法生成 lifter 时返回
//     "failed to generate lifter for field <字段名>: ..."。
func newStructLifter(typ reflect.Type, stack map[reflect.Type]struct{}) (*structLifter, error) {
	// 预生成所有导出字段的 lifter，确保类型组合可行、错误提前暴露
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !isExportedField(field) {
			continue // 未导出字段不参与 Guest 内存布局，跳过
		}
		if _, err := getOrGenerateLifterStack(field.Type, stack); err != nil {
			return nil, fmt.Errorf("failed to generate lifter for field %s: %w", field.Name, err)
		}
	}
	return &structLifter{}, nil
}

// lift 实现 lifter 接口：按 layout 中记录的字段偏移逐个写入导出字段。
//
// 参数示例: type Point struct{ X int32; Y int32 }
//   - val = reflect.ValueOf(Point{X:1, Y:2})
//   - layout 记录 X 偏移 0、Y 偏移 4
//
// 返回示例:
//   - (ptr=0x30000, err=nil)：0x30000 处依次写入 1、2
//   - (ptr=0,       err=非nil)：字段提升/分配失败
func (l *structLifter) lift(ctx context.Context, h *Host, val reflect.Value, layout *TypeLayout) (uint32, error) {
	return liftByWrite(ctx, h, val, layout)
}

// sliceLifter 切片提升器。
// 保留 elemLifter 字段：在构造期用于"校验元素类型可提升 + 判环"。
// Guest 侧切片形如 {ptr,len,cap} 三连字，元素另放连续内存。
type sliceLifter struct{ elemLifter lifter }

// newSliceLifter 构造切片提升器，先递归生成/校验元素类型的 lifter。
//
// 参数:
//   - typ  : 切片类型(如 []int32)。
//   - stack: 循环引用检测栈。
//
// 返回值:
//   - *sliceLifter: 切片提升器，内部缓存元素 lifter。
//   - error        : 元素类型不支持/递归时报错。
func newSliceLifter(typ reflect.Type, stack map[reflect.Type]struct{}) (*sliceLifter, error) {
	el, err := getOrGenerateLifterStack(typ.Elem(), stack)
	if err != nil {
		return nil, err
	}
	return &sliceLifter{elemLifter: el}, nil
}

// lift 实现 lifter 接口。
//
// 参数示例: 提升 []int32{10, 20, 30}
//   - val = reflect.ValueOf([]int32{10, 20, 30})
//
// 返回示例:
//   - (ptr=0x40000, err=nil)：0x40000 是 {数据地址,len=3,cap=3} header
//   - (ptr=0,       err=非nil)：分配/拷贝失败
func (l *sliceLifter) lift(ctx context.Context, h *Host, val reflect.Value, layout *TypeLayout) (uint32, error) {
	return liftByWrite(ctx, h, val, layout)
}

// arrayLifter 定长数组提升器。
// 数组长度编译期已知，元素在内存中连续排列，无需 header，不涉及长度记录。
type arrayLifter struct{ elemLifter lifter }

// newArrayLifter 构造数组提升器，递归校验元素类型。
// 参数与返回值含义同 newSliceLifter。
func newArrayLifter(typ reflect.Type, stack map[reflect.Type]struct{}) (*arrayLifter, error) {
	el, err := getOrGenerateLifterStack(typ.Elem(), stack)
	if err != nil {
		return nil, err
	}
	return &arrayLifter{elemLifter: el}, nil
}

// lift 实现 lifter 接口。
//
// 参数示例: 提升 [2]int32{5, 6}
//
// 返回示例:
//   - (ptr=0x50000, err=nil)：0x50000 起连续 8 字节为 5、6
//   - (ptr=0,       err=非nil)：失败
func (l *arrayLifter) lift(ctx context.Context, h *Host, val reflect.Value, layout *TypeLayout) (uint32, error) {
	return liftByWrite(ctx, h, val, layout)
}

// pointerLifter 指针提升器：解引用后提升被指向的值。
// 注意(与原注释不同，以代码为准)：nil 指针**不**被静默写成 0，
// 而是直接返回错误 "cannot lift nil pointer"——调用方需自行保证非空。
// 这是唯一在 lift 里真正手动递归(调用子 elemLifter.lift)的实现。
type pointerLifter struct{ elemLifter lifter }

// newPointerLifter 构造指针提升器，按指针元素类型递归生成 lifter。
//
// 参数:
//   - typ  : 指针类型(如 *int32 或 *MyStruct)。
//   - stack: 循环引用检测栈(指针是打破递归类型的常见手段)。
//
// 返回值:
//   - *pointerLifter: 指针提升器。
//   - error        : 元素类型不支持/递归时报错。
func newPointerLifter(typ reflect.Type, stack map[reflect.Type]struct{}) (*pointerLifter, error) {
	el, err := getOrGenerateLifterStack(typ.Elem(), stack)
	if err != nil {
		return nil, err
	}
	return &pointerLifter{elemLifter: el}, nil
}

// lift 实现 lifter 接口：先判空，非空则解引用并把元素交给 elemLifter 递归提升。
//
// 参数示例: 提升 &Point{X:1, Y:2}
//   - val = reflect.ValueOf(&Point{X:1,Y:2})
//
// 返回示例:
//   - (ptr=0x60000, err=nil)：元素已按 Point 布局写入
//   - (ptr=0,       err="cannot lift nil pointer")：val 为 nil 指针
func (l *pointerLifter) lift(ctx context.Context, h *Host, val reflect.Value, layout *TypeLayout) (uint32, error) {
	if !val.IsValid() || val.IsNil() {
		return 0, fmt.Errorf("cannot lift nil pointer")
	}
	return l.elemLifter.lift(ctx, h, val.Elem(), layout)
}

// flagsLifter flags(布尔位图)提升器。
// wit 的 flags 是一组布尔开关，Go 侧通常建模成"字段全是 bool 的结构体"，
// 提升时把它们按位压缩编码进一个整数并写入 Guest 内存。具体编码在 liftByWrite。
type flagsLifter struct{}

// newFlagsLifter 构造 flags 提升器(无状态)。
func newFlagsLifter() *flagsLifter { return &flagsLifter{} }

// lift 实现 lifter 接口。
//
// 参数示例: type Perms struct{ Read, Write, Exec bool }
//   - val = reflect.ValueOf(Perms{Read:true, Exec:true})
//
// 返回示例:
//   - (ptr=0x70000, err=nil)：该处整数第0、2位为 1
//   - (ptr=0,       err=非nil)：失败
func (l *flagsLifter) lift(ctx context.Context, h *Host, val reflect.Value, layout *TypeLayout) (uint32, error) {
	return liftByWrite(ctx, h, val, layout)
}

// variantLifter variant/option/result(联合体)提升器。
// variant 在 Guest 内存中的表示是 {discriminant: 判别整数, payload: 当前 case 的载荷}，
// payload 的偏移会按所有 case 的最大对齐做对齐。caseLifters 在构造期逐 case 预生成，
// 主要作用是校验 + 判环；真正的写入由 liftByWrite 按判别值完成。
type variantLifter struct{ caseLifters []lifter }

// newVariantLifter 构造 variant 提升器，为每个 case 递归生成 lifter。
//
// 参数:
//   - typ  : variant 类型(由 isVariant/variantFields 识别)。
//   - stack: 循环引用检测栈。
//
// 说明: 对指针类型的 case 会自动解一层(取 .Elem())来生成 codec——
//
//	因为变体的空 case / 载荷语义与指针解引用等价。
//
// 返回值:
//   - *variantLifter: variant 提升器。
//   - error         : 任一 case 无法生成 lifter 时报错。
func newVariantLifter(typ reflect.Type, stack map[reflect.Type]struct{}) (*variantLifter, error) {
	fields := variantFields(typ)
	lifters := make([]lifter, len(fields))
	for i, field := range fields {
		fieldType := field.Type
		if fieldType.Kind() == reflect.Pointer {
			fieldType = fieldType.Elem()
		}
		cl, err := getOrGenerateLifterStack(fieldType, stack)
		if err != nil {
			return nil, err
		}
		lifters[i] = cl
	}
	return &variantLifter{caseLifters: lifters}, nil
}

// lift 实现 lifter 接口。
//
// 参数示例: 提升 option(Some(42))——内部建模为判别值+载荷的结构体
//   - val = reflect.ValueOf(...)
//
// 返回示例:
//   - (ptr=0x80000, err=nil)：0x80000 处先写判别值，再在 payload 区写 42
//   - (ptr=0,       err=非nil)：失败
func (l *variantLifter) lift(ctx context.Context, h *Host, val reflect.Value, layout *TypeLayout) (uint32, error) {
	return liftByWrite(ctx, h, val, layout)
}

// ============================================================
// 各类型的 Lowerer 实现(Guest 内存 -> Go 值)
// ============================================================

// primitiveLowerer 基本类型降低器：bool / 整数 / 浮点。
// 读取前现算/取布局，然后委托 read 从指定地址把定长字节解码进 outVal。
type primitiveLowerer struct{}

// newPrimitiveLowerer 构造基本类型降低器。
func newPrimitiveLowerer() *primitiveLowerer { return &primitiveLowerer{} }

// lower 实现 lowerer 接口。
//
// 参数示例: 从 Guest 读回一个 uint8
//   - ptr = 0x10000, outVal = reflect.ValueOf(*new(uint8))(可设置)
//
// 返回示例:
//   - err = nil：*out 被填为 Guest 内存 0x10000 处的字节值
//   - err = "invalid host"：h 或 h.module 为空
func (l *primitiveLowerer) lower(ctx context.Context, h *Host, ptr uint32, outVal reflect.Value) error {
	if h == nil || h.module == nil {
		return fmt.Errorf("invalid host")
	}
	layout, err := GetOrCalculateLayout(outVal.Type())
	if err != nil {
		return err
	}
	return read(ctx, h.module.Memory(), ptr, outVal, layout)
}

// stringLowerer 字符串降低器。
// Guest 侧字符串是 {data_ptr,len}，需先读出这对值，再据此从 data_ptr 拷贝 len 字节正文。
type stringLowerer struct{}

// lower 实现 lowerer 接口。
//
// 参数示例: 从 Guest 读回字符串
//   - ptr = 0x20010, outVal = reflect.ValueOf(new(string)).Elem()
//
// 返回示例:
//   - err = nil：outVal 被 SetString 为读回的字符串(如 "hi")
//   - err = "outVal is not settable"：传入的 outVal 不可写(如取的是接口值而非可寻址值)
func (l *stringLowerer) lower(ctx context.Context, h *Host, ptr uint32, outVal reflect.Value) error {
	if h == nil || h.module == nil {
		return fmt.Errorf("invalid host")
	}
	s, err := lowerString(h.module.Memory(), ptr)
	if err != nil {
		return err
	}
	if !outVal.CanSet() {
		return fmt.Errorf("outVal is not settable")
	}
	outVal.SetString(s)
	return nil
}

// structLowerer 结构体降低器。
// 精简实现：字段 codec 列表从未被 lower 方法实际使用(字段读取由 lowerStruct 统一按
// layout 完成)，因此不再维护 fieldCodecs 字段以降低反射与内存开销；
// 只缓存 layout，避免每次 lower 重复计算类型布局。
type structLowerer struct {
	layout *TypeLayout // 缓存布局，避免每次 lower 重复计算
}

// newStructLowerer 构造结构体降低器。
// 仍会递归检查每个导出字段能否生成 lowerer，保证"类型可组合性 + 循环引用检测"
// 的错误同样在构造期暴露(行为与 lifter 侧保持一致)。
//
// 参数:
//   - typ   : 结构体类型。
//   - layout: 已算好的布局(调用方 generateLowerer 已求得并传入)。
//   - stack : 循环引用检测栈。
//
// 返回值:
//   - *structLowerer: 结构体降低器，内部持有 layout 缓存。
//   - error         : 任一导出字段无法生成 lowerer 时报错。
func newStructLowerer(typ reflect.Type, layout *TypeLayout, stack map[reflect.Type]struct{}) (*structLowerer, error) {
	// 仍递归检查字段，保证错误与 cycle 检测
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !isExportedField(field) {
			continue
		}
		if _, err := getOrGenerateLowererStack(field.Type, stack); err != nil {
			return nil, fmt.Errorf("failed to generate lowerer for field %s: %w", field.Name, err)
		}
	}
	return &structLowerer{layout: layout}, nil
}

// lower 实现 lowerer 接口：按缓存的 layout 从 ptr 逐字段读入 outVal。
// (若 layout 因故为 nil，会现场重新计算作为兜底。)
//
// 参数示例: type Point struct{ X int32; Y int32 }
//   - ptr = 0x30000, outVal = reflect.ValueOf(&Point{}).Elem()
//
// 返回示例:
//   - err = nil：Point{X:读到的X, Y:读到的Y}
//   - err = "invalid host" / 读取越界等
func (l *structLowerer) lower(ctx context.Context, h *Host, ptr uint32, outVal reflect.Value) error {
	if h == nil || h.module == nil {
		return fmt.Errorf("invalid host")
	}
	layout := l.layout
	if layout == nil {
		var err error
		layout, err = GetOrCalculateLayout(outVal.Type())
		if err != nil {
			return err
		}
	}
	return lowerStruct(ctx, h.module.Memory(), ptr, outVal, layout)
}

// sliceLowerer 切片降低器。
// elemLowerer 在构造期预生成以校验元素类型，运行期读回由 lowerSlice 统一完成。
type sliceLowerer struct{ elemLowerer lowerer }

// newSliceLowerer 构造切片降低器，递归校验元素类型。
// 参数/返回值含义同 newSliceLifter，方向相反。
func newSliceLowerer(typ reflect.Type, stack map[reflect.Type]struct{}) (*sliceLowerer, error) {
	el, err := getOrGenerateLowererStack(typ.Elem(), stack)
	if err != nil {
		return nil, err
	}
	return &sliceLowerer{elemLowerer: el}, nil
}

// lower 实现 lowerer 接口。
//
// 参数示例: 从 Guest 读回 []int32
//   - ptr = 0x40000(header 地址), outVal = reflect.ValueOf(new([]int32)).Elem()
//
// 返回示例:
//   - err = nil：outVal 变为读回的切片(会按 header.len 分配 Go 切片并拷贝元素)
func (l *sliceLowerer) lower(ctx context.Context, h *Host, ptr uint32, outVal reflect.Value) error {
	if h == nil || h.module == nil {
		return fmt.Errorf("invalid host")
	}
	return lowerSlice(ctx, h.module.Memory(), ptr, outVal)
}

// arrayLowerer 定长数组降低器。元素连续读回即可。
type arrayLowerer struct{ elemLowerer lowerer }

// newArrayLowerer 构造数组降低器，递归校验元素类型。
func newArrayLowerer(typ reflect.Type, stack map[reflect.Type]struct{}) (*arrayLowerer, error) {
	el, err := getOrGenerateLowererStack(typ.Elem(), stack)
	if err != nil {
		return nil, err
	}
	return &arrayLowerer{elemLowerer: el}, nil
}

// lower 实现 lowerer 接口。
//
// 参数示例: 从 Guest 读回 [2]int32
//   - ptr = 0x50000, outVal = reflect.ValueOf(new([2]int32)).Elem()
//
// 返回示例:
//   - err = nil：数组各元素按连续布局填充完成
func (l *arrayLowerer) lower(ctx context.Context, h *Host, ptr uint32, outVal reflect.Value) error {
	if h == nil || h.module == nil {
		return fmt.Errorf("invalid host")
	}
	return lowerArray(ctx, h.module.Memory(), ptr, outVal)
}

// pointerLowerer 指针降低器：按 ptr 是否为空决定结果。
//   - ptr == 0(Guest 空引用) -> 把 Go 指针置 nil(空 -> 空，合法映射)；
//   - ptr != 0              -> 若 Go 指针为 nil 先分配目标，再递归 Lower 填充。
//
// 这是 lower 侧唯一真正手动递归的实现。
type pointerLowerer struct{ elemLowerer lowerer }

// newPointerLowerer 构造指针降低器，递归校验指针元素类型。
// 参数/返回值含义同 newPointerLifter。
func newPointerLowerer(typ reflect.Type, stack map[reflect.Type]struct{}) (*pointerLowerer, error) {
	el, err := getOrGenerateLowererStack(typ.Elem(), stack)
	if err != nil {
		return nil, err
	}
	return &pointerLowerer{elemLowerer: el}, nil
}

// lower 实现 lowerer 接口。
//
// 参数示例: 从 Guest 读回 *Point
//   - ptr = 0x60000, outVal = reflect.ValueOf(new(*Point)).Elem()
//
// 返回示例:
//   - ptr==0    : err=nil 且 outVal 为 nil 指针
//   - ptr!=0    : err=nil 且 outVal 指向新分配的 Point，字段已填充
//   - err 非 nil: outVal.Kind() 不是指针 / 目标不可设置 / 读取失败
func (l *pointerLowerer) lower(ctx context.Context, h *Host, ptr uint32, outVal reflect.Value) error {
	if outVal.Kind() != reflect.Pointer {
		return fmt.Errorf("expected pointer, got %v", outVal.Kind())
	}
	if ptr == 0 {
		// Guest 侧空引用：Go 侧置为 nil 指针
		if outVal.CanSet() {
			outVal.Set(reflect.Zero(outVal.Type()))
		}
		return nil
	}
	if outVal.IsNil() {
		if !outVal.CanSet() {
			return fmt.Errorf("cannot lower into nil pointer of type %v", outVal.Type())
		}
		// 目标还是 nil：先分配一个元素实例再填充
		outVal.Set(reflect.New(outVal.Type().Elem()))
	}
	// 递归降低指针指向的实体(Lower 是包级按值的通用降低入口)
	return Lower(ctx, h, ptr, outVal.Elem())
}

// flagsLowerer flags(布尔位图)降低器。
// 从 Guest 内存读出一个整数，按位解码回 Go 侧布尔集合结构体。
type flagsLowerer struct{}

// newFlagsLowerer 构造 flags 降低器。
func newFlagsLowerer() *flagsLowerer { return &flagsLowerer{} }

// lower 实现 lowerer 接口。
//
// 参数示例: 读回 Perms{Read,Write,Exec}
//   - ptr = 0x70000, outVal = reflect.ValueOf(new(Perms)).Elem()
//
// 返回示例:
//   - err = nil：按该整数各位把 outVal.Read/Write/Exec 布尔字段填好
func (l *flagsLowerer) lower(ctx context.Context, h *Host, ptr uint32, outVal reflect.Value) error {
	if h == nil || h.module == nil {
		return fmt.Errorf("invalid host")
	}
	layout, err := GetOrCalculateLayout(outVal.Type())
	if err != nil {
		return err
	}
	return read(ctx, h.module.Memory(), ptr, outVal, layout)
}

// variantLowerer variant/option/result 降低器。
// caseLowerers 在构造期逐 case 预生成以校验 + 判环；运行期读取由 read 依据
// 判别值选择对应 case 布局完成(本文件中未展开其内联细节)。
type variantLowerer struct{ caseLowerers []lowerer }

// newVariantLowerer 构造 variant 降低器，为每个 case 递归生成 lowerer。
// 参数/返回值含义同 newVariantLifter。
func newVariantLowerer(typ reflect.Type, stack map[reflect.Type]struct{}) (*variantLowerer, error) {
	fields := variantFields(typ)
	lowerers := make([]lowerer, len(fields))
	for i, field := range fields {
		fieldType := field.Type
		if fieldType.Kind() == reflect.Pointer {
			fieldType = fieldType.Elem()
		}
		cl, err := getOrGenerateLowererStack(fieldType, stack)
		if err != nil {
			return nil, err
		}
		lowerers[i] = cl
	}
	return &variantLowerer{caseLowerers: lowerers}, nil
}

// lower 实现 lowerer 接口。
//
// 参数示例: 读回 option，Guest 判别值=1 表示 Some
//   - ptr = 0x80000, outVal = reflect.ValueOf(&someStruct{}).Elem()
//
// 返回示例:
//   - err = nil：判别值与载荷均被正确读回并填充 outVal
func (l *variantLowerer) lower(ctx context.Context, h *Host, ptr uint32, outVal reflect.Value) error {
	if h == nil || h.module == nil {
		return fmt.Errorf("invalid host")
	}
	layout, err := GetOrCalculateLayout(outVal.Type())
	if err != nil {
		return err
	}
	return read(ctx, h.module.Memory(), ptr, outVal, layout)
}
