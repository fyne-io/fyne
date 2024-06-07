package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
)

func Test_NewTheme_checkForUniqueColors(t *testing.T) {
	th := NewTheme().(*configurableTheme)
	seen := map[string]fyne.ThemeColorName{}
	for cn, c := range th.colors {
		r, g, b, a := c.RGBA()
		key := fmt.Sprintf("%d %d %d %d", r, g, b, a)
		assert.True(t, seen[key] == "", "color value %#v for color %s already used for color %s in theme %s", c, cn, seen[key], th.name)
		seen[key] = cn
	}
}

func Test_Theme_checkForUniqueColors(t *testing.T) {
	th := Theme().(*configurableTheme)
	seen := map[string]fyne.ThemeColorName{}
	for cn, c := range th.colors {
		r, g, b, a := c.RGBA()
		key := fmt.Sprintf("%d %d %d %d", r, g, b, a)
		assert.True(t, seen[key] == "", "color value %#v for color %s already used for color %s in theme %s", c, cn, seen[key], th.name)
		seen[key] = cn
	}
}
