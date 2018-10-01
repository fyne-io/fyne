// Package dialog defines standard dialog windows for application GUIs
package dialog

import (
	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/canvas"
	"github.com/fyne-io/fyne/layout"
	"github.com/fyne-io/fyne/theme"
	"github.com/fyne-io/fyne/widget"
)

type dialog struct {
	win      fyne.Window
	callback func(bool)
	message  string
	icon     fyne.Resource

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
	if !d.responded && d.callback != nil {
		d.callback(false)
	}
}

func (d *dialog) setButtons(buttons fyne.CanvasObject) {
	bgIcon := canvas.NewImageFromResource(d.icon)
	bgIcon.Alpha = 0.25
	d.win.SetContent(fyne.NewContainerWithLayout(d,
		newLabel(d.message),
		bgIcon,
		buttons,
	))
}

func (d *dialog) Layout(obj []fyne.CanvasObject, size fyne.Size) {
	middle := size.Height / 2

	// icon
	obj[1].Resize(fyne.NewSize(size.Height*2, size.Height*2))
	obj[1].Move(fyne.NewPos(-size.Height*3/4, -size.Height/2))

	// text
	textMin := obj[0].MinSize()
	obj[0].Move(fyne.NewPos(0, middle-textMin.Height-theme.Padding()/2))
	obj[0].Resize(fyne.NewSize(size.Width, textMin.Height))

	// buttons
	btnMin := obj[2].MinSize()
	obj[2].Resize(btnMin)
	obj[2].Move(fyne.NewPos(size.Width/2-(btnMin.Width/2), middle+theme.Padding()/2))
}

func (d *dialog) MinSize(obj []fyne.CanvasObject) fyne.Size {
	textMin := obj[0].MinSize()
	btnMin := obj[2].MinSize()

	return fyne.NewSize(fyne.Max(textMin.Width, btnMin.Width)+64,
		textMin.Height+btnMin.Height+theme.Padding()+32)
}

func newDialog(title, message string, icon fyne.Resource, callback func(bool), parent fyne.App) *dialog {
	dialog := &dialog{message: message, icon: icon}

	win := parent.NewWindow(title)
	win.SetOnClosed(dialog.closed)
	win.SetFixedSize(true)

	dialog.win = win
	dialog.response = make(chan bool, 1)
	dialog.callback = callback

	return dialog
}

func newLabel(message string) fyne.CanvasObject {
	label := &widget.Label{Text: message, Alignment: fyne.TextAlignCenter}

	return label
}

func newButtonList(buttons ...*widget.Button) fyne.CanvasObject {
	list := fyne.NewContainerWithLayout(layout.NewGridLayout(len(buttons)))

	for _, button := range buttons {
		list.AddObject(button)
	}

	return list
}
