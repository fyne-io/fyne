//go:build js || wasm || web
// +build js wasm web

package gl

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"encoding/binary"

	gl "github.com/fyne-io/gl-js"
	"golang.org/x/mobile/exp/f32"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/theme"
)

// Buffer represents a GL buffer
type Buffer gl.Buffer

// Program represents a compiled GL program
type Program gl.Program

var textureFilterToGL = []int{gl.LINEAR, gl.NEAREST}

func (p *glPainter) newTexture(textureFilter canvas.ImageScale) Texture {
	var texture = gl.CreateTexture()
	logError()

	if int(textureFilter) >= len(textureFilterToGL) {
		fyne.LogError(fmt.Sprintf("Invalid canvas.ImageScale value (%d), using canvas.ImageScaleSmooth as default value", textureFilter), nil)
		textureFilter = canvas.ImageScaleSmooth
	}

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	logError()
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, textureFilterToGL[textureFilter])
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, textureFilterToGL[textureFilter])
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	logError()

	return Texture(texture)
}

func (p *glPainter) getTexture(object fyne.CanvasObject, creator func(canvasObject fyne.CanvasObject) Texture) (Texture, error) {
	texture, ok := cache.GetTexture(object)

	if !ok {
		texture = cache.TextureType(creator(object))
		cache.SetTexture(object, texture, p.canvas)
	}
	if !gl.Texture(texture).Valid() {
		return NoTexture, fmt.Errorf("No texture available.")
	}
	return Texture(texture), nil
}

func (p *glPainter) imgToTexture(img image.Image, textureFilter canvas.ImageScale) Texture {
	switch i := img.(type) {
	case *image.Uniform:
		texture := p.newTexture(textureFilter)
		r, g, b, a := i.RGBA()
		r8, g8, b8, a8 := uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)
		data := []uint8{r8, g8, b8, a8}
		gl.TexImage2D(gl.TEXTURE_2D, 0, 1, 1, gl.RGBA, gl.UNSIGNED_BYTE, data)
		logError()
		return texture
	case *image.RGBA:
		if len(i.Pix) == 0 { // image is empty
			return NoTexture
		}

		texture := p.newTexture(textureFilter)
		gl.TexImage2D(gl.TEXTURE_2D, 0, i.Rect.Size().X, i.Rect.Size().Y,
			gl.RGBA, gl.UNSIGNED_BYTE, i.Pix)
		logError()
		return texture
	default:
		rgba := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
		draw.Draw(rgba, rgba.Rect, img, image.ZP, draw.Over)
		return p.imgToTexture(rgba, textureFilter)
	}
}

func (p *glPainter) SetOutputSize(width, height int) {
	gl.Viewport(0, 0, width, height)
	logError()
}

func (p *glPainter) freeTexture(obj fyne.CanvasObject) {
	texture, ok := cache.GetTexture(obj)
	if !ok {
		return
	}

	gl.DeleteTexture(gl.Texture(texture))
	logError()
	cache.DeleteTexture(obj)
}

func glInit() {
	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	logError()
}

func compileShader(source string, shaderType gl.Enum) (gl.Shader, error) {
	shader := gl.CreateShader(shaderType)

	gl.ShaderSource(shader, source)
	logError()
	gl.CompileShader(shader)
	logError()

	status := gl.GetShaderi(shader, gl.COMPILE_STATUS)
	if status == gl.FALSE {
		info := gl.GetShaderInfoLog(shader)

		return gl.NoShader, fmt.Errorf("failed to compile %v: %v", source, info)
	}

	return shader, nil
}

const (
	vertexShaderSource = `
    #version 100
    attribute vec3 vert;
    attribute vec2 vertTexCoord;
    varying highp vec2 fragTexCoord;

    void main() {
        fragTexCoord = vertTexCoord;

        gl_Position = vec4(vert, 1);
    }`

	fragmentShaderSource = `
    #version 100
    uniform sampler2D tex;

    varying highp vec2 fragTexCoord;

    void main() {
        gl_FragColor = texture2D(tex, fragTexCoord);
    }`

	vertexLineShaderSource = `
    #version 100
    attribute vec2 vert;
    attribute vec2 normal;
    
    uniform float lineWidth;

    varying highp vec2 delta;

    void main() {
        delta = normal * lineWidth;

        gl_Position = vec4(vert + delta, 0, 1);
    }`

	fragmentLineShaderSource = `
    #version 100
    uniform highp vec4 color;
    uniform highp float lineWidth;
    uniform highp float feather;

    varying highp vec2 delta;

    void main() {
        highp float alpha = color.a;
        highp float distance = length(delta);

        if (feather == 0.0 || distance <= lineWidth - feather) {
           gl_FragColor = color;
        } else {
           gl_FragColor = vec4(color.r, color.g, color.b, mix(color.a, 0.0, (distance - (lineWidth - feather)) / feather));
        }
    }`
)

func (p *glPainter) Init() {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	logError()

	p.program = Program(prog)

	vertexLineShader, err := compileShader(vertexLineShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentLineShader, err := compileShader(fragmentLineShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	lineProg := gl.CreateProgram()
	gl.AttachShader(lineProg, vertexLineShader)
	gl.AttachShader(lineProg, fragmentLineShader)
	gl.LinkProgram(lineProg)
	logError()

	p.lineProgram = Program(lineProg)
}

func (p *glPainter) glClearBuffer() {
	gl.UseProgram(gl.Program(p.program))
	logError()

	r, g, b, a := theme.BackgroundColor().RGBA()
	max16bit := float32(255 * 255)
	gl.ClearColor(float32(r)/max16bit, float32(g)/max16bit, float32(b)/max16bit, float32(a)/max16bit)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	logError()
}

func (p *glPainter) glScissorOpen(x, y, w, h int32) {
	gl.Scissor(x, y, w, h)
	gl.Enable(gl.SCISSOR_TEST)
	logError()
}

func (p *glPainter) glScissorClose() {
	gl.Disable(gl.SCISSOR_TEST)
	logError()
}

func (p *glPainter) glCreateBuffer(points []float32) Buffer {
	gl.UseProgram(gl.Program(p.program))

	vbo := gl.CreateBuffer()
	logError()
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	logError()
	gl.BufferData(gl.ARRAY_BUFFER, f32.Bytes(binary.LittleEndian, points...), gl.STATIC_DRAW)
	logError()

	vertAttrib := gl.GetAttribLocation(gl.Program(p.program), "vert")
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, 0)
	logError()

	texCoordAttrib := gl.GetAttribLocation(gl.Program(p.program), "vertTexCoord")
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, 3*4)
	logError()

	return Buffer(vbo)
}

func (p *glPainter) glCreateLineBuffer(points []float32) Buffer {
	gl.UseProgram(gl.Program(p.lineProgram))

	vbo := gl.CreateBuffer()
	logError()
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	logError()
	gl.BufferData(gl.ARRAY_BUFFER, f32.Bytes(binary.LittleEndian, points...), gl.STATIC_DRAW)
	logError()

	vertAttrib := gl.GetAttribLocation(gl.Program(p.lineProgram), "vert")
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 2, gl.FLOAT, false, 4*4, 0)
	logError()

	normalAttrib := gl.GetAttribLocation(gl.Program(p.lineProgram), "normal")
	gl.EnableVertexAttribArray(normalAttrib)
	gl.VertexAttribPointer(normalAttrib, 2, gl.FLOAT, false, 4*4, 2*4)
	logError()

	return Buffer(vbo)
}

func (p *glPainter) glFreeBuffer(vbo Buffer) {
	gl.BindBuffer(gl.ARRAY_BUFFER, gl.NoBuffer)
	logError()
	gl.DeleteBuffer(gl.Buffer(vbo))
	logError()
}

func (p *glPainter) glDrawTexture(texture Texture, alpha float32) {
	gl.UseProgram(gl.Program(p.program))

	// here we have to choose between blending the image alpha or fading it...
	// TODO find a way to support both
	if alpha != 1.0 {
		gl.BlendColor(0, 0, 0, alpha)
		gl.BlendFunc(gl.CONSTANT_ALPHA, gl.ONE_MINUS_CONSTANT_ALPHA)
	} else {
		gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)
	}
	logError()

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, gl.Texture(texture))
	logError()

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	logError()
}

func (p *glPainter) glDrawLine(width float32, col color.Color, feather float32) {
	gl.UseProgram(gl.Program(p.lineProgram))

	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	logError()

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
	logError()
}

func (p *glPainter) glCapture(width, height int32, pixels *[]uint8) {
	gl.ReadPixels(*pixels, 0, 0, int(width), int(height), gl.RGBA, gl.UNSIGNED_BYTE)
	logError()
}

func logError() {
	if fyne.CurrentApp().Settings().BuildType() != fyne.BuildDebug {
		return
	}
	logGLError(uint32(gl.GetError()))
}
