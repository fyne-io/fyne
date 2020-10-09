// Package dialog defines standard dialog windows for application GUIs.
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
	Refresh()
	Resize(size fyne.Size)
}

// Declare conformity to Dialog interface
var _ Dialog = (*dialog)(nil)

type dialog struct {
	callback     func(bool)
	sendResponse bool
	title        string
	icon         fyne.Resource

	win            *widget.PopUp
	bg             *canvas.Rectangle
	content, label fyne.CanvasObject
	dismiss        *widget.Button
	parent         fyne.Window
}

// NewCustom creates and returns a dialog over the specified application using custom
// content. The button will have the dismiss text set.
// The MinSize() of the CanvasObject passed will be used to set the size of the window.
func NewCustom(title, dismiss string, content fyne.CanvasObject, parent fyne.Window) Dialog {
	d := &dialog{content: content, title: title, icon: nil, parent: parent}

	d.dismiss = &widget.Button{Text: dismiss,
		OnTapped: d.Hide,
	}
	d.setButtons(widget.NewHBox(layout.NewSpacer(), d.dismiss, layout.NewSpacer()))

	return d
}

// NewCustomConfirm creates and returns a dialog over the specified application using
// custom content. The cancel button will have the dismiss text set and the "OK" will
// use the confirm text. The response callback is called on user action.
// The MinSize() of the CanvasObject passed will be used to set the size of the window.
func NewCustomConfirm(title, confirm, dismiss string, content fyne.CanvasObject,
	callback func(bool), parent fyne.Window) Dialog {
	d := &dialog{content: content, title: title, icon: nil, parent: parent}
	d.callback = callback
	// TODO: This is required to avoid confusion.
	// Normally this function should only provide the dialog, but currently it is also displayed, which is wrong.
	// For this case the ShowCustomConfirm() method was built.
	d.sendResponse = true

	d.dismiss = &widget.Button{Text: dismiss, Icon: theme.CancelIcon(),
		OnTapped: d.Hide,
	}
	ok := &widget.Button{Text: confirm, Icon: theme.ConfirmIcon(), Style: widget.PrimaryButton,
		OnTapped: func() {
			d.hideWithResponse(true)
		},
	}
	d.setButtons(widget.NewHBox(layout.NewSpacer(), d.dismiss, ok, layout.NewSpacer()))

	return d
}

// ShowCustom shows a dialog over the specified application using custom
// content. The button will have the dismiss text set.
// The MinSize() of the CanvasObject passed will be used to set the size of the window.
func ShowCustom(title, dismiss string, content fyne.CanvasObject, parent fyne.Window) {
	NewCustom(title, dismiss, content, parent).Show()
}

// ShowCustomConfirm shows a dialog over the specified application using custom
// content. The cancel button will have the dismiss text set and the "OK" will use
// the confirm text. The response callback is called on user action.
// The MinSize() of the CanvasObject passed will be used to set the size of the window.
func ShowCustomConfirm(title, confirm, dismiss string, content fyne.CanvasObject,
	callback func(bool), parent fyne.Window) {
	NewCustomConfirm(title, confirm, dismiss, content, callback, parent).Show()
}

func (d *dialog) Hide() {
	d.hideWithResponse(false)
}

func (d *dialog) Show() {
	d.sendResponse = true
	d.win.Show()
}

func (d *dialog) Layout(obj []fyne.CanvasObject, size fyne.Size) {
	d.bg.Move(fyne.NewPos(0, 0))
	d.bg.Resize(size)

	btnMin := obj[3].MinSize().Union(obj[3].Size())

	// icon
	iconHeight := padHeight*2 + d.label.MinSize().Height*2 - theme.Padding()
	obj[0].Resize(fyne.NewSize(iconHeight, iconHeight))
	obj[0].Move(fyne.NewPos(size.Width-iconHeight+theme.Padding(), -theme.Padding()))

	// buttons
	obj[3].Resize(btnMin)
	obj[3].Move(fyne.NewPos(size.Width/2-(btnMin.Width/2), size.Height-padHeight-btnMin.Height))

	// content
	contentStart := d.label.Position().Y + d.label.MinSize().Height + padHeight
	contentEnd := obj[3].Position().Y - theme.Padding()
	obj[2].Move(fyne.NewPos((padWidth / 2), d.label.MinSize().Height+padHeight))
	obj[2].Resize(fyne.NewSize(size.Width-padWidth, contentEnd-contentStart))
}

func (d *dialog) MinSize(obj []fyne.CanvasObject) fyne.Size {
	contentMin := obj[2].MinSize()
	btnMin := obj[3].MinSize().Union(obj[3].Size())

	width := fyne.Max(fyne.Max(contentMin.Width, btnMin.Width), obj[4].MinSize().Width) + padWidth
	height := contentMin.Height + btnMin.Height + d.label.MinSize().Height + theme.Padding() + padHeight*2

	return fyne.NewSize(width, height)
}

func (d *dialog) Refresh() {
	d.applyTheme()
	d.win.Refresh()
}

// Resize dialog, call this function after dialog show
func (d *dialog) Resize(size fyne.Size) {
	maxSize := d.win.Size()
	minSize := d.win.MinSize()
	newWidth := size.Width
	if size.Width > maxSize.Width {
		newWidth = maxSize.Width
	} else if size.Width < minSize.Width {
		newWidth = minSize.Width
	}
	newHeight := size.Height
	if size.Height > maxSize.Height {
		newHeight = maxSize.Height
	} else if size.Height < minSize.Height {
		newHeight = minSize.Height
	}
	d.win.Resize(fyne.NewSize(newWidth, newHeight))
}

// SetDismissText allows custom text to be set in the confirmation button
func (d *dialog) SetDismissText(label string) {
	d.dismiss.SetText(label)
	widget.Refresh(d.win)
}

// SetOnClosed allows to set a callback function that is called when
// the dialog is closed
func (d *dialog) SetOnClosed(closed func()) {
	// if there is already a callback set, remember it and call both
	originalCallback := d.callback

	d.callback = func(response bool) {
		closed()
		if originalCallback != nil {
			originalCallback(response)
		}
	}
}

func (d *dialog) applyTheme() {
	r, g, b, _ := theme.BackgroundColor().RGBA()
	bg := &color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 230}
	d.bg.FillColor = bg
}

func (d *dialog) hideWithResponse(resp bool) {
	d.win.Hide()
	if d.sendResponse && d.callback != nil {
		d.callback(resp)
	}
	d.sendResponse = false
}

func (d *dialog) setButtons(buttons fyne.CanvasObject) {
	d.bg = canvas.NewRectangle(theme.BackgroundColor())
	d.label = newDialogTitle(d.title, d)

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
	d.Refresh()
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

// dialogTitle is really just a normal title but we use the Refresh() hook to update the background rectangle.
type dialogTitle struct {
	widget.Label

	d *dialog
}

// Refresh applies the current theme to the whole dialog before refreshing the underlying label.
func (t *dialogTitle) Refresh() {
	t.d.Refresh()

	t.BaseWidget.Refresh()
}

func newDialogTitle(title string, d *dialog) *dialogTitle {
	l := &dialogTitle{}
	l.Text = title
	l.Alignment = fyne.TextAlignLeading
	l.TextStyle.Bold = true

	l.d = d
	l.ExtendBaseWidget(l)
	return l
}
