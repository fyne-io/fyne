package canvas

import (
	"image/color"

	"fyne.io/fyne/v2"
)

// ShadowType defines the type of shadow to render.
//
// Since: 2.8
type ShadowType int

const (
	DropShadow ShadowType = iota
	BoxShadow
)

// Shadow provides base functionality for objects that can have a Shadow.
// Intended to be embedded in other structs to add Shadow support.
//
// Since: 2.8
type Shadow struct {
	ShadowColor    color.Color   // Color of the shadow.
	ShadowSoftness float32       // Softness (blur radius) of the shadow.
	ShadowOffset   fyne.Position // Offset of the shadow relative to the content.
	ShadowType     ShadowType    // Type of shadow (DropShadow or BoxShadow).
}

// NewShadow creates a new Shadow with the specified properties.
//
// Since: 2.8
func NewShadow(color color.Color, softness float32, offet fyne.Position, shadowType ShadowType) Shadow {
	return Shadow{
		ShadowColor:    color,
		ShadowSoftness: softness,
		ShadowOffset:   offet,
		ShadowType:     shadowType,
	}
}
