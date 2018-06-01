package widget

import "testing"

import "github.com/stretchr/testify/assert"

import _ "github.com/fyne-io/fyne/test"
import "github.com/fyne-io/fyne/api/ui/canvas"
import "github.com/fyne-io/fyne/api/ui/theme"

func TestButtonSize(t *testing.T) {
	button := NewButton("Hi", nil)
	min := button.MinSize()

	assert.True(t, min.Width > theme.Padding()*2)
	assert.True(t, min.Height > theme.Padding()*2)
}

func TestButtonType(t *testing.T) {
	button := NewButton("Hi", nil)
	bg := button.Layout(button.MinSize())[0].(*canvas.Rectangle)
	color := bg.FillColor

	button.Style = PrimaryButton
	button.ApplyTheme()
	assert.NotEqual(t, bg.FillColor, color)
}

func TestButtonNotify(t *testing.T) {
	tapped := false
	button := NewButton("Hi", func() {
		tapped = true
	})
	button.OnMouseDown(nil)

	assert.True(t, tapped)
}
