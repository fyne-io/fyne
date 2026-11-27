//go:build wasm

package repository

import (
	"bytes"
	"io"
	"testing"

	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/storage/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndexDBRepository_ReadLargeFile(t *testing.T) {
	repo, err := NewIndexDBRepository()
	require.NoError(t, err)
	repository.Register("idbfile", repo)

	u, err := storage.ParseURI("idbfile:///fyne/test/large.bin")
	require.NoError(t, err)
	t.Cleanup(func() { _ = storage.Delete(u) })

	// more than a reader will ask for in one go
	want := bytes.Repeat([]byte("0123456789abcdef"), 20000)
	w, err := storage.Writer(u)
	require.NoError(t, err)
	_, err = w.Write(want)
	require.NoError(t, err)
	require.NoError(t, w.Close())

	r, err := storage.Reader(u)
	require.NoError(t, err)
	got, err := io.ReadAll(r)
	require.NoError(t, err)
	require.NoError(t, r.Close())
	assert.True(t, bytes.Equal(want, got), "read %d bytes of %d", len(got), len(want))
}
