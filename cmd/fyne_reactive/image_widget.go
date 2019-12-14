package main

import (
	"image/color"

	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"

	"fyne.io/fyne"
	"fyne.io/fyne/dataapi"
	"fyne.io/fyne/widget"
	_ "image/jpeg"
	_ "image/png"
)

// ImageWidget is a custom widget demo
type ImageWidget struct {
	widget.BaseWidget
	widget.DataListener
	urlStr  string
	imgRes  fyne.Resource
	LoadErr error
}

// NewImageWidget returns a new ImageWidget
func NewImageWidget(urlStr string) *ImageWidget {
	img, err := fyne.LoadResourceFromURLString(urlStr)
	if err != nil {
		println("Failed to load", urlStr, err.Error())
	}
	return &ImageWidget{
		BaseWidget:   widget.BaseWidget{},
		DataListener: widget.DataListener{},
		urlStr:       urlStr,
		imgRes:       img,
		LoadErr:      err,
	}
}

// Bind will link the dataitem to the image
func (w *ImageWidget) Bind(data dataapi.DataItem) *ImageWidget {
	w.DataListener.Bind(data, w)
	return w
}

// SetFromData updates the hyperlink from the bound data
func (w *ImageWidget) SetFromData(data dataapi.DataItem) {
	w.urlStr = data.String()
	w.imgRes, w.LoadErr = fyne.LoadResourceFromURLString(w.urlStr)
	w.Refresh()
}

// CreateRenderer creates and returns the renderer for this widget type
func (w *ImageWidget) CreateRenderer() fyne.WidgetRenderer {
	w.ExtendBaseWidget(w)
	img := canvas.NewImageFromResource(w.imgRes)
	return &imageWidgetRenderer{
		imageWidget: w,
		img:         img,
		objects:     []fyne.CanvasObject{img},
	}
}

/*  Renderer Implementation must

type WidgetRenderer interface {
	Layout(Size)
	MinSize() Size

	Refresh()
	BackgroundColor() color.Color
	Objects() []CanvasObject
	Destroy()
}
*/

type imageWidgetRenderer struct {
	imageWidget *ImageWidget
	img         *canvas.Image
	objects     []fyne.CanvasObject
}

// Layout for the renderer
func (r *imageWidgetRenderer) Layout(size fyne.Size) {
	// dont need to do anything here ? we just have 1 image
}

// MinSize that this widget can be
func (r *imageWidgetRenderer) MinSize() fyne.Size {
	return fyne.NewSize(100,100)
}

// Refresh the canvas
func (r *imageWidgetRenderer) Refresh() {
	canvas.Refresh(r.img)
}

// BackgroundColor for this rendererr
func (r *imageWidgetRenderer) BackgroundColor() color.Color {
	return theme.PrimaryColor()
	//return theme.BackgroundColor()
}

// Objects list for this renderer
func (r *imageWidgetRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// Destroy this renderer
func (r *imageWidgetRenderer) Destroy() {}
