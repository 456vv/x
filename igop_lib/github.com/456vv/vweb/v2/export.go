// export by github.com/goplus/igop/cmd/qexp

//go:build igop_lib
// +build igop_lib

package vweb

import (
	q "github.com/456vv/vweb/v2"

	"go/constant"
	"reflect"

	"github.com/goplus/igop"
)

func init() {
	igop.RegisterPackage(&igop.Package{
		Name: "vweb",
		Path: "github.com/456vv/vweb/v2",
		Deps: map[string]string{
			"bufio":                             "bufio",
			"bytes":                             "bytes",
			"context":                           "context",
			"crypto/rand":                       "rand",
			"crypto/tls":                        "tls",
			"crypto/x509":                       "x509",
			"encoding/gob":                      "gob",
			"errors":                            "errors",
			"fmt":                               "fmt",
			"github.com/456vv/vconnpool/v2":     "vconnpool",
			"github.com/456vv/verror":           "verror",
			"github.com/456vv/vmap/v2":          "vmap",
			"github.com/456vv/vweb/v2/builtin":  "builtin",
			"golang.org/x/crypto/acme":          "acme",
			"golang.org/x/crypto/acme/autocert": "autocert",
			"golang.org/x/net/http/httpguts":    "httpguts",
			"io":                                "io",
			"log":                               "log",
			"math/rand":                         "rand",
			"mime/multipart":                    "multipart",
			"net":                               "net",
			"net/http":                          "http",
			"net/rpc":                           "rpc",
			"net/textproto":                     "textproto",
			"net/url":                           "url",
			"os":                                "os",
			"path":                              "path",
			"path/filepath":                     "filepath",
			"reflect":                           "reflect",
			"regexp":                            "regexp",
			"runtime":                           "runtime",
			"strconv":                           "strconv",
			"strings":                           "strings",
			"sync":                              "sync",
			"sync/atomic":                       "atomic",
			"time":                              "time",
			"unicode/utf8":                      "utf8",
		},
		Interfaces: map[string]reflect.Type{
			"Cookier":          reflect.TypeOf((*q.Cookier)(nil)).Elem(),
			"DotContexter":     reflect.TypeOf((*q.DotContexter)(nil)).Elem(),
			"Doter":            reflect.TypeOf((*q.Doter)(nil)).Elem(),
			"DynamicTemplater": reflect.TypeOf((*q.DynamicTemplater)(nil)).Elem(),
			"Globaler":         reflect.TypeOf((*q.Globaler)(nil)).Elem(),
			"PluginHTTP":       reflect.TypeOf((*q.PluginHTTP)(nil)).Elem(),
			"PluginRPC":        reflect.TypeOf((*q.PluginRPC)(nil)).Elem(),
			"Pluginer":         reflect.TypeOf((*q.Pluginer)(nil)).Elem(),
			"Responser":        reflect.TypeOf((*q.Responser)(nil)).Elem(),
			"Sessioner":        reflect.TypeOf((*q.Sessioner)(nil)).Elem(),
		},
		NamedTypes: map[string]reflect.Type{
			"Cookie":               reflect.TypeOf((*q.Cookie)(nil)).Elem(),
			"Dot":                  reflect.TypeOf((*q.Dot)(nil)).Elem(),
			"DynamicTemplateFunc":  reflect.TypeOf((*q.DynamicTemplateFunc)(nil)).Elem(),
			"ExecCall":             reflect.TypeOf((*q.ExecCall)(nil)).Elem(),
			"ExitCall":             reflect.TypeOf((*q.ExitCall)(nil)).Elem(),
			"Forward":              reflect.TypeOf((*q.Forward)(nil)).Elem(),
			"HandleFunc":           reflect.TypeOf((*q.HandleFunc)(nil)).Elem(),
			"PluginHTTPClient":     reflect.TypeOf((*q.PluginHTTPClient)(nil)).Elem(),
			"PluginRPCClient":      reflect.TypeOf((*q.PluginRPCClient)(nil)).Elem(),
			"PluginType":           reflect.TypeOf((*q.PluginType)(nil)).Elem(),
			"Route":                reflect.TypeOf((*q.Route)(nil)).Elem(),
			"ServerHandlerDynamic": reflect.TypeOf((*q.ServerHandlerDynamic)(nil)).Elem(),
			"ServerHandlerStatic":  reflect.TypeOf((*q.ServerHandlerStatic)(nil)).Elem(),
			"Session":              reflect.TypeOf((*q.Session)(nil)).Elem(),
			"Sessions":             reflect.TypeOf((*q.Sessions)(nil)).Elem(),
			"Site":                 reflect.TypeOf((*q.Site)(nil)).Elem(),
			"SiteMan":              reflect.TypeOf((*q.SiteMan)(nil)).Elem(),
			"SitePool":             reflect.TypeOf((*q.SitePool)(nil)).Elem(),
		},
		AliasTypes: map[string]reflect.Type{},
		Vars: map[string]reflect.Value{
			"ConnContextKey":     reflect.ValueOf(&q.ConnContextKey),
			"DefaultSitePool":    reflect.ValueOf(&q.DefaultSitePool),
			"ListenerContextKey": reflect.ValueOf(&q.ListenerContextKey),
			"PluginContextKey":   reflect.ValueOf(&q.PluginContextKey),
			"SiteContextKey":     reflect.ValueOf(&q.SiteContextKey),
		},
		Funcs: map[string]reflect.Value{
			"AddSalt":              reflect.ValueOf(q.AddSalt),
			"AutoCert":             reflect.ValueOf(q.AutoCert),
			"CopyStruct":           reflect.ValueOf(q.CopyStruct),
			"CopyStructDeep":       reflect.ValueOf(q.CopyStructDeep),
			"DepthField":           reflect.ValueOf(q.DepthField),
			"ExecFunc":             reflect.ValueOf(q.ExecFunc),
			"ForMethod":            reflect.ValueOf(q.ForMethod),
			"ForType":              reflect.ValueOf(q.ForType),
			"GenerateRandom":       reflect.ValueOf(q.GenerateRandom),
			"GenerateRandomId":     reflect.ValueOf(q.GenerateRandomId),
			"GenerateRandomString": reflect.ValueOf(q.GenerateRandomString),
			"InDirect":             reflect.ValueOf(q.InDirect),
			"NewSitePool":          reflect.ValueOf(q.NewSitePool),
			"PagePath":             reflect.ValueOf(q.PagePath),
		},
		TypedConsts: map[string]igop.TypedConst{
			"PluginTypeHTTP": {reflect.TypeOf(q.PluginTypeHTTP), constant.MakeInt64(int64(q.PluginTypeHTTP))},
			"PluginTypeRPC":  {reflect.TypeOf(q.PluginTypeRPC), constant.MakeInt64(int64(q.PluginTypeRPC))},
		},
		UntypedConsts: map[string]igop.UntypedConst{},
	})
}
