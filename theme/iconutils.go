package theme

import (
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"image/color"
	"io"
	"io/ioutil"
)

// SVG represents a Scalable Vector Graphic (SVG)
type SVG struct {
	XMLName  xml.Name   `xml:"svg"`
	XMLNS    string     `xml:"xmlns,attr"`
	Width    string     `xml:"width,attr"`
	Height   string     `xml:"height,attr"`
	ViewBox  string     `xml:"viewBox,attr"`
	Paths    []*Path    `xml:"path"`
	Rects    []*Rect    `xml:"rect"`
	Polygons []*Polygon `xml:"polygon"`
}

// Path represents path objects of an SVG image
type Path struct {
	XMLName xml.Name `xml:"path"`
	Fill    string   `xml:"fill,attr"`
	D       string   `xml:"d,attr"`
}

// Rect stores rect objects of an SVG image
type Rect struct {
	XMLName xml.Name `xml:"rect"`
	Fill    string   `xml:"fill,attr"`
	X       string   `xml:"x,attr"`
	Y       string   `xml:"y,attr"`
	Width   string   `xml:"width,attr"`
	Height  string   `xml:"height,attr"`
}

// Polygon stores polygon objects of an SVG image
type Polygon struct {
	XMLName xml.Name `xml:"polygon"`
	Fill    string   `xml:"fill,attr"`
	Points  string   `xml:"points,attr"`
}

func (svg *SVG) replacePathsFill(hexColor string) {
	for _, path := range svg.Paths {
		if path.Fill != "none" {
			path.Fill = hexColor
		}
	}
}

func (svg *SVG) replaceRectsFill(hexColor string) {
	for _, rect := range svg.Rects {
		if rect.Fill != "none" {
			rect.Fill = hexColor
		}
	}
}

func (svg *SVG) replacePolygonsFill(hexColor string) {
	for _, poly := range svg.Polygons {
		if poly.Fill != "none" {
			poly.Fill = hexColor
		}
	}
}

// ReplaceFillColor changes an SVG fill color to the `hexColor` parameter's argument
func (svg *SVG) ReplaceFillColor(reader io.Reader, hexColor string) error {
	bSlice, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	if err := xml.Unmarshal(bSlice, &svg); err != nil {
		return err
	}

	svg.replacePathsFill(hexColor)
	svg.replaceRectsFill(hexColor)
	svg.replacePolygonsFill(hexColor)

	return nil
}

// ColorToHexString returns a hex color string (i.e. #ffffff) for the color parameter's argument
func ColorToHexString(color color.Color) string {
	r, g, b, _ := color.RGBA()
	cBytes := []byte{byte(r), byte(g), byte(b)}
	return fmt.Sprintf("#%s", hex.EncodeToString(cBytes))
}
