package widget

import (
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
)

func TestRichText_HyperLink(t *testing.T) {
	text := NewRichText(&ParagraphSegment{Texts: []RichTextSegment{
		&TextSegment{Text: "Text"},
		&HyperlinkSegment{Text: "Link"},
	}})
	texts := test.WidgetRenderer(text).Objects()
	assert.Equal(t, "Text", texts[0].(*canvas.Text).Text)
	richLink := test.WidgetRenderer(texts[1].(*fyne.Container).Objects[0].(*Hyperlink)).Objects()[0].(fyne.Widget)
	linkText := test.WidgetRenderer(richLink).Objects()[0].(*canvas.Text)
	assert.Equal(t, "Link", linkText.Text)

	c := test.NewCanvas()
	c.SetContent(text)
	assert.Equal(t, texts[0].Position().Y, linkText.Position().Y)
}

func TestRichText_List(t *testing.T) {
	seg := trailingBoldErrorSegment()
	seg.Text = "Test"
	text := NewRichText(&ListSegment{Items: []RichTextSegment{
		seg,
	}})
	texts := test.WidgetRenderer(text).Objects()
	assert.Equal(t, "•", strings.TrimSpace(texts[0].(*canvas.Text).Text))
	assert.Equal(t, "Test", texts[1].(*canvas.Text).Text)
}

func TestRichText_OrderedList(t *testing.T) {
	text := NewRichText(&ListSegment{Ordered: true, Items: []RichTextSegment{
		&TextSegment{Text: "One"},
		&TextSegment{Text: "Two"},
	}})
	texts := test.WidgetRenderer(text).Objects()
	assert.Equal(t, "1.", strings.TrimSpace(texts[0].(*canvas.Text).Text))
	assert.Equal(t, "One", texts[1].(*canvas.Text).Text)
	assert.Equal(t, "2.", strings.TrimSpace(texts[2].(*canvas.Text).Text))
	assert.Equal(t, "Two", texts[3].(*canvas.Text).Text)
}
