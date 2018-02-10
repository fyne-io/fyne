package app

import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/drivers/efl"

var driver efl.EFLDriver

func NewWindow(title string) ui.Window {
	return driver.CreateWindow(title)
}

func Run() {
	driver.Run()
}
