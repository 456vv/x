// Package smtp provide Go+ "github.com/456vv/x/smtp" package, as "github.com/456vv/x/smtp" package in Go.
package smtp

import (
	reflect "reflect"

	smtp "github.com/456vv/x/smtp"
	gop "github.com/goplus/gop"
)

func execPlainAuth(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := smtp.PlainAuth(args[0].(string), args[1].(string), args[2].(string), args[3].(string))
	p.Ret(4, ret0)
}

func execmSmtpOpenConfig(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*smtp.Smtp).OpenConfig(args[1].(string))
	p.Ret(2, ret0)
}

func execmSmtpClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*smtp.Smtp).Close()
	p.Ret(1, ret0)
}

func execmSmtpSend(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := args[0].(*smtp.Smtp).Send(args[1].([]string), args[2].(string), args[3].(string))
	p.Ret(4, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("github.com/456vv/x/smtp")

func init() {
	I.RegisterFuncs(
		I.Func("PlainAuth", smtp.PlainAuth, execPlainAuth),
		I.Func("(*Smtp).OpenConfig", (*smtp.Smtp).OpenConfig, execmSmtpOpenConfig),
		I.Func("(*Smtp).Close", (*smtp.Smtp).Close, execmSmtpClose),
		I.Func("(*Smtp).Send", (*smtp.Smtp).Send, execmSmtpSend),
	)
	I.RegisterTypes(
		I.Type("Info", reflect.TypeOf((*smtp.Info)(nil)).Elem()),
		I.Type("Smtp", reflect.TypeOf((*smtp.Smtp)(nil)).Elem()),
	)
}
