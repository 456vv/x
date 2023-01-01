package vweb_dynamic

import (
	"io"
	"io/ioutil"
	"log"
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
	rootPath string // 文件目录
	pagePath string // 文件名称
	inited   bool

	options interp.Options
	mainFunc reflect.Value
}

func (T *Yaegi) init() {
	if T.inited {
		return
	}
	T.inited = true
	T.options.GoPath = T.rootPath
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
	return T.parse(content)
}

func (T *Yaegi) ParseFile(path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return T.parse(string(b))
}

func (T *Yaegi) SetPath(root, page string) {
	T.rootPath = root
	T.pagePath = page
}

func (T *Yaegi) Parse(r io.Reader) (err error) {
	contact, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	return T.parse(string(contact))
}

func (T *Yaegi) newInterpre(src string) (*interp.Interpreter, error) {
	i := interp.New(T.options)
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

	_, err := i.Eval(src)
	return i, err
}

func (T *Yaegi) parse(src string) error {
	T.init()

	interpre, err := T.newInterpre(src)
	if err != nil {
		return err
	}

	T.mainFunc, err = interpre.Eval("Main")
	if err != nil {
		return err
	}
	return nil
}

func (T *Yaegi) Execute(out io.Writer, in interface{}) (err error) {
	if !T.mainFunc.IsValid() {
		return errTemplateNotParse
	}

	res := T.mainFunc
	if res.Kind() == reflect.Func {
		rt := res.Type()
		if rt.NumIn() == 1 {
			retn, err := call(res, in)
			if err != nil {
				return err
			}
			if rt.NumOut() == 1 {
				switch rv := retn[0].(type) {
				case string:
					io.WriteString(out, rv)
				case []byte:
					out.Write(rv)
				default:
					// 暂时不显示无法识别类型
					log.Printf("yaegi url(%s) returned unrecognized data type(%s)\n", T.pagePath, reflect.ValueOf(rv).Elem().Type().String())
				}
			}
		}
	}
	return nil
}
