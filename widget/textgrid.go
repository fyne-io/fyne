package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

// TextGrid is a monospaced grid of characters.
// This is designed to be used by a text editor or advanced test presentation.
type TextGrid struct {
	baseWidget
}

// Resize sets a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (t *TextGrid) Resize(size fyne.Size) {
	t.resize(size, t)
}

// Move the widget to a new position, relative to it's parent.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (t *TextGrid) Move(pos fyne.Position) {
	t.move(pos, t)
}

// MinSize returns the smallest size this widget can shrink to
func (t *TextGrid) MinSize() fyne.Size {
	return t.minSize(t)
}

// Show this widget, if it was previously hidden
func (t *TextGrid) Show() {
	t.show(t)
}

// Hide this widget, if it was previously visible
func (t *TextGrid) Hide() {
	t.hide(t)
}

// CreateRenderer is a private method to Fyne which links this widget to it's renderer
func (t *TextGrid) CreateRenderer() fyne.WidgetRenderer {
	render := &textGridRender{text: t}
	render.ensureGrid()

	cell := canvas.NewText("M", color.White)
	cell.TextStyle.Monospace = true
	render.cellSize = cell.MinSize()

	return render
}

// NewTextGrid creates a new textgrid widget with the specified string content.
func NewTextGrid(content string) *TextGrid {
	grid := &TextGrid{}
	return grid
}

type textGridRender struct {
	text *TextGrid

	cols, rows int

	cellSize fyne.Size
	objects  []fyne.CanvasObject
}

func newTextCell(str string) *canvas.Text {
	text := canvas.NewText(str, theme.TextColor())
	text.TextStyle.Monospace = true

	if str == "·" || str == "→" || str == "↵" {
		r, g, b, a := theme.PlaceHolderColor().RGBA()
		text.Color = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a) / 3}
	}

	return text
}

func (t *textGridRender) ensureGrid() {
	t.cols = 8
	t.rows = 3

	t.objects = append(t.objects, newTextCell("S"))
	t.objects = append(t.objects, newTextCell("o"))
	t.objects = append(t.objects, newTextCell("m"))
	t.objects = append(t.objects, newTextCell("e"))
	t.objects = append(t.objects, newTextCell("↵"))
	t.objects = append(t.objects, newTextCell(""))
	t.objects = append(t.objects, newTextCell(""))
	t.objects = append(t.objects, newTextCell(""))

	t.objects = append(t.objects, newTextCell("·"))
	t.objects = append(t.objects, newTextCell("·"))
	t.objects = append(t.objects, newTextCell("T"))
	t.objects = append(t.objects, newTextCell("e"))
	t.objects = append(t.objects, newTextCell("x"))
	t.objects = append(t.objects, newTextCell("t"))
	t.objects = append(t.objects, newTextCell("↵"))
	t.objects = append(t.objects, newTextCell(""))

	t.objects = append(t.objects, newTextCell("→"))
	t.objects = append(t.objects, newTextCell(""))
	t.objects = append(t.objects, newTextCell(""))
	t.objects = append(t.objects, newTextCell(""))
	t.objects = append(t.objects, newTextCell("T"))
	t.objects = append(t.objects, newTextCell("a"))
	t.objects = append(t.objects, newTextCell("b"))
	t.objects = append(t.objects, newTextCell("s"))
}

func (t *textGridRender) Layout(size fyne.Size) {
	i := 0
	cellPos := fyne.NewPos(0, 0)
	for y := 0; y < t.rows; y++ {
		for x := 0; x < t.cols; x++ {
			t.objects[i].Move(cellPos)

			cellPos.X += t.cellSize.Width
			i++
		}

		cellPos.X = 0
		cellPos.Y += t.cellSize.Height
	}
}

func (t *textGridRender) MinSize() fyne.Size {
	return fyne.NewSize(t.cellSize.Width*t.cols, t.cellSize.Height*t.rows)
}

func (t *textGridRender) Refresh() {
}

func (t *textGridRender) ApplyTheme() {
}

func (t *textGridRender) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (t *textGridRender) Objects() []fyne.CanvasObject {
	return t.objects
}
