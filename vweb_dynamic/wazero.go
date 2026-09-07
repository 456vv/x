package vweb_dynamic

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	witgo "github.com/456vv/x/vweb_dynamic/internal/witgo"
	"github.com/OpenListTeam/wazero-wasip2/wasip2"
	wasi_clocks "github.com/OpenListTeam/wazero-wasip2/wasip2/clocks"
	wasi_filesystem "github.com/OpenListTeam/wazero-wasip2/wasip2/filesystem"
	wasi_http "github.com/OpenListTeam/wazero-wasip2/wasip2/http"
	wasi_io "github.com/OpenListTeam/wazero-wasip2/wasip2/io"
	wasi_random "github.com/OpenListTeam/wazero-wasip2/wasip2/random"
	wasi_sockets "github.com/OpenListTeam/wazero-wasip2/wasip2/sockets"
	wasi_tls "github.com/OpenListTeam/wazero-wasip2/wasip2/tls"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"github.com/tetratelabs/wazero/sys"
)

// Wazero 封装 wazero 运行时，用于执行 WASI / Component 风格的 WASM 模板模块。
//
// 并发模型：
//   - 同一实例上多个 Execute 可并发（各自独立 Module 实例）。
//   - Close / Parse 与 Execute 互斥：Execute 持读锁贯穿整个生命周期，
//     Close/Parse 持写锁，保证不会对已关闭的 Runtime/CompiledModule 继续使用。
//   - Runtime 与 CompiledModule 本身支持并发 Instantiate。
type Wazero struct {
	mu        sync.RWMutex // 保护 runtime / cm / rootPath / pagePath / cache
	rootPath  string       // 文件目录（用于 FS 挂载与编译缓存）
	pagePath  string       // 文件路径（用于默认入口名推导）
	entryName string       // 保留字段（与原结构兼容；入口名仍由 Execute 参数与 entryname 推导）
	runtime   wazero.Runtime
	cm        wazero.CompiledModule
	cache     wazero.CompilationCache // 目录缓存时持有，Close 时释放（Runtime.Close 不会关闭共享 cache）
}

// Parse 从 io.Reader 解析并编译 WASM。
// 若此前已解析过，会先安全关闭旧 Runtime 再编译新模块（可重入）。
func (T *Wazero) Parse(r io.Reader) (err error) {
	if r == nil {
		return errors.New("reader is nil")
	}
	contact, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	return T.parse(contact)
}

// ParseFile 从文件路径解析并编译 WASM。
func (T *Wazero) ParseFile(p string) error {
	if p == "" {
		return errors.New("path is empty")
	}
	b, err := os.ReadFile(p)
	if err != nil {
		return err
	}
	return T.parse(b)
}

// ParseText 从字符串内容解析（实现原先缺失的功能）。
func (T *Wazero) ParseText(name, content string) error {
	return T.parse([]byte(content))
}

// SetPath 设置根目录与页面路径（用于 FS 挂载与默认入口名）。
// 并发安全：可与 Execute 并发调用（读侧使用快照）。
func (T *Wazero) SetPath(root, page string) {
	T.mu.Lock()
	T.rootPath = root
	T.pagePath = page
	T.mu.Unlock()
}

// parse 核心解析与编译逻辑。
// 创建带编译缓存的 Runtime，实例化 WASIp2 Host，编译用户 WASM，并导出 TemplateFunc 到 "env" 模块。
// 可重入：若已有 Runtime 会先关闭再重建。
func (T *Wazero) parse(calcWasm []byte) error {
	if len(calcWasm) == 0 {
		return errors.New("wasm binary is empty")
	}

	T.mu.Lock()
	defer T.mu.Unlock()

	// 先关闭旧资源，避免泄漏
	_ = T.closeLocked()

	ctx := context.Background()

	// 编译缓存：优先使用 rootPath 磁盘缓存，失败则降级为纯内存缓存（跨平台安全）
	var cache wazero.CompilationCache
	if T.rootPath != "" {
		root := filepath.Clean(T.rootPath)
		if c, err := wazero.NewCompilationCacheWithDir(filepath.Join(root, ".cache")); err == nil {
			cache = c
		}
		// 失败时忽略，继续使用 Runtime 默认的内存缓存
	}

	config := wazero.NewRuntimeConfig().
		WithDebugInfoEnabled(true).
		WithCloseOnContextDone(true).
		WithCoreFeatures(api.CoreFeaturesV2)
	if cache != nil {
		config = config.WithCompilationCache(cache)
	}

	runtime := wazero.NewRuntimeWithConfig(ctx, config)

	// 兼容旧 WASI preview1
	wasi_snapshot_preview1.MustInstantiate(ctx, runtime)

	// 创建 wasip2.Host，并启用所需模块
	h := wasip2.NewHost(
		wasi_io.Module("0.2.0"),
		wasi_random.Module("0.2.0"),
		wasi_clocks.Module("0.2.0"),
		wasi_filesystem.Module("0.2.0"),
		wasi_sockets.Module("0.2.0"),
		wasi_tls.Module("0.2.0"),
		wasi_http.Module("0.2.0"),
	)
	if err := h.Instantiate(ctx, runtime); err != nil {
		_ = runtime.Close(ctx)
		if cache != nil {
			_ = cache.Close(ctx)
		}
		return fmt.Errorf("instantiate wasip2.Host: %w", err)
	}

	cm, err := runtime.CompileModule(ctx, calcWasm)
	if err != nil {
		_ = runtime.Close(ctx)
		if cache != nil {
			_ = cache.Close(ctx)
		}
		return fmt.Errorf("compile wasm module: %w", err)
	}

	// 导出 TemplateFunc 到 "env" 宿主模块并 Instantiate（必须 Instantiate 才能被 Guest 链接）
	builder := runtime.NewHostModuleBuilder("env")
	exporter := witgo.NewExporter(builder)
	for k, v := range TemplateFunc {
		if err := exporter.Export(k, v); err != nil {
			//log.Printf("export template func %q: %v", k, err)
		}
	}
	if _, err := builder.Instantiate(ctx); err != nil {
		_ = runtime.Close(ctx)
		if cache != nil {
			_ = cache.Close(ctx)
		}
		return fmt.Errorf("instantiate env host module: %w", err)
	}

	T.runtime = runtime
	T.cm = cm
	T.cache = cache
	return nil
}

// Execute 实例化已编译模块并调用入口函数，将结果写入 out。
// 每次调用都会创建独立的 Module 实例，保证隔离与并发安全。
// 同一 Wazero 实例上的多个 Execute 可并发执行；与 Close/Parse 互斥。
func (T *Wazero) Execute(entryName string, out io.Writer, in ...any) (err error) {
	if out == nil {
		return errors.New("out writer is nil")
	}

	// 读锁贯穿整个执行过程，避免 Close/Parse 在中途释放 Runtime/cm 导致 use-after-close。
	T.mu.RLock()
	defer T.mu.RUnlock()

	if T.runtime == nil || T.cm == nil {
		return errTemplateNotParse
	}

	ctx := context.Background()

	// ModuleConfig：
	// - FS 使用 WithDirMount（WASI 路径/权限语义更完整，跨平台更一致）。
	// - Stdout/Stderr 默认 Discard：共享 os.Stdout/Stderr 破坏沙箱且阻止真正并发。
	//   模板主输出走 host.Call → body → out，不依赖 WASI stdout。
	// - 不向 Guest 泄漏宿主机 os.Args。
	fsConfig := wazero.NewFSConfig()
	if T.rootPath != "" {
		root := filepath.Clean(T.rootPath)
		fsConfig = fsConfig.WithDirMount(root, "/")
	}

	config := wazero.NewModuleConfig().
		WithFSConfig(fsConfig).
		WithName(""). // 匿名实例，避免名称冲突，支持并发 Instantiate
		WithWalltime(func() (sec int64, nsec int32) {
			t := time.Now()
			return t.Unix(), int32(t.Nanosecond())
		}, sys.ClockResolution(time.Microsecond.Nanoseconds())).
		WithStdout(io.Discard).
		WithStderr(io.Discard).
		WithRandSource(rand.Reader).
		WithArgs(os.Args...)

	mod, err := T.runtime.InstantiateModule(ctx, T.cm, config)
	if err != nil {
		return fmt.Errorf("instantiate module: %w", err)
	}
	defer mod.Close(ctx)

	// 优先调用 Component 风格的 _initialize；兼容旧 command 模块的 _start。
	// _start 正常结束常返回 *sys.ExitError（code==0），需视为成功。
	if initFn := mod.ExportedFunction("_initialize"); initFn != nil {
		if _, err := initFn.Call(ctx); err != nil {
			return fmt.Errorf("failed to call _initialize: %w", err)
		}
	} else if startFn := mod.ExportedFunction("_start"); startFn != nil {
		if _, err := startFn.Call(ctx); err != nil {
			if exitErr, ok := err.(*sys.ExitError); ok {
				if exitErr.ExitCode() != 0 {
					return fmt.Errorf("module _start exited with code %d: %w", exitErr.ExitCode(), err)
				}
				// code == 0：正常退出，继续执行入口函数
			} else {
				return fmt.Errorf("failed to call _start: %w", err)
			}
		}
	}

	host, err := witgo.NewHost(mod)
	if err != nil {
		return fmt.Errorf("new witgo host: %w", err)
	}

	entryName = entryname(entryName, T.pagePath)
	// 输出缓冲区由 witgo 按需扩容，避免固定 1024 字节截断
	var body []byte
	if err := host.Call(ctx, entryName, &body, in...); err != nil {
		return fmt.Errorf("call entry %q: %w", entryName, err)
	}

	if len(body) > 0 {
		if _, err := out.Write(body); err != nil {
			return fmt.Errorf("write body: %w", err)
		}
	}
	return nil
}

// Close 关闭底层 Runtime 与编译缓存（幂等）。关闭后不应再调用 Execute。
// 会等待当前正在进行的 Execute 结束后再关闭，保证并发安全。
func (T *Wazero) Close() error {
	T.mu.Lock()
	defer T.mu.Unlock()
	return T.closeLocked()
}

// closeLocked 在已持有写锁的前提下关闭资源。
// 顺序：先关 Runtime（会释放其编译产物），再关 CompilationCache。
func (T *Wazero) closeLocked() error {
	var firstErr error
	if T.runtime != nil {
		if err := T.runtime.Close(context.Background()); err != nil && firstErr == nil {
			firstErr = err
		}
		T.runtime = nil
		// Runtime.Close 会关闭其编译出的 CompiledModule，无需再单独 Close cm。
		T.cm = nil
	}
	if T.cache != nil {
		// CompilationCache 生命周期独立于 Runtime，需显式关闭。
		if err := T.cache.Close(context.Background()); err != nil && firstErr == nil {
			firstErr = err
		}
		T.cache = nil
	}
	return firstErr
}
