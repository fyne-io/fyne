package widget

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/canvas"
import "github.com/fyne-io/fyne/layout"

// Box widget is a simple list where the child elements are arranged in a single column
// for vertical or a single row for horizontal arrangement
type Box struct {
	baseWidget

	Horizontal bool
	Children   []fyne.CanvasObject
}

// Prepend inserts a new CanvasObject at the top/left of the box
func (b *Box) Prepend(object fyne.CanvasObject) {
	b.Children = append([]fyne.CanvasObject{object}, b.Children...)

	b.Renderer().Refresh()
}

// Append adds a new CanvasObject to the end/right of the box
func (b *Box) Append(object fyne.CanvasObject) {
	b.Children = append(b.Children, object)

	b.Renderer().Refresh()
}

func (b *Box) createRenderer() fyne.WidgetRenderer {
	var lay fyne.Layout
	if b.Horizontal {
		lay = layout.NewHBoxLayout()
	} else {
		lay = layout.NewVBoxLayout()
	}
	return &boxRenderer{objects: b.Children, layout: lay, box: b}
}

// Renderer is a private method to Fyne which links this widget to it's renderer
func (b *Box) Renderer() fyne.WidgetRenderer {
	if b.renderer == nil {
		b.renderer = b.createRenderer()
	}

	return b.renderer
}

// NewHBox creates a new horizontally aligned box widget with the specified list of child objects
func NewHBox(children ...fyne.CanvasObject) *Box {
	box := &Box{baseWidget{}, true, children}

	box.Renderer().Layout(box.MinSize())
	return box
}

// NewVBox creates a new vertically aligned box widget with the specified list of child objects
func NewVBox(children ...fyne.CanvasObject) *Box {
	box := &Box{baseWidget{}, false, children}

	box.Renderer().Layout(box.MinSize())
	return box
}

type boxRenderer struct {
	layout fyne.Layout

	objects []fyne.CanvasObject
	box     *Box
}

func (b *boxRenderer) MinSize() fyne.Size {
	return b.layout.MinSize(b.objects)
}

func (b *boxRenderer) Layout(size fyne.Size) {
	b.layout.Layout(b.objects, size)
}

// ApplyTheme is a fallback method that applies the new theme to all contained
// objects. Widgets that override this should consider doing similarly.
func (b *boxRenderer) ApplyTheme() {
	for _, child := range b.objects {
		switch themed := child.(type) {
		case fyne.ThemedObject:
			themed.ApplyTheme()
		}
	}
}

func (b *boxRenderer) Objects() []fyne.CanvasObject {
	return b.objects
}

func (b *boxRenderer) Refresh() {
	b.objects = b.box.Children
	b.Layout(b.box.CurrentSize())

	canvas.Refresh(b.box)
}
