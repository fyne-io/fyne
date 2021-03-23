package widget

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/theme"
)

func TestNewScroll(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(10, 10))
	scroll := NewScroll(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	barArea := cache.Renderer(scroll).(*scrollContainerRenderer).vertArea
	bar := cache.Renderer(barArea).(*scrollBarAreaRenderer).bar
	assert.Equal(t, float32(0), scroll.Offset.Y)
	assert.Equal(t, theme.ScrollBarSmallSize()*2, barArea.Size().Width)
	assert.Equal(t, theme.ScrollBarSmallSize(), bar.Size().Width)
	assert.Equal(t, theme.ScrollBarSmallSize(), bar.Position().X)
	assert.Equal(t, fyne.NewPos(100-theme.ScrollBarSmallSize()*2, 0), barArea.Position())
}

func TestScrollContainer_MinSize(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(500, 50))
	scroll := NewScroll(rect)
	assert.Equal(t, fyne.NewSize(32, 32), scroll.MinSize())

	scrollMin := fyne.NewSize(100, 100)
	scroll.SetMinSize(scrollMin)
	cache.Renderer(scroll).Layout(scroll.minSize)

	assert.Equal(t, scrollMin, scroll.MinSize())
	assert.Equal(t, fyne.NewSize(500, 100), rect.Size())
	assert.Equal(t, float32(0), scroll.Offset.X)
	assert.Equal(t, float32(0), scroll.Offset.Y)
}

func TestScrollContainer_ScrollToTop(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(500, 50))
	scroll := NewScroll(rect)
	scroll.ScrollToTop()
	Y := scroll.Offset.Y
	assert.Equal(t, float32(0), Y)
}

func TestScrollContainer_ScrollToBottom(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(500, 50))
	scroll := NewScroll(rect)
	scroll.ScrollToBottom()
	ExpectedY := float32(50)
	Y := scroll.Content.Size().Height - scroll.Size().Height
	assert.Equal(t, ExpectedY, Y)
}

func TestScrollContainer_MinSize_Direction(t *testing.T) {
	t.Run("Both", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(100, 100))
		scroll := NewScroll(rect)
		size := scroll.MinSize()
		assert.Equal(t, float32(32), size.Height)
		assert.Equal(t, float32(32), size.Width)
	})
	t.Run("HorizontalOnly", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(100, 100))
		scroll := NewHScroll(rect)
		size := scroll.MinSize()
		assert.Equal(t, float32(100), size.Height)
		assert.Equal(t, float32(32), size.Width)
	})
	t.Run("VerticalOnly", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(100, 100))
		scroll := NewVScroll(rect)
		size := scroll.MinSize()
		assert.Equal(t, float32(32), size.Height)
		assert.Equal(t, float32(100), size.Width)
	})
}

func TestScrollContainer_SetMinSize_Direction(t *testing.T) {
	t.Run("Both", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(100, 100))
		scroll := NewScroll(rect)
		scroll.SetMinSize(fyne.NewSize(50, 50))
		size := scroll.MinSize()
		assert.Equal(t, float32(50), size.Height)
		assert.Equal(t, float32(50), size.Width)
	})
	t.Run("HorizontalOnly", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(100, 100))
		scroll := NewHScroll(rect)
		scroll.SetMinSize(fyne.NewSize(50, 50))
		size := scroll.MinSize()
		assert.Equal(t, float32(100), size.Height)
		assert.Equal(t, float32(50), size.Width)
	})
	t.Run("VerticalOnly", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(100, 100))
		scroll := NewVScroll(rect)
		scroll.SetMinSize(fyne.NewSize(50, 50))
		size := scroll.MinSize()
		assert.Equal(t, float32(50), size.Height)
		assert.Equal(t, float32(100), size.Width)
	})
}

func TestScrollContainer_Resize_Direction(t *testing.T) {
	t.Run("Both", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(100, 100))
		scroll := NewScroll(rect)
		scroll.Resize(scroll.MinSize())
		size := scroll.Size()
		assert.Equal(t, float32(32), size.Height)
		assert.Equal(t, float32(32), size.Width)
	})
	t.Run("HorizontalOnly", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(100, 100))
		scroll := NewHScroll(rect)
		scroll.Resize(scroll.MinSize())
		size := scroll.Size()
		assert.Equal(t, float32(100), size.Height)
		assert.Equal(t, float32(32), size.Width)
	})
	t.Run("VerticalOnly", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(100, 100))
		scroll := NewVScroll(rect)
		scroll.Resize(scroll.MinSize())
		size := scroll.Size()
		assert.Equal(t, float32(32), size.Height)
		assert.Equal(t, float32(100), size.Width)
	})
}

func TestScrollContainer_Refresh(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(1000, 1000))
	scroll := NewScroll(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	assert.Equal(t, fyne.NewSize(1000, 1000), rect.Size())
	scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(-1000, -1000)})
	assert.Equal(t, float32(900), scroll.Offset.X)
	assert.Equal(t, float32(900), scroll.Offset.Y)
	assert.Equal(t, fyne.NewSize(1000, 1000), rect.Size())
	rect.SetMinSize(fyne.NewSize(500, 500))
	scroll.Refresh()
	assert.Equal(t, float32(400), scroll.Offset.X)
	assert.Equal(t, fyne.NewSize(500, 500), rect.Size())

	rect2 := canvas.NewRectangle(color.White)
	scroll.Content = rect2
	scroll.Refresh()
	assert.Equal(t, rect2, cache.Renderer(scroll).Objects()[0])
}

func TestScrollContainer_Scrolled(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(1000, 1000))
	scroll := NewScroll(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	assert.Equal(t, float32(0), scroll.Offset.X)
	assert.Equal(t, float32(0), scroll.Offset.Y)
	scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(-10, -10)})
	assert.Equal(t, float32(10), scroll.Offset.X)
	assert.Equal(t, float32(10), scroll.Offset.Y)

}

func TestScrollContainer_Scrolled_Limit(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(100, 100))
	scroll := NewScroll(rect)
	scroll.Resize(fyne.NewSize(80, 80))
	scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(-25, -25)})
	assert.Equal(t, float32(20), scroll.Offset.X)
}

func TestScrollContainer_Scrolled_Back(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(1000, 1000))
	scroll := NewScroll(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	scroll.Offset.X = 10
	scroll.Offset.Y = 10
	scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(10, 10)})
	assert.Equal(t, float32(0), scroll.Offset.X)
	assert.Equal(t, float32(0), scroll.Offset.Y)
}

func TestScrollContainer_Scrolled_BackLimit(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	scroll := NewScroll(rect)
	rect.SetMinSize(fyne.NewSize(1000, 1000))
	scroll.Resize(fyne.NewSize(100, 100))
	scroll.Offset.X = 10
	scroll.Offset.Y = 10
	scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(20, 20)})
	assert.Equal(t, float32(0), scroll.Offset.X)
	assert.Equal(t, float32(0), scroll.Offset.Y)

}

func TestScrollContainer_Resize(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	scroll := NewScroll(rect)
	rect.SetMinSize(fyne.NewSize(100, 100))
	scroll.Resize(fyne.NewSize(80, 80))
	scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(-20, -20)})
	scroll.Resize(fyne.NewSize(100, 100))
	assert.Equal(t, float32(0), scroll.Offset.X)
	assert.Equal(t, float32(0), scroll.Offset.Y)

}

func TestScrollContainer_ResizeOffset(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	scroll := NewScroll(rect)
	rect.SetMinSize(fyne.NewSize(100, 100))
	scroll.Resize(fyne.NewSize(80, 80))
	scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(-20, -20)})
	scroll.Resize(fyne.NewSize(90, 90))
	assert.Equal(t, float32(10), scroll.Offset.X)
	assert.Equal(t, float32(10), scroll.Offset.Y)
}

func TestScrollContainer_ResizeExpand(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(100, 100))
	scroll := NewScroll(rect)
	scroll.Resize(fyne.NewSize(120, 140))
	assert.Equal(t, float32(120), rect.Size().Width)
	assert.Equal(t, float32(140), rect.Size().Height)
}

func TestScrollContainer_ScrollBarForSmallContentIsHidden(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(100, 100))
	scroll := NewScroll(rect)
	scroll.Resize(fyne.NewSize(100, 200))
	r := cache.Renderer(scroll).(*scrollContainerRenderer)
	assert.False(t, r.vertArea.Visible())
	assert.False(t, r.horizArea.Visible())
}

func TestScrollContainer_ShowHiddenScrollBarIfContentGrows(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	scroll := NewScroll(rect)
	r := cache.Renderer(scroll).(*scrollContainerRenderer)
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
	scroll := NewScroll(rect)
	r := cache.Renderer(scroll).(*scrollContainerRenderer)
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
	scroll := NewScroll(rect)
	rect.SetMinSize(fyne.NewSize(500, 500))
	scroll.Resize(fyne.NewSize(100, 100))
	areaHoriz := cache.Renderer(scroll).(*scrollContainerRenderer).horizArea
	areaVert := cache.Renderer(scroll).(*scrollContainerRenderer).vertArea
	barHoriz := cache.Renderer(areaHoriz).(*scrollBarAreaRenderer).bar
	barVert := cache.Renderer(areaVert).(*scrollBarAreaRenderer).bar
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
	scroll := NewScroll(rect)
	rect.SetMinSize(fyne.NewSize(500, 500))
	scroll.Resize(fyne.NewSize(100, 100))
	areaHoriz := cache.Renderer(scroll).(*scrollContainerRenderer).horizArea
	areaVert := cache.Renderer(scroll).(*scrollContainerRenderer).vertArea
	barHoriz := cache.Renderer(areaHoriz).(*scrollBarAreaRenderer).bar
	barVert := cache.Renderer(areaVert).(*scrollBarAreaRenderer).bar
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
	assert.Equal(t, float32(0), barHoriz.Position().Y)
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
	assert.Equal(t, float32(0), barVert.Position().X)
	barVert.MouseOut()
	assert.Equal(t, theme.ScrollBarSmallSize()*2, areaVert.Size().Width)
	assert.Equal(t, fyne.NewPos(100-theme.ScrollBarSmallSize()*2, 0), areaVert.Position())
	assert.Equal(t, theme.ScrollBarSmallSize(), barVert.Size().Width)
	assert.Equal(t, theme.ScrollBarSmallSize(), barVert.Position().X)
}

func TestScrollContainer_ShowShadowOnLeftIfContentIsScrolled(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(500, 100))
	scroll := NewScroll(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	r := cache.Renderer(scroll).(*scrollContainerRenderer)
	assert.False(t, r.leftShadow.Visible())
	assert.Equal(t, fyne.NewPos(0, 0), r.leftShadow.Position())

	scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.Delta{DX: -1}})
	assert.True(t, r.leftShadow.Visible())

	scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.Delta{DX: 1}})
	assert.False(t, r.leftShadow.Visible())
}

func TestScrollContainer_ShowShadowOnRightIfContentCanScroll(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(500, 100))
	scroll := NewScroll(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	r := cache.Renderer(scroll).(*scrollContainerRenderer)
	assert.True(t, r.rightShadow.Visible())
	assert.Equal(t, scroll.size.Width, r.rightShadow.Position().X+r.rightShadow.Size().Width)

	scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.Delta{DX: -400}})
	assert.False(t, r.rightShadow.Visible())

	scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.Delta{DX: 100}})
	assert.True(t, r.rightShadow.Visible())
}

func TestScrollContainer_ShowShadowOnTopIfContentIsScrolled(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(100, 500))
	scroll := NewScroll(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	r := cache.Renderer(scroll).(*scrollContainerRenderer)
	assert.False(t, r.topShadow.Visible())
	assert.Equal(t, fyne.NewPos(0, 0), r.topShadow.Position())

	scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.Delta{DY: -1}})
	assert.True(t, r.topShadow.Visible())

	scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.Delta{DY: 1}})
	assert.False(t, r.topShadow.Visible())
}

func TestScrollContainer_ShowShadowOnBottomIfContentCanScroll(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(100, 500))
	scroll := NewScroll(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	r := cache.Renderer(scroll).(*scrollContainerRenderer)
	assert.True(t, r.bottomShadow.Visible())
	assert.Equal(t, scroll.size.Height, r.bottomShadow.Position().Y+r.bottomShadow.Size().Height)

	scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.Delta{DY: -400}})
	assert.False(t, r.bottomShadow.Visible())

	scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.Delta{DY: 100}})
	assert.True(t, r.bottomShadow.Visible())
}

func TestScrollContainer_ScrollHorizontallyWithVerticalMouseScroll(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(1000, 50))
	scroll := NewScroll(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	assert.Equal(t, float32(0), scroll.Offset.X)
	assert.Equal(t, float32(0), scroll.Offset.Y)
	scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(0, -10)})
	assert.Equal(t, float32(10), scroll.Offset.X)
	assert.Equal(t, float32(0), scroll.Offset.Y)

	t.Run("not if scroll event includes horizontal offset", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(1000, 50))
		scroll := NewScroll(rect)
		scroll.Resize(fyne.NewSize(100, 100))
		assert.Equal(t, float32(0), scroll.Offset.X)
		assert.Equal(t, float32(0), scroll.Offset.Y)
		scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(-20, -40)})
		assert.Equal(t, float32(20), scroll.Offset.X)
		assert.Equal(t, float32(0), scroll.Offset.Y)
	})

	t.Run("not if content is vertically scrollable", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(1000, 1000))
		scroll := NewScroll(rect)
		scroll.Resize(fyne.NewSize(100, 100))
		assert.Equal(t, float32(0), scroll.Offset.X)
		assert.Equal(t, float32(0), scroll.Offset.Y)
		scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(0, -10)})
		assert.Equal(t, float32(0), scroll.Offset.X)
		assert.Equal(t, float32(10), scroll.Offset.Y)
	})
}

func TestScrollBarRenderer_BarSize(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(100, 100))
	scroll := NewScroll(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	areaHoriz := cache.Renderer(cache.Renderer(scroll).(*scrollContainerRenderer).horizArea).(*scrollBarAreaRenderer)
	areaVert := cache.Renderer(cache.Renderer(scroll).(*scrollContainerRenderer).vertArea).(*scrollBarAreaRenderer)

	assert.Equal(t, float32(100), areaHoriz.bar.Size().Width)
	assert.Equal(t, float32(100), areaVert.bar.Size().Height)

	// resize so content is twice our size. Bar should therefore be half again.
	scroll.Resize(fyne.NewSize(50, 50))
	assert.Equal(t, float32(25), areaHoriz.bar.Size().Width)
	assert.Equal(t, float32(25), areaVert.bar.Size().Height)
}

func TestScrollContainerRenderer_LimitBarSize(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(100, 100))
	scroll := NewScroll(rect)
	scroll.Resize(fyne.NewSize(120, 120))
	areaHoriz := cache.Renderer(cache.Renderer(scroll).(*scrollContainerRenderer).horizArea).(*scrollBarAreaRenderer)
	areaVert := cache.Renderer(cache.Renderer(scroll).(*scrollContainerRenderer).vertArea).(*scrollBarAreaRenderer)

	assert.Equal(t, float32(120), areaHoriz.bar.Size().Width)
	assert.Equal(t, float32(120), areaVert.bar.Size().Height)
}

func TestScrollContainerRenderer_Direction(t *testing.T) {
	t.Run("Both", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(100, 100))
		scroll := NewScroll(rect)
		r := cache.Renderer(scroll).(*scrollContainerRenderer)
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
		scroll := NewHScroll(rect)
		r := cache.Renderer(scroll).(*scrollContainerRenderer)
		assert.NotNil(t, r.vertArea)
		assert.False(t, r.vertArea.Visible())
		assert.NotNil(t, r.topShadow)
		assert.False(t, r.topShadow.Visible())
		assert.NotNil(t, r.bottomShadow)
		assert.False(t, r.bottomShadow.Visible())
		assert.NotNil(t, r.horizArea)
		assert.NotNil(t, r.leftShadow)
		assert.NotNil(t, r.rightShadow)
	})
	t.Run("VerticalOnly", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(100, 100))
		scroll := NewVScroll(rect)
		r := cache.Renderer(scroll).(*scrollContainerRenderer)
		assert.NotNil(t, r.vertArea)
		assert.NotNil(t, r.topShadow)
		assert.NotNil(t, r.bottomShadow)
		assert.NotNil(t, r.horizArea)
		assert.False(t, r.horizArea.Visible())
		assert.NotNil(t, r.leftShadow)
		assert.False(t, r.leftShadow.Visible())
		assert.NotNil(t, r.rightShadow)
		assert.False(t, r.rightShadow.Visible())
	})
}

func TestScrollContainerRenderer_MinSize_Direction(t *testing.T) {
	t.Run("Both", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(100, 100))
		scroll := NewScroll(rect)
		size := cache.Renderer(scroll).MinSize()
		assert.Equal(t, float32(32), size.Height)
		assert.Equal(t, float32(32), size.Width)
	})
	t.Run("HorizontalOnly", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(100, 100))
		scroll := NewHScroll(rect)
		size := cache.Renderer(scroll).MinSize()
		assert.Equal(t, float32(100), size.Height)
		assert.Equal(t, float32(32), size.Width)
	})
	t.Run("VerticalOnly", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(100, 100))
		scroll := NewVScroll(rect)
		size := cache.Renderer(scroll).MinSize()
		assert.Equal(t, float32(32), size.Height)
		assert.Equal(t, float32(100), size.Width)
	})
}

func TestScrollContainerRenderer_SetMinSize_Direction(t *testing.T) {
	t.Run("Both", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(100, 100))
		scroll := NewScroll(rect)
		scroll.SetMinSize(fyne.NewSize(50, 50))
		size := cache.Renderer(scroll).MinSize()
		assert.Equal(t, float32(50), size.Height)
		assert.Equal(t, float32(50), size.Width)
	})
	t.Run("HorizontalOnly", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(100, 100))
		scroll := NewHScroll(rect)
		scroll.SetMinSize(fyne.NewSize(50, 50))
		size := cache.Renderer(scroll).MinSize()
		assert.Equal(t, float32(100), size.Height)
		assert.Equal(t, float32(50), size.Width)
	})
	t.Run("VerticalOnly", func(t *testing.T) {
		rect := canvas.NewRectangle(color.Black)
		rect.SetMinSize(fyne.NewSize(100, 100))
		scroll := NewVScroll(rect)
		scroll.SetMinSize(fyne.NewSize(50, 50))
		size := cache.Renderer(scroll).MinSize()
		assert.Equal(t, float32(50), size.Height)
		assert.Equal(t, float32(100), size.Width)
	})
}

func TestScrollBar_Dragged_ClickedInside(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(500, 500))
	scroll := NewScroll(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	scrollBarHoriz := cache.Renderer(cache.Renderer(scroll).(*scrollContainerRenderer).horizArea).(*scrollBarAreaRenderer).bar
	scrollBarVert := cache.Renderer(cache.Renderer(scroll).(*scrollContainerRenderer).vertArea).(*scrollBarAreaRenderer).bar

	// Create drag event with starting position inside scroll rectangle area
	dragEvent := fyne.DragEvent{Dragged: fyne.Delta{DX: 20}}
	assert.Equal(t, float32(0), scroll.Offset.X)
	scrollBarHoriz.Dragged(&dragEvent)
	assert.Equal(t, float32(100), scroll.Offset.X)

	dragEvent = fyne.DragEvent{Dragged: fyne.Delta{DY: 20}}
	assert.Equal(t, float32(0), scroll.Offset.Y)
	scrollBarVert.Dragged(&dragEvent)
	assert.Equal(t, float32(100), scroll.Offset.Y)
}

func TestScrollBar_DraggedBack_ClickedInside(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(500, 500))
	scroll := NewScroll(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	scrollBarHoriz := cache.Renderer(cache.Renderer(scroll).(*scrollContainerRenderer).horizArea).(*scrollBarAreaRenderer).bar
	scrollBarVert := cache.Renderer(cache.Renderer(scroll).(*scrollContainerRenderer).vertArea).(*scrollBarAreaRenderer).bar

	// Drag forward
	dragEvent := fyne.DragEvent{Dragged: fyne.Delta{DX: 20}}
	scrollBarHoriz.Dragged(&dragEvent)
	dragEvent = fyne.DragEvent{Dragged: fyne.Delta{DY: 20}}
	scrollBarVert.Dragged(&dragEvent)

	// Drag back
	dragEvent = fyne.DragEvent{Dragged: fyne.Delta{DX: -10}}
	assert.Equal(t, float32(100), scroll.Offset.X)
	scrollBarHoriz.Dragged(&dragEvent)
	assert.Equal(t, float32(50), scroll.Offset.X)

	dragEvent = fyne.DragEvent{Dragged: fyne.Delta{DY: -10}}
	assert.Equal(t, float32(100), scroll.Offset.Y)
	scrollBarVert.Dragged(&dragEvent)
	assert.Equal(t, float32(50), scroll.Offset.Y)
}

func TestScrollBar_Dragged_Limit(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(1000, 1000))
	scroll := NewScroll(rect)
	scroll.Resize(fyne.NewSize(200, 200))
	scrollBarHoriz := cache.Renderer(cache.Renderer(scroll).(*scrollContainerRenderer).horizArea).(*scrollBarAreaRenderer).bar
	scrollBarVert := cache.Renderer(cache.Renderer(scroll).(*scrollContainerRenderer).vertArea).(*scrollBarAreaRenderer).bar

	// Drag over limit
	dragEvent := fyne.DragEvent{Dragged: fyne.Delta{DX: 2000}}
	assert.Equal(t, float32(0), scroll.Offset.X)
	scrollBarHoriz.Dragged(&dragEvent)
	assert.Equal(t, float32(800), scroll.Offset.X)

	dragEvent = fyne.DragEvent{Dragged: fyne.Delta{DY: 2000}}
	assert.Equal(t, float32(0), scroll.Offset.Y)
	scrollBarVert.Dragged(&dragEvent)
	assert.Equal(t, float32(800), scroll.Offset.Y)

	// Drag again
	dragEvent = fyne.DragEvent{Dragged: fyne.Delta{DX: 100}}
	// Offset doesn't go over limit
	assert.Equal(t, float32(800), scroll.Offset.X)
	scrollBarHoriz.Dragged(&dragEvent)
	assert.Equal(t, float32(800), scroll.Offset.X)

	dragEvent = fyne.DragEvent{Dragged: fyne.Delta{DY: 100}}
	// Offset doesn't go over limit
	assert.Equal(t, float32(800), scroll.Offset.Y)
	scrollBarVert.Dragged(&dragEvent)
	assert.Equal(t, float32(800), scroll.Offset.Y)

	// Drag back (still outside limit)
	dragEvent = fyne.DragEvent{Dragged: fyne.Delta{DX: -1000}}
	scrollBarHoriz.Dragged(&dragEvent)
	assert.Equal(t, float32(800), scroll.Offset.X)

	dragEvent = fyne.DragEvent{Dragged: fyne.Delta{DY: -1000}}
	scrollBarVert.Dragged(&dragEvent)
	assert.Equal(t, float32(800), scroll.Offset.Y)

	// Drag back (inside limit)
	dragEvent = fyne.DragEvent{Dragged: fyne.Delta{DX: -1040}}
	scrollBarHoriz.Dragged(&dragEvent)
	assert.Equal(t, float32(300), scroll.Offset.X)

	dragEvent = fyne.DragEvent{Dragged: fyne.Delta{DY: -1040}}
	scrollBarVert.Dragged(&dragEvent)
	assert.Equal(t, float32(300), scroll.Offset.Y)
}

func TestScrollBar_Dragged_BackLimit(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(1000, 1000))
	scroll := NewScroll(rect)
	scroll.Resize(fyne.NewSize(200, 200))
	scrollBarHoriz := cache.Renderer(cache.Renderer(scroll).(*scrollContainerRenderer).horizArea).(*scrollBarAreaRenderer).bar
	scrollBarVert := cache.Renderer(cache.Renderer(scroll).(*scrollContainerRenderer).vertArea).(*scrollBarAreaRenderer).bar

	// Drag over back limit
	dragEvent := fyne.DragEvent{Dragged: fyne.Delta{DX: -1000}}
	// Offset doesn't go over limit
	assert.Equal(t, float32(0), scroll.Offset.X)
	scrollBarHoriz.Dragged(&dragEvent)
	assert.Equal(t, float32(0), scroll.Offset.X)

	dragEvent = fyne.DragEvent{Dragged: fyne.Delta{DY: -1000}}
	// Offset doesn't go over limit
	assert.Equal(t, float32(0), scroll.Offset.Y)
	scrollBarVert.Dragged(&dragEvent)
	assert.Equal(t, float32(0), scroll.Offset.Y)

	// Drag (still outside limit)
	dragEvent = fyne.DragEvent{Dragged: fyne.Delta{DX: 500}}
	scrollBarHoriz.Dragged(&dragEvent)
	assert.Equal(t, float32(0), scroll.Offset.X)

	dragEvent = fyne.DragEvent{Dragged: fyne.Delta{DY: 500}}
	scrollBarVert.Dragged(&dragEvent)
	assert.Equal(t, float32(0), scroll.Offset.Y)

	// Drag (inside limit)
	dragEvent = fyne.DragEvent{Dragged: fyne.Delta{DX: 520}}
	scrollBarHoriz.Dragged(&dragEvent)
	assert.Equal(t, float32(100), scroll.Offset.X)

	dragEvent = fyne.DragEvent{Dragged: fyne.Delta{DY: 520}}
	scrollBarVert.Dragged(&dragEvent)
	assert.Equal(t, float32(100), scroll.Offset.Y)
}

func TestScrollBar_DraggedWithNonZeroStartPosition(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(1000, 1000))
	scroll := NewScroll(rect)
	scroll.Resize(fyne.NewSize(200, 200))
	scrollBarHoriz := cache.Renderer(cache.Renderer(scroll).(*scrollContainerRenderer).horizArea).(*scrollBarAreaRenderer).bar
	scrollBarVert := cache.Renderer(cache.Renderer(scroll).(*scrollContainerRenderer).vertArea).(*scrollBarAreaRenderer).bar

	dragEvent := fyne.DragEvent{Dragged: fyne.Delta{DX: 50}}
	assert.Equal(t, float32(0), scroll.Offset.X)
	scrollBarHoriz.Dragged(&dragEvent)
	assert.Equal(t, float32(250), scroll.Offset.X)

	dragEvent = fyne.DragEvent{Dragged: fyne.Delta{DY: 50}}
	assert.Equal(t, float32(0), scroll.Offset.Y)
	scrollBarVert.Dragged(&dragEvent)
	assert.Equal(t, float32(250), scroll.Offset.Y)

	// Drag again (after releasing mouse button)
	dragEvent = fyne.DragEvent{Dragged: fyne.Delta{DX: 20}}
	scrollBarHoriz.Dragged(&dragEvent)
	assert.Equal(t, float32(350), scroll.Offset.X)

	dragEvent = fyne.DragEvent{Dragged: fyne.Delta{DY: 20}}
	scrollBarVert.Dragged(&dragEvent)
	assert.Equal(t, float32(350), scroll.Offset.Y)
}
