package fyne

import "os"

import "testing"
import "github.com/stretchr/testify/assert"

func TestNewResource(t *testing.T) {
	name := "test"
	content := []byte{1, 2, 3, 4}

	res := NewResource(name, content)
	assert.Equal(t, name, res.Name)
	assert.Equal(t, content, res.Content)
}

func TestResourceFile(t *testing.T) {
	content := []byte{1, 2, 3, 4}

	res := NewResource("file.dat", content)
	path := res.CachePath()

	assert.True(t, path != "")
	assert.True(t, pathExists(path))

	os.Remove(path)
}
