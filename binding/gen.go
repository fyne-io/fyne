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

type bindValues struct {
	Name, Type, Default string
}

func writeFile(f *os.File, t *template.Template, d interface{}) {
	if err := t.Execute(f, d); err != nil {
		fyne.LogError("Unable to write file "+f.Name(), err)
	}
}

func main() {
	_, dirname, _, _ := runtime.Caller(0)
	filepath := path.Join(path.Dir(dirname), "binditems.go")
	os.Remove(filepath)
	f, err := os.Create(filepath)
	if err != nil {
		fyne.LogError("Unable to open file "+f.Name(), err)
		return
	}
	defer f.Close()

	f.WriteString(`// auto-generated
// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //

package binding
`)

	item := template.Must(template.New("item").Parse(itemBindTemplate))
	for _, b := range []bindValues{
		bindValues{Name: "Bool", Type: "bool", Default: "false"},
		bindValues{Name: "Float", Type: "float64", Default: "0.0"},
		bindValues{Name: "Int", Type: "int", Default: "0"},
		bindValues{Name: "Rune", Type: "rune", Default: "rune(0)"},
		bindValues{Name: "String", Type: "string", Default: "\"\""},
	} {
		writeFile(f, item, b)
	}
}
