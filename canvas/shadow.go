package canvas

import (
	"image/color"

	"fyne.io/fyne/v2"
)

// ShadowVariant indicates a variation of a shadow.
//
// Since: 2.8
type ShadowVariant int

const (
	// DropShadow represents a shadow effect that is rendered exclusively outside the boundaries of the object,
	// following the object's shape and not appearing beneath its filled area.
	//
	// Since: 2.8
	DropShadow ShadowVariant = iota
	// BoxShadow represents a shadow effect that is rendered both behind and outside the object,
	// appearing as a blurred rectangle that extends beneath the object's filled area as well as beyond its edges.
	//
	// Since: 2.8
	BoxShadow
)

// Shadow provides base functionality for objects that can have a Shadow.
// Intended to be embedded in other structs to add Shadow support.
//
// Since: 2.8
type Shadow struct {
	Color      color.Color   // Color of the shadow.
	BlurRadius float32       // A value of 0 produces no blur, while larger values produce bigger and lighter shadow.
	Spread     float32       // Spread of the shadow (how far out to draw before fading - negative values make it smaller).
	Offset     fyne.Position // Offset of the shadow relative to the content. Positive values move the shadow to the right (x) and down (y) of the element.
	Variant    ShadowVariant // Variation of shadow (DropShadow or BoxShadow).
}
