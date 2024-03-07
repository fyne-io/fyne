//go:build !wasm

package app

import (
	"io"
	"os"
	"path/filepath"
)

func (p *preferences) storageWriter() (writeSyncCloser, error) {
	return p.storageWriterForPath(p.storagePath())
}

func (p *preferences) storageReader() (io.ReadCloser, error) {
	return p.storageReaderForPath(p.storagePath())
}

func (p *preferences) storageWriterForPath(path string) (writeSyncCloser, error) {
	err := os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil { // this is not an exists error according to docs
		return nil, err
	}
	file, err := os.Create(path)
	if err != nil {
		if !os.IsExist(err) {
			return nil, err
		}
		file, err = os.Open(path) // #nosec
		if err != nil {
			return nil, err
		}
	}
	return file, nil
}

func (p *preferences) storageReaderForPath(path string) (io.ReadCloser, error) {
	file, err := os.Open(path) // #nosec
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
				return nil, err
			}
			return nil, errEmptyPreferencesStore
		}
		return nil, err
	}
	return file, nil
}

// the following are only used in tests to save preferences to a tmp file

func (p *preferences) saveToFile(path string) error {
	file, err := p.storageWriterForPath(path)
	if err != nil {
		return err
	}
	return p.saveToStorage(file)
}

func (p *preferences) loadFromFile(path string) error {
	file, err := p.storageReaderForPath(path)
	if err != nil {
		return err
	}
	return p.loadFromStorage(file)
}
