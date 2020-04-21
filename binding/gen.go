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
	GetRef() *{{ .Type }}
	Set({{ .Type }})
	SetRef(*{{ .Type }})
	Listen() <-chan {{ .Type }}
}

// base{{ .Name }} implements a data binding for a {{ .Type }}.
type base{{ .Name }} struct {
	sync.Mutex
	reference *{{ .Type }}
	channels  []chan {{ .Type }}
	traces    []string
}

// Empty{{ .Name }} creates a new binding with the empty value.
func Empty{{ .Name }}() {{ .Name }} {
	return New{{ .Name }}({{ .Default }})
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

// Get returns the bound reference.
func (b *base{{ .Name }}) GetRef() *{{ .Type }} {
	return b.reference
}

// Set updates the value of the bound reference.
func (b *base{{ .Name }}) Set(value {{ .Type }}) {
	if *b.reference == value {
		return
	}
	*b.reference = value
	b.Update()
}

// Set updates the bound reference.
func (b *base{{ .Name }}) SetRef(reference *{{ .Type }}) {
	if b.reference == reference {
		return
	}
	b.reference = reference
	b.Update()
}

// Listen returns a channel through which updates will be published.
func (b *base{{ .Name }}) Listen() <-chan {{ .Type }} {
	b.Lock()
	defer b.Unlock()
	c := make(chan {{ .Type }}, 16)
	c <- b.Get()
	b.channels = append(b.channels, c)
	_, file, line, ok := runtime.Caller(1)
	if ok {
		b.traces = append(b.traces, fmt.Sprintf("%s#%d", file, line))
	}
	return c
}

// Update notifies all listeners after a change.
func (b *base{{ .Name }}) Update() {
	b.Lock()
	defer b.Unlock()
	value := b.Get()
	for i, c := range b.channels {
		select {
		case c <- value:
		default:
			fmt.Printf("%d %s Channel Full\n", i, b.traces[i])
		}
	}
}
`

const listBindingTemplate = `
// {{ .Name }}List defines a data binding for a list of {{ .Type }}.
type {{ .Name }}List interface {
	List
	GetBinding(int) {{ .Name }}
	Get{{ .Name }}(int) {{ .Type }}
	GetRef(int) *{{ .Type }}
	SetBinding(int, {{ .Name }})
	Set{{ .Name }}(int, {{ .Type }})
	SetRef(int, *{{ .Type }})
	AddBinding({{ .Name }})
	Add{{ .Name }}({{ .Type }})
	AddRef(*{{ .Type }})
}

// base{{ .Name }}List implements a data binding for a list of {{ .Type }}.
type base{{ .Name }}List struct {
	sync.Mutex
	references *[]*{{ .Type }}
	bindings   map[*{{ .Type }}]{{ .Name }}
	channels   []chan int
	traces     []string
}

// New{{ .Name }}List creates a new list binding with the given values.
func New{{ .Name }}List(values ...{{ .Type }}) {{ .Name }}List {
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
	return b.GetBinding(index)
}

// GetBinding returns the {{ .Name }} at the given index.
func (b *base{{ .Name }}List) GetBinding(index int) {{ .Name }} {
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

// SetBinding updates the {{ .Name }} at the given index.
func (b *base{{ .Name }}List) SetBinding(index int, binding {{ .Name }}) {
	if index < 0 && index >= b.Length() {
		return
	}
	reference := (*b.references)[index]
	if b.bindings[reference] == binding {
		return
	}
	(*b.references)[index] = binding.GetRef()
	b.bindings[reference] = binding
	b.Update()
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

// AddBinding appends the {{ .Name }} to the list.
func (b *base{{ .Name }}List) AddBinding(binding {{ .Name }}) {
	index := b.Length()
	*b.references = append(*b.references, binding.GetRef())
	b.SetBinding(index, binding)
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

// Listen returns a channel through which updates will be published.
func (b *base{{ .Name }}List) Listen() <-chan int {
	b.Lock()
	defer b.Unlock()
	c := make(chan int, 16)
	c <- b.Length()
	b.channels = append(b.channels, c)
	_, file, line, ok := runtime.Caller(1)
	if ok {
		index := len(b.traces)
		trace := fmt.Sprintf("%s#%d", file, line)
		fmt.Printf("%d %s Channel Created\n", index, trace)
		b.traces = append(b.traces, trace)
	}
	return c
}

// Update notifies all listeners after a change.
func (b *base{{ .Name }}List) Update() {
	b.Lock()
	defer b.Unlock()
	length := b.Length()
	for i, c := range b.channels {
		select {
		case c <- length:
		default:
			fmt.Printf("%d %s Channel Full\n", i, b.traces[i])
		}
	}
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
		"fmt",
		"net/url",
		"runtime",
		"sync",
		"",
		"fyne.io/fyne",
	})

	et := template.Must(template.New("element").Parse(elementBindingTemplate))
	lt := template.Must(template.New("list").Parse(listBindingTemplate))

	for _, b := range []*BindingTemplate{
		&BindingTemplate{Name: "Bool", Type: "bool", Default: "false"},
		//&BindingTemplate{Name: "Byte", Type: "byte", Default: "byte(0)"},
		//&BindingTemplate{Name:"Float32",Type:"float32", Default: "float32(0.0)"},
		&BindingTemplate{Name: "Float64", Type: "float64", Default: "float64(0.0)"},
		&BindingTemplate{Name: "Int", Type: "int", Default: "int(0)"},
		//&BindingTemplate{Name:"Int8",Type:"int8", Default: "int8(0)"},
		//&BindingTemplate{Name:"Int16",Type:"int16", Default: "int16(0)"},
		//&BindingTemplate{Name:"Int32",Type:"int32", Default: "int32(0)"},
		&BindingTemplate{Name: "Int64", Type: "int64", Default: "int64(0)"},
		//&BindingTemplate{Name: "Uint", Type: "uint", Default: "uint(0)"},
		//&BindingTemplate{Name:"Uint8",Type:"uint8", Default: "uint8(0)"},
		//&BindingTemplate{Name:"Uint16",Type:"uint16", Default: "uint16(0)"},
		//&BindingTemplate{Name:"Uint32",Type:"uint32", Default: "uint32(0)"},
		//&BindingTemplate{Name: "Uint64", Type: "uint64", Default: "uint64(0)"},
		&BindingTemplate{Name: "Position", Type: "fyne.Position", Default: "fyne.Position{}"},
		&BindingTemplate{Name: "Resource", Type: "fyne.Resource", Default: "nil"},
		&BindingTemplate{Name: "Rune", Type: "rune", Default: "rune(0)"},
		&BindingTemplate{Name: "Size", Type: "fyne.Size", Default: "fyne.Size{}"},
		&BindingTemplate{Name: "String", Type: "string", Default: "\"\""},
		&BindingTemplate{Name: "URL", Type: "*url.URL", Default: "nil"},
	} {
		writeFile(f, et, b)
		writeFile(f, lt, b)
	}

	f.WriteString(`
// Toggle flips the value of the bound reference.
func (b *baseBool) Toggle() {
	*b.reference = !*b.reference
	b.Update()
}
`)
	f.Close()
}
