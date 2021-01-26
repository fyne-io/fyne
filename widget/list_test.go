package widget

import (
	"fmt"
	"image/color"
	"math"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"

	"github.com/stretchr/testify/assert"
)

func TestNewList(t *testing.T) {
	list := createList(1000)

	template := newListItem(fyne.NewContainerWithLayout(layout.NewHBoxLayout(), NewIcon(theme.DocumentIcon()), NewLabel("Template Object")), nil)
	firstItemIndex := test.WidgetRenderer(list).(*listRenderer).firstItemIndex
	lastItemIndex := test.WidgetRenderer(list).(*listRenderer).lastItemIndex
	visibleCount := len(test.WidgetRenderer(list).(*listRenderer).children)

	assert.Equal(t, 1000, list.Length())
	assert.GreaterOrEqual(t, list.MinSize().Width, template.MinSize().Width)
	assert.Equal(t, list.MinSize(), template.MinSize().Max(test.WidgetRenderer(list).(*listRenderer).scroller.MinSize()))
	assert.Equal(t, 0, firstItemIndex)
	assert.Equal(t, visibleCount, lastItemIndex-firstItemIndex+1)
}

func TestList_MinSize(t *testing.T) {
	for name, tt := range map[string]struct {
		cellSize        fyne.Size
		expectedMinSize fyne.Size
	}{
		"small": {
			fyne.NewSize(1, 1),
			fyne.NewSize(float32(32), float32(32)),
		},
		"large": {
			fyne.NewSize(100, 100),
			fyne.NewSize(100+3*theme.Padding(), 100+2*theme.Padding()),
		},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.expectedMinSize, NewList(
				func() int { return 5 },
				func() fyne.CanvasObject {
					r := canvas.NewRectangle(color.Black)
					r.SetMinSize(tt.cellSize)
					r.Resize(tt.cellSize)
					return r
				},
				func(ListItemID, fyne.CanvasObject) {}).MinSize())
		})
	}
}

func TestList_Resize(t *testing.T) {
	defer test.NewApp()
	list, w := setupList(t)
	template := newListItem(fyne.NewContainerWithLayout(layout.NewHBoxLayout(), NewIcon(theme.DocumentIcon()), NewLabel("Template Object")), nil)

	firstItemIndex := test.WidgetRenderer(list).(*listRenderer).firstItemIndex
	lastItemIndex := test.WidgetRenderer(list).(*listRenderer).lastItemIndex
	visibleCount := len(test.WidgetRenderer(list).(*listRenderer).children)
	assert.Equal(t, 0, firstItemIndex)
	assert.Equal(t, visibleCount, lastItemIndex-firstItemIndex+1)

	w.Resize(fyne.NewSize(200, 600))

	indexChange := int(math.Floor(float64(200) / float64(template.MinSize().Height)))

	newFirstItemIndex := test.WidgetRenderer(list).(*listRenderer).firstItemIndex
	newLastItemIndex := test.WidgetRenderer(list).(*listRenderer).lastItemIndex
	newVisibleCount := len(test.WidgetRenderer(list).(*listRenderer).children)

	assert.Equal(t, firstItemIndex, newFirstItemIndex)
	assert.NotEqual(t, lastItemIndex, newLastItemIndex)
	assert.Equal(t, newLastItemIndex, lastItemIndex+indexChange)
	assert.NotEqual(t, visibleCount, newVisibleCount)
	assert.Equal(t, newVisibleCount, newLastItemIndex-newFirstItemIndex+1)
	test.AssertRendersToMarkup(t, "list/resized.xml", w.Canvas())
}

func TestList_OffsetChange(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	list := createList(1000)
	w := test.NewWindow(list)
	w.Resize(fyne.NewSize(200, 400))
	template := newListItem(fyne.NewContainerWithLayout(layout.NewHBoxLayout(), NewIcon(theme.DocumentIcon()), NewLabel("Template Object")), nil)

	firstItemIndex := test.WidgetRenderer(list).(*listRenderer).firstItemIndex
	lastItemIndex := test.WidgetRenderer(list).(*listRenderer).lastItemIndex
	visibleCount := test.WidgetRenderer(list).(*listRenderer).visibleItemCount

	assert.Equal(t, 0, firstItemIndex)
	assert.Equal(t, visibleCount, lastItemIndex-firstItemIndex)

	scroll := test.WidgetRenderer(list).(*listRenderer).scroller
	scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(0, -280)})

	indexChange := int(math.Floor(float64(300) / float64(template.MinSize().Height)))

	newFirstItemIndex := test.WidgetRenderer(list).(*listRenderer).firstItemIndex
	newLastItemIndex := test.WidgetRenderer(list).(*listRenderer).lastItemIndex
	newVisibleCount := test.WidgetRenderer(list).(*listRenderer).visibleItemCount

	assert.NotEqual(t, firstItemIndex, newFirstItemIndex)
	assert.Equal(t, newFirstItemIndex, firstItemIndex+indexChange-1)
	assert.NotEqual(t, lastItemIndex, newLastItemIndex)
	assert.Equal(t, newLastItemIndex, lastItemIndex+indexChange-1)
	assert.Equal(t, visibleCount, newVisibleCount)
	assert.Equal(t, newVisibleCount, newLastItemIndex-newFirstItemIndex)
	test.AssertRendersToMarkup(t, "list/offset_changed.xml", w.Canvas())
}

func TestList_Hover(t *testing.T) {
	list := createList(1000)
	children := test.WidgetRenderer(list).(*listRenderer).children

	for i := 0; i < 2; i++ {
		assert.False(t, children[i].(*listItem).statusIndicator.Visible())
		children[i].(*listItem).MouseIn(&desktop.MouseEvent{})
		assert.Equal(t, children[i].(*listItem).statusIndicator.FillColor, theme.HoverColor())
		children[i].(*listItem).MouseOut()
		assert.False(t, children[i].(*listItem).statusIndicator.Visible())
	}
}

func TestList_Selection(t *testing.T) {
	list := createList(1000)
	children := test.WidgetRenderer(list).(*listRenderer).children

	assert.False(t, children[0].(*listItem).statusIndicator.Visible())
	children[0].(*listItem).Tapped(&fyne.PointEvent{})
	assert.Equal(t, children[0].(*listItem).statusIndicator.FillColor, theme.PrimaryColor())
	assert.True(t, children[0].(*listItem).statusIndicator.Visible())
	assert.Equal(t, 1, len(list.selected))
	assert.Equal(t, 0, list.selected[0])
	children[1].(*listItem).Tapped(&fyne.PointEvent{})
	assert.Equal(t, children[1].(*listItem).statusIndicator.FillColor, theme.PrimaryColor())
	assert.True(t, children[1].(*listItem).statusIndicator.Visible())
	assert.Equal(t, 1, len(list.selected))
	assert.Equal(t, 1, list.selected[0])
	assert.False(t, children[0].(*listItem).statusIndicator.Visible())
}

func TestList_Select(t *testing.T) {
	list := createList(1000)

	assert.Equal(t, test.WidgetRenderer(list).(*listRenderer).firstItemIndex, 0)
	list.Select(50)
	assert.Equal(t, test.WidgetRenderer(list).(*listRenderer).lastItemIndex, 50)
	children := test.WidgetRenderer(list).(*listRenderer).children
	assert.Equal(t, children[len(children)-1].(*listItem).statusIndicator.FillColor, theme.PrimaryColor())
	assert.True(t, children[len(children)-1].(*listItem).statusIndicator.Visible())

	list.Select(5)
	assert.Equal(t, test.WidgetRenderer(list).(*listRenderer).firstItemIndex, 5)
	children = test.WidgetRenderer(list).(*listRenderer).children
	assert.Equal(t, children[0].(*listItem).statusIndicator.FillColor, theme.PrimaryColor())
	assert.True(t, children[0].(*listItem).statusIndicator.Visible())

	list.Select(6)
	assert.Equal(t, test.WidgetRenderer(list).(*listRenderer).firstItemIndex, 5)
	children = test.WidgetRenderer(list).(*listRenderer).children
	assert.False(t, children[0].(*listItem).statusIndicator.Visible())
	assert.Equal(t, children[1].(*listItem).statusIndicator.FillColor, theme.PrimaryColor())
	assert.True(t, children[1].(*listItem).statusIndicator.Visible())
}

func TestList_Unselect(t *testing.T) {
	list := createList(1000)

	list.Select(10)
	children := test.WidgetRenderer(list).(*listRenderer).children
	assert.Equal(t, children[10].(*listItem).statusIndicator.FillColor, theme.PrimaryColor())
	assert.True(t, children[10].(*listItem).statusIndicator.Visible())

	list.Unselect(10)
	children = test.WidgetRenderer(list).(*listRenderer).children
	assert.False(t, children[10].(*listItem).statusIndicator.Visible())
	assert.Nil(t, list.selected)
}

func TestList_DataChange(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	list, w := setupList(t)
	children := test.WidgetRenderer(list).(*listRenderer).children

	assert.Equal(t, children[0].(*listItem).child.(*fyne.Container).Objects[1].(*Label).Text, "Test Item 0")
	changeData(list)
	list.Refresh()
	children = test.WidgetRenderer(list).(*listRenderer).children
	assert.Equal(t, children[0].(*listItem).child.(*fyne.Container).Objects[1].(*Label).Text, "a")
	test.AssertRendersToMarkup(t, "list/new_data.xml", w.Canvas())
}

func TestList_ThemeChange(t *testing.T) {
	defer test.NewApp()
	list, w := setupList(t)

	test.AssertImageMatches(t, "list/list_initial.png", w.Canvas().Capture())

	test.WithTestTheme(t, func() {
		time.Sleep(100 * time.Millisecond)
		list.Refresh()
		test.AssertImageMatches(t, "list/list_theme_changed.png", w.Canvas().Capture())
	})
}

func TestList_SmallList(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	var data []string
	data = append(data, "Test Item 0")

	list := NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return fyne.NewContainerWithLayout(layout.NewHBoxLayout(), NewIcon(theme.DocumentIcon()), NewLabel("Template Object"))
		},
		func(id ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*Label).SetText(data[id])
		},
	)
	w := test.NewWindow(list)
	w.Resize(fyne.NewSize(200, 400))

	visibleCount := len(test.WidgetRenderer(list).(*listRenderer).children)
	assert.Equal(t, visibleCount, 1)

	data = append(data, "Test Item 1")
	list.Refresh()

	visibleCount = len(test.WidgetRenderer(list).(*listRenderer).children)
	assert.Equal(t, visibleCount, 2)

	test.AssertRendersToMarkup(t, "list/small.xml", w.Canvas())
}

func TestList_ClearList(t *testing.T) {
	defer test.NewApp()
	list, w := setupList(t)
	assert.Equal(t, 1000, list.Length())

	firstItemIndex := test.WidgetRenderer(list).(*listRenderer).firstItemIndex
	lastItemIndex := test.WidgetRenderer(list).(*listRenderer).lastItemIndex
	visibleCount := len(test.WidgetRenderer(list).(*listRenderer).children)

	assert.Equal(t, visibleCount, lastItemIndex-firstItemIndex+1)

	list.Length = func() int {
		return 0
	}
	list.Refresh()

	visibleCount = len(test.WidgetRenderer(list).(*listRenderer).children)

	assert.Equal(t, visibleCount, 0)

	test.AssertRendersToMarkup(t, "list/cleared.xml", w.Canvas())
}

func TestList_RemoveItem(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	var data []string
	data = append(data, "Test Item 0")
	data = append(data, "Test Item 1")
	data = append(data, "Test Item 2")

	list := NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return fyne.NewContainerWithLayout(layout.NewHBoxLayout(), NewIcon(theme.DocumentIcon()), NewLabel("Template Object"))
		},
		func(id ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*Label).SetText(data[id])
		},
	)
	w := test.NewWindow(list)
	w.Resize(fyne.NewSize(200, 400))

	visibleCount := len(test.WidgetRenderer(list).(*listRenderer).children)
	assert.Equal(t, visibleCount, 3)

	data = data[:len(data)-1]
	list.Refresh()

	visibleCount = len(test.WidgetRenderer(list).(*listRenderer).children)
	assert.Equal(t, visibleCount, 2)
	test.AssertRendersToMarkup(t, "list/item_removed.xml", w.Canvas())
}

func TestList_NoFunctionsSet(t *testing.T) {
	list := &List{}
	w := test.NewWindow(list)
	w.Resize(fyne.NewSize(200, 400))
	list.Refresh()
}

func createList(items int) *List {
	var data []string
	for i := 0; i < items; i++ {
		data = append(data, fmt.Sprintf("Test Item %d", i))
	}

	list := NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return fyne.NewContainerWithLayout(layout.NewHBoxLayout(), NewIcon(theme.DocumentIcon()), NewLabel("Template Object"))
		},
		func(id ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*Label).SetText(data[id])
		},
	)
	list.Resize(fyne.NewSize(200, 1000))
	return list
}

func changeData(list *List) {
	data := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	list.Length = func() int {
		return len(data)
	}
	list.UpdateItem = func(id ListItemID, item fyne.CanvasObject) {
		item.(*fyne.Container).Objects[1].(*Label).SetText(data[id])
	}
}

func setupList(t *testing.T) (*List, fyne.Window) {
	test.NewApp()
	list := createList(1000)
	w := test.NewWindow(list)
	w.Resize(fyne.NewSize(200, 400))
	test.AssertRendersToMarkup(t, "list/initial.xml", w.Canvas())
	return list, w
}
