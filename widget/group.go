package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

// Group widget is list of widgets that contains a visual border around the list and a group title at the top.
type Group struct {
	baseWidget

	Text    string
	box     *Box
	content fyne.CanvasObject
}

// Resize sets a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (g *Group) Resize(size fyne.Size) {
	g.resize(size, g)
}

// Move the widget to a new position, relative to it's parent.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (g *Group) Move(pos fyne.Position) {
	g.move(pos, g)
}

// MinSize returns the smallest size this widget can shrink to
func (g *Group) MinSize() fyne.Size {
	return g.minSize(g)
}

// Show this widget, if it was previously hidden
func (g *Group) Show() {
	g.show(g)
}

// Hide this widget, if it was previously visible
func (g *Group) Hide() {
	g.hide(g)
}

// Prepend inserts a new CanvasObject at the top of the group
func (g *Group) Prepend(object fyne.CanvasObject) {
	g.box.Prepend(object)

	Refresh(g)
}

// Append adds a new CanvasObject to the end of the group
func (g *Group) Append(object fyne.CanvasObject) {
	g.box.Append(object)

	Refresh(g)
}

// CreateRenderer is a private method to Fyne which links this widget to it's renderer
func (g *Group) CreateRenderer() fyne.WidgetRenderer {
	label := NewLabel(g.Text)
	labelBg := canvas.NewRectangle(theme.BackgroundColor())
	line := canvas.NewRectangle(theme.ButtonColor())
	objects := []fyne.CanvasObject{line, labelBg, label, g.content}
	return &groupRenderer{label: label, line: line, labelBg: labelBg,
		objects: objects, group: g}
}

// NewGroup creates a new grouped list widget with a title and the specified list of child objects.
func NewGroup(title string, children ...fyne.CanvasObject) *Group {
	box := NewVBox(children...)
	group := &Group{baseWidget{}, title, box, box}

	Renderer(group).Layout(group.MinSize())
	return group
}

// NewGroupWithScroller creates a new grouped list widget with a title and the specified list of child objects.
// This group will scroll when the available space is less than needed to display the items it contains.
func NewGroupWithScroller(title string, children ...fyne.CanvasObject) *Group {
	box := NewVBox(children...)
	group := &Group{baseWidget{}, title, box, NewScrollContainer(box)}

	Renderer(group).Layout(group.MinSize())
	return group
}

type groupRenderer struct {
	label         *Label
	line, labelBg *canvas.Rectangle

	objects []fyne.CanvasObject
	group   *Group
}

func (g *groupRenderer) MinSize() fyne.Size {
	labelMin := g.label.MinSize()
	groupMin := g.group.content.MinSize()

	return fyne.NewSize(fyne.Max(labelMin.Width, groupMin.Width),
		labelMin.Height+groupMin.Height+theme.Padding())
}

func (g *groupRenderer) Layout(size fyne.Size) {
	labelWidth := g.label.MinSize().Width
	labelHeight := g.label.MinSize().Height

	g.line.Move(fyne.NewPos(0, labelHeight/2))
	g.line.Resize(fyne.NewSize(size.Width, theme.Padding()))

	g.labelBg.Move(fyne.NewPos(size.Width/2-labelWidth/2, 0))
	g.labelBg.Resize(g.label.MinSize())
	g.label.Move(fyne.NewPos(size.Width/2-labelWidth/2, 0))
	g.label.Resize(g.label.MinSize())

	g.group.content.Move(fyne.NewPos(0, labelHeight+theme.Padding()))
	g.group.content.Resize(fyne.NewSize(size.Width, size.Height-labelHeight-theme.Padding()))
}

func (g *groupRenderer) ApplyTheme() {
	Renderer(g.label).ApplyTheme()
	g.line.FillColor = theme.ButtonColor()
	g.labelBg.FillColor = theme.BackgroundColor()
}

func (g *groupRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (g *groupRenderer) Objects() []fyne.CanvasObject {
	return g.objects
}

func (g *groupRenderer) Refresh() {
	g.label.SetText(g.group.Text)
	g.Layout(g.group.Size())

	canvas.Refresh(g.group)
}

func (g *groupRenderer) Destroy() {
}
