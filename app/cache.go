package app

import (
	"encoding/base64"
	"io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

type cache struct {
	a fyne.App

	enc base64.Encoding
}

func makeCache(a fyne.App) fyne.Cache {
	c := &cache{a: a, enc: base64.Encoding{}}

	root := c.RootURI()
	exists, err := storage.Exists(root)
	if !exists || err != nil {
		parent, _ := storage.Parent(root)
		exists, err = storage.Exists(parent)
		if !exists || err != nil {
			err = storage.CreateListable(parent)
			if err != nil {
				fyne.LogError("Failed to create fyne cache space", err)
			}
		}

		err = storage.CreateListable(root)
		if err != nil {
			fyne.LogError("Failed to create app cache space", err)
		}
	}

	return c
}

func (c *cache) RootURI() fyne.URI {
	return rootCacheURI(c.a)
}

func (c *cache) Exists(name string) bool {
	path := c.encodePath(name)

	ok, err := storage.Exists(path)
	return ok && err == nil
}

func (c *cache) Read(name string) (io.ReadCloser, error) {
	path := c.encodePath(name)

	return storage.Reader(path)
}

func (c *cache) Write(name string) (io.WriteCloser, error) {
	path := c.encodePath(name)

	return storage.Writer(path)
}

func (c *cache) Remove(name string) error {
	path := c.encodePath(name)

	return storage.Delete(path)
}

func (c *cache) encodePath(badName string) fyne.URI {
	// URL encoding is used as the standard alphabet contains "/", which is a path separator
	name := base64.URLEncoding.EncodeToString([]byte(badName))
	child, _ := storage.Child(c.a.Cache().RootURI(), name)
	return child
}
