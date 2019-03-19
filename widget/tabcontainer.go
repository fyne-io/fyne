package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
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

// Resize sets a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (t *TabContainer) Resize(size fyne.Size) {
	t.resize(size, t)
}

// Move the widget to a new position, relative to it's parent.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (t *TabContainer) Move(pos fyne.Position) {
	t.move(pos, t)
}

// MinSize returns the smallest size this widget can shrink to
func (t *TabContainer) MinSize() fyne.Size {
	return t.minSize(t)
}

// Show this widget, if it was previously hidden
func (t *TabContainer) Show() {
	t.show(t)
}

// Hide this widget, if it was previously visible
func (t *TabContainer) Hide() {
	t.hide(t)
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
	Refresh(t)
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

	t.CreateRenderer().Refresh()}
}

// Append adds a new CanvasObject to the end of the group
func (t *TabContainer) Append(item TabItem) {
	t.Items = append(t.Items, item)
//	t.tabBar.Append(t.makeButton(item))
//	t.children = append(t.children, item.Content)

	t.CreateRenderer().Refresh()
}
*/

// CreateRenderer is a private method to Fyne which links this widget to it's renderer
func (t *TabContainer) CreateRenderer() fyne.WidgetRenderer {
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

	line := canvas.NewRectangle(theme.ButtonColor())
	objects := append(contents, line, buttons)
	return &tabContainerRenderer{tabBar: buttons, line: line, objects: objects, container: t}
}

// NewTabContainer creates a new tab bar widget that allows the user to choose between different visible containers
func NewTabContainer(items ...*TabItem) *TabContainer {
	tabs := &TabContainer{baseWidget{}, items, 0}

	Renderer(tabs).Layout(tabs.MinSize())
	return tabs
}

type tabContainerRenderer struct {
	tabBar *Box
	line   *canvas.Rectangle

	objects   []fyne.CanvasObject
	container *TabContainer
}

func (t *tabContainerRenderer) MinSize() fyne.Size {
	buttonsMin := t.tabBar.MinSize()

	childMin := fyne.NewSize(0, 0)
	for _, child := range t.container.Items {
		childMin = childMin.Union(child.Content.MinSize())
	}

	return fyne.NewSize(fyne.Max(buttonsMin.Width, childMin.Width),
		buttonsMin.Height+childMin.Height+theme.Padding())
}

func (t *tabContainerRenderer) Layout(size fyne.Size) {
	buttonHeight := t.tabBar.MinSize().Height
	t.tabBar.Resize(fyne.NewSize(size.Width, buttonHeight))
	t.line.Move(fyne.NewPos(0, buttonHeight))
	t.line.Resize(fyne.NewSize(size.Width, theme.Padding()))

	child := t.container.Items[t.container.current].Content
	child.Move(fyne.NewPos(0, buttonHeight+theme.Padding()))
	child.Resize(fyne.NewSize(size.Width, size.Height-buttonHeight-theme.Padding()))
}

func (t *tabContainerRenderer) ApplyTheme() {
	t.line.FillColor = theme.ButtonColor()
}

func (t *tabContainerRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (t *tabContainerRenderer) Objects() []fyne.CanvasObject {
	return t.objects
}

func (t *tabContainerRenderer) Refresh() {
	Renderer(t.tabBar).Refresh()
	t.Layout(t.container.Size())

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
		Renderer(button.(*Button)).Refresh()
	}
}

func (t *tabContainerRenderer) Destroy() {
}
