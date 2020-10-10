package widget

import (
	"fmt"
	"math"
	"testing"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
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
	assert.Equal(t, list.MinSize(), test.WidgetRenderer(list).(*listRenderer).scroller.MinSize())
	assert.Equal(t, 0, firstItemIndex)
	assert.Equal(t, visibleCount, lastItemIndex-firstItemIndex+1)
}

func TestList_Resize(t *testing.T) {
	list := createList(1000)
	w := test.NewWindow(list)
	w.Resize(fyne.NewSize(200, 400))
	template := newListItem(fyne.NewContainerWithLayout(layout.NewHBoxLayout(), NewIcon(theme.DocumentIcon()), NewLabel("Template Object")), nil)

	firstItemIndex := test.WidgetRenderer(list).(*listRenderer).firstItemIndex
	lastItemIndex := test.WidgetRenderer(list).(*listRenderer).lastItemIndex
	visibleCount := len(test.WidgetRenderer(list).(*listRenderer).children)
	assert.Equal(t, 0, firstItemIndex)
	assert.Equal(t, visibleCount, lastItemIndex-firstItemIndex+1)
	test.AssertImageMatches(t, "list/list_initial.png", w.Canvas().Capture())

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
	test.AssertImageMatches(t, "list/list_resized.png", w.Canvas().Capture())
}

func TestList_OffsetChange(t *testing.T) {
	list := createList(1000)
	w := test.NewWindow(list)
	w.Resize(fyne.NewSize(200, 400))
	template := newListItem(fyne.NewContainerWithLayout(layout.NewHBoxLayout(), NewIcon(theme.DocumentIcon()), NewLabel("Template Object")), nil)

	firstItemIndex := test.WidgetRenderer(list).(*listRenderer).firstItemIndex
	lastItemIndex := test.WidgetRenderer(list).(*listRenderer).lastItemIndex
	visibleCount := test.WidgetRenderer(list).(*listRenderer).visibleItemCount

	assert.Equal(t, 0, firstItemIndex)
	assert.Equal(t, visibleCount, lastItemIndex-firstItemIndex)
	test.AssertImageMatches(t, "list/list_initial.png", w.Canvas().Capture())

	scroll := test.WidgetRenderer(list).(*listRenderer).scroller
	scroll.Scrolled(&fyne.ScrollEvent{DeltaX: 0, DeltaY: -300})

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
	test.AssertImageMatches(t, "list/list_offset_changed.png", w.Canvas().Capture())
}

func TestList_Hover(t *testing.T) {
	list := createList(1000)
	children := test.WidgetRenderer(list).(*listRenderer).children

	for i := 0; i < 2; i++ {
		assert.Equal(t, children[i].(*listItem).statusIndicator.FillColor, theme.BackgroundColor())
		children[i].(*listItem).MouseIn(&desktop.MouseEvent{})
		assert.Equal(t, children[i].(*listItem).statusIndicator.FillColor, theme.HoverColor())
		children[i].(*listItem).MouseOut()
		assert.Equal(t, children[i].(*listItem).statusIndicator.FillColor, theme.BackgroundColor())
	}
}

func TestList_Selection(t *testing.T) {
	list := createList(1000)
	children := test.WidgetRenderer(list).(*listRenderer).children

	assert.Equal(t, children[0].(*listItem).statusIndicator.FillColor, theme.BackgroundColor())
	children[0].(*listItem).Tapped(&fyne.PointEvent{})
	assert.Equal(t, children[0].(*listItem).statusIndicator.FillColor, theme.FocusColor())
	assert.Equal(t, list.selectedIndex, 0)
	children[1].(*listItem).Tapped(&fyne.PointEvent{})
	assert.Equal(t, children[1].(*listItem).statusIndicator.FillColor, theme.FocusColor())
	assert.Equal(t, list.selectedIndex, 1)
	assert.Equal(t, children[0].(*listItem).statusIndicator.FillColor, theme.BackgroundColor())
}

func TestList_DataChange(t *testing.T) {
	list := createList(1000)
	w := test.NewWindow(list)
	w.Resize(fyne.NewSize(200, 400))
	children := test.WidgetRenderer(list).(*listRenderer).children

	assert.Equal(t, children[0].(*listItem).child.(*fyne.Container).Objects[1].(*Label).Text, "Test Item 0")
	test.AssertImageMatches(t, "list/list_initial.png", w.Canvas().Capture())
	changeData(list)
	list.Refresh()
	children = test.WidgetRenderer(list).(*listRenderer).children
	assert.Equal(t, children[0].(*listItem).child.(*fyne.Container).Objects[1].(*Label).Text, "a")
	test.AssertImageMatches(t, "list/list_new_data.png", w.Canvas().Capture())
}

func TestList_ThemeChange(t *testing.T) {
	list := createList(1000)
	w := test.NewWindow(list)
	w.Resize(fyne.NewSize(200, 400))

	test.AssertImageMatches(t, "list/list_initial.png", w.Canvas().Capture())

	test.WithTestTheme(t, func() {
		time.Sleep(100 * time.Millisecond)
		list.Refresh()
		test.AssertImageMatches(t, "list/list_theme_changed.png", w.Canvas().Capture())
	})
}

func TestList_SmallList(t *testing.T) {
	var data []string
	data = append(data, "Test Item 0")

	list := NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return fyne.NewContainerWithLayout(layout.NewHBoxLayout(), NewIcon(theme.DocumentIcon()), NewLabel("Template Object"))
		},
		func(index int, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*Label).SetText(data[index])
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

	test.AssertImageMatches(t, "list/list_small_list.png", w.Canvas().Capture())
}

func TestList_ClearList(t *testing.T) {
	list := createList(1000)
	w := test.NewWindow(list)
	w.Resize(fyne.NewSize(200, 400))
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

	test.AssertImageMatches(t, "list/list_cleared.png", w.Canvas().Capture())
}

func TestList_RemoveItem(t *testing.T) {
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
		func(index int, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*Label).SetText(data[index])
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
	test.AssertImageMatches(t, "list/list_item_removed.png", w.Canvas().Capture())
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
		func(index int, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*Label).SetText(data[index])
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
	list.UpdateItem = func(index int, item fyne.CanvasObject) {
		item.(*fyne.Container).Objects[1].(*Label).SetText(data[index])
	}
}
