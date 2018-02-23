package ui

type Size struct {
	Width  int
	Height int
}

func (s1 Size) Add(s2 Size) Size {
	return Size{s1.Width + s2.Width, s1.Height + s2.Height}
}

func (s1 Size) Union(s2 Size) Size {
	maxW := Max(s1.Width, s2.Width)
	maxH := Max(s1.Height, s2.Height)

	return NewSize(maxW, maxH)
}

func NewSize(w int, h int) Size {
	return Size{w, h}
}

type Position struct {
	X int
	Y int
}

func (p1 Position) Add(p2 Position) Position {
	return Position{p1.X + p2.X, p1.Y + p2.Y}
}

func NewPos(x int, y int) Position {
	return Position{x, y}
}
