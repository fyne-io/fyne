package theme

import (
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"image/color"
	"io"
	"io/ioutil"
)

// svg holds the unmarshaled XML from a Scalable Vector Graphic
type svg struct {
	XMLName  xml.Name      `xml:"svg"`
	XMLNS    string        `xml:"xmlns,attr"`
	Width    string        `xml:"width,attr"`
	Height   string        `xml:"height,attr"`
	ViewBox  string        `xml:"viewBox,attr"`
	Paths    []*pathObj    `xml:"path"`
	Rects    []*rectObj    `xml:"rect"`
	Circles  []*circleObj  `xml:"circle"`
	Polygons []*polygonObj `xml:"polygon"`
	Groups   []*objGroup   `xml:"g"`
}

type pathObj struct {
	XMLName xml.Name `xml:"path"`
	Fill    string   `xml:"fill,attr"`
	D       string   `xml:"d,attr"`
}

type rectObj struct {
	XMLName xml.Name `xml:"rect"`
	Fill    string   `xml:"fill,attr"`
	X       string   `xml:"x,attr"`
	Y       string   `xml:"y,attr"`
	Width   string   `xml:"width,attr"`
	Height  string   `xml:"height,attr"`
}

type circleObj struct {
	XMLName xml.Name `xml:"circle"`
	Fill    string   `xml:"fill,attr"`
	CX      string   `xml:"cx,attr"`
	CY      string   `xml:"cy,attr"`
	R       string   `xml:"r,attr"`
}

type polygonObj struct {
	XMLName xml.Name `xml:"polygon"`
	Fill    string   `xml:"fill,attr"`
	Points  string   `xml:"points,attr"`
}

type objGroup struct {
	XMLName  xml.Name      `xml:"g"`
	ID       string        `xml:"id,attr"`
	Paths    []*pathObj    `xml:"path"`
	Rects    []*rectObj    `xml:"rect"`
	Polygons []*polygonObj `xml:"polygon"`
}

func replacePathsFill(paths []*pathObj, hexColor string) {
	for _, path := range paths {
		if path.Fill != "none" {
			path.Fill = hexColor
		}
	}
}

func replaceRectsFill(rects []*rectObj, hexColor string) {
	for _, rect := range rects {
		if rect.Fill != "none" {
			rect.Fill = hexColor
		}
	}
}

func replaceCirclesFill(circles []*circleObj, hexColor string) {
	for _, circle := range circles {
		if circle.Fill != "none" {
			circle.Fill = hexColor
		}
	}
}

func replacePolygonsFill(polys []*polygonObj, hexColor string) {
	for _, poly := range polys {
		if poly.Fill != "none" {
			poly.Fill = hexColor
		}
	}
}

func replaceGroupObjectFill(groups []*objGroup, hexColor string) {
	for _, grp := range groups {
		replacePathsFill(grp.Paths, hexColor)
		replaceRectsFill(grp.Rects, hexColor)
		replacePolygonsFill(grp.Polygons, hexColor)
	}
}

// replaceFillColor alters an svg objects fill color.  Note that if an svg with multiple fill
// colors is being operated upon, all fills will be converted to a single color.  Mostly used
// to recolor Icons to match the theme's IconColor.
func (s *svg) replaceFillColor(reader io.Reader, color color.Color) error {
	replacePathsFill(s.Paths, colorToHexString(color))
	replaceRectsFill(s.Rects, colorToHexString(color))
	replaceCirclesFill(s.Circles, colorToHexString(color))
	replacePolygonsFill(s.Polygons, colorToHexString(color))
	replaceGroupObjectFill(s.Groups, colorToHexString(color))
	return nil
}

func svgFromXML(reader io.Reader) (*svg, error) {
	var s svg
	bSlice, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	if err := xml.Unmarshal(bSlice, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func colorToHexString(color color.Color) string {
	r, g, b, _ := color.RGBA()
	cBytes := []byte{byte(r), byte(g), byte(b)}
	return fmt.Sprintf("#%s", hex.EncodeToString(cBytes))
}
