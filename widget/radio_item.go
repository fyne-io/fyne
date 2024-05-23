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
	txt := canvas.Text{Alignment: fyne.TextAlignLeading}
	txt.TextSize = i.Theme().Size(theme.SizeNameText)
	r := &radioItemRenderer{item: i, label: &txt}
	r.SetObjects([]fyne.CanvasObject{&r.focusIndicator, &r.icon, &r.over, &txt})
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
	if !i.focused {
		focusIfNotMobile(i.super())
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
	item *radioItem

	focusIndicator canvas.Circle
	icon, over     canvas.Image
	label          *canvas.Text
}

func (r *radioItemRenderer) Layout(size fyne.Size) {
	th := r.item.Theme()
	innerPadding := th.Size(theme.SizeNameInnerPadding)
	borderSize := th.Size(theme.SizeNameInputBorder)
	iconInlineSize := th.Size(theme.SizeNameInlineIcon)

	focusIndicatorSize := fyne.NewSquareSize(iconInlineSize + innerPadding)
	r.focusIndicator.Resize(focusIndicatorSize)
	r.focusIndicator.Move(fyne.NewPos(borderSize, (size.Height-focusIndicatorSize.Height)/2))

	labelSize := fyne.NewSize(size.Width, size.Height)
	r.label.Resize(labelSize)
	r.label.Move(fyne.NewPos(focusIndicatorSize.Width+th.Size(theme.SizeNamePadding), 0))

	iconPos := fyne.NewPos(innerPadding/2+borderSize, (size.Height-iconInlineSize)/2)
	iconSize := fyne.NewSquareSize(iconInlineSize)
	r.icon.Resize(iconSize)
	r.icon.Move(iconPos)
	r.over.Resize(iconSize)
	r.over.Move(iconPos)
}

func (r *radioItemRenderer) MinSize() fyne.Size {
	th := r.item.Theme()
	inPad := th.Size(theme.SizeNameInnerPadding) * 2

	return r.label.MinSize().
		AddWidthHeight(inPad+th.Size(theme.SizeNameInlineIcon)+th.Size(theme.SizeNamePadding), inPad)
}

func (r *radioItemRenderer) Refresh() {
	r.update()
	canvas.Refresh(r.item.super())
}

func (r *radioItemRenderer) update() {
	th := r.item.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	r.label.Text = r.item.Label
	r.label.TextSize = th.Size(theme.SizeNameText)
	if r.item.Disabled() {
		r.label.Color = th.Color(theme.ColorNameDisabled, v)
	} else {
		r.label.Color = th.Color(theme.ColorNameForeground, v)
	}

	out := theme.NewThemedResource(th.Icon(theme.IconNameRadioButton))
	out.ColorName = theme.ColorNameInputBorder
	in := theme.NewThemedResource(th.Icon(theme.IconNameRadioButtonFill))
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
		r.focusIndicator.FillColor = th.Color(theme.ColorNameFocus, v)
	} else if r.item.hovered {
		r.focusIndicator.FillColor = th.Color(theme.ColorNameHover, v)
	} else {
		r.focusIndicator.FillColor = color.Transparent
	}
}
