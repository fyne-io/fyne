package driver

import "github.com/fyne-io/fyne/ui"

type Driver interface {
	Run()
	Quit()
	CreateWindow(string) ui.Window
}
