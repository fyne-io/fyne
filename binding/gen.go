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
{{ range . }}{{ if ne . "" }}	"{{ . }}"{{ end }}
{{ end }})
`

const elementBindingTemplate = `
// {{ .Name }} defines a data binding for a {{ .Type }}.
type {{ .Name }} interface {
	Binding
	Get() {{ .Type }}
	Set({{ .Type }})
	Add{{ .Name }}Listener(func({{ .Type }})) *NotifyFunction
}

// base{{ .Name }} implements a data binding for a {{ .Type }}.
type base{{ .Name }} struct {
	Base
	value *{{ .Type }}
}

// New{{ .Name }} creates a new binding with the given value.
func New{{ .Name }}(value *{{ .Type }}) {{ .Name }} {
	return &base{{ .Name }}{value: value}
}

// Get returns the bound value.
func (b *base{{ .Name }}) Get() {{ .Type }} {
	return *b.value
}

// Set updates the bound value.
func (b *base{{ .Name }}) Set(value {{ .Type }}) {
	if *b.value != value {
		*b.value = value
		b.Update()
	}
}

// Add{{ .Name }}Listener adds the given function as a listener to the binding.
// The function is wrapped in the returned NotifyFunction which can be passed to DeleteListener.
func (b *base{{ .Name }}) Add{{ .Name }}Listener(listener func({{ .Type }})) *NotifyFunction {
	n := NewNotifyFunction(func(Binding) {
		listener(*b.value)
	})
	b.AddListener(n)
	return n
}
`

const listBindingTemplate = `
// New{{ .Name }}List creates a new list binding with the given values.
func New{{ .Name }}List(values ...{{ .Type }}) *BaseList {
	list := &BaseList{}
	for _, v := range values {
		w := v
		list.Add(New{{ .Name }}(&w))
	}
	return list
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
		"net/url",
		"",
		"fyne.io/fyne",
	})

	et := template.Must(template.New("element").Parse(elementBindingTemplate))
	lt := template.Must(template.New("list").Parse(listBindingTemplate))

	for _, b := range []*BindingTemplate{
		&BindingTemplate{Name: "Bool", Type: "bool"},
		//&BindingTemplate{Name: "Byte", Type: "byte"},
		//&BindingTemplate{Name:"Float32",Type:"float32"},
		&BindingTemplate{Name: "Float64", Type: "float64"},
		&BindingTemplate{Name: "Int", Type: "int"},
		//&BindingTemplate{Name:"Int8",Type:"int8"},
		//&BindingTemplate{Name:"Int16",Type:"int16"},
		//&BindingTemplate{Name:"Int32",Type:"int32"},
		&BindingTemplate{Name: "Int64", Type: "int64"},
		//&BindingTemplate{Name: "Uint", Type: "uint"},
		//&BindingTemplate{Name:"Uint8",Type:"uint8"},
		//&BindingTemplate{Name:"Uint16",Type:"uint16"},
		//&BindingTemplate{Name:"Uint32",Type:"uint32"},
		//&BindingTemplate{Name: "Uint64", Type: "uint64"},
		&BindingTemplate{Name: "Resource", Type: "fyne.Resource"},
		&BindingTemplate{Name: "Rune", Type: "rune"},
		&BindingTemplate{Name: "String", Type: "string"},
		&BindingTemplate{Name: "URL", Type: "*url.URL"},
	} {
		writeFile(f, et, b)
		writeFile(f, lt, b)
	}

	f.Close()
}
