package db
import(
    "strings"
    "reflect"
)

//过滤指定例
//	m []map[string]interface{}	数据库的数据
//	f func(string) bool			返回true,过滤该例
//	[]map[string]interface{}	过滤后的数据
func ColumnFilter(m []map[string]interface{}, f func(string) bool) []map[string]interface{}{
	var retn []map[string]interface{}
	for _, mv := range m {
		nm := make(map[string]interface{})
		for n, v := range mv {
			if f(n) {
				continue
			}
			nm[n]=v
		}
		retn = append(retn, nm)
	}
	return retn
}

//以数组形式返回key列的值
//	m []map[string]interface{}	数据库的数据
//	key string					key键名
//	[]interface{}				返回key对应的值，不包含nil值
func ColumnArray(m []map[string]interface{}, key string) []interface{}{
	var vals []interface{}
	 for _, mv := range m {
	 	if v, ok := mv[key]; ok {
	 		if v != nil {
	 			vals = append(vals, v)
	 		}
	 	}
	 }
	 return vals
}
//以数组形式返回key列的值
//	m []map[string]interface{}	数据库的数据
//	key string					key键名
//	[]interface{}				返回key对应的值，包含nil值
func ColumnArrayNil(m []map[string]interface{}, key string) []interface{}{
	var vals []interface{}
	 for _, mv := range m {
	 	if v, ok := mv[key]; ok {
	 		vals = append(vals, v)
	 	}
	 }
	 return vals
}

//map键转小写
func keyToLower(key string) string {
	var nkey string = key
	if strings.HasSuffix(key, "ID") {
		nkey = strings.ToLower(key)
	}else{
		bkey 	:= []byte(key)
		length	:= len(bkey)
		for i:=0;i<length;i++ {
			if bkey[i] == 95 {
				continue
			}
			if bkey[i]>=65 && bkey[i]<=90 {
				if i != 0 && i+1 != length && bkey[i+1]>=97 && bkey[i+1]<=122 {
					break
				}
				bkey[i]+=32
			}else{
				break
			}
		}
		nkey = string(bkey)
	}
	return nkey
}

func inDirect(v reflect.Value) reflect.Value {
	for ; v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface; v = v.Elem() {
	}
	return v
}