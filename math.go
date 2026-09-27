package fyne

// Min returns the smaller of the passed values.
//
// Deprecated: use the Go builtin min(x, y) instead
//
//go:fix inline
func Min(x, y float32) float32 {
	return min(x, y)
}

// Max returns the larger of the passed values.
//
// Deprecated: use the Go builtin max(x, y) instead
//
//go:fix inline
func Max(x, y float32) float32 {
	return max(x, y)
}
