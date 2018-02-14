package app

import "github.com/fyne-io/fyne/ui"

var uiDriver ui.Driver

func SetUIDriver(d ui.Driver) {
	uiDriver = d
}

func NewWindow(title string) ui.Window {
	return uiDriver.CreateWindow(title)
}

func Run() {
	uiDriver.Run()
}
