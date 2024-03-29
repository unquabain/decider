//go:build wasm
// +build wasm

package wasm

import (
	"errors"
	"fmt"
	"syscall/js"

	"github.com/Unquabain/decider/app"
	"github.com/Unquabain/decider/list"
)

type appWrapper struct {
	app.App
}

func (w appWrapper) Value() js.Value {
	v := js.ValueOf(make(map[string]any))
	v.Set(`peek`, Promisify(w.Peek))
	v.Set(`add`, Promisify(w.Add))
	v.Set(`complete`, Promisify(w.Complete))
	v.Set(`resort`, Promisify(w.Resort))
	v.Set(`resortAll`, Promisify(w.ResortAll))
	v.Set(`list`, Promisify(w.ShowList))
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

func WrapApp(_ js.Value, args []js.Value) (js.Value, error) {
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
