package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal"
	"fyne.io/fyne/internal/cache"
	"fyne.io/fyne/layout"
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

// TabLocation is the location where the tabs of a tab container should be rendered
type TabLocation int

// TabLocation values
const (
	TabLocationTop TabLocation = iota
	TabLocationLeading
	TabLocationBottom
	TabLocationTrailing
)

// TabContainer widget allows switching visible content from a list of TabItems.
// Each item is represented by a button at the top of the widget.
type TabContainer struct {
	BaseWidget

	Items       []*TabItem
	current     int
	tabLocation TabLocation
}

// Show this widget, if it was previously hidden
func (t *TabContainer) Show() {
	t.BaseWidget.Show()
	t.SelectTabIndex(t.current)
	t.refresh(t)
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

	r := cache.Renderer(t).(*tabContainerRenderer)
	r.Layout(t.size)
	t.refresh(t)
}

// CurrentTabIndex returns the index of the currently selected TabItem.
func (t *TabContainer) CurrentTabIndex() int {
	return t.current
}

func (t *TabContainer) makeButton(item *TabItem) *tabButton {
	return &tabButton{Text: item.Text, Icon: item.Icon, OnTap: func() { t.SelectTab(item) }}
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (t *TabContainer) CreateRenderer() fyne.WidgetRenderer {
	t.ExtendBaseWidget(t)
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
	return &tabContainerRenderer{tabBar: tabBar, line: line, objects: objects, container: t}
}

// MinSize returns the size that this widget should not shrink below
func (t *TabContainer) MinSize() fyne.Size {
	t.ExtendBaseWidget(t)
	return t.BaseWidget.MinSize()
}

func (t *TabContainer) buildTabBar(buttons []fyne.CanvasObject) *fyne.Container {
	var lay fyne.Layout
	if fyne.CurrentDevice().IsMobile() {
		lay = layout.NewGridLayout(len(buttons))
	} else if t.tabLocation == TabLocationLeading || t.tabLocation == TabLocationTrailing {
		lay = layout.NewVBoxLayout()
	} else {
		lay = layout.NewHBoxLayout()
	}

	tabBar := fyne.NewContainerWithLayout(lay)
	for _, button := range buttons {
		tabBar.AddObject(button)
	}
	return tabBar
}

// Append adds a new TabItem to the rightmost side of the tab panel
func (t *TabContainer) Append(item *TabItem) {
	r := cache.Renderer(t).(*tabContainerRenderer)
	t.Items = append(t.Items, item)
	r.objects = append(r.objects, item.Content)
	r.tabBar.Objects = append(r.tabBar.Objects, t.makeButton(item))

	t.Refresh()
}

// Remove tab by value
func (t *TabContainer) Remove(item *TabItem) {
	for index, existingItem := range t.Items {
		if existingItem == item {
			t.RemoveIndex(index)
			break
		}
	}
}

// RemoveIndex removes tab by index
func (t *TabContainer) RemoveIndex(index int) {
	r := cache.Renderer(t).(*tabContainerRenderer)
	t.Items = append(t.Items[:index], t.Items[index+1:]...)
	r.objects = append(r.objects[:index], r.objects[index+1:]...)
	r.tabBar.Objects = append(r.tabBar.Objects[:index], r.tabBar.Objects[index+1:]...)

	t.Refresh()
}

// SetTabLocation sets the location of the tab bar
func (t *TabContainer) SetTabLocation(l TabLocation) {
	t.tabLocation = l
	r := cache.Renderer(t).(*tabContainerRenderer)
	buttons := r.tabBar.Objects
	if fyne.CurrentDevice().IsMobile() || l == TabLocationLeading || l == TabLocationTrailing {
		for _, b := range buttons {
			b.(*tabButton).IconPosition = buttonIconTop
		}
	} else {
		for _, b := range buttons {
			b.(*tabButton).IconPosition = buttonIconInline
		}
	}
	r.tabBar = t.buildTabBar(buttons)

	r.Layout(t.size)
}

func (t *TabContainer) mismatchedContent() bool {
	var hasText, hasIcon bool
	for _, tab := range t.Items {
		hasText = hasText || tab.Text != ""
		hasIcon = hasIcon || tab.Icon != nil
	}

	mismatch := false
	for _, tab := range t.Items {
		if (hasText && tab.Text == "") || (hasIcon && tab.Icon == nil) {
			mismatch = true
			break
		}
	}

	return mismatch
}

// NewTabContainer creates a new tab bar widget that allows the user to choose between different visible containers
func NewTabContainer(items ...*TabItem) *TabContainer {
	tabs := &TabContainer{BaseWidget: BaseWidget{}, Items: items}
	tabs.ExtendBaseWidget(tabs)

	if tabs.mismatchedContent() {
		internal.LogHint("TabContainer items should all have the same type of content (text, icons or both)")
	}

	return tabs
}

type tabContainerRenderer struct {
	tabBar *fyne.Container
	line   *canvas.Rectangle

	// objects holds only the CanvasObject of the tabs' content
	objects   []fyne.CanvasObject
	container *TabContainer
}

func (t *tabContainerRenderer) MinSize() fyne.Size {
	buttonsMin := t.tabBar.MinSize()

	childMin := fyne.NewSize(0, 0)
	for _, child := range t.container.Items {
		childMin = childMin.Union(child.Content.MinSize())
	}

	tabLocation := t.container.tabLocation
	if fyne.CurrentDevice().IsMobile() {
		tabLocation = TabLocationBottom
	}
	switch tabLocation {
	case TabLocationLeading, TabLocationTrailing:
		return fyne.NewSize(buttonsMin.Width+childMin.Width+theme.Padding(),
			fyne.Max(buttonsMin.Height, childMin.Height))
	default:
		return fyne.NewSize(fyne.Max(buttonsMin.Width, childMin.Width),
			buttonsMin.Height+childMin.Height+theme.Padding())
	}
}

func (t *tabContainerRenderer) Layout(size fyne.Size) {
	tabLocation := t.container.tabLocation
	if fyne.CurrentDevice().IsMobile() && (tabLocation == TabLocationLeading || tabLocation == TabLocationTrailing) {
		tabLocation = TabLocationBottom
	}

	switch tabLocation {
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
	case TabLocationLeading:
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
	case TabLocationTrailing:
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

func (t *tabContainerRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (t *tabContainerRenderer) Objects() []fyne.CanvasObject {
	return append(t.objects, t.tabBar, t.line)
}

func (t *tabContainerRenderer) Refresh() {
	t.line.FillColor = theme.ButtonColor()
	t.line.Refresh()

	for i, child := range t.container.Items {
		old := t.objects[i]

		if old == child.Content {
			continue
		}

		old.Hide()
		t.objects[i] = child.Content
		if i == t.container.current {
			child.Content.Show()
		} else {
			child.Content.Hide()
		}
	}

	for i, button := range t.tabBar.Objects {
		if i == t.container.current {
			button.(*tabButton).Style = PrimaryButton
		} else {
			button.(*tabButton).Style = DefaultButton
		}

		button.Refresh()
	}
	canvas.Refresh(t.container)
}

func (t *tabContainerRenderer) Destroy() {
}

type buttonIconPosition int

const (
	buttonIconInline buttonIconPosition = iota
	buttonIconTop
)

var _ fyne.Widget = (*tabButton)(nil)
var _ fyne.Tappable = (*tabButton)(nil)
var _ desktop.Hoverable = (*tabButton)(nil)

type tabButton struct {
	BaseWidget
	hovered      bool
	Icon         fyne.Resource
	IconPosition buttonIconPosition
	OnTap        func()
	Style        ButtonStyle
	Text         string
}

func (b *tabButton) MinSize() fyne.Size {
	b.ExtendBaseWidget(b)
	return b.BaseWidget.MinSize()
}

func (b *tabButton) CreateRenderer() fyne.WidgetRenderer {
	b.ExtendBaseWidget(b)
	var icon *canvas.Image
	if b.Icon != nil {
		icon = canvas.NewImageFromResource(b.Icon)
	}

	label := canvas.NewText(b.Text, theme.TextColor())
	label.Alignment = fyne.TextAlignCenter

	objects := []fyne.CanvasObject{label}
	if icon != nil {
		objects = append(objects, icon)
	}

	return &tabButtonRenderer{
		button:  b,
		icon:    icon,
		label:   label,
		objects: objects,
	}
}

func (b *tabButton) MouseIn(e *desktop.MouseEvent) {
	b.hovered = true
	canvas.Refresh(b)
}

func (b *tabButton) MouseMoved(e *desktop.MouseEvent) {
}

func (b *tabButton) MouseOut() {
	b.hovered = false
	canvas.Refresh(b)
}

func (b *tabButton) Tapped(e *fyne.PointEvent) {
	b.OnTap()
}

func (b *tabButton) TappedSecondary(e *fyne.PointEvent) {
}

type tabButtonRenderer struct {
	button  *tabButton
	icon    *canvas.Image
	label   *canvas.Text
	objects []fyne.CanvasObject
}

func (r *tabButtonRenderer) BackgroundColor() color.Color {
	switch {
	case r.button.Style == PrimaryButton:
		return theme.PrimaryColor()
	case r.button.hovered:
		return theme.HoverColor()
	default:
		return theme.BackgroundColor()
	}
}

func (r *tabButtonRenderer) Destroy() {
}

func (r *tabButtonRenderer) Layout(size fyne.Size) {
	padding := r.padding()
	innerSize := size.Subtract(padding)
	innerOffset := fyne.NewPos(padding.Width/2, padding.Height/2)
	labelShift := 0
	if r.icon != nil {
		var iconOffset fyne.Position
		if r.button.IconPosition == buttonIconTop {
			iconOffset = fyne.NewPos((innerSize.Width-r.iconSize())/2, 0)
		} else {
			iconOffset = fyne.NewPos(0, (innerSize.Height-r.iconSize())/2)
		}
		r.icon.Resize(fyne.NewSize(r.iconSize(), r.iconSize()))
		r.icon.Move(innerOffset.Add(iconOffset))
		labelShift = r.iconSize() + theme.Padding()
	}
	if r.label.Text != "" {
		var labelOffset fyne.Position
		var labelSize fyne.Size
		if r.button.IconPosition == buttonIconTop {
			labelOffset = fyne.NewPos(0, labelShift)
			labelSize = fyne.NewSize(innerSize.Width, r.label.MinSize().Height)
		} else {
			labelOffset = fyne.NewPos(labelShift, 0)
			labelSize = fyne.NewSize(innerSize.Width-labelShift, innerSize.Height)
		}
		r.label.Resize(labelSize)
		r.label.Move(innerOffset.Add(labelOffset))
	}
}

func (r *tabButtonRenderer) MinSize() fyne.Size {
	var contentWidth, contentHeight int
	textSize := r.label.MinSize()
	if r.button.IconPosition == buttonIconTop {
		contentWidth = fyne.Max(textSize.Width, r.iconSize())
		if r.icon != nil {
			contentHeight += r.iconSize()
		}
		if r.label.Text != "" {
			if r.icon != nil {
				contentHeight += theme.Padding()
			}
			contentHeight += textSize.Height
		}
	} else {
		contentHeight = fyne.Max(textSize.Height, r.iconSize())
		if r.icon != nil {
			contentWidth += r.iconSize()
		}
		if r.label.Text != "" {
			if r.icon != nil {
				contentWidth += theme.Padding()
			}
			contentWidth += textSize.Width
		}
	}
	return fyne.NewSize(contentWidth, contentHeight).Add(r.padding())
}

func (r *tabButtonRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *tabButtonRenderer) Refresh() {
	r.label.Color = theme.TextColor()
	r.label.TextSize = theme.TextSize()

	canvas.Refresh(r.button)
}

func (r *tabButtonRenderer) padding() fyne.Size {
	if r.label.Text != "" && r.button.IconPosition == buttonIconInline {
		return fyne.NewSize(theme.Padding()*4, theme.Padding()*2)
	}
	return fyne.NewSize(theme.Padding()*2, theme.Padding()*2)
}

func (r *tabButtonRenderer) iconSize() int {
	switch r.button.IconPosition {
	case buttonIconTop:
		return 2 * theme.IconInlineSize()
	default:
		return theme.IconInlineSize()
	}
}
