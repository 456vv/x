// igop
// entryName=abc
// file=y2.go
package main

import (
	"a"
	"fmt"
)

func abc(in any) string {
	if A() != "A" {
		return "!A"
	}
	return a.Sprint(in)
}

func main() {
	s := a.Sprint("test pkg-main")
	fmt.Println(s)
}
