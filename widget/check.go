package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/theme"
)

type checkRenderer struct {
	icon  *canvas.Image
	label *canvas.Text

	focusIndicator *canvas.Circle

	objects []fyne.CanvasObject
	check   *Check
}

// MinSize calculates the minimum size of a check.
// This is based on the contained text, the check icon and a standard amount of padding added.
func (c *checkRenderer) MinSize() fyne.Size {
	min := c.label.MinSize().Add(fyne.NewSize(theme.Padding()*4, theme.Padding()*2))
	min = min.Add(fyne.NewSize(theme.IconInlineSize()+theme.Padding(), 0))

	return min
}

// Layout the components of the check widget
func (c *checkRenderer) Layout(size fyne.Size) {

	focusIndicatorSize := fyne.NewSize(theme.IconInlineSize()+theme.Padding()*2, theme.IconInlineSize()+theme.Padding()*2)
	c.focusIndicator.Resize(focusIndicatorSize)
	c.focusIndicator.Move(fyne.NewPos(0, (size.Height-focusIndicatorSize.Height)/2))

	offset := fyne.NewSize(focusIndicatorSize.Width, 0)

	labelSize := size.Subtract(offset)
	c.label.Resize(labelSize)
	c.label.Move(fyne.NewPos(offset.Width+theme.Padding(), 0))

	c.icon.Resize(fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
	c.icon.Move(fyne.NewPos(theme.Padding(), (size.Height-theme.IconInlineSize())/2))
}

// ApplyTheme is called when the Check may need to update its look
func (c *checkRenderer) ApplyTheme() {
	c.label.Color = theme.TextColor()
	c.label.TextSize = theme.TextSize()
	if c.check.Disabled() {
		c.label.Color = theme.DisabledTextColor()
	}

	c.Refresh()
}

func (c *checkRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (c *checkRenderer) Refresh() {
	c.label.Text = c.check.Text

	res := theme.CheckButtonIcon()
	if c.check.Checked {
		res = theme.CheckButtonCheckedIcon()
	}
	if c.check.Disabled() {
		res = theme.NewDisabledResource(res)
	}

	c.icon.Resource = res

	if c.check.Disabled() {
		c.focusIndicator.FillColor = theme.BackgroundColor()
	} else if c.check.Focused() {
		c.focusIndicator.FillColor = theme.FocusColor()
	} else if c.check.hovered {
		c.focusIndicator.FillColor = theme.HoverColor()
	} else {
		c.focusIndicator.FillColor = theme.BackgroundColor()
	}

	canvas.Refresh(c.check)
}

func (c *checkRenderer) Objects() []fyne.CanvasObject {
	return c.objects
}

func (c *checkRenderer) Destroy() {
}

// Check widget has a text label and a checked (or unchecked) icon and triggers an event func when toggled
type Check struct {
	baseWidget
	Text    string
	Checked bool

	OnChanged func(bool) `json:"-"`

	focused bool
	hovered bool
}

// SetChecked sets the the checked state and refreshes widget
func (c *Check) SetChecked(checked bool) {
	if checked == c.Checked {
		return
	}

	c.Checked = checked

	if c.OnChanged != nil {
		c.OnChanged(c.Checked)
	}

	Refresh(c)
}

// Resize sets a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (c *Check) Resize(size fyne.Size) {
	c.resize(size, c)
}

// Move the widget to a new position, relative to its parent.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (c *Check) Move(pos fyne.Position) {
	c.move(pos, c)
}

// MinSize returns the smallest size this widget can shrink to
func (c *Check) MinSize() fyne.Size {
	return c.minSize(c)
}

// Show this widget, if it was previously hidden
func (c *Check) Show() {
	c.show(c)
}

// Hide this widget, if it was previously visible
func (c *Check) Hide() {
	if c.Focused() {
		c.FocusLost()
		fyne.CurrentApp().Driver().CanvasForObject(c).Focus(nil)
	}
	c.hide(c)
}

// Enable this widget, if it was previously disabled
func (c *Check) Enable() {
	c.enable(c)
	Renderer(c).ApplyTheme()
}

// Disable this widget, if it was previously enabled
func (c *Check) Disable() {
	c.disable(c)
	Renderer(c).ApplyTheme()
}

// Disabled returns true if the widget is disabled
func (c *Check) Disabled() bool {
	return c.disabled
}

// MouseIn is called when a desktop pointer enters the widget
func (c *Check) MouseIn(*desktop.MouseEvent) {
	if c.Disabled() {
		return
	}
	c.hovered = true
	Refresh(c)
}

// MouseOut is called when a desktop pointer exits the widget
func (c *Check) MouseOut() {
	c.hovered = false
	Refresh(c)
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (c *Check) MouseMoved(*desktop.MouseEvent) {
}

// Tapped is called when a pointer tapped event is captured and triggers any change handler
func (c *Check) Tapped(*fyne.PointEvent) {
	if !c.Focused() {
		c.FocusGained()
	}
	if !c.Disabled() {
		c.SetChecked(!c.Checked)
	}
}

// TappedSecondary is called when a secondary pointer tapped event is captured
func (c *Check) TappedSecondary(*fyne.PointEvent) {
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (c *Check) CreateRenderer() fyne.WidgetRenderer {
	icon := canvas.NewImageFromResource(theme.CheckButtonIcon())

	text := canvas.NewText(c.Text, theme.TextColor())
	text.Alignment = fyne.TextAlignLeading

	focusIndicator := canvas.NewCircle(theme.BackgroundColor())
	return &checkRenderer{icon, text, focusIndicator, []fyne.CanvasObject{focusIndicator, icon, text}, c}
}

// NewCheck creates a new check widget with the set label and change handler
func NewCheck(label string, changed func(bool)) *Check {
	c := &Check{
		baseWidget{},
		label,
		false,
		changed,
		false,
		false,
	}

	Renderer(c).Layout(c.MinSize())
	return c
}

// FocusGained is called when the Check has been given focus.
func (c *Check) FocusGained() {
	if c.Disabled() {
		return
	}
	c.focused = true

	Refresh(c)
}

// FocusLost is called when the Check has had focus removed.
func (c *Check) FocusLost() {
	c.focused = false

	Refresh(c)
}

// Focused returns whether or not this Check has focus.
func (c *Check) Focused() bool {
	if c.Disabled() {
		return false
	}
	return c.focused
}

// TypedRune receives text input events when the Check is focused.
func (c *Check) TypedRune(r rune) {
	if c.Disabled() {
		return
	}
	if r == ' ' {
		c.SetChecked(!c.Checked)
	}
}

// TypedKey receives key input events when the Check is focused.
func (c *Check) TypedKey(key *fyne.KeyEvent) {}
