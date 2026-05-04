// export by github.com/goplus/ixgo/cmd/qexp

//go:build igop_lib
// +build igop_lib

package toml

import (
	q "github.com/pelletier/go-toml/v2"

	"reflect"

	"github.com/goplus/ixgo"
)

func init() {
	ixgo.RegisterPackage(&ixgo.Package{
		Name: "toml",
		Path: "github.com/pelletier/go-toml/v2",
		Deps: map[string]string{
			"bytes":         "bytes",
			"encoding":      "encoding",
			"encoding/json": "json",
			"errors":        "errors",
			"fmt":           "fmt",
			"github.com/pelletier/go-toml/v2/internal/characters": "characters",
			"github.com/pelletier/go-toml/v2/internal/tracker":    "tracker",
			"github.com/pelletier/go-toml/v2/unstable":            "unstable",
			"io":          "io",
			"math":        "math",
			"reflect":     "reflect",
			"slices":      "slices",
			"strconv":     "strconv",
			"strings":     "strings",
			"sync/atomic": "atomic",
			"time":        "time",
			"unicode":     "unicode",
		},
		Interfaces: map[string]reflect.Type{},
		NamedTypes: map[string]reflect.Type{
			"DecodeError":        reflect.TypeOf((*q.DecodeError)(nil)).Elem(),
			"Decoder":            reflect.TypeOf((*q.Decoder)(nil)).Elem(),
			"Encoder":            reflect.TypeOf((*q.Encoder)(nil)).Elem(),
			"Key":                reflect.TypeOf((*q.Key)(nil)).Elem(),
			"LocalDate":          reflect.TypeOf((*q.LocalDate)(nil)).Elem(),
			"LocalDateTime":      reflect.TypeOf((*q.LocalDateTime)(nil)).Elem(),
			"LocalTime":          reflect.TypeOf((*q.LocalTime)(nil)).Elem(),
			"StrictMissingError": reflect.TypeOf((*q.StrictMissingError)(nil)).Elem(),
		},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs: map[string]reflect.Value{
			"Marshal":    reflect.ValueOf(q.Marshal),
			"NewDecoder": reflect.ValueOf(q.NewDecoder),
			"NewEncoder": reflect.ValueOf(q.NewEncoder),
			"Unmarshal":  reflect.ValueOf(q.Unmarshal),
		},
		TypedConsts:   map[string]ixgo.TypedConst{},
		UntypedConsts: map[string]ixgo.UntypedConst{},
	})
}
