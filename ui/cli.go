//go:build !wasm
// +build !wasm

package ui

import (
	"fmt"

	"github.com/Unquabain/decider/list"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
)

type CLI struct{}

func (CLI) Decide(i list.Iterator) error {
	d := NewCLIDecider(i)
	return d.Run()
}

func (CLI) Confirm(prompt, task string) bool {
	ctl := huh.NewConfirm().Title(task).Description(`Do you want to complete this task?`)
	if err := ctl.Run(); err != nil {
		log.With(`err`, err).Error(`could not get user confirmation`)
		return false
	}
	return ctl.GetValue().(bool)
}

func (CLI) Prompt(prompt string) (string, error) {
	input := huh.NewInput().Title(prompt).Prompt(`> `)
	if err := input.Run(); err != nil {
		return ``, err
	}
	return input.GetValue().(string), nil
}

func (CLI) List(items []string) error {
	for i, item := range items {
		fmt.Printf("% 2d: %s\n", i+1, item)
	}
	return nil
}
