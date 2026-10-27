package widget

import (
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/test"
)

func TestEntry_TouchUp_AccountsForScroll(t *testing.T) {
	test.NewApp()
	e := NewMultiLineEntry()
	e.SetText("Line 1\nLine 2\nLine 3\nLine 4\nLine 5")

	e.CreateRenderer()
	e.scroller.Offset.Y = 50

	touch := &mobile.TouchEvent{
		PointEvent: fyne.PointEvent{
			Position: fyne.NewPos(10, 10),
		},
	}

	e.TouchDown(touch)
	e.TouchUp(touch)

	if e.CursorRow == 0 {
		t.Errorf("Expected CursorRow > 0, got %d", e.CursorRow)
	}
}

// A touch that becomes a scroll gets a touch down but never a touch up, and must leave the cursor alone.
func TestEntry_TouchDown_DoesNotMoveCursor(t *testing.T) {
	test.NewApp()
	e := NewMultiLineEntry()
	e.SetText("Line 1\nLine 2\nLine 3\nLine 4\nLine 5")
	e.CreateRenderer()

	e.CursorRow, e.CursorColumn = 2, 3
	e.Refresh()

	e.TouchDown(&mobile.TouchEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(10, 10)}})
	if e.CursorRow != 2 || e.CursorColumn != 3 {
		t.Errorf("Expected cursor to stay at 2,3, got %d,%d", e.CursorRow, e.CursorColumn)
	}
}

// Disabling wrapping detaches the scroller and must drop its offset, otherwise a
// later triple tap selects a lower line than the one under the pointer.
func TestEntry_WrapOff_ClearsScrollOffset(t *testing.T) {
	test.NewApp()
	e := NewMultiLineEntry()
	e.Wrapping = fyne.TextWrapWord
	e.SetText("Line 1\nLine 2\nLine 3\nLine 4\nLine 5")
	w := test.NewTempWindow(t, e)
	w.Resize(fyne.NewSize(200, 80))

	e.scroller.Offset.Y = 50
	e.Wrapping = fyne.TextWrapOff
	e.Scroll = fyne.ScrollNone
	e.Refresh()

	if e.scroller.Content != nil {
		t.Error("expected scroller content to be detached")
	}
	if e.scroller.Offset != fyne.NewPos(0, 0) {
		t.Errorf("expected zero scroll offset, got %v", e.scroller.Offset)
	}

	pos := fyne.NewPos(10, 10)
	e.MouseDown(&desktop.MouseEvent{
		PointEvent: fyne.PointEvent{Position: pos},
		Button:     desktop.MouseButtonPrimary,
	})
	if e.CursorRow != 0 {
		t.Errorf("Expected CursorRow 0 after mouse down, got %d", e.CursorRow)
	}

	touch := &mobile.TouchEvent{PointEvent: fyne.PointEvent{Position: pos}}
	e.TouchDown(touch)
	e.TouchUp(touch)
	if e.CursorRow != 0 {
		t.Errorf("Expected CursorRow 0 after touch, got %d", e.CursorRow)
	}

	e.DoubleTapped(&fyne.PointEvent{Position: pos})
	e.TouchDown(touch)
	if got := e.SelectedText(); got != "Line 1" {
		t.Errorf("Expected selected text %q, got %q", "Line 1", got)
	}

	e.Wrapping = fyne.TextWrapWord
	e.Refresh()
	if e.scroller.Content == nil {
		t.Fatal("expected scroller to reattach")
	}
	if e.scroller.Offset != fyne.NewPos(0, 0) {
		t.Errorf("expected zero scroll offset after reattach, got %v", e.scroller.Offset)
	}
	e.sel.doubleTappedAtUnixMillis = 0
	e.tripleTapped = false
	e.TouchDown(touch)
	e.TouchUp(touch)
	if e.CursorRow != 0 {
		t.Errorf("Expected CursorRow 0 after wrap is restored, got %d", e.CursorRow)
	}
}

// The demo's Scroll mode turns wrapping off but keeps the scroller, so its offset stays.
func TestEntry_WrapOff_ScrollBothKeepsOffset(t *testing.T) {
	test.NewApp()
	e := NewMultiLineEntry()
	e.Wrapping = fyne.TextWrapWord
	e.SetText("Line 1" + strings.Repeat("\nLine", 19))
	w := test.NewTempWindow(t, e)
	w.Resize(fyne.NewSize(220, 160))

	// Park the cursor on a row that is already inside this offset, so refresh
	// keeps the scroll position instead of chasing the cursor back to the top.
	e.scroller.Offset.Y = 50
	e.CursorRow = 4
	e.Wrapping = fyne.TextWrapOff
	e.Scroll = fyne.ScrollBoth
	e.Refresh()

	if e.scroller.Content == nil {
		t.Fatal("expected scroller to stay attached")
	}
	if e.scroller.Offset.Y != 50 {
		t.Errorf("expected offset to stay 50, got %v", e.scroller.Offset)
	}
}
