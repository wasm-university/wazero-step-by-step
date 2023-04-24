// Package main: this is a wasm module
package main

import "unsafe"

func main () {}

//export hostPrintString
func hostPrintString(pos, sisze uint32) uint32

// Print a string
func Print(message string) {
    buffer := []byte(message)
	bufferPtr := &buffer[0]
	unsafePtr := uintptr(unsafe.Pointer(bufferPtr))

	pos := uint32(unsafePtr)
	size := uint32(len(buffer))

	hostPrintString(pos, size)
}


//export hello
func hello(valuePosition *uint32, length int) uint64 {
	
	// read the memory to get the parameter
	valueBytes := readBufferFromMemory(valuePosition, length)

	message := "Hello " + string(valueBytes)

	Print("👋 from the module: " + message)

	// copy the value to memory
	posSizePairValue := copyBufferToMemory([]byte(message))

	// return the position and size
	return posSizePairValue

}

// readBufferFromMemory returns a buffer from WebAssembly
func readBufferFromMemory(bufferPosition *uint32, length int) []byte {
	subjectBuffer := make([]byte, length)
	pointer := uintptr(unsafe.Pointer(bufferPosition))
	for i := 0; i < length; i++ {
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
