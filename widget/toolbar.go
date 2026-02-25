package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/layout"
)

// ToolbarItem represents any interface element that can be added to a toolbar
type ToolbarItem interface {
	ToolbarObject() fyne.CanvasObject
}

// ToolbarAction is push button style of ToolbarItem
type ToolbarAction struct {
	Icon        fyne.Resource
	OnActivated func() `json:"-"`
	button      Button
}

// ToolbarObject gets a button to render this ToolbarAction
func (t *ToolbarAction) ToolbarObject() fyne.CanvasObject {
	t.button.Importance = LowImportance

	// synchronize properties
	t.button.Icon = t.Icon
	t.button.OnTapped = t.OnActivated

	return &t.button
}

// SetIcon updates the icon on a ToolbarItem
//
// Since: 2.2
func (t *ToolbarAction) SetIcon(icon fyne.Resource) {
	t.Icon = icon
	t.button.SetIcon(t.Icon)
}

// Enable this ToolbarAction, updating any style or features appropriately.
//
// Since: 2.5
func (t *ToolbarAction) Enable() {
	t.button.Enable()
}

// Disable this ToolbarAction so that it cannot be interacted with, updating any style appropriately.
//
// Since: 2.5
func (t *ToolbarAction) Disable() {
	t.button.Disable()
}

// Disabled returns true if this ToolbarAction is currently disabled or false if it can currently be interacted with.
//
// Since: 2.5
func (t *ToolbarAction) Disabled() bool {
	return t.button.Disabled()
}

// NewToolbarAction returns a new push button style ToolbarItem
func NewToolbarAction(icon fyne.Resource, onActivated func()) *ToolbarAction {
	return &ToolbarAction{Icon: icon, OnActivated: onActivated}
}

// ToolbarSpacer is a blank, stretchable space for a toolbar.
// This is typically used to assist layout if you wish some left and some right aligned items.
// Space will be split evebly amongst all the spacers on a toolbar.
type ToolbarSpacer struct{}

// ToolbarObject gets the actual spacer object for this ToolbarSpacer
func (t *ToolbarSpacer) ToolbarObject() fyne.CanvasObject {
	return layout.NewSpacer()
}

// NewToolbarSpacer returns a new spacer item for a Toolbar to assist with ToolbarItem alignment
func NewToolbarSpacer() *ToolbarSpacer {
	return &ToolbarSpacer{}
}

// ToolbarSeparator is a thin, visible divide that can be added to a Toolbar.
// This is typically used to assist visual grouping of ToolbarItems.
type ToolbarSeparator struct{}

// ToolbarObject gets the visible line object for this ToolbarSeparator
func (t *ToolbarSeparator) ToolbarObject() fyne.CanvasObject {
	return &Separator{invert: true}
}

// NewToolbarSeparator returns a new separator item for a Toolbar to assist with ToolbarItem grouping
func NewToolbarSeparator() *ToolbarSeparator {
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
	items   []fyne.CanvasObject
	toolbar *Toolbar
}

func (r *toolbarRenderer) MinSize() fyne.Size {
	return r.layout.MinSize(r.items)
}

func (r *toolbarRenderer) Layout(size fyne.Size) {
	r.layout.Layout(r.items, size)
}

func (r *toolbarRenderer) Refresh() {
	r.resetObjects()
	canvas.Refresh(r.toolbar)
}

func (r *toolbarRenderer) resetObjects() {
	r.items = make([]fyne.CanvasObject, 0, len(r.toolbar.Items))
	for _, item := range r.toolbar.Items {
		r.items = append(r.items, item.ToolbarObject())
	}
	r.SetObjects(r.items)
}
