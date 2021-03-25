package vweb_dynamic

import (
	"sync"
	"log"
	"io"
	"os"
	"path/filepath"
	"fmt"
	"errors"
	"reflect"
	"bytes"
    "github.com/456vv/vweb/v2"
	"github.com/goplus/gop"
	"github.com/goplus/gop/exec/bytecode"
	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/token"
	"github.com/goplus/gop/exec.spec"
	_ "github.com/goplus/gop/lib/builtin"
	_ "github.com/goplus/gop/lib/unsafe"
    _ "github.com/456vv/x/goplus_lib"
)

func execierrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(error).Error()
	p.Ret(1, ret0)
}
var gopulusOnce sync.Once
type GoPlus struct{
	rootPath			string																// 文件目录
	pagePath			string																// 文件名称
 	name				string
 	
 	ctx					*bytecode.Context
 	fset				*token.FileSet
 	pkgs				map[string]*ast.Package //存放每个文件的代码
 	mainFn				exec.FuncInfo
 	inited				bool
}
func (T *GoPlus) init(){
	if T.inited {
		return
	}
	if cl.CallBuiltinOp == nil {
		cl.CallBuiltinOp = bytecode.CallBuiltinOp
	}
	if T.fset == nil {
		T.fset = token.NewFileSet()
	}
	if T.pkgs == nil {
		T.pkgs = make(map[string]*ast.Package)
	}
	gopulusOnce.Do(func(){
		gopI := bytecode.FindGoPackage("").(*bytecode.GoPackage)
		if gopI == nil {
			gopI = bytecode.NewGoPackage("")
		}
		gopI.RegisterFuncs(gopI.Func("(error).Error", (error).Error, execierrorError))
		
		for name, fn := range vweb.TemplateFunc {
			tfn := reflect.TypeOf(fn)
			switch tfn.Kind() {
			case reflect.Func:
				fnc := func(name string, tfn reflect.Type, fn interface{}) func(arity int, p *gop.Context) {
					return func(arity int, p *gop.Context){
						args := p.GetArgs(arity)
						retn, err := vweb.ExecFunc(fn, args...)
						if err != nil {
							log.Printf("call %s(%v) (%v, %v)\n", name, args, retn, err)
						}
						p.Ret(arity, retn...)
					}
				}(name, tfn, fn)
				gopI.RegisterFuncvs(gopI.Funcv(name, fn, fnc))
			default:
				log.Printf("导入内置函数，无法识别 %s 类型\n", tfn.Kind().String())
			}
		}
	})
	T.inited = true
}

func (T *GoPlus) newPkg(name string) *ast.Package {
	pkg, found := T.pkgs[name]
	if !found {
		pkg = &ast.Package{
			Name:  name,
			Files: make(map[string]*ast.File),
		}
		T.pkgs[name] = pkg
	}
	return pkg
}

func (T *GoPlus) ParseText(name, content string) error {
	T.init()
	if T.name == "" {
		T.name = name
	}
	
	r := bytes.NewBufferString(content)
	return T.Parse(r)
}

func (T *GoPlus) ParseFile(path string) error {
	T.init()
	if T.name == "" {
		T.name = filepath.Base(path)
	}
	
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return T.Parse(file)
}

func (T *GoPlus) ParseDir(dir string) error {
	T.init()
	pkgs, err := parser.ParseDir(T.fset, dir, nil, 0)
	if err != nil {
		return err
	}
	for name, file := range pkgs {
		T.pkgs[name]=file
	}
	return nil
}

func (T *GoPlus) SetPath(root, page string){
	T.rootPath = root
	T.pagePath = page
    T.name = filepath.Base(T.pagePath)
}

func (T *GoPlus) Parse(r io.Reader) error {
	T.init()
	fileHeader, body, err := vweb.TemplateSeparation(r)
	if err != nil {
		return err
	}
	filesContent, err := fileHeader.OpenFile(T.rootPath, T.pagePath)
	if err != nil {
		return err
	}
	filesContent[T.name]=string(body)
	for filename, filecontent := range filesContent {
		src, err := parser.ParseFile(T.fset, filename, filecontent, 0)
		if err != nil {
			return err
		}
		pkg := T.newPkg(src.Name.Name)
		pkg.Files[filename]=src
	}
	return nil
}

func (T *GoPlus) Execute(out io.Writer, in interface{}) (err error) {
	if !T.inited {
		return errors.New("The template has not been parsed and is not available!")
	}
	
	//第一次
	if T.ctx == nil {
		
		//设置入口标识
		tpkg := getPkg(T.pkgs)
		pkgAct := cl.PkgActClAll
		entrance := "init"
		if tpkg.Name == "main" {
			entrance = tpkg.Name
			pkgAct = cl.PkgActClMain
		}
		//组装一个模块
		builder := bytecode.NewBuilder(nil)
		pkg, err := cl.NewPackage(builder.Interface(), tpkg, T.fset, pkgAct)
		if err != nil {
			return err
		}
		
		//判断模块有没有入口
		kind, sym, ok := pkg.Find(entrance)
		if !ok || kind != cl.SymFunc {
			return fmt.Errorf("function %s.%s() not found", tpkg.Name, entrance)
		}
		
		//创建代码上下文
		code := builder.Resolve()
		T.ctx = bytecode.NewContext(code)
		
		T.mainFn = sym.(*cl.FuncDecl).Compile()
	}
	
	ch := make(chan bool)
	T.ctx.Go(0, func(goctx *bytecode.Context){
		defer close(ch)
		
		if T.mainFn.NumIn() == 1 {
			//非可变参数
			T.mainFn.Args(reflect.TypeOf(in))
			goctx.Push(in)
		}
		goctx.Call(T.mainFn)
		
		if T.mainFn.NumOut() == 1 {
			switch rv := goctx.Get(-1).(type) {
			case string:io.WriteString(out, rv)
			case []byte:out.Write(rv)
			default:
				log.Printf("goplus url(%s) returned unrecognized data type(%s)\n", T.pagePath, reflect.ValueOf(&rv).Elem().Elem().Type().String())
			}
		}
	})
	
	<-ch
	return nil
}

func getPkg(pkgs map[string]*ast.Package) *ast.Package {
	for _, pkg := range pkgs {
		return pkg
	}
	return nil
}















