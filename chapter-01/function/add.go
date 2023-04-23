package main

// main is required for TinyGo to compile to Wasm.
func main() {}

// This exports an add function.
// It takes in two uint32 values
// And returns an uint32 value.
// To make this function callable from host,
// we need to add the: "export add" comment above the function

//export add
func add(x uint32, y uint32) uint32 {
	return x + y
}
