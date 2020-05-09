// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package lifecycle defines an event for an app's lifecycle.
//
// The app lifecycle consists of moving back and forth between an ordered
// sequence of stages. For example, being at a stage greater than or equal to
// StageVisible means that the app is visible on the screen.
//
// A lifecycle event is a change from one stage to another, which crosses every
// intermediate stage. For example, changing from StageAlive to StageFocused
// implicitly crosses StageVisible.
//
// Crosses can be in a positive or negative direction. A positive crossing of
// StageFocused means that the app has gained the focus. A negative crossing
// means it has lost the focus.
//
// See the golang.org/x/mobile/app package for details on the event model.
package lifecycle // import "github.com/fyne-io/mobile/event/lifecycle"

import (
	"fmt"
)

// Cross is whether a lifecycle stage was crossed.
type Cross uint32

func (c Cross) String() string {
	switch c {
	case CrossOn:
		return "on"
	case CrossOff:
		return "off"
	}
	return "none"
}

const (
	CrossNone Cross = 0
	CrossOn   Cross = 1
	CrossOff  Cross = 2
)

// Event is a lifecycle change from an old stage to a new stage.
type Event struct {
	From, To Stage

	// DrawContext is the state used for painting, if any is valid.
	//
	// For OpenGL apps, a non-nil DrawContext is a gl.Context.
	//
	// TODO: make this an App method if we move away from an event channel?
	DrawContext interface{}
}

func (e Event) String() string {
	return fmt.Sprintf("lifecycle.Event{From:%v, To:%v, DrawContext:%v}", e.From, e.To, e.DrawContext)
}

// Crosses reports whether the transition from From to To crosses the stage s:
// 	- It returns CrossOn if it does, and the lifecycle change is positive.
// 	- It returns CrossOff if it does, and the lifecycle change is negative.
//	- Otherwise, it returns CrossNone.
// See the documentation for Stage for more discussion of positive and negative
// crosses.
func (e Event) Crosses(s Stage) Cross {
	switch {
	case e.From < s && e.To >= s:
		return CrossOn
	case e.From >= s && e.To < s:
		return CrossOff
	}
	return CrossNone
}

// Stage is a stage in the app's lifecycle. The values are ordered, so that a
// lifecycle change from stage From to stage To implicitly crosses every stage
// in the range (min, max], exclusive on the low end and inclusive on the high
// end, where min is the minimum of From and To, and max is the maximum.
//
// The documentation for individual stages talk about positive and negative
// crosses. A positive lifecycle change is one where its From stage is less
// than its To stage. Similarly, a negative lifecycle change is one where From
// is greater than To. Thus, a positive lifecycle change crosses every stage in
// the range (From, To] in increasing order, and a negative lifecycle change
// crosses every stage in the range (To, From] in decreasing order.
type Stage uint32

// TODO: how does iOS map to these stages? What do cross-platform mobile
// abstractions do?

const (
	// StageDead is the zero stage. No lifecycle change crosses this stage,
	// but:
	//	- A positive change from this stage is the very first lifecycle change.
	//	- A negative change to this stage is the very last lifecycle change.
	StageDead Stage = iota

	// StageAlive means that the app is alive.
	//	- A positive cross means that the app has been created.
	//	- A negative cross means that the app is being destroyed.
	// Each cross, either from or to StageDead, will occur only once.
	// On Android, these correspond to onCreate and onDestroy.
	StageAlive

	// StageVisible means that the app window is visible.
	//	- A positive cross means that the app window has become visible.
	//	- A negative cross means that the app window has become invisible.
	// On Android, these correspond to onStart and onStop.
	// On Desktop, an app window can become invisible if e.g. it is minimized,
	// unmapped, or not on a visible workspace.
	StageVisible

	// StageFocused means that the app window has the focus.
	//	- A positive cross means that the app window has gained the focus.
	//	- A negative cross means that the app window has lost the focus.
	// On Android, these correspond to onResume and onFreeze.
	StageFocused
)

func (s Stage) String() string {
	switch s {
	case StageDead:
		return "StageDead"
	case StageAlive:
		return "StageAlive"
	case StageVisible:
		return "StageVisible"
	case StageFocused:
		return "StageFocused"
	default:
		return fmt.Sprintf("lifecycle.Stage(%d)", s)
	}
}
