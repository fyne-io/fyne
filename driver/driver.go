package driver

import "github.com/fyne-io/fyne/ui"

type Driver interface {
	CreateWindow(string) ui.Window
	Run()
}
