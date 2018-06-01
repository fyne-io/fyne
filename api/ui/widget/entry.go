package widget

import "fmt"

import "github.com/fyne-io/fyne/api/ui"
import "github.com/fyne-io/fyne/api/ui/canvas"
import "github.com/fyne-io/fyne/api/ui/layout"
import "github.com/fyne-io/fyne/api/ui/theme"

// Entry widget allows simple text to be input when focussed.
type Entry struct {
	baseWidget

	OnChanged func(string)

	focussed bool

	label   *canvas.Text
	bg, box *canvas.Rectangle
}

// MinSize calculates the minimum size of an entry widget.
// This is based on the contained text with a standard amount of padding added.
func (e *Entry) MinSize() ui.Size {
	return e.label.MinSize().Add(ui.NewSize(theme.Padding()*4, theme.Padding()*2))
}

// Layout the components of the entry widget.
func (e *Entry) Layout(size ui.Size) []ui.CanvasObject {
	layout.NewMaxLayout().Layout(e.objects, size)

	e.box.Resize(e.box.CurrentSize().Subtract(ui.NewSize(theme.Padding(), theme.Padding())))
	e.box.Move(ui.NewPos(theme.Padding()/2, theme.Padding()/2))

	e.label.Resize(e.label.CurrentSize().Subtract(ui.NewSize(theme.Padding()*2, theme.Padding()*2)))
	e.label.Move(ui.NewPos(theme.Padding(), theme.Padding()))

	return e.objects
}

// Text returns the current content of the Entry.
func (e *Entry) Text() string {
	return e.label.Text
}

// SetText manually sets the text of the Entry to the given text value.
func (e *Entry) SetText(text string) {
	e.label.Text = text

	if e.OnChanged != nil {
		e.OnChanged(text)
	}

	ui.GetCanvas(e).Refresh(e)
}

// OnFocusGained is called when the Entry has been given focus.
func (e *Entry) OnFocusGained() {
	e.focussed = true
	e.ApplyTheme()
	ui.GetCanvas(e).Refresh(e)
}

// OnFocusLost is called when the Entry has had focus removed.
func (e *Entry) OnFocusLost() {
	e.focussed = false
	e.ApplyTheme()
	ui.GetCanvas(e).Refresh(e)
}

// OnKeyDown receives key input events when the Entry widget is focussed.
func (e *Entry) OnKeyDown(key *ui.KeyEvent) {
	if key.Name == "BackSpace" {
		runes := []rune(e.Text())
		substr := string(runes[0 : len(runes)-1])

		e.SetText(substr)
		return
	}

	if key.String != "" {
		e.SetText(fmt.Sprintf("%s%s", e.Text(), key.String))
	}
}

// ApplyTheme is called when the Entry may need to update it's look.
func (e *Entry) ApplyTheme() {
	e.label.Color = theme.TextColor()
	if e.focussed {
		e.bg.FillColor = theme.FocusColor()
	} else {
		e.bg.FillColor = theme.ButtonColor()
	}
	e.box.FillColor = theme.BackgroundColor()
}

// NewEntry creates a new entry widget.
func NewEntry() *Entry {
	text := canvas.NewText("", theme.TextColor())
	bg := canvas.NewRectangle(theme.ButtonColor())
	box := canvas.NewRectangle(theme.BackgroundColor())

	return &Entry{
		baseWidget{
			objects: []ui.CanvasObject{
				bg,
				box,
				text,
			},
		},
		nil,
		false,
		text,
		bg,
		box,
	}
}
