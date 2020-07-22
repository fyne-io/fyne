// +build !ci !windows

package glfw

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne"
	"fyne.io/fyne/storage"
)

func TestGLDriver_FileReaderForURI(t *testing.T) {
	uri, _ := testURI("text.txt")
	read, err := NewGLDriver().FileReaderForURI(uri)
	defer read.Close()

	assert.Nil(t, err)
	assert.Equal(t, uri, read.URI())

	data, err := ioutil.ReadAll(read)
	assert.Nil(t, err)
	assert.Equal(t, "content", string(data))
}

func TestGLDriver_FileReaderForURI_Again(t *testing.T) { // a bug blanked files
	TestGLDriver_FileReaderForURI(t)
}

func TestGLDriver_FileWriterForURI(t *testing.T) {
	uri, path := testURI("text2.txt")
	defer os.Remove(path)
	write, err := NewGLDriver().FileWriterForURI(uri)
	defer write.Close()

	assert.Nil(t, err)
	assert.Equal(t, uri, write.URI())

	count, err := write.Write([]byte("another"))
	assert.Equal(t, 7, count)
	assert.Nil(t, err)
}

func testURI(file string) (fyne.URI, string) {
	path, _ := filepath.Abs(filepath.Join("testdata", file))
	uri := storage.NewURI("file://" + path)

	return uri, path
}
