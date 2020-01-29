package widget

import (
	"image/color"
	"math"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/theme"
)

const noRadioItemIndex = -1

type radioRenderItem struct {
	icon  *canvas.Image
	label *canvas.Text

	focusIndicator *canvas.Circle
}

type radioRenderer struct {
	items []*radioRenderItem

	objects []fyne.CanvasObject
	radio   *Radio
}

func removeDuplicates(options []string) []string {
	var result []string
	found := make(map[string]bool)

	for _, option := range options {
		if _, ok := found[option]; !ok {
			found[option] = true
			result = append(result, option)
		}
	}

	return result
}

// MinSize calculates the minimum size of a radio item.
// This is based on the contained text, the radio icon and a standard amount of padding
// between each item.
func (r *radioRenderer) MinSize() fyne.Size {
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

// Layout the components of the radio widget
func (r *radioRenderer) Layout(size fyne.Size) {
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

// applyTheme updates this Radio to match the current system theme
func (r *radioRenderer) applyTheme() {
	for _, item := range r.items {
		item.label.Color = theme.TextColor()
		item.label.TextSize = theme.TextSize()
		if r.radio.Disabled() {
			item.label.Color = theme.DisabledTextColor()
		}
	}
}

func (r *radioRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *radioRenderer) Refresh() {
	r.applyTheme()
	r.radio.removeDuplicateOptions()

	if len(r.items) < len(r.radio.Options) {
		for i := len(r.items); i < len(r.radio.Options); i++ {
			option := r.radio.Options[i]
			icon := canvas.NewImageFromResource(theme.RadioButtonIcon())

			text := canvas.NewText(option, theme.TextColor())
			text.Alignment = fyne.TextAlignLeading

			focusIndicator := canvas.NewCircle(theme.BackgroundColor())

			r.objects = append(r.objects, focusIndicator, icon, text)
			r.items = append(r.items, &radioRenderItem{icon, text, focusIndicator})
		}
		r.Layout(r.radio.Size())
	} else if len(r.items) > len(r.radio.Options) {
		total := len(r.radio.Options)
		r.items = r.items[:total]
		r.objects = r.objects[:total*2]
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

	canvas.Refresh(r.radio.super())
}

func (r *radioRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *radioRenderer) Destroy() {
}

// Radio widget has a list of text labels and radio check icons next to each.
// Changing the selection (only one can be selected) will trigger the changed func.
type Radio struct {
	DisableableWidget
	Options  []string
	Selected string

	OnChanged  func(string) `json:"-"`
	Horizontal bool

	hoveredItemIndex int
	hovered          bool
}

// indexByPosition returns the item index for a specified position or noRadioItemIndex if any
func (r *Radio) indexByPosition(pos fyne.Position) int {
	index := 0
	if r.Horizontal {
		index = int(math.Floor(float64(pos.X) / float64(r.itemWidth())))
	} else {
		index = int(math.Floor(float64(pos.Y) / float64(r.itemHeight())))
	}
	if index < 0 || index >= len(r.Options) { // in the padding
		return noRadioItemIndex
	}
	return index
}

// MouseIn is called when a desktop pointer enters the widget
func (r *Radio) MouseIn(event *desktop.MouseEvent) {
	if r.Disabled() {
		return
	}

	r.hoveredItemIndex = r.indexByPosition(event.Position)
	r.hovered = true
	r.Refresh()
}

// MouseOut is called when a desktop pointer exits the widget
func (r *Radio) MouseOut() {
	r.hoveredItemIndex = noRadioItemIndex
	r.hovered = false
	r.Refresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (r *Radio) MouseMoved(event *desktop.MouseEvent) {
	if r.Disabled() {
		return
	}

	r.hoveredItemIndex = r.indexByPosition(event.Position)
	r.hovered = true
	r.Refresh()
}

// Append adds a new option to the end of a Radio widget.
func (r *Radio) Append(option string) {
	r.Options = append(r.Options, option)

	r.Refresh()
}

// Tapped is called when a pointer tapped event is captured and triggers any change handler
func (r *Radio) Tapped(event *fyne.PointEvent) {
	if r.Disabled() {
		return
	}

	index := r.indexByPosition(event.Position)

	if index < 0 || index >= len(r.Options) { // in the padding
		return
	}
	clicked := r.Options[index]

	if r.Selected == clicked {
		r.Selected = ""
	} else {
		r.Selected = clicked
	}

	if r.OnChanged != nil {
		r.OnChanged(r.Selected)
	}
	r.Refresh()
}

// TappedSecondary is called when a secondary pointer tapped event is captured
func (r *Radio) TappedSecondary(*fyne.PointEvent) {
}

// MinSize returns the size that this widget should not shrink below
func (r *Radio) MinSize() fyne.Size {
	r.ExtendBaseWidget(r)
	return r.BaseWidget.MinSize()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (r *Radio) CreateRenderer() fyne.WidgetRenderer {
	r.ExtendBaseWidget(r)
	var items []*radioRenderItem
	var objects []fyne.CanvasObject

	for _, option := range r.Options {
		icon := canvas.NewImageFromResource(theme.RadioButtonIcon())

		text := canvas.NewText(option, theme.TextColor())
		text.Alignment = fyne.TextAlignLeading

		focusIndicator := canvas.NewCircle(theme.BackgroundColor())

		objects = append(objects, focusIndicator, icon, text)
		items = append(items, &radioRenderItem{icon, text, focusIndicator})
	}

	return &radioRenderer{items, objects, r}
}

// SetSelected sets the radio option, it can be used to set a default option.
func (r *Radio) SetSelected(option string) {
	if r.Selected == option {
		return
	}

	r.Selected = option

	r.Refresh()
}

func (r *Radio) itemHeight() int {
	if r.Horizontal {
		return r.MinSize().Height
	}

	count := 1
	if r.Options != nil {
		count = len(r.Options)
	}
	return r.MinSize().Height / count
}

func (r *Radio) itemWidth() int {
	if !r.Horizontal {
		return r.MinSize().Width
	}

	return r.MinSize().Width / len(r.Options)
}

func (r *Radio) removeDuplicateOptions() {
	r.Options = removeDuplicates(r.Options)
}

// NewRadio creates a new radio widget with the set options and change handler
func NewRadio(options []string, changed func(string)) *Radio {
	r := &Radio{
		DisableableWidget: DisableableWidget{},
		Options:           options,
		OnChanged:         changed,
	}

	r.removeDuplicateOptions()
	r.ExtendBaseWidget(r)
	return r
}
