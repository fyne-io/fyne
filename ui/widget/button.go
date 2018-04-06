package widget

import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/canvas"
import "github.com/fyne-io/fyne/ui/event"
import "github.com/fyne-io/fyne/ui/layout"
import "github.com/fyne-io/fyne/ui/theme"

// Button widget has a text label and triggers an event func when clicked
type Button struct {
	baseWidget
	Style ButtonStyle

	OnClicked func(*event.MouseEvent)

	label      *canvas.TextObject
	background *canvas.RectangleObject
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
func (b *Button) Layout() []ui.CanvasObject {
	if b.Style == PrimaryButton {
		b.background.FillColor = theme.PrimaryColor()
	}
	layout.NewMaxLayout().Layout(b.objects, b.Size)

	return b.objects
}

// NewButton creates a new button widget with the set label and click handler
func NewButton(label string, clicked func(*event.MouseEvent)) *Button {
	text := canvas.NewText(label)
	bg := canvas.NewRectangle(theme.ButtonColor())

	return &Button{
		baseWidget{
			objects: []ui.CanvasObject{
				bg,
				text,
			},
		},
		DefaultButton,
		clicked,
		text,
		bg,
	}
}
