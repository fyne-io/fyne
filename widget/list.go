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

	// SelectionMode is the selection mode for the list
	//
	// Since: 2.5
	SelectionMode SelectionMode

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
	visible := lo.visible
	if item, ok := visible[id]; ok {
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

// Select adds the item identified by the given ID to the selection.
func (l *List) Select(id ListItemID) {
	l._select(id)
}

// select; return true if selected + refreshed
func (l *List) _select(id ListItemID) bool {
	sel := l.selected
	for _, selId := range sel {
		if id == selId {
			return false
		}
	}

	length := l.length()
	if id < 0 || id >= length {
		return false
	}
	l.selected = append(sel, id)
	defer func() {
		if f := l.OnSelected; f != nil {
			f(id)
		}
	}()
	l.scrollTo(id)
	l.RefreshItem(id)
	return true
}

// SelectOnly selects only the item identified by the given ID to the selection,
// unselecting any previously-selected items.
//
// Since: 2.5
func (l *List) SelectOnly(id ListItemID) {
	if len(l.selected) == 1 && id == l.selected[0] {
		return
	}
	length := l.length()
	if id < 0 || id >= length {
		return
	}
	old := l.selected
	l.selected = []ListItemID{id}
	wasPrevSelected := false
	defer func() {
		if f := l.OnUnselected; f != nil {
			for _, oldSelId := range old {
				if oldSelId != id {
					f(id)
				} else {
					// item represented by id was in prev. selection set
					wasPrevSelected = true
				}
			}
		}
		if f := l.OnSelected; f != nil && !wasPrevSelected {
			f(id)
		}
	}()
	l.scrollTo(id)
	l.Refresh()
}

// SelectAll selects all items in the list.
//
// Since: 2.5
func (l *List) SelectAll(id ListItemID) {
	length := l.length()
	if length == 0 || len(l.selected) == length {
		return
	}

	prev := l.selected
	l.selected = make([]int, length)
	for i := range l.selected {
		l.selected[i] = i
	}
	l.Refresh()

	// Call OnSelected callback for each newly selected item
	// TODO: this is O(n^2). improve?
	wasPrevSelected := func(id ListItemID) bool {
		for _, selId := range prev {
			if id == selId {
				return true
			}
		}
		return false
	}
	f := l.OnSelected
	if f == nil {
		return
	}
	for i := 0; i < length; i++ {
		l := ListItemID(i)
		if !wasPrevSelected(l) {
			f(l)
		}
	}
}

// SetSelection sets the currently selected items in the list
//
// Since: 2.5
func (l *List) SetSelection(selected []ListItemID) {
	length := l.length()
	if length == 0 {
		return
	}

	oldSel := l.selected
	newSel := make([]ListItemID, 0, len(selected))
	for _, id := range selected {
		if id >= 0 && id < length {
			newSel = append(newSel, id)
		}
	}
	l.selected = newSel
	l.Refresh()

	// Call OnSelected, OnUnselected callbacks for each newly (un)selected item
	// TODO: this is O(n^2). improve
	onSelected := l.OnSelected
	onUnselected := l.OnUnselected
	if onSelected == nil && onUnselected == nil {
		return
	}
	find := func(id ListItemID, idSet []ListItemID) bool {
		for _, selId := range oldSel {
			if id == selId {
				return true
			}
		}
		return false
	}
	for i := 0; i < length; i++ {
		id := ListItemID(i)
		wasSel := find(id, oldSel)
		isSel := find(id, newSel)
		if wasSel && !isSel && onUnselected != nil {
			onUnselected(id)
		} else if isSel && !wasSel && onSelected != nil {
			onSelected(id)
		}
	}
}

// ScrollTo scrolls to the item represented by id
//
// Since: 2.1
func (l *List) ScrollTo(id ListItemID) {
	length := l.length()
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
	length := l.length()
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
	maybeSelect := func() bool {
		d, ok := fyne.CurrentApp().Driver().(desktop.Driver)
		if ok && d.CurrentKeyModifiers()&fyne.KeyModifierShift > 0 {
			return l._select(l.currentFocus)
		}
		return false
	}

	switch event.Name {
	case fyne.KeySpace:
		if sel := l.SelectionMode; sel == SelectionSingle {
			l.SelectOnly(l.currentFocus)
		} else if sel == SelectionMultiple {
			l.handleMultiSelectAction(l.currentFocus)
		}
	case fyne.KeyDown:
		if f := l.Length; f != nil && l.currentFocus >= f()-1 {
			return
		}
		l.RefreshItem(l.currentFocus)
		l.currentFocus++
		if !maybeSelect() {
			l.scrollTo(l.currentFocus)
			l.RefreshItem(l.currentFocus)
		}
	case fyne.KeyUp:
		if l.currentFocus <= 0 {
			return
		}
		l.RefreshItem(l.currentFocus)
		l.currentFocus--
		if !maybeSelect() {
			l.scrollTo(l.currentFocus)
			l.RefreshItem(l.currentFocus)
		}
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
	// check if already not selected
	sel := l.selected
	selected := false
	for _, selID := range sel {
		if selID == id {
			selected = true
			break
		}
	}
	if !selected {
		return
	}

	newSel := make([]ListItemID, 0, len(sel)-1)
	for _, selID := range sel {
		if selID != id {
			newSel = append(newSel, selID)
		}
	}

	l.selected = newSel
	l.Refresh()
	if f := l.OnUnselected; f != nil {
		f(id)
	}
}

// UnselectAll removes all items from the selection.
//
// Since: 2.1
func (l *List) UnselectAll() {
	sel := l.selected
	if len(sel) == 0 {
		return
	}

	l.selected = nil
	l.Refresh()
	if f := l.OnUnselected; f != nil {
		for _, id := range sel {
			f(id)
		}
	}
}

// invariant: all of ids are valid and none are already selected
func (l *List) addToSelection(ids []int) {
	l.selected = append(l.selected, ids...)
	l.Refresh()
	if f := l.OnSelected; f != nil {
		for _, id := range ids {
			f(id)
		}
	}
}

func (l *List) handleMultiSelectAction(id ListItemID) {
	sel := l.selected
	isSelected := false
	for _, selID := range sel {
		if selID == id {
			isSelected = true
			break
		}
	}

	toggleSelect := func() {
		if isSelected {
			l.Unselect(id)
		} else {
			l.Select(id)
		}
	}

	desktopDriver, ok := fyne.CurrentApp().Driver().(desktop.Driver)
	if !ok {
		// simple toggle-select behavior for mobile
		toggleSelect()
		return
	}

	// for desktops:
	// *  (no modifier)  + click = select only
	// * ModifierDefault + click = toggle select
	// * ModifierShift   + click = select range
	mods := desktopDriver.CurrentKeyModifiers()
	if mods&fyne.KeyModifierShortcutDefault > 0 {
		toggleSelect()
	} else if mods&fyne.KeyModifierShift > 0 {
		if !isSelected {
			l.selectRange(id)
		}
	} else {
		if !isSelected {
			l.SelectOnly(id)
		}
	}
}

// select range between id and nearest existing selected item
func (l *List) selectRange(id ListItemID) {
	nearest, dist := l.findNearestSelectedItem(id)
	above := nearest < id
	if nearest == -1 || dist <= 1 {
		// either nothing selected, or something selected right next to id
		l.addToSelection([]int{id})
		return
	}
	selAdd := make([]int, 0, dist-1)
	if above {
		// nearest selected item is above id
		for i := 0; i < dist; i++ {
			selAdd = append(selAdd, id-i)
		}
	} else {
		for i := 0; i < dist; i++ {
			selAdd = append(selAdd, id+i)
		}
	}
	l.addToSelection(selAdd)
}

func (l *List) findNearestSelectedItem(id ListItemID) (nearest ListItemID, dist int) {
	above, below := -1, math.MaxInt
	sel := l.selected
	length := l.length()

	for _, selId := range sel {
		if selId >= 0 && selId < id && selId > above {
			above = selId
		} else if selId < length && selId > id && selId < below {
			below = selId
		}
	}
	if above == -1 && below >= length {
		return -1, math.MaxInt // no selected item
	}
	dAbove, dBelow := id-above, below-id
	if above == -1 {
		dAbove = math.MaxInt
	}
	if dAbove <= dBelow {
		return above, dAbove
	}
	return below, dBelow
}

func (l *List) length() int {
	if f := l.Length; f != nil {
		return f()
	}
	return 0
}

func (l *List) visibleItemHeights(itemHeight float32, length int) (visible []float32, offY float32, minRow int) {
	rowOffset := float32(0)
	isVisible := false
	visible = []float32{}

	if l.scroller.Size().Height <= 0 {
		return
	}

	// theme.Padding is a slow call, so we cache it
	padding := theme.Padding()

	if len(l.itemHeights) == 0 {
		paddedItemHeight := itemHeight + padding

		offY = float32(math.Floor(float64(l.offsetY/paddedItemHeight))) * paddedItemHeight
		minRow = int(math.Floor(float64(offY / paddedItemHeight)))
		maxRow := int(math.Ceil(float64((offY + l.scroller.Size().Height) / paddedItemHeight)))

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

		visible = make([]float32, maxRow-minRow)
		for i := 0; i < maxRow-minRow; i++ {
			visible[i] = itemHeight
		}
		return
	}

	for i := 0; i < length; i++ {
		height := itemHeight
		if h, ok := l.itemHeights[i]; ok {
			height = h
		}

		if rowOffset <= l.offsetY-height-padding {
			// before scroll
		} else if rowOffset <= l.offsetY {
			minRow = i
			offY = rowOffset
			isVisible = true
		}
		if rowOffset >= l.offsetY+l.scroller.Size().Height {
			break
		}

		rowOffset += height + padding
		if isVisible {
			visible = append(visible, height)
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
	sel := l.list.selected
	for _, s := range sel {
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
		if sel := l.list.SelectionMode; sel == SelectionSingle {
			l.list.SelectOnly(id)
		} else if sel == SelectionMultiple {
			l.list.handleMultiSelectAction(id)
		}
	}
}

func (l *listLayout) updateList(newOnly bool) {
	l.renderLock.Lock()
	separatorThickness := theme.Padding()
	width := l.list.Size().Width
	length := l.list.length()
	if l.list.UpdateItem == nil {
		fyne.LogError("Missing UpdateCell callback required for List", nil)
	}

	wasVisible := l.visible

	l.list.propertyLock.Lock()
	visibleRowHeights, offY, minRow := l.list.visibleItemHeights(l.list.itemMin.Height, length)
	l.list.propertyLock.Unlock()
	if len(visibleRowHeights) == 0 && length > 0 { // we can't show anything until we have some dimensions
		l.renderLock.Unlock() // user code should not be locked
		return
	}

	visible := make(map[ListItemID]*listItem, len(visibleRowHeights))
	cells := make([]fyne.CanvasObject, len(visibleRowHeights))

	y := offY
	for index, itemHeight := range visibleRowHeights {
		row := index + minRow
		size := fyne.NewSize(width, itemHeight)

		c, ok := wasVisible[row]
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
		visible[row] = c
		cells[index] = c
	}

	l.visible = visible

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
	l.renderLock.Unlock() // user code should not be locked

	if newOnly {
		for row, obj := range visible {
			if _, ok := wasVisible[row]; !ok {
				l.setupListItem(obj, row, l.list.focused && l.list.currentFocus == row)
			}
		}
	} else {
		for row, obj := range visible {
			l.setupListItem(obj, row, l.list.focused && l.list.currentFocus == row)
		}
	}
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
