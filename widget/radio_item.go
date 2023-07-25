package widget

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
)

var _ fyne.Widget = (*radioItem)(nil)
var _ desktop.Hoverable = (*radioItem)(nil)
var _ fyne.Tappable = (*radioItem)(nil)
var _ fyne.Focusable = (*radioItem)(nil)

func newRadioItem(label string, onTap func(*radioItem)) *radioItem {
	i := &radioItem{Label: label, onTap: onTap}
	i.ExtendBaseWidget(i)
	return i
}

// radioItem is a single radio item to be used by RadioGroup.
type radioItem struct {
	DisableableWidget

	Label    string
	Selected bool

	focused bool
	hovered bool
	onTap   func(item *radioItem)
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer.
//
// Implements: fyne.Widget
func (i *radioItem) CreateRenderer() fyne.WidgetRenderer {
	focusIndicator := canvas.NewCircle(color.Transparent)
	// TODO move to `theme.RadioButtonFillIcon()` when we add it in 2.4
	icon := canvas.NewImageFromResource(fyne.CurrentApp().Settings().Theme().Icon("iconNameRadioButtonFill"))
	over := canvas.NewImageFromResource(theme.NewThemedResource(theme.RadioButtonIcon()))
	label := canvas.NewText(i.Label, theme.ForegroundColor())
	label.Alignment = fyne.TextAlignLeading
	r := &radioItemRenderer{
		BaseRenderer:   widget.NewBaseRenderer([]fyne.CanvasObject{focusIndicator, icon, over, label}),
		focusIndicator: focusIndicator,
		icon:           icon,
		over:           over,
		item:           i,
		label:          label,
	}
	r.update()
	return r
}

// FocusGained is called when this item gained the focus.
//
// Implements: fyne.Focusable
func (i *radioItem) FocusGained() {
	i.focused = true
	i.Refresh()
}

// FocusLost is called when this item lost the focus.
//
// Implements: fyne.Focusable
func (i *radioItem) FocusLost() {
	i.focused = false
	i.Refresh()
}

// MouseIn is called when a desktop pointer enters the widget.
//
// Implements: desktop.Hoverable
func (i *radioItem) MouseIn(_ *desktop.MouseEvent) {
	if i.Disabled() {
		return
	}

	i.hovered = true
	i.Refresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget.
//
// Implements: desktop.Hoverable
func (i *radioItem) MouseMoved(_ *desktop.MouseEvent) {
}

// MouseOut is called when a desktop pointer exits the widget
//
// Implements: desktop.Hoverable
func (i *radioItem) MouseOut() {
	if i.Disabled() {
		return
	}

	i.hovered = false
	i.Refresh()
}

// SetSelected sets whether this radio item is selected or not.
func (i *radioItem) SetSelected(selected bool) {
	if i.Disabled() || i.Selected == selected {
		return
	}

	i.Selected = selected
	i.Refresh()
}

// Tapped is called when a pointer tapped event is captured and triggers any change handler
//
// Implements: fyne.Tappable
func (i *radioItem) Tapped(_ *fyne.PointEvent) {
	if !i.focused && !fyne.CurrentDevice().IsMobile() {
		impl := i.super()

		if c := fyne.CurrentApp().Driver().CanvasForObject(impl); c != nil {
			c.Focus(impl.(fyne.Focusable))
		}
	}
	i.toggle()
}

// TypedKey is called when this item receives a key event.
//
// Implements: fyne.Focusable
func (i *radioItem) TypedKey(_ *fyne.KeyEvent) {
}

// TypedRune is called when this item receives a char event.
//
// Implements: fyne.Focusable
func (i *radioItem) TypedRune(r rune) {
	if r == ' ' {
		i.toggle()
	}
}

func (i *radioItem) toggle() {
	if i.Disabled() || i.onTap == nil {
		return
	}

	i.onTap(i)
}

type radioItemRenderer struct {
	widget.BaseRenderer

	focusIndicator *canvas.Circle
	icon, over     *canvas.Image
	item           *radioItem
	label          *canvas.Text
}

func (r *radioItemRenderer) Layout(size fyne.Size) {
	focusIndicatorSize := fyne.NewSquareSize(theme.IconInlineSize() + theme.InnerPadding())
	r.focusIndicator.Resize(focusIndicatorSize)
	r.focusIndicator.Move(fyne.NewPos(theme.InputBorderSize(), (size.Height-focusIndicatorSize.Height)/2))

	labelSize := fyne.NewSize(size.Width, size.Height)
	r.label.Resize(labelSize)
	r.label.Move(fyne.NewPos(focusIndicatorSize.Width+theme.Padding(), 0))

	iconPos := fyne.NewPos(theme.InnerPadding()/2+theme.InputBorderSize(), (size.Height-theme.IconInlineSize())/2)
	iconSize := fyne.NewSquareSize(theme.IconInlineSize())
	r.icon.Resize(iconSize)
	r.icon.Move(iconPos)
	r.over.Resize(iconSize)
	r.over.Move(iconPos)
}

func (r *radioItemRenderer) MinSize() fyne.Size {
	inPad := theme.InnerPadding() * 2

	return r.label.MinSize().
		Add(fyne.NewSize(inPad, inPad)).
		Add(fyne.NewSize(theme.IconInlineSize()+theme.Padding(), 0))
}

func (r *radioItemRenderer) Refresh() {
	r.update()
	canvas.Refresh(r.item.super())
}

func (r *radioItemRenderer) update() {
	r.label.Text = r.item.Label
	r.label.Color = theme.ForegroundColor()
	r.label.TextSize = theme.TextSize()
	if r.item.Disabled() {
		r.label.Color = theme.DisabledColor()
	}

	out := theme.NewThemedResource(theme.RadioButtonIcon())
	out.ColorName = theme.ColorNameInputBorder
	// TODO move to `theme.RadioButtonFillIcon()` when we add it in 2.4
	in := theme.NewThemedResource(fyne.CurrentApp().Settings().Theme().Icon("iconNameRadioButtonFill"))
	in.ColorName = theme.ColorNameInputBackground
	if r.item.Selected {
		in.ColorName = theme.ColorNamePrimary
		out.ColorName = theme.ColorNameForeground
	}
	if r.item.Disabled() {
		if r.item.Selected {
			in.ColorName = theme.ColorNameDisabled
		} else {
			in.ColorName = theme.ColorNameBackground
		}
		out.ColorName = theme.ColorNameDisabled
	}
	r.icon.Resource = in
	r.over.Resource = out

	if r.item.Disabled() {
		r.focusIndicator.FillColor = color.Transparent
	} else if r.item.focused {
		r.focusIndicator.FillColor = theme.FocusColor()
	} else if r.item.hovered {
		r.focusIndicator.FillColor = theme.HoverColor()
	} else {
		r.focusIndicator.FillColor = color.Transparent
	}
}
