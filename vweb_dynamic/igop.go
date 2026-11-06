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
	"sync/atomic"

	_ "github.com/456vv/x/igop_lib"
	"github.com/goplus/ixgo"
	_ "github.com/goplus/ixgo/pkg"                 // 加入默认包
	_ "github.com/goplus/reflectx/icall/icall2048" // 内置
	"golang.org/x/tools/go/ssa"
)

var ixgoOnce sync.Once

// ixgoRuntime 是 parse 成功后发布的不可变运行快照：
// 解释器上下文、SSA 包、入口函数名，以及串行化并发执行的 execMu。
//
// execMu 串行化 RunFunc：ixgo.Context 内部的解释器运行时状态（栈、临时对象表等）
// 不是并发安全的；同一快照上多个 Execute 必须串行执行 RunFunc。
type ixgoRuntime struct {
	ctx       *ixgo.Context
	mainPkg   *ssa.Package
	entryName string
	execMu    sync.Mutex
}

// Ixgo 动态页面执行引擎（基于 ixgo 解释器）。
//
// 并发安全说明：
//   - SetPath 与 Execute 可并发；Execute 使用快照，不受路径切换影响。
//   - Parse / ParseFile / ParseText 成功后原子发布新快照；失败保留旧快照，
//     已开始的 Execute 继续持有旧快照直到完成，不受影响。
//   - Close 清空快照；已在执行的 Execute 继续持有旧快照直到完成。
//   - 同一快照上的多个 Execute 通过 ixgoRuntime.execMu 串行执行 RunFunc。
type Ixgo struct {
	rootPath string // 文件根目录
	pagePath string // 当前页面文件路径

	mu  sync.RWMutex                // 保护 rootPath / pagePath
	cur atomic.Pointer[ixgoRuntime] // 原子发布当前运行快照；nil 表示未解析/已关闭
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

// parse 核心解析逻辑。
// 解析成功后通过 cur.Store 原子发布新快照；失败时不修改 cur（保留旧快照）。
func (T *Ixgo) parse(filename string, buf *bytes.Buffer) error {
	T.init()

	// 读取 rootPath/pagePath 快照，parse 全程使用局部变量；
	// 与 SetPath 并发时不需持有 Ixgo.mu，也不再串行化整个解析过程。
	T.mu.RLock()
	root := T.rootPath
	page := T.pagePath
	T.mu.RUnlock()

	ctx := ixgo.NewContext(ixgo.EnableNoStrict)
	// Lookup 闭包捕获 parse 时的 root 快照，避免 SetPath 后续切换导致旧 ctx 加载新路径。
	ctx.Lookup = func(_, p string) (dir string, found bool) {
		return T.lookup(root, p)
	}

	th := T.setHeaderLine(buf)
	entryName := th.EntryName

	apkg := &ast.Package{
		Name:  "main",
		Files: map[string]*ast.File{},
	}

	dir := filepath.Dir(page)
	for _, f := range th.File {
		var p string
		if path.IsAbs(f) {
			p = filepath.Join(root, f)
		} else {
			p = filepath.Join(root, dir, f)
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

	T.cur.Store(&ixgoRuntime{
		ctx:       ctx,
		mainPkg:   sPkg,
		entryName: entryName,
	})
	return nil
}

// lookup 查找外部模块目录。
// root 由 parse 闭包传入（解析时快照），不读取 T.rootPath，因此并发安全。
func (T *Ixgo) lookup(root, path string) (dir string, found bool) {
	if root == "" {
		return "", false
	}
	dir = filepath.Join(root, "src", filepath.FromSlash(path))
	if info, err := os.Stat(dir); err == nil && info.IsDir() {
		found = true
	}
	return
}

// Execute 执行入口函数并将结果写入 out。
// 入口处加载一次 cur 快照并全程持有；Parse/Close 并发切换不影响本次调用。
func (T *Ixgo) Execute(entryName string, out io.Writer, in ...any) (err error) {
	rt := T.cur.Load()
	if rt == nil {
		return errTemplateNotParse
	}

	entryName = entryname(rt.entryName, entryName)

	// RunFunc 修改 ixgo.Context 内部状态，非并发安全；串行化同一快照上的执行。
	// 局部函数 + defer 保证 RunFunc panic 时锁必释放，避免后续 Execute 永久阻塞。
	retn, err := func() (any, error) {
		rt.execMu.Lock()
		defer rt.execMu.Unlock()
		return rt.ctx.RunFunc(rt.mainPkg, entryName, in...)
	}()
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

// Close 释放资源（当前 ixgo 无显式释放需求，预留接口）。
// 清空 cur 快照；已在执行的 Execute 继续持有旧快照直到完成。
func (T *Ixgo) Close() error {
	T.cur.Store(nil)
	return nil
}
