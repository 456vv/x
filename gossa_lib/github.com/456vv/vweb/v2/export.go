// export by github.com/goplus/gossa/cmd/qexp

//go:build gossa_lib
// +build gossa_lib

package vweb

import (
	q "github.com/456vv/vweb/v2"

	"go/constant"
	"reflect"

	"github.com/goplus/gossa"
)

func init() {
	gossa.RegisterPackage(&gossa.Package{
		Name: "vweb",
		Path: "github.com/456vv/vweb/v2",
		Deps: map[string]string{
			"bufio":                            "bufio",
			"bytes":                            "bytes",
			"context":                          "context",
			"crypto/rand":                      "rand",
			"crypto/tls":                       "tls",
			"encoding/gob":                     "gob",
			"errors":                           "errors",
			"fmt":                              "fmt",
			"github.com/456vv/vconnpool/v2":    "vconnpool",
			"github.com/456vv/verror":          "verror",
			"github.com/456vv/vmap/v2":         "vmap",
			"github.com/456vv/vweb/v2/builtin": "builtin",
			"golang.org/x/net/http/httpguts":   "httpguts",
			"io":                               "io",
			"io/ioutil":                        "ioutil",
			"log":                              "log",
			"math/rand":                        "rand",
			"mime/multipart":                   "multipart",
			"net":                              "net",
			"net/http":                         "http",
			"net/rpc":                          "rpc",
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
			"Cookier":          reflect.TypeOf((*q.Cookier)(nil)).Elem(),
			"DotContexter":     reflect.TypeOf((*q.DotContexter)(nil)).Elem(),
			"DynamicTemplater": reflect.TypeOf((*q.DynamicTemplater)(nil)).Elem(),
			"Globaler":         reflect.TypeOf((*q.Globaler)(nil)).Elem(),
			"PluginHTTP":       reflect.TypeOf((*q.PluginHTTP)(nil)).Elem(),
			"PluginRPC":        reflect.TypeOf((*q.PluginRPC)(nil)).Elem(),
			"Pluginer":         reflect.TypeOf((*q.Pluginer)(nil)).Elem(),
			"Responser":        reflect.TypeOf((*q.Responser)(nil)).Elem(),
			"Sessioner":        reflect.TypeOf((*q.Sessioner)(nil)).Elem(),
			"TemplateDoter":    reflect.TypeOf((*q.TemplateDoter)(nil)).Elem(),
		},
		NamedTypes: map[string]gossa.NamedType{
			"Cookie":               {reflect.TypeOf((*q.Cookie)(nil)).Elem(), "", "Add,Del,Get,ReadAll,RemoveAll"},
			"DynamicTemplateFunc":  {reflect.TypeOf((*q.DynamicTemplateFunc)(nil)).Elem(), "", ""},
			"ExecCall":             {reflect.TypeOf((*q.ExecCall)(nil)).Elem(), "", "Exec,Func"},
			"ExitCall":             {reflect.TypeOf((*q.ExitCall)(nil)).Elem(), "", "Defer,Free"},
			"Forward":              {reflect.TypeOf((*q.Forward)(nil)).Elem(), "", "Rewrite"},
			"PluginHTTPClient":     {reflect.TypeOf((*q.PluginHTTPClient)(nil)).Elem(), "", "Connection"},
			"PluginRPCClient":      {reflect.TypeOf((*q.PluginRPCClient)(nil)).Elem(), "", "Connection"},
			"PluginType":           {reflect.TypeOf((*q.PluginType)(nil)).Elem(), "", ""},
			"Route":                {reflect.TypeOf((*q.Route)(nil)).Elem(), "", "HandleFunc,ServeHTTP,SetSiteMan"},
			"ServerHandlerDynamic": {reflect.TypeOf((*q.ServerHandlerDynamic)(nil)).Elem(), "", "Execute,Parse,ParseFile,ParseText,ServeHTTP"},
			"ServerHandlerStatic":  {reflect.TypeOf((*q.ServerHandlerStatic)(nil)).Elem(), "", "ServeHTTP,body,header"},
			"Session":              {reflect.TypeOf((*q.Session)(nil)).Elem(), "", "Defer,Free,Token"},
			"Sessions":             {reflect.TypeOf((*q.Sessions)(nil)).Elem(), "", "DelSession,GetSession,Len,NewSession,ProcessDeadAll,Session,SessionId,SetSession,generateRandSessionId,generateSessionId,generateSessionIdSalt,setSession,triggerDeadSession,writeToClient"},
			"Site":                 {reflect.TypeOf((*q.Site)(nil)).Elem(), "", "PoolName"},
			"SiteMan":              {reflect.TypeOf((*q.SiteMan)(nil)).Elem(), "", "Add,Get,Range,derogatoryDomain"},
			"SitePool":             {reflect.TypeOf((*q.SitePool)(nil)).Elem(), "", "Close,DelSite,NewSite,RangeSite,SetRecoverSession,Start,start"},
			"TemplateDot":          {reflect.TypeOf((*q.TemplateDot)(nil)).Elem(), "", "Context,Cookie,Defer,Free,Global,Header,Request,RequestLimitSize,Response,ResponseWriter,RootDir,Session,Swap,WithContext"},
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
		TypedConsts: map[string]gossa.TypedConst{
			"PluginTypeHTTP": {reflect.TypeOf(q.PluginTypeHTTP), constant.MakeInt64(int64(q.PluginTypeHTTP))},
			"PluginTypeRPC":  {reflect.TypeOf(q.PluginTypeRPC), constant.MakeInt64(int64(q.PluginTypeRPC))},
		},
		UntypedConsts: map[string]gossa.UntypedConst{},
	})
}
