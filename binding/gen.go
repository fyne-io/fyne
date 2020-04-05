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

func typeError(e string, a interface{}) {
	fyne.LogError(fmt.Sprintf("Incorrect type: expected '%s', got '%v'", e, a), nil)
}
`

const bindingTemplate = `
type {{ .Name }}Binding struct {
	itemBinding
	value {{ .Type }}
}

func New{{ .Name }}Binding(value {{ .Type }}) *{{ .Name }}Binding {
	return &{{ .Name }}Binding{value: value}
}

func (b *{{ .Name }}Binding) Get() {{ .Type }} {
	return b.value
}

func (b *{{ .Name }}Binding) Set(value {{ .Type }}) {
	if b.value != value {
		b.value = value
		b.notify(value)
	}
}

func (b *{{ .Name }}Binding) AddListener(listener func({{ .Type }})) {
	b.addListener(func(value interface{}) {
		v, ok := value.({{ .Type }})
		if ok {
			listener(v)
		} else {
			typeError("{{ .Type }}", value)
		}
	})
}
`

type BindingTemplate struct {
	Name, Type string
}

func writeFile(f *os.File, t *template.Template, d interface{}) {
	if err := t.Execute(f, d); err != nil {
		fyne.LogError("Unable to write file "+f.Name(), err)
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
	writeFile(f, ht, []string{
		"fmt",
		"fyne.io/fyne",
		"net/url",
	})

	t := template.Must(template.New("binding").Parse(bindingTemplate))

	for _, b := range []*BindingTemplate{
		&BindingTemplate{Name: "Bool", Type: "bool"},
		&BindingTemplate{Name: "Byte", Type: "byte"},
		//		&BindingTemplate{Name:"Float32",Type:"float32"},
		&BindingTemplate{Name: "Float64", Type: "float64"},
		&BindingTemplate{Name: "Int", Type: "int"},
		//		&BindingTemplate{Name:"Int8",Type:"int8"},
		//		&BindingTemplate{Name:"Int16",Type:"int16"},
		//		&BindingTemplate{Name:"Int32",Type:"int32"},
		&BindingTemplate{Name: "Int64", Type: "int64"},
		&BindingTemplate{Name: "Uint", Type: "uint"},
		//		&BindingTemplate{Name:"Uint8",Type:"uint8"},
		//		&BindingTemplate{Name:"Uint16",Type:"uint16"},
		//		&BindingTemplate{Name:"Uint32",Type:"uint32"},
		&BindingTemplate{Name: "Uint64", Type: "uint64"},
		&BindingTemplate{Name: "Resource", Type: "fyne.Resource"},
		&BindingTemplate{Name: "Rune", Type: "rune"},
		&BindingTemplate{Name: "String", Type: "string"},
		&BindingTemplate{Name: "URL", Type: "*url.URL"},
	} {
		writeFile(f, t, b)
	}

	f.Close()
}
