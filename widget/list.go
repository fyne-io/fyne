package widget

import (
	"fmt"
	"math"
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

// Declare conformity with Widget interface.
var _ fyne.Widget = (*List)(nil)

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

	scroller      *widget.Scroll
	selected      []ListItemID
	itemMin       fyne.Size
	offsetY       float32
	offsetUpdated func(fyne.Position)
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

	if f := l.CreateItem; f != nil {
		if l.itemMin.IsZero() {
			l.itemMin = newListItem(f(), nil).MinSize()
		}
	}
	layout := &fyne.Container{}
	l.scroller = widget.NewVScroll(layout)
	layout.Layout = newListLayout(l)
	layout.Resize(layout.MinSize())
	objects := []fyne.CanvasObject{l.scroller}
	lr := newListRenderer(objects, l, l.scroller, layout)
	return lr
}

// MinSize returns the size that this widget should not shrink below.
func (l *List) MinSize() fyne.Size {
	l.ExtendBaseWidget(l)

	return l.BaseWidget.MinSize()
}

func (l *List) scrollTo(id ListItemID) {
	if l.scroller == nil {
		return
	}
	y := (float32(id) * l.itemMin.Height) + (float32(id) * theme.SeparatorThicknessSize())
	if y < l.scroller.Offset.Y {
		l.scroller.Offset.Y = y
	} else if y+l.itemMin.Height > l.scroller.Offset.Y+l.scroller.Size().Height {
		l.scroller.Offset.Y = y + l.itemMin.Height - l.scroller.Size().Height
	}
	l.offsetUpdated(l.scroller.Offset)
}

// Resize is called when this list should change size. We refresh to ensure invisible items are drawn.
func (l *List) Resize(s fyne.Size) {
	l.BaseWidget.Resize(s)
	l.offsetUpdated(l.scroller.Offset)
	l.scroller.Content.(*fyne.Container).Layout.(*listLayout).updateList(true)
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
		l.list.itemMin = newListItem(f(), nil).MinSize()
	}
	l.Layout(l.list.Size())
	l.scroller.Refresh()
	l.layout.Layout.(*listLayout).updateList(true)
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

type listLayout struct {
	list       *List
	separators []fyne.CanvasObject
	children   []fyne.CanvasObject

	itemPool   *syncPool
	visible    map[ListItemID]*listItem
	renderLock sync.Mutex
}

func newListLayout(list *List) fyne.Layout {
	l := &listLayout{list: list, itemPool: &syncPool{}, visible: make(map[ListItemID]*listItem)}
	list.offsetUpdated = l.offsetUpdated
	return l
}

func (l *listLayout) Layout([]fyne.CanvasObject, fyne.Size) {
	l.updateList(true)
}

func (l *listLayout) MinSize([]fyne.CanvasObject) fyne.Size {
	if f := l.list.Length; f != nil {
		separatorThickness := theme.SeparatorThicknessSize()
		return fyne.NewSize(l.list.itemMin.Width,
			(l.list.itemMin.Height+separatorThickness)*float32(f())-separatorThickness)
	}
	return fyne.NewSize(0, 0)
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
	l.updateList(false)
}

func (l *listLayout) setupListItem(li *listItem, id ListItemID) {
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

func (l *listLayout) updateList(refresh bool) {
	l.renderLock.Lock()
	defer l.renderLock.Unlock()
	separatorThickness := theme.SeparatorThicknessSize()
	width := l.list.Size().Width
	length := 0
	if f := l.list.Length; f != nil {
		length = f()
	}
	visibleItemCount := int(math.Ceil(float64(l.list.scroller.Size().Height)/float64(l.list.itemMin.Height+theme.SeparatorThicknessSize()))) + 1
	offY := l.list.offsetY - float32(math.Mod(float64(l.list.offsetY), float64(l.list.itemMin.Height+separatorThickness)))
	minRow := ListItemID(offY / (l.list.itemMin.Height + separatorThickness))
	maxRow := ListItemID(fyne.Min(float32(minRow+visibleItemCount), float32(length)))

	if l.list.UpdateItem == nil {
		fyne.LogError("Missing UpdateCell callback required for List", nil)
	}

	wasVisible := l.visible
	l.visible = make(map[ListItemID]*listItem)
	var cells []fyne.CanvasObject
	y := offY
	size := fyne.NewSize(width, l.list.itemMin.Height)
	for row := minRow; row < maxRow; row++ {
		c, ok := wasVisible[row]
		if !ok {
			c = l.getItem()
			if c == nil {
				continue
			}
			c.Resize(size)
			l.setupListItem(c, row)
		}

		c.Move(fyne.NewPos(0, y))
		if refresh {
			c.Resize(size)
			if ok { // refresh visible
				l.setupListItem(c, row)
			}
		}

		y += l.list.itemMin.Height + separatorThickness
		l.visible[row] = c
		cells = append(cells, c)
	}

	for id, old := range wasVisible {
		if _, ok := l.visible[id]; !ok {
			l.itemPool.Release(old)
		}
	}
	l.children = cells
	l.updateSeparators()

	objects := l.children
	objects = append(objects, l.separators...)
	l.list.scroller.Content.(*fyne.Container).Objects = objects
}

func (l *listLayout) updateSeparators() {
	if len(l.children) > 1 {
		if len(l.separators) > len(l.children) {
			l.separators = l.separators[:len(l.children)]
		} else {
			for i := len(l.separators); i < len(l.children); i++ {
				l.separators = append(l.separators, NewSeparator())
			}
		}
	} else {
		l.separators = nil
	}

	separatorThickness := theme.SeparatorThicknessSize()
	for i, child := range l.children {
		if i == 0 {
			continue
		}
		l.separators[i].Move(fyne.NewPos(0, child.Position().Y-separatorThickness))
		l.separators[i].Resize(fyne.NewSize(l.list.Size().Width, separatorThickness))
		l.separators[i].Show()
	}
}
