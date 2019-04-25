package theme

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"image/color"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func helperLoadBytes(t *testing.T, name string) []byte {
	path := filepath.Join("testdata", name) // relative path
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return bytes
}

func TestColorToHexString(t *testing.T) {
	white := color.RGBA{0xff, 0xff, 0xff, 0xff}
	result := ColorToHexString(white)
	expect := "#ffffff"
	assert.Equal(t, expect, result, "White hex value should be #ffffff.")
}

func TestSingleColorFill_SVGReplacePathsFill(t *testing.T) {
	svg := SVG{
		Paths: []*Path{
			{Fill: "#ffffff"},
		},
	}
	assert.Equal(t, "#ffffff", svg.Paths[0].Fill, "Initial path fill should be #ffffff")
	svg.replacePathsFill("#000000")
	assert.Equal(t, "#000000", svg.Paths[0].Fill, "Replaced path fill should be #000000")
}

func TestMultiColorFill_SVGReplacePathsFill(t *testing.T) {
	svg := SVG{
		Paths: []*Path{
			{Fill: "#ffffff"},
			{Fill: "#eeeeee"},
			{Fill: "none"},
		},
	}
	assert.Equal(t, "#ffffff", svg.Paths[0].Fill, "Initial path[0] fill should be #ffffff")
	assert.Equal(t, "#eeeeee", svg.Paths[1].Fill, "Initial path[1] fill should be #eeeeee")
	assert.Equal(t, "none", svg.Paths[2].Fill, "Initial path[2] fill should be none")
	svg.replacePathsFill("#000000")
	assert.Equal(t, "#000000", svg.Paths[0].Fill, "Replaced path fill should be #000000")
	assert.Equal(t, "#000000", svg.Paths[1].Fill, "Replaced path fill should be #000000")
	assert.Equal(t, "none", svg.Paths[2].Fill, "Replaced path fill should still be none")
}

func TestNoFill_SVGReplacePathsFill(t *testing.T) {
	svg := SVG{
		Paths: []*Path{
			{Fill: "none"},
		},
	}
	assert.Equal(t, "none", svg.Paths[0].Fill, "Initial path fill should be 'none'")
	svg.replacePathsFill("#000000")
	assert.Equal(t, "none", svg.Paths[0].Fill, "Replaced path fill should still be 'none'")
}

func TestSingleColorFill_SVGReplaceRectsFill(t *testing.T) {
	svg := SVG{
		Rects: []*Rect{
			{Fill: "#ffffff"},
		},
	}
	assert.Equal(t, "#ffffff", svg.Rects[0].Fill, "Initial rect fill should be #ffffff")
	svg.replaceRectsFill("#000000")
	assert.Equal(t, "#000000", svg.Rects[0].Fill, "Replaced rect fill should be #000000")
}

func TestMultiColorFill_SVGReplaceRectsFill(t *testing.T) {
	svg := SVG{
		Rects: []*Rect{
			{Fill: "#ffffff"},
			{Fill: "#eeeeee"},
			{Fill: "none"},
		},
	}
	assert.Equal(t, "#ffffff", svg.Rects[0].Fill, "Initial rect[0] fill should be #ffffff")
	assert.Equal(t, "#eeeeee", svg.Rects[1].Fill, "Initial rect[1] fill should be #eeeeee")
	assert.Equal(t, "none", svg.Rects[2].Fill, "Initial rect[2] fill should be none")
	svg.replaceRectsFill("#000000")
	assert.Equal(t, "#000000", svg.Rects[0].Fill, "Replaced rect[0] fill should be #000000")
	assert.Equal(t, "#000000", svg.Rects[1].Fill, "Replaced rect[1] fill should be #000000")
	assert.Equal(t, "none", svg.Rects[2].Fill, "Replaced rect[2] fill should still be none")
}

func TestNoFill_SVGReplaceRectsFill(t *testing.T) {
	svg := SVG{
		Rects: []*Rect{
			{Fill: "none"},
		},
	}
	assert.Equal(t, "none", svg.Rects[0].Fill, "Initial rect fill should be 'none'")
	svg.replaceRectsFill("#000000")
	assert.Equal(t, "none", svg.Rects[0].Fill, "Replaced rect fill should still be 'none'")
}

func TestSingleColorFill_SVGReplacePolygonsFill(t *testing.T) {
	svg := SVG{
		Polygons: []*Polygon{
			{Fill: "#ffffff"},
		},
	}
	assert.Equal(t, "#ffffff", svg.Polygons[0].Fill, "Initial polygon fill should be #ffffff")
	svg.replacePolygonsFill("#000000")
	assert.Equal(t, "#000000", svg.Polygons[0].Fill, "Replaced polygon fill should be #000000")
}

func TestMultiColorFill_SVGReplacePolygonsFill(t *testing.T) {
	svg := SVG{
		Polygons: []*Polygon{
			{Fill: "#ffffff"},
			{Fill: "#eeeeee"},
			{Fill: "none"},
		},
	}
	assert.Equal(t, "#ffffff", svg.Polygons[0].Fill, "Initial polygon[0] fill should be #ffffff")
	assert.Equal(t, "#eeeeee", svg.Polygons[1].Fill, "Initial polygon[1] fill should be #eeeeee")
	assert.Equal(t, "none", svg.Polygons[2].Fill, "Initial polygon[2] fill should still be none")
	svg.replacePolygonsFill("#000000")
	assert.Equal(t, "#000000", svg.Polygons[0].Fill, "Replaced polygon[0] fill should be #000000")
	assert.Equal(t, "#000000", svg.Polygons[1].Fill, "Replaced polygon[1] fill should be #000000")
	assert.Equal(t, "none", svg.Polygons[2].Fill, "Replaced polygon[2] fill should still be none")
}

func TestNoFill_SVGReplacePolygonsFill(t *testing.T) {
	svg := SVG{
		Polygons: []*Polygon{
			{Fill: "none"},
		},
	}
	assert.Equal(t, "none", svg.Polygons[0].Fill, "Initial polygon fill should be 'none'")
	svg.replacePolygonsFill("#000000")
	assert.Equal(t, "none", svg.Polygons[0].Fill, "Replaced polygon fill should still be 'none'")
}

func TestPathIcon_ReplaceFillColor(t *testing.T) {
	img := helperLoadBytes(t, "iconWithPaths.svg")
	rdr := bytes.NewReader(img)
	var svg SVG
	if err := svg.ReplaceFillColor(rdr, "#123456"); err != nil {
		t.Fatal("ReplaceFillColor threw an error:", err)
	}
	assert.Equal(t, "none", svg.Paths[0].Fill, "Fill color for path[0] should have been 'none'")
	assert.Equal(t, "#123456", svg.Paths[1].Fill, "Fill color for path[1] should have been #123456")
}

func TestRectPathIcon_ReplaceFillColor(t *testing.T) {
	img := helperLoadBytes(t, "iconWithRectsAndPaths.svg")
	rdr := bytes.NewReader(img)
	var svg SVG
	if err := svg.ReplaceFillColor(rdr, "#123456"); err != nil {
		t.Fatal("ReplaceFillColor threw an error:", err)
	}
	assert.Equal(t, "#123456", svg.Rects[0].Fill, "Fill color for rects[0] should have been #123456")
	assert.Equal(t, "#123456", svg.Rects[1].Fill, "Fill color for rects[1] should have been #123456")
	assert.Equal(t, "#123456", svg.Paths[0].Fill, "Fill color for paths[0] should have been #123456")
}

func TestPolyPathIcon_ReplaceFillColor(t *testing.T) {
	img := helperLoadBytes(t, "iconWithPolysAndPaths.svg")
	rdr := bytes.NewReader(img)
	var svg SVG
	if err := svg.ReplaceFillColor(rdr, "#123456"); err != nil {
		t.Fatal("ReplaceFillColor threw an error:", err)
	}
	assert.Equal(t, "#123456", svg.Paths[0].Fill, "Fill color for paths[0] should have been #123456")
	assert.Equal(t, "#123456", svg.Polygons[0].Fill, "Fill color for polygons[0] should have been #123456")
	assert.Equal(t, "#123456", svg.Polygons[1].Fill, "Fill color for polygons[1] should have been #123456")
}
