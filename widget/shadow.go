package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

// elevationLevel is the level of elevation of the shadow casting object.
type elevationLevel int

// elevationLevel constants inspired by:
// https://storage.googleapis.com/spec-host/mio-staging%2Fmio-design%2F1584058305895%2Fassets%2F0B6xUSjjSulxceF9udnA4Sk5tdU0%2Fbaselineelevation-chart.png
const (
	baseLevel             elevationLevel = 0
	buttonLevel           elevationLevel = 2
	popUpLevel            elevationLevel = 8
	submergedContentLevel elevationLevel = 8
)

type shadowType int

const (
	shadowAround shadowType = iota
	shadowLeft
	shadowRight
	shadowBottom
	shadowTop
)

func newShadow(typ shadowType, level elevationLevel) *shadow {
	s := &shadow{typ: typ, level: level}
	s.ExtendBaseWidget(s)
	return s
}

var _ fyne.Widget = (*shadow)(nil)

type shadow struct {
	BaseWidget
	typ   shadowType
	level elevationLevel
}

func (s *shadow) CreateRenderer() fyne.WidgetRenderer {
	r := &shadowRenderer{s: s}
	r.createShadows()
	return r
}

type shadowRenderer struct {
	baseRenderer
	b, l, r, t     *canvas.LinearGradient
	bl, br, tl, tr *canvas.RadialGradient
	minSize        fyne.Size
	s              *shadow
}

func (r *shadowRenderer) createShadows() {
	switch r.s.typ {
	case shadowLeft:
		r.l = canvas.NewHorizontalGradient(color.Transparent, theme.ShadowColor())
		r.setObjects([]fyne.CanvasObject{r.l})
	case shadowRight:
		r.r = canvas.NewHorizontalGradient(theme.ShadowColor(), color.Transparent)
		r.setObjects([]fyne.CanvasObject{r.r})
	case shadowBottom:
		r.b = canvas.NewVerticalGradient(theme.ShadowColor(), color.Transparent)
		r.setObjects([]fyne.CanvasObject{r.b})
	case shadowTop:
		r.t = canvas.NewVerticalGradient(color.Transparent, theme.ShadowColor())
		r.setObjects([]fyne.CanvasObject{r.t})
	case shadowAround:
		r.tl = canvas.NewRadialGradient(theme.ShadowColor(), color.Transparent)
		r.tl.CenterOffsetX = 0.5
		r.tl.CenterOffsetY = 0.5
		r.t = canvas.NewVerticalGradient(color.Transparent, theme.ShadowColor())
		r.tr = canvas.NewRadialGradient(theme.ShadowColor(), color.Transparent)
		r.tr.CenterOffsetX = -0.5
		r.tr.CenterOffsetY = 0.5
		r.r = canvas.NewHorizontalGradient(theme.ShadowColor(), color.Transparent)
		r.br = canvas.NewRadialGradient(theme.ShadowColor(), color.Transparent)
		r.br.CenterOffsetX = -0.5
		r.br.CenterOffsetY = -0.5
		r.b = canvas.NewVerticalGradient(theme.ShadowColor(), color.Transparent)
		r.bl = canvas.NewRadialGradient(theme.ShadowColor(), color.Transparent)
		r.bl.CenterOffsetX = 0.5
		r.bl.CenterOffsetY = -0.5
		r.l = canvas.NewHorizontalGradient(color.Transparent, theme.ShadowColor())
		r.setObjects([]fyne.CanvasObject{r.tl, r.t, r.tr, r.r, r.br, r.b, r.bl, r.l})
	}
}

func updateShadowStart(g *canvas.LinearGradient) {
	if g == nil {
		return
	}

	g.StartColor = theme.ShadowColor()
	g.Refresh()
}

func updateShadowEnd(g *canvas.LinearGradient) {
	if g == nil {
		return
	}

	g.EndColor = theme.ShadowColor()
	g.Refresh()
}

func updateShadowRadial(g *canvas.RadialGradient) {
	if g == nil {
		return
	}

	g.StartColor = theme.ShadowColor()
	g.Refresh()
}

func (r *shadowRenderer) refreshShadows() {
	updateShadowEnd(r.l)
	updateShadowStart(r.r)
	updateShadowStart(r.b)
	updateShadowEnd(r.t)

	updateShadowRadial(r.tl)
	updateShadowRadial(r.tr)
	updateShadowRadial(r.bl)
	updateShadowRadial(r.br)
}

func (r *shadowRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *shadowRenderer) Layout(size fyne.Size) {
	depth := int(r.s.level)
	if r.tl != nil {
		r.tl.Resize(fyne.NewSize(depth, depth))
		r.tl.Move(fyne.NewPos(-depth, -depth))
	}
	if r.t != nil {
		r.t.Resize(fyne.NewSize(size.Width, depth))
		r.t.Move(fyne.NewPos(0, -depth))
	}
	if r.tr != nil {
		r.tr.Resize(fyne.NewSize(depth, depth))
		r.tr.Move(fyne.NewPos(size.Width, -depth))
	}
	if r.r != nil {
		r.r.Resize(fyne.NewSize(depth, size.Height))
		r.r.Move(fyne.NewPos(size.Width, 0))
	}
	if r.br != nil {
		r.br.Resize(fyne.NewSize(depth, depth))
		r.br.Move(fyne.NewPos(size.Width, size.Height))
	}
	if r.b != nil {
		r.b.Resize(fyne.NewSize(size.Width, depth))
		r.b.Move(fyne.NewPos(0, size.Height))
	}
	if r.bl != nil {
		r.bl.Resize(fyne.NewSize(depth, depth))
		r.bl.Move(fyne.NewPos(-depth, size.Height))
	}
	if r.l != nil {
		r.l.Resize(fyne.NewSize(depth, size.Height))
		r.l.Move(fyne.NewPos(-depth, 0))
	}
}

func (r *shadowRenderer) MinSize() fyne.Size {
	return r.minSize
}

func (r *shadowRenderer) Refresh() {
	r.refreshShadows()
	r.Layout(r.s.Size())
}
