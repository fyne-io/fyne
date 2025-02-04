package container

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Navigation container is used to provide your application with a control bar and an area for content objects.
// Objects can be any CanvasObject, and only the most recent one will be visible.
//
// Since: 2.6
type Navigation struct {
	widget.BaseWidget
	control *fyne.Container
	stack   *fyne.Container
	button  *widget.Button
	label   *widget.Label
	title   string
	titles  []string
}

// NewNavigation creates a new navigation container allowing to pass in one or more objects.
//
// Since: 2.6
func NewNavigation(s string, objs ...fyne.CanvasObject) *Navigation {
	button := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), nil)
	label := widget.NewLabelWithStyle(s, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	titles := make([]string, len(objs))
	for n := 0; n < len(objs); n++ {
		titles[n] = s
	}
	nav := &Navigation{
		button:  button,
		control: NewStack(NewHBox(button), label),
		stack:   NewStack(objs...),
		label:   label,
		title:   s,
		titles:  titles,
	}
	nav.ExtendBaseWidget(nav)
	nav.button.OnTapped = func() {
		nav.Pop()
	}
	if len(objs) <= 1 {
		nav.button.Disable()
	}

	return nav
}

// Push puts the given object on top of the navigation stack and hides the object below.
//
// Since: 2.6
func (nav *Navigation) Push(obj fyne.CanvasObject) {
	nav.PushWithTitle(obj, nav.titles[len(nav.titles)-1])
}

// PushWithTitle puts the given CanvasObject on top, hides the object below, and uses the given title as label text.
//
// Since: 2.6
func (nav *Navigation) PushWithTitle(obj fyne.CanvasObject, s string) {
	objs := nav.stack.Objects
	if len(objs) > 0 {
		objs[len(objs)-1].Hide()
	}
	nav.stack.Objects = append(objs, obj)
	if len(nav.stack.Objects) > 1 {
		nav.button.Enable()
	}

	nav.titles = append(nav.titles, s)
	nav.label.SetText(s)
}

// Pop return the top level CanvasObject, adjusts the title accordingly, and disabled the back button
// when no more objects are left to go back to.
//
// Since: 2.6
func (nav *Navigation) Pop() fyne.CanvasObject {
	objs := nav.stack.Objects
	if len(objs) == 0 {
		return nil
	}
	obj := objs[len(objs)-1]
	objs = objs[:len(objs)-1]
	if len(objs) > 0 {
		objs[len(objs)-1].Show()
	}
	nav.stack.Objects = objs
	if len(nav.stack.Objects) <= 1 {
		nav.button.Disable()
	}

	if len(nav.titles) > 0 {
		nav.titles = nav.titles[:len(nav.titles)-1]
	}
	title := nav.title
	if len(nav.titles) > 0 {
		title = nav.titles[len(nav.titles)-1]
	}
	nav.label.SetText(title)

	return obj
}

// SetTitle changes the navigation title and the title for the current object.
//
// Since: 2.6
func (nav *Navigation) SetTitle(s string) {
	nav.title = s
	nav.titles[len(nav.titles)-1] = s
	nav.label.SetText(s)
}

func (nav *Navigation) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(NewBorder(nav.control, nil, nil, nil, nav.stack))
}
