// +build !ci

package gl

import (
	"testing"

	"fyne.io/fyne/canvas"

	"github.com/stretchr/testify/assert"
)

func TestGlCanvas_Content(t *testing.T) {
	content := &canvas.Circle{}
	w := d.CreateWindow("Test")
	w.SetContent(content)

	assert.Equal(t, content, w.Content())
}

func TestGlCanvas_NilContent(t *testing.T) {
	w := d.CreateWindow("Test")

	assert.NotNil(t, w.Content()) // never a nil canvas so we have a sensible fallback
}
