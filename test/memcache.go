package test

import (
	"bytes"
	"errors"
	"io"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

type memCache struct {
	memStore map[string][]byte
}

func makeCache() fyne.Cache {
	return &memCache{memStore: make(map[string][]byte)}
}

func (*memCache) RootURI() fyne.URI {
	return storage.NewFileURI(os.TempDir()) // in case anyone wants to manually handle storage
}

func (c *memCache) Exists(name string) bool {
	_, ok := c.memStore[name]
	return ok
}

func (c *memCache) Read(name string) (io.ReadCloser, error) {
	data, ok := c.memStore[name]
	if !ok {
		return nil, errors.New("not found")
	}

	return io.NopCloser(bytes.NewReader(data)), nil
}

func (c *memCache) Write(name string) (io.WriteCloser, error) {
	write := &bytes.Buffer{}

	return writeCloser{write, func() {
		c.memStore[name] = write.Bytes()
	}}, nil
}

func (c *memCache) Remove(name string) error {
	delete(c.memStore, name)
	return nil
}

type writeCloser struct {
	io.Writer

	onClose func()
}

func (c writeCloser) Close() error {
	c.onClose()
	return nil
}
