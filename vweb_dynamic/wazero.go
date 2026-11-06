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

	"github.com/456vv/x/vweb_dynamic/internal/wazero/wasip2"
	wasi_clocks "github.com/456vv/x/vweb_dynamic/internal/wazero/wasip2/clocks"
	wasi_filesystem "github.com/456vv/x/vweb_dynamic/internal/wazero/wasip2/filesystem"
	wasi_http "github.com/456vv/x/vweb_dynamic/internal/wazero/wasip2/http"
	wasi_io "github.com/456vv/x/vweb_dynamic/internal/wazero/wasip2/io"
	wasi_random "github.com/456vv/x/vweb_dynamic/internal/wazero/wasip2/random"
	wasi_sockets "github.com/456vv/x/vweb_dynamic/internal/wazero/wasip2/sockets"
	wasi_tls "github.com/456vv/x/vweb_dynamic/internal/wazero/wasip2/tls"
	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"github.com/tetratelabs/wazero/sys"
)

// Wazero 封装 wazero 运行时，用于执行 WASI / Component 风格的 WASM 模板模块。
//
// 并发模型：
//   - 同一实例上多个 Execute / ExecuteContext 可并发（各自独立 Module 实例）。
//   - Close / Parse 与 Execute* 互斥：Execute* 持读锁贯穿整个生命周期，
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

	// host 是 wasip2.Host 实例，持有 filesystem/http/tls/sockets/stream 等管理器。
	// parse 重入与 Close 时必须显式关闭，否则 *os.File / net.Conn / http.Transport
	// 空闲连接等真实资源不会被释放。
	host *wasip2.Host

	// envTable 是与 Runtime 同寿命的句柄表，承载 TemplateFunc 与 root 对象。
	envTable *envHandleTable

	// RootAsFirstArg 控制 Execute 的 in[0] 与 host.Call 的关系（默认 false）。
	//
	//	false（默认，兼容旧路径）：
	//	  - in 全部传给 host.Call（须全部可被 witgo 扁平化：标量/string/slice 等）
	//	  - 若 in[0] != nil，同时 box 为 root，供 tf_root_handle（旧 guest 不调用则无影响）
	//	  - 若 in[0] 是 struct/map/任意复杂 any，host.Call 仍会扁平化失败——与改前行为一致
	//
	//	true（模板点对象推荐）：
	//	  - in[0] 仅 box 为 root，不传给 host.Call
	//	  - in[1:] 传给 host.Call（应为可扁平类型；无额外参数时入口为 0 形参）
	//	  - guest 用 tf_root_handle() 取点对象
	RootAsFirstArg bool
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
// 创建带编译缓存的 Runtime，实例化 WASIp2 Host，编译用户 WASM，并导出句柄式 env 模块。
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

	// 编译缓存：优先使用 rootPath 磁盘缓存，失败则降级为纯内存缓存
	var cache wazero.CompilationCache
	if T.rootPath != "" {
		root := filepath.Clean(T.rootPath)
		if c, err := wazero.NewCompilationCacheWithDir(filepath.Join(root, ".cache")); err == nil {
			cache = c
		}
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
		h.Close()
		if cache != nil {
			_ = cache.Close(ctx)
		}
		return fmt.Errorf("instantiate wasip2.Host: %w", err)
	}

	if T.rootPath != "" {
		err := registerPreopens(h, map[string]string{"/": filepath.Clean(T.rootPath)})
		if err != nil {
			_ = runtime.Close(ctx)
			h.Close()
			if cache != nil {
				_ = cache.Close(ctx)
			}
			return err
		}
	}
	cm, err := runtime.CompileModule(ctx, calcWasm)
	if err != nil {
		_ = runtime.Close(ctx)
		h.Close()
		if cache != nil {
			_ = cache.Close(ctx)
		}
		return fmt.Errorf("compile wasm module: %w", err)
	}

	// 导出 env：句柄式 TemplateFunc 桥
	builder := runtime.NewHostModuleBuilder("env")
	table := newEnvHandleTable()
	registerEnvTemplateAPI(builder, table)
	if _, err := builder.Instantiate(ctx); err != nil {
		_ = runtime.Close(ctx)
		h.Close()
		if cache != nil {
			_ = cache.Close(ctx)
		}
		return fmt.Errorf("instantiate env host module: %w", err)
	}

	T.runtime = runtime
	T.cm = cm
	T.cache = cache
	T.host = h
	T.envTable = table
	return nil
}

// Execute 是 ExecuteContext 的便捷包装，使用 context.Background()，
// 无法从外部取消；执行期间不会被请求级超时或取消中断。
func (T *Wazero) Execute(entryName string, out io.Writer, in ...any) error {
	return T.ExecuteContext(context.Background(), entryName, out, in...)
}

// ExecuteContext 同 Execute，但支持请求级 ctx 取消与超时。
//
// 依赖 T.parse 中 RuntimeConfig.WithCloseOnContextDone(true)：
// 当 ctx 被取消或超时时，wazero 会强制关闭当前 WASM 实例，
// 并让所有挂起的 Guest 调用返回 context.Canceled / context.DeadlineExceeded。
//
// in[0] 的语义由 T.RootAsFirstArg 控制（见 Wazero 结构体注释）。
func (T *Wazero) ExecuteContext(ctx context.Context, entryName string, out io.Writer, in ...any) (err error) {
	if out == nil {
		return errors.New("out writer is nil")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if ctxErr := ctx.Err(); ctxErr != nil {
		return ctxErr
	}

	T.mu.RLock()
	defer T.mu.RUnlock()

	if T.runtime == nil || T.cm == nil {
		return errTemplateNotParse
	}

	fsConfig := wazero.NewFSConfig()
	if T.rootPath != "" {
		root := filepath.Clean(T.rootPath)
		fsConfig = fsConfig.WithDirMount(root, "/")
	}

	config := wazero.NewModuleConfig().
		WithFSConfig(fsConfig).
		WithName("").
		WithWalltime(func() (sec int64, nsec int32) {
			t := time.Now()
			return t.Unix(), int32(t.Nanosecond())
		}, sys.ClockResolution(time.Microsecond.Nanoseconds())).
		WithStdout(io.Discard).
		WithStderr(io.Discard).
		WithRandSource(rand.Reader).
		WithArgs("vweb_dynamic")

	mod, err := T.runtime.InstantiateModule(ctx, T.cm, config)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return fmt.Errorf("instantiate module cancelled: %w", ctxErr)
		}
		return fmt.Errorf("instantiate module: %w", err)
	}
	// 用独立 ctx 关闭模块：请求 ctx 已取消时 mod.Close(ctx) 会立刻返回而不做清理。
	defer mod.Close(context.Background())
	// 清 witgo hostCache：mod 关闭后不删除会随 module 数量线性泄漏。
	defer witgo.ReleaseHost(mod)

	if T.envTable != nil {
		defer T.envTable.setErr(mod, nil)
		defer T.envTable.releaseModule(mod)
	}

	// root 对象：仅在有 in[0] 时 box/清理
	if T.envTable != nil && len(in) > 0 && in[0] != nil {
		if rootHandle := T.envTable.put(in[0]); rootHandle != 0 {
			T.envTable.setRoot(mod, rootHandle)
			defer func() {
				T.envTable.clearRoot(mod)
				T.envTable.drop(rootHandle)
			}()
		}
	}

	// 优先调用 Component 风格的 _initialize；兼容旧 command 模块的 _start。
	if initFn := mod.ExportedFunction("_initialize"); initFn != nil {
		if _, cerr := initFn.Call(ctx); cerr != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return fmt.Errorf("call _initialize cancelled: %w", ctxErr)
			}
			return fmt.Errorf("failed to call _initialize: %w", cerr)
		}
	} else if startFn := mod.ExportedFunction("_start"); startFn != nil {
		if _, cerr := startFn.Call(ctx); cerr != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return fmt.Errorf("call _start cancelled: %w", ctxErr)
			}
			if exitErr, ok := cerr.(*sys.ExitError); ok {
				if exitErr.ExitCode() != 0 {
					return fmt.Errorf("module _start exited with code %d: %w", exitErr.ExitCode(), cerr)
				}
				// code == 0：正常退出，继续执行入口函数
			} else {
				return fmt.Errorf("failed to call _start: %w", cerr)
			}
		}
	}

	host, err := witgo.NewHost(mod)
	if err != nil {
		return fmt.Errorf("new witgo host: %w", err)
	}

	// flatIn 决定哪些参数传给入口函数。
	// RootAsFirstArg=false（默认）：全部传；in[0] 同时作为 root。
	// RootAsFirstArg=true：in[0] 只作 root，in[1:] 传给入口。
	var flatIn []any
	if T.RootAsFirstArg && len(in) > 0 {
		flatIn = in[1:]
	} else {
		flatIn = in
	}

	// 默认模式下对 flatIn 全体做保守预检：复杂对象（struct/array/map）提前返回明确错误，
	// 让用户改用 RootAsFirstArg=true + tf_root_handle，而不是在 witgo 内部看到难读报错。
	// 注意：isLikelyFlatABIValue 是“保守判定”，不是 witgo.flattenParam 的等价物。
	if !T.RootAsFirstArg {
		for i, v := range flatIn {
			if v == nil {
				continue
			}
			if !isLikelyFlatABIValue(v) {
				return fmt.Errorf("vweb_dynamic: flatIn[%d] type %T is not Wasm-flat-friendly; "+
					"set Wazero.RootAsFirstArg=true and use tf_root_handle in guest, "+
					"or pass only scalar/string/scalar-slice args", i, v)
			}
		}
	}

	entryName = entryname(entryName, T.pagePath)
	var body []byte
	if cerr := host.Call(ctx, entryName, &body, flatIn...); cerr != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return fmt.Errorf("call entry %q cancelled: %w", entryName, ctxErr)
		}
		return fmt.Errorf("call entry %q: %w", entryName, cerr)
	}

	if len(body) > 0 {
		if _, werr := out.Write(body); werr != nil {
			return fmt.Errorf("write body: %w", werr)
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
// 顺序：先关 Runtime（停止所有 guest 调用），再清句柄表，最后关 CompilationCache。
// 先关 Runtime 是关键：否则 clear() 之后仍可能有 in-flight guest 访问失效句柄。
func (T *Wazero) closeLocked() error {
	var firstErr error
	if T.runtime != nil {
		if err := T.runtime.Close(context.Background()); err != nil && firstErr == nil {
			firstErr = err
		}
		T.runtime = nil
		T.cm = nil
	}
	// Runtime 关闭后 Guest 已停止；再关闭 host 释放 manager 资源
	// （filesystem 描述符、socket、http 空闲连接、tls.Conn、stream 后台 goroutine 等）。
	if T.host != nil {
		T.host.Close()
		T.host = nil
	}
	if T.envTable != nil {
		T.envTable.clear()
		T.envTable = nil
	}
	if T.cache != nil {
		if err := T.cache.Close(context.Background()); err != nil && firstErr == nil {
			firstErr = err
		}
		T.cache = nil
	}
	return firstErr
}

func registerPreopens(h *wasip2.Host, dirs map[string]string) error {
	for guestPath, hostPath := range dirs {
		f, err := os.Open(hostPath)
		if err != nil {
			return err
		}
		st, err := f.Stat()
		if err != nil {
			f.Close()
			return err
		}
		if !st.IsDir() {
			f.Close()
			return fmt.Errorf("preopen %q is not a directory", hostPath)
		}
		// AddPreopen 返回 0 表示登记失败（host 或 filesystemManager 为 nil）。
		// 此时 f 未被 host 接管，必须关闭，否则 fd 泄漏。
		if handle := h.AddPreopen(f, guestPath); handle == 0 {
			f.Close()
			return fmt.Errorf("add preopen %q failed", guestPath)
		}
	}
	return nil
}
