package theme

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"image/color"
	"io"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

// FromJSON returns a Theme created from the given JSON metadata.
// Any values not present in the data will fall back to the default theme.
// If a parse error occurs it will be returned along with a default theme.
//
// Since: 2.2
func FromJSON(data string) (fyne.Theme, error) {
	return FromJSONReader(strings.NewReader(data))
}

// FromJSONReader returns a Theme created from the given JSON metadata through the reader.
// Any values not present in the data will fall back to the default theme.
// If a parse error occurs it will be returned along with a default theme.
//
// Since: 2.2
func FromJSONReader(r io.Reader) (fyne.Theme, error) {
	var th *schema
	if err := json.NewDecoder(r).Decode(&th); err != nil {
		return DefaultTheme(), err
	}

	return &jsonTheme{data: th, fallback: DefaultTheme()}, nil
}

type hexColor string

func (h hexColor) color() (color.Color, error) {
	data := h
	switch len([]rune(h)) {
	case 8, 6:
	case 9, 7: // remove # prefix
		data = h[1:]
	case 5: // remove # prefix, then double up
		data = h[1:]
		fallthrough
	case 4: // could be rgba or #rgb
		if data[0] == '#' {
			v := []rune(data[1:])
			data = hexColor([]rune{v[0], v[0], v[1], v[1], v[2], v[2]})
			break
		}

		v := []rune(data)
		data = hexColor([]rune{v[0], v[0], v[1], v[1], v[2], v[2], v[3], v[3]})
	case 3:
		v := []rune(h)
		data = hexColor([]rune{v[0], v[0], v[1], v[1], v[2], v[2]})
	default:
		return color.Transparent, errors.New("invalid color format: " + string(h))
	}

	digits, err := hex.DecodeString(string(data))
	if err != nil {
		return nil, err
	}
	ret := &color.NRGBA{R: digits[0], G: digits[1], B: digits[2]}
	if len(digits) == 4 {
		ret.A = digits[3]
	} else {
		ret.A = 0xff
	}

	return ret, nil
}

type uriString string

func (u uriString) resource() fyne.Resource {
	uri, err := storage.ParseURI(string(u))
	if err != nil {
		fyne.LogError("Failed to parse URI", err)
		return nil
	}
	r, err := storage.LoadResourceFromURI(uri)
	if err != nil {
		fyne.LogError("Failed to load resource from URI", err)
		return nil
	}
	return r
}

type schema struct {
	Colors      map[string]hexColor `json:"Colors,omitempty"`
	DarkColors  map[string]hexColor `json:"Colors-dark,omitempty"`
	LightColors map[string]hexColor `json:"Colors-light,omitempty"`
	Sizes       map[string]float32  `json:"Sizes,omitempty"`

	Fonts map[string]uriString `json:"Fonts,omitempty"`
	Icons map[string]uriString `json:"Icons,omitempty"`
}

type jsonTheme struct {
	data     *schema
	fallback fyne.Theme
}

func (t *jsonTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch variant {
	case VariantLight:
		if val, ok := t.data.LightColors[string(name)]; ok {
			c, err := val.color()
			if err != nil {
				fyne.LogError("Failed to parse color", err)
			} else {
				return c
			}
		}
	case VariantDark:
		if val, ok := t.data.DarkColors[string(name)]; ok {
			c, err := val.color()
			if err != nil {
				fyne.LogError("Failed to parse color", err)
			} else {
				return c
			}
		}
	}

	if val, ok := t.data.Colors[string(name)]; ok {
		c, err := val.color()
		if err != nil {
			fyne.LogError("Failed to parse color", err)
		} else {
			return c
		}
	}

	return t.fallback.Color(name, variant)
}

func (t *jsonTheme) Font(style fyne.TextStyle) fyne.Resource {
	if val, ok := t.data.Fonts[styleString(style)]; ok {
		r := val.resource()
		if r != nil {
			return r
		}
	}
	return t.fallback.Font(style)
}

func (t *jsonTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	if val, ok := t.data.Icons[string(name)]; ok {
		r := val.resource()
		if r != nil {
			return r
		}
	}
	return t.fallback.Icon(name)
}

func (t *jsonTheme) Size(name fyne.ThemeSizeName) float32 {
	if val, ok := t.data.Sizes[string(name)]; ok {
		return val
	}

	return t.fallback.Size(name)
}

func styleString(s fyne.TextStyle) string {
	if s.Bold {
		if s.Italic {
			return "boldItalic"
		}
		return "bold"
	}
	if s.Monospace {
		return "monospace"
	}
	return "regular"
}
