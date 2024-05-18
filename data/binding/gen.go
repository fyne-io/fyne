//go:build ignore

package main

import (
	"os"
	"path"
	"runtime"
	"text/template"

	"fyne.io/fyne/v2"
)

const itemBindTemplate = `
// {{ .Name }} supports binding a {{ .Type }} value.
//
// Since: {{ .Since }}
type {{ .Name }} interface {
	DataItem
	Get() ({{ .Type }}, error)
	Set({{ .Type }}) error
}

// External{{ .Name }} supports binding a {{ .Type }} value to an external value.
//
// Since: {{ .Since }}
type External{{ .Name }} interface {
	{{ .Name }}
	Reload() error
}

// New{{ .Name }} returns a bindable {{ .Type }} value that is managed internally.
//
// Since: {{ .Since }}
func New{{ .Name }}() {{ .Name }} {
	var blank {{ .Type }} = {{ .Default }}
	return &bound{{ .Name }}{val: &blank}
}

// Bind{{ .Name }} returns a new bindable value that controls the contents of the provided {{ .Type }} variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
//
// Since: {{ .Since }}
func Bind{{ .Name }}(v *{{ .Type }}) External{{ .Name }} {
	if v == nil {
		var blank {{ .Type }} = {{ .Default }}
		v = &blank // never allow a nil value pointer
	}
	b := &boundExternal{{ .Name }}{}
	b.val = v
	b.old = *v
	return b
}

type bound{{ .Name }} struct {
	base

	val *{{ .Type }}
}

func (b *bound{{ .Name }}) Get() ({{ .Type }}, error) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if b.val == nil {
		return {{ .Default }}, nil
	}
	return *b.val, nil
}

func (b *bound{{ .Name }}) Set(val {{ .Type }}) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	{{- if eq .Comparator "" }}
	if *b.val == val {
		return nil
	}
	{{- else }}
	if {{ .Comparator }}(*b.val, val) {
		return nil
	}
	{{- end }}
	*b.val = val

	b.trigger()
	return nil
}

type boundExternal{{ .Name }} struct {
	bound{{ .Name }}

	old {{ .Type }}
}

func (b *boundExternal{{ .Name }}) Set(val {{ .Type }}) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	{{- if eq .Comparator "" }}
	if b.old == val {
		return nil
	}
	{{- else }}
	if {{ .Comparator }}(b.old, val) {
		return nil
	}
	{{- end }}
	*b.val = val
	b.old = val

	b.trigger()
	return nil
}

func (b *boundExternal{{ .Name }}) Reload() error {
	return b.Set(*b.val)
}
`

const prefTemplate = `
type prefBound{{ .Name }} struct {
	base
	key   string
	p     fyne.Preferences
	cache atomic.Pointer[{{ .Type }}]
}

// BindPreference{{ .Name }} returns a bindable {{ .Type }} value that is managed by the application preferences.
// Changes to this value will be saved to application storage and when the app starts the previous values will be read.
//
// Since: {{ .Since }}
func BindPreference{{ .Name }}(key string, p fyne.Preferences) {{ .Name }} {
	binds := prefBinds.getBindings(p)
	if binds != nil {
		if listen := binds.getItem(key); listen != nil {
			if l, ok := listen.({{ .Name }}); ok {
				return l
			}
			fyne.LogError(keyTypeMismatchError+key, nil)
		}
	}

	listen := &prefBound{{ .Name }}{key: key, p: p}
	binds = prefBinds.ensurePreferencesAttached(p)
	binds.setItem(key, listen)
	return listen
}

func (b *prefBound{{ .Name }}) Get() ({{ .Type }}, error) {
	cache := b.p.{{ .Name }}(b.key)
	b.cache.Store(&cache)
	return cache, nil
}

func (b *prefBound{{ .Name }}) Set(v {{ .Type }}) error {
	b.p.Set{{ .Name }}(b.key, v)

	b.lock.RLock()
	defer b.lock.RUnlock()
	b.trigger()
	return nil
}

func (b *prefBound{{ .Name }}) checkForChange() {
	val := b.cache.Load()
	if val != nil && b.p.{{ .Name }}(b.key) == *val {
		return
	}
	b.trigger()
}

func (b *prefBound{{ .Name }}) replaceProvider(p fyne.Preferences) {
	b.p = p
}
`

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

	s.DataChanged()
	return nil
}

func (s *stringFrom{{ .Name }}) DataChanged() {
	s.lock.RLock()
	defer s.lock.RUnlock()
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

	s.DataChanged()
	return nil
}

func (s *intFrom{{ .Name }}) DataChanged() {
	s.lock.RLock()
	defer s.lock.RUnlock()
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

	s.DataChanged()
	return nil
}

func (s *intTo{{ .Name }}) DataChanged() {
	s.lock.RLock()
	defer s.lock.RUnlock()
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

	s.DataChanged()
	return nil
}

func (s *stringTo{{ .Name }}) DataChanged() {
	s.lock.RLock()
	defer s.lock.RUnlock()
	s.trigger()
}
`

const listBindTemplate = `
// {{ .Name }}List supports binding a list of {{ .Type }} values.
//
// Since: {{ .Since }}
type {{ .Name }}List interface {
	DataList

	Append(value {{ .Type }}) error
	Get() ([]{{ .Type }}, error)
	GetValue(index int) ({{ .Type }}, error)
	Prepend(value {{ .Type }}) error
	Remove(value {{ .Type }}) error
	Set(list []{{ .Type }}) error
	SetValue(index int, value {{ .Type }}) error
}

// External{{ .Name }}List supports binding a list of {{ .Type }} values from an external variable.
//
// Since: {{ .Since }}
type External{{ .Name }}List interface {
	{{ .Name }}List

	Reload() error
}

// New{{ .Name }}List returns a bindable list of {{ .Type }} values.
//
// Since: {{ .Since }}
func New{{ .Name }}List() {{ .Name }}List {
	return &bound{{ .Name }}List{val: &[]{{ .Type }}{}}
}

// Bind{{ .Name }}List returns a bound list of {{ .Type }} values, based on the contents of the passed slice.
// If your code changes the content of the slice this refers to you should call Reload() to inform the bindings.
//
// Since: {{ .Since }}
func Bind{{ .Name }}List(v *[]{{ .Type }}) External{{ .Name }}List {
	if v == nil {
		return New{{ .Name }}List().(External{{ .Name }}List)
	}

	b := &bound{{ .Name }}List{val: v, updateExternal: true}

	for i := range *v {
		b.appendItem(bind{{ .Name }}ListItem(v, i, b.updateExternal))
	}

	return b
}

type bound{{ .Name }}List struct {
	listBase

	updateExternal bool
	val            *[]{{ .Type }}
}

func (l *bound{{ .Name }}List) Append(val {{ .Type }}) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	*l.val = append(*l.val, val)

	return l.doReload()
}

func (l *bound{{ .Name }}List) Get() ([]{{ .Type }}, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	return *l.val, nil
}

func (l *bound{{ .Name }}List) GetValue(i int) ({{ .Type }}, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	if i < 0 || i >= l.Length() {
		return {{ .Default }}, errOutOfBounds
	}

	return (*l.val)[i], nil
}

func (l *bound{{ .Name }}List) Prepend(val {{ .Type }}) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	*l.val = append([]{{ .Type }}{val}, *l.val...)

	return l.doReload()
}

func (l *bound{{ .Name }}List) Reload() error {
	l.lock.Lock()
	defer l.lock.Unlock()

	return l.doReload()
}

// Remove takes the specified {{ .Type }} out of the list.
//
// Since: 2.5
func (l *bound{{ .Name }}List) Remove(val {{ .Type }}) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	v := *l.val
	if len(v) == 0 {
		return nil
	}

	{{- if eq .Comparator "" }}
	if v[0] == val {
		*l.val = v[1:]
	} else if v[len(v)-1] == val {
		*l.val = v[:len(v)-1]
	} else {
	{{- else }}
	if {{ .Comparator }}(v[0], val) {
		*l.val = v[1:]
	} else if {{ .Comparator }}(v[len(v)-1], val) {
		*l.val = v[:len(v)-1]
	} else {
	{{- end }}
		id := -1
		for i, v := range v {
		{{- if eq .Comparator "" }}
			if v == val {
				id = i
				break
			}
		{{- else }}
			if {{ .Comparator }}(v, val) {
				id = i
				break
			}
		{{- end }}
		}

		if id == -1 {
			return nil
		}
		*l.val = append(v[:id], v[id+1:]...)
	}

	return l.doReload()
}

func (l *bound{{ .Name }}List) Set(v []{{ .Type }}) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	*l.val = v

	return l.doReload()
}

func (l *bound{{ .Name }}List) doReload() (retErr error) {
	oldLen := len(l.items)
	newLen := len(*l.val)
	if oldLen > newLen {
		for i := oldLen - 1; i >= newLen; i-- {
			l.deleteItem(i)
		}
		l.trigger()
	} else if oldLen < newLen {
		for i := oldLen; i < newLen; i++ {
			l.appendItem(bind{{ .Name }}ListItem(l.val, i, l.updateExternal))
		}
		l.trigger()
	}

	for i, item := range l.items {
		if i > oldLen || i > newLen {
			break
		}

		var err error
		if l.updateExternal {
			item.(*boundExternal{{ .Name }}ListItem).lock.Lock()
			err = item.(*boundExternal{{ .Name }}ListItem).setIfChanged((*l.val)[i])
			item.(*boundExternal{{ .Name }}ListItem).lock.Unlock()
		} else {
			item.(*bound{{ .Name }}ListItem).lock.Lock()
			err = item.(*bound{{ .Name }}ListItem).doSet((*l.val)[i])
			item.(*bound{{ .Name }}ListItem).lock.Unlock()
		}
		if err != nil {
			retErr = err
		}
	}
	return
}

func (l *bound{{ .Name }}List) SetValue(i int, v {{ .Type }}) error {
	l.lock.RLock()
	len := l.Length()
	l.lock.RUnlock()

	if i < 0 || i >= len {
		return errOutOfBounds
	}

	l.lock.Lock()
	(*l.val)[i] = v
	l.lock.Unlock()

	item, err := l.GetItem(i)
	if err != nil {
		return err
	}
	return item.({{ .Name }}).Set(v)
}

func bind{{ .Name }}ListItem(v *[]{{ .Type }}, i int, external bool) {{ .Name }} {
	if external {
		ret := &boundExternal{{ .Name }}ListItem{old: (*v)[i]}
		ret.val = v
		ret.index = i
		return ret
	}

	return &bound{{ .Name }}ListItem{val: v, index: i}
}

type bound{{ .Name }}ListItem struct {
	base

	val   *[]{{ .Type }}
	index int
}

func (b *bound{{ .Name }}ListItem) Get() ({{ .Type }}, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.index < 0 || b.index >= len(*b.val) {
		return {{ .Default }}, errOutOfBounds
	}

	return (*b.val)[b.index], nil
}

func (b *bound{{ .Name }}ListItem) Set(val {{ .Type }}) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.doSet(val)
}

func (b *bound{{ .Name }}ListItem) doSet(val {{ .Type }}) error {
	(*b.val)[b.index] = val

	b.trigger()
	return nil
}

type boundExternal{{ .Name }}ListItem struct {
	bound{{ .Name }}ListItem

	old {{ .Type }}
}

func (b *boundExternal{{ .Name }}ListItem) setIfChanged(val {{ .Type }}) error {
	{{- if eq .Comparator "" }}
	if val == b.old {
		return nil
	}
	{{- else }}
	if {{ .Comparator }}(val, b.old) {
		return nil
	}
	{{- end }}
	(*b.val)[b.index] = val
	b.old = val

	b.trigger()
	return nil
}
`

const treeBindTemplate = `
// {{ .Name }}Tree supports binding a tree of {{ .Type }} values.
//
{{ if eq .Name "Untyped" -}}
// Since: 2.5
{{- else -}}
// Since: 2.4
{{- end }}
type {{ .Name }}Tree interface {
	DataTree

	Append(parent, id string, value {{ .Type }}) error
	Get() (map[string][]string, map[string]{{ .Type }}, error)
	GetValue(id string) ({{ .Type }}, error)
	Prepend(parent, id string, value {{ .Type }}) error
	Remove(id string) error
	Set(ids map[string][]string, values map[string]{{ .Type }}) error
	SetValue(id string, value {{ .Type }}) error
}

// External{{ .Name }}Tree supports binding a tree of {{ .Type }} values from an external variable.
//
{{ if eq .Name "Untyped" -}}
// Since: 2.5
{{- else -}}
// Since: 2.4
{{- end }}
type External{{ .Name }}Tree interface {
	{{ .Name }}Tree

	Reload() error
}

// New{{ .Name }}Tree returns a bindable tree of {{ .Type }} values.
//
{{ if eq .Name "Untyped" -}}
// Since: 2.5
{{- else -}}
// Since: 2.4
{{- end }}
func New{{ .Name }}Tree() {{ .Name }}Tree {
	t := &bound{{ .Name }}Tree{val: &map[string]{{ .Type }}{}}
	t.ids = make(map[string][]string)
	t.items = make(map[string]DataItem)
	return t
}

// Bind{{ .Name }}Tree returns a bound tree of {{ .Type }} values, based on the contents of the passed values.
// The ids map specifies how each item relates to its parent (with id ""), with the values being in the v map.
// If your code changes the content of the maps this refers to you should call Reload() to inform the bindings.
//
// Since: 2.4
func Bind{{ .Name }}Tree(ids *map[string][]string, v *map[string]{{ .Type }}) External{{ .Name }}Tree {
	if v == nil {
		return New{{ .Name }}Tree().(External{{ .Name }}Tree)
	}

	t := &bound{{ .Name }}Tree{val: v, updateExternal: true}
	t.ids = make(map[string][]string)
	t.items = make(map[string]DataItem)

	for parent, children := range *ids {
		for _, leaf := range children {
			t.appendItem(bind{{ .Name }}TreeItem(v, leaf, t.updateExternal), leaf, parent)
		}
	}

	return t
}

type bound{{ .Name }}Tree struct {
	treeBase

	updateExternal bool
	val            *map[string]{{ .Type }}
}

func (t *bound{{ .Name }}Tree) Append(parent, id string, val {{ .Type }}) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	ids, ok := t.ids[parent]
	if !ok {
		ids = make([]string, 0)
	}

	t.ids[parent] = append(ids, id)
	v := *t.val
	v[id] = val

	return t.doReload()
}

func (t *bound{{ .Name }}Tree) Get() (map[string][]string, map[string]{{ .Type }}, error) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.ids, *t.val, nil
}

func (t *bound{{ .Name }}Tree) GetValue(id string) ({{ .Type }}, error) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	if item, ok := (*t.val)[id]; ok {
		return item, nil
	}

	return {{ .Default }}, errOutOfBounds
}

func (t *bound{{ .Name }}Tree) Prepend(parent, id string, val {{ .Type }}) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	ids, ok := t.ids[parent]
	if !ok {
		ids = make([]string, 0)
	}

	t.ids[parent] = append([]string{id}, ids...)
	v := *t.val
	v[id] = val

	return t.doReload()
}

// Remove takes the specified id out of the tree.
// It will also remove any child items from the data structure.
//
// Since: 2.5
func (t *bound{{ .Name }}Tree) Remove(id string) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.removeChildren(id)
	delete(t.ids, id)
	v := *t.val
	delete(v, id)

	return t.doReload()
}

func (t *bound{{ .Name }}Tree) removeChildren(id string) {
	for _, cid := range t.ids[id] {
		t.removeChildren(cid)

		delete(t.ids, cid)
		v := *t.val
		delete(v, cid)
	}
}

func (t *bound{{ .Name }}Tree) Reload() error {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.doReload()
}

func (t *bound{{ .Name }}Tree) Set(ids map[string][]string, v map[string]{{ .Type }}) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.ids = ids
	*t.val = v

	return t.doReload()
}

func (t *bound{{ .Name }}Tree) doReload() (retErr error) {
	updated := []string{}
	fire := false
	for id := range *t.val {
		found := false
		for child := range t.items {
			if child == id { // update existing
				updated = append(updated, id)
				found = true
				break
			}
		}
		if found {
			continue
		}

		// append new
		t.appendItem(bind{{ .Name }}TreeItem(t.val, id, t.updateExternal), id, parentIDFor(id, t.ids))
		updated = append(updated, id)
		fire = true
	}

	for id := range t.items {
		remove := true
		for _, done := range updated {
			if done == id {
				remove = false
				break
			}
		}

		if remove { // remove item no longer present
			fire = true
			t.deleteItem(id, parentIDFor(id, t.ids))
		}
	}
	if fire {
		t.trigger()
	}

	for id, item := range t.items {
		var err error
		if t.updateExternal {
			item.(*boundExternal{{ .Name }}TreeItem).lock.Lock()
			err = item.(*boundExternal{{ .Name }}TreeItem).setIfChanged((*t.val)[id])
			item.(*boundExternal{{ .Name }}TreeItem).lock.Unlock()
		} else {
			item.(*bound{{ .Name }}TreeItem).lock.Lock()
			err = item.(*bound{{ .Name }}TreeItem).doSet((*t.val)[id])
			item.(*bound{{ .Name }}TreeItem).lock.Unlock()
		}
		if err != nil {
			retErr = err
		}
	}
	return
}

func (t *bound{{ .Name }}Tree) SetValue(id string, v {{ .Type }}) error {
	t.lock.Lock()
	(*t.val)[id] = v
	t.lock.Unlock()

	item, err := t.GetItem(id)
	if err != nil {
		return err
	}
	return item.({{ .Name }}).Set(v)
}

func bind{{ .Name }}TreeItem(v *map[string]{{ .Type }}, id string, external bool) {{ .Name }} {
	if external {
		ret := &boundExternal{{ .Name }}TreeItem{old: (*v)[id]}
		ret.val = v
		ret.id = id
		return ret
	}

	return &bound{{ .Name }}TreeItem{id: id, val: v}
}

type bound{{ .Name }}TreeItem struct {
	base

	val *map[string]{{ .Type }}
	id  string
}

func (t *bound{{ .Name }}TreeItem) Get() ({{ .Type }}, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	v := *t.val
	if item, ok := v[t.id]; ok {
		return item, nil
	}

	return {{ .Default }}, errOutOfBounds
}

func (t *bound{{ .Name }}TreeItem) Set(val {{ .Type }}) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.doSet(val)
}

func (t *bound{{ .Name }}TreeItem) doSet(val {{ .Type }}) error {
	(*t.val)[t.id] = val

	t.trigger()
	return nil
}

type boundExternal{{ .Name }}TreeItem struct {
	bound{{ .Name }}TreeItem

	old {{ .Type }}
}

func (t *boundExternal{{ .Name }}TreeItem) setIfChanged(val {{ .Type }}) error {
	{{- if eq .Comparator "" }}
	if val == t.old {
		return nil
	}
	{{- else }}
	if {{ .Comparator }}(val, t.old) {
		return nil
	}
	{{- end }}
	(*t.val)[t.id] = val
	t.old = val

	t.trigger()
	return nil
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
	itemFile, err := newFile("binditems")
	if err != nil {
		return
	}
	defer itemFile.Close()
	itemFile.WriteString(`
import (
	"bytes"

	"fyne.io/fyne/v2"
)
`)
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
	prefFile, err := newFile("preference")
	if err != nil {
		return
	}
	defer prefFile.Close()
	prefFile.WriteString(`
import (
	"sync/atomic"

	"fyne.io/fyne/v2"
)

const keyTypeMismatchError = "A previous preference binding exists with different type for key: "
`)

	listFile, err := newFile("bindlists")
	if err != nil {
		return
	}
	defer listFile.Close()
	listFile.WriteString(`
import (
	"bytes"

	"fyne.io/fyne/v2"
)
`)

	treeFile, err := newFile("bindtrees")
	if err != nil {
		return
	}
	defer treeFile.Close()
	treeFile.WriteString(`
import (
	"bytes"

	"fyne.io/fyne/v2"
)
`)

	item := template.Must(template.New("item").Parse(itemBindTemplate))
	fromString := template.Must(template.New("fromString").Parse(fromStringTemplate))
	fromInt := template.Must(template.New("fromInt").Parse(fromIntTemplate))
	toInt := template.Must(template.New("toInt").Parse(toIntTemplate))
	toString := template.Must(template.New("toString").Parse(toStringTemplate))
	preference := template.Must(template.New("preference").Parse(prefTemplate))
	list := template.Must(template.New("list").Parse(listBindTemplate))
	tree := template.Must(template.New("tree").Parse(treeBindTemplate))
	binds := []bindValues{
		{Name: "Bool", Type: "bool", Default: "false", Format: "%t", SupportsPreferences: true},
		{Name: "Bytes", Type: "[]byte", Default: "nil", Since: "2.2", Comparator: "bytes.Equal"},
		{Name: "Float", Type: "float64", Default: "0.0", Format: "%f", SupportsPreferences: true, ToInt: "internalFloatToInt", FromInt: "internalIntToFloat"},
		{Name: "Int", Type: "int", Default: "0", Format: "%d", SupportsPreferences: true},
		{Name: "Rune", Type: "rune", Default: "rune(0)"},
		{Name: "String", Type: "string", Default: "\"\"", SupportsPreferences: true},
		{Name: "Untyped", Type: "any", Default: "nil", Since: "2.1"},
		{Name: "URI", Type: "fyne.URI", Default: "fyne.URI(nil)", Since: "2.1",
			FromString: "uriFromString", ToString: "uriToString", Comparator: "compareURI"},
	}
	for _, b := range binds {
		if b.Since == "" {
			b.Since = "2.0"
		}

		writeFile(listFile, list, b)
		writeFile(treeFile, tree, b)
		if b.Name == "Untyped" {
			continue // any is special, we have it in binding.go instead
		}

		writeFile(itemFile, item, b)
		if b.SupportsPreferences {
			writeFile(prefFile, preference, b)
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
