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
	return &bound{{ .Name }}{val: &blank}
}

// Bind{{ .Name }} returns a new bindable value that controls the contents of the provided {{ .Type }} variable.
func Bind{{ .Name }}(v *{{ .Type }}) {{ .Name }} {
	if v == nil {
		return New{{ .Name }}() // never allow a nil value pointer
	}

	return &bound{{ .Name }}{val: v}
}

type bound{{ .Name }} struct {
	base

	val *{{ .Type }}
}

func (b *bound{{ .Name }}) Get() {{ .Type }} {
	if b.val == nil {
		return {{ .Default }}
	}
	return *b.val
}

func (b *bound{{ .Name }}) Set(val {{ .Type }}) {
	if *b.val == val {
		return
	}
	if b.val == nil { // was not initialized with a blank value, recover
		b.val = &val
	} else {
		*b.val = val
	}

	b.trigger()
}
`

const prefTemplate = `
type prefBound{{ .Name }} struct {
	base
	key string
	p   fyne.Preferences
}

// BindPreference{{ .Name }} returns a bindable {{ .Type }} value that is managed by the application preferences.
// Changes to this value will be saved to application storage and when the app starts the previous values will be read.
func BindPreference{{ .Name }}(key string, p fyne.Preferences) {{ .Name }} {
	if listen, ok := prefBinds[key]; ok {
		if l, ok := listen.({{ .Name }}); ok {
			return l
		}
		fyne.LogError(keyTypeMismatchError+key, nil)
	}

	listen := &prefBound{{ .Name }}{key: key, p: p}
	prefBinds[key] = listen
	return listen
}

func (b *prefBound{{ .Name }}) Get() {{ .Type }} {
	return b.p.{{ .Name }}(b.key)
}

func (b *prefBound{{ .Name }}) Set(v {{ .Type }}) {
	b.p.Set{{ .Name }}(b.key, v)

	b.trigger()
}
`

const toStringTemplate = `
type stringFrom{{ .Name }} struct {
	base

	format string
	from   {{ .Name }}
}

// {{ .Name }}ToString creates a binding that connects a {{ .Name }} data item to a String.
// Changes to the {{ .Name }} will be pushed to the String and setting the string will parse and set the
// {{ .Name }} if the parse was successful.
func {{ .Name }}ToString(v {{ .Name }}) String {
	return {{ .Name }}ToStringWithFormat(v, "{{ .Format }}")
}

// {{ .Name }}ToStringWithFormat creates a binding that connects a {{ .Name }} data item to a String and is
// presented using the specified format. Changes to the {{ .Name }} will be pushed to the String and setting
// the string will parse and set the {{ .Name }} if the string matches the format and its parse was successful.
func {{ .Name }}ToStringWithFormat(v {{ .Name }}, format string) String {
	str := &stringFrom{{ .Name }}{from: v, format: format}
	v.AddListener(str)
	return str
}

func (s *stringFrom{{ .Name }}) Get() string {
	val := s.from.Get()

	return fmt.Sprintf(s.format, val)
}

func (s *stringFrom{{ .Name }}) Set(str string) {
	var val {{ .Type }}
	n, err := fmt.Sscanf(str, s.format, &val)
	if err != nil || n != 1 {
		fyne.LogError("{{ .Type }} parse error", err)
		return
	}
	if val == s.from.Get() {
		return
	}
	s.from.Set(val)

	s.trigger()
}

func (s *stringFrom{{ .Name }}) DataChanged() {
	s.trigger()
}
`

type bindValues struct {
	Name, Type, Default string
	Format              string
	SupportsPreferences bool
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
	prefFile, err := newFile("preference")
	if err != nil {
		return
	}
	defer prefFile.Close()
	prefFile.WriteString(`
import "fyne.io/fyne"

const keyTypeMismatchError = "A previous preference binding exists with different type for key: "

// Because there is no preference listener yet we connect any listeners asking for the same key.
var prefBinds = make(map[string]DataItem)
`)

	item := template.Must(template.New("item").Parse(itemBindTemplate))
	toString := template.Must(template.New("toString").Parse(toStringTemplate))
	preference := template.Must(template.New("preference").Parse(prefTemplate))
	for _, b := range []bindValues{
		bindValues{Name: "Bool", Type: "bool", Default: "false", Format: "%t", SupportsPreferences: true},
		bindValues{Name: "Float", Type: "float64", Default: "0.0", Format: "%f", SupportsPreferences: true},
		bindValues{Name: "Int", Type: "int", Default: "0", Format: "%d", SupportsPreferences: true},
		bindValues{Name: "Rune", Type: "rune", Default: "rune(0)"},
		bindValues{Name: "String", Type: "string", Default: "\"\"", SupportsPreferences: true},
	} {
		writeFile(itemFile, item, b)
		if b.SupportsPreferences {
			writeFile(prefFile, preference, b)
		}
		if b.Format != "" {
			writeFile(toStringFile, toString, b)
		}
	}
}
