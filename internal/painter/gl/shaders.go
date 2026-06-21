//go:build (!gles && !arm && !arm64 && !android && !ios && !mobile && !test_web_driver && !wasm) || (darwin && !mobile && !ios && !wasm && !test_web_driver)

package gl

import (
	_ "embed"

	"fyne.io/fyne/v2/canvas"
)

var (
	//go:embed shaders/blur.frag
	shaderBlurFrag []byte

	//go:embed shaders/blur.vert
	shaderBlurVert []byte

	//go:embed shaders/line.frag
	shaderLineFrag []byte

	//go:embed shaders/line.vert
	shaderLineVert []byte

	//go:embed shaders/rectangle.frag
	shaderRectangleFrag []byte

	//go:embed shaders/rectangle.vert
	shaderRectangleVert []byte

	//go:embed shaders/round_rectangle.frag
	shaderRoundrectangleFrag []byte

	//go:embed shaders/simple.frag
	shaderSimpleFrag []byte

	//go:embed shaders/simple.vert
	shaderSimpleVert []byte

	//go:embed shaders/regular_polygon.frag
	shaderPolygonFrag []byte

	//go:embed shaders/arc.frag
	shaderArcFrag []byte

	//go:embed shaders/bezier_curve.frag
	shaderBezierCurveFrag []byte

	//go:embed shaders/arbitrary_polygon.frag
	shaderArbitraryPolygonFrag []byte

	//go:embed shaders/ellipse.frag
	shaderEllipseFrag []byte
)

func shaderSourceNamed(name string) ([]byte, []byte) {
	switch name {
	case "line":
		return shaderLineVert, shaderLineFrag
	case "simple":
		return shaderSimpleVert, shaderSimpleFrag
	case "rectangle":
		return shaderRectangleVert, shaderRectangleFrag
	case "round_rectangle":
		return shaderRectangleVert, shaderRoundrectangleFrag
	case "blur":
		return shaderBlurVert, shaderBlurFrag
	case "polygon":
		return shaderRectangleVert, shaderPolygonFrag
	case "arc":
		return shaderRectangleVert, shaderArcFrag
	case "bezier_curve":
		return shaderRectangleVert, shaderBezierCurveFrag
	case "arbitrary_polygon":
		return shaderRectangleVert, shaderArbitraryPolygonFrag
	case "ellipse":
		return shaderRectangleVert, shaderEllipseFrag
	}
	return nil, nil
}

// rectangleVertexSource returns the standard vertex shader used to fill a vector
// shape's bounding box. User shaders reuse it, just like the built in shapes.
func rectangleVertexSource() []byte {
	return shaderRectangleVert
}

// userShaderFragment returns the fragment shader source to use for the given
// shader object on this build target.
func userShaderFragment(s *canvas.Shader) []byte {
	return s.Source
}
