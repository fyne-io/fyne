package main

import "github.com/fyne-io/fyne/app"

func main() {
	w := app.NewWindow("Hello")
	w.Canvas().NewText("Hello Fyne!")

	app.Run()
}
