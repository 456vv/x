package vweb_dynamic

import (
	"bytes"
	"fmt"
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

// Ixgo 动态页面执行引擎（基于 ixgo 解释器）
type Ixgo struct {
	rootPath  string // 文件根目录
	pagePath  string // 当前页面文件路径
	entryName string // 入口函数名

	ctx     *ixgo.Context
	mainPkg *ssa.Package

	mu sync.RWMutex // 保护 ctx / mainPkg / entryName 等可变状态
}

// init 注册内置函数到 ixgo（只执行一次）
func (T *Ixgo) init() {
	ixgoOnce.Do(func() {
		// 增加内置函数
		for name, fn := range TemplateFunc {
			ixgo.RegisterExternal(name, fn)
		}
	})
}

// SetPath 设置根目录与页面路径
func (T *Ixgo) SetPath(root, page string) {
	T.mu.Lock()
	T.rootPath = filepath.Clean(root)
	T.pagePath = filepath.Clean(page)
	T.mu.Unlock()
}

// setHeaderLine 解析缓冲区头部的模板头注释
func (T *Ixgo) setHeaderLine(buf *bytes.Buffer) TemplateHeader {
	l := fileHeaderLine(buf)
	return templateHeader(l)
}

// ParseText 从字符串解析动态页面
func (T *Ixgo) ParseText(name, content string) error {
	buf := bytes.NewBufferString(content)
	return T.parse(name, buf)
}

// ParseFile 从文件解析动态页面
func (T *Ixgo) ParseFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return T.Parse(file)
}

// Parse 从 Reader 解析动态页面
func (T *Ixgo) Parse(r io.Reader) (err error) {
	buf := bytes.NewBuffer(nil)
	if _, err = buf.ReadFrom(r); err != nil {
		return err
	}
	return T.parse("main.go", buf)
}

// parse 核心解析逻辑
func (T *Ixgo) parse(filename string, buf *bytes.Buffer) error {
	T.init()

	T.mu.Lock()
	defer T.mu.Unlock()

	ctx := ixgo.NewContext(ixgo.EnableNoStrict)
	// 加载外部模块
	ctx.Lookup = T.lookup

	th := T.setHeaderLine(buf)
	T.entryName = th.EntryName

	apkg := &ast.Package{
		Name:  "main",
		Files: map[string]*ast.File{},
	}

	// 加载扩展文件
	dir := filepath.Dir(T.pagePath)
	for _, f := range th.File {
		var p string
		if path.IsAbs(f) {
			p = filepath.Join(T.rootPath, f)
		} else {
			p = filepath.Join(T.rootPath, dir, f)
		}
		p = filepath.Clean(p)

		file, err := ctx.ParseFile(p, nil)
		if err != nil {
			return fmt.Errorf("parse extend file %q: %w", p, err)
		}
		apkg.Files[f] = file
	}

	// 加载main文件
	file, err := ctx.ParseFile(filename, buf.String())
	if err != nil {
		return fmt.Errorf("parse main file %q: %w", filename, err)
	}
	apkg.Files[filename] = file

	sPkg, err := ctx.LoadAstPackage(apkg.Name, apkg)
	if err != nil {
		return fmt.Errorf("load ast package: %w", err)
	}

	T.ctx = ctx
	T.mainPkg = sPkg
	return nil
}

// lookup 查找外部模块目录（用于 ixgo 加载）
func (T *Ixgo) lookup(root, path string) (dir string, found bool) {
	dir = filepath.Join(T.rootPath, "src", filepath.FromSlash(path))
	if info, err := os.Stat(dir); err == nil && info.IsDir() {
		found = true
	}
	return
}

// Execute 执行入口函数并将结果写入 out
func (T *Ixgo) Execute(entryName string, out io.Writer, in ...any) (err error) {
	T.mu.RLock()
	mainPkg := T.mainPkg
	ctx := T.ctx
	defaultEntry := T.entryName
	T.mu.RUnlock()

	if mainPkg == nil || ctx == nil {
		return errTemplateNotParse
	}

	entryName = entryname(defaultEntry, entryName)
	retn, err := ctx.RunFunc(mainPkg, entryName, in...)
	if err != nil {
		return err
	}

	switch rv := retn.(type) {
	case string:
		_, err = io.WriteString(out, rv)
	case []byte:
		_, err = out.Write(rv)
	case io.Reader:
		_, err = io.Copy(out, rv)
	case fmt.Stringer:
		_, err = io.WriteString(out, rv.String())
	case nil:
		// 无返回值，正常
	default:
		// 暂时不显示无法识别类型
		log.Printf("ixgo call %s returned unrecognized data type(%s)\n", entryName, reflect.TypeOf(rv).String())
		return fmt.Errorf("ixgo call %s returned unrecognized data type: %T", entryName, rv)
	}
	return err
}

// Close 释放资源（当前 ixgo 无显式释放需求，预留接口）
func (T *Ixgo) Close() error {
	T.mu.Lock()
	defer T.mu.Unlock()
	T.ctx = nil
	T.mainPkg = nil
	return nil
}
