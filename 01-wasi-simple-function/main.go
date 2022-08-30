package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func main() {
	// Choose the context to use for function calls.
	ctx := context.Background()

	// Create a new WebAssembly Runtime.
	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx) // This closes everything this Runtime created.

	// TinyGo specific to use WASI
	_, err := wasi_snapshot_preview1.Instantiate(ctx, r)
	if err != nil {
		log.Panicln(err)
	}

	// Load then Instantiate a WebAssembly module
	wasmPath1 := "./functions/hey/target/wasm32-wasi/debug/hey.wasm"
	//wasmPath2 := "./functions/hello/hello.wasm"
	helloWasm, err := os.ReadFile(wasmPath1)

	if err != nil {
		log.Panicln(err)
	}

	mod, err := r.InstantiateModuleFromBinary(ctx, helloWasm)
	if err != nil {
		log.Panicln(err)
	}

	// Get references to WebAssembly function: "add"
	addWasmModuleFunction := mod.ExportedFunction("add")

	// Now, we can call "add", which reads the data we wrote to memory!
	// result []uint64
	result, err := addWasmModuleFunction.Call(ctx, 20, 22)
	if err != nil {
		log.Panicln(err)
	}

	fmt.Println("result:", result[0])

}
