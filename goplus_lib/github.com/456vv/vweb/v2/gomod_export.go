// Package vweb provide Go+ "github.com/456vv/vweb/v2" package, as "github.com/456vv/vweb/v2" package in Go.
package vweb

import (
	context "context"
	io "io"
	http "net/http"
	reflect "reflect"
	template "text/template"
	time "time"

	vweb "github.com/456vv/vweb/v2"
	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execAddSalt(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := vweb.AddSalt(args[0].([]byte), args[1].(string))
	p.Ret(2, ret0)
}

func execmCookieReadAll(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vweb.Cookie).ReadAll()
	p.Ret(1, ret0)
}

func execmCookieRemoveAll(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*vweb.Cookie).RemoveAll()
	p.PopN(1)
}

func execmCookieGet(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*vweb.Cookie).Get(args[1].(string))
	p.Ret(2, ret0)
}

func execmCookieAdd(_ int, p *gop.Context) {
	args := p.GetArgs(9)
	args[0].(*vweb.Cookie).Add(args[1].(string), args[2].(string), args[3].(string), args[4].(string), args[5].(int), args[6].(bool), args[7].(bool), args[8].(http.SameSite))
	p.PopN(9)
}

func execmCookieDel(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*vweb.Cookie).Del(args[1].(string))
	p.PopN(2)
}

func execiCookierAdd(_ int, p *gop.Context) {
	args := p.GetArgs(9)
	args[0].(vweb.Cookier).Add(args[1].(string), args[2].(string), args[3].(string), args[4].(string), args[5].(int), args[6].(bool), args[7].(bool), args[8].(http.SameSite))
	p.PopN(9)
}

func execiCookierDel(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(vweb.Cookier).Del(args[1].(string))
	p.PopN(2)
}

func execiCookierGet(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(vweb.Cookier).Get(args[1].(string))
	p.Ret(2, ret0)
}

func execiCookierReadAll(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(vweb.Cookier).ReadAll()
	p.Ret(1, ret0)
}

func execiCookierRemoveAll(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(vweb.Cookier).RemoveAll()
	p.PopN(1)
}

func execCopyStruct(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := vweb.CopyStruct(args[0], args[1], args[2].(func(name string, dsc reflect.Value, src reflect.Value) bool))
	p.Ret(3, ret0)
}

func execCopyStructDeep(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := vweb.CopyStructDeep(args[0], args[1], args[2].(func(name string, dsc reflect.Value, src reflect.Value) bool))
	p.Ret(3, ret0)
}

func execDepthField(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := vweb.DepthField(args[0], args[1:]...)
	p.Ret(arity, ret0, ret1)
}

func execiDotContexterContext(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(vweb.DotContexter).Context()
	p.Ret(1, ret0)
}

func toType0(v interface{}) context.Context {
	if v == nil {
		return nil
	}
	return v.(context.Context)
}

func execiDotContexterWithContext(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(vweb.DotContexter).WithContext(toType0(args[1]))
	p.PopN(2)
}

func toType1(v interface{}) io.Writer {
	if v == nil {
		return nil
	}
	return v.(io.Writer)
}

func execiDynamicTemplaterExecute(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(vweb.DynamicTemplater).Execute(toType1(args[1]), args[2])
	p.Ret(3, ret0)
}

func toType2(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func execiDynamicTemplaterParse(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(vweb.DynamicTemplater).Parse(toType2(args[1]))
	p.Ret(2, ret0)
}

func execiDynamicTemplaterSetPath(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(vweb.DynamicTemplater).SetPath(args[1].(string), args[2].(string))
	p.PopN(3)
}

func execExecFunc(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := vweb.ExecFunc(args[0], args[1:]...)
	p.Ret(arity, ret0, ret1)
}

func execmExitCallDefer(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*vweb.ExitCall).Defer(args[1], args[2:]...)
	p.Ret(arity, ret0)
}

func execmExitCallFree(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*vweb.ExitCall).Free()
	p.PopN(1)
}

func execExtendTemplatePackage(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	vweb.ExtendTemplatePackage(args[0].(string), args[1].(template.FuncMap))
	p.PopN(2)
}

func execForMethod(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := vweb.ForMethod(args[0])
	p.Ret(1, ret0)
}

func execForType(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := vweb.ForType(args[0], args[1].(bool))
	p.Ret(2, ret0)
}

func execmForwardRewrite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1, ret2 := args[0].(*vweb.Forward).Rewrite(args[1].(string))
	p.Ret(2, ret0, ret1, ret2)
}

func execGenerateRandom(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := vweb.GenerateRandom(args[0].(int))
	p.Ret(1, ret0, ret1)
}

func execGenerateRandomId(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := vweb.GenerateRandomId(args[0].([]byte))
	p.Ret(1, ret0)
}

func execGenerateRandomString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := vweb.GenerateRandomString(args[0].(int))
	p.Ret(1, ret0, ret1)
}

func execiGlobalerDel(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(vweb.Globaler).Del(args[1])
	p.PopN(2)
}

func execiGlobalerGet(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(vweb.Globaler).Get(args[1])
	p.Ret(2, ret0)
}

func execiGlobalerHas(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(vweb.Globaler).Has(args[1])
	p.Ret(2, ret0)
}

func execiGlobalerReset(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(vweb.Globaler).Reset()
	p.PopN(1)
}

func execiGlobalerSet(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(vweb.Globaler).Set(args[1], args[2])
	p.PopN(3)
}

func execiGlobalerSetExpired(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(vweb.Globaler).SetExpired(args[1], args[2].(time.Duration))
	p.PopN(3)
}

func execiGlobalerSetExpiredCall(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	args[0].(vweb.Globaler).SetExpiredCall(args[1], args[2].(time.Duration), args[3].(func(interface{})))
	p.PopN(4)
}

func execInDirect(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := vweb.InDirect(args[0].(reflect.Value))
	p.Ret(1, ret0)
}

func execNewSitePool(_ int, p *gop.Context) {
	ret0 := vweb.NewSitePool()
	p.Ret(0, ret0)
}

func execPagePath(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1, ret2 := vweb.PagePath(args[0].(string), args[1].(string), args[2].([]string))
	p.Ret(3, ret0, ret1, ret2)
}

func execiPluginHTTPCancelRequest(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(vweb.PluginHTTP).CancelRequest(args[1].(*http.Request))
	p.PopN(2)
}

func execiPluginHTTPCloseIdleConnections(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(vweb.PluginHTTP).CloseIdleConnections()
	p.PopN(1)
}

func toType3(v interface{}) http.RoundTripper {
	if v == nil {
		return nil
	}
	return v.(http.RoundTripper)
}

func execiPluginHTTPRegisterProtocol(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(vweb.PluginHTTP).RegisterProtocol(args[1].(string), toType3(args[2]))
	p.PopN(3)
}

func execiPluginHTTPRoundTrip(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(vweb.PluginHTTP).RoundTrip(args[1].(*http.Request))
	p.Ret(2, ret0, ret1)
}

func toType4(v interface{}) http.ResponseWriter {
	if v == nil {
		return nil
	}
	return v.(http.ResponseWriter)
}

func execiPluginHTTPServeHTTP(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(vweb.PluginHTTP).ServeHTTP(toType4(args[1]), args[2].(*http.Request))
	p.PopN(3)
}

func execiPluginHTTPType(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(vweb.PluginHTTP).Type()
	p.Ret(1, ret0)
}

func execmPluginHTTPClientConnection(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*vweb.PluginHTTPClient).Connection()
	p.Ret(1, ret0, ret1)
}

func execiPluginRPCCall(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(vweb.PluginRPC).Call(args[1].(string), args[2])
	p.Ret(3, ret0, ret1)
}

func execiPluginRPCClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(vweb.PluginRPC).Close()
	p.Ret(1, ret0)
}

func execiPluginRPCDiscard(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(vweb.PluginRPC).Discard()
	p.Ret(1, ret0)
}

func execiPluginRPCRegister(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(vweb.PluginRPC).Register(args[1])
	p.PopN(2)
}

func execiPluginRPCType(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(vweb.PluginRPC).Type()
	p.Ret(1, ret0)
}

func execmPluginRPCClientConnection(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*vweb.PluginRPCClient).Connection()
	p.Ret(1, ret0, ret1)
}

func execiResponserError(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(vweb.Responser).Error(args[1].(string), args[2].(int))
	p.PopN(3)
}

func execiResponserFlush(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(vweb.Responser).Flush()
	p.PopN(1)
}

func execiResponserHijack(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2 := args[0].(vweb.Responser).Hijack()
	p.Ret(1, ret0, ret1, ret2)
}

func execiResponserPush(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(vweb.Responser).Push(args[1].(string), args[2].(*http.PushOptions))
	p.Ret(3, ret0)
}

func execiResponserReadFrom(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(vweb.Responser).ReadFrom(toType2(args[1]))
	p.Ret(2, ret0, ret1)
}

func execiResponserRedirect(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(vweb.Responser).Redirect(args[1].(string), args[2].(int))
	p.PopN(3)
}

func execiResponserWrite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(vweb.Responser).Write(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execiResponserWriteHeader(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(vweb.Responser).WriteHeader(args[1].(int))
	p.PopN(2)
}

func execiResponserWriteString(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(vweb.Responser).WriteString(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmRouteHandleFunc(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*vweb.Route).HandleFunc(args[1].(string), args[2].(func(w http.ResponseWriter, r *http.Request)))
	p.PopN(3)
}

func execmRouteServeHTTP(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*vweb.Route).ServeHTTP(toType4(args[1]), args[2].(*http.Request))
	p.PopN(3)
}

func execmServerHandlerDynamicServeHTTP(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*vweb.ServerHandlerDynamic).ServeHTTP(toType4(args[1]), args[2].(*http.Request))
	p.PopN(3)
}

func execmServerHandlerDynamicParseText(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*vweb.ServerHandlerDynamic).ParseText(args[1].(string), args[2].(string))
	p.Ret(3, ret0)
}

func execmServerHandlerDynamicParseFile(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*vweb.ServerHandlerDynamic).ParseFile(args[1].(string))
	p.Ret(2, ret0)
}

func execmServerHandlerDynamicParse(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*vweb.ServerHandlerDynamic).Parse(toType2(args[1]))
	p.Ret(2, ret0)
}

func execmServerHandlerDynamicExecute(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*vweb.ServerHandlerDynamic).Execute(toType1(args[1]), args[2])
	p.Ret(3, ret0)
}

func execmServerHandlerDynamicTemplateExtendNewFunc(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*vweb.ServerHandlerDynamicTemplateExtend).NewFunc(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmServerHandlerDynamicTemplateExtendCall(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*vweb.ServerHandlerDynamicTemplateExtend).Call(args[1].(func([]reflect.Value) []reflect.Value), args[2:]...)
	p.Ret(arity, ret0)
}

func execmServerHandlerDynamicTemplateExtendExecuteTemplate(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := args[0].(*vweb.ServerHandlerDynamicTemplateExtend).ExecuteTemplate(toType1(args[1]), args[2].(string), args[3])
	p.Ret(4, ret0)
}

func execmServerHandlerStaticServeHTTP(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*vweb.ServerHandlerStatic).ServeHTTP(toType4(args[1]), args[2].(*http.Request))
	p.PopN(3)
}

func execmSessionToken(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vweb.Session).Token()
	p.Ret(1, ret0)
}

func execmSessionDefer(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*vweb.Session).Defer(args[1], args[2:]...)
	p.Ret(arity, ret0)
}

func execmSessionFree(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*vweb.Session).Free()
	p.PopN(1)
}

func execiSessionerDefer(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(vweb.Sessioner).Defer(args[1], args[2:]...)
	p.Ret(arity, ret0)
}

func execiSessionerDel(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(vweb.Sessioner).Del(args[1])
	p.PopN(2)
}

func execiSessionerFree(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(vweb.Sessioner).Free()
	p.PopN(1)
}

func execiSessionerGet(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(vweb.Sessioner).Get(args[1])
	p.Ret(2, ret0)
}

func execiSessionerGetHas(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(vweb.Sessioner).GetHas(args[1])
	p.Ret(2, ret0, ret1)
}

func execiSessionerHas(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(vweb.Sessioner).Has(args[1])
	p.Ret(2, ret0)
}

func execiSessionerReset(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(vweb.Sessioner).Reset()
	p.PopN(1)
}

func execiSessionerSet(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(vweb.Sessioner).Set(args[1], args[2])
	p.PopN(3)
}

func execiSessionerSetExpired(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(vweb.Sessioner).SetExpired(args[1], args[2].(time.Duration))
	p.PopN(3)
}

func execiSessionerSetExpiredCall(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	args[0].(vweb.Sessioner).SetExpiredCall(args[1], args[2].(time.Duration), args[3].(func(interface{})))
	p.PopN(4)
}

func execiSessionerToken(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(vweb.Sessioner).Token()
	p.Ret(1, ret0)
}

func execmSessionsLen(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vweb.Sessions).Len()
	p.Ret(1, ret0)
}

func execmSessionsProcessDeadAll(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vweb.Sessions).ProcessDeadAll()
	p.Ret(1, ret0)
}

func execmSessionsSessionId(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*vweb.Sessions).SessionId(args[1].(*http.Request))
	p.Ret(2, ret0, ret1)
}

func execmSessionsNewSession(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*vweb.Sessions).NewSession(args[1].(string))
	p.Ret(2, ret0)
}

func execmSessionsGetSession(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*vweb.Sessions).GetSession(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func toType5(v interface{}) vweb.Sessioner {
	if v == nil {
		return nil
	}
	return v.(vweb.Sessioner)
}

func execmSessionsSetSession(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*vweb.Sessions).SetSession(args[1].(string), toType5(args[2]))
	p.Ret(3, ret0)
}

func execmSessionsDelSession(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*vweb.Sessions).DelSession(args[1].(string))
	p.PopN(2)
}

func execmSessionsSession(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*vweb.Sessions).Session(toType4(args[1]), args[2].(*http.Request))
	p.Ret(3, ret0)
}

func execmSitePoolName(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vweb.Site).PoolName()
	p.Ret(1, ret0)
}

func execmSiteManAdd(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*vweb.SiteMan).Add(args[1].(string), args[2].(*vweb.Site))
	p.PopN(3)
}

func execmSiteManGet(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*vweb.SiteMan).Get(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmSiteManRange(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*vweb.SiteMan).Range(args[1].(func(host string, site *vweb.Site) bool))
	p.PopN(2)
}

func execmSitePoolNewSite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*vweb.SitePool).NewSite(args[1].(string))
	p.Ret(2, ret0)
}

func execmSitePoolDelSite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*vweb.SitePool).DelSite(args[1].(string))
	p.PopN(2)
}

func execmSitePoolRangeSite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*vweb.SitePool).RangeSite(args[1].(func(name string, site *vweb.Site) bool))
	p.PopN(2)
}

func execmSitePoolSetRecoverSession(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*vweb.SitePool).SetRecoverSession(args[1].(time.Duration))
	p.PopN(2)
}

func execmSitePoolStart(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vweb.SitePool).Start()
	p.Ret(1, ret0)
}

func execmSitePoolClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vweb.SitePool).Close()
	p.Ret(1, ret0)
}

func execmTemplateDotRootDir(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*vweb.TemplateDot).RootDir(args[1].(string))
	p.Ret(2, ret0)
}

func execmTemplateDotRequest(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vweb.TemplateDot).Request()
	p.Ret(1, ret0)
}

func execmTemplateDotRequestLimitSize(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*vweb.TemplateDot).RequestLimitSize(args[1].(int64))
	p.Ret(2, ret0)
}

func execmTemplateDotHeader(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vweb.TemplateDot).Header()
	p.Ret(1, ret0)
}

func execmTemplateDotResponse(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vweb.TemplateDot).Response()
	p.Ret(1, ret0)
}

func execmTemplateDotResponseWriter(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vweb.TemplateDot).ResponseWriter()
	p.Ret(1, ret0)
}

func execmTemplateDotSession(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vweb.TemplateDot).Session()
	p.Ret(1, ret0)
}

func execmTemplateDotGlobal(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vweb.TemplateDot).Global()
	p.Ret(1, ret0)
}

func execmTemplateDotCookie(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vweb.TemplateDot).Cookie()
	p.Ret(1, ret0)
}

func execmTemplateDotSwap(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vweb.TemplateDot).Swap()
	p.Ret(1, ret0)
}

func execmTemplateDotDefer(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*vweb.TemplateDot).Defer(args[1], args[2:]...)
	p.Ret(arity, ret0)
}

func execmTemplateDotFree(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*vweb.TemplateDot).Free()
	p.PopN(1)
}

func execmTemplateDotContext(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vweb.TemplateDot).Context()
	p.Ret(1, ret0)
}

func execmTemplateDotWithContext(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*vweb.TemplateDot).WithContext(toType0(args[1]))
	p.PopN(2)
}

func execiTemplateDoterCookie(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(vweb.TemplateDoter).Cookie()
	p.Ret(1, ret0)
}

func execiTemplateDoterDefer(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(vweb.TemplateDoter).Defer(args[1], args[2:]...)
	p.Ret(arity, ret0)
}

func execiTemplateDoterGlobal(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(vweb.TemplateDoter).Global()
	p.Ret(1, ret0)
}

func execiTemplateDoterHeader(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(vweb.TemplateDoter).Header()
	p.Ret(1, ret0)
}

func execiTemplateDoterRequest(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(vweb.TemplateDoter).Request()
	p.Ret(1, ret0)
}

func execiTemplateDoterRequestLimitSize(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(vweb.TemplateDoter).RequestLimitSize(args[1].(int64))
	p.Ret(2, ret0)
}

func execiTemplateDoterResponse(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(vweb.TemplateDoter).Response()
	p.Ret(1, ret0)
}

func execiTemplateDoterResponseWriter(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(vweb.TemplateDoter).ResponseWriter()
	p.Ret(1, ret0)
}

func execiTemplateDoterRootDir(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(vweb.TemplateDoter).RootDir(args[1].(string))
	p.Ret(2, ret0)
}

func execiTemplateDoterSession(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(vweb.TemplateDoter).Session()
	p.Ret(1, ret0)
}

func execiTemplateDoterSwap(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(vweb.TemplateDoter).Swap()
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("github.com/456vv/vweb/v2")

func init() {
	I.RegisterFuncs(
		I.Func("AddSalt", vweb.AddSalt, execAddSalt),
		I.Func("(*Cookie).ReadAll", (*vweb.Cookie).ReadAll, execmCookieReadAll),
		I.Func("(*Cookie).RemoveAll", (*vweb.Cookie).RemoveAll, execmCookieRemoveAll),
		I.Func("(*Cookie).Get", (*vweb.Cookie).Get, execmCookieGet),
		I.Func("(*Cookie).Add", (*vweb.Cookie).Add, execmCookieAdd),
		I.Func("(*Cookie).Del", (*vweb.Cookie).Del, execmCookieDel),
		I.Func("(Cookier).Add", (vweb.Cookier).Add, execiCookierAdd),
		I.Func("(Cookier).Del", (vweb.Cookier).Del, execiCookierDel),
		I.Func("(Cookier).Get", (vweb.Cookier).Get, execiCookierGet),
		I.Func("(Cookier).ReadAll", (vweb.Cookier).ReadAll, execiCookierReadAll),
		I.Func("(Cookier).RemoveAll", (vweb.Cookier).RemoveAll, execiCookierRemoveAll),
		I.Func("CopyStruct", vweb.CopyStruct, execCopyStruct),
		I.Func("CopyStructDeep", vweb.CopyStructDeep, execCopyStructDeep),
		I.Func("(DotContexter).Context", (vweb.DotContexter).Context, execiDotContexterContext),
		I.Func("(DotContexter).WithContext", (vweb.DotContexter).WithContext, execiDotContexterWithContext),
		I.Func("(DynamicTemplater).Execute", (vweb.DynamicTemplater).Execute, execiDynamicTemplaterExecute),
		I.Func("(DynamicTemplater).Parse", (vweb.DynamicTemplater).Parse, execiDynamicTemplaterParse),
		I.Func("(DynamicTemplater).SetPath", (vweb.DynamicTemplater).SetPath, execiDynamicTemplaterSetPath),
		I.Func("(*ExitCall).Free", (*vweb.ExitCall).Free, execmExitCallFree),
		I.Func("ExtendTemplatePackage", vweb.ExtendTemplatePackage, execExtendTemplatePackage),
		I.Func("ForMethod", vweb.ForMethod, execForMethod),
		I.Func("ForType", vweb.ForType, execForType),
		I.Func("(*Forward).Rewrite", (*vweb.Forward).Rewrite, execmForwardRewrite),
		I.Func("GenerateRandom", vweb.GenerateRandom, execGenerateRandom),
		I.Func("GenerateRandomId", vweb.GenerateRandomId, execGenerateRandomId),
		I.Func("GenerateRandomString", vweb.GenerateRandomString, execGenerateRandomString),
		I.Func("(Globaler).Del", (vweb.Globaler).Del, execiGlobalerDel),
		I.Func("(Globaler).Get", (vweb.Globaler).Get, execiGlobalerGet),
		I.Func("(Globaler).Has", (vweb.Globaler).Has, execiGlobalerHas),
		I.Func("(Globaler).Reset", (vweb.Globaler).Reset, execiGlobalerReset),
		I.Func("(Globaler).Set", (vweb.Globaler).Set, execiGlobalerSet),
		I.Func("(Globaler).SetExpired", (vweb.Globaler).SetExpired, execiGlobalerSetExpired),
		I.Func("(Globaler).SetExpiredCall", (vweb.Globaler).SetExpiredCall, execiGlobalerSetExpiredCall),
		I.Func("InDirect", vweb.InDirect, execInDirect),
		I.Func("NewSitePool", vweb.NewSitePool, execNewSitePool),
		I.Func("PagePath", vweb.PagePath, execPagePath),
		I.Func("(PluginHTTP).CancelRequest", (vweb.PluginHTTP).CancelRequest, execiPluginHTTPCancelRequest),
		I.Func("(PluginHTTP).CloseIdleConnections", (vweb.PluginHTTP).CloseIdleConnections, execiPluginHTTPCloseIdleConnections),
		I.Func("(PluginHTTP).RegisterProtocol", (vweb.PluginHTTP).RegisterProtocol, execiPluginHTTPRegisterProtocol),
		I.Func("(PluginHTTP).RoundTrip", (vweb.PluginHTTP).RoundTrip, execiPluginHTTPRoundTrip),
		I.Func("(PluginHTTP).ServeHTTP", (vweb.PluginHTTP).ServeHTTP, execiPluginHTTPServeHTTP),
		I.Func("(PluginHTTP).Type", (vweb.PluginHTTP).Type, execiPluginHTTPType),
		I.Func("(*PluginHTTPClient).Connection", (*vweb.PluginHTTPClient).Connection, execmPluginHTTPClientConnection),
		I.Func("(PluginRPC).Call", (vweb.PluginRPC).Call, execiPluginRPCCall),
		I.Func("(PluginRPC).Close", (vweb.PluginRPC).Close, execiPluginRPCClose),
		I.Func("(PluginRPC).Discard", (vweb.PluginRPC).Discard, execiPluginRPCDiscard),
		I.Func("(PluginRPC).Register", (vweb.PluginRPC).Register, execiPluginRPCRegister),
		I.Func("(PluginRPC).Type", (vweb.PluginRPC).Type, execiPluginRPCType),
		I.Func("(*PluginRPCClient).Connection", (*vweb.PluginRPCClient).Connection, execmPluginRPCClientConnection),
		I.Func("(Responser).Error", (vweb.Responser).Error, execiResponserError),
		I.Func("(Responser).Flush", (vweb.Responser).Flush, execiResponserFlush),
		I.Func("(Responser).Hijack", (vweb.Responser).Hijack, execiResponserHijack),
		I.Func("(Responser).Push", (vweb.Responser).Push, execiResponserPush),
		I.Func("(Responser).ReadFrom", (vweb.Responser).ReadFrom, execiResponserReadFrom),
		I.Func("(Responser).Redirect", (vweb.Responser).Redirect, execiResponserRedirect),
		I.Func("(Responser).Write", (vweb.Responser).Write, execiResponserWrite),
		I.Func("(Responser).WriteHeader", (vweb.Responser).WriteHeader, execiResponserWriteHeader),
		I.Func("(Responser).WriteString", (vweb.Responser).WriteString, execiResponserWriteString),
		I.Func("(*Route).HandleFunc", (*vweb.Route).HandleFunc, execmRouteHandleFunc),
		I.Func("(*Route).ServeHTTP", (*vweb.Route).ServeHTTP, execmRouteServeHTTP),
		I.Func("(*ServerHandlerDynamic).ServeHTTP", (*vweb.ServerHandlerDynamic).ServeHTTP, execmServerHandlerDynamicServeHTTP),
		I.Func("(*ServerHandlerDynamic).ParseText", (*vweb.ServerHandlerDynamic).ParseText, execmServerHandlerDynamicParseText),
		I.Func("(*ServerHandlerDynamic).ParseFile", (*vweb.ServerHandlerDynamic).ParseFile, execmServerHandlerDynamicParseFile),
		I.Func("(*ServerHandlerDynamic).Parse", (*vweb.ServerHandlerDynamic).Parse, execmServerHandlerDynamicParse),
		I.Func("(*ServerHandlerDynamic).Execute", (*vweb.ServerHandlerDynamic).Execute, execmServerHandlerDynamicExecute),
		I.Func("(*ServerHandlerDynamicTemplateExtend).NewFunc", (*vweb.ServerHandlerDynamicTemplateExtend).NewFunc, execmServerHandlerDynamicTemplateExtendNewFunc),
		I.Func("(*ServerHandlerDynamicTemplateExtend).ExecuteTemplate", (*vweb.ServerHandlerDynamicTemplateExtend).ExecuteTemplate, execmServerHandlerDynamicTemplateExtendExecuteTemplate),
		I.Func("(*ServerHandlerStatic).ServeHTTP", (*vweb.ServerHandlerStatic).ServeHTTP, execmServerHandlerStaticServeHTTP),
		I.Func("(*Session).Token", (*vweb.Session).Token, execmSessionToken),
		I.Func("(*Session).Free", (*vweb.Session).Free, execmSessionFree),
		I.Func("(Sessioner).Del", (vweb.Sessioner).Del, execiSessionerDel),
		I.Func("(Sessioner).Free", (vweb.Sessioner).Free, execiSessionerFree),
		I.Func("(Sessioner).Get", (vweb.Sessioner).Get, execiSessionerGet),
		I.Func("(Sessioner).GetHas", (vweb.Sessioner).GetHas, execiSessionerGetHas),
		I.Func("(Sessioner).Has", (vweb.Sessioner).Has, execiSessionerHas),
		I.Func("(Sessioner).Reset", (vweb.Sessioner).Reset, execiSessionerReset),
		I.Func("(Sessioner).Set", (vweb.Sessioner).Set, execiSessionerSet),
		I.Func("(Sessioner).SetExpired", (vweb.Sessioner).SetExpired, execiSessionerSetExpired),
		I.Func("(Sessioner).SetExpiredCall", (vweb.Sessioner).SetExpiredCall, execiSessionerSetExpiredCall),
		I.Func("(Sessioner).Token", (vweb.Sessioner).Token, execiSessionerToken),
		I.Func("(*Sessions).Len", (*vweb.Sessions).Len, execmSessionsLen),
		I.Func("(*Sessions).ProcessDeadAll", (*vweb.Sessions).ProcessDeadAll, execmSessionsProcessDeadAll),
		I.Func("(*Sessions).SessionId", (*vweb.Sessions).SessionId, execmSessionsSessionId),
		I.Func("(*Sessions).NewSession", (*vweb.Sessions).NewSession, execmSessionsNewSession),
		I.Func("(*Sessions).GetSession", (*vweb.Sessions).GetSession, execmSessionsGetSession),
		I.Func("(*Sessions).SetSession", (*vweb.Sessions).SetSession, execmSessionsSetSession),
		I.Func("(*Sessions).DelSession", (*vweb.Sessions).DelSession, execmSessionsDelSession),
		I.Func("(*Sessions).Session", (*vweb.Sessions).Session, execmSessionsSession),
		I.Func("(*Site).PoolName", (*vweb.Site).PoolName, execmSitePoolName),
		I.Func("(*SiteMan).Add", (*vweb.SiteMan).Add, execmSiteManAdd),
		I.Func("(*SiteMan).Get", (*vweb.SiteMan).Get, execmSiteManGet),
		I.Func("(*SiteMan).Range", (*vweb.SiteMan).Range, execmSiteManRange),
		I.Func("(*SitePool).NewSite", (*vweb.SitePool).NewSite, execmSitePoolNewSite),
		I.Func("(*SitePool).DelSite", (*vweb.SitePool).DelSite, execmSitePoolDelSite),
		I.Func("(*SitePool).RangeSite", (*vweb.SitePool).RangeSite, execmSitePoolRangeSite),
		I.Func("(*SitePool).SetRecoverSession", (*vweb.SitePool).SetRecoverSession, execmSitePoolSetRecoverSession),
		I.Func("(*SitePool).Start", (*vweb.SitePool).Start, execmSitePoolStart),
		I.Func("(*SitePool).Close", (*vweb.SitePool).Close, execmSitePoolClose),
		I.Func("(*TemplateDot).RootDir", (*vweb.TemplateDot).RootDir, execmTemplateDotRootDir),
		I.Func("(*TemplateDot).Request", (*vweb.TemplateDot).Request, execmTemplateDotRequest),
		I.Func("(*TemplateDot).RequestLimitSize", (*vweb.TemplateDot).RequestLimitSize, execmTemplateDotRequestLimitSize),
		I.Func("(*TemplateDot).Header", (*vweb.TemplateDot).Header, execmTemplateDotHeader),
		I.Func("(*TemplateDot).Response", (*vweb.TemplateDot).Response, execmTemplateDotResponse),
		I.Func("(*TemplateDot).ResponseWriter", (*vweb.TemplateDot).ResponseWriter, execmTemplateDotResponseWriter),
		I.Func("(*TemplateDot).Session", (*vweb.TemplateDot).Session, execmTemplateDotSession),
		I.Func("(*TemplateDot).Global", (*vweb.TemplateDot).Global, execmTemplateDotGlobal),
		I.Func("(*TemplateDot).Cookie", (*vweb.TemplateDot).Cookie, execmTemplateDotCookie),
		I.Func("(*TemplateDot).Swap", (*vweb.TemplateDot).Swap, execmTemplateDotSwap),
		I.Func("(*TemplateDot).Free", (*vweb.TemplateDot).Free, execmTemplateDotFree),
		I.Func("(*TemplateDot).Context", (*vweb.TemplateDot).Context, execmTemplateDotContext),
		I.Func("(*TemplateDot).WithContext", (*vweb.TemplateDot).WithContext, execmTemplateDotWithContext),
		I.Func("(TemplateDoter).Cookie", (vweb.TemplateDoter).Cookie, execiTemplateDoterCookie),
		I.Func("(TemplateDoter).Global", (vweb.TemplateDoter).Global, execiTemplateDoterGlobal),
		I.Func("(TemplateDoter).Header", (vweb.TemplateDoter).Header, execiTemplateDoterHeader),
		I.Func("(TemplateDoter).Request", (vweb.TemplateDoter).Request, execiTemplateDoterRequest),
		I.Func("(TemplateDoter).RequestLimitSize", (vweb.TemplateDoter).RequestLimitSize, execiTemplateDoterRequestLimitSize),
		I.Func("(TemplateDoter).Response", (vweb.TemplateDoter).Response, execiTemplateDoterResponse),
		I.Func("(TemplateDoter).ResponseWriter", (vweb.TemplateDoter).ResponseWriter, execiTemplateDoterResponseWriter),
		I.Func("(TemplateDoter).RootDir", (vweb.TemplateDoter).RootDir, execiTemplateDoterRootDir),
		I.Func("(TemplateDoter).Session", (vweb.TemplateDoter).Session, execiTemplateDoterSession),
		I.Func("(TemplateDoter).Swap", (vweb.TemplateDoter).Swap, execiTemplateDoterSwap),
	)
	I.RegisterFuncvs(
		I.Funcv("DepthField", vweb.DepthField, execDepthField),
		I.Funcv("ExecFunc", vweb.ExecFunc, execExecFunc),
		I.Funcv("(*ExitCall).Defer", (*vweb.ExitCall).Defer, execmExitCallDefer),
		I.Funcv("(*ServerHandlerDynamicTemplateExtend).Call", (*vweb.ServerHandlerDynamicTemplateExtend).Call, execmServerHandlerDynamicTemplateExtendCall),
		I.Funcv("(*Session).Defer", (*vweb.Session).Defer, execmSessionDefer),
		I.Funcv("(Sessioner).Defer", (vweb.Sessioner).Defer, execiSessionerDefer),
		I.Funcv("(*TemplateDot).Defer", (*vweb.TemplateDot).Defer, execmTemplateDotDefer),
		I.Funcv("(TemplateDoter).Defer", (vweb.TemplateDoter).Defer, execiTemplateDoterDefer),
	)
	I.RegisterVars(
		I.Var("DefaultSitePool", &vweb.DefaultSitePool),
		I.Var("TemplateFunc", &vweb.TemplateFunc),
	)
	I.RegisterTypes(
		I.Type("Cookie", reflect.TypeOf((*vweb.Cookie)(nil)).Elem()),
		I.Type("Cookier", reflect.TypeOf((*vweb.Cookier)(nil)).Elem()),
		I.Type("DotContexter", reflect.TypeOf((*vweb.DotContexter)(nil)).Elem()),
		I.Type("DynamicTemplateFunc", reflect.TypeOf((*vweb.DynamicTemplateFunc)(nil)).Elem()),
		I.Type("DynamicTemplater", reflect.TypeOf((*vweb.DynamicTemplater)(nil)).Elem()),
		I.Type("ExitCall", reflect.TypeOf((*vweb.ExitCall)(nil)).Elem()),
		I.Type("Forward", reflect.TypeOf((*vweb.Forward)(nil)).Elem()),
		I.Type("Globaler", reflect.TypeOf((*vweb.Globaler)(nil)).Elem()),
		I.Type("PluginHTTP", reflect.TypeOf((*vweb.PluginHTTP)(nil)).Elem()),
		I.Type("PluginHTTPClient", reflect.TypeOf((*vweb.PluginHTTPClient)(nil)).Elem()),
		I.Type("PluginRPC", reflect.TypeOf((*vweb.PluginRPC)(nil)).Elem()),
		I.Type("PluginRPCClient", reflect.TypeOf((*vweb.PluginRPCClient)(nil)).Elem()),
		I.Type("PluginType", reflect.TypeOf((*vweb.PluginType)(nil)).Elem()),
		I.Type("Responser", reflect.TypeOf((*vweb.Responser)(nil)).Elem()),
		I.Type("Route", reflect.TypeOf((*vweb.Route)(nil)).Elem()),
		I.Type("ServerHandlerDynamic", reflect.TypeOf((*vweb.ServerHandlerDynamic)(nil)).Elem()),
		I.Type("ServerHandlerDynamicTemplateExtend", reflect.TypeOf((*vweb.ServerHandlerDynamicTemplateExtend)(nil)).Elem()),
		I.Type("ServerHandlerStatic", reflect.TypeOf((*vweb.ServerHandlerStatic)(nil)).Elem()),
		I.Type("Session", reflect.TypeOf((*vweb.Session)(nil)).Elem()),
		I.Type("Sessioner", reflect.TypeOf((*vweb.Sessioner)(nil)).Elem()),
		I.Type("Sessions", reflect.TypeOf((*vweb.Sessions)(nil)).Elem()),
		I.Type("Site", reflect.TypeOf((*vweb.Site)(nil)).Elem()),
		I.Type("SiteMan", reflect.TypeOf((*vweb.SiteMan)(nil)).Elem()),
		I.Type("SitePool", reflect.TypeOf((*vweb.SitePool)(nil)).Elem()),
		I.Type("TemplateDot", reflect.TypeOf((*vweb.TemplateDot)(nil)).Elem()),
		I.Type("TemplateDoter", reflect.TypeOf((*vweb.TemplateDoter)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("PluginTypeHTTP", qspec.Int, vweb.PluginTypeHTTP),
		I.Const("PluginTypeRPC", qspec.Int, vweb.PluginTypeRPC),
		I.Const("Version", qspec.String, vweb.Version),
	)
}
