package widget

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/canvas"
import "github.com/fyne-io/fyne/theme"

const thickness = 2

// Group widget is list of widgets that contains a visual border around the list and a group title at the top.
type Group struct {
	baseWidget

	Text string
	box  *Box
}

// Prepend inserts a new CanvasObject at the top of the group
func (g *Group) Prepend(object fyne.CanvasObject) {
	g.box.Prepend(object)

	g.Renderer().Refresh()
}

// Append adds a new CanvasObject to the end of the group
func (g *Group) Append(object fyne.CanvasObject) {
	g.box.Append(object)

	g.Renderer().Refresh()
}

func (g *Group) createRenderer() fyne.WidgetRenderer {
	label := NewLabel(g.Text)
	border := canvas.NewRectangle(theme.ButtonColor())
	bg := canvas.NewRectangle(theme.BackgroundColor())
	objects := []fyne.CanvasObject{border, bg, label, g.box}
	return &groupRenderer{label: label, border: border, bg: bg, objects: objects, group: g}
}

// Renderer is a private method to Fyne which links this widget to it's renderer
func (g *Group) Renderer() fyne.WidgetRenderer {
	if g.renderer == nil {
		g.renderer = g.createRenderer()
	}

	return g.renderer
}

// NewGroup creates a new grouped list widget with a title and the specified list of child objects
func NewGroup(title string, children ...fyne.CanvasObject) *Group {
	group := &Group{baseWidget{}, title, NewVBox(children...)}

	group.Renderer().Layout(group.MinSize())
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
	g.label.ApplyTheme()
	g.border.FillColor = theme.ButtonColor()
	g.bg.FillColor = theme.BackgroundColor()

	g.group.box.ApplyTheme()
}

func (g *groupRenderer) Objects() []fyne.CanvasObject {
	return g.objects
}

func (g *groupRenderer) Refresh() {
	g.label.Text = g.group.Text
	g.Layout(g.group.CurrentSize())

	canvas.Refresh(g.group)
}
