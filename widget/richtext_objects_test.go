package widget

import (
	"strings"
	"testing"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
)

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
