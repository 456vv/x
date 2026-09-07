package witgo

// ============================================================
// WIT 类型定义：Option, Result, Tuple
// ============================================================

// Unit 是一个空结构体，用于表示无载荷的 variant case。
// 在 Option 的 None  case 中使用，占位但不占用额外空间。
type Unit struct{}

// Option[T] 表示可选值，None 或 Some(T)。
// 对应 WIT 的 option 类型，实现 Optioner 接口以便识别。
//
// 内存布局: {discriminant: u8, payload: T}
// - discriminant 0 = None（Unit）
// - discriminant 1 = Some(T)
type Option[T any] struct {
	None *Unit `wit:"case(0)"` // None case
	Some *T    `wit:"case(1)"` // Some case
}

// IsOption 实现 Optioner 接口，标记此类型为 Option。
func (o Option[T]) IsOption() {}

// Some 创建一个包含值的 Option。
//
// 参数:
//   - value: 要包装的值
//
// 返回:
//   - Option[T]: 包含 value 的 Option
//
// 示例:
//
//	opt := Some(42)
func Some[T any](value T) Option[T] {
	return Option[T]{Some: &value}
}

// SomePtr 返回 Option 的指针。
func SomePtr[T any](value T) *Option[T] {
	return &Option[T]{Some: &value}
}

// None 返回空 Option。
func None[T any]() Option[T] {
	return Option[T]{None: &Unit{}}
}

// NonePtr 返回空 Option 的指针。
func NonePtr[T any]() *Option[T] {
	return &Option[T]{None: &Unit{}}
}

// IsSome 返回是否为 Some。
func (o Option[T]) IsSome() bool { return o.Some != nil }

// IsNone 返回是否为 None。
func (o Option[T]) IsNone() bool { return o.Some == nil }

// Value 返回 Some 的值和 true，否则返回零值和 false。
func (o Option[T]) Value() (T, bool) {
	if o.IsSome() {
		return *o.Some, true
	}
	var zero T
	return zero, false
}

// Expect 若为 Some 返回值，否则 panic 并显示消息。
func (o Option[T]) Expect(msg string) T {
	if o.IsSome() {
		return *o.Some
	}
	panic(msg)
}

// Unwrap 若为 Some 返回值，否则 panic 显示标准消息。
func (o Option[T]) Unwrap() T {
	return o.Expect("called `Unwrap()` on a `None` value")
}

// UnwrapOr 若为 Some 返回值，否则返回 defaultValue。
func (o Option[T]) UnwrapOr(defaultValue T) T {
	if o.IsSome() {
		return *o.Some
	}
	return defaultValue
}

// UnwrapOrElse 若为 Some 返回值，否则调用 f 并返回其值。
func (o Option[T]) UnwrapOrElse(f func() T) T {
	if o.IsSome() {
		return *o.Some
	}
	if f == nil {
		var zero T
		return zero
	}
	return f()
}

// Result[T, E] 表示操作结果，Ok(T) 或 Err(E)。
// 对应 WIT 的 result 类型，实现 Resulter 接口。
//
// 内存布局: {discriminant: u8, payload: T 或 E}
// - discriminant 0 = Ok(T)
// - discriminant 1 = Err(E)
type Result[T, E any] struct {
	Ok  *T `wit:"case(0)"` // 成功情况
	Err *E `wit:"case(1)"` // 错误情况
}

// IsResult 实现 Resulter 接口。
func (r Result[T, E]) IsResult() {}

// Ok 创建一个成功的 Result。
func Ok[T any, E any](value T) Result[T, E] {
	return Result[T, E]{Ok: &value}
}

// Err 创建一个错误的 Result。
func Err[T any, E any](err E) Result[T, E] {
	return Result[T, E]{Err: &err}
}

// IsOk 返回是否为 Ok。
func (r Result[T, E]) IsOk() bool { return r.Ok != nil }

// IsErr 返回是否为 Err。
func (r Result[T, E]) IsErr() bool { return r.Ok == nil && r.Err != nil }

// Value 若为 Ok 返回值，否则返回零值和 false。
func (r Result[T, E]) Value() (T, bool) {
	if r.IsOk() {
		return *r.Ok, true
	}
	var zero T
	return zero, false
}

// ErrValue 若为 Err 返回错误值，否则返回零值和 false。
func (r Result[T, E]) ErrValue() (E, bool) {
	if r.IsErr() {
		return *r.Err, true
	}
	var zero E
	return zero, false
}

// Expect 若为 Ok 返回值，否则 panic。
func (r Result[T, E]) Expect(msg string) T {
	if r.IsOk() {
		return *r.Ok
	}
	panic(msg)
}

// Unwrap 若为 Ok 返回值，否则 panic。
func (r Result[T, E]) Unwrap() T {
	return r.Expect("called `Unwrap()` on an `Err` value")
}

// UnwrapOr 若为 Ok 返回值，否则返回 defaultValue。
func (r Result[T, E]) UnwrapOr(defaultValue T) T {
	if r.IsOk() {
		return *r.Ok
	}
	return defaultValue
}

// UnwrapOrElse 若为 Ok 返回值，否则调用 f 并返回其值。
func (r Result[T, E]) UnwrapOrElse(f func() T) T {
	if r.IsOk() {
		return *r.Ok
	}
	if f == nil {
		var zero T
		return zero
	}
	return f()
}

// Tuple[T0, T1] 表示二元组。
type Tuple[T0, T1 any] struct {
	F0 T0
	F1 T1
}

// Tuple3[T0, T1, T2] 表示三元组。
type Tuple3[T0, T1, T2 any] struct {
	F0 T0
	F1 T1
	F2 T2
}

// String 返回字符串指针，便于构造 Option[String] 等。
func String(s string) *string { return &s }

// UnitResult 是用于无载荷结果的类型别名。
// 对应 WIT 的 unit result，成功为 0，失败为 1。
type UnitResult = uint32

// UintOk 返回表示成功的 UnitResult。
func UintOk() UnitResult { return 0 }

// UintErr 返回表示失败的 UnitResult。
func UintErr() UnitResult { return 1 }
