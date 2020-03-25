package software

import (
	"image"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/goki/freetype"
	"github.com/goki/freetype/truetype"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
)

func drawImage(c fyne.Canvas, img *canvas.Image, pos fyne.Position, frame fyne.Size, base *image.RGBA) {
	bounds := img.Size()
	width := internal.ScaleInt(c, bounds.Width)
	height := internal.ScaleInt(c, bounds.Height)

	if img.FillMode == canvas.ImageFillStretch {
		width = internal.ScaleInt(c, img.Size().Width)
		height = internal.ScaleInt(c, img.Size().Height)
	} // TODO support contain

	scaledX, scaledY := internal.ScaleInt(c, pos.X), internal.ScaleInt(c, pos.Y)
	outBounds := image.Rect(scaledX, scaledY, scaledX+width, scaledY+height)
	switch img.TextureFilter {
	case canvas.LinearFilter:
		draw.ApproxBiLinear.Scale(base, outBounds, img.Image, img.Image.Bounds(), draw.Over, nil)
	case canvas.NearestFilter:
		draw.NearestNeighbor.Scale(base, outBounds, img.Image, img.Image.Bounds(), draw.Over, nil)
	}
}

func drawText(c fyne.Canvas, text *canvas.Text, pos fyne.Position, frame fyne.Size, base *image.RGBA) {
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

func drawRectangle(c fyne.Canvas, rect *canvas.Rectangle, pos fyne.Position, frame fyne.Size, base *image.RGBA) {
	scaledWidth := internal.ScaleInt(c, rect.Size().Width)
	scaledHeight := internal.ScaleInt(c, rect.Size().Height)
	scaledX, scaledY := internal.ScaleInt(c, pos.X), internal.ScaleInt(c, pos.Y)
	bounds := image.Rect(scaledX, scaledY, scaledX+scaledWidth, scaledY+scaledHeight)
	draw.Draw(base, bounds, image.NewUniform(rect.FillColor), image.ZP, draw.Over)
}

func drawWidget(c fyne.Canvas, wid fyne.Widget, pos fyne.Position, frame fyne.Size, base *image.RGBA) {
	scaledWidth := internal.ScaleInt(c, wid.Size().Width)
	scaledHeight := internal.ScaleInt(c, wid.Size().Height)
	scaledX, scaledY := internal.ScaleInt(c, pos.X), internal.ScaleInt(c, pos.Y)
	bounds := image.Rect(scaledX, scaledY, scaledX+scaledWidth, scaledY+scaledHeight)
	draw.Draw(base, bounds, image.NewUniform(widget.Renderer(wid).BackgroundColor()), image.ZP, draw.Over)
}
