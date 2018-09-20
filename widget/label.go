package widget

import "bufio"
import "strings"

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/canvas"
import "github.com/fyne-io/fyne/theme"

type labelRenderer struct {
	objects []fyne.CanvasObject

	background *canvas.Rectangle
	texts      []fyne.CanvasObject

	label *Label
}

func (l *labelRenderer) parseText(text string) []fyne.CanvasObject {
	if strings.Contains(text, "\n") {
		texts := []fyne.CanvasObject{}
		s := bufio.NewScanner(strings.NewReader(text))
		for s.Scan() {
			texts = append(texts, canvas.NewText(s.Text(), theme.TextColor()))
		}

		return texts
	}

	return []fyne.CanvasObject{canvas.NewText(text, theme.TextColor())}
}

// MinSize calculates the minimum size of a label.
// This is based on the contained text with a standard amount of padding added.
func (l *labelRenderer) MinSize() fyne.Size {
	height := 0
	width := 0
	for _, text := range l.texts {
		min := text.MinSize()
		height += min.Height
		width = fyne.Max(width, min.Width)
	}

	return fyne.NewSize(width, height).Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2))
}

func (l *labelRenderer) Layout(size fyne.Size) {
	yPos := 0
	lineHeight := size.Height
	if len(l.texts) > 1 {
		lineHeight = size.Height / len(l.texts)
	}
	lineSize := fyne.NewSize(size.Width, lineHeight)
	for _, text := range l.texts {
		text.Resize(lineSize)
		text.Move(fyne.NewPos(0, yPos))
		yPos += lineHeight
	}

	l.background.Resize(size)
}

func (l *labelRenderer) Objects() []fyne.CanvasObject {
	return l.objects
}

// ApplyTheme is called when the Label may need to update it's look
func (l *labelRenderer) ApplyTheme() {
	l.background.FillColor = theme.BackgroundColor()

	for _, text := range l.texts {
		text.(*canvas.Text).Color = theme.TextColor()
	}
}

func (l *labelRenderer) Refresh() {
	for _, text := range l.texts {
		text.(*canvas.Text).Alignment = l.label.Alignment
		text.(*canvas.Text).TextStyle = l.label.TextStyle
		text.(*canvas.Text).Text = l.label.Text
	}

	fyne.RefreshObject(l.label)
}

// Label widget is a basic text component with appropriate padding and layout.
type Label struct {
	baseWidget

	Text      string         // The content of the label
	Alignment fyne.TextAlign // The alignment of the Text
	TextStyle fyne.TextStyle // The style of the label text
}

// SetText updates the text of the label widget
func (l *Label) SetText(text string) {
	l.Text = text

	l.Renderer().Refresh()
}

func (l *Label) createRenderer() fyne.WidgetRenderer {
	render := &labelRenderer{label: l}

	// TODO move this to a renderer method and call on setText too
	texts := render.parseText(l.Text)
	bg := canvas.NewRectangle(theme.ButtonColor())
	objects := []fyne.CanvasObject{bg}
	objects = append(objects, texts...)

	render.objects = objects
	render.background = bg
	render.texts = texts

	return render
}

// Renderer is a private method to Fyne which links this widget to it's renderer
func (l *Label) Renderer() fyne.WidgetRenderer {
	if l.renderer == nil {
		l.renderer = l.createRenderer()
	}

	return l.renderer
}

// NewLabel creates a new layout widget with the set text content
func NewLabel(text string) *Label {
	var style fyne.TextStyle

	l := &Label{
		baseWidget{},
		text,
		fyne.TextAlignLeading,
		style,
	}

	l.Renderer().Layout(l.MinSize())
	return l
}
