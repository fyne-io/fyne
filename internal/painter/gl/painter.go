// Package gl provides a full Fyne render implementation using system OpenGL libraries.
package gl

import (
	"fmt"
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal"
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
	Paint(fyne.CanvasObject, fyne.Position, fyne.Size, *internal.ClipItem)
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
	canvas                  fyne.Canvas
	ctx                     context
	contextProvider         driver.WithContext
	program                 ProgramState
	blurProgram             ProgramState
	lineProgram             ProgramState
	rectangleProgram        ProgramState
	roundRectangleProgram   ProgramState
	polygonProgram          ProgramState
	arcProgram              ProgramState
	bezierCurveProgram      ProgramState
	arbitraryPolygonProgram ProgramState
	texScale                float32
	pixScale                float32 // pre-calculate scale*texScale for each draw
	blurSnapTex             Texture // cached texture for GPU-side blur snapshot
	blurSnapTexValid        bool    // whether blurSnapTex has been allocated
	blurSnapW, blurSnapH    int     // size of blurSnapTex in pixels
	fbHeight                int     // current framebuffer height in pixels
}

type ProgramState struct {
	ref        Program
	buff       Buffer
	uniforms   map[string]*UniformState
	attributes map[string]Attribute
}

type UniformState struct {
	ref   Uniform
	prev  [4]float32
	prevv []float32
}

func (p *painter) SetUniform1f(pState ProgramState, name string, v float32) {
	u := p.getUniformLocation(pState, name)
	if u.prev[0] == v {
		return
	}
	u.prev[0] = v
	p.ctx.Uniform1f(u.ref, v)
}

func (p *painter) SetUniform1fv(pState ProgramState, name string, v []float32) {
	u := p.getUniformLocation(pState, name)
	if float32SlicesEqual(u.prevv, v) {
		return
	}
	u.prevv = append(u.prevv[:0], v...)
	p.ctx.Uniform1fv(u.ref, v)
}

func (p *painter) SetUniform2f(pState ProgramState, name string, v0, v1 float32) {
	u := p.getUniformLocation(pState, name)
	if u.prev[0] == v0 && u.prev[1] == v1 {
		return
	}
	u.prev[0] = v0
	u.prev[1] = v1
	p.ctx.Uniform2f(u.ref, v0, v1)
}

func (p *painter) SetUniform2fv(pState ProgramState, name string, v []float32) {
	u := p.getUniformLocation(pState, name)
	if float32SlicesEqual(u.prevv, v) {
		return
	}
	u.prevv = append(u.prevv[:0], v...)
	p.ctx.Uniform2fv(u.ref, v)
}

func (p *painter) SetUniform4f(pState ProgramState, name string, v0, v1, v2, v3 float32) {
	u := p.getUniformLocation(pState, name)
	if u.prev[0] == v0 && u.prev[1] == v1 && u.prev[2] == v2 && u.prev[3] == v3 {
		return
	}
	u.prev[0] = v0
	u.prev[1] = v1
	u.prev[2] = v2
	u.prev[3] = v3
	p.ctx.Uniform4f(u.ref, v0, v1, v2, v3)
}

func (p *painter) UpdateVertexArray(pState ProgramState, name string, size, stride, offset int) {
	a := p.enableAttribArray(pState, name)

	p.ctx.VertexAttribPointerWithOffset(a, size, float, false, stride*floatSize, offset*floatSize)
	p.logError()
}

// Declare conformity to Painter interface
var _ Painter = (*painter)(nil)

func (p *painter) Clear() {
	r, g, b, a := theme.Color(theme.ColorNameBackground).RGBA()
	p.ctx.ClearColor(float32(r)/max16bit, float32(g)/max16bit, float32(b)/max16bit, float32(a)/max16bit)
	p.ctx.Clear(bitColorBuffer | bitDepthBuffer)
	p.logError()
}

func (p *painter) Free(obj fyne.CanvasObject) {
	p.freeTexture(obj)
}

func (p *painter) Paint(obj fyne.CanvasObject, pos fyne.Position, frame fyne.Size, clip *internal.ClipItem) {
	if obj.Visible() {
		p.drawObject(obj, pos, frame, clip)
	}
}

func (p *painter) SetFrameBufferScale(scale float32) {
	p.texScale = scale
	p.pixScale = p.canvas.Scale() * p.texScale
}

func (p *painter) SetOutputSize(width, height int) {
	p.ctx.Viewport(0, 0, width, height)
	p.fbHeight = height
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
	// An empty info has been seen as "\x00" or "\x00\x00".
	if len(info) > 0 && info != "\x00" && info != "\x00\x00" {
		fmt.Printf("OpenGL shader compilation output:\n%s\n>>> SHADER SOURCE\n%s\n<<< SHADER SOURCE\n", info, source)
	}

	return shader, nil
}

func (p *painter) createProgram(shaderFilename string) Program {
	// Why a switch over a filename?
	// Because this allows for a minimal change, once we reach Go 1.16 and use go:embed instead of
	// fyne bundle.
	vertexSrc, fragmentSrc := shaderSourceNamed(shaderFilename)
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
	// An empty info has been seen as "\x00" or "\x00\x00".
	if len(info) > 0 && info != "\x00" && info != "\x00\x00" {
		fmt.Printf("OpenGL program linking output:\n%s\n", info)
	}

	if glErr := p.ctx.GetError(); glErr != 0 {
		panic(fmt.Sprintf("failed to link OpenGL program; error code: %x", glErr))
	}

	p.ctx.UseProgram(prog)

	return prog
}

func (p *painter) enableAttribArray(pState ProgramState, name string) Attribute {
	a, ok := pState.attributes[name]
	if !ok {
		a = p.ctx.GetAttribLocation(pState.ref, name)
		p.ctx.EnableVertexAttribArray(a)
		pState.attributes[name] = a
	}

	return a
}

func (p *painter) getUniformLocation(pState ProgramState, name string) *UniformState {
	u, ok := pState.uniforms[name]
	if !ok {
		u = &UniformState{ref: p.ctx.GetUniformLocation(pState.ref, name)}
		pState.uniforms[name] = u
	}

	return u
}

func (p *painter) logError() {
	logGLError(p.ctx.GetError)
}

func float32SlicesEqual(a, b []float32) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
