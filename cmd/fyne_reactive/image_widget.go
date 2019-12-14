package main

import (
	"image/color"

	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"

	_ "image/jpeg"
	_ "image/png"

	"fyne.io/fyne"
	"fyne.io/fyne/dataapi"
	"fyne.io/fyne/widget"
)

// ImageWidget is a custom widget demo
type ImageWidget struct {
	widget.BaseWidget
	widget.DataListener
	urlStr  string
	imgRes  fyne.Resource
	img     *canvas.Image
	LoadErr error
	OnBind  func(string)
}

// NewImageWidget returns a new ImageWidget
func NewImageWidget(urlStr string) *ImageWidget {
	var img *canvas.Image
	imgRes, err := fyne.LoadResourceFromURLString(urlStr)
	if err != nil {
		img = canvas.NewImageFromResource(theme.CancelIcon())
		println("Failed to load", urlStr, err.Error())
	} else {
		img = canvas.NewImageFromResource(imgRes)
	}
	return &ImageWidget{
		BaseWidget:   widget.BaseWidget{},
		DataListener: widget.DataListener{},
		urlStr:       urlStr,
		imgRes:       imgRes,
		LoadErr:      err,
		img:          img,
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
	if w.LoadErr != nil {
		w.img.Resource = theme.CancelIcon()
	} else {
		w.img.Resource = w.imgRes
	}
	w.Refresh()
}

// CreateRenderer creates and returns the renderer for this widget type
func (w *ImageWidget) CreateRenderer() fyne.WidgetRenderer {
	w.ExtendBaseWidget(w)
	if w.img == nil {
		w.img = canvas.NewImageFromResource(theme.CancelIcon())
	}
	return &imageWidgetRenderer{
		imageWidget: w,
		objects:     []fyne.CanvasObject{w.img},
	}
}

func (w *ImageWidget) Img() *canvas.Image {
	if w == nil {
		return nil
	}
	return w.img
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
	objects     []fyne.CanvasObject
}

func (r *imageWidgetRenderer) setImg(img *canvas.Image) {
	r.objects[0] = img
}

// Layout for the renderer
func (r *imageWidgetRenderer) Layout(size fyne.Size) {
	// dont need to do anything here ? we just have 1 image
	if size.Width > 320 {
		size.Width = 320
	}
	if size.Height > 200 {
		size.Height = 200
	}

	r.imageWidget.img.Resize(size)
}

// MinSize that this widget can be
func (r *imageWidgetRenderer) MinSize() fyne.Size {
	return fyne.NewSize(320, 200)
}

// Refresh the canvas
func (r *imageWidgetRenderer) Refresh() {
	if i := r.imageWidget.Img(); i != nil {
		canvas.Refresh(i)
	}
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
