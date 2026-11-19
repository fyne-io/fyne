package widget

import (
	"fyne.io/fyne/v2"
	internalWidget "fyne.io/fyne/v2/internal/widget"
)

type Scroll struct {
	BaseWidget

	fakeScroll      *internalWidget.Scroll
	displayedScroll *cachedScroll
}

func NewHScroll(content fyne.CanvasObject) *Scroll {
	s := &Scroll{
		fakeScroll:      internalWidget.NewHScroll(content),
		displayedScroll: newCachedHScroll(content),
	}
	s.ExtendBaseWidget(s)
	return s
}

func NewVScroll(content fyne.CanvasObject) *Scroll {
	s := &Scroll{
		fakeScroll:      internalWidget.NewVScroll(content),
		displayedScroll: newCachedVScroll(content),
	}
	s.ExtendBaseWidget(s)
	return s
}

func NewScroll(content fyne.CanvasObject) *Scroll {
	s := &Scroll{
		fakeScroll:      internalWidget.NewScroll(content),
		displayedScroll: newCachedScroll(content),
	}
	s.ExtendBaseWidget(s)
	return s
}

func (s *Scroll) CreateRenderer() fyne.WidgetRenderer {
	return newScrollRenderer(s)
}

type cachedScroll struct {
	*internalWidget.Scroll

	content      fyne.CanvasObject
	hasContainer bool
}

func newCachedHScroll(content fyne.CanvasObject) *cachedScroll {
	s := &cachedScroll{
		Scroll:       internalWidget.NewHScroll(newScrollLayout(content)),
		content:      content,
		hasContainer: isContainer(content),
	}
	s.OnScrolled = func(fyne.Position) { s.update() }
	return s
}

func newCachedVScroll(content fyne.CanvasObject) *cachedScroll {
	s := &cachedScroll{
		Scroll:       internalWidget.NewVScroll(newScrollLayout(content)),
		content:      content,
		hasContainer: isContainer(content),
	}
	s.OnScrolled = func(fyne.Position) { s.update() }
	return s
}

func newCachedScroll(content fyne.CanvasObject) *cachedScroll {
	s := &cachedScroll{
		Scroll:       internalWidget.NewScroll(newScrollLayout(content)),
		content:      content,
		hasContainer: isContainer(content),
	}
	s.OnScrolled = func(fyne.Position) { s.update() }
	return s
}

func isContainer(content fyne.CanvasObject) bool {
	switch content.(type) {
	case fyne.Widget:
		return false
	case *fyne.Container:
		return true
	}
	panic("the canvas object is nor a widget nor a container")
}

func (s *cachedScroll) update() {
	if s.hasContainer == true {
		s.Content.(*fyne.Container).Objects = s.Content.(*fyne.Container).Objects[:0]
	}
	s.Content.(*fyne.Container).Objects = s.Content.(*fyne.Container).Objects[:0]
	s.updateCache(s.content)
}

func (s *cachedScroll) updateCache(object fyne.CanvasObject) {
	if s.canBeRendered(object) == false {
		return
	}
	switch o := object.(type) {
	case fyne.Widget:
		s.Content.(*fyne.Container).Objects = append(s.Content.(*fyne.Container).Objects, o)
		return
	case *fyne.Container:
		for _, v := range o.Objects {
			s.updateCache(v)
		}
		return
	}
	panic("the canvas object is nor a widget nor a container")
}

func (s *cachedScroll) canBeRendered(object fyne.CanvasObject) bool {
	op := fyne.Position{
		X: object.Position().X - s.Offset.X,
		Y: object.Position().Y - s.Offset.Y,
	}
	return object.Visible() &&
		op.X <= s.Size().Width &&
		op.Y <= s.Size().Height &&
		op.X+object.Size().Width >= 0 &&
		op.Y+object.Size().Height >= 0
}

type scrollLayout struct {
	content fyne.CanvasObject
}

func newScrollLayout(content fyne.CanvasObject) *fyne.Container {
	return &fyne.Container{
		Layout: &scrollLayout{content: content},
	}
}

func (l *scrollLayout) MinSize([]fyne.CanvasObject) fyne.Size {
	return l.content.Size()
}

func (l *scrollLayout) Layout([]fyne.CanvasObject, fyne.Size) {
}

type scrollRenderer struct {
	scrollWidget *Scroll
	layoutSize   fyne.Size
}

func newScrollRenderer(scrollWidget *Scroll) *scrollRenderer {
	r := &scrollRenderer{
		scrollWidget: scrollWidget,
	}
	return r
}

func (r *scrollRenderer) Destroy() {
}

func (r *scrollRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.scrollWidget.displayedScroll.Scroll}
}

func (r *scrollRenderer) SetObjects(objects []fyne.CanvasObject) {
}

func (r *scrollRenderer) Layout(size fyne.Size) {
	if r.layoutSize.Height == size.Height &&
		r.layoutSize.Width == size.Width {
		return
	}
	r.layoutSize = size

	r.scrollWidget.fakeScroll.Resize(size)
	r.scrollWidget.displayedScroll.Resize(size)
	r.scrollWidget.displayedScroll.update()
}

func (r *scrollRenderer) MinSize() fyne.Size {
	return r.scrollWidget.fakeScroll.MinSize()
}

func (r *scrollRenderer) Refresh() {
}
