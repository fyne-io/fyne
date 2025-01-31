//go:build ignore

package main

import (
	"os"
	"path"
	"runtime"
	"text/template"

	"fyne.io/fyne/v2"
)

const fromStringTemplate = `
type stringTo{{ .Name }} struct {
	base
{{ if .Format }}
	format string
{{ end }}
	from String
}

// StringTo{{ .Name }} creates a binding that connects a String data item to a {{ .Name }}.
// Changes to the String will be parsed and pushed to the {{ .Name }} if the parse was successful, and setting
// the {{ .Name }} update the String binding.
//
// Since: {{ .Since }}
func StringTo{{ .Name }}(str String) {{ .Name }} {
	v := &stringTo{{ .Name }}{from: str}
	str.AddListener(v)
	return v
}
{{ if .Format }}
// StringTo{{ .Name }}WithFormat creates a binding that connects a String data item to a {{ .Name }} and is
// presented using the specified format. Changes to the {{ .Name }} will be parsed and if the format matches and
// the parse is successful it will be pushed to the String. Setting the {{ .Name }} will push a formatted value
// into the String.
//
// Since: {{ .Since }}
func StringTo{{ .Name }}WithFormat(str String, format string) {{ .Name }} {
	if format == "{{ .Format }}" { // Same as not using custom format.
		return StringTo{{ .Name }}(str)
	}

	v := &stringTo{{ .Name }}{from: str, format: format}
	str.AddListener(v)
	return v
}
{{ end }}
func (s *stringTo{{ .Name }}) Get() ({{ .Type }}, error) {
	str, err := s.from.Get()
	if str == "" || err != nil {
		return {{ .Default }}, err
	}
{{ if .FromString }}
	return {{ .FromString }}(str)
{{- else }}
	var val {{ .Type }}
	if s.format != "" {
		n, err := fmt.Sscanf(str, s.format+" ", &val) // " " denotes match to end of string
		if err != nil {
			return {{ .Default }}, err
		}
		if n != 1 {
			return {{ .Default }}, errParseFailed
		}
	} else {
		new, err := parse{{ .Name }}(str)
		if err != nil {
			return {{ .Default }}, err
		}
		val = new
	}

	return val, nil
{{- end }}
}

func (s *stringTo{{ .Name }}) Set(val {{ .Type }}) error {
{{- if .ToString }}
	str, err := {{ .ToString }}(val)
	if err != nil {
		return err
	}
{{- else }}
	var str string
	if s.format != "" {
		str = fmt.Sprintf(s.format, val)
	} else {
		str = format{{ .Name }}(val)
	}
{{ end }}
	old, err := s.from.Get()
	if str == old {
		return err
	}

	if err = s.from.Set(str); err != nil {
		return err
	}

	queueItem(s.DataChanged)
	return nil
}

func (s *stringTo{{ .Name }}) DataChanged() {
	s.trigger()
}
`

type bindValues struct {
	Name, Type, Default  string
	Format, Since        string
	FromString, ToString string // function names...
	Comparator           string // comparator function name
	FromInt, ToInt       string // function names...
}

func newFile(name string) (*os.File, error) {
	_, dirname, _, _ := runtime.Caller(0)
	filepath := path.Join(path.Dir(dirname), name+".go")
	os.Remove(filepath)
	f, err := os.Create(filepath)
	if err != nil {
		fyne.LogError("Unable to open file "+f.Name(), err)
		return nil, err
	}

	f.WriteString(`// auto-generated
// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //

package binding
`)
	return f, nil
}

func writeFile(f *os.File, t *template.Template, d any) {
	if err := t.Execute(f, d); err != nil {
		fyne.LogError("Unable to write file "+f.Name(), err)
	}
}

func main() {
	convertFile, err := newFile("convert")
	if err != nil {
		return
	}
	defer convertFile.Close()
	convertFile.WriteString(`
import (
	"fmt"

	"fyne.io/fyne/v2"
)

func internalFloatToInt(val float64) (int, error) {
	return int(val), nil
}

func internalIntToFloat(val int) (float64, error) {
	return float64(val), nil
}
`)

	fromString := template.Must(template.New("fromString").Parse(fromStringTemplate))
	binds := []bindValues{
		{Name: "Bool", Type: "bool", Default: "false", Format: "%t"},
		{Name: "Bytes", Type: "[]byte", Default: "nil", Since: "2.2", Comparator: "bytes.Equal"},
		{Name: "Float", Type: "float64", Default: "0.0", Format: "%f", ToInt: "internalFloatToInt", FromInt: "internalIntToFloat"},
		{Name: "Int", Type: "int", Default: "0", Format: "%d"},
		{Name: "Rune", Type: "rune", Default: "rune(0)"},
		{Name: "String", Type: "string", Default: "\"\""},
		{Name: "URI", Type: "fyne.URI", Default: "fyne.URI(nil)", Since: "2.1",
			FromString: "uriFromString", ToString: "uriToString", Comparator: "compareURI"},
	}

	// add StringTo... at the bottom of the convertFile for correct ordering
	for _, b := range binds {
		if b.Since == "" {
			b.Since = "2.0"
		}

		if b.Format != "" || b.FromString != "" {
			writeFile(convertFile, fromString, b)
		}
	}
}
