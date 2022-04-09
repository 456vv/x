package vweb_dynamic

import (
	"sync"
	"io"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"errors"
	"log"
	"go/token"
	"golang.org/x/tools/go/ssa"
	"github.com/goplus/gossa"
	_ "github.com/goplus/gossa/pkg" //加入默认包
	_ "github.com/456vv/x/gossa_lib"
)

var ssaOnce sync.Once

type Gossa struct{
	rootPath			string																// 文件目录
	pagePath			string																// 文件名称
 	name				string
 	inited				bool
 	
 	c					*gossa.Context
 	mainPkg				*ssa.Package
}

func (T *Gossa) init(){
	if T.inited {
		return
	}
	ssaOnce.Do(func (){
		//增加内置函数
		for name, fn := range TemplateFunc {
			gossa.RegisterExternal(name, fn)
		}
	})
	T.inited = true
}

func (T *Gossa) ParseText(name, content string) error {
	T.name = name
	return T.parse(content)
}

func (T *Gossa) ParseFile(path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	T.name = filepath.Base(path)
	script := string(b)
	return T.parse(script)
}

func (T *Gossa) SetPath(root, page string){
	T.rootPath = root
	T.pagePath = page
    T.name = filepath.Base(page)
}

func (T *Gossa) Parse(r io.Reader) (err error) {
	contact, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	
	script := string(contact)
	return T.parse(script)
}

func (T *Gossa) parse(script string) error {
	T.init()
	
	T.c = gossa.NewContext(0)
	fset := token.NewFileSet()
	pkg, err := T.c.LoadFile(fset, T.name,  script)
	if err != nil {
		return err
	}
	T.mainPkg = pkg
	return nil
}

func (T *Gossa) Execute(out io.Writer, in interface{}) (err error) {
	if T.mainPkg == nil {
		return errors.New("The template has not been parsed and is not available!")
	}
	
	interp, err := T.c.NewInterp(T.mainPkg)
	if err != nil {
		return err
	}
	retn, err := interp.RunFunc("Main", in)
	if err != nil {
		return err
	}
	
	switch rv := retn.(type) {
	case string:io.WriteString(out, rv)
	case []byte:out.Write(rv)
	default:
		//暂时不显示无法识别类型
		log.Printf("goass url(%s) returned unrecognized data type(%s)\n", T.pagePath, reflect.ValueOf(&rv).Elem().Elem().Type().String())
	}
	
	return nil
}

