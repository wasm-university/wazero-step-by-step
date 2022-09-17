fn main() {
    println!("Hello, world!");

}


#[link(wasm_import_module = "env")]
extern "C" {
  #[link_name = "hostLogUint32"]
  fn hostLogUint32(value: u32);
}


#[no_mangle]
pub unsafe fn add(a: u32, b: u32) -> u32 {
  hostLogUint32(1968);
  hostLogUint32(a+b);
  return a + b;
}
