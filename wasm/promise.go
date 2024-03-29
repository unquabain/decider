//go:build wasm
// +build wasm

package wasm

import (
	"fmt"
	"syscall/js"
)

type Promise js.Value

func (p Promise) Value() js.Value {
	return js.Value(p)
}
func (p Promise) Then() (<-chan js.Value, error) {
	f := p.Value().Get(`then`)
	if t := f.Type(); t != js.TypeFunction {
		return nil, fmt.Errorf(`expected then to be a function: %v`, t)
	}

	c := make(chan js.Value)

	thenfunc := func(_ js.Value, args []js.Value) any {
		if len(args) == 0 {
			c <- js.Undefined()
			return nil
		}
		c <- args[0]
		return nil
	}
	p.Value().Call(`then`, js.FuncOf(thenfunc))

	return c, nil
}

func (p Promise) Catch() (<-chan js.Value, error) {
	f := p.Value().Get(`catch`)
	if t := f.Type(); t != js.TypeFunction {
		return nil, fmt.Errorf(`expected catch to be a function: %v`, t)
	}

	c := make(chan js.Value)

	thenfunc := func(_ js.Value, args []js.Value) any {
		if len(args) == 0 {
			c <- js.Undefined()
			return nil
		}
		c <- args[0]
		return nil
	}
	p.Value().Call(`catch`, js.FuncOf(thenfunc))

	return c, nil
}

func (prom Promise) Wait() (js.Value, error) {
	then, err := prom.Then()
	if err != nil {
		return js.Undefined(), err
	}
	catch, err := prom.Catch()
	if err != nil {
		return js.Undefined(), err
	}
	select {
	case jsErr := <-catch:
		return js.Undefined(), fmt.Errorf(`catching JS error: %v`, jsErr)
	case ret := <-then:
		return ret, nil
	}
}

func Promisify(f func(js.Value, []js.Value) (js.Value, error)) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		handler := js.FuncOf(func(_ js.Value, promArgs []js.Value) any {
			resolve := promArgs[0]
			reject := promArgs[1]

			go func() {
				data, err := f(this, args)
				if err != nil {
					// err should be an instance of `error`, eg `errors.New("some error")`
					errorConstructor := js.Global().Get("Error")
					errorObject := errorConstructor.New(err.Error())
					reject.Invoke(errorObject)
				} else {
					resolve.Invoke(data)
				}
			}()

			return nil
		})
		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}
