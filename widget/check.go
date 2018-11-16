package widget

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/canvas"
import "github.com/fyne-io/fyne/theme"

type checkRenderer struct {
	background *canvas.Rectangle
	icon       *canvas.Image
	label      *canvas.Text

	objects []fyne.CanvasObject
	check   *Check
}

// MinSize calculates the minimum size of a check.
// This is based on the contained text, the check icon and a standard amount of padding added.
func (c *checkRenderer) MinSize() fyne.Size {
	min := c.label.MinSize().Add(fyne.NewSize(theme.Padding()*4, theme.Padding()*2))
	min = min.Add(fyne.NewSize(theme.IconInlineSize()+theme.Padding(), 0))

	return min
}

// Layout the components of the check widget
func (c *checkRenderer) Layout(size fyne.Size) {
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

// ApplyTheme is called when the Check may need to update it's look
func (c *checkRenderer) ApplyTheme() {
	c.label.Color = theme.TextColor()
	c.background.FillColor = theme.BackgroundColor()

	c.Refresh()
}

func (c *checkRenderer) Refresh() {
	c.label.Text = c.check.Text

	if c.check.Checked {
		c.icon.File = theme.CheckedIcon().CachePath()
	} else {
		c.icon.File = theme.UncheckedIcon().CachePath()
	}

	canvas.Refresh(c.check)
}

func (c *checkRenderer) Objects() []fyne.CanvasObject {
	return c.objects
}

// Check widget has a text label and a checked (or unchecked) icon and triggers an event func when toggled
type Check struct {
	baseWidget
	Text    string
	Checked bool

	OnChanged func(bool) `json:"-"`
}

// OnMouseDown is called when a mouse down event is captured and triggers any change handler
func (c *Check) OnMouseDown(*fyne.MouseEvent) {
	c.Checked = !c.Checked

	if c.OnChanged != nil {
		c.OnChanged(c.Checked)
	}
	c.Renderer().Refresh()
}

func (c *Check) createRenderer() fyne.WidgetRenderer {
	icon := canvas.NewImageFromResource(theme.UncheckedIcon())

	text := canvas.NewText(c.Text, theme.TextColor())
	text.Alignment = fyne.TextAlignCenter
	bg := canvas.NewRectangle(theme.BackgroundColor())

	return &checkRenderer{bg, icon, text, []fyne.CanvasObject{bg, icon, text}, c}
}

// Renderer is a private method to Fyne which links this widget to it's renderer
func (c *Check) Renderer() fyne.WidgetRenderer {
	if c.renderer == nil {
		c.renderer = c.createRenderer()
	}

	return c.renderer
}

// NewCheck creates a new check widget with the set label and change handler
func NewCheck(label string, changed func(bool)) *Check {
	c := &Check{
		baseWidget{},
		label,
		false,
		changed,
	}

	c.Renderer().Layout(c.MinSize())
	return c
}
