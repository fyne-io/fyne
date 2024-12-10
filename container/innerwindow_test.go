package container

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
)

func TestInnerWindow_Title(t *testing.T) {
	w := NewInnerWindow("Thing", widget.NewLabel("Content"))
	w.SetTitle("New Title 123")
	assert.Equal(t, "New Title 123", w.Title())
}

func TestInnerWindowIcon_Tap_Left(t *testing.T) {
	w := NewInnerWindow("Thing", widget.NewLabel("Content"))
	w.Icon = theme.GridIcon()

	var testValue bool
	w.OnTappedIcon = func() {
		testValue = true
	}
	w.ButtonAlignment = widget.ButtonAlignLeading

	outer := test.NewTempWindow(t, w)
	outer.SetPadded(false)
	outer.Resize(w.MinSize())
	assert.True(t, w.Visible())

	iconPos := fyne.NewPos(w.Size().Width-10, 10)
	test.TapCanvas(outer.Canvas(), iconPos)
	assert.True(t, testValue)

}

func TestInnerWindowIcon_Tap_Right(t *testing.T) {
	w := NewInnerWindow("Thing", widget.NewLabel("Content"))
	w.Icon = theme.GridIcon()

	var testValue bool
	w.OnTappedIcon = func() {
		testValue = true
	}
	w.ButtonAlignment = widget.ButtonAlignTrailing

	outer := test.NewTempWindow(t, w)
	outer.SetPadded(false)
	outer.Resize(w.MinSize())
	assert.True(t, w.Visible())

	iconPos := fyne.NewPos(10, 10)
	test.TapCanvas(outer.Canvas(), iconPos)
	assert.True(t, testValue)

}

func TestInnerWindow_Close_Left(t *testing.T) {
	w := NewInnerWindow("Thing", widget.NewLabel("Content"))
	w.ButtonAlignment = widget.ButtonAlignLeading
	outer := test.NewTempWindow(t, w)
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

func TestInnerWindow_Close_Right(t *testing.T) {
	w := NewInnerWindow("Thing", widget.NewLabel("Content"))
	w.ButtonAlignment = widget.ButtonAlignTrailing
	outer := test.NewTempWindow(t, w)
	outer.SetPadded(false)
	outer.Resize(w.MinSize())
	assert.True(t, w.Visible())

	closePos := fyne.NewPos(w.Size().Width-10, 10)
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
	assert.Equal(t, btnMin.Height+labelMin.Height+theme.Padding()*4, winMin.Height)
	assert.Greater(t, winMin.Width, btnMin.Width*3+theme.Padding()*5)

	w2 := NewInnerWindow("Much longer title that will truncate", widget.NewLabel("Content"))
	assert.Equal(t, winMin, w2.MinSize())
}

func TestInnerWindow_SetContent(t *testing.T) {
	w := NewInnerWindow("Title", widget.NewLabel("Content"))
	r := cache.Renderer(w).(*innerWindowRenderer)
	title := r.Objects()[4].(*fyne.Container)
	assert.Equal(t, "Content", title.Objects[0].(*widget.Label).Text)

	w.SetContent(widget.NewLabel("Content2"))
	assert.Equal(t, "Content2", title.Objects[0].(*widget.Label).Text)
}

func TestInnerWindow_SetPadded(t *testing.T) {
	w := NewInnerWindow("Title", widget.NewLabel("Content"))
	minPadded := w.MinSize()

	w.SetPadded(false)
	assert.Less(t, w.MinSize().Height, minPadded.Height)

	w.SetPadded(true)
	assert.Equal(t, minPadded, w.MinSize())
}

func TestInnerWindow_SetTitle(t *testing.T) {
	w := NewInnerWindow("Title1", widget.NewLabel("Content"))
	r := cache.Renderer(w).(*innerWindowRenderer)
	title := r.bar.Objects[0].(*draggableLabel)
	assert.Equal(t, "Title1", title.Text)

	w.SetTitle("Title2")
	assert.Equal(t, "Title2", title.Text)
}
