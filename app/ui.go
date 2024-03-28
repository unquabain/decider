package app

import "github.com/Unquabain/decider/list"

type UI interface {
	Decide(list.Iterator) error
	Confirm(prompt, task string) bool
	Prompt(prompt string) (string, error)
	List([]string) error
}
