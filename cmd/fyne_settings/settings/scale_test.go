package settings

import (
	"testing"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
)

func TestChooseScale(t *testing.T) {
	s := &Settings{}
	s.fyneSettings.Scale = 1.0
	buttons := s.makeScaleButtons()

	test.Tap(buttons[0].(*widget.Button))
	assert.Equal(t, float32(0.5), s.fyneSettings.Scale)
	assert.Equal(t, widget.HighImportance, buttons[0].(*widget.Button).Importance)
	assert.Equal(t, widget.MediumImportance, buttons[2].(*widget.Button).Importance)
}

func TestMakeScaleButtons(t *testing.T) {
	s := &Settings{}
	s.fyneSettings.Scale = 1.0
	buttons := s.makeScaleButtons()

	assert.Equal(t, 5, len(buttons))
	assert.Equal(t, widget.MediumImportance, buttons[0].(*widget.Button).Importance)
	assert.Equal(t, widget.HighImportance, buttons[2].(*widget.Button).Importance)
}

func TestMakeScalePreviews(t *testing.T) {
	s := &Settings{}
	s.fyneSettings.Scale = 1.0
	previews := s.makeScalePreviews(1.0)

	assert.Equal(t, 5, len(previews))
	assert.Equal(t, theme.TextSize(), previews[2].(*canvas.Text).TextSize)

	s.appliedScale(1.5)
	assert.Equal(t, theme.TextSize()/1.5, previews[2].(*canvas.Text).TextSize)
}
