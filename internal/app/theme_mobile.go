//go:build android || ios || mobile

package app

import (
	"fyne.io/fyne/v2"
)

// SystemTheme contains the system’s theme variant.
// It is intended for internal use, only!
var SystemTheme fyne.ThemeVariant

// DefaultVariant returns the systems default fyne.ThemeVariant.
// Normally, you should not need this. It is extracted out of the root app package to give the
// settings app access to it.
func DefaultVariant() fyne.ThemeVariant {
	return SystemTheme
}
