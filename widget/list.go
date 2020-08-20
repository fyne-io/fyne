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

type listItem struct {
	DisableableWidget

	background        color.Color
	onTapped          func()
	statusIndicator   *canvas.Rectangle
	child             fyne.CanvasObject
	hovered, selected bool
}

// Tapped is called when a pointer tapped event is captured and triggers any tap handler
func (li *listItem) Tapped(*fyne.PointEvent) {
	if li.onTapped != nil && !li.Disabled() {
		li.selected = true
		li.Refresh()
		li.onTapped()
	}
}

// MouseIn is called when a desktop pointer enters the widget
func (li *listItem) MouseIn(*desktop.MouseEvent) {
	li.hovered = true
	li.Refresh()
}

// MouseOut is called when a desktop pointer exits the widget
func (li *listItem) MouseOut() {
	li.hovered = false
	li.Refresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (li *listItem) MouseMoved(*desktop.MouseEvent) {
}

// MinSize returns the size that this widget should not shrink below
func (li *listItem) MinSize() fyne.Size {
	li.ExtendBaseWidget(li)
	return li.BaseWidget.MinSize()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (li *listItem) CreateRenderer() fyne.WidgetRenderer {
	li.ExtendBaseWidget(li)

	li.statusIndicator = canvas.NewRectangle(theme.BackgroundColor())

	objects := []fyne.CanvasObject{li.statusIndicator, li.child}

	return &listItemRenderer{widget.NewShadowingRenderer(objects, widget.ButtonLevel), li, layout.NewHBoxLayout()}
}

func newListItem(child fyne.CanvasObject, tapped func()) *listItem {
	li := &listItem{
		child:    child,
		onTapped: tapped,
	}

	li.ExtendBaseWidget(li)
	return li
}

type listLayout struct {
	itemMin          fyne.Size
	itemCount        int
	offsetX, offsetY int
}

var _ fyne.Layout = (*listLayout)(nil)

func (l *listLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	total := 0
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		total += l.itemMin.Height
	}

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

	return fyne.NewSize(l.itemMin.Width,
		l.itemMin.Height*l.itemCount)
}

func newListLayout(itemMin fyne.Size, itemCount int) fyne.Layout {
	return &listLayout{itemMin: itemMin, itemCount: itemCount}
}

type listRenderer struct {
	widget.BaseRenderer
	list      *List
	itemCache []fyne.CanvasObject
}

func (l *listRenderer) MinSize() fyne.Size {
	return l.list.scroller.MinSize()
}

func (l *listRenderer) Layout(size fyne.Size) {
	l.list.scroller.Resize(size)
	if l.list.Length == 0 {
		return
	}
	firstItemIndex := 0
	if l.list.scroller.Offset.Y > 0 {
		firstItemIndex = int(math.Ceil(float64(l.list.scroller.Offset.Y) / float64(l.list.itemMin.Height)))
	}
	lastItemIndex := int(math.Ceil(float64(l.list.scroller.Offset.Y+size.Height) / float64(l.list.itemMin.Height)))
	if lastItemIndex >= l.list.Length {
		lastItemIndex = l.list.Length - 1
	}
	visibleItemCount := int(math.Ceil(float64(size.Height) / float64(l.list.itemMin.Height)))
	if lastItemIndex-firstItemIndex < visibleItemCount {
		visibleItemCount = lastItemIndex - firstItemIndex
	}
	if len(l.itemCache) < visibleItemCount && len(l.itemCache) < l.list.Length {
		for i := len(l.itemCache); i < visibleItemCount && i < l.list.Length; i++ {
			item := newListItem(l.list.OnNewItem(), nil)
			l.itemCache = append(l.itemCache, item)
		}
	} else if len(l.itemCache) > l.list.Length {
		for i := l.list.Length - 1; i < len(l.itemCache); i++ {
			l.itemCache[i].Hide()
		}
		l.itemCache = l.itemCache[:l.list.Length]
	} else if len(l.itemCache) > visibleItemCount {
		for i := visibleItemCount + 1; i < len(l.itemCache); i++ {
			l.itemCache[i].Hide()
		}
		l.itemCache = l.itemCache[:visibleItemCount+1]
	}
	l.list.listLayout.Objects = l.itemCache
	l.list.listLayout.Layout.Layout(l.list.listLayout.Objects, l.list.listLayout.Layout.MinSize(nil))

	j := 0
	for i := firstItemIndex; i <= lastItemIndex && j < len(l.itemCache); i++ {
		item := l.itemCache[j].(*listItem)
		l.list.OnUpdateItem(i, item.child)
		index := i
		item.onTapped = func() {
			if l.list.selected != nil {
				l.list.selected.selected = false
				l.list.selected.Refresh()
			}
			l.list.selected = item
			l.list.selected.selected = true
			l.list.selected.Refresh()
			if l.list.OnItemSelected != nil {
				l.list.OnItemSelected(index, item.child)
			}
		}
		item.Refresh()
		j++
	}
}

func (l *listRenderer) BackgroundColor() color.Color {
	if l.list.background == nil {
		return theme.BackgroundColor()
	}

	return l.list.background
}

func (l *listRenderer) Refresh() {
	l.Layout(l.list.scroller.Size())
	canvas.Refresh(l.list.super())
}

// List is a widget that caches list items for performance and
// lays the items out in a vertical direction inside of a scroller
type List struct {
	BaseWidget
	background color.Color

	Length         int
	OnNewItem      func() fyne.CanvasObject
	OnUpdateItem   func(index int, item fyne.CanvasObject)
	OnItemSelected func(index int, item fyne.CanvasObject)

	listLayout *fyne.Container
	scroller   *ScrollContainer
	selected   *listItem
	itemMin    fyne.Size
}

// Refresh updates list to match the current theme
func (l *List) Refresh() {
	if l.background != nil {
		l.background = theme.BackgroundColor()
	}

	l.BaseWidget.Refresh()
}

// MinSize returns the size that this widget should not shrink below
func (l *List) MinSize() fyne.Size {
	l.ExtendBaseWidget(l)

	return l.BaseWidget.MinSize()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (l *List) CreateRenderer() fyne.WidgetRenderer {
	l.ExtendBaseWidget(l)

	if l.itemMin.Width == 0 && l.itemMin.Height == 0 && l.OnNewItem != nil {
		l.itemMin = newListItem(l.OnNewItem(), nil).MinSize()
	}
	l.listLayout = fyne.NewContainerWithLayout(newListLayout(l.itemMin, l.Length))
	l.scroller = NewVScrollContainer(l.listLayout)

	l.scroller.onOffsetChanged = func() {
		ll := l.listLayout.Layout.(*listLayout)
		if ll.offsetX == l.scroller.Offset.X && ll.offsetY == l.scroller.Offset.Y {
			return
		}
		ll.offsetX = l.scroller.Offset.X
		ll.offsetY = l.scroller.Offset.Y
		l.listLayout.Layout.Layout(l.listLayout.Objects, l.listLayout.Layout.MinSize(nil))
		l.BaseWidget.Refresh()
	}

	objects := []fyne.CanvasObject{l.listLayout, l.scroller}
	return &listRenderer{BaseRenderer: widget.NewBaseRenderer(objects), list: l}
}

// NewList creates and returns a list widget for displaying items in
// a vertical layout with scrolling and caching for performance
func NewList() *List {
	return &List{BaseWidget: BaseWidget{}}
}
