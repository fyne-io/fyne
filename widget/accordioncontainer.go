package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/theme"
)

var _ fyne.Widget = (*AccordionContainer)(nil)

// AccordionContainer displays a list of AccordionItems.
// Each item is represented by a button that reveals a detailed view when tapped.
type AccordionContainer struct {
	BaseWidget
	Items     []*AccordionItem
	MultiOpen bool
}

// NewAccordionContainer creates a new accordion widget.
func NewAccordionContainer(items ...*AccordionItem) *AccordionContainer {
	a := &AccordionContainer{
		Items: items,
	}
	a.ExtendBaseWidget(a)
	return a
}

// Append adds the given item to this AccordionContainer.
func (a *AccordionContainer) Append(item *AccordionItem) {
	a.Items = append(a.Items, item)
	a.Refresh()
}

// Close collapses the item at the given index.
func (a *AccordionContainer) Close(index int) {
	if index < 0 || index >= len(a.Items) {
		return
	}
	a.Items[index].Open = false
	a.Refresh()
}

// CloseAll collapses all items.
func (a *AccordionContainer) CloseAll() {
	for _, i := range a.Items {
		i.Open = false
	}
	a.Refresh()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (a *AccordionContainer) CreateRenderer() fyne.WidgetRenderer {
	a.ExtendBaseWidget(a)
	r := &accordionContainerRenderer{
		container: a,
	}
	r.updateObjects()
	return r
}

// MinSize returns the size that this widget should not shrink below.
func (a *AccordionContainer) MinSize() fyne.Size {
	a.ExtendBaseWidget(a)
	return a.BaseWidget.MinSize()
}

// Open expands the item at the given index.
func (a *AccordionContainer) Open(index int) {
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
func (a *AccordionContainer) OpenAll() {
	if !a.MultiOpen {
		return
	}
	for _, i := range a.Items {
		i.Open = true
	}
	a.Refresh()
}

// Remove deletes the given item from this AccordionContainer.
func (a *AccordionContainer) Remove(item *AccordionItem) {
	for i, ai := range a.Items {
		if ai == item {
			a.RemoveIndex(i)
			break
		}
	}
}

// RemoveIndex deletes the item at the given index from this AccordionContainer.
func (a *AccordionContainer) RemoveIndex(index int) {
	if index < 0 || index >= len(a.Items) {
		return
	}
	a.Items = append(a.Items[:index], a.Items[index+1:]...)
	a.Refresh()
}

type accordionContainerRenderer struct {
	widget.BaseRenderer
	container *AccordionContainer
	headers   []*Button
}

func (r *accordionContainerRenderer) Layout(size fyne.Size) {
	x := 0
	y := 0
	for i, ai := range r.container.Items {
		if i != 0 {
			y += theme.Padding()
		}
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
	}
}

func (r *accordionContainerRenderer) MinSize() (size fyne.Size) {
	for i, ai := range r.container.Items {
		if i != 0 {
			size.Height += theme.Padding()
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

func (r *accordionContainerRenderer) Refresh() {
	r.updateObjects()
	r.Layout(r.container.Size())
	canvas.Refresh(r.container)
}

func (r *accordionContainerRenderer) updateObjects() {
	is := len(r.container.Items)
	hs := len(r.headers)
	i := 0
	for ; i < is; i++ {
		ai := r.container.Items[i]
		var h *Button
		if i < hs {
			h = r.headers[i]
		} else {
			h = &Button{}
			r.headers = append(r.headers, h)
		}
		h.Alignment = ButtonAlignLeading
		h.IconPlacement = ButtonIconLeadingText
		h.Hidden = false
		h.HideShadow = true
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
	r.SetObjects(objects)
}

// AccordionItem represents a single item in an AccordionContainer.
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
