package widget

import (
	"testing"

	"fyne.io/fyne/binding"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"github.com/stretchr/testify/assert"
)

func TestNewIcon(t *testing.T) {
	icon := NewIcon(theme.ConfirmIcon())
	render := test.WidgetRenderer(icon)

	assert.Equal(t, 1, len(render.Objects()))
	obj := render.Objects()[0]
	img, ok := obj.(*canvas.Image)
	if !ok {
		t.Fail()
	}
	assert.Equal(t, theme.ConfirmIcon(), img.Resource)
}

func TestIcon_Nil(t *testing.T) {
	icon := NewIcon(nil)
	render := test.WidgetRenderer(icon)

	assert.Equal(t, 0, len(render.Objects()))
}

func TestIcon_MinSize(t *testing.T) {
	icon := NewIcon(theme.CancelIcon())
	min := icon.MinSize()

	assert.Equal(t, theme.IconInlineSize(), min.Width)
	assert.Equal(t, theme.IconInlineSize(), min.Height)
}

func TestIcon_BindResource(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	done := make(chan bool)
	icon := NewIcon(theme.WarningIcon())

	resource := theme.QuestionIcon()
	data := binding.NewResourceRef(&resource)
	icon.BindResource(data)
	data.AddListener(binding.NewNotifyFunction(func(binding.Binding) {
		done <- true
	}))
	timedWait(t, done)
	assert.Equal(t, theme.QuestionIcon(), icon.Resource)

	// Set directly
	resource = theme.SearchIcon()
	data.Update()
	timedWait(t, done)
	assert.Equal(t, theme.SearchIcon(), icon.Resource)

	// Set by binding
	data.Set(theme.InfoIcon())
	timedWait(t, done)
	assert.Equal(t, theme.InfoIcon(), icon.Resource)
}

func TestIconRenderer_ApplyTheme(t *testing.T) {
	icon := NewIcon(theme.CancelIcon())
	render := test.WidgetRenderer(icon).(*iconRenderer)
	visible := render.objects[0].Visible()

	render.Refresh()
	assert.Equal(t, visible, render.objects[0].Visible())
}
