package widget

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/canvas"
import "github.com/fyne-io/fyne/layout"
import "github.com/fyne-io/fyne/theme"

// Button widget has a text label and triggers an event func when clicked
type Button struct {
	baseWidget
	Style ButtonStyle

	OnTapped func()

	label      *canvas.Text
	background *canvas.Rectangle
}

// ButtonStyle determines the behaviour and rendering of a button.
type ButtonStyle int

const (
	// DefaultButton is the standard button style
	DefaultButton ButtonStyle = iota
	// PrimaryButton that should be more prominent to the user
	PrimaryButton
)

// MinSize calculates the minimum size of a button.
// This is based on the contained text with a standard amount of padding added.
func (b *Button) MinSize() fyne.Size {
	return b.label.MinSize().Add(fyne.NewSize(theme.Padding()*4, theme.Padding()*2))
}

// Layout the components of the button widget
func (b *Button) Layout(size fyne.Size) []fyne.CanvasObject {
	layout.NewMaxLayout().Layout(b.objects, size)

	return b.objects
}

// OnMouseDown is called when a mouse down event is captured and triggers any tap handler
func (b *Button) OnMouseDown(*fyne.MouseEvent) {
	if b.OnTapped != nil {
		b.OnTapped()
	}
}

// ApplyTheme is called when the Button may need to update it's look
func (b *Button) ApplyTheme() {
	b.label.Color = theme.TextColor()

	if b.Style == PrimaryButton {
		b.background.FillColor = theme.PrimaryColor()
	} else {
		b.background.FillColor = theme.ButtonColor()
	}
}

// NewButton creates a new button widget with the set label and tap handler
func NewButton(label string, tapped func()) *Button {
	text := canvas.NewText(label, theme.TextColor())
	text.Alignment = fyne.TextAlignCenter
	bg := canvas.NewRectangle(theme.ButtonColor())

	return &Button{
		baseWidget{
			objects: []fyne.CanvasObject{
				bg,
				text,
			},
		},
		DefaultButton,
		tapped,
		text,
		bg,
	}
}
