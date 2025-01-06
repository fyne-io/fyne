package mobile

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/storage/repository"
)

func TestFileRepositoryChild(t *testing.T) {
	f := &mobileFileRepo{}
	repository.Register("file", f)

	p, _ := storage.Child(storage.NewFileURI("/foo/bar"), "baz")
	assert.Equal(t, "file:///foo/bar/baz", p.String())

	p, _ = storage.Child(storage.NewFileURI("/foo/bar/"), "baz")
	assert.Equal(t, "file:///foo/bar/baz", p.String())

	uri, _ := storage.ParseURI("content://thing")
	p, err := storage.Child(uri, "new")
	require.Error(t, err)
	assert.Nil(t, p)
}

func TestFileRepositoryParent(t *testing.T) {
	f := &mobileFileRepo{}
	repository.Register("file", f)

	p, _ := storage.Parent(storage.NewFileURI("/foo/bar"))
	assert.Equal(t, "file:///foo", p.String())

	p, _ = storage.Parent(storage.NewFileURI("/foo/bar/"))
	assert.Equal(t, "file:///foo", p.String())

	uri, _ := storage.ParseURI("content://thing")
	p, err := storage.Parent(uri)
	require.Error(t, err)
	assert.Nil(t, p)
}
