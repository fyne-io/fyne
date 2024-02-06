// Package dialog defines standard dialog windows for application GUIs.
package dialog // import "fyne.io/fyne/v2/dialog"

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	col "fyne.io/fyne/v2/internal/color"
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

	// Since: 2.1
	MinSize() fyne.Size
}

// Declare conformity to Dialog interface
var _ Dialog = (*dialog)(nil)

type dialog struct {
	callback    func(bool)
	title       string
	icon        fyne.Resource
	desiredSize fyne.Size

	win     *widget.PopUp
	content fyne.CanvasObject
	dismiss *widget.Button
	parent  fyne.Window
}

func (d *dialog) Hide() {
	d.hideWithResponse(false)
}

// MinSize returns the size that this dialog should not shrink below
//
// Since: 2.1
func (d *dialog) MinSize() fyne.Size {
	return d.win.MinSize()
}

func (d *dialog) Show() {
	if !d.desiredSize.IsZero() {
		d.win.Resize(d.desiredSize)
	}
	d.win.Show()
}

func (d *dialog) Refresh() {
	d.win.Refresh()
}

// Resize dialog, call this function after dialog show
func (d *dialog) Resize(size fyne.Size) {
	d.desiredSize = size
	d.win.Resize(size)
}

// SetDismissText allows custom text to be set in the dismiss button
// This is a no-op for dialogs without dismiss buttons.
func (d *dialog) SetDismissText(label string) {
	if d.dismiss == nil {
		return
	}

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

func (d *dialog) create(buttons fyne.CanvasObject) {
	label := widget.NewLabelWithStyle(d.title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	var image fyne.CanvasObject
	if d.icon != nil {
		image = &canvas.Image{Resource: d.icon}
	} else {
		image = &layout.Spacer{}
	}

	content := container.New(&dialogLayout{d: d},
		image,
		newThemedBackground(),
		d.content,
		buttons,
		label,
	)

	d.win = widget.NewModalPopUp(content, d.parent.Canvas())
}

func (d *dialog) setButtons(buttons fyne.CanvasObject) {
	d.win.Content.(*fyne.Container).Objects[3] = buttons
	d.win.Refresh()
}

// The method .create() needs to be called before the dialog cna be shown.
func newDialog(title, message string, icon fyne.Resource, callback func(bool), parent fyne.Window) *dialog {
	d := &dialog{content: newCenterLabel(message), title: title, icon: icon, parent: parent}
	d.callback = callback

	return d
}

func newCenterLabel(message string) fyne.CanvasObject {
	return &widget.Label{Text: message, Alignment: fyne.TextAlignCenter}
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
	rect := canvas.NewRectangle(theme.OverlayBackgroundColor())
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
	r, g, b, _ := col.ToNRGBA(theme.OverlayBackgroundColor())
	bg := &color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 230}
	renderer.rect.FillColor = bg
}

// ===============================================================
// DialogLayout
// ===============================================================

type dialogLayout struct {
	d *dialog
}

func (l *dialogLayout) Layout(obj []fyne.CanvasObject, size fyne.Size) {
	btnMin := obj[3].MinSize()
	labelMin := obj[4].MinSize()

	// icon
	iconHeight := padHeight*2 + labelMin.Height*2 - theme.Padding()
	obj[0].Resize(fyne.NewSize(iconHeight, iconHeight))
	obj[0].Move(fyne.NewPos(size.Width-iconHeight+theme.Padding(), -theme.Padding()))

	// background
	obj[1].Move(fyne.NewPos(0, 0))
	obj[1].Resize(size)

	// content
	contentStart := obj[4].Position().Y + labelMin.Height + padHeight
	contentEnd := obj[3].Position().Y - theme.Padding()
	obj[2].Move(fyne.NewPos(padWidth/2, labelMin.Height+padHeight))
	obj[2].Resize(fyne.NewSize(size.Width-padWidth, contentEnd-contentStart))

	// buttons
	obj[3].Resize(btnMin)
	obj[3].Move(fyne.NewPos(size.Width/2-(btnMin.Width/2), size.Height-padHeight-btnMin.Height))
}

func (l *dialogLayout) MinSize(obj []fyne.CanvasObject) fyne.Size {
	contentMin := obj[2].MinSize()
	btnMin := obj[3].MinSize()
	labelMin := obj[4].MinSize()

	width := fyne.Max(fyne.Max(contentMin.Width, btnMin.Width), labelMin.Width) + padWidth
	height := contentMin.Height + btnMin.Height + labelMin.Height + theme.Padding() + padHeight*2

	return fyne.NewSize(width, height)
}
