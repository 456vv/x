package db
import(
    "strings"
)
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