package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tetratelabs/wazero"
	//"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/wasi_snapshot_preview1"
)

func main() {
	// Choose the context to use for function calls.
	ctx := context.Background()

	// Create a new WebAssembly Runtime.
	wasmRuntime := wazero.NewRuntimeWithConfig(wazero.NewRuntimeConfig().WithWasmCore2())
	defer wasmRuntime.Close(ctx) // This closes everything this Runtime created.

	// 👋 Add a Host Function
	// Instantiate a Go-defined module named "env"
	// that exports a function (host_log_uint32) from the host to the wasm module
	_, errEnv := wasmRuntime.NewModuleBuilder("env").
		ExportFunction("hostLogUint32", func(value uint32) {
			fmt.Println("🤖:", value)
		}).
		Instantiate(ctx, wasmRuntime)

	if errEnv != nil {
		log.Panicln("🔴 Error with env module and host function(s):", errEnv)
	}

	_, errInstantiate := wasi_snapshot_preview1.Instantiate(ctx, wasmRuntime)
	if errInstantiate != nil {
		log.Panicln("🔴 Error with Instantiate:", errInstantiate)
	}

	// Load then Instantiate a WebAssembly module
	helloWasm, errLoadWasmModule := os.ReadFile("./function/hello.wasm")
	if errLoadWasmModule != nil {
		log.Panicln("🔴 Error while loading the wasm module", errLoadWasmModule)
	}

	mod, errInstanceWasmModule := wasmRuntime.InstantiateModuleFromBinary(ctx, helloWasm)
	if errInstanceWasmModule != nil {
		log.Panicln("🔴 Error while creating module instance ", errInstanceWasmModule)
	}

	// Get references to WebAssembly function: "add"
	addWasmModuleFunction := mod.ExportedFunction("add")

	// Now, we can call "add", which reads the string we wrote to memory!
	// result []uint64
	result, errCallFunction := addWasmModuleFunction.Call(ctx, 20, 22)
	if errCallFunction != nil {
		log.Panicln("🔴 Error while calling the function ", errCallFunction)
	}

	fmt.Println("result:", result[0])

}
