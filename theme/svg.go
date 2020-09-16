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
	XMLName         xml.Name `xml:"path"`
	Fill            string   `xml:"fill,attr,omitempty"`
	FillOpacity     string   `xml:"fill-opacity,attr,omitempty"`
	Stroke          string   `xml:"stroke,attr,omitempty"`
	StrokeWidth     string   `xml:"stroke-width,attr,omitempty"`
	StrokeLineCap   string   `xml:"stroke-linecap,attr,omitempty"`
	StrokeLineJoin  string   `xml:"stroke-linejoin,attr,omitempty"`
	StrokeDashArray string   `xml:"stroke-dasharray,attr,omitempty"`
	D               string   `xml:"d,attr"`
}

type rectObj struct {
	XMLName         xml.Name `xml:"rect"`
	Fill            string   `xml:"fill,attr,omitempty"`
	FillOpacity     string   `xml:"fill-opacity,attr,omitempty"`
	Stroke          string   `xml:"stroke,attr,omitempty"`
	StrokeWidth     string   `xml:"stroke-width,attr,omitempty"`
	StrokeLineCap   string   `xml:"stroke-linecap,attr,omitempty"`
	StrokeLineJoin  string   `xml:"stroke-linejoin,attr,omitempty"`
	StrokeDashArray string   `xml:"stroke-dasharray,attr,omitempty"`
	X               string   `xml:"x,attr,omitempty"`
	Y               string   `xml:"y,attr,omitempty"`
	Width           string   `xml:"width,attr,omitempty"`
	Height          string   `xml:"height,attr,omitempty"`
}

type circleObj struct {
	XMLName         xml.Name `xml:"circle"`
	Fill            string   `xml:"fill,attr,omitempty"`
	FillOpacity     string   `xml:"fill-opacity,attr,omitempty"`
	Stroke          string   `xml:"stroke,attr,omitempty"`
	StrokeWidth     string   `xml:"stroke-width,attr,omitempty"`
	StrokeLineCap   string   `xml:"stroke-linecap,attr,omitempty"`
	StrokeLineJoin  string   `xml:"stroke-linejoin,attr,omitempty"`
	StrokeDashArray string   `xml:"stroke-dasharray,attr,omitempty"`
	CX              string   `xml:"cx,attr,omitempty"`
	CY              string   `xml:"cy,attr,omitempty"`
	R               string   `xml:"r,attr,omitempty"`
}

type polygonObj struct {
	XMLName         xml.Name `xml:"polygon"`
	Fill            string   `xml:"fill,attr,omitempty"`
	FillOpacity     string   `xml:"fill-opacity,attr,omitempty"`
	Stroke          string   `xml:"stroke,attr,omitempty"`
	StrokeWidth     string   `xml:"stroke-width,attr,omitempty"`
	StrokeLineCap   string   `xml:"stroke-linecap,attr,omitempty"`
	StrokeLineJoin  string   `xml:"stroke-linejoin,attr,omitempty"`
	StrokeDashArray string   `xml:"stroke-dasharray,attr,omitempty"`
	Points          string   `xml:"points,attr"`
}

type objGroup struct {
	XMLName         xml.Name      `xml:"g"`
	ID              string        `xml:"id,attr,omitempty"`
	Fill            string        `xml:"fill,attr,omitempty"`
	Stroke          string        `xml:"stroke,attr,omitempty"`
	StrokeWidth     string        `xml:"stroke-width,attr,omitempty"`
	StrokeLineCap   string        `xml:"stroke-linecap,attr,omitempty"`
	StrokeLineJoin  string        `xml:"stroke-linejoin,attr,omitempty"`
	StrokeDashArray string        `xml:"stroke-dasharray,attr,omitempty"`
	Paths           []*pathObj    `xml:"path"`
	Circles         []*circleObj  `xml:"circle"`
	Rects           []*rectObj    `xml:"rect"`
	Polygons        []*polygonObj `xml:"polygon"`
}

func replacePathsFill(paths []*pathObj, hexColor string, opacity string) {
	for _, path := range paths {
		if path.Fill != "none" {
			path.Fill = hexColor
			path.FillOpacity = opacity
		}
	}
}

func replaceRectsFill(rects []*rectObj, hexColor string, opacity string) {
	for _, rect := range rects {
		if rect.Fill != "none" {
			rect.Fill = hexColor
			rect.FillOpacity = opacity
		}
	}
}

func replaceCirclesFill(circles []*circleObj, hexColor string, opacity string) {
	for _, circle := range circles {
		if circle.Fill != "none" {
			circle.Fill = hexColor
			circle.FillOpacity = opacity
		}
	}
}

func replacePolygonsFill(polys []*polygonObj, hexColor string, opacity string) {
	for _, poly := range polys {
		if poly.Fill != "none" {
			poly.Fill = hexColor
			poly.FillOpacity = opacity
		}
	}
}

func replaceGroupObjectFill(groups []*objGroup, hexColor string, opacity string) {
	for _, grp := range groups {
		replaceCirclesFill(grp.Circles, hexColor, opacity)
		replacePathsFill(grp.Paths, hexColor, opacity)
		replaceRectsFill(grp.Rects, hexColor, opacity)
		replacePolygonsFill(grp.Polygons, hexColor, opacity)
	}
}

// replaceFillColor alters an svg objects fill color.  Note that if an svg with multiple fill
// colors is being operated upon, all fills will be converted to a single color.  Mostly used
// to recolor Icons to match the theme's IconColor.
func (s *svg) replaceFillColor(color color.Color) error {
	hexColor, opacity := colorToHexAndOpacity(color)
	replacePathsFill(s.Paths, hexColor, opacity)
	replaceRectsFill(s.Rects, hexColor, opacity)
	replaceCirclesFill(s.Circles, hexColor, opacity)
	replacePolygonsFill(s.Polygons, hexColor, opacity)
	replaceGroupObjectFill(s.Groups, hexColor, opacity)
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

func colorToHexAndOpacity(color color.Color) (string, string) {
	r, g, b, a := color.RGBA()
	cBytes := []byte{byte(r), byte(g), byte(b)}
	return fmt.Sprintf("#%s", hex.EncodeToString(cBytes)), fmt.Sprintf("%f", float64(a)/0xffff)
}
