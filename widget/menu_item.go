package widget

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/svg"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
)

var _ fyne.Widget = (*menuItem)(nil)

// menuItem is a widget for displaying a fyne.menuItem.
type menuItem struct {
	widget.Base
	Item *fyne.MenuItem

	alignment     fyne.TextAlign
	child, parent *Menu
}

// newMenuItem creates a new menuItem.
func newMenuItem(item *fyne.MenuItem, parent *Menu) *menuItem {
	i := &menuItem{Item: item, parent: parent}
	i.alignment = parent.alignment
	i.ExtendBaseWidget(i)
	return i
}

func (i *menuItem) Child() *Menu {
	if i.Item.ChildMenu != nil && i.child == nil {
		child := NewMenu(i.Item.ChildMenu)
		child.Hide()
		child.OnDismiss = i.parent.Dismiss
		i.child = child
	}
	return i.child
}

// CreateRenderer returns a new renderer for the menu item.
//
// Implements: fyne.Widget
func (i *menuItem) CreateRenderer() fyne.WidgetRenderer {
	background := canvas.NewRectangle(theme.ColorForWidget(theme.ColorNameHover, i.parent))
	background.CornerRadius = theme.SizeForWidget(theme.SizeNameSelectionRadius, i.parent)
	background.Hide()
	text := canvas.NewText(i.Item.Label, theme.ColorForWidget(theme.ColorNameForeground, i.parent))
	text.Alignment = i.alignment
	objects := []fyne.CanvasObject{background, text}
	var expandIcon *canvas.Image
	if i.Item.ChildMenu != nil {
		expandIcon = canvas.NewImageFromResource(theme.IconForWidget(theme.IconNameMenuExpand, i.parent))
		objects = append(objects, expandIcon)
	}
	checkIcon := canvas.NewImageFromResource(theme.IconForWidget(theme.IconNameConfirm, i.parent))
	if !i.Item.Checked {
		checkIcon.Hide()
	}
	var icon *canvas.Image
	if i.Item.Icon != nil {
		icon = canvas.NewImageFromResource(i.Item.Icon)
		objects = append(objects, icon)
	}
	var shortcutTexts []*canvas.Text
	if s, ok := i.Item.Shortcut.(fyne.KeyboardShortcut); ok {
		shortcutTexts = i.textsForShortcut(s)
		for _, t := range shortcutTexts {
			objects = append(objects, t)
		}
	}

	objects = append(objects, checkIcon)
	r := &menuItemRenderer{
		BaseRenderer:  widget.NewBaseRenderer(objects),
		i:             i,
		expandIcon:    expandIcon,
		checkIcon:     checkIcon,
		icon:          icon,
		shortcutTexts: shortcutTexts,
		text:          text,
		background:    background,
	}
	r.updateVisuals()
	return r
}

// MouseIn activates the item which shows the submenu if the item has one.
// The submenu of any sibling of the item will be hidden.
//
// Implements: desktop.Hoverable
func (i *menuItem) MouseIn(*desktop.MouseEvent) {
	i.activate()
}

// MouseMoved does nothing.
//
// Implements: desktop.Hoverable
func (i *menuItem) MouseMoved(*desktop.MouseEvent) {
}

// MouseOut deactivates the item unless it has an open submenu.
//
// Implements: desktop.Hoverable
func (i *menuItem) MouseOut() {
	if !i.isSubmenuOpen() {
		i.deactivate()
	}
}

// Tapped performs the action of the item and dismisses the menu.
// It does nothing if the item doesn’t have an action.
//
// Implements: fyne.Tappable
func (i *menuItem) Tapped(*fyne.PointEvent) {
	if i.Item.Disabled {
		return
	}
	if i.Item.Action == nil {
		if fyne.CurrentDevice().IsMobile() {
			i.activate()
		}

		return
	} else if i.Item.ChildMenu != nil {
		if fyne.CurrentDevice().IsMobile() {
			i.activate()
		}

		return
	}

	i.trigger()
}

func (i *menuItem) activate() {
	if i.Item.Disabled {
		return
	}
	if i.Child() != nil {
		i.Child().Show()
	}
	i.parent.activateItem(i)
}

func (i *menuItem) activateLastSubmenu() bool {
	if i.Child() == nil {
		return false
	}
	if i.isSubmenuOpen() {
		return i.Child().ActivateLastSubmenu()
	}
	i.Child().Show()
	i.Child().ActivateNext()
	return true
}

func (i *menuItem) deactivate() {
	if i.Child() != nil {
		i.Child().Hide()
	}
	i.parent.DeactivateChild()
}

func (i *menuItem) deactivateLastSubmenu() bool {
	if !i.isSubmenuOpen() {
		return false
	}
	if !i.Child().DeactivateLastSubmenu() {
		i.Child().DeactivateChild()
		i.Child().Hide()
	}
	return true
}

func (i *menuItem) isActive() bool {
	return i.parent.activeItem == i
}

func (i *menuItem) isSubmenuOpen() bool {
	return i.Child() != nil && i.Child().Visible()
}

func (i *menuItem) trigger() {
	i.parent.Dismiss()
	if i.Item.Action != nil {
		i.Item.Action()
	}
}

func (i *menuItem) triggerLast() {
	if i.isSubmenuOpen() {
		i.Child().TriggerLast()
		return
	}
	i.trigger()
}

type menuItemRenderer struct {
	widget.BaseRenderer
	i                *menuItem
	background       *canvas.Rectangle
	checkIcon        *canvas.Image
	expandIcon       *canvas.Image
	icon             *canvas.Image
	lastThemePadding float32
	minSize          fyne.Size
	shortcutTexts    []*canvas.Text
	text             *canvas.Text
}

func (r *menuItemRenderer) Layout(size fyne.Size) {
	th := r.i.parent.Theme()
	innerPad := th.Size(theme.SizeNameInnerPadding)
	inlineIcon := th.Size(theme.SizeNameInlineIcon)

	leftOffset := innerPad + r.checkSpace()
	rightOffset := size.Width
	iconSize := fyne.NewSquareSize(inlineIcon)
	iconTopOffset := (size.Height - inlineIcon) / 2

	if r.expandIcon != nil {
		rightOffset -= inlineIcon
		r.expandIcon.Resize(iconSize)
		r.expandIcon.Move(fyne.NewPos(rightOffset, iconTopOffset))
	}

	rightOffset -= innerPad
	textHeight := r.text.MinSize().Height
	for i := len(r.shortcutTexts) - 1; i >= 0; i-- {
		text := r.shortcutTexts[i]
		text.Resize(text.MinSize())
		rightOffset -= text.MinSize().Width
		text.Move(fyne.NewPos(rightOffset, innerPad+(textHeight-text.Size().Height)))

		if i == 0 {
			rightOffset -= innerPad
		}
	}

	r.checkIcon.Resize(iconSize)
	r.checkIcon.Move(fyne.NewPos(innerPad, iconTopOffset))

	if r.icon != nil {
		r.icon.Resize(iconSize)
		r.icon.Move(fyne.NewPos(leftOffset, iconTopOffset))
		leftOffset += inlineIcon
		leftOffset += innerPad
	}

	r.text.Resize(fyne.NewSize(rightOffset-leftOffset, textHeight))
	r.text.Move(fyne.NewPos(leftOffset, innerPad))

	r.background.Resize(size)
}

func (r *menuItemRenderer) MinSize() fyne.Size {
	if r.minSizeUnchanged() {
		return r.minSize
	}

	th := r.i.parent.Theme()
	innerPad := th.Size(theme.SizeNameInnerPadding)
	inlineIcon := th.Size(theme.SizeNameInlineIcon)
	innerPad2 := innerPad * 2

	minSize := r.text.MinSize().AddWidthHeight(innerPad2+r.checkSpace(), innerPad2)
	if r.expandIcon != nil {
		minSize = minSize.AddWidthHeight(inlineIcon, 0)
	}
	if r.icon != nil {
		minSize = minSize.AddWidthHeight(inlineIcon+innerPad, 0)
	}
	if r.shortcutTexts != nil {
		var textWidth float32
		for _, text := range r.shortcutTexts {
			textWidth += text.MinSize().Width
		}
		minSize = minSize.AddWidthHeight(textWidth+innerPad, 0)
	}
	r.minSize = minSize
	return r.minSize
}

func (r *menuItemRenderer) updateVisuals() {
	r.background.CornerRadius = theme.SizeForWidget(theme.SizeNameSelectionRadius, r.i.parent)
	if fyne.CurrentDevice().IsMobile() {
		r.background.Hide()
	} else if r.i.isActive() {
		r.background.FillColor = theme.ColorForWidget(theme.ColorNameFocus, r.i.parent)
		r.background.Show()
	} else {
		r.background.Hide()
	}
	r.background.Refresh()
	r.text.Alignment = r.i.alignment
	r.refreshText(r.text, false)
	for _, text := range r.shortcutTexts {
		r.refreshText(text, true)
	}

	if r.i.Item.Checked {
		r.checkIcon.Show()
	} else {
		r.checkIcon.Hide()
	}
	r.updateIcon(r.checkIcon, theme.IconForWidget(theme.IconNameConfirm, r.i.parent))
	r.updateIcon(r.expandIcon, theme.IconForWidget(theme.IconNameMenuExpand, r.i.parent))
	r.updateIcon(r.icon, r.i.Item.Icon)
}

func (r *menuItemRenderer) Refresh() {
	r.updateVisuals()
	canvas.Refresh(r.i)
}

func (r *menuItemRenderer) checkSpace() float32 {
	if r.i.parent.containsCheck {
		return theme.IconInlineSize() + theme.InnerPadding()
	}
	return 0
}

func (r *menuItemRenderer) minSizeUnchanged() bool {
	th := r.i.parent.Theme()

	return !r.minSize.IsZero() &&
		r.text.TextSize == th.Size(theme.SizeNameText) &&
		(r.expandIcon == nil || r.expandIcon.Size().Width == th.Size(theme.SizeNameInlineIcon)) &&
		r.lastThemePadding == th.Size(theme.SizeNameInnerPadding)
}

func (r *menuItemRenderer) updateIcon(img *canvas.Image, rsc fyne.Resource) {
	if img == nil {
		return
	}
	if r.i.Item.Disabled && svg.IsResourceSVG(rsc) {
		img.Resource = theme.NewDisabledResource(rsc)
	} else {
		img.Resource = rsc
	}
}

func (r *menuItemRenderer) refreshText(text *canvas.Text, shortcut bool) {
	text.TextSize = theme.SizeForWidget(theme.SizeNameText, r.i.parent)
	if r.i.Item.Disabled {
		text.Color = theme.ColorForWidget(theme.ColorNameDisabled, r.i.parent)
	} else {
		if shortcut {
			text.Color = r.i.shortcutColor()
		} else {
			text.Color = theme.ColorForWidget(theme.ColorNameForeground, r.i.parent)
		}
	}
	text.Refresh()
}

func (i *menuItem) shortcutColor() color.Color {
	r, g, b, a := theme.ColorForWidget(theme.ColorNameForeground, i.parent).RGBA()
	a = uint32(float32(a) * 0.95)
	return color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
}

func (i *menuItem) textsForShortcut(sc fyne.KeyboardShortcut) (texts []*canvas.Text) {
	// add modifier
	b := strings.Builder{}
	mods := sc.Mod()
	if mods&fyne.KeyModifierControl != 0 {
		b.WriteString(textModifierControl)
	}
	if mods&fyne.KeyModifierAlt != 0 {
		b.WriteString(textModifierAlt)
	}
	if mods&fyne.KeyModifierShift != 0 {
		b.WriteString(textModifierShift)
	}
	if mods&fyne.KeyModifierSuper != 0 {
		b.WriteString(textModifierSuper)
	}
	shortColor := i.shortcutColor()
	if b.Len() > 0 {
		t := canvas.NewText(b.String(), shortColor)
		t.TextStyle = styleModifiers
		texts = append(texts, t)
	}
	// add key
	style := defaultStyleKeys
	s, ok := keyTexts[sc.Key()]
	if !ok {
		s = string(sc.Key())
	} else if len(s) == 1 {
		style = fyne.TextStyle{Symbol: true}
	}
	t := canvas.NewText(s, shortColor)
	t.TextStyle = style
	texts = append(texts, t)
	return
}
