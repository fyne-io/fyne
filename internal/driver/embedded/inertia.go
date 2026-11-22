package embedded

import (
	"math"
	"time"

	"fyne.io/fyne/v2"
)

// inertiaTarget is a draggable and a snapshot of its drag event.
type inertiaTarget struct {
	wid   fyne.Draggable
	event fyne.DragEvent
}

// inertiaState is the active flick and its window.
type inertiaState struct {
	window  *noosWindow
	targets []inertiaTarget
}

// startInertia stores the flick to decay each frame.
func (n *noosDriver) startInertia(w *noosWindow, targets []inertiaTarget) {
	if len(targets) == 0 {
		n.stopDrawTicker()
		n.inertia = nil
		return
	}
	n.startDrawTicker()
	n.inertia = &inertiaState{window: w, targets: targets}
}

// stopInertia ends a running flick and its drag.
func (n *noosDriver) stopInertia() {
	n.stopDrawTicker()

	state := n.inertia
	if state == nil {
		return
	}
	n.inertia = nil
	for _, target := range state.targets {
		target.wid.DragEnd()
	}
}

// startDrawTicker starts the frame ticker for momentum.
func (n *noosDriver) startDrawTicker() {
	if n.drawTicker == nil {
		n.drawTicker = time.NewTicker(time.Second / 60)
	}
}

func (n *noosDriver) stopDrawTicker() {
	if n.drawTicker != nil {
		n.drawTicker.Stop()
		n.drawTicker = nil
	}
}

// stepInertia applies one frame of momentum and returns the window to draw.
func (n *noosDriver) stepInertia() *noosWindow {
	state := n.inertia
	if state == nil {
		return nil
	}
	n.inertia = nil

	active := make([]inertiaTarget, 0, len(state.targets))
	for _, target := range state.targets {
		event := target.event
		event.Dragged.DX *= tapMoveDecay
		event.Dragged.DY *= tapMoveDecay

		if math.Abs(float64(event.Dragged.DX)) <= tapMoveEndThreshold &&
			math.Abs(float64(event.Dragged.DY)) <= tapMoveEndThreshold {
			target.wid.DragEnd()
			continue
		}

		target.event = event
		active = append(active, target)
		target.wid.Dragged(&event)
	}

	if len(active) > 0 {
		n.inertia = &inertiaState{window: state.window, targets: active}
	} else {
		n.stopDrawTicker()
	}
	return state.window
}
