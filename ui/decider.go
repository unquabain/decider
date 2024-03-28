package ui

import (
	"github.com/Unquabain/decider/list"
	"github.com/charmbracelet/huh"
)

type Decider struct {
	s *huh.Select[int]
	i list.Iterator
}

func NewDecider(i list.Iterator) *Decider {
	return &Decider{
		i: i,
		s: huh.NewSelect[int]().
			Title(`Select the most urgent of these tasks.`).
			Description(`You will be asked the relative urgency of a subset of available tasks in order to find their proper, semi-sorted position`),
	}
}

func (d *Decider) Run() error {
	for d.i != nil {
		options := make([]huh.Option[int], 0, 3)
		for idx, task := range d.i.Tasks() {
			options = append(options, huh.Option[int]{Value: idx, Key: task})
		}
		d.s.Options(options...)
		if err := d.s.Run(); err != nil {
			return err
		}
		idx := d.s.GetValue().(int)
		i, err := d.i.Greatest(idx)
		if err != nil {
			return err
		}
		d.i = i
	}
	return nil
}
