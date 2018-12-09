package widget

import (
	"image/color"

	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/canvas"
	"github.com/fyne-io/fyne/theme"
)

const thickness = 2

// Group widget is list of widgets that contains a visual border around the list and a group title at the top.
type Group struct {
	baseWidget

	Text string
	box  *Box
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

	Renderer(g).Refresh()
}

// Append adds a new CanvasObject to the end of the group
func (g *Group) Append(object fyne.CanvasObject) {
	g.box.Append(object)

	Renderer(g).Refresh()
}

// CreateRenderer is a private method to Fyne which links this widget to it's renderer
func (g *Group) CreateRenderer() fyne.WidgetRenderer {
	label := NewLabel(g.Text)
	border := canvas.NewRectangle(theme.ButtonColor())
	bg := canvas.NewRectangle(theme.BackgroundColor())
	objects := []fyne.CanvasObject{border, bg, label, g.box}
	return &groupRenderer{label: label, border: border, bg: bg, objects: objects, group: g}
}

// NewGroup creates a new grouped list widget with a title and the specified list of child objects
func NewGroup(title string, children ...fyne.CanvasObject) *Group {
	group := &Group{baseWidget{}, title, NewVBox(children...)}

	Renderer(group).Layout(group.MinSize())
	return group
}

type groupRenderer struct {
	label      *Label
	border, bg *canvas.Rectangle

	objects []fyne.CanvasObject
	group   *Group
}

func (g *groupRenderer) MinSize() fyne.Size {
	labelMin := g.label.MinSize()
	groupMin := g.group.box.MinSize()

	return fyne.NewSize(fyne.Max(labelMin.Width, groupMin.Width)+(theme.Padding()*2)+(thickness*2),
		labelMin.Height+groupMin.Height+theme.Padding()*2+thickness)
}

func (g *groupRenderer) Layout(size fyne.Size) {
	labelHeight := g.label.MinSize().Height
	halfHeight := labelHeight / 2

	g.border.Move(fyne.NewPos(0, halfHeight))
	g.border.Resize(fyne.NewSize(size.Width, size.Height-halfHeight))
	g.bg.Move(fyne.NewPos(thickness, halfHeight+thickness))
	g.bg.Resize((fyne.NewSize(size.Width-thickness*2, size.Height-halfHeight-thickness*2)))

	g.label.Move(fyne.NewPos(theme.Padding()+thickness, 0))
	g.label.Resize(g.label.MinSize())

	g.group.box.Move(fyne.NewPos(theme.Padding()+thickness, labelHeight+theme.Padding()))
	g.group.box.Resize(fyne.NewSize(size.Width-(theme.Padding()*2)-(thickness*2),
		size.Height-labelHeight-theme.Padding()*2-thickness))
}

func (g *groupRenderer) ApplyTheme() {
	Renderer(g.label).ApplyTheme()
	g.border.FillColor = theme.ButtonColor()
	g.bg.FillColor = theme.BackgroundColor()

	Renderer(g.group.box).ApplyTheme()
}

func (g *groupRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (g *groupRenderer) Objects() []fyne.CanvasObject {
	return g.objects
}

func (g *groupRenderer) Refresh() {
	g.label.Text = g.group.Text
	g.Layout(g.group.Size())

	canvas.Refresh(g.group)
}
