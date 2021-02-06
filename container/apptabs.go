package container

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// AppTabs container is used to split your application into various different areas identified by tabs.
// The tabs contain text and/or an icon and allow the user to switch between the content specified in each TabItem.
// Each item is represented by a button at the edge of the container.
//
// Since: 1.4
type AppTabs struct {
	widget.BaseWidget

	Items       []*TabItem
	OnChanged   func(tab *TabItem)
	current     int
	tabLocation TabLocation
}

// TabItem represents a single view in a AppTabs.
// The Text and Icon are used for the tab button and the Content is shown when the corresponding tab is active.
//
// Since: 1.4
type TabItem struct {
	Text    string
	Icon    fyne.Resource
	Content fyne.CanvasObject
}

// TabLocation is the location where the tabs of a tab container should be rendered
//
// Since: 1.4
type TabLocation int

// TabLocation values
const (
	TabLocationTop TabLocation = iota
	TabLocationLeading
	TabLocationBottom
	TabLocationTrailing
)

// NewAppTabs creates a new tab container that allows the user to choose between different areas of an app.
//
// Since: 1.4
func NewAppTabs(items ...*TabItem) *AppTabs {
	tabs := &AppTabs{BaseWidget: widget.BaseWidget{}, Items: items, current: -1}
	if len(items) > 0 {
		// Current is first tab item
		tabs.current = 0
	}
	tabs.BaseWidget.ExtendBaseWidget(tabs)

	if tabs.mismatchedContent() {
		internal.LogHint("AppTabs items should all have the same type of content (text, icons or both)")
	}

	return tabs
}

// NewTabItem creates a new item for a tabbed widget - each item specifies the content and a label for its tab.
//
// Since: 1.4
func NewTabItem(text string, content fyne.CanvasObject) *TabItem {
	return &TabItem{Text: text, Content: content}
}

// NewTabItemWithIcon creates a new item for a tabbed widget - each item specifies the content and a label with an icon for its tab.
//
// Since: 1.4
func NewTabItemWithIcon(text string, icon fyne.Resource, content fyne.CanvasObject) *TabItem {
	return &TabItem{Text: text, Icon: icon, Content: content}
}

// Append adds a new TabItem to the rightmost side of the tab panel
func (c *AppTabs) Append(item *TabItem) {
	c.SetItems(append(c.Items, item))
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (c *AppTabs) CreateRenderer() fyne.WidgetRenderer {
	c.BaseWidget.ExtendBaseWidget(c)
	r := &appTabsRenderer{line: canvas.NewRectangle(theme.ShadowColor()),
		underline: canvas.NewRectangle(theme.PrimaryColor()), container: c}
	r.updateTabs()
	return r
}

// CurrentTab returns the currently selected TabItem.
func (c *AppTabs) CurrentTab() *TabItem {
	if c.current < 0 || c.current >= len(c.Items) {
		return nil
	}
	return c.Items[c.current]
}

// CurrentTabIndex returns the index of the currently selected TabItem.
func (c *AppTabs) CurrentTabIndex() int {
	return c.current
}

// ExtendBaseWidget is used by an extending widget to make use of BaseWidget functionality.
//
// Deprecated: Support for extending containers is being removed
func (c *AppTabs) ExtendBaseWidget(wid fyne.Widget) {
	c.BaseWidget.ExtendBaseWidget(wid)
}

// MinSize returns the size that this widget should not shrink below
func (c *AppTabs) MinSize() fyne.Size {
	c.BaseWidget.ExtendBaseWidget(c)
	return c.BaseWidget.MinSize()
}

// Remove tab by value
func (c *AppTabs) Remove(item *TabItem) {
	for index, existingItem := range c.Items {
		if existingItem == item {
			c.RemoveIndex(index)
			break
		}
	}
}

// RemoveIndex removes tab by index
func (c *AppTabs) RemoveIndex(index int) {
	c.SetItems(append(c.Items[:index], c.Items[index+1:]...))
}

// SetItems sets the container’s items and refreshes.
func (c *AppTabs) SetItems(items []*TabItem) {
	c.Items = items
	if l := len(c.Items); c.current >= l {
		c.current = l - 1
	}
	c.Refresh()
}

// SelectTab sets the specified TabItem to be selected and its content visible.
func (c *AppTabs) SelectTab(item *TabItem) {
	for i, child := range c.Items {
		if child == item {
			c.SelectTabIndex(i)
			return
		}
	}
}

// SelectTabIndex sets the TabItem at the specific index to be selected and its content visible.
func (c *AppTabs) SelectTabIndex(index int) {
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
func (c *AppTabs) SetTabLocation(l TabLocation) {
	c.tabLocation = l
	c.Refresh()
}

// Show this widget, if it was previously hidden
func (c *AppTabs) Show() {
	c.BaseWidget.Show()
	c.SelectTabIndex(c.current)
	c.Refresh()
}

func (c *AppTabs) mismatchedContent() bool {
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

type appTabsRenderer struct {
	animation       *fyne.Animation
	container       *AppTabs
	tabLoc          TabLocation
	line, underline *canvas.Rectangle
	objects         []fyne.CanvasObject // holds only the CanvasObject of the tabs' content
	tabBar          *fyne.Container
}

func (r *appTabsRenderer) Destroy() {
}

func (r *appTabsRenderer) Layout(size fyne.Size) {
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

func (r *appTabsRenderer) MinSize() fyne.Size {
	buttonsMin := r.tabBar.MinSize()

	childMin := fyne.NewSize(0, 0)
	for _, child := range r.container.Items {
		childMin = childMin.Max(child.Content.MinSize())
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

func (r *appTabsRenderer) Objects() []fyne.CanvasObject {
	return append(r.objects, r.tabBar, r.line, r.underline)
}

func (r *appTabsRenderer) Refresh() {
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
				button.(*tabButton).Importance = widget.HighImportance
			} else {
				button.(*tabButton).Importance = widget.MediumImportance
			}

			button.Refresh()
		}
	}
	r.moveSelection()
	canvas.Refresh(r.container)
}

func (r *appTabsRenderer) adaptedLocation() TabLocation {
	tabLocation := r.container.tabLocation
	if fyne.CurrentDevice().IsMobile() && (tabLocation == TabLocationLeading || tabLocation == TabLocationTrailing) {
		return TabLocationBottom
	}

	return r.container.tabLocation
}

func (r *appTabsRenderer) buildButton(item *TabItem, iconPos buttonIconPosition) *tabButton {
	return &tabButton{
		Text:         item.Text,
		Icon:         item.Icon,
		IconPosition: iconPos,
		OnTap:        func() { r.container.SelectTab(item) },
	}
}

func (r *appTabsRenderer) buildTabBar(buttons []fyne.CanvasObject) *fyne.Container {
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

	return fyne.NewContainerWithLayout(lay, buttons...)
}

func (r *appTabsRenderer) moveSelection() {
	if r.container.current < 0 {
		r.underline.Hide()
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

	r.underline.Show()
	if r.underline.Position().IsZero() {
		r.underline.Move(underlinePos)
		r.underline.Resize(underlineSize)
		return
	}

	if r.animation != nil {
		r.animation.Stop()
	}
	r.animation = canvas.NewPositionAnimation(r.underline.Position(), underlinePos, canvas.DurationShort, func(p fyne.Position) {
		r.underline.Move(p)
		canvas.Refresh(r.underline)
		if p == underlinePos {
			r.animation = nil
		}
	})
	r.animation.Start()

	canvas.NewSizeAnimation(r.underline.Size(), underlineSize, canvas.DurationShort, func(s fyne.Size) {
		r.underline.Resize(s)
		canvas.Refresh(r.underline)
	}).Start()
}

func (r *appTabsRenderer) tabsInSync() bool {
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

func (r *appTabsRenderer) updateTabs() bool {
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

	length := len(r.container.Items)
	buttons := make([]fyne.CanvasObject, length)
	objects := make([]fyne.CanvasObject, length)
	for i, item := range r.container.Items {
		button := r.buildButton(item, iconPos)
		if i == r.container.current {
			button.Importance = widget.HighImportance
			item.Content.Show()
		} else {
			item.Content.Hide()
		}
		buttons[i] = button
		objects[i] = item.Content
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
	widget.BaseWidget
	hovered      bool
	Icon         fyne.Resource
	IconPosition buttonIconPosition
	Importance   widget.ButtonImportance
	OnTap        func()
	Text         string
}

func (b *tabButton) CreateRenderer() fyne.WidgetRenderer {
	b.ExtendBaseWidget(b)
	background := canvas.NewRectangle(theme.HoverColor())
	background.Hide()
	var icon *canvas.Image
	if b.Icon != nil {
		icon = canvas.NewImageFromResource(b.Icon)
	}

	label := canvas.NewText(b.Text, theme.ForegroundColor())
	label.TextStyle.Bold = true
	label.Alignment = fyne.TextAlignCenter

	objects := []fyne.CanvasObject{background, label}
	if icon != nil {
		objects = append(objects, icon)
	}

	r := &tabButtonRenderer{
		button:     b,
		background: background,
		icon:       icon,
		label:      label,
		objects:    objects,
	}
	r.Refresh()
	return r
}

func (b *tabButton) MinSize() fyne.Size {
	b.ExtendBaseWidget(b)
	return b.BaseWidget.MinSize()
}

func (b *tabButton) MouseIn(e *desktop.MouseEvent) {
	b.hovered = true
	b.Refresh()
}

func (b *tabButton) MouseMoved(e *desktop.MouseEvent) {
}

func (b *tabButton) MouseOut() {
	b.hovered = false
	b.Refresh()
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
	button     *tabButton
	background *canvas.Rectangle
	icon       *canvas.Image
	label      *canvas.Text
	objects    []fyne.CanvasObject
}

func (r *tabButtonRenderer) Destroy() {
}

func (r *tabButtonRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)
	padding := r.padding()
	innerSize := size.Subtract(padding)
	innerOffset := fyne.NewPos(padding.Width/2, padding.Height/2)
	labelShift := float32(0)
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
	var contentWidth, contentHeight float32
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
	if r.button.hovered {
		r.background.FillColor = theme.HoverColor()
		r.background.Show()
	} else {
		r.background.Hide()
	}
	r.background.Refresh()

	r.label.Text = r.button.Text
	if r.button.Importance == widget.HighImportance {
		r.label.Color = theme.PrimaryColor()
	} else {
		r.label.Color = theme.ForegroundColor()
	}
	r.label.TextSize = theme.TextSize()
	if r.button.Text == "" {
		r.label.Hide()
	} else {
		r.label.Show()
	}

	if r.icon != nil && r.icon.Resource != nil {
		switch res := r.icon.Resource.(type) {
		case *theme.ThemedResource:
			if r.button.Importance == widget.HighImportance {
				r.icon.Resource = theme.NewPrimaryThemedResource(res)
				r.icon.Refresh()
			}
		case *theme.PrimaryThemedResource:
			if r.button.Importance != widget.HighImportance {
				r.icon.Resource = res.Original()
				r.icon.Refresh()
			}
		}
	}

	canvas.Refresh(r.button)
}

func (r *tabButtonRenderer) iconSize() float32 {
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
