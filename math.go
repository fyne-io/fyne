package fyne

// Min returns the smaller of the passed values.
func Min(x, y float32) float32 {
	if x < y {
		return x
	}
	return y
}

// Max returns the largest of the passed values.
func Max(nums ...float32) float32 {
	if len(nums) == 0 {
		return 0
	}
	if len(nums) == 1 {
		return nums[0]
	}
	max := nums[0]
	for i := 1; i < len(nums); i++ {
		if nums[i] > max {
			max = nums[i]
		}
	}
	return max
}
