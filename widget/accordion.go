package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
)

var _ fyne.Widget = (*Accordion)(nil)

// Accordion displays a list of AccordionItems.
// Each item is represented by a button that reveals a detailed view when tapped.
type Accordion struct {
	BaseWidget
	Items     []*AccordionItem
	MultiOpen bool
}

// NewAccordion creates a new accordion widget.
func NewAccordion(items ...*AccordionItem) *Accordion {
	a := &Accordion{
		Items: items,
	}
	a.ExtendBaseWidget(a)
	return a
}

// Append adds the given item to this Accordion.
func (a *Accordion) Append(item *AccordionItem) {
	a.Items = append(a.Items, item)

	a.Refresh()
}

// Close collapses the item at the given index.
func (a *Accordion) Close(index int) {
	if index < 0 || index >= len(a.Items) {
		return
	}
	a.Items[index].Open = false

	a.Refresh()
}

// CloseAll collapses all items.
func (a *Accordion) CloseAll() {
	for _, i := range a.Items {
		i.Open = false
	}

	a.Refresh()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (a *Accordion) CreateRenderer() fyne.WidgetRenderer {
	a.ExtendBaseWidget(a)
	r := &accordionRenderer{container: a}
	r.updateObjects()
	return r
}

// MinSize returns the size that this widget should not shrink below.
func (a *Accordion) MinSize() fyne.Size {
	a.ExtendBaseWidget(a)
	return a.BaseWidget.MinSize()
}

// Open expands the item at the given index.
func (a *Accordion) Open(index int) {
	if index < 0 || index >= len(a.Items) {
		return
	}

	for i, ai := range a.Items {
		if i == index {
			ai.Open = true
		} else if !a.MultiOpen {
			ai.Open = false
		}
	}

	a.Refresh()
}

// OpenAll expands all items, note that your Accordion should have [MultiOpen] set to `true` for this to operate as
// expected. For single-open accordions it will open only the first item.
func (a *Accordion) OpenAll() {
	if !a.MultiOpen {
		a.Open(0)
		return
	}

	for _, i := range a.Items {
		i.Open = true
	}

	a.Refresh()
}

// Prepend adds the given item to the beginning of this Accordion.
//
// Since: 2.6
func (a *Accordion) Prepend(item *AccordionItem) {
	a.Items = append([]*AccordionItem{item}, a.Items...)

	a.Refresh()
}

// Remove deletes the given item from this Accordion.
func (a *Accordion) Remove(item *AccordionItem) {
	for i, ai := range a.Items {
		if ai == item {
			a.Items = append(a.Items[:i], a.Items[i+1:]...)
			return
		}
	}
}

// RemoveIndex deletes the item at the given index from this Accordion.
func (a *Accordion) RemoveIndex(index int) {
	if index < 0 || index >= len(a.Items) {
		return
	}
	a.Items = append(a.Items[:index], a.Items[index+1:]...)

	a.Refresh()
}

type accordionRenderer struct {
	widget.BaseRenderer
	container *Accordion
	headers   []*Button
	dividers  []fyne.CanvasObject
}

func (r *accordionRenderer) Layout(size fyne.Size) {
	r.updateObjects()
	th := r.container.Theme()
	pad := th.Size(theme.SizeNamePadding)
	separator := th.Size(theme.SizeNameSeparatorThickness)
	dividerOff := (pad + separator) / 2
	x := float32(0)
	y := float32(0)
	hasOpen := 0

	for i, ai := range r.container.Items {
		h := r.headers[i]
		min := h.MinSize().Height
		y += min

		if ai.Open {
			y += pad + ai.Detail.MinSize().Height
			hasOpen++
		}
		if i < len(r.container.Items)-1 {
			y += pad
		}
	}

	extra := (size.Height - y) / float32(hasOpen)
	if extra < 0 {
		extra = 0
	}
	y = 0
	for i, ai := range r.container.Items {
		if i != 0 {
			div := r.dividers[i-1]
			if i > 0 {
				div.Move(fyne.NewPos(x, y-dividerOff))
			}
			div.Resize(fyne.NewSize(size.Width, separator))
		}

		h := r.headers[i]
		h.Move(fyne.NewPos(x, y))
		min := h.MinSize().Height
		h.Resize(fyne.NewSize(size.Width, min))
		y += min

		if ai.Open {
			y += pad
			d := ai.Detail
			d.Move(fyne.NewPos(x, y))

			openSize := ai.Detail.MinSize().Height + extra
			d.Resize(fyne.NewSize(size.Width, openSize))
			y += openSize
		}
		if i < len(r.container.Items)-1 {
			y += pad
		}
	}
}

func (r *accordionRenderer) MinSize() fyne.Size {
	r.updateObjects()
	th := r.container.Theme()
	pad := th.Size(theme.SizeNamePadding)
	size := fyne.Size{}

	for i, ai := range r.container.Items {
		if i != 0 {
			size.Height += pad
		}
		min := r.headers[i].MinSize()
		size.Width = fyne.Max(size.Width, min.Width)
		size.Height += min.Height
		min = ai.Detail.MinSize()
		size.Width = fyne.Max(size.Width, min.Width)
		if ai.Open {
			size.Height += min.Height
			size.Height += pad
		}
	}

	return size
}

func (r *accordionRenderer) Refresh() {
	r.updateObjects()
	r.Layout(r.container.Size())
	canvas.Refresh(r.container)
}

func (r *accordionRenderer) updateObjects() {
	th := r.container.Theme()
	is := len(r.container.Items)
	hs := len(r.headers)
	ds := len(r.dividers)
	i := 0
	for ; i < is; i++ {
		ai := r.container.Items[i]
		var h *Button
		if i < hs {
			h = r.headers[i]
			h.Show()
		} else {
			h = &Button{}
			r.headers = append(r.headers, h)
			hs++
		}
		h.Alignment = ButtonAlignLeading
		h.IconPlacement = ButtonIconLeadingText
		h.Hidden = false
		h.Importance = LowImportance
		h.Text = ai.Title
		index := i // capture
		h.OnTapped = func() {
			if ai.Open {
				r.container.Close(index)
			} else {
				r.container.Open(index)
			}
		}
		if ai.Open {
			h.Icon = th.Icon(theme.IconNameArrowDropUp)
			ai.Detail.Show()
		} else {
			h.Icon = th.Icon(theme.IconNameArrowDropDown)
			ai.Detail.Hide()
		}
		h.Refresh()
	}
	// Hide extras
	for ; i < hs; i++ {
		r.headers[i].Hide()
	}
	// Set objects
	objects := make([]fyne.CanvasObject, hs+is+ds)
	for i, header := range r.headers {
		objects[i] = header
	}
	for i, item := range r.container.Items {
		objects[hs+i] = item.Detail
	}
	// add dividers
	for i = 0; i < ds; i++ {
		if i < len(r.container.Items)-1 {
			r.dividers[i].Show()
		} else {
			r.dividers[i].Hide()
		}
		objects[hs+is+i] = r.dividers[i]
	}
	// make new dividers
	for ; i < is-1; i++ {
		div := NewSeparator()
		r.dividers = append(r.dividers, div)
		objects = append(objects, div)
	}

	r.SetObjects(objects)
}

// AccordionItem represents a single item in an Acc rdion.
type AccordionItem struct {
	Title  string
	Detail fyne.CanvasObject
	Open   bool
}

// NewAccordionItem creates a new item for an Accordion.
func NewAccordionItem(title string, detail fyne.CanvasObject) *AccordionItem {
	return &AccordionItem{
		Title:  title,
		Detail: detail,
	}
}
