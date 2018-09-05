package dialog

import (
	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/layout"
	"github.com/fyne-io/fyne/widget"
)

type dialog struct {
	win      fyne.Window
	callback func(bool)

	response  chan bool
	responded bool
}

func (d *dialog) wait() {
	select {
	case response := <-d.response:
		d.responded = true
		d.win.Close()
		if d.callback != nil {
			d.callback(response)
		}
	}
}

func (d *dialog) closed() {
	if !d.responded {
		d.callback(false)
	}
}

func newDialog(title string, callback func(bool), parent fyne.App) *dialog {
	dialog := &dialog{}

	win := parent.NewWindow(title) // TODO don't depend on the main loop calling this!
	win.SetOnClosed(dialog.closed)

	dialog.win = win
	dialog.response = make(chan bool, 1)
	dialog.callback = callback

	return dialog
}

func newLabel(message string) fyne.CanvasObject {
	label := widget.NewLabel(message)
	label.Alignment = fyne.TextAlignCenter
	label.SetText(message) // TODO fix issue where align is not respected

	return label
}

func newButtonList(buttons ...*widget.Button) fyne.CanvasObject {
	list := fyne.NewContainerWithLayout(layout.NewGridLayout(len(buttons)+2), layout.NewSpacer())

	for _, button := range buttons {
		list.AddObject(button)
	}

	list.AddObject(layout.NewSpacer())
	return list
}
