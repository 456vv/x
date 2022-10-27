// export by github.com/goplus/igop/cmd/qexp

//go:build igop_lib
// +build igop_lib

package vcipher

import (
	q "github.com/456vv/vcipher"

	"reflect"

	"github.com/goplus/igop"
)

func init() {
	igop.RegisterPackage(&igop.Package{
		Name: "vcipher",
		Path: "github.com/456vv/vcipher",
		Deps: map[string]string{
			"bytes":         "bytes",
			"crypto/aes":    "aes",
			"crypto/cipher": "cipher",
			"crypto/rand":   "rand",
		},
		Interfaces: map[string]reflect.Type{},
		NamedTypes: map[string]reflect.Type{
			"Cipher": reflect.TypeOf((*q.Cipher)(nil)).Elem(),
		},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs: map[string]reflect.Value{
			"AES":       reflect.ValueOf(q.AES),
			"NewCipher": reflect.ValueOf(q.NewCipher),
		},
		TypedConsts:   map[string]igop.TypedConst{},
		UntypedConsts: map[string]igop.UntypedConst{},
	})
}
