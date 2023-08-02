package widget

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
)

const (
	runeModifierAlt     = '⌥'
	runeModifierControl = '⌃'
	runeModifierShift   = '⇧'
)

var keySymbols = map[fyne.KeyName]rune{
	fyne.KeyBackspace: '⌫',
	fyne.KeyDelete:    '⌦',
	fyne.KeyDown:      '↓',
	fyne.KeyEnd:       '↘',
	fyne.KeyEnter:     '↩',
	fyne.KeyEscape:    '⎋',
	fyne.KeyHome:      '↖',
	fyne.KeyLeft:      '←',
	fyne.KeyPageDown:  '⇟',
	fyne.KeyPageUp:    '⇞',
	fyne.KeyReturn:    '↩',
	fyne.KeyRight:     '→',
	fyne.KeySpace:     '␣',
	fyne.KeyTab:       '⇥',
	fyne.KeyUp:        '↑',
}

var _ fyne.Widget = (*menuItem)(nil)

// menuItem is a widget for displaying a fyne.menuItem.
type menuItem struct {
	widget.Base
	Item   *fyne.MenuItem
	Parent *Menu

	alignment fyne.TextAlign
	child     *Menu
}

// newMenuItem creates a new menuItem.
func newMenuItem(item *fyne.MenuItem, parent *Menu) *menuItem {
	i := &menuItem{Item: item, Parent: parent}
	i.alignment = parent.alignment
	i.ExtendBaseWidget(i)
	return i
}

func (i *menuItem) Child() *Menu {
	if i.Item.ChildMenu != nil && i.child == nil {
		child := NewMenu(i.Item.ChildMenu)
		child.Hide()
		child.OnDismiss = i.Parent.Dismiss
		i.child = child
	}
	return i.child
}

// CreateRenderer returns a new renderer for the menu item.
//
// Implements: fyne.Widget
func (i *menuItem) CreateRenderer() fyne.WidgetRenderer {
	background := canvas.NewRectangle(theme.HoverColor())
	background.CornerRadius = theme.SelectionRadiusSize()
	background.Hide()
	text := canvas.NewText(i.Item.Label, theme.ForegroundColor())
	text.Alignment = i.alignment
	objects := []fyne.CanvasObject{background, text}
	var expandIcon *canvas.Image
	if i.Item.ChildMenu != nil {
		expandIcon = canvas.NewImageFromResource(theme.MenuExpandIcon())
		objects = append(objects, expandIcon)
	}
	checkIcon := canvas.NewImageFromResource(theme.ConfirmIcon())
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
		shortcutTexts = textsForShortcut(s)
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
	i.Parent.activateItem(i)
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
	i.Parent.DeactivateChild()
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
	return i.Parent.activeItem == i
}

func (i *menuItem) isSubmenuOpen() bool {
	return i.Child() != nil && i.Child().Visible()
}

func (i *menuItem) trigger() {
	i.Parent.Dismiss()
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
	leftOffset := theme.InnerPadding() + r.checkSpace()
	rightOffset := size.Width
	iconSize := fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize())
	iconTopOffset := (size.Height - theme.IconInlineSize()) / 2

	if r.expandIcon != nil {
		rightOffset -= theme.IconInlineSize()
		r.expandIcon.Resize(iconSize)
		r.expandIcon.Move(fyne.NewPos(rightOffset, iconTopOffset))
	}

	rightOffset -= theme.InnerPadding()
	textHeight := r.text.MinSize().Height
	for i := len(r.shortcutTexts) - 1; i >= 0; i-- {
		text := r.shortcutTexts[i]
		text.Resize(text.MinSize())
		rightOffset -= text.MinSize().Width
		text.Move(fyne.NewPos(rightOffset, theme.InnerPadding()+(textHeight-text.Size().Height)))

		if i == 0 {
			rightOffset -= theme.InnerPadding()
		}
	}

	r.checkIcon.Resize(iconSize)
	r.checkIcon.Move(fyne.NewPos(theme.InnerPadding(), iconTopOffset))

	if r.icon != nil {
		r.icon.Resize(iconSize)
		r.icon.Move(fyne.NewPos(leftOffset, iconTopOffset))
		leftOffset += theme.IconInlineSize()
		leftOffset += theme.InnerPadding()
	}

	r.text.Resize(fyne.NewSize(rightOffset-leftOffset, textHeight))
	r.text.Move(fyne.NewPos(leftOffset, theme.InnerPadding()))

	r.background.Resize(size)
}

func (r *menuItemRenderer) MinSize() fyne.Size {
	if r.minSizeUnchanged() {
		return r.minSize
	}

	minSize := r.text.MinSize().AddWidthHeight(theme.InnerPadding()*2+r.checkSpace(), theme.InnerPadding()*2)
	if r.expandIcon != nil {
		minSize = minSize.AddWidthHeight(theme.IconInlineSize(), 0)
	}
	if r.icon != nil {
		minSize = minSize.AddWidthHeight(theme.IconInlineSize()+theme.InnerPadding(), 0)
	}
	if r.shortcutTexts != nil {
		var textWidth float32
		for _, text := range r.shortcutTexts {
			textWidth += text.MinSize().Width
		}
		minSize = minSize.AddWidthHeight(textWidth+theme.InnerPadding(), 0)
	}
	r.minSize = minSize
	return r.minSize
}

func (r *menuItemRenderer) updateVisuals() {
	r.background.CornerRadius = theme.SelectionRadiusSize()
	if fyne.CurrentDevice().IsMobile() {
		r.background.Hide()
	} else if r.i.isActive() {
		r.background.FillColor = theme.FocusColor()
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
	r.updateIcon(r.checkIcon, theme.ConfirmIcon())
	r.updateIcon(r.expandIcon, theme.MenuExpandIcon())
	r.updateIcon(r.icon, r.i.Item.Icon)
}

func (r *menuItemRenderer) Refresh() {
	r.updateVisuals()
	canvas.Refresh(r.i)
}

func (r *menuItemRenderer) checkSpace() float32 {
	if r.i.Parent.containsCheck {
		return theme.IconInlineSize() + theme.InnerPadding()
	}
	return 0
}

func (r *menuItemRenderer) minSizeUnchanged() bool {
	return !r.minSize.IsZero() &&
		r.text.TextSize == theme.TextSize() &&
		(r.expandIcon == nil || r.expandIcon.Size().Width == theme.IconInlineSize()) &&
		r.lastThemePadding == theme.InnerPadding()
}

func (r *menuItemRenderer) updateIcon(img *canvas.Image, rsc fyne.Resource) {
	if img == nil {
		return
	}
	if r.i.Item.Disabled {
		img.Resource = theme.NewDisabledResource(rsc)
	} else {
		img.Resource = rsc
	}
}

func (r *menuItemRenderer) refreshText(text *canvas.Text, shortcut bool) {
	text.TextSize = theme.TextSize()
	if r.i.Item.Disabled {
		text.Color = theme.DisabledColor()
	} else {
		if shortcut {
			text.Color = shortcutColor()
		} else {
			text.Color = theme.ForegroundColor()
		}
	}
	text.Refresh()
}

func shortcutColor() color.Color {
	r, g, b, a := theme.ForegroundColor().RGBA()
	a = uint32(float32(a) * 0.95)
	return color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
}

func textsForShortcut(s fyne.KeyboardShortcut) (texts []*canvas.Text) {
	b := strings.Builder{}
	mods := s.Mod()
	if mods&fyne.KeyModifierControl != 0 {
		b.WriteRune(runeModifierControl)
	}
	if mods&fyne.KeyModifierAlt != 0 {
		b.WriteRune(runeModifierAlt)
	}
	if mods&fyne.KeyModifierShift != 0 {
		b.WriteRune(runeModifierShift)
	}
	if mods&fyne.KeyModifierSuper != 0 {
		b.WriteRune(runeModifierSuper)
	}
	r := keySymbols[s.Key()]
	if r != 0 {
		b.WriteRune(r)
	}
	shortColor := shortcutColor()
	t := canvas.NewText(b.String(), shortColor)
	t.TextStyle.Symbol = true
	texts = append(texts, t)
	if r == 0 {
		text := canvas.NewText(string(s.Key()), shortColor)
		text.TextStyle.Monospace = true
		texts = append(texts, text)
	}
	return
}
