//go:build wasm
// +build wasm

package main

import (
	"syscall/js"

	"github.com/Unquabain/decider/wasm"
)

func main() {
	freeze := make(chan struct{})

	js.Global().Set(`newApp`, wasm.Promisify(wasm.WrapApp))
	js.Global().Call(`init`)

	<-freeze // deadlock
}
