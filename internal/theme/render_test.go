package theme

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"github.com/stretchr/testify/assert"
)

func TestRenderingTheme(t *testing.T) {
	current := &simpleTheme{onlyColor: color.Black}
	assert.Equal(t, current, CurrentlyRenderingWithFallback(current))

	custom := &simpleTheme{onlyColor: color.White}
	PushRenderingTheme(custom)
	assert.Equal(t, custom, CurrentlyRenderingWithFallback(current))

	PopRenderingTheme()
	assert.Equal(t, current, CurrentlyRenderingWithFallback(current))
}

type simpleTheme struct {
	fyne.Theme

	onlyColor color.Color
}

func (s *simpleTheme) Color(fyne.ThemeColorName, fyne.ThemeVariant) color.Color {
	return s.onlyColor
}
