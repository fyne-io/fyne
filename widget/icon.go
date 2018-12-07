package widget

import (
	"image/color"

	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/canvas"
	"github.com/fyne-io/fyne/theme"
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

	var raster *canvas.Image
	if i.image.Resource != nil {
		raster = canvas.NewImageFromResource(i.image.Resource)
	}
	i.objects = append(i.objects, raster)
	i.Layout(i.image.CurrentSize())

	canvas.Refresh(i.image)
}

// Icon widget is a basic image component that load's its resource to match the theme.
type Icon struct {
	baseWidget

	Resource fyne.Resource // The resource for this icon
}

// SetResource updates the resource rendered in this icon widget
func (i *Icon) SetResource(res fyne.Resource) {
	i.Resource = res

	i.Renderer().Refresh()
}

func (i *Icon) createRenderer() fyne.WidgetRenderer {
	render := &iconRenderer{image: i}

	render.objects = []fyne.CanvasObject{}

	return render
}

// Renderer is a private method to Fyne which links this widget to it's renderer
func (i *Icon) Renderer() fyne.WidgetRenderer {
	if i.renderer == nil {
		i.renderer = i.createRenderer()
	}

	return i.renderer
}

// NewIcon returns a new icon widget that displays a themed icon resource
func NewIcon(res fyne.Resource) *Icon {
	icon := &Icon{}
	icon.SetResource(res) // force the image conversion

	return icon
}
