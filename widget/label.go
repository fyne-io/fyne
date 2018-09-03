package widget

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/canvas"
import "github.com/fyne-io/fyne/theme"

type labelRenderer struct {
	objects []fyne.CanvasObject

	background *canvas.Rectangle
	text       *canvas.Text

	label *Label
}

// MinSize calculates the minimum size of a label.
// This is based on the contained text with a standard amount of padding added.
func (l *labelRenderer) MinSize() fyne.Size {
	return l.text.MinSize().Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2))
}

func (l *labelRenderer) Layout(size fyne.Size) {
	l.text.Resize(size)
	l.background.Resize(size)
}

func (l *labelRenderer) Objects() []fyne.CanvasObject {
	return l.objects
}

// ApplyTheme is called when the Label may need to update it's look
func (l *labelRenderer) ApplyTheme() {
	l.background.FillColor = theme.BackgroundColor()

	l.text.Color = theme.TextColor()
}

func (l *labelRenderer) Refresh() {
	l.text.Alignment = l.label.Alignment
	l.text.TextStyle = l.label.TextStyle
	l.text.Text = l.label.Text

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
	obj := canvas.NewText(l.Text, theme.TextColor())
	bg := canvas.NewRectangle(theme.ButtonColor())

	objects := []fyne.CanvasObject{
		bg,
		obj,
	}

	return &labelRenderer{objects, bg, obj, l}
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
