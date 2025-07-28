package container

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func TestInnerWindow_Alignment(t *testing.T) {
	w := NewInnerWindow("Title", widget.NewLabel("Content"))
	w.Resize(fyne.NewSize(150, 100))
	assert.Equal(t, widget.ButtonAlignCenter, w.Alignment)
	assert.NotEqual(t, widget.ButtonAlignCenter, w.buttonPosition())

	buttons := test.WidgetRenderer(w).(*innerWindowRenderer).buttonBox
	w.Alignment = widget.ButtonAlignLeading
	w.Refresh()
	assert.Zero(t, buttons.Position().X)

	w.Alignment = widget.ButtonAlignTrailing
	w.Refresh()
	assert.Greater(t, buttons.Position().X, float32(100))
}

func TestInnerWindow_Close(t *testing.T) {
	w := NewInnerWindow("Thing", widget.NewLabel("Content"))

	outer := test.NewTempWindow(t, w)
	outer.SetPadded(false)
	outer.Resize(w.MinSize())
	assert.True(t, w.Visible())

	closePos := fyne.NewPos(10, 10)
	if w.buttonPosition() == widget.ButtonAlignTrailing {
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
	w := NewInnerWindow("Thing", content)

	btnMin := theme.Size(theme.SizeNameWindowButtonHeight)

	winMin := w.MinSize()
	assert.Equal(t, content.MinSize().Height+theme.Size(theme.SizeNameWindowTitleBarHeight)+theme.Padding()*3, winMin.Height)
	assert.Greater(t, winMin.Width, btnMin*3+theme.Padding()*5)

	w2 := NewInnerWindow("Much longer title that will truncate", widget.NewLabel("Content"))
	assert.Equal(t, winMin, w2.MinSize())
}

func TestInnerWindow_SetContent(t *testing.T) {
	w := NewInnerWindow("Title", widget.NewLabel("Content"))
	r := cache.Renderer(w).(*innerWindowRenderer)
	title := r.Objects()[3].(*fyne.Container)
	assert.Equal(t, "Content", title.Objects[0].(*widget.Label).Text)

	w.SetContent(widget.NewLabel("Content2"))
	assert.Equal(t, "Content2", title.Objects[0].(*widget.Label).Text)
}

func TestInnerWindow_SetMaximized(t *testing.T) {
	w := NewInnerWindow("Title", widget.NewLabel("Content"))

	icon := test.WidgetRenderer(w).(*innerWindowRenderer).buttons[2]
	assert.Equal(t, "foreground_maximize.svg", icon.b.Icon.Name())

	w.SetMaximized(true)
	assert.Equal(t, "foreground_view-zoom-fit.svg", icon.b.Icon.Name())
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
	title := r.bar.Objects[2].(*fyne.Container).Objects[0].(*draggableLabel)
	assert.Equal(t, "Title1", title.Text)

	w.SetTitle("Title2")
	assert.Equal(t, "Title2", title.Text)
}
