// export by github.com/goplus/gossa/cmd/qexp

//go:build gossa_lib
// +build gossa_lib

package viot

import (
	q "github.com/456vv/viot/v2"

	"go/constant"
	"reflect"

	"github.com/goplus/gossa"
)

func init() {
	gossa.RegisterPackage(&gossa.Package{
		Name: "viot",
		Path: "github.com/456vv/viot/v2",
		Deps: map[string]string{
			"bufio":                            "bufio",
			"bytes":                            "bytes",
			"context":                          "context",
			"crypto/rand":                      "rand",
			"crypto/tls":                       "tls",
			"encoding/base64":                  "base64",
			"encoding/json":                    "json",
			"errors":                           "errors",
			"fmt":                              "fmt",
			"github.com/456vv/vconn":           "vconn",
			"github.com/456vv/vconnpool/v2":    "vconnpool",
			"github.com/456vv/verror":          "verror",
			"github.com/456vv/vmap/v2":         "vmap",
			"github.com/456vv/vweb/v2":         "vweb",
			"github.com/456vv/vweb/v2/builtin": "builtin",
			"golang.org/x/net/http/httpguts":   "httpguts",
			"io":                               "io",
			"io/ioutil":                        "ioutil",
			"log":                              "log",
			"math":                             "math",
			"math/big":                         "big",
			"net":                              "net",
			"net/http":                         "http",
			"net/textproto":                    "textproto",
			"net/url":                          "url",
			"os":                               "os",
			"path":                             "path",
			"path/filepath":                    "filepath",
			"reflect":                          "reflect",
			"regexp":                           "regexp",
			"runtime":                          "runtime",
			"strconv":                          "strconv",
			"strings":                          "strings",
			"sync":                             "sync",
			"sync/atomic":                      "atomic",
			"time":                             "time",
		},
		Interfaces: map[string]reflect.Type{
			"CloseNotifier":    reflect.TypeOf((*q.CloseNotifier)(nil)).Elem(),
			"DotContexter":     reflect.TypeOf((*q.DotContexter)(nil)).Elem(),
			"DynamicTemplater": reflect.TypeOf((*q.DynamicTemplater)(nil)).Elem(),
			"Flusher":          reflect.TypeOf((*q.Flusher)(nil)).Elem(),
			"Handler":          reflect.TypeOf((*q.Handler)(nil)).Elem(),
			"Hijacker":         reflect.TypeOf((*q.Hijacker)(nil)).Elem(),
			"Launcher":         reflect.TypeOf((*q.Launcher)(nil)).Elem(),
			"RawControler":     reflect.TypeOf((*q.RawControler)(nil)).Elem(),
			"ResponseWriter":   reflect.TypeOf((*q.ResponseWriter)(nil)).Elem(),
			"RoundTripper":     reflect.TypeOf((*q.RoundTripper)(nil)).Elem(),
			"TemplateDoter":    reflect.TypeOf((*q.TemplateDoter)(nil)).Elem(),
		},
		NamedTypes: map[string]gossa.NamedType{
			"Client":               {reflect.TypeOf((*q.Client)(nil)).Elem(), "", "Do,DoCtx,Get,GetCtx,Post,PostCtx"},
			"ConnState":            {reflect.TypeOf((*q.ConnState)(nil)).Elem(), "String", ""},
			"DynamicTemplateFunc":  {reflect.TypeOf((*q.DynamicTemplateFunc)(nil)).Elem(), "", ""},
			"HandlerFunc":          {reflect.TypeOf((*q.HandlerFunc)(nil)).Elem(), "ServeIOT", ""},
			"Header":               {reflect.TypeOf((*q.Header)(nil)).Elem(), "Clone,Del,Get,Set", ""},
			"LogLevel":             {reflect.TypeOf((*q.LogLevel)(nil)).Elem(), "", ""},
			"Request":              {reflect.TypeOf((*q.Request)(nil)).Elem(), "", "Context,GetBasicAuth,GetBody,GetNonce,GetTokenAuth,ProtoAtLeast,RequestConfig,SetBasicAuth,SetBody,SetTokenAuth,WithContext,wantsClose"},
			"RequestConfig":        {reflect.TypeOf((*q.RequestConfig)(nil)).Elem(), "", "GetBody,Marshal,SetBody,Unmarshal"},
			"Response":             {reflect.TypeOf((*q.Response)(nil)).Elem(), "", "ResponseConfig,SetNonce,WriteAt,WriteTo"},
			"ResponseConfig":       {reflect.TypeOf((*q.ResponseConfig)(nil)).Elem(), "", "Marshal"},
			"Route":                {reflect.TypeOf((*q.Route)(nil)).Elem(), "", "HandleFunc,ServeIOT,SetSiteMan"},
			"Server":               {reflect.TypeOf((*q.Server)(nil)).Elem(), "", "Close,ListenAndServe,RegisterOnShutdown,Serve,SetKeepAlivesEnabled,Shutdown,closeConns,closeDoneChan,closeIdleConns,closeListeners,doKeepAlives,idleTimeout,init,logDebugReadData,logDebugWriteData,logf,maxLineBytes,shuttingDown,trackConn,trackListener"},
			"ServerHandlerDynamic": {reflect.TypeOf((*q.ServerHandlerDynamic)(nil)).Elem(), "", "Execute,Parse,ParseFile,ParseText,ServeIOT"},
			"SitePool":             {reflect.TypeOf((*q.SitePool)(nil)).Elem(), "", "NewSite"},
			"TemplateDot":          {reflect.TypeOf((*q.TemplateDot)(nil)).Elem(), "", "Context,Defer,Free,Global,Header,Hijack,Launch,Request,ResponseWriter,RootDir,Session,Swap,WithContext"},
		},
		AliasTypes: map[string]reflect.Type{
			"Globaler":  reflect.TypeOf((*q.Globaler)(nil)).Elem(),
			"Session":   reflect.TypeOf((*q.Session)(nil)).Elem(),
			"Sessioner": reflect.TypeOf((*q.Sessioner)(nil)).Elem(),
			"Sessions":  reflect.TypeOf((*q.Sessions)(nil)).Elem(),
			"Site":      reflect.TypeOf((*q.Site)(nil)).Elem(),
			"SiteMan":   reflect.TypeOf((*q.SiteMan)(nil)).Elem(),
		},
		Vars: map[string]reflect.Value{
			"ConnContextKey":      reflect.ValueOf(&q.ConnContextKey),
			"DefaultSitePool":     reflect.ValueOf(&q.DefaultSitePool),
			"ErrAbortHandler":     reflect.ValueOf(&q.ErrAbortHandler),
			"ErrBodyNotAllowed":   reflect.ValueOf(&q.ErrBodyNotAllowed),
			"ErrConnClose":        reflect.ValueOf(&q.ErrConnClose),
			"ErrDoned":            reflect.ValueOf(&q.ErrDoned),
			"ErrGetBodyed":        reflect.ValueOf(&q.ErrGetBodyed),
			"ErrHijacked":         reflect.ValueOf(&q.ErrHijacked),
			"ErrHostInvalid":      reflect.ValueOf(&q.ErrHostInvalid),
			"ErrLaunched":         reflect.ValueOf(&q.ErrLaunched),
			"ErrMethodInvalid":    reflect.ValueOf(&q.ErrMethodInvalid),
			"ErrProtoInvalid":     reflect.ValueOf(&q.ErrProtoInvalid),
			"ErrReqUnavailable":   reflect.ValueOf(&q.ErrReqUnavailable),
			"ErrRespUnavailable":  reflect.ValueOf(&q.ErrRespUnavailable),
			"ErrRwaControl":       reflect.ValueOf(&q.ErrRwaControl),
			"ErrServerClosed":     reflect.ValueOf(&q.ErrServerClosed),
			"ErrURIInvalid":       reflect.ValueOf(&q.ErrURIInvalid),
			"ListenerContextKey":  reflect.ValueOf(&q.ListenerContextKey),
			"LocalAddrContextKey": reflect.ValueOf(&q.LocalAddrContextKey),
			"ServerContextKey":    reflect.ValueOf(&q.ServerContextKey),
			"SiteContextKey":      reflect.ValueOf(&q.SiteContextKey),
		},
		Funcs: map[string]reflect.Value{
			"Error":                 reflect.ValueOf(q.Error),
			"NewRequest":            reflect.ValueOf(q.NewRequest),
			"NewRequestWithContext": reflect.ValueOf(q.NewRequestWithContext),
			"NewSitePool":           reflect.ValueOf(q.NewSitePool),
			"Nonce":                 reflect.ValueOf(q.Nonce),
			"ParseIOTVersion":       reflect.ValueOf(q.ParseIOTVersion),
			"ReadRequest":           reflect.ValueOf(q.ReadRequest),
			"ReadResponse":          reflect.ValueOf(q.ReadResponse),
		},
		TypedConsts: map[string]gossa.TypedConst{
			"LogDebug":      {reflect.TypeOf(q.LogDebug), constant.MakeInt64(int64(q.LogDebug))},
			"LogErr":        {reflect.TypeOf(q.LogErr), constant.MakeInt64(int64(q.LogErr))},
			"LogNone":       {reflect.TypeOf(q.LogNone), constant.MakeInt64(int64(q.LogNone))},
			"StateActive":   {reflect.TypeOf(q.StateActive), constant.MakeInt64(int64(q.StateActive))},
			"StateClosed":   {reflect.TypeOf(q.StateClosed), constant.MakeInt64(int64(q.StateClosed))},
			"StateHijacked": {reflect.TypeOf(q.StateHijacked), constant.MakeInt64(int64(q.StateHijacked))},
			"StateIdle":     {reflect.TypeOf(q.StateIdle), constant.MakeInt64(int64(q.StateIdle))},
			"StateNew":      {reflect.TypeOf(q.StateNew), constant.MakeInt64(int64(q.StateNew))},
		},
		UntypedConsts: map[string]gossa.UntypedConst{
			"DefaultLineBytes": {"untyped int", constant.MakeInt64(int64(q.DefaultLineBytes))},
		},
	})
}
