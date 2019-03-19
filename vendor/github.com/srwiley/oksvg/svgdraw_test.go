// Copyright 2018 The oksvg Authors. All rights reserved.
// created: 2018 by S.R.Wiley
package oksvg_test

import (
	"bufio"
	"fmt"
	"image"
	"os"

	"image/png"
	"strings"
	"testing"

	. "github.com/srwiley/oksvg"
	. "github.com/srwiley/rasterx"
	//"github.com/srwiley/go/scanFT"
)

const testArco = `M150,350 l 50,-55 
           a25,25 -30 0,1 50,-25 l 50,-25 
           a25,50 -30 0,1 50,-25 l 50,-25 
           a25,75 -30 0,1 50,-25 l 50,-25 
           a25,100 -30 0,1 50,-25 l 50,15z`

const testArco2 = `M150,350 l 50,-55 
           a35,25 -30 0,0 50,-25 l 50,-25 
           a25,50 -30 0,1 50,-25 l 50,-25 
           a25,75 -30 0,1 50,-25 l 50,-25 
           a25,100 -30 0,1 50,-25, l 50,15z`

const testArcoS = `M150,350 l 50,-55 
           a35,25 -30 0,0 50,-25,
           25,50 -30 0,1 50,-25
           a25,75 -30 0,1 50,-25 l 50,-25 
           a25,100 -30 0,1 50,-25 l 50,15,0,25,-15,-15  z`

// Explicitly call each command in abs and rel mode and concatenated forms
const testSVG0 = `m20,20,0,400,400,0z`
const testSVG1 = `M20,20 L500,800 L800,200z`
const testSVG2 = `M20,20 Q200,800 800,800z`
const testSVG3 = `M20,50 C200,200 800,200 800,500z`
const testSVG4 = `M20,50 S200,1400 400,500 S700,800 800,400z`
const testSVG5 = `M50,20 Q 800,500 500,800z`
const testSVG6 = `M20,50 c200,200 800,200 400,300z`
const testSVG7 = `M20,20 c0,500 500,0 500,500z`
const testSVG8 = `M20,50 c200,200 800,200 400,300c200,200 800,200 400,300z`
const testSVG9 = `M20,50 c200,200 800,200 400,300,200,200 800,200 400,300z`
const testSVG10 = `M20,50 c200,200 800,200 400,300,200,200 800,200 400,300s500,300 200,200s600,300 200,200z`
const testSVG11 = `M20,50 c200,200 800,200 400,300,200,200 800,200 400,300s500,300 200,200,600,300 200,200z`
const testSVG12 = `M100,100 Q400,100 250,250 T400,400z`
const testSVG13 = `M100,100 Q400,100 250,250 t150,150,150,150z`

func TestTransform(t *testing.T) {
	icon, errSvg := ReadIcon("testdata/landscapeIcons/sea.svg", WarnErrorMode)
	if errSvg != nil {
		t.Error(errSvg)
	}
	w, h := int(icon.ViewBox.W), int(icon.ViewBox.H)
	img := image.NewRGBA(image.Rect(0, 0, w*3, h*3))

	scannerGV := NewScannerGV(w*3, h*3, img, img.Bounds())

	raster := NewDasher(w*3, h*3, scannerGV)
	icon.Draw(raster, 1.0)
	icon.Transform = Identity.Translate(float64(w), float64(h))
	icon.Draw(raster, 1.0)

	icon.SetTarget(float64(w), float64(0), float64(w), float64(h)*.5)
	icon.Draw(raster, 1.0)

	err := SaveToPngFile(fmt.Sprintf("testdata/transform.png"), img)
	if err != nil {
		t.Error(err)
	}
}

func DrawIcon(t *testing.T, iconPath string) image.Image {
	icon, errSvg := ReadIcon(iconPath, WarnErrorMode)
	if errSvg != nil {
		t.Error(errSvg)
		return nil
	}
	w, h := int(icon.ViewBox.W), int(icon.ViewBox.H)
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	// Uncomment the next three lines and comment the three after to use ScannerFT
	//	painter := scanFT.NewRGBAPainter(img)
	//	scannerFT := scanFT.NewScannerFT(w, h, painter)
	//	raster := NewDasher(w, h, scannerFT)
	//tb := img.Bounds()
	//tb.Max.X /= 2
	scannerGV := NewScannerGV(w, h, img, img.Bounds())
	raster := NewDasher(w, h, scannerGV)

	icon.Draw(raster, 1.0)
	return img
}

func SaveIcon(t *testing.T, iconPath string) {
	img := DrawIcon(t, iconPath)
	if img != nil {
		p := strings.Split(iconPath, "/")
		err := SaveToPngFile(fmt.Sprintf("testdata/%s.png", p[len(p)-1]), img)
		if err != nil {
			t.Error(err)
		}
	}
}

func SaveToPngFile(filePath string, m image.Image) error {
	// Create the file
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	// Create Writer from file
	b := bufio.NewWriter(f)
	// Write the image into the buffer
	err = png.Encode(b, m)
	if err != nil {
		return err
	}
	err = b.Flush()
	if err != nil {
		return err
	}
	return nil
}

func TestSvgPathsStroke(t *testing.T) {
	for i, p := range []string{testArco, testArco2, testArcoS,
		testSVG0, testSVG1, testSVG2, testSVG3, testSVG4, testSVG5,
		testSVG6, testSVG7, testSVG8, testSVG9, testSVG10,
		testSVG11, testSVG12, testSVG13,
	} {
		w := 1600
		img := image.NewRGBA(image.Rect(0, 0, w, w))

		scannerGV := NewScannerGV(w, w, img, img.Bounds())
		raster := NewDasher(w, w, scannerGV)

		c := &PathCursor{}
		d := DefaultStyle
		icon := SvgIcon{}

		err := c.CompilePath(p)
		if err != nil {
			t.Error(err)
		}
		icon.SVGPaths = append(icon.SVGPaths, SvgPath{PathStyle: d, Path: c.Path})
		icon.Draw(raster, 1)

		err = SaveToPngFile(fmt.Sprintf("testdata/fill_%d.png", i), img)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestLandscapeIcons(t *testing.T) {
	for _, p := range []string{
		"beach", "cape", "iceberg", "island",
		"mountains", "sea", "trees", "village"} {
		SaveIcon(t, "testdata/landscapeIcons/"+p+".svg")
	}
}

func TestTestIcons(t *testing.T) {
	for _, p := range []string{
		"astronaut", "jupiter", "lander", "school-bus", "telescope", "content-cut-light"} {
		SaveIcon(t, "testdata/testIcons/"+p+".svg")
	}
}

func TestStrokeIcons(t *testing.T) {
	for _, p := range []string{
		"OpacityStrokeDashTest.svg",
		"OpacityStrokeDashTest2.svg",
		"OpacityStrokeDashTest3.svg",
		"TestShapes.svg",
		"TestShapes2.svg",
		"TestShapes3.svg",
		"TestShapes4.svg",
		"TestShapes5.svg",
		"TestShapes6.svg",
	} {
		t.Log("reading ", p)
		SaveIcon(t, "testdata/"+p)
	}
}
