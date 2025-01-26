package container

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	intWidget "fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type titleBarButtonMode int

const (
	modeClose titleBarButtonMode = iota
	modeMinimize
	modeMaximize
	modeIcon
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
	th := w.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	min := newBorderButton(theme.WindowMinimizeIcon(), modeMinimize, th, w.OnMinimized)
	if w.OnMinimized == nil {
		min.Disable()
	}
	max := newBorderButton(theme.WindowMaximizeIcon(), modeMaximize, th, w.OnMaximized)
	if w.OnMaximized == nil {
		max.Disable()
	}

	close := newBorderButton(theme.WindowCloseIcon(), modeClose, th, func() {
		if f := w.CloseIntercept; f != nil {
			f()
		} else {
			w.Close()
		}
	})
	buttons := NewCenter(NewHBox(close, min, max))

	var icon fyne.CanvasObject

	if w.Icon != nil {
		icon = newBorderButton(w.Icon, modeIcon, th, func() {
			if f := w.OnTappedIcon; f != nil {
				f()
			}
		})
		if w.OnTappedIcon == nil {
			icon.(*borderButton).Disable()
		}
	}
	title := newDraggableLabel(w.title, w)
	title.Truncation = fyne.TextTruncateEllipsis

	height := w.Theme().Size(theme.SizeNameWindowTitleBarHeight)
	off := (height - title.labelMinSize().Height) / 2
	bar := NewBorder(nil, nil, buttons, icon,
		New(layout.NewCustomPaddedLayout(off, 0, 0, 0), title))
	bg := canvas.NewRectangle(th.Color(theme.ColorNameOverlayBackground, v))
	contentBG := canvas.NewRectangle(th.Color(theme.ColorNameBackground, v))
	corner := newDraggableCorner(w)

	objects := []fyne.CanvasObject{bg, contentBG, bar, w.content, corner}
	return &innerWindowRenderer{ShadowingRenderer: intWidget.NewShadowingRenderer(objects, intWidget.DialogLevel),
		win: w, bar: bar, buttons: []*borderButton{min, max, close}, bg: bg, corner: corner, contentBG: contentBG}
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
	buttons       []*borderButton
	bg, contentBG *canvas.Rectangle
	corner        fyne.CanvasObject
}

func (i *innerWindowRenderer) Layout(size fyne.Size) {
	th := i.win.Theme()
	pad := th.Size(theme.SizeNamePadding)

	i.LayoutShadow(size, fyne.Position{})
	i.bg.Resize(size)

	barHeight := i.win.Theme().Size(theme.SizeNameWindowTitleBarHeight)
	i.bar.Move(fyne.NewPos(pad, 0))
	i.bar.Resize(fyne.NewSize(size.Width-pad*2, barHeight))

	innerPos := fyne.NewPos(pad, barHeight)
	innerSize := fyne.NewSize(size.Width-pad*2, size.Height-pad-barHeight)
	i.contentBG.Move(innerPos)
	i.contentBG.Resize(innerSize)
	i.win.content.Move(innerPos)
	i.win.content.Resize(innerSize)

	cornerSize := i.corner.MinSize()
	i.corner.Move(fyne.NewPos(size.Components()).Subtract(cornerSize).AddXY(1, 1))
	i.corner.Resize(cornerSize)
}

func (i *innerWindowRenderer) MinSize() fyne.Size {
	th := i.win.Theme()
	pad := th.Size(theme.SizeNamePadding)
	contentMin := i.win.content.MinSize()
	barHeight := th.Size(theme.SizeNameWindowTitleBarHeight)

	innerWidth := fyne.Max(i.bar.MinSize().Width, contentMin.Width)

	return fyne.NewSize(innerWidth+pad*2, contentMin.Height+pad+barHeight)
}

func (i *innerWindowRenderer) Refresh() {
	th := i.win.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()
	i.bg.FillColor = th.Color(theme.ColorNameOverlayBackground, v)
	i.bg.Refresh()
	i.contentBG.FillColor = th.Color(theme.ColorNameBackground, v)
	i.contentBG.Refresh()

	for _, b := range i.buttons {
		b.setTheme(th)
	}
	i.bar.Refresh()

	title := i.bar.Objects[0].(*fyne.Container).Objects[0].(*draggableLabel)
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

func (d *draggableLabel) MinSize() fyne.Size {
	width := d.Label.MinSize().Width
	height := d.Label.Theme().Size(theme.SizeNameWindowButtonHeight)
	return fyne.NewSize(width, height)
}

func (d *draggableLabel) Tapped(_ *fyne.PointEvent) {
	if f := d.win.OnTappedBar; f != nil {
		f()
	}
}

func (d *draggableLabel) labelMinSize() fyne.Size {
	return d.Label.MinSize()
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

type borderButton struct {
	widget.BaseWidget

	b    *widget.Button
	c    *ThemeOverride
	mode titleBarButtonMode
}

func newBorderButton(icon fyne.Resource, mode titleBarButtonMode, th fyne.Theme, fn func()) *borderButton {
	buttonImportance := widget.DangerImportance
	if mode == modeIcon {
		buttonImportance = widget.LowImportance
	}
	b := &widget.Button{Icon: icon, Importance: buttonImportance, OnTapped: fn}
	c := NewThemeOverride(b, &buttonTheme{Theme: th, mode: mode})

	ret := &borderButton{b: b, c: c, mode: mode}
	ret.ExtendBaseWidget(ret)
	return ret
}

func (b *borderButton) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(b.c)
}

func (b *borderButton) Disable() {
	b.b.Disable()
}

func (b *borderButton) MinSize() fyne.Size {
	height := b.Theme().Size(theme.SizeNameWindowButtonHeight)
	return fyne.NewSquareSize(height)
}

func (b *borderButton) setTheme(th fyne.Theme) {
	b.c.Theme = &buttonTheme{Theme: th, mode: b.mode}
}

type buttonTheme struct {
	fyne.Theme
	mode titleBarButtonMode
}

func (b *buttonTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	switch n {
	case theme.ColorNameError:
		n = theme.ColorNameWindowButton
	case theme.ColorNameHover:
		if b.mode == modeClose {
			n = theme.ColorNameError
		}
	}
	return b.Theme.Color(n, v)
}

func (b *buttonTheme) Size(n fyne.ThemeSizeName) float32 {
	switch n {
	case theme.SizeNameInputRadius:
		if b.mode == modeIcon {
			return 0
		}
		n = theme.SizeNameWindowButtonRadius
	case theme.SizeNameInlineIcon:
		n = theme.SizeNameWindowButtonIcon
	}

	return b.Theme.Size(n)
}
