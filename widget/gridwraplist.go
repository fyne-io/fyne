package widget

import (
	"fmt"
	"math"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
)

// Declare conformity with Widget interface.
var _ fyne.Widget = (*GridWrapList)(nil)

// List is a widget that pools list items for performance and
// lays the items out in a vertical direction inside of a scroller.
// List requires that all items are the same size.
//
// Since: TODO version
type GridWrapList struct {
	BaseWidget

	Length       func() int                                  `json:"-"`
	CreateItem   func() fyne.CanvasObject                    `json:"-"`
	UpdateItem   func(id ListItemID, item fyne.CanvasObject) `json:"-"`
	OnSelected   func(id ListItemID)                         `json:"-"`
	OnUnselected func(id ListItemID)                         `json:"-"`

	scroller      *widget.Scroll
	selected      []ListItemID
	itemMin       fyne.Size
	offsetY       float32
	offsetUpdated func(fyne.Position)
}

// NewList creates and returns a list widget for displaying items in
// a vertical layout with scrolling and caching for performance.
//
// Since: TODO version
func NewGridWrapList(length func() int, createItem func() fyne.CanvasObject, updateItem func(ListItemID, fyne.CanvasObject)) *GridWrapList {
	gwList := &GridWrapList{BaseWidget: BaseWidget{}, Length: length, CreateItem: createItem, UpdateItem: updateItem}
	gwList.ExtendBaseWidget(gwList)
	return gwList
}

// NewListWithData creates a new list widget that will display the contents of the provided data.
//
// Since: TODO version
func NewGridWrapListWithData(data binding.DataList, createItem func() fyne.CanvasObject, updateItem func(binding.DataItem, fyne.CanvasObject)) *GridWrapList {
	gwList := NewGridWrapList(
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

	data.AddListener(binding.NewDataListener(gwList.Refresh))
	return gwList
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer.
func (l *GridWrapList) CreateRenderer() fyne.WidgetRenderer {
	l.ExtendBaseWidget(l)

	if f := l.CreateItem; f != nil {
		if l.itemMin.IsZero() {
			l.itemMin = newListItem(f(), nil).MinSize()
		}
	}
	layout := &fyne.Container{}
	l.scroller = widget.NewVScroll(layout)
	layout.Layout = newGridWrapListLayout(l)
	layout.Resize(layout.MinSize())
	objects := []fyne.CanvasObject{l.scroller}
	lr := newGridWrapListRenderer(objects, l, l.scroller, layout)
	return lr
}

// MinSize returns the size that this widget should not shrink below.
func (l *GridWrapList) MinSize() fyne.Size {
	l.ExtendBaseWidget(l)

	return l.BaseWidget.MinSize()
}

func (l *GridWrapList) scrollTo(id ListItemID) {
	if l.scroller == nil {
		return
	}
	row := math.Floor(float64(id) / float64(l.getColCount()))
	y := (float32(row) * l.itemMin.Height) + (float32(row) * theme.Padding())
	if y < l.scroller.Offset.Y {
		l.scroller.Offset.Y = y
	} else if y+l.itemMin.Height > l.scroller.Offset.Y+l.scroller.Size().Height {
		l.scroller.Offset.Y = y + l.itemMin.Height - l.scroller.Size().Height
	}
	l.offsetUpdated(l.scroller.Offset)
}

// Resize is called when this list should change size. We refresh to ensure invisible items are drawn.
func (l *GridWrapList) Resize(s fyne.Size) {
	l.BaseWidget.Resize(s)
	l.offsetUpdated(l.scroller.Offset)
	l.scroller.Content.(*fyne.Container).Layout.(*gridWrapListLayout).updateList(true)
}

// Select add the item identified by the given ID to the selection.
func (l *GridWrapList) Select(id ListItemID) {
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
func (l *GridWrapList) ScrollTo(id ListItemID) {
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
func (l *GridWrapList) ScrollToBottom() {
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
func (l *GridWrapList) ScrollToTop() {
	l.scrollTo(0)
	l.Refresh()
}

// ScrollToOffset scrolls the list to the given offset position
func (l *GridWrapList) ScrollToOffset(offset float32) {
	// TODO: bounds checking
	l.scroller.Offset.Y = offset
	l.offsetUpdated(l.scroller.Offset)
}

// GetScrollOffset returns the current scroll offset position
func (l *GridWrapList) GetScrollOffset() float32 {
	return l.offsetY
}

// Unselect removes the item identified by the given ID from the selection.
func (l *GridWrapList) Unselect(id ListItemID) {
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
func (l *GridWrapList) UnselectAll() {
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
var _ fyne.WidgetRenderer = (*gridWrapListRenderer)(nil)

type gridWrapListRenderer struct {
	widget.BaseRenderer

	list     *GridWrapList
	scroller *widget.Scroll
	layout   *fyne.Container
}

func newGridWrapListRenderer(objects []fyne.CanvasObject, l *GridWrapList, scroller *widget.Scroll, layout *fyne.Container) *gridWrapListRenderer {
	lr := &gridWrapListRenderer{BaseRenderer: widget.NewBaseRenderer(objects), list: l, scroller: scroller, layout: layout}
	lr.scroller.OnScrolled = l.offsetUpdated
	return lr
}

func (l *gridWrapListRenderer) Layout(size fyne.Size) {
	l.scroller.Resize(size)
}

func (l *gridWrapListRenderer) MinSize() fyne.Size {
	return l.scroller.MinSize().Max(l.list.itemMin)
}

func (l *gridWrapListRenderer) Refresh() {
	if f := l.list.CreateItem; f != nil {
		l.list.itemMin = newListItem(f(), nil).MinSize()
	}
	l.Layout(l.list.Size())
	l.scroller.Refresh()
	l.layout.Layout.(*gridWrapListLayout).updateList(true)
	canvas.Refresh(l.list.super())
}

// Declare conformity with Layout interface.
var _ fyne.Layout = (*gridWrapListLayout)(nil)

type gridWrapListLayout struct {
	list     *GridWrapList
	children []fyne.CanvasObject

	itemPool   *syncPool
	visible    map[ListItemID]*listItem
	renderLock sync.Mutex
}

func newGridWrapListLayout(list *GridWrapList) fyne.Layout {
	l := &gridWrapListLayout{list: list, itemPool: &syncPool{}, visible: make(map[ListItemID]*listItem)}
	list.offsetUpdated = l.offsetUpdated
	return l
}

func (l *gridWrapListLayout) Layout([]fyne.CanvasObject, fyne.Size) {
	l.updateList(true)
}

func (l *gridWrapListLayout) MinSize([]fyne.CanvasObject) fyne.Size {
	if lenF := l.list.Length; lenF != nil {
		cols := l.list.getColCount()
		rows := float32(math.Ceil(float64(lenF()) / float64(cols)))
		return fyne.NewSize(l.list.itemMin.Width,
			(l.list.itemMin.Height+theme.Padding())*rows-theme.Padding())
	}
	return fyne.NewSize(0, 0)
}

func (l *gridWrapListLayout) getItem() *listItem {
	item := l.itemPool.Obtain()
	if item == nil {
		if f := l.list.CreateItem; f != nil {
			item = newListItem(f(), nil)
		}
	}
	return item.(*listItem)
}

func (l *gridWrapListLayout) offsetUpdated(pos fyne.Position) {
	if l.list.offsetY == pos.Y {
		return
	}
	l.list.offsetY = pos.Y
	l.updateList(false)
}

func (l *gridWrapListLayout) setupListItem(li *listItem, id ListItemID) {

	previousIndicator := li.selected
	li.selected = false
	for _, s := range l.list.selected {
		if id == s {
			li.selected = true
			break
		}
	}
	if previousIndicator != li.selected {
		li.Refresh()
	}
	if f := l.list.UpdateItem; f != nil {
		f(id, li.child)
	}
	li.onTapped = func() {
		l.list.Select(id)
	}
}

func (l *GridWrapList) getColCount() int {
	colCount := 1
	width := l.Size().Width
	if width > l.itemMin.Width {
		colCount = int(math.Floor(float64(width+theme.Padding()) / float64(l.itemMin.Width+theme.Padding())))
	}
	return colCount
}

func (l *gridWrapListLayout) updateList(refresh bool) {
	// code here is a mashup of listLayout.updateList and gridWrapLayout.Layout

	l.renderLock.Lock()
	defer l.renderLock.Unlock()
	length := 0
	if f := l.list.Length; f != nil {
		length = f()
	}

	colCount := l.list.getColCount()
	visibleRowsCount := int(math.Ceil(float64(l.list.scroller.Size().Height)/float64(l.list.itemMin.Height+theme.Padding()))) + 1

	offY := l.list.offsetY - float32(math.Mod(float64(l.list.offsetY), float64(l.list.itemMin.Height+theme.Padding())))
	minRow := int(offY / (l.list.itemMin.Height + theme.Padding()))
	minItem := ListItemID(minRow * colCount)
	maxRow := int(math.Min(float64(minRow+visibleRowsCount), math.Ceil(float64(length)/float64(colCount))))
	maxItem := ListItemID(math.Min(float64(maxRow*colCount), float64(length-1)))

	if l.list.UpdateItem == nil {
		fyne.LogError("Missing UpdateCell callback required for GridWrapList", nil)
	}

	wasVisible := l.visible
	l.visible = make(map[ListItemID]*listItem)
	var cells []fyne.CanvasObject
	y := offY
	curItem := minItem
	for row := minRow; row <= maxRow && curItem <= maxItem; row++ {
		x := float32(0)
		for col := 0; col < colCount && curItem <= maxItem; col++ {
			c, ok := wasVisible[curItem]
			if !ok {
				c = l.getItem()
				if c == nil {
					continue
				}
				c.Resize(l.list.itemMin)
				l.setupListItem(c, curItem)
			}

			c.Move(fyne.NewPos(x, y))
			if refresh {
				c.Resize(l.list.itemMin)
				if ok { // refresh visible
					l.setupListItem(c, curItem)
				}
			}

			x += l.list.itemMin.Width + theme.Padding()
			l.visible[curItem] = c
			cells = append(cells, c)
			curItem++
		}
		y += l.list.itemMin.Height + theme.Padding()
	}

	for id, old := range wasVisible {
		if _, ok := l.visible[id]; !ok {
			l.itemPool.Release(old)
		}
	}
	l.children = cells

	objects := l.children
	l.list.scroller.Content.(*fyne.Container).Objects = objects
}
