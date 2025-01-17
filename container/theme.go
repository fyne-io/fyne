package container

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/cache"
	intTheme "fyne.io/fyne/v2/internal/theme"
	"fyne.io/fyne/v2/widget"
)

// ThemeOverride is a container where the child widgets are themed by the specified theme.
// Containers will be traversed and all child widgets will reflect the theme in this container.
// This should be used sparingly to avoid a jarring user experience.
//
// Since: 2.5
type ThemeOverride struct {
	widget.BaseWidget

	Content fyne.CanvasObject
	Theme   fyne.Theme

	holder *fyne.Container

	mobile bool
}

// NewThemeOverride provides a container where the child widgets are themed by the specified theme.
// Containers will be traversed and all child widgets will reflect the theme in this container.
// This should be used sparingly to avoid a jarring user experience.
//
// If the content `obj` of this theme override is a container and items are later added to the container or any
// sub-containers ensure that you call `Refresh()` on this `ThemeOverride` to ensure the new items match the theme.
//
// Since: 2.5
func NewThemeOverride(obj fyne.CanvasObject, th fyne.Theme) *ThemeOverride {
	t := &ThemeOverride{Content: obj, Theme: th, holder: NewStack(obj)}
	t.ExtendBaseWidget(t)

	cache.OverrideTheme(obj, addFeatures(th, t))
	obj.Refresh() // required as the widgets passed in could have been initially rendered with default theme
	return t
}

func (t *ThemeOverride) CreateRenderer() fyne.WidgetRenderer {
	cache.OverrideTheme(t.Content, addFeatures(t.Theme, t))

	return widget.NewSimpleRenderer(t.holder)
}

func (t *ThemeOverride) Refresh() {
	if t.holder.Objects[0] != t.Content {
		t.holder.Objects[0] = t.Content
		t.holder.Refresh()
	}

	cache.OverrideTheme(t.Content, addFeatures(t.Theme, t))
	t.BaseWidget.Refresh()
}

// SetDeviceIsMobile allows a ThemeOverride container to shape the contained widgets as a mobile device.
// This will impact containers such as AppTabs and DocTabs, and more in the future, to display a layout
// that would automatically be used for a mobile device runtime.
//
// Since: 2.6
func (t *ThemeOverride) SetDeviceIsMobile(on bool) {
	t.mobile = on
	t.BaseWidget.Refresh()
}

type featureTheme struct {
	fyne.Theme

	over *ThemeOverride
}

func addFeatures(th fyne.Theme, o *ThemeOverride) fyne.Theme {
	return &featureTheme{Theme: th, over: o}
}

func (f *featureTheme) Feature(n intTheme.FeatureName) any {
	if n == intTheme.FeatureNameDeviceIsMobile {
		return f.over.mobile
	}

	return nil
}
