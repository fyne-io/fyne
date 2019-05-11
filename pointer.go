package fyne

// Pointer Events implementation based on W3C specs and various implementations
// of mouse, touch and tablet handling thoughout platforms and frameworks
// See W3C spec for more info https://www.w3.org/TR/pointerevents

import (
	"fmt"
	"strings"
	"sync"
)

type pointerActiveButtonState struct {
	// set to a pointable target if the pointer is in capturing mode
	capturingTarget Pointable
}

var (
	// mutex guarding these pointer-related variables
	pointerMutex sync.Mutex
	// assigned app-unique PointerIDs for original raw platform ID for each device type
	allocatedPointerIDs [numPointerDeviceTypes]map[uint64]PointerID
	// last allocated app-unique PointerID
	lastAllocatedPointerID PointerID

	// Mapping of PointerID of the pointer in the 'active buttons state' to an optinal
	// capturing Pointable. For mouse, active buttons stat is when the device has
	// at least one button depressed, for a pen it is when at least one button
	// is depressed while hovering.
	buttonActivePointers = make(map[PointerID]Pointable)
)

func init() {
	// prepare maps for every device type
	for n := 0; n < int(numPointerDeviceTypes); n++ {
		allocatedPointerIDs[n] = make(map[uint64]PointerID)
	}
}

// PointerIDForPlatformPointerID returns application-unique PointerID for provided device type and raw platform ID.
// Raw platform ID is expected to be unique among all active pointers for that device type.
// For example for a touch device, each active finger must have unique ID but that ID doesn't have to be unique among
// IDs currently used by a tablet device type.
func PointerIDForPlatformPointerID(pointerDevType PointerDeviceType, rawPlatformID uint64) PointerID {
	pointerMutex.Lock()
	defer pointerMutex.Unlock()

	if pointerID, allocated := allocatedPointerIDs[pointerDevType][rawPlatformID]; allocated {
		return pointerID
	}

	lastAllocatedPointerID++
	allocatedPointerIDs[pointerDevType][rawPlatformID] = lastAllocatedPointerID

	return lastAllocatedPointerID
}

// SetPointerCapture sets pointer capture for the pointer identified by the argument
// pointerID to the capturingTarget. For subsequent events of the pointer,
// the capturing target will substitute the normal hit testing result as if the pointer
// is always over the capturing target, and they MUST always be targeted at this
// element until capture is released. The pointer MUST be in its active buttons state
// for this method to be effective, otherwise it fails and returns false.
func SetPointerCapture(pointerID PointerID, capturingTarget Pointable) bool {
	pointerMutex.Lock()
	defer pointerMutex.Unlock()

	if existingCapturingTarget, exists := buttonActivePointers[pointerID]; exists {
		ev := PointerEvent{ID: pointerID}

		if existingCapturingTarget != nil && existingCapturingTarget != capturingTarget {
			SendPointerEvent(existingCapturingTarget, PointerLostCapture, &ev)
		}

		buttonActivePointers[pointerID] = capturingTarget
		SendPointerEvent(capturingTarget, PointerGotCapture, &ev)
		return true
	}
	return false
}

// GetPointerCapture returns the element requesting pointer capture for the pointer identified
// by the argument pointerID or nil if the pointer is not in the capturing mode.
func GetPointerCapture(pointerID PointerID) Pointable {
	pointerMutex.Lock()
	defer pointerMutex.Unlock()

	return buttonActivePointers[pointerID]
}

// ReleasePointerCapture releases pointer capture for the pointer identified
// by the argument pointerID. Subsequent events for the pointer follow normal
// hit testing mechanisms.
func ReleasePointerCapture(pointerID PointerID) {
	pointerMutex.Lock()
	defer pointerMutex.Unlock()

	if tgt, exists := buttonActivePointers[pointerID]; exists {
		if tgt != nil {
			ev := PointerEvent{ID: pointerID}
			SendPointerEvent(tgt, PointerLostCapture, &ev)
		}
	}
}

func setPointerActiveButtonsState(pointerID PointerID, active bool) {
	pointerMutex.Lock()
	defer pointerMutex.Unlock()

	if active {
		buttonActivePointers[pointerID] = nil
	} else {
		delete(buttonActivePointers, pointerID)
	}
}

// SendPointerEvent sends event to a pointable object.
// It also automatically handles the pointer capture and release.
// This is only called by drivers to send input events to the Fyne.
// Provided ev pointer must not be altered until return from this function.
func SendPointerEvent(tgt Pointable, phase PointerEventPhase, ev *PointerEvent) bool {
	// Event is called synchronously on the target.
	//
	// Originally I used sync.Pool, made a copy of ev and called the target
	// handler in a gorountine with the copy. That didn't work well because
	// pointer events has to be sent serialized, in the original order.
	//
	// The receiver is expected not to block the execution. It is their
	// responsibility to use go-rountines for long-running operations!

	if phase == PointerDown {
		// now in buttons active sate, add to the list
		setPointerActiveButtonsState(ev.ID, true)
	} else if phase == PointerUp {
		// release the capture if set
		ReleasePointerCapture(ev.ID)
		// remove from the list of buttons active pointers
		setPointerActiveButtonsState(ev.ID, false)
	}

	return tgt.PointerEvent(phase, ev)
}

// PointerDeviceType is a type of the pointer device
type PointerDeviceType uint8

func (pdt PointerDeviceType) String() string {
	switch pdt {
	case MousePointerDevice:
		return "MousePointerDevice"
	case TouchPointerDevice:
		return "TouchPointerDevice"
	case TouchpadPointerDevice:
		return "TouchpadPointerDevice"
	case StylusPointerDevice:
		return "StylusPointerDevice"
	default:
		return "(unknown)"
	}
}

// Pointer device types
const (
	// Common mouse pointer
	MousePointerDevice PointerDeviceType = iota
	// Finger on a touch screen
	TouchPointerDevice
	// Finger on a touchpad device (Windows 8.1+ and Android are able to differentiate)
	TouchpadPointerDevice
	// A stylus of a tablet device
	StylusPointerDevice
	// A special stylus that knows about rotation
	//RotationStylusPointerDevice
	// A mouse with Z axis
	//FourDMouse

	// Number of defined pointer device types
	numPointerDeviceTypes
)

// PointerEventPhase describes a phase of pointing device op (up, down, stationary, ...)
type PointerEventPhase uint8

func (pep PointerEventPhase) String() string {
	switch pep {
	case PointerDown:
		return "PointerDown"
	case PointerUpdated:
		return "PointerUpdated"
	case PointerUp:
		return "PointerUp"
	case PointerCancelled:
		return "PointerCancelled"
	case PointerEntered:
		return "PointerEntered"
	case PointerLeft:
		return "PointerLeft"
	case PointerGotCapture:
		return "PointerGotCapture"
	case PointerLostCapture:
		return "PointerLostCapture"
	default:
		return "(unknown)"
	}
}

// Pointing device phases
const (
	// Finger or pen made a physical contact or a pointer entered the active buttons state.
	// For mouse, this is when the device has at least one button depressed, for pen it is
	// when at least one button is depressed while hovering.
	PointerDown PointerEventPhase = iota
	// Finger, pen or mouse pointer changed coordinates, button states, pressure or other
	// properties defining the pointer state.
	PointerUpdated
	// Point tracking ended (finger or pen lifted)
	PointerUp
	// Point tracking unexpectedly ended by an interruption.
	// (a popup displayed, device disconnected or other platform and device-specific reason)
	PointerCancelled
	// Pointer entered the CanvasObject area
	PointerEntered
	// Pointer left the CanvasObject area
	PointerLeft
	// Pointer gained capturing mode for the target and subsequent pointer events will
	// be sent to that target, until the mode is lost.
	PointerGotCapture
	// Pointer lost capturing mode for the target and subsequent pointer events will
	// honor normal hit testing.
	PointerLostCapture
)

// PointerButton holds the state of a button or a set of buttons as bits.
type PointerButton uint32

// String returns string representation of buttons in the set
func (pb PointerButton) String() string {
	if pb == NoButton {
		return "NoButton"
	} else if pb == AnyButton {
		return "AnyButton"
	}

	var btns []string

	if pb&LeftButton != 0 {
		btns = append(btns, "LeftButton")
	}
	if pb&RightButton != 0 {
		btns = append(btns, "RightButton")
	}
	if pb&MiddleButton != 0 {
		btns = append(btns, "MiddleButton")
	}
	if pb&BackButton != 0 {
		btns = append(btns, "BackButton")
	}
	if pb&ForwardButton != 0 {
		btns = append(btns, "ForwardButton")
	}
	if pb&TaskButton != 0 {
		btns = append(btns, "TaskButton")
	}

	bit := pb >> 7
	curBtn := 4

	for bit > 0 {
		if bit&1 != 0 {
			btns = append(btns, fmt.Sprintf("ExtraButton%d", curBtn))
			bit >>= 1
			curBtn++
		}
	}

	return strings.Join(btns, "|")
}

// Pressed returns thrue if the button btn is present in the set
// That indicates the button is pressed or has been changed (depending on the context)
func (pb PointerButton) Pressed(btn PointerButton) bool {
	return pb&btn != 0
}

// Count returns the number of buttons active in the set
func (pb PointerButton) Count() int {
	num := 0
	for i := uint(0); i < 32; i++ {
		if pb&PointerButton(1<<i) != 0 {
			num++
		}
	}
	return num
}

// Pointer device buttons
const (
	NoButton        PointerButton = 0
	AnyButton       PointerButton = 0x0fffffff
	LeftButton      PointerButton = 0x1
	RightButton     PointerButton = 0x2
	PenBarrelButton PointerButton = RightButton
	MiddleButton    PointerButton = 0x4
	// Back button on the side of a mouse
	BackButton    PointerButton = 0x8
	ExtraButton1                = BackButton
	ForwardButton PointerButton = 0x10
	// Forward button on the side of a mouse
	ExtraButton2   PointerButton = ForwardButton
	TaskButton     PointerButton = 0x20
	PenEaserButton PointerButton = TaskButton
	ExtraButton3   PointerButton = TaskButton
	ExtraButton4   PointerButton = 0x40
	ExtraButton5   PointerButton = 0x80
	ExtraButton6   PointerButton = 0x100
	ExtraButton7   PointerButton = 0x200
	ExtraButton8   PointerButton = 0x400
	ExtraButton9   PointerButton = 0x800
	ExtraButton10  PointerButton = 0x1000
	ExtraButton11  PointerButton = 0x2000
	ExtraButton12  PointerButton = 0x4000
	ExtraButton13  PointerButton = 0x8000
	ExtraButton14  PointerButton = 0x10000
	ExtraButton15  PointerButton = 0x20000
	ExtraButton16  PointerButton = 0x40000
	ExtraButton17  PointerButton = 0x80000
	ExtraButton18  PointerButton = 0x100000
	ExtraButton19  PointerButton = 0x200000
	ExtraButton20  PointerButton = 0x400000
	ExtraButton21  PointerButton = 0x800000
	ExtraButton22  PointerButton = 0x1000000
	ExtraButton23  PointerButton = 0x2000000
	ExtraButton24  PointerButton = 0x4000000
)
