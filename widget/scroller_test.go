package widget

import (
	"image/color"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewScrollContainer(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(10, 10))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	barArea := test.WidgetRenderer(scroll).(*scrollContainerRenderer).vertArea
	bar := test.WidgetRenderer(barArea).(*scrollBarAreaRenderer).bar
	assert.Equal(t, 0, scroll.Offset.Y)
	assert.Equal(t, theme.ScrollBarSmallSize()*2, barArea.Size().Width)
	assert.Equal(t, theme.ScrollBarSmallSize(), bar.Size().Width)
	assert.Equal(t, theme.ScrollBarSmallSize(), bar.Position().X)
	assert.Equal(t, fyne.NewPos(100-theme.ScrollBarSmallSize()*2, 0), barArea.Position())
}

func TestScrollContainer_MinSize(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(500, 50))
	scroll := NewScrollContainer(rect)
	assert.Equal(t, fyne.NewSize(32, 32), scroll.MinSize())

	scrollMin := fyne.NewSize(100, 100)
	scroll.SetMinSize(scrollMin)
	test.WidgetRenderer(scroll).Layout(scroll.minSize)

	assert.Equal(t, scrollMin, scroll.MinSize())
	assert.Equal(t, fyne.NewSize(500, 100), rect.Size())
	assert.Equal(t, 0, scroll.Offset.X)
	assert.Equal(t, 0, scroll.Offset.Y)
}

func TestScrollContainer_ScrollToTop(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(500, 50))
	scroll := NewScrollContainer(rect)
	scroll.ScrollToTop()
	Y := scroll.Offset.Y
	assert.Equal(t, 0, Y)
}

func TestScrollContainer_ScrollToBottom(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(500, 50))
	scroll := NewScrollContainer(rect)
	scroll.ScrollToBottom()
	ExpectedY := 50
	Y := scroll.Content.Size().Height - scroll.Size().Height
	assert.Equal(t, ExpectedY, Y)
}

func TestScrollContainer_MinSize_Direction(t *testing.T) {
	t.Run("Both", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(100, 100))
		scroll := NewScrollContainer(rect)
		size := scroll.MinSize()
		assert.Equal(t, 32, size.Height)
		assert.Equal(t, 32, size.Width)
	})
	t.Run("HorizontalOnly", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(100, 100))
		scroll := NewHScrollContainer(rect)
		size := scroll.MinSize()
		assert.Equal(t, 100, size.Height)
		assert.Equal(t, 32, size.Width)
	})
	t.Run("VerticalOnly", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(100, 100))
		scroll := NewVScrollContainer(rect)
		size := scroll.MinSize()
		assert.Equal(t, 32, size.Height)
		assert.Equal(t, 100, size.Width)
	})
}

func TestScrollContainer_SetMinSize_Direction(t *testing.T) {
	t.Run("Both", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(100, 100))
		scroll := NewScrollContainer(rect)
		scroll.SetMinSize(fyne.NewSize(50, 50))
		size := scroll.MinSize()
		assert.Equal(t, 50, size.Height)
		assert.Equal(t, 50, size.Width)
	})
	t.Run("HorizontalOnly", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(100, 100))
		scroll := NewHScrollContainer(rect)
		scroll.SetMinSize(fyne.NewSize(50, 50))
		size := scroll.MinSize()
		assert.Equal(t, 100, size.Height)
		assert.Equal(t, 50, size.Width)
	})
	t.Run("VerticalOnly", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(100, 100))
		scroll := NewVScrollContainer(rect)
		scroll.SetMinSize(fyne.NewSize(50, 50))
		size := scroll.MinSize()
		assert.Equal(t, 50, size.Height)
		assert.Equal(t, 100, size.Width)
	})
}

func TestScrollContainer_Resize_Direction(t *testing.T) {
	t.Run("Both", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(100, 100))
		scroll := NewScrollContainer(rect)
		scroll.Resize(scroll.MinSize())
		size := scroll.Size()
		assert.Equal(t, 32, size.Height)
		assert.Equal(t, 32, size.Width)
	})
	t.Run("HorizontalOnly", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(100, 100))
		scroll := NewHScrollContainer(rect)
		scroll.Resize(scroll.MinSize())
		size := scroll.Size()
		assert.Equal(t, 100, size.Height)
		assert.Equal(t, 32, size.Width)
	})
	t.Run("VerticalOnly", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(100, 100))
		scroll := NewVScrollContainer(rect)
		scroll.Resize(scroll.MinSize())
		size := scroll.Size()
		assert.Equal(t, 32, size.Height)
		assert.Equal(t, 100, size.Width)
	})
}

func TestScrollContainer_Refresh(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(1000, 1000))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	assert.Equal(t, fyne.NewSize(1000, 1000), rect.Size())
	scroll.Scrolled(&fyne.ScrollEvent{DeltaX: -1000, DeltaY: -1000})
	assert.Equal(t, 900, scroll.Offset.X)
	assert.Equal(t, 900, scroll.Offset.Y)
	assert.Equal(t, fyne.NewSize(1000, 1000), rect.Size())
	rect.SetMinSize(fyne.NewSize(500, 500))
	Refresh(scroll)
	assert.Equal(t, 400, scroll.Offset.X)
	assert.Equal(t, fyne.NewSize(500, 500), rect.Size())

	rect2 := canvas.NewRectangle(color.White)
	scroll.Content = rect2
	scroll.Refresh()
	assert.Equal(t, rect2, test.WidgetRenderer(scroll).Objects()[0])
}

func TestScrollContainer_Scrolled(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(1000, 1000))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	assert.Equal(t, 0, scroll.Offset.X)
	assert.Equal(t, 0, scroll.Offset.Y)
	scroll.Scrolled(&fyne.ScrollEvent{DeltaX: -10, DeltaY: -10})
	assert.Equal(t, 10, scroll.Offset.X)
	assert.Equal(t, 10, scroll.Offset.Y)

}

func TestScrollContainer_Scrolled_Limit(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(100, 100))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(80, 80))
	scroll.Scrolled(&fyne.ScrollEvent{DeltaX: -25, DeltaY: -25})
	assert.Equal(t, 20, scroll.Offset.X)
}

func TestScrollContainer_Scrolled_Back(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(1000, 1000))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	scroll.Offset.X = 10
	scroll.Offset.Y = 10
	scroll.Scrolled(&fyne.ScrollEvent{DeltaX: 10, DeltaY: 10})
	assert.Equal(t, 0, scroll.Offset.X)
	assert.Equal(t, 0, scroll.Offset.Y)
}

func TestScrollContainer_Scrolled_BackLimit(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	scroll := NewScrollContainer(rect)
	rect.SetMinSize(fyne.NewSize(1000, 1000))
	scroll.Resize(fyne.NewSize(100, 100))
	scroll.Offset.X = 10
	scroll.Offset.Y = 10
	scroll.Scrolled(&fyne.ScrollEvent{DeltaX: 20, DeltaY: 20})
	assert.Equal(t, 0, scroll.Offset.X)
	assert.Equal(t, 0, scroll.Offset.Y)

}

func TestScrollContainer_Resize(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	scroll := NewScrollContainer(rect)
	rect.SetMinSize(fyne.NewSize(100, 100))
	scroll.Resize(fyne.NewSize(80, 80))
	scroll.Scrolled(&fyne.ScrollEvent{DeltaX: -20, DeltaY: -20})
	scroll.Resize(fyne.NewSize(100, 100))
	assert.Equal(t, 0, scroll.Offset.X)
	assert.Equal(t, 0, scroll.Offset.Y)

}

func TestScrollContainer_ResizeOffset(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	scroll := NewScrollContainer(rect)
	rect.SetMinSize(fyne.NewSize(100, 100))
	scroll.Resize(fyne.NewSize(80, 80))
	scroll.Scrolled(&fyne.ScrollEvent{DeltaX: -20, DeltaY: -20})
	scroll.Resize(fyne.NewSize(90, 90))
	assert.Equal(t, 10, scroll.Offset.X)
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
	r := test.WidgetRenderer(scroll).(*scrollContainerRenderer)
	assert.False(t, r.vertArea.Visible())
	assert.False(t, r.horizArea.Visible())
}

func TestScrollContainer_ShowHiddenScrollBarIfContentGrows(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	scroll := NewScrollContainer(rect)
	r := test.WidgetRenderer(scroll).(*scrollContainerRenderer)
	rect.SetMinSize(fyne.NewSize(100, 100))
	scroll.Resize(fyne.NewSize(200, 200))
	require.False(t, r.horizArea.Visible())
	require.False(t, r.vertArea.Visible())
	rect.SetMinSize(fyne.NewSize(300, 300))
	r.Layout(scroll.Size())
	assert.True(t, r.horizArea.Visible())
	assert.True(t, r.vertArea.Visible())
}

func TestScrollContainer_HideScrollBarIfContentShrinks(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	scroll := NewScrollContainer(rect)
	r := test.WidgetRenderer(scroll).(*scrollContainerRenderer)
	rect.SetMinSize(fyne.NewSize(300, 300))
	scroll.Resize(fyne.NewSize(200, 200))
	require.True(t, r.horizArea.Visible())
	require.True(t, r.vertArea.Visible())
	rect.SetMinSize(fyne.NewSize(200, 200))
	r.Layout(scroll.Size())
	assert.False(t, r.horizArea.Visible())
	assert.False(t, r.vertArea.Visible())
}

func TestScrollContainer_ScrollBarIsSmall(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	scroll := NewScrollContainer(rect)
	rect.SetMinSize(fyne.NewSize(500, 500))
	scroll.Resize(fyne.NewSize(100, 100))
	areaHoriz := test.WidgetRenderer(scroll).(*scrollContainerRenderer).horizArea
	areaVert := test.WidgetRenderer(scroll).(*scrollContainerRenderer).vertArea
	barHoriz := test.WidgetRenderer(areaHoriz).(*scrollBarAreaRenderer).bar
	barVert := test.WidgetRenderer(areaVert).(*scrollBarAreaRenderer).bar
	require.True(t, areaHoriz.Visible())
	require.True(t, areaVert.Visible())
	require.Less(t, theme.ScrollBarSmallSize(), theme.ScrollBarSize())

	assert.Equal(t, theme.ScrollBarSmallSize()*2, areaHoriz.Size().Height)
	assert.Equal(t, theme.ScrollBarSmallSize()*2, areaVert.Size().Width)
	assert.Equal(t, fyne.NewPos(0, 100-theme.ScrollBarSmallSize()*2), areaHoriz.Position())
	assert.Equal(t, fyne.NewPos(100-theme.ScrollBarSmallSize()*2, 0), areaVert.Position())
	assert.Equal(t, theme.ScrollBarSmallSize(), barHoriz.Size().Height)
	assert.Equal(t, theme.ScrollBarSmallSize(), barVert.Size().Width)
	assert.Equal(t, theme.ScrollBarSmallSize(), barHoriz.Position().Y)
	assert.Equal(t, theme.ScrollBarSmallSize(), barVert.Position().X)

}

func TestScrollContainer_ScrollBarGrowsAndShrinksOnMouseInAndMouseOut(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	scroll := NewScrollContainer(rect)
	rect.SetMinSize(fyne.NewSize(500, 500))
	scroll.Resize(fyne.NewSize(100, 100))
	areaHoriz := test.WidgetRenderer(scroll).(*scrollContainerRenderer).horizArea
	areaVert := test.WidgetRenderer(scroll).(*scrollContainerRenderer).vertArea
	barHoriz := test.WidgetRenderer(areaHoriz).(*scrollBarAreaRenderer).bar
	barVert := test.WidgetRenderer(areaVert).(*scrollBarAreaRenderer).bar
	require.True(t, barHoriz.Visible())
	require.Less(t, theme.ScrollBarSmallSize(), theme.ScrollBarSize())
	require.Equal(t, theme.ScrollBarSmallSize()*2, areaHoriz.Size().Height)
	require.Equal(t, fyne.NewPos(0, 100-theme.ScrollBarSmallSize()*2), areaHoriz.Position())
	require.Equal(t, theme.ScrollBarSmallSize(), barHoriz.Size().Height)
	require.Equal(t, theme.ScrollBarSmallSize(), barHoriz.Position().Y)
	barHoriz.MouseIn(&desktop.MouseEvent{})
	assert.Equal(t, theme.ScrollBarSize(), areaHoriz.Size().Height)
	assert.Equal(t, fyne.NewPos(0, 100-theme.ScrollBarSize()), areaHoriz.Position())
	assert.Equal(t, theme.ScrollBarSize(), barHoriz.Size().Height)
	assert.Equal(t, 0, barHoriz.Position().Y)
	barHoriz.MouseOut()
	assert.Equal(t, theme.ScrollBarSmallSize()*2, areaHoriz.Size().Height)
	assert.Equal(t, fyne.NewPos(0, 100-theme.ScrollBarSmallSize()*2), areaHoriz.Position())
	assert.Equal(t, theme.ScrollBarSmallSize(), barHoriz.Size().Height)
	assert.Equal(t, theme.ScrollBarSmallSize(), barHoriz.Position().Y)
	require.True(t, barVert.Visible())
	require.Less(t, theme.ScrollBarSmallSize(), theme.ScrollBarSize())
	require.Equal(t, theme.ScrollBarSmallSize()*2, areaVert.Size().Width)
	require.Equal(t, fyne.NewPos(100-theme.ScrollBarSmallSize()*2, 0), areaVert.Position())
	require.Equal(t, theme.ScrollBarSmallSize(), barVert.Size().Width)
	require.Equal(t, theme.ScrollBarSmallSize(), barVert.Position().X)
	barVert.MouseIn(&desktop.MouseEvent{})
	assert.Equal(t, theme.ScrollBarSize(), areaVert.Size().Width)
	assert.Equal(t, fyne.NewPos(100-theme.ScrollBarSize(), 0), areaVert.Position())
	assert.Equal(t, theme.ScrollBarSize(), barVert.Size().Width)
	assert.Equal(t, 0, barVert.Position().X)
	barVert.MouseOut()
	assert.Equal(t, theme.ScrollBarSmallSize()*2, areaVert.Size().Width)
	assert.Equal(t, fyne.NewPos(100-theme.ScrollBarSmallSize()*2, 0), areaVert.Position())
	assert.Equal(t, theme.ScrollBarSmallSize(), barVert.Size().Width)
	assert.Equal(t, theme.ScrollBarSmallSize(), barVert.Position().X)
}

func TestScrollContainer_ShowShadowOnLeftIfContentIsScrolled(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(500, 100))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	r := test.WidgetRenderer(scroll).(*scrollContainerRenderer)
	assert.False(t, r.leftShadow.Visible())
	assert.Equal(t, fyne.NewPos(0, 0), r.leftShadow.Position())

	scroll.Scrolled(&fyne.ScrollEvent{DeltaX: -1})
	assert.True(t, r.leftShadow.Visible())

	scroll.Scrolled(&fyne.ScrollEvent{DeltaX: 1})
	assert.False(t, r.leftShadow.Visible())
}

func TestScrollContainer_ShowShadowOnRightIfContentCanScroll(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(500, 100))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	r := test.WidgetRenderer(scroll).(*scrollContainerRenderer)
	assert.True(t, r.rightShadow.Visible())
	assert.Equal(t, scroll.size.Width, r.rightShadow.Position().X+r.rightShadow.Size().Width)

	scroll.Scrolled(&fyne.ScrollEvent{DeltaX: -400})
	assert.False(t, r.rightShadow.Visible())

	scroll.Scrolled(&fyne.ScrollEvent{DeltaX: 100})
	assert.True(t, r.rightShadow.Visible())
}

func TestScrollContainer_ShowShadowOnTopIfContentIsScrolled(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(100, 500))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	r := test.WidgetRenderer(scroll).(*scrollContainerRenderer)
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
	r := test.WidgetRenderer(scroll).(*scrollContainerRenderer)
	assert.True(t, r.bottomShadow.Visible())
	assert.Equal(t, scroll.size.Height, r.bottomShadow.Position().Y+r.bottomShadow.Size().Height)

	scroll.Scrolled(&fyne.ScrollEvent{DeltaY: -400})
	assert.False(t, r.bottomShadow.Visible())

	scroll.Scrolled(&fyne.ScrollEvent{DeltaY: 100})
	assert.True(t, r.bottomShadow.Visible())
}

func TestScrollContainer_ScrollHorizontallyWithVerticalMouseScroll(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(1000, 50))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	assert.Equal(t, 0, scroll.Offset.X)
	assert.Equal(t, 0, scroll.Offset.Y)
	scroll.Scrolled(&fyne.ScrollEvent{DeltaX: 0, DeltaY: -10})
	assert.Equal(t, 10, scroll.Offset.X)
	assert.Equal(t, 0, scroll.Offset.Y)

	t.Run("not if scroll event includes horizontal offset", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(1000, 50))
		scroll := NewScrollContainer(rect)
		scroll.Resize(fyne.NewSize(100, 100))
		assert.Equal(t, 0, scroll.Offset.X)
		assert.Equal(t, 0, scroll.Offset.Y)
		scroll.Scrolled(&fyne.ScrollEvent{DeltaX: -20, DeltaY: -40})
		assert.Equal(t, 20, scroll.Offset.X)
		assert.Equal(t, 0, scroll.Offset.Y)
	})

	t.Run("not if content is vertically scrollable", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(1000, 1000))
		scroll := NewScrollContainer(rect)
		scroll.Resize(fyne.NewSize(100, 100))
		assert.Equal(t, 0, scroll.Offset.X)
		assert.Equal(t, 0, scroll.Offset.Y)
		scroll.Scrolled(&fyne.ScrollEvent{DeltaX: 0, DeltaY: -10})
		assert.Equal(t, 0, scroll.Offset.X)
		assert.Equal(t, 10, scroll.Offset.Y)
	})
}

func TestScrollBarRenderer_BarSize(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(100, 100))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	areaHoriz := test.WidgetRenderer(test.WidgetRenderer(scroll).(*scrollContainerRenderer).horizArea).(*scrollBarAreaRenderer)
	areaVert := test.WidgetRenderer(test.WidgetRenderer(scroll).(*scrollContainerRenderer).vertArea).(*scrollBarAreaRenderer)

	assert.Equal(t, 100, areaHoriz.bar.Size().Width)
	assert.Equal(t, 100, areaVert.bar.Size().Height)

	// resize so content is twice our size. Bar should therefore be half again.
	scroll.Resize(fyne.NewSize(50, 50))
	assert.Equal(t, 25, areaHoriz.bar.Size().Width)
	assert.Equal(t, 25, areaVert.bar.Size().Height)
}

func TestScrollContainerRenderer_LimitBarSize(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(100, 100))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(120, 120))
	areaHoriz := test.WidgetRenderer(test.WidgetRenderer(scroll).(*scrollContainerRenderer).horizArea).(*scrollBarAreaRenderer)
	areaVert := test.WidgetRenderer(test.WidgetRenderer(scroll).(*scrollContainerRenderer).vertArea).(*scrollBarAreaRenderer)

	assert.Equal(t, 120, areaHoriz.bar.Size().Width)
	assert.Equal(t, 120, areaVert.bar.Size().Height)
}

func TestScrollContainerRenderer_Direction(t *testing.T) {
	t.Run("Both", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(100, 100))
		scroll := NewScrollContainer(rect)
		r := test.WidgetRenderer(scroll).(*scrollContainerRenderer)
		assert.NotNil(t, r.vertArea)
		assert.NotNil(t, r.topShadow)
		assert.NotNil(t, r.bottomShadow)
		assert.NotNil(t, r.horizArea)
		assert.NotNil(t, r.leftShadow)
		assert.NotNil(t, r.rightShadow)
	})
	t.Run("HorizontalOnly", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(100, 100))
		scroll := NewHScrollContainer(rect)
		r := test.WidgetRenderer(scroll).(*scrollContainerRenderer)
		assert.Nil(t, r.vertArea)
		assert.Nil(t, r.topShadow)
		assert.Nil(t, r.bottomShadow)
		assert.NotNil(t, r.horizArea)
		assert.NotNil(t, r.leftShadow)
		assert.NotNil(t, r.rightShadow)
	})
	t.Run("VerticalOnly", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(100, 100))
		scroll := NewVScrollContainer(rect)
		r := test.WidgetRenderer(scroll).(*scrollContainerRenderer)
		assert.NotNil(t, r.vertArea)
		assert.NotNil(t, r.topShadow)
		assert.NotNil(t, r.bottomShadow)
		assert.Nil(t, r.horizArea)
		assert.Nil(t, r.leftShadow)
		assert.Nil(t, r.rightShadow)
	})
}

func TestScrollContainerRenderer_MinSize_Direction(t *testing.T) {
	t.Run("Both", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(100, 100))
		scroll := NewScrollContainer(rect)
		size := test.WidgetRenderer(scroll).MinSize()
		assert.Equal(t, 32, size.Height)
		assert.Equal(t, 32, size.Width)
	})
	t.Run("HorizontalOnly", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(100, 100))
		scroll := NewHScrollContainer(rect)
		size := test.WidgetRenderer(scroll).MinSize()
		assert.Equal(t, 100, size.Height)
		assert.Equal(t, 32, size.Width)
	})
	t.Run("VerticalOnly", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(100, 100))
		scroll := NewVScrollContainer(rect)
		size := test.WidgetRenderer(scroll).MinSize()
		assert.Equal(t, 32, size.Height)
		assert.Equal(t, 100, size.Width)
	})
}

func TestScrollContainerRenderer_SetMinSize_Direction(t *testing.T) {
	t.Run("Both", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(100, 100))
		scroll := NewScrollContainer(rect)
		scroll.SetMinSize(fyne.NewSize(50, 50))
		size := test.WidgetRenderer(scroll).MinSize()
		assert.Equal(t, 50, size.Height)
		assert.Equal(t, 50, size.Width)
	})
	t.Run("HorizontalOnly", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(100, 100))
		scroll := NewHScrollContainer(rect)
		scroll.SetMinSize(fyne.NewSize(50, 50))
		size := test.WidgetRenderer(scroll).MinSize()
		assert.Equal(t, 100, size.Height)
		assert.Equal(t, 50, size.Width)
	})
	t.Run("VerticalOnly", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(100, 100))
		scroll := NewVScrollContainer(rect)
		scroll.SetMinSize(fyne.NewSize(50, 50))
		size := test.WidgetRenderer(scroll).MinSize()
		assert.Equal(t, 50, size.Height)
		assert.Equal(t, 100, size.Width)
	})
}

func TestScrollBar_Dragged_ClickedInside(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(500, 500))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	scrollBarHoriz := test.WidgetRenderer(test.WidgetRenderer(scroll).(*scrollContainerRenderer).horizArea).(*scrollBarAreaRenderer).bar
	scrollBarVert := test.WidgetRenderer(test.WidgetRenderer(scroll).(*scrollContainerRenderer).vertArea).(*scrollBarAreaRenderer).bar

	// Create drag event with starting position inside scroll rectangle area
	dragEvent := fyne.DragEvent{DraggedX: 20}
	assert.Equal(t, 0, scroll.Offset.X)
	scrollBarHoriz.Dragged(&dragEvent)
	assert.Equal(t, 100, scroll.Offset.X)

	dragEvent = fyne.DragEvent{DraggedY: 20}
	assert.Equal(t, 0, scroll.Offset.Y)
	scrollBarVert.Dragged(&dragEvent)
	assert.Equal(t, 100, scroll.Offset.Y)
}

func TestScrollBar_DraggedBack_ClickedInside(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(500, 500))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	scrollBarHoriz := test.WidgetRenderer(test.WidgetRenderer(scroll).(*scrollContainerRenderer).horizArea).(*scrollBarAreaRenderer).bar
	scrollBarVert := test.WidgetRenderer(test.WidgetRenderer(scroll).(*scrollContainerRenderer).vertArea).(*scrollBarAreaRenderer).bar

	// Drag forward
	dragEvent := fyne.DragEvent{DraggedX: 20}
	scrollBarHoriz.Dragged(&dragEvent)
	dragEvent = fyne.DragEvent{DraggedY: 20}
	scrollBarVert.Dragged(&dragEvent)

	// Drag back
	dragEvent = fyne.DragEvent{DraggedX: -10}
	assert.Equal(t, 100, scroll.Offset.X)
	scrollBarHoriz.Dragged(&dragEvent)
	assert.Equal(t, 50, scroll.Offset.X)

	dragEvent = fyne.DragEvent{DraggedY: -10}
	assert.Equal(t, 100, scroll.Offset.Y)
	scrollBarVert.Dragged(&dragEvent)
	assert.Equal(t, 50, scroll.Offset.Y)
}

func TestScrollBar_Dragged_Limit(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(1000, 1000))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(200, 200))
	scrollBarHoriz := test.WidgetRenderer(test.WidgetRenderer(scroll).(*scrollContainerRenderer).horizArea).(*scrollBarAreaRenderer).bar
	scrollBarVert := test.WidgetRenderer(test.WidgetRenderer(scroll).(*scrollContainerRenderer).vertArea).(*scrollBarAreaRenderer).bar

	// Drag over limit
	dragEvent := fyne.DragEvent{DraggedX: 2000}
	assert.Equal(t, 0, scroll.Offset.X)
	scrollBarHoriz.Dragged(&dragEvent)
	assert.Equal(t, 800, scroll.Offset.X)

	dragEvent = fyne.DragEvent{DraggedY: 2000}
	assert.Equal(t, 0, scroll.Offset.Y)
	scrollBarVert.Dragged(&dragEvent)
	assert.Equal(t, 800, scroll.Offset.Y)

	// Drag again
	dragEvent = fyne.DragEvent{DraggedX: 100}
	// Offset doesn't go over limit
	assert.Equal(t, 800, scroll.Offset.X)
	scrollBarHoriz.Dragged(&dragEvent)
	assert.Equal(t, 800, scroll.Offset.X)

	dragEvent = fyne.DragEvent{DraggedY: 100}
	// Offset doesn't go over limit
	assert.Equal(t, 800, scroll.Offset.Y)
	scrollBarVert.Dragged(&dragEvent)
	assert.Equal(t, 800, scroll.Offset.Y)

	// Drag back (still outside limit)
	dragEvent = fyne.DragEvent{DraggedX: -1000}
	scrollBarHoriz.Dragged(&dragEvent)
	assert.Equal(t, 800, scroll.Offset.X)

	dragEvent = fyne.DragEvent{DraggedY: -1000}
	scrollBarVert.Dragged(&dragEvent)
	assert.Equal(t, 800, scroll.Offset.Y)

	// Drag back (inside limit)
	dragEvent = fyne.DragEvent{DraggedX: -1040}
	scrollBarHoriz.Dragged(&dragEvent)
	assert.Equal(t, 300, scroll.Offset.X)

	dragEvent = fyne.DragEvent{DraggedY: -1040}
	scrollBarVert.Dragged(&dragEvent)
	assert.Equal(t, 300, scroll.Offset.Y)
}

func TestScrollBar_Dragged_BackLimit(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(1000, 1000))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(200, 200))
	scrollBarHoriz := test.WidgetRenderer(test.WidgetRenderer(scroll).(*scrollContainerRenderer).horizArea).(*scrollBarAreaRenderer).bar
	scrollBarVert := test.WidgetRenderer(test.WidgetRenderer(scroll).(*scrollContainerRenderer).vertArea).(*scrollBarAreaRenderer).bar

	// Drag over back limit
	dragEvent := fyne.DragEvent{DraggedX: -1000}
	// Offset doesn't go over limit
	assert.Equal(t, 0, scroll.Offset.X)
	scrollBarHoriz.Dragged(&dragEvent)
	assert.Equal(t, 0, scroll.Offset.X)

	dragEvent = fyne.DragEvent{DraggedY: -1000}
	// Offset doesn't go over limit
	assert.Equal(t, 0, scroll.Offset.Y)
	scrollBarVert.Dragged(&dragEvent)
	assert.Equal(t, 0, scroll.Offset.Y)

	// Drag (still outside limit)
	dragEvent = fyne.DragEvent{DraggedX: 500}
	scrollBarHoriz.Dragged(&dragEvent)
	assert.Equal(t, 0, scroll.Offset.X)

	dragEvent = fyne.DragEvent{DraggedY: 500}
	scrollBarVert.Dragged(&dragEvent)
	assert.Equal(t, 0, scroll.Offset.Y)

	// Drag (inside limit)
	dragEvent = fyne.DragEvent{DraggedX: 520}
	scrollBarHoriz.Dragged(&dragEvent)
	assert.Equal(t, 100, scroll.Offset.X)

	dragEvent = fyne.DragEvent{DraggedY: 520}
	scrollBarVert.Dragged(&dragEvent)
	assert.Equal(t, 100, scroll.Offset.Y)
}

func TestScrollBar_DraggedWithNonZeroStartPosition(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(1000, 1000))
	scroll := NewScrollContainer(rect)
	scroll.Resize(fyne.NewSize(200, 200))
	scrollBarHoriz := test.WidgetRenderer(test.WidgetRenderer(scroll).(*scrollContainerRenderer).horizArea).(*scrollBarAreaRenderer).bar
	scrollBarVert := test.WidgetRenderer(test.WidgetRenderer(scroll).(*scrollContainerRenderer).vertArea).(*scrollBarAreaRenderer).bar

	dragEvent := fyne.DragEvent{DraggedX: 50}
	assert.Equal(t, 0, scroll.Offset.X)
	scrollBarHoriz.Dragged(&dragEvent)
	assert.Equal(t, 250, scroll.Offset.X)

	dragEvent = fyne.DragEvent{DraggedY: 50}
	assert.Equal(t, 0, scroll.Offset.Y)
	scrollBarVert.Dragged(&dragEvent)
	assert.Equal(t, 250, scroll.Offset.Y)

	// Drag again (after releasing mouse button)
	dragEvent = fyne.DragEvent{DraggedX: 20}
	scrollBarHoriz.Dragged(&dragEvent)
	assert.Equal(t, 350, scroll.Offset.X)

	dragEvent = fyne.DragEvent{DraggedY: 20}
	scrollBarVert.Dragged(&dragEvent)
	assert.Equal(t, 350, scroll.Offset.Y)
}
