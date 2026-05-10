package vweb_dynamic

import (
	"reflect"
	"testing"

	"github.com/issue9/assert/v4"
)

func Test_pkgTypes_new(t *testing.T) {
	tests := []struct {
		T      *pkgTypes
		args   []any
		result func(v reflect.Value) bool
	}{
		{
			T:    &pkgTypes{t: reflect.TypeOf((*testing.T)(nil))},
			args: []any{(testing.TB)(t)},
			result: func(v reflect.Value) bool {
				return v.Type().String() == "*testing.T" && v.IsValid()
			},
		}, {
			T: &pkgTypes{t: reflect.TypeOf((*testing.T)(nil)).Elem()},
			result: func(v reflect.Value) bool {
				return v.IsValid() && v.Kind() == reflect.Struct
			},
		},
	}
	for _, tt := range tests {
		got := tt.T.new(tt.args...)
		gotval := reflect.ValueOf(got)
		if !tt.result(gotval) {
			t.Errorf("pkgTypes.Exec() = %v", got)
		}
	}
}

func Test_pkgInterface_to(t *testing.T) {
	tests := []struct {
		T      pkgInterface
		result func(reflect.Value) bool
	}{
		{
			T: pkgInterface{t: reflect.TypeOf((*testing.TB)(nil)).Elem()},
			result: func(a reflect.Value) bool {
				return a.Type().String() == "*testing.T"
			},
		},
	}

	for _, tt := range tests {
		got := tt.T.to(t)
		gotval := reflect.ValueOf(got)
		if !tt.result(gotval) {
			t.Errorf("pkgInterface.to() = %v", got)
		}
	}
}

func Test_entryname(t *testing.T) {
	at := assert.New(t, true)

	tests := []struct {
		name   string
		result string
	}{
		{
			name:   "index.html",
			result: "Main",
		}, {
			name:   "",
			result: "Main",
		}, {
			name:   "/",
			result: "Main",
		}, {
			name:   ".",
			result: "Main",
		}, {
			name:   "_a",
			result: "Main",
		}, {
			name:   "@a",
			result: "Main",
		}, {
			name:   "A",
			result: "A",
		}, {
			name:   "a",
			result: "A",
		}, {
			name:   "A_",
			result: "Main",
		},
	}
	for _, tt := range tests {
		at.Equal(tt.result, entryname("", tt.name))
	}
}
