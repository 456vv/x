package vweb_dynamic

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"
	"sync"

	_ "github.com/456vv/x/gossa_lib"
	"github.com/goplus/gossa"
	_ "github.com/goplus/gossa/pkg" // 加入默认包
	"golang.org/x/tools/go/ssa"
)

var ssaOnce sync.Once

type Gossa struct {
	rootPath string // 文件目录
	pagePath string // 文件名称
	name     string
	inited   bool

	c       *gossa.Context
	mainPkg *ssa.Package
}

func (T *Gossa) init() {
	if T.inited {
		return
	}
	ssaOnce.Do(func() {
		// 增加内置函数
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

func (T *Gossa) SetPath(root, page string) {
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

	ctx := gossa.NewContext(0)
	fset := token.NewFileSet()

	// 加载main文件
	file, err := parser.ParseFile(fset, T.name, script, ctx.ParserMode)
	if err != nil {
		return err
	}
	sPackage, err := T.loadAstFile(ctx, fset, file)
	if err != nil {
		return err
	}
	T.c = ctx
	T.mainPkg = sPackage
	return nil
}

func (T *Gossa) loadAstFile(ctx *gossa.Context, fset *token.FileSet, file *ast.File) (*ssa.Package, error) {
	pkg := types.NewPackage(file.Name.Name, "")
	gp := &gossaPackage{
		p:        T,
		ctx:      ctx,
		fset:     fset,
		packages: make(map[string]*types.Package),
		prog:     ssa.NewProgram(fset, ctx.BuilderMode),
	}

	ssapkg, err := gp.buildPackage(pkg, []*ast.File{file})
	if err != nil {
		return nil, err
	}
	ssapkg.Build()
	return ssapkg, nil
}

func (T *gossaPackage) buildPackage(pkg *types.Package, files []*ast.File) (*ssa.Package, error) {
	if pkg.Path() == "" {
		return nil, errors.New("package has no import path")
	}

	// 检查代码，该过程将会调用Import
	info := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Scopes:     make(map[ast.Node]*types.Scope),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
	}
	tc := &types.Config{
		Importer: gossa.NewImporter(T.ctx.Loader, T),
		Sizes:    T.ctx.Sizes,
	}
	if err := types.NewChecker(tc, T.fset, pkg, info).Files(files); err != nil {
		return nil, err
	}

	// 导入
	var (
		created   = make(map[*types.Package]bool)
		createAll func(pkgs []*types.Package)
	)
	createAll = func(pkgs []*types.Package) {
		for _, p := range pkgs {
			if !created[p] {
				created[p] = true
				if !p.Complete() {
					p.MarkComplete()
				}

				// 往内层继续导入模块
				createAll(p.Imports())

				// 生成ssa包-模块
				// 不用调用 Build() 因为这是内置模块
				if T.prog.ImportedPackage(p.Path()) == nil {
					T.prog.CreatePackage(p, nil, nil, true)
				}
			}
		}
	}
	// 导入模块
	createAll(pkg.Imports()) // 这里的内容由 SetImports 增加进去
	createAll(T.ctx.Loader.Packages())
	// 生成ssa包
	ssapkg := T.prog.CreatePackage(pkg, files, info, true)
	return ssapkg, nil
}

func (T *Gossa) Execute(out io.Writer, in interface{}) (err error) {
	if T.mainPkg == nil {
		return errors.New("the template has not been parsed and is not available")
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
	case string:
		io.WriteString(out, rv)
	case []byte:
		out.Write(rv)
	default:
		// 暂时不显示无法识别类型
		log.Printf("goass url(%s) returned unrecognized data type(%s)\n", T.pagePath, reflect.ValueOf(&rv).Elem().Elem().Type().String())
	}

	return nil
}

type gossaPackage struct {
	p        *Gossa
	fset     *token.FileSet
	ctx      *gossa.Context
	packages map[string]*types.Package
	prog     *ssa.Program
}

func (T *gossaPackage) imports(imports []*ast.ImportSpec) ([]*types.Package, error) {
	var tPackages []*types.Package
	for _, spec := range imports {
		p := spec.Path.Value
		p = p[1 : len(p)-1]
		// 先读取内置模块，再生成新模块
		tPackage, err := gossa.NewImporter(T.ctx.Loader, T).Import(p)
		if err != nil {
			return nil, err
		}
		tPackages = append(tPackages, tPackage)
	}
	return tPackages, nil
}

func (T *gossaPackage) Import(path string) (*types.Package, error) {
	// types.NewChecker 调用会检查模块
	if tPackage, ok := T.packages[path]; ok {
		return tPackage, nil
	}
	dirPath := filepath.Join(T.p.rootPath, "vendor", path)
	pkgs, err := parser.ParseDir(T.fset, dirPath, nil, T.ctx.ParserMode)
	if err != nil {
		return nil, err
	}

	// 创建模块
	for name, apkg := range pkgs {
		pkg := types.NewPackage(path, name)
		T.packages[path] = pkg

		var list []*types.Package
		var files []*ast.File
		for _, f := range apkg.Files {
			files = append(files, f)
			// 模块由多文件组成，每个文件有多个导入
			tPackages, err := T.imports(f.Imports)
			if err != nil {
				return nil, err
			}
			list = append(list, tPackages...)
		}
		// 模块需要加入子模块
		pkg.SetImports(list)
		// 生成ssa包
		ssaPackage, err := T.buildPackage(pkg, files)
		if err != nil {
			return nil, err
		}
		ssaPackage.Build()
		return pkg, nil
	}
	return nil, fmt.Errorf("(%s) package not code files", dirPath)
}
