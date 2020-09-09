package widget

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/theme"
)

// RGBAHolder defines getters and setters for RGBA colors.
type RGBAHolder interface {
	RGBA() (float64, float64, float64, float64)
	SetRGBA(float64, float64, float64, float64)
}

// NewRGBAColorArea returns a new RGBA color area using the given holder.
func NewRGBAColorArea(holder RGBAHolder) *ColorArea {
	lookup := func(x, y, w, h int) color.Color {
		// TODO translate x,y into Red/Green/Blue
		return nil
	}
	selection := func(w, h int) (x, y int) {
		// TODO translate current color into x,y coords
		return
	}
	return NewColorArea(lookup, selection, func(x, y, w, h int) {
		// TODO if in bounds
		// TODO   holder.SetRGBA
	})
}

// HSLAHolder defines getters and setters for HSLA colors.
type HSLAHolder interface {
	HSLA() (float64, float64, float64, float64)
	SetHSLA(float64, float64, float64, float64)
}

// NewHSLAColorArea returns a new HSLA color area using the given holder.
func NewHSLAColorArea(holder HSLAHolder) *ColorArea {
	lookup := func(x, y, w, h int) color.Color {
		angle, radius, limit := cartesianToPolar(float64(x), float64(y), float64(w), float64(h))
		if radius > limit {
			// Out of bounds
			return theme.BackgroundColor()
		}
		_, _, lightness, alpha := holder.HSLA()
		hue, saturation := polarToHS(angle, radius, limit)
		red, green, blue := hslToRgb(hue, saturation, lightness)
		return &color.NRGBA{
			R: uint8(red * 255.0),
			G: uint8(green * 255.0),
			B: uint8(blue * 255.0),
			A: uint8(alpha * 255.0),
		}
	}
	selection := func(w, h int) (int, int) {
		hue, saturation, _, _ := holder.HSLA()
		angle := hue * 2 * math.Pi
		limit := math.Min(float64(w), float64(h)) / 2.0
		radius := saturation * limit
		x, y := polarToCartesian(angle, radius, float64(w), float64(h))
		return int(x), int(y)
	}
	return NewColorArea(lookup, selection, func(x, y, w, h int) {
		angle, radius, limit := cartesianToPolar(float64(x), float64(y), float64(w), float64(h))
		if radius > limit {
			// Out of bounds
			return
		}
		_, _, lightness, alpha := holder.HSLA()
		hue, saturation := polarToHS(angle, radius, limit)
		holder.SetHSLA(hue, saturation, lightness, alpha)
	})
}

var _ fyne.Widget = (*ColorArea)(nil)
var _ fyne.Tappable = (*ColorArea)(nil)
var _ fyne.Draggable = (*ColorArea)(nil)

// ColorArea displays a color gradient and triggers the callback when tapped.
type ColorArea struct {
	BaseWidget
	generator func(w, h int) image.Image
	cache     draw.Image
	minSize   fyne.Size
	selection func(int, int) (int, int)
	onChange  func(int, int, int, int)
}

// NewColorArea returns a new color area with the given lookup and selection callbacks, and triggers the given onChange callback.
func NewColorArea(lookup func(int, int, int, int) color.Color, selection func(int, int) (int, int), onChange func(int, int, int, int)) *ColorArea {
	a := &ColorArea{
		selection: selection,
		onChange:  onChange,
	}
	a.generator = func(w, h int) image.Image {
		if a.cache == nil || a.cache.Bounds().Dx() != w || a.cache.Bounds().Dy() != h {
			rect := image.Rect(0, 0, w, h)
			a.cache = image.NewRGBA(rect)
		}
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				if c := lookup(x, y, w, h); c != nil {
					a.cache.Set(x, y, c)
				}
			}
		}
		return a.cache
	}
	a.ExtendBaseWidget(a)
	return a
}

// Cursor returns the cursor type of this widget
func (a *ColorArea) Cursor() desktop.Cursor {
	return desktop.CrosshairCursor
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (a *ColorArea) CreateRenderer() fyne.WidgetRenderer {
	raster := &canvas.Raster{
		Generator: a.generator,
	}
	x := canvas.NewLine(color.Black)
	y := canvas.NewLine(color.Black)
	return &colorAreaRenderer{
		BaseRenderer: widget.NewBaseRenderer([]fyne.CanvasObject{raster, x, y}),
		area:         a,
		raster:       raster,
		x:            x,
		y:            y,
	}
}

// MinSize returns the size that this widget should not shrink below
func (a *ColorArea) MinSize() (min fyne.Size) {
	a.ExtendBaseWidget(a)
	a.propertyLock.RLock()
	min = a.minSize
	a.propertyLock.RUnlock()
	if min.IsZero() {
		min = a.BaseWidget.MinSize()
	}
	return
}

// SetMinSize specifies the smallest size this object should be
func (a *ColorArea) SetMinSize(size fyne.Size) {
	a.propertyLock.Lock()
	defer a.propertyLock.Unlock()

	a.minSize = size
}

// Tapped is called when a pointer tapped event is captured and triggers any change handler
func (a *ColorArea) Tapped(event *fyne.PointEvent) {
	x, y := a.locationForPosition(event.Position)
	if c := a.cache; c != nil {
		b := c.Bounds()
		if f := a.onChange; f != nil {
			f(x, y, b.Dx(), b.Dy())
		}
	}
}

// Dragged is called when a pointer drag event is captured and triggers any change handler
func (a *ColorArea) Dragged(event *fyne.DragEvent) {
	x, y := a.locationForPosition(event.Position)
	if c := a.cache; c != nil {
		b := c.Bounds()
		if f := a.onChange; f != nil {
			f(x, y, b.Dx(), b.Dy())
		}
	}
}

// DragEnd is called when a pointer drag ends
func (a *ColorArea) DragEnd() {
}

func (a *ColorArea) locationForPosition(pos fyne.Position) (x, y int) {
	can := fyne.CurrentApp().Driver().CanvasForObject(a)
	x, y = pos.X, pos.Y
	if can != nil {
		x, y = can.PixelCoordinateForPosition(pos)
	}
	return
}

type colorAreaRenderer struct {
	widget.BaseRenderer
	area   *ColorArea
	raster *canvas.Raster
	x, y   *canvas.Line
}

func (r *colorAreaRenderer) Layout(size fyne.Size) {
	p := theme.Padding()
	w := size.Width - 2*p
	h := size.Height - 2*p
	if f := r.area.selection; f != nil {
		x, y := f(w, h)
		r.x.Position1 = fyne.NewPos(p, y+p)
		r.x.Position2 = fyne.NewPos(w+p, y+p)
		r.y.Position1 = fyne.NewPos(x+p, p)
		r.y.Position2 = fyne.NewPos(x+p, h+p)
	}
	r.raster.Move(fyne.NewPos(p, p))
	r.raster.Resize(fyne.NewSize(w, h))
	return
}

func (r *colorAreaRenderer) MinSize() (min fyne.Size) {
	min = r.raster.MinSize()
	size := 2 * theme.IconInlineSize()
	min = min.Max(fyne.NewSize(size, size))
	min = min.Add(fyne.NewSize(2*theme.Padding(), 2*theme.Padding()))
	return
}

func (r *colorAreaRenderer) Refresh() {
	s := r.area.Size()
	if s.IsZero() {
		r.area.Resize(r.area.MinSize())
	} else {
		r.Layout(s)
	}
	r.x.StrokeColor = theme.IconColor()
	r.x.Refresh()
	r.y.StrokeColor = theme.IconColor()
	r.y.Refresh()
	r.raster.Refresh()
	canvas.Refresh(r.area.super())
}
