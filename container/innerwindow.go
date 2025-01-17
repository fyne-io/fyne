package container

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	intWidget "fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/layout"
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

	CloseIntercept                                      func()                `json:"-"`
	OnDragged, OnResized                                func(*fyne.DragEvent) `json:"-"`
	OnMinimized, OnMaximized, OnTappedBar, OnTappedIcon func()                `json:"-"`
	Icon                                                fyne.Resource

	title   string
	content *fyne.Container
}

// NewInnerWindow creates a new window border around the given `content`, displaying the `title` along the top.
// This will behave like a normal contain and will probably want to be added to a `MultipleWindows` parent.
//
// Since: 2.5
func NewInnerWindow(title string, content fyne.CanvasObject) *InnerWindow {
	w := &InnerWindow{title: title, content: NewPadded(content)}
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
	title := newDraggableLabel(w.title, w)
	title.Truncation = fyne.TextTruncateEllipsis
	th := w.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	bar := NewBorder(nil, nil, buttons, icon, title)
	bg := canvas.NewRectangle(th.Color(theme.ColorNameOverlayBackground, v))
	contentBG := canvas.NewRectangle(th.Color(theme.ColorNameBackground, v))
	corner := newDraggableCorner(w)

	objects := []fyne.CanvasObject{bg, contentBG, bar, w.content, corner}
	return &innerWindowRenderer{ShadowingRenderer: intWidget.NewShadowingRenderer(objects, intWidget.DialogLevel),
		win: w, bar: bar, bg: bg, corner: corner, contentBG: contentBG}
}

func (w *InnerWindow) SetContent(obj fyne.CanvasObject) {
	w.content.Objects[0] = obj

	w.content.Refresh()
}

func (w *InnerWindow) SetPadded(pad bool) {
	if pad {
		w.content.Layout = layout.NewPaddedLayout()
	} else {
		w.content.Layout = layout.NewStackLayout()
	}
	w.content.Refresh()
}

func (w *InnerWindow) SetTitle(title string) {
	w.title = title
	w.Refresh()
}

var _ fyne.WidgetRenderer = (*innerWindowRenderer)(nil)

type innerWindowRenderer struct {
	*intWidget.ShadowingRenderer

	win           *InnerWindow
	bar           *fyne.Container
	bg, contentBG *canvas.Rectangle
	corner        fyne.CanvasObject
}

func (i *innerWindowRenderer) Layout(size fyne.Size) {
	th := i.win.Theme()
	pad := th.Size(theme.SizeNamePadding)

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
	i.corner.Move(pos.Add(size).Subtract(cornerSize).AddXY(1, 1))
	i.corner.Resize(cornerSize)
}

func (i *innerWindowRenderer) MinSize() fyne.Size {
	th := i.win.Theme()
	pad := th.Size(theme.SizeNamePadding)
	contentMin := i.win.content.MinSize()
	barMin := i.bar.MinSize()

	innerWidth := fyne.Max(barMin.Width, contentMin.Width)

	return fyne.NewSize(innerWidth+pad*2, contentMin.Height+pad+barMin.Height).Add(fyne.NewSquareSize(pad))
}

func (i *innerWindowRenderer) Refresh() {
	th := i.win.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()
	i.bg.FillColor = th.Color(theme.ColorNameOverlayBackground, v)
	i.bg.Refresh()
	i.contentBG.FillColor = th.Color(theme.ColorNameBackground, v)
	i.contentBG.Refresh()
	i.bar.Refresh()

	title := i.bar.Objects[0].(*draggableLabel)
	title.SetText(i.win.title)
	i.ShadowingRenderer.RefreshShadow()
}

type draggableLabel struct {
	widget.Label
	win *InnerWindow
}

func newDraggableLabel(title string, win *InnerWindow) *draggableLabel {
	d := &draggableLabel{win: win}
	d.ExtendBaseWidget(d)
	d.Text = title
	return d
}

func (d *draggableLabel) Dragged(ev *fyne.DragEvent) {
	if f := d.win.OnDragged; f != nil {
		f(ev)
	}
}

func (d *draggableLabel) DragEnd() {
}

func (d *draggableLabel) Tapped(ev *fyne.PointEvent) {
	if f := d.win.OnTappedBar; f != nil {
		f()
	}
}

type draggableCorner struct {
	widget.BaseWidget
	win *InnerWindow
}

func newDraggableCorner(w *InnerWindow) *draggableCorner {
	d := &draggableCorner{win: w}
	d.ExtendBaseWidget(d)
	return d
}

func (c *draggableCorner) CreateRenderer() fyne.WidgetRenderer {
	prop := canvas.NewImageFromResource(fyne.CurrentApp().Settings().Theme().Icon(theme.IconNameDragCornerIndicator))
	prop.SetMinSize(fyne.NewSquareSize(16))
	return widget.NewSimpleRenderer(prop)
}

func (c *draggableCorner) Dragged(ev *fyne.DragEvent) {
	if f := c.win.OnResized; f != nil {
		c.win.OnResized(ev)
	}
}

func (c *draggableCorner) DragEnd() {
}
