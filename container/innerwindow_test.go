package container

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
)

func TestInnerWindow_Close(t *testing.T) {
	w := NewInnerWindow("Thing", widget.NewLabel("Content"))

	outer := test.NewWindow(w)
	outer.SetPadded(false)
	outer.Resize(w.MinSize())
	assert.True(t, w.Visible())

	closePos := fyne.NewPos(10, 10)
	test.TapCanvas(outer.Canvas(), closePos)
	assert.False(t, w.Visible())

	w.Show()
	assert.True(t, w.Visible())

	closing := true
	w.CloseIntercept = func() {
		closing = true
	}

	test.TapCanvas(outer.Canvas(), closePos)
	assert.True(t, closing)
	assert.True(t, w.Visible())
}

func TestInnerWindow_MinSize(t *testing.T) {
	w := NewInnerWindow("Thing", widget.NewLabel("Content"))

	btnMin := widget.NewButtonWithIcon("", theme.WindowCloseIcon(), func() {}).MinSize()
	labelMin := widget.NewLabel("Inner").MinSize()

	winMin := w.MinSize()
	assert.Equal(t, btnMin.Height+labelMin.Height+theme.Padding()*2, winMin.Height)
	assert.Greater(t, winMin.Width, btnMin.Width*3+theme.Padding()*3)

	w2 := NewInnerWindow("Much longer title that will truncate", widget.NewLabel("Content"))
	assert.Equal(t, winMin, w2.MinSize())
}
