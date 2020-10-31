package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
)

// TabContainer widget allows switching visible content from a list of TabItems.
// Each item is represented by a button at the top of the widget.
//
// Deprecated: use container.Tabs instead.
type TabContainer struct {
	BaseWidget

	Items       []*TabItem
	OnChanged   func(tab *TabItem)
	current     int
	tabLocation TabLocation
}

// TabItem represents a single view in a TabContainer.
// The Text and Icon are used for the tab button and the Content is shown when the corresponding tab is active.
//
// Deprecated: use container.TabItem instead.
type TabItem struct {
	Text    string
	Icon    fyne.Resource
	Content fyne.CanvasObject
}

// TabLocation is the location where the tabs of a tab container should be rendered.
//
// Deprecated: use container.TabLocation instead.
type TabLocation int

// TabLocation values
const (
	// Deprecated: use container.TabLocationTop
	TabLocationTop TabLocation = iota
	// Deprecated: use container.TabLocationLeading
	TabLocationLeading
	// Deprecated: use container.TabLocationBottom
	TabLocationBottom
	// Deprecated: use container.TabLocationTrailing
	TabLocationTrailing
)

// NewTabContainer creates a new tab bar widget that allows the user to choose between different visible containers
func NewTabContainer(items ...*TabItem) *TabContainer {
	tabs := &TabContainer{BaseWidget: BaseWidget{}, Items: items, current: -1}
	if len(items) > 0 {
		// Current is first tab item
		tabs.current = 0
	}
	tabs.ExtendBaseWidget(tabs)

	if tabs.mismatchedContent() {
		internal.LogHint("TabContainer items should all have the same type of content (text, icons or both)")
	}

	return tabs
}

// NewTabItem creates a new item for a tabbed widget - each item specifies the content and a label for its tab.
func NewTabItem(text string, content fyne.CanvasObject) *TabItem {
	return &TabItem{Text: text, Content: content}
}

// NewTabItemWithIcon creates a new item for a tabbed widget - each item specifies the content and a label with an icon for its tab.
func NewTabItemWithIcon(text string, icon fyne.Resource, content fyne.CanvasObject) *TabItem {
	return &TabItem{Text: text, Icon: icon, Content: content}
}

// Append adds a new TabItem to the rightmost side of the tab panel
func (c *TabContainer) Append(item *TabItem) {
	c.SetItems(append(c.Items, item))
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (c *TabContainer) CreateRenderer() fyne.WidgetRenderer {
	c.ExtendBaseWidget(c)
	r := &tabContainerRenderer{line: canvas.NewRectangle(theme.ShadowColor()),
		underline: canvas.NewRectangle(theme.PrimaryColor()), container: c}
	r.updateTabs()
	return r
}

// CurrentTab returns the currently selected TabItem.
func (c *TabContainer) CurrentTab() *TabItem {
	if c.current < 0 || c.current >= len(c.Items) {
		return nil
	}
	return c.Items[c.current]
}

// CurrentTabIndex returns the index of the currently selected TabItem.
func (c *TabContainer) CurrentTabIndex() int {
	return c.current
}

// MinSize returns the size that this widget should not shrink below
func (c *TabContainer) MinSize() fyne.Size {
	c.ExtendBaseWidget(c)
	return c.BaseWidget.MinSize()
}

// Remove tab by value
func (c *TabContainer) Remove(item *TabItem) {
	for index, existingItem := range c.Items {
		if existingItem == item {
			c.RemoveIndex(index)
			break
		}
	}
}

// RemoveIndex removes tab by index
func (c *TabContainer) RemoveIndex(index int) {
	c.SetItems(append(c.Items[:index], c.Items[index+1:]...))
}

// SetItems sets the container’s items and refreshes.
func (c *TabContainer) SetItems(items []*TabItem) {
	c.Items = items
	c.Refresh()
}

// SelectTab sets the specified TabItem to be selected and its content visible.
func (c *TabContainer) SelectTab(item *TabItem) {
	for i, child := range c.Items {
		if child == item {
			c.SelectTabIndex(i)
			return
		}
	}
}

// SelectTabIndex sets the TabItem at the specific index to be selected and its content visible.
func (c *TabContainer) SelectTabIndex(index int) {
	if index < 0 || index >= len(c.Items) || c.current == index {
		return
	}

	c.current = index
	c.Refresh()

	if c.OnChanged != nil {
		c.OnChanged(c.Items[c.current])
	}
}

// SetTabLocation sets the location of the tab bar
func (c *TabContainer) SetTabLocation(l TabLocation) {
	c.tabLocation = l
	c.Refresh()
}

// Show this widget, if it was previously hidden
func (c *TabContainer) Show() {
	c.BaseWidget.Show()
	c.SelectTabIndex(c.current)
	c.Refresh()
}

func (c *TabContainer) mismatchedContent() bool {
	var hasText, hasIcon bool
	for _, tab := range c.Items {
		hasText = hasText || tab.Text != ""
		hasIcon = hasIcon || tab.Icon != nil
	}

	mismatch := false
	for _, tab := range c.Items {
		if (hasText && tab.Text == "") || (hasIcon && tab.Icon == nil) {
			mismatch = true
			break
		}
	}

	return mismatch
}

type tabContainerRenderer struct {
	container       *TabContainer
	tabLoc          TabLocation
	line, underline *canvas.Rectangle
	objects         []fyne.CanvasObject // holds only the CanvasObject of the tabs' content
	tabBar          *fyne.Container
}

func (r *tabContainerRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *tabContainerRenderer) Destroy() {
}

func (r *tabContainerRenderer) Layout(size fyne.Size) {
	tabBarMinSize := r.tabBar.MinSize()
	var tabBarPos fyne.Position
	var tabBarSize fyne.Size
	var linePos fyne.Position
	var lineSize fyne.Size
	var childPos fyne.Position
	var childSize fyne.Size
	switch r.adaptedLocation() {
	case TabLocationTop:
		buttonHeight := tabBarMinSize.Height
		tabBarPos = fyne.NewPos(0, 0)
		tabBarSize = fyne.NewSize(size.Width, buttonHeight)
		linePos = fyne.NewPos(0, buttonHeight)
		lineSize = fyne.NewSize(size.Width, theme.Padding())
		barHeight := buttonHeight + theme.Padding()
		childPos = fyne.NewPos(0, barHeight)
		childSize = fyne.NewSize(size.Width, size.Height-barHeight)
	case TabLocationLeading:
		buttonWidth := tabBarMinSize.Width
		tabBarPos = fyne.NewPos(0, 0)
		tabBarSize = fyne.NewSize(buttonWidth, size.Height)
		linePos = fyne.NewPos(buttonWidth, 0)
		lineSize = fyne.NewSize(theme.Padding(), size.Height)
		barWidth := buttonWidth + theme.Padding()
		childPos = fyne.NewPos(barWidth, 0)
		childSize = fyne.NewSize(size.Width-barWidth, size.Height)
	case TabLocationBottom:
		buttonHeight := tabBarMinSize.Height
		tabBarPos = fyne.NewPos(0, size.Height-buttonHeight)
		tabBarSize = fyne.NewSize(size.Width, buttonHeight)
		barHeight := buttonHeight + theme.Padding()
		linePos = fyne.NewPos(0, size.Height-barHeight)
		lineSize = fyne.NewSize(size.Width, theme.Padding())
		childPos = fyne.NewPos(0, 0)
		childSize = fyne.NewSize(size.Width, size.Height-barHeight)
	case TabLocationTrailing:
		buttonWidth := tabBarMinSize.Width
		tabBarPos = fyne.NewPos(size.Width-buttonWidth, 0)
		tabBarSize = fyne.NewSize(buttonWidth, size.Height)
		barWidth := buttonWidth + theme.Padding()
		linePos = fyne.NewPos(size.Width-barWidth, 0)
		lineSize = fyne.NewSize(theme.Padding(), size.Height)
		childPos = fyne.NewPos(0, 0)
		childSize = fyne.NewSize(size.Width-barWidth, size.Height)
	}

	r.tabBar.Move(tabBarPos)
	r.tabBar.Resize(tabBarSize)
	r.line.Move(linePos)
	r.line.Resize(lineSize)
	if r.container.current >= 0 && r.container.current < len(r.container.Items) {
		child := r.container.Items[r.container.current].Content
		child.Move(childPos)
		child.Resize(childSize)
	}
	r.moveSelection()
}

func (r *tabContainerRenderer) MinSize() fyne.Size {
	buttonsMin := r.tabBar.MinSize()

	childMin := fyne.NewSize(0, 0)
	for _, child := range r.container.Items {
		childMin = childMin.Union(child.Content.MinSize())
	}

	tabLocation := r.container.tabLocation
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

func (r *tabContainerRenderer) Objects() []fyne.CanvasObject {
	return append(r.objects, r.tabBar, r.line, r.underline)
}

func (r *tabContainerRenderer) Refresh() {
	r.line.FillColor = theme.ShadowColor()
	r.line.Refresh()
	r.underline.FillColor = theme.PrimaryColor()

	if r.updateTabs() {
		r.Layout(r.container.Size())
	} else {
		current := r.container.current
		if current >= 0 && current < len(r.objects) && !r.objects[current].Visible() {
			r.Layout(r.container.Size())
			for i, o := range r.objects {
				if i == current {
					o.Show()
				} else {
					o.Hide()
				}
			}
		}
		for i, button := range r.tabBar.Objects {
			if i == current {
				button.(*tabButton).Importance = HighImportance
			} else {
				button.(*tabButton).Importance = MediumImportance
			}

			button.Refresh()
		}
	}
	r.moveSelection()
	canvas.Refresh(r.container)
}

func (r *tabContainerRenderer) adaptedLocation() TabLocation {
	tabLocation := r.container.tabLocation
	if fyne.CurrentDevice().IsMobile() && (tabLocation == TabLocationLeading || tabLocation == TabLocationTrailing) {
		return TabLocationBottom
	}

	return r.container.tabLocation
}

func (r *tabContainerRenderer) buildButton(item *TabItem, iconPos buttonIconPosition) *tabButton {
	return &tabButton{
		Text:         item.Text,
		Icon:         item.Icon,
		IconPosition: iconPos,
		OnTap:        func() { r.container.SelectTab(item) },
	}
}

func (r *tabContainerRenderer) buildTabBar(buttons []fyne.CanvasObject) *fyne.Container {
	var lay fyne.Layout
	if fyne.CurrentDevice().IsMobile() {
		cells := len(buttons)
		if cells == 0 {
			cells = 1
		}
		lay = layout.NewGridLayout(cells)
	} else if r.container.tabLocation == TabLocationLeading || r.container.tabLocation == TabLocationTrailing {
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

func (r *tabContainerRenderer) moveSelection() {
	if r.container.current < 0 {
		return
	}
	selected := r.tabBar.Objects[r.container.current]

	var underlinePos fyne.Position
	var underlineSize fyne.Size
	switch r.adaptedLocation() {
	case TabLocationTop:
		underlinePos = fyne.NewPos(selected.Position().X, r.tabBar.MinSize().Height)
		underlineSize = fyne.NewSize(selected.Size().Width, theme.Padding())
	case TabLocationLeading:
		underlinePos = fyne.NewPos(r.tabBar.MinSize().Width, selected.Position().Y)
		underlineSize = fyne.NewSize(theme.Padding(), selected.Size().Height)
	case TabLocationBottom:
		underlinePos = fyne.NewPos(selected.Position().X, r.tabBar.Position().Y-theme.Padding())
		underlineSize = fyne.NewSize(selected.Size().Width, theme.Padding())
	case TabLocationTrailing:
		underlinePos = fyne.NewPos(r.tabBar.Position().X-theme.Padding(), selected.Position().Y)
		underlineSize = fyne.NewSize(theme.Padding(), selected.Size().Height)
	}
	r.underline.Resize(underlineSize)
	r.underline.Move(underlinePos)
}

func (r *tabContainerRenderer) tabsInSync() bool {
	if r.tabBar == nil {
		return false
	}
	if r.tabLoc != r.container.tabLocation {
		return false
	}
	if len(r.objects) != len(r.container.Items) {
		return false
	}
	if len(r.tabBar.Objects) != len(r.container.Items) {
		return false
	}
	for i, item := range r.container.Items {
		if item.Content != r.objects[i] {
			return false
		}
		button := r.tabBar.Objects[i].(*tabButton)
		if item.Text != button.Text {
			return false
		}
		if item.Icon != button.Icon {
			return false
		}
	}
	return true
}

func (r *tabContainerRenderer) updateTabs() bool {
	if r.tabsInSync() {
		return false
	}

	r.tabLoc = r.container.tabLocation
	var iconPos buttonIconPosition
	if fyne.CurrentDevice().IsMobile() || r.container.tabLocation == TabLocationLeading || r.container.tabLocation == TabLocationTrailing {
		iconPos = buttonIconTop
	} else {
		iconPos = buttonIconInline
	}
	var buttons, objects []fyne.CanvasObject
	for i, item := range r.container.Items {
		button := r.buildButton(item, iconPos)
		if i == r.container.current {
			button.Importance = HighImportance
			item.Content.Show()
		} else {
			item.Content.Hide()
		}
		buttons = append(buttons, button)
		objects = append(objects, item.Content)
	}
	r.tabBar = r.buildTabBar(buttons)
	r.objects = objects
	r.moveSelection()
	return true
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
	Importance   ButtonImportance
	OnTap        func()
	Text         string
}

func (b *tabButton) CreateRenderer() fyne.WidgetRenderer {
	b.ExtendBaseWidget(b)
	var icon *canvas.Image
	if b.Icon != nil {
		icon = canvas.NewImageFromResource(b.Icon)
		if b.Importance == HighImportance {
			icon.Resource = theme.NewPrimaryThemedResource(b.Icon)
		}
	}

	label := canvas.NewText(b.Text, theme.TextColor())
	if b.Importance == HighImportance {
		label.Color = theme.PrimaryColor()
	}
	label.TextStyle.Bold = true
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

func (b *tabButton) MinSize() fyne.Size {
	b.ExtendBaseWidget(b)
	return b.BaseWidget.MinSize()
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

func (b *tabButton) setText(text string) {
	if text == b.Text {
		return
	}

	b.Text = text
	b.Refresh()
}

type tabButtonRenderer struct {
	button  *tabButton
	icon    *canvas.Image
	label   *canvas.Text
	objects []fyne.CanvasObject
}

func (r *tabButtonRenderer) BackgroundColor() color.Color {
	switch {
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
	r.label.Text = r.button.Text
	if r.button.Importance == HighImportance {
		r.label.Color = theme.PrimaryColor()
	} else {
		r.label.Color = theme.TextColor()
	}
	r.label.TextSize = theme.TextSize()

	if r.icon != nil && r.icon.Resource != nil {
		switch res := r.icon.Resource.(type) {
		case *theme.ThemedResource:
			if r.button.Importance == HighImportance {
				r.icon.Resource = theme.NewPrimaryThemedResource(res)
				r.icon.Refresh()
			}
		case *theme.PrimaryThemedResource:
			if r.button.Importance != HighImportance {
				r.icon.Resource = res.Original()
				r.icon.Refresh()
			}
		}
	}

	canvas.Refresh(r.button)
}

func (r *tabButtonRenderer) iconSize() int {
	switch r.button.IconPosition {
	case buttonIconTop:
		return 2 * theme.IconInlineSize()
	default:
		return theme.IconInlineSize()
	}
}

func (r *tabButtonRenderer) padding() fyne.Size {
	if r.label.Text != "" && r.button.IconPosition == buttonIconInline {
		return fyne.NewSize(theme.Padding()*4, theme.Padding()*2)
	}
	return fyne.NewSize(theme.Padding()*2, theme.Padding()*2)
}
