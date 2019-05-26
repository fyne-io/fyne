package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

// TabItem represents a single view in a TabContainer.
// The Text and Icon are used for the tab button and the Content is shown when the corresponding tab is active.
type TabItem struct {
	Text    string
	Icon    fyne.Resource
	Content fyne.CanvasObject
}

// NewTabItem creates a new item for a tabbed widget - each item specifies the content and a label for its tab.
func NewTabItem(text string, content fyne.CanvasObject) *TabItem {
	return &TabItem{Text: text, Content: content}
}

// NewTabItemWithIcon creates a new item for a tabbed widget - each item specifies the content and a label with an icon for its tab.
func NewTabItemWithIcon(text string, icon fyne.Resource, content fyne.CanvasObject) *TabItem {
	return &TabItem{Text: text, Icon: icon, Content: content}
}

// TabLocation ist the location where the tabs of a tab container should be rendered
type TabLocation int

// TabLocation values
const (
	TabLocationTop TabLocation = iota
	TabLocationLeft
	TabLocationBottom
	TabLocationRight
)

// TabContainer widget allows switching visible content from a list of TabItems.
// Each item is represented by a button at the top of the widget.
type TabContainer struct {
	baseWidget

	Items       []*TabItem
	current     int
	tabLocation TabLocation
}

// Resize sets a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (t *TabContainer) Resize(size fyne.Size) {
	t.resize(size, t)
}

// Move the widget to a new position, relative to its parent.
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
	t.SelectTabIndex(t.current)
	Renderer(t).Refresh()
}

// Hide this widget, if it was previously visible
func (t *TabContainer) Hide() {
	t.hide(t)
}

// SelectTab sets the specified TabItem to be selected and its content visible.
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

// SelectTabIndex sets the TabItem at the specific index to be selected and its content visible.
func (t *TabContainer) SelectTabIndex(index int) {
	if index < 0 || index >= len(t.Items) {
		return
	}

	t.current = index

	for i, child := range t.Items {
		if i == t.current {
			child.Content.Show()
		} else {
			child.Content.Hide()
		}
	}

	Refresh(t)
}

// CurrentTabIndex returns the index of the currently selected TabItem.
func (t *TabContainer) CurrentTabIndex() int {
	return t.current
}

func (t *TabContainer) makeButton(item *TabItem) *Button {
	it := item
	return NewButtonWithIcon(item.Text, item.Icon, func() { t.SelectTab(it) })
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

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (t *TabContainer) CreateRenderer() fyne.WidgetRenderer {
	var buttons, objects []fyne.CanvasObject
	for i, item := range t.Items {
		button := t.makeButton(item)
		if i == t.current {
			button.Style = PrimaryButton
		} else {
			item.Content.Hide()
		}
		buttons = append(buttons, button)
		objects = append(objects, item.Content)
	}
	tabBar := t.buildTabBar(buttons)
	line := canvas.NewRectangle(theme.ButtonColor())
	objects = append(objects, line, tabBar)
	return &tabContainerRenderer{tabBar: tabBar, line: line, objects: objects, container: t}
}

func (t *TabContainer) buildTabBar(buttons []fyne.CanvasObject) *Box {
	var tabBar *Box
	switch t.tabLocation {
	case TabLocationLeft, TabLocationRight:
		tabBar = NewVBox()
	default:
		tabBar = NewHBox()
	}
	for _, button := range buttons {
		tabBar.Append(button)
	}
	return tabBar
}

// SetTabLocation sets the location of the tab bar
func (t *TabContainer) SetTabLocation(l TabLocation) {
	t.tabLocation = l
	r := Renderer(t).(*tabContainerRenderer)
	buttons := r.tabBar.Children
	r.tabBar.Children = nil
	r.tabBar = t.buildTabBar(buttons)
	r.objects[len(r.objects)-1] = r.tabBar
	r.Refresh()
}

// NewTabContainer creates a new tab bar widget that allows the user to choose between different visible containers
func NewTabContainer(items ...*TabItem) *TabContainer {
	tabs := &TabContainer{baseWidget: baseWidget{}, Items: items}

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

	switch t.container.tabLocation {
	case TabLocationLeft, TabLocationRight:
		return fyne.NewSize(buttonsMin.Width+childMin.Width+theme.Padding(),
			fyne.Max(buttonsMin.Height, childMin.Height))
	default:
		return fyne.NewSize(fyne.Max(buttonsMin.Width, childMin.Width),
			buttonsMin.Height+childMin.Height+theme.Padding())
	}
}

func (t *tabContainerRenderer) Layout(size fyne.Size) {
	switch t.container.tabLocation {
	case TabLocationTop:
		buttonHeight := t.tabBar.MinSize().Height
		t.tabBar.Move(fyne.NewPos(0, 0))
		t.tabBar.Resize(fyne.NewSize(size.Width, buttonHeight))
		t.line.Move(fyne.NewPos(0, buttonHeight))
		t.line.Resize(fyne.NewSize(size.Width, theme.Padding()))

		child := t.container.Items[t.container.current].Content
		barHeight := buttonHeight + theme.Padding()
		child.Move(fyne.NewPos(0, barHeight))
		child.Resize(fyne.NewSize(size.Width, size.Height-barHeight))
	case TabLocationLeft:
		buttonWidth := t.tabBar.MinSize().Width
		t.tabBar.Move(fyne.NewPos(0, 0))
		t.tabBar.Resize(fyne.NewSize(buttonWidth, size.Height))
		t.line.Move(fyne.NewPos(buttonWidth, 0))
		t.line.Resize(fyne.NewSize(theme.Padding(), size.Height))

		child := t.container.Items[t.container.current].Content
		barWidth := buttonWidth + theme.Padding()
		child.Move(fyne.NewPos(barWidth, 0))
		child.Resize(fyne.NewSize(size.Width-barWidth, size.Height))
	case TabLocationBottom:
		buttonHeight := t.tabBar.MinSize().Height
		t.tabBar.Move(fyne.NewPos(0, size.Height-buttonHeight))
		t.tabBar.Resize(fyne.NewSize(size.Width, buttonHeight))
		barHeight := buttonHeight + theme.Padding()
		t.line.Move(fyne.NewPos(0, size.Height-barHeight))
		t.line.Resize(fyne.NewSize(size.Width, theme.Padding()))

		child := t.container.Items[t.container.current].Content
		child.Move(fyne.NewPos(0, 0))
		child.Resize(fyne.NewSize(size.Width, size.Height-barHeight))
	case TabLocationRight:
		buttonWidth := t.tabBar.MinSize().Width
		t.tabBar.Move(fyne.NewPos(size.Width-buttonWidth, 0))
		t.tabBar.Resize(fyne.NewSize(buttonWidth, size.Height))
		barWidth := buttonWidth + theme.Padding()
		t.line.Move(fyne.NewPos(size.Width-barWidth, 0))
		t.line.Resize(fyne.NewSize(theme.Padding(), size.Height))

		child := t.container.Items[t.container.current].Content
		child.Move(fyne.NewPos(0, 0))
		child.Resize(fyne.NewSize(size.Width-barWidth, size.Height))
	}
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
	t.Layout(t.container.Size().Union(t.container.MinSize()))

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
