//go:build !wasm
// +build !wasm

package list

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"os"
)

func New(filename string) *Model {
	return &Model{
		Filename: filename,
		tasks:    nil,
	}
}

func (m *Model) Open() error {
	f, err := os.OpenFile(m.Filename, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf(`could not open %q: %w`, m.Filename, err)
	}
	defer f.Close()
	if err := gob.NewDecoder(f).Decode(&m.tasks); err != nil {
		if errors.Is(err, os.ErrExist) {
			return nil
		}
		if errors.Is(err, io.EOF) {
			return nil
		}
		return fmt.Errorf(`could not decode %q: %w`, m.Filename, err)
	}
	return nil
}

func (m *Model) Save() error {
	f, err := os.OpenFile(m.Filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf(`could not open %q: %w`, m.Filename, err)
	}
	defer f.Close()
	if err := gob.NewEncoder(f).Encode(m.tasks); err != nil {
		return fmt.Errorf(`could not encode %q: %w`, m.Filename, err)
	}
	return nil
}
