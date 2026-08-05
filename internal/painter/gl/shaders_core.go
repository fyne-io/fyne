//go:build (!gles && !arm && !arm64 && !android && !ios && !mobile && !test_web_driver && !wasm) || (darwin && !mobile && !ios && !wasm && !test_web_driver)

package gl

import (
	_ "embed"

	"fyne.io/fyne/v2/canvas"
)

var (
	//go:embed shaders/arbitrary_polygon.frag
	shaderFragArbitraryPolygon []byte
	//go:embed shaders/arc.frag
	shaderFragArc []byte
	//go:embed shaders/bezier_curve.frag
	shaderFragBezierCurve []byte
	//go:embed shaders/blur.frag
	shaderFragBlur []byte
	//go:embed shaders/ellipse.frag
	shaderFragEllipse []byte
	//go:embed shaders/line.frag
	shaderFragLine []byte
	//go:embed shaders/regular_polygon.frag
	shaderFragPolygon []byte
	//go:embed shaders/rectangle.frag
	shaderFragRectangle []byte
	//go:embed shaders/round_rectangle.frag
	shaderFragRoundRectangle []byte
	//go:embed shaders/simple.frag
	shaderFragSimple []byte

	//go:embed shaders/line.vert
	shaderVertLine []byte
	//go:embed shaders/passthrough_2d.vert
	shaderVertPassthrough2D []byte
	//go:embed shaders/textured_passthrough_2d.vert
	shaderVertTexturedPassthrough2D []byte
)

// userShaderFragment returns the fragment shader source to use for the given
// shader object on this build target.
func userShaderFragment(s *canvas.Shader) []byte {
	return s.Source
}
