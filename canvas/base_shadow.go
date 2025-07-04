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
// It provides methods to calculate the scaled size and position including shadow paddings,
// retrieve the shadow paddings, and obtain the content size and position excluding the shadow.
// Content refers to the area without shadow, while the shadow is rendered outside this area.
// The canvas object size is reduced to include the shadow within the allocated frame.
//
// Since: 2.7
type Shadowable interface {
	// SizeAndPositionWithShadow returns the total size and position needed to render the shadow around the content
	// with the input size. Total size is larger than the input size.
	SizeAndPositionWithShadow(size fyne.Size) (fyne.Size, fyne.Position)
	// ShadowPaddings returns the paddings (left, top, right, bottom) for the shadow.
	ShadowPaddings() [4]float32
	// ContentSize returns the size of the content excluding shadow paddings.
	ContentSize() fyne.Size
	// ContentPos returns the position of the content excluding shadow paddings.
	ContentPos() fyne.Position
}

// Ensure baseShadow implements Shadowable.
var _ Shadowable = (*baseShadow)(nil)

// baseShadow provides base functionality for objects that can have a shadow.
// Intended to be embedded in other structs to add shadow support.
type baseShadow struct {
	baseObject

	ShadowColor    color.Color   // Color of the shadow.
	ShadowSoftness float32       // Softness (blur radius) of the shadow.
	ShadowOffset   fyne.Position // Offset of the shadow relative to the content.
	ShadowType     ShadowType    // Type of shadow (DropShadow or BoxShadow).
}

// SizeAndPositionWithShadow calculates the total size and adjusted position needed to accommodate the shadow around the content.
// The returned size includes the shadow paddings on all sides, while the position shifts the shadow to ensure content is correctly aligned within the shadow area.
// The input size parameter represents the original content size, excluding any shadow.
func (r *baseShadow) SizeAndPositionWithShadow(size fyne.Size) (fyne.Size, fyne.Position) {
	paddings := r.ShadowPaddings()
	return fyne.NewSize(size.Width+paddings[0]+paddings[2], size.Height+paddings[1]+paddings[3]), fyne.NewPos(-paddings[0], -paddings[1])
}

// ShadowPaddings calculates the shadow paddings (left, top, right, bottom) based on offset and softness.
func (r *baseShadow) ShadowPaddings() [4]float32 {
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

// ContentSize returns the size of the content area, excluding any space added for shadow paddings.
func (r *baseShadow) ContentSize() fyne.Size {
	paddings := r.ShadowPaddings()
	size := r.baseObject.Size()
	return fyne.NewSize(size.Width-paddings[0]-paddings[2], size.Height-paddings[1]-paddings[3])
}

// ContentPos returns the top-left position of the content area, adjusted to exclude the shadow paddings.
// This gives the position where the actual content (without shadow) is rendered within the shadowed object.
func (r *baseShadow) ContentPos() fyne.Position {
	paddings := r.ShadowPaddings()
	position := r.baseObject.Position()
	return fyne.NewPos(position.X+paddings[0], position.Y+paddings[1])
}