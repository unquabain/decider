package server

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
)

//go:embed fs/**/* fs/*
var embedFS embed.FS

func ListenAndServe(port int) error {
	docs, err := fs.Sub(embedFS, `fs`)
	if err != nil {
		// This should all be worked out at compile time.
		panic(err)
	}
	http.Handle(`/`, http.FileServer(http.FS(docs)))
	err = http.ListenAndServe(fmt.Sprintf(`:%d`, port), nil)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
