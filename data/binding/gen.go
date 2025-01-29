//go:build ignore

package main

import (
	"os"
	"path"
	"runtime"
	"text/template"

	"fyne.io/fyne/v2"
)

const toStringTemplate = `
type stringFrom{{ .Name }} struct {
	base
{{ if .Format }}
	format string
{{ end }}
	from {{ .Name }}
}

// {{ .Name }}ToString creates a binding that connects a {{ .Name }} data item to a String.
// Changes to the {{ .Name }} will be pushed to the String and setting the string will parse and set the
// {{ .Name }} if the parse was successful.
//
// Since: {{ .Since }}
func {{ .Name }}ToString(v {{ .Name }}) String {
	str := &stringFrom{{ .Name }}{from: v}
	v.AddListener(str)
	return str
}
{{ if .Format }}
// {{ .Name }}ToStringWithFormat creates a binding that connects a {{ .Name }} data item to a String and is
// presented using the specified format. Changes to the {{ .Name }} will be pushed to the String and setting
// the string will parse and set the {{ .Name }} if the string matches the format and its parse was successful.
//
// Since: {{ .Since }}
func {{ .Name }}ToStringWithFormat(v {{ .Name }}, format string) String {
	if format == "{{ .Format }}" { // Same as not using custom formatting.
		return {{ .Name }}ToString(v)
	}

	str := &stringFrom{{ .Name }}{from: v, format: format}
	v.AddListener(str)
	return str
}
{{ end }}
func (s *stringFrom{{ .Name }}) Get() (string, error) {
	val, err := s.from.Get()
	if err != nil {
		return "", err
	}
{{ if .ToString }}
	return {{ .ToString }}(val)
{{- else }}
	if s.format != "" {
		return fmt.Sprintf(s.format, val), nil
	}

	return format{{ .Name }}(val), nil
{{- end }}
}

func (s *stringFrom{{ .Name }}) Set(str string) error {
{{- if .FromString }}
	val, err := {{ .FromString }}(str)
	if err != nil {
		return err
	}
{{ else }}
	var val {{ .Type }}
	if s.format != "" {
		safe := stripFormatPrecision(s.format)
		n, err := fmt.Sscanf(str, safe+" ", &val) // " " denotes match to end of string
		if err != nil {
			return err
		}
		if n != 1 {
			return errParseFailed
		}
	} else {
		new, err := parse{{ .Name }}(str)
		if err != nil {
			return err
		}
		val = new
	}
{{ end }}
	old, err := s.from.Get()
	if err != nil {
		return err
	}
	if val == old {
		return nil
	}
	if err = s.from.Set(val); err != nil {
		return err
	}

	queueItem(s.DataChanged)
	return nil
}

func (s *stringFrom{{ .Name }}) DataChanged() {
	s.trigger()
}
`
const toIntTemplate = `
type intFrom{{ .Name }} struct {
	base
	from {{ .Name }}
}

// {{ .Name }}ToInt creates a binding that connects a {{ .Name }} data item to an Int.
//
// Since: 2.5
func {{ .Name }}ToInt(v {{ .Name }}) Int {
	i := &intFrom{{ .Name }}{from: v}
	v.AddListener(i)
	return i
}

func (s *intFrom{{ .Name }}) Get() (int, error) {
	val, err := s.from.Get()
	if err != nil {
		return 0, err
	}
	return {{ .ToInt }}(val)
}

func (s *intFrom{{ .Name }}) Set(v int) error {
	val, err := {{ .FromInt }}(v)
	if err != nil {
		return err
	}

	old, err := s.from.Get()
	if err != nil {
		return err
	}
	if val == old {
		return nil
	}
	if err = s.from.Set(val); err != nil {
		return err
	}

	queueItem(s.DataChanged)
	return nil
}

func (s *intFrom{{ .Name }}) DataChanged() {
	s.trigger()
}
`
const fromIntTemplate = `
type intTo{{ .Name }} struct {
	base
	from Int
}

// IntTo{{ .Name }} creates a binding that connects an Int data item to a {{ .Name }}.
//
// Since: 2.5
func IntTo{{ .Name }}(val Int) {{ .Name }} {
	v := &intTo{{ .Name }}{from: val}
	val.AddListener(v)
	return v
}

func (s *intTo{{ .Name }}) Get() ({{ .Type }}, error) {
	val, err := s.from.Get()
	if err != nil {
		return {{ .Default }}, err
	}
	return {{ .FromInt }}(val)
}

func (s *intTo{{ .Name }}) Set(val {{ .Type }}) error {
	i, err := {{ .ToInt }}(val)
	if err != nil {
		return err
	}
	old, err := s.from.Get()
	if i == old {
		return nil
	}
	if err != nil {
		return err
	}
	if err = s.from.Set(i); err != nil {
		return err
	}

	queueItem(s.DataChanged)
	return nil
}

func (s *intTo{{ .Name }}) DataChanged() {
	s.trigger()
}
`
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
	SupportsPreferences  bool
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
	fromInt := template.Must(template.New("fromInt").Parse(fromIntTemplate))
	toInt := template.Must(template.New("toInt").Parse(toIntTemplate))
	toString := template.Must(template.New("toString").Parse(toStringTemplate))
	binds := []bindValues{
		{Name: "Bool", Type: "bool", Default: "false", Format: "%t", SupportsPreferences: true},
		{Name: "Bytes", Type: "[]byte", Default: "nil", Since: "2.2", Comparator: "bytes.Equal"},
		{Name: "Float", Type: "float64", Default: "0.0", Format: "%f", SupportsPreferences: true, ToInt: "internalFloatToInt", FromInt: "internalIntToFloat"},
		{Name: "Int", Type: "int", Default: "0", Format: "%d", SupportsPreferences: true},
		{Name: "Rune", Type: "rune", Default: "rune(0)"},
		{Name: "String", Type: "string", Default: "\"\"", SupportsPreferences: true},
		{Name: "URI", Type: "fyne.URI", Default: "fyne.URI(nil)", Since: "2.1",
			FromString: "uriFromString", ToString: "uriToString", Comparator: "compareURI"},
	}
	for _, b := range binds {
		if b.Since == "" {
			b.Since = "2.0"
		}

		if b.Format != "" || b.ToString != "" {
			writeFile(convertFile, toString, b)
		}
		if b.FromInt != "" {
			writeFile(convertFile, fromInt, b)
		}
		if b.ToInt != "" {
			writeFile(convertFile, toInt, b)
		}
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
