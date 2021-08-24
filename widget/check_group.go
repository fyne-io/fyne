package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/widget"
)

// CheckGroup widget has a list of text labels and checkbox icons next to each.
// Changing the selection (any number can be selected) will trigger the changed func.
//
// Since: 2.1
type CheckGroup struct {
	DisableableWidget
	Horizontal bool
	Required   bool
	OnChanged  func([]string) `json:"-"`
	Options    []string
	Selected   []string

	items []*Check
}

var _ fyne.Widget = (*CheckGroup)(nil)

// NewCheckGroup creates a new check group widget with the set options and change handler
//
// Since: 2.1
func NewCheckGroup(options []string, changed func([]string)) *CheckGroup {
	r := &CheckGroup{
		DisableableWidget: DisableableWidget{},
		Options:           options,
		OnChanged:         changed,
	}
	r.ExtendBaseWidget(r)
	r.update()
	return r
}

// Append adds a new option to the end of a CheckGroup widget.
func (r *CheckGroup) Append(option string) {
	r.Options = append(r.Options, option)

	r.Refresh()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (r *CheckGroup) CreateRenderer() fyne.WidgetRenderer {
	r.ExtendBaseWidget(r)
	r.propertyLock.Lock()
	defer r.propertyLock.Unlock()

	r.update()
	objects := make([]fyne.CanvasObject, len(r.items))
	for i, item := range r.items {
		objects[i] = item
	}

	return &checkGroupRenderer{widget.NewBaseRenderer(objects), r.items, r}
}

// MinSize returns the size that this widget should not shrink below
func (r *CheckGroup) MinSize() fyne.Size {
	r.ExtendBaseWidget(r)
	return r.BaseWidget.MinSize()
}

// Refresh causes this widget to be redrawn in it's current state.
//
// Implements: fyne.CanvasObject
func (r *CheckGroup) Refresh() {
	r.propertyLock.Lock()
	r.update()
	r.propertyLock.Unlock()
	r.BaseWidget.Refresh()
}

// SetSelected sets the checked options, it can be used to set a default option.
func (r *CheckGroup) SetSelected(options []string) {
	//if r.Selected == options {
	//	return
	//}

	r.Selected = options

	if r.OnChanged != nil {
		r.OnChanged(options)
	}

	r.Refresh()
}

func (r *CheckGroup) itemTapped(item *Check) {
	if r.Disabled() {
		return
	}

	contains := false
	for i, s := range r.Selected {
		if s == item.Text {
			contains = true
			if len(r.Selected) <= 1 {
				if r.Required {
					item.SetChecked(true)
					return
				}
				r.Selected = nil
			} else {
				r.Selected = append(r.Selected[:i], r.Selected[i+1:]...)
			}
			break
		}
	}

	if !contains {
		r.Selected = append(r.Selected, item.Text)
	}

	if r.OnChanged != nil {
		r.OnChanged(r.Selected)
	}
	r.Refresh()
}

func (r *CheckGroup) update() {
	r.Options = removeDuplicates(r.Options)
	if len(r.items) < len(r.Options) {
		for i := len(r.items); i < len(r.Options); i++ {
			var item *Check
			item = NewCheck(r.Options[i], func(bool) {
				r.itemTapped(item)
			})
			r.items = append(r.items, item)
		}
	} else if len(r.items) > len(r.Options) {
		r.items = r.items[:len(r.Options)]
	}
	for i, item := range r.items {
		contains := false
		for _, s := range r.Selected {
			if s == item.Text {
				contains = true
				break
			}
		}

		item.Text = r.Options[i]
		item.Checked = contains
		item.DisableableWidget.disabled = r.disabled
		item.Refresh()
	}
}

type checkGroupRenderer struct {
	widget.BaseRenderer
	items  []*Check
	checks *CheckGroup
}

// Layout the components of the checks widget
func (r *checkGroupRenderer) Layout(_ fyne.Size) {
	count := 1
	if r.items != nil && len(r.items) > 0 {
		count = len(r.items)
	}
	var itemHeight, itemWidth float32
	minSize := r.checks.MinSize()
	if r.checks.Horizontal {
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
		if r.checks.Horizontal {
			x += itemWidth
		} else {
			y += itemHeight
		}
	}
}

// MinSize calculates the minimum size of a checks item.
// This is based on the contained text, the checks icon and a standard amount of padding
// between each item.
func (r *checkGroupRenderer) MinSize() fyne.Size {
	width := float32(0)
	height := float32(0)
	for _, item := range r.items {
		itemMin := item.MinSize()
		if r.checks.Horizontal {
			height = fyne.Max(height, itemMin.Height)
			width += itemMin.Width
		} else {
			width = fyne.Max(width, itemMin.Width)
			height += itemMin.Height
		}
	}

	return fyne.NewSize(width, height)
}

func (r *checkGroupRenderer) Refresh() {
	r.updateItems()
	canvas.Refresh(r.checks.super())
}

func (r *checkGroupRenderer) updateItems() {
	if len(r.items) < len(r.checks.Options) {
		for i := len(r.items); i < len(r.checks.Options); i++ {
			var item *Check
			item = NewCheck(r.checks.Options[i], func(bool) {
				r.checks.itemTapped(item)
			})
			r.SetObjects(append(r.Objects(), item))
			r.items = append(r.items, item)
		}
		r.Layout(r.checks.Size())
	} else if len(r.items) > len(r.checks.Options) {
		total := len(r.checks.Options)
		r.items = r.items[:total]
		r.SetObjects(r.Objects()[:total])
	}
	for i, item := range r.items {
		contains := false
		for _, s := range r.checks.Selected {
			if s == item.Text {
				contains = true
				break
			}
		}
		item.Text = r.checks.Options[i]
		item.Checked = contains
		item.disabled = r.checks.disabled
		item.Refresh()
	}
}
