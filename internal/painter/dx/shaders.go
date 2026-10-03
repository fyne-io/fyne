//go:build windows && directx

// Package dx provides a full Fyne render implementation using Direct3D 11.
// It has no cgo: the COM calls go through syscall, and the HLSL below is
// compiled at startup by d3dcompiler_47.dll, which ships with Windows 8.1 and
// later.
//
// Every canvas primitive is drawn on the GPU: rectangles (with per-corner
// radii, strokes and shadows), circles, ellipses, arcs, regular and arbitrary
// polygons, Bezier curves, blur, lines, text, images, rasters and both
// gradient types. User canvas.Shader objects are supported through
// Shader.SourceHLSL; a shader that carries no HLSL source is not drawn.
//
// The shader sources live in shaders/, one file per primitive plus a shared
// prelude (common.hlsl) that declares the frame constants and the vertex
// output struct. All pixel positions in the shaders are top-origin: both
// SV_Position and the `bounds` uniform put y=0 at the top of the frame, and
// where a signed-distance computation needs a centred vector with y pointing
// up (the rounded-rectangle, ellipse and arc shaders, whose corner and radius
// selection is sign-sensitive), y is negated when forming that vector.
//
// The package builds only for Windows targets that opt in with the `directx`
// build tag.
package dx

import _ "embed"

// Raw sources. Common is prepended to the others rather than being #include'd:
// D3DCompile only resolves includes when the caller supplies an ID3DInclude COM
// handler, and concatenating is a great deal less machinery than implementing one.
var (
	//go:embed shaders/common.hlsl
	Common string

	//go:embed shaders/passthrough_2d.vs.hlsl
	vertPassthrough2D string
	//go:embed shaders/line.vs.hlsl
	vertLine string
	//go:embed shaders/glyph.vs.hlsl
	vertGlyph string

	//go:embed shaders/simple.ps.hlsl
	pixelSimple string
	//go:embed shaders/rectangle.ps.hlsl
	pixelRectangle string
	//go:embed shaders/round_rectangle.ps.hlsl
	pixelRoundRectangle string
	//go:embed shaders/ellipse.ps.hlsl
	pixelEllipse string
	//go:embed shaders/arc.ps.hlsl
	pixelArc string
	//go:embed shaders/regular_polygon.ps.hlsl
	pixelRegularPolygon string
	//go:embed shaders/arbitrary_polygon.ps.hlsl
	pixelArbitraryPolygon string
	//go:embed shaders/bezier_curve.ps.hlsl
	pixelBezierCurve string
	//go:embed shaders/blur.ps.hlsl
	pixelBlur string
	//go:embed shaders/line.ps.hlsl
	pixelLine string
	//go:embed shaders/glyph.ps.hlsl
	pixelGlyph string
)

// Compilable sources, each already carrying the shared prelude.
var (
	VertexPassthrough2D = Common + vertPassthrough2D
	VertexLine          = Common + vertLine
	VertexGlyph         = Common + vertGlyph

	PixelSimple         = Common + pixelSimple
	PixelRectangle      = Common + pixelRectangle
	PixelRoundRectangle = Common + pixelRoundRectangle
	PixelEllipse        = Common + pixelEllipse
	PixelArc            = Common + pixelArc
	PixelRegularPolygon = Common + pixelRegularPolygon

	PixelArbitraryPolygon = Common + pixelArbitraryPolygon
	PixelBezierCurve      = Common + pixelBezierCurve
	PixelBlur             = Common + pixelBlur
	PixelLine             = Common + pixelLine
	PixelGlyph            = Common + pixelGlyph
)
