package container_test

import (
	"runtime"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func TestInnerWindow_Close(t *testing.T) {
	w := container.NewInnerWindow("Thing", widget.NewLabel("Content"))

	outer := test.NewTempWindow(t, w)
	outer.SetPadded(false)
	outer.Resize(w.MinSize())
	assert.True(t, w.Visible())

	closePos := fyne.NewPos(10, 10)
	if runtime.GOOS != "darwin" {
		closePos = fyne.NewPos(w.Size().Width-10, 10)
	}
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
	content := widget.NewLabel("Content")
	w := container.NewInnerWindow("Thing", content)

	btnMin := theme.Size(theme.SizeNameWindowButtonHeight)

	winMin := w.MinSize()
	assert.Equal(t, content.MinSize().Height+theme.Size(theme.SizeNameWindowTitleBarHeight)+theme.Padding()*3, winMin.Height)
	assert.Greater(t, winMin.Width, btnMin*3+theme.Padding()*5)

	w2 := container.NewInnerWindow("Much longer title that will truncate", widget.NewLabel("Content"))
	assert.Equal(t, winMin, w2.MinSize())
}
