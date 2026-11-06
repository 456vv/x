package vweb_dynamic

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"sync"
	"sync/atomic"

	"github.com/456vv/x/yaegi_lib"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/traefik/yaegi/stdlib/syscall"
	"github.com/traefik/yaegi/stdlib/unsafe"
)

var (
	yaegiOnce sync.Once
	yaegiFunc interp.Exports
)

// Yaegi 动态页面解释器，基于 Yaegi 实现 Go 模板/脚本执行。
//
// 并发安全说明：
//   - 同一 *Yaegi 实例在 Parse / ParseFile / ParseText / Close 等写操作期间，
//     不得被其他 goroutine 并发调用 Execute。
//   - 解析完成后，多个 goroutine 可安全并发调用 Execute 执行不同入口函数：
//     Execute 在入口处通过 atomic.Pointer 加载解释器快照，并持有该快照直到本次调用结束；
//     后续 Parse 替换解释器时，旧快照由 mainFunc 闭包引用保持存活，不会被提前回收。
//   - mainFunc 缓存与解释器快照强相关：Parse 更新 interpreRef 时同步清空 mainFunc，
//     确保新调用不会复用旧解释器产生的入口函数。
type Yaegi struct {
	rootPath  string // 文件根目录（GoPath）
	pagePath  string // 当前页面文件路径
	entryName string // 默认入口函数名
	inited    bool
	m         sync.Mutex
	mainFunc  map[string]reflect.Value // 入口函数缓存（与当前 interpreRef 对应）

	// interpreRef 原子保存当前解释器快照；Execute 在入口处 Load 一次并全程持有，
	// 避免读写 T.interpre 字段的 data race，也让 parse/Close 与 Execute 无需长期互斥。
	interpreRef atomic.Pointer[interp.Interpreter]
}

// init 初始化解释器内部状态（全程加锁，消除竞态）
func (T *Yaegi) init() {
	T.m.Lock()
	defer T.m.Unlock()

	if T.inited {
		return
	}
	T.inited = true
	T.mainFunc = make(map[string]reflect.Value)

	yaegiOnce.Do(func() {
		// 增加内置函数
		yaegiFunc = make(interp.Exports)
		builtinMap := make(map[string]reflect.Value)
		for name, fn := range TemplateFunc {
			builtinMap[name] = reflect.ValueOf(fn)
		}
		builtinMap["Symbols"] = reflect.ValueOf(yaegiFunc)
		yaegiFunc["this/this"] = builtinMap
	})
}

// ParseText 从字符串解析
func (T *Yaegi) ParseText(name, content string) error {
	return T.parse([]byte(content))
}

// ParseFile 从文件路径解析
func (T *Yaegi) ParseFile(p string) error {
	file, err := os.Open(p)
	if err != nil {
		return err
	}
	defer file.Close()
	return T.Parse(file)
}

// SetPath 设置根目录与当前页面路径
func (T *Yaegi) SetPath(root, page string) {
	T.rootPath = root
	T.pagePath = page
}

func (T *Yaegi) setHeaderLine(buf *bytes.Buffer) TemplateHeader {
	l := fileHeaderLine(buf)
	return templateHeader(l)
}

// Parse 从 Reader 解析模板源码
func (T *Yaegi) Parse(r io.Reader) (err error) {
	contact, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	return T.parse(contact)
}

func (T *Yaegi) newInterpre() (*interp.Interpreter, error) {
	options := interp.Options{
		GoPath: T.rootPath,
	}
	i := interp.New(options)
	// 内置标准库
	if err := i.Use(stdlib.Symbols); err != nil {
		return nil, err
	}
	if err := i.Use(unsafe.Symbols); err != nil {
		return nil, err
	}
	if err := i.Use(syscall.Symbols); err != nil {
		return nil, err
	}
	if err := i.Use(interp.Symbols); err != nil {
		return nil, err
	}
	if err := i.Use(yaegi_lib.Symbols); err != nil {
		return nil, err
	}
	// 自定函数
	if err := i.Use(yaegiFunc); err != nil {
		return nil, err
	}

	i.ImportUsed()

	return i, nil
}

func (T *Yaegi) parse(src []byte) error {
	T.init()

	buf := bytes.NewBuffer(src)
	th := T.setHeaderLine(buf)

	interpre, err := T.newInterpre()
	if err != nil {
		return err
	}

	// 加载额外的文件
	dir := filepath.Dir(T.pagePath)
	for _, f := range th.File {
		var p string
		if path.IsAbs(f) {
			p = filepath.Join(T.rootPath, f)
		} else {
			p = filepath.Join(T.rootPath, dir, f)
		}
		p = filepath.Clean(p)
		if _, err = interpre.EvalPath(p); err != nil {
			return err
		}
	}

	// 加载主文件
	if _, err := interpre.Eval(buf.String()); err != nil {
		return err
	}

	T.m.Lock()
	T.entryName = th.EntryName
	// 切换解释器快照并让旧 mainFunc 缓存作废：
	// 缓存项与旧解释器强绑定，必须整体替换，避免 Execute 在新快照下复用旧入口函数。
	T.interpreRef.Store(interpre)
	T.mainFunc = make(map[string]reflect.Value)
	T.m.Unlock()
	return nil
}

// Execute 执行指定入口函数，并将结果写入 out。
// 支持返回 string / []byte / io.Reader / fmt.Stringer / error 等常见类型。
// 若函数返回 (T, error) 且 error 非 nil，则直接返回该 error。
//
// 执行期间持有 interpreRef 的一次快照，Parse/Close 并发切换解释器不会影响本次调用；
// 旧解释器由本地 mainFunc 的闭包引用保持存活，直到本次调用结束。
func (T *Yaegi) Execute(entryName string, out io.Writer, in ...any) (err error) {
	// 快照解释器：把"解释器生命周期"从互斥锁保护中解耦，
	// 后续 Parse/Close 通过 Store 覆盖引用不影响本 goroutine 已持有的旧快照。
	interpre := T.interpreRef.Load()
	if interpre == nil {
		return errTemplateNotParse
	}

	T.m.Lock()
	entryName = entryname(T.entryName, entryName)
	mainFunc, ok := T.mainFunc[entryName]
	if !ok {
		mainFunc, err = interpre.Eval(entryName)
		if err != nil {
			T.m.Unlock()
			return err
		}
		if T.mainFunc == nil {
			T.mainFunc = make(map[string]reflect.Value)
		}
		T.mainFunc[entryName] = mainFunc
	}
	T.m.Unlock()

	if mainFunc.Kind() != reflect.Func {
		return fmt.Errorf("yaegi: entry %q is not a function", entryName)
	}

	retn, err := call(mainFunc, in...)
	if err != nil {
		return err
	}

	if len(retn) == 0 {
		return nil
	}

	// 优先检查最后一个返回值是否为 error（支持 (T, error) 模式）
	last := retn[len(retn)-1]
	if e, ok := last.(error); ok {
		if e != nil {
			return e
		}
		// 唯一返回值且为 nil error，直接成功
		if len(retn) == 1 {
			return nil
		}
	}

	// 取第一个有效返回值进行输出
	rv := retn[0]
	if rv == nil {
		return nil
	}
	switch v := rv.(type) {
	case string:
		_, err = io.WriteString(out, v)
	case []byte:
		_, err = out.Write(v)
	case io.Reader:
		_, err = io.Copy(out, v)
	case fmt.Stringer:
		_, err = io.WriteString(out, v.String())
	case error:
		// 已在上方处理，此处仅兜底
		if v != nil {
			return v
		}
	default:
		// 无法识别的类型仅记录日志，不中断
		log.Printf("yaegi call %s returned unrecognized data type(%s)\n", entryName, reflect.TypeOf(rv).String())
	}
	return err
}

// Close 关闭解释器（当前无状态资源，保留接口兼容性）。
// 关闭后 interpreRef 置空、mainFunc 缓存清空；已在执行中的 Execute 继续持有旧快照直到完成。
func (T *Yaegi) Close() error {
	T.m.Lock()
	defer T.m.Unlock()
	T.interpreRef.Store(nil)
	T.mainFunc = nil
	T.inited = false
	return nil
}
