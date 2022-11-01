package main

import (
	"a"
	"fmt"
)

func Main(in interface{}) string {
	return a.Sprint("test pkg")
}

func main() {
	s := a.Sprint("test pkg")
	fmt.Println(s)
}
