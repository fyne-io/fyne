package gomobile

import (
	"testing"

	"github.com/stretchr/testify/assert"

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

	p, err := storage.Child(storage.NewURI("content://thing"), "new")
	assert.NotNil(t, err)
	assert.Nil(t, p)
}

func TestFileRepositoryParent(t *testing.T) {
	f := &mobileFileRepo{}
	repository.Register("file", f)

	p, _ := storage.Parent(storage.NewFileURI("/foo/bar"))
	assert.Equal(t, "file:///foo", p.String())

	p, _ = storage.Parent(storage.NewFileURI("/foo/bar/"))
	assert.Equal(t, "file:///foo", p.String())

	p, err := storage.Parent(storage.NewURI("content://thing"))
	assert.NotNil(t, err)
	assert.Nil(t, p)
}
