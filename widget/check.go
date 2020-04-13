package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/binding"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/theme"
)

type checkRenderer struct {
	baseRenderer
	icon           *canvas.Image
	label          *canvas.Text
	focusIndicator *canvas.Circle
	check          *Check
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

// applyTheme updates this Check to the current theme
func (c *checkRenderer) applyTheme() {
	c.label.Color = theme.TextColor()
	c.label.TextSize = theme.TextSize()
	if c.check.Disabled() {
		c.label.Color = theme.DisabledTextColor()
	}
}

func (c *checkRenderer) Refresh() {
	c.applyTheme()
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
	} else if c.check.focused {
		c.focusIndicator.FillColor = theme.FocusColor()
	} else if c.check.hovered {
		c.focusIndicator.FillColor = theme.HoverColor()
	} else {
		c.focusIndicator.FillColor = theme.BackgroundColor()
	}

	canvas.Refresh(c.check.super())
}

// Check widget has a text label and a checked (or unchecked) icon and triggers an event func when toggled
type Check struct {
	DisableableWidget
	Text    string
	Checked bool

	OnChanged func(bool) `json:"-"`

	focused bool
	hovered bool

	changeBind  *binding.Bool
	textBind    *binding.String
	checkNotify *binding.NotifyFunction
	textNotify  *binding.NotifyFunction
}

// SetChecked sets the the checked state and refreshes widget
func (c *Check) SetChecked(checked bool) {
	if checked == c.Checked {
		return
	}

	c.Checked = checked

	if c.changeBind != nil {
		c.changeBind.Set(c.Checked)
	}

	if c.OnChanged != nil {
		c.OnChanged(c.Checked)
	}

	c.Refresh()
}

// SetText allows the check label to be changed
func (c *Check) SetText(text string) {
	if c.Text != text {
		c.Text = text
		c.Refresh()
	}
}

// Hide this widget, if it was previously visible
func (c *Check) Hide() {
	if c.focused {
		c.FocusLost()
		fyne.CurrentApp().Driver().CanvasForObject(c).Focus(nil)
	}

	c.BaseWidget.Hide()
}

// MouseIn is called when a desktop pointer enters the widget
func (c *Check) MouseIn(*desktop.MouseEvent) {
	if c.Disabled() {
		return
	}
	c.hovered = true
	c.Refresh()
}

// MouseOut is called when a desktop pointer exits the widget
func (c *Check) MouseOut() {
	c.hovered = false
	c.Refresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (c *Check) MouseMoved(*desktop.MouseEvent) {
}

// Tapped is called when a pointer tapped event is captured and triggers any change handler
func (c *Check) Tapped(*fyne.PointEvent) {
	if !c.focused {
		c.FocusGained()
	}
	if !c.Disabled() {
		c.SetChecked(!c.Checked)
	}
}

// MinSize returns the size that this widget should not shrink below
func (c *Check) MinSize() fyne.Size {
	c.ExtendBaseWidget(c)
	return c.BaseWidget.MinSize()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (c *Check) CreateRenderer() fyne.WidgetRenderer {
	c.ExtendBaseWidget(c)
	icon := canvas.NewImageFromResource(theme.CheckButtonIcon())

	text := canvas.NewText(c.Text, theme.TextColor())
	text.Alignment = fyne.TextAlignLeading

	focusIndicator := canvas.NewCircle(theme.BackgroundColor())
	return &checkRenderer{baseRenderer{[]fyne.CanvasObject{focusIndicator, icon, text}}, icon, text, focusIndicator, c}
}

// BindChecked binds the Check's OnChanged to the given data binding.
// Returns the Check for chaining.
func (c *Check) BindChecked(data *binding.Bool) *Check {
	c.changeBind = data
	c.checkNotify = data.AddBoolListener(c.SetChecked)
	return c
}

// UnbindChecked unbinds the Check's OnChanged from the data binding (if any).
// Returns the Check for chaining.
func (c *Check) UnbindChecked() *Check {
	c.changeBind.DeleteListener(c.checkNotify)
	c.changeBind = nil
	c.checkNotify = nil
	return c
}

// BindText binds the Check's Text to the given data binding.
// Returns the Check for chaining.
func (c *Check) BindText(data *binding.String) *Check {
	c.textBind = data
	c.textNotify = data.AddStringListener(c.SetText)
	return c
}

// UnbindText unbinds the Check's Text from the data binding (if any).
// Returns the Check for chaining.
func (c *Check) UnbindText() *Check {
	c.textBind.DeleteListener(c.textNotify)
	c.textBind = nil
	c.textNotify = nil
	return c
}

// NewCheck creates a new check widget with the set label and change handler
func NewCheck(label string, changed func(bool)) *Check {
	c := &Check{
		DisableableWidget: DisableableWidget{},
		Text:              label,
		Checked:           false,
		OnChanged:         changed,
		focused:           false,
		hovered:           false,
	}

	c.ExtendBaseWidget(c)
	return c
}

// FocusGained is called when the Check has been given focus.
func (c *Check) FocusGained() {
	if c.Disabled() {
		return
	}
	c.focused = true

	c.Refresh()
}

// FocusLost is called when the Check has had focus removed.
func (c *Check) FocusLost() {
	c.focused = false

	c.Refresh()
}

// Focused returns whether or not this Check has focus.
// Deprecated: this method will be removed as it is no longer required, widgets do not expose their focus state.
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
