// Package main loads a very basic Hello World graphical application.
package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello")

	hello := widget.NewLabel("Hello Fyne!")
	w.SetContent(container.NewVBox(
		hello,
		widget.NewButton("Hi!", func() {
			hello.SetText("Welcome :)")
		}),
	))

	if desk, ok := a.(desktop.App); ok {
		menu := fyne.NewMenu("Hello World",
			fyne.NewMenuItem("Hello", func() {
				log.Println("System tray menu tapped")
			}))
		desk.SetSystemTrayMenu(menu)
	}

	w.ShowAndRun()
}
