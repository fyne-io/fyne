package widget

import (
	"fmt"
	"image/color"
	"math"
	"testing"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
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
			fyne.NewSize(scrollContainerMinSize, scrollContainerMinSize),
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
	test.AssertRendersToMarkup(t, `
		<canvas padded size="200x600">
			<content>
				<widget pos="4,4" size="192x592" type="*widget.List">
					<widget size="192x592" type="*widget.ScrollContainer">
						<container size="192x37999">
							<widget size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">Test Item 0</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,38" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">Test Item 1</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,76" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">Test Item 2</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,114" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">Test Item 3</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,152" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">Test Item 4</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,190" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">Test Item 5</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,228" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">Test Item 6</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,266" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">Test Item 7</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,304" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">Test Item 8</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,342" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">Test Item 9</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,380" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="97x29" type="*widget.Label">
										<text pos="4,4" size="89x21">Test Item 10</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,418" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="97x29" type="*widget.Label">
										<text pos="4,4" size="89x21">Test Item 11</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,456" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="97x29" type="*widget.Label">
										<text pos="4,4" size="89x21">Test Item 12</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,494" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="97x29" type="*widget.Label">
										<text pos="4,4" size="89x21">Test Item 13</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,532" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="97x29" type="*widget.Label">
										<text pos="4,4" size="89x21">Test Item 14</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,570" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="97x29" type="*widget.Label">
										<text pos="4,4" size="89x21">Test Item 15</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,608" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="97x29" type="*widget.Label">
										<text pos="4,4" size="89x21">Test Item 16</text>
									</widget>
								</container>
							</widget>
							<widget size="0x0" type="*widget.Separator">
								<rectangle fillColor="disabled" size="0x0"/>
							</widget>
							<widget pos="4,37" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,75" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,113" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,151" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,189" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,227" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,265" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,303" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,341" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,379" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,417" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,455" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,493" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,531" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,569" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,607" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
						</container>
						<widget pos="186,0" size="6x592" type="*widget.scrollBarArea">
							<widget backgroundColor="scrollbar" pos="3,0" size="3x16" type="*widget.scrollBar">
							</widget>
						</widget>
						<widget pos="0,592" size="192x0" type="*widget.Shadow">
							<linearGradient endColor="shadow" pos="0,-8" size="192x8"/>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, w.Canvas())
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
	scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(0, -300)})

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
	test.AssertRendersToMarkup(t, `
		<canvas padded size="200x400">
			<content>
				<widget pos="4,4" size="192x392" type="*widget.List">
					<widget size="192x392" type="*widget.ScrollContainer">
						<container pos="0,-300" size="192x37999">
							<widget pos="0,266" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">Test Item 7</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,304" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">Test Item 8</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,342" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">Test Item 9</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,380" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="97x29" type="*widget.Label">
										<text pos="4,4" size="89x21">Test Item 10</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,418" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="97x29" type="*widget.Label">
										<text pos="4,4" size="89x21">Test Item 11</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,456" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">Test Item 12</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,494" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">Test Item 13</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,532" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">Test Item 14</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,570" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">Test Item 15</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,608" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">Test Item 16</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,646" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">Test Item 17</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,684" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">Test Item 18</text>
									</widget>
								</container>
							</widget>
							<widget size="0x0" type="*widget.Separator">
								<rectangle fillColor="disabled" size="0x0"/>
							</widget>
							<widget pos="4,303" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,341" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,379" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,417" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,455" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,493" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,531" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,569" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,607" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,645" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,683" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
						</container>
						<widget pos="186,0" size="6x392" type="*widget.scrollBarArea">
							<widget backgroundColor="scrollbar" pos="3,2" size="3x16" type="*widget.scrollBar">
							</widget>
						</widget>
						<widget size="192x0" type="*widget.Shadow">
							<linearGradient size="192x8" startColor="shadow"/>
						</widget>
						<widget pos="0,392" size="192x0" type="*widget.Shadow">
							<linearGradient endColor="shadow" pos="0,-8" size="192x8"/>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, w.Canvas())
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
	assert.Equal(t, 1, len(list.selected))
	assert.Equal(t, 0, list.selected[0])
	children[1].(*listItem).Tapped(&fyne.PointEvent{})
	assert.Equal(t, children[1].(*listItem).statusIndicator.FillColor, theme.FocusColor())
	assert.Equal(t, 1, len(list.selected))
	assert.Equal(t, 1, list.selected[0])
	assert.Equal(t, children[0].(*listItem).statusIndicator.FillColor, theme.BackgroundColor())
}

func TestList_Select(t *testing.T) {
	list := createList(1000)

	assert.Equal(t, test.WidgetRenderer(list).(*listRenderer).firstItemIndex, 0)
	list.Select(50)
	assert.Equal(t, test.WidgetRenderer(list).(*listRenderer).lastItemIndex, 50)
	children := test.WidgetRenderer(list).(*listRenderer).children
	assert.Equal(t, children[len(children)-1].(*listItem).statusIndicator.FillColor, theme.FocusColor())

	list.Select(5)
	assert.Equal(t, test.WidgetRenderer(list).(*listRenderer).firstItemIndex, 5)
	children = test.WidgetRenderer(list).(*listRenderer).children
	assert.Equal(t, children[0].(*listItem).statusIndicator.FillColor, theme.FocusColor())

	list.Select(6)
	assert.Equal(t, test.WidgetRenderer(list).(*listRenderer).firstItemIndex, 5)
	children = test.WidgetRenderer(list).(*listRenderer).children
	assert.Equal(t, children[0].(*listItem).statusIndicator.FillColor, theme.BackgroundColor())
	assert.Equal(t, children[1].(*listItem).statusIndicator.FillColor, theme.FocusColor())
}

func TestList_Unselect(t *testing.T) {
	list := createList(1000)

	list.Select(10)
	children := test.WidgetRenderer(list).(*listRenderer).children
	assert.Equal(t, children[10].(*listItem).statusIndicator.FillColor, theme.FocusColor())

	list.Unselect(10)
	children = test.WidgetRenderer(list).(*listRenderer).children
	assert.Equal(t, children[10].(*listItem).statusIndicator.FillColor, theme.BackgroundColor())
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
	test.AssertRendersToMarkup(t, `
		<canvas padded size="200x400">
			<content>
				<widget pos="4,4" size="192x392" type="*widget.List">
					<widget size="192x392" type="*widget.ScrollContainer">
						<container size="192x987">
							<widget size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">a</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,38" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">b</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,76" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">c</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,114" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">d</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,152" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">e</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,190" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">f</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,228" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">g</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,266" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">h</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,304" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">i</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,342" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">j</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,380" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="97x29" type="*widget.Label">
										<text pos="4,4" size="89x21">k</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,418" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="97x29" type="*widget.Label">
										<text pos="4,4" size="89x21">l</text>
									</widget>
								</container>
							</widget>
							<widget size="0x0" type="*widget.Separator">
								<rectangle fillColor="disabled" size="0x0"/>
							</widget>
							<widget pos="4,37" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,75" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,113" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,151" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,189" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,227" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,265" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,303" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,341" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,379" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,417" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
						</container>
						<widget pos="186,0" size="6x392" type="*widget.scrollBarArea">
							<widget backgroundColor="scrollbar" pos="3,0" size="3x155" type="*widget.scrollBar">
							</widget>
						</widget>
						<widget pos="0,392" size="192x0" type="*widget.Shadow">
							<linearGradient endColor="shadow" pos="0,-8" size="192x8"/>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, w.Canvas())
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

	test.AssertRendersToMarkup(t, `
		<canvas padded size="200x400">
			<content>
				<widget pos="4,4" size="192x392" type="*widget.List">
					<widget size="192x392" type="*widget.ScrollContainer">
						<container size="192x392">
							<widget size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">Test Item 0</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,38" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">Test Item 1</text>
									</widget>
								</container>
							</widget>
							<widget size="0x0" type="*widget.Separator">
								<rectangle fillColor="disabled" size="0x0"/>
							</widget>
							<widget pos="4,37" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
						</container>
					</widget>
				</widget>
			</content>
		</canvas>
	`, w.Canvas())
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

	test.AssertRendersToMarkup(t, `
		<canvas padded size="200x400">
			<content>
				<widget pos="4,4" size="192x392" type="*widget.List">
					<widget size="192x392" type="*widget.ScrollContainer">
						<container size="192x392">
						</container>
					</widget>
				</widget>
			</content>
		</canvas>
	`, w.Canvas())
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
	test.AssertRendersToMarkup(t, `
		<canvas padded size="200x400">
			<content>
				<widget pos="4,4" size="192x392" type="*widget.List">
					<widget size="192x392" type="*widget.ScrollContainer">
						<container size="192x392">
							<widget size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">Test Item 0</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,38" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">Test Item 1</text>
									</widget>
								</container>
							</widget>
							<widget size="0x0" type="*widget.Separator">
								<rectangle fillColor="disabled" size="0x0"/>
							</widget>
							<widget pos="4,37" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
						</container>
					</widget>
				</widget>
			</content>
		</canvas>
	`, w.Canvas())
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
	test.AssertRendersToMarkup(t, `
		<canvas padded size="200x400">
			<content>
				<widget pos="4,4" size="192x392" type="*widget.List">
					<widget size="192x392" type="*widget.ScrollContainer">
						<container size="192x37999">
							<widget size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">Test Item 0</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,38" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">Test Item 1</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,76" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">Test Item 2</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,114" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">Test Item 3</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,152" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">Test Item 4</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,190" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">Test Item 5</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,228" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">Test Item 6</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,266" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">Test Item 7</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,304" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">Test Item 8</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,342" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="88x29" type="*widget.Label">
										<text pos="4,4" size="80x21">Test Item 9</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,380" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="97x29" type="*widget.Label">
										<text pos="4,4" size="89x21">Test Item 10</text>
									</widget>
								</container>
							</widget>
							<widget pos="0,418" size="192x37" type="*widget.listItem">
								<rectangle fillColor="background" size="4x37"/>
								<container pos="8,4" size="180x29">
									<widget size="20x29" type="*widget.Icon">
										<image fillMode="contain" rsc="documentIcon" size="20x29"/>
									</widget>
									<widget pos="24,0" size="97x29" type="*widget.Label">
										<text pos="4,4" size="89x21">Test Item 11</text>
									</widget>
								</container>
							</widget>
							<widget size="0x0" type="*widget.Separator">
								<rectangle fillColor="disabled" size="0x0"/>
							</widget>
							<widget pos="4,37" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,75" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,113" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,151" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,189" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,227" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,265" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,303" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,341" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,379" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
							<widget pos="4,417" size="184x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="184x1"/>
							</widget>
						</container>
						<widget pos="186,0" size="6x392" type="*widget.scrollBarArea">
							<widget backgroundColor="scrollbar" pos="3,0" size="3x16" type="*widget.scrollBar">
							</widget>
						</widget>
						<widget pos="0,392" size="192x0" type="*widget.Shadow">
							<linearGradient endColor="shadow" pos="0,-8" size="192x8"/>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, w.Canvas())
	return list, w
}
