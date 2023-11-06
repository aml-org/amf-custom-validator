package main

// #include <stdlib.h>
import "C"
import "unsafe"

// main is required for TinyGo to compile to Wasm.
func main() {}

//go:wasmimport env log
func _log(ptr, size uint32)

func stringToPtr(s string) (uint32, uint32) {
	ptr := unsafe.Pointer(unsafe.StringData(s))
	return uint32(uintptr(ptr)), uint32(len(s))
}

//export greet
func greet() {
	ptr, size := stringToPtr("ASD FROM GO")
	_log(ptr, size)
}
