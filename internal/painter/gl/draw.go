package gl

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
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
	fill canvas.ImageFill, aspect float32, pad int) ([]float32, uint32, uint32) {
	size, pos = rectInnerCoords(size, pos, fill, aspect)

	xPos := float32(pos.X-pad) / float32(frame.Width)
	x1 := -1 + xPos*2
	x2Pos := float32(pos.X+size.Width+pad) / float32(frame.Width)
	x2 := -1 + x2Pos*2

	yPos := float32(pos.Y-pad) / float32(frame.Height)
	y1 := 1 - yPos*2
	y2Pos := float32(pos.Y+size.Height+pad) / float32(frame.Height)
	y2 := 1 - y2Pos*2

	points := []float32{
		// coord x, y, x texture x, y
		x1, y2, 0, 0.0, 1.0, // top left
		x1, y1, 0, 0.0, 0.0, // bottom left
		x2, y2, 0, 1.0, 1.0, // top right
		x2, y1, 0, 1.0, 0.0, // bottom right
	}

	vao, vbo := p.glCreateBuffer(points)
	return points, vao, vbo
}

func (p *glPainter) freeCoords(vao, vbo uint32) {
	p.glFreeBuffer(vao, vbo)
}

func (p *glPainter) drawWidget(wid fyne.Widget, pos fyne.Position, frame fyne.Size) {
	if widget.Renderer(wid).BackgroundColor() == color.Transparent {
		return
	}

	points, vao, vbo := p.rectCoords(wid.Size(), pos, frame, canvas.ImageFillStretch, 0.0, 0)
	texture := getTexture(wid, p.newGlRectTexture)

	p.glDrawTexture(texture, points, 1.0)
	p.freeCoords(vao, vbo)
}

func (p *glPainter) drawCircle(circle *canvas.Circle, pos fyne.Position, frame fyne.Size) {
	points, vao, vbo := p.rectCoords(circle.Size(), pos, frame, canvas.ImageFillStretch, 0.0, vectorPad)
	texture := getTexture(circle, p.newGlCircleTexture)

	p.glDrawTexture(texture, points, 1.0)
	p.freeCoords(vao, vbo)
}

func (p *glPainter) drawLine(line *canvas.Line, pos fyne.Position, frame fyne.Size) {
	points, vao, vbo := p.rectCoords(line.Size(), pos, frame, canvas.ImageFillStretch, 0.0, vectorPad)
	texture := getTexture(line, p.newGlLineTexture)

	p.glDrawTexture(texture, points, 1.0)
	p.freeCoords(vao, vbo)
}

func (p *glPainter) drawImage(img *canvas.Image, pos fyne.Position, frame fyne.Size) {
	texture := getTexture(img, p.newGlImageTexture)
	if texture == 0 {
		return
	}

	aspect := aspects[img.Resource]
	if aspect == 0 {
		aspect = aspects[img]
	}
	points, vao, vbo := p.rectCoords(img.Size(), pos, frame, img.FillMode, aspect, 0)
	p.glDrawTexture(texture, points, float32(img.Alpha()))
	p.freeCoords(vao, vbo)
}

func (p *glPainter) drawRaster(img *canvas.Raster, pos fyne.Position, frame fyne.Size) {
	texture := getTexture(img, p.newGlRasterTexture)
	if texture == 0 {
		return
	}

	points, vao, vbo := p.rectCoords(img.Size(), pos, frame, canvas.ImageFillStretch, 0.0, 0)
	p.glDrawTexture(texture, points, float32(img.Alpha()))
	p.freeCoords(vao, vbo)
}

func (p *glPainter) drawGradient(o fyne.CanvasObject, texCreator func(fyne.CanvasObject) uint32, pos fyne.Position, frame fyne.Size) {
	texture := getTexture(o, texCreator)
	if texture == 0 {
		return
	}

	points, vao, vbo := p.rectCoords(o.Size(), pos, frame, canvas.ImageFillStretch, 0.0, 0)
	p.glDrawTexture(texture, points, 1.0)
	p.freeCoords(vao, vbo)
}

func (p *glPainter) drawRectangle(rect *canvas.Rectangle, pos fyne.Position, frame fyne.Size) {
	points, vao, vbo := p.rectCoords(rect.Size(), pos, frame, canvas.ImageFillStretch, 0.0, 0)
	texture := getTexture(rect, p.newGlRectTexture)

	p.glDrawTexture(texture, points, 1.0)
	p.freeCoords(vao, vbo)
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

	points, vao, vbo := p.rectCoords(size, pos, frame, canvas.ImageFillStretch, 0.0, 0)
	texture := getTexture(text, p.newGlTextTexture)

	p.glDrawTexture(texture, points, 1.0)
	p.freeCoords(vao, vbo)
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
