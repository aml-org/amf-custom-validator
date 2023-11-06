package main

import (
	"fmt"
	"strings"
	"unsafe"
)

// main is required for TinyGo to compile to Wasm.
func main() {}

//go:wasmimport env assignResult
func _assignResult(ptr, size uint32)

func stringToPtr(s string) (uint32, uint32) {
	ptr := unsafe.Pointer(unsafe.StringData(s))
	return uint32(uintptr(ptr)), uint32(len(s))
}

func ptrToString(ptr uint32, size uint32) string {
	return unsafe.String((*byte)(unsafe.Pointer(uintptr(ptr))), size)
}

func greetImpl(name string) {
	message := fmt.Sprintf("Hello %s", name)
	ptr, size := stringToPtr(message)
	//runtime.KeepAlive(message)
	_assignResult(ptr, size)
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

//export greet
func greet(subject *uint32, size int) {
	name := ptrToString2(subject, size)
	greetImpl(name)
}

//export alloc
func alloc(size uint32) *byte {
	buf := make([]byte, size)
	return &buf[0]
}
