package widget

import (
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"github.com/stretchr/testify/assert"
)

func TestRichTextMarkdown_Blockquote(t *testing.T) {
	r := NewRichTextFromMarkdown("p1\n\n> quote\n\np2")

	assert.Len(t, r.Segments, 6)
	if text, ok := r.Segments[2].(*TextSegment); ok {
		assert.Equal(t, "quote", text.Text)
		assert.Equal(t, 1, text.Style.QuotingDepth)
	} else {
		t.Error("Segment should be Text")
	}

	r = NewRichTextFromMarkdown("> # Head1\n> > # Head2")

	assert.Len(t, r.Segments, 4)
	if text, ok := r.Segments[0].(*TextSegment); ok {
		assert.Equal(t, "Head1", text.Text)
		assert.Equal(t, RichTextStyleHeading.SizeName, text.Style.SizeName)
		assert.Equal(t, true, text.Style.TextStyle.Bold)
		assert.Equal(t, true, text.Style.TextStyle.Italic)
		assert.Equal(t, 1, text.Style.QuotingDepth)
	} else {
		t.Error("Segment should be blockquoted Heading")
	}
	if text, ok := r.Segments[2].(*TextSegment); ok {
		assert.Equal(t, "Head2", text.Text)
		assert.Equal(t, RichTextStyleHeading.SizeName, text.Style.SizeName)
		assert.Equal(t, true, text.Style.TextStyle.Bold)
		assert.Equal(t, true, text.Style.TextStyle.Italic)
		assert.Equal(t, 2, text.Style.QuotingDepth)
	} else {
		t.Error("Segment should be a blockquoted Heading")
	}

	r = NewRichTextFromMarkdown("> - quoted list item")

	assert.Len(t, r.Segments, 1)
	if list, ok := r.Segments[0].(*ListSegment); ok {
		assert.Len(t, list.Items[0].(*ParagraphSegment).Texts, 1)
		assert.Equal(t, "quoted list item", list.Items[0].(*ParagraphSegment).Texts[0].(*TextSegment).Text)
		assert.Equal(t, 1, list.Items[0].(*ParagraphSegment).Texts[0].(*TextSegment).Style.QuotingDepth)
		assert.Equal(t, 1, list.quotingLevel)
	} else {
		t.Error("Segment should be a List")
	}

	r = NewRichTextFromMarkdown("> [fyne.io](https://fyne.io/)")

	assert.Len(t, r.Segments, 2)
	if link, ok := r.Segments[0].(*HyperlinkSegment); ok {
		assert.Equal(t, "fyne.io", link.Text)
		assert.Equal(t, 1, link.quotingLevel)
	} else {
		t.Error("Segment should be a blockquoted hyperlink")
	}

	r = NewRichTextFromMarkdown("> ```go\n> package main\n> ```")

	assert.Len(t, r.Segments, 1)
	if block, ok := r.Segments[0].(*CodeBlockSegment); ok {
		assert.Equal(t, "package main", block.Text)
		assert.Equal(t, 1, block.quotingLevel)
	} else {
		t.Error("Segment should be a blockquoted code block")
	}
}

func TestRichTextMarkdown_NestedBlockquote(t *testing.T) {
	r := NewRichTextFromMarkdown("p1\n\n> quote\n> > nested quote\n\np2")

	assert.Len(t, r.Segments, 8)
	if text, ok := r.Segments[4].(*TextSegment); ok {
		assert.Equal(t, "nested quote", text.Text)
		assert.Equal(t, 2, text.Style.QuotingDepth)
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
	if code, ok := r.Segments[0].(*CodeBlockSegment); ok {
		assert.Equal(t, "code\nblock", code.Text)
	} else {
		t.Error("Segment should be CodeBlock")
	}
}

func TestRichTextMarkdown_CodeBlockScrolls(t *testing.T) {
	long := strings.Repeat("abcdefghij", 50) // one ~500-char line
	cb := newRichCodeBlock(long)
	test.TempWidgetRenderer(t, cb)
	min := cb.MinSize()

	assert.Less(t, min.Width, float32(200)) // scrolls rather than demanding full width
	assert.Greater(t, min.Height, float32(10))
}

func TestRichTextMarkdown_Table(t *testing.T) {
	r := NewRichTextFromMarkdown("| Feature | Value |\n| :--- | ---: |\n| **bold** | `code` |")

	assert.Len(t, r.Segments, 1)
	table, ok := r.Segments[0].(*TableSegment)
	if !ok {
		t.Fatal("Segment should be a Table")
	}

	assert.Len(t, table.Headers, 2)
	assert.Len(t, table.Rows, 1)
	assert.Len(t, table.Rows[0], 2)

	// Per-column alignment from the delimiter row.
	assert.Equal(t, fyne.TextAlignLeading, table.Alignments[0])
	assert.Equal(t, fyne.TextAlignTrailing, table.Alignments[1])

	// Header text.
	if text, ok := table.Headers[0][0].(*TextSegment); ok {
		assert.Equal(t, "Feature", text.Text)
	} else {
		t.Error("header cell should be Text")
	}

	// Inline formatting inside body cells is preserved.
	if text, ok := table.Rows[0][0][0].(*TextSegment); ok {
		assert.True(t, text.Style.TextStyle.Bold, "bold cell should be bold")
	} else {
		t.Error("bold cell should be Text")
	}
	if text, ok := table.Rows[0][1][0].(*TextSegment); ok {
		assert.True(t, text.Style.TextStyle.Monospace, "code cell should be monospace")
	} else {
		t.Error("code cell should be Text")
	}
}

func TestRichTextMarkdown_TableAlignment(t *testing.T) {
	// The four GFM delimiter forms: left, center, right and none (no marker).
	r := NewRichTextFromMarkdown("| L | C | R | D |\n| :--- | :---: | ---: | --- |\n| a | b | c | d |")

	table, ok := r.Segments[0].(*TableSegment)
	if !ok {
		t.Fatal("Segment should be a Table")
	}

	want := []fyne.TextAlign{
		fyne.TextAlignLeading,  // :---
		fyne.TextAlignCenter,   // :---:
		fyne.TextAlignTrailing, // ---:
		fyne.TextAlignLeading,  // --- (no marker defaults to leading)
	}
	assert.Equal(t, want, table.Alignments)

	// alignFor reports each column's alignment and defaults to leading for
	// columns beyond the delimiter row.
	for col, a := range want {
		assert.Equal(t, a, table.alignFor(col))
	}
	assert.Equal(t, fyne.TextAlignLeading, table.alignFor(len(want)))
}

func TestTableCellAppliesAlignment(t *testing.T) {
	test.NewTempApp(t)

	// A body cell applies the column alignment to its text without bolding it.
	body := &TextSegment{Style: RichTextStyleInline, Text: "x"}
	newTableCell([]RichTextSegment{body}, fyne.TextAlignTrailing, false)
	assert.Equal(t, fyne.TextAlignTrailing, body.Style.Alignment)
	assert.False(t, body.Style.TextStyle.Bold)

	// A header cell applies the alignment and makes the text bold.
	header := &TextSegment{Style: RichTextStyleInline, Text: "h"}
	newTableCell([]RichTextSegment{header}, fyne.TextAlignCenter, true)
	assert.Equal(t, fyne.TextAlignCenter, header.Style.Alignment)
	assert.True(t, header.Style.TextStyle.Bold)

	// Hyperlink cells are aligned via their own Alignment field.
	link := &HyperlinkSegment{Text: "l"}
	newTableCell([]RichTextSegment{link}, fyne.TextAlignTrailing, false)
	assert.Equal(t, fyne.TextAlignTrailing, link.Alignment)
}

func TestRichTextMarkdown_TableAlignmentAppliedThroughVisual(t *testing.T) {
	test.NewTempApp(t)

	r := NewRichTextFromMarkdown("| L | R |\n| :--- | ---: |\n| a | b |")
	table := r.Segments[0].(*TableSegment)

	// Visual builds the cells, feeding alignFor(col) into each via newTableCell.
	table.Visual()

	assert.Equal(t, fyne.TextAlignLeading, table.Rows[0][0][0].(*TextSegment).Style.Alignment)
	assert.Equal(t, fyne.TextAlignTrailing, table.Rows[0][1][0].(*TextSegment).Style.Alignment)
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

	r.ParseMarkdown("~c~.")

	assert.Len(t, r.Segments, 3)
	if text, ok := r.Segments[0].(*TextSegment); ok {
		assert.Equal(t, "c", text.Text)
		assert.False(t, text.Style.TextStyle.Bold)
		assert.True(t, text.Style.TextStyle.Strikethrough)
	} else {
		t.Error("Segment should be text")
	}

	r.ParseMarkdown("~~d~~.")

	assert.Len(t, r.Segments, 3)
	if text, ok := r.Segments[0].(*TextSegment); ok {
		assert.Equal(t, "d", text.Text)
		assert.False(t, text.Style.TextStyle.Bold)
		assert.True(t, text.Style.TextStyle.Strikethrough)
	} else {
		t.Error("Segment should be text")
	}

	r.ParseMarkdown("~~**e**~~.")

	assert.Len(t, r.Segments, 3)
	if text, ok := r.Segments[0].(*TextSegment); ok {
		assert.Equal(t, "e", text.Text)
		assert.True(t, text.Style.TextStyle.Bold)
		assert.True(t, text.Style.TextStyle.Strikethrough)
	} else {
		t.Error("Segment should be text")
	}

	r.ParseMarkdown("**~~f~~**.")

	assert.Len(t, r.Segments, 3)
	if text, ok := r.Segments[0].(*TextSegment); ok {
		assert.Equal(t, "f", text.Text)
		assert.True(t, text.Style.TextStyle.Bold)
		assert.True(t, text.Style.TextStyle.Strikethrough)
	} else {
		t.Error("Segment should be text")
	}
}

func TestRichTextMarkdown_Heading(t *testing.T) {
	r := NewRichTextFromMarkdown("# Head1\n\n## Head2!\n\n### Head3\n\n## *Head4*\n\n## ~Head5~")

	assert.Len(t, r.Segments, 11)
	if text, ok := r.Segments[0].(*TextSegment); ok {
		assert.Equal(t, "Head1", text.Text)
		assert.Equal(t, RichTextStyleHeading, text.Style)
	} else {
		t.Error("Segment should be Heading")
	}
	if text, ok := r.Segments[2].(*TextSegment); ok {
		assert.Equal(t, "Head2", text.Text)
		assert.Equal(t, theme.SizeNameSubHeadingText, text.Style.SizeName)
	} else {
		t.Error("Segment should be SubHeading")
	}
	if text, ok := r.Segments[3].(*TextSegment); ok {
		assert.Equal(t, "!", text.Text)
		assert.Equal(t, RichTextStyleSubHeading, text.Style)
	} else {
		t.Error("Segment should be SubHeading")
	}

	if text, ok := r.Segments[5].(*TextSegment); ok {
		assert.Equal(t, "Head3", text.Text)
		assert.True(t, text.Style.TextStyle.Bold) // we don't have 6 levels of heading so just bold others
	} else {
		t.Error("Segment should be Strong")
	}

	if text, ok := r.Segments[7].(*TextSegment); ok {
		assert.Equal(t, "Head4", text.Text)
		assert.True(t, text.Style.TextStyle.Italic)
	} else {
		t.Error("Segment should be Italic")
	}

	if text, ok := r.Segments[9].(*TextSegment); ok {
		assert.Equal(t, "Head5", text.Text)
		assert.True(t, text.Style.TextStyle.Strikethrough)
	} else {
		t.Error("Segment should be Struck through")
	}
}

func TestRichTextMarkdown_Heading_Blank(t *testing.T) {
	r := NewRichTextFromMarkdown("#")

	assert.Len(t, r.Segments, 2)
	if text, ok := r.Segments[0].(*TextSegment); ok {
		assert.Equal(t, "", text.Text)
		assert.Equal(t, RichTextStyleHeading, text.Style)
	} else {
		t.Error("Segment should be Text")
	}

	r = NewRichTextFromMarkdown("# ")

	assert.Len(t, r.Segments, 2)
	if text, ok := r.Segments[0].(*TextSegment); ok {
		assert.Equal(t, "", text.Text)
		assert.Equal(t, RichTextStyleHeading, text.Style)
	} else {
		t.Error("Segment should be Text")
	}
}

func TestRichTextMarkdown_HeadingWithHyperlink(t *testing.T) {
	r := NewRichTextFromMarkdown("# Head1 [link](https://fyne.io/)\n\n## Head2 <https://fyne.io>\n### [Head3](https://fyne.io/)\n")

	assert.Len(t, r.Segments, 8)
	if text, ok := r.Segments[0].(*TextSegment); ok {
		assert.Equal(t, "Head1 ", text.Text)
		assert.Equal(t, RichTextStyleHeading, text.Style)
	} else {
		t.Error("Segment should be Heading text")
	}
	if link, ok := r.Segments[1].(*HyperlinkSegment); ok {
		assert.Equal(t, "link", link.Text)
		assert.Equal(t, RichTextStyleHeading.TextStyle, link.TextStyle)
	} else {
		t.Error("Segment should be Heading hyperlink")
	}
	if text, ok := r.Segments[2].(*TextSegment); ok {
		assert.Equal(t, "", text.Text)
		assert.Equal(t, RichTextStyleParagraph, text.Style)
	} else {
		t.Error("Segment should be Paragraph (linebreak)")
	}
	if link, ok := r.Segments[4].(*HyperlinkSegment); ok {
		assert.Equal(t, "https://fyne.io", link.Text)
		assert.Equal(t, RichTextStyleSubHeading.TextStyle, link.TextStyle)
	} else {
		t.Error("Segment should be SubHeading hyperlink")
	}
	if link, ok := r.Segments[6].(*HyperlinkSegment); ok {
		assert.Equal(t, "Head3", link.Text)
		assert.Equal(t, RichTextStyleStrong.TextStyle, link.TextStyle) // we don't have 6 levels of heading so just bold others
	} else {
		t.Error("Segment should be Strong hyperlink")
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

func TestRichTextMarkdown_HyperlinkInAngleBrackets(t *testing.T) {
	r := NewRichTextFromMarkdown("<https://fyne.io/>")

	assert.Len(t, r.Segments, 2)
	if link, ok := r.Segments[0].(*HyperlinkSegment); ok {
		assert.Equal(t, "https://fyne.io/", link.Text)
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

func TestRichTextMarkdown_TrailingExclamationMark(t *testing.T) {
	r := NewRichTextFromMarkdown("Hello! Hi!")

	assert.Len(t, r.Segments, 3)
	if text, ok := r.Segments[0].(*TextSegment); ok {
		assert.Equal(t, "Hello! Hi", text.Text)
		assert.True(t, text.Inline())
	} else {
		t.Error("Segment should be Text")
	}
	if text, ok := r.Segments[1].(*TextSegment); ok {
		assert.Equal(t, "!", text.Text)
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
				assert.Equal(t, tt.Start, list.StartNumber())
				assert.Equal(t, tt.startIndex, list.startIndex)
			} else {
				t.Error("Segment should be a List")
			}
		})
	}
}

func TestRichTextMarkdown_NestedList(t *testing.T) {
	r := NewRichTextFromMarkdown("* item 1\n  * item 1.1\n* item 2")

	assert.Len(t, r.Segments, 1)
	if list, ok := r.Segments[0].(*ListSegment); ok {
		assert.Len(t, list.Items, 3)
		assert.Len(t, list.Items[0].(*ParagraphSegment).Texts, 1)
		assert.Equal(t, "item 1", list.Items[0].(*ParagraphSegment).Texts[0].(*TextSegment).Text)
		if sublist, ok := list.Items[1].(*ListSegment); ok {
			assert.Len(t, sublist.Items, 1)
			assert.False(t, sublist.Ordered)
			assert.Equal(t, "item 1.1", sublist.Items[0].(*ParagraphSegment).Texts[0].(*TextSegment).Text)
		} else {
			t.Error("Segment should be a List")
		}
		assert.Len(t, list.Items[1].(*ListSegment).Items, 1)
	} else {
		t.Error("Segment should be a List")
	}

	r.ParseMarkdown("1. item 1\n   - item 1.1\n2. item 2")

	assert.Len(t, r.Segments, 1)
	if list, ok := r.Segments[0].(*ListSegment); ok {
		assert.True(t, list.Ordered)
		assert.Len(t, list.Items, 3)
		assert.Len(t, list.Items[0].(*ParagraphSegment).Texts, 1)
		if sublist, ok := list.Items[1].(*ListSegment); ok {
			assert.Len(t, sublist.Items, 1)
			assert.False(t, sublist.Ordered)
			assert.Equal(t, "item 1.1", sublist.Items[0].(*ParagraphSegment).Texts[0].(*TextSegment).Text)
		} else {
			t.Error("Segment should be a List")
		}
		assert.Equal(t, "item 2", list.Items[2].(*ParagraphSegment).Texts[0].(*TextSegment).Text)
	} else {
		t.Error("Segment should be a List")
	}
}

func TestRichTextMarkdown_Separator(t *testing.T) {
	r := NewRichTextFromMarkdown("---\n")

	assert.Len(t, r.Segments, 1)
	if _, ok := r.Segments[0].(*SeparatorSegment); !ok {
		t.Error("Segment should be a separator")
	}
}

func TestRichTextMarkdown_TextEncode(t *testing.T) {
	r := NewRichTextFromMarkdown("a&mdash;b")

	assert.GreaterOrEqual(t, len(r.Segments), 1)
	if _, ok := r.Segments[0].(*TextSegment); !ok {
		t.Error("Segment should be a text segment")
	}

	assert.Equal(t, "a—b", r.Segments[0].(*TextSegment).Text)
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
