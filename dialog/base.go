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
	Hide()
	SetDismissText(label string)
	SetOnClosed(closed func())
}

// Declare conformity to Dialog interface
var _ Dialog = (*dialog)(nil)

type dialog struct {
	callback         func(bool)
	originalCallback func(bool) // used by SetOnClosed() if there is already a callback set
	sendResponse     bool
	title            string
	icon             fyne.Resource

	win            *widget.PopUp
	bg             *canvas.Rectangle
	content, label fyne.CanvasObject
	dismiss        *widget.Button
	parent         fyne.Window
}

// SetOnClosed allows to set a callback function that is called when
// the dialog is closed
func (d *dialog) SetOnClosed(closed func()) {
	// if there is already a callback set, remember it and call both
	if d.callback != nil {
		d.originalCallback = d.callback
	}

	d.callback = func(response bool) {
		closed()
		if d.originalCallback != nil {
			d.originalCallback(response)
		}
	}
}

func (d *dialog) setButtons(buttons fyne.CanvasObject) {
	d.bg = canvas.NewRectangle(theme.BackgroundColor())
	d.label = widget.NewLabelWithStyle(d.title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	var content fyne.CanvasObject
	if d.icon == nil {
		content = fyne.NewContainerWithLayout(d,
			&canvas.Image{},
			d.bg,
			d.content,
			buttons,
			d.label,
		)
	} else {
		bgIcon := canvas.NewImageFromResource(d.icon)
		content = fyne.NewContainerWithLayout(d,
			bgIcon,
			d.bg,
			d.content,
			buttons,
			d.label,
		)
	}

	d.win = widget.NewModalPopUp(content, d.parent.Canvas())
	d.applyTheme()
}

func (d *dialog) Layout(obj []fyne.CanvasObject, size fyne.Size) {
	d.bg.Move(fyne.NewPos(-theme.Padding(), -theme.Padding()))
	d.bg.Resize(size.Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)))

	textMin := obj[2].MinSize()
	btnMin := obj[3].MinSize().Union(obj[3].Size())

	// icon
	iconHeight := padHeight*2 + textMin.Height + d.label.MinSize().Height - theme.Padding()
	obj[0].Resize(fyne.NewSize(iconHeight, iconHeight))
	obj[0].Move(fyne.NewPos(size.Width-iconHeight+theme.Padding(), -theme.Padding()))

	// content (text)
	obj[2].Move(fyne.NewPos(size.Width/2-(textMin.Width/2), size.Height-padHeight-btnMin.Height-textMin.Height-theme.Padding()))
	obj[2].Resize(fyne.NewSize(textMin.Width, textMin.Height))

	// buttons
	obj[3].Resize(btnMin)
	obj[3].Move(fyne.NewPos(size.Width/2-(btnMin.Width/2), size.Height-padHeight-btnMin.Height))
}

func (d *dialog) MinSize(obj []fyne.CanvasObject) fyne.Size {
	textMin := obj[2].MinSize()
	btnMin := obj[3].MinSize().Union(obj[3].Size())

	width := fyne.Max(fyne.Max(textMin.Width, btnMin.Width), obj[4].MinSize().Width) + padWidth*2
	height := textMin.Height + btnMin.Height + d.label.MinSize().Height + theme.Padding() + padHeight*2

	return fyne.NewSize(width, height)
}

func (d *dialog) applyTheme() {
	r, g, b, _ := theme.BackgroundColor().RGBA()
	bg := &color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 230}
	d.bg.FillColor = bg
}

func newDialog(title, message string, icon fyne.Resource, callback func(bool), parent fyne.Window) *dialog {
	d := &dialog{content: newLabel(message), title: title, icon: icon, parent: parent}

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
	d.sendResponse = true
	d.win.Show()
}

func (d *dialog) Hide() {
	d.hideWithResponse(false)
}

// SetDismissText allows custom text to be set in the confirmation button
func (d *dialog) SetDismissText(label string) {
	d.dismiss.SetText(label)
	widget.Refresh(d.win)
}

func (d *dialog) hideWithResponse(resp bool) {
	d.win.Hide()
	if d.sendResponse && d.callback != nil {
		d.callback(resp)
	}
	d.sendResponse = false
}

// ShowCustom shows a dialog over the specified application using custom
// content. The button will have the dismiss text set.
// The MinSize() of the CanvasObject passed will be used to set the size of the window.
func ShowCustom(title, dismiss string, content fyne.CanvasObject, parent fyne.Window) {
	d := &dialog{content: content, title: title, icon: nil, parent: parent}

	d.dismiss = &widget.Button{Text: dismiss,
		OnTapped: d.Hide,
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
	d := &dialog{content: content, title: title, icon: nil, parent: parent}
	d.callback = callback

	d.dismiss = &widget.Button{Text: dismiss, Icon: theme.CancelIcon(),
		OnTapped: d.Hide,
	}
	ok := &widget.Button{Text: confirm, Icon: theme.ConfirmIcon(), Style: widget.PrimaryButton,
		OnTapped: func() {
			d.hideWithResponse(true)
		},
	}
	d.setButtons(widget.NewHBox(layout.NewSpacer(), d.dismiss, ok, layout.NewSpacer()))

	d.Show()
}
