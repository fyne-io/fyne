package widget

import (
	"image/color"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/theme"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewScrollContainer(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(10, 10))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	barArea := Renderer(scroll).(*scrollRenderer).vertArea
	render := Renderer(barArea).(*scrollBarAreaRenderer)

	assert.Equal(t, 0, scroll.Offset.Y)
	assert.Equal(t, fyne.NewSize(theme.ScrollBarSmallSize(), 100), render.barSizeVertical())
	assert.Equal(t, fyne.NewPos(100-theme.ScrollBarSmallSize(), 0), barArea.Position())
}

func TestScrollContainer_Refresh(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(1000, 1000))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	scroll.Scrolled(&fyne.ScrollEvent{DeltaY: -1000})

	assert.Equal(t, 900, scroll.Offset.Y)
	assert.Equal(t, fyne.NewSize(1000, 1000), rect.Size())
	rect.SetMinSize(fyne.NewSize(1000, 500))
	Refresh(scroll)
	assert.Equal(t, 400, scroll.Offset.Y)
	assert.Equal(t, fyne.NewSize(1000, 500), rect.Size())
}

func TestScrollContainer_Scrolled(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(1000, 1000))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(100, 100))

	assert.Equal(t, 0, scroll.Offset.Y)
	scroll.Scrolled(&fyne.ScrollEvent{DeltaY: -10})
	assert.Equal(t, 10, scroll.Offset.Y)
}

func TestScrollContainer_Scrolled_Limit(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(100, 100))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(80, 80))

	scroll.Scrolled(&fyne.ScrollEvent{DeltaY: -25})
	assert.Equal(t, 20, scroll.Offset.Y)
}

func TestScrollContainer_Scrolled_Back(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(1000, 1000))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	scroll.Offset.Y = 10

	scroll.Scrolled(&fyne.ScrollEvent{DeltaY: 10})
	assert.Equal(t, 0, scroll.Offset.Y)
}

func TestScrollContainer_Scrolled_BackLimit(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(1000, 1000))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	scroll.Offset.Y = 10

	scroll.Scrolled(&fyne.ScrollEvent{DeltaY: 20})
	assert.Equal(t, 0, scroll.Offset.Y)
}

func TestScrollContainer_Resize(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(100, 100))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(80, 80))

	scroll.Scrolled(&fyne.ScrollEvent{DeltaY: -20})
	scroll.Resize(fyne.NewSize(80, 100))
	assert.Equal(t, 0, scroll.Offset.Y)
}

func TestScrollContainer_ResizeOffset(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(100, 100))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(80, 80))

	scroll.Scrolled(&fyne.ScrollEvent{DeltaY: -20})
	scroll.Resize(fyne.NewSize(80, 90))
	assert.Equal(t, 10, scroll.Offset.Y)
}

func TestScrollContainer_ResizeExpand(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(100, 100))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(120, 140))

	assert.Equal(t, 120, rect.Size().Width)
	assert.Equal(t, 140, rect.Size().Height)
}

func TestScrollContainer_ScrollBarForSmallContentIsHidden(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(100, 100))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(100, 200))

	r := Renderer(scroll).(*scrollRenderer)
	assert.False(t, r.vertArea.Visible())
}

func TestScrollContainer_ShowHiddenScrollBarIfContentGrows(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(100, 100))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(100, 200))
	r := Renderer(scroll).(*scrollRenderer)
	require.False(t, r.vertArea.Visible())

	rect.SetMinSize(fyne.NewSize(100, 300))
	r.Layout(scroll.Size())
	assert.True(t, r.vertArea.Visible())
}

func TestScrollContainer_HideScrollBarIfContentShrinks(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(100, 300))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(100, 200))
	r := Renderer(scroll).(*scrollRenderer)
	require.True(t, r.vertArea.Visible())

	rect.SetMinSize(fyne.NewSize(100, 100))
	r.Layout(scroll.Size())
	assert.False(t, r.vertArea.Visible())
}

func TestScrollContainer_ScrollBarIsSmall(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(100, 500))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	r := Renderer(scroll).(*scrollRenderer)
	require.True(t, r.vertArea.Visible())
	require.Less(t, theme.ScrollBarSmallSize(), theme.ScrollBarSize())

	assert.Equal(t, theme.ScrollBarSmallSize(), r.vertArea.Size().Width)
	assert.Equal(t, theme.ScrollBarSmallSize(), Renderer(r.vertArea).(*scrollBarAreaRenderer).barSizeVertical().Width)
}

func TestScrollContainer_ScrollBarGrowsAndShrinksOnMouseInAndMouseOut(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(100, 500))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	area := Renderer(scroll).(*scrollRenderer).vertArea
	ar := Renderer(area).(*scrollBarAreaRenderer)
	bar := ar.bar
	require.True(t, bar.Visible())
	require.Less(t, theme.ScrollBarSmallSize(), theme.ScrollBarSize())
	require.Equal(t, theme.ScrollBarSmallSize(), bar.Size().Width)
	require.Equal(t, theme.ScrollBarSmallSize(), area.Size().Width)
	require.Equal(t, theme.ScrollBarSmallSize(), ar.barSizeVertical().Width)

	bar.MouseIn(&desktop.MouseEvent{})
	assert.Equal(t, theme.ScrollBarSize(), bar.Size().Width)
	assert.Equal(t, theme.ScrollBarSize(), ar.barSizeVertical().Width)

	bar.MouseOut()
	assert.Equal(t, theme.ScrollBarSmallSize(), bar.Size().Width)
	assert.Equal(t, theme.ScrollBarSmallSize(), ar.barSizeVertical().Width)
}

func TestScrollContainer_ShowShadowOnTopIfContentIsScrolled(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(100, 500))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	r := Renderer(scroll).(*scrollRenderer)
	assert.False(t, r.topShadow.Visible())
	assert.Equal(t, fyne.NewPos(0, 0), r.topShadow.Position())

	scroll.Scrolled(&fyne.ScrollEvent{DeltaY: -1})
	assert.True(t, r.topShadow.Visible())

	scroll.Scrolled(&fyne.ScrollEvent{DeltaY: 1})
	assert.False(t, r.topShadow.Visible())
}

func TestScrollContainer_ShowShadowOnBottomIfContentCanScroll(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(100, 500))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	r := Renderer(scroll).(*scrollRenderer)
	assert.True(t, r.bottomShadow.Visible())
	assert.Equal(t, scroll.size.Height, r.bottomShadow.Position().Y+r.bottomShadow.Size().Height)

	scroll.Scrolled(&fyne.ScrollEvent{DeltaY: -400})
	assert.False(t, r.bottomShadow.Visible())

	scroll.Scrolled(&fyne.ScrollEvent{DeltaY: 100})
	assert.True(t, r.bottomShadow.Visible())
}

func TestScrollBarRenderer_BarSize(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(100, 100))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	ar := Renderer(Renderer(scroll).(*scrollRenderer).vertArea).(*scrollBarAreaRenderer)

	assert.Equal(t, fyne.NewSize(theme.ScrollBarSmallSize(), 100), ar.barSizeVertical())

	// resize so content is twice our size. Bar should therefore be half again.
	scroll.Resize(fyne.NewSize(50, 50))
	assert.Equal(t, fyne.NewSize(theme.ScrollBarSmallSize(), 25), ar.barSizeVertical())
}

func TestScrollContainerRenderer_LimitBarSize(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(100, 100))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(120, 120))
	ar := Renderer(Renderer(scroll).(*scrollRenderer).vertArea).(*scrollBarAreaRenderer)

	assert.Equal(t, fyne.NewSize(theme.ScrollBarSmallSize(), 120), ar.barSizeVertical())
}

func TestScrollBar_Dragged_ClickedInside(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(1000, 1000))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	scrollBar := Renderer(Renderer(scroll).(*scrollRenderer).vertArea).(*scrollBarAreaRenderer).bar

	// Create drag event with starting position inside scroll rectangle area
	dragEvent := fyne.DragEvent{DraggedY: 20}
	dragEvent.Position = fyne.NewPos(10, 10)

	assert.Equal(t, 0, scroll.Offset.Y)
	scrollBar.Dragged(&dragEvent)
	assert.Equal(t, 200, scroll.Offset.Y)
}

func TestScrollBar_DraggedBack_ClickedInside(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(1000, 1000))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	scrollBar := Renderer(Renderer(scroll).(*scrollRenderer).vertArea).(*scrollBarAreaRenderer).bar

	// Drag forward
	dragEvent := fyne.DragEvent{DraggedY: 20}
	dragEvent.Position = fyne.NewPos(10, 10)
	scrollBar.Dragged(&dragEvent)

	// Drag back
	dragEvent = fyne.DragEvent{DraggedY: -10}
	dragEvent.Position = fyne.NewPos(10, 20)

	assert.Equal(t, 200, scroll.Offset.Y)
	scrollBar.Dragged(&dragEvent)
	assert.Equal(t, 100, scroll.Offset.Y)
}

func TestScrollBar_Dragged_Limit(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(1000, 1000))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	scrollBar := Renderer(Renderer(scroll).(*scrollRenderer).vertArea).(*scrollBarAreaRenderer).bar

	// Drag over limit
	dragEvent := fyne.DragEvent{DraggedY: 2000}
	dragEvent.Position = fyne.NewPos(10, 10)

	assert.Equal(t, 0, scroll.Offset.Y)
	scrollBar.Dragged(&dragEvent)
	assert.Equal(t, 900, scroll.Offset.Y)

	// Drag again
	dragEvent = fyne.DragEvent{DraggedY: 100}
	dragEvent.Position = fyne.NewPos(10, 90)

	// Offset doesn't go over limit
	assert.Equal(t, 900, scroll.Offset.Y)
	scrollBar.Dragged(&dragEvent)
	assert.Equal(t, 900, scroll.Offset.Y)
}

func TestScrollBar_Dragged_BackLimit(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(1000, 1000))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	scrollBar := Renderer(Renderer(scroll).(*scrollRenderer).vertArea).(*scrollBarAreaRenderer).bar

	// Drag over back limit
	dragEvent := fyne.DragEvent{DraggedY: -1000}
	dragEvent.Position = fyne.NewPos(10, 10)

	// Offset doesn't go over limit
	assert.Equal(t, 0, scroll.Offset.Y)
	scrollBar.Dragged(&dragEvent)
	assert.Equal(t, 0, scroll.Offset.Y)
}
