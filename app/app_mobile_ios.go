// +build !ci

// +build darwin,ios

package app

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

func defaultTheme() fyne.Theme {
	// TODO read the iOS setting when 10.13 arrives in 2019
	return theme.LightTheme()
}
