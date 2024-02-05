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

// Declare conformity with interfaces.
var _ fyne.Widget = (*GridWrap)(nil)
var _ fyne.Focusable = (*GridWrap)(nil)

// GridWrapItemID is the ID of an individual item in the GridWrap widget.
//
// Since: 2.4
type GridWrapItemID = int

// GridWrap is a widget with an API very similar to widget.List,
// that lays out items in a scrollable wrapping grid similar to container.NewGridWrap.
// It caches and reuses widgets for performance.
//
// Since: 2.4
type GridWrap struct {
	BaseWidget

	Length       func() int                                      `json:"-"`
	CreateItem   func() fyne.CanvasObject                        `json:"-"`
	UpdateItem   func(id GridWrapItemID, item fyne.CanvasObject) `json:"-"`
	OnSelected   func(id GridWrapItemID)                         `json:"-"`
	OnUnselected func(id GridWrapItemID)                         `json:"-"`

	currentFocus  ListItemID
	focused       bool
	scroller      *widget.Scroll
	selected      []GridWrapItemID
	itemMin       fyne.Size
	offsetY       float32
	offsetUpdated func(fyne.Position)
	colCountCache int
}

// NewGridWrap creates and returns a GridWrap widget for displaying items in
// a wrapping grid layout with scrolling and caching for performance.
//
// Since: 2.4
func NewGridWrap(length func() int, createItem func() fyne.CanvasObject, updateItem func(GridWrapItemID, fyne.CanvasObject)) *GridWrap {
	gwList := &GridWrap{Length: length, CreateItem: createItem, UpdateItem: updateItem}
	gwList.ExtendBaseWidget(gwList)
	return gwList
}

// NewGridWrapWithData creates a new GridWrap widget that will display the contents of the provided data.
//
// Since: 2.4
func NewGridWrapWithData(data binding.DataList, createItem func() fyne.CanvasObject, updateItem func(binding.DataItem, fyne.CanvasObject)) *GridWrap {
	gwList := NewGridWrap(
		data.Length,
		createItem,
		func(i GridWrapItemID, o fyne.CanvasObject) {
			item, err := data.GetItem(int(i))
			if err != nil {
				fyne.LogError(fmt.Sprintf("Error getting data item %d", i), err)
				return
			}
			updateItem(item, o)
		})

	data.AddListener(binding.NewDataListener(gwList.Refresh))
	return gwList
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer.
func (l *GridWrap) CreateRenderer() fyne.WidgetRenderer {
	l.ExtendBaseWidget(l)

	if f := l.CreateItem; f != nil && l.itemMin.IsZero() {
		l.itemMin = f().MinSize()
	}

	layout := &fyne.Container{Layout: newGridWrapLayout(l)}
	l.scroller = widget.NewVScroll(layout)
	layout.Resize(layout.MinSize())

	return newGridWrapRenderer([]fyne.CanvasObject{l.scroller}, l, l.scroller, layout)
}

// FocusGained is called after this GridWrap has gained focus.
//
// Implements: fyne.Focusable
func (l *GridWrap) FocusGained() {
	l.focused = true
	l.scrollTo(l.currentFocus)
	l.RefreshItem(l.currentFocus)
}

// FocusLost is called after this GridWrap has lost focus.
//
// Implements: fyne.Focusable
func (l *GridWrap) FocusLost() {
	l.focused = false
	l.RefreshItem(l.currentFocus)
}

// MinSize returns the size that this widget should not shrink below.
func (l *GridWrap) MinSize() fyne.Size {
	l.ExtendBaseWidget(l)

	return l.BaseWidget.MinSize()
}

func (l *GridWrap) scrollTo(id GridWrapItemID) {
	if l.scroller == nil {
		return
	}
	row := math.Floor(float64(id) / float64(l.getColCount()))
	y := float32(row)*l.itemMin.Height + float32(row)*theme.Padding()
	if y < l.scroller.Offset.Y {
		l.scroller.Offset.Y = y
	} else if size := l.scroller.Size(); y+l.itemMin.Height > l.scroller.Offset.Y+size.Height {
		l.scroller.Offset.Y = y + l.itemMin.Height - size.Height
	}
	l.offsetUpdated(l.scroller.Offset)
}

// RefreshItem refreshes a single item, specified by the item ID passed in.
//
// Since: 2.4
func (l *GridWrap) RefreshItem(id GridWrapItemID) {
	if l.scroller == nil {
		return
	}
	l.BaseWidget.Refresh()
	lo := l.scroller.Content.(*fyne.Container).Layout.(*gridWrapLayout)
	lo.renderLock.Lock() // ensures we are not changing visible info in render code during the search
	item, ok := lo.searchVisible(lo.visible, id)
	lo.renderLock.Unlock()
	if ok {
		lo.setupGridItem(item, id, l.focused && l.currentFocus == id)
	}
}

// Resize is called when this GridWrap should change size. We refresh to ensure invisible items are drawn.
func (l *GridWrap) Resize(s fyne.Size) {
	l.colCountCache = 0
	l.BaseWidget.Resize(s)
	l.offsetUpdated(l.scroller.Offset)
	l.scroller.Content.(*fyne.Container).Layout.(*gridWrapLayout).updateGrid(true)
}

// Select adds the item identified by the given ID to the selection.
func (l *GridWrap) Select(id GridWrapItemID) {
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
	l.selected = []GridWrapItemID{id}
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
func (l *GridWrap) ScrollTo(id GridWrapItemID) {
	length := 0
	if f := l.Length; f != nil {
		length = f()
	}
	if id < 0 || int(id) >= length {
		return
	}
	l.scrollTo(id)
	l.Refresh()
}

// ScrollToBottom scrolls to the end of the list
func (l *GridWrap) ScrollToBottom() {
	length := 0
	if f := l.Length; f != nil {
		length = f()
	}
	if length > 0 {
		length--
	}
	l.scrollTo(GridWrapItemID(length))
	l.Refresh()
}

// ScrollToTop scrolls to the start of the list
func (l *GridWrap) ScrollToTop() {
	l.scrollTo(0)
	l.Refresh()
}

// ScrollToOffset scrolls the list to the given offset position
func (l *GridWrap) ScrollToOffset(offset float32) {
	l.scroller.Offset.Y = offset
	l.offsetUpdated(l.scroller.Offset)
}

// TypedKey is called if a key event happens while this GridWrap is focused.
//
// Implements: fyne.Focusable
func (l *GridWrap) TypedKey(event *fyne.KeyEvent) {
	switch event.Name {
	case fyne.KeySpace:
		l.Select(l.currentFocus)
	case fyne.KeyDown:
		count := 0
		if f := l.Length; f != nil {
			count = f()
		}
		l.RefreshItem(l.currentFocus)
		l.currentFocus += l.getColCount()
		if l.currentFocus >= count-1 {
			l.currentFocus = count - 1
		}
		l.scrollTo(l.currentFocus)
		l.RefreshItem(l.currentFocus)
	case fyne.KeyLeft:
		if l.currentFocus <= 0 {
			return
		}
		if l.currentFocus%l.getColCount() == 0 {
			return
		}

		l.RefreshItem(l.currentFocus)
		l.currentFocus--
		l.scrollTo(l.currentFocus)
		l.RefreshItem(l.currentFocus)
	case fyne.KeyRight:
		if f := l.Length; f != nil && l.currentFocus >= f()-1 {
			return
		}
		if (l.currentFocus+1)%l.getColCount() == 0 {
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
		l.currentFocus -= l.getColCount()
		if l.currentFocus < 0 {
			l.currentFocus = 0
		}
		l.scrollTo(l.currentFocus)
		l.RefreshItem(l.currentFocus)
	}
}

// TypedRune is called if a text event happens while this GridWrap is focused.
//
// Implements: fyne.Focusable
func (l *GridWrap) TypedRune(_ rune) {
	// intentionally left blank
}

// GetScrollOffset returns the current scroll offset position
func (l *GridWrap) GetScrollOffset() float32 {
	return l.offsetY
}

// Unselect removes the item identified by the given ID from the selection.
func (l *GridWrap) Unselect(id GridWrapItemID) {
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
func (l *GridWrap) UnselectAll() {
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

// Declare conformity with WidgetRenderer interface.
var _ fyne.WidgetRenderer = (*gridWrapRenderer)(nil)

type gridWrapRenderer struct {
	objects []fyne.CanvasObject

	list     *GridWrap
	scroller *widget.Scroll
	layout   *fyne.Container
}

func newGridWrapRenderer(objects []fyne.CanvasObject, l *GridWrap, scroller *widget.Scroll, layout *fyne.Container) *gridWrapRenderer {
	lr := &gridWrapRenderer{objects: objects, list: l, scroller: scroller, layout: layout}
	lr.scroller.OnScrolled = l.offsetUpdated
	return lr
}

func (l *gridWrapRenderer) Layout(size fyne.Size) {
	l.scroller.Resize(size)
}

func (l *gridWrapRenderer) MinSize() fyne.Size {
	return l.scroller.MinSize().Max(l.list.itemMin)
}

func (l *gridWrapRenderer) Refresh() {
	if f := l.list.CreateItem; f != nil {
		l.list.itemMin = f().MinSize()
	}
	l.Layout(l.list.Size())
	l.scroller.Refresh()
	l.layout.Layout.(*gridWrapLayout).updateGrid(true)
	canvas.Refresh(l.list)
}

func (l *gridWrapRenderer) Destroy() {
}

func (l *gridWrapRenderer) Objects() []fyne.CanvasObject {
	return l.objects
}

// Declare conformity with interfaces.
var _ fyne.Widget = (*gridWrapItem)(nil)
var _ fyne.Tappable = (*gridWrapItem)(nil)
var _ desktop.Hoverable = (*gridWrapItem)(nil)

type gridWrapItem struct {
	BaseWidget

	onTapped          func()
	background        *canvas.Rectangle
	child             fyne.CanvasObject
	hovered, selected bool
}

func newGridWrapItem(child fyne.CanvasObject, tapped func()) *gridWrapItem {
	gw := &gridWrapItem{
		child:    child,
		onTapped: tapped,
	}

	gw.ExtendBaseWidget(gw)
	return gw
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer.
func (gw *gridWrapItem) CreateRenderer() fyne.WidgetRenderer {
	gw.ExtendBaseWidget(gw)

	gw.background = canvas.NewRectangle(theme.HoverColor())
	gw.background.CornerRadius = theme.SelectionRadiusSize()
	gw.background.Hide()

	objects := []fyne.CanvasObject{gw.background, gw.child}

	return &gridWrapItemRenderer{widget.NewBaseRenderer(objects), gw}
}

// MinSize returns the size that this widget should not shrink below.
func (gw *gridWrapItem) MinSize() fyne.Size {
	gw.ExtendBaseWidget(gw)
	return gw.BaseWidget.MinSize()
}

// MouseIn is called when a desktop pointer enters the widget.
func (gw *gridWrapItem) MouseIn(*desktop.MouseEvent) {
	gw.hovered = true
	gw.Refresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget.
func (gw *gridWrapItem) MouseMoved(*desktop.MouseEvent) {
}

// MouseOut is called when a desktop pointer exits the widget.
func (gw *gridWrapItem) MouseOut() {
	gw.hovered = false
	gw.Refresh()
}

// Tapped is called when a pointer tapped event is captured and triggers any tap handler.
func (gw *gridWrapItem) Tapped(*fyne.PointEvent) {
	if gw.onTapped != nil {
		gw.selected = true
		gw.Refresh()
		gw.onTapped()
	}
}

// Declare conformity with the WidgetRenderer interface.
var _ fyne.WidgetRenderer = (*gridWrapItemRenderer)(nil)

type gridWrapItemRenderer struct {
	widget.BaseRenderer

	item *gridWrapItem
}

// MinSize calculates the minimum size of a listItem.
// This is based on the size of the status indicator and the size of the child object.
func (gw *gridWrapItemRenderer) MinSize() fyne.Size {
	return gw.item.child.MinSize()
}

// Layout the components of the listItem widget.
func (gw *gridWrapItemRenderer) Layout(size fyne.Size) {
	gw.item.background.Resize(size)
	gw.item.child.Resize(size)
}

func (gw *gridWrapItemRenderer) Refresh() {
	gw.item.background.CornerRadius = theme.SelectionRadiusSize()
	if gw.item.selected {
		gw.item.background.FillColor = theme.SelectionColor()
		gw.item.background.Show()
	} else if gw.item.hovered {
		gw.item.background.FillColor = theme.HoverColor()
		gw.item.background.Show()
	} else {
		gw.item.background.Hide()
	}
	gw.item.background.Refresh()
	canvas.Refresh(gw.item.super())
}

// Declare conformity with Layout interface.
var _ fyne.Layout = (*gridWrapLayout)(nil)

type gridItemAndID struct {
	item *gridWrapItem
	id   GridWrapItemID
}

type gridWrapLayout struct {
	list *GridWrap

	itemPool   syncPool
	slicePool  sync.Pool // *[]itemAndID
	visible    []gridItemAndID
	renderLock sync.Mutex
}

func newGridWrapLayout(list *GridWrap) fyne.Layout {
	l := &gridWrapLayout{list: list}
	l.slicePool.New = func() interface{} {
		s := make([]gridItemAndID, 0)
		return &s
	}
	list.offsetUpdated = l.offsetUpdated
	return l
}

func (l *gridWrapLayout) Layout(_ []fyne.CanvasObject, _ fyne.Size) {
	l.updateGrid(true)
}

func (l *gridWrapLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	padding := theme.Padding()
	if lenF := l.list.Length; lenF != nil {
		cols := l.list.getColCount()
		rows := float32(math.Ceil(float64(lenF()) / float64(cols)))
		return fyne.NewSize(l.list.itemMin.Width,
			(l.list.itemMin.Height+padding)*rows-padding)
	}
	return fyne.NewSize(0, 0)
}

func (l *gridWrapLayout) getItem() *gridWrapItem {
	item := l.itemPool.Obtain()
	if item == nil {
		if f := l.list.CreateItem; f != nil {
			item = newGridWrapItem(f(), nil)
		}
	}
	return item.(*gridWrapItem)
}

func (l *gridWrapLayout) offsetUpdated(pos fyne.Position) {
	if l.list.offsetY == pos.Y {
		return
	}
	l.list.offsetY = pos.Y
	l.updateGrid(false)
}

func (l *gridWrapLayout) setupGridItem(li *gridWrapItem, id GridWrapItemID, focus bool) {
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
			l.list.RefreshItem(l.list.currentFocus)
			canvas := fyne.CurrentApp().Driver().CanvasForObject(l.list)
			if canvas != nil {
				canvas.Focus(l.list)
			}

			l.list.currentFocus = id
		}

		l.list.Select(id)
	}
}

func (l *GridWrap) getColCount() int {
	if l.colCountCache < 1 {
		padding := theme.Padding()
		l.colCountCache = 1
		width := l.Size().Width
		if width > l.itemMin.Width {
			l.colCountCache = int(math.Floor(float64(width+padding) / float64(l.itemMin.Width+padding)))
		}
	}
	return l.colCountCache
}

func (l *gridWrapLayout) updateGrid(refresh bool) {
	// code here is a mashup of listLayout.updateList and gridWrapLayout.Layout
	padding := theme.Padding()

	l.renderLock.Lock()
	length := 0
	if f := l.list.Length; f != nil {
		length = f()
	}

	colCount := l.list.getColCount()
	visibleRowsCount := int(math.Ceil(float64(l.list.scroller.Size().Height)/float64(l.list.itemMin.Height+padding))) + 1

	offY := l.list.offsetY - float32(math.Mod(float64(l.list.offsetY), float64(l.list.itemMin.Height+padding)))
	minRow := int(offY / (l.list.itemMin.Height + padding))
	minItem := GridWrapItemID(minRow * colCount)
	maxRow := int(math.Min(float64(minRow+visibleRowsCount), math.Ceil(float64(length)/float64(colCount))))
	maxItem := GridWrapItemID(math.Min(float64(maxRow*colCount), float64(length-1)))

	if l.list.UpdateItem == nil {
		fyne.LogError("Missing UpdateCell callback required for GridWrap", nil)
	}

	// Keep pointer reference for copying slice header when returning to the pool
	// https://blog.mike.norgate.xyz/unlocking-go-slice-performance-navigating-sync-pool-for-enhanced-efficiency-7cb63b0b453e
	wasVisiblePtr := l.slicePool.Get().(*[]gridItemAndID)
	wasVisible := (*wasVisiblePtr)[:0]
	wasVisible = append(wasVisible, l.visible...)

	oldVisibleLen := len(l.visible)
	l.visible = l.visible[:0]

	c := l.list.scroller.Content.(*fyne.Container)
	oldObjLen := len(c.Objects)
	c.Objects = c.Objects[:0]
	y := offY
	curItemID := minItem
	for row := minRow; row <= maxRow && curItemID <= maxItem; row++ {
		x := float32(0)
		for col := 0; col < colCount && curItemID <= maxItem; col++ {
			item, ok := l.searchVisible(wasVisible, curItemID)
			if !ok {
				item = l.getItem()
				if item == nil {
					continue
				}
				item.Resize(l.list.itemMin)
			}

			item.Move(fyne.NewPos(x, y))
			if refresh {
				item.Resize(l.list.itemMin)
			}

			x += l.list.itemMin.Width + padding
			l.visible = append(l.visible, gridItemAndID{item: item, id: curItemID})
			c.Objects = append(c.Objects, item)
			curItemID++
		}
		y += l.list.itemMin.Height + padding
	}
	l.nilOldSliceData(c.Objects, len(c.Objects), oldObjLen)
	l.nilOldVisibleSliceData(l.visible, len(l.visible), oldVisibleLen)

	for _, old := range wasVisible {
		if _, ok := l.searchVisible(l.visible, old.id); !ok {
			l.itemPool.Release(old.item)
		}
	}

	// make a local deep copy of l.visible since rest of this function is unlocked
	// and cannot safely access l.visible
	visiblePtr := l.slicePool.Get().(*[]gridItemAndID)
	visible := (*visiblePtr)[:0]
	visible = append(visible, l.visible...)
	l.renderLock.Unlock() // user code should not be locked

	for _, obj := range visible {
		l.setupGridItem(obj.item, obj.id, l.list.focused && l.list.currentFocus == obj.id)
	}

	// nil out all references before returning slices to pool
	for i := 0; i < len(wasVisible); i++ {
		wasVisible[i].item = nil
	}
	for i := 0; i < len(visible); i++ {
		visible[i].item = nil
	}
	*wasVisiblePtr = wasVisible // Copy the slice header over to the heap
	*visiblePtr = visible
	l.slicePool.Put(wasVisiblePtr)
	l.slicePool.Put(visiblePtr)
}

// invariant: visible is in ascending order of IDs
func (l *gridWrapLayout) searchVisible(visible []gridItemAndID, id GridWrapItemID) (*gridWrapItem, bool) {
	ln := len(visible)
	idx := sort.Search(ln, func(i int) bool { return visible[i].id >= id })
	if idx < ln && visible[idx].id == id {
		return visible[idx].item, true
	}
	return nil, false
}

func (l *gridWrapLayout) nilOldSliceData(objs []fyne.CanvasObject, len, oldLen int) {
	if oldLen > len {
		objs = objs[:oldLen] // gain view into old data
		for i := len; i < oldLen; i++ {
			objs[i] = nil
		}
	}
}

func (l *gridWrapLayout) nilOldVisibleSliceData(objs []gridItemAndID, len, oldLen int) {
	if oldLen > len {
		objs = objs[:oldLen] // gain view into old data
		for i := len; i < oldLen; i++ {
			objs[i].item = nil
		}
	}
}
