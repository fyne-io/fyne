//go:build ((gles || arm || arm64) && !android && !ios && !mobile && !darwin && !wasm && !test_web_driver) || ((android || ios || mobile) && (!wasm || !test_web_driver)) || wasm || test_web_driver

package gl

import _ "embed"

var (
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

	//go:embed shaders/polygon_es.frag
	shaderPolygonesFrag []byte

	//go:embed shaders/arc_es.frag
	shaderArcesFrag []byte
)

func shaderSourceNamed(name string) ([]byte, []byte) {
	switch name {
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
	}
	return nil, nil
}
