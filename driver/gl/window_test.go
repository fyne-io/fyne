// +build !ci

package gl

import (
	"testing"

	_ "fyne.io/fyne/test"

	"github.com/stretchr/testify/assert"
)

func TestWindow_SetTitle(t *testing.T) {
	d := NewGLDriver()
	w := d.CreateWindow("Test")

	title := "My title"
	w.SetTitle(title)

	assert.Equal(t, title, w.Title())
}

func TestWindow_Clipboard(t *testing.T) {
	d := NewGLDriver()
	w := d.CreateWindow("Test")

	text := "My content from test window"
	cb := w.Clipboard()
	assert.Empty(t, cb.Content())

	cb.SetContent(text)
	assert.Equal(t, text, cb.Content())
}
