// export by github.com/goplus/gossa/cmd/qexp

//go:build gossa_lib
// +build gossa_lib

package vcipher

import (
	q "github.com/456vv/vcipher"

	"reflect"

	"github.com/goplus/gossa"
)

func init() {
	gossa.RegisterPackage(&gossa.Package{
		Name: "vcipher",
		Path: "github.com/456vv/vcipher",
		Deps: map[string]string{
			"bytes":         "bytes",
			"crypto/aes":    "aes",
			"crypto/cipher": "cipher",
			"crypto/rand":   "rand",
		},
		Interfaces: map[string]reflect.Type{},
		NamedTypes: map[string]gossa.NamedType{
			"Cipher": {reflect.TypeOf((*q.Cipher)(nil)).Elem(), "", "BlockSize,CBCDecrypt,CBCEncrypt,CFBDecrypt,CFBEncrypt,CTR,Decrypt,Encrypt,OFB,Padding,Unpadding"},
		},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs: map[string]reflect.Value{
			"AES":       reflect.ValueOf(q.AES),
			"NewCipher": reflect.ValueOf(q.NewCipher),
		},
		TypedConsts:   map[string]gossa.TypedConst{},
		UntypedConsts: map[string]gossa.UntypedConst{},
	})
}
