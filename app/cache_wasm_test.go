//go:build wasm

package app

import (
	"io"
	"syscall/js"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRootCacheURIWasm(t *testing.T) {
	uri := rootCacheURI(nil)

	assert.Equal(t, "idbfile", uri.Scheme())
	assert.Equal(t, "/fyne/cache", uri.Path())
}

func TestCacheWasm(t *testing.T) {
	global := js.Global()
	if global.Get("indexedDB").IsUndefined() || global.Get("localStorage").IsUndefined() {
		t.Skip("browser storage APIs are not available")
	}

	c := NewWithID("io.fyne.cache.wasm.test").Cache()
	require.Equal(t, "idbfile", c.RootURI().Scheme())

	w, err := c.Write("entry")
	require.NoError(t, err)
	_, err = w.Write([]byte("data"))
	require.NoError(t, err)
	require.NoError(t, w.Close())

	assert.True(t, c.Exists("entry"))

	r, err := c.Read("entry")
	require.NoError(t, err)
	data, err := io.ReadAll(r)
	require.NoError(t, err)
	require.NoError(t, r.Close())
	assert.Equal(t, []byte("data"), data)

	assert.NoError(t, c.Remove("entry"))
	assert.False(t, c.Exists("entry"))
}
