package container

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Navigation container is used to provide your application with a control bar and an area for content objects.
// Objects can be any CanvasObject, and only the most recent one will be visible.
//
// Since: 2.7
type Navigation struct {
	widget.BaseWidget

	Root      fyne.CanvasObject
	Title     string
	OnBack    func()
	OnForward func()

	level  int
	stack  fyne.Container
	titles []string
}

// NewNavigation creates a new navigation container with a given root object.
//
// Since: 2.7
func NewNavigation(root fyne.CanvasObject) *Navigation {
	return NewNavigationWithTitle(root, "")
}

// NewNavigationWithTitle creates a new navigation container with a given root object and a default title.
//
// Since: 2.7
func NewNavigationWithTitle(root fyne.CanvasObject, s string) *Navigation {
	var nav *Navigation
	nav = &Navigation{
		Root:      root,
		Title:     s,
		OnBack:    func() { _ = nav.Back() },
		OnForward: func() { _ = nav.Forward() },
	}
	return nav
}

// Push puts the given object on top of the navigation stack and hides the object below.
//
// Since: 2.7
func (nav *Navigation) Push(obj fyne.CanvasObject) {
	nav.PushWithTitle(obj, nav.Title)
}

// PushWithTitle puts the given CanvasObject on top, hides the object below, and uses the given title as label text.
//
// Since: 2.7
func (nav *Navigation) PushWithTitle(obj fyne.CanvasObject, s string) {
	obj.Show()
	objs := nav.stack.Objects[:nav.level]
	if len(objs) > 0 {
		objs[len(objs)-1].Hide()
	}
	nav.stack.Objects = append(objs, obj)
	nav.titles = append(nav.titles[:nav.level], s)
	nav.level++
	nav.Refresh()
}

// Back returns the top level CanvasObject, adjusts the title accordingly, and disabled the back button
// when no more objects are left to go back to.
//
// Since: 2.7
func (nav *Navigation) Back() fyne.CanvasObject {
	if nav.level == 0 || nav.level == 1 && nav.Root != nil {
		return nil
	}

	objs := nav.stack.Objects
	objs[nav.level-1].Hide()
	if nav.level > 1 {
		objs[nav.level-2].Show()
	}

	nav.level--
	nav.Refresh()

	return objs[nav.level]
}

// Forward shows the next object in the stack again.
//
// Since: 2.7
func (nav *Navigation) Forward() fyne.CanvasObject {
	if nav.level >= len(nav.stack.Objects) {
		return nil
	}

	nav.stack.Objects[nav.level-1].Hide()
	nav.stack.Objects[nav.level].Show()
	nav.level++

	return nav.stack.Objects[nav.level-1]
}

// SetTitle changes the root navigation title shown by default.
//
// Since: 2.7
func (nav *Navigation) SetTitle(s string) {
	nav.Title = s
	nav.Refresh()
}

// SetCurrentTitle changes the navigation title for the current level.
//
// Since: 2.7
func (nav *Navigation) SetCurrentTitle(s string) {
	if nav.level > 1 && nav.level-1 < len(nav.titles) {
		nav.titles[nav.level-1] = s
		nav.Refresh()
	}
}

func (nav *Navigation) setup() {
	objs := []fyne.CanvasObject{}
	titles := []string{}
	if nav.Root != nil {
		objs = append(objs, nav.Root)
		titles = append(titles, nav.Title)
	}
	nav.level = len(objs)
	nav.stack.Layout = layout.NewStackLayout()
	nav.stack.Objects = objs
	nav.titles = titles
	nav.ExtendBaseWidget(nav)
}

var _ fyne.WidgetRenderer = (*navigatorRenderer)(nil)

type navigatorRenderer struct {
	nav     *Navigation
	back    widget.Button
	forward widget.Button
	title   widget.Label
	object  fyne.CanvasObject
}

func (nav *Navigation) CreateRenderer() fyne.WidgetRenderer {
	r := &navigatorRenderer{
		nav: nav,
		title: widget.Label{
			Text:      nav.Title,
			Alignment: fyne.TextAlignCenter,
		},
		back: widget.Button{
			Icon:     theme.NavigateBackIcon(),
			OnTapped: nav.OnBack,
		},
		forward: widget.Button{
			Icon:     theme.NavigateNextIcon(),
			OnTapped: nav.OnForward,
		},
	}
	r.back.Disable()
	r.forward.Disable()

	nav.setup()

	r.object = NewBorder(
		NewStack(NewHBox(&r.back, layout.NewSpacer(), &r.forward), &r.title),
		nil,
		nil,
		nil,
		&nav.stack,
	)

	return r
}

func (r *navigatorRenderer) Destroy() {
}

func (r *navigatorRenderer) Layout(s fyne.Size) {
	r.object.Resize(s)
}

func (r *navigatorRenderer) MinSize() fyne.Size {
	return r.object.MinSize()
}

func (r *navigatorRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.object}
}

func (r *navigatorRenderer) Refresh() {
	if r.nav.level < 1 || r.nav.level == 1 && r.nav.Root != nil {
		r.back.Disable()
	} else {
		r.back.Enable()
	}

	if r.nav.level == len(r.nav.stack.Objects) {
		r.forward.Disable()
	} else {
		r.forward.Enable()
	}

	if r.nav.level-1 >= 0 && r.nav.level-1 < len(r.nav.titles) {
		r.title.Text = r.nav.titles[r.nav.level-1]
	} else {
		r.title.Text = r.nav.Title
	}

	r.object.Refresh()
}
