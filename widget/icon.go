package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
)

type iconRenderer struct {
	widget.BaseRenderer
	raster *canvas.Image

	image *Icon
}

func (i *iconRenderer) MinSize() fyne.Size {
	return fyne.NewSquareSize(theme.IconInlineSize())
}

func (i *iconRenderer) Layout(size fyne.Size) {
	if len(i.Objects()) == 0 {
		return
	}

	i.Objects()[0].Resize(size)
}

func (i *iconRenderer) Refresh() {
	if i.image.Resource == i.image.cachedRes {
		return
	}

	i.image.propertyLock.RLock()
	i.raster.Resource = i.image.Resource
	i.image.cachedRes = i.image.Resource

	if i.image.Resource == nil {
		i.raster.Image = nil // reset the internal caching too...
	}
	i.image.propertyLock.RUnlock()

	i.raster.Refresh()
}

// Icon widget is a basic image component that load's its resource to match the theme.
type Icon struct {
	BaseWidget

	Resource  fyne.Resource // The resource for this icon
	cachedRes fyne.Resource
}

// SetResource updates the resource rendered in this icon widget
func (i *Icon) SetResource(res fyne.Resource) {
	i.Resource = res
	i.Refresh()
}

// MinSize returns the size that this widget should not shrink below
func (i *Icon) MinSize() fyne.Size {
	i.ExtendBaseWidget(i)
	return i.BaseWidget.MinSize()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (i *Icon) CreateRenderer() fyne.WidgetRenderer {
	i.ExtendBaseWidget(i)
	i.propertyLock.RLock()
	defer i.propertyLock.RUnlock()

	img := canvas.NewImageFromResource(i.Resource)
	img.FillMode = canvas.ImageFillContain
	r := &iconRenderer{image: i, raster: img}
	r.SetObjects([]fyne.CanvasObject{img})
	i.cachedRes = i.Resource
	return r
}

// NewIcon returns a new icon widget that displays a themed icon resource
func NewIcon(res fyne.Resource) *Icon {
	icon := &Icon{}
	icon.ExtendBaseWidget(icon)
	icon.SetResource(res) // force the image conversion

	return icon
}
