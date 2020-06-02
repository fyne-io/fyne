package widget

import (
	"image/color"
	"strings"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
)

const buttonTapDuration = 250

type buttonRenderer struct {
	*widget.ShadowingRenderer

	icon   *canvas.Image
	label  *canvas.Text
	button *Button
	layout fyne.Layout
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
func (b *buttonRenderer) MinSize() (size fyne.Size) {
	hasIcon := b.icon != nil
	hasLabel := b.label.Text != ""
	iconSize := fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize())
	labelSize := b.label.MinSize()
	if hasLabel {
		size.Width = labelSize.Width
	}
	if hasIcon {
		if hasLabel {
			size.Width += theme.Padding()
		}
		size.Width += iconSize.Width
	}
	size.Height = fyne.Max(labelSize.Height, iconSize.Height)
	size = size.Add(b.padding())
	return
}

// Layout the components of the button widget
func (b *buttonRenderer) Layout(size fyne.Size) {
	b.LayoutShadow(size, fyne.NewPos(0, 0))
	hasIcon := b.icon != nil
	hasLabel := b.label.Text != ""
	if !hasIcon && !hasLabel {
		// Nothing to layout
		return
	}
	iconSize := fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize())
	labelSize := b.label.MinSize()
	padding := b.padding()
	if hasLabel {
		if hasIcon {
			// Both
			var objects []fyne.CanvasObject
			if b.button.IconPlacement == ButtonIconLeadingText {
				objects = append(objects, b.icon, b.label)
			} else {
				objects = append(objects, b.label, b.icon)
			}
			b.icon.SetMinSize(fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
			min := b.layout.MinSize(objects)
			b.layout.Layout(objects, min)
			pos := alignedPosition(b.button.Alignment, padding, min, size)
			b.label.Move(b.label.Position().Add(pos))
			b.icon.Move(b.icon.Position().Add(pos))
		} else {
			// Label Only
			b.label.Move(alignedPosition(b.button.Alignment, padding, labelSize, size))
			b.label.Resize(labelSize)
		}
	} else {
		// Icon Only
		b.icon.Move(alignedPosition(b.button.Alignment, padding, iconSize, size))
		b.icon.Resize(iconSize)
	}
}

func alignedPosition(align ButtonAlign, padding, objectSize, layoutSize fyne.Size) (pos fyne.Position) {
	pos.Y = (layoutSize.Height - objectSize.Height) / 2
	switch align {
	case ButtonAlignCenter:
		pos.X = (layoutSize.Width - objectSize.Width) / 2
	case ButtonAlignLeading:
		pos.X = padding.Width / 2
	case ButtonAlignTrailing:
		pos.X = layoutSize.Width - objectSize.Width - padding.Width/2
	}
	return
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
	case b.button.hovered, b.button.tapped: // TODO tapped will be different to hovered when we have animation
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
			b.icon.FillMode = canvas.ImageFillContain
			b.SetObjects(append(b.Objects(), b.icon))
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
	Text          string
	Style         ButtonStyle
	Icon          fyne.Resource
	disabledIcon  fyne.Resource
	Alignment     ButtonAlign
	IconPlacement ButtonIconPlacement

	OnTapped   func() `json:"-"`
	HideShadow bool

	hovered, tapped bool
}

// ButtonStyle determines the behaviour and rendering of a button.
type ButtonStyle int

const (
	// DefaultButton is the standard button style
	DefaultButton ButtonStyle = iota
	// PrimaryButton that should be more prominent to the user
	PrimaryButton
)

// ButtonAlign represents the horizontal alignment of a button.
type ButtonAlign int

const (
	// ButtonAlignCenter aligns the icon and the text centrally.
	ButtonAlignCenter ButtonAlign = iota
	// ButtonAlignLeading aligns the icon and the text with the leading edge.
	ButtonAlignLeading
	// ButtonAlignTrailing aligns the icon and the text with the trailing edge.
	ButtonAlignTrailing
)

// ButtonIconPlacement represents the ordering of icon & text within a button.
type ButtonIconPlacement int

const (
	// ButtonIconLeadingText aligns the icon on the leading edge of the text.
	ButtonIconLeadingText ButtonIconPlacement = iota
	// ButtonIconTrailingText aligns the icon on the trailing edge of the text.
	ButtonIconTrailingText
)

// Tapped is called when a pointer tapped event is captured and triggers any tap handler
func (b *Button) Tapped(*fyne.PointEvent) {
	b.tapped = true
	defer func() { // TODO move to a real animation
		time.Sleep(time.Millisecond * buttonTapDuration)
		b.tapped = false
		b.Refresh()
	}()
	b.Refresh()

	if b.OnTapped != nil && !b.Disabled() {
		b.OnTapped()
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
		icon.FillMode = canvas.ImageFillContain
	}

	text := canvas.NewText(b.Text, theme.TextColor())

	objects := []fyne.CanvasObject{
		text,
	}
	shadowLevel := widget.ButtonLevel
	if b.HideShadow {
		shadowLevel = widget.BaseLevel
	}
	if icon != nil {
		objects = append(objects, icon)
	}

	return &buttonRenderer{widget.NewShadowingRenderer(objects, shadowLevel), icon, text, b, layout.NewHBoxLayout()}
}

// SetText allows the button label to be changed
func (b *Button) SetText(text string) {
	b.Text = text

	b.Refresh()
}

// Cursor returns the cursor type of this widget
func (b *Button) Cursor() desktop.Cursor {
	return desktop.DefaultCursor
}

// SetIcon updates the icon on a label - pass nil to hide an icon
func (b *Button) SetIcon(icon fyne.Resource) {
	b.Icon = icon

	if icon != nil {
		b.disabledIcon = theme.NewDisabledResource(icon)
	} else {
		b.disabledIcon = nil
	}

	b.Refresh()
}

// NewButton creates a new button widget with the set label and tap handler
func NewButton(label string, tapped func()) *Button {
	button := &Button{
		Text:     label,
		OnTapped: tapped,
	}

	button.ExtendBaseWidget(button)
	return button
}

// NewButtonWithIcon creates a new button widget with the specified label, themed icon and tap handler
func NewButtonWithIcon(label string, icon fyne.Resource, tapped func()) *Button {
	button := &Button{
		Text:         label,
		Icon:         icon,
		disabledIcon: theme.NewDisabledResource(icon),
		OnTapped:     tapped,
	}

	button.ExtendBaseWidget(button)
	return button
}
