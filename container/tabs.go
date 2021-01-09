package container

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// TabItem represents a single view in a TabContainer.
// The Text and Icon are used for the tab button and the Content is shown when the corresponding tab is active.
//
// Since: 1.4
type TabItem = widget.TabItem

// TabLocation is the location where the tabs of a tab container should be rendered
//
// Since: 1.4
type TabLocation = widget.TabLocation

// TabLocation values
const (
	TabLocationTop TabLocation = iota
	TabLocationLeading
	TabLocationBottom
	TabLocationTrailing
)

// NewTabItem creates a new item for a tabbed widget - each item specifies the content and a label for its tab.
//
// Since: 1.4
func NewTabItem(text string, content fyne.CanvasObject) *TabItem {
	return widget.NewTabItem(text, content)
}

// NewTabItemWithIcon creates a new item for a tabbed widget - each item specifies the content and a label with an icon for its tab.
//
// Since: 1.4
func NewTabItemWithIcon(text string, icon fyne.Resource, content fyne.CanvasObject) *TabItem {
	return widget.NewTabItemWithIcon(text, icon, content)
}

// TODO move the implementation into here in 2.0 when we delete the old API.
// we cannot do that right now due to Scroll dependency order.

type baseTabs struct {
	widget.BaseWidget

	Items              []*TabItem
	OnSelectionChanged func(tab *TabItem)

	current     int
	tabLocation TabLocation

	popUp *widget.PopUpMenu
}

// Append adds a new TabItem to the end of the tab panel
func (t *baseTabs) Append(item *TabItem) {
	t.SetItems(append(t.Items, item))
}

// Hide hides the widget.
//
// Implements: fyne.Widget
func (t *baseTabs) Hide() {
	if t.popUp != nil {
		t.popUp.Hide()
		t.popUp = nil
	}
	t.BaseWidget.Hide()
}

// Remove tab by value
func (t *baseTabs) Remove(item *TabItem) {
	for index, existingItem := range t.Items {
		if existingItem == item {
			t.RemoveIndex(index)
			break
		}
	}
}

// RemoveIndex removes tab by index
func (t *baseTabs) RemoveIndex(index int) {
	if index < 0 || index >= len(t.Items) {
		return
	}
	t.SetItems(append(t.Items[:index], t.Items[index+1:]...))
}

// Select sets the specified TabItem to be selected and its content visible.
func (t *baseTabs) Select(item *TabItem) {
	for i, child := range t.Items {
		if child == item {
			t.SelectIndex(i)
			return
		}
	}
}

// SelectIndex sets the TabItem at the specific index to be selected and its content visible.
func (t *baseTabs) SelectIndex(index int) {
	if index < 0 || index >= len(t.Items) || t.current == index {
		return
	}

	t.current = index
	t.Refresh()

	if f := t.OnSelectionChanged; f != nil {
		f(t.Items[t.current])
	}
}

// Selection returns the currently selected TabItem.
func (t *baseTabs) Selection() *TabItem {
	if t.current < 0 || t.current >= len(t.Items) {
		return nil
	}
	return t.Items[t.current]
}

// SelectionIndex returns the index of the currently selected TabItem.
func (t *baseTabs) SelectionIndex() int {
	return t.current
}

// SetItems sets the container’s items and refreshes.
func (t *baseTabs) SetItems(items []*TabItem) {
	if mismatchedTabItems(items) {
		internal.LogHint("Tab items should all have the same type of content (text, icons or both)")
	}
	t.Items = items
	if len(items) == 0 {
		// No items available to be current
		t.current = -1
	} else if t.current < 0 {
		// Current is first tab item
		t.current = 0
	}
	t.Refresh()
}

// SetTabLocation sets the location of the tab bar
func (t *baseTabs) SetTabLocation(l TabLocation) {
	t.tabLocation = l
	t.Refresh()
}

// Show this widget, if it was previously hidden
func (t *baseTabs) Show() {
	t.BaseWidget.Show()
	t.SelectIndex(t.current)
	t.Refresh()
}

func (t *baseTabs) showPopUp(button *widget.Button, items []*fyne.MenuItem) {
	d := fyne.CurrentApp().Driver()
	c := d.CanvasForObject(button)
	t.popUp = widget.NewPopUpMenu(fyne.NewMenu("", items...), c)
	buttonPos := d.AbsolutePositionForObject(button)
	buttonSize := button.Size()
	popUpMin := t.popUp.MinSize()
	var popUpPos fyne.Position
	switch t.tabLocation {
	case TabLocationLeading:
		popUpPos.X = buttonPos.X + buttonSize.Width
		popUpPos.Y = buttonPos.Y + buttonSize.Height - popUpMin.Height
	case TabLocationTrailing:
		popUpPos.X = buttonPos.X - popUpMin.Width
		popUpPos.Y = buttonPos.Y + buttonSize.Height - popUpMin.Height
	case TabLocationTop:
		popUpPos.X = buttonPos.X + buttonSize.Width - popUpMin.Width
		popUpPos.Y = buttonPos.Y + buttonSize.Height
	case TabLocationBottom:
		popUpPos.X = buttonPos.X + buttonSize.Width - popUpMin.Width
		popUpPos.Y = buttonPos.Y - popUpMin.Height
	}
	t.popUp.ShowAtPosition(popUpPos)
}

type baseTabsRenderer struct {
	animation *fyne.Animation

	action             *widget.Button
	bar                *fyne.Container
	divider, indicator *canvas.Rectangle

	buttonCache map[*TabItem]*tabButton
}

func (r *baseTabsRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *baseTabsRenderer) Destroy() {
}

func (r *baseTabsRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.bar, r.divider, r.indicator}
}

func (r *baseTabsRenderer) animateIndicator(p fyne.Position, s fyne.Size) {
	r.indicator.Show()
	if r.indicator.Position().IsZero() || r.indicator.Position() == p {
		r.indicator.Move(p)
		r.indicator.Resize(s)
	} else if r.animation == nil {
		r.animation = canvas.NewPositionAnimation(r.indicator.Position(), p, canvas.DurationShort, func(p fyne.Position) {
			r.indicator.Move(p)
			canvas.Refresh(r.indicator)
			if p == p {
				r.animation = nil
			}
		})
		r.animation.Start()

		canvas.NewSizeAnimation(r.indicator.Size(), s, canvas.DurationShort, func(s fyne.Size) {
			r.indicator.Resize(s)
			canvas.Refresh(r.indicator)
		}).Start()
	}
}

func (r *baseTabsRenderer) buildTabButtons(t *baseTabs, count int) *fyne.Container {
	buttons := &fyne.Container{}

	var iconPos buttonIconPosition
	if fyne.CurrentDevice().IsMobile() {
		cells := count
		if cells == 0 {
			cells = 1
		}
		buttons.Layout = layout.NewGridLayout(cells)
		iconPos = buttonIconTop
	} else if t.tabLocation == TabLocationLeading || t.tabLocation == TabLocationTrailing {
		buttons.Layout = layout.NewVBoxLayout()
		iconPos = buttonIconTop
	} else {
		buttons.Layout = layout.NewHBoxLayout()
		iconPos = buttonIconInline
	}

	for i := 0; i < count; i++ {
		item := t.Items[i]
		button, ok := r.buttonCache[item]
		if !ok {
			button = &tabButton{
				OnTap: func() { t.Select(item) },
			}
			r.buttonCache[item] = button
		}
		button.Text = item.Text
		button.Icon = item.Icon
		button.IconPosition = iconPos
		if i == t.current {
			button.Importance = widget.HighImportance
		} else {
			button.Importance = widget.MediumImportance
		}
		button.Refresh()
		buttons.Objects = append(buttons.Objects, button)
	}
	return buttons
}

func (r *baseTabsRenderer) layout(t *baseTabs, size fyne.Size) {
	var (
		barPos, dividerPos, contentPos    fyne.Position
		barSize, dividerSize, contentSize fyne.Size
	)

	barMin := r.bar.MinSize()

	switch t.tabLocation {
	case TabLocationTop:
		barHeight := barMin.Height
		barPos = fyne.NewPos(0, 0)
		barSize = fyne.NewSize(size.Width, barHeight)
		dividerPos = fyne.NewPos(0, barHeight)
		dividerSize = fyne.NewSize(size.Width, theme.Padding())
		contentPos = fyne.NewPos(0, barHeight+theme.Padding())
		contentSize = fyne.NewSize(size.Width, size.Height-barHeight-theme.Padding())
	case TabLocationLeading:
		barWidth := barMin.Width
		barPos = fyne.NewPos(0, 0)
		barSize = fyne.NewSize(barWidth, size.Height)
		dividerPos = fyne.NewPos(barWidth, 0)
		dividerSize = fyne.NewSize(theme.Padding(), size.Height)
		contentPos = fyne.NewPos(barWidth+theme.Padding(), 0)
		contentSize = fyne.NewSize(size.Width-barWidth-theme.Padding(), size.Height)
	case TabLocationBottom:
		barHeight := barMin.Height
		barPos = fyne.NewPos(0, size.Height-barHeight)
		barSize = fyne.NewSize(size.Width, barHeight)
		dividerPos = fyne.NewPos(0, size.Height-barHeight-theme.Padding())
		dividerSize = fyne.NewSize(size.Width, theme.Padding())
		contentPos = fyne.NewPos(0, 0)
		contentSize = fyne.NewSize(size.Width, size.Height-barHeight-theme.Padding())
	case TabLocationTrailing:
		barWidth := barMin.Width
		barPos = fyne.NewPos(size.Width-barWidth, 0)
		barSize = fyne.NewSize(barWidth, size.Height)
		dividerPos = fyne.NewPos(size.Width-barWidth-theme.Padding(), 0)
		dividerSize = fyne.NewSize(theme.Padding(), size.Height)
		contentPos = fyne.NewPos(0, 0)
		contentSize = fyne.NewSize(size.Width-barWidth-theme.Padding(), size.Height)
	}

	r.bar.Move(barPos)
	r.bar.Resize(barSize)
	r.divider.Move(dividerPos)
	r.divider.Resize(dividerSize)
	if t.current >= 0 && t.current < len(t.Items) {
		content := t.Items[t.current].Content
		content.Move(contentPos)
		content.Resize(contentSize)
	}
}

func (r *baseTabsRenderer) minSize(t *baseTabs) fyne.Size {
	barMin := r.bar.MinSize()

	contentMin := fyne.NewSize(0, 0)
	for _, content := range t.Items {
		contentMin = contentMin.Max(content.Content.MinSize())
	}

	switch t.tabLocation {
	case TabLocationLeading, TabLocationTrailing:
		return fyne.NewSize(barMin.Width+contentMin.Width+theme.Padding(), contentMin.Height)
	default:
		return fyne.NewSize(contentMin.Width, barMin.Height+contentMin.Height+theme.Padding())
	}
}

func (r *baseTabsRenderer) refresh(t *baseTabs) {
	if r.action != nil {
		if t.tabLocation == TabLocationLeading || t.tabLocation == TabLocationTrailing {
			r.action.SetIcon(theme.DotsVerticalIcon())
		} else {
			r.action.SetIcon(theme.DotsHorizontalIcon())
		}
	}

	r.bar.Refresh()

	r.divider.FillColor = theme.ShadowColor()
	r.divider.Refresh()

	r.indicator.FillColor = theme.PrimaryColor()
	r.indicator.Refresh()
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
	var icon *canvas.Image
	if b.Icon != nil {
		icon = canvas.NewImageFromResource(b.Icon)
	}

	label := canvas.NewText(b.Text, theme.TextColor())
	label.TextStyle.Bold = true
	label.Alignment = fyne.TextAlignCenter

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
	r.Refresh()
	return r
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
	r.label.Text = r.button.Text
	if r.button.Importance == widget.HighImportance {
		r.label.Color = theme.PrimaryColor()
	} else {
		r.label.Color = theme.TextColor()
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

func mismatchedTabItems(items []*TabItem) bool {
	var hasText, hasIcon bool
	for _, tab := range items {
		hasText = hasText || tab.Text != ""
		hasIcon = hasIcon || tab.Icon != nil
	}

	mismatch := false
	for _, tab := range items {
		if (hasText && tab.Text == "") || (hasIcon && tab.Icon == nil) {
			mismatch = true
			break
		}
	}

	return mismatch
}
