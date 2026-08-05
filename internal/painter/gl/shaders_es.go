//go:build ((gles || arm || arm64) && !android && !ios && !mobile && !darwin && !wasm && !test_web_driver) || ((android || ios || mobile) && (!wasm || !test_web_driver)) || wasm || test_web_driver

package gl

import (
	_ "embed"

	"fyne.io/fyne/v2/canvas"
)

var (
	//go:embed shaders/arbitrary_polygon_es.frag
	shaderFragArbitraryPolygon []byte
	//go:embed shaders/arc_es.frag
	shaderFragArc []byte
	//go:embed shaders/bezier_curve_es.frag
	shaderFragBezierCurve []byte
	//go:embed shaders/blur_es.frag
	shaderFragBlur []byte
	//go:embed shaders/ellipse_es.frag
	shaderFragEllipse []byte
	//go:embed shaders/line_es.frag
	shaderFragLine []byte
	//go:embed shaders/regular_polygon_es.frag
	shaderFragPolygon []byte
	//go:embed shaders/rectangle_es.frag
	shaderFragRectangle []byte
	//go:embed shaders/round_rectangle_es.frag
	shaderFragRoundRectangle []byte
	//go:embed shaders/simple_es.frag
	shaderFragSimple []byte

	//go:embed shaders/line_es.vert
	shaderVertLine []byte
	//go:embed shaders/passthrough_2d_es.vert
	shaderVertPassthrough2D []byte
	//go:embed shaders/textured_passthrough_2d_es.vert
	shaderVertTexturedPassthrough2D []byte
)

// userShaderFragment returns the fragment shader source to use for the given
// shader object on this build target.
func userShaderFragment(s *canvas.Shader) []byte {
	return s.SourceES
}
