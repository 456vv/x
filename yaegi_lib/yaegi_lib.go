package yaegi_lib

import "reflect"

// Symbols variable stores the map of stdlib symbols per package.
var Symbols = map[string]map[string]reflect.Value{}

func init() {
	Symbols["github.com/456vv/x/yaegi_lib"] = map[string]reflect.Value{
		"Symbols": reflect.ValueOf(Symbols),
	}
}
