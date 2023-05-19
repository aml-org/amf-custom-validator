package generator

import (
	"context"
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/config"
	"github.com/aml-org/amf-custom-validator/internal/parser"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
	"github.com/aml-org/amf-custom-validator/test"
	"github.com/open-policy-agent/opa/rego"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestGenerated(t *testing.T) {

	for _, fix := range test.Fixtures("../../test/data/basic") {
		profile.GenReset()
		prof, err := parser.Parse(fix.ReadProfile())
		if err != nil {
			panic(err)
		}
		profile.GenReset()
		generated := Generate(*prof)
		success, err := validateRegoUnit(generated)
		if !success {
			println(generated.Code)
			t.Error(err)
		}
		if config.Override {
			test.ForceWrite(fix.Generated, generated.Code)
		} else {
			actual := strings.TrimSpace(generated.Code)
			expected := strings.TrimSpace(fix.ReadGenerated())

			if actual != expected {
				t.Errorf("%s > Actual did not match expected", fix.Profile)
			}
		}

	}
}

func validateRegoUnit(module RegoUnit) (bool, error) {
	validator := rego.New(
		rego.Query("data."+module.Name+"."+module.Entrypoint),
		rego.Module(module.Name+".rego", module.Code),
	)
	ctx := context.Background()
	_, err := validator.Eval(ctx)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func TestGeneratedWASM(t *testing.T) {

	// read & parse profile
	bytes, err := ioutil.ReadFile("../../test/data/basic/profile1.yaml")
	if err != nil {
		panic(err)
	}
	profileText := string(bytes)
	fmt.Println(profileText)
	profileModel, err := parser.Parse(profileText)

	// generate Rego from parsed profile
	ru := Generate(*profileModel)
	fmt.Println(ru)
	validator := rego.New(
		rego.Query("data."+ru.Name+"."+ru.Entrypoint),
		rego.Module(ru.Name+".rego", ru.Code),
	)

	// generate WASM from Rego & save it to a file
	ctx := context.Background()
	compileResult, err := validator.Compile(ctx)
	module, err := js.Global().Get("WebAssembly").Call("compile", compileResult.Bytes) // Create a Wasm module from the byte array
	if err != nil {
		fmt.Println("Failed to compile Wasm module:", err)
		return
	}
	err = ioutil.WriteFile("example.wasm", module.Get("exports").Get("memory").Call("buffer").Interface().([]byte), os.ModePerm) // Save the Wasm module to a file
	if err != nil {
		fmt.Println("Failed to write Wasm file:", err)
		return
	}

	// execute the wasm exported function
	wasmBytes, err := ioutil.ReadFile("example.wasm") // Load the Wasm module from the file
	if err != nil {
		fmt.Println("Failed to read Wasm file:", err)
		return
	}
	wasmModule, err := js.Global().Get("WebAssembly").Call("instantiate", wasmBytes) // Create a new WebAssembly instance
	if err != nil {
		fmt.Println("Failed to instantiate Wasm module:", err)
		return
	}
	result := wasmModule.Get("exports").Call("add", 3, 4) // Call a function exported by the Wasm module
	fmt.Println("Result of Wasm execution:", result.Int())
}

func TestGeneratedWASM2(t *testing.T) {

	// read & parse profile
	bytes, err := ioutil.ReadFile("../../test/data/basic/profile1.yaml")
	if err != nil {
		panic(err)
	}
	profileText := string(bytes)
	fmt.Println(profileText)
	profileModel, err := parser.Parse(profileText)

	// generate Rego from parsed profile
	ru := Generate(*profileModel)
	fmt.Println(ru)
	validator := rego.New(
		rego.Query("data."+ru.Name+"."+ru.Entrypoint),
		rego.Module(ru.Name+".rego", ru.Code),
	)

	// generate WASM from Rego & save it to a file
	ctx := context.Background()
	compileResult, err := validator.Compile(ctx)
	module, err := js.Global().Get("WebAssembly").Call("compile", compileResult.Bytes) // Create a Wasm module from the byte array
	if err != nil {
		fmt.Println("Failed to compile Wasm module:", err)
		return
	}
	err = ioutil.WriteFile("example.wasm", module.Get("exports").Get("memory").Call("buffer").Interface().([]byte), os.ModePerm) // Save the Wasm module to a file
	if err != nil {
		fmt.Println("Failed to write Wasm file:", err)
		return
	}

	// execute the wasm exported function
	wasmBytes, err := ioutil.ReadFile("example.wasm") // Load the Wasm module from the file
	if err != nil {
		fmt.Println("Failed to read Wasm file:", err)
		return
	}
	wasmModule, err := js.Global().Get("WebAssembly").Call("instantiate", wasmBytes) // Create a new WebAssembly instance
	if err != nil {
		fmt.Println("Failed to instantiate Wasm module:", err)
		return
	}
	result := wasmModule.Get("exports").Call("add", 3, 4) // Call a function exported by the Wasm module
	fmt.Println("Result of Wasm execution:", result.Int())
}
