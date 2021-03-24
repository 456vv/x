// Package viot provide Go+ "github.com/456vv/viot/v2" package, as "github.com/456vv/viot/v2" package in Go.
package viot

import (
	context "context"
	io "io"
	net "net"
	http "net/http"
	reflect "reflect"
	template "text/template"
	time "time"

	viot "github.com/456vv/viot/v2"
	vweb "github.com/456vv/vweb/v2"
	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execmClientGet(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*viot.Client).Get(args[1].(string), args[2].(viot.Header))
	p.Ret(3, ret0, ret1)
}

func toType0(v interface{}) context.Context {
	if v == nil {
		return nil
	}
	return v.(context.Context)
}

func execmClientGetCtx(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0, ret1 := args[0].(*viot.Client).GetCtx(toType0(args[1]), args[2].(string), args[3].(viot.Header))
	p.Ret(4, ret0, ret1)
}

func execmClientDo(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*viot.Client).Do(args[1].(*viot.Request))
	p.Ret(2, ret0, ret1)
}

func execmClientDoCtx(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*viot.Client).DoCtx(toType0(args[1]), args[2].(*viot.Request))
	p.Ret(3, ret0, ret1)
}

func execmClientPost(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0, ret1 := args[0].(*viot.Client).Post(args[1].(string), args[2].(viot.Header), args[3])
	p.Ret(4, ret0, ret1)
}

func execmClientPostCtx(_ int, p *gop.Context) {
	args := p.GetArgs(5)
	ret0, ret1 := args[0].(*viot.Client).PostCtx(toType0(args[1]), args[2].(string), args[3].(viot.Header), args[4])
	p.Ret(5, ret0, ret1)
}

func execiCloseNotifierCloseNotify(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(viot.CloseNotifier).CloseNotify()
	p.Ret(1, ret0)
}

func execmConnStateString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(viot.ConnState).String()
	p.Ret(1, ret0)
}

func execiDotContexterContext(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(viot.DotContexter).Context()
	p.Ret(1, ret0)
}

func execiDotContexterWithContext(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(viot.DotContexter).WithContext(toType0(args[1]))
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
	ret0 := args[0].(viot.DynamicTemplater).Execute(toType1(args[1]), args[2])
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
	ret0 := args[0].(viot.DynamicTemplater).Parse(toType2(args[1]))
	p.Ret(2, ret0)
}

func execiDynamicTemplaterSetPath(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(viot.DynamicTemplater).SetPath(args[1].(string), args[2].(string))
	p.PopN(3)
}

func toType3(v interface{}) viot.ResponseWriter {
	if v == nil {
		return nil
	}
	return v.(viot.ResponseWriter)
}

func execError(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	viot.Error(toType3(args[0]), args[1].(string), args[2].(int))
	p.PopN(3)
}

func execExtendTemplatePackage(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	viot.ExtendTemplatePackage(args[0].(string), args[1].(template.FuncMap))
	p.PopN(2)
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

func execiHandlerServeIOT(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(viot.Handler).ServeIOT(toType3(args[1]), args[2].(*viot.Request))
	p.PopN(3)
}

func execmHandlerFuncServeIOT(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(viot.HandlerFunc).ServeIOT(toType3(args[1]), args[2].(*viot.Request))
	p.PopN(3)
}

func execmHeaderSet(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(viot.Header).Set(args[1].(string), args[2].(string))
	p.PopN(3)
}

func execmHeaderGet(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(viot.Header).Get(args[1].(string))
	p.Ret(2, ret0)
}

func execmHeaderDel(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(viot.Header).Del(args[1].(string))
	p.PopN(2)
}

func execmHeaderClone(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(viot.Header).Clone()
	p.Ret(1, ret0)
}

func execiHijackerHijack(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2 := args[0].(viot.Hijacker).Hijack()
	p.Ret(1, ret0, ret1, ret2)
}

func execmHomePoolName(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*viot.Home).PoolName()
	p.Ret(1, ret0)
}

func execmHomeManAdd(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*viot.HomeMan).Add(args[1].(string), args[2].(*viot.Home))
	p.PopN(3)
}

func execmHomeManGet(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*viot.HomeMan).Get(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmHomeManRange(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*viot.HomeMan).Range(args[1].(func(host string, home *viot.Home) bool))
	p.PopN(2)
}

func execmHomePoolNewHome(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*viot.HomePool).NewHome(args[1].(string))
	p.Ret(2, ret0)
}

func execmHomePoolDelHome(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*viot.HomePool).DelHome(args[1].(string))
	p.PopN(2)
}

func execmHomePoolRangeHome(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*viot.HomePool).RangeHome(args[1].(func(name string, home *viot.Home) bool))
	p.PopN(2)
}

func execmHomePoolSetRecoverSession(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*viot.HomePool).SetRecoverSession(args[1].(time.Duration))
	p.PopN(2)
}

func execmHomePoolStart(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*viot.HomePool).Start()
	p.Ret(1, ret0)
}

func execmHomePoolClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*viot.HomePool).Close()
	p.Ret(1, ret0)
}

func execiLauncherLaunch(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(viot.Launcher).Launch()
	p.Ret(1, ret0)
}

func execNewHomePool(_ int, p *gop.Context) {
	ret0 := viot.NewHomePool()
	p.Ret(0, ret0)
}

func execNonce(_ int, p *gop.Context) {
	ret0, ret1 := viot.Nonce()
	p.Ret(0, ret0, ret1)
}

func execParseIOTVersion(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2 := viot.ParseIOTVersion(args[0].(string))
	p.Ret(1, ret0, ret1, ret2)
}

func execReadRequest(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := viot.ReadRequest(toType2(args[0]))
	p.Ret(1, ret0, ret1)
}

func execReadResponse(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := viot.ReadResponse(toType2(args[0]), args[1].(*viot.Request))
	p.Ret(2, ret0, ret1)
}

func execmRequestGetNonce(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*viot.Request).GetNonce()
	p.Ret(1, ret0)
}

func execmRequestGetBody(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*viot.Request).GetBody(args[1])
	p.Ret(2, ret0)
}

func execmRequestSetBody(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*viot.Request).SetBody(args[1])
	p.Ret(2, ret0)
}

func execmRequestProtoAtLeast(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*viot.Request).ProtoAtLeast(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmRequestContext(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*viot.Request).Context()
	p.Ret(1, ret0)
}

func execmRequestWithContext(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*viot.Request).WithContext(toType0(args[1]))
	p.Ret(2, ret0)
}

func execmRequestGetBasicAuth(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2 := args[0].(*viot.Request).GetBasicAuth()
	p.Ret(1, ret0, ret1, ret2)
}

func execmRequestSetBasicAuth(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*viot.Request).SetBasicAuth(args[1].(string), args[2].(string))
	p.PopN(3)
}

func execmRequestGetTokenAuth(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*viot.Request).GetTokenAuth()
	p.Ret(1, ret0, ret1)
}

func execmRequestSetTokenAuth(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*viot.Request).SetTokenAuth(args[1].(string))
	p.PopN(2)
}

func execmRequestRequestConfig(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*viot.Request).RequestConfig(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmRequestConfigSetBody(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*viot.RequestConfig).SetBody(args[1])
	p.PopN(2)
}

func execmRequestConfigGetBody(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*viot.RequestConfig).GetBody()
	p.Ret(1, ret0)
}

func execmRequestConfigMarshal(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*viot.RequestConfig).Marshal()
	p.Ret(1, ret0, ret1)
}

func execmRequestConfigUnmarshal(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*viot.RequestConfig).Unmarshal(args[1].([]byte))
	p.Ret(2, ret0)
}

func execmResponseSetNonce(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*viot.Response).SetNonce(args[1].(string))
	p.PopN(2)
}

func execmResponseWriteTo(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*viot.Response).WriteTo(toType3(args[1]))
	p.PopN(2)
}

func execmResponseWrite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*viot.Response).Write(toType1(args[1]))
	p.Ret(2, ret0)
}

func execmResponseResponseConfig(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*viot.Response).ResponseConfig(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execiResponseWriterHeader(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(viot.ResponseWriter).Header()
	p.Ret(1, ret0)
}

func execiResponseWriterSetBody(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(viot.ResponseWriter).SetBody(args[1])
	p.Ret(2, ret0)
}

func execiResponseWriterStatus(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(viot.ResponseWriter).Status(args[1].(int))
	p.PopN(2)
}

func execiRoundTripperRoundTrip(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(viot.RoundTripper).RoundTrip(args[1].(*viot.Request))
	p.Ret(2, ret0, ret1)
}

func execiRoundTripperRoundTripContext(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(viot.RoundTripper).RoundTripContext(toType0(args[1]), args[2].(*viot.Request))
	p.Ret(3, ret0, ret1)
}

func execmRouteHandleFunc(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*viot.Route).HandleFunc(args[1].(string), args[2].(func(w viot.ResponseWriter, r *viot.Request)))
	p.PopN(3)
}

func execmRouteServeIOT(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*viot.Route).ServeIOT(toType3(args[1]), args[2].(*viot.Request))
	p.PopN(3)
}

func execmServerListenAndServe(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*viot.Server).ListenAndServe()
	p.Ret(1, ret0)
}

func toType4(v interface{}) net.Listener {
	if v == nil {
		return nil
	}
	return v.(net.Listener)
}

func execmServerServe(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*viot.Server).Serve(toType4(args[1]))
	p.Ret(2, ret0)
}

func execmServerClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*viot.Server).Close()
	p.Ret(1, ret0)
}

func execmServerShutdown(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*viot.Server).Shutdown(toType0(args[1]))
	p.Ret(2, ret0)
}

func execmServerRegisterOnShutdown(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*viot.Server).RegisterOnShutdown(args[1].(func()))
	p.PopN(2)
}

func execmServerSetKeepAlivesEnabled(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*viot.Server).SetKeepAlivesEnabled(args[1].(bool))
	p.PopN(2)
}

func execmServerHandlerDynamicServeIOT(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*viot.ServerHandlerDynamic).ServeIOT(toType3(args[1]), args[2].(*viot.Request))
	p.PopN(3)
}

func execmServerHandlerDynamicParseText(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*viot.ServerHandlerDynamic).ParseText(args[1].(string), args[2].(string))
	p.Ret(3, ret0)
}

func execmServerHandlerDynamicParseFile(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*viot.ServerHandlerDynamic).ParseFile(args[1].(string))
	p.Ret(2, ret0)
}

func execmServerHandlerDynamicParse(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*viot.ServerHandlerDynamic).Parse(toType2(args[1]))
	p.Ret(2, ret0)
}

func execmServerHandlerDynamicExecute(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*viot.ServerHandlerDynamic).Execute(toType1(args[1]), args[2])
	p.Ret(3, ret0)
}

func execmServerHandlerDynamicTemplateExtendNewFunc(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*viot.ServerHandlerDynamicTemplateExtend).NewFunc(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmServerHandlerDynamicTemplateExtendCall(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*viot.ServerHandlerDynamicTemplateExtend).Call(args[1].(func([]reflect.Value) []reflect.Value), args[2:]...)
	p.Ret(arity, ret0)
}

func execmServerHandlerDynamicTemplateExtendExecuteTemplate(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := args[0].(*viot.ServerHandlerDynamicTemplateExtend).ExecuteTemplate(toType1(args[1]), args[2].(string), args[3])
	p.Ret(4, ret0)
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

func toType6(v interface{}) http.ResponseWriter {
	if v == nil {
		return nil
	}
	return v.(http.ResponseWriter)
}

func execmSessionsSession(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*vweb.Sessions).Session(toType6(args[1]), args[2].(*http.Request))
	p.Ret(3, ret0)
}

func execmTemplateDotRootDir(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*viot.TemplateDot).RootDir(args[1].(string))
	p.Ret(2, ret0)
}

func execmTemplateDotRequest(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*viot.TemplateDot).Request()
	p.Ret(1, ret0)
}

func execmTemplateDotHeader(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*viot.TemplateDot).Header()
	p.Ret(1, ret0)
}

func execmTemplateDotResponseWriter(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*viot.TemplateDot).ResponseWriter()
	p.Ret(1, ret0)
}

func execmTemplateDotLaunch(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*viot.TemplateDot).Launch()
	p.Ret(1, ret0)
}

func execmTemplateDotHijack(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2 := args[0].(*viot.TemplateDot).Hijack()
	p.Ret(1, ret0, ret1, ret2)
}

func execmTemplateDotSession(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*viot.TemplateDot).Session(args[1].(string))
	p.Ret(2, ret0)
}

func execmTemplateDotGlobal(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*viot.TemplateDot).Global()
	p.Ret(1, ret0)
}

func execmTemplateDotSwap(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*viot.TemplateDot).Swap()
	p.Ret(1, ret0)
}

func execmTemplateDotDefer(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*viot.TemplateDot).Defer(args[1], args[2:]...)
	p.Ret(arity, ret0)
}

func execmTemplateDotFree(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*viot.TemplateDot).Free()
	p.PopN(1)
}

func execmTemplateDotContext(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*viot.TemplateDot).Context()
	p.Ret(1, ret0)
}

func execmTemplateDotWithContext(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*viot.TemplateDot).WithContext(toType0(args[1]))
	p.PopN(2)
}

func execiTemplateDoterDefer(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(viot.TemplateDoter).Defer(args[1], args[2:]...)
	p.Ret(arity, ret0)
}

func execiTemplateDoterGlobal(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(viot.TemplateDoter).Global()
	p.Ret(1, ret0)
}

func execiTemplateDoterHeader(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(viot.TemplateDoter).Header()
	p.Ret(1, ret0)
}

func execiTemplateDoterHijack(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2 := args[0].(viot.TemplateDoter).Hijack()
	p.Ret(1, ret0, ret1, ret2)
}

func execiTemplateDoterLaunch(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(viot.TemplateDoter).Launch()
	p.Ret(1, ret0)
}

func execiTemplateDoterRequest(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(viot.TemplateDoter).Request()
	p.Ret(1, ret0)
}

func execiTemplateDoterResponseWriter(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(viot.TemplateDoter).ResponseWriter()
	p.Ret(1, ret0)
}

func execiTemplateDoterRootDir(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(viot.TemplateDoter).RootDir(args[1].(string))
	p.Ret(2, ret0)
}

func execiTemplateDoterSession(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(viot.TemplateDoter).Session(args[1].(string))
	p.Ret(2, ret0)
}

func execiTemplateDoterSwap(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(viot.TemplateDoter).Swap()
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("github.com/456vv/viot/v2")

func init() {
	I.RegisterFuncs(
		I.Func("(*Client).Get", (*viot.Client).Get, execmClientGet),
		I.Func("(*Client).GetCtx", (*viot.Client).GetCtx, execmClientGetCtx),
		I.Func("(*Client).Do", (*viot.Client).Do, execmClientDo),
		I.Func("(*Client).DoCtx", (*viot.Client).DoCtx, execmClientDoCtx),
		I.Func("(*Client).Post", (*viot.Client).Post, execmClientPost),
		I.Func("(*Client).PostCtx", (*viot.Client).PostCtx, execmClientPostCtx),
		I.Func("(CloseNotifier).CloseNotify", (viot.CloseNotifier).CloseNotify, execiCloseNotifierCloseNotify),
		I.Func("(ConnState).String", (viot.ConnState).String, execmConnStateString),
		I.Func("(DotContexter).Context", (viot.DotContexter).Context, execiDotContexterContext),
		I.Func("(DotContexter).WithContext", (viot.DotContexter).WithContext, execiDotContexterWithContext),
		I.Func("(DynamicTemplater).Execute", (viot.DynamicTemplater).Execute, execiDynamicTemplaterExecute),
		I.Func("(DynamicTemplater).Parse", (viot.DynamicTemplater).Parse, execiDynamicTemplaterParse),
		I.Func("(DynamicTemplater).SetPath", (viot.DynamicTemplater).SetPath, execiDynamicTemplaterSetPath),
		I.Func("Error", viot.Error, execError),
		I.Func("ExtendTemplatePackage", viot.ExtendTemplatePackage, execExtendTemplatePackage),
		I.Func("(Globaler).Del", (vweb.Globaler).Del, execiGlobalerDel),
		I.Func("(Globaler).Get", (vweb.Globaler).Get, execiGlobalerGet),
		I.Func("(Globaler).Has", (vweb.Globaler).Has, execiGlobalerHas),
		I.Func("(Globaler).Reset", (vweb.Globaler).Reset, execiGlobalerReset),
		I.Func("(Globaler).Set", (vweb.Globaler).Set, execiGlobalerSet),
		I.Func("(Globaler).SetExpired", (vweb.Globaler).SetExpired, execiGlobalerSetExpired),
		I.Func("(Globaler).SetExpiredCall", (vweb.Globaler).SetExpiredCall, execiGlobalerSetExpiredCall),
		I.Func("(Handler).ServeIOT", (viot.Handler).ServeIOT, execiHandlerServeIOT),
		I.Func("(HandlerFunc).ServeIOT", (viot.HandlerFunc).ServeIOT, execmHandlerFuncServeIOT),
		I.Func("(Header).Set", (viot.Header).Set, execmHeaderSet),
		I.Func("(Header).Get", (viot.Header).Get, execmHeaderGet),
		I.Func("(Header).Del", (viot.Header).Del, execmHeaderDel),
		I.Func("(Header).Clone", (viot.Header).Clone, execmHeaderClone),
		I.Func("(Hijacker).Hijack", (viot.Hijacker).Hijack, execiHijackerHijack),
		I.Func("(*Home).PoolName", (*viot.Home).PoolName, execmHomePoolName),
		I.Func("(*HomeMan).Add", (*viot.HomeMan).Add, execmHomeManAdd),
		I.Func("(*HomeMan).Get", (*viot.HomeMan).Get, execmHomeManGet),
		I.Func("(*HomeMan).Range", (*viot.HomeMan).Range, execmHomeManRange),
		I.Func("(*HomePool).NewHome", (*viot.HomePool).NewHome, execmHomePoolNewHome),
		I.Func("(*HomePool).DelHome", (*viot.HomePool).DelHome, execmHomePoolDelHome),
		I.Func("(*HomePool).RangeHome", (*viot.HomePool).RangeHome, execmHomePoolRangeHome),
		I.Func("(*HomePool).SetRecoverSession", (*viot.HomePool).SetRecoverSession, execmHomePoolSetRecoverSession),
		I.Func("(*HomePool).Start", (*viot.HomePool).Start, execmHomePoolStart),
		I.Func("(*HomePool).Close", (*viot.HomePool).Close, execmHomePoolClose),
		I.Func("(Launcher).Launch", (viot.Launcher).Launch, execiLauncherLaunch),
		I.Func("NewHomePool", viot.NewHomePool, execNewHomePool),
		I.Func("Nonce", viot.Nonce, execNonce),
		I.Func("ParseIOTVersion", viot.ParseIOTVersion, execParseIOTVersion),
		I.Func("ReadRequest", viot.ReadRequest, execReadRequest),
		I.Func("ReadResponse", viot.ReadResponse, execReadResponse),
		I.Func("(*Request).GetNonce", (*viot.Request).GetNonce, execmRequestGetNonce),
		I.Func("(*Request).GetBody", (*viot.Request).GetBody, execmRequestGetBody),
		I.Func("(*Request).SetBody", (*viot.Request).SetBody, execmRequestSetBody),
		I.Func("(*Request).ProtoAtLeast", (*viot.Request).ProtoAtLeast, execmRequestProtoAtLeast),
		I.Func("(*Request).Context", (*viot.Request).Context, execmRequestContext),
		I.Func("(*Request).WithContext", (*viot.Request).WithContext, execmRequestWithContext),
		I.Func("(*Request).GetBasicAuth", (*viot.Request).GetBasicAuth, execmRequestGetBasicAuth),
		I.Func("(*Request).SetBasicAuth", (*viot.Request).SetBasicAuth, execmRequestSetBasicAuth),
		I.Func("(*Request).GetTokenAuth", (*viot.Request).GetTokenAuth, execmRequestGetTokenAuth),
		I.Func("(*Request).SetTokenAuth", (*viot.Request).SetTokenAuth, execmRequestSetTokenAuth),
		I.Func("(*Request).RequestConfig", (*viot.Request).RequestConfig, execmRequestRequestConfig),
		I.Func("(*RequestConfig).SetBody", (*viot.RequestConfig).SetBody, execmRequestConfigSetBody),
		I.Func("(*RequestConfig).GetBody", (*viot.RequestConfig).GetBody, execmRequestConfigGetBody),
		I.Func("(*RequestConfig).Marshal", (*viot.RequestConfig).Marshal, execmRequestConfigMarshal),
		I.Func("(*RequestConfig).Unmarshal", (*viot.RequestConfig).Unmarshal, execmRequestConfigUnmarshal),
		I.Func("(*Response).SetNonce", (*viot.Response).SetNonce, execmResponseSetNonce),
		I.Func("(*Response).WriteTo", (*viot.Response).WriteTo, execmResponseWriteTo),
		I.Func("(*Response).Write", (*viot.Response).Write, execmResponseWrite),
		I.Func("(*Response).ResponseConfig", (*viot.Response).ResponseConfig, execmResponseResponseConfig),
		I.Func("(ResponseWriter).Header", (viot.ResponseWriter).Header, execiResponseWriterHeader),
		I.Func("(ResponseWriter).SetBody", (viot.ResponseWriter).SetBody, execiResponseWriterSetBody),
		I.Func("(ResponseWriter).Status", (viot.ResponseWriter).Status, execiResponseWriterStatus),
		I.Func("(RoundTripper).RoundTrip", (viot.RoundTripper).RoundTrip, execiRoundTripperRoundTrip),
		I.Func("(RoundTripper).RoundTripContext", (viot.RoundTripper).RoundTripContext, execiRoundTripperRoundTripContext),
		I.Func("(*Route).HandleFunc", (*viot.Route).HandleFunc, execmRouteHandleFunc),
		I.Func("(*Route).ServeIOT", (*viot.Route).ServeIOT, execmRouteServeIOT),
		I.Func("(*Server).ListenAndServe", (*viot.Server).ListenAndServe, execmServerListenAndServe),
		I.Func("(*Server).Serve", (*viot.Server).Serve, execmServerServe),
		I.Func("(*Server).Close", (*viot.Server).Close, execmServerClose),
		I.Func("(*Server).Shutdown", (*viot.Server).Shutdown, execmServerShutdown),
		I.Func("(*Server).RegisterOnShutdown", (*viot.Server).RegisterOnShutdown, execmServerRegisterOnShutdown),
		I.Func("(*Server).SetKeepAlivesEnabled", (*viot.Server).SetKeepAlivesEnabled, execmServerSetKeepAlivesEnabled),
		I.Func("(*ServerHandlerDynamic).ServeIOT", (*viot.ServerHandlerDynamic).ServeIOT, execmServerHandlerDynamicServeIOT),
		I.Func("(*ServerHandlerDynamic).ParseText", (*viot.ServerHandlerDynamic).ParseText, execmServerHandlerDynamicParseText),
		I.Func("(*ServerHandlerDynamic).ParseFile", (*viot.ServerHandlerDynamic).ParseFile, execmServerHandlerDynamicParseFile),
		I.Func("(*ServerHandlerDynamic).Parse", (*viot.ServerHandlerDynamic).Parse, execmServerHandlerDynamicParse),
		I.Func("(*ServerHandlerDynamic).Execute", (*viot.ServerHandlerDynamic).Execute, execmServerHandlerDynamicExecute),
		I.Func("(*ServerHandlerDynamicTemplateExtend).NewFunc", (*viot.ServerHandlerDynamicTemplateExtend).NewFunc, execmServerHandlerDynamicTemplateExtendNewFunc),
		I.Func("(*ServerHandlerDynamicTemplateExtend).ExecuteTemplate", (*viot.ServerHandlerDynamicTemplateExtend).ExecuteTemplate, execmServerHandlerDynamicTemplateExtendExecuteTemplate),
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
		I.Func("(*TemplateDot).RootDir", (*viot.TemplateDot).RootDir, execmTemplateDotRootDir),
		I.Func("(*TemplateDot).Request", (*viot.TemplateDot).Request, execmTemplateDotRequest),
		I.Func("(*TemplateDot).Header", (*viot.TemplateDot).Header, execmTemplateDotHeader),
		I.Func("(*TemplateDot).ResponseWriter", (*viot.TemplateDot).ResponseWriter, execmTemplateDotResponseWriter),
		I.Func("(*TemplateDot).Launch", (*viot.TemplateDot).Launch, execmTemplateDotLaunch),
		I.Func("(*TemplateDot).Hijack", (*viot.TemplateDot).Hijack, execmTemplateDotHijack),
		I.Func("(*TemplateDot).Session", (*viot.TemplateDot).Session, execmTemplateDotSession),
		I.Func("(*TemplateDot).Global", (*viot.TemplateDot).Global, execmTemplateDotGlobal),
		I.Func("(*TemplateDot).Swap", (*viot.TemplateDot).Swap, execmTemplateDotSwap),
		I.Func("(*TemplateDot).Free", (*viot.TemplateDot).Free, execmTemplateDotFree),
		I.Func("(*TemplateDot).Context", (*viot.TemplateDot).Context, execmTemplateDotContext),
		I.Func("(*TemplateDot).WithContext", (*viot.TemplateDot).WithContext, execmTemplateDotWithContext),
		I.Func("(TemplateDoter).Global", (viot.TemplateDoter).Global, execiTemplateDoterGlobal),
		I.Func("(TemplateDoter).Header", (viot.TemplateDoter).Header, execiTemplateDoterHeader),
		I.Func("(TemplateDoter).Hijack", (viot.TemplateDoter).Hijack, execiTemplateDoterHijack),
		I.Func("(TemplateDoter).Launch", (viot.TemplateDoter).Launch, execiTemplateDoterLaunch),
		I.Func("(TemplateDoter).Request", (viot.TemplateDoter).Request, execiTemplateDoterRequest),
		I.Func("(TemplateDoter).ResponseWriter", (viot.TemplateDoter).ResponseWriter, execiTemplateDoterResponseWriter),
		I.Func("(TemplateDoter).RootDir", (viot.TemplateDoter).RootDir, execiTemplateDoterRootDir),
		I.Func("(TemplateDoter).Session", (viot.TemplateDoter).Session, execiTemplateDoterSession),
		I.Func("(TemplateDoter).Swap", (viot.TemplateDoter).Swap, execiTemplateDoterSwap),
	)
	I.RegisterFuncvs(
		I.Funcv("(*ServerHandlerDynamicTemplateExtend).Call", (*viot.ServerHandlerDynamicTemplateExtend).Call, execmServerHandlerDynamicTemplateExtendCall),
		I.Funcv("(*Session).Defer", (*vweb.Session).Defer, execmSessionDefer),
		I.Funcv("(Sessioner).Defer", (vweb.Sessioner).Defer, execiSessionerDefer),
		I.Funcv("(*TemplateDot).Defer", (*viot.TemplateDot).Defer, execmTemplateDotDefer),
		I.Funcv("(TemplateDoter).Defer", (viot.TemplateDoter).Defer, execiTemplateDoterDefer),
	)
	I.RegisterVars(
		I.Var("DefaultHomePool", &viot.DefaultHomePool),
		I.Var("ErrAbortHandler", &viot.ErrAbortHandler),
		I.Var("ErrBodyNotAllowed", &viot.ErrBodyNotAllowed),
		I.Var("ErrConnClose", &viot.ErrConnClose),
		I.Var("ErrDoned", &viot.ErrDoned),
		I.Var("ErrGetBodyed", &viot.ErrGetBodyed),
		I.Var("ErrHijacked", &viot.ErrHijacked),
		I.Var("ErrHomeInvalid", &viot.ErrHomeInvalid),
		I.Var("ErrLaunched", &viot.ErrLaunched),
		I.Var("ErrMethodInvalid", &viot.ErrMethodInvalid),
		I.Var("ErrProtoInvalid", &viot.ErrProtoInvalid),
		I.Var("ErrReqUnavailable", &viot.ErrReqUnavailable),
		I.Var("ErrServerClosed", &viot.ErrServerClosed),
		I.Var("ErrURIInvalid", &viot.ErrURIInvalid),
		I.Var("LocalAddrContextKey", &viot.LocalAddrContextKey),
		I.Var("ServerContextKey", &viot.ServerContextKey),
		I.Var("TemplateFunc", &viot.TemplateFunc),
	)
	I.RegisterTypes(
		I.Type("Client", reflect.TypeOf((*viot.Client)(nil)).Elem()),
		I.Type("CloseNotifier", reflect.TypeOf((*viot.CloseNotifier)(nil)).Elem()),
		I.Type("ConnState", reflect.TypeOf((*viot.ConnState)(nil)).Elem()),
		I.Type("DotContexter", reflect.TypeOf((*viot.DotContexter)(nil)).Elem()),
		I.Type("DynamicTemplateFunc", reflect.TypeOf((*viot.DynamicTemplateFunc)(nil)).Elem()),
		I.Type("DynamicTemplater", reflect.TypeOf((*viot.DynamicTemplater)(nil)).Elem()),
		I.Type("Globaler", reflect.TypeOf((*viot.Globaler)(nil)).Elem()),
		I.Type("Handler", reflect.TypeOf((*viot.Handler)(nil)).Elem()),
		I.Type("HandlerFunc", reflect.TypeOf((*viot.HandlerFunc)(nil)).Elem()),
		I.Type("Header", reflect.TypeOf((*viot.Header)(nil)).Elem()),
		I.Type("Hijacker", reflect.TypeOf((*viot.Hijacker)(nil)).Elem()),
		I.Type("Home", reflect.TypeOf((*viot.Home)(nil)).Elem()),
		I.Type("HomeMan", reflect.TypeOf((*viot.HomeMan)(nil)).Elem()),
		I.Type("HomePool", reflect.TypeOf((*viot.HomePool)(nil)).Elem()),
		I.Type("Launcher", reflect.TypeOf((*viot.Launcher)(nil)).Elem()),
		I.Type("Request", reflect.TypeOf((*viot.Request)(nil)).Elem()),
		I.Type("RequestConfig", reflect.TypeOf((*viot.RequestConfig)(nil)).Elem()),
		I.Type("Response", reflect.TypeOf((*viot.Response)(nil)).Elem()),
		I.Type("ResponseConfig", reflect.TypeOf((*viot.ResponseConfig)(nil)).Elem()),
		I.Type("ResponseWriter", reflect.TypeOf((*viot.ResponseWriter)(nil)).Elem()),
		I.Type("RoundTripper", reflect.TypeOf((*viot.RoundTripper)(nil)).Elem()),
		I.Type("Route", reflect.TypeOf((*viot.Route)(nil)).Elem()),
		I.Type("Server", reflect.TypeOf((*viot.Server)(nil)).Elem()),
		I.Type("ServerHandlerDynamic", reflect.TypeOf((*viot.ServerHandlerDynamic)(nil)).Elem()),
		I.Type("ServerHandlerDynamicTemplateExtend", reflect.TypeOf((*viot.ServerHandlerDynamicTemplateExtend)(nil)).Elem()),
		I.Type("Session", reflect.TypeOf((*viot.Session)(nil)).Elem()),
		I.Type("Sessioner", reflect.TypeOf((*viot.Sessioner)(nil)).Elem()),
		I.Type("Sessions", reflect.TypeOf((*viot.Sessions)(nil)).Elem()),
		I.Type("TemplateDot", reflect.TypeOf((*viot.TemplateDot)(nil)).Elem()),
		I.Type("TemplateDoter", reflect.TypeOf((*viot.TemplateDoter)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("DefaultLineBytes", qspec.ConstUnboundInt, viot.DefaultLineBytes),
		I.Const("StateActive", qspec.Int, viot.StateActive),
		I.Const("StateClosed", qspec.Int, viot.StateClosed),
		I.Const("StateHijacked", qspec.Int, viot.StateHijacked),
		I.Const("StateIdle", qspec.Int, viot.StateIdle),
		I.Const("StateNew", qspec.Int, viot.StateNew),
	)
}
