package vweb_dynamic

import (
	"bytes"
	"go/ast"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"sync"

	_ "github.com/456vv/x/igop_lib"
	"github.com/goplus/ixgo"
	_ "github.com/goplus/ixgo/pkg"                 // 加入默认包
	_ "github.com/goplus/reflectx/icall/icall2048" // 内置
	"golang.org/x/tools/go/ssa"
)

var ixgoOnce sync.Once

type Ixgo struct {
	rootPath  string // 文件目录
	pagePath  string // 文件名称
	entryName string

	ctx     *ixgo.Context
	mainPkg *ssa.Package
}

func (T *Ixgo) init() {
	ixgoOnce.Do(func() {
		// 增加内置函数
		for name, fn := range TemplateFunc {
			ixgo.RegisterExternal(name, fn)
		}
	})
}

func (T *Ixgo) SetPath(root, page string) {
	T.rootPath = root
	T.pagePath = page
}

func (T *Ixgo) setHeaderLine(buf *bytes.Buffer) TemplateHeader {
	l := fileHeaderLine(buf)
	return templateHeader(l)
}

func (T *Ixgo) ParseText(name, content string) error {
	buf := bytes.NewBufferString(content)
	return T.parse(name, buf)
}

func (T *Ixgo) ParseFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return T.Parse(file)
}

func (T *Ixgo) Parse(r io.Reader) (err error) {
	buf := bytes.NewBuffer(nil)
	buf.ReadFrom(r)
	return T.parse("main.go", buf)
}

func (T *Ixgo) parse(filename string, buf *bytes.Buffer) error {
	T.init()

	ctx := ixgo.NewContext(ixgo.EnableNoStrict)
	// 加载外部模块
	ctx.Lookup = T.lookup

	th := T.setHeaderLine(buf)
	T.entryName = th.EntryName

	var (
		err  error
		file *ast.File
	)
	apkg := &ast.Package{
		Name:  "main",
		Files: map[string]*ast.File{},
	}

	// 加载扩展文件
	dir := filepath.Dir(T.pagePath)
	for _, f := range th.File {
		if path.IsAbs(f) {
			p := filepath.Join(T.rootPath, f)
			file, err = ctx.ParseFile(p, nil)
		} else {
			p := filepath.Join(T.rootPath, dir, f)
			file, err = ctx.ParseFile(p, nil)
		}
		if err != nil {
			return err
		}
		apkg.Files[f] = file
	}

	// 加载main文件
	file, err = ctx.ParseFile(filename, buf.String())
	if err != nil {
		return err
	}
	apkg.Files[filename] = file

	sPkg, err := ctx.LoadAstPackage(apkg.Name, apkg)
	if err != nil {
		return err
	}

	T.ctx = ctx
	T.mainPkg = sPkg
	return nil
}

func (T *Ixgo) lookup(root, path string) (dir string, found bool) {
	dir = filepath.Join(T.rootPath, "src", path)
	if info, err := os.Stat(dir); err == nil && info.IsDir() {
		found = true
	}
	return
}

func (T *Ixgo) Execute(entryName string, out io.Writer, in ...any) (err error) {
	if T.mainPkg == nil {
		return errTemplateNotParse
	}

	entryName = entryname(T.entryName, entryName)
	retn, err := T.ctx.RunFunc(T.mainPkg, entryName, in...)
	if err != nil {
		return err
	}

	switch rv := retn.(type) {
	case string:
		io.WriteString(out, rv)
	case []byte:
		out.Write(rv)
	case io.Reader:
		io.Copy(out, rv)
	default:
		// 暂时不显示无法识别类型
		log.Printf("ixgo call %s returned unrecognized data type(%s)\n", entryName, reflect.ValueOf(rv).Type().String())
	}

	return nil
}

func (T *Ixgo) Close() error {
	return nil
}
