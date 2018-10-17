package widget

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/canvas"
import "github.com/fyne-io/fyne/theme"

type buttonRenderer struct {
	background *canvas.Rectangle
	icon       *canvas.Image
	label      *canvas.Text

	objects []fyne.CanvasObject
	button  *Button
}

// MinSize calculates the minimum size of a button.
// This is based on the contained text, any icon that is set and a standard
// amount of padding added.
func (b *buttonRenderer) MinSize() fyne.Size {
	min := b.label.MinSize().Add(fyne.NewSize(theme.Padding()*4, theme.Padding()*2))

	if b.icon != nil {
		min = min.Add(fyne.NewSize(theme.IconInlineSize()+theme.Padding(), 0))
	}
	return min
}

// Layout the components of the button widget
func (b *buttonRenderer) Layout(size fyne.Size) {
	b.background.Resize(size)

	if b.button.Icon == nil {
		b.label.Resize(size)
		b.label.Move(fyne.NewPos(0, 0))
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

// ApplyTheme is called when the Button may need to update it's look
func (b *buttonRenderer) ApplyTheme() {
	b.label.Color = theme.TextColor()

	b.Refresh()
}

func (b *buttonRenderer) Refresh() {
	b.label.Text = b.button.Text

	if b.button.Style == PrimaryButton {
		b.background.FillColor = theme.PrimaryColor()
	} else {
		b.background.FillColor = theme.ButtonColor()
	}

	// TODO don't tweak alpha, add visibility instead
	if b.button.Icon != nil {
		if b.icon == nil {
			b.icon = canvas.NewImageFromResource(b.button.Icon)
			b.objects = append(b.objects, b.icon)
		}
		b.icon.Translucency = 0
		b.icon.File = b.button.Icon.CachePath()
	} else if b.icon != nil {
		b.icon.Translucency = 1
	}

	b.Layout(b.button.CurrentSize())
	fyne.RefreshObject(b.button)
}

func (b *buttonRenderer) Objects() []fyne.CanvasObject {
	return b.objects
}

// Button widget has a text label and triggers an event func when clicked
type Button struct {
	baseWidget
	Text  string
	Style ButtonStyle
	Icon  fyne.Resource

	OnTapped func() `json:"-"`
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
		go b.OnTapped()
	}
}

func (b *Button) createRenderer() fyne.WidgetRenderer {
	var icon *canvas.Image
	if b.Icon != nil {
		icon = canvas.NewImageFromResource(b.Icon)
	}

	text := canvas.NewText(b.Text, theme.TextColor())
	text.Alignment = fyne.TextAlignCenter
	bg := canvas.NewRectangle(theme.ButtonColor())

	objects := []fyne.CanvasObject{
		bg,
		text,
	}
	if icon != nil {
		objects = append(objects, icon)
	}

	return &buttonRenderer{bg, icon, text, objects, b}
}

// Renderer is a private method to Fyne which links this widget to it's renderer
func (b *Button) Renderer() fyne.WidgetRenderer {
	if b.renderer == nil {
		b.renderer = b.createRenderer()
	}

	return b.renderer
}

// SetText allows the button label to be changed
func (b *Button) SetText(text string) {
	b.Text = text

	b.Renderer().Refresh()
}

// SetIcon updates the icon on a label - pass nil to hide an icon
func (b *Button) SetIcon(icon fyne.Resource) {
	b.Icon = icon

	b.Renderer().Refresh()
}

// NewButton creates a new button widget with the set label and tap handler
func NewButton(label string, tapped func()) *Button {
	button := &Button{baseWidget{}, label, DefaultButton, nil, tapped}

	button.Renderer().Layout(button.MinSize())
	return button
}

// NewButtonWithIcon creates a new button widget with the specified label,
// themed icon and tap handler
func NewButtonWithIcon(label string, icon fyne.Resource, tapped func()) *Button {
	button := &Button{baseWidget{}, label, DefaultButton, icon, tapped}

	button.Renderer().Layout(button.MinSize())
	return button
}
