package widget

import (
	"image/color"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
)

// ButtonAlign represents the horizontal alignment of a button.
type ButtonAlign int

// ButtonIconPlacement represents the ordering of icon & text within a button.
type ButtonIconPlacement int

// ButtonImportance represents how prominent the button should appear
//
// Since: 1.4
type ButtonImportance int

// ButtonStyle determines the behaviour and rendering of a button.
type ButtonStyle int

const (
	// ButtonAlignCenter aligns the icon and the text centrally.
	ButtonAlignCenter ButtonAlign = iota
	// ButtonAlignLeading aligns the icon and the text with the leading edge.
	ButtonAlignLeading
	// ButtonAlignTrailing aligns the icon and the text with the trailing edge.
	ButtonAlignTrailing
)

const (
	// ButtonIconLeadingText aligns the icon on the leading edge of the text.
	ButtonIconLeadingText ButtonIconPlacement = iota
	// ButtonIconTrailingText aligns the icon on the trailing edge of the text.
	ButtonIconTrailingText
)

const (
	// MediumImportance applies a standard appearance.
	MediumImportance ButtonImportance = iota
	// HighImportance applies a prominent appearance.
	HighImportance
	// LowImportance applies a subtle appearance.
	LowImportance
)

const (
	// DefaultButton is the standard button style.
	// Deprecated: use Importance = MediumImportance instead.
	DefaultButton ButtonStyle = iota
	// PrimaryButton that should be more prominent to the user.
	// Deprecated: use Importance = HighImportance instead.
	PrimaryButton

	buttonTapDuration = 250
)

// Button widget has a text label and triggers an event func when clicked
type Button struct {
	DisableableWidget
	Text string
	// Deprecated, use Importance instead, where HighImportance matches PrimaryButton
	Style ButtonStyle
	Icon  fyne.Resource
	// Specify how prominent the button should be, High will highlight the button and Low will remove some decoration.
	//
	// Since: 1.4
	Importance    ButtonImportance
	Alignment     ButtonAlign
	IconPlacement ButtonIconPlacement

	OnTapped func() `json:"-"`
	// Deprecated: use Importance = LowImportance instead of HideShadow = true.
	HideShadow bool

	hovered, tapped bool
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
		Text:     label,
		Icon:     icon,
		OnTapped: tapped,
	}

	button.ExtendBaseWidget(button)
	return button
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (b *Button) CreateRenderer() fyne.WidgetRenderer {
	b.ExtendBaseWidget(b)
	text := canvas.NewText(b.Text, theme.TextColor())
	text.TextStyle.Bold = true

	bg := canvas.NewRectangle(color.Transparent)
	objects := []fyne.CanvasObject{
		bg,
		text,
	}
	shadowLevel := widget.ButtonLevel
	if b.HideShadow || b.Importance == LowImportance {
		shadowLevel = widget.BaseLevel
	}
	r := &buttonRenderer{
		ShadowingRenderer: widget.NewShadowingRenderer(objects, shadowLevel),
		bg:                bg,
		button:            b,
		label:             text,
		layout:            layout.NewHBoxLayout(),
	}
	bg.FillColor = r.buttonColor()
	r.updateIconAndText()
	r.applyTheme()
	return r
}

// Cursor returns the cursor type of this widget
func (b *Button) Cursor() desktop.Cursor {
	return desktop.DefaultCursor
}

// MinSize returns the size that this widget should not shrink below
func (b *Button) MinSize() fyne.Size {
	b.ExtendBaseWidget(b)
	return b.BaseWidget.MinSize()
}

// MouseIn is called when a desktop pointer enters the widget
func (b *Button) MouseIn(*desktop.MouseEvent) {
	b.hovered = true
	b.Refresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (b *Button) MouseMoved(*desktop.MouseEvent) {
}

// MouseOut is called when a desktop pointer exits the widget
func (b *Button) MouseOut() {
	b.hovered = false
	b.Refresh()
}

// SetIcon updates the icon on a label - pass nil to hide an icon
func (b *Button) SetIcon(icon fyne.Resource) {
	b.Icon = icon

	b.Refresh()
}

// SetText allows the button label to be changed
func (b *Button) SetText(text string) {
	b.Text = text

	b.Refresh()
}

// Tapped is called when a pointer tapped event is captured and triggers any tap handler
func (b *Button) Tapped(*fyne.PointEvent) {
	if b.Disabled() {
		return
	}

	b.tapped = true
	defer func() { // TODO move to a real animation
		time.Sleep(time.Millisecond * buttonTapDuration)
		b.tapped = false
		b.Refresh()
	}()
	b.Refresh()

	if b.OnTapped != nil {
		b.OnTapped()
	}
}

type buttonRenderer struct {
	*widget.ShadowingRenderer

	icon   *canvas.Image
	label  *canvas.Text
	bg     *canvas.Rectangle
	button *Button
	layout fyne.Layout
}

func (r *buttonRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

// Layout the components of the button widget
func (r *buttonRenderer) Layout(size fyne.Size) {
	var inset fyne.Position
	bgSize := size
	if !r.button.HideShadow || r.button.Importance != LowImportance {
		inset = fyne.NewPos(theme.Padding()/2, theme.Padding()/2)
		bgSize = size.Subtract(fyne.NewSize(theme.Padding(), theme.Padding()))
	}
	r.LayoutShadow(bgSize, inset)

	r.bg.Move(inset)
	r.bg.Resize(bgSize)

	hasIcon := r.icon != nil
	hasLabel := r.label.Text != ""
	if !hasIcon && !hasLabel {
		// Nothing to layout
		return
	}
	iconSize := fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize())
	labelSize := r.label.MinSize()
	padding := r.padding()
	if hasLabel {
		if hasIcon {
			// Both
			var objects []fyne.CanvasObject
			if r.button.IconPlacement == ButtonIconLeadingText {
				objects = append(objects, r.icon, r.label)
			} else {
				objects = append(objects, r.label, r.icon)
			}
			r.icon.SetMinSize(iconSize)
			min := r.layout.MinSize(objects)
			r.layout.Layout(objects, min)
			pos := alignedPosition(r.button.Alignment, padding, min, size)
			r.label.Move(r.label.Position().Add(pos))
			r.icon.Move(r.icon.Position().Add(pos))
		} else {
			// Label Only
			r.label.Move(alignedPosition(r.button.Alignment, padding, labelSize, size))
			r.label.Resize(labelSize)
		}
	} else {
		// Icon Only
		r.icon.Move(alignedPosition(r.button.Alignment, padding, iconSize, size))
		r.icon.Resize(iconSize)
	}
}

// MinSize calculates the minimum size of a button.
// This is based on the contained text, any icon that is set and a standard
// amount of padding added.
func (r *buttonRenderer) MinSize() (size fyne.Size) {
	hasIcon := r.icon != nil
	hasLabel := r.label.Text != ""
	iconSize := fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize())
	labelSize := r.label.MinSize()
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
	size = size.Add(r.padding())
	return
}

func (r *buttonRenderer) Refresh() {
	r.label.Text = r.button.Text
	r.bg.Refresh()
	r.updateIconAndText()
	r.applyTheme()
	r.Layout(r.button.Size())
	canvas.Refresh(r.button.super())
}

// applyTheme updates this button to match the current theme
func (r *buttonRenderer) applyTheme() {
	r.bg.FillColor = r.buttonColor()
	r.label.TextSize = theme.TextSize()
	r.label.Color = theme.TextColor()
	switch {
	case r.button.disabled:
		r.label.Color = theme.DisabledTextColor()
	case r.button.Style == PrimaryButton || r.button.Importance == HighImportance:
		r.label.Color = theme.BackgroundColor()
	}
	if r.icon != nil && r.icon.Resource != nil {
		switch res := r.icon.Resource.(type) {
		case *theme.ThemedResource:
			if r.button.Style == PrimaryButton || r.button.Importance == HighImportance {
				r.icon.Resource = theme.NewInvertedThemedResource(res)
				r.icon.Refresh()
			}
		case *theme.InvertedThemedResource:
			if r.button.Style != PrimaryButton || r.button.Importance != HighImportance {
				r.icon.Resource = res.Original()
				r.icon.Refresh()
			}
		}
	}
}

func (r *buttonRenderer) buttonColor() color.Color {
	switch {
	case r.button.Disabled():
		return theme.DisabledButtonColor()
	case r.button.Style == PrimaryButton, r.button.Importance == HighImportance:
		return theme.PrimaryColor()
	case r.button.hovered, r.button.tapped: // TODO tapped will be different to hovered when we have animation
		return theme.HoverColor()
	default:
		return theme.ButtonColor()
	}
}

func (r *buttonRenderer) padding() fyne.Size {
	if r.button.HideShadow || r.button.Importance == LowImportance {
		return fyne.NewSize(theme.Padding()*2, theme.Padding()*2)
	}
	if r.button.Text == "" {
		return fyne.NewSize(theme.Padding()*4, theme.Padding()*4)
	}
	return fyne.NewSize(theme.Padding()*6, theme.Padding()*4)
}

func (r *buttonRenderer) updateIconAndText() {
	if r.button.Icon != nil && r.button.Visible() {
		if r.icon == nil {
			r.icon = canvas.NewImageFromResource(r.button.Icon)
			r.icon.FillMode = canvas.ImageFillContain
			r.SetObjects([]fyne.CanvasObject{r.bg, r.label, r.icon})
		}
		if r.button.Disabled() {
			r.icon.Resource = theme.NewDisabledResource(r.button.Icon)
		} else {
			r.icon.Resource = r.button.Icon
		}
		r.icon.Refresh()
		r.icon.Show()
	} else if r.icon != nil {
		r.icon.Hide()
	}
	if r.button.Text == "" {
		r.label.Hide()
	} else {
		r.label.Show()
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
