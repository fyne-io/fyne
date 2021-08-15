// Package dialog defines standard dialog windows for application GUIs.
package dialog // import "fyne.io/fyne/v2/dialog"

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
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
	callback    func(bool)
	title       string
	icon        fyne.Resource
	desiredSize fyne.Size

	win            *widget.PopUp
	bg             *themedBackground
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
	d.setButtons(container.NewHBox(layout.NewSpacer(), d.dismiss, layout.NewSpacer()))

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

	d.dismiss = &widget.Button{Text: dismiss, Icon: theme.CancelIcon(),
		OnTapped: d.Hide,
	}
	ok := &widget.Button{Text: confirm, Icon: theme.ConfirmIcon(), Importance: widget.HighImportance,
		OnTapped: func() {
			d.hideWithResponse(true)
		},
	}
	d.setButtons(container.NewHBox(layout.NewSpacer(), d.dismiss, ok, layout.NewSpacer()))

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
	if !d.desiredSize.IsZero() {
		d.win.Resize(d.desiredSize)
	}
	d.win.Show()
}

func (d *dialog) Layout(obj []fyne.CanvasObject, size fyne.Size) {
	d.bg.Move(fyne.NewPos(0, 0))
	d.bg.Resize(size)

	btnMin := obj[3].MinSize()

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
	obj[2].Move(fyne.NewPos(padWidth/2, d.label.MinSize().Height+padHeight))
	obj[2].Resize(fyne.NewSize(size.Width-padWidth, contentEnd-contentStart))
}

func (d *dialog) MinSize(obj []fyne.CanvasObject) fyne.Size {
	contentMin := obj[2].MinSize()
	btnMin := obj[3].MinSize()

	width := fyne.Max(fyne.Max(contentMin.Width, btnMin.Width), obj[4].MinSize().Width) + padWidth
	height := contentMin.Height + btnMin.Height + d.label.MinSize().Height + theme.Padding() + padHeight*2

	return fyne.NewSize(width, height)
}

func (d *dialog) Refresh() {
	d.win.Refresh()
}

// Resize dialog, call this function after dialog show
func (d *dialog) Resize(size fyne.Size) {
	d.desiredSize = size
	d.win.Resize(size)
}

// SetDismissText allows custom text to be set in the confirmation button
func (d *dialog) SetDismissText(label string) {
	d.dismiss.SetText(label)
	d.win.Refresh()
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

func (d *dialog) hideWithResponse(resp bool) {
	d.win.Hide()
	if d.callback != nil {
		d.callback(resp)
	}
}

func (d *dialog) setButtons(buttons fyne.CanvasObject) {
	d.bg = newThemedBackground()
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
		list.Add(button)
	}

	return list
}

// ===============================================================
// ThemedBackground
// ===============================================================

type themedBackground struct {
	widget.BaseWidget
}

func newThemedBackground() *themedBackground {
	t := &themedBackground{}
	t.ExtendBaseWidget(t)
	return t
}

func (t *themedBackground) CreateRenderer() fyne.WidgetRenderer {
	t.ExtendBaseWidget(t)
	rect := canvas.NewRectangle(theme.BackgroundColor())
	return &themedBackgroundRenderer{rect, []fyne.CanvasObject{rect}}
}

type themedBackgroundRenderer struct {
	rect    *canvas.Rectangle
	objects []fyne.CanvasObject
}

func (renderer *themedBackgroundRenderer) Destroy() {
}

func (renderer *themedBackgroundRenderer) Layout(size fyne.Size) {
	renderer.rect.Resize(size)
}

func (renderer *themedBackgroundRenderer) MinSize() fyne.Size {
	return renderer.rect.MinSize()
}

func (renderer *themedBackgroundRenderer) Objects() []fyne.CanvasObject {
	return renderer.objects
}

func (renderer *themedBackgroundRenderer) Refresh() {
	r, g, b, _ := theme.BackgroundColor().RGBA()
	bg := &color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 230}
	renderer.rect.FillColor = bg
}

// ===============================================================
// DialogLayout
// ===============================================================
