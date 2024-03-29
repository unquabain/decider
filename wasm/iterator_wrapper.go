//go:build wasm
// +build wasm

package wasm

import (
	"fmt"
	"syscall/js"

	"github.com/Unquabain/decider/list"
)

type iteratorWrapper struct {
	i list.Iterator
}

func (i iteratorWrapper) Value() js.Value {
	v := js.ValueOf(make(map[string]any))
	v.Set(`tasks`, Promisify(i.Tasks))
	v.Set(`greatest`, Promisify(i.Greatest))
	return v
}

func (i iteratorWrapper) Tasks(_ js.Value, _ []js.Value) (js.Value, error) {
	if i.i == nil {
		return js.Undefined(), nil
	}
	return stringArrayToValue(i.i.Tasks()), nil
}

func (i iteratorWrapper) Greatest(_ js.Value, args []js.Value) (js.Value, error) {
	if len(args) < 1 {
		return js.Undefined(), fmt.Errorf(`not enough args to Greatest: Expected 1`)
	}
	vIdx := args[0]
	if t := vIdx.Type(); t != js.TypeNumber {
		return js.Undefined(), fmt.Errorf(`argument to Greatest should be a number. Got %v`, t)
	}

	idx := vIdx.Int()
	iter, err := i.i.Greatest(idx)
	if err != nil {
		return js.Undefined(), err
	}
	i.i = iter
	return i.Value(), nil
}
