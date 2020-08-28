package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/theme"
)

// CardContainer widget groups title, subtitle with content and a header image
type CardContainer struct {
	BaseWidget
	Title, SubTitle string
	Image           *canvas.Image
	Content         fyne.CanvasObject
}

// NewCardContainer creates a new card widget with the specified title, subtitle and content (all optional).
func NewCardContainer(title, subtitle string, content fyne.CanvasObject) *CardContainer {
	card := &CardContainer{
		Title:    title,
		SubTitle: subtitle,
		Content:  content,
	}

	card.ExtendBaseWidget(card)
	return card
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (c *CardContainer) CreateRenderer() fyne.WidgetRenderer {
	c.ExtendBaseWidget(c)

	header := canvas.NewText(c.Title, theme.TextColor())
	header.TextStyle.Bold = true
	subHeader := canvas.NewText(c.SubTitle, theme.TextColor())

	objects := []fyne.CanvasObject{header, subHeader}
	if c.Image != nil {
		objects = append(objects, c.Image)
	}
	if c.Content != nil {
		objects = append(objects, c.Content)
	}
	r := &cardRenderer{widget.NewShadowingRenderer(objects, widget.CardLevel),
		header, subHeader, c}
	r.applyTheme()
	return r
}

// MinSize returns the size that this widget should not shrink below
func (c *CardContainer) MinSize() fyne.Size {
	c.ExtendBaseWidget(c)
	return c.BaseWidget.MinSize()
}

// SetContent changes the body of this card to have the specified content.
func (c *CardContainer) SetContent(obj fyne.CanvasObject) {
	c.Content = obj

	c.Refresh()
}

// SetImage changes the image displayed above the title for this card.
func (c *CardContainer) SetImage(img *canvas.Image) {
	c.Image = img

	c.Refresh()
}

// SetSubTitle updates the secondary title for this card.
func (c *CardContainer) SetSubTitle(text string) {
	c.SubTitle = text

	c.Refresh()
}

// SetTitle updates the main title for this card.
func (c *CardContainer) SetTitle(text string) {
	c.Title = text

	c.Refresh()
}

type cardRenderer struct {
	*widget.ShadowingRenderer

	header, subHeader *canvas.Text

	card *CardContainer
}

const (
	cardMediaHeight = 128
)

func (c *cardRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

// Layout the components of the card container.
func (c *cardRenderer) Layout(size fyne.Size) {
	pos := fyne.NewPos(theme.Padding()/2, theme.Padding()/2)
	size = size.Subtract(fyne.NewSize(theme.Padding(), theme.Padding()))
	c.LayoutShadow(size, pos)

	if c.card.Image != nil {
		c.card.Image.Move(pos)
		c.card.Image.Resize(fyne.NewSize(size.Width, cardMediaHeight))
		pos.Y += cardMediaHeight
	}

	size.Width -= theme.Padding() * 4
	pos.X += theme.Padding() * 2
	pos.Y += theme.Padding()

	if c.card.Title != "" {
		height := c.header.MinSize().Height
		c.header.Move(pos)
		c.header.Resize(fyne.NewSize(size.Width, height))
		pos.Y += height + theme.Padding()
	}

	if c.card.SubTitle != "" {
		height := c.subHeader.MinSize().Height
		c.subHeader.Move(pos)
		c.subHeader.Resize(fyne.NewSize(size.Width, height))
		pos.Y += height + theme.Padding()
	}

	size.Width += theme.Padding() * 2
	pos.X -= theme.Padding()
	pos.Y += theme.Padding()

	if c.card.Content != nil {
		height := size.Height - pos.Y - theme.Padding()*3
		c.card.Content.Move(pos)
		c.card.Content.Resize(fyne.NewSize(size.Width, height))
	}
}

// MinSize calculates the minimum size of a card.
// This is based on the contained text, image and content.
func (c *cardRenderer) MinSize() fyne.Size {
	hasHeader := c.card.Title != ""
	hasSubHeader := c.card.SubTitle != ""
	hasImage := c.card.Image != nil
	hasContent := c.card.Content != nil

	if !hasHeader && !hasSubHeader && !hasContent {
		if c.card.Image == nil {
			return fyne.NewSize(theme.Padding(), theme.Padding()) // empty, just space for border
		}
		return fyne.NewSize(c.card.Image.MinSize().Width+theme.Padding(), cardMediaHeight+theme.Padding())
	}

	min := fyne.NewSize(theme.Padding()*5, theme.Padding()*5) // content padding plus 1 pad border
	if hasImage {
		min = fyne.NewSize(min.Width, min.Height+cardMediaHeight)
	}

	if hasHeader {
		min = fyne.NewSize(fyne.Max(min.Width, c.header.MinSize().Width+theme.Padding()*5),
			min.Height+c.header.MinSize().Height)
		if hasSubHeader || hasContent {
			min.Height += theme.Padding()
		}
	}
	if hasSubHeader {
		min = fyne.NewSize(fyne.Max(min.Width, c.subHeader.MinSize().Width+theme.Padding()*5),
			min.Height+c.subHeader.MinSize().Height)
		if hasContent {
			min.Height += theme.Padding()
		}
	}
	if hasContent {
		min = fyne.NewSize(fyne.Max(min.Width, c.card.Content.MinSize().Width+theme.Padding()*3),
			min.Height+c.card.Content.MinSize().Height)
	}

	return min
}

func (c *cardRenderer) Refresh() {
	c.header.Text = c.card.Title
	c.header.Refresh()
	c.subHeader.Text = c.card.SubTitle
	c.subHeader.Refresh()

	c.applyTheme()
	c.Layout(c.card.Size())
	canvas.Refresh(c.card.super())
}

// applyTheme updates this button to match the current theme
func (c *cardRenderer) applyTheme() {
	if c.header != nil {
		c.header.TextSize = int(float32(theme.TextSize()) * 1.7)
		c.header.Color = theme.TextColor()
	}
	if c.subHeader != nil {
		c.subHeader.TextSize = theme.TextSize()
		c.subHeader.Color = theme.TextColor()
	}
}
