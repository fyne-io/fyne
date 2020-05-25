package widget

import (
	"fmt"
	"image/color"
	"log"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
)

// TabContainer widget allows switching visible content from a list of TabItems.
// Each item is represented by a button at the top of the widget.
type TabContainer struct {
	BaseWidget

	Items            []*TabItem
	OnChanged        func(tab *TabItem)
	current          int
	tabLocation      TabLocation
	overflowExpanded bool
}

// TabItem represents a single view in a TabContainer.
// The Text and Icon are used for the tab button and the Content is shown when the corresponding tab is active.
type TabItem struct {
	Text    string
	Icon    fyne.Resource
	Content fyne.CanvasObject
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
	var items []*TabItem
	c.propertyLock.RLock()
	items = append(items, c.Items...)
	c.propertyLock.RUnlock()
	c.SetItems(append(items, item))
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
// Implements: fyne.Widget
func (c *TabContainer) CreateRenderer() fyne.WidgetRenderer {
	c.ExtendBaseWidget(c)
	r := &tabContainerRenderer{
		line:      canvas.NewRectangle(theme.ButtonColor()),
		container: c,
	}
	r.updateTabs()
	return r
}

// CurrentTab returns the currently selected TabItem.
func (c *TabContainer) CurrentTab() *TabItem {
	if c.current < 0 || c.current >= len(c.Items) {
		return nil
	}
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	return c.Items[c.current]
}

// CurrentTabIndex returns the index of the currently selected TabItem.
func (c *TabContainer) CurrentTabIndex() int {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	return c.current
}

// MinSize returns the size that this widget should not shrink below
// Implements: fyne.Widget
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
	var items []*TabItem
	c.propertyLock.RLock()
	items = append(items, c.Items[:index]...)
	items = append(items, c.Items[index+1:]...)
	c.propertyLock.RUnlock()
	c.SetItems(items)
}

// SetItems sets the container’s items and refreshes.
func (c *TabContainer) SetItems(items []*TabItem) {
	c.propertyLock.Lock()
	c.Items = items
	c.propertyLock.Unlock()
	c.Refresh()
}

// SelectTab sets the specified TabItem to be selected and its content visible.
func (c *TabContainer) SelectTab(item *TabItem) {
	for i, tab := range c.Items {
		if tab == item {
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

	c.propertyLock.Lock()
	c.current = index
	item := c.Items[c.current]
	c.propertyLock.Unlock()

	c.Refresh()

	if c.OnChanged != nil {
		c.OnChanged(item)
	}
}

// SetTabLocation sets the location of the tab bar
func (c *TabContainer) SetTabLocation(l TabLocation) {
	if fyne.CurrentDevice().IsMobile() && (l == TabLocationLeading || l == TabLocationTrailing) {
		l = TabLocationBottom
	}
	c.propertyLock.Lock()
	c.tabLocation = l
	c.propertyLock.Unlock()
	c.Refresh()
}

// Show this widget, if it was previously hidden
func (c *TabContainer) Show() {
	c.BaseWidget.Show()
	c.SelectTabIndex(c.current)
	c.Refresh()
}

// TabLocation gets the location of the tab bar
func (c *TabContainer) TabLocation() TabLocation {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	return c.tabLocation
}

// ToggleOverflow switches whether the overflow popup is displayed.
func (c *TabContainer) ToggleOverflow() {
	c.propertyLock.Lock()
	c.overflowExpanded = !c.overflowExpanded
	log.Println("OverflowExpanded:", c.overflowExpanded)
	c.propertyLock.Unlock()
	c.Refresh()
}

func (c *TabContainer) contentTypes() (hasText, hasIcon bool) {
	for _, tab := range c.Items {
		hasText = hasText || tab.Text != ""
		hasIcon = hasIcon || tab.Icon != nil
	}
	return
}

func (c *TabContainer) mismatchedContent() bool {
	hasText, hasIcon := c.contentTypes()

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
	container   *TabContainer
	tabLocation TabLocation

	overflowButton *tabOverflowButton
	line           *canvas.Rectangle

	tabButtons  []fyne.CanvasObject
	tabContents []fyne.CanvasObject
}

func (r *tabContainerRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *tabContainerRenderer) Destroy() {
}

func (r *tabContainerRenderer) Layout(size fyne.Size) {
	log.Println("tabContainerRenderer.Layout:", size)
	if size.IsZero() {
		return
	}
	var linePos fyne.Position
	var lineSize fyne.Size

	var tabBarPos fyne.Position
	var tabBarSize fyne.Size
	var tabBarSlots int // FIXME this implementation uses the maximum minsize of all tab buttons to calculate how many tabs can fit, however if one tab is big and the rest are small then the overflow button is used even though all tabs might fit in the given space.
	tabButtonMin := r.tabButtonMinSize()
	var tabContentPos fyne.Position
	var tabContentSize fyne.Size

	switch r.tabLocation {
	case TabLocationTop:
		tabBarPos = fyne.NewPos(0, 0)
		tabBarSize = fyne.NewSize(size.Width, tabButtonMin.Height)
		tabBarSlots = int((tabBarSize.Width + theme.Padding()) / (tabButtonMin.Width + theme.Padding()))
		linePos = fyne.NewPos(0, tabButtonMin.Height)
		lineSize = fyne.NewSize(size.Width, theme.Padding())
		barHeight := tabButtonMin.Height + theme.Padding()
		tabContentPos = fyne.NewPos(0, barHeight)
		tabContentSize = fyne.NewSize(size.Width, size.Height-barHeight)
	case TabLocationLeading:
		tabBarPos = fyne.NewPos(0, 0)
		tabBarSize = fyne.NewSize(tabButtonMin.Width, size.Height)
		tabBarSlots = int((tabBarSize.Height + theme.Padding()) / (tabButtonMin.Height + theme.Padding()))
		linePos = fyne.NewPos(tabButtonMin.Width, 0)
		lineSize = fyne.NewSize(theme.Padding(), size.Height)
		barWidth := tabButtonMin.Width + theme.Padding()
		tabContentPos = fyne.NewPos(barWidth, 0)
		tabContentSize = fyne.NewSize(size.Width-barWidth, size.Height)
	case TabLocationBottom:
		tabBarPos = fyne.NewPos(0, size.Height-tabButtonMin.Height)
		tabBarSize = fyne.NewSize(size.Width, tabButtonMin.Height)
		tabBarSlots = int((tabBarSize.Width + theme.Padding()) / (tabButtonMin.Width + theme.Padding()))
		barHeight := tabButtonMin.Height + theme.Padding()
		linePos = fyne.NewPos(0, size.Height-barHeight)
		lineSize = fyne.NewSize(size.Width, theme.Padding())
		tabContentPos = fyne.NewPos(0, 0)
		tabContentSize = fyne.NewSize(size.Width, size.Height-barHeight)
	case TabLocationTrailing:
		tabBarPos = fyne.NewPos(size.Width-tabButtonMin.Width, 0)
		tabBarSize = fyne.NewSize(tabButtonMin.Width, size.Height)
		tabBarSlots = int((tabBarSize.Height + theme.Padding()) / (tabButtonMin.Height + theme.Padding()))
		barWidth := tabButtonMin.Width + theme.Padding()
		linePos = fyne.NewPos(size.Width-barWidth, 0)
		lineSize = fyne.NewSize(theme.Padding(), size.Height)
		tabContentPos = fyne.NewPos(0, 0)
		tabContentSize = fyne.NewSize(size.Width-barWidth, size.Height)
	}
	log.Println("TabBarPos:", tabBarPos)
	log.Println("TabBarSize:", tabBarSize)
	log.Println("TabBarSlots:", tabBarSlots)

	r.line.Move(linePos)
	r.line.Resize(lineSize)

	var tabBarObjects []fyne.CanvasObject
	if tabBarSlots < len(r.tabButtons) {
		if tabBarSlots > 0 {
			index := 0
			for ; index < tabBarSlots-1; index++ {
				o := r.tabButtons[index]
				o.Show()
				tabBarObjects = append(tabBarObjects, o)
			}
			for ; index < len(r.tabButtons); index++ {
				o := r.tabButtons[index]
				o.Hide()
			}
		}
		r.overflowButton.Show()
		tabBarObjects = append(tabBarObjects, r.overflowButton)
	} else {
		for _, o := range r.tabButtons {
			o.Show()
			tabBarObjects = append(tabBarObjects, o)
		}
		r.overflowButton.Hide()
	}
	tabLayout := layout.NewHBoxLayout()
	if fyne.CurrentDevice().IsMobile() {
		tabLayout = layout.NewGridLayout(tabBarSlots)
	} else if r.tabLocation == TabLocationLeading || r.tabLocation == TabLocationTrailing {
		tabLayout = layout.NewVBoxLayout()
	}
	tabBar := fyne.NewContainerWithLayout(tabLayout, tabBarObjects...)
	tabBar.Resize(tabBarSize)
	for _, o := range tabBarObjects {
		o.Move(o.Position().Add(tabBarPos))
	}

	if r.container.current >= 0 && r.container.current < len(r.container.Items) {
		tabContent := r.container.Items[r.container.current].Content
		tabContent.Move(tabContentPos)
		tabContent.Resize(tabContentSize)
	}
}

func (r *tabContainerRenderer) MinSize() (min fyne.Size) {
	tabButtonMin := r.tabButtonMinSize()
	tabContentMin := r.tabContentMinSize()

	switch r.tabLocation {
	case TabLocationLeading, TabLocationTrailing:
		min.Width = tabButtonMin.Width + theme.Padding() + tabContentMin.Width
		min.Height = fyne.Max(tabButtonMin.Height, tabContentMin.Height)
	default:
		min.Width = fyne.Max(tabButtonMin.Width, tabContentMin.Width)
		min.Height = tabButtonMin.Height + theme.Padding() + tabContentMin.Height
	}
	log.Println("tabContainerRenderer.MinSize:", min)
	return
}

func (r *tabContainerRenderer) Objects() (objects []fyne.CanvasObject) {
	objects = append(objects, r.line)
	objects = append(objects, r.overflowButton)
	objects = append(objects, r.tabButtons...)
	objects = append(objects, r.tabContents...)
	return
}

func (r *tabContainerRenderer) Refresh() {
	r.container.propertyLock.RLock()
	r.updateLine()
	r.updateOverflowButton()

	if r.updateTabs() {
		r.Layout(r.container.size)
	} else {
		current := r.container.CurrentTabIndex()
		if current >= 0 && current < len(r.tabContents) && !r.tabContents[current].Visible() {
			for i, tc := range r.tabContents {
				if i == current {
					tc.Show()
				} else {
					tc.Hide()
				}
			}
			r.Layout(r.container.size)
		}
		for i, button := range r.tabButtons {
			if i == current {
				button.(*tabButton).Style = PrimaryButton
			} else {
				button.(*tabButton).Style = DefaultButton
			}

			button.Refresh()
		}
	}
	r.container.propertyLock.RUnlock()
	canvas.Refresh(r.container.super())
}

func (r *tabContainerRenderer) buildTabButton(index int, item *TabItem, iconPos buttonIconPosition) *tabButton {
	return &tabButton{
		Text:         item.Text,
		Icon:         item.Icon,
		IconPosition: iconPos,
		OnTapped:     func() { r.container.SelectTabIndex(index) },
	}
}

func (r *tabContainerRenderer) buildTabOverflowButton(iconPos buttonIconPosition) *tabOverflowButton {
	return &tabOverflowButton{
		tabButton: tabButton{
			IconPosition: iconPos,
			OnTapped:     func() { r.container.ToggleOverflow() },
		},
		container: r.container,
	}
}

func (r *tabContainerRenderer) tabButtonMinSize() (min fyne.Size) {
	min = r.overflowButton.MinSize()
	for _, b := range r.tabButtons {
		m := b.MinSize()
		min = min.Max(m)
	}
	return
}

func (r *tabContainerRenderer) tabContentMinSize() (min fyne.Size) {
	for _, c := range r.tabContents {
		min = min.Max(c.MinSize())
	}
	return
}

func (r *tabContainerRenderer) tabsInSync() bool {
	if r.tabButtons == nil {
		return false
	}
	if len(r.tabButtons) != len(r.container.Items) {
		return false
	}
	if len(r.tabContents) != len(r.container.Items) {
		return false
	}
	if r.tabLocation != r.container.TabLocation() {
		return false
	}
	for i, item := range r.container.Items {
		if item.Content != r.tabContents[i] {
			return false
		}
		button := r.tabButtons[i].(*tabButton)
		if item.Text != button.Text {
			return false
		}
		if item.Icon != button.Icon {
			return false
		}
	}
	return true
}

func (r *tabContainerRenderer) updateLine() {
	r.line.FillColor = theme.ButtonColor()
	r.line.Refresh()
}

func (r *tabContainerRenderer) updateOverflowButton() {
	r.overflowButton.Refresh()
}

func (r *tabContainerRenderer) updateTabs() bool {
	if r.tabsInSync() {
		return false
	}

	r.tabLocation = r.container.TabLocation()
	var iconPos buttonIconPosition
	if fyne.CurrentDevice().IsMobile() || r.tabLocation == TabLocationLeading || r.tabLocation == TabLocationTrailing {
		iconPos = buttonIconTop
	} else {
		iconPos = buttonIconInline
	}
	r.overflowButton = r.buildTabOverflowButton(iconPos)
	var buttons, contents []fyne.CanvasObject
	for i, item := range r.container.Items {
		button := r.buildTabButton(i, item, iconPos)
		if i == r.container.current {
			button.Style = PrimaryButton
			item.Content.Show()
		} else {
			item.Content.Hide()
		}
		buttons = append(buttons, button)
		contents = append(contents, item.Content)
	}
	r.tabButtons = buttons
	r.tabContents = contents
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
	OnTapped     func()
	Style        ButtonStyle
	Text         string
}

func (b *tabButton) CreateRenderer() fyne.WidgetRenderer {
	b.ExtendBaseWidget(b)
	var icon *canvas.Image
	if b.Icon != nil {
		icon = canvas.NewImageFromResource(b.Icon)
	}

	label := &canvas.Text{
		Alignment: fyne.TextAlignCenter,
	}

	objects := []fyne.CanvasObject{label}
	if icon != nil {
		objects = append(objects, icon)
	}

	r := &tabButtonRenderer{
		button:  b,
		icon:    icon,
		label:   label,
		objects: objects,
	}
	r.updateLabel()
	return r
}

func (b *tabButton) MinSize() fyne.Size {
	b.ExtendBaseWidget(b)
	return b.BaseWidget.MinSize()
}

func (b *tabButton) MouseIn(*desktop.MouseEvent) {
	b.hovered = true
	canvas.Refresh(b.super())
}

func (b *tabButton) MouseMoved(*desktop.MouseEvent) {
}

func (b *tabButton) MouseOut() {
	b.hovered = false
	canvas.Refresh(b.super())
}

func (b *tabButton) Tapped(*fyne.PointEvent) {
	b.OnTapped()
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
		iconSize := r.iconSize()
		var iconOffset fyne.Position
		if r.button.IconPosition == buttonIconTop {
			iconOffset = fyne.NewPos((innerSize.Width-iconSize)/2, 0)
		} else {
			iconOffset = fyne.NewPos(0, (innerSize.Height-iconSize)/2)
		}
		r.icon.Resize(fyne.NewSize(iconSize, iconSize))
		r.icon.Move(innerOffset.Add(iconOffset))
		labelShift = iconSize + theme.Padding()
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

func (r *tabButtonRenderer) MinSize() (min fyne.Size) {
	iconMinSize := r.iconSize()
	labelMinSize := r.label.MinSize()
	if r.button.IconPosition == buttonIconTop {
		if r.icon != nil {
			min.Width = fyne.Max(min.Width, iconMinSize)
			min.Height += iconMinSize
		}
		if r.label.Text != "" {
			min.Width = fyne.Max(min.Width, labelMinSize.Width)
			if r.icon != nil {
				min.Height += theme.Padding()
			}
			min.Height += labelMinSize.Height
		}
	} else {
		if r.icon != nil {
			min.Width += iconMinSize
			min.Height = fyne.Max(min.Height, iconMinSize)
		}
		if r.label.Text != "" {
			if r.icon != nil {
				min.Width += theme.Padding()
			}
			min.Width += labelMinSize.Width
			min.Height = fyne.Max(min.Height, labelMinSize.Height)
		}
	}
	min = min.Add(r.padding())
	//log.Println("tabButtonRenderer.MinSize:", min)
	return
}

func (r *tabButtonRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *tabButtonRenderer) Refresh() {
	r.updateLabel()

	canvas.Refresh(r.button.super())
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

func (r *tabButtonRenderer) updateLabel() {
	r.label.Text = r.button.Text
	r.label.Color = theme.TextColor()
	r.label.TextSize = theme.TextSize()
}

var _ fyne.Widget = (*tabOverflowButton)(nil)
var _ fyne.Tappable = (*tabOverflowButton)(nil)
var _ desktop.Hoverable = (*tabOverflowButton)(nil)

type tabOverflowButton struct {
	tabButton
	container *TabContainer
	popUp     *PopUpMenu
}

func (b *tabOverflowButton) CreateRenderer() fyne.WidgetRenderer {
	b.ExtendBaseWidget(b)
	b.updateContent()
	b.updateOverflowPopUp()
	return b.tabButton.CreateRenderer()
}

// Hide hides the tabOverflowButton.
// Implements: fyne.Widget
func (b *tabOverflowButton) Hide() {
	b.ExtendBaseWidget(b)
	b.hideOverflowPopUp()
	b.BaseWidget.Hide()
}

// MinSize returns the size that this widget should not shrink below
// Implements: fyne.Widget
func (b *tabOverflowButton) MinSize() fyne.Size {
	b.ExtendBaseWidget(b)
	return b.BaseWidget.MinSize()
}

// Move changes the relative position of the tabOverflowButton.
// Implements: fyne.Widget
func (b *tabOverflowButton) Move(pos fyne.Position) {
	b.ExtendBaseWidget(b)
	b.BaseWidget.Move(pos)
	b.Refresh()
}

// Refresh updates the Icon, Text, and PopUp based on TabContainer.tabLocation.
// Implements: fyne.Widget
func (b *tabOverflowButton) Refresh() {
	b.ExtendBaseWidget(b)
	b.updateContent()
	b.updateOverflowPopUp()
	b.BaseWidget.Refresh()
}

func (b *tabOverflowButton) createOverflowPopUp() {
	log.Println("Creating Overflow PopUp")
	b.container.propertyLock.RLock()
	tabs := b.container.Items
	b.container.propertyLock.RUnlock()
	var items []*fyne.MenuItem
	for i, item := range tabs {
		index := i // capture
		text := item.Text
		if text == "" {
			text = fmt.Sprintf("Tab %d", i+1)
		}
		items = append(items, fyne.NewMenuItem(text, func() {
			b.container.propertyLock.Lock()
			b.container.overflowExpanded = false
			b.container.propertyLock.Unlock()
			b.container.SelectTabIndex(index)
		}))
	}

	c := fyne.CurrentApp().Driver().CanvasForObject(b.super())
	b.popUp = newPopUpMenu(fyne.NewMenu("", items...), c)
	b.popUp.Show()
}

func (b *tabOverflowButton) hideOverflowPopUp() {
	if b.popUp != nil {
		b.popUp.Dismiss()
		b.popUp.Hide()
		b.popUp = nil
	}
}

func (b *tabOverflowButton) updateContent() {
	hasText, hasIcon := b.container.contentTypes()
	if hasText {
		b.Text = "More"
	} else {
		b.Text = ""
	}
	if hasIcon {
		switch b.container.TabLocation() {
		case TabLocationLeading, TabLocationTrailing:
			b.Icon = theme.MoreVerticalIcon()
		default:
			b.Icon = theme.MoreHorizontalIcon()
		}
	} else {
		b.Icon = nil
	}
}

func (b *tabOverflowButton) updateOverflowPopUp() {
	b.container.propertyLock.RLock()
	expanded := b.container.overflowExpanded
	b.container.propertyLock.RUnlock()
	if expanded {
		log.Println("Overflow Expanded")
		if b.popUp == nil {
			b.createOverflowPopUp()
		}

		tabOverflowButtonPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(b.super())
		log.Println("TabOverflowButtonPos:", tabOverflowButtonPos)
		tabOverflowButtonSize := b.Size()
		log.Println("TabOverflowButtonSize:", tabOverflowButtonSize)
		popUpSize := b.popUp.MinSize()
		popUpPos := tabOverflowButtonPos
		switch b.container.TabLocation() {
		case TabLocationTop:
			popUpPos.X += tabOverflowButtonSize.Width - popUpSize.Width
			popUpPos.Y += tabOverflowButtonPos.Y + tabOverflowButtonSize.Height
		case TabLocationBottom:
			popUpPos.X += tabOverflowButtonSize.Width - popUpSize.Width
			popUpPos.Y -= popUpSize.Height
		case TabLocationLeading:
			popUpPos.X += tabOverflowButtonSize.Width
			popUpPos.Y += tabOverflowButtonSize.Height - popUpSize.Height
		case TabLocationTrailing:
			popUpPos.X -= popUpSize.Width
			popUpPos.Y += tabOverflowButtonSize.Height - popUpSize.Height
		}
		log.Println("PopUpPos:", popUpPos)
		b.popUp.Move(popUpPos)
		log.Println("PopUpSize:", popUpSize)
		b.popUp.Resize(popUpSize)
	} else {
		log.Println("Overflow Collapsed")
		b.hideOverflowPopUp()
	}
}
