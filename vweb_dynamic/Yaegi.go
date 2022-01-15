package vweb_dynamic

import (
	"sync"
	"io"
	"io/ioutil"
	"path/filepath"
	"errors"
	"reflect"
	"log"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/traefik/yaegi/stdlib/unsafe"
	"github.com/traefik/yaegi/stdlib/syscall"
	"github.com/456vv/x/yaegi_lib"
)

var yaegiOnce sync.Once
var yaegiFunc interp.Exports

type Yaegi struct{
	rootPath			string																// 文件目录
	pagePath			string																// 文件名称
 	name				string
 	inited				bool
 	
 	options			interp.Options
 	program			*interp.Program
 	interpreter		*interp.Interpreter
}

func (T *Yaegi) init(){
	if T.inited {
		return
	}
	yaegiOnce.Do(func (){
		//增加内置函数
		yaegiFunc = make(interp.Exports)
		builtin := make(map[string]reflect.Value)
		for name, fn := range TemplateFunc {
			builtin[name] = reflect.ValueOf(fn)
		}
		builtin["Symbols"] = reflect.ValueOf(yaegiFunc)
		yaegiFunc["this/this"] = builtin
		
		T.options.GoPath = T.rootPath
	})
	T.inited = true
}

func (T *Yaegi) ParseText(name, content string) error {
	T.name = name
	return T.parse(content)
}

func (T *Yaegi) ParseFile(path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	T.name = filepath.Base(path)
	script := string(b)
	return T.parse(script)
}

func (T *Yaegi) SetPath(root, page string){
	T.rootPath = root
	T.pagePath = page
    T.name = filepath.Base(page)
}

func (T *Yaegi) Parse(r io.Reader) (err error) {
	contact, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	
	script := string(contact)
	return T.parse(script)
}

func (T *Yaegi) parse(script string) error {
	T.init()
	
	i := interp.New(T.options)
	//内置标准库
	if err := i.Use(stdlib.Symbols); err != nil {
		return err
	}
	if err := i.Use(unsafe.Symbols); err != nil {
		return err
	}
	if err := i.Use(syscall.Symbols); err != nil {
		return err
	}
	if err := i.Use(interp.Symbols); err != nil {
		return err
	}
	if err := i.Use(yaegi_lib.Symbols); err != nil {
		return err
	}
	//自定函数
	if err := i.Use(yaegiFunc); err != nil {
		return err
	}
	
	i.ImportUsed()
	_, err := i.Eval(script)
	if err != nil {
		return err
	}
	T.interpreter = i
	return nil
}

func (T *Yaegi) Execute(out io.Writer, in interface{}) (err error) {
	if T.interpreter == nil {
		return errors.New("The template has not been parsed and is not available!")
	}
	
	var res reflect.Value
   	res, err = T.interpreter.Eval("Main")
	if err != nil {
		return err
	}
	
	if res.Kind() == reflect.Func {
		rt := res.Type()
		if rt.NumIn() == 1  {
			retn, err := call(res, in)
			if err != nil {
				return err
			}
			if rt.NumOut() == 1 {
				switch rv := retn[0].(type) {
				case string:io.WriteString(out, rv)
				case []byte:out.Write(rv)
				default:
					//暂时不显示无法识别类型
					log.Printf("yaegi url(%s) returned unrecognized data type(%s)\n", T.pagePath, reflect.ValueOf(&rv).Elem().Elem().Type().String())
				}
			}
		}
	}
	return nil
}

