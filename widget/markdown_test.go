package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2/storage"
)

func TestRichTextMarkdown_Blockquote(t *testing.T) {
	r := NewRichTextFromMarkdown("p1\n\n> quote\n\np2")

	assert.Len(t, r.Segments, 5)
	if text, ok := r.Segments[2].(*TextSegment); ok {
		assert.Equal(t, "quote", text.Text)
		assert.Equal(t, RichTextStyleBlockquote, text.Style)
	} else {
		t.Error("Segment should be Text")
	}
}

func TestRichTextMarkdown_Code(t *testing.T) {
	r := NewRichTextFromMarkdown("a `code` inline")

	assert.Len(t, r.Segments, 4)
	if text, ok := r.Segments[1].(*TextSegment); ok {
		assert.Equal(t, "code", text.Text)
		assert.Equal(t, RichTextStyleCodeInline, text.Style)
	} else {
		t.Error("Segment should be Text")
	}

	r.ParseMarkdown("``` go\ncode\nblock\n```")
	assert.Len(t, r.Segments, 1)
	if text, ok := r.Segments[0].(*TextSegment); ok {
		assert.Equal(t, "code\nblock", text.Text)
		assert.Equal(t, RichTextStyleCodeBlock, text.Style)
	} else {
		t.Error("Segment should be Text")
	}
}

func TestRichTextMarkdown_Code_Incomplete(t *testing.T) {
	r := NewRichTextFromMarkdown("` ")

	assert.Len(t, r.Segments, 2)
	if text, ok := r.Segments[0].(*TextSegment); ok {
		assert.Equal(t, "`", text.Text)
		assert.Equal(t, RichTextStyleInline, text.Style)
	} else {
		t.Error("Segment should be Text")
	}

	r.ParseMarkdown("``` ")
	assert.Empty(t, r.Segments)

	r.ParseMarkdown("~~~ ")
	assert.Empty(t, r.Segments)
}

func TestRichTextMarkdown_Emphasis(t *testing.T) {
	r := NewRichTextFromMarkdown("*a*")

	assert.Len(t, r.Segments, 2)
	if text, ok := r.Segments[0].(*TextSegment); ok {
		assert.Equal(t, "a", text.Text)
		assert.True(t, text.Style.TextStyle.Italic)
	} else {
		t.Error("Segment should be text")
	}

	r.ParseMarkdown("**b**.")

	assert.Len(t, r.Segments, 3)
	if text, ok := r.Segments[0].(*TextSegment); ok {
		assert.Equal(t, "b", text.Text)
		assert.True(t, text.Style.TextStyle.Bold)
	} else {
		t.Error("Segment should be text")
	}
}

func TestRichTextMarkdown_Heading(t *testing.T) {
	r := NewRichTextFromMarkdown("# Head1\n\n## Head2!\n### Head3\n")

	assert.Len(t, r.Segments, 3)
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
		assert.True(t, text.Style.TextStyle.Bold) // we don't have 6 levels of heading so just bold others
	} else {
		t.Error("Segment should be Strong")
	}
}

func TestRichTextMarkdown_Heading_Blank(t *testing.T) {
	r := NewRichTextFromMarkdown("#")

	assert.Len(t, r.Segments, 1)
	if text, ok := r.Segments[0].(*TextSegment); ok {
		assert.Equal(t, "", text.Text)
		assert.Equal(t, RichTextStyleHeading, text.Style)
	} else {
		t.Error("Segment should be Text")
	}

	r = NewRichTextFromMarkdown("# ")

	assert.Len(t, r.Segments, 1)
	if text, ok := r.Segments[0].(*TextSegment); ok {
		assert.Equal(t, "", text.Text)
		assert.Equal(t, RichTextStyleHeading, text.Style)
	} else {
		t.Error("Segment should be Text")
	}
}

func TestRichTextMarkdown_Hyperlink(t *testing.T) {
	r := NewRichTextFromMarkdown("[title](https://fyne.io/)")

	assert.Len(t, r.Segments, 2)
	if link, ok := r.Segments[0].(*HyperlinkSegment); ok {
		assert.Equal(t, "title", link.Text)
		assert.Equal(t, "fyne.io", link.URL.Host)
	} else {
		t.Error("Segment should be a Hyperlink")
	}
}

func TestRichTextMarkdown_Image(t *testing.T) {
	r := NewRichTextFromMarkdown("![title](../../theme/icons/fyne.png)")

	assert.Len(t, r.Segments, 2)
	if img, ok := r.Segments[0].(*ImageSegment); ok {
		assert.Equal(t, storage.NewFileURI("../../theme/icons/fyne.png"), img.Source)
	} else {
		t.Error("Segment should be a Image")
	}

	r = NewRichTextFromMarkdown("![](../../theme/icons/fyne.png)")

	assert.Len(t, r.Segments, 2)
	if img, ok := r.Segments[0].(*ImageSegment); ok {
		assert.Equal(t, storage.NewFileURI("../../theme/icons/fyne.png"), img.Source)
	} else {
		t.Error("Segment should be a Image")
	}
}

func TestRichTextMarkdown_Lines(t *testing.T) {
	r := NewRichTextFromMarkdown("line1\nline2\n") // a single newline is not a new paragraph

	assert.Len(t, r.Segments, 3)
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

	assert.Len(t, r.Segments, 4)
	if text, ok := r.Segments[0].(*TextSegment); ok {
		assert.Equal(t, "line1", text.Text)
		assert.True(t, text.Inline())
	} else {
		t.Error("Segment should be Text")
	}
	if text, ok := r.Segments[2].(*TextSegment); ok {
		assert.Equal(t, "line2", text.Text)
	} else {
		t.Error("Segment should be Text")
	}
}

func TestRichTextMarkdown_List(t *testing.T) {
	r := NewRichTextFromMarkdown("* line1 in _three_ segments\n* line2")

	assert.Len(t, r.Segments, 1)
	if list, ok := r.Segments[0].(*ListSegment); ok {
		assert.Len(t, list.Items, 2)
		assert.Len(t, list.Items[0].(*ParagraphSegment).Texts, 3)
		assert.Equal(t, "line1 in ", list.Items[0].(*ParagraphSegment).Texts[0].(*TextSegment).Text)
	} else {
		t.Error("Segment should be a List")
	}

	r.ParseMarkdown("1. line1\n2. line2")

	assert.Len(t, r.Segments, 1)
	if list, ok := r.Segments[0].(*ListSegment); ok {
		assert.True(t, list.Ordered)
		assert.Len(t, list.Items, 2)
		assert.Len(t, list.Items[1].(*ParagraphSegment).Texts, 1)
		assert.Equal(t, "line2", list.Items[1].(*ParagraphSegment).Texts[0].(*TextSegment).Text)
	} else {
		t.Error("Segment should be a List")
	}
}

func TestRichTextMarkdown_ListWithDifferentStartingIndex(t *testing.T) {
	for name, tt := range map[string]struct {
		text       string
		Start      int
		startIndex int
	}{
		"Start at 2": {"2. line1\n3. line2", 2, 1},
		"Start at 1": {"1. line1\n2. line2", 1, 0},
		"Start at 0": {"0. line1\n1. line2", 0, -1},
	} {
		t.Run(name, func(t *testing.T) {
			r := NewRichTextFromMarkdown(tt.text)

			assert.Len(t, r.Segments, 1)
			if list, ok := r.Segments[0].(*ListSegment); ok {
				assert.Equal(t, tt.Start, list.Start())
				assert.Equal(t, tt.startIndex, list.startIndex)
			} else {
				t.Error("Segment should be a List")
			}
		})
	}
}

func TestRichTextMarkdown_Separator(t *testing.T) {
	r := NewRichTextFromMarkdown("---\n")

	assert.Len(t, r.Segments, 1)
	if _, ok := r.Segments[0].(*SeparatorSegment); !ok {
		t.Error("Segment should be a separator")
	}
}

func TestRichTextMarkdown_Paragraph(t *testing.T) {
	r := NewRichTextFromMarkdown("foo")

	assert.Len(t, r.Segments, 2)
	if text, ok := r.Segments[1].(*TextSegment); ok {
		assert.Equal(t, "", text.Text)
		assert.Equal(t, RichTextStyleParagraph, text.Style)
	} else {
		t.Error("Last segment of paragraph should be empty text")
	}
}

func TestRichTextMarkdown_NewlinesAroundEmphasis(t *testing.T) {
	r := NewRichTextFromMarkdown("foo\n*bar*\nbaz")
	assert.Equal(t, "foo bar baz", r.String())
}

func TestRichTextMarkdown_NewlinesAroundStrong(t *testing.T) {
	r := NewRichTextFromMarkdown("foo\n**bar**\nbaz")
	assert.Equal(t, "foo bar baz", r.String())
}

func TestRichTextMarkdown_NewlinesAroundHyperlink(t *testing.T) {
	r := NewRichTextFromMarkdown("foo\n[bar](https://fyne.io/)\nbaz")
	assert.Equal(t, "foo bar baz", r.String())
}

func TestRichTextMarkdown_SpacesAroundHyperlink(t *testing.T) {
	r := NewRichTextFromMarkdown("foo [bar](https://fyne.io/) baz")
	assert.Equal(t, "foo bar baz", r.String())
}

func TestRichTextMarkdown_NewlineInHyperlink(t *testing.T) {
	r := NewRichTextFromMarkdown("[foo\nbar](https://fyne.io/)")
	assert.Equal(t, "foo bar", r.String())
}

func BenchmarkMarkdownParsing(b *testing.B) {
	md := `# Test heading
This is some test markdown. It contains some different markdown
elements in an effort to put a realistic load onto the parser.

> Here is a quote.
> It stretches over two lines and contains some ` + "`inline code`" + `.

Moreover, we've got a list:
- foo
- bar
- baz
`
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		NewRichTextFromMarkdown(md)
	}
}
