package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/theme"
	"image/color"
	"strings"
)

type buttonRenderer struct {
	icon  *canvas.Image
	label *canvas.Text

	objects []fyne.CanvasObject
	button  *Button
}

// MinSize calculates the minimum size of a button.
// This is based on the contained text, any icon that is set and a standard
// amount of padding added.
func (b *buttonRenderer) MinSize() fyne.Size {
	var min fyne.Size

	if b.button.Text != "" {
		min = b.label.MinSize().Add(fyne.NewSize(theme.Padding()*4, theme.Padding()*2))
		if b.icon != nil {
			min = min.Add(fyne.NewSize(theme.IconInlineSize()+theme.Padding(), 0))
		}
	} else if b.icon != nil {
		min = fyne.NewSize(theme.IconInlineSize()+theme.Padding()*2, theme.IconInlineSize()+theme.Padding()*2)
	}

	return min
}

// Layout the components of the button widget
func (b *buttonRenderer) Layout(size fyne.Size) {
	if b.button.Text != "" {
		inner := size.Subtract(fyne.NewSize(theme.Padding()*4, theme.Padding()*2))

		if b.button.Icon == nil {
			b.label.Resize(inner)
			b.label.Move(fyne.NewPos(theme.Padding()*2, theme.Padding()))
		} else {
			offset := fyne.NewSize(theme.IconInlineSize(), 0)
			labelSize := inner.Subtract(offset)
			b.label.Resize(labelSize)
			b.label.Move(fyne.NewPos(theme.IconInlineSize()+theme.Padding()*2, theme.Padding()))

			b.icon.Resize(fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
			b.icon.Move(fyne.NewPos(
				(size.Width-theme.IconInlineSize()-b.label.MinSize().Width-theme.Padding())/2,
				(size.Height-theme.IconInlineSize())/2))
		}
	} else {
		b.icon.Resize(fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
		b.icon.Move(fyne.NewPos(theme.Padding(), theme.Padding()))
	}
}

// ApplyTheme is called when the Button may need to update its look
func (b *buttonRenderer) ApplyTheme() {
	b.label.Color = theme.TextColor()
	if b.button.disabled {
		b.label.Color = theme.DisabledTextColor()
	}

	b.Refresh()
}

func (b *buttonRenderer) BackgroundColor() color.Color {
	if b.button.Style == PrimaryButton && !b.button.disabled {
		return theme.PrimaryColor()
	} else if b.button.disabled {
		return theme.DisabledButtonColor()
	}

	if b.button.hovered && !b.button.disabled {
		return theme.HoverColor()
	}
	return theme.ButtonColor()
}

func (b *buttonRenderer) Refresh() {
	b.label.Text = b.button.Text

	if b.button.Icon != nil {
		if b.icon == nil {
			b.icon = canvas.NewImageFromResource(b.button.Icon)
			b.objects = append(b.objects, b.icon)
		} else {
			if b.button.disabled {
				// if the icon has changed, create a new disabled version
				// if we could be sure that button.Icon is only ever set through the button.SetIcon method, we could remove this
				if !strings.HasSuffix(b.button.disabledIcon.Name(), b.button.Icon.Name()){
					b.icon.Resource = theme.NewDisabledResource(b.button.Icon)
				} else {
					b.icon.Resource = b.button.disabledIcon
				}
			} else {
				b.icon.Resource = b.button.Icon
			}
		}
		b.icon.Hidden = false
	} else if b.icon != nil {
		b.icon.Hidden = true
	}

	b.Layout(b.button.Size())
	canvas.Refresh(b.button)
}

func (b *buttonRenderer) Objects() []fyne.CanvasObject {
	return b.objects
}

func (b *buttonRenderer) Destroy() {
}

// Button widget has a text label and triggers an event func when clicked
type Button struct {
	baseWidget
	Text  string
	Style ButtonStyle
	Icon  fyne.Resource
	disabledIcon fyne.Resource

	OnTapped func() `json:"-"`
	hovered  bool
}

// ButtonStyle determines the behaviour and rendering of a button.
type ButtonStyle int

const (
	// DefaultButton is the standard button style
	DefaultButton ButtonStyle = iota
	// PrimaryButton that should be more prominent to the user
	PrimaryButton
)

// Resize sets a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (b *Button) Resize(size fyne.Size) {
	b.resize(size, b)
}

// Move the widget to a new position, relative to its parent.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (b *Button) Move(pos fyne.Position) {
	b.move(pos, b)
}

// MinSize returns the smallest size this widget can shrink to
func (b *Button) MinSize() fyne.Size {
	return b.minSize(b)
}

// Show this widget, if it was previously hidden
func (b *Button) Show() {
	b.show(b)
}

// Hide this widget, if it was previously visible
func (b *Button) Hide() {
	b.hide(b)
}

// Enable this widget, if it was previously disabled
func (b *Button) Enable() {
	b.enable(b)
	Renderer(b).ApplyTheme()
}

// Disable this widget, if it was previously enabled
func (b *Button) Disable() {
	b.disable(b)
	Renderer(b).ApplyTheme()
}

// Tapped is called when a pointer tapped event is captured and triggers any tap handler
func (b *Button) Tapped(*fyne.PointEvent) {
	if b.OnTapped != nil && !b.disabled {
		b.OnTapped()
	}
}

// TappedSecondary is called when a secondary pointer tapped event is captured
func (b *Button) TappedSecondary(*fyne.PointEvent) {
}

// MouseIn is called when a desktop pointer enters the widget
func (b *Button) MouseIn(*desktop.MouseEvent) {
	b.hovered = true
	Refresh(b)
}

// MouseOut is called when a desktop pointer exits the widget
func (b *Button) MouseOut() {
	b.hovered = false
	Refresh(b)
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (b *Button) MouseMoved(*desktop.MouseEvent) {
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (b *Button) CreateRenderer() fyne.WidgetRenderer {
	var icon *canvas.Image
	if b.Icon != nil {
		icon = canvas.NewImageFromResource(b.Icon)
	}

	text := canvas.NewText(b.Text, theme.TextColor())
	text.Alignment = fyne.TextAlignCenter

	objects := []fyne.CanvasObject{
		text,
	}
	if icon != nil {
		objects = append(objects, icon)
	}

	return &buttonRenderer{icon, text, objects, b}
}

// SetText allows the button label to be changed
func (b *Button) SetText(text string) {
	b.Text = text

	Refresh(b)
}

// SetIcon updates the icon on a label - pass nil to hide an icon
func (b *Button) SetIcon(icon fyne.Resource) {
	b.Icon = icon

	if icon != nil {
		b.disabledIcon = theme.NewDisabledResource(icon)
	} else {
		b.disabledIcon = nil
	}

	Refresh(b)
}

// NewButton creates a new button widget with the set label and tap handler
func NewButton(label string, tapped func()) *Button {
	button := &Button{baseWidget{}, label, DefaultButton, nil, nil, tapped, false}

	Renderer(button).Layout(button.MinSize())
	return button
}

// NewButtonWithIcon creates a new button widget with the specified label, themed icon and tap handler
func NewButtonWithIcon(label string, icon fyne.Resource, tapped func()) *Button {
	button := &Button{baseWidget{}, label, DefaultButton, icon, theme.NewDisabledResource(icon), tapped, false}

	Renderer(button).Layout(button.MinSize())
	return button
}
