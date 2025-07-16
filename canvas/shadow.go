package canvas

import (
	"image/color"

	"fyne.io/fyne/v2"
)

// ShadowType defines the type of shadow to render.
type ShadowType int

const (
	DropShadow ShadowType = iota
	BoxShadow
)

// Shadowable defines an interface for objects that can render a shadow.
// It provides methods to retrieve the shadow paddings.
//
// Since: 2.7
type Shadowable interface {
	// ShadowPaddings returns the paddings (left, top, right, bottom) of the shadow.
	ShadowPaddings() [4]float32
}

// Ensure shadow implements Shadowable.
var _ Shadowable = (*shadow)(nil)

// shadow provides base functionality for objects that can have a shadow.
// Intended to be embedded in other structs to add shadow support.
type shadow struct {
	ShadowColor    color.Color   // Color of the shadow.
	ShadowSoftness float32       // Softness (blur radius) of the shadow.
	ShadowOffset   fyne.Position // Offset of the shadow relative to the content.
	ShadowType     ShadowType    // Type of shadow (DropShadow or BoxShadow).
}

// ShadowPaddings calculates the shadow paddings (left, top, right, bottom) based on offset and softness.
func (r *shadow) ShadowPaddings() [4]float32 {
	offsetX := r.ShadowOffset.X
	offsetY := r.ShadowOffset.Y
	softness := r.ShadowSoftness

	rightReach := -offsetX + softness
	leftReach := offsetX + softness
	topReach := -offsetY + softness
	bottomReach := offsetY + softness

	var padLeft, padRight, padTop, padBottom float32

	if leftReach > 0 {
		padLeft = leftReach
	}
	if rightReach > 0 {
		padRight = rightReach
	}
	if topReach > 0 {
		padTop = topReach
	}
	if bottomReach > 0 {
		padBottom = bottomReach
	}

	// Returns paddings in order: left, top, right, bottom.
	return [4]float32{padLeft, padTop, padRight, padBottom}
}
