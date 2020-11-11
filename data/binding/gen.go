// +build ignore

package main

import (
	"os"
	"path"
	"runtime"
	"text/template"

	"fyne.io/fyne"
)

const itemBindTemplate = `
// {{ .Name }} supports binding a {{ .Type }} value in a Fyne application
type {{ .Name }} interface {
	DataItem
	Get() {{ .Type }}
	Set({{ .Type }})
}

// New{{ .Name }} returns a bindable {{ .Type }} value that is managed internally.
func New{{ .Name }}() {{ .Name }} {
	blank := {{ .Default }}
	return &bind{{ .Name }}{val: &blank}
}

// Bind{{ .Name }} returns a new bindable value that controls the contents of the provided {{ .Type }} variable.
func Bind{{ .Name }}(v *{{ .Type }}) {{ .Name }} {
	if v == nil {
		return New{{ .Name }}() // never allow a nil value pointer
	}

	return &bind{{ .Name }}{val: v}
}

type bind{{ .Name }} struct {
	base

	val *{{ .Type }}
}

func (b *bind{{ .Name }}) Get() {{ .Type }} {
	if b.val == nil {
		return {{ .Default }}
	}
	return *b.val
}

func (b *bind{{ .Name }}) Set(val {{ .Type }}) {
	if *b.val == val {
		return
	}
	if b.val == nil { // was not initialized with a blank value, recover
		b.val = &val
	} else {
		*b.val = val
	}

	b.trigger(b)
}
`

const toStringTemplate = `
type stringFrom{{ .Name }} struct {
	base

	from {{ .Name }}
}

// {{ .Name }}ToString creates a binding that connects a {{ .Name }} data item to a String.
// Changes to the {{ .Name }} will be pushed to the String and setting the string will parse and set the
// {{ .Name }} if the parse was successful.
func {{ .Name }}ToString(v {{ .Name }}) String {
	str := &stringFrom{{ .Name }}{from: v}
	v.AddListener(str)
	return str
}

func (s *stringFrom{{ .Name }}) Get() string {
	val := s.from.Get()

	return fmt.Sprintf("{{ .Format }}", val)
}

func (s *stringFrom{{ .Name }}) Set(str string) {
	var val {{ .Type }}
	n, err := fmt.Sscanf(str, "{{ .Format }}", &val)
	if err != nil || n != 1 {
		fyne.LogError("{{ .Type }} parse error", err)
		return
	}
	if val == s.from.Get() {
		return
	}
	s.from.Set(val)

	s.trigger(s)
}

func (s *stringFrom{{ .Name }}) DataChanged(_ DataItem) {
	s.trigger(s)
}
`

type bindValues struct {
	Name, Type, Default string
	Format              string
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

func writeFile(f *os.File, t *template.Template, d interface{}) {
	if err := t.Execute(f, d); err != nil {
		fyne.LogError("Unable to write file "+f.Name(), err)
	}
}

func main() {
	itemFile, err := newFile("binditems")
	if err != nil {
		return
	}
	defer itemFile.Close()
	toStringFile, err := newFile("tostring")
	if err != nil {
		return
	}
	defer itemFile.Close()
	toStringFile.WriteString(`
import (
	"fmt"

	"fyne.io/fyne"
)
`)

	item := template.Must(template.New("item").Parse(itemBindTemplate))
	toString := template.Must(template.New("toString").Parse(toStringTemplate))
	for _, b := range []bindValues{
		bindValues{Name: "Bool", Type: "bool", Default: "false", Format: "%t"},
		bindValues{Name: "Float", Type: "float64", Default: "0.0", Format: "%f"},
		bindValues{Name: "Int", Type: "int", Default: "0", Format: "%d"},
		bindValues{Name: "Rune", Type: "rune", Default: "rune(0)"},
		bindValues{Name: "String", Type: "string", Default: "\"\""},
	} {
		writeFile(itemFile, item, b)
		if b.Type == "string" || b.Type == "rune" {
			continue
		}
		writeFile(toStringFile, toString, b)
	}
}
