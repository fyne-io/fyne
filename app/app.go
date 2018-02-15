package app

import "github.com/fyne-io/fyne/ui"

type App interface {
	NewWindow(title string) ui.Window
}

