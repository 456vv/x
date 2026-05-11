package vweb_dynamic

import (
	"bytes"
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

type Yaegi struct {
	rootPath  string // 文件目录
	pagePath  string // 文件路径
	entryName string
	inited    bool
	m         sync.Mutex
	mainFunc  map[string]reflect.Value
	interpre  *interp.Interpreter
}

func (T *Yaegi) init() {
	if T.inited {
		return
	}
	T.inited = true

	T.m.Lock()
	defer T.m.Unlock()

	T.mainFunc = make(map[string]reflect.Value)

	yaegiOnce.Do(func() {
		// 增加内置函数
		yaegiFunc = make(interp.Exports)
		builtin := make(map[string]reflect.Value)
		for name, fn := range TemplateFunc {
			builtin[name] = reflect.ValueOf(fn)
		}
		builtin["Symbols"] = reflect.ValueOf(yaegiFunc)
		yaegiFunc["this/this"] = builtin
	})
}

func (T *Yaegi) ParseText(name, content string) error {
	return T.parse([]byte(content))
}

func (T *Yaegi) ParseFile(p string) error {
	file, err := os.Open(p)
	if err != nil {
		return err
	}
	defer file.Close()
	return T.Parse(file)
}

func (T *Yaegi) SetPath(root, page string) {
	T.rootPath = root
	T.pagePath = page
}

func (T *Yaegi) setHeaderLine(buf *bytes.Buffer) TemplateHeader {
	l := fileHeaderLine(buf)
	hm := templateHeader(l)
	return hm
}

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
	T.entryName = th.EntryName

	interpre, err := T.newInterpre()
	if err != nil {
		return err
	}

	// 加载额外的文件
	dir := filepath.Dir(T.pagePath)
	for _, f := range th.File {
		if path.IsAbs(f) {
			_, err = interpre.EvalPath(filepath.Join(T.rootPath, f))
		} else {
			p := filepath.Join(T.rootPath, dir, f)
			_, err = interpre.EvalPath(p)
		}
		if err != nil {
			return err
		}
	}

	// 加载主文件
	if _, err := interpre.Eval(buf.String()); err != nil {
		return err
	}

	T.interpre = interpre
	return nil
}

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
		T.mainFunc[entryName] = mainFunc
	}
	T.m.Unlock()

	if mainFunc.Kind() == reflect.Func {
		retn, err := call(mainFunc, in...)
		if err != nil {
			return err
		}

		if mainFunc.Type().NumOut() == 1 {
			switch rv := retn[0].(type) {
			case string:
				io.WriteString(out, rv)
			case []byte:
				out.Write(rv)
			case io.Reader:
				io.Copy(out, rv)
			default:
				// 暂时不显示无法识别类型
				log.Printf("yaegi call %s returned unrecognized data type(%s)\n", entryName, reflect.ValueOf(rv).Type().String())
			}
		}
	}
	return nil
}

func (T *Yaegi) Close() error {
	return nil
}
