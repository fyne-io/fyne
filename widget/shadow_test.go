package widget

import (
	"image/color"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/internal/cache"
	"fyne.io/fyne/theme"
	"github.com/stretchr/testify/assert"
)

var shadowWidth = 5

func TestShadow_TopShadow(t *testing.T) {
	s := newShadow(shadowTop, shadowWidth)
	r := Renderer(s).(*shadowRenderer)
	r.Layout(fyne.NewSize(100, 100))

	assert.Equal(t, []fyne.CanvasObject{r.t}, r.Objects())
	assert.Equal(t, fyne.NewSize(100, shadowWidth), r.t.Size())
	assert.Equal(t, fyne.NewPos(0, -shadowWidth), r.t.Position())
	assert.Equal(t, 0.0, r.t.Angle)
	assert.Equal(t, color.Transparent, r.t.StartColor)
	assert.Equal(t, theme.ShadowColor(), r.t.EndColor)
}

func TestShadow_BottomShadow(t *testing.T) {
	s := newShadow(shadowBottom, shadowWidth)
	r := Renderer(s).(*shadowRenderer)
	r.Layout(fyne.NewSize(100, 100))

	assert.Equal(t, []fyne.CanvasObject{r.b}, r.Objects())
	assert.Equal(t, fyne.NewSize(100, shadowWidth), r.b.Size())
	assert.Equal(t, fyne.NewPos(0, 100), r.b.Position())
	assert.Equal(t, 0.0, r.b.Angle)
	assert.Equal(t, theme.ShadowColor(), r.b.StartColor)
	assert.Equal(t, color.Transparent, r.b.EndColor)
}

func TestShadow_AroundShadow(t *testing.T) {
	s := newShadow(shadowAround, shadowWidth)
	r := Renderer(s).(*shadowRenderer)
	r.Layout(fyne.NewSize(100, 100))

	assert.Equal(t, []fyne.CanvasObject{r.tl, r.t, r.tr, r.r, r.br, r.b, r.bl, r.l}, r.Objects())

	cornerSize := fyne.NewSize(shadowWidth, shadowWidth)
	horizontalSize := fyne.NewSize(100, shadowWidth)
	verticalSize := fyne.NewSize(shadowWidth, 100)

	assert.Equal(t, cornerSize, r.tl.Size())
	assert.Equal(t, fyne.NewPos(-shadowWidth, -shadowWidth), r.tl.Position())
	assert.Equal(t, 0.5, r.tl.CenterOffsetX)
	assert.Equal(t, 0.5, r.tl.CenterOffsetY)
	assert.Equal(t, theme.ShadowColor(), r.tl.StartColor)
	assert.Equal(t, color.Transparent, r.tl.EndColor)

	assert.Equal(t, horizontalSize, r.t.Size())
	assert.Equal(t, fyne.NewPos(0, -shadowWidth), r.t.Position())
	assert.Equal(t, 0.0, r.t.Angle)
	assert.Equal(t, color.Transparent, r.t.StartColor)
	assert.Equal(t, theme.ShadowColor(), r.t.EndColor)

	assert.Equal(t, cornerSize, r.tr.Size())
	assert.Equal(t, fyne.NewPos(100, -shadowWidth), r.tr.Position())
	assert.Equal(t, -0.5, r.tr.CenterOffsetX)
	assert.Equal(t, 0.5, r.tr.CenterOffsetY)
	assert.Equal(t, theme.ShadowColor(), r.tr.StartColor)
	assert.Equal(t, color.Transparent, r.tr.EndColor)

	assert.Equal(t, verticalSize, r.r.Size())
	assert.Equal(t, fyne.NewPos(100, 0), r.r.Position())
	assert.Equal(t, 270.0, r.r.Angle)
	assert.Equal(t, theme.ShadowColor(), r.r.StartColor)
	assert.Equal(t, color.Transparent, r.r.EndColor)

	assert.Equal(t, cornerSize, r.br.Size())
	assert.Equal(t, fyne.NewPos(100, 100), r.br.Position())
	assert.Equal(t, -0.5, r.br.CenterOffsetX)
	assert.Equal(t, -0.5, r.br.CenterOffsetY)
	assert.Equal(t, theme.ShadowColor(), r.br.StartColor)
	assert.Equal(t, color.Transparent, r.br.EndColor)

	assert.Equal(t, fyne.NewSize(100, shadowWidth), r.b.Size())
	assert.Equal(t, fyne.NewPos(0, 100), r.b.Position())
	assert.Equal(t, 0.0, r.b.Angle)
	assert.Equal(t, theme.ShadowColor(), r.b.StartColor)
	assert.Equal(t, color.Transparent, r.b.EndColor)

	assert.Equal(t, cornerSize, r.bl.Size())
	assert.Equal(t, fyne.NewPos(-shadowWidth, 100), r.bl.Position())
	assert.Equal(t, 0.5, r.bl.CenterOffsetX)
	assert.Equal(t, -0.5, r.bl.CenterOffsetY)
	assert.Equal(t, theme.ShadowColor(), r.bl.StartColor)
	assert.Equal(t, color.Transparent, r.bl.EndColor)

	assert.Equal(t, verticalSize, r.l.Size())
	assert.Equal(t, fyne.NewPos(-shadowWidth, 0), r.l.Position())
	assert.Equal(t, 270.0, r.l.Angle)
	assert.Equal(t, color.Transparent, r.l.StartColor)
	assert.Equal(t, theme.ShadowColor(), r.l.EndColor)
}

func TestShadow_ApplyTheme(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	s := newShadow(shadowAround, shadowWidth)
	r := Renderer(s).(*shadowRenderer)
	assert.Equal(t, theme.ShadowColor(), r.b.StartColor)

	fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
	r.Refresh()
	assert.Equal(t, theme.ShadowColor(), r.b.StartColor)
}

func TestShadow_BackgroundColor(t *testing.T) {
	assert.Equal(t, color.Transparent, Renderer(newShadow(shadowAround, theme.Padding())).BackgroundColor())
}

func TestShadow_MinSize(t *testing.T) {
	assert.Equal(t, fyne.NewSize(0, 0), newShadow(shadowAround, theme.Padding()).MinSize())
}

func TestShadow_Theme(t *testing.T) {
	shadow := newShadow(shadowAround, theme.Padding())
	light := theme.LightTheme()
	fyne.CurrentApp().Settings().SetTheme(light)
	shadow.Refresh()
	assert.Equal(t, light.ShadowColor(), cache.Renderer(shadow).(*shadowRenderer).l.EndColor)
	assert.Equal(t, light.ShadowColor(), cache.Renderer(shadow).(*shadowRenderer).r.StartColor)
	assert.Equal(t, light.ShadowColor(), cache.Renderer(shadow).(*shadowRenderer).tr.StartColor)

	dark := theme.DarkTheme()
	fyne.CurrentApp().Settings().SetTheme(dark)
	shadow.Refresh()
	assert.Equal(t, dark.ShadowColor(), cache.Renderer(shadow).(*shadowRenderer).r.StartColor)
	assert.Equal(t, dark.ShadowColor(), cache.Renderer(shadow).(*shadowRenderer).r.StartColor)
	assert.Equal(t, dark.ShadowColor(), cache.Renderer(shadow).(*shadowRenderer).tr.StartColor)
}
