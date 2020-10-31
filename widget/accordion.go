package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/theme"
)

const accordionDividerHeight = 1

var _ fyne.Widget = (*Accordion)(nil)

// AccordionContainer displays a list of AccordionItems.
// Each item is represented by a button that reveals a detailed view when tapped.
//
// Deprecated: This has been renamed to Accordion
type AccordionContainer = Accordion

// NewAccordionContainer creates a new accordion widget.
//
// Deprecated: Use NewAccordion instead
func NewAccordionContainer(items ...*AccordionItem) *AccordionContainer {
	a := &Accordion{
		Items: items,
	}
	a.ExtendBaseWidget(a)
	return a
}

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
	r := &accordionRenderer{
		container: a,
	}
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

// OpenAll expands all items.
func (a *Accordion) OpenAll() {
	if !a.MultiOpen {
		return
	}
	for _, i := range a.Items {
		i.Open = true
	}
	a.Refresh()
}

// Remove deletes the given item from this Accordion.
func (a *Accordion) Remove(item *AccordionItem) {
	for i, ai := range a.Items {
		if ai == item {
			a.RemoveIndex(i)
			break
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
	x := 0
	y := 0
	for i, ai := range r.container.Items {
		if i != 0 {
			div := r.dividers[i-1]
			div.Move(fyne.NewPos(x, y))
			div.Resize(fyne.NewSize(size.Width, accordionDividerHeight))
			y += accordionDividerHeight
		}
		y += theme.Padding()

		h := r.headers[i]
		h.Move(fyne.NewPos(x, y))
		min := h.MinSize().Height
		h.Resize(fyne.NewSize(size.Width, min))
		y += min
		if ai.Open {
			y += theme.Padding()
			d := ai.Detail
			d.Move(fyne.NewPos(x, y))
			min := d.MinSize().Height
			d.Resize(fyne.NewSize(size.Width, min))
			y += min
		}

		y += theme.Padding()
	}
}

func (r *accordionRenderer) MinSize() (size fyne.Size) {
	for i, ai := range r.container.Items {
		size.Height += theme.Padding() * 2
		if i != 0 {
			size.Height += accordionDividerHeight
		}
		min := r.headers[i].MinSize()
		size.Width = fyne.Max(size.Width, min.Width)
		size.Height += min.Height
		min = ai.Detail.MinSize()
		size.Width = fyne.Max(size.Width, min.Width)
		if ai.Open {
			size.Height += min.Height
			size.Height += theme.Padding()
		}
	}
	return
}

func (r *accordionRenderer) Refresh() {
	r.updateObjects()
	r.Layout(r.container.Size())
	canvas.Refresh(r.container)
}

func (r *accordionRenderer) updateObjects() {
	is := len(r.container.Items)
	hs := len(r.headers)
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
			h.Icon = theme.MenuDropUpIcon()
			ai.Detail.Show()
		} else {
			h.Icon = theme.MenuDropDownIcon()
			ai.Detail.Hide()
		}
		h.Refresh()
	}
	// Hide extras
	for ; i < hs; i++ {
		r.headers[i].Hide()
	}
	// Set objects
	var objects []fyne.CanvasObject
	for _, h := range r.headers {
		objects = append(objects, h)
	}
	for _, i := range r.container.Items {
		objects = append(objects, i.Detail)
	}
	// add dividers
	for i = 0; i < len(r.dividers); i++ {
		if i < len(r.container.Items)-1 {
			r.dividers[i].Show()
		} else {
			r.dividers[i].Hide()
		}
		objects = append(objects, r.dividers[i])
	}
	// make new dividers
	for ; i < is-1; i++ {
		div := NewSeparator()
		r.dividers = append(r.dividers, div)
		objects = append(objects, div)
	}
	r.SetObjects(objects)
}

// AccordionItem represents a single item in an Accordion.
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
