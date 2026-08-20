package widget

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/test"
)

func TestEntry_TouchDown_AccountsForScroll(t *testing.T) {
	test.NewApp()
	e := NewMultiLineEntry()
	e.SetText("Line 1\nLine 2\nLine 3\nLine 4\nLine 5")

	e.CreateRenderer()
	e.scroll.Offset.Y = 50

	touch := &mobile.TouchEvent{
		PointEvent: fyne.PointEvent{
			Position: fyne.NewPos(10, 10),
		},
	}

	e.TouchDown(touch)

	if e.CursorRow == 0 {
		t.Errorf("Expected CursorRow > 0, got %d", e.CursorRow)
	}
}
