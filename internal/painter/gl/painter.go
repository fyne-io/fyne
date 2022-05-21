// Package gl provides a full Fyne render implementation using system OpenGL libraries.
package gl

import (
	"fmt"
	"image"
	"image/draw"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/driver"
	"fyne.io/fyne/v2/theme"
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

var shaderSources = map[string][2][]byte{
	"line":      {shaderLineVert.StaticContent, shaderLineFrag.StaticContent},
	"line_es":   {shaderLineesVert.StaticContent, shaderLineesFrag.StaticContent},
	"simple":    {shaderSimpleVert.StaticContent, shaderSimpleFrag.StaticContent},
	"simple_es": {shaderSimpleesVert.StaticContent, shaderSimpleesFrag.StaticContent},
}

func (p *painter) SetFrameBufferScale(scale float32) {
	p.texScale = scale
	p.pixScale = p.canvas.Scale() * p.texScale
}

func (p *painter) Clear() {
	r, g, b, a := theme.BackgroundColor().RGBA()
	p.ctx.ClearColor(float32(r)/max16bit, float32(g)/max16bit, float32(b)/max16bit, float32(a)/max16bit)
	p.ctx.Clear(bitColorBuffer | bitDepthBuffer)
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

func (p *painter) Paint(obj fyne.CanvasObject, pos fyne.Position, frame fyne.Size) {
	if obj.Visible() {
		p.drawObject(obj, pos, frame)
	}
}

func (p *painter) Free(obj fyne.CanvasObject) {
	p.freeTexture(obj)
}

func (p *painter) SetOutputSize(width, height int) {
	p.ctx.Viewport(0, 0, width, height)
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

func (p *painter) createBuffer(points []float32) Buffer {
	vbo := p.ctx.CreateBuffer()
	p.logError()
	p.ctx.BindBuffer(arrayBuffer, vbo)
	p.logError()
	p.ctx.BufferData(arrayBuffer, points, staticDraw)
	p.logError()
	return vbo
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

func (p *painter) defineVertexArray(prog Program, name string, size, stride, offset int) {
	vertAttrib := p.ctx.GetAttribLocation(prog, name)
	p.ctx.EnableVertexAttribArray(vertAttrib)
	p.ctx.VertexAttribPointerWithOffset(vertAttrib, size, float, false, stride, offset)
	p.logError()
}

func (p *painter) freeBuffer(vbo Buffer) {
	p.ctx.BindBuffer(arrayBuffer, noBuffer)
	p.logError()
	p.ctx.DeleteBuffer(vbo)
	p.logError()
}

func (p *painter) freeTexture(obj fyne.CanvasObject) {
	texture, ok := cache.GetTexture(obj)
	if !ok {
		return
	}

	p.ctx.DeleteTexture(Texture(texture))
	p.logError()
	cache.DeleteTexture(obj)
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
	return p
}
