package widget

import (
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
)

func TestEntry_UndoMerge(t *testing.T) {
	text := "abc def"
	entry, fakeTime := newEntryWithFakeTimestamper()

	// Successive TypedRune() events should be merged into a single action.
	entry.SetText(text)
	entry.CursorColumn = 3
	typeToEntry(entry, "XYZ", fakeTime)
	entry.Undo()
	assert.Equal(t, text, entry.Text)

	// Not all TypedRune events are merged into a single action.
	entry.SetText(text)
	entry.CursorColumn = 3
	typeToEntry(entry, " NO", fakeTime)
	fakeTime.add(10 * time.Second)
	typeToEntry(entry, "YES ", fakeTime) // only this should be cancelled due to a large delay
	entry.Undo()
	assert.Equal(t, "abc NO def", entry.Text)
	assert.Equal(t, 1, entry.redoOffset)
	assert.Equal(t, 2, len(entry.actionLog))

	// Sequence of deletions is also merged into a single action.
	entry.SetText(text)
	entry.CursorColumn = 3
	for i := 0; i < 3; i++ {
		entry.TypedKey(&fyne.KeyEvent{
			Name: fyne.KeyBackspace,
		})
	}
	entry.Undo()
	assert.Equal(t, text, entry.Text)

	// Pressing Enter is merged with TypedRune() actions
	entry.SetText(text)
	entry.MultiLine = true
	entry.CursorColumn = 3
	entry.TypedKey(&fyne.KeyEvent{
		Name: fyne.KeyEnter,
	})
	typeToEntry(entry, "Second line", fakeTime)
	entry.Undo()
	assert.Equal(t, text, entry.Text)
}

func TestEntry_UndoUndoables(t *testing.T) {
	text := "abc def"
	entry, fakeTime := newEntryWithFakeTimestamper()

	// SetTextUndoable() is undoable.
	entry.SetText("")
	typeToEntry(entry, text, fakeTime)
	entry.SetTextUndoable("WOW")
	entry.CursorRow = 0
	entry.CursorColumn = 3
	typeToEntry(entry, " OWO", fakeTime)
	entry.Undo()
	entry.Undo()
	assert.Equal(t, text, entry.Text)
	entry.Undo()
	assert.Equal(t, "", entry.Text)

	// Cutting is undoable.
	clipboard := test.NewClipboard()
	entry.SetText(text)
	entry.selectAll()
	entry.TypedShortcut(&fyne.ShortcutCut{
		Clipboard: clipboard,
	})
	entry.Undo()
	assert.Equal(t, text, entry.Text)

	// Pasting is undoable.
	entry.SetText(text)
	clipboard.SetContent("qwerty")
	entry.CursorColumn = 7
	entry.TypedShortcut(&fyne.ShortcutPaste{
		Clipboard: clipboard,
	})
	entry.Undo()
	assert.Equal(t, text, entry.Text)
}

func TestEntry_Undo(t *testing.T) {
	text := "abc def"
	entry, _ := newEntryWithFakeTimestamper()

	// Empty action history should result in no changes (and no crashes).
	entry.SetText(text)
	assert.Equal(t, false, entry.CanUndo())
	assert.Equal(t, false, entry.CanRedo())
	entry.Undo()
	entry.Undo()
	entry.Redo()
	assert.False(t, entry.CanUndo())
	assert.False(t, entry.CanRedo())
	assert.Equal(t, text, entry.Text)
	assert.Equal(t, 0, entry.redoOffset)
	assert.Equal(t, 0, len(entry.actionLog))

	// Basic Undo operation.
	entry.SetText(text)
	entry.CursorRow = 0
	entry.CursorColumn = 3
	test.Type(entry, "Z")
	entry.Undo()
	assert.Equal(t, text, entry.Text)
	assert.Equal(t, 1, entry.redoOffset)
	assert.False(t, entry.CanUndo())
	assert.True(t, entry.CanRedo())

	// Basic Redo operation.
	entry.Redo()
	assert.Equal(t, "abcZ def", entry.Text)
	assert.Equal(t, 0, entry.redoOffset)
}

func clickEntry(e *Entry, ev *fyne.PointEvent) {
	mouseEvent := &desktop.MouseEvent{PointEvent: *ev, Button: desktop.MouseButtonPrimary}
	e.MouseDown(mouseEvent)
	e.MouseUp(mouseEvent)
	e.Tapped(ev)
}

func typeToEntry(e *Entry, text string, fakeTime *fakeTime) {
	for c := range text {
		e.TypedRune(rune(text[c]))
		if fakeTime != nil {
			fakeTime.add(time.Millisecond)
		}
	}
}

type fakeTime struct {
	now time.Time
}

func (ft *fakeTime) add(delta time.Duration) {
	ft.now = ft.now.Add(delta)
}

func newEntryWithFakeTimestamper() (*Entry, *fakeTime) {
	entry := NewEntry()
	entry.Wrapping = fyne.TextWrapOff
	entry.HistoryDisabled = false
	entry.Refresh()

	fakeTimestamper := fakeTime{
		now: time.Unix(1633720000, 0),
	}

	entry.timestamper = func() time.Time {
		return fakeTimestamper.now
	}

	return entry, &fakeTimestamper
}
