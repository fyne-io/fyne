//go:build tamago || noos || tinygo

package app

import (
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/embedded"
	intNoos "fyne.io/fyne/v2/internal/driver/embedded"
	"fyne.io/fyne/v2/theme"
)

func NewWithID(id string) fyne.App {
	return newAppWithDriver(nil, nil, id)
}

// SetNoOSDriver provides the required information to our app to run without a standard
// driver. This is useful for embedded devices like GOOS=tamago or GOOS=noos.
//
// Since: 2.7
func SetDriverDetails(render func(img image.Image), events chan embedded.Event) {
	a := fyne.CurrentApp().(*fyneApp)

	a.Settings().SetTheme(theme.DefaultTheme())
	a.driver = intNoos.NewNoOSDriver(render, events)
}
