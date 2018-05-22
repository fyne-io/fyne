package widget

import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/canvas"
import "github.com/fyne-io/fyne/ui/layout"
import "github.com/fyne-io/fyne/ui/theme"

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
func (b *Button) MinSize() ui.Size {
	return b.label.MinSize().Add(ui.NewSize(theme.Padding()*4, theme.Padding()*2))
}

// Layout the components of the button widget
func (b *Button) Layout(size ui.Size) []ui.CanvasObject {
	if b.Style == PrimaryButton {
		b.background.FillColor = theme.PrimaryColor()
	}
	layout.NewMaxLayout().Layout(b.objects, size)

	return b.objects
}

func (b *Button) OnMouseDown(*ui.MouseEvent) {
	if b.OnTapped != nil {
		b.OnTapped()
	}
}

func (b *Button) ApplyTheme() {
	b.label.Color = theme.TextColor()
	b.background.FillColor = theme.ButtonColor()
}

// NewButton creates a new button widget with the set label and tap handler
func NewButton(label string, tapped func()) *Button {
	text := canvas.NewText(label, theme.TextColor())
	bg := canvas.NewRectangle(theme.ButtonColor())

	return &Button{
		baseWidget{
			objects: []ui.CanvasObject{
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
