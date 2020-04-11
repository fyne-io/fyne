package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal/cache"
	"fyne.io/fyne/theme"
)

// AccordionContainer displays a list of AccordionItems.
// Each item is represented by a button that reveals a detailed view when tapped.
type AccordionContainer struct {
	BaseWidget

	items []*accordionItem
}

// NewAccordionContainer creates a new accordion widget.
func NewAccordionContainer() *AccordionContainer {
	a := &AccordionContainer{
		BaseWidget: BaseWidget{},
	}
	a.ExtendBaseWidget(a)
	return a
}

// MinSize returns the size that this widget should not shrink below.
func (a *AccordionContainer) MinSize() fyne.Size {
	a.ExtendBaseWidget(a)
	return a.BaseWidget.MinSize()
}

// Append adds the given header and detail view to this AccordionContainer.
func (a *AccordionContainer) Append(header string, detail fyne.CanvasObject) {
	item := &accordionItem{
		BaseWidget: BaseWidget{},
		container:  a,
		detail:     detail,
		open:       true,
	}
	item.header = &accordionItemHeader{
		BaseWidget: BaseWidget{},
		item:       item,
		text:       header,
	}
	item.header.ExtendBaseWidget(item.header)
	item.ExtendBaseWidget(item)

	r := cache.Renderer(a).(*accordionContainerRenderer)
	a.items = append(a.items, item)
	r.objects = append(r.objects, item)
	a.Refresh()
}

// Remove deletes the item at the given index from this AccordionContainer.
func (a *AccordionContainer) Remove(index int) {
	r := cache.Renderer(a).(*accordionContainerRenderer)
	a.items = append(a.items[:index], a.items[index+1:]...)
	r.objects = append(r.objects[:index], r.objects[index+1:]...)
	a.Refresh()
}

// Open expands the item at the given index.
func (a *AccordionContainer) Open(index int) {
	if index < 0 || index >= len(a.items) {
		return
	}
	a.items[index].setOpen(true)
	a.Refresh()
}

// OpenAll expands all items.
func (a *AccordionContainer) OpenAll() {
	for _, i := range a.items {
		i.setOpen(true)
	}
	a.Refresh()
}

// Close collapses the item at the given index.
func (a *AccordionContainer) Close(index int) {
	if index < 0 || index >= len(a.items) {
		return
	}
	a.items[index].setOpen(false)
	a.Refresh()
}

// CloseAll collapses all items.
func (a *AccordionContainer) CloseAll() {
	for _, i := range a.items {
		i.setOpen(false)
	}
	a.Refresh()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (a *AccordionContainer) CreateRenderer() fyne.WidgetRenderer {
	var os []fyne.CanvasObject
	for _, i := range a.items {
		os = append(os, i)
	}
	return &accordionContainerRenderer{
		container: a,
		objects:   os,
	}
}

type accordionContainerRenderer struct {
	container *AccordionContainer
	objects   []fyne.CanvasObject
}

func (r *accordionContainerRenderer) MinSize() fyne.Size {
	width := 0
	height := 0
	for _, i := range r.container.items {
		min := i.MinSize()
		width = fyne.Max(width, min.Width)
		height += min.Height
	}
	if len(r.container.items) > 0 {
		width += theme.Padding() * 2
		height += theme.Padding() * 2
	}
	return fyne.NewSize(width, height)
}

func (r *accordionContainerRenderer) Layout(size fyne.Size) {
	x := theme.Padding()
	y := theme.Padding()
	w := size.Width - theme.Padding()*2
	for _, i := range r.container.items {
		h := i.MinSize().Height
		i.Move(fyne.NewPos(x, y))
		i.Resize(fyne.NewSize(w, h))
		y += h
	}
}

func (r *accordionContainerRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *accordionContainerRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *accordionContainerRenderer) Refresh() {
	for _, i := range r.container.items {
		i.Refresh()
	}
	canvas.Refresh(r.container)
}

func (r *accordionContainerRenderer) Destroy() {
}

var _ fyne.Widget = (*accordionItem)(nil)

type accordionItem struct {
	BaseWidget
	container *AccordionContainer
	header    *accordionItemHeader
	detail    fyne.CanvasObject
	open      bool
}

func (i *accordionItem) setOpen(open bool) {
	if i.open == open {
		return
	}
	if i.open = open; i.open {
		i.detail.Show()
	} else {
		i.detail.Hide()
	}
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (i *accordionItem) CreateRenderer() fyne.WidgetRenderer {
	return &accordionItemRenderer{
		item: i,
		objects: []fyne.CanvasObject{
			i.header,
			i.detail,
		},
	}
}

type accordionItemRenderer struct {
	item    *accordionItem
	objects []fyne.CanvasObject
}

func (r *accordionItemRenderer) MinSize() fyne.Size {
	min := r.item.header.MinSize()
	width := min.Width
	height := min.Height
	if r.item.open {
		min := r.item.detail.MinSize()
		width = fyne.Max(width, min.Width)
		height += min.Height
	}
	return fyne.NewSize(width, height)
}

func (r *accordionItemRenderer) Layout(size fyne.Size) {
	height := r.item.header.MinSize().Height
	r.item.header.Move(fyne.NewPos(0, 0))
	r.item.header.Resize(fyne.NewSize(size.Width, height))
	if r.item.open {
		r.item.detail.Move(fyne.NewPos(0, height))
		r.item.detail.Resize(fyne.NewSize(size.Width, r.item.detail.MinSize().Height))
	}
}

func (r *accordionItemRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *accordionItemRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *accordionItemRenderer) Refresh() {
	r.item.header.Refresh()
	r.item.detail.Refresh()
	canvas.Refresh(r.item)
}

func (r *accordionItemRenderer) Destroy() {
}

var _ fyne.Tappable = (*accordionItemHeader)(nil)
var _ fyne.Widget = (*accordionItemHeader)(nil)
var _ desktop.Hoverable = (*accordionItemHeader)(nil)

type accordionItemHeader struct {
	BaseWidget
	item    *accordionItem
	text    string
	hovered bool
}

// Tapped is called when a pointer tapped event is captured and triggers any tap handler
func (h *accordionItemHeader) Tapped(*fyne.PointEvent) {
	h.item.setOpen(!h.item.open)
	h.item.container.Refresh()
}

// MouseIn is called when a desktop pointer enters the widget
func (h *accordionItemHeader) MouseIn(*desktop.MouseEvent) {
	h.hovered = true
	h.item.container.Refresh()
}

// MouseOut is called when a desktop pointer exits the widget
func (h *accordionItemHeader) MouseOut() {
	h.hovered = false
	h.item.container.Refresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (h *accordionItemHeader) MouseMoved(*desktop.MouseEvent) {
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (h *accordionItemHeader) CreateRenderer() fyne.WidgetRenderer {
	r := &accordionItemHeaderRenderer{
		header: h,
		images: []*canvas.Image{
			canvas.NewImageFromResource(theme.ContentRemoveIcon()),
			canvas.NewImageFromResource(theme.ContentAddIcon()),
		},
		text: &canvas.Text{},
	}
	r.updateCanvasObjects()
	return r
}

type accordionItemHeaderRenderer struct {
	header *accordionItemHeader
	images []*canvas.Image
	image  *canvas.Image
	text   *canvas.Text
}

func (r *accordionItemHeaderRenderer) MinSize() fyne.Size {
	min := r.text.MinSize()
	width := theme.IconInlineSize() + theme.Padding() + min.Width
	height := fyne.Max(theme.IconInlineSize(), min.Height)
	return fyne.NewSize(width, height)
}

func (r *accordionItemHeaderRenderer) Layout(size fyne.Size) {
	r.image.Move(fyne.NewPos(0, 0))
	r.image.Resize(fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
	r.text.Move(fyne.NewPos(theme.IconInlineSize()+theme.Padding(), 0))
	r.text.Resize(fyne.NewSize(size.Width, r.text.MinSize().Height))
}

func (r *accordionItemHeaderRenderer) BackgroundColor() color.Color {
	if r.header.hovered {
		return theme.HoverColor()
	}
	return theme.BackgroundColor()
}

func (r *accordionItemHeaderRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.image, r.text}
}

func (r *accordionItemHeaderRenderer) Refresh() {
	r.updateCanvasObjects()
	canvas.Refresh(r.header)
}

func (r *accordionItemHeaderRenderer) updateCanvasObjects() {
	if r.header.item.open {
		r.image = r.images[0]
	} else {
		r.image = r.images[1]
	}
	r.text.Text = r.header.text
	r.text.Color = theme.TextColor()
	r.text.TextSize = theme.TextSize()
	r.text.TextStyle = fyne.TextStyle{}
}

func (r *accordionItemHeaderRenderer) Destroy() {
}
