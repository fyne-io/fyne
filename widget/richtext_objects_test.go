package widget

import (
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/test"
)

func TestRichText_Image(t *testing.T) {
	img := &ImageSegment{Title: "test", Source: storage.NewFileURI("./testdata/richtext/richtext_multiline.png")}
	text := NewRichText(img)
	texts := test.TempWidgetRenderer(t, text).Objects()
	drawn := texts[0].(*richImage).img

	text.Resize(fyne.NewSize(200, 200))
	assert.Equal(t, float32(0), drawn.Position().X)

	img.Alignment = fyne.TextAlignCenter
	text.Refresh()
	assert.Less(t, float32(0), drawn.Position().X)
	assert.Less(t, drawn.Position().X, text.Size().Width/2)

	img.Alignment = fyne.TextAlignTrailing
	text.Refresh()
	assert.Greater(t, float32(200), drawn.Position().X)
	assert.Greater(t, drawn.Position().X, text.Size().Width/2)
}

func TestRichText_HyperLink(t *testing.T) {
	text := NewRichText(&ParagraphSegment{Texts: []RichTextSegment{
		&TextSegment{Text: "Text"},
		&HyperlinkSegment{Text: "Link"},
	}})
	texts := test.TempWidgetRenderer(t, text).Objects()
	assert.Equal(t, "Text", texts[0].(*canvas.Text).Text)
	richLink := test.TempWidgetRenderer(t, texts[1].(*fyne.Container).Objects[0].(*Hyperlink)).Objects()[0].(fyne.Widget)
	linkText := test.TempWidgetRenderer(t, richLink).Objects()[0].(*canvas.Text)
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
	texts := test.TempWidgetRenderer(t, text).Objects()
	assert.Equal(t, "•", strings.TrimSpace(texts[0].(*canvas.Text).Text))
	assert.Equal(t, "Test", texts[1].(*canvas.Text).Text)
}

func TestRichText_OrderedList(t *testing.T) {
	text := NewRichText(&ListSegment{Ordered: true, Items: []RichTextSegment{
		&TextSegment{Text: "One"},
		&TextSegment{Text: "Two"},
	}})
	texts := test.TempWidgetRenderer(t, text).Objects()
	assert.Equal(t, "1.", strings.TrimSpace(texts[0].(*canvas.Text).Text))
	assert.Equal(t, "One", texts[1].(*canvas.Text).Text)
	assert.Equal(t, "2.", strings.TrimSpace(texts[2].(*canvas.Text).Text))
	assert.Equal(t, "Two", texts[3].(*canvas.Text).Text)
}

func TestRichText_HyperLink_WrappedHoverSynced(t *testing.T) {
	u, _ := url.Parse("https://fyne.io")
	seg := &HyperlinkSegment{Text: "this is a long hyperlink that wraps", URL: u}
	rt := NewRichText(seg)
	rt.Wrapping = fyne.TextWrapWord
	// Render narrow enough to force wrapping into at least 2 rows.
	rt.Resize(fyne.NewSize(100, 200))

	objs := test.TempWidgetRenderer(t, rt).Objects()
	// Collect all Hyperlink visuals for the segment (one per wrapped row).
	var links []*Hyperlink
	for _, obj := range objs {
		if c, ok := obj.(*fyne.Container); ok {
			if hl, ok := c.Objects[0].(*Hyperlink); ok {
				links = append(links, hl)
			}
		}
	}
	assert.GreaterOrEqual(t, len(links), 2, "expected hyperlink to wrap into multiple segments")

	// MouseIn on the first segment — choose a position inside the hyperlink's bounds.
	center := fyne.NewPos(links[0].Size().Width/2, links[0].Size().Height/2)
	inside := &desktop.MouseEvent{PointEvent: fyne.PointEvent{Position: center}}
	links[0].MouseIn(inside)
	for i, hl := range links {
		assert.True(t, hl.hovered, "expected link[%d] to be hovered after MouseIn on link[0]", i)
	}

	// MouseOut on the first segment — all siblings should be unhovered too.
	links[0].MouseOut()
	for i, hl := range links {
		assert.False(t, hl.hovered, "expected link[%d] to be unhovered after MouseOut on link[0]", i)
	}
}

func TestRichText_OrderedListDifferentIndex(t *testing.T) {
	for name, tt := range map[string]struct {
		index        int
		text1, text2 string
	}{
		"Start at -1": {index: -1, text1: "-1.", text2: "0."},
		"Start at 0":  {index: 0, text1: "0.", text2: "1."},
		"Start at 1":  {index: 1, text1: "1.", text2: "2."},
		"Start at 2":  {index: 2, text1: "2.", text2: "3."},
	} {
		t.Run(name, func(t *testing.T) {
			listSegment := &ListSegment{Ordered: true, Items: []RichTextSegment{
				&TextSegment{Text: "One"},
				&TextSegment{Text: "Two"},
			}}
			listSegment.SetStartNumber(tt.index)
			text := NewRichText(listSegment)
			texts := test.TempWidgetRenderer(t, text).Objects()

			assert.Equal(t, tt.text1, strings.TrimSpace(texts[0].(*canvas.Text).Text))
			assert.Equal(t, tt.text2, strings.TrimSpace(texts[2].(*canvas.Text).Text))
		})
	}
}

func TestRichText_List_WrappedIndent(t *testing.T) {
	text := NewRichText(&ListSegment{Items: []RichTextSegment{
		&TextSegment{Text: "A very long line that will wrap around and test the indentation of the second line"},
	}})
	text.Wrapping = fyne.TextWrapWord
	text.Resize(fyne.NewSize(100, 200))

	objs := test.TempWidgetRenderer(t, text).Objects()
	assert.GreaterOrEqual(t, len(objs), 3, "expected text to wrap into multiple segments")

	bullet := objs[0].(*canvas.Text)
	assert.Equal(t, "•", strings.TrimSpace(bullet.Text))

	line1 := objs[1].(*canvas.Text)
	line2 := objs[2].(*canvas.Text)

	assert.Equal(t, line1.Position().X, line2.Position().X, "wrapped line should align with start of first line")
	assert.Greater(t, line1.Position().X, bullet.Position().X, "text should be indented from bullet")
}
