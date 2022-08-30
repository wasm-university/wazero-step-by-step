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
		NewFunctionBuilder().WithFunc(giveMeString).Export("host_get_string").
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

	// ======================================================
	// Get a string from wasm
	// ======================================================
	helloWorldWasmModuleFunction := mod.ExportedFunction("helloWorld")

	ptrSize, errCallFunction := helloWorldWasmModuleFunction.Call(ctx)
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
		fmt.Println("😃 the string message is:", string(bytes))
	}

	// ======================================================
	// Pass a string param and get a string result
	// ======================================================
	// Let's use the argument to this main function in Wasm.
	name := "Bob Morane"
	nameSize := uint64(len(name))
	// Get references to WebAssembly functions we'll use in this example.

	sayHelloWasmModuleFunction := mod.ExportedFunction("sayHello")

	// These are undocumented, but exported. See tinygo-org/tinygo#2788
	malloc := mod.ExportedFunction("malloc")
	free := mod.ExportedFunction("free")

	// Instead of an arbitrary memory offset, use TinyGo's allocator. Notice
	// there is nothing string-specific in this allocation function. The same
	// function could be used to pass binary serialized data to Wasm.
	results, err := malloc.Call(ctx, nameSize)
	if err != nil {
		log.Panicln(err)
	}
	namePtr := results[0]
	// This pointer is managed by TinyGo, but TinyGo is unaware of external usage.
	// So, we have to free it when finished
	defer free.Call(ctx, namePtr)

	// The pointer is a linear memory offset, which is where we write the name.
	if !mod.Memory().Write(ctx, uint32(namePtr), []byte(name)) {
		log.Panicf("🟥 Memory.Write(%d, %d) out of range of memory size %d",
			namePtr, nameSize, mod.Memory().Size(ctx))
	}

	// Finally, we get the greeting message "greet" printed. This shows how to
	// read-back something allocated by TinyGo.
	sayHelloPtrSize, err := sayHelloWasmModuleFunction.Call(ctx, namePtr, nameSize)
	if err != nil {
		log.Panicln(err)
	}
	// Note: This pointer is still owned by TinyGo, so don't try to free it!
	sayHelloPtr := uint32(sayHelloPtrSize[0] >> 32)
	sayHelloSize := uint32(sayHelloPtrSize[0])
	// The pointer is a linear memory offset, which is where we write the name.
	if bytes, ok := mod.Memory().Read(ctx, sayHelloPtr, sayHelloSize); !ok {
		log.Panicf("Memory.Read(%d, %d) out of range of memory size %d",
			sayHelloPtr, sayHelloSize, mod.Memory().Size(ctx))
	} else {
		fmt.Println("👋 saying hello :", string(bytes))
	}

}

func logString(ctx context.Context, module api.Module, offset, byteCount uint32) {
	buf, ok := module.Memory().Read(ctx, offset, byteCount)
	if !ok {
		log.Panicf("🟥 Memory.Read(%d, %d) out of range", offset, byteCount)
	}
	fmt.Println("👽:", string(buf))
}

func giveMeString(ctx context.Context, module api.Module, retBufPtr uint32, retBufSize uint32) {

	message := "🚀 this is a string coming from the host"
	lengthOfTheMessage := len(message)
	results, err := module.ExportedFunction("allocate_buffer").Call(ctx, uint64(lengthOfTheMessage))
	if err != nil {
		log.Panicln(err)
	}

	offset := uint32(results[0])
	module.Memory().WriteUint32Le(ctx, retBufPtr, offset)
	module.Memory().WriteUint32Le(ctx, retBufSize, uint32(lengthOfTheMessage))

	// add the message to the memory of the module
	module.Memory().Write(ctx, offset, []byte(message))

}
