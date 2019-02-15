package fyne

// Min returns the smaller of the passed values.
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// Max returns the larger of the passed values.
func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
