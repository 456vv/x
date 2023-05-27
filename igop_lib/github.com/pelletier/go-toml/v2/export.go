// export by github.com/goplus/igop/cmd/qexp

//go:build igop_lib
// +build igop_lib

package toml

import (
	q "github.com/pelletier/go-toml/v2"

	"reflect"

	"github.com/goplus/igop"
)

func init() {
	igop.RegisterPackage(&igop.Package{
		Name: "toml",
		Path: "github.com/pelletier/go-toml/v2",
		Deps: map[string]string{
			"bytes":    "bytes",
			"encoding": "encoding",
			"errors":   "errors",
			"fmt":      "fmt",
			"github.com/pelletier/go-toml/v2/internal/characters": "characters",
			"github.com/pelletier/go-toml/v2/internal/danger":     "danger",
			"github.com/pelletier/go-toml/v2/internal/tracker":    "tracker",
			"github.com/pelletier/go-toml/v2/unstable":            "unstable",
			"io":          "io",
			"io/ioutil":   "ioutil",
			"math":        "math",
			"reflect":     "reflect",
			"sort":        "sort",
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
		TypedConsts:   map[string]igop.TypedConst{},
		UntypedConsts: map[string]igop.UntypedConst{},
	})
}
