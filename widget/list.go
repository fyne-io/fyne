package widget

import (
	"math"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/theme"
)

// List is a widget that pools list items for performance and
// lays the items out in a vertical direction inside of a scroller.
// List requires that all items are the same size
type List struct {
	BaseWidget

	Length         func() int
	CreateItem     func() fyne.CanvasObject
	UpdateItem     func(index int, item fyne.CanvasObject)
	OnItemSelected func(index int, item fyne.CanvasObject)

	offsetY         int
	previousOffsetY int
	selectedItem    *listItem
	selectedIndex   int
	itemMin         fyne.Size
}

// NewList creates and returns a list widget for displaying items in
// a vertical layout with scrolling and caching for performance
func NewList(length func() int, createItem func() fyne.CanvasObject, updateItem func(index int, item fyne.CanvasObject)) *List {
	list := &List{BaseWidget: BaseWidget{}, Length: length, CreateItem: createItem, UpdateItem: updateItem, selectedIndex: -1}
	list.ExtendBaseWidget(list)
	return list
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (l *List) CreateRenderer() fyne.WidgetRenderer {
	l.ExtendBaseWidget(l)

	if l.itemMin.IsZero() && l.CreateItem != nil {
		l.itemMin = newListItem(l.CreateItem(), nil).MinSize()
	}
	layout := fyne.NewContainerWithLayout(newListLayout(l))
	layout.Resize(layout.MinSize())

	scroller := NewVScrollContainer(layout)

	objects := []fyne.CanvasObject{scroller}
	return newListRenderer(objects, l, scroller, layout)
}

// MinSize returns the size that this widget should not shrink below
func (l *List) MinSize() fyne.Size {
	l.ExtendBaseWidget(l)

	return l.BaseWidget.MinSize()
}

type listRenderer struct {
	widget.BaseRenderer
	list             *List
	scroller         *ScrollContainer
	layout           *fyne.Container
	itemPool         *listPool
	children         []fyne.CanvasObject
	visibleItemCount int
	firstItemIndex   int
	lastItemIndex    int
	size             fyne.Size
}

func newListRenderer(objects []fyne.CanvasObject, l *List, scroller *ScrollContainer, layout *fyne.Container) *listRenderer {
	lr := &listRenderer{BaseRenderer: widget.NewBaseRenderer(objects), list: l, scroller: scroller, layout: layout}
	lr.scroller.onOffsetChanged = func() {
		if l.offsetY == lr.scroller.Offset.Y {
			return
		}
		l.offsetY = lr.scroller.Offset.Y
		lr.offsetChanged()
		//l.BaseWidget.Refresh()
	}
	return lr
}

func (l *listRenderer) Layout(size fyne.Size) {
	if size != l.size {
		if size.Width != l.size.Width {
			for _, child := range l.children {
				child.Resize(fyne.NewSize(size.Width, l.list.itemMin.Height))
			}
		}
		l.scroller.Resize(size)
		l.size = size
	}
	if l.list.Length() == 0 {
		return
	}
	if l.itemPool == nil {
		l.itemPool = &listPool{}
	}

	// Relayout What Is Visible - no scroll change - initial layout or possibly from a resize
	l.visibleItemCount = int(math.Ceil(float64(l.scroller.size.Height) / float64(l.list.itemMin.Height)))
	if len(l.children) >= l.visibleItemCount || len(l.children) >= l.list.Length() {
		return
	}
	for i := len(l.children) + l.firstItemIndex; len(l.children) <= l.visibleItemCount && i < l.list.Length(); i++ {
		l.appendItem(i)
	}
	l.layout.Objects = l.children
	l.layout.Layout.Layout(l.layout.Objects, l.list.itemMin)
	l.lastItemIndex = l.firstItemIndex + len(l.children) - 1

	poolCapacity := 5
	if l.list.Length()-len(l.children) < 5 {
		poolCapacity = l.list.Length() - len(l.children)
	}
	for i := l.itemPool.Count(); i < poolCapacity; i++ {
		// Make sure to keep extra items in the pool
		item := newListItem(l.list.CreateItem(), nil)
		l.itemPool.Release(item)
	}
	i := l.firstItemIndex
	for _, child := range l.children {
		l.list.UpdateItem(i, child.(*listItem).child)
		i++
	}
}

func (l *listRenderer) MinSize() fyne.Size {
	return l.scroller.MinSize()
}

func (l *listRenderer) Refresh() {
	if l.list.CreateItem != nil {
		l.list.itemMin = newListItem(l.list.CreateItem(), nil).MinSize()
	}
	l.Layout(l.list.Size())
	canvas.Refresh(l.list.super())
}

func (l *listRenderer) appendItem(index int) {
	item := l.getItem()
	l.children = append(l.children, item)
	l.setupListItem(item, index)
	l.layout.Objects = l.children
	l.layout.Layout.(*listLayout).itemAppended(l.layout.Objects)
}

func (l *listRenderer) getItem() fyne.CanvasObject {
	item := l.itemPool.Obtain()
	if item == nil {
		item = newListItem(l.list.CreateItem(), nil)
	}
	return item
}

func (l *listRenderer) offsetChanged() {
	offsetChange := int(math.Abs(float64(l.list.previousOffsetY - l.list.offsetY)))

	if l.list.previousOffsetY < l.list.offsetY {
		// Scrolling Down
		l.scrollDown(offsetChange)
		return

	} else if l.list.previousOffsetY > l.list.offsetY {
		// Scrolling Up
		l.scrollUp(offsetChange)
		return
	}
}

func (l *listRenderer) prependItem(index int) {
	item := l.getItem()
	l.children = append([]fyne.CanvasObject{item}, l.children...)
	l.setupListItem(item, index)
	l.layout.Objects = l.children
	l.layout.Layout.(*listLayout).itemPrepended(l.layout.Objects)
}

func (l *listRenderer) scrollDown(offsetChange int) {
	itemChange := 0
	layoutEndY := l.children[len(l.children)-1].Position().Y + l.list.itemMin.Height
	scrollerEndY := l.scroller.Offset.Y + l.scroller.Size().Height
	if layoutEndY < scrollerEndY {
		itemChange = int(math.Ceil(float64(scrollerEndY-layoutEndY) / float64(l.list.itemMin.Height)))
	} else if offsetChange < l.list.itemMin.Height {
		return
	} else {
		itemChange = int(math.Floor(float64(offsetChange) / float64(l.list.itemMin.Height)))
	}
	l.list.previousOffsetY = l.list.offsetY
	for i := 0; i < itemChange && l.lastItemIndex != l.list.Length()-1; i++ {
		l.itemPool.Release(l.children[0])
		l.children = l.children[1:]
		l.firstItemIndex++
		l.lastItemIndex++
		l.appendItem(l.lastItemIndex)
	}
}

func (l *listRenderer) scrollUp(offsetChange int) {
	itemChange := 0
	layoutStartY := l.children[0].Position().Y
	if layoutStartY > l.scroller.Offset.Y {
		itemChange = int(math.Ceil(float64(layoutStartY-l.scroller.Offset.Y) / float64(l.list.itemMin.Height)))
	} else if offsetChange < l.list.itemMin.Height {
		return
	} else {
		itemChange = int(math.Floor(float64(offsetChange) / float64(l.list.itemMin.Height)))
	}
	l.list.previousOffsetY = l.list.offsetY
	for i := 0; i < itemChange && l.firstItemIndex != 0; i++ {
		l.itemPool.Release(l.children[len(l.children)-1])
		l.children = l.children[:len(l.children)-1]
		l.firstItemIndex--
		l.lastItemIndex--
		l.prependItem(l.firstItemIndex)
	}
}

func (l *listRenderer) setupListItem(item fyne.CanvasObject, index int) {
	if index != l.list.selectedIndex {
		item.(*listItem).selected = false
	} else {
		item.(*listItem).selected = true
		l.list.selectedItem = item.(*listItem)
	}
	l.list.UpdateItem(index, item.(*listItem).child)
	item.(*listItem).onTapped = func() {
		if l.list.selectedItem != nil && l.list.selectedIndex >= l.firstItemIndex && l.list.selectedIndex <= l.lastItemIndex {
			l.list.selectedItem.selected = false
			l.list.selectedItem.Refresh()
		}
		l.list.selectedItem = item.(*listItem)
		l.list.selectedIndex = index
		l.list.selectedItem.selected = true
		l.list.selectedItem.Refresh()
		if l.list.OnItemSelected != nil {
			l.list.OnItemSelected(index, item.(*listItem).child)
		}
	}
}

type listItem struct {
	BaseWidget

	onTapped          func()
	statusIndicator   *canvas.Rectangle
	child             fyne.CanvasObject
	divider           *canvas.Rectangle
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
	li.divider = canvas.NewRectangle(theme.ShadowColor())

	objects := []fyne.CanvasObject{li.statusIndicator, li.child, li.divider}

	return &listItemRenderer{widget.NewBaseRenderer(objects), li}
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
	if li.onTapped != nil {
		li.selected = true
		li.Refresh()
		li.onTapped()
	}
}

type listItemRenderer struct {
	widget.BaseRenderer

	item *listItem
}

// MinSize calculates the minimum size of a listItem.
// This is based on the size of the status indicator and the size of the child object.
func (li *listItemRenderer) MinSize() (size fyne.Size) {
	size = fyne.NewSize((theme.Padding()*2)+li.item.child.Size().Width,
		li.item.child.Size().Height+theme.Padding()*2)
	return
}

// Layout the components of the listItem widget
func (li *listItemRenderer) Layout(size fyne.Size) {
	li.item.statusIndicator.Move(fyne.NewPos(0, 0))
	s := fyne.NewSize(theme.Padding(), size.Height-1)
	li.item.statusIndicator.SetMinSize(s)
	li.item.statusIndicator.Resize(s)

	x := theme.Padding() * 2
	li.item.child.Move(fyne.NewPos(x, theme.Padding()))
	li.item.child.Resize(fyne.NewSize(size.Width-x, size.Height-(2*theme.Padding())))

	li.item.divider.Move(fyne.NewPos(0, size.Height-1))
	s = fyne.NewSize(size.Width, 1)
	li.item.divider.SetMinSize(s)
	li.item.divider.Resize(s)
}

func (li *listItemRenderer) Refresh() {
	if li.item.selected {
		li.item.statusIndicator.FillColor = theme.FocusColor()
		canvas.Refresh(li.item.super())
		return
	}
	if li.item.hovered {
		li.item.statusIndicator.FillColor = theme.HoverColor()
	} else {
		li.item.statusIndicator.FillColor = theme.BackgroundColor()
	}
	canvas.Refresh(li.item.super())
}

type listLayout struct {
	list      *List
	itemCount int
}

var _ fyne.Layout = (*listLayout)(nil)

func newListLayout(list *List) fyne.Layout {
	return &listLayout{list: list}
}

func (l *listLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if l.list.offsetY != 0 {
		return
	}
	y := 0
	for _, child := range objects {
		child.Move(fyne.NewPos(theme.Padding(), y))
		y += l.list.itemMin.Height
		child.Resize(fyne.NewSize(l.list.size.Width-theme.Padding()*2, l.list.itemMin.Height))
	}
}

func (l *listLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {

	return fyne.NewSize(l.list.itemMin.Width+theme.Padding(),
		l.list.itemMin.Height*l.list.Length())
}

func (l *listLayout) itemAppended(objects []fyne.CanvasObject) {
	if len(objects) > 1 {
		objects[len(objects)-1].Move(fyne.NewPos(theme.Padding(), objects[len(objects)-2].Position().Y+l.list.itemMin.Height))
	} else {
		objects[len(objects)-1].Move(fyne.NewPos(theme.Padding(), 0))
	}
	objects[len(objects)-1].Resize(fyne.NewSize(l.list.size.Width-theme.Padding()*2, l.list.itemMin.Height))
}

func (l *listLayout) itemPrepended(objects []fyne.CanvasObject) {
	objects[0].Move(fyne.NewPos(theme.Padding(), objects[1].Position().Y-l.list.itemMin.Height))
	objects[0].Resize(fyne.NewSize(l.list.size.Width-theme.Padding()*2, l.list.itemMin.Height))
}

type pool interface {
	Count() int
	Obtain() fyne.CanvasObject
	Release(fyne.CanvasObject)
}

type listPool struct {
	contents []fyne.CanvasObject
}

// Count returns the number of available items in the pool
func (lc *listPool) Count() int {
	return len(lc.contents)
}

// Obtain returns an item from the pool for use
func (lc *listPool) Obtain() fyne.CanvasObject {
	if len(lc.contents) == 0 {
		return nil
	}
	item := lc.contents[0]
	lc.contents = lc.contents[1:]
	return item
}

// Release adds an item into the pool to be used later
func (lc *listPool) Release(item fyne.CanvasObject) {
	lc.contents = append(lc.contents, item)
}
