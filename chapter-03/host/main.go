// Package main of the host application
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

// PrintString : print a string to the console
var PrintString = api.GoModuleFunc(func(ctx context.Context, module api.Module, stack []uint64) {

	position := uint32(stack[0])
	length := uint32(stack[1])

	buffer, ok := module.Memory().Read(position, length)
	if !ok {
		log.Panicf("Memory.Read(%d, %d) out of range", position, length)
	}
	fmt.Println(string(buffer))

	stack[0] = 0 // return 0
})

func main() {
	// Choose the context to use for function calls.
	ctx := context.Background()
	
	// Create a new WebAssembly Runtime.
	runtime := wazero.NewRuntime(ctx)

    // This closes everything this Runtime created.
	defer runtime.Close(ctx) 

	// START: Host functions
	builder := runtime.NewHostModuleBuilder("env")

	// hostPrintString
	builder.NewFunctionBuilder().
		WithGoModuleFunction(PrintString, 
			[]api.ValueType{
				api.ValueTypeI32, // string position
				api.ValueTypeI32, // string length
			}, 
			[]api.ValueType{api.ValueTypeI32}).
		Export("hostPrintString")

	_, err := builder.Instantiate(ctx)
	if err != nil {
		log.Panicln("Error with env module and host function(s):", err)
	}
	// END: Host functions
	
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
