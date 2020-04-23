package widget

import (
	"image/color"
	"log"

	"fyne.io/fyne"
	"fyne.io/fyne/binding"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/theme"
)

// ButtonStyle determines the behaviour and rendering of a button.
type ButtonStyle int

const (
	// DefaultButton is the standard button style
	DefaultButton ButtonStyle = iota
	// PrimaryButton that should be more prominent to the user
	PrimaryButton
)

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

	DisabledProperty binding.Bool
	HiddenProperty   binding.Bool
	HoveredProperty  binding.Bool
	IconProperty     binding.Resource
	ShadowedProperty binding.Bool
	StyleProperty    binding.Int
	TextProperty     binding.String
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
	log.Println("Button.CreateRenderer")
	// Ensure all Properties are set.
	if b.DisabledProperty == nil {
		b.DisabledProperty = binding.EmptyBool()
	}
	if b.HiddenProperty == nil {
		b.HiddenProperty = binding.EmptyBool()
	}
	if b.HoveredProperty == nil {
		b.HoveredProperty = binding.EmptyBool()
	}
	if b.IconProperty == nil {
		b.IconProperty = binding.EmptyResource()
	}
	if b.ShadowedProperty == nil {
		b.ShadowedProperty = binding.EmptyBool()
	}
	if b.StyleProperty == nil {
		b.StyleProperty = binding.EmptyInt()
	}
	if b.TextProperty == nil {
		b.TextProperty = binding.EmptyString()
	}
	image := &canvas.Image{}
	imageDisabled := &canvas.Image{}
	label := &canvas.Text{
		Alignment: fyne.TextAlignCenter,
	}
	objects := []fyne.CanvasObject{
		image,
		imageDisabled,
		label,
	}
	r := &buttonRenderer{
		ShadowingRenderer: widget.NewShadowingRenderer(objects, widget.ButtonLevel), //Temporary - shadowedChan will trigger update
		button:            b,
		image:             image,
		imageDisabled:     imageDisabled,
		label:             label,
		done:              make(chan bool),
	}
	// Create goroutine to listen to each property, respond to changes, and trigger relayout
	go func() {
		disabledChan := b.DisabledProperty.Listen()
		hoveredChan := b.HoveredProperty.Listen()
		hiddenChan := b.HoveredProperty.Listen()
		iconChan := b.IconProperty.Listen()
		shadowedChan := b.ShadowedProperty.Listen()
		styleChan := b.StyleProperty.Listen()
		textChan := b.TextProperty.Listen()
		for {
			relayout := false
			select {
			case d := <-disabledChan:
				log.Println("disabled:", d)
				if d {
					r.imageDisabled.Show()
					r.image.Hide()
				} else {
					r.image.Show()
					r.imageDisabled.Hide()
				}
				relayout = true
			case h := <-hiddenChan:
				log.Println("hidden:", h)
				if h {
					continue
				}
			case h := <-hoveredChan:
				log.Println("hovered:", h)
			case i := <-iconChan:
				log.Println("icon:", i)
				r.image.Resource = i
				if i == nil {
					r.imageDisabled.Resource = nil
				} else {
					r.imageDisabled.Resource = theme.NewDisabledResource(i)
				}
				relayout = true
			case s := <-shadowedChan:
				log.Println("shadowed:", s)
				shadowLevel := widget.ButtonLevel
				if r.button.HideShadow {
					shadowLevel = widget.BaseLevel
				}
				r.ShadowingRenderer = widget.NewShadowingRenderer(objects, shadowLevel)
			case s := <-styleChan:
				log.Println("style:", s)
			case t := <-textChan:
				log.Println("text:", t)
				r.label.Text = t
				relayout = true
			case <-r.done:
				return
			}
			if relayout {
				r.Layout(b.Size())
			}
			canvas.Refresh(b)
		}
	}()

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
	if b.Disabled() {
		return
	}
	b.hovered = true
	b.Refresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (b *Button) MouseMoved(*desktop.MouseEvent) {
}

// MouseOut is called when a desktop pointer exits the widget
func (b *Button) MouseOut() {
	if b.Disabled() {
		return
	}
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
	if b.OnTapped != nil {
		b.OnTapped()
	}
}

type buttonRenderer struct {
	*widget.ShadowingRenderer
	button *Button

	image         *canvas.Image
	imageDisabled *canvas.Image
	label         *canvas.Text

	done chan bool
}

func (r *buttonRenderer) BackgroundColor() color.Color {
	switch {
	case r.button.DisabledProperty.Get():
		return theme.DisabledButtonColor()
	case r.button.StyleProperty.Get() == int(PrimaryButton):
		return theme.PrimaryColor()
	case r.button.HoveredProperty.Get():
		return theme.HoverColor()
	default:
		return theme.ButtonColor()
	}
}

// Layout the components of the button widget
func (r *buttonRenderer) Layout(size fyne.Size) {
	r.LayoutShadow(size, fyne.NewPos(0, 0))
	hasIcon := r.button.IconProperty.Get() != nil
	if r.button.TextProperty.Get() != "" {
		padding := r.padding()
		innerSize := size.Subtract(padding)
		innerOffset := fyne.NewPos(padding.Width/2, padding.Height/2)

		labelSize := r.label.MinSize()
		contentWidth := labelSize.Width

		if hasIcon {
			contentWidth += theme.Padding() + theme.IconInlineSize()
			imageOffset := fyne.NewPos((innerSize.Width-contentWidth)/2, (innerSize.Height-theme.IconInlineSize())/2)
			imageSize := fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize())
			imagePos := innerOffset.Add(imageOffset)
			r.image.Resize(imageSize)
			r.image.Move(imagePos)
			r.imageDisabled.Resize(imageSize)
			r.imageDisabled.Move(imagePos)
		}
		labelOffset := fyne.NewPos((innerSize.Width+contentWidth)/2-labelSize.Width, (innerSize.Height-labelSize.Height)/2)
		r.label.Resize(labelSize)
		r.label.Move(innerOffset.Add(labelOffset))
	} else if hasIcon {
		imageSize := fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize())
		imagePos := fyne.NewPos((size.Width-theme.IconInlineSize())/2, (size.Height-theme.IconInlineSize())/2)
		r.image.Resize(imageSize)
		r.image.Move(imagePos)
		r.imageDisabled.Resize(imageSize)
		r.imageDisabled.Move(imagePos)
	}
}

// MinSize calculates the minimum size of a button.
// This is based on the contained text, any icon that is set and a standard
// amount of padding added.
func (r *buttonRenderer) MinSize() fyne.Size {
	labelSize := r.label.MinSize()
	contentHeight := fyne.Max(labelSize.Height, theme.IconInlineSize())
	contentWidth := 0
	hasIcon := r.button.IconProperty.Get() != nil
	if hasIcon {
		contentWidth += theme.IconInlineSize()
	}
	if r.button.TextProperty.Get() != "" {
		if hasIcon {
			contentWidth += theme.Padding()
		}
		contentWidth += labelSize.Width
	}
	return fyne.NewSize(contentWidth, contentHeight).Add(r.padding())
}

func (r *buttonRenderer) Refresh() {
	r.applyTheme()
	// Push button state through channel bindings
	// Listeners only triggered if value is different.
	r.button.DisabledProperty.Set(r.button.Disabled())
	r.button.HiddenProperty.Set(r.button.Hidden)
	r.button.HoveredProperty.Set(r.button.hovered)
	r.button.IconProperty.Set(r.button.Icon)
	r.button.ShadowedProperty.Set(r.button.HideShadow)
	r.button.StyleProperty.Set(int(r.button.Style))
	r.button.TextProperty.Set(r.button.Text)
}

// applyTheme updates this button to match the current theme
func (r *buttonRenderer) applyTheme() {
	r.label.TextSize = theme.TextSize()
	r.label.Color = theme.TextColor()
	if r.button.DisabledProperty.Get() {
		r.label.Color = theme.DisabledTextColor()
	}
}

func (r *buttonRenderer) padding() fyne.Size {
	if r.button.TextProperty.Get() == "" {
		return fyne.NewSize(theme.Padding()*2, theme.Padding()*2)
	}
	return fyne.NewSize(theme.Padding()*4, theme.Padding()*2)
}
