#!/bin/bash
wasmedge --reactor add.wasm add 18 24

wasmer add.wasm --invoke add 30 12
