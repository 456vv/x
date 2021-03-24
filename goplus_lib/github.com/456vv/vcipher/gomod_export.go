// Package vcipher provide Go+ "github.com/456vv/vcipher" package, as "github.com/456vv/vcipher" package in Go.
package vcipher

import (
	cipher "crypto/cipher"
	reflect "reflect"

	vcipher "github.com/456vv/vcipher"
	gop "github.com/goplus/gop"
)

func execAES(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2 := vcipher.AES(args[0].(int))
	p.Ret(1, ret0, ret1, ret2)
}

func execmCipherBlockSize(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vcipher.Cipher).BlockSize()
	p.Ret(1, ret0)
}

func execmCipherEncrypt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*vcipher.Cipher).Encrypt(args[1].([]byte), args[2].([]byte))
	p.PopN(3)
}

func execmCipherDecrypt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*vcipher.Cipher).Decrypt(args[1].([]byte), args[2].([]byte))
	p.PopN(3)
}

func execmCipherCBCEncrypt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*vcipher.Cipher).CBCEncrypt(args[1].([]byte), args[2].([]byte))
	p.PopN(3)
}

func execmCipherCBCDecrypt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*vcipher.Cipher).CBCDecrypt(args[1].([]byte), args[2].([]byte))
	p.PopN(3)
}

func execmCipherPadding(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*vcipher.Cipher).Padding(args[1].([]byte))
	p.Ret(2, ret0)
}

func execmCipherUnpadding(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*vcipher.Cipher).Unpadding(args[1].([]byte))
	p.Ret(2, ret0)
}

func execmCipherCFBEncrypt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*vcipher.Cipher).CFBEncrypt(args[1].([]byte), args[2].([]byte))
	p.PopN(3)
}

func execmCipherCFBDecrypt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*vcipher.Cipher).CFBDecrypt(args[1].([]byte), args[2].([]byte))
	p.PopN(3)
}

func execmCipherOFB(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*vcipher.Cipher).OFB(args[1].([]byte), args[2].([]byte))
	p.PopN(3)
}

func execmCipherCTR(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*vcipher.Cipher).CTR(args[1].([]byte), args[2].([]byte))
	p.PopN(3)
}

func toType0(v interface{}) cipher.Block {
	if v == nil {
		return nil
	}
	return v.(cipher.Block)
}

func execNewCipher(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := vcipher.NewCipher(toType0(args[0]), args[1].([]byte))
	p.Ret(2, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("github.com/456vv/vcipher")

func init() {
	I.RegisterFuncs(
		I.Func("AES", vcipher.AES, execAES),
		I.Func("(*Cipher).BlockSize", (*vcipher.Cipher).BlockSize, execmCipherBlockSize),
		I.Func("(*Cipher).Encrypt", (*vcipher.Cipher).Encrypt, execmCipherEncrypt),
		I.Func("(*Cipher).Decrypt", (*vcipher.Cipher).Decrypt, execmCipherDecrypt),
		I.Func("(*Cipher).CBCEncrypt", (*vcipher.Cipher).CBCEncrypt, execmCipherCBCEncrypt),
		I.Func("(*Cipher).CBCDecrypt", (*vcipher.Cipher).CBCDecrypt, execmCipherCBCDecrypt),
		I.Func("(*Cipher).Padding", (*vcipher.Cipher).Padding, execmCipherPadding),
		I.Func("(*Cipher).Unpadding", (*vcipher.Cipher).Unpadding, execmCipherUnpadding),
		I.Func("(*Cipher).CFBEncrypt", (*vcipher.Cipher).CFBEncrypt, execmCipherCFBEncrypt),
		I.Func("(*Cipher).CFBDecrypt", (*vcipher.Cipher).CFBDecrypt, execmCipherCFBDecrypt),
		I.Func("(*Cipher).OFB", (*vcipher.Cipher).OFB, execmCipherOFB),
		I.Func("(*Cipher).CTR", (*vcipher.Cipher).CTR, execmCipherCTR),
		I.Func("NewCipher", vcipher.NewCipher, execNewCipher),
	)
	I.RegisterTypes(
		I.Type("Cipher", reflect.TypeOf((*vcipher.Cipher)(nil)).Elem()),
	)
}
