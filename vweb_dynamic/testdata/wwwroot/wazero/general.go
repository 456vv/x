// render.go — vweb_dynamic WASI env 桥使用示例（TinyGo）
//
// 编译：
//   tinygo build -o render.wasm -target=wasip2 general.go
//
// 宿主侧：
//   w := &vweb_dynamic.Wazero{RootAsFirstArg: true}
//   w.ParseFile("render.wasm")
//   w.Execute("Main", out, pointObject)
//
// 入口 Main 返回 (ptr, len)，匹配 host.Call(&body []byte) 的 unflatten 路径。
// 宿主程序的完整示例见文件末尾第 9 节。

package main

import (
	"encoding/json"
	"errors"
	"strconv"
	"unsafe"
)

// # 编译
// tinygo build -o testdata/wwwroot/wazero/render.wasm -target=wasip2 testdata/wwwroot/wazero/general.go

// ===========================================================================
// 1. host 导出符号（模块名固定为 "env"）
// ===========================================================================

//go:wasmimport env tf_protocol_version
func tfProtocolVersion() int32

//go:wasmimport env tf_lookup
func tfLookup(namePtr, nameLen uint32) uint64

//go:wasmimport env tf_call
func tfCall(fnHandle uint64, argsPtr, argsLen uint32) uint64

//go:wasmimport env tf_drop
func tfDrop(handle uint64)

//go:wasmimport env tf_stringify
func tfStringify(handle uint64, outPtr, outCap uint32) uint32

//go:wasmimport env tf_last_error
func tfLastError(outPtr, outCap uint32) uint32

//go:wasmimport env tf_root_handle
func tfRootHandle() uint64

//go:wasmimport env tf_root_get
func tfRootGet(rootHandle uint64, keyPtr, keyLen uint32) uint64

//go:wasmimport env tf_root_set
func tfRootSet(rootHandle uint64, keyPtr, keyLen, valPtr, valLen uint32) int32

//go:wasmimport env tf_root_has
func tfRootHas(rootHandle uint64, keyPtr, keyLen uint32) int32

//go:wasmimport env tf_root_len
func tfRootLen(rootHandle uint64, keyPtr, keyLen uint32) int32

//go:wasmimport env tf_stringify_json
func tfStringifyJSON(handle uint64, outPtr, outCap uint32) uint32

//go:wasmimport env tf_root_get_path
func tfRootGetPath(rootHandle uint64, pathPtr, pathLen uint32) uint64

//go:wasmimport env tf_root_set_path
func tfRootSetPath(rootHandle uint64, pathPtr, pathLen, valPtr, valLen uint32) int32

//go:wasmimport env tf_root_call
func tfRootCall(rootHandle uint64, methodPtr, methodLen, argsPtr, argsLen uint32) uint64

// ===========================================================================
// 2. 基础 helper
// ===========================================================================

func ptrOf(b []byte) uint32 {
	if len(b) == 0 {
		return 0
	}
	return uint32(uintptr(unsafe.Pointer(&b[0])))
}

func lastError() error {
	buf := make([]byte, 256)
	n := tfLastError(ptrOf(buf), uint32(len(buf)))
	switch {
	case n == 0:
		return nil
	case n == 0xFFFFFFFF:
		return errors.New("tf_last_error: host write failed")
	case int(n) > len(buf):
		buf = make([]byte, n)
		n = tfLastError(ptrOf(buf), uint32(len(buf)))
	}
	return errors.New(string(buf[:n]))
}

func lookup(name string) (uint64, error) {
	b := []byte(name)
	h := tfLookup(ptrOf(b), uint32(len(b)))
	if h == 0 {
		return 0, lastError()
	}
	return h, nil
}

func stringifyHandle(h uint64) (string, error) {
	n := tfStringify(h, 0, 0)
	if n == 0 {
		if err := lastError(); err != nil {
			return "", err
		}
		return "", nil
	}
	buf := make([]byte, n)
	ret := tfStringify(h, ptrOf(buf), n)
	if ret > n {
		buf = make([]byte, ret)
		ret = tfStringify(h, ptrOf(buf), ret)
	}
	if ret > uint32(len(buf)) {
		return "", errors.New("tf_stringify: need more space")
	}
	return string(buf[:ret]), nil
}

func stringifyJSON(h uint64) ([]byte, error) {
	n := tfStringifyJSON(h, 0, 0)
	if n == 0 {
		return nil, lastError()
	}
	buf := make([]byte, n)
	ret := tfStringifyJSON(h, ptrOf(buf), n)
	if ret > n {
		buf = make([]byte, ret)
		ret = tfStringifyJSON(h, ptrOf(buf), ret)
	}
	if ret == 0 || ret > uint32(len(buf)) {
		return nil, lastError()
	}
	return buf[:ret], nil
}

func callRaw(fn uint64, args ...any) (string, error) {
	raw, err := json.Marshal(args)
	if err != nil {
		return "", err
	}
	res := tfCall(fn, ptrOf(raw), uint32(len(raw)))
	if res == 0 {
		return "", lastError()
	}
	defer tfDrop(res)
	return stringifyHandle(res)
}

func callByName(name string, args ...any) (string, error) {
	fn, err := lookup(name)
	if err != nil {
		return "", err
	}
	defer tfDrop(fn)
	return callRaw(fn, args...)
}

func callJSON(name string, args ...any) ([]byte, error) {
	fn, err := lookup(name)
	if err != nil {
		return nil, err
	}
	defer tfDrop(fn)
	raw, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}
	res := tfCall(fn, ptrOf(raw), uint32(len(raw)))
	if res == 0 {
		return nil, lastError()
	}
	defer tfDrop(res)
	return stringifyJSON(res)
}

func callInt(name string, args ...any) (int, error) {
	s, err := callByName(name, args...)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(s)
}

func callBool(name string, args ...any) (bool, error) {
	s, err := callByName(name, args...)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(s)
}

func rootCall(root uint64, method string, args ...any) ([]byte, error) {
	mb := []byte(method)
	ab, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}
	h := tfRootCall(root, ptrOf(mb), uint32(len(mb)), ptrOf(ab), uint32(len(ab)))
	if h == 0 {
		return nil, lastError()
	}
	defer tfDrop(h)
	return stringifyJSON(h)
}

func rootGet(root uint64, key string) (string, error) {
	kb := []byte(key)
	h := tfRootGet(root, ptrOf(kb), uint32(len(kb)))
	if h == 0 {
		return "", lastError()
	}
	defer tfDrop(h)
	return stringifyHandle(h)
}

func rootSet(root uint64, key string, val any) error {
	kb := []byte(key)
	vb, err := json.Marshal(val)
	if err != nil {
		return err
	}
	if ret := tfRootSet(root, ptrOf(kb), uint32(len(kb)), ptrOf(vb), uint32(len(vb))); ret == -1 {
		return lastError()
	}
	return nil
}

func rootHas(root uint64, key string) bool {
	kb := []byte(key)
	return tfRootHas(root, ptrOf(kb), uint32(len(kb))) == 1
}

func rootLen(root uint64, key string) (int, error) {
	kb := []byte(key)
	n := tfRootLen(root, ptrOf(kb), uint32(len(kb)))
	if n == -1 {
		return 0, lastError()
	}
	return int(n), nil
}

func rootGetPathJSON(root uint64, path string) ([]byte, error) {
	pb := []byte(path)
	h := tfRootGetPath(root, ptrOf(pb), uint32(len(pb)))
	if h == 0 {
		return nil, lastError()
	}
	defer tfDrop(h)
	return stringifyJSON(h)
}

func rootSetPath(root uint64, path string, val any) error {
	pb := []byte(path)
	vb, err := json.Marshal(val)
	if err != nil {
		return err
	}
	if ret := tfRootSetPath(root, ptrOf(pb), uint32(len(pb)), ptrOf(vb), uint32(len(vb))); ret == -1 {
		return lastError()
	}
	return nil
}

// ===========================================================================
// 3. 示例用类型（字段 JSON tag 与宿主侧对齐）
// ===========================================================================

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Order struct {
	ID    int     `json:"id"`
	Items []Point `json:"items"`
	Total float64 `json:"total"`
}

type S struct {
	A int `json:"a"`
}

type Config struct {
	Name   string         `json:"name"`
	Tags   []string       `json:"tags"`
	Nested map[string]any `json:"nested"`
}

// ===========================================================================
// 4. 示例函数（保留作参考；不影响 Main 入口）
// ===========================================================================

func demoInt() (string, error)        { return callByName("Compute", 1, "+", 2) }
func demoFloatCompare() (bool, error) { return callBool("LT", 1.5, 2.25) }
func demoLenSlice() (int, error)      { return callInt("Len", []int{10, 20, 30}) }
func demoBigUint() (string, error)    { return callByName("Compute", "9007199254740993", "+", "0") }
func demoOr() (bool, error)           { return callBool("Or", true, false) }
func demoNot() (bool, error)          { return callBool("Not", true) }
func demoAnd() (bool, error)          { return callBool("And", 1, true) }
func demoOrMany() (bool, error)       { return callBool("Or", nil, 0, "", false, 1) }

func demoSliceOfStruct() (int, error) {
	pts := []Point{{X: 1, Y: 2}, {X: 3, Y: 4}}
	return callInt("SumPoints", pts)
}

func demoMapArg() (string, error) {
	return callByName("RenderUser", map[string]any{
		"name": "Alice", "age": 30, "admin": true,
	})
}

func demoNested() (string, error) {
	o := Order{ID: 1001, Items: []Point{{1, 2}, {3, 4}}, Total: 99.5}
	return callByName("FormatOrder", o)
}

func demoPointerArg() (string, error) {
	p := Point{X: 7, Y: 8}
	return callByName("PointToString", p)
}

func demoNil() (string, error) { return callByName("IsNil", nil) }

func demoOptional(name string, age *int) (string, error) {
	var ageArg any
	if age != nil {
		ageArg = *age
	}
	return callByName("RenderProfile", name, ageArg)
}

func demoMixed() (string, error) {
	return callByName("BuildQuery", "users", map[string]any{"status": "active", "score": 80}, 100, true)
}

func demoComputeMixed() (string, error) { return callByName("Compute", 10, "*", 3) }
func demoStructReturn() (string, error) { return callByName("SomeFuncReturningStruct") }

func demoMultiReturn() {
	s, _ := callByName("Divide", 10, 2)
	_ = s
}

func batchLen(items [][]int) ([]int, error) {
	fn, err := lookup("Len")
	if err != nil {
		return nil, err
	}
	defer tfDrop(fn)
	out := make([]int, len(items))
	for i, it := range items {
		s, err := callRaw(fn, it)
		if err != nil {
			return nil, err
		}
		out[i], _ = strconv.Atoi(s)
	}
	return out, nil
}

func demoGetMap() (int, error)   { return callInt("Get", map[string]any{"a": 1, "b": 2}, "a") }
func demoGetSlice() (int, error) { return callInt("Get", []int{10, 20, 30}, 1) }
func demoGetStruct() (int, error) {
	return callInt("Get", S{A: 7}, "A")
}

func demoSetMap() (string, error)   { return callByName("Set", map[string]any{}, "k", 42) }
func demoSetSlice() (string, error) { return callByName("Set", []int{0, 0, 0}, 1, 99) }
func demoStringUpper() (string, error) {
	return callByName("StringUpper", "hello")
}

func demoAppend() (string, error) {
	return callByName("Append", []int{1, 2}, 3, 4, 5)
}

func demoGetConfig() (Config, error) {
	var cfg Config
	b, err := callJSON("GetConfig")
	if err != nil {
		return cfg, err
	}
	if err := json.Unmarshal(b, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// ===========================================================================
// 5. 页面入口（唯一的 output / main / Main）
// ===========================================================================

// output 必须是包级变量，避免 render 返回后局部切片被 GC 回收，
// 导致 host 在 fn.Call 返回后读到脏数据。
var output []byte

//go:wasmexport Main
func Main() (uint32, uint32) {
	output = []byte(render())
	return ptrOf(output), uint32(len(output))
}

//go:wasmexport MainHTML
func MainHTML() (uint32, uint32) {
	output = []byte(renderHTML())
	return ptrOf(output), uint32(len(output))
}

// main 是 TinyGo wasip2 target 要求的空函数（不被 host 直接调用）。
func main() {}

// ===========================================================================
// 6. 页面渲染逻辑
// ===========================================================================

// render 演示三种典型路径：root 直接 stringify / 调 TemplateFunc / 兜底。
func render() string {
	if v := tfProtocolVersion(); v != 1 {
		return "<h1>unsupported protocol</h1>"
	}
	if root := tfRootHandle(); root != 0 {
		if s, err := stringifyHandle(root); err == nil {
			return s
		}
	}
	if s, err := callStringByName("Get", "Title"); err == nil {
		return "Title=" + s
	}
	return "hello from guest"
}

// renderHTML 演示非 root 的直接模板渲染路径。
func renderHTML() string {
	if tfProtocolVersion() != 1 {
		return "<h1>unsupported protocol</h1>"
	}
	title, err := callStringByName("Get", "Title")
	if err != nil {
		return "<h1>error: " + err.Error() + "</h1>"
	}
	escTitle, _ := callStringByName("Escape", title)
	count, _ := callStringByName("Len", "Items")
	return "<h1>" + escTitle + "</h1><p>items: " + count + "</p>"
}

// callStringByName 语义与 callByName 相同，便于从 render/renderHTML 内调用。
func callStringByName(name string, args ...any) (string, error) {
	return callByName(name, args...)
}

// renderWithRoot 演示 root 句柄（RootAsFirstArg=true 时宿主会把 in[0] 绑成 root）。
func renderWithRoot() string {
	root := tfRootHandle()
	if root == 0 {
		return "no root"
	}
	title, _ := callByName("RootGet", root, "Title")
	return title
}

// renderSet 演示"读 root → 改 root → 再读回"。
func renderSet() string {
	if tfProtocolVersion() != 1 {
		return "<h1>unsupported protocol</h1>"
	}
	root := tfRootHandle()
	if root == 0 {
		return "<h1>no root</h1>"
	}
	title, _ := rootGet(root, "Title")
	n, _ := rootLen(root, "Items")
	_ = rootSet(root, "Title", "<"+title+">")
	title2, _ := rootGet(root, "Title")
	hasItems := rootHas(root, "Items")
	return title2 + " items=" + strconv.Itoa(n) +
		" hasItems=" + strconv.FormatBool(hasItems)
}

// renderPage 演示按 "." 路径修改嵌套字段（root 必须是 *T 才能改 struct 值字段）。
func renderPage() string {
	root := tfRootHandle()
	_ = rootSetPath(root, "User.Name", "Bob")
	_ = rootSetPath(root, "User.Age", 31)
	_ = rootSetPath(root, "Meta.lang", "zh-CN")
	_ = rootSetPath(root, "Tags.0", "first")

	nameB, _ := rootGetPathJSON(root, "User.Name")
	ageB, _ := rootGetPathJSON(root, "User.Age")
	var name string
	var age int
	_ = json.Unmarshal(nameB, &name)
	_ = json.Unmarshal(ageB, &age)
	return name + "/" + strconv.Itoa(age)
}

// renderStruct 演示 struct 方法调用 + 嵌套字段修改。
func renderStruct() string {
	root := tfRootHandle()
	_ = rootSetPath(root, "Title", "new")
	_ = rootSetPath(root, "User.Name", "Bob")
	_ = rootSetPath(root, "User.Age", 42)
	_ = rootSetPath(root, "Meta.lang", "zh-CN")
	_, _ = rootCall(root, "SetTitle", "from-method")
	b, _ := rootGetPathJSON(root, "User")
	return string(b)
}

/*
===============================================================================
9. 宿主侧代码示例（Go 程序）
===============================================================================

以下内容不属于 guest 模块，只是展示如何从宿主进程（普通 Go 程序）调用
上面编译出的 render.wasm，以及如何注册 guest 用到的 TemplateFunc 与 root
对象。TinyGo 编译时整块注释被忽略。

编译 guest：
    tinygo build -o render.wasm -target=wasip2 general.go

下面的示例代码可以直接存为 host.go 运行（需 `go mod tidy` 拉取
github.com/456vv/x/vweb_dynamic）。

===============================================================================
9.1 最小骨架 —— 无 root，直接调 TemplateFunc（对应 MainHTML）
===============================================================================

guest 侧 MainHTML 入口（renderHTML）只调三个 TemplateFunc：

    Get("Title")        → 取标题
    Escape(title)       → HTML 转义
    Len("Items")        → 取长度

宿主程序：

    // host.go
    package main

    import (
        "bytes"
        "fmt"
        "html"
        "os"
        "reflect"

        "github.com/456vv/x/vweb_dynamic"
    )

    func main() {
        // 1. 注册 guest 用到的 TemplateFunc（必须在 Parse/Execute 前完成）。
        //    TemplateFunc 是 template.FuncMap，全局可写；写入需在 Execute 前串行完成。
        state := struct {
            Title string
            Items []int
        }{Title: "<Hello>", Items: []int{1, 2, 3}}

        // 覆盖 builtin.Get 仅为演示：只接 key，从 state 取值。
        // 生产环境请保留 builtin.Get，或在 guest 侧用 tf_root_get 直接取。
        vweb_dynamic.TemplateFunc["Get"] = func(key string) any {
            switch key {
            case "Title":
                return state.Title
            case "Items":
                return state.Items
            }
            return nil
        }

        // 覆盖 builtin.Len 仅为演示：字符串当 key 从 state 取长度，
        // 其他值按容器/字符串长度返回。
        vweb_dynamic.TemplateFunc["Len"] = func(v any) int {
            if key, ok := v.(string); ok {
                switch key {
                case "Items":
                    return len(state.Items)
                }
            }
            rv := reflect.ValueOf(v)
            switch rv.Kind() {
            case reflect.Map, reflect.Slice, reflect.Array, reflect.String:
                return rv.Len()
            }
            return 0
        }

        // Escape 是 guest 侧 renderHTML 里 callStringByName("Escape", title) 调用的。
        vweb_dynamic.TemplateFunc["Escape"] = html.EscapeString

        // 2. 创建 Wazero 并加载 wasm。
        w := &vweb_dynamic.Wazero{
            // true：in[0] 仅作 root；本例不传 in，二者等价。
            // 若改传 root，MainHTML 也能跑（它不读 tf_root_handle）。
            RootAsFirstArg: true,
        }
        w.SetPath(".", "render.wasm")
        if err := w.ParseFile("render.wasm"); err != nil {
            fmt.Fprintf(os.Stderr, "parse: %v\n", err)
            os.Exit(1)
        }
        defer w.Close()

        // 3. 执行 MainHTML（无 root）。
        var out bytes.Buffer
        if err := w.Execute("MainHTML", &out); err != nil {
            fmt.Fprintf(os.Stderr, "execute MainHTML: %v\n", err)
            os.Exit(1)
        }
        fmt.Println(out.String())
        // 输出：<h1>&lt;Hello&gt;</h1><p>items: 3</p>
    }

===============================================================================
9.2 带 root 对象 —— RootAsFirstArg=true（对应 Main / renderPage / renderStruct）
===============================================================================

guest 侧 Main 入口（render）会先看 tf_root_handle()：

    root != 0   → 直接 stringifyHandle(root) 返回（本例不演示，见下方说明）
    root == 0   → 调 Get("Title")

renderPage / renderStruct 则直接用 tf_root_* 按 "." 路径改字段：

    tfRootSetPath(root, "User.Name", "Bob")
    tfRootSetPath(root, "User.Age",  31)
    tfRootSetPath(root, "Meta.lang", "zh-CN")
    tfRootSetPath(root, "Tags.0",    "first")
    tfRootCall(root, "SetTitle", "from-method")

root 对象必须是 *Page（指针）：
  - renderPage 改 struct 值字段（Page.User.Name）需要字段可寻址、可 Set；
  - renderStruct 调 (*Page).SetTitle 需要指针接收者；
  - 传 Page 值 → tf_root_set_path 在 struct 值上会返回 "not settable"。

宿主程序：

    // host.go
    package main

    import (
        "bytes"
        "fmt"
        "html"
        "os"
        "reflect"

        "github.com/456vv/x/vweb_dynamic"
    )

    // Page 是传给 guest 的 root 对象。
    // 字段 tag 与 guest 侧 renderPage / renderStruct 用的 "." 路径一致：
    //   "User.Name" / "User.Age" / "Meta.lang" / "Tags.0"
    type Page struct {
        Title string         `json:"title"`
        Items []int          `json:"items"`
        User  User           `json:"user"`
        Meta  map[string]any `json:"meta"`
        Tags  []string       `json:"tags"`
    }

    type User struct {
        Name string `json:"name"`
        Age  int    `json:"age"`
    }

    // SetTitle 供 guest 的 tfRootCall(root, "SetTitle", "...") 调用。
    // 必须是 *Page 的方法（指针接收者），值接收者改不到 p。
    func (p *Page) SetTitle(s string) {
        p.Title = s
    }

    func main() {
        // 1. 注册 guest 用到的 TemplateFunc。
        //    render / renderWithRoot 可能走 Get / RootGet / RootHas / RootLen / EscapeHTML。
        //    这些是 demo 用包装函数，签名与 guest 侧 callByName(...) 的参数对齐。
        registerTemplateFuncs()

        // 2. 创建 Wazero。
        w := &vweb_dynamic.Wazero{RootAsFirstArg: true}
        w.SetPath(".", "render.wasm")
        if err := w.ParseFile("render.wasm"); err != nil {
            fmt.Fprintf(os.Stderr, "parse: %v\n", err)
            os.Exit(1)
        }
        defer w.Close()

        // 3. 构造 root 对象（必须是指针）。
        root := &Page{
            Title: "<Hello>",
            Items: []int{1, 2, 3},
            User:  User{Name: "Alice", Age: 30},
            Meta:  map[string]any{"lang": "en"},
            Tags:  []string{"a", "b"},
        }

        // 4. 执行 Main（root 非 0 时 render 会直接 stringify root 返回）。
        var out bytes.Buffer
        if err := w.Execute("Main", &out, root); err != nil {
            fmt.Fprintf(os.Stderr, "execute Main: %v\n", err)
            os.Exit(1)
        }
        fmt.Println("Main:", out.String())

        // 5. 检查 root 是否被 guest 修改（如果走了 renderPage / renderStruct）。
        //    这里 Main 走的是 stringifyHandle(root)，不会改字段；
        //    想验证修改，把 guest 侧 Main 改成调 renderPage/renderStruct。
        fmt.Printf("after Main: Title=%q User=%+v Meta=%v Tags=%v\n",
            root.Title, root.User, root.Meta, root.Tags)
    }

    func registerTemplateFuncs() {
        // renderWithRoot / render 用到的包装函数。
        // 注意签名：guest 侧 callByName("RootGet", root, "Title") 会把
        // root(uint64) 和 key(string) 一起 JSON 化后再调用，所以这里直接
        // 声明 (uint64, string) 形参即可。

        vweb_dynamic.TemplateFunc["RootGet"] = func(root uint64, key string) string {
            // 真实实现：从 root 对象取 key。这里仅示意签名，
            // 因为 root 句柄只在 guest 的 env 桥内部可达，宿主侧拿不到对象；
            // 若要在宿主侧实现，需另建映射或改用 tf_root_get。
            return ""
        }
        vweb_dynamic.TemplateFunc["RootHas"] = func(root uint64, key string) bool { return false }
        vweb_dynamic.TemplateFunc["RootLen"] = func(root uint64, key string) int { return 0 }
        vweb_dynamic.TemplateFunc["EscapeHTML"] = html.EscapeString

        // render 兜底分支用到的 Get（仅在没有 root 时走到）。
        // 见 9.1 的覆盖说明；这里保持与 9.1 一致。
        vweb_dynamic.TemplateFunc["Get"] = func(key string) any {
            return nil
        }
        vweb_dynamic.TemplateFunc["Len"] = func(v any) int {
            rv := reflect.ValueOf(v)
            switch rv.Kind() {
            case reflect.Map, reflect.Slice, reflect.Array, reflect.String:
                return rv.Len()
            }
            return 0
        }
    }

说明：
  - 上面 RootGet / RootHas / RootLen 的 "root uint64" 形参只是 demo 签名。
    guest 侧 tf_call 的 JSON 参数会把 root 句柄当作 uint64 序列化，宿主侧
    函数签名 (uint64, string) 能收到这个数，但宿主侧无法把它反解回 Page 对象
    —— 句柄表在 wazero_env.go 内部。
  - 真正需要"guest 按 key 取 root 字段"时，请让 guest 直接用 tf_root_get，
    而不是绕道 TemplateFunc。renderWithRoot 的示例只演示调用链路。
  - 若要验证 renderPage / renderStruct 的效果，把 guest 的 Main 改为：

        //go:wasmexport Main
        func Main() (uint32, uint32) {
            output = []byte(renderStruct())   // 或 renderPage()
            return ptrOf(output), uint32(len(output))
        }

    重新编译后执行同一宿主程序，最后打印的 root 字段就会是 guest 改过的值：
        Title="from-method" User={Name:Bob Age:42} Meta=map[lang:zh-CN] Tags=[first b]

===============================================================================
9.3 可选：demo 用 TemplateFunc 清单
===============================================================================

第 4 节的 demoXxx 函数演示了各类参数/返回值形态。要让它们可运行，宿主侧
需注册对应 TemplateFunc（生产环境大多应复用 builtin 包，不要手写）：

    vweb_dynamic.TemplateFunc["Compute"] = builtin.Compute                  // (any, string, any) any
    vweb_dynamic.TemplateFunc["LT"]      = builtin.LT                       // (any, any) bool
    vweb_dynamic.TemplateFunc["Or"]      = builtin.Or                       // (...any) bool
    vweb_dynamic.TemplateFunc["Not"]     = builtin.Not                      // (any) bool
    vweb_dynamic.TemplateFunc["And"]     = builtin.And                      // (...any) bool
    vweb_dynamic.TemplateFunc["Get"]     = builtin.Get                      // (any, any) any
    vweb_dynamic.TemplateFunc["Set"]     = builtin.Set                      // (any, any, any) any
    vweb_dynamic.TemplateFunc["Append"]  = builtin.Append                   // (any, ...any) any
    vweb_dynamic.TemplateFunc["Len"]     = builtin.Len                      // (any) int

    // 业务函数（宿主自己定义）
    vweb_dynamic.TemplateFunc["SumPoints"]  = func(pts []Point) int { ... }
    vweb_dynamic.TemplateFunc["RenderUser"] = func(m map[string]any) string { ... }
    vweb_dynamic.TemplateFunc["FormatOrder"] = func(o Order) string { ... }
    vweb_dynamic.TemplateFunc["PointToString"] = func(p Point) string { ... }
    vweb_dynamic.TemplateFunc["IsNil"]      = func(v any) bool { return v == nil }
    vweb_dynamic.TemplateFunc["RenderProfile"] = func(name string, age *int) string { ... }
    vweb_dynamic.TemplateFunc["BuildQuery"] = func(table string, where map[string]any, limit int, asc bool) string { ... }
    vweb_dynamic.TemplateFunc["SomeFuncReturningStruct"] = func() S { return S{A: 1} }
    vweb_dynamic.TemplateFunc["Divide"]     = func(a, b int) (int, error) {
        if b == 0 { return 0, errors.New("divide by zero") }
        return a / b, nil
    }
    vweb_dynamic.TemplateFunc["StringUpper"] = strings.ToUpper
    vweb_dynamic.TemplateFunc["GetConfig"]   = func() Config { ... }

注册时机：
  - TemplateFunc 是全局 template.FuncMap，写入本身不并发安全；
  - 但 tf_lookup 只在 Execute 时读，故只需在首次 Execute 前完成注册；
  - 多 goroutine 并发 Execute 时，注册动作应集中在启动阶段完成。

===============================================================================
9.4 注意事项
===============================================================================

1. root 必须是指针
   - 传 Page（值）时，wazero_env.go 的 setByPath 会在 struct 值上返回
     "not settable（root 不可寻址？改用 *T 作 root）"；
   - 传 *Page 时，setByPath 能走到 FieldByName().Set()。

2. RootAsFirstArg 的选择
   - false（默认）：in[0] 既作 root 又作为入口函数的第一参数。
     只适合 in[0] 是标量/string/标量 slice 的场景；
     struct/map 会被 isLikelyFlatABIValue 提前拒绝，并在 witgo 内部报错。
   - true：in[0] 仅作 root，不传给入口函数。
     guest 用 tf_root_handle() 取，适合模板点对象。

3. 句柄生命周期
   - tf_call / tf_root_get* 创建的结果句柄由 ExecuteContext 在入口返回后
     通过 releaseModule 兜底回收；guest 显式 tf_drop 只是提前释放。
   - root 句柄由 ExecuteContext 在入口返回后清理；root 对象本身仍归宿主所有。

4. 编译缓存
   - w.SetPath(root, page) 的 root 非空时，Wazero.parse 会在
     <root>/.cache 下建磁盘编译缓存；切换 root 会切换缓存目录。
*/
