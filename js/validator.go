package main

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/validator"
	"strings"
	"unsafe"
)

// main is required for TinyGo to compile to Wasm.
func main() {

}

//go:wasmimport env assignReport
func _assignReport(ptr, size uint32)

func stringToPtr(s string) (uint32, uint32) {
	ptr := unsafe.Pointer(unsafe.StringData(s))
	return uint32(uintptr(ptr)), uint32(len(s))
}

func validateImpl(ruleset, data string) {
	report, _ := validator.Validate(ruleset, data, false, nil)
	ptr, size := stringToPtr(report)
	_assignReport(ptr, size)
}

func ptrToString2(subject *uint32, size int) string {
	var subjectStr strings.Builder
	pointer := uintptr(unsafe.Pointer(subject))
	for i := 0; i < size; i++ {
		s := *(*int32)(unsafe.Pointer(pointer + uintptr(i)))
		subjectStr.WriteByte(byte(s))
	}

	return subjectStr.String()
}

//export validate
func validate(rulesetPtr *uint32, rulesetSize int, dataPtr *uint32, dataSize int) {
	//var m runtime.MemStats
	//runtime.ReadMemStats(&m)
	//fmt.Printf("GO %v\n", m.HeapInuse)
	//runtime.GC()
	//runtime.ReadMemStats(&m)
	//fmt.Printf("GO %v\n", m.HeapInuse)

	fmt.Println("Reading ruleset")
	ruleset := ptrToString2(rulesetPtr, rulesetSize)
	fmt.Println("Reading data")
	data := ptrToString2(dataPtr, dataSize)

	//fmt.Println("=============================================ruleset============================================")
	//fmt.Printf("%s", ruleset)
	//fmt.Println("=============================================data============================================")
	//fmt.Printf("%s", data)
	//fmt.Printf("Ruleset ptr: %p\n", rulesetPtr)
	//fmt.Printf("Ruleset address: %p\n", &ruleset)
	//fmt.Printf("Data ptr: %p\n", dataPtr)
	//fmt.Printf("Data address: %p\n", &data)

	//runtime.KeepAlive(ruleset)
	//runtime.KeepAlive(data)

	validateImpl(ruleset, data)

}

//export alloc
func alloc(size uint32) *byte {
	buf := make([]byte, size)
	return &buf[0]
}
