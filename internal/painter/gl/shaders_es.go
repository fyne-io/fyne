//go:build ((gles || arm || arm64) && !android && !ios && !mobile && !darwin && !wasm && !test_web_driver) || ((android || ios || mobile) && (!wasm || !test_web_driver)) || wasm || test_web_driver

package gl

import (
	_ "embed"

	"fyne.io/fyne/v2/canvas"
)

var (
	//go:embed shaders/blur_es.frag
	shaderBluresFrag []byte

	//go:embed shaders/blur_es.vert
	shaderBluresVert []byte

	//go:embed shaders/line_es.frag
	shaderLineesFrag []byte

	//go:embed shaders/line_es.vert
	shaderLineesVert []byte

	//go:embed shaders/rectangle_es.frag
	shaderRectangleesFrag []byte

	//go:embed shaders/rectangle_es.vert
	shaderRectangleesVert []byte

	//go:embed shaders/round_rectangle_es.frag
	shaderRoundrectangleesFrag []byte

	//go:embed shaders/simple_es.frag
	shaderSimpleesFrag []byte

	//go:embed shaders/simple_es.vert
	shaderSimpleesVert []byte

	//go:embed shaders/regular_polygon_es.frag
	shaderPolygonesFrag []byte

	//go:embed shaders/arc_es.frag
	shaderArcesFrag []byte

	//go:embed shaders/bezier_curve_es.frag
	shaderBezierCurveesFrag []byte

	//go:embed shaders/arbitrary_polygon_es.frag
	shaderArbitraryPolygonesFrag []byte
)

func shaderSourceNamed(name string) ([]byte, []byte) {
	switch name {
	case "blur_es":
		return shaderBluresVert, shaderBluresFrag
	case "line_es":
		return shaderLineesVert, shaderLineesFrag
	case "simple_es":
		return shaderSimpleesVert, shaderSimpleesFrag
	case "rectangle_es":
		return shaderRectangleesVert, shaderRectangleesFrag
	case "round_rectangle_es":
		return shaderRectangleesVert, shaderRoundrectangleesFrag
	case "polygon_es":
		return shaderRectangleesVert, shaderPolygonesFrag
	case "arc_es":
		return shaderRectangleesVert, shaderArcesFrag
	case "bezier_curve_es":
		return shaderRectangleesVert, shaderBezierCurveesFrag
	case "arbitrary_polygon_es":
		return shaderRectangleesVert, shaderArbitraryPolygonesFrag
	}
	return nil, nil
}

// rectangleVertexSource returns the standard vertex shader used to fill a vector
// shape's bounding box. User shaders reuse it, just like the built in shapes.
func rectangleVertexSource() []byte {
	return shaderRectangleesVert
}

// userShaderFragment returns the fragment shader source to use for the given
// shader object on this build target.
func userShaderFragment(s *canvas.Shader) []byte {
	return s.SourceES
}
