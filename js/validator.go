package main

// #include <stdlib.h>
import "C"
import (
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

func stringToLeakedPtr(s string) (uint32, uint32) {
	size := C.ulong(len(s))
	ptr := unsafe.Pointer(C.malloc(size))
	copy(unsafe.Slice((*byte)(ptr), size), s)
	return uint32(uintptr(ptr)), uint32(size)
}

//export greet
func greet(name string) {
	message := "ASD FROM GO"
	ptr, size := stringToPtr(message)
	//runtime.KeepAlive(message)
	_assignResult(ptr, size)
}
