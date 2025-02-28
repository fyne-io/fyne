package glfw

import (
	"fyne.io/fyne/v2/internal/driver/common"
)

type eventType int

const (
	emptyEvent eventType = iota
	moveEvent
	closeEvent
	mouseMoveEvent
	mouseClickEvent
	mouseScrollEvent
	keyPressEvent
	charInputEvent
	focusEvent
	dropEvent
)

// union of different GLFW events we receive
type event struct {
	eventType eventType

	window *window

	intArg1 int
	intArg2 int
	intArg3 int
	intArg4 int

	float64Arg1 float64
	float64Arg2 float64
}

var eventQueue = common.NewRingBuffer[event](1024)

func newCloseEvent(win *window) event {
	return event{eventType: closeEvent, window: win}
}

func newMoveEvent(win *window, x, y int) event {
	return event{
		eventType: moveEvent,
		window:    win,
		intArg1:   x,
		intArg2:   y,
	}
}

func newMouseMoveEvent(win *window, xpos, ypos float64) event {
	return event{
		eventType:   mouseMoveEvent,
		window:      win,
		float64Arg1: xpos,
		float64Arg2: ypos,
	}
}

func newMouseClickEvent(win *window, btn, action, mods int) event {
	return event{
		eventType: mouseClickEvent,
		window:    win,
		intArg1:   btn,
		intArg2:   action,
		intArg3:   mods,
	}
}

func newMouseScrollEvent(win *window, xoff, yoff float64) event {
	return event{
		eventType:   mouseScrollEvent,
		window:      win,
		float64Arg1: xoff,
		float64Arg2: yoff,
	}
}

func newKeyPressEvent(win *window, key, scancode, action, mods int) event {
	return event{
		eventType: keyPressEvent,
		window:    win,
		intArg1:   key,
		intArg2:   scancode,
		intArg3:   action,
		intArg4:   mods,
	}
}

func newCharInputEvent(win *window, char rune) event {
	return event{
		eventType: charInputEvent,
		window:    win,
		intArg1:   int(char),
	}
}

func newFocusEvent(win *window, focused bool) event {
	f := 0
	if focused {
		f = 1
	}
	return event{
		eventType: focusEvent,
		window:    win,
		intArg1:   f,
	}
}

func processEvent(e event) {
	switch e.eventType {
	case closeEvent:
		e.window.closed(e.window.viewport)
	case moveEvent:
		e.window.moved(e.window.viewport, e.intArg1, e.intArg2)
	case mouseMoveEvent:
		e.window.processMouseMoved(e.float64Arg1, e.float64Arg2)
	case mouseClickEvent:
		e.window.mouseClicked(e.window.viewport, e.intArg1, e.intArg2, e.intArg3)
	case mouseScrollEvent:
		e.window.mouseScrolled(e.window.viewport, e.float64Arg1, e.float64Arg2)
	case keyPressEvent:
		e.window.keyPressed(e.window.viewport, e.intArg1, e.intArg2, e.intArg3, e.intArg4)
	case charInputEvent:
		e.window.charInput(e.window.viewport, rune(e.intArg1))
	case focusEvent:
		e.window.focused(e.window.viewport, e.intArg1 > 0)
	}
}
