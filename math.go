package fyne

// Min returns the smaller of the passed values.
//
// Deprecated: Use the builtin [min] instead.
func Min(x, y float32) float32 {
	return min(x, y)
}

// Max returns the larger of the passed values.
//
// Deprecated: Use the builtin [max] instead.
func Max(x, y float32) float32 {
	return max(x, y)
}
