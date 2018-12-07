package widget

import (
	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/canvas"
	"github.com/fyne-io/fyne/layout"
	"github.com/fyne-io/fyne/theme"
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
		t.renderer.(*boxRenderer).setBackgroundColor(theme.ButtonColor())
	}

	return t.renderer
}

// ApplyTheme updates this widget's visuals to reflect the current theme
func (t *Toolbar) ApplyTheme() {
	for i, item := range t.Items {
		if _, ok := item.(*ToolbarSeparator); ok {
			rect := t.renderer.(*boxRenderer).objects[i].(*canvas.Rectangle)
			rect.FillColor = theme.TextColor()
		}
	}

	t.renderer.ApplyTheme()
	t.renderer.(*boxRenderer).setBackgroundColor(theme.ButtonColor())
	t.renderer.Refresh()
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
