//go:build wasm

package app

import (
	"errors"
	"io"
)

func (p *preferences) storageReader() (io.ReadCloser, error) {
	return nil, errors.New("not implemented")
}

func (p *preferences) storageWriter() (writeSyncCloser, error) {
	return nil, errors.New("not implemented")
}
