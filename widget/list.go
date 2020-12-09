package widget

import (
	"math"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/theme"
)

// ListItemID uniquely identifies an item within a list.
type ListItemID = int

const listDividerHeight = 1

// Declare conformity with Widget interface.
var _ fyne.Widget = (*List)(nil)

// List is a widget that pools list items for performance and
// lays the items out in a vertical direction inside of a scroller.
// List requires that all items are the same size.
//
// Since: 1.4
type List struct {
	BaseWidget

	Length       func() int
	CreateItem   func() fyne.CanvasObject
	UpdateItem   func(id ListItemID, item fyne.CanvasObject)
	OnSelected   func(id ListItemID)
	OnUnselected func(id ListItemID)

	scroller *ScrollContainer
	selected []ListItemID
	itemMin  fyne.Size
	offsetY  int
}

// NewList creates and returns a list widget for displaying items in
// a vertical layout with scrolling and caching for performance.
//
// Since: 1.4
func NewList(length func() int, createItem func() fyne.CanvasObject, updateItem func(ListItemID, fyne.CanvasObject)) *List {
	list := &List{BaseWidget: BaseWidget{}, Length: length, CreateItem: createItem, UpdateItem: updateItem}
	list.ExtendBaseWidget(list)
	return list
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer.
func (l *List) CreateRenderer() fyne.WidgetRenderer {
	l.ExtendBaseWidget(l)

	if f := l.CreateItem; f != nil {
		if l.itemMin.IsZero() {
			l.itemMin = newListItem(f(), nil).MinSize()
		}
	}
	layout := fyne.NewContainerWithLayout(newListLayout(l))
	layout.Resize(layout.MinSize())
	l.scroller = NewVScrollContainer(layout)
	objects := []fyne.CanvasObject{l.scroller}
	return newListRenderer(objects, l, l.scroller, layout)
}

// MinSize returns the size that this widget should not shrink below.
func (l *List) MinSize() fyne.Size {
	l.ExtendBaseWidget(l)

	return l.BaseWidget.MinSize()
}

// Select add the item identified by the given ID to the selection.
func (l *List) Select(id ListItemID) {
	if len(l.selected) > 0 && id == l.selected[0] {
		return
	}
	length := 0
	if f := l.Length; f != nil {
		length = f()
	}
	if id < 0 || id >= length {
		return
	}
	old := l.selected
	l.selected = []ListItemID{id}
	defer func() {
		if f := l.OnUnselected; f != nil && len(old) > 0 {
			f(old[0])
		}
		if f := l.OnSelected; f != nil {
			f(id)
		}
	}()
	if l.scroller == nil {
		return
	}
	y := (id * l.itemMin.Height) + (id * listDividerHeight)
	if y < l.scroller.Offset.Y {
		l.scroller.Offset.Y = y
	} else if y+l.itemMin.Height > l.scroller.Offset.Y+l.scroller.Size().Height {
		l.scroller.Offset.Y = y + l.itemMin.Height - l.scroller.Size().Height
	}
	l.scroller.onOffsetChanged()
	l.Refresh()
}

// Unselect removes the item identified by the given ID from the selection.
func (l *List) Unselect(id ListItemID) {
	if len(l.selected) == 0 {
		return
	}

	l.selected = nil
	l.Refresh()
	if f := l.OnUnselected; f != nil {
		f(id)
	}
}

// Declare conformity with WidgetRenderer interface.
var _ fyne.WidgetRenderer = (*listRenderer)(nil)

type listRenderer struct {
	widget.BaseRenderer

	list             *List
	scroller         *ScrollContainer
	layout           *fyne.Container
	itemPool         *syncPool
	children         []fyne.CanvasObject
	size             fyne.Size
	visibleItemCount int
	firstItemIndex   ListItemID
	lastItemIndex    ListItemID
	previousOffsetY  int
}

func newListRenderer(objects []fyne.CanvasObject, l *List, scroller *ScrollContainer, layout *fyne.Container) *listRenderer {
	lr := &listRenderer{BaseRenderer: widget.NewBaseRenderer(objects), list: l, scroller: scroller, layout: layout}
	lr.scroller.onOffsetChanged = func() {
		if lr.list.offsetY == lr.scroller.Offset.Y {
			return
		}
		lr.list.offsetY = lr.scroller.Offset.Y
		lr.offsetChanged()
	}
	return lr
}

func (l *listRenderer) Layout(size fyne.Size) {
	length := 0
	if f := l.list.Length; f != nil {
		length = f()
	}
	if length <= 0 {
		if len(l.children) > 0 {
			for _, child := range l.children {
				l.itemPool.Release(child)
			}
			l.previousOffsetY = 0
			l.firstItemIndex = 0
			l.lastItemIndex = 0
			l.visibleItemCount = 0
			l.list.offsetY = 0
			l.layout.Layout.(*listLayout).layoutEndY = 0
			l.children = nil
			l.layout.Objects = nil
			l.list.Refresh()
		}
		return
	}
	if size != l.size {
		if size.Width != l.size.Width {
			for _, child := range l.children {
				child.Resize(fyne.NewSize(size.Width, l.list.itemMin.Height))
			}
		}
		l.scroller.Resize(size)
		l.size = size
	}
	if l.itemPool == nil {
		l.itemPool = &syncPool{}
	}

	// Relayout What Is Visible - no scroll change - initial layout or possibly from a resize.
	l.visibleItemCount = int(math.Ceil(float64(l.scroller.size.Height) / float64(l.list.itemMin.Height+listDividerHeight)))
	if l.visibleItemCount <= 0 {
		return
	}
	min := fyne.Min(length, l.visibleItemCount)
	if len(l.children) > min {
		for i := len(l.children); i >= min; i-- {
			l.itemPool.Release(l.children[i-1])
		}
		l.children = l.children[:min-1]
	}
	for i := len(l.children) + l.firstItemIndex; len(l.children) <= l.visibleItemCount && i < length; i++ {
		l.appendItem(i)
	}
	l.layout.Layout.(*listLayout).children = l.children
	l.layout.Layout.Layout(l.children, l.list.itemMin)
	l.layout.Objects = l.layout.Layout.(*listLayout).getObjects()
	l.lastItemIndex = l.firstItemIndex + len(l.children) - 1

	i := l.firstItemIndex
	for _, child := range l.children {
		if f := l.list.UpdateItem; f != nil {
			f(i, child.(*listItem).child)
		}
		l.setupListItem(child, i)
		i++
	}
}

func (l *listRenderer) MinSize() fyne.Size {
	return l.scroller.MinSize().Max(l.list.itemMin)
}

func (l *listRenderer) Refresh() {
	if f := l.list.CreateItem; f != nil {
		l.list.itemMin = newListItem(f(), nil).MinSize()
	}
	l.Layout(l.list.Size())
	l.scroller.Refresh()
	canvas.Refresh(l.list.super())
}

func (l *listRenderer) appendItem(id ListItemID) {
	item := l.getItem()
	l.children = append(l.children, item)
	l.setupListItem(item, id)
	l.layout.Layout.(*listLayout).children = l.children
	l.layout.Layout.(*listLayout).appendedItem(l.children)
	l.layout.Objects = l.layout.Layout.(*listLayout).getObjects()
}

func (l *listRenderer) getItem() fyne.CanvasObject {
	item := l.itemPool.Obtain()
	if item == nil {
		if f := l.list.CreateItem; f != nil {
			item = newListItem(f(), nil)
		}
	}
	return item
}

func (l *listRenderer) offsetChanged() {
	offsetChange := int(math.Abs(float64(l.previousOffsetY - l.list.offsetY)))

	if l.previousOffsetY < l.list.offsetY {
		// Scrolling Down.
		l.scrollDown(offsetChange)
	} else if l.previousOffsetY > l.list.offsetY {
		// Scrolling Up.
		l.scrollUp(offsetChange)
	}
	l.layout.Layout.(*listLayout).updateDividers()
}

func (l *listRenderer) prependItem(id ListItemID) {
	item := l.getItem()
	l.children = append([]fyne.CanvasObject{item}, l.children...)
	l.setupListItem(item, id)
	l.layout.Layout.(*listLayout).children = l.children
	l.layout.Layout.(*listLayout).prependedItem(l.children)
	l.layout.Objects = l.layout.Layout.(*listLayout).getObjects()
}

func (l *listRenderer) scrollDown(offsetChange int) {
	itemChange := 0
	layoutEndY := l.children[len(l.children)-1].Position().Y + l.list.itemMin.Height + listDividerHeight
	scrollerEndY := l.scroller.Offset.Y + l.scroller.Size().Height
	if layoutEndY < scrollerEndY {
		itemChange = int(math.Ceil(float64(scrollerEndY-layoutEndY) / float64(l.list.itemMin.Height+listDividerHeight)))
	} else if offsetChange < l.list.itemMin.Height+listDividerHeight {
		return
	} else {
		itemChange = int(math.Floor(float64(offsetChange) / float64(l.list.itemMin.Height+listDividerHeight)))
	}
	l.previousOffsetY = l.list.offsetY
	length := 0
	if f := l.list.Length; f != nil {
		length = f()
	}
	if length == 0 {
		return
	}
	for i := 0; i < itemChange && l.lastItemIndex != length-1; i++ {
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
		itemChange = int(math.Ceil(float64(layoutStartY-l.scroller.Offset.Y) / float64(l.list.itemMin.Height+listDividerHeight)))
	} else if offsetChange < l.list.itemMin.Height+listDividerHeight {
		return
	} else {
		itemChange = int(math.Floor(float64(offsetChange) / float64(l.list.itemMin.Height+listDividerHeight)))
	}
	l.previousOffsetY = l.list.offsetY
	for i := 0; i < itemChange && l.firstItemIndex != 0; i++ {
		l.itemPool.Release(l.children[len(l.children)-1])
		l.children = l.children[:len(l.children)-1]
		l.firstItemIndex--
		l.lastItemIndex--
		l.prependItem(l.firstItemIndex)
	}
}

func (l *listRenderer) setupListItem(item fyne.CanvasObject, id ListItemID) {
	li := item.(*listItem)
	previousIndicator := li.selected
	li.selected = false
	for _, s := range l.list.selected {
		if id == s {
			li.selected = true
		}
	}
	if previousIndicator != li.selected {
		item.Refresh()
	}
	if f := l.list.UpdateItem; f != nil {
		f(id, li.child)
	}
	li.onTapped = func() {
		l.list.Select(id)
	}
}

// Declare conformity with interfaces.
var _ fyne.Widget = (*listItem)(nil)
var _ fyne.Tappable = (*listItem)(nil)
var _ desktop.Hoverable = (*listItem)(nil)

type listItem struct {
	BaseWidget

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

// CreateRenderer is a private method to Fyne which links this widget to its renderer.
func (li *listItem) CreateRenderer() fyne.WidgetRenderer {
	li.ExtendBaseWidget(li)

	li.statusIndicator = canvas.NewRectangle(theme.BackgroundColor())

	objects := []fyne.CanvasObject{li.statusIndicator, li.child}

	return &listItemRenderer{widget.NewBaseRenderer(objects), li}
}

// MinSize returns the size that this widget should not shrink below.
func (li *listItem) MinSize() fyne.Size {
	li.ExtendBaseWidget(li)
	return li.BaseWidget.MinSize()
}

// MouseIn is called when a desktop pointer enters the widget.
func (li *listItem) MouseIn(*desktop.MouseEvent) {
	li.hovered = true
	li.Refresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget.
func (li *listItem) MouseMoved(*desktop.MouseEvent) {
}

// MouseOut is called when a desktop pointer exits the widget.
func (li *listItem) MouseOut() {
	li.hovered = false
	li.Refresh()
}

// Tapped is called when a pointer tapped event is captured and triggers any tap handler.
func (li *listItem) Tapped(*fyne.PointEvent) {
	if li.onTapped != nil {
		li.selected = true
		li.Refresh()
		li.onTapped()
	}
}

// Declare conformity with the WidgetRenderer interface.
var _ fyne.WidgetRenderer = (*listItemRenderer)(nil)

type listItemRenderer struct {
	widget.BaseRenderer

	item *listItem
}

// MinSize calculates the minimum size of a listItem.
// This is based on the size of the status indicator and the size of the child object.
func (li *listItemRenderer) MinSize() (size fyne.Size) {
	itemSize := li.item.child.MinSize()
	size = fyne.NewSize(itemSize.Width+theme.Padding()*3,
		itemSize.Height+theme.Padding()*2)
	return
}

// Layout the components of the listItem widget.
func (li *listItemRenderer) Layout(size fyne.Size) {
	li.item.statusIndicator.Move(fyne.NewPos(0, 0))
	s := fyne.NewSize(theme.Padding(), size.Height)
	li.item.statusIndicator.SetMinSize(s)
	li.item.statusIndicator.Resize(s)

	li.item.child.Move(fyne.NewPos(theme.Padding()*2, theme.Padding()))
	li.item.child.Resize(fyne.NewSize(size.Width-theme.Padding()*3, size.Height-theme.Padding()*2))
}

func (li *listItemRenderer) Refresh() {
	if li.item.selected {
		li.item.statusIndicator.FillColor = theme.FocusColor()
	} else if li.item.hovered {
		li.item.statusIndicator.FillColor = theme.HoverColor()
	} else {
		li.item.statusIndicator.FillColor = theme.BackgroundColor()
	}
	canvas.Refresh(li.item.super())
}

// Declare conformity with Layout interface.
var _ fyne.Layout = (*listLayout)(nil)

type listLayout struct {
	list       *List
	dividers   []fyne.CanvasObject
	children   []fyne.CanvasObject
	layoutEndY int
}

func newListLayout(list *List) fyne.Layout {
	return &listLayout{list: list}
}

func (l *listLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if l.list.offsetY != 0 {
		return
	}
	y := 0
	for _, child := range l.children {
		child.Move(fyne.NewPos(0, y))
		y += l.list.itemMin.Height + listDividerHeight
		child.Resize(fyne.NewSize(l.list.size.Width, l.list.itemMin.Height))
	}
	l.layoutEndY = y
	l.updateDividers()
}

func (l *listLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if f := l.list.Length; f != nil {
		return fyne.NewSize(l.list.itemMin.Width,
			(l.list.itemMin.Height+listDividerHeight)*f()-listDividerHeight)
	}
	return fyne.NewSize(0, 0)
}

func (l *listLayout) getObjects() []fyne.CanvasObject {
	objects := l.children
	objects = append(objects, l.dividers...)
	return objects
}

func (l *listLayout) appendedItem(objects []fyne.CanvasObject) {
	if len(objects) > 1 {
		objects[len(objects)-1].Move(fyne.NewPos(0, objects[len(objects)-2].Position().Y+l.list.itemMin.Height+listDividerHeight))
	} else {
		objects[len(objects)-1].Move(fyne.NewPos(0, 0))
	}
	objects[len(objects)-1].Resize(fyne.NewSize(l.list.size.Width, l.list.itemMin.Height))
}

func (l *listLayout) prependedItem(objects []fyne.CanvasObject) {
	objects[0].Move(fyne.NewPos(0, objects[1].Position().Y-l.list.itemMin.Height-listDividerHeight))
	objects[0].Resize(fyne.NewSize(l.list.size.Width, l.list.itemMin.Height))
}

func (l *listLayout) updateDividers() {
	if len(l.children) > 1 {
		if len(l.dividers) > len(l.children) {
			l.dividers = l.dividers[:len(l.children)]
		} else {
			for i := len(l.dividers); i < len(l.children); i++ {
				l.dividers = append(l.dividers, NewSeparator())
			}
		}
	} else {
		l.dividers = nil
	}
	for i, child := range l.children {
		if i == 0 {
			continue
		}
		l.dividers[i].Move(fyne.NewPos(theme.Padding(), child.Position().Y-listDividerHeight))
		l.dividers[i].Resize(fyne.NewSize(l.list.Size().Width-(theme.Padding()*2), listDividerHeight))
		l.dividers[i].Show()
	}
}
