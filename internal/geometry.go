package internal

import "fyne.io/fyne/v2"

// MaxSizes returns the element-wise maximum of the two sizes. It is what
// fyne.Size.Max does, minus boxing the argument into the Vector2 parameter,
// which heap-allocates in per-frame layout code.
func MaxSizes(a, b fyne.Size) fyne.Size {
	return fyne.Size{Width: fyne.Max(a.Width, b.Width), Height: fyne.Max(a.Height, b.Height)}
}

// MinSizes returns the element-wise minimum of the two sizes, allocation-free
// like MaxSizes.
func MinSizes(a, b fyne.Size) fyne.Size {
	return fyne.Size{Width: fyne.Min(a.Width, b.Width), Height: fyne.Min(a.Height, b.Height)}
}
