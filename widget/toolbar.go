package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
)

// ToolbarItem represents any interface element that can be added to a toolbar
type ToolbarItem interface {
	ToolbarObject() fyne.CanvasObject
}

// ToolbarAction is push button style of ToolbarItem
type ToolbarAction struct {
	Icon        fyne.Resource
	OnActivated func()
}

// ToolbarObject gets a button to render this ToolbarAction
func (t *ToolbarAction) ToolbarObject() fyne.CanvasObject {
	//button := NewButtonWithIcon("", t.Icon, t.OnActivated)
	//button.HideShadow = true
	button := newToolbarButton(t.Icon, t.OnActivated)
	return button
}

// NewToolbarAction returns a new push button style ToolbarItem
func NewToolbarAction(icon fyne.Resource, onActivated func()) ToolbarItem {
	return &ToolbarAction{icon, onActivated}
}

// ToolbarSpacer is a blank, stretchable space for a toolbar.
// This is typically used to assist layout if you wish some left and some right aligned items.
// Space will be split evebly amongst all the spacers on a toolbar.
type ToolbarSpacer struct {
}

// ToolbarObject gets the actual spacer object for this ToolbarSpacer
func (t *ToolbarSpacer) ToolbarObject() fyne.CanvasObject {
	return layout.NewSpacer()
}

// NewToolbarSpacer returns a new spacer item for a Toolbar to assist with ToolbarItem alignment
func NewToolbarSpacer() ToolbarItem {
	return &ToolbarSpacer{}
}

// ToolbarSeparator is a thin, visible divide that can be added to a Toolbar.
// This is typically used to assist visual grouping of ToolbarItems.
type ToolbarSeparator struct {
}

// ToolbarObject gets the visible line object for this ToolbarSeparator
func (t *ToolbarSeparator) ToolbarObject() fyne.CanvasObject {
	return canvas.NewRectangle(theme.TextColor())
}

// NewToolbarSeparator returns a new separator item for a Toolbar to assist with ToolbarItem grouping
func NewToolbarSeparator() ToolbarItem {
	return &ToolbarSeparator{}
}

// Toolbar widget creates a horizontal list of tool buttons
type Toolbar struct {
	BaseWidget
	Items    []ToolbarItem
	focused  bool
	current  int
	renderer *toolbarRenderer
}

// FocusGained is called when the Entry has been given focus.
func (b *Toolbar) FocusGained() {
	b.focused = true
	if b.current < len(b.Items) {
		if o, ok := b.renderer.objs[b.current].(*ToolbarButton); ok {
			o.focused = true
			o.Refresh()
		}
	}
	b.Refresh()
}

// FocusLost is called when the Entry has had focus removed.
func (b *Toolbar) FocusLost() {
	b.focused = false
	if b.current < len(b.renderer.objs) {
		bt := b.renderer.objs[b.current]
		if o, ok := bt.(*ToolbarButton); ok {
			o.focused = false
			o.Refresh()
		}
	}
	b.Refresh()
}

// Focused returns whether or not this Entry has focus.
func (b *Toolbar) Focused() bool {
	return b.focused
}

func (b *Toolbar) TypedRune(r rune) {
}

func (t *Toolbar) changeFocusedButton(delta int) {
	i := t.current
	ok := false
	for i+delta >= 0 && i+delta < len(t.renderer.objs) && !ok {
		i = i + delta
		_, ok = t.renderer.objs[i].(*ToolbarButton)
	}
	if ok {
		t.current = i
	}
}

func (b *Toolbar) TypedKey(key *fyne.KeyEvent) {
	if bt, ok := b.renderer.objs[b.current].(*ToolbarButton); ok {
		bt.focused = false
		bt.Refresh()
		if key.Name == fyne.KeyReturn || key.Name == fyne.KeyEnter || key.Name == fyne.KeySpace {
			bt.OnTap()
		}
	}
	if key.Name == fyne.KeyLeft || key.Name == fyne.KeyUp {
		b.changeFocusedButton(-1)
	} else if key.Name == fyne.KeyRight || key.Name == fyne.KeyDown {
		b.changeFocusedButton(+1)
	}

	if bt, ok := b.renderer.objs[b.current].(*ToolbarButton); ok {
		bt.focused = true
		bt.Refresh()
	}
	b.Refresh()
}

func (b *Toolbar) KeyUp(key *fyne.KeyEvent) {
	if key.Name == fyne.KeyReturn || key.Name == fyne.KeyEnter || key.Name == fyne.KeySpace {
		if btn, ok := b.renderer.objs[b.current].(*ToolbarButton); ok {
			btn.pressed = false
			btn.Refresh()
		}
	}
}

func (b *Toolbar) KeyDown(key *fyne.KeyEvent) {
	if key.Name == fyne.KeyReturn || key.Name == fyne.KeyEnter || key.Name == fyne.KeySpace {
		if btn, ok := b.renderer.objs[b.current].(*ToolbarButton); ok {
			btn.pressed = true
			btn.Refresh()
		}
	}
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (t *Toolbar) CreateRenderer() fyne.WidgetRenderer {
	t.ExtendBaseWidget(t)
	r := &toolbarRenderer{toolbar: t, layout: layout.NewHBoxLayout()}
	r.resetObjects()
	t.renderer = r
	return r
}

// Append a new ToolbarItem to the end of this Toolbar
func (t *Toolbar) Append(item ToolbarItem) {
	t.Items = append(t.Items, item)
	t.Refresh()
}

// Prepend a new ToolbarItem to the start of this Toolbar
func (t *Toolbar) Prepend(item ToolbarItem) {
	t.Items = append([]ToolbarItem{item}, t.Items...)
	t.Refresh()
}

// MinSize returns the size that this widget should not shrink below
func (t *Toolbar) MinSize() fyne.Size {
	t.ExtendBaseWidget(t)
	return t.BaseWidget.MinSize()
}

// NewToolbar creates a new toolbar widget.
func NewToolbar(items ...ToolbarItem) *Toolbar {
	t := &Toolbar{Items: items}
	t.ExtendBaseWidget(t)
	t.Refresh()
	return t
}

type toolbarRenderer struct {
	widget.BaseRenderer
	layout  fyne.Layout
	objs    []fyne.CanvasObject
	toolbar *Toolbar
}

func (r *toolbarRenderer) MinSize() fyne.Size {
	return r.layout.MinSize(r.Objects())
}

func (r *toolbarRenderer) Layout(size fyne.Size) {
	r.layout.Layout(r.Objects(), size)
}

func (r *toolbarRenderer) BackgroundColor() color.Color {
	return theme.ButtonColor()
}

func (r *toolbarRenderer) Refresh() {
	r.resetObjects()
	for i, item := range r.toolbar.Items {
		if _, ok := item.(*ToolbarSeparator); ok {
			rect := r.Objects()[i].(*canvas.Rectangle)
			rect.FillColor = theme.TextColor()
		}
	}

	canvas.Refresh(r.toolbar)
}

func (r *toolbarRenderer) resetObjects() {
	if len(r.objs) != len(r.toolbar.Items) {
		r.objs = make([]fyne.CanvasObject, 0, len(r.toolbar.Items))
		for _, item := range r.toolbar.Items {
			o := item.ToolbarObject()
			if b,ok := o.(*ToolbarButton); ok {
				b.toolbar = r.toolbar
			}
			r.objs = append(r.objs, o)
		}
	}
	r.SetObjects(r.objs)
}
