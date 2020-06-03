package widget

import (
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal/widget"
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

var _ fyne.Tappable = (*ToolbarButton)(nil)
var _ desktop.Hoverable = (*ToolbarButton)(nil)

// ToolbarButton is mostly like an ordinary button, but it can not be focused, only tapped
// The focusing is done by the toolbar. This is needed to avoid tabbing through all toolbar
// buttons. Instead the arrow keys can be used to select the toolbar button when the toolbar
// has focus.
type ToolbarButton struct {
	BaseWidget
	toolbar      *Toolbar
	hovered      bool
	focused      bool
	pressed      bool
	Icon         fyne.Resource
	IconPosition buttonIconPosition
	OnTap        func()
	Text         string
}

// MinSize returns the size that this widget should not shrink below
func (b *ToolbarButton) MinSize() fyne.Size {
	b.ExtendBaseWidget(b)
	return b.BaseWidget.MinSize()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (b *ToolbarButton) CreateRenderer() fyne.WidgetRenderer {
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
	return &ToolbarButtonRenderer{
		BaseRenderer: widget.NewBaseRenderer(objects),
		button:  b,
		icon:    icon,
		label:   label,
	}
}

// Tapped is called when a pointer tapped event is captured and triggers any tap handler
func (b *ToolbarButton) Tapped(e *fyne.PointEvent) {
	b.OnTap()
}

// MouseIn is called when a desktop pointer enters the widget
func (b *ToolbarButton) MouseIn(e *desktop.MouseEvent) {
	b.hovered = true
	canvas.Refresh(b)
}

// MouseOut is called when a desktop pointer exits the widget
func (b *ToolbarButton) MouseOut() {
	b.hovered = false
	canvas.Refresh(b)
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (b *ToolbarButton) MouseMoved(e *desktop.MouseEvent) {
}

// MouseDown called on mouse click
func (b *ToolbarButton) MouseDown(m *desktop.MouseEvent) {
	b.pressed = true
	b.focused = true
	for j,o:=range b.toolbar.renderer.objs {
		if tb,ok:=o.(*ToolbarButton); ok {
			if tb==b {
				tb.focused = true
				b.toolbar.current = j
			} else {
				tb.focused = false
			}
		}
	}
	b.Refresh()
}

// MouseUp called on mouse release
// If a mouse drag event has completed then check to see if it has resulted in an empty selection,
// if so, and if a text select key isn't held, then disable selecting
func (b *ToolbarButton) MouseUp(m *desktop.MouseEvent) {
	b.pressed = false
	b.focused = false
	b.Refresh()
}

type ToolbarButtonRenderer struct {
	widget.BaseRenderer
	button  *ToolbarButton
	icon    *canvas.Image
	label   *canvas.Text
}

// Refresh updates the widget state when requested.
func (r *ToolbarButtonRenderer) Refresh() {
	r.label.Color = theme.TextColor()
	r.label.TextSize = theme.TextSize()
	canvas.Refresh(r.button)
}

func (r *ToolbarButtonRenderer) padding() fyne.Size {
	return fyne.NewSize(theme.Padding()*2, theme.Padding()*2)
}

// BackgroundColor returns the theme background color.
// Implements: fyne.WidgetRenderer
func (r *ToolbarButtonRenderer) BackgroundColor() color.Color {
	switch {
	case r.button.pressed:
		return theme.PressedColor()
	case r.button.hovered:
		return theme.HoverColor()
	case r.button.focused:
		return theme.FocusColor()
	default:
		return theme.ButtonColor()
	}
}

// Destroy does nothing
func (r *ToolbarButtonRenderer) Destroy() {
}

// Layout the components of the button widget
func (r *ToolbarButtonRenderer) Layout(size fyne.Size) {
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

// MinSize calculates the smallest size that will fit the listed
func (r *ToolbarButtonRenderer) MinSize() fyne.Size {
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

func (r *ToolbarButtonRenderer) iconSize() int {
	return theme.IconInlineSize()
}

// newToolbarButton creates a new button widget with the specified label, themed icon and tap handler
func newToolbarButton(icon fyne.Resource, tapped func()) *ToolbarButton {
	button := &ToolbarButton{OnTap: tapped, Icon: icon}
	button.ExtendBaseWidget(button)
	return button
}
