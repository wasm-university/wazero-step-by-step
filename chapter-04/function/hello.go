// Package main: this is a wasm module
package main

import (
	"reflect"
	"unsafe"
)

func main() {}

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


//export hostTalk
func hostTalk(messagePosition, messageLength uint32, returnValuePosition **byte, returnValueLength *int) uint32


func stringToPtr(s string) (uint32, uint32) {
	buf := []byte(s)
	ptr := &buf[0]
	unsafePtr := uintptr(unsafe.Pointer(ptr))
	return uint32(unsafePtr), uint32(len(buf))
}

func getStringResult(buffPtr *byte, buffSize int) string {

	result := *(*string)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(buffPtr)),
		Len:  uintptr(buffSize),
		Cap:  uintptr(buffSize),
	}))

	//Log("🟧 getStringResult: " + result)

	return result
}



//export hello
func hello(valuePosition *uint32, length int) uint64 {

	// read the memory to get the parameter
	valueBytes := readBufferFromMemory(valuePosition, length)

	message := "Hello " + string(valueBytes)

	Print("👋 from the module: " + message)

	var buffPtr *byte
	var buffSize int
	posMessage, lengthMessage := stringToPtr("this is a message")
	
	hostTalk(posMessage, lengthMessage, &buffPtr, &buffSize)

	valueFromHost := getStringResult(buffPtr, buffSize) 

	Print("From host: " + valueFromHost)

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
