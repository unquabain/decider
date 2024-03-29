//go:build !wasm
// +build !wasm

package main

//go:generate sh -c "GOOS=js GOARCH=wasm go build -tags wasm -o server/fs/lib/decider.wasm"
//go:generate sh -c "cp $(go env GOROOT)/misc/wasm/wasm_exec.js server/fs/lib/"
