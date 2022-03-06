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
 	pool			sync.Pool
 	src				string
}

func (T *Yaegi) init(){
	if T.inited {
		return
	}
	T.options.GoPath = T.rootPath
	yaegiOnce.Do(func (){
		//增加内置函数
		yaegiFunc = make(interp.Exports)
		builtin := make(map[string]reflect.Value)
		for name, fn := range TemplateFunc {
			builtin[name] = reflect.ValueOf(fn)
		}
		builtin["Symbols"] = reflect.ValueOf(yaegiFunc)
		yaegiFunc["this/this"] = builtin
		
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

func (T *Yaegi) newInterpre() (*interp.Interpreter, error){
	if T.src == "" {
		return nil, errors.New("The template has not been parsed and is not available!")
	}
	i := interp.New(T.options)
	//内置标准库
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
	//自定函数
	if err := i.Use(yaegiFunc); err != nil {
		return nil, err
	}
	
	i.ImportUsed()
	
	_, err := i.Eval(T.src)
	return i, err
}
func (T *Yaegi) parse(script string) error {
	T.src = script
	return nil
}

func (T *Yaegi) Execute(out io.Writer, in interface{}) (err error) {
	T.init()

	interpre, ok := T.pool.Get().(*interp.Interpreter)
	if !ok {
		//新建一个
		interpre, err = T.newInterpre()
		if err != nil {
			return err
		}
	}
	defer T.pool.Put(interpre)
	
	var res reflect.Value
   	res, err = interpre.Eval("Main")
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

