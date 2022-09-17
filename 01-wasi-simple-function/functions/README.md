
## Rust

https://bytecodealliance.github.io/cargo-wasi/install.html

Installing RUST (https://www.rust-lang.org/tools/install)
```bash
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
```

```bash
cargo install cargo-wasi
cargo wasi --version
```


``` bash
cargo new hey
cd hey
cargo wasi run

cargo build --target wasm32-wasi
wasmedge --reactor target/wasm32-wasi/debug/hey.wasm add 23 19

wasmtime --invoke add target/wasm32-wasi/debug/hey.wasm 23 19
wasmtime --invoke _start target/wasm32-wasi/debug/hey.wasm

??
wasmer run --invoke add target/wasm32-wasi/debug/hey.wasm 23 19
```
