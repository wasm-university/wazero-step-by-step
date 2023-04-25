// Package main: this is a wasm module
package main

import (
	"fmt"
	"unsafe"
)

func main () {}

//export hello
func hello(valuePosition *uint32, length uint32) uint64 {
	
	// read the memory to get the parameter
	valueBytes := readBufferFromMemory(valuePosition, length)

	message := "Hello " + string(valueBytes)

	fmt.Println("🎃" + message)

	// copy the value to memory
	posSizePairValue := copyBufferToMemory([]byte(message))

	// return the position and size
	return posSizePairValue

}

// readBufferFromMemory returns a buffer from WebAssembly
func readBufferFromMemory(bufferPosition *uint32, length uint32) []byte {
	subjectBuffer := make([]byte, length)
	pointer := uintptr(unsafe.Pointer(bufferPosition))
	for i := 0; i < int(length); i++ {
		s := *(*int32)(unsafe.Pointer(pointer + uintptr(i)))
		subjectBuffer[i] = byte(s)
	}
	return subjectBuffer
}

// copyBufferToMemory returns a single value (a kind of pair with position and length)
func copyBufferToMemory(buffer []byte) uint64 {
	bufferPtr := &buffer[0]
	unsafePtr := uintptr(unsafe.Pointer(bufferPtr))

	ptr := uint32(unsafePtr)
	size := uint32(len(buffer))

	return (uint64(ptr) << uint64(32)) | uint64(size)
}
