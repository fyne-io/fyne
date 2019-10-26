// +build android ios mobile

package app

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/theme"
)

func TestDefaultTheme(t *testing.T) {
	assert.Equal(t, theme.LightTheme(), defaultTheme())
}
