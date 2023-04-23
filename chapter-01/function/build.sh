#!/bin/bash
tinygo build -o add.wasm \
  -scheduler=none --no-debug \
  -target wasi ./add.go

ls -lh *.wasm
