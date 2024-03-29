//go:build wasm
// +build wasm

package wasm

import "syscall/js"

func stringArrayToValue(ss []string) js.Value {
	aa := make([]any, len(ss))
	for i, s := range ss {
		aa[i] = s
	}
	return js.ValueOf(aa)
}
