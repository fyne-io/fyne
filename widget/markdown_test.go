package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2/storage"
)

func TestRichTextMarkdown_Blockquote(t *testing.T) {
	r := NewRichTextFromMarkdown("p1\n\n> quote\n\np2")

	assert.Equal(t, 3, len(r.Segments))
	if text, ok := r.Segments[1].(*TextSegment); ok {
		assert.Equal(t, "quote", text.Text)
		assert.Equal(t, RichTextStyleBlockquote, text.Style)
	} else {
		t.Error("Segment should be Text")
	}
}

func TestRichTextMarkdown_Code(t *testing.T) {
	r := NewRichTextFromMarkdown("a `code` inline")

	assert.Equal(t, 3, len(r.Segments))
	if text, ok := r.Segments[1].(*TextSegment); ok {
		assert.Equal(t, "code", text.Text)
		assert.Equal(t, RichTextStyleCodeInline, text.Style)
	} else {
		t.Error("Segment should be Text")
	}

	r.ParseMarkdown("``` go\ncode\nblock\n```")
	assert.Equal(t, 1, len(r.Segments))
	if text, ok := r.Segments[0].(*TextSegment); ok {
		assert.Equal(t, "code\nblock", text.Text)
		assert.Equal(t, RichTextStyleCodeBlock, text.Style)
	} else {
		t.Error("Segment should be Text")
	}
}

func TestRichTextMarkdown_Code_Incomplete(t *testing.T) {
	r := NewRichTextFromMarkdown("` ")

	assert.Equal(t, 1, len(r.Segments))
	if text, ok := r.Segments[0].(*TextSegment); ok {
		assert.Equal(t, "`", text.Text)
		assert.Equal(t, RichTextStyleParagraph, text.Style)
	} else {
		t.Error("Segment should be Text")
	}

	r.ParseMarkdown("``` ")
	assert.Equal(t, 0, len(r.Segments))

	r.ParseMarkdown("~~~ ")
	assert.Equal(t, 0, len(r.Segments))
}

func TestRichTextMarkdown_Emphasis(t *testing.T) {
	r := NewRichTextFromMarkdown("*a*")

	assert.Equal(t, 1, len(r.Segments))
	if text, ok := r.Segments[0].(*TextSegment); ok {
		assert.Equal(t, "a", text.Text)
		assert.True(t, text.Style.TextStyle.Italic)
	} else {
		t.Error("Segment should be text")
	}

	r.ParseMarkdown("**b**.")

	assert.Equal(t, 2, len(r.Segments))
	if text, ok := r.Segments[0].(*TextSegment); ok {
		assert.Equal(t, "b", text.Text)
		assert.True(t, text.Style.TextStyle.Bold)
	} else {
		t.Error("Segment should be text")
	}
}

func TestRichTextMarkdown_Heading(t *testing.T) {
	r := NewRichTextFromMarkdown("# Head1\n\n## Head2!\n### Head3\n")

	assert.Equal(t, 3, len(r.Segments))
	if text, ok := r.Segments[0].(*TextSegment); ok {
		assert.Equal(t, "Head1", text.Text)
		assert.Equal(t, RichTextStyleHeading, text.Style)
	} else {
		t.Error("Segment should be Heading")
	}
	if text, ok := r.Segments[1].(*TextSegment); ok {
		assert.Equal(t, "Head2!", text.Text)
		assert.Equal(t, RichTextStyleSubHeading, text.Style)
	} else {
		t.Error("Segment should be SubHeading")
	}

	if text, ok := r.Segments[2].(*TextSegment); ok {
		assert.Equal(t, "Head3", text.Text)
		assert.Equal(t, true, text.Style.TextStyle.Bold) // we don't have 6 levels of heading so just bold others
	} else {
		t.Error("Segment should be Strong")
	}
}

func TestRichTextMarkdown_Heading_Blank(t *testing.T) {
	r := NewRichTextFromMarkdown("#")

	assert.Equal(t, 1, len(r.Segments))
	if text, ok := r.Segments[0].(*TextSegment); ok {
		assert.Equal(t, "#", text.Text)
		assert.Equal(t, RichTextStyleParagraph, text.Style)
	} else {
		t.Error("Segment should be Text")
	}

	r = NewRichTextFromMarkdown("# ")

	assert.Equal(t, 1, len(r.Segments))
	if text, ok := r.Segments[0].(*TextSegment); ok {
		assert.Equal(t, "", text.Text)
		assert.Equal(t, RichTextStyleHeading, text.Style)
	} else {
		t.Error("Segment should be Text")
	}
}

func TestRichTextMarkdown_Hyperlink(t *testing.T) {
	r := NewRichTextFromMarkdown("[title](https://fyne.io/)")

	assert.Equal(t, 1, len(r.Segments))
	if link, ok := r.Segments[0].(*HyperlinkSegment); ok {
		assert.Equal(t, "title", link.Text)
		assert.Equal(t, "fyne.io", link.URL.Host)
	} else {
		t.Error("Segment should be a Hyperlink")
	}
}

func TestRichTextMarkdown_Image(t *testing.T) {
	r := NewRichTextFromMarkdown("![title](../../theme/icons/fyne.png)")

	assert.Equal(t, 1, len(r.Segments))
	if img, ok := r.Segments[0].(*ImageSegment); ok {
		assert.Equal(t, storage.NewFileURI("../../theme/icons/fyne.png"), img.Source)
	} else {
		t.Error("Segment should be a Image")
	}
}

func TestRichTextMarkdown_Lines(t *testing.T) {
	r := NewRichTextFromMarkdown("line1\nline2\n") // a single newline is not a new paragraph

	assert.Equal(t, 2, len(r.Segments))
	if text, ok := r.Segments[0].(*TextSegment); ok {
		assert.Equal(t, "line1 ", text.Text)
		assert.True(t, text.Inline())
	} else {
		t.Error("Segment should be Text")
	}
	if text, ok := r.Segments[1].(*TextSegment); ok {
		assert.Equal(t, "line2", text.Text)
	} else {
		t.Error("Segment should be Text")
	}

	r.ParseMarkdown("line1\n\nline2\n")

	assert.Equal(t, 2, len(r.Segments))
	if text, ok := r.Segments[0].(*TextSegment); ok {
		assert.Equal(t, "line1", text.Text)
		assert.False(t, text.Inline())
	} else {
		t.Error("Segment should be Text")
	}
	if text, ok := r.Segments[1].(*TextSegment); ok {
		assert.Equal(t, "line2", text.Text)
	} else {
		t.Error("Segment should be Text")
	}
}

func TestRichTextMarkdown_List(t *testing.T) {
	r := NewRichTextFromMarkdown("* line1 in _three_ segments\n* line2")

	assert.Equal(t, 1, len(r.Segments))
	if list, ok := r.Segments[0].(*ListSegment); ok {
		assert.Equal(t, 2, len(list.Items))
		assert.Equal(t, 3, len(list.Items[0].(*ParagraphSegment).Texts))
		assert.Equal(t, "line1 in ", list.Items[0].(*ParagraphSegment).Texts[0].(*TextSegment).Text)
	} else {
		t.Error("Segment should be a List")
	}

	r.ParseMarkdown("1. line1\n2. line2")

	assert.Equal(t, 1, len(r.Segments))
	if list, ok := r.Segments[0].(*ListSegment); ok {
		assert.True(t, list.Ordered)
		assert.Equal(t, 2, len(list.Items))
		assert.Equal(t, 1, len(list.Items[1].(*ParagraphSegment).Texts))
		assert.Equal(t, "line2", list.Items[1].(*ParagraphSegment).Texts[0].(*TextSegment).Text)
	} else {
		t.Error("Segment should be a List")
	}
}

func TestRichTextMarkdown_Separator(t *testing.T) {
	r := NewRichTextFromMarkdown("---\n")

	assert.Equal(t, 1, len(r.Segments))
	if _, ok := r.Segments[0].(*SeparatorSegment); !ok {
		t.Error("Segment should be a separator")
	}
}
