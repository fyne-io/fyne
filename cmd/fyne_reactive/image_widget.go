package main

import (
	"image/color"
	"sync"

	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"

	_ "image/jpeg"
	_ "image/png"

	"fyne.io/fyne"
	"fyne.io/fyne/dataapi"
	"fyne.io/fyne/widget"
)

// ImageCache is a simple sync map of URL name to contents
type ImageCache struct {
	res sync.Map
}

func newImageCache() *ImageCache {
	return &ImageCache{}
}

// Get fetches the given resource from the URL, with caching
func (c *ImageCache) Get(urlStr string) (fyne.Resource, error) {
	// if there, return it
	if v, ok := c.res.Load(urlStr); ok {
		return v.(fyne.Resource), nil
	}
	// get it over the network
	img, err := fyne.LoadResourceFromURLString(urlStr)
	if err != nil {
		return theme.CancelIcon(), err
	}
	// save it
	c.res.Store(urlStr, img)
	return img, nil
}

// ImageWidget is a custom widget demo
type ImageWidget struct {
	widget.BaseWidget
	widget.DataListener
	cache         *ImageCache
	urlStr        string
	imgRes        fyne.Resource
	img           *canvas.Image
	LoadErr       error
	UpdateBinding func(string)
	hovered       bool
}

// NewImageWidget returns a new ImageWidget
func NewImageWidget(cache *ImageCache) *ImageWidget {
	return &ImageWidget{
		BaseWidget:   widget.BaseWidget{},
		DataListener: widget.DataListener{},
		cache:        cache,
		img:          canvas.NewImageFromResource(theme.FyneLogo()),
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
	w.load()
}

func (w *ImageWidget) load() {
	if !w.BaseWidget.Hidden {
		w.imgRes, w.LoadErr = w.cache.Get(w.urlStr)
		w.img.Resource = w.imgRes
		w.Refresh()
	}
}

// SetURL updates the url, and propogates the change to the bound data
func (w *ImageWidget) SetURL(urlStr string) {
	w.urlStr = urlStr
	if w.UpdateBinding != nil {
		w.UpdateBinding(urlStr)
	}
	w.load()
}

// Tapped for clicks with the main button
func (w *ImageWidget) Tapped(ev *fyne.PointEvent) {
	println("tapped main")
	w.SetURL(FyneAvatarKangaroo)
}

// TappedSecondary for clicks with the other button
func (w *ImageWidget) TappedSecondary(ev *fyne.PointEvent) {
	println("tapped 2nd")
	w.SetURL("invalid")
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

// Img safely returns the image for the renderer to use
func (w *ImageWidget) Img() *canvas.Image {
	if w == nil {
		return nil
	}
	return w.img
}

type imageWidgetRenderer struct {
	imageWidget *ImageWidget
	objects     []fyne.CanvasObject
}

func (r *imageWidgetRenderer) setImg(img *canvas.Image) {
	r.objects[0] = img
}

// Layout for the renderer
func (r *imageWidgetRenderer) Layout(size fyne.Size) {
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
