package widget

import (
	"image/color"

	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/canvas"
	"github.com/fyne-io/fyne/theme"
)

// TabItem represents a single view in a TabContainer.
// The Text is used for the tab button and the Content is shown when the corresponding tab is active.
type TabItem struct {
	Text    string
	Content fyne.CanvasObject
}

// NewTabItem creates a new item for a tabbed widget - each item specifies the content an a label for it's tab.
func NewTabItem(text string, content fyne.CanvasObject) *TabItem {
	return &TabItem{Text: text, Content: content}
}

// TabContainer widget allows switching visible content from a list of TabItems.
// Each item is represented by a button at the top of the widget.
type TabContainer struct {
	baseWidget

	Items   []*TabItem
	current int
}

// SelectTab sets the specified TabItem to be selected and it's content visible.
func (t *TabContainer) SelectTab(item *TabItem) {
	for i, child := range t.Items {
		if child == item {
			t.SelectTabIndex(i)
			return
		}
	}
}

// CurrentTab returns the currently selected TabItem.
func (t *TabContainer) CurrentTab() *TabItem {
	return t.Items[t.current]
}

// SelectTabIndex sets the TabItem at the specifie index to be selected and it's content visible.
func (t *TabContainer) SelectTabIndex(index int) {
	if index < 0 || index >= len(t.Items) {
		return
	}

	t.current = index
	t.Renderer().Refresh()
}

// CurrentTabIndex returns the index of the currently selected TabItem.
func (t *TabContainer) CurrentTabIndex() int {
	return t.current
}

func (t *TabContainer) makeButton(item *TabItem) *Button {
	it := item
	return NewButton(item.Text, func() { t.SelectTab(it) })
}

/*
// Prepend inserts a new CanvasObject at the top of the group
func (t *TabContainer) Prepend(object fyne.CanvasObject) {
	t.Items = append(t.Items, item)
	t.tabBar.Append(t.makeButton(item))
	t.children = append(t.children, item.Content)

	t.Renderer().Refresh()}
}

// Append adds a new CanvasObject to the end of the group
func (t *TabContainer) Append(item TabItem) {
	t.Items = append(t.Items, item)
//	t.tabBar.Append(t.makeButton(item))
//	t.children = append(t.children, item.Content)

	t.Renderer().Refresh()
}
*/

func (t *TabContainer) createRenderer() fyne.WidgetRenderer {
	var contents []fyne.CanvasObject

	buttons := NewHBox()
	for i, item := range t.Items {
		button := t.makeButton(item)
		if i == t.current {
			button.Style = PrimaryButton
		} else {
			item.Content.Hide()
		}
		buttons.Append(button)
		contents = append(contents, item.Content)
	}

	objects := append(contents, buttons)
	return &tabContainerRenderer{tabBar: buttons, objects: objects, container: t}
}

// Renderer is a private method to Fyne which links this widget to it's renderer
func (t *TabContainer) Renderer() fyne.WidgetRenderer {
	if t.renderer == nil {
		t.renderer = t.createRenderer()
	}

	return t.renderer
}

// NewTabContainer creates a new tab bar widget that allows the user to choose between different visible containers
func NewTabContainer(items ...*TabItem) *TabContainer {
	tabs := &TabContainer{baseWidget{}, items, 0}

	tabs.Renderer().Layout(tabs.MinSize())
	return tabs
}

type tabContainerRenderer struct {
	tabBar *Box

	objects   []fyne.CanvasObject
	container *TabContainer
}

func (t *tabContainerRenderer) MinSize() fyne.Size {
	buttonsMin := t.tabBar.MinSize()

	childMin := fyne.NewSize(0, 0)
	for _, child := range t.container.Items {
		childMin = childMin.Union(child.Content.MinSize())
	}

	return fyne.NewSize(fyne.Max(buttonsMin.Width, childMin.Width), buttonsMin.Height+childMin.Height)
}

func (t *tabContainerRenderer) Layout(size fyne.Size) {
	buttonHeight := t.tabBar.MinSize().Height
	t.tabBar.Resize(fyne.NewSize(size.Width, buttonHeight))

	child := t.container.Items[t.container.current].Content
	child.Move(fyne.NewPos(0, buttonHeight))
	child.Resize(fyne.NewSize(size.Width, size.Height-buttonHeight))
}

func (t *tabContainerRenderer) ApplyTheme() {
	t.tabBar.ApplyTheme()

	for _, child := range t.container.Items {
		if wid, ok := child.Content.(fyne.ThemedObject); ok {
			wid.ApplyTheme()
		}
	}
}

func (t *tabContainerRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (t *tabContainerRenderer) Objects() []fyne.CanvasObject {
	return t.objects
}

func (t *tabContainerRenderer) Refresh() {
	t.tabBar.Renderer().Refresh()
	t.Layout(t.container.CurrentSize())

	for i, child := range t.container.Items {
		if i == t.container.current {
			child.Content.Show()
		} else {
			child.Content.Hide()
		}
	}

	canvas.Refresh(t.container)

	for i, button := range t.tabBar.Children {
		if i == t.container.current {
			button.(*Button).Style = PrimaryButton
		} else {
			button.(*Button).Style = DefaultButton
		}
		button.(*Button).Renderer().Refresh()
	}
}
