package font

import (
	"github.com/benoitkugler/textlayout/fonts"
	"github.com/benoitkugler/textlayout/fonts/truetype"
)

type Resource = fonts.Resource

func ParseTTF(file Resource) (Face, error) {
	return truetype.Parse(file)
}
