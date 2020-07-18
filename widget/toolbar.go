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
	button := NewButtonWithIcon("", t.Icon, t.OnActivated)
	button.HideShadow = true

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
	Items []ToolbarItem
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (t *Toolbar) CreateRenderer() fyne.WidgetRenderer {
	t.ExtendBaseWidget(t)
	r := &toolbarRenderer{toolbar: t, layout: layout.NewHBoxLayout()}
	r.resetObjects()
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
			r.objs = append(r.objs, item.ToolbarObject())
		}
	}
	r.SetObjects(r.objs)
}
