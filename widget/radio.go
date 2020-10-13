package widget

import (
	"math"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/widget"
)

const noRadioItemIndex = -1

// Radio widget has a list of text labels and radio check icons next to each.
// Changing the selection (only one can be selected) will trigger the changed func.
type Radio struct {
	DisableableWidget
	Horizontal bool
	Required   bool
	OnChanged  func(string) `json:"-"`
	Options    []string
	Selected   string

	items []*radioItem
}

var _ fyne.Widget = (*Radio)(nil)

// NewRadio creates a new radio widget with the set options and change handler
func NewRadio(options []string, changed func(string)) *Radio {
	r := &Radio{
		DisableableWidget: DisableableWidget{},
		Options:           options,
		OnChanged:         changed,
	}
	r.ExtendBaseWidget(r)
	r.update()
	return r
}

// Append adds a new option to the end of a Radio widget.
func (r *Radio) Append(option string) {
	r.Options = append(r.Options, option)

	r.Refresh()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (r *Radio) CreateRenderer() fyne.WidgetRenderer {
	r.ExtendBaseWidget(r)
	r.propertyLock.Lock()
	r.propertyLock.Unlock()

	r.update()
	var objects []fyne.CanvasObject
	for _, item := range r.items {
		objects = append(objects, item)
	}
	return &radioRenderer{widget.NewBaseRenderer(objects), r.items, r}
}

// MinSize returns the size that this widget should not shrink below
func (r *Radio) MinSize() fyne.Size {
	r.ExtendBaseWidget(r)
	return r.BaseWidget.MinSize()
}

// Refresh causes this widget to be redrawn in it's current state.
//
// Implements: fyne.CanvasObject
func (r *Radio) Refresh() {
	r.propertyLock.Lock()
	r.update()
	r.propertyLock.Unlock()
	r.BaseWidget.Refresh()
}

// SetSelected sets the radio option, it can be used to set a default option.
func (r *Radio) SetSelected(option string) {
	if r.Selected == option {
		return
	}

	r.Selected = option

	if r.OnChanged != nil {
		r.OnChanged(option)
	}

	r.Refresh()
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

func (r *Radio) itemHeight() int {
	if r.Horizontal {
		return r.MinSize().Height
	}

	count := 1
	if r.Options != nil && len(r.Options) > 0 {
		count = len(r.Options)
	}
	return r.MinSize().Height / count
}

func (r *Radio) itemWidth() int {
	if !r.Horizontal {
		return r.MinSize().Width
	}

	count := 1
	if r.Options != nil && len(r.Options) > 0 {
		count = len(r.Options)
	}
	return r.MinSize().Width / count
}

func (r *Radio) itemTapped(item *radioItem) {
	if r.Disabled() {
		return
	}

	if r.Selected == item.Label {
		if r.Required {
			return
		}
		r.Selected = ""
		item.SetSelected(false)
	} else {
		r.Selected = item.Label
		item.SetSelected(true)
	}

	if r.OnChanged != nil {
		r.OnChanged(r.Selected)
	}
	r.Refresh()
}

func (r *Radio) update() {
	r.Options = removeDuplicates(r.Options)
	if len(r.items) < len(r.Options) {
		for i := len(r.items); i < len(r.Options); i++ {
			item := newRadioItem(r.Options[i], r.itemTapped)
			r.items = append(r.items, item)
		}
	} else if len(r.items) > len(r.Options) {
		r.items = r.items[:len(r.Options)]
	}
	for i, item := range r.items {
		item.Label = r.Options[i]
		item.Selected = item.Label == r.Selected
		item.DisableableWidget.disabled = r.disabled
		item.Refresh()
	}
}

type radioRenderer struct {
	widget.BaseRenderer
	items []*radioItem
	radio *Radio
}

// Layout the components of the radio widget
func (r *radioRenderer) Layout(_ fyne.Size) {
	itemWidth := r.radio.itemWidth()
	itemHeight := r.radio.itemHeight()
	itemSize := fyne.NewSize(itemWidth, itemHeight)
	x, y := 0, 0
	for _, item := range r.items {
		item.Resize(itemSize)
		item.Move(fyne.NewPos(x, y))
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
func (r *radioRenderer) MinSize() fyne.Size {
	width := 0
	height := 0
	for _, item := range r.items {
		itemMin := item.MinSize()
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

func (r *radioRenderer) Refresh() {
	r.updateItems()
	canvas.Refresh(r.radio.super())
}

func (r *radioRenderer) updateItems() {
	if len(r.items) < len(r.radio.Options) {
		for i := len(r.items); i < len(r.radio.Options); i++ {
			item := newRadioItem(r.radio.Options[i], r.radio.itemTapped)
			r.SetObjects(append(r.Objects(), item))
			r.items = append(r.items, item)
		}
		r.Layout(r.radio.Size())
	} else if len(r.items) > len(r.radio.Options) {
		total := len(r.radio.Options)
		r.items = r.items[:total]
		r.SetObjects(r.Objects()[:total])
	}
	for i, item := range r.items {
		item.Label = r.radio.Options[i]
		item.Selected = item.Label == r.radio.Selected
		item.Refresh()
	}
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
