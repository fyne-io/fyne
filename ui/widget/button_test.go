package widget

import "testing"

import "github.com/stretchr/testify/assert"

import _ "github.com/fyne-io/fyne/test"
import "github.com/fyne-io/fyne/ui/theme"

func TestButtonTestSize(t *testing.T) {
	button := NewButton("Hi", nil)
	min := button.MinSize()

	assert.True(t, min.Width >= theme.Padding()*2)
	assert.True(t, min.Height >= theme.Padding()*2)
}

func TestButtonTestNotify(t *testing.T) {
	tapped := false
	button := NewButton("Hi", func() {
		tapped = true
	})
	button.OnMouseDown(nil)

	assert.True(t, tapped)
}
