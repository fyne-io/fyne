//go:build !hints
// +build !hints

package theme

import (
	"image/color"

	"fyne.io/fyne/v2"
)

var (
	fallbackColor = color.Transparent
	fallbackIcon  = &fyne.StaticResource{}
)
