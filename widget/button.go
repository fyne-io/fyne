package widget

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/canvas"
import "github.com/fyne-io/fyne/theme"

type buttonLayout struct {
	background *canvas.Rectangle
	icon       *canvas.Image
	label      *canvas.Text
}

// MinSize calculates the minimum size of a button.
// This is based on the contained text, any icon that is set and a standard
// amount of padding added.
func (b *buttonLayout) MinSize([]fyne.CanvasObject) fyne.Size {
	min := b.label.MinSize().Add(fyne.NewSize(theme.Padding()*4, theme.Padding()*2))

	if b.icon != nil {
		min = min.Add(fyne.NewSize(theme.IconInlineSize()+theme.Padding(), 0))
	}
	return min
}

// Layout the components of the button widget
func (b *buttonLayout) Layout(_ []fyne.CanvasObject, size fyne.Size) {
	b.background.Resize(size)

	if b.icon == nil {
		b.label.Resize(size)
	} else {
		offset := fyne.NewSize(theme.IconInlineSize()+theme.Padding(), 0)
		labelSize := size.Subtract(offset)
		b.label.Resize(labelSize)
		b.label.Move(fyne.NewPos(theme.IconInlineSize()+theme.Padding(), 0))

		b.icon.Resize(fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
		b.icon.Move(fyne.NewPos(
			(size.Width-theme.IconInlineSize()-b.label.MinSize().Width)/2,
			(size.Height-theme.IconInlineSize())/2))
	}
}

// Button widget has a text label and triggers an event func when clicked
type Button struct {
	baseWidget
	Style ButtonStyle

	OnTapped func()
	layout   *buttonLayout
	resource fyne.ThemedResource
}

// ButtonStyle determines the behaviour and rendering of a button.
type ButtonStyle int

const (
	// DefaultButton is the standard button style
	DefaultButton ButtonStyle = iota
	// PrimaryButton that should be more prominent to the user
	PrimaryButton
)

// OnMouseDown is called when a mouse down event is captured and triggers any tap handler
func (b *Button) OnMouseDown(*fyne.MouseEvent) {
	if b.OnTapped != nil {
		b.OnTapped()
	}
}

// ApplyTheme is called when the Button may need to update it's look
func (b *Button) ApplyTheme() {
	b.layout.label.Color = theme.TextColor()

	if b.Style == PrimaryButton {
		b.layout.background.FillColor = theme.PrimaryColor()
	} else {
		b.layout.background.FillColor = theme.ButtonColor()
	}

	if b.resource != nil {
		b.layout.icon.File = b.resource.CurrentResource().CachePath()
	}
}

func constructButton(label string, resource fyne.ThemedResource, tapped func()) *Button {
	var icon *canvas.Image
	if resource != nil {
		icon = canvas.NewImageFromResource(resource.CurrentResource())
	}

	text := canvas.NewText(label, theme.TextColor())
	text.Alignment = fyne.TextAlignCenter
	bg := canvas.NewRectangle(theme.ButtonColor())
	layout := &buttonLayout{bg, icon, text}

	objects := []fyne.CanvasObject{
		bg,
		text,
	}
	if icon != nil {
		objects = append(objects, icon)
	}

	b := &Button{
		baseWidget{
			objects: objects,
			layout:  layout,
		},
		DefaultButton,
		tapped,
		layout,
		resource,
	}

	b.Layout(b.MinSize())
	return b
}

// NewButton creates a new button widget with the set label and tap handler
func NewButton(label string, tapped func()) *Button {
	return constructButton(label, nil, tapped)
}

// NewButtonWithIcon creates a new button widget with the specified label,
// themed icon and tap handler
func NewButtonWithIcon(label string, icon fyne.ThemedResource, tapped func()) *Button {
	return constructButton(label, icon, tapped)
}
