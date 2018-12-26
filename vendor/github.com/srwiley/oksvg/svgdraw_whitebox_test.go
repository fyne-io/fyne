// Copyright 2018 The oksvg Authors. All rights reserved.
// created: 2018 by S.R.Wiley
package oksvg

import "testing"

func TestReadFloat(t *testing.T) {
	c := new(PathCursor)

	fStr := "23.4.56"
	c.ReadFloat(fStr)
	if len(c.points) != 2 || c.points[0] != 23.4 || c.points[1] != 0.56 {
		t.Error("read float failed", fStr, c.points)
	}

	c.points = nil

	fStr = "23.4"
	c.ReadFloat(fStr)
	if len(c.points) != 1 || c.points[0] != 23.4 {
		t.Error("read float failed", fStr)
	}
	c.points = nil

	fStr = "23"
	c.ReadFloat(fStr)
	if len(c.points) != 1 || c.points[0] != 23.0 {
		t.Error("read float failed", fStr, c.points)
	}

	c.points = nil
	fStr = ".4"
	c.ReadFloat(fStr)
	if len(c.points) != 1 || c.points[0] != 0.4 {
		t.Error("read float failed", fStr, c.points)
	}

	c.points = nil
	fStr = "23.4.56.67.32"
	c.ReadFloat(fStr)
	if len(c.points) != 4 || c.points[0] != 23.4 || c.points[1] != 0.56 ||
		c.points[2] != 0.67 || c.points[3] != 0.32 {
		t.Error("read float failed", fStr, c.points)
	}

}
