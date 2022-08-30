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

func main() {
	// Choose the context to use for function calls.
	ctx := context.Background()

	// Create a new WebAssembly Runtime.
	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx) // This closes everything this Runtime created.

	_, errEnv := r.NewHostModuleBuilder("env").
		NewFunctionBuilder().
		WithFunc(func(value uint32) {
			fmt.Println("🤖:", value)
		}).
		Export("host_log_uint32").
		NewFunctionBuilder().WithFunc(logString).Export("host_log_string").
		Instantiate(ctx, r)

	if errEnv != nil {
		log.Panicln("🔴 Error with env module and host function(s):", errEnv)
	}

	_, errInstantiate := wasi_snapshot_preview1.Instantiate(ctx, r)
	if errInstantiate != nil {
		log.Panicln("🔴 Error with Instantiate:", errInstantiate)
	}

	// Load then Instantiate a WebAssembly module
	helloWasm, errLoadWasmModule := os.ReadFile("./function/hello.wasm")
	if errLoadWasmModule != nil {
		log.Panicln("🔴 Error while loading the wasm module", errLoadWasmModule)
	}

	mod, errInstanceWasmModule := r.InstantiateModuleFromBinary(ctx, helloWasm)
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

	// Get a string from wasm
	helloWorldwasmModuleFunction := mod.ExportedFunction("helloWorld")

	ptrSize, errCallFunction := helloWorldwasmModuleFunction.Call(ctx)
	if errCallFunction != nil {
		log.Panicln("🔴 Error while calling the function ", errCallFunction)
	}
	// Note: This pointer is still owned by TinyGo, so don't try to free it!
	helloWorldPtr := uint32(ptrSize[0] >> 32)
	helloWorldSize := uint32(ptrSize[0])

	// The pointer is a linear memory offset, which is where we write the name.
	if bytes, ok := mod.Memory().Read(ctx, helloWorldPtr, helloWorldSize); !ok {
		log.Panicf("🟥 Memory.Read(%d, %d) out of range of memory size %d",
			helloWorldPtr, helloWorldSize, mod.Memory().Size(ctx))
	} else {
		fmt.Println(bytes)
		fmt.Println("😃 the string message is:", string(bytes))
	}

}

func logString(ctx context.Context, module api.Module, offset, byteCount uint32) {
	buf, ok := module.Memory().Read(ctx, offset, byteCount)
	if !ok {
		log.Panicf("🟥 Memory.Read(%d, %d) out of range", offset, byteCount)
	}
	fmt.Println("👽:", string(buf))
}
