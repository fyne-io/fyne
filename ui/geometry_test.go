package ui

import "testing"

func TestSizeAdd(t *testing.T) {
	size1 := NewSize(10, 10)
	size2 := NewSize(25, 25)

	size3 := size1.Add(size2)

	if size3.Width != size1.Width+size2.Width || size3.Height != size1.Height+size2.Height {
		t.Fatal("Expected", NewPos(35, 35), "but got", size3)
	}
}

func TestSizeUnion(t *testing.T) {
	size1 := NewSize(10, 100)
	size2 := NewSize(100, 10)

	size3 := size1.Union(size2)

	maxW := Max(size1.Width, size2.Width)
	maxH := Max(size1.Height, size2.Height)

	if size3.Width != maxW || size3.Height != maxH {
		t.Fatal("Expected", NewSize(100, 100), "but got", size3)
	}
}

func TestPositionAdd(t *testing.T) {
	pos1 := NewPos(10, 10)
	pos2 := NewPos(25, 25)

	pos3 := pos1.Add(pos2)

	if pos3.X != pos1.X+pos2.X || pos3.Y != pos1.Y+pos2.Y {
		t.Fatal("Expected", NewPos(35, 35), "but got", pos3)
	}
}
