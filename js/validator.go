package main

import (
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
	ruleset := ptrToString2(rulesetPtr, rulesetSize)
	data := ptrToString2(dataPtr, dataSize)
	validateImpl(ruleset, data)
}

//export alloc
func alloc(size uint32) *byte {
	buf := make([]byte, size)
	return &buf[0]
}
