package vweb_dynamic

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sync"

	_ "github.com/456vv/x/igop_lib"
	"github.com/goplus/igop"
	_ "github.com/goplus/igop/pkg"                 // 加入默认包
	_ "github.com/goplus/reflectx/icall/icall2048" // 内置
	"golang.org/x/tools/go/ssa"
)

var igopOnce sync.Once

type Igop struct {
	rootPath  string // 文件目录
	pagePath  string // 文件名称
	entryName string

	ctx     *igop.Context
	mainPkg *ssa.Package
}

func (T *Igop) init() {
	igopOnce.Do(func() {
		// 增加内置函数
		for name, fn := range TemplateFunc {
			igop.RegisterExternal(name, fn)
		}
	})
}

func (T *Igop) SetPath(root, page string) {
	T.rootPath = root
	T.pagePath = page
}

func (T *Igop) SetEntryName(name string) {
	T.entryName = name
}

func (T *Igop) setHeaderLine(l []string) {
	hm := headerMap(l)
	if T.entryName == "" && len(hm["entryname"]) > 0 {
		T.entryName = hm["entryname"][0]
	}
}

func (T *Igop) ParseText(name, content string) error {
	buf := bytes.NewBufferString(content)
	l := fileHeaderLine(buf)
	T.setHeaderLine(l)

	return T.parse(name, buf.String())
}

func (T *Igop) ParseFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return T.Parse(file)
}

func (T *Igop) Parse(r io.Reader) (err error) {
	buf := bytes.NewBuffer(nil)
	buf.ReadFrom(r)
	l := fileHeaderLine(buf)
	T.setHeaderLine(l)

	return T.parse("", buf.String())
}

func (T *Igop) parse(filename string, src interface{}) error {
	T.init()

	ctx := igop.NewContext(igop.EnableNoStrict)
	// 加载外部模块
	ctx.Lookup = T.lookup

	// 加载main文件
	file, err := ctx.ParseFile(filename, src)
	if err != nil {
		return err
	}
	sPkg, err := ctx.LoadAstFile(file.Name.Name, file)
	if err != nil {
		return err
	}
	T.ctx = ctx
	T.mainPkg = sPkg
	return nil
}

func (T *Igop) lookup(root, path string) (dir string, found bool) {
	dir = filepath.Join(T.rootPath, "src", path)
	if info, err := os.Stat(dir); err == nil && info.IsDir() {
		found = true
	}
	return
}

func (T *Igop) Execute(out io.Writer, in interface{}) (err error) {
	if T.mainPkg == nil {
		return errTemplateNotParse
	}

	// 调用Main函数
	if T.entryName == "" {
		T.entryName = entryname(T.pagePath)
	}
	retn, err := T.ctx.RunFunc(T.mainPkg, T.entryName, in)
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
		log.Printf("igop url(%s) returned unrecognized data type(%s)\n", T.pagePath, reflect.ValueOf(rv).Elem().Type().String())
	}

	return nil
}
