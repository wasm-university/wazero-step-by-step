// Package main: this is a wasm module
package main

import (
	"unsafe"
)

func main() {}

//export hostPrintString
func hostPrintString(pos, sisze uint32) uint32

// Print a string
func Print(message string) {
	/*
	buffer := []byte(message)
	bufferPtr := &buffer[0]
	unsafePtr := uintptr(unsafe.Pointer(bufferPtr))

	pos := uint32(unsafePtr)
	size := uint32(len(buffer))
	*/

	pos, size := getBufferPosSize([]byte(message))

	hostPrintString(pos, size)
}


//export hostTalk
func hostTalk(messagePosition, messageLength uint32, returnValuePosition **uint32, returnValueLength *uint32) uint32

// Talk is an helper to use the hostTalk function
func Talk(messageToHost string) string {

	messagePosition, messageSize := getBufferPosSize([]byte(messageToHost))
	
	// This will be use to get the response from the host
	var responseBufferPtr *uint32
	var responseBufferSize uint32

	// Send the lessage to the host
	hostTalk(messagePosition, messageSize, &responseBufferPtr, &responseBufferSize)

	responseFromHost := readBufferFromMemory(responseBufferPtr, responseBufferSize)

	return string(responseFromHost)
}


//export hello
func hello(valuePosition *uint32, length uint32) uint64 {

	// read the memory to get the parameter
	valueBytes := readBufferFromMemory(valuePosition, length)

	message := "Hello " + string(valueBytes)

	Print("👋 from the module: " + message)

	Print(Talk("Hey there!"))

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

// getStringPosSize returns the memory position and size of the string
func getStringPosSize(s string) (uint32, uint32) {
	buff := []byte(s)
	ptr := &buff[0]
	unsafePtr := uintptr(unsafe.Pointer(ptr))
	return uint32(unsafePtr), uint32(len(buff))
}

// getBufferPosSize returns the memory position and size of the buffer
func getBufferPosSize(buff []byte) (uint32, uint32) {
	ptr := &buff[0]
	unsafePtr := uintptr(unsafe.Pointer(ptr))
	return uint32(unsafePtr), uint32(len(buff))
}

/* 
Allocate the in-Wasm memory region and returns its pointer to hosts.
The region is supposed to store random strings generated in hosts
*/
//export allocateBuffer
func allocateBuffer(size uint32) *byte {
	buf := make([]byte, size)
	return &buf[0]
}


/*
func getStringResult(buffPtr *byte, buffSize int) string {

	result := *(*string)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(buffPtr)),
		Len:  uintptr(buffSize),
		Cap:  uintptr(buffSize),
	}))
	return result
}
*/