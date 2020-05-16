package layout

import (
	"fyne.io/fyne"
	"fyne.io/fyne/internal/layout"
)

// NewHBoxLayout returns a horizontal box layout for stacking a number of child
// canvas objects or widgets left to right.
func NewHBoxLayout() fyne.Layout {
	return &layout.Box{Horizontal: true}
}

// NewVBoxLayout returns a vertical box layout for stacking a number of child
// canvas objects or widgets top to bottom.
func NewVBoxLayout() fyne.Layout {
	return &layout.Box{Horizontal: false}
}
