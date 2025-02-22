package common

var framecounter uint64 = 1

func IncrementFrameCounter() {
	framecounter += 1
}

func CurrentFrameCounter() uint64 {
	return framecounter
}
