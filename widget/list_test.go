package widget

import (
	"fmt"
	"math"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"github.com/stretchr/testify/assert"
)

func TestNewList(t *testing.T) {
	list := createList()

	template := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), NewIcon(theme.DocumentIcon()), NewLabel("Template Object"))
	firstItemIndex := test.WidgetRenderer(list).(*listRenderer).firstItemIndex
	lastItemIndex := test.WidgetRenderer(list).(*listRenderer).lastItemIndex
	visibleCount := len(test.WidgetRenderer(list).(*listRenderer).children)

	assert.Equal(t, 1000, list.Length())
	assert.GreaterOrEqual(t, list.MinSize().Width, template.MinSize().Width)
	assert.Equal(t, list.MinSize(), test.WidgetRenderer(list).(*listRenderer).scroller.MinSize())
	assert.Equal(t, 0, firstItemIndex)
	assert.Equal(t, visibleCount, lastItemIndex-firstItemIndex+1)
}

func TestListResize(t *testing.T) {
	list := createList()
	template := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), NewIcon(theme.DocumentIcon()), NewLabel("Template Object"))

	firstItemIndex := test.WidgetRenderer(list).(*listRenderer).firstItemIndex
	lastItemIndex := test.WidgetRenderer(list).(*listRenderer).lastItemIndex
	visibleCount := len(test.WidgetRenderer(list).(*listRenderer).children)
	poolCount := test.WidgetRenderer(list).(*listRenderer).itemPool.Count()

	assert.Equal(t, 0, firstItemIndex)
	assert.Equal(t, visibleCount, lastItemIndex-firstItemIndex+1)
	assert.Equal(t, poolCount, 5)

	list.Resize(fyne.NewSize(200, 2000))

	indexChange := int(math.Floor(float64(1000) / float64(template.MinSize().Height)))

	newFirstItemIndex := test.WidgetRenderer(list).(*listRenderer).firstItemIndex
	newLastItemIndex := test.WidgetRenderer(list).(*listRenderer).lastItemIndex
	newVisibleCount := len(test.WidgetRenderer(list).(*listRenderer).children)
	poolCount = test.WidgetRenderer(list).(*listRenderer).itemPool.Count()

	assert.Equal(t, firstItemIndex, newFirstItemIndex)
	assert.NotEqual(t, lastItemIndex, newLastItemIndex)
	assert.Equal(t, newLastItemIndex, lastItemIndex+indexChange)
	assert.NotEqual(t, visibleCount, newVisibleCount)
	assert.Equal(t, newVisibleCount, newLastItemIndex-newFirstItemIndex+1)
	assert.Equal(t, poolCount, 5)
}

func TestListOffsetChange(t *testing.T) {
	list := createList()
	template := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), NewIcon(theme.DocumentIcon()), NewLabel("Template Object"))
	test.WidgetRenderer(list).(*listRenderer).Layout(fyne.NewSize(100, 1000))

	firstItemIndex := test.WidgetRenderer(list).(*listRenderer).firstItemIndex
	lastItemIndex := test.WidgetRenderer(list).(*listRenderer).lastItemIndex
	visibleCount := test.WidgetRenderer(list).(*listRenderer).visibleItemCount
	poolCount := test.WidgetRenderer(list).(*listRenderer).itemPool.Count()

	assert.Equal(t, 0, firstItemIndex)
	assert.Equal(t, visibleCount, lastItemIndex-firstItemIndex)
	assert.Equal(t, poolCount, 5)

	scroll := test.WidgetRenderer(list).(*listRenderer).scroller
	scroll.Scrolled(&fyne.ScrollEvent{DeltaX: 0, DeltaY: -300})

	indexChange := int(math.Floor(float64(300) / float64(template.MinSize().Height)))

	newFirstItemIndex := test.WidgetRenderer(list).(*listRenderer).firstItemIndex
	newLastItemIndex := test.WidgetRenderer(list).(*listRenderer).lastItemIndex
	newVisibleCount := test.WidgetRenderer(list).(*listRenderer).visibleItemCount
	poolCount = test.WidgetRenderer(list).(*listRenderer).itemPool.Count()

	assert.NotEqual(t, firstItemIndex, newFirstItemIndex)
	assert.Equal(t, newFirstItemIndex, firstItemIndex+indexChange-1)
	assert.NotEqual(t, lastItemIndex, newLastItemIndex)
	assert.Equal(t, newLastItemIndex, lastItemIndex+indexChange-1)
	assert.Equal(t, visibleCount, newVisibleCount)
	assert.Equal(t, newVisibleCount, newLastItemIndex-newFirstItemIndex)
	assert.Equal(t, poolCount, 5)
}

func TestListHover(t *testing.T) {
	list := createList()
	children := test.WidgetRenderer(list).(*listRenderer).children

	for i := 0; i < 2; i++ {
		assert.Equal(t, children[i].(*listItem).statusIndicator.FillColor, theme.BackgroundColor())
		children[i].(*listItem).MouseIn(&desktop.MouseEvent{})
		assert.Equal(t, children[i].(*listItem).statusIndicator.FillColor, theme.HoverColor())
		children[i].(*listItem).MouseOut()
		assert.Equal(t, children[i].(*listItem).statusIndicator.FillColor, theme.BackgroundColor())
	}
}

func TestListSelection(t *testing.T) {
	list := createList()
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

func TestListDataChange(t *testing.T) {
	list := createList()
	children := test.WidgetRenderer(list).(*listRenderer).children

	assert.Equal(t, children[0].(*listItem).child.(*fyne.Container).Objects[1].(*Label).Text, "Test Item 0")
	changeData(list)
	list.Refresh()
	children = test.WidgetRenderer(list).(*listRenderer).children
	assert.Equal(t, children[0].(*listItem).child.(*fyne.Container).Objects[1].(*Label).Text, "a")
}

func createList() *List {
	var data []string
	for i := 0; i < 1000; i++ {
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
	var data []string
	data = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	list.Length = func() int {
		return len(data)
	}
	list.UpdateItem = func(index int, item fyne.CanvasObject) {
		item.(*fyne.Container).Objects[1].(*Label).SetText(data[index])
	}
}
