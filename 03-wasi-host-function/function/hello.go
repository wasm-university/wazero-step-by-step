package main

import (
  "unsafe"
  "strconv"
)

// main is required for TinyGo to compile to Wasm.
func main() {}

//export host_log_string
func host_log_string(ptr uint32, size uint32)

//export host_log_uint32
func host_log_uint32(value uint32)


//export add
func add(x uint32, y uint32) uint32 {
  // 🖐 a wasm module cannot print something
  //fmt.Println(x,y)
  res := x + y

  host_log_uint32(res)

	ptr, size := stringToPtr("from wasm: " + strconv.FormatUint(uint64(res), 10))
	host_log_string(ptr, size)

  return res;
}

// stringToPtr returns a pointer and size pair for the given string in a way
// compatible with WebAssembly numeric types.
func stringToPtr(s string) (uint32, uint32) {
	buf := []byte(s)
	ptr := &buf[0]
	unsafePtr := uintptr(unsafe.Pointer(ptr))
	return uint32(unsafePtr), uint32(len(buf))
}
