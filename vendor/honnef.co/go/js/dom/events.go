package dom

import (
	"time"

	"github.com/gopherjs/gopherjs/js"
)

func WrapEvent(o *js.Object) Event {
	return wrapEvent(o)
}

func wrapEvent(o *js.Object) Event {
	if o == nil || o == js.Undefined {
		return nil
	}
	ev := &BasicEvent{o}
	c := o.Get("constructor")
	switch c {
	case js.Global.Get("AnimationEvent"):
		return &AnimationEvent{ev}
	case js.Global.Get("AudioProcessingEvent"):
		return &AudioProcessingEvent{ev}
	case js.Global.Get("BeforeInputEvent"):
		return &BeforeInputEvent{ev}
	case js.Global.Get("BeforeUnloadEvent"):
		return &BeforeUnloadEvent{ev}
	case js.Global.Get("BlobEvent"):
		return &BlobEvent{ev}
	case js.Global.Get("ClipboardEvent"):
		return &ClipboardEvent{ev}
	case js.Global.Get("CloseEvent"):
		return &CloseEvent{BasicEvent: ev}
	case js.Global.Get("CompositionEvent"):
		return &CompositionEvent{ev}
	case js.Global.Get("CSSFontFaceLoadEvent"):
		return &CSSFontFaceLoadEvent{ev}
	case js.Global.Get("CustomEvent"):
		return &CustomEvent{ev}
	case js.Global.Get("DeviceLightEvent"):
		return &DeviceLightEvent{ev}
	case js.Global.Get("DeviceMotionEvent"):
		return &DeviceMotionEvent{ev}
	case js.Global.Get("DeviceOrientationEvent"):
		return &DeviceOrientationEvent{ev}
	case js.Global.Get("DeviceProximityEvent"):
		return &DeviceProximityEvent{ev}
	case js.Global.Get("DOMTransactionEvent"):
		return &DOMTransactionEvent{ev}
	case js.Global.Get("DragEvent"):
		return &DragEvent{ev}
	case js.Global.Get("EditingBeforeInputEvent"):
		return &EditingBeforeInputEvent{ev}
	case js.Global.Get("ErrorEvent"):
		return &ErrorEvent{ev}
	case js.Global.Get("FocusEvent"):
		return &FocusEvent{ev}
	case js.Global.Get("GamepadEvent"):
		return &GamepadEvent{ev}
	case js.Global.Get("HashChangeEvent"):
		return &HashChangeEvent{ev}
	case js.Global.Get("IDBVersionChangeEvent"):
		return &IDBVersionChangeEvent{ev}
	case js.Global.Get("KeyboardEvent"):
		return &KeyboardEvent{BasicEvent: ev}
	case js.Global.Get("MediaStreamEvent"):
		return &MediaStreamEvent{ev}
	case js.Global.Get("MessageEvent"):
		return &MessageEvent{BasicEvent: ev}
	case js.Global.Get("MouseEvent"):
		return &MouseEvent{UIEvent: &UIEvent{ev}}
	case js.Global.Get("MutationEvent"):
		return &MutationEvent{ev}
	case js.Global.Get("OfflineAudioCompletionEvent"):
		return &OfflineAudioCompletionEvent{ev}
	case js.Global.Get("PageTransitionEvent"):
		return &PageTransitionEvent{ev}
	case js.Global.Get("PointerEvent"):
		return &PointerEvent{&MouseEvent{UIEvent: &UIEvent{ev}}}
	case js.Global.Get("PopStateEvent"):
		return &PopStateEvent{ev}
	case js.Global.Get("ProgressEvent"):
		return &ProgressEvent{ev}
	case js.Global.Get("RelatedEvent"):
		return &RelatedEvent{ev}
	case js.Global.Get("RTCPeerConnectionIceEvent"):
		return &RTCPeerConnectionIceEvent{ev}
	case js.Global.Get("SensorEvent"):
		return &SensorEvent{ev}
	case js.Global.Get("StorageEvent"):
		return &StorageEvent{ev}
	case js.Global.Get("SVGEvent"):
		return &SVGEvent{ev}
	case js.Global.Get("SVGZoomEvent"):
		return &SVGZoomEvent{ev}
	case js.Global.Get("TimeEvent"):
		return &TimeEvent{ev}
	case js.Global.Get("TouchEvent"):
		return &TouchEvent{BasicEvent: ev}
	case js.Global.Get("TrackEvent"):
		return &TrackEvent{ev}
	case js.Global.Get("TransitionEvent"):
		return &TransitionEvent{ev}
	case js.Global.Get("UIEvent"):
		return &UIEvent{ev}
	case js.Global.Get("UserProximityEvent"):
		return &UserProximityEvent{ev}
	case js.Global.Get("WheelEvent"):
		return &WheelEvent{MouseEvent: &MouseEvent{UIEvent: &UIEvent{ev}}}
	default:
		return ev
	}
}

const (
	EvPhaseNone      = 0
	EvPhaseCapturing = 1
	EvPhaseAtTarget  = 2
	EvPhaseBubbling  = 3
)

type Event interface {
	Bubbles() bool
	Cancelable() bool
	CurrentTarget() Element
	DefaultPrevented() bool
	EventPhase() int
	Target() Element
	Timestamp() time.Time
	Type() string
	PreventDefault()
	StopImmediatePropagation()
	StopPropagation()
	Underlying() *js.Object
}

// Type BasicEvent implements the Event interface and is embedded by
// concrete event types.
type BasicEvent struct{ *js.Object }

type EventOptions struct {
	Bubbles    bool
	Cancelable bool
}

func CreateEvent(typ string, opts EventOptions) *BasicEvent {
	var event = js.Global.Get("Event").New(typ, js.M{
		"bubbles":    opts.Bubbles,
		"cancelable": opts.Cancelable,
	})
	return &BasicEvent{event}
}

func (ev *BasicEvent) Bubbles() bool {
	return ev.Get("bubbles").Bool()
}

func (ev *BasicEvent) Cancelable() bool {
	return ev.Get("cancelable").Bool()
}

func (ev *BasicEvent) CurrentTarget() Element {
	return wrapElement(ev.Get("currentTarget"))
}

func (ev *BasicEvent) DefaultPrevented() bool {
	return ev.Get("defaultPrevented").Bool()
}

func (ev *BasicEvent) EventPhase() int {
	return ev.Get("eventPhase").Int()
}

func (ev *BasicEvent) Target() Element {
	return wrapElement(ev.Get("target"))
}

func (ev *BasicEvent) Timestamp() time.Time {
	ms := ev.Get("timeStamp").Int()
	s := ms / 1000
	ns := (ms % 1000 * 1e6)
	return time.Unix(int64(s), int64(ns))
}

func (ev *BasicEvent) Type() string {
	return ev.Get("type").String()
}

func (ev *BasicEvent) PreventDefault() {
	ev.Call("preventDefault")
}

func (ev *BasicEvent) StopImmediatePropagation() {
	ev.Call("stopImmediatePropagation")
}

func (ev *BasicEvent) StopPropagation() {
	ev.Call("stopPropagation")
}

func (ev *BasicEvent) Underlying() *js.Object {
	return ev.Object
}

type AnimationEvent struct{ *BasicEvent }
type AudioProcessingEvent struct{ *BasicEvent }
type BeforeInputEvent struct{ *BasicEvent }
type BeforeUnloadEvent struct{ *BasicEvent }
type BlobEvent struct{ *BasicEvent }
type ClipboardEvent struct{ *BasicEvent }

type CloseEvent struct {
	*BasicEvent
	Code     int    `js:"code"`
	Reason   string `js:"reason"`
	WasClean bool   `js:"wasClean"`
}

type CompositionEvent struct{ *BasicEvent }
type CSSFontFaceLoadEvent struct{ *BasicEvent }
type CustomEvent struct{ *BasicEvent }
type DeviceLightEvent struct{ *BasicEvent }
type DeviceMotionEvent struct{ *BasicEvent }
type DeviceOrientationEvent struct{ *BasicEvent }
type DeviceProximityEvent struct{ *BasicEvent }
type DOMTransactionEvent struct{ *BasicEvent }
type DragEvent struct{ *BasicEvent }
type EditingBeforeInputEvent struct{ *BasicEvent }
type ErrorEvent struct{ *BasicEvent }

type FocusEvent struct{ *BasicEvent }

func (ev *FocusEvent) RelatedTarget() Element {
	return wrapElement(ev.Get("relatedTarget"))
}

type GamepadEvent struct{ *BasicEvent }
type HashChangeEvent struct{ *BasicEvent }
type IDBVersionChangeEvent struct{ *BasicEvent }

const (
	KeyLocationStandard = 0
	KeyLocationLeft     = 1
	KeyLocationRight    = 2
	KeyLocationNumpad   = 3
)

type KeyboardEvent struct {
	*BasicEvent
	AltKey        bool   `js:"altKey"`
	CharCode      int    `js:"charCode"`
	CtrlKey       bool   `js:"ctrlKey"`
	Key           string `js:"key"`
	KeyIdentifier string `js:"keyIdentifier"`
	KeyCode       int    `js:"keyCode"`
	Locale        string `js:"locale"`
	Location      int    `js:"location"`
	KeyLocation   int    `js:"keyLocation"`
	MetaKey       bool   `js:"metaKey"`
	Repeat        bool   `js:"repeat"`
	ShiftKey      bool   `js:"shiftKey"`
}

func (ev *KeyboardEvent) ModifierState(mod string) bool {
	return ev.Call("getModifierState", mod).Bool()
}

type MediaStreamEvent struct{ *BasicEvent }

type MessageEvent struct {
	*BasicEvent
	Data *js.Object `js:"data"`
}

type MouseEvent struct {
	*UIEvent
	AltKey    bool `js:"altKey"`
	Button    int  `js:"button"`
	ClientX   int  `js:"clientX"`
	ClientY   int  `js:"clientY"`
	CtrlKey   bool `js:"ctrlKey"`
	MetaKey   bool `js:"metaKey"`
	MovementX int  `js:"movementX"`
	MovementY int  `js:"movementY"`
	ScreenX   int  `js:"screenX"`
	ScreenY   int  `js:"screenY"`
	ShiftKey  bool `js:"shiftKey"`
}

func (ev *MouseEvent) RelatedTarget() Element {
	return wrapElement(ev.Get("relatedTarget"))
}

func (ev *MouseEvent) ModifierState(mod string) bool {
	return ev.Call("getModifierState", mod).Bool()
}

type MutationEvent struct{ *BasicEvent }
type OfflineAudioCompletionEvent struct{ *BasicEvent }
type PageTransitionEvent struct{ *BasicEvent }
type PointerEvent struct{ *MouseEvent }
type PopStateEvent struct{ *BasicEvent }
type ProgressEvent struct{ *BasicEvent }
type RelatedEvent struct{ *BasicEvent }
type RTCPeerConnectionIceEvent struct{ *BasicEvent }
type SensorEvent struct{ *BasicEvent }
type StorageEvent struct{ *BasicEvent }
type SVGEvent struct{ *BasicEvent }
type SVGZoomEvent struct{ *BasicEvent }
type TimeEvent struct{ *BasicEvent }

// TouchEvent represents an event sent when the state of contacts with a touch-sensitive
// surface changes. This surface can be a touch screen or trackpad, for example. The event
// can describe one or more points of contact with the screen and includes support for
// detecting movement, addition and removal of contact points, and so forth.
//
// Reference: https://developer.mozilla.org/en-US/docs/Web/API/TouchEvent.
type TouchEvent struct {
	*BasicEvent
	AltKey   bool `js:"altKey"`
	CtrlKey  bool `js:"ctrlKey"`
	MetaKey  bool `js:"metaKey"`
	ShiftKey bool `js:"shiftKey"`
}

// ChangedTouches lists all individual points of contact whose states changed between
// the previous touch event and this one.
//
// Reference: https://developer.mozilla.org/en-US/docs/Web/API/TouchEvent/changedTouches.
func (ev *TouchEvent) ChangedTouches() []*Touch {
	return touchListToTouches(ev.Get("changedTouches"))
}

// TargetTouches lists all points of contact that are both currently in contact with the
// touch surface and were also started on the same element that is the target of the event.
//
// Reference: https://developer.mozilla.org/en-US/docs/Web/API/TouchEvent/targetTouches.
func (ev *TouchEvent) TargetTouches() []*Touch {
	return touchListToTouches(ev.Get("targetTouches"))
}

// Touches lists all current points of contact with the surface, regardless of target
// or changed status.
//
// Reference: https://developer.mozilla.org/en-US/docs/Web/API/TouchEvent/touches.
func (ev *TouchEvent) Touches() []*Touch {
	return touchListToTouches(ev.Get("touches"))
}

func touchListToTouches(tl *js.Object) []*Touch {
	out := make([]*Touch, tl.Length())
	for i := range out {
		out[i] = &Touch{Object: tl.Index(i)}
	}
	return out
}

// Touch represents a single contact point on a touch-sensitive device. The contact point
// is commonly a finger or stylus and the device may be a touchscreen or trackpad.
//
// Reference: https://developer.mozilla.org/en-US/docs/Web/API/Touch.
type Touch struct {
	*js.Object
	Identifier    int     `js:"identifier"`
	ScreenX       float64 `js:"screenX"`
	ScreenY       float64 `js:"screenY"`
	ClientX       float64 `js:"clientX"`
	ClientY       float64 `js:"clientY"`
	PageX         float64 `js:"pageX"`
	PageY         float64 `js:"pageY"`
	RadiusX       float64 `js:"radiusX"`
	RadiusY       float64 `js:"radiusY"`
	RotationAngle float64 `js:"rotationAngle"`
	Force         float64 `js:"force"`
}

// Target returns the Element on which the touch point started when it was first placed
// on the surface, even if the touch point has since moved outside the interactive area
// of that element or even been removed from the document.
//
// Reference: https://developer.mozilla.org/en-US/docs/Web/API/Touch/target.
func (t *Touch) Target() Element {
	return wrapElement(t.Get("target"))
}

type TrackEvent struct{ *BasicEvent }
type TransitionEvent struct{ *BasicEvent }
type UIEvent struct{ *BasicEvent }
type UserProximityEvent struct{ *BasicEvent }

const (
	DeltaPixel = 0
	DeltaLine  = 1
	DeltaPage  = 2
)

type WheelEvent struct {
	*MouseEvent
	DeltaX    float64 `js:"deltaX"`
	DeltaY    float64 `js:"deltaY"`
	DeltaZ    float64 `js:"deltaZ"`
	DeltaMode int     `js:"deltaMode"`
}

type EventTarget interface {
	// AddEventListener adds a new event listener and returns the
	// wrapper function it generated. If using RemoveEventListener,
	// that wrapper has to be used.
	AddEventListener(typ string, useCapture bool, listener func(Event)) func(*js.Object)
	RemoveEventListener(typ string, useCapture bool, listener func(*js.Object))
	DispatchEvent(event Event) bool
}
