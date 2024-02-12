//go:build wasm

package app

import (
	"errors"
	"io"
)

func (a *fyneApp) storageRoot() string {
	return "" // no storage root for web driver yet
}

func (p *preferences) storageReader() (io.ReadCloser, error) {
	return nil, errors.New("not implemented")
}

func (p *preferences) storageWriter() (writeSyncCloser, error) {
	return nil, errors.New("not implemented")
}

func (p *preferences) watch() {
	// no-op for web driver
}
