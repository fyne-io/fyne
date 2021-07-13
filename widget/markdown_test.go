package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

	r = NewRichTextFromMarkdown("``` go\ncode\nblock\n```")
	assert.Equal(t, 1, len(r.Segments))
	if text, ok := r.Segments[0].(*TextSegment); ok {
		assert.Equal(t, "code\nblock", text.Text)
		assert.Equal(t, RichTextStyleCodeBlock, text.Style)
	} else {
		t.Error("Segment should be Text")
	}
}

func TestRichTextMarkdown_Heading(t *testing.T) {
	r := NewRichTextFromMarkdown("# Head1\n\n## Head2\n")

	assert.Equal(t, 2, len(r.Segments))
	if text, ok := r.Segments[0].(*TextSegment); ok {
		assert.Equal(t, "Head1", text.Text)
		assert.Equal(t, RichTextStyleHeading, text.Style)
	} else {
		t.Error("Segment should be Text")
	}
	if text, ok := r.Segments[1].(*TextSegment); ok {
		assert.Equal(t, "Head2", text.Text)
		assert.Equal(t, RichTextStyleSubHeading, text.Style)
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

func TestRichTextMarkdown_Lines(t *testing.T) {
	r := NewRichTextFromMarkdown("line1\nline2\n") // a single newline is not a new paragraph

	assert.Equal(t, 2, len(r.Segments))
	if text, ok := r.Segments[0].(*TextSegment); ok {
		assert.Equal(t, "line1", text.Text)
		assert.True(t, text.Inline())
	} else {
		t.Error("Segment should be Text")
	}
	if text, ok := r.Segments[1].(*TextSegment); ok {
		assert.Equal(t, "line2", text.Text)
	} else {
		t.Error("Segment should be Text")
	}

	r = NewRichTextFromMarkdown("line1\n\nline2\n")

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

func TestRichTextMarkdown_Separator(t *testing.T) {
	r := NewRichTextFromMarkdown("---\n")

	assert.Equal(t, 1, len(r.Segments))
	if _, ok := r.Segments[0].(*SeparatorSegment); !ok {
		t.Error("Segment should be a separator")
	}
}
