package gl

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

var (
	//go:embed shaders/line.frag
	lineFrag []byte

	//go:embed shaders/line.vert
	lineVert []byte

	//go:embed shaders/line_es.frag
	lineesFrag []byte

	//go:embed shaders/line_es.vert
	lineesVert []byte

	//go:embed shaders/rectangle.frag
	rectangleFrag []byte

	//go:embed shaders/rectangle.vert
	rectangleVert []byte

	//go:embed shaders/rectangle_es.frag
	rectangleesFrag []byte

	//go:embed shaders/rectangle_es.vert
	rectangleesVert []byte

	//go:embed shaders/round_rectangle.frag
	roundrectangleFrag []byte

	//go:embed shaders/round_rectangle_es.frag
	roundrectangleesFrag []byte

	//go:embed shaders/simple.frag
	simpleFrag []byte

	//go:embed shaders/simple.vert
	simpleVert []byte

	//go:embed shaders/simple_es.frag
	simpleesFrag []byte

	//go:embed shaders/simple_es.vert
	simpleesVert []byte
)

var (
	shaderLineFrag = &fyne.StaticResource{
		StaticName:    "line.frag",
		StaticContent: lineFrag,
	}

	shaderLineVert = &fyne.StaticResource{
		StaticName:    "line.vert",
		StaticContent: lineVert,
	}

	shaderLineesFrag = &fyne.StaticResource{
		StaticName:    "line_es.frag",
		StaticContent: lineesFrag,
	}

	shaderLineesVert = &fyne.StaticResource{
		StaticName:    "line_es.vert",
		StaticContent: lineesVert,
	}

	shaderRectangleFrag = &fyne.StaticResource{
		StaticName:    "rectangle.frag",
		StaticContent: rectangleFrag,
	}

	shaderRectangleVert = &fyne.StaticResource{
		StaticName:    "rectangle.vert",
		StaticContent: rectangleVert,
	}

	shaderRectangleesFrag = &fyne.StaticResource{
		StaticName:    "rectangle_es.frag",
		StaticContent: rectangleesFrag,
	}

	shaderRectangleesVert = &fyne.StaticResource{
		StaticName:    "rectangle_es.vert",
		StaticContent: rectangleesVert,
	}

	shaderRoundrectangleFrag = &fyne.StaticResource{
		StaticName:    "round_rectangle.frag",
		StaticContent: roundrectangleFrag,
	}

	shaderRoundrectangleesFrag = &fyne.StaticResource{
		StaticName:    "round_rectangle_es.frag",
		StaticContent: roundrectangleesFrag,
	}

	shaderSimpleFrag = &fyne.StaticResource{
		StaticName:    "simple.frag",
		StaticContent: simpleFrag,
	}

	shaderSimpleVert = &fyne.StaticResource{
		StaticName:    "simple.vert",
		StaticContent: simpleVert,
	}

	shaderSimpleesFrag = &fyne.StaticResource{
		StaticName:    "simple_es.frag",
		StaticContent: simpleesFrag,
	}

	shaderSimpleesVert = &fyne.StaticResource{
		StaticName:    "simple_es.vert",
		StaticContent: simpleesVert,
	}
)
