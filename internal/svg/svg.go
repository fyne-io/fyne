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
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fyne-io/oksvg"
	"github.com/srwiley/rasterx"

	"fyne.io/fyne/v2"
	col "fyne.io/fyne/v2/internal/color"
)

// Colorize creates a new SVG from a given one by replacing all fill colors by the given color.
func Colorize(src []byte, clr color.Color) ([]byte, error) {
	rdr := bytes.NewReader(src)
	s, err := svgFromXML(rdr)
	if err != nil {
		return src, fmt.Errorf("could not load SVG, falling back to static content: %v", err)
	}
	if err := s.replaceFillColor(clr); err != nil {
		return src, fmt.Errorf("could not replace fill color, falling back to static content: %v", err)
	}
	colorized, err := xml.Marshal(s)
	if err != nil {
		return src, fmt.Errorf("could not marshal svg, falling back to static content: %v", err)
	}
	return colorized, nil
}

type Decoder struct {
	icon *oksvg.SvgIcon
}

type Config struct {
	Width  int
	Height int
	Aspect float32
}

func NewDecoder(stream io.Reader) (*Decoder, error) {
	icon, err := oksvg.ReadIconStream(stream)
	if err != nil {
		return nil, err
	}

	return &Decoder{
		icon: icon,
	}, nil
}

func (d *Decoder) Config() Config {
	return Config{
		int(d.icon.ViewBox.W),
		int(d.icon.ViewBox.H),
		float32(d.icon.ViewBox.W / d.icon.ViewBox.H),
	}
}

func (d *Decoder) Draw(width, height int) (*image.NRGBA, error) {
	config := d.Config()

	viewAspect := float32(width) / float32(height)
	imgW, imgH := width, height
	if viewAspect > config.Aspect {
		imgW = int(float32(height) * config.Aspect)
	} else if viewAspect < config.Aspect {
		imgH = int(float32(width) / config.Aspect)
	}

	x, y := svgOffset(d.icon, imgW, imgH)
	d.icon.SetTarget(x, y, float64(imgW), float64(imgH))

	img := image.NewNRGBA(image.Rect(0, 0, imgW, imgH))
	scanner := rasterx.NewScannerGV(config.Width, config.Height, img, img.Bounds())
	raster := rasterx.NewDasher(width, height, scanner)

	err := drawSVGSafely(d.icon, raster)
	if err != nil {
		err = fmt.Errorf("SVG render error: %w", err)
		return nil, err
	}
	return img, nil
}

func IsFileSVG(path string) bool {
	return strings.EqualFold(filepath.Ext(path), ".svg")
}

// IsResourceSVG checks if the resource is an SVG or not.
func IsResourceSVG(res fyne.Resource) bool {
	if IsFileSVG(res.Name()) {
		return true
	}

	if len(res.Content()) < 5 {
		return false
	}

	switch strings.ToLower(string(res.Content()[:5])) {
	case "<!doc", "<?xml", "<svg ":
		return true
	}
	return false
}

func svgOffset(icon *oksvg.SvgIcon, _, height int) (x, y float64) {
	if icon.ViewBox.Y < 0 { // adjust so our positive offset calculations work
		y = icon.ViewBox.Y + (-icon.ViewBox.Y/icon.ViewBox.H)*float64(height)
	}

	return 0, y
}

// svg holds the unmarshaled XML from a Scalable Vector Graphic
type svg struct {
	XMLName  xml.Name      `xml:"svg"`
	XMLNS    string        `xml:"xmlns,attr"`
	Width    string        `xml:"width,attr,omitempty"`
	Height   string        `xml:"height,attr,omitempty"`
	ViewBox  string        `xml:"viewBox,attr,omitempty"`
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
	Transform       string   `xml:"transform,attr,omitempty"`
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
	Transform       string   `xml:"transform,attr,omitempty"`
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
	Transform       string   `xml:"transform,attr,omitempty"`
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
	Transform       string   `xml:"transform,attr,omitempty"`
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
	Transform       string   `xml:"transform,attr,omitempty"`
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
	Transform       string        `xml:"transform,attr,omitempty"`
	Paths           []*pathObj    `xml:"path"`
	Circles         []*circleObj  `xml:"circle"`
	Ellipses        []*ellipseObj `xml:"ellipse"`
	Rects           []*rectObj    `xml:"rect"`
	Polygons        []*polygonObj `xml:"polygon"`
	Groups          []*objGroup   `xml:"g"`
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
		replaceGroupObjectFill(grp.Groups, hexColor, opacity)
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
	bSlice, err := io.ReadAll(reader)
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
