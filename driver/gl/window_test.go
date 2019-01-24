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
