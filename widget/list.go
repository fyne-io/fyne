package widget

import (
	"image/color"
	"log"
	"sync"

	"fyne.io/fyne"
	"fyne.io/fyne/binding"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal/cache"
	"fyne.io/fyne/theme"
)

var (
	listMinSize    = 32 // TODO consider the smallest useful list view?
	listMinVisible = 1.75
)

// Declare conformity with interfaces
var _ fyne.Scrollable = (*List)(nil)
var _ fyne.Tappable = (*List)(nil)
var _ fyne.Widget = (*List)(nil)

// List widget has a list of text labels and list check icons next to each.
// Changing the selection (only one can be selected) will trigger the changed func.
type List struct {
	CellSize   binding.Size
	Disabled   binding.Bool
	Hidden     binding.Bool
	Horizontal binding.Bool
	Items      binding.List
	Padding    binding.Int
	Pos        binding.Position
	Selected   binding.Int
	Siz        binding.Size

	index        binding.Int // Index of first visible item
	offsetItem   binding.Int // Offset from top-left of visible area to top-left of cell.
	offsetScroll binding.Int // Offset from top-left of visible area to top-left of widget.

	// TODO Add ScrollBars

	OnCreateCell func() fyne.CanvasObject
	OnBindCell   func(fyne.CanvasObject, binding.Binding)

	OnSelected func(int, binding.Binding)

	minSize  fyne.Size
	minSizes []fyne.Size
}

// NewList creates a new list widget with the set items and change handler
func NewList(items []string, selected func(int, binding.Binding)) *List {
	l := &List{
		Items:      binding.NewStringList(items...),
		OnSelected: selected,
	}
	return l
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (l *List) CreateRenderer() fyne.WidgetRenderer {
	log.Println("List.CreateRenderer")
	// Ensure all fields are set
	if l.CellSize == nil {
		l.CellSize = binding.EmptySize()
	}
	if l.Hidden == nil {
		l.Hidden = binding.EmptyBool()
	}
	if l.Horizontal == nil {
		l.Horizontal = binding.EmptyBool()
	}
	if l.index == nil {
		l.index = binding.EmptyInt()
	}
	if l.Items == nil {
		l.Items = binding.NewStringList()
	}
	if l.offsetItem == nil {
		l.offsetItem = binding.EmptyInt()
	}
	if l.offsetScroll == nil {
		l.offsetScroll = binding.EmptyInt()
	}
	if l.Padding == nil {
		l.Padding = binding.NewInt(theme.Padding())
	}
	if l.Pos == nil {
		l.Pos = binding.EmptyPosition()
	}
	if l.Selected == nil {
		l.Selected = binding.NewInt(-1) // Start with nothing selected
	}
	if l.Siz == nil {
		l.Siz = binding.EmptySize()
	}
	if l.OnCreateCell == nil {
		l.OnCreateCell = func() fyne.CanvasObject {
			return &canvas.Text{}
		}
	}
	if l.OnBindCell == nil {
		l.OnBindCell = func(object fyne.CanvasObject, data binding.Binding) {
			t, ok := object.(*canvas.Text)
			if ok {
				s, ok := data.(binding.String)
				if ok {
					t.Text = s.Get()
				}
				// TODO should theme use bindings?
				t.Color = theme.TextColor()
				t.TextSize = theme.TextSize()
				t.Show()
			}
		}
	}
	if l.Size().IsZero() {
		l.Resize(fyne.NewSize(listMinSize, listMinSize))
	}

	r := &listRenderer{
		list: l,
		done: make(chan bool),
		pool: sync.Pool{
			New: func() interface{} {
				return l.OnCreateCell()
			},
		},
	}
	// Create goroutine to listen to each field, respond to changes, and trigger refresh
	go func() {
		cellSizeChan := l.CellSize.Listen()
		hiddenChan := l.Hidden.Listen()
		horizontalChan := l.Horizontal.Listen()
		indexChan := l.index.Listen()
		itemsChan := l.Items.Listen()
		offsetItemChan := l.offsetItem.Listen()
		offsetScrollChan := l.offsetScroll.Listen()
		paddingChan := l.Padding.Listen()
		positionChan := l.Pos.Listen()
		selectedChan := l.Selected.Listen()
		sizeChan := l.Siz.Listen()
		for {
			select {
			case c := <-cellSizeChan:
				log.Println("CellSize:", c)
			case h := <-hiddenChan:
				log.Println("Hidden:", h)
				if h {
					continue
				}
			case h := <-horizontalChan:
				log.Println("Horizontal:", h)
			case i := <-indexChan:
				log.Println("IndexItem:", i)
				r.updateItems(i, l.Items.Length())
			case items := <-itemsChan:
				log.Println("Items:", items)
				r.updateItems(l.index.Get(), items)
			case o := <-offsetItemChan:
				log.Println("OffsetItem:", o)
				r.updateItems(l.index.Get(), l.Items.Length())
			case o := <-offsetScrollChan:
				log.Println("OffsetScroll:", o)
				r.updateItems(l.index.Get(), l.Items.Length())
			case p := <-paddingChan:
				log.Println("Padding:", p)
			case p := <-positionChan:
				log.Println("Position:", p)
			case s := <-selectedChan:
				log.Println("Selected:", s)
				if s >= 0 && s < l.Items.Length() && l.OnSelected != nil {
					l.OnSelected(s, l.Items.Get(s))
				}
			case s := <-sizeChan:
				log.Println("Size:", s)
				r.updateItems(l.index.Get(), l.Items.Length())
			case <-r.done:
				return
			}
			r.Layout(l.Size())
			canvas.Refresh(l)
		}
	}()

	return r
}

func (l *List) Enable() {
	if l.Disabled == nil {
		l.Disabled = binding.NewBool(false)
	} else {
		l.Disabled.Set(false)
	}
}

func (l *List) Disable() {
	if l.Disabled == nil {
		l.Disabled = binding.NewBool(true)
	} else {
		l.Disabled.Set(true)
	}
}

func (l *List) IsDisabled() (disabled bool) {
	if l.Disabled != nil {
		disabled = l.Disabled.Get()
	}
	return
}

func (l *List) Hide() {
	if l.Hidden == nil {
		l.Hidden = binding.NewBool(true)
	} else {
		l.Hidden.Set(true)
	}
}

// MinSize returns the size that this widget should not shrink below
func (l *List) MinSize() fyne.Size {
	min := cache.Renderer(l).(*listRenderer).MinSize()
	if l.Horizontal.Get() {
		// Ensure List is Wide Enough to Show Minimum Visible Items at Average Width
		min.Width = int(float64(min.Width) * listMinVisible)
	} else {
		// Ensure List is Tall Enough to Show Minimum Visible Items at Average Height
		min.Height = int(float64(min.Height) * listMinVisible)
	}
	l.minSize = l.minSize.Max(min)
	log.Println("List.MinSize:", l.minSize)
	return l.minSize
}

// MouseIn is called when a desktop pointer enters the widget
func (l *List) MouseIn(event *desktop.MouseEvent) {
	if l.IsDisabled() {
		return
	}
	// log.Println("List.MouseIn:", event)
	// TODO
}

// MouseOut is called when a desktop pointer exits the widget
func (l *List) MouseOut() {
	if l.IsDisabled() {
		return
	}
	// log.Println("List.MouseOut")
	// TODO
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (l *List) MouseMoved(event *desktop.MouseEvent) {
	if l.IsDisabled() {
		return
	}
	// log.Println("List.MouseMoved:", event)
	// TODO
}

func (l *List) Move(pos fyne.Position) {
	log.Println("List.Move:", pos)
	if l.Pos == nil {
		l.Pos = binding.NewPosition(pos)
	} else {
		l.Pos.Set(pos)
	}
}

func (l *List) Position() (pos fyne.Position) {
	if l.Pos != nil {
		pos = l.Pos.Get()
	}
	//log.Println("List.Position:", pos)
	return
}

func (l *List) Refresh() {
	log.Println("List.Refresh")
	// no-op
}

func (l *List) Resize(size fyne.Size) {
	log.Println("List.Resize:", size)
	if l.Siz == nil {
		l.Siz = binding.NewSize(size)
	} else {
		l.Siz.Set(size)
	}
}

// Scrolled is called when an input device triggers a scroll event
func (l *List) Scrolled(event *fyne.ScrollEvent) {
	log.Println("List.Scrolled:", event)

	var limit int
	horizontal := l.Horizontal.Get()
	padding := l.Padding.Get()
	index := l.index.Get()
	offsetItem := l.offsetItem.Get()
	offsetScroll := l.offsetScroll.Get()

	// Calculate limit and change offsets based on orientation
	if horizontal {
		limit = l.Size().Width
		offsetItem -= event.DeltaX
		offsetScroll -= event.DeltaX
	} else {
		limit = l.Size().Height
		offsetItem -= event.DeltaY
		offsetScroll -= event.DeltaY
	}

	// Scroll Down
	for s := l.sizeOf(index); offsetItem > s; s = l.sizeOf(index) {
		offsetItem -= s
		offsetItem -= padding
		index++
	}
	// Cap
	length := l.Items.Length()
	if index >= length {
		index = length - 1
	}
	usage := -offsetItem
	for i := index; i < length; i++ {
		usage += l.sizeOf(i)
		usage += padding
	}
	if usage < limit {
		// Gap after last element
		gap := limit - usage
		log.Println("Gap after last element:", gap)
		offsetItem -= gap
		offsetScroll -= gap
	}

	// Scroll Up
	for offsetItem < 0 {
		index--
		offsetItem += l.sizeOf(index)
		offsetItem += padding
	}
	// Cap
	if index < 0 {
		index = 0
		offsetItem = 0
		offsetScroll = 0
	}

	l.index.Set(index)
	l.offsetItem.Set(offsetItem)
	l.offsetScroll.Set(offsetScroll)
}

func (l *List) Show() {
	if l.Hidden == nil {
		l.Hidden = binding.NewBool(false)
	} else {
		l.Hidden.Set(false)
	}
}

func (l *List) Size() (size fyne.Size) {
	if l.Siz != nil {
		size = l.Siz.Get()
	}
	//log.Println("List.Size:", size)
	return
}

// Tapped is called when a pointer tapped event is captured and triggers any change handler
func (l *List) Tapped(event *fyne.PointEvent) {
	if l.IsDisabled() {
		return
	}
	log.Println("List.Tapped:", event)
	if l.OnSelected != nil {
		index := l.indexOf(event.Position.X, event.Position.Y)
		l.OnSelected(index, l.Items.Get(index))
	}
}

func (l *List) Visible() (visible bool) {
	if l.Hidden != nil {
		visible = !l.Hidden.Get()
	}
	//log.Println("List.Visible:", visible)
	return
}

func (l *List) indexOf(x, y int) (index int) {
	// TODO
	log.Println("List.indexOf:", x, y, index)
	return
}

func (l *List) sizeOf(index int) (size int) {
	if l.CellSize != nil {
		min := l.CellSize.Get()
		if l.Horizontal.Get() {
			size = min.Width
		} else {
			size = min.Height
		}
	} else {
		if index >= 0 && index < len(l.minSizes) {
			min := l.minSizes[index]
			if l.Horizontal.Get() {
				size = min.Width
			} else {
				size = min.Height
			}
		}
	}
	if size <= 0 {
		size = listMinSize
	}
	//log.Println("List.sizeOf:", index, size)
	return
}

type listRenderer struct {
	list     *List
	done     chan bool
	cells    []fyne.CanvasObject
	bindings []binding.Binding
	pool     sync.Pool
}

// BackgroundColor satisfies the fyne.WidgetRenderer interface.
func (r *listRenderer) BackgroundColor() color.Color {
	// log.Println("listRenderer.BackgroundColor")
	return theme.BackgroundColor()
}

// Destroy satisfies the fyne.WidgetRenderer interface.
func (r *listRenderer) Destroy() {
	// log.Println("listRenderer.Destroy")
	r.done <- true
}

// Layout the visible items of the list widget
func (r *listRenderer) Layout(size fyne.Size) {
	// log.Println("listRenderer.Layout:", size)
	horizontal := r.list.Horizontal.Get()
	index := r.list.index.Get()
	padding := r.list.Padding.Get()
	offsetItem := r.list.offsetItem.Get()
	// index := r.list.index.Get()
	var cellSize fyne.Size
	if r.list.CellSize != nil {
		cellSize = r.list.CellSize.Get()
	}
	var min fyne.Size
	x := 0
	y := 0
	if horizontal {
		x -= offsetItem
		for i, c := range r.cells {
			if cellSize.IsZero() {
				if index+i < len(r.list.minSizes) {
					min = r.list.minSizes[index+i]
				} else {
					min = c.MinSize()
				}
			} else {
				min = cellSize
			}
			pos := fyne.NewPos(x, y)
			siz := fyne.NewSize(min.Width, size.Height)
			// log.Println(i+index, pos, siz)
			c.Move(pos)
			c.Resize(siz)
			x += min.Width + padding
		}
	} else {
		y -= offsetItem
		for i, c := range r.cells {
			if cellSize.IsZero() {
				if index+i < len(r.list.minSizes) {
					min = r.list.minSizes[index+i]
				} else {
					min = c.MinSize()
				}
			} else {
				min = cellSize
			}
			pos := fyne.NewPos(x, y)
			siz := fyne.NewSize(size.Width, min.Height)
			// log.Println(i+index, pos, siz)
			c.Move(pos)
			c.Resize(siz)
			y += min.Height + padding
		}
	}
}

// MinSize calculates the largest minimum size of the visible list items.
func (r *listRenderer) MinSize() (size fyne.Size) {
	size.Width = listMinSize
	size.Height = listMinSize
	for _, c := range r.cells {
		size = size.Max(c.MinSize())
	}
	// log.Println("listRenderer.MinSize:", size)
	return
}

// Objects satisfies the fyne.WidgetRenderer interface.
func (r *listRenderer) Objects() []fyne.CanvasObject {
	return r.cells
}

// Refresh satisfies the fyne.WidgetRenderer interface.
func (r *listRenderer) Refresh() {
	// log.Println("listRenderer.Refresh")
	// no-op
}

func (r *listRenderer) updateItems(start, length int) {
	size := r.list.Size()
	// log.Println("listRenderer.updateItems:", start, length, size)
	// Create and bind cell for each visible item
	horizontal := r.list.Horizontal.Get()
	offset := r.list.offsetItem.Get()
	indexCell := 0
	indexItem := start
	x := -offset
	y := -offset
	var cellSize fyne.Size
	if r.list.CellSize != nil {
		cellSize = r.list.CellSize.Get()
	}
	var c fyne.CanvasObject
	var min fyne.Size
	for ; indexItem < length && x < size.Width && y < size.Height; indexCell, indexItem = indexCell+1, indexItem+1 {
		if indexCell < len(r.cells) {
			c = r.cells[indexCell]
		} else {
			c = r.pool.Get().(fyne.CanvasObject)
			r.cells = append(r.cells, c)
		}
		b := r.list.Items.Get(indexItem)
		if indexCell < len(r.bindings) {
			if bdg := r.bindings[indexCell]; b != bdg {
				r.list.OnBindCell(c, b)
				c.Refresh()
				r.bindings[indexCell] = b
			}
		} else {
			r.list.OnBindCell(c, b)
			c.Refresh()
			r.bindings = append(r.bindings, b)
		}
		if cellSize.IsZero() {
			min = c.MinSize()
			// log.Println(indexItem, min)
			if indexItem < len(r.list.minSizes) {
				r.list.minSizes[indexItem] = min
			} else {
				r.list.minSizes = append(r.list.minSizes, min)
			}
		} else {
			min = cellSize
		}
		if horizontal {
			x += min.Width
		} else {
			y += min.Height
		}
	}
	// Recycle unused cells
	lastCell := indexCell
	for ; indexCell < len(r.cells); indexCell++ {
		c := r.cells[indexCell]
		c.Hide()
		r.pool.Put(c)
	}
	r.cells = r.cells[:lastCell]
}
