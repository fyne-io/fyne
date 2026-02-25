package widget

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
)

// Check widget has a text label and a checked (or unchecked) icon and triggers an event func when toggled
type Check struct {
	DisableableWidget
	Text    string
	Checked bool

	// Partial check is when there is an indeterminate state (usually meaning that child items are some-what checked).
	// Turning this on will override the checked state and show a dash icon (neither checked nor unchecked).
	// The user interaction cannot turn this on, tapping a partial check state will set `Checked` to true.
	//
	// Since: 2.6
	Partial bool

	OnChanged func(bool) `json:"-"`

	focused bool
	hovered bool

	binder basicBinder

	minSize fyne.Size // cached for hover/tap position calculations
}

// NewCheck creates a new check widget with the set label and change handler
func NewCheck(label string, changed func(bool)) *Check {
	c := &Check{
		Text:      label,
		OnChanged: changed,
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

// SetChecked sets the checked state and refreshes widget
// If the `Partial` state is set this will be turned off to respect the `checked` bool passed in here.
func (c *Check) SetChecked(checked bool) {
	if checked == c.Checked && !c.Partial {
		return
	}

	c.Partial = false
	c.Checked = checked
	onChanged := c.OnChanged

	if onChanged != nil {
		onChanged(checked)
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
func (c *Check) MouseIn(me *desktop.MouseEvent) {
	c.MouseMoved(me)
}

// MouseOut is called when a desktop pointer exits the widget
func (c *Check) MouseOut() {
	if c.hovered {
		c.hovered = false
		c.Refresh()
	}
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (c *Check) MouseMoved(me *desktop.MouseEvent) {
	if c.Disabled() {
		return
	}

	oldHovered := c.hovered

	// only hovered if cached minSize has not been initialized (test code)
	// or the pointer is within the "active" area of the widget (its minSize)
	c.hovered = c.minSize.IsZero() ||
		(me.Position.X <= c.minSize.Width && me.Position.Y <= c.minSize.Height)

	if oldHovered != c.hovered {
		c.Refresh()
	}
}

// Tapped is called when a pointer tapped event is captured and triggers any change handler
func (c *Check) Tapped(pe *fyne.PointEvent) {
	if c.Disabled() {
		return
	}

	minHeight := c.minSize.Height
	minY := (c.Size().Height - minHeight) / 2
	if !c.minSize.IsZero() &&
		(pe.Position.X > c.minSize.Width || pe.Position.Y < minY || pe.Position.Y > minY+minHeight) {
		// tapped outside the active area of the widget
		return
	}

	if !c.focused {
		focusIfNotMobile(c.super())
	}
	c.SetChecked(!c.Checked)
}

// MinSize returns the size that this widget should not shrink below
func (c *Check) MinSize() fyne.Size {
	c.ExtendBaseWidget(c)
	c.minSize = c.BaseWidget.MinSize()
	return c.minSize
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (c *Check) CreateRenderer() fyne.WidgetRenderer {
	th := c.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	c.ExtendBaseWidget(c)
	bg := canvas.NewImageFromResource(th.Icon(theme.IconNameCheckButtonFill))
	icon := canvas.NewImageFromResource(th.Icon(theme.IconNameCheckButton))

	text := canvas.NewText(c.Text, th.Color(theme.ColorNameForeground, v))
	text.Alignment = fyne.TextAlignLeading

	focusIndicator := canvas.NewCircle(th.Color(theme.ColorNameBackground, v))
	r := &checkRenderer{
		widget.NewBaseRenderer([]fyne.CanvasObject{focusIndicator, bg, icon, text}),
		bg,
		icon,
		text,
		focusIndicator,
		c,
	}
	r.applyTheme(th, v)
	r.updateLabel()
	r.updateResource(th)
	r.updateFocusIndicator(th, v)
	return r
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

// SetText sets the text of the Check
//
// Since: 2.4
func (c *Check) SetText(text string) {
	c.Text = text
	c.Refresh()
}

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

type checkRenderer struct {
	widget.BaseRenderer
	bg, icon       *canvas.Image
	label          *canvas.Text
	focusIndicator *canvas.Circle
	check          *Check
}

// MinSize calculates the minimum size of a check.
// This is based on the contained text, the check icon and a standard amount of padding added.
func (c *checkRenderer) MinSize() fyne.Size {
	th := c.check.Theme()

	pad4 := th.Size(theme.SizeNameInnerPadding) * 2
	min := c.label.MinSize().Add(fyne.NewSize(th.Size(theme.SizeNameInlineIcon)+pad4, pad4))

	if c.check.Text != "" {
		min.Add(fyne.NewSize(th.Size(theme.SizeNamePadding), 0))
	}

	return min
}

// Layout the components of the check widget
func (c *checkRenderer) Layout(size fyne.Size) {
	th := c.check.Theme()
	innerPadding := th.Size(theme.SizeNameInnerPadding)
	borderSize := th.Size(theme.SizeNameInputBorder)
	iconInlineSize := th.Size(theme.SizeNameInlineIcon)

	focusIndicatorSize := fyne.NewSquareSize(iconInlineSize + innerPadding)
	c.focusIndicator.Resize(focusIndicatorSize)
	c.focusIndicator.Move(fyne.NewPos(borderSize, (size.Height-focusIndicatorSize.Height)/2))

	xOff := focusIndicatorSize.Width + borderSize*2
	labelSize := size.SubtractWidthHeight(xOff, 0)
	c.label.Resize(labelSize)
	c.label.Move(fyne.NewPos(xOff, 0))

	iconPos := fyne.NewPos(innerPadding/2+borderSize, (size.Height-iconInlineSize)/2)
	iconSize := fyne.NewSquareSize(iconInlineSize)
	c.bg.Move(iconPos)
	c.bg.Resize(iconSize)
	c.icon.Move(iconPos)
	c.icon.Resize(iconSize)
}

// applyTheme updates this Check to the current theme
func (c *checkRenderer) applyTheme(th fyne.Theme, v fyne.ThemeVariant) {
	c.label.Color = th.Color(theme.ColorNameForeground, v)
	c.label.TextSize = th.Size(theme.SizeNameText)
	if c.check.Disabled() {
		c.label.Color = th.Color(theme.ColorNameDisabled, v)
	}
}

func (c *checkRenderer) Refresh() {
	th := c.check.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	c.applyTheme(th, v)
	c.updateLabel()
	c.updateResource(th)
	c.updateFocusIndicator(th, v)
	canvas.Refresh(c.check.super())
}

// must be called while holding c.check.propertyLock for reading
func (c *checkRenderer) updateLabel() {
	c.label.Text = c.check.Text
}

// must be called while holding c.check.propertyLock for reading
func (c *checkRenderer) updateResource(th fyne.Theme) {
	res := theme.NewThemedResource(th.Icon(theme.IconNameCheckButton))
	res.ColorName = theme.ColorNameInputBorder
	bgRes := theme.NewThemedResource(th.Icon(theme.IconNameCheckButtonFill))
	bgRes.ColorName = theme.ColorNameInputBackground

	if c.check.Partial {
		res = theme.NewThemedResource(th.Icon(theme.IconNameCheckButtonPartial))
		res.ColorName = theme.ColorNamePrimary
		bgRes.ColorName = theme.ColorNameBackground
	} else if c.check.Checked {
		res = theme.NewThemedResource(th.Icon(theme.IconNameCheckButtonChecked))
		res.ColorName = theme.ColorNamePrimary
		bgRes.ColorName = theme.ColorNameBackground
	}
	if c.check.Disabled() {
		if c.check.Checked {
			res = theme.NewThemedResource(theme.CheckButtonCheckedIcon())
		}
		res.ColorName = theme.ColorNameDisabled
		bgRes.ColorName = theme.ColorNameBackground
	}
	c.icon.Resource = res
	c.icon.Refresh()
	c.bg.Resource = bgRes
	c.bg.Refresh()
}

// must be called while holding c.check.propertyLock for reading
func (c *checkRenderer) updateFocusIndicator(th fyne.Theme, v fyne.ThemeVariant) {
	if c.check.Disabled() {
		c.focusIndicator.FillColor = color.Transparent
	} else if c.check.focused {
		c.focusIndicator.FillColor = th.Color(theme.ColorNameFocus, v)
	} else if c.check.hovered {
		c.focusIndicator.FillColor = th.Color(theme.ColorNameHover, v)
	} else {
		c.focusIndicator.FillColor = color.Transparent
	}
}

func focusIfNotMobile(w fyne.Widget) {
	if w == nil {
		return
	}

	if !fyne.CurrentDevice().IsMobile() {
		if c := fyne.CurrentApp().Driver().CanvasForObject(w); c != nil {
			c.Focus(w.(fyne.Focusable))
		}
	}
}
