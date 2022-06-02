// Package gl provides a full Fyne render implementation using system OpenGL libraries.
package gl

import (
	"fmt"
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/driver"
	"fyne.io/fyne/v2/theme"
)

var shaderSources = map[string][2][]byte{
	"line":      {shaderLineVert.StaticContent, shaderLineFrag.StaticContent},
	"line_es":   {shaderLineesVert.StaticContent, shaderLineesFrag.StaticContent},
	"simple":    {shaderSimpleVert.StaticContent, shaderSimpleFrag.StaticContent},
	"simple_es": {shaderSimpleesVert.StaticContent, shaderSimpleesFrag.StaticContent},
}

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

// NewPainter creates a new GL based renderer for the provided canvas.
// If it is a master painter it will also initialise OpenGL
func NewPainter(c fyne.Canvas, ctx driver.WithContext) Painter {
	p := &painter{canvas: c, contextProvider: ctx}
	p.SetFrameBufferScale(1.0)
	return p
}

type painter struct {
	canvas          fyne.Canvas
	ctx             context
	contextProvider driver.WithContext
	program         Program
	lineProgram     Program
	texScale        float32
	pixScale        float32 // pre-calculate scale*texScale for each draw
}

// Declare conformity to Painter interface
var _ Painter = (*painter)(nil)

func (p *painter) Clear() {
	r, g, b, a := theme.BackgroundColor().RGBA()
	p.ctx.ClearColor(float32(r)/max16bit, float32(g)/max16bit, float32(b)/max16bit, float32(a)/max16bit)
	p.ctx.Clear(bitColorBuffer | bitDepthBuffer)
	p.logError()
}

func (p *painter) Free(obj fyne.CanvasObject) {
	p.freeTexture(obj)
}

func (p *painter) Paint(obj fyne.CanvasObject, pos fyne.Position, frame fyne.Size) {
	if obj.Visible() {
		p.drawObject(obj, pos, frame)
	}
}

func (p *painter) SetFrameBufferScale(scale float32) {
	p.texScale = scale
	p.pixScale = p.canvas.Scale() * p.texScale
}

func (p *painter) SetOutputSize(width, height int) {
	p.ctx.Viewport(0, 0, width, height)
	p.logError()
}

func (p *painter) StartClipping(pos fyne.Position, size fyne.Size) {
	x := p.textureScale(pos.X)
	y := p.textureScale(p.canvas.Size().Height - pos.Y - size.Height)
	w := p.textureScale(size.Width)
	h := p.textureScale(size.Height)
	p.ctx.Scissor(int32(x), int32(y), int32(w), int32(h))
	p.ctx.Enable(scissorTest)
	p.logError()
}

func (p *painter) StopClipping() {
	p.ctx.Disable(scissorTest)
	p.logError()
}

func (p *painter) compileShader(source string, shaderType uint32) (Shader, error) {
	shader := p.ctx.CreateShader(shaderType)

	p.ctx.ShaderSource(shader, source)
	p.logError()
	p.ctx.CompileShader(shader)
	p.logError()

	info := p.ctx.GetShaderInfoLog(shader)
	if p.ctx.GetShaderi(shader, compileStatus) == glFalse {
		return noShader, fmt.Errorf("failed to compile OpenGL shader:\n%s\n>>> SHADER SOURCE\n%s\n<<< SHADER SOURCE", info, source)
	}

	// The info is probably a null terminated string.
	// An empty info has been seen as "\x00".
	if len(info) > 0 && info != "\x00" {
		fmt.Printf("OpenGL shader compilation output:\n%s\n>>> SHADER SOURCE\n%s\n<<< SHADER SOURCE\n", info, source)
	}

	return shader, nil
}

func (p *painter) createProgram(shaderFilename string) Program {
	// Why a switch over a filename?
	// Because this allows for a minimal change, once we reach Go 1.16 and use go:embed instead of
	// fyne bundle.
	sources := shaderSources[shaderFilename]
	vertexSrc, fragmentSrc := sources[0], sources[1]
	if vertexSrc == nil {
		panic("shader not found: " + shaderFilename)
	}

	vertShader, err := p.compileShader(string(vertexSrc), vertexShader)
	if err != nil {
		panic(err)
	}
	fragShader, err := p.compileShader(string(fragmentSrc), fragmentShader)
	if err != nil {
		panic(err)
	}

	prog := p.ctx.CreateProgram()
	p.ctx.AttachShader(prog, vertShader)
	p.ctx.AttachShader(prog, fragShader)
	p.ctx.LinkProgram(prog)

	info := p.ctx.GetProgramInfoLog(prog)
	if p.ctx.GetProgrami(prog, linkStatus) == glFalse {
		panic(fmt.Errorf("failed to link OpenGL program:\n%s", info))
	}

	// The info is probably a null terminated string.
	// An empty info has been seen as "\x00".
	if len(info) > 0 && info != "\x00" {
		fmt.Printf("OpenGL program linking output:\n%s\n", info)
	}

	if glErr := p.ctx.GetError(); glErr != 0 {
		panic(fmt.Sprintf("failed to link OpenGL program; error code: %x", glErr))
	}

	return prog
}

func (p *painter) logError() {
	logGLError(p.ctx.GetError())
}
