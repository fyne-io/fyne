// Package gl provides a full Fyne render implementation using system OpenGL libraries.
package gl

import (
	"fmt"
	"image"
	"image/draw"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/driver"
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
var _ Painter = (*painter)(nil)

type painter struct {
	canvas          fyne.Canvas
	ctx             context
	contextProvider driver.WithContext
	program         Program
	lineProgram     Program
	texScale        float32
	pixScale        float32 // pre-calculate scale*texScale for each draw
}

func (p *painter) SetFrameBufferScale(scale float32) {
	p.texScale = scale
	p.pixScale = p.canvas.Scale() * p.texScale
}

func (p *painter) Clear() {
	p.glClearBuffer()
}

func (p *painter) StartClipping(pos fyne.Position, size fyne.Size) {
	x := p.textureScale(pos.X)
	y := p.textureScale(p.canvas.Size().Height - pos.Y - size.Height)
	w := p.textureScale(size.Width)
	h := p.textureScale(size.Height)
	p.glScissorOpen(int32(x), int32(y), int32(w), int32(h))
}

func (p *painter) StopClipping() {
	p.glScissorClose()
}

func (p *painter) Paint(obj fyne.CanvasObject, pos fyne.Position, frame fyne.Size) {
	if obj.Visible() {
		p.drawObject(obj, pos, frame)
	}
}

func (p *painter) Free(obj fyne.CanvasObject) {
	p.freeTexture(obj)
}

func (p *painter) imgToTexture(img image.Image, textureFilter canvas.ImageScale) Texture {
	switch i := img.(type) {
	case *image.Uniform:
		texture := p.newTexture(textureFilter)
		r, g, b, a := i.RGBA()
		r8, g8, b8, a8 := uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)
		data := []uint8{r8, g8, b8, a8}
		p.ctx.TexImage2D(
			texture2D,
			0,
			1,
			1,
			colorFormatRGBA,
			unsignedByte,
			data,
		)
		p.logError()
		return texture
	case *image.RGBA:
		if len(i.Pix) == 0 { // image is empty
			return noTexture
		}

		texture := p.newTexture(textureFilter)
		p.ctx.TexImage2D(
			texture2D,
			0,
			i.Rect.Size().X,
			i.Rect.Size().Y,
			colorFormatRGBA,
			unsignedByte,
			i.Pix,
		)
		p.logError()
		return texture
	default:
		rgba := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
		draw.Draw(rgba, rgba.Rect, img, image.Point{}, draw.Over)
		return p.imgToTexture(rgba, textureFilter)
	}
}

func (p *painter) logError() {
	logGLError(p.ctx.GetError())
}

func (p *painter) newTexture(textureFilter canvas.ImageScale) Texture {
	if int(textureFilter) >= len(textureFilterToGL) {
		fyne.LogError(fmt.Sprintf("Invalid canvas.ImageScale value (%d), using canvas.ImageScaleSmooth as default value", textureFilter), nil)
		textureFilter = canvas.ImageScaleSmooth
	}

	texture := p.ctx.CreateTexture()
	p.logError()
	p.ctx.ActiveTexture(texture0)
	p.ctx.BindTexture(texture2D, texture)
	p.logError()
	p.ctx.TexParameteri(texture2D, textureMinFilter, textureFilterToGL[textureFilter])
	p.ctx.TexParameteri(texture2D, textureMagFilter, textureFilterToGL[textureFilter])
	p.ctx.TexParameteri(texture2D, textureWrapS, clampToEdge)
	p.ctx.TexParameteri(texture2D, textureWrapT, clampToEdge)
	p.logError()

	return texture
}

func (p *painter) textureScale(v float32) float32 {
	if p.pixScale == 1.0 {
		return float32(math.Round(float64(v)))
	}

	return float32(math.Round(float64(v * p.pixScale)))
}

// NewPainter creates a new GL based renderer for the provided canvas.
// If it is a master painter it will also initialise OpenGL
func NewPainter(c fyne.Canvas, ctx driver.WithContext) Painter {
	p := &painter{canvas: c, contextProvider: ctx}
	p.SetFrameBufferScale(1.0)

	p.glInit()

	return p
}
