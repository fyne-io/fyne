package copy

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type CopyCanvasObject interface {
	fyne.CanvasObject

	// Give a reference to the CanvasObject used as a source, should only be used for the pointer value, but not for any of its field
	Ref() fyne.CanvasObject
}

type CopyReference struct {
	ref fyne.CanvasObject
}

func (copy *CopyReference) Ref() fyne.CanvasObject {
	return copy.ref
}

var _ CopyCanvasObject = (*CopyCircle)(nil)
var _ CopyCanvasObject = (*CopyContainer)(nil)
var _ CopyCanvasObject = (*CopyImage)(nil)
var _ CopyCanvasObject = (*CopyLine)(nil)
var _ CopyCanvasObject = (*CopyLinearGradient)(nil)
var _ CopyCanvasObject = (*CopyRadialGradient)(nil)
var _ CopyCanvasObject = (*CopyRaster)(nil)
var _ CopyCanvasObject = (*CopyRectangle)(nil)
var _ CopyCanvasObject = (*CopyText)(nil)

type CopyCircle struct {
	canvas.Circle
	CopyReference
}

func NewCopyCircle(source *canvas.Circle) *CopyCircle {
	return &CopyCircle{Circle: *source, CopyReference: CopyReference{source}}
}

type CopyContainer struct {
	fyne.Container
	CopyReference
}

func NewCopyContainer(source *fyne.Container) *CopyContainer {
	return &CopyContainer{Container: *source, CopyReference: CopyReference{source}}
}

type CopyImage struct {
	canvas.Image
	CopyReference
}

func NewCopyImage(source *canvas.Image) *CopyImage {
	return &CopyImage{Image: *source, CopyReference: CopyReference{source}}
}

type CopyLine struct {
	canvas.Line
	CopyReference
}

func NewCopyLine(source *canvas.Line) *CopyLine {
	return &CopyLine{Line: *source, CopyReference: CopyReference{source}}
}

type CopyLinearGradient struct {
	canvas.LinearGradient
	CopyReference
}

func NewCopyLinearGradient(source *canvas.LinearGradient) *CopyLinearGradient {
	return &CopyLinearGradient{LinearGradient: *source, CopyReference: CopyReference{source}}
}

type CopyRadialGradient struct {
	canvas.RadialGradient
	CopyReference
}

func NewCopyRadialGradient(source *canvas.RadialGradient) *CopyRadialGradient {
	return &CopyRadialGradient{RadialGradient: *source, CopyReference: CopyReference{source}}
}

type CopyRaster struct {
	canvas.Raster
	CopyReference
}

func NewCopyRaster(source *canvas.Raster) *CopyRaster {
	return &CopyRaster{Raster: *source, CopyReference: CopyReference{source}}
}

type CopyRectangle struct {
	canvas.Rectangle
	CopyReference
}

func NewCopyRectangle(source *canvas.Rectangle) *CopyRectangle {
	return &CopyRectangle{Rectangle: *source, CopyReference: CopyReference{source}}
}

type CopyText struct {
	canvas.Text
	CopyReference
}

func NewCopyText(source *canvas.Text) *CopyText {
	return &CopyText{Text: *source, CopyReference: CopyReference{source}}
}
