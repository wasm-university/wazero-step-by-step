
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
	wasmPath := "../function/add.wasm"
	addWasm, err := os.ReadFile(wasmPath)

	if err != nil {
		log.Panicln(err)
	}

	// Instantiate the guest Wasm into the same runtime. 
    // It exports the `add` function, 
    // implemented in WebAssembly.
	mod, err := runtime.Instantiate(ctx, addWasm)
	if err != nil {
		log.Panicln(err)
	}

	// Get the reference to the WebAssembly function: "add"
	addFunction := mod.ExportedFunction("add")

	// Now, we can call "add"
	// result []uint64
	result, err := addFunction.Call(ctx, 20, 22)
	if err != nil {
		log.Panicln(err)
	}

	fmt.Println("result:", result[0])	
}