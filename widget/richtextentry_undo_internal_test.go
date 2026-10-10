package widget

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

const undoTestMarkdown = "# Title\n\nSome **bold** text\n\n- one\n- two\n\n```\ncode\n```\n\nend"

func typeKey(e *RichTextEntry, key fyne.KeyName, count int) {
	for i := 0; i < count; i++ {
		e.TypedKey(&fyne.KeyEvent{Name: key})
	}
}

func undo(e *RichTextEntry) {
	e.TypedShortcut(&fyne.ShortcutUndo{})
}

func redo(e *RichTextEntry) {
	e.TypedShortcut(&fyne.ShortcutRedo{})
}

func TestRichTextEntry_Undo_DeletedStyle(t *testing.T) {
	e := NewRichTextEntryFromMarkdown("a **bold** c")
	e.setCursorOffset(6)

	typeKey(e, fyne.KeyBackspace, 4)
	assert.Equal(t, "a  c", e.Text)

	undo(e)
	assert.Equal(t, "a bold c", e.Text)
	assert.Equal(t, []string{"a |", "bold|b", " c|"}, segmentDump(e))
	assert.Equal(t, "a **bold** c", e.Markdown())
	assert.Equal(t, 6, e.CursorTextOffset())

	redo(e)
	assert.Equal(t, "a  c", e.Text)
	assert.Equal(t, 2, e.CursorTextOffset())

	undo(e)
	assert.Equal(t, "a **bold** c", e.Markdown())
}

func TestRichTextEntry_Undo_DeleteKey(t *testing.T) {
	e := NewRichTextEntryFromMarkdown("a **bold** c")
	e.setCursorOffset(2)

	typeKey(e, fyne.KeyDelete, 4)
	assert.Equal(t, "a  c", e.Text)

	// as in an Entry, each use of the delete key is a change of its own
	undo(e)
	assert.Equal(t, "a **d** c", e.Markdown())
	for i := 0; i < 3; i++ {
		undo(e)
	}
	assert.Equal(t, "a **bold** c", e.Markdown())
	assert.Equal(t, 2, e.CursorTextOffset())
}

func TestRichTextEntry_Undo_SelectAll(t *testing.T) {
	e := NewRichTextEntryFromMarkdown(undoTestMarkdown)
	segments := segmentDump(e)
	text := e.Text

	e.TypedShortcut(&fyne.ShortcutSelectAll{})
	typeKey(e, fyne.KeyBackspace, 1)
	assert.Equal(t, "", e.Text)

	undo(e)
	assert.Equal(t, text, e.Text)
	assert.Equal(t, segments, segmentDump(e))
	assert.Equal(t, undoTestMarkdown, e.Markdown())

	redo(e)
	assert.Equal(t, "", e.Text)

	undo(e)
	assert.Equal(t, text, e.Text)
	assert.Equal(t, segments, segmentDump(e))
	assert.Equal(t, undoTestMarkdown, e.Markdown())
}

func TestRichTextEntry_Undo_SelectAllReplace(t *testing.T) {
	e := NewRichTextEntryFromMarkdown(undoTestMarkdown)

	e.TypedShortcut(&fyne.ShortcutSelectAll{})
	test.Type(e, "new")
	assert.Equal(t, "new", e.Text)

	// the text that was replaced comes back with the typing taken out
	undo(e)
	assert.Equal(t, undoTestMarkdown, e.Markdown())

	redo(e)
	assert.Equal(t, "new", e.Text)
}

func TestRichTextEntry_Undo_Cut(t *testing.T) {
	e := NewRichTextEntryFromMarkdown("a **bold** c")
	selectRange(e, 2, 7)

	clipboard := test.NewClipboard()
	e.TypedShortcut(&fyne.ShortcutCut{Clipboard: clipboard})
	assert.Equal(t, "a c", e.Text)
	assert.Equal(t, "bold ", clipboard.Content())

	undo(e)
	assert.Equal(t, "a **bold** c", e.Markdown())

	e.setCursorOffset(8)
	e.TypedShortcut(&fyne.ShortcutPaste{Clipboard: clipboard})
	assert.Equal(t, "a bold cbold ", e.Text)

	undo(e)
	assert.Equal(t, "a **bold** c", e.Markdown())
	assert.Equal(t, 8, e.CursorTextOffset())

	redo(e)
	assert.Equal(t, "a bold cbold ", e.Text)
	assert.Equal(t, 13, e.CursorTextOffset())
}

func TestRichTextEntry_Undo_Typing(t *testing.T) {
	e := NewRichTextEntryFromMarkdown("a **b** c")
	e.setCursorOffset(3)

	test.Type(e, "old")
	assert.Equal(t, []string{"a |", "bold|b", " c|"}, segmentDump(e))

	// a word is typed as a single change
	undo(e)
	assert.Equal(t, []string{"a |", "b|b", " c|"}, segmentDump(e))
	assert.Equal(t, 3, e.CursorTextOffset())
	assert.False(t, e.undoStack.CanUndo())

	redo(e)
	assert.Equal(t, []string{"a |", "bold|b", " c|"}, segmentDump(e))
	assert.Equal(t, 6, e.CursorTextOffset())
	assert.False(t, e.undoStack.CanRedo())
}

func TestRichTextEntry_Undo_TypingAfterUndo(t *testing.T) {
	e := NewRichTextEntry()
	test.Type(e, "one")
	undo(e)
	assert.Equal(t, "", e.Text)

	test.Type(e, "two")
	assert.False(t, e.undoStack.CanRedo())
	undo(e)
	assert.Equal(t, "", e.Text)
	redo(e)
	assert.Equal(t, "two", e.Text)
}

func TestRichTextEntry_Undo_Markdown(t *testing.T) {
	e := NewRichTextEntry()
	e.TypeMarkdown = true

	test.Type(e, "a **b**")
	assert.Equal(t, "a b", e.Text)
	assert.Equal(t, "a **b**", e.Markdown())

	// the markdown that was typed is there again once the style is undone
	undo(e)
	assert.Equal(t, "a **b*", e.Text)
	assert.Equal(t, []string{"a **b*|"}, segmentDump(e))

	redo(e)
	assert.Equal(t, "a b", e.Text)
	assert.Equal(t, "a **b**", e.Markdown())

	test.Type(e, " c")
	assert.Equal(t, "a **b** c", e.Markdown())
	undo(e)
	assert.Equal(t, "a **b**", e.Markdown())
}

func TestRichTextEntry_Undo_Heading(t *testing.T) {
	e := NewRichTextEntry()
	e.TypeMarkdown = true

	test.Type(e, "# Title\nbody")
	assert.Equal(t, "# Title\n\nbody", e.Markdown())

	undo(e)
	assert.Equal(t, "# Title", e.Markdown())

	// text typed on the line after a heading is not part of it
	test.Type(e, "\nmore")
	assert.Equal(t, "# Title\n\nmore", e.Markdown())
}

func TestRichTextEntry_Undo_List(t *testing.T) {
	e := NewRichTextEntryFromMarkdown("- one\n- two")
	e.setCursorOffset(3)

	typeKey(e, fyne.KeyReturn, 1)
	assert.Equal(t, "- one\n- \n- two", e.Markdown())

	undo(e)
	assert.Equal(t, "- one\n- two", e.Markdown())
	assert.Equal(t, 3, e.CursorTextOffset())

	redo(e)
	assert.Equal(t, "- one\n- \n- two", e.Markdown())
	assert.Equal(t, 4, e.CursorTextOffset())
}

func TestRichTextEntry_Undo_TypingStyle(t *testing.T) {
	e := NewRichTextEntry()
	test.Type(e, "ab")
	e.setCursorOffset(1)

	e.TypedShortcut(styleShortcut(fyne.KeyB))
	test.Type(e, "X")
	assert.Equal(t, "a**X**b", e.Markdown())

	undo(e)
	assert.Equal(t, "ab", e.Markdown())
	assert.True(t, e.StyleAtCursor().TextStyle.Bold)

	// choosing the style is not a change that can be seen, so it is not a step of its own
	undo(e)
	assert.Equal(t, "", e.Text)

	redo(e)
	assert.Equal(t, "ab", e.Markdown())
	redo(e)
	assert.Equal(t, "a**X**b", e.Markdown())
}

func TestRichTextEntry_Undo_SetStyle(t *testing.T) {
	e := NewRichTextEntry()
	e.SetText("hello world")

	style := RichTextStyleInline
	style.TextStyle.Italic = true
	e.SetStyleForRange(0, 5, style)
	assert.Equal(t, "*hello* world", e.Markdown())

	undo(e)
	assert.Equal(t, "hello world", e.Markdown())
	redo(e)
	assert.Equal(t, "*hello* world", e.Markdown())
}

func TestRichTextEntry_Undo_ContentReplaced(t *testing.T) {
	e := NewRichTextEntry()
	test.Type(e, "one")

	e.AppendMarkdown("**two**")
	test.Type(e, "!")

	// the history does not reach back past content that was set from code
	undo(e)
	assert.Equal(t, "one\n\n**two**", e.Markdown())
	undo(e)
	assert.Equal(t, "one\n\n**two**", e.Markdown())

	e.ParseMarkdown("*new*")
	undo(e)
	assert.Equal(t, "*new*", e.Markdown())
}

func TestRichTextEntry_Undo_OnChanged(t *testing.T) {
	e := NewRichTextEntry()
	test.Type(e, "one")

	var changed []string
	e.OnChanged = func(s string) { changed = append(changed, s) }
	undo(e)
	redo(e)
	redo(e) // nothing more to apply
	assert.Equal(t, []string{"", "one"}, changed)
}

func TestRichTextEntry_Undo_Disabled(t *testing.T) {
	e := NewRichTextEntryFromMarkdown("a **b** c")
	e.setCursorOffset(3)
	test.Type(e, "X")

	e.Disable()
	e.Enable()
	undo(e)
	assert.Equal(t, "a **b** c", e.Markdown())
}

// TestRichTextEntry_Undo_Random makes edits at random, checking that the history
// follows each of them and that it leads back to the content that was there before.
func TestRichTextEntry_Undo_Random(t *testing.T) {
	test.NewTempApp(t)
	keys := []fyne.KeyName{fyne.KeyBackspace, fyne.KeyDelete, fyne.KeyReturn, fyne.KeyLeft, fyne.KeyUp}
	runes := []rune("ab c*_#->1. ~")
	clipboard := test.NewClipboard()
	clipboard.SetContent("pa**st**e\nd")

	for seed := int64(0); seed < 50; seed++ {
		r := rand.New(rand.NewSource(seed)) //nolint:gosec // the same edits are wanted on every run
		e := NewRichTextEntryFromMarkdown(undoTestMarkdown)
		e.TypeMarkdown = seed%2 == 0
		start := cloneSegments(e.Segments())

		for i := 0; i < 80; i++ {
			switch r.Intn(10) {
			case 0, 1, 2:
				e.TypedRune(runes[r.Intn(len(runes))])
			case 3, 4:
				typeKey(e, keys[r.Intn(len(keys))], 1)
			case 5:
				length := len([]rune(e.Text))
				selectRange(e, r.Intn(length+1), r.Intn(length+1))
			case 6:
				e.TypedShortcut(styleShortcuts[r.Intn(len(styleShortcuts))].shortcut)
			case 7:
				e.TypedShortcut(&fyne.ShortcutPaste{Clipboard: clipboard})
			case 8:
				undo(e)
			case 9:
				redo(e)
			}

			require.True(t, sameSegments(e.undoState, e.Segments()), "seed %d edit %d was not recorded", seed, i)
			for _, action := range e.undoStack.items {
				require.IsType(t, &richTextUndoAction{}, action, "seed %d edit %d", seed, i)
			}
		}

		for e.undoStack.CanRedo() {
			redo(e)
		}
		end := cloneSegments(e.Segments())

		for e.undoStack.CanUndo() {
			undo(e)
		}
		require.True(t, sameSegments(start, e.Segments()), "seed %d undo gave %q", seed, e.Markdown())

		for e.undoStack.CanRedo() {
			redo(e)
		}
		require.True(t, sameSegments(end, e.Segments()), "seed %d redo gave %q", seed, e.Markdown())
	}
}

func TestDiffSegments(t *testing.T) {
	bold := RichTextStyleInline
	bold.TextStyle.Bold = true
	before := []RichTextSegment{
		&TextSegment{Style: RichTextStyleInline, Text: "first "},
		&TextSegment{Style: bold, Text: "some long text"},
		&TextSegment{Style: RichTextStyleInline, Text: " last"},
	}
	after := cloneSegments(before)
	after[1].(*TextSegment).Text = "some very long text"

	change := diffSegments(before, after)
	assert.Equal(t, 1, change.at)
	assert.Equal(t, 5, change.head)
	assert.Equal(t, 9, change.tail)
	assert.Equal(t, "", change.before[0].Textual())
	assert.Equal(t, "very ", change.after[0].Textual())

	assert.True(t, sameSegments(after, change.apply(before)))
	assert.True(t, sameSegments(before, change.revert(after)))
	assert.Equal(t, "some long text", before[1].Textual())
}

func TestDiffSegments_Split(t *testing.T) {
	bold := RichTextStyleInline
	bold.TextStyle.Bold = true
	before := []RichTextSegment{&TextSegment{Style: RichTextStyleInline, Text: "hello world"}}
	after := []RichTextSegment{
		&TextSegment{Style: RichTextStyleInline, Text: "hello "},
		&TextSegment{Style: bold, Text: "big"},
		&TextSegment{Style: RichTextStyleInline, Text: " world"},
	}

	change := diffSegments(before, after)
	assert.True(t, sameSegments(after, change.apply(before)))
	assert.True(t, sameSegments(before, change.revert(after)))

	change = diffSegments(after, before)
	assert.True(t, sameSegments(before, change.apply(after)))
	assert.True(t, sameSegments(after, change.revert(before)))
}

func TestDiffSegments_Repeated(t *testing.T) {
	before := []RichTextSegment{&TextSegment{Style: RichTextStyleInline, Text: "aaa"}}
	after := []RichTextSegment{&TextSegment{Style: RichTextStyleInline, Text: "aaaa"}}

	change := diffSegments(before, after)
	assert.True(t, sameSegments(after, change.apply(before)))
	assert.True(t, sameSegments(before, change.revert(after)))

	change = diffSegments(after, before)
	assert.True(t, sameSegments(before, change.apply(after)))
	assert.True(t, sameSegments(after, change.revert(before)))
}
