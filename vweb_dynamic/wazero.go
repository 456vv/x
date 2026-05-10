package vweb_dynamic

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"time"
	"unsafe"

	"github.com/OpenListTeam/wazero-wasip2/wasip2"
	wasi_clocks "github.com/OpenListTeam/wazero-wasip2/wasip2/clocks"
	wasi_filesystem "github.com/OpenListTeam/wazero-wasip2/wasip2/filesystem"
	wasi_http "github.com/OpenListTeam/wazero-wasip2/wasip2/http"
	wasi_io "github.com/OpenListTeam/wazero-wasip2/wasip2/io"
	wasi_random "github.com/OpenListTeam/wazero-wasip2/wasip2/random"
	wasi_sockets "github.com/OpenListTeam/wazero-wasip2/wasip2/sockets"
	wasi_tls "github.com/OpenListTeam/wazero-wasip2/wasip2/tls"
	witgo "github.com/OpenListTeam/wazero-wasip2/wit-go"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"github.com/tetratelabs/wazero/sys"
)

type Wazero struct {
	rootPath  string // 文件目录
	pagePath  string // 文件路径
	entryName string
	runtime   wazero.Runtime
	cm        wazero.CompiledModule
}

func (T *Wazero) ParseText(name, content string) error {
	return errors.New("not implemented yet")
}

func (T *Wazero) ParseFile(p string) error {
	b, err := os.ReadFile(p)
	if err != nil {
		return err
	}
	return T.parse(b)
}

func (T *Wazero) SetPath(root, page string) {
	T.rootPath = root
	T.pagePath = page
}

func (T *Wazero) Parse(r io.Reader) (err error) {
	contact, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	return T.parse(contact)
}

func (T *Wazero) parse(calcWasm []byte) error {
	// 创建共享编译缓存目录的运行时配置
	cache, err := wazero.NewCompilationCacheWithDir(T.rootPath)
	if err != nil {
		return err
	}

	// 使用相同的 wazero.CompilationCache 实例允许内存中缓存共享。
	config := wazero.NewRuntimeConfig().
		WithCompilationCache(cache).
		WithDebugInfoEnabled(true).
		WithCloseOnContextDone(true).
		WithCoreFeatures(api.CoreFeaturesV2)
	ctx := context.Background()
	runtime := wazero.NewRuntimeWithConfig(ctx, config)

	wasi_snapshot_preview1.MustInstantiate(ctx, runtime)
	// 创建我们的 wasip2.Host 实例，并启用所有需要的模块。
	h := wasip2.NewHost(
		wasi_io.Module("0.2.0"),
		wasi_random.Module("0.2.0"),
		wasi_clocks.Module("0.2.0"),
		wasi_filesystem.Module("0.2.0"),
		wasi_sockets.Module("0.2.0"),
		wasi_tls.Module("0.2.0"),
		wasi_http.Module("0.2.0"),
	)
	// 将我们的 Host 实现实例化到 wazero 运行时中。
	err = h.Instantiate(ctx, runtime)
	if err != nil {
		return fmt.Errorf("instantiate wasip2.Host: %v", err)
	}

	cm, err := runtime.CompileModule(ctx, calcWasm)
	if err != nil {
		return fmt.Errorf("compile wasm module: %v", err)
	}

	exporter := witgo.NewExporter(runtime.NewHostModuleBuilder("env"))
	for k, v := range TemplateFunc {
		exporter.Export(k, v)
	}

	T.runtime = runtime
	T.cm = cm

	return nil
}

func (T *Wazero) Execute(entryName string, out io.Writer, in ...any) (err error) {
	if T.runtime == nil {
		return errTemplateNotParse
	}

	ctx := context.Background()
	config := wazero.NewModuleConfig().
		WithFSConfig(wazero.NewFSConfig().WithFSMount(os.DirFS(T.rootPath), "/")).
		WithName("").
		WithWalltime(func() (sec int64, nsec int32) {
			t := time.Now()
			return t.Unix(), int32(t.Nanosecond())
		}, sys.ClockResolution(time.Microsecond.Nanoseconds())).
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		WithRandSource(rand.Reader).
		WithArgs(os.Args...)

	mod, err := T.runtime.InstantiateModule(ctx, T.cm, config)
	if err != nil {
		return err
	}
	defer mod.Close(ctx)

	if initFn := mod.ExportedFunction("_initialize"); initFn != nil {
		if _, err := initFn.Call(ctx); err != nil {
			return fmt.Errorf("failed to call _initialize: %w", err)
		}
	} else {
		// 兼容旧的 command 模块（_start），但 reactor 模式下通常不需要
		if startFn := mod.ExportedFunction("_start"); startFn != nil {
			if _, err := startFn.Call(ctx); err != nil {
				// _start 可能返回错误或直接退出，这里根据需要处理
				return fmt.Errorf("warning: _start returned: %v", err)
			}
		}
	}
	host, err := witgo.NewHost(mod)
	if err != nil {
		return err
	}

	entryName = entryname(entryName, T.pagePath)
	var body string
	if err := host.Call(ctx, entryName, &body, in...); err != nil {
		return err
	}

	io.WriteString(out, body)

	return nil
}

func (T *Wazero) Close() error {
	if T.runtime != nil {
		T.runtime.Close(context.Background())
	}
	return nil
}

func toPointer(a any) uint64 {
	return uint64(reflect.ValueOf(a).Pointer())
}

func read(offset unsafe.Pointer, dlen uint32) ([]byte, bool) {
	return unsafe.Slice((*byte)(offset), dlen), true
}

func write(offset unsafe.Pointer, v []byte) bool {
	dst := unsafe.Slice((*byte)(offset), len(v))
	copy(dst, v)
	return true
}
