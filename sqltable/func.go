package sqltable

import (
	"reflect"
)

func quoteValue(vals []interface{}) []interface{} {
	for index, val := range vals {
		rt := reflect.TypeOf(val)
		if rt.Kind() != reflect.Ptr {
			vals[index] = &val
		}
	}
	return vals
}
func sliceContainString(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}