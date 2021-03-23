package fyne

// Min returns the smaller of the passed values.
func Min(x, y float32) float32 {
	if x < y {
		return x
	}
	return y
}

// Max returns the larger of the passed values.
func Max(x, y float32) float32 {
	if x > y {
		return x
	}
	return y
}
