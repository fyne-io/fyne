package ui

import "testing"

func TestMin(t *testing.T) {
	if Min(1, 3) != 1 {
		t.Fatal("Expected 1 but got", Min(1, 3))
	}
	if Min(1, -3) != -3 {
		t.Fatal("Expected -3 but got", Min(1, -3))
	}
}

func TestMax(t *testing.T) {
	if Max(1, 3) != 3 {
		t.Fatal("Expected 3 but got", Max(1, 3))
	}
	if Max(1, -3) != 1 {
		t.Fatal("Expected 1 but got", Max(1, -3))
	}
}
