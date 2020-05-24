package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/theme"
)

// PopUpMenu is a Menu which displays itself in an OverlayContainer.
type PopUpMenu struct {
	*Menu
	canvas  fyne.Canvas
	overlay *widget.OverlayContainer
}

// ShowPopUpMenuAtPosition creates a PopUp menu populated with items from the passed menu structure.
// It will automatically be positioned at the provided location and shown as an overlay on the specified canvas.
func ShowPopUpMenuAtPosition(menu *fyne.Menu, c fyne.Canvas, pos fyne.Position) {
	m := newPopUpMenu(menu, c)
	m.ShowAtPosition(pos)
}

func newPopUpMenu(menu *fyne.Menu, c fyne.Canvas) *PopUpMenu {
	p := &PopUpMenu{Menu: NewMenu(menu), canvas: c}
	p.Menu.Resize(p.Menu.MinSize())
	p.Menu.customSized = true
	o := widget.NewOverlayContainer(p.Menu, c, p.Dismiss)
	o.Resize(o.MinSize())
	p.overlay = o

	focused := c.Focused()
	p.OnDismiss = func() {
		if c.Focused() == nil {
			c.Focus(focused)
		}
		p.Hide()
	}
	return p
}

// CreateRenderer returns a new renderer for the pop-up menu.
// Implements: fyne.Widget
func (p *PopUpMenu) CreateRenderer() fyne.WidgetRenderer {
	return p.overlay.CreateRenderer()
}

// Hide hides the pop-up menu.
// Implements: fyne.Widget
func (p *PopUpMenu) Hide() {
	p.overlay.Hide()
	p.Menu.Hide()
}

// Move moves the pop-up menu.
// The position is absolute because pop-up menus are shown in an overlay which covers the whole canvas.
// Implements: fyne.Widget
func (p *PopUpMenu) Move(pos fyne.Position) {
	widget.MoveWidget(&p.Base, p, p.adjustedPosition(pos, p.Size()))
}

// Resize changes the size of the pop-up menu.
// Implements: fyne.Widget
func (p *PopUpMenu) Resize(size fyne.Size) {
	widget.MoveWidget(&p.Base, p, p.adjustedPosition(p.Position(), size))
	p.Menu.Resize(size)
}

// Show makes the pop-up menu visible.
// Implements: fyne.Widget
func (p *PopUpMenu) Show() {
	p.overlay.Show()
	p.Menu.Show()
}

// ShowAtPosition shows the pop-up menu at the specified position.
func (p *PopUpMenu) ShowAtPosition(pos fyne.Position) {
	p.Move(pos)
	p.Show()
}

func (p *PopUpMenu) adjustedPosition(pos fyne.Position, size fyne.Size) fyne.Position {
	x := pos.X
	y := pos.Y
	if x+size.Width > p.canvas.Size().Width {
		x = p.canvas.Size().Width - size.Width
		if x < 0 {
			x = 0 // TODO here we may need a scroller as it's wider than our canvas
		}
	}
	if y+size.Height > p.canvas.Size().Height {
		y = p.canvas.Size().Height - size.Height
		if y < 0 {
			y = 0 // TODO here we may need a scroller as it's longer than our canvas
		}
	}
	return fyne.NewPos(x, y)
}

//
// Deprecated pop-up menu implementation
//

// NewPopUpMenuAtPosition creates a PopUp widget populated with menu items from the passed menu structure.
// It will automatically be positioned at the provided location and shown as an overlay on the specified canvas.
// Deprecated: Use ShowPopUpMenuAtPosition instead.
func NewPopUpMenuAtPosition(menu *fyne.Menu, c fyne.Canvas, pos fyne.Position) *PopUp {
	options := NewVBox()
	for _, option := range menu.Items {
		opt := option // capture value
		if opt.IsSeparator {
			options.Append(newSeparator())
		} else {
			options.Append(newMenuItemWidget(opt.Label))
		}
	}
	pop := NewPopUpAtPosition(options, c, pos)
	focused := c.Focused()
	for i, o := range options.Children {
		if label, ok := o.(*menuItemWidget); ok {
			item := menu.Items[i]
			label.OnTapped = func() {
				if c.Focused() == nil {
					c.Focus(focused)
				}
				pop.Hide()
				item.Action()
			}
		}
	}
	return pop
}

// NewPopUpMenu creates a PopUp widget populated with menu items from the passed menu structure.
// It will automatically be shown as an overlay on the specified canvas.
// Deprecated: Use ShowPopUpMenuAtPosition instead.
func NewPopUpMenu(menu *fyne.Menu, c fyne.Canvas) *PopUp {
	return NewPopUpMenuAtPosition(menu, c, fyne.NewPos(0, 0))
}

type menuItemWidget struct {
	*Label
	OnTapped func()
	hovered  bool
}

func (t *menuItemWidget) Tapped(*fyne.PointEvent) {
	t.OnTapped()
}

func (t *menuItemWidget) CreateRenderer() fyne.WidgetRenderer {
	return &menuItemWidgetRenderer{t.Label.CreateRenderer().(*textRenderer), t}
}

// MouseIn is called when a desktop pointer enters the widget
func (t *menuItemWidget) MouseIn(*desktop.MouseEvent) {
	t.hovered = true

	canvas.Refresh(t)
}

// MouseOut is called when a desktop pointer exits the widget
func (t *menuItemWidget) MouseOut() {
	t.hovered = false

	canvas.Refresh(t)
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (t *menuItemWidget) MouseMoved(*desktop.MouseEvent) {
}

func newMenuItemWidget(label string) *menuItemWidget {
	ret := &menuItemWidget{Label: NewLabel(label)}
	ret.ExtendBaseWidget(ret)
	return ret
}

type menuItemWidgetRenderer struct {
	*textRenderer
	label *menuItemWidget
}

func (h *menuItemWidgetRenderer) BackgroundColor() color.Color {
	if h.label.hovered {
		return theme.HoverColor()
	}

	return color.Transparent
}

func newSeparator() fyne.CanvasObject {
	s := canvas.NewRectangle(theme.DisabledTextColor())
	s.SetMinSize(fyne.NewSize(1, 2))
	return s
}
