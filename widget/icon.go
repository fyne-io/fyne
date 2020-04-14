package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/binding"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

type iconRenderer struct {
	baseRenderer
	image *Icon
}

func (i *iconRenderer) MinSize() fyne.Size {
	size := theme.IconInlineSize()
	return fyne.NewSize(size, size)
}

func (i *iconRenderer) Layout(size fyne.Size) {
	if len(i.Objects()) == 0 {
		return
	}

	i.Objects()[0].Resize(size)
}

func (i *iconRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (i *iconRenderer) Refresh() {
	i.setObjects(nil)

	if i.image.Resource != nil {
		raster := canvas.NewImageFromResource(i.image.Resource)
		raster.FillMode = canvas.ImageFillContain

		i.setObjects([]fyne.CanvasObject{raster})
	}
	i.Layout(i.image.Size())

	canvas.Refresh(i.image.super())
}

// Icon widget is a basic image component that load's its resource to match the theme.
type Icon struct {
	BaseWidget

	Resource fyne.Resource // The resource for this icon

	resourceBind   binding.Resource
	resourceNotify binding.Notifiable
}

// SetResource updates the resource rendered in this icon widget
func (i *Icon) SetResource(res fyne.Resource) {
	if i.Resource != res {
		i.Resource = res
		i.Refresh()
	}
}

// MinSize returns the size that this widget should not shrink below
func (i *Icon) MinSize() fyne.Size {
	i.ExtendBaseWidget(i)
	return i.BaseWidget.MinSize()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (i *Icon) CreateRenderer() fyne.WidgetRenderer {
	i.ExtendBaseWidget(i)
	return &iconRenderer{image: i}
}

// BindResource binds the Icon's Resource to the given data binding.
// Returns the Icon for chaining.
func (i *Icon) BindResource(data binding.Resource) *Icon {
	i.resourceBind = data
	i.resourceNotify = data.AddResourceListener(i.SetResource)
	return i
}

// UnbindResource unbinds the Icon's Resource from the data binding (if any).
// Returns the Icon for chaining.
func (i *Icon) UnbindResource() *Icon {
	if i.resourceBind != nil {
		i.resourceBind.DeleteListener(i.resourceNotify)
	}
	i.resourceBind = nil
	i.resourceNotify = nil
	return i
}

// NewIcon returns a new icon widget that displays a themed icon resource
func NewIcon(res fyne.Resource) *Icon {
	icon := &Icon{}
	icon.ExtendBaseWidget(icon)
	icon.SetResource(res) // force the image conversion

	return icon
}
