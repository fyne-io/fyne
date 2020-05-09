package software

import (
	"fmt"
	"image"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal"
	"fyne.io/fyne/internal/painter"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/goki/freetype"
	"github.com/goki/freetype/truetype"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
)

type gradient interface {
	Generate(int, int) image.Image
	Size() fyne.Size
}

func drawGradient(c fyne.Canvas, g gradient, pos fyne.Position, frame fyne.Size, base *image.NRGBA) {
	bounds := g.Size()
	width := internal.ScaleInt(c, bounds.Width)
	height := internal.ScaleInt(c, bounds.Height)
	tex := g.Generate(width, height)
	drawTex(c, pos, width, height, base, tex)
}

func drawImage(c fyne.Canvas, img *canvas.Image, pos fyne.Position, frame fyne.Size, base *image.NRGBA) {
	bounds := img.Size()
	width := internal.ScaleInt(c, bounds.Width)
	height := internal.ScaleInt(c, bounds.Height)
	scaledX, scaledY := internal.ScaleInt(c, pos.X), internal.ScaleInt(c, pos.Y)

	tmpImg := painter.PaintImage(img, c, width, height)

	if img.FillMode == canvas.ImageFillContain {
		imgAspect := painter.GetAspect(img)
		objAspect := float32(width) / float32(height)

		fmt.Printf("imgAspect:%f - objAspect:%f\n", imgAspect, objAspect)

		if objAspect < 1 {
			newHeight := int(float32(width) / imgAspect)
			scaledY += (height - newHeight) / 2
			height = internal.ScaleInt(c, newHeight)
		} else if objAspect > 1 {
			newWidth := int(float32(height) * imgAspect)
			scaledX += (width - newWidth) / 2
			width = internal.ScaleInt(c, newWidth)
		}
	}

	outBounds := image.Rect(scaledX, scaledY, scaledX+width, scaledY+height)
	switch img.ScalingFilter {
	case canvas.LinearFilter:
		draw.CatmullRom.Scale(base, outBounds, tmpImg, img.Image.Bounds(), draw.Over, nil)
	case canvas.NearestFilter:
		draw.NearestNeighbor.Scale(base, outBounds, tmpImg, img.Image.Bounds(), draw.Over, nil)
	}
}

func drawTex(c fyne.Canvas, pos fyne.Position, width int, height int, base *image.NRGBA, tex image.Image) {
	scaledX, scaledY := internal.ScaleInt(c, pos.X), internal.ScaleInt(c, pos.Y)
	outBounds := image.Rect(scaledX, scaledY, scaledX+width, scaledY+height)
	draw.Draw(base, outBounds, tex, image.ZP, draw.Over)
}

func drawText(c fyne.Canvas, text *canvas.Text, pos fyne.Position, frame fyne.Size, base *image.NRGBA) {
	bounds := text.MinSize()
	width := internal.ScaleInt(c, bounds.Width)
	height := internal.ScaleInt(c, bounds.Height)
	txtImg := image.NewRGBA(image.Rect(0, 0, width, height))

	var opts truetype.Options
	fontSize := float64(text.TextSize) * float64(c.Scale())
	opts.Size = fontSize
	opts.DPI = 78.0
	f, _ := truetype.Parse(theme.TextFont().Content())
	face := truetype.NewFace(f, &opts)

	d := font.Drawer{}
	d.Dst = txtImg
	d.Src = &image.Uniform{C: text.Color}
	d.Face = face
	d.Dot = freetype.Pt(0, height-face.Metrics().Descent.Ceil())
	d.DrawString(text.Text)

	scaledX, scaledY := internal.ScaleInt(c, pos.X), internal.ScaleInt(c, pos.Y)
	imgBounds := image.Rect(scaledX, scaledY, scaledX+width, scaledY+height)
	draw.Draw(base, imgBounds, txtImg, image.ZP, draw.Over)
}

func drawRectangle(c fyne.Canvas, rect *canvas.Rectangle, pos fyne.Position, frame fyne.Size, base *image.NRGBA) {
	scaledWidth := internal.ScaleInt(c, rect.Size().Width)
	scaledHeight := internal.ScaleInt(c, rect.Size().Height)
	scaledX, scaledY := internal.ScaleInt(c, pos.X), internal.ScaleInt(c, pos.Y)
	bounds := image.Rect(scaledX, scaledY, scaledX+scaledWidth, scaledY+scaledHeight)
	draw.Draw(base, bounds, image.NewUniform(rect.FillColor), image.ZP, draw.Over)
}

func drawWidget(c fyne.Canvas, wid fyne.Widget, pos fyne.Position, frame fyne.Size, base *image.NRGBA) {
	scaledWidth := internal.ScaleInt(c, wid.Size().Width)
	scaledHeight := internal.ScaleInt(c, wid.Size().Height)
	scaledX, scaledY := internal.ScaleInt(c, pos.X), internal.ScaleInt(c, pos.Y)
	bounds := image.Rect(scaledX, scaledY, scaledX+scaledWidth, scaledY+scaledHeight)
	draw.Draw(base, bounds, image.NewUniform(widget.Renderer(wid).BackgroundColor()), image.ZP, draw.Over)
}
