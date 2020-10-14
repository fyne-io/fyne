package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/theme"
)

var _ fyne.Widget = (*RadioItem)(nil)
var _ desktop.Hoverable = (*RadioItem)(nil)
var _ fyne.Tappable = (*RadioItem)(nil)

func newRadioItem(label string, onTap func(*RadioItem)) *RadioItem {
	i := &RadioItem{Label: label, onTap: onTap}
	i.ExtendBaseWidget(i)
	return i
}

// RadioItem is a single radio item to be used by RadioGroup.
type RadioItem struct {
	DisableableWidget

	Label    string
	Selected bool

	hovered bool
	onTap   func(item *RadioItem)
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer.
//
// Implements: fyne.Widget
func (i *RadioItem) CreateRenderer() fyne.WidgetRenderer {
	focusIndicator := canvas.NewCircle(theme.BackgroundColor())
	icon := canvas.NewImageFromResource(theme.RadioButtonIcon())
	label := canvas.NewText(i.Label, theme.TextColor())
	label.Alignment = fyne.TextAlignLeading
	r := &radioItemRenderer{
		BaseRenderer:   widget.NewBaseRenderer([]fyne.CanvasObject{focusIndicator, icon, label}),
		focusIndicator: focusIndicator,
		icon:           icon,
		item:           i,
		label:          label,
	}
	r.update()
	return r
}

// MouseIn is called when a desktop pointer enters the widget.
//
// Implements: desktop.Hoverable
func (i *RadioItem) MouseIn(_ *desktop.MouseEvent) {
	if i.Disabled() {
		return
	}

	i.hovered = true
	i.Refresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget.
//
// Implements: desktop.Hoverable
func (i *RadioItem) MouseMoved(_ *desktop.MouseEvent) {
}

// MouseOut is called when a desktop pointer exits the widget
//
// Implements: desktop.Hoverable
func (i *RadioItem) MouseOut() {
	if i.Disabled() {
		return
	}

	i.hovered = false
	i.Refresh()
}

// SetSelected sets whether this radio item is selected or not.
func (i *RadioItem) SetSelected(selected bool) {
	if i.Disabled() {
		return
	}
	if i.Selected == selected {
		return
	}
	i.Selected = selected
	i.Refresh()
}

// Tapped is called when a pointer tapped event is captured and triggers any change handler
//
// Implements: fyne.Tappable
func (i *RadioItem) Tapped(_ *fyne.PointEvent) {
	if i.Disabled() {
		return
	}
	if i.onTap == nil {
		return
	}

	i.onTap(i)
}

type radioItemRenderer struct {
	widget.BaseRenderer

	focusIndicator *canvas.Circle
	icon           *canvas.Image
	item           *RadioItem
	label          *canvas.Text
}

func (r *radioItemRenderer) Layout(size fyne.Size) {
	labelSize := fyne.NewSize(size.Width, size.Height)
	focusIndicatorSize := fyne.NewSize(theme.IconInlineSize()+theme.Padding()*2, theme.IconInlineSize()+theme.Padding()*2)

	r.focusIndicator.Resize(focusIndicatorSize)
	r.focusIndicator.Move(fyne.NewPos(0, (size.Height-focusIndicatorSize.Height)/2))

	r.label.Resize(labelSize)
	r.label.Move(fyne.NewPos(focusIndicatorSize.Width+theme.Padding(), 0))

	r.icon.Resize(fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
	r.icon.Move(fyne.NewPos(theme.Padding(), (labelSize.Height-theme.IconInlineSize())/2))
}

func (r *radioItemRenderer) MinSize() fyne.Size {
	return r.label.MinSize().
		Add(fyne.NewSize(theme.Padding()*4, theme.Padding()*2)).
		Add(fyne.NewSize(theme.IconInlineSize()+theme.Padding(), 0))
}

func (r *radioItemRenderer) Refresh() {
	r.update()
	canvas.Refresh(r.item.super())
}

func (r *radioItemRenderer) update() {
	r.label.Text = r.item.Label
	r.label.Color = theme.TextColor()
	r.label.TextSize = theme.TextSize()
	if r.item.Disabled() {
		r.label.Color = theme.DisabledTextColor()
	}

	res := theme.RadioButtonIcon()
	if r.item.Selected {
		res = theme.RadioButtonCheckedIcon()
	}
	if r.item.Disabled() {
		res = theme.NewDisabledResource(res)
	}
	r.icon.Resource = res

	if r.item.Disabled() {
		r.focusIndicator.FillColor = theme.BackgroundColor()
	} else if r.item.hovered {
		r.focusIndicator.FillColor = theme.HoverColor()
	} else {
		r.focusIndicator.FillColor = theme.BackgroundColor()
	}
}
