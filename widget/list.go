package widget

import (
	"image/color"
	"math"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
)

// List is a widget that caches list items for performance and
// lays the items out in a vertical direction inside of a scroller
type List struct {
	BaseWidget

	Length         func() int
	CreateItem     func() fyne.CanvasObject
	UpdateItem     func(index int, item fyne.CanvasObject)
	OnItemSelected func(index int, item fyne.CanvasObject)

	background    color.Color
	selectedItem  *listItem
	selectedIndex int
	itemMin       fyne.Size
}

// NewList creates and returns a list widget for displaying items in
// a vertical layout with scrolling and caching for performance
func NewList(length func() int, createItem func() fyne.CanvasObject, updateItem func(index int, item fyne.CanvasObject)) *List {
	return &List{BaseWidget: BaseWidget{}, Length: length, CreateItem: createItem, UpdateItem: updateItem, selectedIndex: -1}
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (l *List) CreateRenderer() fyne.WidgetRenderer {
	l.ExtendBaseWidget(l)

	if l.itemMin.Width == 0 && l.itemMin.Height == 0 && l.CreateItem != nil {
		l.itemMin = newListItem(l.CreateItem(), nil).MinSize()
	}
	layout := fyne.NewContainerWithLayout(newListLayout(l.itemMin, l.Length()))
	scroller := NewVScrollContainer(layout)

	scroller.onOffsetChanged = func() {
		ll := layout.Layout.(*listLayout)
		if ll.offsetX == scroller.Offset.X && ll.offsetY == scroller.Offset.Y {
			return
		}
		ll.offsetX = scroller.Offset.X
		ll.offsetY = scroller.Offset.Y
		l.BaseWidget.Refresh()
	}

	objects := []fyne.CanvasObject{scroller}
	return &listRenderer{BaseRenderer: widget.NewBaseRenderer(objects), list: l, scroller: scroller, layout: layout}
}

// MinSize returns the size that this widget should not shrink below
func (l *List) MinSize() fyne.Size {
	l.ExtendBaseWidget(l)

	return l.BaseWidget.MinSize()
}

// Refresh updates list to match the current theme
func (l *List) Refresh() {
	if l.background != nil {
		l.background = theme.BackgroundColor()
	}

	l.BaseWidget.Refresh()
}

type listRenderer struct {
	widget.BaseRenderer
	list      *List
	scroller  *ScrollContainer
	layout    *fyne.Container
	itemCache []fyne.CanvasObject
}

func (l *listRenderer) BackgroundColor() color.Color {
	if l.list.background == nil {
		return theme.BackgroundColor()
	}

	return l.list.background
}

func (l *listRenderer) Layout(size fyne.Size) {
	l.scroller.Resize(size)
	if l.list.Length() == 0 {
		return
	}
	firstItemIndex := 0
	if l.scroller.Offset.Y > 0 {
		firstItemIndex = int(math.Ceil(float64(l.scroller.Offset.Y) / float64(l.list.itemMin.Height)))
	}
	lastItemIndex := int(math.Ceil(float64(l.scroller.Offset.Y+size.Height) / float64(l.list.itemMin.Height)))
	if lastItemIndex >= l.list.Length() {
		lastItemIndex = l.list.Length() - 1
	}
	visibleItemCount := int(math.Ceil(float64(size.Height) / float64(l.list.itemMin.Height)))

	if len(l.itemCache) < visibleItemCount && len(l.itemCache) < l.list.Length() {
		for i := len(l.itemCache); i < visibleItemCount && i < l.list.Length(); i++ {
			item := newListItem(l.list.CreateItem(), nil)
			l.itemCache = append(l.itemCache, item)
		}
	} else if len(l.itemCache) > l.list.Length() {
		for i := l.list.Length() - 1; i < len(l.itemCache); i++ {
			l.itemCache[i].Hide()
		}
		l.itemCache = l.itemCache[:l.list.Length()]
	} else if len(l.itemCache) > visibleItemCount {
		for i := visibleItemCount + 1; i < len(l.itemCache); i++ {
			l.itemCache[i].Hide()
		}
		l.itemCache = l.itemCache[:visibleItemCount+1]
	}
	l.layout.Objects = l.itemCache
	l.layout.Layout.Layout(l.layout.Objects, l.layout.Layout.MinSize(nil))

	j := 0
	for i := firstItemIndex; i <= lastItemIndex && j < len(l.itemCache); i++ {
		item := l.itemCache[j].(*listItem)
		index := i
		if index != l.list.selectedIndex {
			item.selected = false
		} else {
			item.selected = true
			l.list.selectedItem = item
		}
		l.list.UpdateItem(index, item.child)
		item.onTapped = func() {
			if l.list.selectedItem != nil && l.list.selectedIndex >= firstItemIndex && l.list.selectedIndex <= lastItemIndex {
				l.list.selectedItem.selected = false
				l.list.selectedItem.Refresh()
			}
			l.list.selectedItem = item
			l.list.selectedIndex = index
			l.list.selectedItem.selected = true
			l.list.selectedItem.Refresh()
			if l.list.OnItemSelected != nil {
				l.list.OnItemSelected(index, item.child)
			}
		}
		item.Refresh()
		j++
	}
}

func (l *listRenderer) MinSize() fyne.Size {
	return l.scroller.MinSize()
}

func (l *listRenderer) Refresh() {
	l.Layout(l.scroller.Size())
	canvas.Refresh(l.list.super())
}

type listItem struct {
	DisableableWidget

	background        color.Color
	onTapped          func()
	statusIndicator   *canvas.Rectangle
	child             fyne.CanvasObject
	hovered, selected bool
}

func newListItem(child fyne.CanvasObject, tapped func()) *listItem {
	li := &listItem{
		child:    child,
		onTapped: tapped,
	}

	li.ExtendBaseWidget(li)
	return li
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (li *listItem) CreateRenderer() fyne.WidgetRenderer {
	li.ExtendBaseWidget(li)

	li.statusIndicator = canvas.NewRectangle(theme.BackgroundColor())

	objects := []fyne.CanvasObject{li.statusIndicator, li.child}

	return &listItemRenderer{widget.NewShadowingRenderer(objects, widget.ButtonLevel), li, layout.NewHBoxLayout()}
}

// MinSize returns the size that this widget should not shrink below
func (li *listItem) MinSize() fyne.Size {
	li.ExtendBaseWidget(li)
	return li.BaseWidget.MinSize()
}

// MouseIn is called when a desktop pointer enters the widget
func (li *listItem) MouseIn(*desktop.MouseEvent) {
	li.hovered = true
	li.Refresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (li *listItem) MouseMoved(*desktop.MouseEvent) {
}

// MouseOut is called when a desktop pointer exits the widget
func (li *listItem) MouseOut() {
	li.hovered = false
	li.Refresh()
}

// Tapped is called when a pointer tapped event is captured and triggers any tap handler
func (li *listItem) Tapped(*fyne.PointEvent) {
	if li.onTapped != nil && !li.Disabled() {
		li.selected = true
		li.Refresh()
		li.onTapped()
	}
}

type listItemRenderer struct {
	*widget.ShadowingRenderer

	item   *listItem
	layout fyne.Layout
}

// MinSize calculates the minimum size of a listItem.
// This is based on the size of the status indicator and the size of the child object.
func (li *listItemRenderer) MinSize() (size fyne.Size) {
	li.item.statusIndicator.SetMinSize(fyne.NewSize(theme.Padding(), li.item.child.Size().Height))
	size = fyne.NewSize(theme.Padding()+li.item.child.Size().Width, li.item.child.Size().Height)
	return
}

// Layout the components of the button widget
func (li *listItemRenderer) Layout(size fyne.Size) {
	li.LayoutShadow(size, fyne.NewPos(0, 0))
	li.layout.Layout([]fyne.CanvasObject{li.item.statusIndicator, li.item.child}, size)
}

func (li *listItemRenderer) BackgroundColor() color.Color {
	if li.item.background == nil {
		return theme.BackgroundColor()
	}

	return li.item.background
}

func (li *listItemRenderer) Refresh() {
	if li.item.selected == true {
		li.item.statusIndicator.FillColor = theme.FocusColor()
		canvas.Refresh(li.item.super())
		return
	}
	if li.item.hovered == true {
		li.item.statusIndicator.FillColor = theme.HoverColor()
	} else {
		li.item.statusIndicator.FillColor = theme.BackgroundColor()
	}
	canvas.Refresh(li.item.super())
}

type listLayout struct {
	itemMin          fyne.Size
	itemCount        int
	offsetX, offsetY int
}

var _ fyne.Layout = (*listLayout)(nil)

func (l *listLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	x, y := l.offsetX, l.offsetY

	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		child.Move(fyne.NewPos(x, y))
		y += l.itemMin.Height
		child.Resize(fyne.NewSize(l.itemMin.Width, l.itemMin.Height))
	}
}

func (l *listLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {

	return fyne.NewSize(l.itemMin.Width+theme.Padding(),
		l.itemMin.Height*l.itemCount)
}

func newListLayout(itemMin fyne.Size, itemCount int) fyne.Layout {
	return &listLayout{itemMin: itemMin, itemCount: itemCount}
}
