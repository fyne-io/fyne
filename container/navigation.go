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
// Since: 2.6
type Navigation struct {
	widget.BaseWidget

	Root  fyne.CanvasObject
	Title string
	Back  NavigationObject
	Next  NavigationObject
	Label *widget.Label

	level  int
	stack  *fyne.Container
	titles []string
}

// NavigationObject allows using any object implementing the `fyne.Disableable`,
// `fyne.CanvasObject`, and `fyne.Tappable` interfaces to be used as Back/Next
// navigation control.
//
// Since: 2.6
type NavigationObject interface {
	fyne.CanvasObject
	fyne.Disableable
	fyne.Tappable
}

// NewNavigation creates a new navigation container with a given root object.
//
// Since: 2.6
func NewNavigation(root fyne.CanvasObject) *Navigation {
	return NewNavigationWithTitle(root, "")
}

// NewNavigation creates a new navigation container with a given root object and a default title.
//
// Since: 2.6
func NewNavigationWithTitle(root fyne.CanvasObject, s string) *Navigation {
	label := widget.NewLabelWithStyle(s, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	objs := []fyne.CanvasObject{}
	titles := []string{}
	if root != nil {
		objs = append(objs, root)
		titles = append(titles, s)
	}

	nav := &Navigation{
		Root:   root,
		Title:  s,
		Label:  label,
		level:  len(objs),
		stack:  NewStack(objs...),
		titles: titles,
	}
	nav.ExtendBaseWidget(nav)
	nav.Back = widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() { _ = nav.Pop() })
	nav.Back.Disable()
	nav.Next = widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() { _ = nav.Forward() })
	nav.Next.Disable()

	return nav
}

// Push puts the given object on top of the navigation stack and hides the object below.
//
// Since: 2.6
func (nav *Navigation) Push(obj fyne.CanvasObject) {
	s := nav.Title
	if nav.level > 0 {
		s = nav.titles[nav.level-1]
	}
	nav.PushWithTitle(obj, s)
}

// PushWithTitle puts the given CanvasObject on top, hides the object below, and uses the given title as label text.
//
// Since: 2.6
func (nav *Navigation) PushWithTitle(obj fyne.CanvasObject, s string) {
	objs := nav.stack.Objects[:nav.level]
	if len(objs) > 0 {
		objs[len(objs)-1].Hide()
	}
	nav.stack.Objects = append(objs, obj)
	if len(nav.stack.Objects) > 0 {
		if nav.Back != nil {
			nav.Back.Enable()
		}
	}
	nav.titles = append(nav.titles[:nav.level], s)
	nav.level++
	nav.Label.SetText(s)
	if nav.Next != nil {
		nav.Next.Disable()
	}
	obj.Show()
}

// Pop return the top level CanvasObject, adjusts the title accordingly, and disabled the back button
// when no more objects are left to go back to.
//
// Since: 2.6
func (nav *Navigation) Pop() fyne.CanvasObject {
	if nav.level == 0 || nav.level == 1 && nav.Root != nil {
		return nil
	}

	title := nav.Title
	objs := nav.stack.Objects
	objs[nav.level-1].Hide()
	if nav.level > 1 {
		objs[nav.level-2].Show()
		title = nav.titles[nav.level-2]
	}
	nav.Label.SetText(title)

	if nav.level == 1 || nav.level == 2 && nav.Root != nil {
		if nav.Back != nil {
			nav.Back.Disable()
		}
	}

	if nav.Next != nil {
		nav.Next.Enable()
	}
	nav.level--

	return objs[nav.level]
}

// Forward shows the next object in the again.
//
// Since: 2.6
func (nav *Navigation) Forward() fyne.CanvasObject {
	nav.stack.Objects[nav.level-1].Hide()
	nav.stack.Objects[nav.level].Show()

	s := nav.Title
	if len(nav.titles) > 0 {
		s = nav.titles[nav.level]
	}
	nav.level++
	nav.Label.SetText(s)

	if nav.level == len(nav.stack.Objects) {
		if nav.Next != nil {
			nav.Next.Disable()
		}
	}
	if nav.level > 0 {
		if nav.Back != nil {
			nav.Back.Enable()
		}
	}

	return nav.stack.Objects[nav.level-1]
}

// SetTitle changes the navigation title and the title for the current object.
//
// Since: 2.6
func (nav *Navigation) SetTitle(s string) {
	nav.titles[nav.level] = s
	nav.Label.SetText(s)
}

func (nav *Navigation) CreateRenderer() fyne.WidgetRenderer {
	control := NewStack(NewHBox(nav.Back, layout.NewSpacer(), nav.Next), nav.Label)
	return widget.NewSimpleRenderer(NewBorder(control, nil, nil, nil, nav.stack))
}
