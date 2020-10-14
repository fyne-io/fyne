package widget

import (
	"math"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/theme"
)

const noRadioGroupItemIndex = -1

// RadioGroup widget has a list of text labels and radio check icons next to each.
// Changing the selection (only one can be selected) will trigger the changed func.
type RadioGroup struct {
	DisableableWidget
	Horizontal bool
	Required   bool
	OnChanged  func(string) `json:"-"`
	Options    []string
	Selected   string

	hoveredItemIndex int
	hovered          bool
}

var _ fyne.Widget = (*RadioGroup)(nil)

// NewRadioGroup creates a new radio widget with the set options and change handler
func NewRadioGroup(options []string, changed func(string)) *RadioGroup {
	r := &RadioGroup{
		DisableableWidget: DisableableWidget{},
		Options:           options,
		OnChanged:         changed,
	}

	r.removeDuplicateOptions()
	r.ExtendBaseWidget(r)
	return r
}

// Append adds a new option to the end of a RadioGroup widget.
func (r *RadioGroup) Append(option string) {
	r.Options = append(r.Options, option)

	r.Refresh()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (r *RadioGroup) CreateRenderer() fyne.WidgetRenderer {
	r.ExtendBaseWidget(r)
	r.propertyLock.Lock()
	defer r.propertyLock.Unlock()
	var items []*radioGroupRenderItem
	var objects []fyne.CanvasObject

	for _, option := range r.Options {
		icon := canvas.NewImageFromResource(theme.RadioButtonIcon())

		text := canvas.NewText(option, theme.TextColor())
		text.Alignment = fyne.TextAlignLeading

		focusIndicator := canvas.NewCircle(theme.BackgroundColor())

		objects = append(objects, focusIndicator, icon, text)
		items = append(items, &radioGroupRenderItem{icon, text, focusIndicator})
	}

	rr := &radioGroupRenderer{widget.NewBaseRenderer(objects), items, r}
	rr.applyTheme()
	r.removeDuplicateOptions()
	rr.updateItems()
	return rr
}

// MinSize returns the size that this widget should not shrink below
func (r *RadioGroup) MinSize() fyne.Size {
	r.ExtendBaseWidget(r)
	return r.BaseWidget.MinSize()
}

// MouseIn is called when a desktop pointer enters the widget
func (r *RadioGroup) MouseIn(event *desktop.MouseEvent) {
	if r.Disabled() {
		return
	}

	r.hoveredItemIndex = r.indexByPosition(event.Position)
	r.hovered = true
	r.Refresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (r *RadioGroup) MouseMoved(event *desktop.MouseEvent) {
	if r.Disabled() {
		return
	}

	r.hoveredItemIndex = r.indexByPosition(event.Position)
	r.hovered = true
	r.Refresh()
}

// MouseOut is called when a desktop pointer exits the widget
func (r *RadioGroup) MouseOut() {
	r.hoveredItemIndex = noRadioGroupItemIndex
	r.hovered = false
	r.Refresh()
}

// SetSelected sets the radio option, it can be used to set a default option.
func (r *RadioGroup) SetSelected(option string) {
	if r.Selected == option {
		return
	}

	r.Selected = option

	if r.OnChanged != nil {
		r.OnChanged(option)
	}

	r.Refresh()
}

// Tapped is called when a pointer tapped event is captured and triggers any change handler
func (r *RadioGroup) Tapped(event *fyne.PointEvent) {
	if r.Disabled() {
		return
	}

	index := r.indexByPosition(event.Position)

	if index < 0 || index >= len(r.Options) { // in the padding
		return
	}
	clicked := r.Options[index]

	if r.Selected == clicked {
		if r.Required {
			return
		}
		r.Selected = ""
	} else {
		r.Selected = clicked
	}

	if r.OnChanged != nil {
		r.OnChanged(r.Selected)
	}
	r.Refresh()
}

// indexByPosition returns the item index for a specified position or noRadioGroupItemIndex if any
func (r *RadioGroup) indexByPosition(pos fyne.Position) int {
	index := 0
	if r.Horizontal {
		index = int(math.Floor(float64(pos.X) / float64(r.itemWidth())))
	} else {
		index = int(math.Floor(float64(pos.Y) / float64(r.itemHeight())))
	}
	if index < 0 || index >= len(r.Options) { // in the padding
		return noRadioGroupItemIndex
	}
	return index
}

func (r *RadioGroup) itemHeight() int {
	if r.Horizontal {
		return r.MinSize().Height
	}

	count := 1
	if r.Options != nil && len(r.Options) > 0 {
		count = len(r.Options)
	}
	return r.MinSize().Height / count
}

func (r *RadioGroup) itemWidth() int {
	if !r.Horizontal {
		return r.MinSize().Width
	}

	count := 1
	if r.Options != nil && len(r.Options) > 0 {
		count = len(r.Options)
	}
	return r.MinSize().Width / count
}

func (r *RadioGroup) removeDuplicateOptions() {
	r.Options = removeDuplicates(r.Options)
}

type radioGroupRenderItem struct {
	icon  *canvas.Image
	label *canvas.Text

	focusIndicator *canvas.Circle
}

type radioGroupRenderer struct {
	widget.BaseRenderer
	items []*radioGroupRenderItem
	radio *RadioGroup
}

// Layout the components of the radio widget
func (r *radioGroupRenderer) Layout(size fyne.Size) {
	itemWidth := r.radio.itemWidth()
	itemHeight := r.radio.itemHeight()
	labelSize := fyne.NewSize(itemWidth, itemHeight)
	focusIndicatorSize := fyne.NewSize(theme.IconInlineSize()+theme.Padding()*2, theme.IconInlineSize()+theme.Padding()*2)

	x, y := 0, 0
	for _, item := range r.items {
		item.focusIndicator.Resize(focusIndicatorSize)
		item.focusIndicator.Move(fyne.NewPos(x, y+(itemHeight-focusIndicatorSize.Height)/2))

		item.label.Resize(labelSize)
		item.label.Move(fyne.NewPos(x+focusIndicatorSize.Width+theme.Padding(), y))

		item.icon.Resize(fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
		item.icon.Move(fyne.NewPos(x+theme.Padding(),
			y+(labelSize.Height-theme.IconInlineSize())/2))

		if r.radio.Horizontal {
			x += itemWidth
		} else {
			y += itemHeight
		}
	}
}

// MinSize calculates the minimum size of a radio item.
// This is based on the contained text, the radio icon and a standard amount of padding
// between each item.
func (r *radioGroupRenderer) MinSize() fyne.Size {
	width := 0
	height := 0
	for _, item := range r.items {
		itemMin := item.label.MinSize().Add(fyne.NewSize(theme.Padding()*4, theme.Padding()*2))
		itemMin = itemMin.Add(fyne.NewSize(theme.IconInlineSize()+theme.Padding(), 0))

		if r.radio.Horizontal {
			height = fyne.Max(height, itemMin.Height)
			width += itemMin.Width
		} else {
			width = fyne.Max(width, itemMin.Width)
			height += itemMin.Height
		}
	}

	return fyne.NewSize(width, height)
}

func (r *radioGroupRenderer) Refresh() {
	r.radio.propertyLock.Lock()
	r.applyTheme()
	r.radio.removeDuplicateOptions()
	r.updateItems()
	r.radio.propertyLock.Unlock()
	canvas.Refresh(r.radio.super())
}

// applyTheme updates this RadioGroup to match the current system theme
func (r *radioGroupRenderer) applyTheme() {
	for _, item := range r.items {
		item.label.Color = theme.TextColor()
		item.label.TextSize = theme.TextSize()
		if r.radio.Disabled() {
			item.label.Color = theme.DisabledTextColor()
		}
	}
}

func (r *radioGroupRenderer) updateItems() {
	if len(r.items) < len(r.radio.Options) {
		for i := len(r.items); i < len(r.radio.Options); i++ {
			option := r.radio.Options[i]
			icon := canvas.NewImageFromResource(theme.RadioButtonIcon())

			text := canvas.NewText(option, theme.TextColor())
			text.Alignment = fyne.TextAlignLeading

			focusIndicator := canvas.NewCircle(theme.BackgroundColor())

			r.SetObjects(append(r.Objects(), focusIndicator, icon, text))
			r.items = append(r.items, &radioGroupRenderItem{icon, text, focusIndicator})
		}
		r.Layout(r.radio.Size())
	} else if len(r.items) > len(r.radio.Options) {
		total := len(r.radio.Options)
		r.items = r.items[:total]
		r.SetObjects(r.Objects()[:total*2])
	}

	for i, item := range r.items {
		option := r.radio.Options[i]
		item.label.Text = option

		res := theme.RadioButtonIcon()
		if r.radio.Selected == option {
			res = theme.RadioButtonCheckedIcon()
		}
		if r.radio.Disabled() {
			res = theme.NewDisabledResource(res)
		}
		item.icon.Resource = res

		if r.radio.Disabled() {
			item.focusIndicator.FillColor = theme.BackgroundColor()
		} else if r.radio.hovered && r.radio.hoveredItemIndex == i {
			item.focusIndicator.FillColor = theme.HoverColor()
		} else {
			item.focusIndicator.FillColor = theme.BackgroundColor()
		}
	}
}
