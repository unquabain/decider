//go:build wasm
// +build wasm

package wasm

import (
	"fmt"
	"syscall/js"

	"github.com/Unquabain/decider/list"
	"github.com/charmbracelet/log"
)

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
