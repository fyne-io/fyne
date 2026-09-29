package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type dummyCanvas struct {
	fyne.Canvas
}

func TestOverlayContainer_AccessibilityChildren(t *testing.T) {
	content := canvas.NewRectangle(nil)
	overlay := &OverlayContainer{Content: content, Background: canvas.NewRectangle(nil)}
	assert.Equal(t, []fyne.CanvasObject{content}, overlay.AccessibilityChildren())
	overlay.Content = nil
	assert.Nil(t, overlay.AccessibilityChildren())
}

func TestOverlayContainer_Tapped_Dismiss(t *testing.T) {
	dismissed := false
	o := NewOverlayContainer(nil, &dummyCanvas{}, func() { dismissed = true })

	o.Tapped(&fyne.PointEvent{})
	assert.True(t, dismissed)
}

func TestOverlayContainer_TappedSecondary_Dismiss(t *testing.T) {
	dismissed := false
	o := NewOverlayContainer(nil, &dummyCanvas{}, func() { dismissed = true })

	o.TappedSecondary(&fyne.PointEvent{})
	assert.True(t, dismissed)
}

func TestOverlayContainer_Tapped_NilDismiss(t *testing.T) {
	o := NewOverlayContainer(nil, &dummyCanvas{}, nil)

	assert.NotPanics(t, func() {
		o.Tapped(&fyne.PointEvent{})
	})
}

func TestOverlayContainer_TappedSecondary_NilDismiss(t *testing.T) {
	o := NewOverlayContainer(nil, &dummyCanvas{}, nil)

	assert.NotPanics(t, func() {
		o.TappedSecondary(&fyne.PointEvent{})
	})
}
