package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

type iconRenderer struct {
	objects []fyne.CanvasObject

	image *Icon
}

func (i *iconRenderer) MinSize() fyne.Size {
	size := theme.IconInlineSize()
	return fyne.NewSize(size, size)
}

func (i *iconRenderer) Layout(size fyne.Size) {
	if len(i.objects) == 0 {
		return
	}

	i.objects[0].Resize(size)
}

func (i *iconRenderer) Objects() []fyne.CanvasObject {
	return i.objects
}

func (i *iconRenderer) ApplyTheme() {
	i.Refresh()
}

func (i *iconRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (i *iconRenderer) Refresh() {
	i.objects = nil

	if i.image.Resource != nil {
		raster := canvas.NewImageFromResource(i.image.Resource)
		raster.FillMode = canvas.ImageFillContain

		i.objects = append(i.objects, raster)
	}
	i.Layout(i.image.Size())

	canvas.Refresh(i.image)
}

func (i *iconRenderer) Destroy() {
}

// Icon widget is a basic image component that load's its resource to match the theme.
type Icon struct {
	baseWidget

	Resource fyne.Resource // The resource for this icon
}

// Resize sets a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (i *Icon) Resize(size fyne.Size) {
	i.resize(size, i)
}

// Move the widget to a new position, relative to it's parent.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (i *Icon) Move(pos fyne.Position) {
	i.move(pos, i)
}

// MinSize returns the smallest size this widget can shrink to
func (i *Icon) MinSize() fyne.Size {
	return i.minSize(i)
}

// Show this widget, if it was previously hidden
func (i *Icon) Show() {
	i.show(i)
}

// Hide this widget, if it was previously visible
func (i *Icon) Hide() {
	i.hide(i)
}

// SetResource updates the resource rendered in this icon widget
func (i *Icon) SetResource(res fyne.Resource) {
	i.Resource = res

	Refresh(i)
}

// CreateRenderer is a private method to Fyne which links this widget to it's renderer
func (i *Icon) CreateRenderer() fyne.WidgetRenderer {
	render := &iconRenderer{image: i}

	render.objects = []fyne.CanvasObject{}

	return render
}

// NewIcon returns a new icon widget that displays a themed icon resource
func NewIcon(res fyne.Resource) *Icon {
	icon := &Icon{}
	icon.SetResource(res) // force the image conversion

	return icon
}
