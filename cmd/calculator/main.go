// Package main launches the calculator example directly
package main

import "github.com/fyne-io/fyne/examples/apps"

func main() {
	app := apps.NewApp()

	apps.Calculator(app)
}
