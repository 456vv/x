package main

import(
	"this"
	"fmt"
	"time"
	pkg "github.com/foo/pkg"
)

func Main(a interface{}) string {
	b := pkg.NewSample()()
	if b != "root Fromage Cheese" {
		panic("error")
	}
	if !this.IsNil(nil) {
		panic("error")
	}

	time.Sleep(time.Millisecond)
	return fmt.Sprint(a)
}

