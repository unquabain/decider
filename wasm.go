//go:build wasm
// +build wasm

package main

import (
	"errors"
	"fmt"
	"syscall/js"

	"github.com/Unquabain/decider/app"
	"github.com/Unquabain/decider/list"
	"github.com/charmbracelet/log"
)

func stringArrayToValue(ss []string) js.Value {
	aa := make([]any, len(ss))
	for i, s := range ss {
		aa[i] = s
	}
	return js.ValueOf(aa)
}

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

func promisify(f func(js.Value, []js.Value) (js.Value, error)) js.Func {
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

type iteratorWrapper struct {
	i list.Iterator
}

func (i iteratorWrapper) Value() js.Value {
	v := js.ValueOf(make(map[string]any))
	v.Set(`tasks`, promisify(i.Tasks))
	v.Set(`greatest`, promisify(i.Greatest))
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

type JSUI js.Value

func (ui JSUI) Value() js.Value {
	return js.Value(ui)
}

func (ui JSUI) isObject() bool {
	return ui.Value().Type() == js.TypeObject
}

func (ui JSUI) getFunc(funcname string) (js.Value, error) {
	v := ui.Value()
	if !ui.isObject() {
		return js.Undefined(), fmt.Errorf(`JS ui is not an object`)
	}
	f := v.Get(funcname)
	if f.Type() != js.TypeFunction {
		return js.Undefined(), fmt.Errorf(`JS ui.%s is not a function`, funcname)
	}
	return f, nil
}

func (ui JSUI) Decide(i list.Iterator) error {
	f, err := ui.getFunc(`decide`)
	if err != nil {
		return err
	}
	for i != nil {
		ret, err := Promise(f.Invoke(iteratorWrapper{i}.Value())).Wait()
		if err != nil {
			return err
		}
		if ret.Type() != js.TypeNumber {
			return fmt.Errorf(`unexpected return from decide function: %T, %+v`, ret, ret)
		}
		i, err = i.Greatest(ret.Int())
		if err != nil {
			return err
		}
	}
	return nil
}

func (ui JSUI) Confirm(prompt, task string) bool {
	f, err := ui.getFunc(`confirm`)
	if err != nil {
		log.With(`err`, err).Error(`could not find confirm function in JS`)
		return false
	}
	ret, err := Promise(f.Invoke(js.ValueOf(prompt), js.ValueOf(task))).Wait()
	if err != nil {
		log.With(`err`, err).Error(`could not get confirmation`)
		return false
	}
	if ret.Type() != js.TypeBoolean {
		log.With(`ret`, ret).Errorf(`Unexpected return from JS confirm: %T`, ret)
		return false
	}
	return ret.Bool()
}
func (ui JSUI) Prompt(prompt string) (string, error) {
	f, err := ui.getFunc(`prompt`)
	if err != nil {
		return ``, fmt.Errorf(`could not find prompt function in js: %w`, err)
	}
	ret, err := Promise(f.Invoke(js.ValueOf(prompt))).Wait()
	if err != nil {
		return ``, fmt.Errorf(`could not get prompt response: %w`, err)
	}
	if ret.Type() != js.TypeString {
		return ``, fmt.Errorf(`Unexpected return from JS prompt: %T`, ret)
	}
	return ret.String(), nil
}
func (ui JSUI) List(tasks []string) error {
	f, err := ui.getFunc(`list`)
	if err != nil {
		return fmt.Errorf(`could not find list function in js: %w`, err)
	}
	ret := f.Invoke(stringArrayToValue(tasks))
	if ret.Type() != js.TypeUndefined {
		return fmt.Errorf(`Unexpected return from JS list: %T`, ret)
	}
	return nil
}

type appWrapper struct {
	app.App
}

func (w appWrapper) Value() js.Value {
	v := js.ValueOf(make(map[string]any))
	v.Set(`peek`, promisify(w.Peek))
	v.Set(`add`, promisify(w.Add))
	v.Set(`complete`, promisify(w.Complete))
	v.Set(`resort`, promisify(w.Resort))
	v.Set(`resortAll`, promisify(w.ResortAll))
	v.Set(`list`, promisify(w.ShowList))
	v.Set(`tasks`, js.FuncOf(func(_ js.Value, _ []js.Value) any {
		return stringArrayToValue(w.App.List.Tasks())
	}))
	return v
}

func (w appWrapper) Peek(_ js.Value, _ []js.Value) (js.Value, error) {
	task, err := w.App.Peek()
	if errors.Is(err, list.ErrCaughtUp) {
		return js.ValueOf("All caught up!"), nil
	}

	return js.ValueOf(task), err
}

func (w appWrapper) Add(_ js.Value, args []js.Value) (js.Value, error) {
	task := ``
	if len(args) > 0 && args[0].Type() == js.TypeString {
		task = args[0].String()
	}
	return js.Undefined(), w.App.Add(task)
}

func (w appWrapper) Complete(_ js.Value, args []js.Value) (js.Value, error) {
	prompt := false
	if len(args) > 0 {
		prompt = args[0].Truthy()
	}
	return js.Undefined(), w.App.Complete(prompt)
}

func (w appWrapper) Resort(_ js.Value, _ []js.Value) (js.Value, error) {
	return js.Undefined(), w.App.Resort()
}

func (w appWrapper) ResortAll(_ js.Value, _ []js.Value) (js.Value, error) {
	return js.Undefined(), w.App.ResortAll()
}
func (w appWrapper) ShowList(_ js.Value, _ []js.Value) (js.Value, error) {
	return js.Undefined(), w.App.ShowList()
}

func wrapApp(_ js.Value, args []js.Value) (js.Value, error) {
	if len(args) < 2 {
		return js.Undefined(), fmt.Errorf(`expected 2 arguments: a list of tasks and a UI object`)
	}

	valTasks := args[0]
	valUI := args[1]

	if t := valTasks.Type(); t != js.TypeObject {
		return js.Undefined(), fmt.Errorf(`expected argument 1 to be an array of strings, got %v %v`, t, valTasks)
	}
	if t := valUI.Type(); t != js.TypeObject {
		return js.Undefined(), fmt.Errorf(`expected argument 2 to be a UI interface, got %v`, t)
	}
	tasks := make([]string, valTasks.Length())
	for i := range tasks {
		valTask := valTasks.Index(i)
		if t := valTask.Type(); t != js.TypeString {
			return js.Undefined(), fmt.Errorf(`expected argument 1 to be an array of strings, but element %d was a %v`, i, t)
		}
		tasks[i] = valTask.String()
	}
	model := list.NewFromList(tasks)
	ui := JSUI(valUI)
	app := app.App{
		List: model,
		UI:   ui,
	}
	return appWrapper{app}.Value(), nil
}

func main() {
	freeze := make(chan struct{})

	js.Global().Set(`newApp`, promisify(wrapApp))
	js.Global().Call(`init`)

	<-freeze // deadlock
}
