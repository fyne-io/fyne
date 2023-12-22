package widget

import (
	"fmt"
	"math"
	"sort"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
)

// ListItemID uniquely identifies an item within a list.
type ListItemID = int

// Declare conformity with interfaces.
var _ fyne.Widget = (*List)(nil)
var _ fyne.Focusable = (*List)(nil)

// List is a widget that pools list items for performance and
// lays the items out in a vertical direction inside of a scroller.
// List requires that all items are the same size.
//
// Since: 1.4
type List struct {
	BaseWidget

	Length       func() int                                  `json:"-"`
	CreateItem   func() fyne.CanvasObject                    `json:"-"`
	UpdateItem   func(id ListItemID, item fyne.CanvasObject) `json:"-"`
	OnSelected   func(id ListItemID)                         `json:"-"`
	OnUnselected func(id ListItemID)                         `json:"-"`

	currentFocus  ListItemID
	focused       bool
	scroller      *widget.Scroll
	selected      []ListItemID
	itemMin       fyne.Size
	itemHeights   map[ListItemID]float32
	offsetY       float32
	offsetUpdated func(fyne.Position)
}

// NewList creates and returns a list widget for displaying items in
// a vertical layout with scrolling and caching for performance.
//
// Since: 1.4
func NewList(length func() int, createItem func() fyne.CanvasObject, updateItem func(ListItemID, fyne.CanvasObject)) *List {
	list := &List{Length: length, CreateItem: createItem, UpdateItem: updateItem}
	list.ExtendBaseWidget(list)
	return list
}

// NewListWithData creates a new list widget that will display the contents of the provided data.
//
// Since: 2.0
func NewListWithData(data binding.DataList, createItem func() fyne.CanvasObject, updateItem func(binding.DataItem, fyne.CanvasObject)) *List {
	l := NewList(
		data.Length,
		createItem,
		func(i ListItemID, o fyne.CanvasObject) {
			item, err := data.GetItem(i)
			if err != nil {
				fyne.LogError(fmt.Sprintf("Error getting data item %d", i), err)
				return
			}
			updateItem(item, o)
		})

	data.AddListener(binding.NewDataListener(l.Refresh))
	return l
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer.
func (l *List) CreateRenderer() fyne.WidgetRenderer {
	l.ExtendBaseWidget(l)

	if f := l.CreateItem; f != nil && l.itemMin.IsZero() {
		l.itemMin = f().MinSize()
	}

	layout := &fyne.Container{Layout: newListLayout(l)}
	l.scroller = widget.NewVScroll(layout)
	layout.Resize(layout.MinSize())
	objects := []fyne.CanvasObject{l.scroller}
	return newListRenderer(objects, l, l.scroller, layout)
}

// FocusGained is called after this List has gained focus.
//
// Implements: fyne.Focusable
func (l *List) FocusGained() {
	l.focused = true
	l.scrollTo(l.currentFocus)
	l.RefreshItem(l.currentFocus)
}

// FocusLost is called after this List has lost focus.
//
// Implements: fyne.Focusable
func (l *List) FocusLost() {
	l.focused = false
	l.RefreshItem(l.currentFocus)
}

// MinSize returns the size that this widget should not shrink below.
func (l *List) MinSize() fyne.Size {
	l.ExtendBaseWidget(l)

	return l.BaseWidget.MinSize()
}

// RefreshItem refreshes a single item, specified by the item ID passed in.
//
// Since: 2.4
func (l *List) RefreshItem(id ListItemID) {
	if l.scroller == nil {
		return
	}
	l.BaseWidget.Refresh()
	lo := l.scroller.Content.(*fyne.Container).Layout.(*listLayout)
	lo.renderLock.RLock() // ensures we are not changing visible info in render code during the search
	item, ok := lo.searchVisible(lo.visible, id)
	lo.renderLock.RUnlock()
	if ok {
		lo.setupListItem(item, id, l.focused && l.currentFocus == id)
	}
}

// SetItemHeight supports changing the height of the specified list item. Items normally take the height of the template
// returned from the CreateItem callback. The height parameter uses the same units as a fyne.Size type and refers
// to the internal content height not including the divider size.
//
// Since: 2.3
func (l *List) SetItemHeight(id ListItemID, height float32) {
	l.propertyLock.Lock()

	if l.itemHeights == nil {
		l.itemHeights = make(map[ListItemID]float32)
	}

	refresh := l.itemHeights[id] != height
	l.itemHeights[id] = height
	l.propertyLock.Unlock()

	if refresh {
		l.RefreshItem(id)
	}
}

func (l *List) scrollTo(id ListItemID) {
	if l.scroller == nil {
		return
	}

	separatorThickness := theme.Padding()
	y := float32(0)
	lastItemHeight := l.itemMin.Height
	if l.itemHeights == nil || len(l.itemHeights) == 0 {
		y = (float32(id) * l.itemMin.Height) + (float32(id) * separatorThickness)
	} else {
		for i := 0; i < id; i++ {
			height := l.itemMin.Height
			if h, ok := l.itemHeights[i]; ok {
				height = h
			}

			y += height + separatorThickness
			lastItemHeight = height
		}
	}

	if y < l.scroller.Offset.Y {
		l.scroller.Offset.Y = y
	} else if y+l.itemMin.Height > l.scroller.Offset.Y+l.scroller.Size().Height {
		l.scroller.Offset.Y = y + lastItemHeight - l.scroller.Size().Height
	}
	l.offsetUpdated(l.scroller.Offset)
}

// Resize is called when this list should change size. We refresh to ensure invisible items are drawn.
func (l *List) Resize(s fyne.Size) {
	l.BaseWidget.Resize(s)
	if l.scroller == nil {
		return
	}

	l.offsetUpdated(l.scroller.Offset)
	l.scroller.Content.(*fyne.Container).Layout.(*listLayout).updateList(false)
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
	l.scrollTo(id)
	l.Refresh()
}

// ScrollTo scrolls to the item represented by id
//
// Since: 2.1
func (l *List) ScrollTo(id ListItemID) {
	length := 0
	if f := l.Length; f != nil {
		length = f()
	}
	if id < 0 || id >= length {
		return
	}
	l.scrollTo(id)
	l.Refresh()
}

// ScrollToBottom scrolls to the end of the list
//
// Since: 2.1
func (l *List) ScrollToBottom() {
	length := 0
	if f := l.Length; f != nil {
		length = f()
	}
	if length > 0 {
		length--
	}
	l.scrollTo(length)
	l.Refresh()
}

// ScrollToTop scrolls to the start of the list
//
// Since: 2.1
func (l *List) ScrollToTop() {
	l.scrollTo(0)
	l.Refresh()
}

// TypedKey is called if a key event happens while this List is focused.
//
// Implements: fyne.Focusable
func (l *List) TypedKey(event *fyne.KeyEvent) {
	switch event.Name {
	case fyne.KeySpace:
		l.Select(l.currentFocus)
	case fyne.KeyDown:
		if f := l.Length; f != nil && l.currentFocus >= f()-1 {
			return
		}
		l.RefreshItem(l.currentFocus)
		l.currentFocus++
		l.scrollTo(l.currentFocus)
		l.RefreshItem(l.currentFocus)
	case fyne.KeyUp:
		if l.currentFocus <= 0 {
			return
		}
		l.RefreshItem(l.currentFocus)
		l.currentFocus--
		l.scrollTo(l.currentFocus)
		l.RefreshItem(l.currentFocus)
	}
}

// TypedRune is called if a text event happens while this List is focused.
//
// Implements: fyne.Focusable
func (l *List) TypedRune(_ rune) {
	// intentionally left blank
}

// Unselect removes the item identified by the given ID from the selection.
func (l *List) Unselect(id ListItemID) {
	if len(l.selected) == 0 || l.selected[0] != id {
		return
	}

	l.selected = nil
	l.Refresh()
	if f := l.OnUnselected; f != nil {
		f(id)
	}
}

// UnselectAll removes all items from the selection.
//
// Since: 2.1
func (l *List) UnselectAll() {
	if len(l.selected) == 0 {
		return
	}

	selected := l.selected
	l.selected = nil
	l.Refresh()
	if f := l.OnUnselected; f != nil {
		for _, id := range selected {
			f(id)
		}
	}
}

// fills l.visibleRowHeights and also returns offY and minRow
func (l *listLayout) calculateVisibleRowHeights(itemHeight float32, length int) (offY float32, minRow int) {
	rowOffset := float32(0)
	isVisible := false
	l.visibleRowHeights = l.visibleRowHeights[:0]

	if l.list.scroller.Size().Height <= 0 {
		return
	}

	// theme.Padding is a slow call, so we cache it
	padding := theme.Padding()

	if len(l.list.itemHeights) == 0 {
		paddedItemHeight := itemHeight + padding

		offY = float32(math.Floor(float64(l.list.offsetY/paddedItemHeight))) * paddedItemHeight
		minRow = int(math.Floor(float64(offY / paddedItemHeight)))
		maxRow := int(math.Ceil(float64((offY + l.list.scroller.Size().Height) / paddedItemHeight)))

		if minRow > length-1 {
			minRow = length - 1
		}
		if minRow < 0 {
			minRow = 0
			offY = 0
		}

		if maxRow > length {
			maxRow = length
		}

		for i := 0; i < maxRow-minRow; i++ {
			l.visibleRowHeights = append(l.visibleRowHeights, itemHeight)
		}
		return
	}

	for i := 0; i < length; i++ {
		height := itemHeight
		if h, ok := l.list.itemHeights[i]; ok {
			height = h
		}

		if rowOffset <= l.list.offsetY-height-padding {
			// before scroll
		} else if rowOffset <= l.list.offsetY {
			minRow = i
			offY = rowOffset
			isVisible = true
		}
		if rowOffset >= l.list.offsetY+l.list.scroller.Size().Height {
			break
		}

		rowOffset += height + padding
		if isVisible {
			l.visibleRowHeights = append(l.visibleRowHeights, height)
		}
	}
	return
}

// Declare conformity with WidgetRenderer interface.
var _ fyne.WidgetRenderer = (*listRenderer)(nil)

type listRenderer struct {
	widget.BaseRenderer

	list     *List
	scroller *widget.Scroll
	layout   *fyne.Container
}

func newListRenderer(objects []fyne.CanvasObject, l *List, scroller *widget.Scroll, layout *fyne.Container) *listRenderer {
	lr := &listRenderer{BaseRenderer: widget.NewBaseRenderer(objects), list: l, scroller: scroller, layout: layout}
	lr.scroller.OnScrolled = l.offsetUpdated
	return lr
}

func (l *listRenderer) Layout(size fyne.Size) {
	l.scroller.Resize(size)
}

func (l *listRenderer) MinSize() fyne.Size {
	return l.scroller.MinSize().Max(l.list.itemMin)
}

func (l *listRenderer) Refresh() {
	if f := l.list.CreateItem; f != nil {
		l.list.itemMin = f().MinSize()
	}
	l.Layout(l.list.Size())
	l.scroller.Refresh()
	l.layout.Layout.(*listLayout).updateList(false)
	canvas.Refresh(l.list.super())
}

// Declare conformity with interfaces.
var _ fyne.Widget = (*listItem)(nil)
var _ fyne.Tappable = (*listItem)(nil)
var _ desktop.Hoverable = (*listItem)(nil)

type listItem struct {
	BaseWidget

	onTapped          func()
	background        *canvas.Rectangle
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

	li.background = canvas.NewRectangle(theme.HoverColor())
	li.background.CornerRadius = theme.SelectionRadiusSize()
	li.background.Hide()

	objects := []fyne.CanvasObject{li.background, li.child}

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
func (li *listItemRenderer) MinSize() fyne.Size {
	return li.item.child.MinSize()
}

// Layout the components of the listItem widget.
func (li *listItemRenderer) Layout(size fyne.Size) {
	li.item.background.Resize(size)
	li.item.child.Resize(size)
}

func (li *listItemRenderer) Refresh() {
	li.item.background.CornerRadius = theme.SelectionRadiusSize()
	if li.item.selected {
		li.item.background.FillColor = theme.SelectionColor()
		li.item.background.Show()
	} else if li.item.hovered {
		li.item.background.FillColor = theme.HoverColor()
		li.item.background.Show()
	} else {
		li.item.background.Hide()
	}
	li.item.background.Refresh()
	canvas.Refresh(li.item.super())
}

// Declare conformity with Layout interface.
var _ fyne.Layout = (*listLayout)(nil)

type listItemAndID struct {
	item *listItem
	id   ListItemID
}

type listLayout struct {
	list       *List
	separators []fyne.CanvasObject
	children   []fyne.CanvasObject

	itemPool          syncPool
	visible           []listItemAndID
	slicePool         sync.Pool // *[]itemAndID
	visibleRowHeights []float32
	renderLock        sync.RWMutex
}

func newListLayout(list *List) fyne.Layout {
	l := &listLayout{list: list}
	l.slicePool.New = func() interface{} {
		s := make([]listItemAndID, 0)
		return &s
	}
	list.offsetUpdated = l.offsetUpdated
	return l
}

func (l *listLayout) Layout([]fyne.CanvasObject, fyne.Size) {
	l.updateList(true)
}

func (l *listLayout) MinSize([]fyne.CanvasObject) fyne.Size {
	l.list.propertyLock.Lock()
	defer l.list.propertyLock.Unlock()
	items := 0
	if f := l.list.Length; f == nil {
		return fyne.NewSize(0, 0)
	} else {
		items = f()
	}

	separatorThickness := theme.Padding()
	if l.list.itemHeights == nil || len(l.list.itemHeights) == 0 {
		return fyne.NewSize(l.list.itemMin.Width,
			(l.list.itemMin.Height+separatorThickness)*float32(items)-separatorThickness)
	}

	height := float32(0)
	templateHeight := l.list.itemMin.Height
	for item := 0; item < items; item++ {
		itemHeight, ok := l.list.itemHeights[item]
		if ok {
			height += itemHeight
		} else {
			height += templateHeight
		}
	}

	return fyne.NewSize(l.list.itemMin.Width, height+separatorThickness*float32(items-1))
}

func (l *listLayout) getItem() *listItem {
	item := l.itemPool.Obtain()
	if item == nil {
		if f := l.list.CreateItem; f != nil {
			item = newListItem(f(), nil)
		}
	}
	return item.(*listItem)
}

func (l *listLayout) offsetUpdated(pos fyne.Position) {
	if l.list.offsetY == pos.Y {
		return
	}
	l.list.offsetY = pos.Y
	l.updateList(true)
}

func (l *listLayout) setupListItem(li *listItem, id ListItemID, focus bool) {
	previousIndicator := li.selected
	li.selected = false
	for _, s := range l.list.selected {
		if id == s {
			li.selected = true
			break
		}
	}
	if focus {
		li.hovered = true
		li.Refresh()
	} else if previousIndicator != li.selected || li.hovered {
		li.hovered = false
		li.Refresh()
	}
	if f := l.list.UpdateItem; f != nil {
		f(id, li.child)
	}
	li.onTapped = func() {
		if !fyne.CurrentDevice().IsMobile() {
			canvas := fyne.CurrentApp().Driver().CanvasForObject(l.list)
			if canvas != nil {
				canvas.Focus(l.list)
			}

			l.list.currentFocus = id
		}

		l.list.Select(id)
	}
}

func (l *listLayout) updateList(newOnly bool) {
	l.renderLock.Lock()
	separatorThickness := theme.Padding()
	width := l.list.Size().Width
	length := 0
	if f := l.list.Length; f != nil {
		length = f()
	}
	if l.list.UpdateItem == nil {
		fyne.LogError("Missing UpdateCell callback required for List", nil)
	}

	// Keep pointer reference for copying slice header when returning to the pool
	// https://blog.mike.norgate.xyz/unlocking-go-slice-performance-navigating-sync-pool-for-enhanced-efficiency-7cb63b0b453e
	wasVisiblePtr := l.slicePool.Get().(*[]listItemAndID)
	wasVisible := (*wasVisiblePtr)[:0]
	wasVisible = append(wasVisible, l.visible...)

	l.list.propertyLock.Lock()
	offY, minRow := l.calculateVisibleRowHeights(l.list.itemMin.Height, length)
	l.list.propertyLock.Unlock()
	if len(l.visibleRowHeights) == 0 && length > 0 { // we can't show anything until we have some dimensions
		l.renderLock.Unlock() // user code should not be locked
		return
	}

	oldVisibleLen := len(l.visible)
	l.visible = l.visible[:0]
	oldChildrenLen := len(l.children)
	l.children = l.children[:0]

	y := offY
	for index, itemHeight := range l.visibleRowHeights {
		row := index + minRow
		size := fyne.NewSize(width, itemHeight)

		c, ok := l.searchVisible(wasVisible, row)
		if !ok {
			c = l.getItem()
			if c == nil {
				continue
			}
			c.Resize(size)
		}

		c.Move(fyne.NewPos(0, y))
		c.Resize(size)

		y += itemHeight + separatorThickness
		l.visible = append(l.visible, listItemAndID{id: row, item: c})
		l.children = append(l.children, c)
	}
	l.nilOldSliceData(l.children, len(l.children), oldChildrenLen)
	l.nilOldVisibleSliceData(l.visible, len(l.visible), oldVisibleLen)

	for _, wasVis := range wasVisible {
		if _, ok := l.searchVisible(l.visible, wasVis.id); !ok {
			l.itemPool.Release(wasVis.item)
		}
	}

	l.updateSeparators()

	c := l.list.scroller.Content.(*fyne.Container)
	oldObjLen := len(c.Objects)
	c.Objects = c.Objects[:0]
	c.Objects = append(c.Objects, l.children...)
	c.Objects = append(c.Objects, l.separators...)
	l.nilOldSliceData(c.Objects, len(c.Objects), oldObjLen)

	// make a local deep copy of l.visible since rest of this function is unlocked
	// and cannot safely access l.visible
	visiblePtr := l.slicePool.Get().(*[]listItemAndID)
	visible := (*visiblePtr)[:0]
	visible = append(visible, l.visible...)
	l.renderLock.Unlock() // user code should not be locked

	if newOnly {
		for _, vis := range visible {
			if _, ok := l.searchVisible(wasVisible, vis.id); !ok {
				l.setupListItem(vis.item, vis.id, l.list.focused && l.list.currentFocus == vis.id)
			}
		}
	} else {
		for _, vis := range visible {
			l.setupListItem(vis.item, vis.id, l.list.focused && l.list.currentFocus == vis.id)
		}
	}

	// nil out all references before returning slices to pool
	for i := 0; i < len(wasVisible); i++ {
		wasVisible[i].item = nil
	}
	for i := 0; i < len(visible); i++ {
		visible[i].item = nil
	}
	*wasVisiblePtr = wasVisible // Copy the stack header over to the heap
	*visiblePtr = visible
	l.slicePool.Put(wasVisiblePtr)
	l.slicePool.Put(visiblePtr)
}

func (l *listLayout) updateSeparators() {
	if lenChildren := len(l.children); lenChildren > 1 {
		if lenSep := len(l.separators); lenSep > lenChildren {
			l.separators = l.separators[:lenChildren]
		} else {
			for i := lenSep; i < lenChildren; i++ {
				l.separators = append(l.separators, NewSeparator())
			}
		}
	} else {
		l.separators = nil
	}

	separatorThickness := theme.SeparatorThicknessSize()
	dividerOff := (theme.Padding() + separatorThickness) / 2
	for i, child := range l.children {
		if i == 0 {
			continue
		}
		l.separators[i].Move(fyne.NewPos(0, child.Position().Y-dividerOff))
		l.separators[i].Resize(fyne.NewSize(l.list.Size().Width, separatorThickness))
		l.separators[i].Show()
	}
}

// invariant: visible is in ascending order of IDs
func (l *listLayout) searchVisible(visible []listItemAndID, id ListItemID) (*listItem, bool) {
	ln := len(visible)
	idx := sort.Search(ln, func(i int) bool { return visible[i].id >= id })
	if idx < ln && visible[idx].id == id {
		return visible[idx].item, true
	}
	return nil, false
}

func (l *listLayout) nilOldSliceData(objs []fyne.CanvasObject, len, oldLen int) {
	if oldLen > len {
		objs = objs[:oldLen] // gain view into old data
		for i := len; i < oldLen; i++ {
			objs[i] = nil
		}
	}
}

func (l *listLayout) nilOldVisibleSliceData(objs []listItemAndID, len, oldLen int) {
	if oldLen > len {
		objs = objs[:oldLen] // gain view into old data
		for i := len; i < oldLen; i++ {
			objs[i].item = nil
		}
	}
}
