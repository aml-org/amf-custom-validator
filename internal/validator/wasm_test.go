package validator

import (
	"fmt"
	"reflect"
	"testing"
)

const basicPath relativePath = "../../test/data/basic/profile1.yaml"

func TestCompileProfileWASMFrom(t *testing.T) {
	wasm, err := ProcessProfileWASM(read(basicPath), false, nil)

	if err != nil {
		panic(err)
	}

	fmt.Println(reflect.TypeOf(wasm))
	fmt.Println(wasm)
}
