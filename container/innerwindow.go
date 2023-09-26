package container

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	intWidget "fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var _ fyne.Widget = (*InnerWindow)(nil)

// InnerWindow defines a container that wraps content in a window border - that can then be placed inside
// a regular container/canvas.
//
// Since: 2.5
type InnerWindow struct {
	widget.BaseWidget

	CloseIntercept                                      func()
	OnDragged, OnResized                                func(*fyne.DragEvent)
	OnMinimized, OnMaximized, OnTappedBar, OnTappedIcon func()
	Icon                                                fyne.Resource

	title   string
	content fyne.CanvasObject
}

// NewInnerWindow creates a new window border around the given `content`, displaying the `title` along the top.
// This will behave like a normal contain and will probably want to be added to a `MultipleWindows` parent.
//
// Since: 2.5
func NewInnerWindow(title string, content fyne.CanvasObject) *InnerWindow {
	w := &InnerWindow{title: title, content: content}
	w.ExtendBaseWidget(w)
	return w
}

func (w *InnerWindow) Close() {
	w.Hide()
}

func (w *InnerWindow) CreateRenderer() fyne.WidgetRenderer {
	w.ExtendBaseWidget(w)

	min := &widget.Button{Icon: theme.WindowMinimizeIcon(), Importance: widget.LowImportance, OnTapped: w.OnMinimized}
	if w.OnMinimized == nil {
		min.Disable()
	}
	max := &widget.Button{Icon: theme.WindowMaximizeIcon(), Importance: widget.LowImportance, OnTapped: w.OnMaximized}
	if w.OnMaximized == nil {
		max.Disable()
	}

	buttons := NewHBox(
		&widget.Button{Icon: theme.WindowCloseIcon(), Importance: widget.DangerImportance, OnTapped: func() {
			if f := w.CloseIntercept; f != nil {
				f()
			} else {
				w.Close()
			}
		}},
		min, max)

	var icon fyne.CanvasObject
	if w.Icon != nil {
		icon = &widget.Button{Icon: w.Icon, Importance: widget.LowImportance, OnTapped: func() {
			if f := w.OnTappedIcon; f != nil {
				f()
			}
		}}
		if w.OnTappedIcon == nil {
			icon.(*widget.Button).Disable()
		}
	}
	title := newDraggableLabel(w.title, w.OnDragged, w.OnTappedBar)
	title.Truncation = fyne.TextTruncateEllipsis

	bar := NewBorder(nil, nil, buttons, icon, title)
	bg := canvas.NewRectangle(theme.OverlayBackgroundColor())
	contentBG := canvas.NewRectangle(theme.BackgroundColor())
	corner := newDraggableCorner(w.OnResized)

	objects := []fyne.CanvasObject{bg, contentBG, bar, corner, w.content}
	return &innerWindowRenderer{ShadowingRenderer: intWidget.NewShadowingRenderer(objects, intWidget.DialogLevel),
		win: w, bar: bar, bg: bg, corner: corner, contentBG: contentBG}
}

var _ fyne.WidgetRenderer = (*innerWindowRenderer)(nil)

type innerWindowRenderer struct {
	*intWidget.ShadowingRenderer
	min fyne.Size

	win           *InnerWindow
	bar           *fyne.Container
	bg, contentBG *canvas.Rectangle
	corner        fyne.CanvasObject
}

func (i *innerWindowRenderer) Layout(size fyne.Size) {
	pad := theme.Padding()
	pos := fyne.NewSquareOffsetPos(pad / 2)
	size = size.Subtract(fyne.NewSquareSize(pad))
	i.LayoutShadow(size, pos)

	i.bg.Move(pos)
	i.bg.Resize(size)

	barHeight := i.bar.MinSize().Height
	i.bar.Move(pos.AddXY(pad, 0))
	i.bar.Resize(fyne.NewSize(size.Width-pad*2, barHeight))

	innerPos := pos.AddXY(pad, barHeight)
	innerSize := fyne.NewSize(size.Width-pad*2, size.Height-pad-barHeight)
	i.contentBG.Move(innerPos)
	i.contentBG.Resize(innerSize)
	i.win.content.Move(innerPos)
	i.win.content.Resize(innerSize)

	cornerSize := i.corner.MinSize()
	i.corner.Move(pos.Add(size).Subtract(cornerSize))
	i.corner.Resize(cornerSize)
}

func (i *innerWindowRenderer) MinSize() fyne.Size {
	pad := theme.Padding()
	contentMin := i.win.content.MinSize()
	barMin := i.bar.MinSize()

	// only allow windows to grow, as per normal windows
	contentMin = contentMin.Max(i.min)
	i.min = contentMin
	innerWidth := fyne.Max(barMin.Width, contentMin.Width)

	return fyne.NewSize(innerWidth+pad*2, contentMin.Height+pad+barMin.Height).Add(fyne.NewSquareSize(pad))
}

func (i *innerWindowRenderer) Refresh() {
	i.bg.FillColor = theme.OverlayBackgroundColor()
	i.bg.Refresh()
	i.contentBG.FillColor = theme.BackgroundColor()
	i.contentBG.Refresh()
	i.bar.Refresh()

	i.ShadowingRenderer.RefreshShadow()
}

type draggableLabel struct {
	widget.Label
	drag func(*fyne.DragEvent)
	tap  func()
}

func newDraggableLabel(title string, fn func(*fyne.DragEvent), tap func()) *draggableLabel {
	d := &draggableLabel{drag: fn, tap: tap}
	d.ExtendBaseWidget(d)
	d.Text = title
	return d
}

func (d *draggableLabel) Dragged(ev *fyne.DragEvent) {
	if f := d.drag; f != nil {
		d.drag(ev)
	}
}

func (d *draggableLabel) DragEnd() {
}

func (d *draggableLabel) Tapped(ev *fyne.PointEvent) {
	if f := d.tap; f != nil {
		d.tap()
	}
}

type draggableCorner struct {
	widget.BaseWidget
	drag func(*fyne.DragEvent)
}

func newDraggableCorner(fn func(*fyne.DragEvent)) *draggableCorner {
	d := &draggableCorner{drag: fn}
	d.ExtendBaseWidget(d)
	return d
}

func (c *draggableCorner) CreateRenderer() fyne.WidgetRenderer {
	prop := canvas.NewRectangle(color.Transparent)
	prop.SetMinSize(fyne.NewSquareSize(20))
	return widget.NewSimpleRenderer(prop)
}

func (c *draggableCorner) Dragged(ev *fyne.DragEvent) {
	if f := c.drag; f != nil {
		c.drag(ev)
	}
}

func (c *draggableCorner) DragEnd() {
}
