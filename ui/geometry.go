package ui

type Size struct {
	Width  int
	Height int
}

func NewSize(w int, h int) Size {
	return Size{w, h}
}

type Position struct {
	X int
	Y int
}

func NewPos(x int, y int) Position {
	return Position{x, y}
}
