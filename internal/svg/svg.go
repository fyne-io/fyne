package svg

import (
	"bytes"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"io/ioutil"
	"strconv"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"

	"fyne.io/fyne/v2"
	col "fyne.io/fyne/v2/internal/color"
)

// Colorize creates a new SVG from a given one by replacing all fill colors by the given color.
func Colorize(src []byte, clr color.Color) []byte {
	rdr := bytes.NewReader(src)
	s, err := svgFromXML(rdr)
	if err != nil {
		fyne.LogError("could not load SVG, falling back to static content:", err)
		return src
	}
	if err := s.replaceFillColor(clr); err != nil {
		fyne.LogError("could not replace fill color, falling back to static content:", err)
		return src
	}
	colorized, err := xml.Marshal(s)
	if err != nil {
		fyne.LogError("could not marshal svg, falling back to static content:", err)
		return src
	}
	return colorized
}

// ToImage reads an SVG from an io.Reader and renders it into an image.NRGBA using the requested width and height.
// The optional `validateSize` callback can be used to cancel the rendering depending on the SVG’s original size.
// In this case `nil` is returned.
func ToImage(file io.Reader, width, height int, validateSize func(origW, origH int) bool) (*image.NRGBA, error) {
	icon, err := oksvg.ReadIconStream(file)
	if err != nil {
		return nil, fmt.Errorf("SVG Load error: %w", err)
	}

	origW, origH := int(icon.ViewBox.W), int(icon.ViewBox.H)
	if validateSize != nil && !validateSize(origW, origH) {
		return nil, nil
	}

	aspect := float32(origW) / float32(origH)
	viewAspect := float32(width) / float32(height)
	imgW, imgH := width, height
	if viewAspect > aspect {
		imgW = int(float32(height) * aspect)
	} else if viewAspect < aspect {
		imgH = int(float32(width) / aspect)
	}
	icon.SetTarget(0, 0, float64(imgW), float64(imgH))

	img := image.NewNRGBA(image.Rect(0, 0, imgW, imgH))
	scanner := rasterx.NewScannerGV(origW, origH, img, img.Bounds())
	raster := rasterx.NewDasher(width, height, scanner)

	err = drawSVGSafely(icon, raster)
	if err != nil {
		err = fmt.Errorf("SVG render error: %w", err)
		return nil, err
	}
	return img, nil
}

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
	Ellipses []*ellipseObj `xml:"ellipse"`
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

type ellipseObj struct {
	XMLName         xml.Name `xml:"ellipse"`
	Fill            string   `xml:"fill,attr,omitempty"`
	FillOpacity     string   `xml:"fill-opacity,attr,omitempty"`
	Stroke          string   `xml:"stroke,attr,omitempty"`
	StrokeWidth     string   `xml:"stroke-width,attr,omitempty"`
	StrokeLineCap   string   `xml:"stroke-linecap,attr,omitempty"`
	StrokeLineJoin  string   `xml:"stroke-linejoin,attr,omitempty"`
	StrokeDashArray string   `xml:"stroke-dasharray,attr,omitempty"`
	CX              string   `xml:"cx,attr,omitempty"`
	CY              string   `xml:"cy,attr,omitempty"`
	RX              string   `xml:"rx,attr,omitempty"`
	RY              string   `xml:"ry,attr,omitempty"`
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
	Ellipses        []*ellipseObj `xml:"ellipse"`
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

func replaceEllipsesFill(ellipses []*ellipseObj, hexColor string, opacity string) {
	for _, ellipse := range ellipses {
		if ellipse.Fill != "none" {
			ellipse.Fill = hexColor
			ellipse.FillOpacity = opacity
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
		replaceEllipsesFill(grp.Ellipses, hexColor, opacity)
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
	replaceEllipsesFill(s.Ellipses, hexColor, opacity)
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

func colorToHexAndOpacity(color color.Color) (hexStr, aStr string) {
	r, g, b, a := col.ToNRGBA(color)
	cBytes := []byte{byte(r), byte(g), byte(b)}
	hexStr, aStr = "#"+hex.EncodeToString(cBytes), strconv.FormatFloat(float64(a)/0xff, 'f', 6, 64)
	return
}

func drawSVGSafely(icon *oksvg.SvgIcon, raster *rasterx.Dasher) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("crash when rendering svg")
		}
	}()
	icon.Draw(raster, 1)

	return err
}
