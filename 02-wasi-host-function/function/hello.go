package main

// main is required for TinyGo to compile to Wasm.
func main() {}


//export host_log_uint32
func host_log_uint32(value uint32)

// This exports an add function.
// It takes in two uint32 values
// And returns a uint32 value.
// To make this function callable from host,
// add the: "export add" comment above the function

//export add
func add(x uint32, y uint32) uint32 {
  // 🖐 a wasm module cannot print something
  //fmt.Println(x,y)
  res := x + y

  host_log_uint32(res)

  return res;
}
