package widget

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
)

// segmentDump renders the segments as "text|style" pairs so that tests can
// assert on the structure without depending on exact style values.
func segmentDump(e *RichTextEntry) []string {
	var out []string
	for _, seg := range e.Segments() {
		text, ok := seg.(*TextSegment)
		if !ok {
			out = append(out, seg.Textual()+"|object")
			continue
		}

		style := ""
		if text.Style.TextStyle.Bold {
			style += "b"
		}
		if text.Style.TextStyle.Italic {
			style += "i"
		}
		if text.Style.TextStyle.Monospace {
			style += "m"
		}
		if text.Style.TextStyle.Strikethrough {
			style += "s"
		}
		switch text.Style.SizeName {
		case theme.SizeNameHeadingText:
			style += "h1"
		case theme.SizeNameSubHeadingText:
			style += "h2"
		}
		out = append(out, text.Text+"|"+style)
	}
	return out
}

func typeString(e *RichTextEntry, s string) {
	for _, r := range s {
		e.TypedRune(r)
	}
}

func TestRichTextEntry_ParseMarkdown(t *testing.T) {
	e := NewRichTextEntryFromMarkdown("# Title\n\nSome **bold** and *italic* text.")

	assert.Equal(t, "Title\nSome bold and italic text.", e.Text)
	assert.Equal(t, []string{"Title|bh1", "\nSome |", "bold|b", " and |", "italic|i", " text.|"},
		segmentDump(e))
}

func TestRichTextEntry_ParseMarkdown_Blocks(t *testing.T) {
	e := NewRichTextEntryFromMarkdown("- one\n- two\n\n> quoted")

	// the bullets are drawn by the list, they are not part of the content
	assert.Equal(t, "one\ntwo\nquoted", e.Text)
}

func TestRichTextEntry_SetText(t *testing.T) {
	e := NewRichTextEntryFromMarkdown("# Title")
	e.SetText("plain")

	assert.Equal(t, "plain", e.Text)
	assert.Equal(t, []string{"plain|"}, segmentDump(e))
}

func TestRichTextEntry_TypedRune_KeepsStyles(t *testing.T) {
	e := NewRichTextEntryFromMarkdown("a **b** c")
	e.CursorRow, e.CursorColumn = 0, 3 // just after the bold "b"

	e.TypedRune('X')

	assert.Equal(t, "a bX c", e.Text)
	assert.Equal(t, []string{"a |", "bX|b", " c|"}, segmentDump(e))
	assert.Equal(t, 4, e.CursorColumn)
}

func TestRichTextEntry_TypedRune_AtSegmentStart(t *testing.T) {
	e := NewRichTextEntryFromMarkdown("**b** c")
	e.CursorRow, e.CursorColumn = 0, 0

	e.TypedRune('X')

	assert.Equal(t, "Xb c", e.Text)
	assert.Equal(t, []string{"Xb|b", " c|"}, segmentDump(e))
}

func TestRichTextEntry_Backspace_AcrossSegments(t *testing.T) {
	e := NewRichTextEntryFromMarkdown("a **b** c")
	e.CursorRow, e.CursorColumn = 0, 3

	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyBackspace})

	assert.Equal(t, "a  c", e.Text)
	// the emptied bold segment stays at the cursor, so typing continues in bold
	assert.Equal(t, []string{"a |", "|b", " c|"}, segmentDump(e))
	assert.Equal(t, 2, e.CursorColumn)
}

func TestRichTextEntry_Rows(t *testing.T) {
	e := NewRichTextEntryFromMarkdown("# Title\n\nBody text")
	provider := e.textProvider()

	assert.Equal(t, 2, provider.rows())
	assert.Equal(t, "Title", string(provider.row(0)))
	assert.Equal(t, "Body text", string(provider.row(1)))

	// the heading row is taller than the body row
	_, headingHeight := provider.rowGeometry(0)
	bodyY, bodyHeight := provider.rowGeometry(1)
	assert.Greater(t, headingHeight, bodyHeight)
	assert.Equal(t, headingHeight, bodyY)
}

func TestRichTextEntry_CursorOffsetAcrossRows(t *testing.T) {
	e := NewRichTextEntryFromMarkdown("# Title\n\nBody")

	e.CursorRow, e.CursorColumn = 1, 2
	assert.Equal(t, 8, e.CursorTextOffset()) // "Title\n" is 6 runes

	row, col := e.rowColFromTextPos(8)
	assert.Equal(t, 1, row)
	assert.Equal(t, 2, col)
}

func TestRichTextEntry_SetStyleForRange(t *testing.T) {
	e := NewRichTextEntry()
	e.SetText("hello world")

	style := RichTextStyleInline
	style.TextStyle.Bold = true
	e.SetStyleForRange(6, 11, style)

	assert.Equal(t, "hello world", e.Text)
	assert.Equal(t, []string{"hello |", "world|b"}, segmentDump(e))
}

func TestRichTextEntry_SetStyleForSelection(t *testing.T) {
	e := NewRichTextEntry()
	e.SetText("hello world")
	e.syncSelectable()
	e.sel.selecting = true
	e.sel.selectRow, e.sel.selectColumn = 0, 0
	e.CursorRow, e.CursorColumn = 0, 5
	e.syncSelectable()

	style := RichTextStyleInline
	style.TextStyle.Italic = true
	e.SetStyleForSelection(style)

	assert.Equal(t, []string{"hello|i", " world|"}, segmentDump(e))
}

func TestRichTextEntry_MarkdownMode_Bold(t *testing.T) {
	e := NewRichTextEntry()
	e.TypeMarkdown = true

	typeString(e, "a **b** c")

	assert.Equal(t, "a b c", e.Text)
	assert.Equal(t, []string{"a |", "b|b", " c|"}, segmentDump(e))
}

func TestRichTextEntry_MarkdownMode_Italic(t *testing.T) {
	e := NewRichTextEntry()
	e.TypeMarkdown = true

	typeString(e, "*hi* there")

	assert.Equal(t, "hi there", e.Text)
	assert.Equal(t, []string{"hi|i", " there|"}, segmentDump(e))
}

func TestRichTextEntry_MarkdownMode_Code(t *testing.T) {
	e := NewRichTextEntry()
	e.TypeMarkdown = true

	typeString(e, "run `go test` now")

	assert.Equal(t, "run go test now", e.Text)
	assert.Equal(t, []string{"run |", "go test|m", " now|"}, segmentDump(e))
}

func TestRichTextEntry_MarkdownMode_Heading(t *testing.T) {
	e := NewRichTextEntry()
	e.TypeMarkdown = true

	typeString(e, "# Title")

	assert.Equal(t, "Title", e.Text)
	assert.Equal(t, []string{"Title|bh1"}, segmentDump(e))
}

func TestRichTextEntry_MarkdownMode_NoFalseMatch(t *testing.T) {
	e := NewRichTextEntry()
	e.TypeMarkdown = true

	typeString(e, "2 * 3 * 4")

	assert.Equal(t, "2 * 3 * 4", e.Text)
	assert.Equal(t, []string{"2 * 3 * 4|"}, segmentDump(e))
}

func TestRichTextEntry_Markdown(t *testing.T) {
	source := "# Title\n\nSome **bold** and *italic* text."
	e := NewRichTextEntryFromMarkdown(source)

	assert.Equal(t, source, e.Markdown())
}

func TestRichTextEntry_Undo(t *testing.T) {
	e := NewRichTextEntryFromMarkdown("a **b** c")
	e.CursorRow, e.CursorColumn = 0, 3

	e.TypedRune('X')
	assert.Equal(t, "a bX c", e.Text)

	e.Undo()
	assert.Equal(t, "a b c", e.Text)
	assert.Equal(t, []string{"a |", "b|b", " c|"}, segmentDump(e))

	e.Redo()
	assert.Equal(t, "a bX c", e.Text)
	assert.Equal(t, []string{"a |", "bX|b", " c|"}, segmentDump(e))
}

func TestRichTextEntry_Selection(t *testing.T) {
	e := NewRichTextEntryFromMarkdown("a **b** c")
	e.syncSelectable()
	e.sel.selecting = true
	e.sel.selectRow, e.sel.selectColumn = 0, 0
	e.CursorRow, e.CursorColumn = 0, 4
	e.syncSelectable()

	assert.Equal(t, "a b ", e.SelectedText())
}

func TestRichTextEntry_Renders(t *testing.T) {
	test.NewTempApp(t)

	e := NewRichTextEntryFromMarkdown("# Title\n\nBody **text**")
	w := test.NewTempWindow(t, e)
	w.Resize(fyne.NewSize(200, 100))

	e.FocusGained()
	e.TypedRune('!')
	assert.Equal(t, "!Title\nBody text", e.Text)
}

func TestRichTextEntry_MouseRowAcrossHeights(t *testing.T) {
	test.NewTempApp(t)

	e := NewRichTextEntryFromMarkdown("# Heading\n\nBody\n\nMore")
	w := test.NewTempWindow(t, e)
	w.Resize(fyne.NewSize(300, 200))

	provider := e.textProvider()
	e.syncSelectable()
	innerPad := e.Theme().Size(theme.SizeNameInnerPadding)
	lineSpace := e.Theme().Size(theme.SizeNameLineSpacing)

	for row := 0; row < provider.rows(); row++ {
		y, height := provider.rowGeometry(row)
		mid := y + height/2 + innerPad - lineSpace
		got, _ := e.sel.getRowCol(fyne.NewPos(1, mid))
		assert.Equal(t, row, got, "row %d at y %v", row, mid)
	}
}

func TestRichTextEntry_MouseColumnUsesSegmentSize(t *testing.T) {
	test.NewTempApp(t)

	e := NewRichTextEntryFromMarkdown("# Heading")
	w := test.NewTempWindow(t, e)
	w.Resize(fyne.NewSize(300, 200))
	e.syncSelectable()

	th := e.Theme()
	provider := e.textProvider()
	// clicking just past the fourth character of the large heading text
	x := provider.lineSizeToColumn(4, 0, th.Size(theme.SizeNameText), th.Size(theme.SizeNameInnerPadding)).Width
	_, col := e.sel.getRowCol(fyne.NewPos(x+1, 2))
	assert.Equal(t, 4, col)
}

func TestRichTextEntry_ReturnLeavesHeading(t *testing.T) {
	e := NewRichTextEntryFromMarkdown("# Title")
	e.CursorRow, e.CursorColumn = 0, 5

	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	typeString(e, "body")

	assert.Equal(t, "Title\nbody", e.Text)
	assert.Equal(t, []string{"Title\n|bh1", "body|"}, segmentDump(e))
}

func TestRichTextEntry_ReturnKeepsInlineStyle(t *testing.T) {
	e := NewRichTextEntryFromMarkdown("**bold**")
	e.CursorRow, e.CursorColumn = 0, 4

	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	typeString(e, "more")

	assert.Equal(t, "bold\nmore", e.Text)
	assert.Equal(t, []string{"bold\nmore|b"}, segmentDump(e))
}

func TestRichTextEntry_Wrapping(t *testing.T) {
	test.NewTempApp(t)

	e := NewRichTextEntry()
	e.SetText("one two three four five six seven eight nine ten")
	w := test.NewTempWindow(t, e)
	w.Resize(fyne.NewSize(120, 200))
	e.Refresh()

	provider := e.textProvider()
	assert.Greater(t, provider.rows(), 1)

	// every row must map back to the offset that it starts at
	for row := 0; row < provider.rows(); row++ {
		start := textPosFromRowCol(row, 0, provider)
		gotRow, gotCol := e.rowColFromTextPos(start)
		assert.Equal(t, 0, gotCol, "row %d", row)
		assert.Equal(t, row, gotRow, "row %d", row)
	}
}

func TestRichTextEntry_AppendMarkdown(t *testing.T) {
	e := NewRichTextEntryFromMarkdown("first")
	e.AppendMarkdown("\n\n**second**")

	assert.Equal(t, "first\nsecond", e.Text)
	assert.Equal(t, []string{"first\n|", "second|b"}, segmentDump(e))
}

func TestRichTextEntry_MarkdownRoundTrip(t *testing.T) {
	source := "# Title\n\nSome **bold** text\n\n## Second\n\nWith `code` and *emphasis*"
	e := NewRichTextEntryFromMarkdown(source)
	again := NewRichTextEntryFromMarkdown(e.Markdown())

	assert.Equal(t, e.Text, again.Text)
	assert.Equal(t, segmentDump(e), segmentDump(again))
}

func TestRichTextEntry_Disabled(t *testing.T) {
	e := NewRichTextEntryFromMarkdown("**bold** text")
	e.Disable()
	e.Refresh()

	for _, seg := range e.Segments() {
		text := seg.(*TextSegment)
		assert.Equal(t, theme.ColorNameDisabled, text.Style.ColorName)
	}
	assert.True(t, e.Segments()[0].(*TextSegment).Style.TextStyle.Bold)

	e.Enable()
	e.Refresh()
	assert.Equal(t, theme.ColorNameForeground, e.Segments()[0].(*TextSegment).Style.ColorName)
}

func TestRichTextEntry_QuoteIndent(t *testing.T) {
	test.NewTempApp(t)

	e := NewRichTextEntryFromMarkdown("body\n\n> quoted\n\nafter")
	w := test.NewTempWindow(t, e)
	w.Resize(fyne.NewSize(300, 200))
	e.Refresh()

	th := e.Theme()
	provider := e.textProvider()
	textSize, innerPad := th.Size(theme.SizeNameText), th.Size(theme.SizeNameInnerPadding)
	x := func(row int) float32 {
		return provider.lineSizeToColumn(0, row, textSize, innerPad).Width
	}

	// the quote is indented, the lines either side of it are not
	assert.Equal(t, x(0), x(2), "the line after a quote must not inherit its indent")
	assert.Greater(t, x(1), x(0), "a quoted line is indented")
}

func TestRichTextEntry_DeleteWordAcrossRows(t *testing.T) {
	e := NewRichTextEntryFromMarkdown("# Title\n\nalpha beta gamma")
	e.CursorRow, e.CursorColumn = e.rowColFromTextPos(e.richProvider().len())
	e.syncSelectable()

	e.deleteWord(false)

	assert.Equal(t, "Title\nalpha beta ", e.Text)
}

// itemTexts returns the text of each item in the list segment at the given index.
func itemTexts(t *testing.T, e *RichTextEntry, index int) []string {
	t.Helper()

	list, ok := e.Segments()[index].(*ListSegment)
	if !ok {
		t.Fatalf("segment %d is a %T, not a list", index, e.Segments()[index])
	}

	var out []string
	for _, item := range list.Items {
		text := &strings.Builder{}
		for _, seg := range appendContentSegments([]RichTextSegment{item}, nil) {
			text.WriteString(seg.Textual())
		}
		out = append(out, text.String())
	}
	return out
}

func TestRichTextEntry_ParseMarkdown_List(t *testing.T) {
	e := NewRichTextEntryFromMarkdown("- one\n- two")

	// the list is kept, so that it draws the bullets that are not content
	assert.Equal(t, "one\ntwo", e.Text)
	assert.Len(t, e.Segments(), 1)
	assert.Equal(t, []string{"one\n", "two"}, itemTexts(t, e, 0))
}

func TestRichTextEntry_ListReturnAddsItem(t *testing.T) {
	e := NewRichTextEntryFromMarkdown("- one\n- two")
	e.CursorRow, e.CursorColumn = 0, 3 // the end of the first item

	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	typeString(e, "new")

	assert.Equal(t, "one\nnew\ntwo", e.Text)
	assert.Equal(t, []string{"one\n", "new\n", "two"}, itemTexts(t, e, 0))
	assert.Equal(t, 1, e.CursorRow)
}

func TestRichTextEntry_ListReturnSplitsItem(t *testing.T) {
	e := NewRichTextEntryFromMarkdown("- one two")
	e.CursorRow, e.CursorColumn = 0, 3

	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})

	assert.Equal(t, "one\n two", e.Text)
	assert.Equal(t, []string{"one\n", " two"}, itemTexts(t, e, 0))
}

func TestRichTextEntry_ListReturnOnEmptyItemLeaves(t *testing.T) {
	e := NewRichTextEntryFromMarkdown("- one")
	e.CursorRow, e.CursorColumn = 0, 3

	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn}) // an item to type the next entry in
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn}) // nothing was typed, so the list ends
	typeString(e, "after")

	assert.Equal(t, "one\nafter", e.Text)
	assert.Equal(t, []string{"one\n"}, itemTexts(t, e, 0))
	assert.Equal(t, "after", e.Segments()[1].Textual()) // outside of the list
}

func TestRichTextEntry_ListBackspaceJoinsItems(t *testing.T) {
	e := NewRichTextEntryFromMarkdown("- one\n- two")
	e.CursorRow, e.CursorColumn = 1, 0 // the start of the second item

	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyBackspace})

	assert.Equal(t, "onetwo", e.Text)
	assert.Equal(t, []string{"onetwo"}, itemTexts(t, e, 0))
}

func TestRichTextEntry_ListBackspaceLeavesFirstItem(t *testing.T) {
	e := NewRichTextEntryFromMarkdown("- one\n- two")
	e.CursorRow, e.CursorColumn = 0, 0

	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyBackspace})

	assert.Equal(t, "one\ntwo", e.Text)
	assert.Equal(t, "one\n", e.Segments()[0].Textual()) // the item is out of the list now
	assert.Equal(t, []string{"two"}, itemTexts(t, e, 1))
}

func TestRichTextEntry_ListDeleteSelectionJoinsItems(t *testing.T) {
	e := NewRichTextEntryFromMarkdown("- one\n- two\n- three")
	e.CursorRow, e.CursorColumn = 0, 2
	e.sel.selectRow, e.sel.selectColumn, e.sel.selecting = 0, 2, true
	e.CursorRow, e.CursorColumn = 2, 3
	e.syncSelectable()

	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyDelete})

	assert.Equal(t, "onee", e.Text)
	assert.Equal(t, []string{"onee"}, itemTexts(t, e, 0))
}

func TestRichTextEntry_MarkdownMode_List(t *testing.T) {
	e := NewRichTextEntry()
	e.TypeMarkdown = true

	typeString(e, "- one")
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	typeString(e, "two")

	assert.Equal(t, "one\ntwo\n", e.Text) // the bullets are drawn, not typed
	assert.Equal(t, []string{"one\n", "two\n"}, itemTexts(t, e, 0))
}

func TestRichTextEntry_MarkdownMode_OrderedList(t *testing.T) {
	e := NewRichTextEntry()
	e.TypeMarkdown = true

	typeString(e, "1. first")
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	typeString(e, "second")

	list := e.Segments()[0].(*ListSegment)
	assert.True(t, list.Ordered)
	assert.Equal(t, []string{"first\n", "second\n"}, itemTexts(t, e, 0))

	list.Segments() // the numbers follow the order of the items
	assert.Equal(t, "1. ", list.markers[0].marker())
	assert.Equal(t, "2. ", list.markers[1].marker())
}

func TestRichTextEntry_MarkdownMode_CodeFence(t *testing.T) {
	e := NewRichTextEntry()
	e.TypeMarkdown = true

	typeString(e, "```")
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	typeString(e, "x := 1")

	code, ok := e.Segments()[0].(*CodeBlockSegment)
	assert.True(t, ok, "typing a fence opens a code block")
	assert.Equal(t, "x := 1\n", code.Text)

	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	typeString(e, "```")
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	typeString(e, "after")

	assert.Equal(t, "x := 1\n", code.Text) // the fence closed the block
	assert.Equal(t, "x := 1\nafter", e.Text)
	assert.Equal(t, "after", e.Segments()[1].Textual())
}

func TestRichTextEntry_CodeBlockEdits(t *testing.T) {
	e := NewRichTextEntryFromMarkdown("intro\n\n```\ncode\n```\n\nafter")
	assert.Equal(t, "intro\ncode\nafter", e.Text)

	code := e.Segments()[1].(*CodeBlockSegment)
	e.CursorRow, e.CursorColumn = 1, 4 // the end of the code

	typeString(e, "X")
	assert.Equal(t, "codeX\n", code.Text)

	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyBackspace})
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyBackspace})
	assert.Equal(t, "cod\n", code.Text)
	assert.Equal(t, "intro\ncod\nafter", e.Text)
}

func TestRichTextEntry_BlocksLayOutInRows(t *testing.T) {
	test.NewTempApp(t)

	e := NewRichTextEntryFromMarkdown("intro\n\n- one\n- two\n\n```\ncode\nhere\n```\n\nafter")
	w := test.NewTempWindow(t, e)
	w.Resize(fyne.NewSize(300, 300))
	e.Refresh()

	provider := e.richProvider()
	var rows []string
	var panels []bool
	for i := 0; i < provider.rows(); i++ {
		rows = append(rows, string(provider.row(i)))
		panels = append(panels, provider.rowBounds[i].panel != nil)
	}

	// the bullets are not part of any row, and the code lines sit on the panel
	assert.Equal(t, []string{"intro", "one", "two", "code", "here", "", "after"}, rows)
	assert.Equal(t, []bool{false, false, false, true, true, true, false}, panels)
}

func TestRichTextEntry_ListItemWrapIndents(t *testing.T) {
	test.NewTempApp(t)

	e := NewRichTextEntryFromMarkdown("- an item with enough words in it to wrap onto a second line")
	w := test.NewTempWindow(t, e)
	w.Resize(fyne.NewSize(150, 200))
	e.Refresh()

	provider := e.richProvider()
	assert.Greater(t, provider.rows(), 1, "the item should wrap")
	assert.Greater(t, provider.rowBounds[1].indent, float32(0), "a wrapped item lines up with its text")
}

func TestRichTextEntry_ListDeleteJoinsItems(t *testing.T) {
	e := NewRichTextEntryFromMarkdown("- one\n- two")
	e.CursorRow, e.CursorColumn = 0, 3 // the end of the first item

	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyDelete})

	assert.Equal(t, "onetwo", e.Text)
	assert.Equal(t, []string{"onetwo"}, itemTexts(t, e, 0))
	assert.Equal(t, 3, e.CursorColumn)
}

func TestRichTextEntry_BackspaceRemovesEmptyCodeBlock(t *testing.T) {
	test.NewTempApp(t)

	e := NewRichTextEntry()
	e.TypeMarkdown = true
	w := test.NewTempWindow(t, e)
	w.Resize(fyne.NewSize(250, 200))

	typeString(e, "intro")
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	typeString(e, "```")
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	assert.IsType(t, &CodeBlockSegment{}, e.Segments()[1], "the fence opened a block")

	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyBackspace})
	typeString(e, "x")

	// the block is gone, leaving the line it was on to be typed in plain text
	assert.Equal(t, "intro\nx", e.Text)
	assert.Equal(t, []string{"intro\n|", "x|"}, segmentDump(e))
	assert.Nil(t, e.richProvider().rowBounds[1].panel)
}

func TestRichTextEntry_BackspaceKeepsCodeBlockWithContent(t *testing.T) {
	e := NewRichTextEntryFromMarkdown("intro\n\n```\ncode\n```\n\nafter")
	e.CursorRow, e.CursorColumn = 1, 4

	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn}) // an empty line inside the block
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyBackspace})

	code, ok := e.Segments()[1].(*CodeBlockSegment)
	assert.True(t, ok, "a block with code in it stays")
	assert.Equal(t, "code\n", code.Text)
	assert.Equal(t, "intro\ncode\nafter", e.Text)
}

func TestRichTextEntry_BackspaceClearsEmptyStyle(t *testing.T) {
	e := NewRichTextEntryFromMarkdown("intro\n\n# Title")
	e.CursorRow, e.CursorColumn = 1, 5
	for range "Title" {
		e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyBackspace})
	}
	assert.Equal(t, []string{"intro\n|", "|bh1"}, segmentDump(e))

	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyBackspace}) // the empty heading loses its style
	assert.Equal(t, "intro\n", e.Text)
	assert.Equal(t, []string{"intro\n|", "|"}, segmentDump(e))

	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyBackspace}) // and then the line goes as usual
	assert.Equal(t, "intro", e.Text)
}

func TestRichTextEntry_DeleteClearsEmptyStyle(t *testing.T) {
	e := NewRichTextEntry()
	e.TypeMarkdown = true

	typeString(e, "> ")
	assert.True(t, e.StyleAtCursor().QuotingDepth > 0)

	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyDelete})
	typeString(e, "plain")

	assert.Equal(t, "plain", e.Text)
	assert.Equal(t, RichTextStyleInline, e.StyleAtCursor())
}

func TestRichTextEntry_BackspaceJoinsBlankLineToItem(t *testing.T) {
	test.NewTempApp(t)

	e := NewRichTextEntry()
	e.TypeMarkdown = true
	w := test.NewTempWindow(t, e)
	w.Resize(fyne.NewSize(250, 200))

	typeString(e, "- one")
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn}) // an empty line below the list
	assert.Equal(t, 1, e.CursorRow)

	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyBackspace})
	typeString(e, "X")

	// the blank line is gone and the cursor carried on at the end of the item
	assert.Equal(t, "oneX", e.Text)
	assert.Equal(t, []string{"oneX"}, itemTexts(t, e, 0))
	assert.Equal(t, 0, e.CursorRow)
	assert.Equal(t, 1, e.richProvider().rows())
}

func TestRichTextEntry_BackspaceJoinsLineToItem(t *testing.T) {
	test.NewTempApp(t)

	e := NewRichTextEntryFromMarkdown("- one\n\nafter")
	w := test.NewTempWindow(t, e)
	w.Resize(fyne.NewSize(250, 200))
	e.CursorRow, e.CursorColumn = 1, 0

	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyBackspace})

	assert.Equal(t, "oneafter", e.Text)
	assert.Equal(t, []string{"oneafter"}, itemTexts(t, e, 0))
	assert.Equal(t, 0, e.CursorRow)
	assert.Equal(t, 3, e.CursorColumn)
}

func TestRichTextEntry_BackspaceJoinsLineToCodeBlock(t *testing.T) {
	test.NewTempApp(t)

	e := NewRichTextEntryFromMarkdown("```\ncode\n```\n\nafter")
	w := test.NewTempWindow(t, e)
	w.Resize(fyne.NewSize(250, 200))
	e.CursorRow, e.CursorColumn = 2, 0 // the line below the block

	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyBackspace})

	code, ok := e.Segments()[0].(*CodeBlockSegment)
	assert.True(t, ok)
	assert.Equal(t, "codeafter", code.Text) // the line joined the code above it
	assert.Equal(t, "codeafter", e.Text)
	assert.Equal(t, 4, e.CursorColumn)
}

func TestRichTextEntry_SelectionDrawsOverCodePanel(t *testing.T) {
	test.NewTempApp(t)

	e := NewRichTextEntryFromMarkdown("intro\n\n```\ncode\n```")
	w := test.NewTempWindow(t, e)
	w.Resize(fyne.NewSize(250, 200))
	e.FocusGained()
	e.selectAll()
	e.sel.Refresh()

	provider := e.richProvider()
	assert.True(t, len(provider.decor) > 0, "the text draws the selection over its panel")
	assert.Empty(t, cache.Renderer(e.sel).Objects(), "so the selection widget draws nothing")

	panels, highlights, texts := -1, -1, -1
	for i, obj := range cache.Renderer(provider).Objects() {
		switch o := obj.(type) {
		case *richCodeBlock:
			panels = i
		case *canvas.Rectangle:
			if highlights == -1 {
				highlights = i
			}
		case *canvas.Text:
			if texts == -1 && o.Text != "" {
				texts = i
			}
		}
	}

	// the panel is behind the highlight, which is in turn behind the text
	assert.NotEqual(t, -1, panels)
	assert.NotEqual(t, -1, highlights)
	assert.NotEqual(t, -1, texts)
	assert.Less(t, panels, highlights)
	assert.Less(t, highlights, texts)
}

func TestRichTextEntry_SelectionEndsWithoutHighlights(t *testing.T) {
	test.NewTempApp(t)

	e := NewRichTextEntryFromMarkdown("```\ncode\n```")
	w := test.NewTempWindow(t, e)
	w.Resize(fyne.NewSize(250, 200))
	e.FocusGained()
	e.selectAll()
	e.sel.Refresh()
	assert.NotEmpty(t, e.richProvider().highlightObjects())

	e.ClearSelection()
	assert.Empty(t, e.richProvider().highlightObjects(), "an ended selection leaves nothing drawn")
}

func TestRichTextEntry_MarkdownBlocks(t *testing.T) {
	for _, source := range []string{
		"- one\n- two",
		"1. first\n2. second",
		"- one\n    - inner\n- two",
		"intro\n\n- one\n- two\n\nafter",
		"```\ncode here\n```",
		"intro\n\n```\ncode\nlines\n```\n\nafter",
		"> quoted\n\n> - a\n> - b",
		"# Title\n\n- item with **bold**\n\n```\nx := 1\n```",
	} {
		e := NewRichTextEntryFromMarkdown(source)
		assert.Equal(t, source, e.Markdown(), "round trip of %q", source)
	}
}

func TestRichTextEntry_MarkdownOfTypedBlocks(t *testing.T) {
	e := NewRichTextEntry()
	e.TypeMarkdown = true

	typeString(e, "Notes")
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	typeString(e, "- one")
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	typeString(e, "two with **bold**")
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn}) // an empty item closes the list
	typeString(e, "```")
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	typeString(e, "x := 1")
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	typeString(e, "```")
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	typeString(e, "after")

	assert.Equal(t, "Notes\n\n- one\n- two with **bold**\n\n```\nx := 1\n```\n\nafter", e.Markdown())
}

func TestRichTextEntry_MarkdownOfTypedOrderedList(t *testing.T) {
	e := NewRichTextEntry()
	e.TypeMarkdown = true

	typeString(e, "3. third")
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	typeString(e, "fourth")

	assert.Equal(t, "3. third\n4. fourth", e.Markdown())
}
