package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/theme"
)

// Group widget contains a list of widgets that are grouped under a dividing line and title at the top.
type Group struct {
	BaseWidget

	Text    string
	box     *Box
	content fyne.CanvasObject
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

// MinSize returns the size that this widget should not shrink below
func (g *Group) MinSize() fyne.Size {
	g.ExtendBaseWidget(g)
	return g.BaseWidget.MinSize()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (g *Group) CreateRenderer() fyne.WidgetRenderer {
	g.ExtendBaseWidget(g)
	label := NewLabel(g.Text)
	labelBg := canvas.NewRectangle(theme.BackgroundColor())
	line := canvas.NewRectangle(theme.ShadowColor())
	return &groupRenderer{
		BaseRenderer: widget.NewBaseRenderer([]fyne.CanvasObject{line, labelBg, label, g.content}),
		label:        label,
		line:         line,
		labelBg:      labelBg,
		group:        g,
	}
}

// NewGroup creates a new grouped list widget with a title and the specified list of child objects.
// Deprecated: Consider using the Card instead.
func NewGroup(title string, children ...fyne.CanvasObject) *Group {
	box := NewVBox(children...)
	group := &Group{BaseWidget{}, title, box, box}

	Renderer(group).Layout(group.MinSize())
	return group
}

// NewGroupWithScroller creates a new grouped list widget with a title and the specified list of child objects.
// This group will scroll when the available space is less than needed to display the items it contains.
// Deprecated: Consider using the Card instead.
func NewGroupWithScroller(title string, children ...fyne.CanvasObject) *Group {
	box := NewVBox(children...)
	group := &Group{BaseWidget{}, title, box, NewVScrollContainer(box)}

	Renderer(group).Layout(group.MinSize())
	return group
}

type groupRenderer struct {
	widget.BaseRenderer
	label         *Label
	line, labelBg *canvas.Rectangle
	group         *Group
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

func (g *groupRenderer) Refresh() {
	g.line.FillColor = theme.ShadowColor()
	g.labelBg.FillColor = theme.BackgroundColor()

	g.label.SetText(g.group.Text)
	g.Layout(g.group.Size())

	g.line.Refresh()
	g.labelBg.Refresh()
}
