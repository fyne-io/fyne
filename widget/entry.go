package widget

import "fmt"

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/canvas"
import "github.com/fyne-io/fyne/theme"

type entryRenderer struct {
	label   *canvas.Text
	bg, box *canvas.Rectangle

	objects []fyne.CanvasObject
	entry   *Entry
}

// MinSize calculates the minimum size of an entry widget.
// This is based on the contained text with a standard amount of padding added.
func (e *entryRenderer) MinSize() fyne.Size {
	var textSize fyne.Size
	if e.label.Text == "" {
		textSize = fyne.GetDriver().RenderedTextSize("M", e.label.TextSize, e.label.TextStyle)
	} else {
		textSize = e.label.MinSize()
	}

	return textSize.Add(fyne.NewSize(theme.Padding()*4, theme.Padding()*2))
}

// Layout the components of the entry widget.
func (e *entryRenderer) Layout(size fyne.Size) {
	e.bg.Resize(size)

	e.box.Resize(size.Subtract(fyne.NewSize(theme.Padding(), theme.Padding())))
	e.box.Move(fyne.NewPos(theme.Padding()/2, theme.Padding()/2))

	e.label.Resize(size.Subtract(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)))
	e.label.Move(fyne.NewPos(theme.Padding(), theme.Padding()))
}

// ApplyTheme is called when the Entry may need to update it's look.
func (e *entryRenderer) ApplyTheme() {
	e.label.Color = theme.TextColor()
	e.box.FillColor = theme.BackgroundColor()
	e.Refresh()
}

func (e *entryRenderer) Refresh() {
	e.label.Text = e.entry.Text

	if e.entry.focused {
		e.bg.FillColor = theme.FocusColor()
	} else {
		e.bg.FillColor = theme.ButtonColor()
	}

	fyne.RefreshObject(e.entry)
}

func (e *entryRenderer) Objects() []fyne.CanvasObject {
	return e.objects
}

// Entry widget allows simple text to be input when focussed.
type Entry struct {
	baseWidget

	Text      string
	OnChanged func(string) `json:"-"`

	focused bool
}

// SetText manually sets the text of the Entry to the given text value.
func (e *Entry) SetText(text string) {
	e.Text = text

	if e.OnChanged != nil {
		e.OnChanged(text)
	}

	e.Renderer().Refresh()
}

// OnFocusGained is called when the Entry has been given focus.
func (e *Entry) OnFocusGained() {
	e.focused = true

	e.Renderer().Refresh()
}

// OnFocusLost is called when the Entry has had focus removed.
func (e *Entry) OnFocusLost() {
	e.focused = false

	e.Renderer().Refresh()
}

// Focused returns whether or not this Entry has focus.
func (e *Entry) Focused() bool {
	return e.focused
}

// OnKeyDown receives key input events when the Entry widget is focussed.
func (e *Entry) OnKeyDown(key *fyne.KeyEvent) {
	if key.Name == "BackSpace" {
		runes := []rune(e.Text)
		substr := string(runes[0 : len(runes)-1])

		e.SetText(substr)
		return
	}

	if key.String != "" {
		e.SetText(fmt.Sprintf("%s%s", e.Text, key.String))
	}
}

func (e *Entry) createRenderer() fyne.WidgetRenderer {
	text := canvas.NewText(e.Text, theme.TextColor())
	bg := canvas.NewRectangle(theme.ButtonColor())
	box := canvas.NewRectangle(theme.BackgroundColor())

	return &entryRenderer{text, bg, box, []fyne.CanvasObject{bg, box, text}, e}
}

// Renderer is a private method to Fyne which links this widget to it's renderer
func (e *Entry) Renderer() fyne.WidgetRenderer {
	if e.renderer == nil {
		e.renderer = e.createRenderer()
	}

	return e.renderer
}

// NewEntry creates a new entry widget.
func NewEntry() *Entry {
	e := &Entry{
		baseWidget{},
		"",
		nil,
		false,
	}

	e.Renderer().Layout(e.MinSize())
	return e
}
