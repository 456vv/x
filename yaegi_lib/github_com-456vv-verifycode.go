// Code generated by 'yaegi extract github.com/456vv/verifycode'. DO NOT EDIT.

package yaegi_lib

import (
	"github.com/456vv/verifycode"
	"reflect"
)

func init() {
	Symbols["github.com/456vv/verifycode/verifycode"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"Rand":       reflect.ValueOf(verifycode.Rand),
		"RandRange":  reflect.ValueOf(verifycode.RandRange),
		"RandomText": reflect.ValueOf(verifycode.RandomText),

		// type definitions
		"Color":      reflect.ValueOf((*verifycode.Color)(nil)),
		"Font":       reflect.ValueOf((*verifycode.Font)(nil)),
		"Glyph":      reflect.ValueOf((*verifycode.Glyph)(nil)),
		"Style":      reflect.ValueOf((*verifycode.Style)(nil)),
		"VerifyCode": reflect.ValueOf((*verifycode.VerifyCode)(nil)),
	}
}
