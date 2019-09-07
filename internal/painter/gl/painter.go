// Package gl provides a full Fyne render implementation using system OpenGL libraries.
package gl

import (
	"image"
	"math"

	"fyne.io/fyne"
	"fyne.io/fyne/internal/driver"
	"fyne.io/fyne/widget"
)

// Painter defines the functionality of our OpenGL based renderer
type Painter interface {
	// Init tell a new painter to initialise, usually called after a context is available
	Init()
	// SetOutputSize is used to change the resolution of our output viewport
	SetOutputSize(int, int)
	// SetFrameBufferScale tells us when we have more than 1 framebuffer pixel for each output pixel
	SetFrameBufferScale(float32)
	// Paint is the main render method for this painter
	Paint(fyne.CanvasObject, fyne.Canvas, fyne.Size)
	// Clear tells our painter to prepare a fresh paint
	Clear()
	// Free is used to indicate that a certain canvas object is no longer needed
	Free(fyne.CanvasObject)
	// Capture requests that the specified canvas be drawn to an in-memory image
	Capture(fyne.Canvas) image.Image
}

// Declare conformity to Painter interface
var _ Painter = (*glPainter)(nil)

type glPainter struct {
	canvas   fyne.Canvas
	context  driver.WithContext
	program  Program
	texScale float32
}

func (p *glPainter) SetFrameBufferScale(scale float32) {
	p.texScale = scale
}

func (p *glPainter) Clear() {
	p.glClearBuffer()
}

func (p *glPainter) Paint(co fyne.CanvasObject, c fyne.Canvas, size fyne.Size) {
	if co == nil {
		return
	}

	paint := func(obj fyne.CanvasObject, pos fyne.Position, _ fyne.Position, _ fyne.Size) bool {
		// TODO should this be somehow not scroll container specific?
		if _, ok := obj.(*widget.ScrollContainer); ok {
			scrollX := p.textureScaleInt(pos.X)
			scrollY := p.textureScaleInt(pos.Y)
			scrollWidth := p.textureScaleInt(obj.Size().Width)
			scrollHeight := p.textureScaleInt(obj.Size().Height)
			pixHeight := p.textureScaleInt(co.Size().Height)
			p.glScissorOpen(int32(scrollX), int32(pixHeight-scrollY-scrollHeight), int32(scrollWidth), int32(scrollHeight))
		}
		if obj.Visible() {
			p.drawObject(obj, pos, size)
		}
		return false
	}
	afterPaint := func(obj, _ fyne.CanvasObject) {
		if _, ok := obj.(*widget.ScrollContainer); ok {
			p.glScissorClose()
		}
	}

	driver.WalkVisibleObjectTree(co, paint, afterPaint)
}

func (p *glPainter) Free(obj fyne.CanvasObject) {
	p.freeTexture(obj)
}

func (p *glPainter) textureScaleInt(v int) int {
	if p.canvas.Scale() == 1.0 && p.texScale == 1.0 {
		return v
	}
	return int(math.Round(float64(v) * float64(p.canvas.Scale()*p.texScale)))
}

func unscaleInt(c fyne.Canvas, v int) int {
	switch c.Scale() {
	case 1.0:
		return v
	default:
		return int(float32(v) / c.Scale())
	}
}

// NewPainter creates a new GL based renderer for the provided canvas.
// If it is a master painter it will also initialise OpenGL
func NewPainter(c fyne.Canvas, ctx driver.WithContext) Painter {
	p := &glPainter{canvas: c, context: ctx}
	p.texScale = 1.0
	go svgCacheMonitorTheme()

	glInit()
	return p
}
