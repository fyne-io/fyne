package software

import (
	"image"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/driver"
	"fyne.io/fyne/v2/internal/scale"
)

// Painter is a simple software painter that can paint a canvas in memory.
type Painter struct{}

// NewPainter creates a new Painter.
func NewPainter() *Painter {
	return &Painter{}
}

// Paint is the main entry point for a simple software painter.
// The canvas to be drawn is passed in as a parameter and the return is an
// image containing the result of rendering.
func (p *Painter) Paint(c fyne.Canvas) image.Image {
	bounds := image.Rect(0, 0, scale.ToScreenCoordinate(c, c.Size().Width), scale.ToScreenCoordinate(c, c.Size().Height))
	base := image.NewNRGBA(bounds)
	p.paintInto(c, base)
	return base
}

// PaintScratch renders the canvas like Paint but into a frame buffer taken
// from an internal pool. The caller must hand the image back with
// ReleaseScratch once it has been consumed and must not retain it afterwards.
// It is intended for callers that consume the frame immediately, such as a
// software canvas capture that composites it onto its own image.
func (p *Painter) PaintScratch(c fyne.Canvas) *image.NRGBA {
	bounds := image.Rect(0, 0, scale.ToScreenCoordinate(c, c.Size().Width), scale.ToScreenCoordinate(c, c.Size().Height))
	need := bounds.Dx() * bounds.Dy() * 4
	base, _ := framePool.Get().(*image.NRGBA)
	if base == nil || cap(base.Pix) < need {
		base = image.NewNRGBA(bounds)
	} else {
		base.Pix = base.Pix[:need]
		clear(base.Pix)
		base.Stride = bounds.Dx() * 4
		base.Rect = bounds
	}
	p.paintInto(c, base)
	return base
}

// ReleaseScratch returns a frame obtained from PaintScratch to the pool.
// Frames larger than maxPooledFrameBytes are dropped so that one oversized
// capture does not pin a large buffer in memory for the life of the process.
func (*Painter) ReleaseScratch(img *image.NRGBA) {
	if img == nil || cap(img.Pix) > maxPooledFrameBytes {
		return
	}
	framePool.Put(img)
}

// maxPooledFrameBytes bounds the size of frames kept in the pool
// (a 3840x2160 NRGBA frame).
const maxPooledFrameBytes = 3840 * 2160 * 4

var framePool sync.Pool

func (*Painter) paintInto(c fyne.Canvas, base *image.NRGBA) {

	paint := func(obj fyne.CanvasObject, pos, clipPos fyne.Position, clipSize fyne.Size) bool {
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
			drawImage(c, o, pos, base, clip)
		case *canvas.Text:
			drawText(c, o, pos, base, clip)
		case gradient:
			drawGradient(c, o, pos, base, clip)
		case *canvas.Circle:
			drawCircle(c, o, pos, base, clip)
		case *canvas.Line:
			drawLine(c, o, pos, base, clip)
		case *canvas.RegularPolygon:
			drawPolygon(c, o, pos, base, clip)
		case *canvas.ArbitraryPolygon:
			drawArbitraryPolygon(c, o, pos, base, clip)
		case *canvas.Raster:
			drawRaster(c, o, pos, base, clip)
		case *canvas.Rectangle:
			drawRectangle(c, o, pos, base, clip)
		case *canvas.Blur:
			drawBlur(c, o, pos, base, clip)
		case *canvas.Arc:
			drawArc(c, o, pos, base, clip)
		case *canvas.BezierCurve:
			drawBezierCurve(c, o, pos, base, clip)
		case *canvas.Ellipse:
			drawEllipse(c, o, pos, base, clip)
		}

		return false
	}

	driver.WalkVisibleObjectTree(c.Content(), paint, nil)
	for _, o := range c.Overlays().List() {
		driver.WalkVisibleObjectTree(o, paint, nil)
	}
}
