package svg

import (
	"bytes"
	"encoding/xml"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"fyne.io/fyne/v2/internal/test"
)

func TestColorize(t *testing.T) {
	tests := map[string]struct {
		svgFile   string
		color     color.Color
		wantImage string
	}{
		"paths": {
			svgFile:   "cancel_Paths.svg",
			color:     color.NRGBA{R: 100, G: 100, A: 200},
			wantImage: "colorized/paths.png",
		},
		"circles": {
			svgFile:   "circles.svg",
			color:     color.NRGBA{R: 100, B: 100, A: 200},
			wantImage: "colorized/circles.png",
		},
		"polygons": {
			svgFile:   "polygons.svg",
			color:     color.NRGBA{G: 100, B: 100, A: 200},
			wantImage: "colorized/polygons.png",
		},
		"rects": {
			svgFile:   "rects.svg",
			color:     color.NRGBA{R: 100, G: 100, B: 100, A: 200},
			wantImage: "colorized/rects.png",
		},
		"negative rects": {
			svgFile:   "rects-negative.svg",
			color:     color.NRGBA{R: 100, G: 100, B: 100, A: 200},
			wantImage: "colorized/rects.png",
		},
		"group of paths": {
			svgFile:   "check_GroupPaths.svg",
			color:     color.NRGBA{R: 100, G: 100, A: 100},
			wantImage: "colorized/group_paths.png",
		},
		"group of circles": {
			svgFile:   "group_circles.svg",
			color:     color.NRGBA{R: 100, B: 100, A: 100},
			wantImage: "colorized/group_circles.png",
		},
		"group of polygons": {
			svgFile:   "warning_GroupPolygons.svg",
			color:     color.NRGBA{G: 100, B: 100, A: 100},
			wantImage: "colorized/group_polygons.png",
		},
		"group of rects": {
			svgFile:   "info_GroupRects.svg",
			color:     color.NRGBA{R: 100, G: 100, B: 100, A: 100},
			wantImage: "colorized/group_rects.png",
		},
		"NRGBA64": {
			svgFile: "circles.svg",
			// If the low 8 bits of each component were used, this would look cyan instead of yellow.
			// When the MSB is used instead, it correctly looks yellow.
			color:     color.NRGBA64{R: 0xff00, G: 0xffff, B: 0x00ff, A: 0xffff},
			wantImage: "colorized/circles_yellow.png",
		},
		"translucent NRGBA64": {
			svgFile:   "circles.svg",
			color:     color.NRGBA64{R: 0xff00, G: 0xffff, B: 0x00ff, A: 0x7fff},
			wantImage: "colorized/circles_yellow_translucent.png",
		},
		"RGBA": {
			svgFile:   "circles.svg",
			color:     color.RGBAModel.Convert(color.NRGBA{R: 0xff, G: 0xff, B: 0x00, A: 0xff}),
			wantImage: "colorized/circles_yellow.png",
		},
		"transluscent RGBA": {
			svgFile:   "circles.svg",
			color:     color.RGBAModel.Convert(color.NRGBA{R: 0xff, G: 0xff, B: 0x00, A: 0x7f}),
			wantImage: "colorized/circles_yellow_translucent.png",
		},
		"RGBA64": {
			svgFile: "circles.svg",
			// If the least significant byte of each component was being used, this would look cyan instead of yellow.
			// Since alpha=0xffff, unmultiplyAlpha knows it does not need to unmultiply anything, and so it just
			// returns the MSB of each component.
			color:     color.RGBA64Model.Convert(color.NRGBA64{R: 0xff00, G: 0xffff, B: 0x00ff, A: 0xffff}),
			wantImage: "colorized/circles_yellow.png",
		},
		"transluscent RGBA64": {
			svgFile: "circles.svg",
			// Since alpha!=0xffff, if we were to use R:0xff00, G:0xffff, B:0x00ff like before,
			// this would end up being drawn with 0xfeff00 instead of 0xffff00, and we would need a separate image to test for that.
			// Instead, we use R:0xfff0, G:0xfff0, B:0x000f, A:0x7fff instead, which unmultiplyAlpha returns as 0xff, 0xff, 0x00, 0x7f,
			// so that we correctly get 0xffff00 with alpha 0x7f when ToRGBA is used.
			// The RGBA64's contents are 0x7ff7, 0x7ff7, 0x0007, 0x7fff, so:
			// If ToRGBA wasn't being called and instead the LSB of each component was being read, this would show up as 0xf7f707 with alpha 0xff.
			// If the MSB was being read without umultiplication, this would show up as 0x7f7f00 with alpha 0x7f.
			color:     color.RGBA64Model.Convert(color.NRGBA64{R: 0xfff0, G: 0xfff0, B: 0x000f, A: 0x7fff}),
			wantImage: "colorized/circles_yellow_translucent.png",
		},
		"Alpha": {
			svgFile:   "circles.svg",
			color:     color.Alpha{A: 0x7f},
			wantImage: "colorized/circles_white_translucent.png",
		},
		"Alpha16": {
			svgFile: "circles.svg",
			// If the LSB from components returned by RGBA() was being used, this would be black.
			// If the MSB from components returned by RGBA() was being used, this would be grey.
			// It is white when either we bypass RGBA() and directly make a 0xffffff color with the alpha's MSB (which is what ToRGBA does),
			// or if we call RBGA(), un-multiply the alpha from the non-alpha components, and use their MSB to get white (Or something very near it like 0xfefefe).
			color:     color.Alpha16{A: 0x7f00},
			wantImage: "colorized/circles_white_translucent.png",
		},
		"Gray": {
			svgFile:   "circles.svg",
			color:     color.Gray{Y: 0xff},
			wantImage: "colorized/circles_white.png",
		},
		"Gray16": {
			svgFile: "circles.svg",
			// If the LSB from components returned by RGBA() was being used, this would be black.
			// It is white when either we bypass RGBA() and directly make a 0xffffff color with the alpha's MSB (which is what ToRGBA does),
			// or if we call RBGA(), un-multiply the alpha from the non-alpha components, and use their MSB to get white (Or something very near it like 0xfefefe),
			// or if the MSB from components returned by RGBA() was being used (because Gray and Gray16 do not have alpha values).
			color:     color.Gray16{Y: 0xff00},
			wantImage: "colorized/circles_white.png",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			bytes, err := os.ReadFile(filepath.Join("testdata", tt.svgFile))
			require.NoError(t, err)
			content, _ := Colorize(bytes, tt.color)
			got := helperDrawSVG(t, content)
			test.AssertImageMatches(t, tt.wantImage, got)
		})
	}
}

func TestSVG_ReplaceFillColor(t *testing.T) {
	src, err := os.ReadFile("testdata/cancel_Paths.svg")
	if err != nil {
		t.Fatal(err)
	}
	red := color.NRGBA{0xff, 0x00, 0x00, 0xff}
	rdr := bytes.NewReader(src)
	s, err := svgFromXML(rdr)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.replaceFillColor(red); err != nil {
		t.Fatal(err)
	}
	res, err := xml.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEqual(t, string(src), string(res))
	assert.Contains(t, string(res), "#ff0000")
}

func TestSVG_ReplaceFillColor_Ellipse(t *testing.T) {
	src, err := os.ReadFile("testdata/ellipse.svg")
	if err != nil {
		t.Fatal(err)
	}
	red := color.NRGBA{0xff, 0x00, 0x00, 0xff}
	rdr := bytes.NewReader(src)
	s, err := svgFromXML(rdr)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.replaceFillColor(red); err != nil {
		t.Fatal(err)
	}
	res, err := xml.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEqual(t, string(src), string(res))
	assert.Contains(t, string(res), "#ff0000")
}

func helperDrawSVG(t *testing.T, data []byte) image.Image {
	icon, err := oksvg.ReadIconStream(bytes.NewReader(data))
	require.NoError(t, err, "failed to read SVG data")

	width := int(icon.ViewBox.W) * 2
	height := int(icon.ViewBox.H) * 2
	x, y := svgOffset(icon, width, height)
	icon.SetTarget(x, y, float64(width), float64(height))
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	scanner := rasterx.NewScannerGV(width, height, img, img.Bounds())
	raster := rasterx.NewDasher(width, height, scanner)
	icon.Draw(raster, 1)
	return img
}
