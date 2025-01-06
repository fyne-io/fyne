package internal

import (
	"testing"

	intRepo "fyne.io/fyne/v2/internal/repository"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/storage/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDocs_Create(t *testing.T) {
	repository.Register("file", intRepo.NewInMemoryRepository("file"))
	docs := &Docs{storage.NewFileURI("/tmp/docs/create")}
	exist, err := storage.Exists(docs.RootDocURI)
	require.NoError(t, err)
	assert.False(t, exist)

	w, err := docs.Create("test")
	require.NoError(t, err)
	_ = w.Close()

	exist, err = storage.Exists(docs.RootDocURI)
	require.NoError(t, err)
	assert.True(t, exist)
}

func TestDocs_Create_ErrAlreadyExists(t *testing.T) {
	repository.Register("file", intRepo.NewInMemoryRepository("file"))
	docs := &Docs{storage.NewFileURI("/tmp/docs/create")}
	exist, err := storage.Exists(docs.RootDocURI)
	require.NoError(t, err)
	assert.False(t, exist)

	w, err := docs.Create("test")
	require.NoError(t, err)
	_, _ = w.Write([]byte{})
	_ = w.Close()
	_, err = docs.Create("test")
	assert.Equal(t, storage.ErrAlreadyExists, err)
}

func TestDocs_List(t *testing.T) {
	repository.Register("file", intRepo.NewInMemoryRepository("file"))
	docs := &Docs{storage.NewFileURI("/tmp/docs/list")}
	list := docs.List()
	assert.Empty(t, list)

	w, _ := docs.Create("1")
	_, _ = w.Write([]byte{})
	_ = w.Close()
	w, _ = docs.Create("2")
	_, _ = w.Write([]byte{})
	_ = w.Close()

	list = docs.List()
	assert.Len(t, list, 2)
}

func TestDocs_Save(t *testing.T) {
	r := intRepo.NewInMemoryRepository("file")
	repository.Register("file", r)
	docs := &Docs{storage.NewFileURI("/tmp/docs/save")}
	w, err := docs.Create("save.txt")
	require.NoError(t, err)
	_, _ = w.Write([]byte{})
	_ = w.Close()
	u := w.URI()

	exist, err := r.Exists(u)
	require.NoError(t, err)
	assert.True(t, exist)

	w, err = docs.Save("save.txt")
	require.NoError(t, err)
	n, err := w.Write([]byte("save"))
	require.NoError(t, err)
	assert.Equal(t, 4, n)
	err = w.Close()
	require.NoError(t, err)
}

func TestDocs_Save_ErrNotExists(t *testing.T) {
	r := intRepo.NewInMemoryRepository("file")
	repository.Register("file", r)
	docs := &Docs{storage.NewFileURI("/tmp/docs/save")}

	_, err := docs.Save("save.txt")
	assert.Equal(t, storage.ErrNotExists, err)
}
