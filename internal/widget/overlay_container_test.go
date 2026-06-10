package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
)

type dummyCanvas struct {
	fyne.Canvas
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
