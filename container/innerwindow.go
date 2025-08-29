package container

import (
	"image/color"
	"runtime"
	"strings"

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

	// Alignment allows an inner window to specify if the buttons should be on the left
	// (`ButtonAlignLeading`) or right of the window border.
	//
	// Since: 2.6
	Alignment widget.ButtonAlign

	title     string
	content   *fyne.Container
	maximized bool
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

	borderIcon := newBorderButton(w.Icon, modeIcon, th, func() {
		if f := w.OnTappedIcon; f != nil {
			f()
		}
	})
	if w.OnTappedIcon == nil {
		borderIcon.Disable()
	}

	if w.Icon == nil {
		borderIcon.Hide()
	}
	title := newDraggableLabel(w.title, w)
	title.Truncation = fyne.TextTruncateEllipsis

	height := w.Theme().Size(theme.SizeNameWindowTitleBarHeight)
	off := (height - title.labelMinSize().Height) / 2
	barMid := New(layout.NewCustomPaddedLayout(off, 0, 0, 0), title)
	if w.buttonPosition() == widget.ButtonAlignTrailing {
		buttons = NewCenter(NewHBox(min, max, close))
	}

	bg := canvas.NewRectangle(th.Color(theme.ColorNameOverlayBackground, v))
	contentBG := canvas.NewRectangle(th.Color(theme.ColorNameBackground, v))
	corner := newDraggableCorner(w)
	bar := New(&titleBarLayout{buttons: buttons, icon: borderIcon, title: barMid, win: w},
		buttons, borderIcon, barMid)

	if w.content == nil {
		w.content = NewPadded(canvas.NewRectangle(color.Transparent))
	}
	objects := []fyne.CanvasObject{bg, contentBG, bar, w.content, corner}
	r := &innerWindowRenderer{
		ShadowingRenderer: intWidget.NewShadowingRenderer(objects, intWidget.DialogLevel),
		win:               w, bar: bar, buttonBox: buttons, buttons: []*borderButton{close, min, max}, bg: bg,
		corner: corner, contentBG: contentBG, icon: borderIcon,
	}
	r.Layout(w.Size())
	return r
}

func (w *InnerWindow) SetContent(obj fyne.CanvasObject) {
	w.content.Objects[0] = obj

	w.content.Refresh()
}

// SetMaximized tells the window if the maximized state should be set or not.
//
// Since: 2.6
func (w *InnerWindow) SetMaximized(max bool) {
	w.maximized = max
	w.Refresh()
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

func (w *InnerWindow) buttonPosition() widget.ButtonAlign {
	if w.Alignment != widget.ButtonAlignCenter {
		return w.Alignment
	}

	if runtime.GOOS == "windows" || runtime.GOOS == "linux" || strings.Contains(runtime.GOOS, "bsd") {
		return widget.ButtonAlignTrailing
	}
	// macOS
	return widget.ButtonAlignLeading
}

var _ fyne.WidgetRenderer = (*innerWindowRenderer)(nil)

type innerWindowRenderer struct {
	*intWidget.ShadowingRenderer

	win            *InnerWindow
	bar, buttonBox *fyne.Container
	buttons        []*borderButton
	icon           *borderButton
	bg, contentBG  *canvas.Rectangle
	corner         fyne.CanvasObject
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

	if i.win.buttonPosition() == widget.ButtonAlignTrailing {
		i.buttonBox.Objects[0].(*fyne.Container).Objects = []fyne.CanvasObject{i.buttons[1], i.buttons[2], i.buttons[0]}
	} else {
		i.buttonBox.Objects[0].(*fyne.Container).Objects = []fyne.CanvasObject{i.buttons[0], i.buttons[1], i.buttons[2]}
	}
	for _, b := range i.buttons {
		b.setTheme(th)
	}
	i.bar.Refresh()

	if i.win.OnMinimized == nil {
		i.buttons[1].Disable()
	} else {
		i.buttons[1].SetOnTapped(i.win.OnMinimized)
		i.buttons[1].Enable()
	}

	max := i.buttons[2]
	if i.win.OnMaximized == nil {
		i.buttons[2].Disable()
	} else {
		max.SetOnTapped(i.win.OnMaximized)
		max.Enable()
	}
	if i.win.maximized {
		max.b.SetIcon(theme.ViewRestoreIcon())
	} else {
		max.b.SetIcon(theme.WindowMaximizeIcon())
	}

	title := i.bar.Objects[2].(*fyne.Container).Objects[0].(*draggableLabel)
	title.SetText(i.win.title)
	i.ShadowingRenderer.RefreshShadow()
	if i.win.OnTappedIcon == nil {
		i.icon.Disable()
	} else {
		i.icon.Enable()
	}
	if i.win.Icon != nil {
		i.icon.b.SetIcon(i.win.Icon)
		i.icon.Show()
	} else {
		i.icon.Hide()
	}
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
	buttonImportance := widget.MediumImportance
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

func (b *borderButton) Enable() {
	b.b.Enable()
}

func (b *borderButton) SetOnTapped(fn func()) {
	b.b.OnTapped = fn
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

type titleBarLayout struct {
	win                  *InnerWindow
	buttons, icon, title fyne.CanvasObject
}

func (t *titleBarLayout) Layout(_ []fyne.CanvasObject, s fyne.Size) {
	buttonMinWidth := t.buttons.MinSize().Width
	t.buttons.Resize(fyne.NewSize(buttonMinWidth, s.Height))
	t.icon.Resize(fyne.NewSquareSize(s.Height))
	usedWidth := buttonMinWidth
	if t.icon.Visible() {
		usedWidth += s.Height
	}
	t.title.Resize(fyne.NewSize(s.Width-usedWidth, s.Height))

	if t.win.buttonPosition() == widget.ButtonAlignTrailing {
		t.buttons.Move(fyne.NewPos(s.Width-buttonMinWidth, 0))
		t.icon.Move(fyne.Position{})
		if t.icon.Visible() {
			t.title.Move(fyne.NewPos(s.Height, 0))
		} else {
			t.title.Move(fyne.Position{})
		}
	} else {
		t.buttons.Move(fyne.NewPos(0, 0))
		t.icon.Move(fyne.NewPos(s.Width-s.Height, 0))
		t.title.Move(fyne.NewPos(buttonMinWidth, 0))
	}
}

func (t *titleBarLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	buttonMin := t.buttons.MinSize()
	iconMin := t.icon.MinSize()
	titleMin := t.title.MinSize() // can truncate

	return fyne.NewSize(buttonMin.Width+iconMin.Width+titleMin.Width,
		fyne.Max(fyne.Max(buttonMin.Height, iconMin.Height), titleMin.Height))
}
