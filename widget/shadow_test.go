package widget

import (
	"image/color"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"github.com/stretchr/testify/assert"
)

func TestShadow_TopShadow(t *testing.T) {
	shadowWidth := 2 * theme.Padding()
	s := newShadow(shadowTop, shadowWidth)
	r := Renderer(s).(*shadowRenderer)
	r.Layout(fyne.NewSize(100, 100))

	assert.Equal(t, []fyne.CanvasObject{r.t}, r.Objects())
	assert.Equal(t, fyne.NewSize(100, shadowWidth), r.t.Size())
	assert.Equal(t, fyne.NewPos(0, -shadowWidth), r.t.Position())
	assert.Equal(t, canvas.GradientDirectionVertical, r.t.Direction)
	assert.Equal(t, color.Transparent, r.t.StartColor)
	assert.Equal(t, theme.ShadowColor(), r.t.EndColor)
}

func TestShadow_BottomShadow(t *testing.T) {
	shadowWidth := 2 * theme.Padding()
	s := newShadow(shadowBottom, shadowWidth)
	r := Renderer(s).(*shadowRenderer)
	r.Layout(fyne.NewSize(100, 100))

	assert.Equal(t, []fyne.CanvasObject{r.b}, r.Objects())
	assert.Equal(t, fyne.NewSize(100, shadowWidth), r.b.Size())
	assert.Equal(t, fyne.NewPos(0, 100), r.b.Position())
	assert.Equal(t, canvas.GradientDirectionVertical, r.b.Direction)
	assert.Equal(t, theme.ShadowColor(), r.b.StartColor)
	assert.Equal(t, color.Transparent, r.b.EndColor)
}

func TestShadow_AroundShadow(t *testing.T) {
	shadowWidth := 2 * theme.Padding()
	s := newShadow(shadowAround, shadowWidth)
	r := Renderer(s).(*shadowRenderer)
	r.Layout(fyne.NewSize(100, 100))

	assert.Equal(t, []fyne.CanvasObject{r.tl, r.t, r.tr, r.r, r.br, r.b, r.bl, r.l}, r.Objects())

	cornerSize := fyne.NewSize(shadowWidth, shadowWidth)
	horizontalSize := fyne.NewSize(100, shadowWidth)
	verticalSize := fyne.NewSize(shadowWidth, 100)

	assert.Equal(t, cornerSize, r.tl.Size())
	assert.Equal(t, fyne.NewPos(-shadowWidth, -shadowWidth), r.tl.Position())
	assert.Equal(t, canvas.GradientDirectionCircular, r.tl.Direction)
	assert.Equal(t, fyne.NewPos(shadowWidth/2, shadowWidth/2), r.tl.CenterOffset)
	assert.Equal(t, theme.ShadowColor(), r.tl.StartColor)
	assert.Equal(t, color.Transparent, r.tl.EndColor)

	assert.Equal(t, horizontalSize, r.t.Size())
	assert.Equal(t, fyne.NewPos(0, -shadowWidth), r.t.Position())
	assert.Equal(t, canvas.GradientDirectionVertical, r.t.Direction)
	assert.Equal(t, color.Transparent, r.t.StartColor)
	assert.Equal(t, theme.ShadowColor(), r.t.EndColor)

	assert.Equal(t, cornerSize, r.tr.Size())
	assert.Equal(t, fyne.NewPos(100, -shadowWidth), r.tr.Position())
	assert.Equal(t, canvas.GradientDirectionCircular, r.tr.Direction)
	assert.Equal(t, fyne.NewPos(-shadowWidth/2, shadowWidth/2), r.tr.CenterOffset)
	assert.Equal(t, theme.ShadowColor(), r.tr.StartColor)
	assert.Equal(t, color.Transparent, r.tr.EndColor)

	assert.Equal(t, verticalSize, r.r.Size())
	assert.Equal(t, fyne.NewPos(100, 0), r.r.Position())
	assert.Equal(t, canvas.GradientDirectionHorizontal, r.r.Direction)
	assert.Equal(t, theme.ShadowColor(), r.r.StartColor)
	assert.Equal(t, color.Transparent, r.r.EndColor)

	assert.Equal(t, cornerSize, r.br.Size())
	assert.Equal(t, fyne.NewPos(100, 100), r.br.Position())
	assert.Equal(t, canvas.GradientDirectionCircular, r.br.Direction)
	assert.Equal(t, fyne.NewPos(-shadowWidth/2, -shadowWidth/2), r.br.CenterOffset)
	assert.Equal(t, theme.ShadowColor(), r.br.StartColor)
	assert.Equal(t, color.Transparent, r.br.EndColor)

	assert.Equal(t, fyne.NewSize(100, shadowWidth), r.b.Size())
	assert.Equal(t, fyne.NewPos(0, 100), r.b.Position())
	assert.Equal(t, canvas.GradientDirectionVertical, r.b.Direction)
	assert.Equal(t, theme.ShadowColor(), r.b.StartColor)
	assert.Equal(t, color.Transparent, r.b.EndColor)

	assert.Equal(t, cornerSize, r.bl.Size())
	assert.Equal(t, fyne.NewPos(-shadowWidth, 100), r.bl.Position())
	assert.Equal(t, canvas.GradientDirectionCircular, r.bl.Direction)
	assert.Equal(t, fyne.NewPos(shadowWidth/2, -shadowWidth/2), r.bl.CenterOffset)
	assert.Equal(t, theme.ShadowColor(), r.bl.StartColor)
	assert.Equal(t, color.Transparent, r.bl.EndColor)

	assert.Equal(t, verticalSize, r.l.Size())
	assert.Equal(t, fyne.NewPos(-shadowWidth, 0), r.l.Position())
	assert.Equal(t, canvas.GradientDirectionHorizontal, r.l.Direction)
	assert.Equal(t, color.Transparent, r.l.StartColor)
	assert.Equal(t, theme.ShadowColor(), r.l.EndColor)
}

func TestShadow_ApplyTheme(t *testing.T) {
	shadowWidth := 2 * theme.Padding()
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	s := newShadow(shadowAround, shadowWidth)
	r := Renderer(s).(*shadowRenderer)
	assert.Equal(t, theme.ShadowColor(), r.b.StartColor)

	fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
	r.ApplyTheme()
	assert.Equal(t, theme.ShadowColor(), r.b.StartColor)
}

func TestShadow_BackgroundColor(t *testing.T) {
	assert.Equal(t, color.Transparent, Renderer(newShadow(shadowAround, theme.Padding())).BackgroundColor())
}

func TestShadow_MinSize(t *testing.T) {
	assert.Equal(t, fyne.NewSize(0, 0), newShadow(shadowAround, theme.Padding()).MinSize())
}
