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

// FromJSONWithFallback returns a Theme created from the given JSON metadata.
// Any values not present in the data will fall back to the specified theme.
// If a parse error occurs it will be returned along with a specified fallback theme.
//
// Since: 2.7
func FromJSONWithFallback(data string, fallback fyne.Theme) (fyne.Theme, error) {
	return fromJSONWithFallback(strings.NewReader(data), fallback)
}

// FromJSONReader returns a Theme created from the given JSON metadata through the reader.
// Any values not present in the data will fall back to the default theme.
// If a parse error occurs it will be returned along with a default theme.
//
// Since: 2.2
func FromJSONReader(r io.Reader) (fyne.Theme, error) {
	return fromJSONWithFallback(r, DefaultTheme())
}

// FromJSONReaderWithFallback returns a Theme created from the given JSON metadata through the reader.
// Any values not present in the data will fall back to the specified theme.
// If a parse error occurs it will be returned along with a specified fallback theme.
//
// Since: 2.7
func FromJSONReaderWithFallback(r io.Reader, fallback fyne.Theme) (fyne.Theme, error) {
	return fromJSONWithFallback(r, fallback)
}

func fromJSONWithFallback(r io.Reader, fallback fyne.Theme) (fyne.Theme, error) {
	var th *schema
	if err := json.NewDecoder(r).Decode(&th); err != nil {
		return fallback, err
	}

	return &jsonTheme{data: th, fallback: fallback}, nil
}

type jsonColor struct {
	color color.Color
}

func (h *jsonColor) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	return h.parseColor(str)
}

func (h *jsonColor) parseColor(str string) error {
	data := str
	switch len([]rune(str)) {
	case 8, 6:
	case 9, 7: // remove # prefix
		data = str[1:]
	case 5: // remove # prefix, then double up
		data = str[1:]
		fallthrough
	case 4: // could be rgba or #rgb
		if data[0] == '#' {
			v := []rune(data[1:])
			data = string([]rune{v[0], v[0], v[1], v[1], v[2], v[2]})
			break
		}

		v := []rune(data)
		data = string([]rune{v[0], v[0], v[1], v[1], v[2], v[2], v[3], v[3]})
	case 3:
		v := []rune(str)
		data = string([]rune{v[0], v[0], v[1], v[1], v[2], v[2]})
	default:
		h.color = color.Transparent
		return errors.New("invalid color format: " + str)
	}

	digits, err := hex.DecodeString(data)
	if err != nil {
		return err
	}
	ret := &color.NRGBA{R: digits[0], G: digits[1], B: digits[2]}
	if len(digits) == 4 {
		ret.A = digits[3]
	} else {
		ret.A = 0xff
	}

	h.color = ret
	return nil
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
	Colors      map[string]jsonColor `json:"Colors,omitempty"`
	DarkColors  map[string]jsonColor `json:"Colors-dark,omitempty"`
	LightColors map[string]jsonColor `json:"Colors-light,omitempty"`
	Sizes       map[string]float32   `json:"Sizes,omitempty"`

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
			return val.color
		}
	case VariantDark:
		if val, ok := t.data.DarkColors[string(name)]; ok {
			return val.color
		}
	}

	if val, ok := t.data.Colors[string(name)]; ok {
		return val.color
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
