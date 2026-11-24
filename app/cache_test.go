package app

import (
	"io"
	"testing"

	"fyne.io/fyne/v2/test"

	"github.com/stretchr/testify/assert"
)

func TestCache_Exists(t *testing.T) {
	c := NewWithID("io.fyne.test").Cache()
	testName := "entry"

	assert.False(t, c.Exists(testName))

	w, err := c.Write(testName)
	assert.NoError(t, err)
	n, _ := w.Write([]byte("data"))
	assert.Equal(t, 4, n)
	_ = w.Close()

	assert.True(t, c.Exists(testName))

	// and verify it was stored
	c = NewWithID("io.fyne.test").Cache()
	assert.True(t, c.Exists(testName))

	_ = c.Remove(testName)
}

func TestCache_ExistsTest(t *testing.T) {
	c := test.NewApp().Cache()
	testName := "entry"

	assert.False(t, c.Exists(testName))

	w, err := c.Write(testName)
	assert.NoError(t, err)
	n, _ := w.Write([]byte("data"))
	assert.Equal(t, 4, n)
	_ = w.Close()

	assert.True(t, c.Exists(testName))
	// no need to remove a test cache, it resets
}

func TestCache_Set(t *testing.T) {
	c := NewWithID("io.fyne.test").Cache()
	testName := "entry2"

	w, err := c.Write(testName)
	assert.NoError(t, err)
	n, _ := w.Write([]byte("data"))
	assert.Equal(t, 4, n)
	_ = w.Close()

	read, err := c.Read(testName)
	assert.NoError(t, err)
	data, err := io.ReadAll(read)
	assert.NoError(t, err)
	_ = read.Close()
	assert.NoError(t, err)
	assert.Equal(t, []byte("data"), data)

	// and verify it was stored
	c = NewWithID("io.fyne.test").Cache()

	read, err = c.Read(testName)
	assert.NoError(t, err)
	data, err = io.ReadAll(read)
	assert.NoError(t, err)
	_ = read.Close()
	assert.NoError(t, err)
	assert.Equal(t, []byte("data"), data)

	_ = c.Remove(testName)
}

func TestCache_SetTest(t *testing.T) {
	c := test.NewApp().Cache()
	testName := "entry2"

	w, err := c.Write(testName)
	assert.NoError(t, err)
	n, _ := w.Write([]byte("data"))
	assert.Equal(t, 4, n)
	_ = w.Close()

	read, err := c.Read(testName)
	assert.NoError(t, err)
	data, err := io.ReadAll(read)
	assert.NoError(t, err)
	_ = read.Close()
	assert.NoError(t, err)
	assert.Equal(t, []byte("data"), data)
}

func TestCache_URLName(t *testing.T) {
	c := NewWithID("io.fyne.test").Cache()
	testName := "https://example.com/icons.png?size=64" // standard base64 of this contains "/"

	_, err := c.Read(testName)
	assert.Error(t, err)
	assert.False(t, c.Exists(testName))

	w, err := c.Write(testName)
	assert.NoError(t, err)
	n, _ := w.Write([]byte("data"))
	assert.Equal(t, 4, n)
	_ = w.Close()

	assert.True(t, c.Exists(testName))
	read, err := c.Read(testName)
	assert.NoError(t, err)
	data, err := io.ReadAll(read)
	_ = read.Close()
	assert.NoError(t, err)
	assert.Equal(t, []byte("data"), data)

	_ = c.Remove(testName)
	assert.False(t, c.Exists(testName))
}
