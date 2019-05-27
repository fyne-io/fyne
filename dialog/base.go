// Package dialog defines standard dialog windows for application GUIs
package dialog // import "fyne.io/fyne/dialog"

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

const (
	padWidth  = 32
	padHeight = 16
)

// Dialog is the common API for any dialog window with a single dismiss button
type Dialog interface {
	Show()
	SetDismissText(label string)
}

// Declare confirmity to Dialog interface
var _ Dialog = (*dialog)(nil)

type dialog struct {
	win      fyne.Window
	callback func(bool)
	bg       *canvas.Rectangle
	content  fyne.CanvasObject
	icon     fyne.Resource

	dismiss *widget.Button

	response  chan bool
	responded bool
	parent    fyne.Window
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

	if d.parent != nil {
		d.parent.RequestFocus()
	}
}

func (d *dialog) setButtons(buttons fyne.CanvasObject) {
	d.bg = canvas.NewRectangle(theme.BackgroundColor())

	if d.icon == nil {
		d.win.SetContent(fyne.NewContainerWithLayout(d,
			&canvas.Image{},
			d.bg,
			d.content,
			buttons,
		))
	} else {
		bgIcon := canvas.NewImageFromResource(d.icon)
		d.win.SetContent(fyne.NewContainerWithLayout(d,
			bgIcon,
			d.bg,
			d.content,
			buttons,
		))
	}
}

func (d *dialog) Layout(obj []fyne.CanvasObject, size fyne.Size) {
	d.ApplyTheme() // we are not really a widget so simulate the applyTheme call
	d.bg.Move(fyne.NewPos(-theme.Padding(), -theme.Padding()))
	d.bg.Resize(size.Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)))
	// icon
	obj[0].Resize(fyne.NewSize(size.Height*2, size.Height*2))
	obj[0].Move(fyne.NewPos(-size.Height*3/4, -size.Height/2))

	// content (text)
	textMin := obj[2].MinSize()
	obj[2].Move(fyne.NewPos(size.Width/2-(textMin.Width/2), padHeight))
	obj[2].Resize(fyne.NewSize(textMin.Width, textMin.Height))

	// buttons
	btnMin := obj[3].MinSize().Union(obj[3].Size())
	obj[3].Resize(btnMin)
	obj[3].Move(fyne.NewPos(size.Width/2-(btnMin.Width/2), size.Height-padHeight-btnMin.Height))
}

func (d *dialog) MinSize(obj []fyne.CanvasObject) fyne.Size {
	textMin := obj[2].MinSize()
	btnMin := obj[3].MinSize().Union(obj[3].Size())

	return fyne.NewSize(fyne.Max(textMin.Width, btnMin.Width)+padWidth*2,
		textMin.Height+btnMin.Height+theme.Padding()+padHeight*2)
}

func (d *dialog) ApplyTheme() {
	r, g, b, _ := theme.BackgroundColor().RGBA()
	bg := &color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 230}
	d.bg.FillColor = bg
}

func newDialogWin(title string, _ fyne.Window) fyne.Window {
	win := fyne.CurrentApp().Driver().CreateWindow(title)
	win.SetFixedSize(true)
	//	win.SetParent(parent)

	return win
}

func newDialog(title, message string, icon fyne.Resource, callback func(bool), parent fyne.Window) *dialog {
	d := &dialog{content: newLabel(message), icon: icon, parent: parent}

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
	d.Layout(d.win.Content().(*fyne.Container).Objects, d.win.Content().MinSize())
}

// ShowCustom shows a dialog over the specified application using custom
// content. The button will have the dismiss text set.
// The MinSize() of the CanvasObject passed will be used to set the size of the window.
func ShowCustom(title, dismiss string, content fyne.CanvasObject, parent fyne.Window) {
	d := &dialog{content: content, icon: nil, parent: parent}

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

// ShowCustomConfirm shows a dialog over the specified application using custom
// content. The cancel button will have the dismiss text set and the "OK" will use
// the confirm text. The response callback is called on user action.
// The MinSize() of the CanvasObject passed will be used to set the size of the window.
func ShowCustomConfirm(title, confirm, dismiss string, content fyne.CanvasObject,
	callback func(bool), parent fyne.Window) {
	d := &dialog{content: content, icon: nil, parent: parent}

	win := newDialogWin(title, parent)
	win.SetOnClosed(d.closed)
	d.win = win
	d.response = make(chan bool, 1)
	d.callback = callback

	d.dismiss = &widget.Button{Text: dismiss, Icon: theme.CancelIcon(),
		OnTapped: func() {
			d.response <- false
		},
	}
	ok := &widget.Button{Text: confirm, Icon: theme.ConfirmIcon(), Style: widget.PrimaryButton,
		OnTapped: func() {
			d.response <- true
		},
	}
	d.setButtons(widget.NewHBox(layout.NewSpacer(), d.dismiss, ok, layout.NewSpacer()))

	d.Show()
}
