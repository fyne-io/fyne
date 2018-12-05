// +build !ci,gl

package gl

import (
	"fmt"
	"image"
	"image/draw"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/canvas"
	"github.com/fyne-io/fyne/theme"
	"github.com/fyne-io/fyne/widget"
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

var textures = make(map[fyne.CanvasObject]uint32)
var refreshQueue = make(chan fyne.CanvasObject, 1024)

func getTexture(object fyne.CanvasObject, creator func(canvasObject fyne.CanvasObject) uint32) uint32 {
	texture := textures[object]
	if texture != 0 {
		return texture
	}

	texture = creator(object)
	textures[object] = texture
	return texture
}

func newTexture() uint32 {
	var texture uint32

	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	return texture
}

func (c *glCanvas) newGlRectTexture(rect fyne.CanvasObject) uint32 {
	texture := newTexture()

	col := theme.BackgroundColor()
	if _, ok := rect.(*widget.Toolbar); ok { // TODO don't make this a special case
		col = theme.ButtonColor()
	} else if _, ok := rect.(*canvas.Rectangle); ok {
		col = rect.(*canvas.Rectangle).FillColor
	}

	r, g, b, a := col.RGBA()
	data := []uint8{uint8(r), uint8(g), uint8(b), uint8(a)}
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, 1, 1, 0, gl.RGBA,
		gl.UNSIGNED_BYTE, gl.Ptr(data))

	return texture
}

func (c *glCanvas) newGlTextTexture(text fyne.CanvasObject) uint32 {
	texture := newTexture()
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	bounds := text.MinSize()
	width := scaleInt(c, bounds.Width*4)
	height := scaleInt(c, bounds.Height*4)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	dpi := float64(textDPI) * 4

	var opts truetype.Options
	fc := fontCache(text.(*canvas.Text).TextStyle)
	fontSize := float64(text.(*canvas.Text).TextSize) * float64(c.Scale())
	opts.Size = fontSize
	opts.DPI = dpi
	face := truetype.NewFace(fc, &opts)

	ctx := freetype.NewContext()
	ctx.SetDPI(dpi)
	ctx.SetFont(fc)
	ctx.SetFontSize(fontSize)
	ctx.SetClip(img.Bounds())
	ctx.SetDst(img)
	ctx.SetSrc(&image.Uniform{theme.TextColor()})

	ctx.DrawString(text.(*canvas.Text).Text, freetype.Pt(0, height+2-face.Metrics().Descent.Ceil()))

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(img.Rect.Size().X), int32(img.Rect.Size().Y),
		0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))

	return texture
}

func (c *glCanvas) newGlImageTexture(img *canvas.Image) uint32 {
	texture := newTexture()
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	width := scaleInt(c, img.Size.Width)
	height := scaleInt(c, img.Size.Height)
	raw := image.NewRGBA(image.Rect(0, 0, width, height))

	if img.File != "" {
		if strings.ToLower(filepath.Ext(img.File)) == ".svg" {
			icon, err := oksvg.ReadIcon(img.File)
			icon.SetTarget(0, 0, float64(width), float64(height))

			w, h := int(icon.ViewBox.W), int(icon.ViewBox.H)
			if err != nil {
				log.Println("SVG Load error:", err, img.File)

				return 0
			}
			raw = image.NewRGBA(image.Rect(0, 0, width, height))
			scanner := rasterx.NewScannerGV(w, h, raw, raw.Bounds())
			raster := rasterx.NewDasher(width, height, scanner)

			icon.Draw(raster, img.Alpha())
		} else {
			file, _ := os.Open(img.File)
			pixels, _, err := image.Decode(file)

			if err != nil {
				log.Println("image err", err)

				return 0
			} else {
				raw = image.NewRGBA(pixels.Bounds())

				draw.Draw(raw, raw.Bounds(), pixels, image.ZP, draw.Src)
			}
		}
	} else if img.PixelColor != nil {
		pixels := NewPixelImage(img, c.Scale())
		draw.Draw(raw, raw.Bounds(), pixels, image.ZP, draw.Src)
	}

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(raw.Rect.Size().X), int32(raw.Rect.Size().Y),
		0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(raw.Pix))

	return texture
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		info := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(info))

		return 0, fmt.Errorf("failed to compile %v: %v", source, info)
	}

	return shader, nil
}

const (
	vertexShaderSource = `
    #version 150
    in vec3 vert;
    in vec2 vertTexCoord;
    out vec2 fragTexCoord;

    void main() {
        fragTexCoord = vertTexCoord;

        gl_Position = vec4(vert, 1);
    }
` + "\x00"

	fragmentShaderSource = `
    #version 150
    uniform sampler2D tex;

    in vec2 fragTexCoord;
    out vec4 frag_colour;
    
    void main() {
        vec4 color = texture(tex, fragTexCoord);
        if(color.a < 0.01)
            discard;

        frag_colour = color;
    }
` + "\x00"
)

func (c *glCanvas) initOpenGL() {
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

	c.program = prog
}
