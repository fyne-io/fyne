package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/widget"
)

// RadioGroup widget has a list of text labels and checks check icons next to each.
// Changing the selection (only one can be selected) will trigger the changed func.
//
// Since: 1.4
type RadioGroup struct {
	DisableableWidget
	Horizontal bool
	Required   bool
	OnChanged  func(string) `json:"-"`
	Options    []string
	Selected   string

	items []*radioItem
}

var _ fyne.Widget = (*RadioGroup)(nil)

// NewRadioGroup creates a new radio group widget with the set options and change handler
//
// Since: 1.4
func NewRadioGroup(options []string, changed func(string)) *RadioGroup {
	r := &RadioGroup{
		DisableableWidget: DisableableWidget{},
		Options:           options,
		OnChanged:         changed,
	}
	r.ExtendBaseWidget(r)
	r.update()
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

	r.update()
	objects := make([]fyne.CanvasObject, len(r.items))
	for i, item := range r.items {
		objects[i] = item
	}

	return &radioGroupRenderer{widget.NewBaseRenderer(objects), r.items, r}
}

// MinSize returns the size that this widget should not shrink below
func (r *RadioGroup) MinSize() fyne.Size {
	r.ExtendBaseWidget(r)
	return r.BaseWidget.MinSize()
}

// Refresh causes this widget to be redrawn in it's current state.
//
// Implements: fyne.CanvasObject
func (r *RadioGroup) Refresh() {
	r.propertyLock.Lock()
	r.update()
	r.propertyLock.Unlock()
	r.BaseWidget.Refresh()
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

func (r *RadioGroup) itemTapped(item *radioItem) {
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

func (r *RadioGroup) update() {
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

type radioGroupRenderer struct {
	widget.BaseRenderer
	items []*radioItem
	radio *RadioGroup
}

// Layout the components of the radio widget
func (r *radioGroupRenderer) Layout(_ fyne.Size) {
	count := 1
	if r.items != nil && len(r.items) > 0 {
		count = len(r.items)
	}
	var itemHeight, itemWidth float32
	minSize := r.radio.MinSize()
	if r.radio.Horizontal {
		itemHeight = minSize.Height
		itemWidth = minSize.Width / float32(count)
	} else {
		itemHeight = minSize.Height / float32(count)
		itemWidth = minSize.Width
	}

	itemSize := fyne.NewSize(itemWidth, itemHeight)
	x, y := float32(0), float32(0)
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
func (r *radioGroupRenderer) MinSize() fyne.Size {
	width := float32(0)
	height := float32(0)
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

func (r *radioGroupRenderer) Refresh() {
	r.updateItems()
	canvas.Refresh(r.radio.super())
}

func (r *radioGroupRenderer) updateItems() {
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
		item.disabled = r.radio.disabled
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
