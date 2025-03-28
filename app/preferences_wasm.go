//go:build wasm

package app

import (
	"bytes"
	"io"
	"strings"
	"syscall/js"
)

const preferencesLocalStorageKey = "fyne-preferences.json"

func (a *fyneApp) storageRoot() string {
	return "idbfile:///fyne/"
}

func (p *preferences) storageReader() (io.ReadCloser, error) {
	key := js.ValueOf(preferencesLocalStorageKey)
	data := js.Global().Get("localStorage").Call("getItem", key)
	if data.IsNull() || data.IsUndefined() {
		return nil, errEmptyPreferencesStore
	}

	return readerNopCloser{reader: strings.NewReader(data.String())}, nil
}

func (p *preferences) storageWriter() (writeSyncCloser, error) {
	return &localStorageWriter{key: preferencesLocalStorageKey}, nil
}

func (p *preferences) watch() {
	// no-op for web driver
}

type readerNopCloser struct {
	reader io.Reader
}

func (r readerNopCloser) Read(b []byte) (int, error) {
	return r.reader.Read(b)
}

func (r readerNopCloser) Close() error {
	return nil
}

type localStorageWriter struct {
	bytes.Buffer
	key string
}

func (s *localStorageWriter) Sync() error {
	text := s.String()
	s.Reset()
	js.Global().Get("localStorage").Call("setItem", js.ValueOf(s.key), js.ValueOf(text))
	return nil
}

func (s *localStorageWriter) Close() error {
	return nil
}
