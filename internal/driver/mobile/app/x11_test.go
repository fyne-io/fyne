//go:build (linux && !android) || freebsd || openbsd

package app

import (
	"testing"
	"time"

	"fyne.io/fyne/v2/internal/driver/mobile/event/key"
	"fyne.io/fyne/v2/internal/driver/mobile/event/mouse"
	"fyne.io/fyne/v2/internal/driver/mobile/event/touch"
)

func TestX11KeySymToFyneKeyCode(t *testing.T) {
	tests := map[int]key.Code{
		// letters and digits
		0x0061: key.CodeA, // XK_a
		0x007a: key.CodeZ, // XK_z
		0x0041: key.CodeA, // XK_A
		0x005a: key.CodeZ, // XK_Z
		0x0030: key.Code0, // XK_0
		0x0039: key.Code9, // XK_9
		// punctuation
		0x0020: key.CodeSpacebar,    // XK_space
		0x0027: key.CodeApostrophe,  // XK_apostrophe
		0x002c: key.CodeComma,       // XK_comma
		0x002d: key.CodeHyphenMinus, // XK_minus
		0x002e: key.CodeFullStop,    // XK_period
		0x002f: key.CodeSlash,       // XK_slash
		0x003b: key.CodeSemicolon,   // XK_semicolon
		0x003d: key.CodeEqualSign,   // XK_equal
		0x005b: key.CodeLeftSquareBracket,
		0x005c: key.CodeBackslash,
		0x005d: key.CodeRightSquareBracket,
		0x0060: key.CodeGraveAccent, // XK_grave
		// control and navigation
		0xff08: key.CodeDeleteBackspace, // XK_BackSpace
		0xff09: key.CodeTab,             // XK_Tab
		0xff0d: key.CodeReturnEnter,     // XK_Return
		0xff1b: key.CodeEscape,          // XK_Escape
		0xff50: key.CodeHome,            // XK_Home
		0xff51: key.CodeLeftArrow,       // XK_Left
		0xff52: key.CodeUpArrow,         // XK_Up
		0xff53: key.CodeRightArrow,      // XK_Right
		0xff54: key.CodeDownArrow,       // XK_Down
		0xff55: key.CodePageUp,          // XK_Page_Up
		0xff56: key.CodePageDown,        // XK_Page_Down
		0xff57: key.CodeEnd,             // XK_End
		0xff63: key.CodeInsert,          // XK_Insert
		0xffff: key.CodeDeleteForward,   // XK_Delete
		// function keys
		0xffbe: key.CodeF1,
		0xffc9: key.CodeF12,
		0xffd5: key.CodeF24,
		// keypad (number keys with num lock on)
		0xff7f: key.CodeKeypadNumLock,     // XK_Num_Lock
		0xff8d: key.CodeKeypadEnter,       // XK_KP_Enter
		0xffaa: key.CodeKeypadAsterisk,    // XK_KP_Multiply
		0xffab: key.CodeKeypadPlusSign,    // XK_KP_Add
		0xffad: key.CodeKeypadHyphenMinus, // XK_KP_Subtract
		0xffae: key.CodeKeypadFullStop,    // XK_KP_Decimal
		0xffaf: key.CodeKeypadSlash,       // XK_KP_Divide
		0xffb0: key.CodeKeypad0,           // XK_KP_0
		0xffb1: key.CodeKeypad1,           // XK_KP_1
		0xffb2: key.CodeKeypad2,           // XK_KP_2
		0xffb3: key.CodeKeypad3,           // XK_KP_3
		0xffb4: key.CodeKeypad4,           // XK_KP_4
		0xffb5: key.CodeKeypad5,           // XK_KP_5
		0xffb6: key.CodeKeypad6,           // XK_KP_6
		0xffb7: key.CodeKeypad7,           // XK_KP_7
		0xffb8: key.CodeKeypad8,           // XK_KP_8
		0xffb9: key.CodeKeypad9,           // XK_KP_9
		0xffbd: key.CodeKeypadEqualSign,   // XK_KP_Equal
		// keypad (navigation with num lock off)
		0xff95: key.CodeHome,          // XK_KP_Home
		0xff96: key.CodeLeftArrow,     // XK_KP_Left
		0xff97: key.CodeUpArrow,       // XK_KP_Up
		0xff98: key.CodeRightArrow,    // XK_KP_Right
		0xff99: key.CodeDownArrow,     // XK_KP_Down
		0xff9a: key.CodePageUp,        // XK_KP_Page_Up
		0xff9b: key.CodePageDown,      // XK_KP_Page_Down
		0xff9c: key.CodeEnd,           // XK_KP_End
		0xff9d: key.CodeHome,          // XK_KP_Begin
		0xff9e: key.CodeInsert,        // XK_KP_Insert
		0xff9f: key.CodeDeleteForward, // XK_KP_Delete
		// modifiers
		0xffe1: key.CodeLeftShift,
		0xffe2: key.CodeRightShift,
		0xffe3: key.CodeLeftControl,
		0xffe4: key.CodeRightControl,
		0xffe5: key.CodeCapsLock,
		0xffe9: key.CodeLeftAlt,
		0xffea: key.CodeRightAlt,
		0xffeb: key.CodeLeftGUI,
		0xffec: key.CodeRightGUI,
		// media keys
		0x1008ff11: key.CodeVolumeDown,
		0x1008ff12: key.CodeMute,
		0x1008ff13: key.CodeVolumeUp,
	}

	for keysym, want := range tests {
		if got := x11KeySymToFyneKeyCode(keysym); got != want {
			t.Errorf("x11KeySymToFyneKeyCode(0x%x) = %v, want %v", keysym, got, want)
		}
	}

	if got := x11KeySymToFyneKeyCode(0x12345); got != key.CodeUnknown {
		t.Errorf("unmapped keysym should be CodeUnknown, got %v", got)
	}
}

// markerEvent is sent onto the event queue to mark the end of the events that
// a mouse button has produced, the queue keeps its order.
type markerEvent struct{}

func TestX11MouseButtons(t *testing.T) {
	tests := []struct {
		name   string
		button int
		want   []touch.Type
	}{
		{name: "left", button: x11ButtonLeft, want: []touch.Type{touch.TypeBegin, touch.TypeEnd}},
		{name: "right", button: x11ButtonRight, want: []touch.Type{touch.TypeBegin, touch.TypeEnd}},
		// the middle button does not produce any event
		{name: "middle", button: x11ButtonMiddle},
		// the wheel tilt buttons 6 and 7 are not supported, no event
		{name: "wheel tilt left", button: 6},
		{name: "wheel tilt right", button: 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			onTouchBegin(10, 20, tt.button)
			onTouchEnd(10, 20, tt.button)
			theApp.Send(markerEvent{})

			var got []touch.Event
			for _, event := range collectEvents(t) {
				touchEvent, ok := event.(touch.Event)
				if !ok {
					t.Fatalf("expected a touch event for button %d, got %#v", tt.button, event)
				}
				got = append(got, touchEvent)
			}

			if len(got) != len(tt.want) {
				t.Fatalf("expected %d touch events for button %d, got %d: %v",
					len(tt.want), tt.button, len(got), got)
			}

			for i, want := range tt.want {
				if got[i].Type != want {
					t.Errorf("expected touch event %d for button %d to be %s, was %s",
						i, tt.button, want, got[i].Type)
				}
				if got[i].X != 10 || got[i].Y != 20 {
					t.Errorf("expected touch event %d for button %d at (10, 20), was (%v, %v)",
						i, tt.button, got[i].X, got[i].Y)
				}
			}
		})
	}
}

func TestX11ScrollWheel(t *testing.T) {
	tests := []struct {
		name   string
		button int
		want   mouse.Event
	}{
		{name: "up", button: x11ButtonWheelUp, want: mouse.Event{X: 10, Y: 20, ScrollY: 1}},
		{name: "down", button: x11ButtonWheelDown, want: mouse.Event{X: 10, Y: 20, ScrollY: -1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// a scroll wheel reports a press and a release, only the press scrolls
			onTouchBegin(10, 20, tt.button)
			onTouchEnd(10, 20, tt.button)
			theApp.Send(markerEvent{})

			events := collectEvents(t)
			if len(events) != 1 {
				t.Fatalf("expected 1 event for button %d, got %d: %#v", tt.button, len(events), events)
			}

			got, ok := events[0].(mouse.Event)
			if !ok {
				t.Fatalf("expected a scroll event for button %d, got %#v", tt.button, events[0])
			}
			if got != tt.want {
				t.Errorf("expected scroll event %#v for button %d, got %#v", tt.want, tt.button, got)
			}
		})
	}
}

func TestX11ButtonIsTouch(t *testing.T) {
	for button := 1; button <= 10; button++ {
		want := button == x11ButtonLeft || button == x11ButtonRight
		if got := x11ButtonIsTouch(button); got != want {
			t.Errorf("expected button %d to be a touch: %t, was %t", button, want, got)
		}
	}
}

// collectEvents gathers the events a mouse button produced, up to the marker
// event - the queue keeps their order.
func collectEvents(t *testing.T) []any {
	t.Helper()

	var events []any
	for {
		event := nextEvent(t)
		if _, ok := event.(markerEvent); ok {
			return events
		}

		events = append(events, event)
	}
}

func nextEvent(t *testing.T) any {
	t.Helper()

	select {
	case event := <-theApp.Events():
		return event
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for an event")

		return nil
	}
}
