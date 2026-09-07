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
type Yaegi struct {
	rootPath  string // 文件根目录（GoPath）
	pagePath  string // 当前页面文件路径
	entryName string // 默认入口函数名
	inited    bool
	m         sync.Mutex
	mainFunc  map[string]reflect.Value // 入口函数缓存（并发安全）
	interpre  *interp.Interpreter
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
	T.interpre = interpre
	T.m.Unlock()
	return nil
}

// Execute 执行指定入口函数，并将结果写入 out。
// 支持返回 string / []byte / io.Reader / fmt.Stringer / error 等常见类型。
// 若函数返回 (T, error) 且 error 非 nil，则直接返回该 error。
func (T *Yaegi) Execute(entryName string, out io.Writer, in ...any) (err error) {
	if T.interpre == nil {
		return errTemplateNotParse
	}

	T.m.Lock()
	entryName = entryname(T.entryName, entryName)
	mainFunc, ok := T.mainFunc[entryName]
	if !ok {
		mainFunc, err = T.interpre.Eval(entryName)
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

// Close 关闭解释器（当前无状态资源，保留接口兼容性）
func (T *Yaegi) Close() error {
	T.m.Lock()
	defer T.m.Unlock()
	T.interpre = nil
	T.mainFunc = nil
	T.inited = false
	return nil
}
