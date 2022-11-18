//go:build !linux
// +build !linux

package desktop

import (
	"fyne.io/fyne/v2"
)

// NewGnomeTheme returns the GNOME theme. If the current OS is not Linux, it returns the default theme.
func NewGnomeTheme() fyne.Theme {
	return DefaultTheme()
}
