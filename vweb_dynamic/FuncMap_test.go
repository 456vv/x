package vweb_dynamic

import (
	"reflect"
	"testing"
)

func Test_pkgTypes_new(t *testing.T) {
	tests := []struct {
		T      *pkgTypes
		args   []any
		result func(v reflect.Value) bool
	}{
		{
			T:    &pkgTypes{t: reflect.TypeOf((*testing.T)(nil)).Elem()},
			args: []any{(testing.TB)(t)},
			result: func(v reflect.Value) bool {
				// testing.TB.(*testing.T)
				return v.Type().String() != "*testing.T"
			},
		}, {
			T: &pkgTypes{t: reflect.TypeOf((*testing.T)(nil)).Elem()},
			result: func(v reflect.Value) bool {
				// new(testing.T)
				return !v.IsValid() && v.Kind() != reflect.Pointer && v.String() == "&{[]}"
			},
		}, {
			T:    &pkgTypes{t: reflect.TypeOf((*testing.T)(nil)).Elem()},
			args: []any{nil},
			result: func(v reflect.Value) bool {
				//(*testing.T)(nil)
				return !v.IsValid() && v.Elem().Kind() == reflect.Invalid && v.Interface() == nil
			},
		},
	}
	for _, tt := range tests {
		got := tt.T.new(tt.args...)
		gotval := reflect.ValueOf(got)
		if tt.result(gotval) {
			t.Errorf("pkgTypes.Exec() = %v", got)
		}
	}
}

func Test_pkgInterface_to(t *testing.T) {
	tests := []struct {
		T      *pkgInterface
		arg    any
		result func(reflect.Value) bool
	}{
		{
			T:   &pkgInterface{t: reflect.TypeOf((*testing.TB)(nil)).Elem()},
			arg: t,
			result: func(a reflect.Value) bool {
				// 转为 testing.TB 接口后，返回是 interface{}
				return a.Type().String() != "*testing.T"
			},
		},
	}
	for _, tt := range tests {
		got := tt.T.to(tt.arg)
		gotval := reflect.ValueOf(got)
		if tt.result(gotval) {
			t.Errorf("pkgInterface.to() = %v", got)
		}
	}
}
