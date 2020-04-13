package widget

import (
	"image/color"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/binding"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/theme"
)

type buttonRenderer struct {
	*shadowingRenderer

	icon   *canvas.Image
	label  *canvas.Text
	button *Button
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
	b.layoutShadow(size, fyne.NewPos(0, 0))
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

// applyTheme updates this button to match the current theme
func (b *buttonRenderer) applyTheme() {
	b.label.TextSize = theme.TextSize()
	b.label.Color = theme.TextColor()
	if b.button.Disabled() {
		b.label.Color = theme.DisabledTextColor()
	}
}

func (b *buttonRenderer) BackgroundColor() color.Color {
	switch {
	case b.button.Disabled():
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
	b.applyTheme()
	b.label.Text = b.button.Text

	if b.button.Icon != nil && b.button.Visible() {
		if b.icon == nil {
			b.icon = canvas.NewImageFromResource(b.button.Icon)
			b.setObjects(append(b.Objects(), b.icon))
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
		b.icon.Show()
	} else if b.icon != nil {
		b.icon.Hide()
	}

	b.Layout(b.button.Size())
	canvas.Refresh(b.button.super())
}

// Button widget has a text label and triggers an event func when clicked
type Button struct {
	DisableableWidget
	Text         string
	Style        ButtonStyle
	Icon         fyne.Resource
	disabledIcon fyne.Resource

	OnTapped   func() `json:"-"`
	hovered    bool
	HideShadow bool

	tapBind    *binding.Bool
	textBind   *binding.String
	iconBind   *binding.Resource
	textNotify *binding.NotifyFunction
	iconNotify *binding.NotifyFunction
}

// ButtonStyle determines the behaviour and rendering of a button.
type ButtonStyle int

const (
	// DefaultButton is the standard button style
	DefaultButton ButtonStyle = iota
	// PrimaryButton that should be more prominent to the user
	PrimaryButton
)

// Tapped is called when a pointer tapped event is captured and triggers any tap handler
func (b *Button) Tapped(*fyne.PointEvent) {
	if !b.Disabled() {
		if b.tapBind != nil {
			b.tapBind.Set(true)
		}
		if b.OnTapped != nil {
			b.OnTapped()
		}
	}
}

// MouseIn is called when a desktop pointer enters the widget
func (b *Button) MouseIn(*desktop.MouseEvent) {
	b.hovered = true
	b.Refresh()
}

// MouseOut is called when a desktop pointer exits the widget
func (b *Button) MouseOut() {
	b.hovered = false
	b.Refresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (b *Button) MouseMoved(*desktop.MouseEvent) {
}

// MinSize returns the size that this widget should not shrink below
func (b *Button) MinSize() fyne.Size {
	b.ExtendBaseWidget(b)
	return b.BaseWidget.MinSize()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (b *Button) CreateRenderer() fyne.WidgetRenderer {
	b.ExtendBaseWidget(b)
	var icon *canvas.Image
	if b.Icon != nil {
		icon = canvas.NewImageFromResource(b.Icon)
	}

	text := canvas.NewText(b.Text, theme.TextColor())
	text.Alignment = fyne.TextAlignCenter

	objects := []fyne.CanvasObject{
		text,
	}
	shadowLevel := buttonLevel
	if b.HideShadow {
		shadowLevel = baseLevel
	}
	if icon != nil {
		objects = append(objects, icon)
	}

	return &buttonRenderer{newShadowingRenderer(objects, shadowLevel), icon, text, b}
}

// SetText allows the button label to be changed
func (b *Button) SetText(text string) {
	if b.Text != text {
		b.Text = text
		b.Refresh()
	}
}

// Cursor returns the cursor type of this widget
func (b *Button) Cursor() desktop.Cursor {
	return desktop.DefaultCursor
}

// SetIcon updates the icon on a label - pass nil to hide an icon
func (b *Button) SetIcon(icon fyne.Resource) {
	if b.Icon != icon {
		b.Icon = icon

		if icon != nil {
			b.disabledIcon = theme.NewDisabledResource(icon)
		} else {
			b.disabledIcon = nil
		}

		b.Refresh()
	}
}

// BindTapped binds the Button's OnTap to the given data binding.
// Returns the Button for chaining.
func (b *Button) BindTapped(data *binding.Bool) *Button {
	b.tapBind = data
	return b
}

// UnbindTapped unbinds the Button's OnTap from the data binding (if any).
// Returns the Button for chaining.
func (b *Button) UnbindTapped() *Button {
	b.tapBind = nil
	return b
}

// BindText binds the Button's Text to the given data binding.
// Returns the Button for chaining.
func (b *Button) BindText(data *binding.String) *Button {
	b.textBind = data
	b.textNotify = data.AddStringListener(b.SetText)
	return b
}

// UnbindText unbinds the Button's Text from the data binding (if any).
// Returns the Button for chaining.
func (b *Button) UnbindText() *Button {
	b.textBind.DeleteListener(b.textNotify)
	b.textBind = nil
	b.textNotify = nil
	return b
}

// BindIcon binds the Button's Icon to the given data binding.
// Returns the Button for chaining.
func (b *Button) BindIcon(data *binding.Resource) *Button {
	b.iconBind = data
	b.iconNotify = data.AddResourceListener(b.SetIcon)
	return b
}

// UnbindIcon unbinds the Button's Icon from the data binding (if any).
// Returns the Button for chaining.
func (b *Button) UnbindIcon() *Button {
	b.iconBind.DeleteListener(b.iconNotify)
	b.iconBind = nil
	b.iconNotify = nil
	return b
}

// NewButton creates a new button widget with the set label and tap handler
func NewButton(label string, tapped func()) *Button {
	button := &Button{
		DisableableWidget: DisableableWidget{},
		Text:              label,
		Style:             DefaultButton,
		OnTapped:          tapped,
		hovered:           false,
		HideShadow:        false,
	}

	button.ExtendBaseWidget(button)
	return button
}

// NewButtonWithIcon creates a new button widget with the specified label, themed icon and tap handler
func NewButtonWithIcon(label string, icon fyne.Resource, tapped func()) *Button {
	button := &Button{
		DisableableWidget: DisableableWidget{},
		Text:              label,
		Style:             DefaultButton,
		Icon:              icon,
		disabledIcon:      theme.NewDisabledResource(icon),
		OnTapped:          tapped,
		hovered:           false,
		HideShadow:        false,
	}

	button.ExtendBaseWidget(button)
	return button
}
