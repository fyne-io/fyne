package ui

import "testing"

func TestGeometry(t *testing.T) {
	pos1 := NewPos(10, 10)
	pos2 := NewPos(25, 25)

	pos3 := pos1.Add(pos2)

	if pos3.X != pos1.X+pos2.X || pos3.Y != pos1.Y+pos2.Y {
		t.Fatal("Expected", NewPos(25, 25), "but got", pos3)
	}
}
