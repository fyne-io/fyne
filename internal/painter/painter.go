package painter

import (
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal"
)

// Painter defines the functionality of a Fyne renderer.
//
// It lives here rather than beside a particular backend so that a canvas can be
// wired to any renderer without dragging that renderer's platform dependencies
// (cgo OpenGL, Direct3D, ...) into everything that merely holds a painter.
type Painter interface {
	// Init tell a new painter to initialize, usually called after a context is available
	Init()
	// Capture requests that the specified canvas be drawn to an in-memory image
	Capture(fyne.Canvas) image.Image
	// Clear tells our painter to prepare a fresh paint
	Clear()
	// Free is used to indicate that a certain canvas object is no longer needed
	Free(fyne.CanvasObject)
	// Paint a single fyne.CanvasObject but not its children.
	Paint(fyne.CanvasObject, fyne.Position, fyne.Size, *internal.ClipItem)
	// SetFrameBufferScale tells us when we have more than 1 framebuffer pixel for each output pixel
	SetFrameBufferScale(float32)
	// SetOutputSize is used to change the resolution of our output viewport
	SetOutputSize(int, int)
	// StartClipping tells us that the following paint actions should be clipped to the specified area.
	StartClipping(fyne.Position, fyne.Size)
	// StopClipping stops clipping paint actions.
	StopClipping()
}
