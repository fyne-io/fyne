package settings

import (
	"testing"

	"fyne.io/fyne/canvas"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/stretchr/testify/assert"
)

func TestChooseScale(t *testing.T) {
	s := &Settings{}
	s.fyneSettings.Scale = 1.0
	buttons := s.makeScaleButtons()

	test.Tap(buttons[0].(*widget.Button))
	assert.Equal(t, float32(0.5), s.fyneSettings.Scale)
	assert.Equal(t, widget.PrimaryButton, buttons[0].(*widget.Button).Style)
	assert.Equal(t, widget.DefaultButton, buttons[2].(*widget.Button).Style)
}

func TestMakeScaleButtons(t *testing.T) {
	s := &Settings{}
	s.fyneSettings.Scale = 1.0
	buttons := s.makeScaleButtons()

	assert.Equal(t, 5, len(buttons))
	assert.Equal(t, widget.DefaultButton, buttons[0].(*widget.Button).Style)
	assert.Equal(t, widget.PrimaryButton, buttons[2].(*widget.Button).Style)
}

func TestMakeScalePreviews(t *testing.T) {
	s := &Settings{}
	s.fyneSettings.Scale = 1.0
	previews := s.makeScalePreviews(1.0)

	assert.Equal(t, 5, len(previews))
	assert.Equal(t, theme.TextSize(), previews[2].(*canvas.Text).TextSize)

	s.appliedScale(1.5)
	assert.Equal(t, int(float32(theme.TextSize())/1.5), previews[2].(*canvas.Text).TextSize)
}
