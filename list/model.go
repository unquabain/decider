package list

import (
	"errors"
	"fmt"
)

var ErrCaughtUp = errors.New(`all caught up`)

type Model struct {
	Filename string
	tasks    []string
}

func NewFromList(tasks []string) *Model {
	return &Model{
		tasks: tasks,
	}
}

func (m *Model) Tasks() []string {
	return m.tasks
}

type index int

func (i index) parent() index {
	return (i - 1) / 2
}
func (i index) leftChild() index {
	return 2*i + 1
}
func (i index) rightChild() index {
	return 2*i + 2
}

func (m *Model) Len() int {
	return len(m.tasks)
}

func (m *Model) Swap(i, j int) error {
	if i < 0 || i >= m.Len() {
		return fmt.Errorf(`invalid first index: %d`, i)
	}
	if j < 0 || j >= m.Len() {
		return fmt.Errorf(`invalid second index: %d`, j)
	}
	m.tasks[i], m.tasks[j] = m.tasks[j], m.tasks[i]
	return nil
}

func (m *Model) Peek() (string, error) {
	if m.Len() == 0 {
		return ``, ErrCaughtUp
	}
	return m.tasks[0], nil
}

func (m *Model) Pop() (Iterator, error) {
	if m.Len() == 0 {
		return nil, ErrCaughtUp
	}
	if m.Len() == 1 {
		m.tasks = nil
		return nil, nil
	}
	last := m.Len() - 1
	if last == 1 {
		m.tasks = m.tasks[1:]
		return nil, nil
	}
	if err := m.Swap(0, last); err != nil {
		return nil, fmt.Errorf(`could not promote last task: %w`, err)
	}
	m.tasks = m.tasks[:last]

	return &floatUpIterator{
		m: m,
		i: 0,
	}, nil
}

func (m *Model) Push(task string) (Iterator, error) {
	if m.Len() == 0 {
		m.tasks = []string{task}
		return nil, nil
	}

	last := m.Len()
	m.tasks = append(m.tasks, task)
	return &sinkDownIterator{
		m: m,
		i: index(last),
	}, nil
}

type Iterator interface {
	Tasks() []string
	Greatest(int) (Iterator, error)
}

type floatUpIterator struct {
	m *Model
	i index
}

func (i *floatUpIterator) Tasks() []string {
	tasks := make([]string, 2, 3)
	tasks[0] = i.m.tasks[i.i]
	tasks[1] = i.m.tasks[i.i.leftChild()]
	if rc := i.i.rightChild(); int(rc) < i.m.Len() {
		tasks = append(tasks, i.m.tasks[rc])
	}
	return tasks
}

func (i *floatUpIterator) Greatest(idx int) (Iterator, error) {
	var child index
	switch idx {
	case 0:
		return nil, nil
	case 1:
		child = i.i.leftChild()
	case 2:
		child = i.i.rightChild()
	default:
		return nil, fmt.Errorf(`invalid index %d`, idx)
	}
	if int(child) >= i.m.Len() {
		return nil, fmt.Errorf(`invalid index %d`, idx)
	}
	if err := i.m.Swap(int(i.i), int(child)); err != nil {
		return nil, fmt.Errorf(`could not float up: %w`, err)
	}
	i.i = child
	if int(i.i.leftChild()) >= i.m.Len() {
		return nil, nil
	}
	return i, nil

}

type sinkDownIterator struct {
	m *Model
	i index
}

func (i *sinkDownIterator) Tasks() []string {
	return []string{
		i.m.tasks[i.i],
		i.m.tasks[i.i.parent()],
	}
}

func (i *sinkDownIterator) Greatest(idx int) (Iterator, error) {
	switch idx {
	case 0:
		if err := i.m.Swap(int(i.i), int(i.i.parent())); err != nil {
			return nil, fmt.Errorf(`could not float up: %w`, err)
		}
		i.i = i.i.parent()
		if i.i == 0 {
			return nil, nil
		}
		return i, nil
	case 1:
		return nil, nil
	default:
		return nil, fmt.Errorf(`invalid index: %d`, idx)
	}
}
