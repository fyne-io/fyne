package theme

import (
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"image/color"
	"io"
	"io/ioutil"
)

// SVG holds the unmarshaled XML from a Scalable Vector Graphic
type SVG struct {
	XMLName  xml.Name   `xml:"svg"`
	XMLNS    string     `xml:"xmlns,attr"`
	Width    string     `xml:"width,attr"`
	Height   string     `xml:"height,attr"`
	ViewBox  string     `xml:"viewBox,attr"`
	Paths    []*path    `xml:"path"`
	Rects    []*rect    `xml:"rect"`
	Polygons []*polygon `xml:"polygon"`
	Groups   []*group   `xml:"g"`
}

type path struct {
	XMLName xml.Name `xml:"path"`
	Fill    string   `xml:"fill,attr"`
	D       string   `xml:"d,attr"`
}

type rect struct {
	XMLName xml.Name `xml:"rect"`
	Fill    string   `xml:"fill,attr"`
	X       string   `xml:"x,attr"`
	Y       string   `xml:"y,attr"`
	Width   string   `xml:"width,attr"`
	Height  string   `xml:"height,attr"`
}

type polygon struct {
	XMLName xml.Name `xml:"polygon"`
	Fill    string   `xml:"fill,attr"`
	Points  string   `xml:"points,attr"`
}

type group struct {
	XMLName  xml.Name   `xml:"g"`
	Id       string     `xml:"id,attr"`
	Paths    []*path    `xml:"path"`
	Rects    []*rect    `xml:"rect"`
	Polygons []*polygon `xml:"polygon"`
}

func replacePathsFill(paths []*path, hexColor string) {
	for _, path := range paths {
		if path.Fill != "none" {
			path.Fill = hexColor
		}
	}
}

func replaceRectsFill(rects []*rect, hexColor string) {
	for _, rect := range rects {
		if rect.Fill != "none" {
			rect.Fill = hexColor
		}
	}
}

func replacePolygonsFill(polys []*polygon, hexColor string) {
	for _, poly := range polys {
		if poly.Fill != "none" {
			poly.Fill = hexColor
		}
	}
}

func replaceGroupObjectFill(groups []*group, hexColor string) {
	for _, grp := range groups {
		replacePathsFill(grp.Paths, hexColor)
		replaceRectsFill(grp.Rects, hexColor)
		replacePolygonsFill(grp.Polygons, hexColor)
	}
}

// ReplaceFillColor alters an SVG objects fill color.  Note that if an SVG with multiple fill
// colors is being operated upon, all fills will be converted to a single color.  Mostly used
// to recolor Icons to match the theme's IconColor.
func (s *SVG) ReplaceFillColor(reader io.Reader, color color.Color) error {
	bSlice, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	if err := xml.Unmarshal(bSlice, &s); err != nil {
		return err
	}

	replacePathsFill(s.Paths, colorToHexString(color))
	replaceRectsFill(s.Rects, colorToHexString(color))
	replacePolygonsFill(s.Polygons, colorToHexString(color))
	replaceGroupObjectFill(s.Groups, colorToHexString(color))

	return nil
}

func colorToHexString(color color.Color) string {
	r, g, b, _ := color.RGBA()
	cBytes := []byte{byte(r), byte(g), byte(b)}
	return fmt.Sprintf("#%s", hex.EncodeToString(cBytes))
}
