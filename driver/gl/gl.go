package gl

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg" // avoid users having to import when using image widget
	_ "image/png"  // avoid the same for PNG images
	"io"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/goki/freetype"
	"github.com/goki/freetype/truetype"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var textures = make(map[fyne.CanvasObject]uint32)

const vectorPad = 10

func getTexture(object fyne.CanvasObject, creator func(canvasObject fyne.CanvasObject) uint32) uint32 {
	if _, skipCache := object.(*canvas.Raster); skipCache {
		return creator(object)
	}
	texture := textures[object]

	if texture == 0 {
		texture = creator(object)
		textures[object] = texture
	}
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

func (c *glCanvas) newGlCircleTexture(obj fyne.CanvasObject) uint32 {
	circle := obj.(*canvas.Circle)
	radius := fyne.Min(circle.Size().Width, circle.Size().Height) / 2

	width := textureScaleInt(c, circle.Size().Width+vectorPad*2)
	height := textureScaleInt(c, circle.Size().Height+vectorPad*2)
	stroke := circle.StrokeWidth * c.scale * c.texScale

	raw := image.NewRGBA(image.Rect(0, 0, width, height))
	scanner := rasterx.NewScannerGV(circle.Size().Width, circle.Size().Height, raw, raw.Bounds())

	if circle.FillColor != nil {
		filler := rasterx.NewFiller(width, height, scanner)
		filler.SetColor(circle.FillColor)
		rasterx.AddCircle(float64(width/2), float64(height/2), float64(textureScaleInt(c, radius)), filler)
		filler.Draw()
	}

	dasher := rasterx.NewDasher(width, height, scanner)
	dasher.SetColor(circle.StrokeColor)
	dasher.SetStroke(fixed.Int26_6(float64(stroke)*64), 0, nil, nil, nil, 0, nil, 0)
	rasterx.AddCircle(float64(width/2), float64(height/2), float64(textureScaleInt(c, radius)), dasher)
	dasher.Draw()

	return c.imgToTexture(raw)
}

func (c *glCanvas) newGlLineTexture(obj fyne.CanvasObject) uint32 {
	line := obj.(*canvas.Line)

	col := line.StrokeColor
	width := textureScaleInt(c, line.Size().Width+vectorPad*2)
	height := textureScaleInt(c, line.Size().Height+vectorPad*2)
	stroke := line.StrokeWidth * c.scale * c.texScale

	raw := image.NewRGBA(image.Rect(0, 0, width, height))
	scanner := rasterx.NewScannerGV(line.Size().Width, line.Size().Height, raw, raw.Bounds())
	dasher := rasterx.NewDasher(width, height, scanner)
	dasher.SetColor(col)
	dasher.SetStroke(fixed.Int26_6(float64(stroke)*64), 0, nil, nil, nil, 0, nil, 0)
	p1x, p1y := textureScaleInt(c, line.Position1.X-line.Position().X+vectorPad), textureScaleInt(c, line.Position1.Y-line.Position().Y+vectorPad)
	p2x, p2y := textureScaleInt(c, line.Position2.X-line.Position().X+vectorPad), textureScaleInt(c, line.Position2.Y-line.Position().Y+vectorPad)

	dasher.Start(rasterx.ToFixedP(float64(p1x), float64(p1y)))
	dasher.Line(rasterx.ToFixedP(float64(p2x), float64(p2y)))
	dasher.Stop(true)
	dasher.Draw()

	return c.imgToTexture(raw)
}

func (c *glCanvas) newGlRectTexture(rect fyne.CanvasObject) uint32 {
	col := theme.BackgroundColor()
	if wid, ok := rect.(fyne.Widget); ok {
		widCol := widget.Renderer(wid).BackgroundColor()
		if widCol != nil {
			col = widCol
		}
	} else if rect, ok := rect.(*canvas.Rectangle); ok {
		if rect.FillColor != nil {
			col = rect.FillColor
		}
	}

	return c.imgToTexture(image.NewUniform(col))
}

func (c *glCanvas) newGlTextTexture(obj fyne.CanvasObject) uint32 {
	text := obj.(*canvas.Text)

	bounds := text.MinSize()
	width := textureScaleInt(c, bounds.Width)
	height := textureScaleInt(c, bounds.Height)
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	var opts truetype.Options
	fontSize := float64(text.TextSize) * float64(c.Scale())
	opts.Size = fontSize
	opts.DPI = float64(textDPI * c.texScale)
	face := cachedFontFace(text.TextStyle, &opts)

	d := font.Drawer{}
	d.Dst = img
	d.Src = &image.Uniform{C: text.Color}
	d.Face = face
	d.Dot = freetype.Pt(0, height-face.Metrics().Descent.Ceil())
	d.DrawString(text.Text)

	return c.imgToTexture(img)
}

func (c *glCanvas) newGlImageTexture(obj fyne.CanvasObject) uint32 {
	img := obj.(*canvas.Image)

	width := textureScaleInt(c, img.Size().Width)
	height := textureScaleInt(c, img.Size().Height)
	if width <= 0 || height <= 0 {
		return 0
	}

	switch {
	case img.File != "" || img.Resource != nil:
		var file io.Reader
		var name string
		if img.Resource != nil {
			name = img.Resource.Name()
			file = bytes.NewReader(img.Resource.Content())
		} else {
			name = img.File
			handle, _ := os.Open(img.File)
			defer handle.Close()
			file = handle
		}

		if strings.ToLower(filepath.Ext(name)) == ".svg" {
			tex := svgCacheGet(img.Resource, width, height)
			if tex == nil {
				// Not in cache, so load the item and add to cache

				icon, err := oksvg.ReadIconStream(file)
				if err != nil {
					fyne.LogError("SVG Load error:", err)

					return 0
				}
				icon.SetTarget(0, 0, float64(width), float64(height))

				w, h := int(icon.ViewBox.W), int(icon.ViewBox.H)
				// this is used by our render code, so let's set it to the file aspect
				aspects[img.Resource] = float32(w) / float32(h)
				// if the image specifies it should be original size we need at least that many pixels on screen
				if img.FillMode == canvas.ImageFillOriginal {
					pixSize := fyne.NewSize(unscaleInt(c, w), unscaleInt(c, h))
					img.SetMinSize(pixSize)
				}

				tex = image.NewRGBA(image.Rect(0, 0, width, height))
				scanner := rasterx.NewScannerGV(w, h, tex, tex.Bounds())
				raster := rasterx.NewDasher(width, height, scanner)

				icon.Draw(raster, 1)
				svgCachePut(img.Resource, tex, width, height)
			}

			return c.imgToTexture(tex)
		}

		pixels, _, err := image.Decode(file)

		if err != nil {
			fyne.LogError("image err", err)

			return 0
		}
		origSize := pixels.Bounds().Size()
		// this is used by our render code, so let's set it to the file aspect
		aspects[img] = float32(origSize.X) / float32(origSize.Y)
		// if the image specifies it should be original size we need at least that many pixels on screen
		if img.FillMode == canvas.ImageFillOriginal {
			pixSize := fyne.NewSize(unscaleInt(c, origSize.X), unscaleInt(c, origSize.Y))
			img.SetMinSize(pixSize)
		}

		tex := image.NewRGBA(pixels.Bounds())
		draw.Draw(tex, pixels.Bounds(), pixels, image.ZP, draw.Src)

		return c.imgToTexture(tex)
	case img.Image != nil:
		origSize := img.Image.Bounds().Size()
		// this is used by our render code, so let's set it to the file aspect
		aspects[img] = float32(origSize.X) / float32(origSize.Y)
		// if the image specifies it should be original size we need at least that many pixels on screen
		if img.FillMode == canvas.ImageFillOriginal {
			pixSize := fyne.NewSize(unscaleInt(c, origSize.X), unscaleInt(c, origSize.Y))
			img.SetMinSize(pixSize)
		}

		tex := image.NewRGBA(img.Image.Bounds())
		draw.Draw(tex, img.Image.Bounds(), img.Image, image.ZP, draw.Src)

		return c.imgToTexture(tex)
	default:
		return c.imgToTexture(image.NewRGBA(image.Rect(0, 0, 1, 1)))
	}
}

func (c *glCanvas) newGlRasterTexture(obj fyne.CanvasObject) uint32 {
	rast := obj.(*canvas.Raster)

	width := textureScaleInt(c, rast.Size().Width)
	height := textureScaleInt(c, rast.Size().Height)

	return c.imgToTexture(rast.Generator(width, height))
}

func (c *glCanvas) imgToTexture(img image.Image) uint32 {
	switch i := img.(type) {
	case *image.Uniform:
		texture := newTexture()
		r, g, b, a := i.RGBA()
		r8, g8, b8, a8 := uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)
		data := []uint8{r8, g8, b8, a8}
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, 1, 1, 0, gl.RGBA,
			gl.UNSIGNED_BYTE, gl.Ptr(data))
		return texture
	case *image.RGBA:
		if len(i.Pix) == 0 { // image is empty
			return 0
		}

		var texture uint32
		texture = newTexture()
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(i.Rect.Size().X), int32(i.Rect.Size().Y),
			0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(i.Pix))
		return texture
	default:
		rgba := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
		draw.Draw(rgba, rgba.Rect, img, image.ZP, draw.Over)
		return c.imgToTexture(rgba)
	}
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
    #version 110
    attribute vec3 vert;
    attribute vec2 vertTexCoord;
    varying vec2 fragTexCoord;

    void main() {
        fragTexCoord = vertTexCoord;

        gl_Position = vec4(vert, 1);
    }
` + "\x00"

	fragmentShaderSource = `
    #version 110
    uniform sampler2D tex;

    varying vec2 fragTexCoord;

    void main() {
        gl_FragColor = texture2D(tex, fragTexCoord);
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
