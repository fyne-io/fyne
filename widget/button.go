package widget

import (
	"image/color"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/theme"
)

type buttonRenderer struct {
	icon   *canvas.Image
	label  *canvas.Text
	shadow fyne.CanvasObject

	objects []fyne.CanvasObject
	button  *Button
}

func (b *buttonRenderer) padding() fyne.Size {
	if b.button.Text == "" {
		return fyne.NewSize(theme.Padding()*2, theme.Padding()*2)
	}
	return fyne.NewSize(theme.Padding()*4, theme.Padding()*2)
}

// MinSize calculates the minimum size of a button.
// This is based on the contained text, any icon that is set and a standard
// amount of padding added.
func (b *buttonRenderer) MinSize() fyne.Size {
	labelSize := b.label.MinSize()
	contentHeight := fyne.Max(labelSize.Height, theme.IconInlineSize())
	contentWidth := 0
	if b.icon != nil {
		contentWidth += theme.IconInlineSize()
	}
	if b.button.Text != "" {
		if b.icon != nil {
			contentWidth += theme.Padding()
		}
		contentWidth += labelSize.Width
	}
	return fyne.NewSize(contentWidth, contentHeight).Add(b.padding())
}

// Layout the components of the button widget
func (b *buttonRenderer) Layout(size fyne.Size) {
	if b.shadow != nil {
		if b.button.HideShadow {
			b.shadow.Hide()
		} else {
			b.shadow.Resize(size)
		}
	}
	if b.button.Text != "" {
		padding := b.padding()
		innerSize := size.Subtract(padding)
		innerOffset := fyne.NewPos(padding.Width/2, padding.Height/2)

		labelSize := b.label.MinSize()
		contentWidth := labelSize.Width

		if b.icon != nil {
			contentWidth += theme.Padding() + theme.IconInlineSize()
			iconOffset := fyne.NewPos((innerSize.Width-contentWidth)/2, (innerSize.Height-theme.IconInlineSize())/2)
			b.icon.Resize(fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
			b.icon.Move(innerOffset.Add(iconOffset))
		}
		labelOffset := fyne.NewPos((innerSize.Width+contentWidth)/2-labelSize.Width, (innerSize.Height-labelSize.Height)/2)
		b.label.Resize(labelSize)
		b.label.Move(innerOffset.Add(labelOffset))
	} else if b.icon != nil {
		b.icon.Resize(fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
		b.icon.Move(fyne.NewPos((size.Width-theme.IconInlineSize())/2, (size.Height-theme.IconInlineSize())/2))
	}
}

// ApplyTheme is called when the Button may need to update its look
func (b *buttonRenderer) ApplyTheme() {
	b.label.TextSize = theme.TextSize()
	b.label.Color = theme.TextColor()
	if b.button.Disabled() {
		b.label.Color = theme.DisabledTextColor()
	}

	b.Refresh()
}

func (b *buttonRenderer) BackgroundColor() color.Color {
	switch {
	case b.button.disabled:
		return theme.DisabledButtonColor()
	case b.button.Style == PrimaryButton:
		return theme.PrimaryColor()
	case b.button.hovered:
		return theme.HoverColor()
	default:
		return theme.ButtonColor()
	}
}

func (b *buttonRenderer) Refresh() {
	b.label.Text = b.button.Text

	if b.button.Icon != nil && b.button.Visible() {
		if b.icon == nil {
			b.icon = canvas.NewImageFromResource(b.button.Icon)
			b.objects = append(b.objects, b.icon)
		} else {
			if b.button.Disabled() {
				// if the icon has changed, create a new disabled version
				// if we could be sure that button.Icon is only ever set through the button.SetIcon method, we could remove this
				if !strings.HasSuffix(b.button.disabledIcon.Name(), b.button.Icon.Name()) {
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
	Text         string
	Style        ButtonStyle
	Icon         fyne.Resource
	disabledIcon fyne.Resource

	OnTapped   func() `json:"-"`
	hovered    bool
	HideShadow bool
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

// Disabled returns true if the widget is disabled
func (b *Button) Disabled() bool {
	return b.disabled
}

// Tapped is called when a pointer tapped event is captured and triggers any tap handler
func (b *Button) Tapped(*fyne.PointEvent) {
	if b.OnTapped != nil && !b.Disabled() {
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
	var shadow fyne.CanvasObject
	if !b.HideShadow {
		shadow = newShadow(shadowAround, theme.Padding()/2)
		objects = append(objects, shadow)
	}
	if icon != nil {
		objects = append(objects, icon)
	}

	return &buttonRenderer{icon, text, shadow, objects, b}
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
	button := &Button{baseWidget{}, label, DefaultButton, nil, nil,
		tapped, false, false}

	Renderer(button).Layout(button.MinSize())
	return button
}

// NewButtonWithIcon creates a new button widget with the specified label, themed icon and tap handler
func NewButtonWithIcon(label string, icon fyne.Resource, tapped func()) *Button {
	button := &Button{baseWidget{}, label, DefaultButton, icon, theme.NewDisabledResource(icon),
		tapped, false, false}

	Renderer(button).Layout(button.MinSize())
	return button
}
