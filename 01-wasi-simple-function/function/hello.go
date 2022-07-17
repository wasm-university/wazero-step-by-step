package main

// main is required for TinyGo to compile to Wasm.
func main() {}

// This exports an add function.
// It takes in two uint32 values
// And returns a uint32 value.
// To make this function callable from host,
// we need to add the: "export add" comment above the function

//export add
func add(x uint32, y uint32) uint32 {
  // 🖐 a wasm module cannot print something
  //fmt.Println(x,y)
  return x + y;
}
