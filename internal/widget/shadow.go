package widget

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

var _ fyne.Widget = (*Shadow)(nil)

// Shadow is a widget that renders a shadow.
type Shadow struct {
	Base
	level ElevationLevel
	typ   ShadowType
}

// ElevationLevel is the level of elevation of the shadow casting object.
type ElevationLevel int

// ElevationLevel constants inspired by:
// https://storage.googleapis.com/spec-host/mio-staging%2Fmio-design%2F1584058305895%2Fassets%2F0B6xUSjjSulxceF9udnA4Sk5tdU0%2Fbaselineelevation-chart.png
const (
	BaseLevel             ElevationLevel = 0
	CardLevel             ElevationLevel = 1
	ButtonLevel           ElevationLevel = 2
	MenuLevel             ElevationLevel = 4
	PopUpLevel            ElevationLevel = 8
	SubmergedContentLevel ElevationLevel = 8
	DialogLevel           ElevationLevel = 24
)

// ShadowType specifies the type of the shadow.
type ShadowType int

// ShadowType constants
const (
	ShadowAround ShadowType = iota
	ShadowLeft
	ShadowRight
	ShadowBottom
	ShadowTop
)

// NewShadow create a new Shadow.
func NewShadow(typ ShadowType, level ElevationLevel) *Shadow {
	s := &Shadow{typ: typ, level: level}
	s.ExtendBaseWidget(s)
	return s
}

// CreateRenderer returns a new renderer for the shadow.
func (s *Shadow) CreateRenderer() fyne.WidgetRenderer {
	r := &shadowRenderer{s: s}
	r.createShadows()
	return r
}

type shadowRenderer struct {
	BaseRenderer
	b, l, r, t     *canvas.LinearGradient
	bl, br, tl, tr *canvas.RadialGradient
	minSize        fyne.Size
	s              *Shadow
}

func (r *shadowRenderer) Layout(size fyne.Size) {
	depth := float32(r.s.level)
	sideOff, topOff := float32(0.0), float32(0.0)
	if r.s.typ == ShadowAround {
		sideOff = depth * 0.2
		topOff = sideOff * 2
	}

	if r.tl != nil {
		r.tl.Resize(fyne.NewSize(depth, depth))
		r.tl.Move(fyne.NewPos(-depth+sideOff, -depth+topOff))
	}
	if r.t != nil {
		r.t.Resize(fyne.NewSize(size.Width-sideOff*2, depth))
		r.t.Move(fyne.NewPos(sideOff, -depth+topOff))
	}
	if r.tr != nil {
		r.tr.Resize(fyne.NewSize(depth, depth))
		r.tr.Move(fyne.NewPos(size.Width-sideOff, -depth+topOff))
	}
	if r.r != nil {
		r.r.Resize(fyne.NewSize(depth, size.Height-topOff))
		r.r.Move(fyne.NewPos(size.Width-sideOff, topOff))
	}
	if r.br != nil {
		r.br.Resize(fyne.NewSize(depth, depth))
		r.br.Move(fyne.NewPos(size.Width-sideOff, size.Height))
	}
	if r.b != nil {
		r.b.Resize(fyne.NewSize(size.Width-sideOff*2, depth))
		r.b.Move(fyne.NewPos(sideOff, size.Height))
	}
	if r.bl != nil {
		r.bl.Resize(fyne.NewSize(depth, depth))
		r.bl.Move(fyne.NewPos(-depth+sideOff, size.Height))
	}
	if r.l != nil {
		r.l.Resize(fyne.NewSize(depth, size.Height-topOff))
		r.l.Move(fyne.NewPos(-depth+sideOff, topOff))
	}
}

func (r *shadowRenderer) MinSize() fyne.Size {
	return r.minSize
}

func (r *shadowRenderer) Refresh() {
	r.refreshShadows()
	r.Layout(r.s.Size())
	canvas.Refresh(r.s)
}

func (r *shadowRenderer) createShadows() {
	th := theme.CurrentForWidget(r.s)
	v := fyne.CurrentApp().Settings().ThemeVariant()
	fg := th.Color(theme.ColorNameShadow, v)

	switch r.s.typ {
	case ShadowLeft:
		r.l = canvas.NewHorizontalGradient(color.Transparent, fg)
		r.SetObjects([]fyne.CanvasObject{r.l})
	case ShadowRight:
		r.r = canvas.NewHorizontalGradient(fg, color.Transparent)
		r.SetObjects([]fyne.CanvasObject{r.r})
	case ShadowBottom:
		r.b = canvas.NewVerticalGradient(fg, color.Transparent)
		r.SetObjects([]fyne.CanvasObject{r.b})
	case ShadowTop:
		r.t = canvas.NewVerticalGradient(color.Transparent, fg)
		r.SetObjects([]fyne.CanvasObject{r.t})
	case ShadowAround:
		r.tl = canvas.NewRadialGradient(fg, color.Transparent)
		r.tl.CenterOffsetX = 0.5
		r.tl.CenterOffsetY = 0.5
		r.t = canvas.NewVerticalGradient(color.Transparent, fg)
		r.tr = canvas.NewRadialGradient(fg, color.Transparent)
		r.tr.CenterOffsetX = -0.5
		r.tr.CenterOffsetY = 0.5
		r.r = canvas.NewHorizontalGradient(fg, color.Transparent)
		r.br = canvas.NewRadialGradient(fg, color.Transparent)
		r.br.CenterOffsetX = -0.5
		r.br.CenterOffsetY = -0.5
		r.b = canvas.NewVerticalGradient(fg, color.Transparent)
		r.bl = canvas.NewRadialGradient(fg, color.Transparent)
		r.bl.CenterOffsetX = 0.5
		r.bl.CenterOffsetY = -0.5
		r.l = canvas.NewHorizontalGradient(color.Transparent, fg)
		r.SetObjects([]fyne.CanvasObject{r.tl, r.t, r.tr, r.r, r.br, r.b, r.bl, r.l})
	}
}

func (r *shadowRenderer) refreshShadows() {
	th := theme.CurrentForWidget(r.s)
	v := fyne.CurrentApp().Settings().ThemeVariant()
	fg := th.Color(theme.ColorNameShadow, v)

	updateShadowEnd(r.l, fg)
	updateShadowStart(r.r, fg)
	updateShadowStart(r.b, fg)
	updateShadowEnd(r.t, fg)

	updateShadowRadial(r.tl, fg)
	updateShadowRadial(r.tr, fg)
	updateShadowRadial(r.bl, fg)
	updateShadowRadial(r.br, fg)
}

func updateShadowEnd(g *canvas.LinearGradient, fg color.Color) {
	if g == nil {
		return
	}

	g.EndColor = fg
	g.Refresh()
}

func updateShadowRadial(g *canvas.RadialGradient, fg color.Color) {
	if g == nil {
		return
	}

	g.StartColor = fg
	g.Refresh()
}

func updateShadowStart(g *canvas.LinearGradient, fg color.Color) {
	if g == nil {
		return
	}

	g.StartColor = fg
	g.Refresh()
}
