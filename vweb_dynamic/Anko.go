package vweb_dynamic

import (
	"sync"
	"io"
	"io/ioutil"
	"path/filepath"
	"errors"
	"log"
	"reflect"
	"github.com/mattn/anko/env"
	"github.com/mattn/anko/parser"
	"github.com/mattn/anko/vm"
	"github.com/mattn/anko/ast"
	"github.com/mattn/anko/core"
	_ "github.com/mattn/anko/packages" //加入默认包
	"github.com/456vv/vweb/v2"
)

var anko_env *env.Env
var ankoOnce sync.Once

type Anko struct{
	rootPath			string																// 文件目录
	pagePath			string																// 文件名称
 	name				string
 	stmt				ast.Stmt
 	inited				bool
}

func (T *Anko) init(){
	if T.inited {
		return
	}
	ankoOnce.Do(func (){
		//增加anko 模块包
		parser.EnableErrorVerbose()	//解析错误详细信息
		anko_env = env.NewEnv()
		core.Import(anko_env) 		//加载内置的一些函数
		
		//增加内置函数
		for name, fn := range vweb.TemplateFunc {
			anko_env.Define(name, fn)
		}
	})
	T.inited = true
}

func (T *Anko) ParseText(name, content string) error {
	T.name = name
	return T.parse(content)
}

func (T *Anko) ParseFile(path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	T.name = filepath.Base(path)
	script := string(b)
	return T.parse(script)
}

func (T *Anko) SetPath(root, page string){
	T.rootPath = root
	T.pagePath = page
    T.name = filepath.Base(page)
}

func (T *Anko) Parse(r io.Reader) (err error) {
	contact, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	
	script := string(contact)
	return T.parse(script)
}

func (T *Anko) parse(script string) error {
	T.init()
	stmt, err := parser.ParseSrc(script)
	if err != nil {
		if pe, ok := err.(*parser.Error); ok {
			pe.Filename = filepath.Join(T.rootPath, T.pagePath)
		}
		return err
	}
	T.stmt = stmt
	return nil
}

func (T *Anko) Execute(out io.Writer, in interface{}) (err error) {
	if T.stmt == nil {
		return errors.New("The template has not been parsed and is not available!")
	}
	
	env := anko_env.NewEnv()
	env.Define("T", in)

	var retn interface{}
    if tdot, ok := in.(vweb.DotContexter); ok {
		retn, err = vm.RunContext(tdot.Context(), env, nil, T.stmt)
    }else{
    	retn, err = vm.Run(env, nil, T.stmt)
    }
	if err != nil {
		//排除中断的错误
		//可能用户关闭连接
		if err.Error() == vm.ErrInterrupt.Error() {
			return nil
		}
		return err
	}
	if out != nil && retn != nil {
		switch rv := retn.(type) {
		case string:io.WriteString(out, rv)
		case []byte:out.Write(rv)
		default:log.Printf("anko returned unrecognized data type %s\n", reflect.TypeOf(&rv).Elem().String())
		}
	}
	return nil
}

