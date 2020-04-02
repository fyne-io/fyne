// +build ignore

package main

import (
	"os"
	"path"
	"runtime"
	"text/template"

	"fyne.io/fyne"
)

const headerTemplate = `// auto-generated
// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //

package binding

import (
{{ range . }}	"{{ . }}"
{{ end }})
`
const unitBindingTemplate = `
type {{ .Name }}Binding struct {
	BaseBinding
	Value {{ .Type }}
}

func (b *{{ .Name }}Binding) Get{{ .Name }}() {{ .Type }} {
	return b.Value
}

func (b *{{ .Name }}Binding) Set(value interface{}) {
	v, ok := value.({{ .Type }})
	if ok {
		b.Set{{ .Name }}(v)
	} else {
		fyne.LogError(fmt.Sprintf("Incorrect type: expected '{{ .Type }}', got '%v'", value), nil)
	}
}

func (b *{{ .Name }}Binding) Set{{ .Name }}(value {{ .Type }}) {
	if b.Value != value {
		b.Value = value
		b.notify(value)
	}
}

func (b *{{ .Name }}Binding) Add{{ .Name }}Listener(listener func({{ .Type }})) {
	b.addListener(func(value interface{}) {
		v, ok := value.({{ .Type }})
		if ok {
			listener(v)
		} else {
			fyne.LogError(fmt.Sprintf("Incorrect type: expected '{{ .Type }}', got '%v'", value), nil)
		}
	})
}
`

func writeHeader(f *os.File, t *template.Template, imports []string) {
	if err := t.Execute(f, imports); err != nil {
		fyne.LogError("Unable to write file "+f.Name(), err)
	}
}

func writeUnitBinding(f *os.File, t *template.Template, name, typ string) {
	data := struct {
		Name, Type string
	}{
		Name: name,
		Type: typ,
	}
	if err := t.Execute(f, data); err != nil {
		fyne.LogError("Unable to write file "+f.Name(), err)
		return
	}
}

func main() {
	_, dirname, _, _ := runtime.Caller(0)
	filepath := path.Join(path.Dir(dirname), "bindings.go")
	os.Remove(filepath)
	f, err := os.Create(filepath)
	if err != nil {
		fyne.LogError("Unable to open file "+f.Name(), err)
		return
	}

	ht := template.Must(template.New("header").Parse(headerTemplate))
	writeHeader(f, ht, []string{
		"fmt",
		"fyne.io/fyne",
		"net/url",
	})

	ubt := template.Must(template.New("binding").Parse(unitBindingTemplate))

	writeUnitBinding(f, ubt, "Bool", "bool")
	writeUnitBinding(f, ubt, "Byte", "byte")

	writeUnitBinding(f, ubt, "Float32", "float32")
	writeUnitBinding(f, ubt, "Float64", "float64")

	writeUnitBinding(f, ubt, "Int", "int")
	writeUnitBinding(f, ubt, "Int8", "int8")
	writeUnitBinding(f, ubt, "Int16", "int16")
	writeUnitBinding(f, ubt, "Int32", "int32")
	writeUnitBinding(f, ubt, "Int64", "int64")

	writeUnitBinding(f, ubt, "Uint", "uint")
	writeUnitBinding(f, ubt, "Uint8", "uint8")
	writeUnitBinding(f, ubt, "Uint16", "uint16")
	writeUnitBinding(f, ubt, "Uint32", "uint32")
	writeUnitBinding(f, ubt, "Uint64", "uint64")

	writeUnitBinding(f, ubt, "Resource", "fyne.Resource")

	writeUnitBinding(f, ubt, "Rune", "rune")

	writeUnitBinding(f, ubt, "String", "string")

	writeUnitBinding(f, ubt, "URL", "*url.URL")

	f.Close()
}
