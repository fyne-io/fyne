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
	reference *{{ .Type }}
}

// New{{ .Name }} creates a new binding with the given value.
func New{{ .Name }}(value {{ .Type }}) {{ .Name }} {
	return New{{ .Name }}Ref(&value)
}

// New{{ .Name }}Ref creates a new binding with the given reference.
func New{{ .Name }}Ref(reference *{{ .Type }}) {{ .Name }} {
	return &base{{ .Name }}{reference: reference}
}

// Get returns the value of the bound reference.
func (b *base{{ .Name }}) Get() {{ .Type }} {
	return *b.reference
}

// Set updates the value of the bound reference.
func (b *base{{ .Name }}) Set(value {{ .Type }}) {
	if *b.reference == value {
		return
	}
	*b.reference = value
	b.Update()
}

// Add{{ .Name }}Listener adds the given function as a listener to the binding.
// The function is wrapped in the returned NotifyFunction which can be passed to DeleteListener.
func (b *base{{ .Name }}) Add{{ .Name }}Listener(listener func({{ .Type }})) *NotifyFunction {
	n := NewNotifyFunction(func(Binding) {
		listener(*b.reference)
	})
	b.AddListener(n)
	return n
}
`

const listBindingTemplate = `
// {{ .Name }}List defines a data binding for a list of {{ .Type }}.
type {{ .Name }}List interface {
	List
	Get{{ .Name }}(int) {{ .Type }}
	GetRef(int) *{{ .Type }}
	Set{{ .Name }}(int, {{ .Type }})
	SetRef(int, *{{ .Type }})
	Add{{ .Name }}({{ .Type }})
	AddRef(*{{ .Type }})
}

// base{{ .Name }}List implements a data binding for a list of {{ .Type }}.
type base{{ .Name }}List struct {
	Base
	sync.Mutex
	references *[]*{{ .Type }}
	bindings   map[*{{ .Type }}]{{ .Name }}
}

// New{{ .Name }}List creates a new list binding with the given values.
func New{{ .Name }}List(values []{{ .Type }}) {{ .Name }}List {
	var references []*{{ .Type }}
	for i := 0; i < len(values); i++ {
		references = append(references, &values[i])
	}
	return New{{ .Name }}ListRefs(&references)
}

// New{{ .Name }}ListRefs creates a new list binding with the given references.
func New{{ .Name }}ListRefs(references *[]*{{ .Type }}) {{ .Name }}List {
	return &base{{ .Name }}List{
		references: references,
		bindings:   make(map[*{{ .Type }}]{{ .Name }}),
	}
}

// Length returns the number of elements in the list.
func (b *base{{ .Name }}List) Length() int {
	return len(*b.references)
}

// Get returns the binding at the given index.
func (b *base{{ .Name }}List) Get(index int) Binding {
	reference := b.GetRef(index)
	if reference == nil {
		return nil
	}
	b.Lock()
	defer b.Unlock()
	binding, ok := b.bindings[reference]
	if !ok {
		binding = New{{ .Name }}Ref(reference)
		b.bindings[reference] = binding
	}
	return binding
}

// Get{{ .Name }} returns the {{ .Type }} at the given index.
func (b *base{{ .Name }}List) Get{{ .Name }}(index int) {{ .Type }} {
	if index < 0 && index >= b.Length() {
		return {{ .Default }}
	}
	return *(*b.references)[index]
}

// GetRef returns the reference at the given index.
func (b *base{{ .Name }}List) GetRef(index int) *{{ .Type }} {
	if index < 0 && index >= b.Length() {
		return nil
	}
	return (*b.references)[index]
}

// Set{{ .Name }} updates the {{ .Type }} at the given index.
func (b *base{{ .Name }}List) Set{{ .Name }}(index int, value {{ .Type }}) {
	if index < 0 && index >= b.Length() {
		return
	}
	if *(*b.references)[index] == value {
		return
	}
	b.Get(index).({{ .Name }}).Set(value)
}

// SetRef updates the {{ .Type }} at the given index.
func (b *base{{ .Name }}List) SetRef(index int, reference *{{ .Type }}) {
	if index < 0 && index >= b.Length() {
		return
	}
	if (*b.references)[index] == reference {
		return
	}
	(*b.references)[index] = reference
	b.Update()
}

// Add{{ .Name }} appends the {{ .Type }} to the list.
func (b *base{{ .Name }}List) Add{{ .Name }}(value {{ .Type }}) {
	b.AddRef(&value)
}

// AddRef appends the {{ .Type }} to the list.
func (b *base{{ .Name }}List) AddRef(reference *{{ .Type }}) {
	*b.references = append(*b.references, reference)
	b.Update()
}
`

type BindingTemplate struct {
	Name, Type, Default string
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
		"sync",
		"",
		"fyne.io/fyne",
	})

	et := template.Must(template.New("element").Parse(elementBindingTemplate))
	lt := template.Must(template.New("list").Parse(listBindingTemplate))

	for _, b := range []*BindingTemplate{
		&BindingTemplate{Name: "Bool", Type: "bool", Default: "false"},
		//&BindingTemplate{Name: "Byte", Type: "byte", Default: "0"},
		//&BindingTemplate{Name:"Float32",Type:"float32", Default: "0.0"},
		&BindingTemplate{Name: "Float64", Type: "float64", Default: "0.0"},
		&BindingTemplate{Name: "Int", Type: "int", Default: "0"},
		//&BindingTemplate{Name:"Int8",Type:"int8", Default: "0"},
		//&BindingTemplate{Name:"Int16",Type:"int16", Default: "0"},
		//&BindingTemplate{Name:"Int32",Type:"int32", Default: "0"},
		&BindingTemplate{Name: "Int64", Type: "int64", Default: "0"},
		//&BindingTemplate{Name: "Uint", Type: "uint", Default: "0"},
		//&BindingTemplate{Name:"Uint8",Type:"uint8", Default: "0"},
		//&BindingTemplate{Name:"Uint16",Type:"uint16", Default: "0"},
		//&BindingTemplate{Name:"Uint32",Type:"uint32", Default: "0"},
		//&BindingTemplate{Name: "Uint64", Type: "uint64", Default: "0"},
		&BindingTemplate{Name: "Resource", Type: "fyne.Resource", Default: "nil"},
		&BindingTemplate{Name: "Rune", Type: "rune", Default: "0"},
		&BindingTemplate{Name: "String", Type: "string", Default: "\"\""},
		&BindingTemplate{Name: "URL", Type: "*url.URL", Default: "nil"},
	} {
		writeFile(f, et, b)
		writeFile(f, lt, b)
	}

	f.Close()
}
