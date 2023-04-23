// Package main of the host application
package main

import (
	"context"
	"fmt"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"log"
	"os"
)

func main() {
	// Choose the context to use for function calls.
	ctx := context.Background()
	
	// Create a new WebAssembly Runtime.
	runtime := wazero.NewRuntime(ctx)

    // This closes everything this Runtime created.
	defer runtime.Close(ctx) 

	// Instantiate WASI
	wasi_snapshot_preview1.MustInstantiate(ctx, runtime)

	// Load the WebAssembly module
	wasmPath := "../function/hello.wasm"
	helloWasm, err := os.ReadFile(wasmPath)

	if err != nil {
		log.Panicln(err)
	}

	// Instantiate the guest Wasm into the same runtime. 
    // It exports the `hello` function, 
    // implemented in WebAssembly.
	mod, err := runtime.Instantiate(ctx, helloWasm)
	if err != nil {
		log.Panicln(err)
	}

	// Get the reference to the WebAssembly function: "hello"
	helloFunction := mod.ExportedFunction("hello")

	// Function parameter
	name := "Bob Morane"
	nameSize := uint64(len(name))

	// These function are exported by TinyGo
	malloc := mod.ExportedFunction("malloc")
	free := mod.ExportedFunction("free")

	// Allocate Memory for "Bob Morane"
	results, err := malloc.Call(ctx, nameSize)
	if err != nil {
		log.Panicln(err)
	}
	namePosition := results[0]
	
	// This pointer is managed by TinyGo, 
	// but TinyGo is unaware of external usage.
	// So, we have to free it when finished
	defer free.Call(ctx, namePosition)

	// Copy "Bob Morane" to memory
	if !mod.Memory().Write(uint32(namePosition), []byte(name)) {
		log.Panicf("out of range of memory size")
	}

	// Now, we can call "hello" with the position and the size of "Bob Morane"
	// the result type is []uint64
	result, err := helloFunction.Call(ctx, namePosition, nameSize)
	if err != nil {
		log.Panicln(err)
	}

	// Extract the position and size of the returned value
	valuePosition := uint32(result[0] >> 32)
	valueSize := uint32(result[0])

	// Read the value from the memory
	if bytes, ok := mod.Memory().Read(valuePosition, valueSize); !ok {
		log.Panicf("out of range of memory size")
	} else {
		fmt.Println("Returned value :", string(bytes))
	}
}
