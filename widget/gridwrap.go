package widget

import (
	"fmt"
	"math"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/async"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
)

// Declare conformity with interfaces.
var (
	_ fyne.Widget    = (*GridWrap)(nil)
	_ fyne.Focusable = (*GridWrap)(nil)
)

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

	// Length is a callback for returning the number of items in the GridWrap.
	Length func() int `json:"-"`

	// CreateItem is a callback invoked to create a new widget to render
	// an item in the GridWrap.
	CreateItem func() fyne.CanvasObject `json:"-"`

	// UpdateItem is a callback invoked to update a GridWrap item widget
	// to display a new item in the list. The UpdateItem callback should
	// only update the given item, it should not invoke APIs that would
	// change other properties of the GridWrap itself.
	UpdateItem func(id GridWrapItemID, item fyne.CanvasObject) `json:"-"`

	// OnSelected is a callback to be notified when a given item
	// in the GridWrap has been selected.
	OnSelected func(id GridWrapItemID) `json:"-"`

	// OnSelected is a callback to be notified when a given item
	// in the GridWrap has been unselected.
	OnUnselected func(id GridWrapItemID) `json:"-"`

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
			item, err := data.GetItem(i)
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
		item := createItemAndApplyThemeScope(f, l)

		l.itemMin = item.MinSize()
	}

	layout := &fyne.Container{Layout: newGridWrapLayout(l)}
	l.scroller = widget.NewVScroll(layout)
	layout.Resize(layout.MinSize())

	return newGridWrapRenderer([]fyne.CanvasObject{l.scroller}, l, l.scroller, layout)
}

// FocusGained is called after this GridWrap has gained focus.
func (l *GridWrap) FocusGained() {
	l.focused = true
	l.RefreshItem(l.currentFocus)
}

// FocusLost is called after this GridWrap has lost focus.
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

	pad := l.Theme().Size(theme.SizeNamePadding)
	row := math.Floor(float64(id) / float64(l.ColumnCount()))
	y := float32(row)*l.itemMin.Height + float32(row)*pad
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
	item, ok := lo.searchVisible(lo.visible, id)
	if ok {
		lo.setupGridItem(item, id, l.focused && l.currentFocus == id)
	}
}

// Resize is called when this GridWrap should change size. We refresh to ensure invisible items are drawn.
func (l *GridWrap) Resize(s fyne.Size) {
	oldColCount := l.ColumnCount()
	oldHeight := l.size.Height
	l.colCountCache = 0
	l.BaseWidget.Resize(s)
	newColCount := l.ColumnCount()

	if oldColCount == newColCount && oldHeight == s.Height {
		// no content update needed if resizing only horizontally and col count is unchanged
		return
	}

	if l.scroller != nil {
		l.offsetUpdated(l.scroller.Offset)
		l.scroller.Content.(*fyne.Container).Layout.(*gridWrapLayout).updateGrid(true)
	}
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
	if id < 0 || id >= length {
		return
	}
	l.scrollTo(id)
	l.Refresh()
}

// ScrollToBottom scrolls to the end of the list
func (l *GridWrap) ScrollToBottom() {
	l.scroller.ScrollToBottom()
	l.offsetUpdated(l.scroller.Offset)
}

// ScrollToTop scrolls to the start of the list
func (l *GridWrap) ScrollToTop() {
	l.scroller.ScrollToTop()
	l.offsetUpdated(l.scroller.Offset)
}

// ScrollToOffset scrolls the list to the given offset position
func (l *GridWrap) ScrollToOffset(offset float32) {
	if l.scroller == nil {
		return
	}
	if offset < 0 {
		offset = 0
	}
	contentHeight := l.contentMinSize().Height
	if l.Size().Height >= contentHeight {
		return // content fully visible - no need to scroll
	}
	if offset > contentHeight {
		offset = contentHeight
	}
	l.scroller.ScrollToOffset(fyne.NewPos(0, offset))
	l.offsetUpdated(l.scroller.Offset)
}

// TypedKey is called if a key event happens while this GridWrap is focused.
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
		l.currentFocus += l.ColumnCount()
		if l.currentFocus >= count-1 {
			l.currentFocus = count - 1
		}
		l.scrollTo(l.currentFocus)
		l.RefreshItem(l.currentFocus)
	case fyne.KeyLeft:
		if l.currentFocus <= 0 {
			return
		}
		if l.currentFocus%l.ColumnCount() == 0 {
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
		if (l.currentFocus+1)%l.ColumnCount() == 0 {
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
		l.currentFocus -= l.ColumnCount()
		if l.currentFocus < 0 {
			l.currentFocus = 0
		}
		l.scrollTo(l.currentFocus)
		l.RefreshItem(l.currentFocus)
	}
}

// TypedRune is called if a text event happens while this GridWrap is focused.
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

func (l *GridWrap) contentMinSize() fyne.Size {
	padding := l.Theme().Size(theme.SizeNamePadding)
	if l.Length == nil {
		return fyne.NewSize(0, 0)
	}

	cols := l.ColumnCount()
	rows := float32(math.Ceil(float64(l.Length()) / float64(cols)))
	return fyne.NewSize(l.itemMin.Width, (l.itemMin.Height+padding)*rows-padding)
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
		item := createItemAndApplyThemeScope(f, l.list)

		l.list.itemMin = item.MinSize()
	}
	l.Layout(l.list.Size())
	l.scroller.Refresh()
	l.layout.Layout.(*gridWrapLayout).updateGrid(false)
	canvas.Refresh(l.list)
}

func (l *gridWrapRenderer) Destroy() {
}

func (l *gridWrapRenderer) Objects() []fyne.CanvasObject {
	return l.objects
}

// Declare conformity with interfaces.
var (
	_ fyne.Widget       = (*gridWrapItem)(nil)
	_ fyne.Tappable     = (*gridWrapItem)(nil)
	_ desktop.Hoverable = (*gridWrapItem)(nil)
)

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
	th := gw.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	gw.background = canvas.NewRectangle(th.Color(theme.ColorNameHover, v))
	gw.background.CornerRadius = th.Size(theme.SizeNameSelectionRadius)
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
	th := gw.item.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	gw.item.background.CornerRadius = th.Size(theme.SizeNameSelectionRadius)
	if gw.item.selected {
		gw.item.background.FillColor = th.Color(theme.ColorNameSelection, v)
		gw.item.background.Show()
	} else if gw.item.hovered {
		gw.item.background.FillColor = th.Color(theme.ColorNameHover, v)
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
	gw *GridWrap

	itemPool   async.Pool[fyne.CanvasObject]
	visible    []gridItemAndID
	wasVisible []gridItemAndID
}

func newGridWrapLayout(gw *GridWrap) fyne.Layout {
	gwl := &gridWrapLayout{gw: gw}
	gw.offsetUpdated = gwl.offsetUpdated
	return gwl
}

func (l *gridWrapLayout) Layout(_ []fyne.CanvasObject, _ fyne.Size) {
	l.updateGrid(true)
}

func (l *gridWrapLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return l.gw.contentMinSize()
}

func (l *gridWrapLayout) getItem() *gridWrapItem {
	item := l.itemPool.Get()
	if item == nil {
		if f := l.gw.CreateItem; f != nil {
			child := createItemAndApplyThemeScope(f, l.gw)

			item = newGridWrapItem(child, nil)
		}
	}
	return item.(*gridWrapItem)
}

func (l *gridWrapLayout) offsetUpdated(pos fyne.Position) {
	if l.gw.offsetY == pos.Y {
		return
	}
	l.gw.offsetY = pos.Y
	l.updateGrid(true)
}

func (l *gridWrapLayout) setupGridItem(li *gridWrapItem, id GridWrapItemID, focus bool) {
	previousIndicator := li.selected
	li.selected = false
	for _, s := range l.gw.selected {
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
	if f := l.gw.UpdateItem; f != nil {
		f(id, li.child)
	}
	li.onTapped = func() {
		if !fyne.CurrentDevice().IsMobile() {
			l.gw.RefreshItem(l.gw.currentFocus)
			canvas := fyne.CurrentApp().Driver().CanvasForObject(l.gw)
			if canvas != nil {
				canvas.Focus(l.gw.impl.(fyne.Focusable))
			}

			l.gw.currentFocus = id
		}

		l.gw.Select(id)
	}
}

// ColumnCount returns the number of columns that are/will be shown
// in this GridWrap, based on the widget's current width.
//
// Since: 2.5
func (l *GridWrap) ColumnCount() int {
	if l.colCountCache < 1 {
		padding := l.Theme().Size(theme.SizeNamePadding)
		l.colCountCache = 1
		width := l.Size().Width
		if width > l.itemMin.Width {
			l.colCountCache = int(math.Floor(float64(width+padding) / float64(l.itemMin.Width+padding)))
		}
	}
	return l.colCountCache
}

func (l *gridWrapLayout) updateGrid(newOnly bool) {
	// code here is a mashup of listLayout.updateList and gridWrapLayout.Layout
	padding := l.gw.Theme().Size(theme.SizeNamePadding)

	length := 0
	if f := l.gw.Length; f != nil {
		length = f()
	}

	colCount := l.gw.ColumnCount()
	visibleRowsCount := int(math.Ceil(float64(l.gw.scroller.Size().Height)/float64(l.gw.itemMin.Height+padding))) + 1

	offY := l.gw.offsetY - float32(math.Mod(float64(l.gw.offsetY), float64(l.gw.itemMin.Height+padding)))
	minRow := int(offY / (l.gw.itemMin.Height + padding))
	minItem := minRow * colCount
	maxRow := int(math.Min(float64(minRow+visibleRowsCount), math.Ceil(float64(length)/float64(colCount))))
	maxItem := GridWrapItemID(math.Min(float64(maxRow*colCount), float64(length-1)))

	if l.gw.UpdateItem == nil {
		fyne.LogError("Missing UpdateCell callback required for GridWrap", nil)
	}

	// l.wasVisible now represents the currently visible items, while
	// l.visible will be updated to represent what is visible *after* the update
	l.wasVisible = append(l.wasVisible, l.visible...)
	l.visible = l.visible[:0]

	c := l.gw.scroller.Content.(*fyne.Container)
	oldObjLen := len(c.Objects)
	c.Objects = c.Objects[:0]
	y := offY
	curItemID := minItem
	for row := minRow; row <= maxRow && curItemID <= maxItem; row++ {
		x := float32(0)
		for col := 0; col < colCount && curItemID <= maxItem; col++ {
			item, ok := l.searchVisible(l.wasVisible, curItemID)
			if !ok {
				item = l.getItem()
				if item == nil {
					continue
				}
				item.Resize(l.gw.itemMin)
			}

			item.Move(fyne.NewPos(x, y))
			item.Resize(l.gw.itemMin)

			x += l.gw.itemMin.Width + padding
			l.visible = append(l.visible, gridItemAndID{item: item, id: curItemID})
			c.Objects = append(c.Objects, item)
			curItemID++
		}
		y += l.gw.itemMin.Height + padding
	}
	l.nilOldSliceData(c.Objects, len(c.Objects), oldObjLen)

	for _, old := range l.wasVisible {
		if _, ok := l.searchVisible(l.visible, old.id); !ok {
			l.itemPool.Put(old.item)
		}
	}

	if newOnly {
		for _, obj := range l.visible {
			if _, ok := l.searchVisible(l.wasVisible, obj.id); !ok {
				l.setupGridItem(obj.item, obj.id, l.gw.focused && l.gw.currentFocus == obj.id)
			}
		}
	} else {
		for _, obj := range l.visible {
			l.setupGridItem(obj.item, obj.id, l.gw.focused && l.gw.currentFocus == obj.id)
		}
	}

	// we don't need wasVisible now until next call to update
	// nil out all references before truncating slice
	for i := 0; i < len(l.wasVisible); i++ {
		l.wasVisible[i].item = nil
	}
	l.wasVisible = l.wasVisible[:0]
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
