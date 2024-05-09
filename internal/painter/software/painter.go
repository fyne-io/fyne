package software

import (
	"fmt"
	"image"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/driver"
	"fyne.io/fyne/v2/internal/painter"
	"fyne.io/fyne/v2/internal/scale"
)

var _ painter.Painter = (*Painter)(nil)

// Painter is a simple software painter that can paint a canvas in memory.
type Painter struct {
	canvas     fyne.Canvas
	dirtyRects []image.Rectangle
}

// NewPainter creates a new Painter.
func NewPainter() *Painter {
	return &Painter{}
}

// NewPainterWithCanvas creates a new Painter with an existing fyne.Canvas.
func NewPainterWithCanvas(canvas fyne.Canvas) *Painter {
	return &Painter{
		canvas: canvas,
	}
}

// Capture is the main entry point for a simple software painter.
// The canvas to be drawn is passed in as a parameter and the return is an
// image containing the result of rendering.
func (p *Painter) Capture(c fyne.Canvas) image.Image {
	t := time.Now()

	bounds := image.Rect(0, 0, scale.ToScreenCoordinate(c, c.Size().Width), scale.ToScreenCoordinate(c, c.Size().Height))
	base := image.NewNRGBA(bounds)

	paint := func(obj fyne.CanvasObject, pos, clipPos fyne.Position, clipSize fyne.Size) bool {
		shouldTest := true
		shouldPaint := true

		switch obj.(type) {
		case *fyne.Container, fyne.Widget:
			shouldTest = true
		default:
			if _, ok := cache.GetTexture(obj); !ok {
				shouldTest = false
				shouldPaint = true
			}
		}

		if shouldTest {
			// TODO: This breaks a bunch of tests, because by default it'll return mostly blank images
			shouldPaint = driver.WalkVisibleObjectTree(obj, func(obj fyne.CanvasObject, _, _ fyne.Position, _ fyne.Size) bool {
				switch obj.(type) {
				case *fyne.Container, fyne.Widget:
					return false
				}
				if _, ok := cache.GetTexture(obj); !ok {
					return true
				}
				return false
			}, nil)
		}

		if shouldPaint {
			w := fyne.Min(clipPos.X+clipSize.Width, c.Size().Width)
			h := fyne.Min(clipPos.Y+clipSize.Height, c.Size().Height)
			clip := image.Rect(
				scale.ToScreenCoordinate(c, clipPos.X),
				scale.ToScreenCoordinate(c, clipPos.Y),
				scale.ToScreenCoordinate(c, w),
				scale.ToScreenCoordinate(c, h),
			)

			switch o := obj.(type) {
			case *canvas.Image:
				p.drawImage(c, o, pos, base, clip)
			case *canvas.Text:
				p.drawText(c, o, pos, base, clip)
			case gradient:
				p.drawGradient(c, o, pos, base, clip)
			case *canvas.Circle:
				p.drawCircle(c, o, pos, base, clip)
			case *canvas.Line:
				p.drawLine(c, o, pos, base, clip)
			case *canvas.Raster:
				p.drawRaster(c, o, pos, base, clip)
			case *canvas.Rectangle:
				p.drawRectangle(c, o, pos, base, clip)
			}
		}

		return !shouldPaint
	}

	driver.WalkVisibleObjectTreeIgnoreSibs(c.Content(), paint, nil)
	for _, o := range c.Overlays().List() {
		driver.WalkVisibleObjectTreeIgnoreSibs(o, paint, nil)
	}

	fmt.Println("Capture:", time.Since(t))
	return base
}

func (p *Painter) Init() {
	// TODO implement me
	panic("implement me")
}

func (p *Painter) Clear() {
	// TODO implement me
	panic("implement me")
}

func (p *Painter) Free(object fyne.CanvasObject) {
	p.freeTexture(object)
}

func (p *Painter) Paint(object fyne.CanvasObject, position fyne.Position, size fyne.Size) {
	// TODO implement me
	panic("implement me")
}

func (p *Painter) SetFrameBufferScale(f float32) {
	// TODO implement me
	panic("implement me")
}

func (p *Painter) SetOutputSize(i int, i2 int) {
	// TODO implement me
	panic("implement me")
}

func (p *Painter) StartClipping(position fyne.Position, size fyne.Size) {
	// TODO implement me
	panic("implement me")
}

func (p *Painter) StopClipping() {
	// TODO implement me
	panic("implement me")
}

func (p *Painter) ResetDirtyRects() {
	p.dirtyRects = nil
}

func (p *Painter) DirtyRects() []image.Rectangle {
	return p.dirtyRects
}
