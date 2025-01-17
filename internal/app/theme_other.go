//go:build !linux && !darwin && !windows && !freebsd && !openbsd && !netbsd && !wasm && !test_web_driver

package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/theme"
)

// DefaultVariant returns the systems default fyne.ThemeVariant.
// Normally, you should not need this. It is extracted out of the root app package to give the
// settings app access to it.
func DefaultVariant() fyne.ThemeVariant {
	return theme.VariantDark
}
