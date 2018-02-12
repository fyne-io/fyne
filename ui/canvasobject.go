package ui

import "image/color"

type CanvasObject interface {
	SetColor(color.RGBA)
	Canvas() Canvas
}
