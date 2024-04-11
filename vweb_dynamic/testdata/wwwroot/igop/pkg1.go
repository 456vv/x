// igop
// entryname=abc
package main

import (
	"a"
	"fmt"
)

func abc(in interface{}) string {
	return a.Sprint("test pkg-abc")
}

func main() {
	s := a.Sprint("test pkg-main")
	fmt.Println(s)
}
