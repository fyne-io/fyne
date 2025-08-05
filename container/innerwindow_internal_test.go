package container

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/widget"
)

func TestInnerWindow_Alignment(t *testing.T) {
	w := NewInnerWindow("Title", widget.NewLabel("Content"))
	w.Resize(fyne.NewSize(150, 100))
	assert.Equal(t, widget.ButtonAlignCenter, w.Alignment)
	assert.NotEqual(t, widget.ButtonAlignCenter, w.buttonPosition())

	buttons := cache.Renderer(w).(*innerWindowRenderer).buttonBox
	w.Alignment = widget.ButtonAlignLeading
	w.Refresh()
	assert.Zero(t, buttons.Position().X)

	w.Alignment = widget.ButtonAlignTrailing
	w.Refresh()
	assert.Greater(t, buttons.Position().X, float32(100))
}

func TestInnerWindow_SetContent(t *testing.T) {
	w := NewInnerWindow("Title", widget.NewLabel("Content"))
	r := cache.Renderer(w).(*innerWindowRenderer)
	title := r.Objects()[4].(*fyne.Container)
	assert.Equal(t, "Content", title.Objects[0].(*widget.Label).Text)

	w.SetContent(widget.NewLabel("Content2"))
	assert.Equal(t, "Content2", title.Objects[0].(*widget.Label).Text)
}

func TestInnerWindow_SetMaximized(t *testing.T) {
	w := NewInnerWindow("Title", widget.NewLabel("Content"))

	icon := cache.Renderer(w).(*innerWindowRenderer).buttons[2]
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
