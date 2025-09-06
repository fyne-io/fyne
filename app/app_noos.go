//go:build tamago || noos || tinygo

package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/embedded"
	intNoos "fyne.io/fyne/v2/internal/driver/embedded"
	"fyne.io/fyne/v2/theme"
)

func NewWithID(id string) fyne.App {
	return newAppWithDriver(nil, nil, id)
}

// SetDriverDetails provides the required information to our app to run without a standard
// driver. This is useful for embedded devices like GOOS=tamago, tinygo or noos.
//
// Since: 2.7
func SetDriverDetails(a fyne.App, d embedded.Driver) {
	a.(*fyneApp).Settings().SetTheme(theme.DefaultTheme())
	a.(*fyneApp).driver = intNoos.NewNoOSDriver(d.Render, d.Run, d.Queue(), d.ScreenSize)
}
