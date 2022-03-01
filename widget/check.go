package widget

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
)

type checkRenderer struct {
	widget.BaseRenderer
	icon           *canvas.Image
	label          *canvas.Text
	focusIndicator *canvas.Circle
	check          *Check
}

// MinSize calculates the minimum size of a check.
// This is based on the contained text, the check icon and a standard amount of padding added.
func (c *checkRenderer) MinSize() fyne.Size {
	pad4 := theme.Padding() * 4
	min := c.label.MinSize().Add(fyne.NewSize(theme.IconInlineSize()+pad4, pad4))
	if c.check.Text != "" {
		min.Add(fyne.NewSize(theme.Padding(), 0))
	}

	return min
}

// Layout the components of the check widget
func (c *checkRenderer) Layout(size fyne.Size) {

	focusIndicatorSize := fyne.NewSize(theme.IconInlineSize()+theme.Padding()*2, theme.IconInlineSize()+theme.Padding()*2)
	c.focusIndicator.Resize(focusIndicatorSize)
	c.focusIndicator.Move(fyne.NewPos(theme.Padding()*0.5, (size.Height-focusIndicatorSize.Height)/2))

	offset := fyne.NewSize(focusIndicatorSize.Width, 0)

	labelSize := size.Subtract(offset)
	c.label.Resize(labelSize)
	c.label.Move(fyne.NewPos(offset.Width+theme.Padding(), 0))

	c.icon.Resize(fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
	c.icon.Move(fyne.NewPos(theme.Padding()*1.5, (size.Height-theme.IconInlineSize())/2))
}

// applyTheme updates this Check to the current theme
func (c *checkRenderer) applyTheme() {
	c.label.Color = theme.ForegroundColor()
	c.label.TextSize = theme.TextSize()
	if c.check.disabled {
		c.label.Color = theme.DisabledColor()
	}
}

func (c *checkRenderer) Refresh() {
	c.check.propertyLock.RLock()
	c.applyTheme()
	c.updateLabel()
	c.updateResource()
	c.updateFocusIndicator()
	c.check.propertyLock.RUnlock()
	canvas.Refresh(c.check.super())
}

func (c *checkRenderer) updateLabel() {
	c.label.Text = c.check.Text
}

func (c *checkRenderer) updateResource() {
	res := theme.CheckButtonIcon()
	if c.check.Checked {
		res = theme.NewPrimaryThemedResource(theme.CheckButtonCheckedIcon())
	}
	if c.check.Disabled() {
		if c.check.Checked {
			res = theme.NewDisabledResource(theme.CheckButtonCheckedIcon())
		} else {
			res = theme.NewDisabledResource(res)
		}
	}
	c.icon.Resource = res
}

func (c *checkRenderer) updateFocusIndicator() {
	if c.check.Disabled() {
		c.focusIndicator.FillColor = theme.BackgroundColor()
	} else if c.check.focused {
		c.focusIndicator.FillColor = theme.FocusColor()
	} else if c.check.hovered {
		c.focusIndicator.FillColor = theme.HoverColor()
	} else {
		c.focusIndicator.FillColor = theme.BackgroundColor()
	}
}

// Check widget has a text label and a checked (or unchecked) icon and triggers an event func when toggled
type Check struct {
	DisableableWidget
	Text    string
	Checked bool

	OnChanged func(bool) `json:"-"`

	focused bool
	hovered bool

	binder basicBinder
}

// Bind connects the specified data source to this Check.
// The current value will be displayed and any changes in the data will cause the widget to update.
// User interactions with this Check will set the value into the data source.
//
// Since: 2.0
func (c *Check) Bind(data binding.Bool) {
	c.binder.SetCallback(c.updateFromData)
	c.binder.Bind(data)

	c.OnChanged = func(_ bool) {
		c.binder.CallWithData(c.writeData)
	}
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

	c.Refresh()
}

// Hide this widget, if it was previously visible
func (c *Check) Hide() {
	if c.focused {
		c.FocusLost()
		impl := c.super()

		if c := fyne.CurrentApp().Driver().CanvasForObject(impl); c != nil {
			c.Focus(nil)
		}
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
	if !c.focused && !fyne.CurrentDevice().IsMobile() {
		impl := c.super()

		if c := fyne.CurrentApp().Driver().CanvasForObject(impl); c != nil {
			c.Focus(impl.(fyne.Focusable))
		}
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
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	icon := canvas.NewImageFromResource(theme.CheckButtonIcon())

	text := canvas.NewText(c.Text, theme.ForegroundColor())
	text.Alignment = fyne.TextAlignLeading

	focusIndicator := canvas.NewCircle(theme.BackgroundColor())
	r := &checkRenderer{
		widget.NewBaseRenderer([]fyne.CanvasObject{focusIndicator, icon, text}),
		icon,
		text,
		focusIndicator,
		c,
	}
	r.applyTheme()
	r.updateLabel()
	r.updateResource()
	r.updateFocusIndicator()
	return r
}

// NewCheck creates a new check widget with the set label and change handler
func NewCheck(label string, changed func(bool)) *Check {
	c := &Check{
		DisableableWidget: DisableableWidget{},
		Text:              label,
		OnChanged:         changed,
	}

	c.ExtendBaseWidget(c)
	return c
}

// NewCheckWithData returns a check widget connected with the specified data source.
//
// Since: 2.0
func NewCheckWithData(label string, data binding.Bool) *Check {
	check := NewCheck(label, nil)
	check.Bind(data)

	return check
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

// Unbind disconnects any configured data source from this Check.
// The current value will remain at the last value of the data source.
//
// Since: 2.0
func (c *Check) Unbind() {
	c.OnChanged = nil
	c.binder.Unbind()
}

func (c *Check) updateFromData(data binding.DataItem) {
	if data == nil {
		return
	}
	boolSource, ok := data.(binding.Bool)
	if !ok {
		return
	}
	val, err := boolSource.Get()
	if err != nil {
		fyne.LogError("Error getting current data value", err)
		return
	}
	c.SetChecked(val) // if val != c.Checked, this will call updateFromData again, but only once
}

func (c *Check) writeData(data binding.DataItem) {
	if data == nil {
		return
	}
	boolTarget, ok := data.(binding.Bool)
	if !ok {
		return
	}
	currentValue, err := boolTarget.Get()
	if err != nil {
		return
	}
	if currentValue != c.Checked {
		err := boolTarget.Set(c.Checked)
		if err != nil {
			fyne.LogError(fmt.Sprintf("Failed to set binding value to %t", c.Checked), err)
		}
	}
}
