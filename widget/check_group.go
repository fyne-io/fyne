package widget

import (
	"math"
	"strings"

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
	numColumns int
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
		Options:   options,
		OnChanged: changed,
	}
	r.ExtendBaseWidget(r)
	r.update()
	return r
}

func (r *CheckGroup) SetColumns(columns int) {
	r.numColumns = columns
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

// Remove removes the first occurrence of the specified option found from a CheckGroup widget.
// Return true if an option was removed.
//
// Since: 2.3
func (r *CheckGroup) Remove(option string) bool {
	for i, o := range r.Options {
		if strings.EqualFold(option, o) {
			r.Options = append(r.Options[:i], r.Options[i+1:]...)
			for j, s := range r.Selected {
				if strings.EqualFold(option, s) {
					r.Selected = append(r.Selected[:j], r.Selected[j+1:]...)
					break
				}
			}
			r.Refresh()
			return true
		}
	}
	return false
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
		item.DisableableWidget.disabled.Store(r.disabled.Load())
		item.Refresh()
	}
}

type checkGroupRenderer struct {
	widget.BaseRenderer
	items  []*Check
	checks *CheckGroup
}

func (r *CheckGroup) countColumns() int {
	if r.numColumns < 1 {
		r.numColumns = 1
	}
	return r.numColumns
}

func (r *CheckGroup) countRows() int {
	return int(math.Ceil(float64(len(r.items)) / float64(r.countColumns())))
}

func getLeading(size float64, offset int) float32 {
	return float32(size * float64(offset))
}

func getTrailing(size float64, offset int) float32 {
	return getLeading(size, offset+1)
}

// Layout the components of the checks widget
func (r *checkGroupRenderer) Layout(_ fyne.Size) {
	cols := r.checks.countColumns()

	primaryObjects := cols
	secondaryObjects := r.checks.countRows()
	if r.checks.Horizontal {
		primaryObjects, secondaryObjects = secondaryObjects, primaryObjects
	}

	size := r.checks.MinSize()
	cellWidth := float64(size.Width) / float64(primaryObjects)
	cellHeight := float64(size.Height) / float64(secondaryObjects)

	row, col := 0, 0
	i := 0
	for _, child := range r.items {
		x1 := getLeading(cellWidth, col)
		y1 := getLeading(cellHeight, row)
		x2 := getTrailing(cellWidth, col)
		y2 := getTrailing(cellHeight, row)

		child.Move(fyne.NewPos(x1, y1))
		child.Resize(fyne.NewSize(x2-x1, y2-y1))

		if r.checks.Horizontal {
			if (i+1)%cols == 0 {
				col++
				row = 0
			} else {
				row++
			}
		} else {
			if (i+1)%cols == 0 {
				row++
				col = 0
			} else {
				col++
			}
		}
		i++
	}
}

// MinSize calculates the minimum size of a checks item.
// This is based on the contained text, the checks icon and a standard amount of padding
// between each item.
func (r *checkGroupRenderer) MinSize() fyne.Size {
	minSize := fyne.NewSize(0, 0)
	for _, child := range r.items {
		minSize = minSize.Max(child.MinSize())
	}

	primaryObjects := r.checks.countColumns()
	secondaryObjects := r.checks.countRows()
	if r.checks.Horizontal {
		primaryObjects, secondaryObjects = secondaryObjects, primaryObjects
	}

	width := minSize.Width * float32(primaryObjects)
	height := minSize.Height * float32(secondaryObjects)

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
		r.Layout(r.checks.Size()) // argument is ignored
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
		item.disabled.Store(r.checks.disabled.Load())
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
