// Package gl provides a full Fyne render implementation using system OpenGL libraries.
package gl

import (
	"image"
	"math"
	"sync"

	"fyne.io/fyne"
	"fyne.io/fyne/internal/driver"
	"fyne.io/fyne/internal/painter"
)

// Painter defines the functionality of our OpenGL based renderer
type Painter interface {
	// Init tell a new painter to initialise, usually called after a context is available
	Init()
	// Capture requests that the specified canvas be drawn to an in-memory image
	Capture(fyne.Canvas) image.Image
	// Clear tells our painter to prepare a fresh paint
	Clear()
	// Free is used to indicate that a certain canvas object is no longer needed
	Free(fyne.CanvasObject)
	// Paint a single fyne.CanvasObject but not its children.
	Paint(fyne.CanvasObject, fyne.Position, fyne.Size)
	// SetFrameBufferScale tells us when we have more than 1 framebuffer pixel for each output pixel
	SetFrameBufferScale(float32)
	// SetOutputSize is used to change the resolution of our output viewport
	SetOutputSize(int, int)
	// StartClipping tells us that the following paint actions should be clipped to the specified area.
	StartClipping(fyne.Position, fyne.Size)
	// StopClipping stops clipping paint actions.
	StopClipping()
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

func (p *glPainter) StartClipping(pos fyne.Position, size fyne.Size) {
	x := p.textureScaleInt(pos.X)
	y := p.textureScaleInt(pos.Y)
	w := p.textureScaleInt(size.Width)
	h := p.textureScaleInt(size.Height)
	p.glScissorOpen(int32(x), int32(y), int32(w), int32(h))
}

func (p *glPainter) StopClipping() {
	p.glScissorClose()
}

func (p *glPainter) Paint(obj fyne.CanvasObject, pos fyne.Position, frame fyne.Size) {
	if obj.Visible() {
		p.drawObject(obj, pos, frame)
	}
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

func (p *glPainter) textureScale(v float32) int {
	if p.canvas.Scale() == 1.0 && p.texScale == 1.0 {
		return int(v)
	}
	return int(math.Round(float64(v) * float64(p.canvas.Scale()*p.texScale)))
}

var startCacheMonitor = &sync.Once{}

// NewPainter creates a new GL based renderer for the provided canvas.
// If it is a master painter it will also initialise OpenGL
func NewPainter(c fyne.Canvas, ctx driver.WithContext) Painter {
	p := &glPainter{canvas: c, context: ctx}
	p.texScale = 1.0

	glInit()
	startCacheMonitor.Do(func() {
		go painter.SvgCacheMonitorTheme()
	})

	return p
}
