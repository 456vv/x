// Code generated by 'yaegi extract github.com/456vv/vweb/v2'. DO NOT EDIT.

//go:build yaegi_lib
// +build yaegi_lib

package yaegi_lib

import (
	"bufio"
	"context"
	"github.com/456vv/vmap/v2"
	"github.com/456vv/vweb/v2"
	"io"
	"net"
	"net/http"
	"reflect"
	"time"
)

func init() {
	Symbols["github.com/456vv/vweb/v2/vweb"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"AddSalt":               reflect.ValueOf(vweb.AddSalt),
		"CopyStruct":            reflect.ValueOf(vweb.CopyStruct),
		"CopyStructDeep":        reflect.ValueOf(vweb.CopyStructDeep),
		"DefaultSitePool":       reflect.ValueOf(&vweb.DefaultSitePool).Elem(),
		"DepthField":            reflect.ValueOf(vweb.DepthField),
		"ExecFunc":              reflect.ValueOf(vweb.ExecFunc),
		"ExtendTemplatePackage": reflect.ValueOf(vweb.ExtendTemplatePackage),
		"ForMethod":             reflect.ValueOf(vweb.ForMethod),
		"ForType":               reflect.ValueOf(vweb.ForType),
		"GenerateRandom":        reflect.ValueOf(vweb.GenerateRandom),
		"GenerateRandomId":      reflect.ValueOf(vweb.GenerateRandomId),
		"GenerateRandomString":  reflect.ValueOf(vweb.GenerateRandomString),
		"InDirect":              reflect.ValueOf(vweb.InDirect),
		"NewSitePool":           reflect.ValueOf(vweb.NewSitePool),
		"PagePath":              reflect.ValueOf(vweb.PagePath),
		"PluginTypeHTTP":        reflect.ValueOf(vweb.PluginTypeHTTP),
		"PluginTypeRPC":         reflect.ValueOf(vweb.PluginTypeRPC),
		"Version":               reflect.ValueOf(vweb.Version),

		// type definitions
		"Cookie":               reflect.ValueOf((*vweb.Cookie)(nil)),
		"Cookier":              reflect.ValueOf((*vweb.Cookier)(nil)),
		"DotContexter":         reflect.ValueOf((*vweb.DotContexter)(nil)),
		"DynamicTemplateFunc":  reflect.ValueOf((*vweb.DynamicTemplateFunc)(nil)),
		"DynamicTemplater":     reflect.ValueOf((*vweb.DynamicTemplater)(nil)),
		"ExecCall":             reflect.ValueOf((*vweb.ExecCall)(nil)),
		"ExitCall":             reflect.ValueOf((*vweb.ExitCall)(nil)),
		"Forward":              reflect.ValueOf((*vweb.Forward)(nil)),
		"Globaler":             reflect.ValueOf((*vweb.Globaler)(nil)),
		"PluginHTTP":           reflect.ValueOf((*vweb.PluginHTTP)(nil)),
		"PluginHTTPClient":     reflect.ValueOf((*vweb.PluginHTTPClient)(nil)),
		"PluginRPC":            reflect.ValueOf((*vweb.PluginRPC)(nil)),
		"PluginRPCClient":      reflect.ValueOf((*vweb.PluginRPCClient)(nil)),
		"PluginType":           reflect.ValueOf((*vweb.PluginType)(nil)),
		"Responser":            reflect.ValueOf((*vweb.Responser)(nil)),
		"Route":                reflect.ValueOf((*vweb.Route)(nil)),
		"ServerHandlerDynamic": reflect.ValueOf((*vweb.ServerHandlerDynamic)(nil)),
		"ServerHandlerStatic":  reflect.ValueOf((*vweb.ServerHandlerStatic)(nil)),
		"Session":              reflect.ValueOf((*vweb.Session)(nil)),
		"Sessioner":            reflect.ValueOf((*vweb.Sessioner)(nil)),
		"Sessions":             reflect.ValueOf((*vweb.Sessions)(nil)),
		"Site":                 reflect.ValueOf((*vweb.Site)(nil)),
		"SiteMan":              reflect.ValueOf((*vweb.SiteMan)(nil)),
		"SitePool":             reflect.ValueOf((*vweb.SitePool)(nil)),
		"TemplateDot":          reflect.ValueOf((*vweb.TemplateDot)(nil)),
		"TemplateDoter":        reflect.ValueOf((*vweb.TemplateDoter)(nil)),

		// interface wrapper definitions
		"_Cookier":          reflect.ValueOf((*_github_com_456vv_vweb_v2_Cookier)(nil)),
		"_DotContexter":     reflect.ValueOf((*_github_com_456vv_vweb_v2_DotContexter)(nil)),
		"_DynamicTemplater": reflect.ValueOf((*_github_com_456vv_vweb_v2_DynamicTemplater)(nil)),
		"_Globaler":         reflect.ValueOf((*_github_com_456vv_vweb_v2_Globaler)(nil)),
		"_PluginHTTP":       reflect.ValueOf((*_github_com_456vv_vweb_v2_PluginHTTP)(nil)),
		"_PluginRPC":        reflect.ValueOf((*_github_com_456vv_vweb_v2_PluginRPC)(nil)),
		"_Responser":        reflect.ValueOf((*_github_com_456vv_vweb_v2_Responser)(nil)),
		"_Sessioner":        reflect.ValueOf((*_github_com_456vv_vweb_v2_Sessioner)(nil)),
		"_TemplateDoter":    reflect.ValueOf((*_github_com_456vv_vweb_v2_TemplateDoter)(nil)),
	}
}

// _github_com_456vv_vweb_v2_Cookier is an interface wrapper for Cookier type
type _github_com_456vv_vweb_v2_Cookier struct {
	IValue     interface{}
	WAdd       func(name string, value string, path string, domain string, maxAge int, secure bool, only bool, sameSite http.SameSite)
	WDel       func(name string)
	WGet       func(name string) string
	WReadAll   func() map[string]string
	WRemoveAll func()
}

func (W _github_com_456vv_vweb_v2_Cookier) Add(name string, value string, path string, domain string, maxAge int, secure bool, only bool, sameSite http.SameSite) {
	W.WAdd(name, value, path, domain, maxAge, secure, only, sameSite)
}
func (W _github_com_456vv_vweb_v2_Cookier) Del(name string) {
	W.WDel(name)
}
func (W _github_com_456vv_vweb_v2_Cookier) Get(name string) string {
	return W.WGet(name)
}
func (W _github_com_456vv_vweb_v2_Cookier) ReadAll() map[string]string {
	return W.WReadAll()
}
func (W _github_com_456vv_vweb_v2_Cookier) RemoveAll() {
	W.WRemoveAll()
}

// _github_com_456vv_vweb_v2_DotContexter is an interface wrapper for DotContexter type
type _github_com_456vv_vweb_v2_DotContexter struct {
	IValue       interface{}
	WContext     func() context.Context
	WWithContext func(ctx context.Context)
}

func (W _github_com_456vv_vweb_v2_DotContexter) Context() context.Context {
	return W.WContext()
}
func (W _github_com_456vv_vweb_v2_DotContexter) WithContext(ctx context.Context) {
	W.WWithContext(ctx)
}

// _github_com_456vv_vweb_v2_DynamicTemplater is an interface wrapper for DynamicTemplater type
type _github_com_456vv_vweb_v2_DynamicTemplater struct {
	IValue   interface{}
	WExecute func(out io.Writer, dot interface{}) error
	WParse   func(r io.Reader) (err error)
	WSetPath func(rootPath string, pagePath string)
}

func (W _github_com_456vv_vweb_v2_DynamicTemplater) Execute(out io.Writer, dot interface{}) error {
	return W.WExecute(out, dot)
}
func (W _github_com_456vv_vweb_v2_DynamicTemplater) Parse(r io.Reader) (err error) {
	return W.WParse(r)
}
func (W _github_com_456vv_vweb_v2_DynamicTemplater) SetPath(rootPath string, pagePath string) {
	W.WSetPath(rootPath, pagePath)
}

// _github_com_456vv_vweb_v2_Globaler is an interface wrapper for Globaler type
type _github_com_456vv_vweb_v2_Globaler struct {
	IValue          interface{}
	WDel            func(key interface{})
	WGet            func(key interface{}) interface{}
	WHas            func(key interface{}) bool
	WReset          func()
	WSet            func(key interface{}, val interface{})
	WSetExpired     func(key interface{}, d time.Duration)
	WSetExpiredCall func(key interface{}, d time.Duration, f func(interface{}))
}

func (W _github_com_456vv_vweb_v2_Globaler) Del(key interface{}) {
	W.WDel(key)
}
func (W _github_com_456vv_vweb_v2_Globaler) Get(key interface{}) interface{} {
	return W.WGet(key)
}
func (W _github_com_456vv_vweb_v2_Globaler) Has(key interface{}) bool {
	return W.WHas(key)
}
func (W _github_com_456vv_vweb_v2_Globaler) Reset() {
	W.WReset()
}
func (W _github_com_456vv_vweb_v2_Globaler) Set(key interface{}, val interface{}) {
	W.WSet(key, val)
}
func (W _github_com_456vv_vweb_v2_Globaler) SetExpired(key interface{}, d time.Duration) {
	W.WSetExpired(key, d)
}
func (W _github_com_456vv_vweb_v2_Globaler) SetExpiredCall(key interface{}, d time.Duration, f func(interface{})) {
	W.WSetExpiredCall(key, d, f)
}

// _github_com_456vv_vweb_v2_PluginHTTP is an interface wrapper for PluginHTTP type
type _github_com_456vv_vweb_v2_PluginHTTP struct {
	IValue                interface{}
	WCancelRequest        func(req *http.Request)
	WCloseIdleConnections func()
	WRegisterProtocol     func(scheme string, rt http.RoundTripper)
	WRoundTrip            func(r *http.Request) (resp *http.Response, err error)
	WServeHTTP            func(w http.ResponseWriter, r *http.Request)
	WType                 func() vweb.PluginType
}

func (W _github_com_456vv_vweb_v2_PluginHTTP) CancelRequest(req *http.Request) {
	W.WCancelRequest(req)
}
func (W _github_com_456vv_vweb_v2_PluginHTTP) CloseIdleConnections() {
	W.WCloseIdleConnections()
}
func (W _github_com_456vv_vweb_v2_PluginHTTP) RegisterProtocol(scheme string, rt http.RoundTripper) {
	W.WRegisterProtocol(scheme, rt)
}
func (W _github_com_456vv_vweb_v2_PluginHTTP) RoundTrip(r *http.Request) (resp *http.Response, err error) {
	return W.WRoundTrip(r)
}
func (W _github_com_456vv_vweb_v2_PluginHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	W.WServeHTTP(w, r)
}
func (W _github_com_456vv_vweb_v2_PluginHTTP) Type() vweb.PluginType {
	return W.WType()
}

// _github_com_456vv_vweb_v2_PluginRPC is an interface wrapper for PluginRPC type
type _github_com_456vv_vweb_v2_PluginRPC struct {
	IValue    interface{}
	WCall     func(name string, arg interface{}) (interface{}, error)
	WClose    func() error
	WDiscard  func() error
	WRegister func(value interface{})
	WType     func() vweb.PluginType
}

func (W _github_com_456vv_vweb_v2_PluginRPC) Call(name string, arg interface{}) (interface{}, error) {
	return W.WCall(name, arg)
}
func (W _github_com_456vv_vweb_v2_PluginRPC) Close() error {
	return W.WClose()
}
func (W _github_com_456vv_vweb_v2_PluginRPC) Discard() error {
	return W.WDiscard()
}
func (W _github_com_456vv_vweb_v2_PluginRPC) Register(value interface{}) {
	W.WRegister(value)
}
func (W _github_com_456vv_vweb_v2_PluginRPC) Type() vweb.PluginType {
	return W.WType()
}

// _github_com_456vv_vweb_v2_Responser is an interface wrapper for Responser type
type _github_com_456vv_vweb_v2_Responser struct {
	IValue       interface{}
	WError       func(a0 string, a1 int)
	WFlush       func()
	WHijack      func() (net.Conn, *bufio.ReadWriter, error)
	WPush        func(target string, opts *http.PushOptions) error
	WReadFrom    func(a0 io.Reader) (int64, error)
	WRedirect    func(a0 string, a1 int)
	WWrite       func(a0 []byte) (int, error)
	WWriteHeader func(a0 int)
	WWriteString func(a0 string) (int, error)
}

func (W _github_com_456vv_vweb_v2_Responser) Error(a0 string, a1 int) {
	W.WError(a0, a1)
}
func (W _github_com_456vv_vweb_v2_Responser) Flush() {
	W.WFlush()
}
func (W _github_com_456vv_vweb_v2_Responser) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return W.WHijack()
}
func (W _github_com_456vv_vweb_v2_Responser) Push(target string, opts *http.PushOptions) error {
	return W.WPush(target, opts)
}
func (W _github_com_456vv_vweb_v2_Responser) ReadFrom(a0 io.Reader) (int64, error) {
	return W.WReadFrom(a0)
}
func (W _github_com_456vv_vweb_v2_Responser) Redirect(a0 string, a1 int) {
	W.WRedirect(a0, a1)
}
func (W _github_com_456vv_vweb_v2_Responser) Write(a0 []byte) (int, error) {
	return W.WWrite(a0)
}
func (W _github_com_456vv_vweb_v2_Responser) WriteHeader(a0 int) {
	W.WWriteHeader(a0)
}
func (W _github_com_456vv_vweb_v2_Responser) WriteString(a0 string) (int, error) {
	return W.WWriteString(a0)
}

// _github_com_456vv_vweb_v2_Sessioner is an interface wrapper for Sessioner type
type _github_com_456vv_vweb_v2_Sessioner struct {
	IValue          interface{}
	WDefer          func(call interface{}, args ...interface{}) error
	WDel            func(key interface{})
	WFree           func()
	WGet            func(key interface{}) interface{}
	WGetHas         func(key interface{}) (val interface{}, ok bool)
	WHas            func(key interface{}) bool
	WReset          func()
	WSet            func(key interface{}, val interface{})
	WSetExpired     func(key interface{}, d time.Duration)
	WSetExpiredCall func(key interface{}, d time.Duration, f func(interface{}))
	WToken          func() string
}

func (W _github_com_456vv_vweb_v2_Sessioner) Defer(call interface{}, args ...interface{}) error {
	return W.WDefer(call, args...)
}
func (W _github_com_456vv_vweb_v2_Sessioner) Del(key interface{}) {
	W.WDel(key)
}
func (W _github_com_456vv_vweb_v2_Sessioner) Free() {
	W.WFree()
}
func (W _github_com_456vv_vweb_v2_Sessioner) Get(key interface{}) interface{} {
	return W.WGet(key)
}
func (W _github_com_456vv_vweb_v2_Sessioner) GetHas(key interface{}) (val interface{}, ok bool) {
	return W.WGetHas(key)
}
func (W _github_com_456vv_vweb_v2_Sessioner) Has(key interface{}) bool {
	return W.WHas(key)
}
func (W _github_com_456vv_vweb_v2_Sessioner) Reset() {
	W.WReset()
}
func (W _github_com_456vv_vweb_v2_Sessioner) Set(key interface{}, val interface{}) {
	W.WSet(key, val)
}
func (W _github_com_456vv_vweb_v2_Sessioner) SetExpired(key interface{}, d time.Duration) {
	W.WSetExpired(key, d)
}
func (W _github_com_456vv_vweb_v2_Sessioner) SetExpiredCall(key interface{}, d time.Duration, f func(interface{})) {
	W.WSetExpiredCall(key, d, f)
}
func (W _github_com_456vv_vweb_v2_Sessioner) Token() string {
	return W.WToken()
}

// _github_com_456vv_vweb_v2_TemplateDoter is an interface wrapper for TemplateDoter type
type _github_com_456vv_vweb_v2_TemplateDoter struct {
	IValue            interface{}
	WContext          func() context.Context
	WCookie           func() vweb.Cookier
	WDefer            func(call interface{}, args ...interface{}) error
	WGlobal           func() vweb.Globaler
	WHeader           func() http.Header
	WRequest          func() *http.Request
	WRequestLimitSize func(l int64) *http.Request
	WResponse         func() vweb.Responser
	WResponseWriter   func() http.ResponseWriter
	WRootDir          func(path string) string
	WSession          func() vweb.Sessioner
	WSwap             func() *vmap.Map
	WWithContext      func(ctx context.Context)
}

func (W _github_com_456vv_vweb_v2_TemplateDoter) Context() context.Context {
	return W.WContext()
}
func (W _github_com_456vv_vweb_v2_TemplateDoter) Cookie() vweb.Cookier {
	return W.WCookie()
}
func (W _github_com_456vv_vweb_v2_TemplateDoter) Defer(call interface{}, args ...interface{}) error {
	return W.WDefer(call, args...)
}
func (W _github_com_456vv_vweb_v2_TemplateDoter) Global() vweb.Globaler {
	return W.WGlobal()
}
func (W _github_com_456vv_vweb_v2_TemplateDoter) Header() http.Header {
	return W.WHeader()
}
func (W _github_com_456vv_vweb_v2_TemplateDoter) Request() *http.Request {
	return W.WRequest()
}
func (W _github_com_456vv_vweb_v2_TemplateDoter) RequestLimitSize(l int64) *http.Request {
	return W.WRequestLimitSize(l)
}
func (W _github_com_456vv_vweb_v2_TemplateDoter) Response() vweb.Responser {
	return W.WResponse()
}
func (W _github_com_456vv_vweb_v2_TemplateDoter) ResponseWriter() http.ResponseWriter {
	return W.WResponseWriter()
}
func (W _github_com_456vv_vweb_v2_TemplateDoter) RootDir(path string) string {
	return W.WRootDir(path)
}
func (W _github_com_456vv_vweb_v2_TemplateDoter) Session() vweb.Sessioner {
	return W.WSession()
}
func (W _github_com_456vv_vweb_v2_TemplateDoter) Swap() *vmap.Map {
	return W.WSwap()
}
func (W _github_com_456vv_vweb_v2_TemplateDoter) WithContext(ctx context.Context) {
	W.WWithContext(ctx)
}
