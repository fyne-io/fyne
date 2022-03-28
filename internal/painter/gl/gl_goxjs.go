//go:build js || wasm || test_web_driver
// +build js wasm test_web_driver

package gl

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/fyne-io/gl-js"
	"golang.org/x/mobile/exp/f32"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/theme"
)

const (
	clampToEdge      = gl.CLAMP_TO_EDGE
	colorFormatRGBA  = gl.RGBA
	texture0         = gl.TEXTURE0
	texture2D        = gl.TEXTURE_2D
	textureMinFilter = gl.TEXTURE_MIN_FILTER
	textureMagFilter = gl.TEXTURE_MAG_FILTER
	textureWrapS     = gl.TEXTURE_WRAP_S
	textureWrapT     = gl.TEXTURE_WRAP_T
	unsignedByte     = gl.UNSIGNED_BYTE
)

// Buffer represents a GL buffer
type Buffer gl.Buffer

// Program represents a compiled GL program
type Program gl.Program

var textureFilterToGL = []int32{gl.LINEAR, gl.NEAREST}

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
		draw.Draw(rgba, rgba.Rect, img, image.ZP, draw.Over)
		return p.imgToTexture(rgba, textureFilter)
	}
}

func (p *painter) SetOutputSize(width, height int) {
	gl.Viewport(0, 0, width, height)
	p.logError()
}

func (p *painter) freeTexture(obj fyne.CanvasObject) {
	texture, ok := cache.GetTexture(obj)
	if !ok {
		return
	}

	gl.DeleteTexture(gl.Texture(texture))
	p.logError()
	cache.DeleteTexture(obj)
}

func (p *painter) glInit() {
	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	p.logError()
}

func (p *painter) compileShader(source string, shaderType gl.Enum) (gl.Shader, error) {
	shader := gl.CreateShader(shaderType)

	gl.ShaderSource(shader, source)
	p.logError()
	gl.CompileShader(shader)
	p.logError()

	status := gl.GetShaderi(shader, gl.COMPILE_STATUS)
	if status == gl.FALSE {
		info := gl.GetShaderInfoLog(shader)

		return gl.NoShader, fmt.Errorf("failed to compile %v: %v", source, info)
	}

	return shader, nil
}

var vertexShaderSource = string(shaderSimpleesVert.StaticContent)
var fragmentShaderSource = string(shaderSimpleesFrag.StaticContent)
var vertexLineShaderSource = string(shaderLineesVert.StaticContent)
var fragmentLineShaderSource = string(shaderLineesFrag.StaticContent)

func (p *painter) Init() {
	p.ctx = &xjsContext{}
	vertexShader, err := p.compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := p.compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	p.logError()

	p.program = Program(prog)

	vertexLineShader, err := p.compileShader(vertexLineShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentLineShader, err := p.compileShader(fragmentLineShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	lineProg := gl.CreateProgram()
	gl.AttachShader(lineProg, vertexLineShader)
	gl.AttachShader(lineProg, fragmentLineShader)
	gl.LinkProgram(lineProg)
	p.logError()

	p.lineProgram = Program(lineProg)
}

func (p *painter) glClearBuffer() {
	gl.UseProgram(gl.Program(p.program))
	p.logError()

	r, g, b, a := theme.BackgroundColor().RGBA()
	max16bit := float32(255 * 255)
	gl.ClearColor(float32(r)/max16bit, float32(g)/max16bit, float32(b)/max16bit, float32(a)/max16bit)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	p.logError()
}

func (p *painter) glScissorOpen(x, y, w, h int32) {
	gl.Scissor(x, y, w, h)
	gl.Enable(gl.SCISSOR_TEST)
	p.logError()
}

func (p *painter) glScissorClose() {
	gl.Disable(gl.SCISSOR_TEST)
	p.logError()
}

func (p *painter) glCreateBuffer(points []float32) Buffer {
	gl.UseProgram(gl.Program(p.program))

	vbo := gl.CreateBuffer()
	p.logError()
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	p.logError()
	gl.BufferData(gl.ARRAY_BUFFER, f32.Bytes(binary.LittleEndian, points...), gl.STATIC_DRAW)
	p.logError()

	vertAttrib := gl.GetAttribLocation(gl.Program(p.program), "vert")
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, 0)
	p.logError()

	texCoordAttrib := gl.GetAttribLocation(gl.Program(p.program), "vertTexCoord")
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, 3*4)
	p.logError()

	return Buffer(vbo)
}

func (p *painter) glCreateLineBuffer(points []float32) Buffer {
	gl.UseProgram(gl.Program(p.lineProgram))

	vbo := gl.CreateBuffer()
	p.logError()
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	p.logError()
	gl.BufferData(gl.ARRAY_BUFFER, f32.Bytes(binary.LittleEndian, points...), gl.STATIC_DRAW)
	p.logError()

	vertAttrib := gl.GetAttribLocation(gl.Program(p.lineProgram), "vert")
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 2, gl.FLOAT, false, 4*4, 0)
	p.logError()

	normalAttrib := gl.GetAttribLocation(gl.Program(p.lineProgram), "normal")
	gl.EnableVertexAttribArray(normalAttrib)
	gl.VertexAttribPointer(normalAttrib, 2, gl.FLOAT, false, 4*4, 2*4)
	p.logError()

	return Buffer(vbo)
}

func (p *painter) glFreeBuffer(vbo Buffer) {
	gl.BindBuffer(gl.ARRAY_BUFFER, gl.NoBuffer)
	p.logError()
	gl.DeleteBuffer(gl.Buffer(vbo))
	p.logError()
}

func (p *painter) glDrawTexture(texture Texture, alpha float32) {
	gl.UseProgram(gl.Program(p.program))

	// here we have to choose between blending the image alpha or fading it...
	// TODO find a way to support both
	if alpha != 1.0 {
		gl.BlendColor(0, 0, 0, alpha)
		gl.BlendFunc(gl.CONSTANT_ALPHA, gl.ONE_MINUS_CONSTANT_ALPHA)
	} else {
		gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)
	}
	p.logError()

	p.ctx.ActiveTexture(texture0)
	p.ctx.BindTexture(texture2D, texture)
	p.logError()

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	p.logError()
}

func (p *painter) glDrawLine(width float32, col color.Color, feather float32) {
	gl.UseProgram(gl.Program(p.lineProgram))

	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	p.logError()

	colorUniform := gl.GetUniformLocation(gl.Program(p.lineProgram), "color")
	r, g, b, a := col.RGBA()
	if a == 0 {
		gl.Uniform4f(colorUniform, 0, 0, 0, 0)
	} else {
		alpha := float32(a)
		col := []float32{float32(r) / alpha, float32(g) / alpha, float32(b) / alpha, alpha / 0xffff}
		gl.Uniform4fv(colorUniform, col)
	}
	lineWidthUniform := gl.GetUniformLocation(gl.Program(p.lineProgram), "lineWidth")
	gl.Uniform1f(lineWidthUniform, width)

	featherUniform := gl.GetUniformLocation(gl.Program(p.lineProgram), "feather")
	gl.Uniform1f(featherUniform, feather)

	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	p.logError()
}

func (p *painter) glCapture(width, height int32, pixels *[]uint8) {
	gl.ReadPixels(*pixels, 0, 0, int(width), int(height), gl.RGBA, gl.UNSIGNED_BYTE)
	p.logError()
}

type xjsContext struct{}

var _ context = (*xjsContext)(nil)

func (c *xjsContext) ActiveTexture(textureUnit uint32) {
	gl.ActiveTexture(gl.Enum(textureUnit))
}

func (c *xjsContext) BindTexture(target uint32, texture Texture) {
	gl.BindTexture(gl.Enum(target), gl.Texture(texture))
}

func (c *xjsContext) CreateTexture() (texture Texture) {
	return Texture(gl.CreateTexture())
}

func (c *xjsContext) GetError() uint32 {
	return uint32(gl.GetError())
}

func (c *xjsContext) TexImage2D(target uint32, level, width, height int, colorFormat, typ uint32, data []uint8) {
	gl.TexImage2D(
		gl.Enum(target),
		level,
		width,
		height,
		gl.Enum(colorFormat),
		gl.Enum(typ),
		data,
	)
}

func (c *xjsContext) TexParameteri(target, param uint32, value int32) {
	gl.TexParameteri(gl.Enum(target), gl.Enum(param), int(value))
}
