package vweb_dynamic

import (
	"os"
	"testing"
)

func Test_wazero(t *testing.T) {
	var wz Wazero
	wz.SetPath("testdata/wwwroot/", "wazero/add.wasm")
	err := wz.ParseFile("testdata/wwwroot/wazero/add.wasm")
	if err != nil {
		t.Fatal(err)
	}
	defer wz.Close()

	err = wz.Execute("add", os.Stdout, uint32(1), uint32(2))
	if err != nil {
		t.Fatal(err)
	}
}
