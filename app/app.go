package app

import (
	"fmt"

	"github.com/Unquabain/decider/list"
)

type App struct {
	UI
	List *list.Model
}

func (a App) Peek() (string, error) {
	return a.List.Peek()
}

func (a App) Add(task string) error {
	if task == `` {
		var err error
		task, err = a.UI.Prompt(`What is the new task you want to perform?`)
		if err != nil {
			return fmt.Errorf(`could not get task description: %w`, err)
		}
	}
	i, err := a.List.Push(task)
	if err != nil {
		return fmt.Errorf(`could not add task: %w`, err)
	}
	if err := a.UI.Decide(i); err != nil {
		return fmt.Errorf(`could not prioritize task: %w`, err)
	}
	return nil
}

func (a App) completePrompt() (bool, error) {
	task, err := a.Peek()
	if err != nil {
		return false, fmt.Errorf(`could not get current task: %w`, err)
	}
	return a.UI.Confirm(`Do you want to complete this task?`, task), nil
}

func (a App) Complete(prompt bool) error {
	if a.List.Len() == 0 {
		return nil
	}
	if prompt {
		cont, err := a.completePrompt()
		if err != nil || !cont {
			return err
		}
	}
	i, err := a.List.Pop()
	if err != nil {
		return fmt.Errorf(`could not pop task: %w`, err)
	}
	if err := a.UI.Decide(i); err != nil {
		return fmt.Errorf(`could not prioritize task: %w`, err)
	}
	return nil
}

func (a App) Resort() error {
	if a.List.Len() <= 1 {
		return nil
	}
	task, err := a.Peek()
	if err != nil {
		return fmt.Errorf(`couldn't save first task: %w`, err)
	}
	if err := a.Complete(false); err != nil {
		return fmt.Errorf(`couldn't remove first task: %w`, err)
	}
	if err := a.Add(task); err != nil {
		return fmt.Errorf(`couldn't add back task: %w`, err)
	}
	return nil
}

func (a App) ResortAll() error {
	tasks := a.List.Tasks()
	a.List = list.NewFromList(nil)
	for _, task := range tasks {
		if err := a.Add(task); err != nil {
			return fmt.Errorf(`could not resort all tasks: %w`, err)
		}
	}
	return nil
}

func (a App) ShowList() error {
	return a.UI.List(a.List.Tasks())
}
