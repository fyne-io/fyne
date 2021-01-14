package repository

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"fyne.io/fyne/storage"
	"fyne.io/fyne/storage/repository"

	"github.com/stretchr/testify/assert"
)

func TestFileRepositoryRegistration(t *testing.T) {
	f := NewFileRepository("file")
	repository.Register("file", f)

	// this should never fail, and we assume it doesn't in other tests here
	// for brevity
	foo, err := storage.ParseURI("file:///foo")
	assert.Nil(t, err)

	// make sure we get the same repo back
	repo, err := repository.ForURI(foo)
	assert.Nil(t, err)
	assert.Equal(t, f, repo)
}

func TestFileRepositoryExists(t *testing.T) {
	dir, err := ioutil.TempDir("", "FyneInternalRepositoryFileTest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	existsPath := path.Join(dir, "exists")
	notExistsPath := path.Join(dir, "notExists")

	err = ioutil.WriteFile(existsPath, []byte{1, 2, 3, 4}, 0755)
	if err != nil {
		t.Fatal(err)
	}

	ex, err := storage.Exists(storage.NewFileURI(existsPath))
	assert.Nil(t, err)
	assert.True(t, ex)

	ex, err = storage.Exists(storage.NewFileURI(notExistsPath))
	assert.Nil(t, err)
	assert.False(t, ex)
}

func TestFileRepositoryReader(t *testing.T) {
	// Set up a temporary directory.
	dir, err := ioutil.TempDir("", "FyneInternalRepositoryFileTest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// Create some files to tests with.
	fooPath := path.Join(dir, "foo")
	barPath := path.Join(dir, "bar")
	bazPath := path.Join(dir, "baz")
	err = ioutil.WriteFile(fooPath, []byte{}, 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(barPath, []byte{1, 2, 3}, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Set up our repository - it's OK if we already registered it...
	f := NewFileRepository("file")
	repository.Register("file", f)

	// ...and some URIs - we know that they will not fail parsing
	foo := storage.NewFileURI(fooPath)
	bar := storage.NewFileURI(barPath)
	baz := storage.NewFileURI(bazPath)

	// Make sure we can read the empty file.
	fooReader, err := storage.Reader(foo)
	assert.Nil(t, err)
	fooData, err := ioutil.ReadAll(fooReader)
	assert.Equal(t, []byte{}, fooData)
	assert.Nil(t, err)

	// Make sure we can read the file with data.
	barReader, err := storage.Reader(bar)
	assert.Nil(t, err)
	barData, err := ioutil.ReadAll(barReader)
	assert.Equal(t, []byte{1, 2, 3}, barData)
	assert.Nil(t, err)

	// Make sure we get an error if the file doesn't exist.
	_, err = storage.Reader(baz)
	assert.NotNil(t, err)
}
