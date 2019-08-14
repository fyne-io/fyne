package software

import (
	"image"
	"image/draw"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/painter"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/goki/freetype"
	"github.com/goki/freetype/truetype"
	"golang.org/x/image/font"
)

func drawImage(c fyne.Canvas, img *canvas.Image, pos fyne.Position, frame fyne.Size, base *image.RGBA) {
	bounds := img.Size()
	width := bounds.Width   //textureScaleInt(c, bounds.Width)
	height := bounds.Height //textureScaleInt(c, bounds.Height)
	tex := painter.PaintImage(img, c, width, height)

	outBounds := image.Rect(pos.X, pos.Y, bounds.Width+pos.X, bounds.Height+pos.Y)
	draw.Draw(base, outBounds, tex, image.ZP, draw.Over)
}

func drawText(c fyne.Canvas, text *canvas.Text, pos fyne.Position, frame fyne.Size, base *image.RGBA) {
	bounds := text.MinSize()
	width := bounds.Width   //textureScaleInt(c, bounds.Width)
	height := bounds.Height //textureScaleInt(c, bounds.Height)
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

	imgBounds := image.Rect(pos.X, pos.Y, text.Size().Width+pos.X, text.Size().Height+pos.Y)
	draw.Draw(base, imgBounds, txtImg, image.ZP, draw.Over)
}

func drawRectangle(c fyne.Canvas, rect *canvas.Rectangle, pos fyne.Position, frame fyne.Size, base *image.RGBA) {
	bounds := image.Rect(pos.X, pos.Y, rect.Size().Width+pos.X, rect.Size().Height+pos.Y)
	draw.Draw(base, bounds, image.NewUniform(rect.FillColor), image.ZP, draw.Over)
}

func drawWidget(c fyne.Canvas, wid fyne.Widget, pos fyne.Position, frame fyne.Size, base *image.RGBA) {
	bounds := image.Rect(pos.X, pos.Y, wid.Size().Width+pos.X, wid.Size().Height+pos.Y)
	draw.Draw(base, bounds, image.NewUniform(widget.Renderer(wid).BackgroundColor()), image.ZP, draw.Over)
}
