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
// {{ .Name }} supports binding a {{ .Type }} value.
//
// Since: 2.0.0
type {{ .Name }} interface {
	DataItem
	Get() ({{ .Type }}, error)
	Set({{ .Type }}) error
}

// New{{ .Name }} returns a bindable {{ .Type }} value that is managed internally.
//
// Since: 2.0.0
func New{{ .Name }}() {{ .Name }} {
	blank := {{ .Default }}
	return &bound{{ .Name }}{val: &blank}
}

// Bind{{ .Name }} returns a new bindable value that controls the contents of the provided {{ .Type }} variable.
//
// Since: 2.0.0
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

func (b *bound{{ .Name }}) Get() ({{ .Type }}, error) {
	if b.val == nil {
		return {{ .Default }}, nil
	}
	return *b.val, nil
}

func (b *bound{{ .Name }}) Set(val {{ .Type }}) error {
	if *b.val == val {
		return nil
	}
	if b.val == nil { // was not initialized with a blank value, recover
		b.val = &val
	} else {
		*b.val = val
	}

	b.trigger()
	return nil
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
//
// Since: 2.0.0
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

func (b *prefBound{{ .Name }}) Get() ({{ .Type }}, error) {
	return b.p.{{ .Name }}(b.key), nil
}

func (b *prefBound{{ .Name }}) Set(v {{ .Type }}) error {
	b.p.Set{{ .Name }}(b.key, v)

	b.trigger()
	return nil
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
//
// Since: 2.0.0
func {{ .Name }}ToString(v {{ .Name }}) String {
	return {{ .Name }}ToStringWithFormat(v, "{{ .Format }}")
}

// {{ .Name }}ToStringWithFormat creates a binding that connects a {{ .Name }} data item to a String and is
// presented using the specified format. Changes to the {{ .Name }} will be pushed to the String and setting
// the string will parse and set the {{ .Name }} if the string matches the format and its parse was successful.
//
// Since: 2.0.0
func {{ .Name }}ToStringWithFormat(v {{ .Name }}, format string) String {
	str := &stringFrom{{ .Name }}{from: v, format: format}
	v.AddListener(str)
	return str
}

func (s *stringFrom{{ .Name }}) Get() (string, error) {
	val, err := s.from.Get()

	return fmt.Sprintf(s.format, val), err
}

func (s *stringFrom{{ .Name }}) Set(str string) error {
	var val {{ .Type }}
	n, err := fmt.Sscanf(str, s.format, &val)
	if err != nil || n != 1 {
		return err
	}
	old, err := s.from.Get()
	if val == old {
		return err
	}
	err = s.from.Set(val)

	s.trigger()
	return err
}

func (s *stringFrom{{ .Name }}) DataChanged() {
	s.trigger()
}
`

const fromStringTemplate = `
type stringTo{{ .Name }} struct {
	base

	format string
	from   String
}

// StringTo{{ .Name }} creates a binding that connects a String data item to a {{ .Name }}.
// Changes to the String will be parsed and pushed to the {{ .Name }} if the parse was successful, and setting
// the {{ .Name }} update the String binding.
//
// Since: 2.0.0
func StringTo{{ .Name }}(str String) {{ .Name }} {
	return StringTo{{ .Name }}WithFormat(str, "{{ .Format }}")
}

// StringTo{{ .Name }}WithFormat creates a binding that connects a String data item to a {{ .Name }} and is
// presented using the specified format. Changes to the {{ .Name }} will be parsed and if the format matches and
// the parse is successful it will be pushed to the String. Setting the {{ .Name }} will push a formatted value
// into the String.
//
// Since: 2.0.0
func StringTo{{ .Name }}WithFormat(str String, format string) {{ .Name }} {
	v := &stringTo{{ .Name }}{from: str, format: format}
	str.AddListener(v)
	return v
}

func (s *stringTo{{ .Name }}) Get() ({{ .Type }}, error) {
	str, err := s.from.Get()
	if str == "" || err != nil {
		return {{ .Default }}, nil
	}

	var val {{ .Type }}
	n, err := fmt.Sscanf(str, s.format, &val)
	if err != nil || n != 1 {
		fyne.LogError("{{ .Type }} parse error", err)
		return {{ .Default }}, err
	}

	return val, nil
}

func (s *stringTo{{ .Name }}) Set(val {{ .Type }}) error {
	str := fmt.Sprintf(s.format, val)
	old, err := s.from.Get()
	if str == old {
		return err
	}

	err = s.from.Set(str)
	s.trigger()
	return err
}

func (s *stringTo{{ .Name }}) DataChanged() {
	s.trigger()
}
`

const listBindTemplate = `
// {{ .Name }}List supports binding a list of {{ .Type }} values.
//
// Since: 2.0.0
type {{ .Name }}List interface {
	DataList

	Append({{ .Type }}) error
	Get(int) ({{ .Type }}, error)
	Prepend({{ .Type }}) error
	Set(int, {{ .Type }}) error
}

// New{{ .Name }}List returns a bindable list of {{ .Type }} values.
//
// Since: 2.0.0
func New{{ .Name }}List() {{ .Name }}List {
	return &bound{{ .Name }}List{}
}

// Bind{{ .Name }}List returns a bound list of {{ .Type }} values, based on the contents of the passed slice.
//
// Since: 2.0.0
func Bind{{ .Name }}List(v *[]{{ .Type }}) {{ .Name }}List {
	if v == nil {
		return New{{ .Name }}List()
	}

	b := &bound{{ .Name }}List{val: v}

	for i := range *v {
		b.appendItem(Bind{{ .Name }}(&((*v)[i])))
	}

	return b
}

type bound{{ .Name }}List struct {
	listBase

	val *[]{{ .Type }}
}

func (l *bound{{ .Name }}List) Append(val {{ .Type }}) error {
	if l.val != nil {
		*l.val = append(*l.val, val)
	}

	l.appendItem(Bind{{ .Name }}(&val))
	return nil
}

func (l *bound{{ .Name }}List) Get(i int) ({{ .Type }}, error) {
	if i < 0 || i >= l.Length() {
		return {{ .Default }}, errOutOfBounds
	}
	if l.val != nil {
		return (*l.val)[i], nil
	}

	item, ok := l.GetItem(i)
	if !ok {
		return {{ .Default }}, errOutOfBounds
	}
	return item.({{ .Name }}).Get()
}

func (l *bound{{ .Name }}List) Prepend(val {{ .Type }}) error {
	if l.val != nil {
		*l.val = append([]{{ .Type }}{val}, *l.val...)
	}

	l.prependItem(Bind{{ .Name }}(&val))
	return nil
}

func (l *bound{{ .Name }}List) Set(i int, v {{ .Type }}) error {
	if i < 0 || i >= l.Length() {
		return errOutOfBounds
	}
	if l.val != nil {
		(*l.val)[i] = v
	}

	item, ok := l.GetItem(i)
	if !ok {
		return errOutOfBounds
	}
	return item.({{ .Name }}).Set(v)
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
	convertFile, err := newFile("convert")
	if err != nil {
		return
	}
	defer convertFile.Close()
	convertFile.WriteString(`
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

	listFile, err := newFile("bindlists")
	if err != nil {
		return
	}
	defer listFile.Close()

	item := template.Must(template.New("item").Parse(itemBindTemplate))
	fromString := template.Must(template.New("fromString").Parse(fromStringTemplate))
	toString := template.Must(template.New("toString").Parse(toStringTemplate))
	preference := template.Must(template.New("preference").Parse(prefTemplate))
	list := template.Must(template.New("list").Parse(listBindTemplate))
	binds := []bindValues{
		bindValues{Name: "Bool", Type: "bool", Default: "false", Format: "%t", SupportsPreferences: true},
		bindValues{Name: "Float", Type: "float64", Default: "0.0", Format: "%f", SupportsPreferences: true},
		bindValues{Name: "Int", Type: "int", Default: "0", Format: "%d", SupportsPreferences: true},
		bindValues{Name: "Rune", Type: "rune", Default: "rune(0)"},
		bindValues{Name: "String", Type: "string", Default: "\"\"", SupportsPreferences: true},
	}
	for _, b := range binds {
		writeFile(itemFile, item, b)
		if b.SupportsPreferences {
			writeFile(prefFile, preference, b)
		}
		if b.Format != "" {
			writeFile(convertFile, toString, b)
		}
		writeFile(listFile, list, b)
	}
	// add StringTo... at the bottom of the convertFile for correct ordering
	for _, b := range binds {
		if b.Format != "" {
			writeFile(convertFile, fromString, b)
		}
	}
}
