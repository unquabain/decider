//go:build server
// +build server

package server

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"
)

//go:embed fs/**/* fs/*
var embedFS embed.FS

var cacheHash []byte

func init() {
	cacheHash = []byte(uuid.New().String())
}

type cacheAgeMiddleware struct {
	http.Handler
}

func (m cacheAgeMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case `/`, `/index.html`, `/main.css`:
		recorder := httptest.NewRecorder()
		m.Handler.ServeHTTP(recorder, r)
		for key, val := range recorder.Header() {
			w.Header()[key] = val
		}
		b := recorder.Body.Bytes()
		b = bytes.ReplaceAll(b, []byte(`CACHEHASH`), cacheHash)
		w.Header().Set(`Content-Length`, fmt.Sprint(len(b)))
		w.WriteHeader(recorder.Code)
		w.Write(b)
		return
	}
	w.Header().Add(`Cache-Control`, `public, max-age=7776000, immutable`)
	m.Handler.ServeHTTP(w, r)
}

func ListenAndServe(port int) error {
	docs, err := fs.Sub(embedFS, `fs`)
	if err != nil {
		// This should all be worked out at compile time.
		panic(err)
	}
	http.Handle(`/`, cacheAgeMiddleware{http.FileServer(http.FS(docs))})
	err = http.ListenAndServe(fmt.Sprintf(`:%d`, port), nil)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
