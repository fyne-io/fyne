package widget

import (
	"testing"

	"fyne.io/fyne/v2"
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
