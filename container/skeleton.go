package container

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// AppSkeleton is a container that assists in laying out a complete application.
//
// since: 2.3
type AppSkeleton struct {
	widget.BaseWidget
	content *fyne.Container
}

// NewAppSkeleton creates a widget that lays out a standard application container.
// The content passed in will fill the space left after any accessories are packed in.
//
// Since: 2.3
func NewAppSkeleton(content fyne.CanvasObject) *AppSkeleton {
	a := &AppSkeleton{content: NewBorder(nil, nil, nil, nil, content)}
	a.ExtendBaseWidget(a)
	return a
}

// CreateRenderer is an internal method for returning the information for rendering this skeleton.
func (a *AppSkeleton) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(a.content)
}

// SetToolbar is used to specify that the app skeleton should have a tool bar placed in the content.
// It will appear in the top of the application, above any other content but under a menu, if set.
func (a *AppSkeleton) SetToolbar(t *widget.Toolbar) {
	a.content.Add(t)
	a.content.Layout = layout.NewBorderLayout(t, nil, nil, nil)
	a.content.Refresh()
}
