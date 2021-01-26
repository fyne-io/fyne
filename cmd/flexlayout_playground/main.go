// Package main loads a very basic Hello World graphical application.
package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello")

	w.Resize(fyne.NewSize(500, 110))

	ml := newMyLabel("Hello", 30)
	e := newMyEntry("Hello world")
	hello := widget.NewLabel("hello2")

	c := container.NewRow(
		fyne.AxisAlignment{
			MainAxisAlignment:  fyne.MainAxisAlignmentCenter,
			CrossAxisAlignment: fyne.CrossAxisAlignmentCenter,
		},
		ml, widget.NewExpanded(e), hello,
	)
	// c := fyne.NewContainerWithLayout(layout.NewRow(
	// 	layout.MainAxisAlignmentSpaceAround, layout.CrossAxisAlignmentCenter,
	// ), ml, widget.NewExpanded(e), hello)

	// c2 := fyne.NewContainerWithLayout(layout.NewColumn(
	// 	layout.MainAxisAlignmentCenter, layout.CrossAxisAlignmentCenter,
	// ), c, widget.NewButton("HolaBtn", func() {}), widget.NewLabel("Hola3"))

	w.SetContent(c)

	w.ShowAndRun()
}

// ===============================================================
// MyLabel (Implements Baseline interface???)
// ===============================================================

type myLabel struct {
	widget.BaseWidget
	text *canvas.Text
}

func newMyLabel(s string, size float32) *myLabel {
	l := &myLabel{}
	l.text = &canvas.Text{
		Text:     s,
		Color:    color.Black,
		TextSize: size,
	}
	l.ExtendBaseWidget(l)

	return l
}

func (l *myLabel) DistanceToTextBaseline() float32 {
	d := l.text.MinSize().Height - 6
	return d
}

func (l *myLabel) CreateRenderer() fyne.WidgetRenderer {
	l.ExtendBaseWidget(l)
	bg := canvas.NewRectangle(theme.ErrorColor())
	line := canvas.NewRectangle(theme.PrimaryColor())
	return &lrenderer{l, bg, line}
}

type lrenderer struct {
	l        *myLabel
	bg, line *canvas.Rectangle
}

func (*lrenderer) Destroy() {}

func (r *lrenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.bg, r.l.text, r.line}
}

func (r *lrenderer) MinSize() fyne.Size {
	return r.l.text.MinSize().Add(fyne.NewSize(0, theme.Padding()*2))
}

func (r *lrenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	r.bg.Move(fyne.NewPos(0, 0))

	r.line.Resize(fyne.NewSize(size.Width, theme.InputBorderSize()))
	r.line.Move(fyne.NewPos(0, r.l.DistanceToTextBaseline()))

	r.l.text.Move(fyne.NewPos(0, theme.Padding()))
	r.l.text.Resize(r.l.text.MinSize())
}

func (r *lrenderer) Refresh() {
	r.l.text.Refresh()
}

// ===============================================================
// MyEntry (Implements Baseline interface???)
// ===============================================================

type myEntry struct {
	widget.BaseWidget
	entry *widget.Entry
}

func newMyEntry(placeholder string) *myEntry {
	e := &myEntry{}
	e.ExtendBaseWidget(e)
	e.entry = widget.NewEntry()
	e.entry.PlaceHolder = placeholder
	return e
}

func (e *myEntry) DistanceToTextBaseline() float32 {
	return fyne.MeasureText("M", theme.TextSize(), fyne.TextStyle{}).Height + theme.Padding()
}

func (e *myEntry) CreateRenderer() fyne.WidgetRenderer {
	e.ExtendBaseWidget(e)
	objects := []fyne.CanvasObject{e.entry}
	return &myEntryRenderer{objects, e}
}

type myEntryRenderer struct {
	objects []fyne.CanvasObject
	me      *myEntry
}

func (r *myEntryRenderer) Destroy() {}

func (r *myEntryRenderer) MinSize() fyne.Size {
	return r.me.entry.MinSize()
}

func (r *myEntryRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *myEntryRenderer) Layout(size fyne.Size) {
	r.me.entry.Move(fyne.NewPos(0, 0))
	r.me.entry.Resize(size)
}

func (r *myEntryRenderer) Refresh() {
	r.me.entry.Refresh()
}
