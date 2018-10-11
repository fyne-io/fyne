package widget

import (
	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/layout"
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
	return NewButtonWithIcon(" ", t.Icon, t.OnActivated)
}

// NewToolbarAction returns a new push button style ToolbarItem
func NewToolbarAction(icon fyne.Resource, onActivated func()) ToolbarItem {
	return &ToolbarAction{icon, onActivated}
}

// ToolbarSpacer is a blank, stretchable space for a toolbar.
// This is typically used to assit layout if you wish some left and some right aligned items.
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

// Toolbar widget creates a horizontal list of tool buttons
type Toolbar struct {
	baseWidget

	Items []ToolbarItem

	box *Box
}

func (t *Toolbar) append(item ToolbarItem) {
	if t.box == nil { // TODO fix smell
		t.Renderer()
	}

	t.box.Append(item.ToolbarObject())
}

func (t *Toolbar) prepend(item ToolbarItem) {
	if t.box == nil { // TODO fix smell
		t.Renderer()
	}

	t.box.Prepend(item.ToolbarObject())
}

// Renderer is a private method to Fyne which links this widget to it's renderer
func (t *Toolbar) Renderer() fyne.WidgetRenderer {
	if t.box == nil {
		t.box = NewHBox()
		for _, item := range t.Items {
			t.append(item)
		}
		t.renderer = t.box.Renderer()
	}

	return t.renderer
}

// Append a new ToolbarItem to the end of this Toolbar
func (t *Toolbar) Append(item ToolbarItem) {
	t.Items = append(t.Items, item)

	t.append(item)
}

// Prepend a new ToolbarItem to the start of this Toolbar
func (t *Toolbar) Prepend(item ToolbarItem) {
	t.Items = append([]ToolbarItem{item}, t.Items...)

	t.prepend(item)
}

// NewToolbar creates a new toolbar widget.
func NewToolbar(items ...ToolbarItem) *Toolbar {

	t := &Toolbar{Items: items}

	t.Renderer().Layout(t.MinSize())
	return t
}
