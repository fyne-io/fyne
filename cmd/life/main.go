// Package main launches the game of life example directly
package main

import "github.com/fyne-io/fyne/examples/apps"

func main() {
	app := apps.NewApp()

	apps.Life(app)
}
