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

	// this index is ONE-BASED so the default zero-value is unselected
	// use r.selectedIndex(), r.setSelectedIndex(int) to maniupulate this field
	// as if it were a zero-based index (with -1 == nothing selected)
	_selIdx int
}

var _ fyne.Widget = (*RadioGroup)(nil)

// NewRadioGroup creates a new radio group widget with the set options and change handler
//
// Since: 1.4
func NewRadioGroup(options []string, changed func(string)) *RadioGroup {
	r := &RadioGroup{
		Options:   options,
		OnChanged: changed,
	}
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

	items := make([]fyne.CanvasObject, len(r.Options))
	for i, option := range r.Options {
		idx := i
		items[idx] = newRadioItem(option, func(item *radioItem) {
			r.itemTapped(item, idx)
		})
	}

	render := &radioGroupRenderer{widget.NewBaseRenderer(items), items, r}
	r.updateSelectedIndex()
	render.updateItems(false)
	return render
}

// MinSize returns the size that this widget should not shrink below
func (r *RadioGroup) MinSize() fyne.Size {
	r.ExtendBaseWidget(r)
	return r.BaseWidget.MinSize()
}

// SetSelected sets the radio option, it can be used to set a default option.
func (r *RadioGroup) SetSelected(option string) {
	if r.Selected == option {
		return
	}

	r.Selected = option
	// selectedIndex will be updated on refresh to the first matching item

	if r.OnChanged != nil {
		r.OnChanged(option)
	}

	r.Refresh()
}

func (r *RadioGroup) itemTapped(item *radioItem, idx int) {
	if r.Disabled() {
		return
	}

	if r.selectedIndex() == idx {
		if r.Required {
			return
		}
		r.Selected = ""
		r.setSelectedIndex(-1)
		item.SetSelected(false)
	} else {
		r.Selected = item.Label
		r.setSelectedIndex(idx)
		item.SetSelected(true)
	}

	if r.OnChanged != nil {
		r.OnChanged(r.Selected)
	}
	r.Refresh()
}

func (r *RadioGroup) Refresh() {
	r.updateSelectedIndex()
	r.BaseWidget.Refresh()
}

func (r *RadioGroup) selectedIndex() int {
	return r._selIdx - 1
}

func (r *RadioGroup) setSelectedIndex(idx int) {
	r._selIdx = idx + 1
}

// if selectedIndex does not match the public Selected property,
// set it to the index of the first radio item whose label matches Selected
func (r *RadioGroup) updateSelectedIndex() {
	sel := r.Selected
	sIdx := r.selectedIndex()
	if sIdx >= 0 && sIdx < len(r.Options) && r.Options[sIdx] == sel {
		return // selected index matches Selected
	}
	if sIdx == -1 && sel == "" {
		return // nothing selected
	}

	sIdx = -1
	for i, opt := range r.Options {
		if sel == opt {
			sIdx = i
			break
		}
	}
	r.setSelectedIndex(sIdx)
}

type radioGroupRenderer struct {
	widget.BaseRenderer

	// slice of *radioItem, but using fyne.CanvasObject as the type
	// so we can directly set it to the BaseRenderer's objects slice
	items []fyne.CanvasObject
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

		width = fyne.Max(width, itemMin.Width)
		height = fyne.Max(height, itemMin.Height)
	}

	if r.radio.Horizontal {
		width = width * float32(len(r.items))
	} else {
		height = height * float32(len(r.items))
	}

	return fyne.NewSize(width, height)
}

func (r *radioGroupRenderer) Refresh() {
	r.updateItems(true)
	canvas.Refresh(r.radio.super())
}

func (r *radioGroupRenderer) updateItems(refresh bool) {
	if len(r.items) < len(r.radio.Options) {
		for i := len(r.items); i < len(r.radio.Options); i++ {
			idx := i
			item := newRadioItem(r.radio.Options[idx], func(item *radioItem) {
				r.radio.itemTapped(item, idx)
			})
			r.items = append(r.items, item)
		}
		r.Layout(r.radio.Size())
	} else if len(r.items) > len(r.radio.Options) {
		total := len(r.radio.Options)
		r.items = r.items[:total]
	}
	r.SetObjects(r.items)

	for i, item := range r.items {
		item := item.(*radioItem)
		changed := false
		if l := r.radio.Options[i]; l != item.Label {
			item.Label = r.radio.Options[i]
			changed = true
		}
		if sel := i == r.radio.selectedIndex(); sel != item.Selected {
			item.Selected = sel
			changed = true
		}
		if d := r.radio.disabled.Load(); d != item.disabled.Load() {
			item.disabled.Store(d)
			changed = true
		}

		if refresh || changed {
			item.Refresh()
		}
	}
}
