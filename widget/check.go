package widget

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/canvas"
import "github.com/fyne-io/fyne/theme"

type checkLayout struct {
	background *canvas.Rectangle
	icon       *canvas.Image
	label      *canvas.Text
}

// MinSize calculates the minimum size of a check.
// This is based on the contained text, the check icon and a standard amount of padding added.
func (c *checkLayout) MinSize([]fyne.CanvasObject) fyne.Size {
	min := c.label.MinSize().Add(fyne.NewSize(theme.Padding()*4, theme.Padding()*2))
	min = min.Add(fyne.NewSize(theme.IconInlineSize()+theme.Padding(), 0))

	return min
}

// Layout the components of the check widget
func (c *checkLayout) Layout(_ []fyne.CanvasObject, size fyne.Size) {
	c.background.Resize(size)

	offset := fyne.NewSize(theme.IconInlineSize()+theme.Padding(), 0)
	labelSize := size.Subtract(offset)
	c.label.Resize(labelSize)
	c.label.Move(fyne.NewPos(theme.IconInlineSize()+theme.Padding(), 0))

	c.icon.Resize(fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
	c.icon.Move(fyne.NewPos(
		(size.Width-theme.IconInlineSize()-c.label.MinSize().Width)/2,
		(size.Height-theme.IconInlineSize())/2))
}

// Check widget has a text label and a checked (or unchecked) icon and triggers an event func when toggled
type Check struct {
	baseWidget
	Checked bool

	OnChanged func(bool)
	layout    *checkLayout
}

// OnMouseDown is called when a mouse down event is captured and triggers any change handler
func (c *Check) OnMouseDown(*fyne.MouseEvent) {
	c.Checked = !c.Checked
	c.ApplyTheme()

	if c.OnChanged != nil {
		c.OnChanged(c.Checked)
	}
	fyne.GetCanvas(c).Refresh(c)
}

// ApplyTheme is called when the Check may need to update it's look
func (c *Check) ApplyTheme() {
	c.layout.label.Color = theme.TextColor()
	c.layout.background.FillColor = theme.BackgroundColor()

	if c.Checked {
		c.layout.icon.File = theme.CheckedIcon().CachePath()
	} else {
		c.layout.icon.File = theme.UncheckedIcon().CachePath()
	}
}

// NewCheck creates a new check widget with the set label and change handler
func NewCheck(label string, changed func(bool)) *Check {
	icon := canvas.NewImageFromResource(theme.UncheckedIcon())

	text := canvas.NewText(label, theme.TextColor())
	text.Alignment = fyne.TextAlignCenter
	bg := canvas.NewRectangle(theme.BackgroundColor())
	layout := &checkLayout{bg, icon, text}

	objects := []fyne.CanvasObject{
		bg,
		text,
		icon,
	}

	c := &Check{
		baseWidget{
			objects: objects,
			layout:  layout,
		},
		false,
		changed,
		layout,
	}

	c.Layout(c.MinSize())
	return c
}
