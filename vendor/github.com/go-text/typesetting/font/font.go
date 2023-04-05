package font

import (
	"github.com/go-text/typesetting/opentype/api/font"
	"github.com/go-text/typesetting/opentype/loader"
)

type Resource = loader.Resource

func ParseTTF(file Resource) (Face, error) {
	ld, err := loader.NewLoader(file)
	if err != nil {
		return nil, err
	}
	ft, err := font.NewFont(ld)
	if err != nil {
		return nil, err
	}
	return &font.Face{Font: ft}, nil
}
