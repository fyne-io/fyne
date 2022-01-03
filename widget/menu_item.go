package widget

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
)

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
	background.Hide()
	text := canvas.NewText(i.Item.Label, theme.ForegroundColor())
	text.Alignment = i.alignment
	objects := []fyne.CanvasObject{background, text}
	var icon *canvas.Image
	if i.Item.ChildMenu != nil {
		icon = canvas.NewImageFromResource(theme.MenuExpandIcon())
		objects = append(objects, icon)
	}
	checkIcon := canvas.NewImageFromResource(theme.ConfirmIcon())
	if !i.Item.Checked {
		checkIcon.Hide()
	}
	var shortcutTexts []*canvas.Text
	if s, ok := i.Item.Shortcut.(fyne.KeyboardShortcut); ok {
		shortcutTexts = textsForShortcut(s)
		for _, t := range shortcutTexts {
			objects = append(objects, t)
		}
	}

	objects = append(objects, checkIcon)
	return &menuItemRenderer{
		BaseRenderer:  widget.NewBaseRenderer(objects),
		i:             i,
		icon:          icon,
		checkIcon:     checkIcon,
		shortcutTexts: shortcutTexts,
		text:          text,
		background:    background,
	}
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
	icon             *canvas.Image
	checkIcon        *canvas.Image
	lastThemePadding float32
	minSize          fyne.Size
	shortcutTexts    []*canvas.Text
	text             *canvas.Text
	background       *canvas.Rectangle
}

func (r *menuItemRenderer) Layout(size fyne.Size) {
	padding := r.itemPadding()

	r.text.TextSize = theme.TextSize()
	r.text.Color = theme.ForegroundColor()
	if r.i.Item.Disabled {
		r.text.Color = theme.DisabledColor()
	}
	r.text.Resize(size.Subtract(fyne.NewSize(theme.Padding()*4, theme.Padding()*2)))
	r.text.Move(fyne.NewPos(padding.Width/2+r.checkSpace(), padding.Height/2))

	widthWithoutIcon := size.Width
	if r.icon != nil {
		widthWithoutIcon -= theme.IconInlineSize()
		r.icon.Resize(fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
		r.icon.Move(fyne.NewPos(widthWithoutIcon, (size.Height-theme.IconInlineSize())/2))
	}
	{
		offset := widthWithoutIcon - padding.Width/2
		for i := len(r.shortcutTexts) - 1; i >= 0; i-- {
			text := r.shortcutTexts[i]
			text.TextSize = theme.TextSize()
			text.Resize(text.MinSize())
			offset -= text.MinSize().Width
			text.Move(fyne.NewPos(offset, padding.Height/2))
		}
	}
	r.checkIcon.Resize(fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
	r.checkIcon.Move(fyne.NewPos(padding.Width/4, (size.Height-theme.IconInlineSize())/2))

	r.background.Resize(size)
}

func (r *menuItemRenderer) MinSize() fyne.Size {
	if r.minSizeUnchanged() {
		return r.minSize
	}

	minSize := r.text.MinSize().Add(r.itemPadding()).Add(fyne.NewSize(r.checkSpace(), 0))
	if r.icon != nil {
		minSize = minSize.Add(fyne.NewSize(theme.IconInlineSize(), 0))
	}
	if r.shortcutTexts != nil {
		var textWidth float32
		for _, text := range r.shortcutTexts {
			textWidth += text.MinSize().Width
		}
		minSize = minSize.AddWidthHeight(textWidth+theme.Padding()*2, 0)
	}
	r.minSize = minSize
	return r.minSize
}

func (r *menuItemRenderer) Refresh() {
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
	r.refreshText(r.text)
	for _, text := range r.shortcutTexts {
		r.refreshText(text)
	}

	if r.i.Item.Disabled {
		r.checkIcon.Resource = theme.NewDisabledResource(theme.ConfirmIcon())
	} else {
		r.checkIcon.Resource = theme.ConfirmIcon()
	}
	if r.i.Item.Checked {
		r.checkIcon.Show()
	} else {
		r.checkIcon.Hide()
	}
	r.checkIcon.Refresh()
	canvas.Refresh(r.i)
}

func (r *menuItemRenderer) checkSpace() float32 {
	if r.i.Parent.containsCheck {
		return theme.IconInlineSize()
	}
	return 0
}

func (r *menuItemRenderer) itemPadding() fyne.Size {
	return fyne.NewSize(theme.Padding()*4, theme.Padding()*2)
}

func (r *menuItemRenderer) minSizeUnchanged() bool {
	return !r.minSize.IsZero() &&
		r.text.TextSize == theme.TextSize() &&
		(r.icon == nil || r.icon.Size().Width == theme.IconInlineSize()) &&
		r.lastThemePadding == theme.Padding()
}

func (r *menuItemRenderer) refreshText(text *canvas.Text) {
	if r.i.Item.Disabled {
		text.Color = theme.DisabledColor()
	} else {
		text.Color = theme.ForegroundColor()
	}
	text.Refresh()
}

var keySymbols map[fyne.KeyName]rune = map[fyne.KeyName]rune{
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

func textsForShortcut(s fyne.KeyboardShortcut) (texts []*canvas.Text) {
	b := strings.Builder{}
	mods := s.Mod()
	if mods&fyne.KeyModifierControl != 0 {
		b.WriteRune('⌃')
	}
	if mods&fyne.KeyModifierAlt != 0 {
		b.WriteRune('⌥')
	}
	if mods&fyne.KeyModifierShift != 0 {
		b.WriteRune('⇧')
	}
	if mods&desktop.SuperModifier != 0 {
		b.WriteRune('⌘')
	}
	r := keySymbols[s.Key()]
	if r != 0 {
		b.WriteRune(r)
	}
	t := canvas.NewText(b.String(), theme.ForegroundColor())
	t.TextStyle.Symbol = true
	texts = append(texts, t)
	if r == 0 {
		texts = append(texts, canvas.NewText(string(s.Key()), theme.ForegroundColor()))
	}
	return
}
