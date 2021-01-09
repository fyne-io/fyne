package container

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// Declare conformity with Widget interface.
var _ fyne.Widget = (*DocTabs)(nil)

// DocTabs container is used to display various pieces of content identified by tabs.
// The tabs contain text and/or an icon and allow the user to switch between the content specified in each TabItem.
// Each item is represented by a button at the edge of the container.
//
// Since: 2.0.0
type DocTabs struct {
	baseTabs
}

// NewDocTabs creates a new tab container that allows the user to choose between various pieces of content.
//
// Since: 2.0.0
func NewDocTabs(items ...*TabItem) *DocTabs {
	tabs := &DocTabs{
		baseTabs: baseTabs{
			BaseWidget: widget.BaseWidget{},
			current:    -1,
		},
	}
	tabs.ExtendBaseWidget(tabs)
	tabs.SetItems(items)
	return tabs
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (t *DocTabs) CreateRenderer() fyne.WidgetRenderer {
	t.ExtendBaseWidget(t)
	r := &docTabsRenderer{
		baseTabsRenderer: baseTabsRenderer{
			bar:         &fyne.Container{},
			buttonCache: make(map[*TabItem]*tabButton),
			divider:     canvas.NewRectangle(theme.ShadowColor()),
			indicator:   canvas.NewRectangle(theme.PrimaryColor()),
		},
		docTabs: t,
	}
	r.action = r.buildAllTabsButton()
	r.updateTabs()
	r.moveIndicator()
	return r
}

// MinSize returns the size that this widget should not shrink below
func (t *DocTabs) MinSize() fyne.Size {
	t.ExtendBaseWidget(t)
	return t.BaseWidget.MinSize()
}

// Declare conformity with WidgetRenderer interface.
var _ fyne.WidgetRenderer = (*docTabsRenderer)(nil)

type docTabsRenderer struct {
	baseTabsRenderer
	docTabs *DocTabs
}

func (r *docTabsRenderer) Layout(size fyne.Size) {
	r.updateTabs()
	r.layout(&r.docTabs.baseTabs, size)
	r.moveIndicator()
}

func (r *docTabsRenderer) MinSize() fyne.Size {
	return r.minSize(&r.docTabs.baseTabs)
}

func (r *docTabsRenderer) Objects() []fyne.CanvasObject {
	objects := r.baseTabsRenderer.Objects()
	if i, is := r.docTabs.current, r.docTabs.Items; i >= 0 && i < len(is) {
		objects = append(objects, is[i].Content)
	}
	return objects
}

func (r *docTabsRenderer) Refresh() {
	// TODO Offset Scroller so current tab is always visible

	r.Layout(r.docTabs.Size())

	r.refresh(&r.docTabs.baseTabs)

	canvas.Refresh(r.docTabs)
}

func (r *docTabsRenderer) buildAllTabsButton() (all *widget.Button) {
	all = widget.NewButton("", func() {
		// Show pop up containing all tabs

		var items []*fyne.MenuItem
		for i := 0; i < len(r.docTabs.Items); i++ {
			item := r.docTabs.Items[i]
			// FIXME MenuItem doesn't support icons (#1752)
			// FIXME MenuItem can't show if it is the currently selected tab (#1753
			items = append(items, fyne.NewMenuItem(item.Text, func() {
				r.docTabs.Select(item)
				r.docTabs.popUp = nil
			}))
		}

		r.docTabs.showPopUp(all, items)
	})
	all.Importance = widget.LowImportance
	return
}

func (r *docTabsRenderer) moveIndicator() {
	var selectedPos fyne.Position
	var selectedSize fyne.Size

	scroller := r.bar.Objects[0].(*widget.ScrollContainer)
	buttons := scroller.Content.(*fyne.Container).Objects

	if r.docTabs.current >= 0 {
		selected := buttons[r.docTabs.current]
		selectedPos = selected.Position()
		selectedSize = selected.Size()
	}

	offset := scroller.Offset

	var indicatorPos fyne.Position
	var indicatorSize fyne.Size

	switch r.docTabs.tabLocation {
	case TabLocationTop:
		indicatorPos = fyne.NewPos(selectedPos.X-offset.X, r.bar.MinSize().Height)
		indicatorSize = fyne.NewSize(selectedSize.Width, theme.Padding())
	case TabLocationLeading:
		indicatorPos = fyne.NewPos(r.bar.MinSize().Width, selectedPos.Y-offset.Y)
		indicatorSize = fyne.NewSize(theme.Padding(), selectedSize.Height)
	case TabLocationBottom:
		indicatorPos = fyne.NewPos(selectedPos.X-offset.X, r.bar.Position().Y-theme.Padding())
		indicatorSize = fyne.NewSize(selectedSize.Width, theme.Padding())
	case TabLocationTrailing:
		indicatorPos = fyne.NewPos(r.bar.Position().X-theme.Padding(), selectedPos.Y-offset.Y)
		indicatorSize = fyne.NewSize(theme.Padding(), selectedSize.Height)
	}

	r.animateIndicator(indicatorPos, indicatorSize)
}

func (r *docTabsRenderer) updateTabs() {
	tabCount := len(r.docTabs.Items)

	buttons := r.buildTabButtons(&r.docTabs.baseTabs, tabCount)

	// FIXME Need to create a new scroller each time because scroller can't switch directions (#1751)
	var scroller *widget.ScrollContainer

	// Set layout of tab bar containing tab buttons and overflow action
	if r.docTabs.tabLocation == TabLocationLeading || r.docTabs.tabLocation == TabLocationTrailing {
		r.bar.Layout = layout.NewBorderLayout(nil, r.action, nil, nil)
		scroller = NewVScroll(buttons)
	} else {
		r.bar.Layout = layout.NewBorderLayout(nil, nil, nil, r.action)
		scroller = NewHScroll(buttons)
	}

	// TODO Need to attach an onOffsetChanged callback to scroller to trigger moveIndicator()

	r.bar.Objects = []fyne.CanvasObject{scroller}
	if a := r.action; a != nil {
		r.bar.Objects = append(r.bar.Objects, a)
	}

	r.bar.Refresh()
}
