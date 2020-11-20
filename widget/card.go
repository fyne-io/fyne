package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/theme"
)

// Card widget groups title, subtitle with content and a header image
//
// Since: 1.4
type Card struct {
	BaseWidget
	Title, Subtitle string
	Image           *canvas.Image
	Content         fyne.CanvasObject
}

// NewCard creates a new card widget with the specified title, subtitle and content (all optional).
//
// Since: 1.4
func NewCard(title, subtitle string, content fyne.CanvasObject) *Card {
	card := &Card{
		Title:    title,
		Subtitle: subtitle,
		Content:  content,
	}

	card.ExtendBaseWidget(card)
	return card
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (c *Card) CreateRenderer() fyne.WidgetRenderer {
	c.ExtendBaseWidget(c)

	header := canvas.NewText(c.Title, theme.TextColor())
	header.TextStyle.Bold = true
	subHeader := canvas.NewText(c.Subtitle, theme.TextColor())

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
func (c *Card) MinSize() fyne.Size {
	c.ExtendBaseWidget(c)
	return c.BaseWidget.MinSize()
}

// SetContent changes the body of this card to have the specified content.
func (c *Card) SetContent(obj fyne.CanvasObject) {
	c.Content = obj

	c.Refresh()
}

// SetImage changes the image displayed above the title for this card.
func (c *Card) SetImage(img *canvas.Image) {
	c.Image = img

	c.Refresh()
}

// SetSubTitle updates the secondary title for this card.
func (c *Card) SetSubTitle(text string) {
	c.Subtitle = text

	c.Refresh()
}

// SetTitle updates the main title for this card.
func (c *Card) SetTitle(text string) {
	c.Title = text

	c.Refresh()
}

type cardRenderer struct {
	*widget.ShadowingRenderer

	header, subHeader *canvas.Text

	card *Card
}

const (
	cardMediaHeight = 128
)

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

	contentPad := theme.Padding()
	if c.card.Title != "" || c.card.Subtitle != "" {
		titlePad := theme.Padding() * 2
		size.Width -= titlePad * 2
		pos.X += titlePad
		pos.Y += titlePad

		if c.card.Title != "" {
			height := c.header.MinSize().Height
			c.header.Move(pos)
			c.header.Resize(fyne.NewSize(size.Width, height))
			pos.Y += height + theme.Padding()
		}

		if c.card.Subtitle != "" {
			height := c.subHeader.MinSize().Height
			c.subHeader.Move(pos)
			c.subHeader.Resize(fyne.NewSize(size.Width, height))
			pos.Y += height + theme.Padding()
		}

		size.Width = size.Width + titlePad*2
		pos.X = pos.X - titlePad
		pos.Y += titlePad
	}

	size.Width -= contentPad * 2
	pos.X += contentPad
	if c.card.Content != nil {
		height := size.Height - contentPad*2 - (pos.Y - theme.Padding()/2) // adjust for content and initial offset
		if c.card.Title != "" || c.card.Subtitle != "" {
			height += contentPad
			pos.Y -= contentPad
		}
		c.card.Content.Move(pos.Add(fyne.NewPos(0, contentPad)))
		c.card.Content.Resize(fyne.NewSize(size.Width, height))
	}
}

// MinSize calculates the minimum size of a card.
// This is based on the contained text, image and content.
func (c *cardRenderer) MinSize() fyne.Size {
	hasHeader := c.card.Title != ""
	hasSubHeader := c.card.Subtitle != ""
	hasImage := c.card.Image != nil
	hasContent := c.card.Content != nil

	if !hasHeader && !hasSubHeader && !hasContent { // just image, or nothing
		if c.card.Image == nil {
			return fyne.NewSize(theme.Padding(), theme.Padding()) // empty, just space for border
		}
		return fyne.NewSize(c.card.Image.MinSize().Width+theme.Padding(), cardMediaHeight+theme.Padding())
	}

	contentPad := theme.Padding()
	min := fyne.NewSize(theme.Padding(), theme.Padding())
	if hasImage {
		min = fyne.NewSize(min.Width, min.Height+cardMediaHeight)
	}

	if hasHeader || hasSubHeader {
		titlePad := theme.Padding() * 2
		min = min.Add(fyne.NewSize(0, titlePad*2))
		if hasHeader {
			headerMin := c.header.MinSize()
			min = fyne.NewSize(fyne.Max(min.Width, headerMin.Width+titlePad*2+theme.Padding()),
				min.Height+headerMin.Height)
			if hasSubHeader {
				min.Height += theme.Padding()
			}
		}
		if hasSubHeader {
			subHeaderMin := c.subHeader.MinSize()
			min = fyne.NewSize(fyne.Max(min.Width, subHeaderMin.Width+titlePad*2+theme.Padding()),
				min.Height+subHeaderMin.Height)
		}
	}

	if hasContent {
		contentMin := c.card.Content.MinSize()
		min = fyne.NewSize(fyne.Max(min.Width, contentMin.Width+contentPad*2+theme.Padding()),
			min.Height+contentMin.Height+contentPad*2)
	}

	return min
}

func (c *cardRenderer) Refresh() {
	c.header.Text = c.card.Title
	c.header.Refresh()
	c.subHeader.Text = c.card.Subtitle
	c.subHeader.Refresh()

	objects := []fyne.CanvasObject{c.header, c.subHeader}
	if c.card.Image != nil {
		objects = append(objects, c.card.Image)
	}
	if c.card.Content != nil {
		objects = append(objects, c.card.Content)
	}
	c.ShadowingRenderer.SetObjects(objects)

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
