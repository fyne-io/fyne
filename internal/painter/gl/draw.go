package gl

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/cache"
	"fyne.io/fyne/internal/painter"
)

func rectInnerCoords(size fyne.Size, pos fyne.Position, fill canvas.ImageFill, aspect float32) (fyne.Size, fyne.Position) {
	if fill == canvas.ImageFillContain || fill == canvas.ImageFillOriginal {
		// change pos and size accordingly

		viewAspect := float32(size.Width) / float32(size.Height)

		newWidth, newHeight := size.Width, size.Height
		widthPad, heightPad := 0, 0
		if viewAspect > aspect {
			newWidth = int(float32(size.Height) * aspect)
			widthPad = (size.Width - newWidth) / 2
		} else if viewAspect < aspect {
			newHeight = int(float32(size.Width) / aspect)
			heightPad = (size.Height - newHeight) / 2
		}

		return fyne.NewSize(newWidth, newHeight), fyne.NewPos(pos.X+widthPad, pos.Y+heightPad)
	}

	return size, pos
}

// rectCoords calculates the openGL coordinate space of a rectangle
func (p *glPainter) rectCoords(size fyne.Size, pos fyne.Position, frame fyne.Size,
	fill canvas.ImageFill, aspect float32, pad int) []float32 {
	size, pos = rectInnerCoords(size, pos, fill, aspect)

	xPos := float32(pos.X-pad) / float32(frame.Width)
	x1 := -1 + xPos*2
	x2Pos := float32(pos.X+size.Width+pad) / float32(frame.Width)
	x2 := -1 + x2Pos*2

	yPos := float32(pos.Y-pad) / float32(frame.Height)
	y1 := 1 - yPos*2
	y2Pos := float32(pos.Y+size.Height+pad) / float32(frame.Height)
	y2 := 1 - y2Pos*2

	return []float32{
		// coord x, y, z texture x, y
		x1, y2, 0, 0.0, 1.0, // top left
		x1, y1, 0, 0.0, 0.0, // bottom left
		x2, y2, 0, 1.0, 1.0, // top right
		x2, y1, 0, 1.0, 0.0, // bottom right
	}
}

func (p *glPainter) drawTextureWithDetails(o fyne.CanvasObject, creator func(canvasObject fyne.CanvasObject) Texture,
	pos fyne.Position, size, frame fyne.Size, fill canvas.ImageFill, alpha float32, pad int) {

	texture := getTexture(o, creator)
	if texture == NoTexture {
		return
	}

	aspect := float32(0)
	if img, ok := o.(*canvas.Image); ok {
		aspect = painter.GetAspect(img)
		if aspect == 0 {
			aspect = 1 // fallback, should not occur - normally an image load error
		}
	}
	points := p.rectCoords(size, pos, frame, fill, aspect, pad)
	vbo := p.glCreateBuffer(points)

	p.glDrawTexture(texture, alpha)
	p.glFreeBuffer(vbo)
}

func (p *glPainter) drawWidget(wid fyne.Widget, pos fyne.Position, frame fyne.Size) {
	if cache.Renderer(wid).BackgroundColor() == color.Transparent {
		return
	}

	p.drawTextureWithDetails(wid, p.newGlRectTexture, pos, wid.Size(), frame, canvas.ImageFillStretch, 1.0, 0)
}

func (p *glPainter) drawCircle(circle *canvas.Circle, pos fyne.Position, frame fyne.Size) {
	p.drawTextureWithDetails(circle, p.newGlCircleTexture, pos, circle.Size(), frame, canvas.ImageFillStretch,
		1.0, painter.VectorPad(circle))
}

func (p *glPainter) drawLine(line *canvas.Line, pos fyne.Position, frame fyne.Size) {
	p.drawTextureWithDetails(line, p.newGlLineTexture, pos, line.Size(), frame, canvas.ImageFillStretch,
		1.0, painter.VectorPad(line))
}

func (p *glPainter) drawImage(img *canvas.Image, pos fyne.Position, frame fyne.Size) {
	p.drawTextureWithDetails(img, p.newGlImageTexture, pos, img.Size(), frame, img.FillMode, float32(img.Alpha()), 0)
}

func (p *glPainter) drawRaster(img *canvas.Raster, pos fyne.Position, frame fyne.Size) {
	p.drawTextureWithDetails(img, p.newGlRasterTexture, pos, img.Size(), frame, canvas.ImageFillStretch, float32(img.Alpha()), 0)
}

func (p *glPainter) drawGradient(o fyne.CanvasObject, texCreator func(fyne.CanvasObject) Texture, pos fyne.Position, frame fyne.Size) {
	p.drawTextureWithDetails(o, texCreator, pos, o.Size(), frame, canvas.ImageFillStretch, 1.0, 0)
}

func (p *glPainter) drawRectangle(rect *canvas.Rectangle, pos fyne.Position, frame fyne.Size) {
	p.drawTextureWithDetails(rect, p.newGlRectTexture, pos, rect.Size(), frame, canvas.ImageFillStretch,
		1.0, painter.VectorPad(rect))
}

func (p *glPainter) drawText(text *canvas.Text, pos fyne.Position, frame fyne.Size) {
	if text.Text == "" {
		return
	}

	size := text.MinSize()
	containerSize := text.Size()
	switch text.Alignment {
	case fyne.TextAlignTrailing:
		pos = fyne.NewPos(pos.X+containerSize.Width-size.Width, pos.Y)
	case fyne.TextAlignCenter:
		pos = fyne.NewPos(pos.X+(containerSize.Width-size.Width)/2, pos.Y)
	}

	if text.Size().Height > text.MinSize().Height {
		pos = fyne.NewPos(pos.X, pos.Y+(text.Size().Height-text.MinSize().Height)/2)
	}

	p.drawTextureWithDetails(text, p.newGlTextTexture, pos, size, frame, canvas.ImageFillStretch, 1.0, 0)
}

func (p *glPainter) drawObject(o fyne.CanvasObject, pos fyne.Position, frame fyne.Size) {
	if !o.Visible() {
		return
	}
	switch obj := o.(type) {
	case *canvas.Circle:
		p.drawCircle(obj, pos, frame)
	case *canvas.Line:
		p.drawLine(obj, pos, frame)
	case *canvas.Image:
		p.drawImage(obj, pos, frame)
	case *canvas.Raster:
		p.drawRaster(obj, pos, frame)
	case *canvas.Rectangle:
		p.drawRectangle(obj, pos, frame)
	case *canvas.Text:
		p.drawText(obj, pos, frame)
	case *canvas.LinearGradient:
		p.drawGradient(obj, p.newGlLinearGradientTexture, pos, frame)
	case *canvas.RadialGradient:
		p.drawGradient(obj, p.newGlRadialGradientTexture, pos, frame)
	case fyne.Widget:
		p.drawWidget(obj, pos, frame)
	}
}
