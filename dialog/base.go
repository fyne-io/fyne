// Package dialog defines standard dialog windows for application GUIs
package dialog // import "fyne.io/fyne/dialog"

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// Dialog is the common API for any dialog window with a single dismiss button
type Dialog interface {
	Show()
	SetDismissText(label string)
}

type dialog struct {
	win      fyne.Window
	callback func(bool)
	content  fyne.CanvasObject
	icon     fyne.Resource

	dismiss *widget.Button

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
	if d.icon == nil {
		d.win.SetContent(fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
			d.content,
			buttons,
		))
	} else {
		bgIcon := canvas.NewImageFromResource(d.icon)
		bgIcon.Translucency = 0.75
		d.win.SetContent(fyne.NewContainerWithLayout(d,
			d.content,
			bgIcon,
			buttons,
		))
	}
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
	btnMin := obj[2].MinSize().Union(obj[2].Size())
	obj[2].Resize(btnMin)
	obj[2].Move(fyne.NewPos(size.Width/2-(btnMin.Width/2), middle+theme.Padding()/2))
}

func (d *dialog) MinSize(obj []fyne.CanvasObject) fyne.Size {
	textMin := obj[0].MinSize()
	btnMin := obj[2].MinSize().Union(obj[2].Size())

	return fyne.NewSize(fyne.Max(textMin.Width, btnMin.Width)+64,
		textMin.Height+btnMin.Height+theme.Padding()+32)
}

func newDialogWin(title string, _ fyne.Window) fyne.Window {
	win := fyne.CurrentApp().Driver().CreateWindow(title)
	win.SetFixedSize(true)
	//	win.SetParent(parent)

	return win
}

func newDialog(title, message string, icon fyne.Resource, callback func(bool), parent fyne.Window) *dialog {
	d := &dialog{content: newLabel(message), icon: icon}

	win := newDialogWin(title, parent)
	win.SetOnClosed(d.closed)

	d.win = win
	d.response = make(chan bool, 1)
	d.callback = callback

	return d
}

func newLabel(message string) fyne.CanvasObject {
	return widget.NewLabelWithStyle(message, fyne.TextAlignCenter, fyne.TextStyle{})
}

func newButtonList(buttons ...*widget.Button) fyne.CanvasObject {
	list := fyne.NewContainerWithLayout(layout.NewGridLayout(len(buttons)))

	for _, button := range buttons {
		list.AddObject(button)
	}

	return list
}

func (d *dialog) Show() {
	go d.wait()
	d.win.Show()
}

// SetDismissText allows custom text to be set in the confirmation button
func (d *dialog) SetDismissText(label string) {
	d.dismiss.SetText(label)
}

// ShowCustom shows a dialog over the specified application using custom
// content. The button will have th dismiss text set.
// The MinSize() of the CanvasObject passed will be used to set the size of the window.
func ShowCustom(title, dismiss string, content fyne.CanvasObject, parent fyne.Window) {
	d := &dialog{content: content, icon: nil}

	win := newDialogWin(title, parent)
	win.SetOnClosed(d.closed)
	d.win = win
	d.response = make(chan bool, 1)

	d.dismiss = &widget.Button{Text: dismiss,
		OnTapped: func() {
			d.response <- false
		},
	}
	d.setButtons(widget.NewHBox(layout.NewSpacer(), d.dismiss, layout.NewSpacer()))

	d.Show()
}
