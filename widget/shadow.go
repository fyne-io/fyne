package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

type shadowType int

const (
	shadowAround shadowType = iota
	shadowLeft
	shadowRight
	shadowBottom
	shadowTop
)

func newShadow(typ shadowType, depth int) *shadow {
	s := &shadow{typ: typ, depth: depth}
	s.ExtendBaseWidget(s)
	return s
}

var _ fyne.Widget = (*shadow)(nil)

type shadow struct {
	BaseWidget
	typ   shadowType
	depth int
}

func (s *shadow) CreateRenderer() fyne.WidgetRenderer {
	r := &shadowRenderer{s: s}
	r.createShadows()
	return r
}

type shadowRenderer struct {
	b, l, r, t     *canvas.LinearGradient
	bl, br, tl, tr *canvas.RadialGradient
	minSize        fyne.Size
	objects        []fyne.CanvasObject
	s              *shadow
}

func (r *shadowRenderer) createShadows() {
	switch r.s.typ {
	case shadowLeft:
		r.l = canvas.NewHorizontalGradient(color.Transparent, theme.ShadowColor())
		r.objects = []fyne.CanvasObject{r.l}
	case shadowRight:
		r.r = canvas.NewHorizontalGradient(theme.ShadowColor(), color.Transparent)
		r.objects = []fyne.CanvasObject{r.r}
	case shadowBottom:
		r.b = canvas.NewVerticalGradient(theme.ShadowColor(), color.Transparent)
		r.objects = []fyne.CanvasObject{r.b}
	case shadowTop:
		r.t = canvas.NewVerticalGradient(color.Transparent, theme.ShadowColor())
		r.objects = []fyne.CanvasObject{r.t}
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
		r.objects = []fyne.CanvasObject{r.tl, r.t, r.tr, r.r, r.br, r.b, r.bl, r.l}
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

func (r *shadowRenderer) Destroy() {
}

func (r *shadowRenderer) Layout(size fyne.Size) {
	if r.tl != nil {
		r.tl.Resize(fyne.NewSize(r.s.depth, r.s.depth))
		r.tl.Move(fyne.NewPos(-r.s.depth, -r.s.depth))
	}
	if r.t != nil {
		r.t.Resize(fyne.NewSize(size.Width, r.s.depth))
		r.t.Move(fyne.NewPos(0, -r.s.depth))
	}
	if r.tr != nil {
		r.tr.Resize(fyne.NewSize(r.s.depth, r.s.depth))
		r.tr.Move(fyne.NewPos(size.Width, -r.s.depth))
	}
	if r.r != nil {
		r.r.Resize(fyne.NewSize(r.s.depth, size.Height))
		r.r.Move(fyne.NewPos(size.Width, 0))
	}
	if r.br != nil {
		r.br.Resize(fyne.NewSize(r.s.depth, r.s.depth))
		r.br.Move(fyne.NewPos(size.Width, size.Height))
	}
	if r.b != nil {
		r.b.Resize(fyne.NewSize(size.Width, r.s.depth))
		r.b.Move(fyne.NewPos(0, size.Height))
	}
	if r.bl != nil {
		r.bl.Resize(fyne.NewSize(r.s.depth, r.s.depth))
		r.bl.Move(fyne.NewPos(-r.s.depth, size.Height))
	}
	if r.l != nil {
		r.l.Resize(fyne.NewSize(r.s.depth, size.Height))
		r.l.Move(fyne.NewPos(-r.s.depth, 0))
	}
}

func (r *shadowRenderer) MinSize() fyne.Size {
	return r.minSize
}

func (r *shadowRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *shadowRenderer) Refresh() {
	r.refreshShadows()
}
