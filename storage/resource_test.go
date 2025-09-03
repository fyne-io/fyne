package storage

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadResourceFromURI(t *testing.T) {
	path := filepath.Join("testdata", "test.txt")
	abs, err := filepath.Abs(path)
	assert.NoError(t, err)

	uri := NewFileURI(abs)
	res, err := LoadResourceFromURI(uri)
	assert.NoError(t, err)
	assert.Equal(t, "test.txt", res.Name())
	assert.Equal(t, []byte("Text content"), res.Content())
}
