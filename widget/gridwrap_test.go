package widget

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"github.com/stretchr/testify/assert"
)

func TestGridWrap_Focus(t *testing.T) {
	test.NewTempApp(t)
	list := createGridWrap(100)
	window := test.NewWindow(list)
	defer window.Close()
	window.Resize(list.MinSize().Max(fyne.NewSize(150, 200)))

	canvas := window.Canvas().(test.WindowlessCanvas)
	assert.Nil(t, canvas.Focused())

	canvas.FocusNext()
	assert.NotNil(t, canvas.Focused())
	assert.Equal(t, 0, canvas.Focused().(*GridWrap).currentHighlight)

	children := list.scroller.Content.(*fyne.Container).Objects
	assert.True(t, children[0].(*gridWrapItem).hovered)
	assert.False(t, children[1].(*gridWrapItem).hovered)
	assert.False(t, children[6].(*gridWrapItem).hovered)
	assert.False(t, children[7].(*gridWrapItem).hovered)

	list.TypedKey(&fyne.KeyEvent{Name: fyne.KeyDown})
	assert.False(t, children[0].(*gridWrapItem).hovered)
	assert.False(t, children[1].(*gridWrapItem).hovered)
	assert.True(t, children[6].(*gridWrapItem).hovered)
	assert.False(t, children[7].(*gridWrapItem).hovered)

	list.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
	assert.False(t, children[0].(*gridWrapItem).hovered)
	assert.False(t, children[1].(*gridWrapItem).hovered)
	assert.False(t, children[6].(*gridWrapItem).hovered)
	assert.True(t, children[7].(*gridWrapItem).hovered)

	list.TypedKey(&fyne.KeyEvent{Name: fyne.KeyLeft})
	assert.False(t, children[0].(*gridWrapItem).hovered)
	assert.False(t, children[1].(*gridWrapItem).hovered)
	assert.True(t, children[6].(*gridWrapItem).hovered)
	assert.False(t, children[7].(*gridWrapItem).hovered)

	list.TypedKey(&fyne.KeyEvent{Name: fyne.KeyUp})
	assert.True(t, children[0].(*gridWrapItem).hovered)
	assert.False(t, children[1].(*gridWrapItem).hovered)
	assert.False(t, children[6].(*gridWrapItem).hovered)
	assert.False(t, children[7].(*gridWrapItem).hovered)

	canvas.Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeySpace})
	assert.True(t, children[0].(*gridWrapItem).selected)
}

func TestGridWrap_New(t *testing.T) {
	g := createGridWrap(1000)
	template := NewIcon(theme.AccountIcon())

	assert.Equal(t, 1000, g.Length())
	assert.GreaterOrEqual(t, g.MinSize().Width, template.MinSize().Width)
	assert.Equal(t, float32(0), g.offsetY)
}

func TestGridWrap_OffsetChange(t *testing.T) {
	g := createGridWrap(1000)

	assert.Equal(t, float32(0), g.offsetY)

	g.scroller.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(0, -280)})

	assert.NotEqual(t, 0, g.offsetY)
}

func TestGridWrap_ScrollTo(t *testing.T) {
	g := createGridWrap(1000)

	// override update item to keep track of greatest item rendered
	oldUpdateFunc := g.UpdateItem
	var greatest GridWrapItemID = -1
	g.UpdateItem = func(id GridWrapItemID, item fyne.CanvasObject) {
		if id > greatest {
			greatest = id
		}
		oldUpdateFunc(id, item)
	}

	g.ScrollTo(650)
	assert.GreaterOrEqual(t, greatest, 650)

	g.ScrollTo(800)
	assert.GreaterOrEqual(t, greatest, 800)

	g.ScrollToBottom()
	assert.Equal(t, greatest, GridWrapItemID(999))
}

func TestGridWrap_ScrollToItem(t *testing.T) {
	g := createGridWrap(1000)

	// override update item to keep track of greatest item rendered
	oldUpdateFunc := g.UpdateItem
	var greatest GridWrapItemID = -1
	g.UpdateItem = func(id GridWrapItemID, item fyne.CanvasObject) {
		if id > greatest {
			greatest = id
		}
		oldUpdateFunc(id, item)
	}

	g.ScrollToItem(650)
	assert.GreaterOrEqual(t, greatest, 650)
	assert.Equal(t, 650, g.currentHighlight)

	g.ScrollToItem(800)
	assert.GreaterOrEqual(t, greatest, 800)
	assert.Equal(t, 800, g.currentHighlight)

	g.ScrollToItem(999)
	assert.Equal(t, greatest, 999)
	assert.Equal(t, 999, g.currentHighlight)

	g.ScrollToItem(1001)
	assert.Equal(t, greatest, 999)
	assert.Equal(t, 999, g.currentHighlight)

	g.ScrollToItem(0)
	assert.GreaterOrEqual(t, greatest, 0)
	assert.Equal(t, 0, g.currentHighlight)

	g.ScrollToItem(-1)
	assert.GreaterOrEqual(t, greatest, 0)
	assert.Equal(t, 0, g.currentHighlight)
}

func TestGridWrap_ScrollToOffset(t *testing.T) {
	g := createGridWrap(10)
	g.Resize(fyne.NewSize(10, 10))

	g.ScrollToOffset(2)
	assert.Equal(t, float32(2), g.GetScrollOffset())

	g.ScrollToOffset(-20)
	assert.Equal(t, float32(0), g.GetScrollOffset())

	g.ScrollToOffset(10000)
	assert.LessOrEqual(t, g.GetScrollOffset(), float32(500) /*upper bound on content height*/)

	// GridWrap viewport is larger than content size
	g.Resize(fyne.NewSize(50, 250))
	g.ScrollToOffset(20)
	assert.Equal(t, float32(0), g.GetScrollOffset()) // doesn't scroll
}

func TestGridWrap_ScrollToTop(t *testing.T) {
	g := createGridWrap(1000)
	g.ScrollTo(750)
	assert.NotEqual(t, float32(0), g.offsetY)
	g.ScrollToTop()
	assert.Equal(t, float32(0), g.offsetY)
}

func createGridWrap(items int) *GridWrap {
	data := make([]fyne.Resource, items)
	for i := 0; i < items; i++ {
		switch i % 10 {
		case 0:
			data[i] = theme.AccountIcon()
		case 1:
			data[i] = theme.CancelIcon()
		case 2:
			data[i] = theme.CheckButtonIcon()
		case 3:
			data[i] = theme.FileApplicationIcon()
		case 4:
			data[i] = theme.FileVideoIcon()
		case 5:
			data[i] = theme.DocumentIcon()
		case 6:
			data[i] = theme.MediaPlayIcon()
		case 7:
			data[i] = theme.MediaRecordIcon()
		case 8:
			data[i] = theme.FolderIcon()
		case 9:
			data[i] = theme.FolderOpenIcon()
		}
	}

	list := NewGridWrap(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			icon := NewIcon(theme.DocumentIcon())
			return icon
		},
		func(id GridWrapItemID, item fyne.CanvasObject) {
			item.(*Icon).SetResource(data[id])
		},
	)
	list.Resize(fyne.NewSize(200, 400))
	return list
}

func TestGridWrap_IndexIsInt(t *testing.T) {
	gw := &GridWrap{}

	// Both of these should be allowed to match List behaviour.
	// It allows the same update item function to be shared between both widgets if necessary.
	gw.UpdateItem = func(id GridWrapItemID, item fyne.CanvasObject) {}
	gw.UpdateItem = func(id int, item fyne.CanvasObject) {}
}

func TestGridWrap_RefreshItem(t *testing.T) {
	data := make([]string, 5)
	for i := 0; i < 5; i++ {
		data[i] = "Text"
	}

	list := NewGridWrap(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			icon := NewLabel("dummy")
			return icon
		},
		func(id GridWrapItemID, item fyne.CanvasObject) {
			item.(*Label).SetText(data[id])
		},
	)
	list.Resize(fyne.NewSize(50, 100))

	data[2] = "Replace"
	list.RefreshItem(2)

	children := list.scroller.Content.(*fyne.Container).Objects
	assert.Equal(t, "Text", children[1].(*gridWrapItem).child.(*Label).Text)
	assert.Equal(t, "Replace", children[2].(*gridWrapItem).child.(*Label).Text)
}

func TestGridWrap_Selection(t *testing.T) {
	g := createGridWrap(10)
	assert.Empty(t, g.selected)

	selected := -1
	unselected := -1
	g.OnSelected = func(id GridWrapItemID) {
		selected = id
		unselected = -1
	}

	g.OnUnselected = func(id GridWrapItemID) {
		selected = -1
		unselected = id
	}

	g.Select(0)
	assert.Len(t, g.selected, 1)
	assert.Zero(t, selected)
	assert.Equal(t, -1, unselected)

	g.UnselectAll()
	assert.Empty(t, g.selected)
	assert.Equal(t, -1, selected)
	assert.Zero(t, unselected)

	g.Select(9)
	assert.Len(t, g.selected, 1)
	assert.Equal(t, 9, selected)
	assert.Equal(t, -1, unselected)

	g.Unselect(9)
	assert.Empty(t, g.selected)
	assert.Equal(t, -1, selected)
	assert.Equal(t, 9, unselected)
}

func TestGridWrap_ResizeToSameSizeBeforeRender(t *testing.T) {
	g := NewGridWrap(
		func() int { return 1 },
		func() fyne.CanvasObject { return NewLabel("") },
		func(gwii GridWrapItemID, co fyne.CanvasObject) { co.(*Label).SetText("foo") },
	)
	// will not create renderer.
	// will crash if GridWrap scroller (not yet created) is accessed
	g.Resize(fyne.NewSize(0, 0))
}

func TestGridWrap_TypedKey(t *testing.T) {
	gridWrap := createGridWrap(20)
	window := test.NewWindow(gridWrap)
	defer window.Close()
	window.Resize(fyne.NewSize(80, 100))

	// want 3 columns to make assert navigaiton behavior
	assert.Equal(t, 3, gridWrap.ColumnCount())

	canvas := window.Canvas().(test.WindowlessCanvas)
	canvas.FocusNext()
	gridWrap.TypedKey(&fyne.KeyEvent{Name: fyne.KeyDown})
	assert.Equal(t, 3, gridWrap.currentHighlight)

	gridWrap.TypedKey(&fyne.KeyEvent{Name: fyne.KeyUp})
	assert.Equal(t, 0, gridWrap.currentHighlight)

	gridWrap.TypedKey(&fyne.KeyEvent{Name: fyne.KeyLeft})
	assert.Equal(t, 0, gridWrap.currentHighlight)

	gridWrap.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
	assert.Equal(t, 1, gridWrap.currentHighlight)

	gridWrap.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
	assert.Equal(t, 2, gridWrap.currentHighlight)

	gridWrap.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
	assert.Equal(t, 3, gridWrap.currentHighlight)

	gridWrap.TypedKey(&fyne.KeyEvent{Name: fyne.KeyLeft})
	assert.Equal(t, 2, gridWrap.currentHighlight)

	gridWrap.TypedKey(&fyne.KeyEvent{Name: fyne.KeySpace})
	assert.Equal(t, []int{2}, gridWrap.selected)

	gridWrap.currentHighlight = 20
	gridWrap.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
	assert.Equal(t, 20, gridWrap.currentHighlight)
}
