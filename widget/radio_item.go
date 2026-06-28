package widget

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
)

var (
	_ fyne.Widget       = (*radioItem)(nil)
	_ desktop.Hoverable = (*radioItem)(nil)
	_ fyne.Tappable     = (*radioItem)(nil)
	_ fyne.Focusable    = (*radioItem)(nil)
)

func newRadioItem(label string, wrap fyne.TextWrap, onTap func(*radioItem)) *radioItem {
	i := &radioItem{Label: label, wrapping: wrap, onTap: onTap}
	i.ExtendBaseWidget(i)
	return i
}

// radioItem is a single radio item to be used by RadioGroup.
type radioItem struct {
	DisableableWidget

	Label    string
	Selected bool

	wrapping fyne.TextWrap

	focused bool
	hovered bool
	onTap   func(item *radioItem)
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer.
func (i *radioItem) CreateRenderer() fyne.WidgetRenderer {
	th := i.Theme()
	label := NewRichTextWithText(i.Label)
	label.Wrapping = i.wrapping
	label.inset = fyne.NewSquareSize(th.Size(theme.SizeNameInnerPadding))
	r := &radioItemRenderer{item: i, label: label}
	r.SetObjects([]fyne.CanvasObject{&r.focusIndicator, &r.icon, &r.over, label})
	r.update()
	return r
}

// FocusGained is called when this item gained the focus.
func (i *radioItem) FocusGained() {
	i.focused = true
	i.Refresh()
}

// FocusLost is called when this item lost the focus.
func (i *radioItem) FocusLost() {
	i.focused = false
	i.Refresh()
}

// MouseIn is called when a desktop pointer enters the widget.
func (i *radioItem) MouseIn(_ *desktop.MouseEvent) {
	if i.Disabled() {
		return
	}

	i.hovered = true
	i.Refresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget.
func (i *radioItem) MouseMoved(_ *desktop.MouseEvent) {
}

// MouseOut is called when a desktop pointer exits the widget
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
func (i *radioItem) Tapped(_ *fyne.PointEvent) {
	if !i.focused {
		focusIfNotMobile(i.super())
	}
	i.toggle()
}

// TypedKey is called when this item receives a key event.
func (i *radioItem) TypedKey(_ *fyne.KeyEvent) {
}

// TypedRune is called when this item receives a char event.
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
	label          *RichText
}

// labelOffset returns the x offset of the label inside the radio item.
func (r *radioItemRenderer) labelOffset() float32 {
	th := r.item.Theme()
	focusIndicatorWidth := th.Size(theme.SizeNameInlineIcon) + th.Size(theme.SizeNameInnerPadding)
	return focusIndicatorWidth + th.Size(theme.SizeNamePadding)
}

func (r *radioItemRenderer) Layout(size fyne.Size) {
	th := r.item.Theme()
	innerPadding := th.Size(theme.SizeNameInnerPadding)
	borderSize := th.Size(theme.SizeNameInputBorder)
	iconInlineSize := th.Size(theme.SizeNameInlineIcon)

	focusIndicatorSize := fyne.NewSquareSize(iconInlineSize + innerPadding)
	r.focusIndicator.Resize(focusIndicatorSize)
	r.focusIndicator.Move(fyne.NewPos(borderSize, (size.Height-focusIndicatorSize.Height)/2))

	labelX := r.labelOffset()
	labelWidth := size.Width
	if r.item.wrapping != fyne.TextWrapOff {
		// constrain wrap to the space remaining after the icon column
		labelWidth = size.Width - labelX
		if labelWidth < 0 {
			labelWidth = 0
		}
	}
	r.label.Resize(fyne.NewSize(labelWidth, size.Height))
	r.label.Move(fyne.NewPos(labelX, 0))

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

	seg := r.label.Segments[0].(*TextSegment)
	seg.Text = r.item.Label
	seg.Style.SizeName = theme.SizeNameText
	if r.item.Disabled() {
		seg.Style.ColorName = theme.ColorNameDisabled
	} else {
		seg.Style.ColorName = theme.ColorNameForeground
	}
	r.label.Wrapping = r.item.wrapping
	r.label.inset = fyne.NewSquareSize(th.Size(theme.SizeNameInnerPadding))
	r.label.Refresh()

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
	r.icon.Refresh()
	r.over.Resource = out
	r.over.Refresh()

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
