package main

import (
	"reflect"
	"strconv"
	"unsafe"
)

// main is required for TinyGo to compile to Wasm.
func main() {

	ptr, size := stringToPtr("from the main method: " + getStringFromHost())
	host_log_string(ptr, size)

	msg := talk("hello, are you there ?")

	ptr2, size2 := stringToPtr(msg)
	host_log_string(ptr2, size2)

}

//export allocate_buffer
func allocateBuffer(size uint32) *byte {
	// Allocate the in-Wasm memory region and returns its pointer to hosts.
	// The region is supposed to store random strings generated in hosts,
	// meaning that this is called "inside" of host_get_string.
	buf := make([]byte, size)
	return &buf[0]
}

//export host_talk
//go:linkname host_talk
func host_talk(ptr uint32, size uint32, retBufPtr **byte, retBufSize *int)

//export host_get_string
//go:linkname host_get_string
func host_get_string(retBufPtr **byte, retBufSize *int)

// Get the string from the hosts.
func getStringFromHost() string {
	var bufPtr *byte
	var bufSize int
	host_get_string(&bufPtr, &bufSize)
	//nolint
	return *(*string)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(bufPtr)),
		Len:  uintptr(bufSize),
		Cap:  uintptr(bufSize),
	}))
}

func talk(message string) string {

	ptr, size := stringToPtr(message)

	var bufPtr *byte
	var bufSize int

	host_talk(ptr, size, &bufPtr, &bufSize)

	return *(*string)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(bufPtr)),
		Len:  uintptr(bufSize),
		Cap:  uintptr(bufSize),
	}))

}

//export host_log_string
//go:linkname host_log_string
func host_log_string(ptr uint32, size uint32)

//export host_log_uint32
//go:linkname host_log_uint32
func host_log_uint32(value uint32)

//export add
func add(x uint32, y uint32) uint32 {
	// 🖐 a wasm module cannot print something
	//fmt.Println(x,y)
	res := x + y

	host_log_uint32(res)

	ptr, size := stringToPtr("from wasm: " + strconv.FormatUint(uint64(res), 10))
	host_log_string(ptr, size)

	return res
}

// 🖐️ returns a pointer/size pair packed into a uint64.
// Note: This uses a uint64 instead of two result values for compatibility with
// WebAssembly 1.0.
// https://stackoverflow.com/questions/5801008/go-and-operators
// https://stackoverflow.com/questions/41790574/bitmask-multiple-values-in-int64

//export helloWorld
func helloWorld() (ptrAndSize uint64) {
	ptr, size := stringToPtr("👋 hello world, I'm very happy to meet you, I love what you are doing my friend")
	return (uint64(ptr) << uint64(32)) | uint64(size)
}

//export sayHello
func sayHello(ptr, size uint32) (ptrAndSize uint64) {
	// get the parameter
	name := ptrToString(ptr, size)

	ptr, size = stringToPtr("👋 hello " + name + " 😃")
	return (uint64(ptr) << uint64(32)) | uint64(size)
}

// stringToPtr returns a pointer and size pair for the given string in a way
// compatible with WebAssembly numeric types.
func stringToPtr(s string) (uint32, uint32) {
	buf := []byte(s)
	ptr := &buf[0]
	unsafePtr := uintptr(unsafe.Pointer(ptr))
	return uint32(unsafePtr), uint32(len(buf))
}

// ptrToString returns a string from WebAssembly compatible numeric types
// representing its pointer and length.
func ptrToString(ptr uint32, size uint32) string {
	// Get a slice view of the underlying bytes in the stream. We use SliceHeader, not StringHeader
	// as it allows us to fix the capacity to what was allocated.
	return *(*string)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(ptr),
		Len:  uintptr(size), // Tinygo requires these as uintptrs even if they are int fields.
		Cap:  uintptr(size), // ^^ See https://github.com/tinygo-org/tinygo/issues/1284
	}))
}
