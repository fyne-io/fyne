package binding

import (
	"fmt"
	"sync/atomic"

	"fyne.io/fyne/v2"
)

type bindableItem[T any] interface {
	DataItem
	Get() (T, error)
	Set(T) error
}

func newBaseItem[T any](comparator func(T, T) bool) *baseItem[T] {
	return &baseItem[T]{val: new(T), comparator: comparator}
}

func newBaseItemComparable[T bool | float64 | int | rune | string]() *baseItem[T] {
	return newBaseItem[T](func(a, b T) bool { return a == b })
}

type baseItem[T any] struct {
	base

	comparator func(T, T) bool
	val        *T
}

func (b *baseItem[T]) Get() (T, error) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if b.val == nil {
		return *new(T), nil
	}
	return *b.val, nil
}

func (b *baseItem[T]) Set(val T) error {
	b.lock.Lock()
	equal := b.comparator(*b.val, val)
	*b.val = val
	b.lock.Unlock()

	if !equal {
		b.trigger()
	}

	return nil
}

func baseBindExternal[T any](val *T, comparator func(T, T) bool) *baseExternalItem[T] {
	if val == nil {
		val = new(T) // never allow a nil value pointer
	}
	b := &baseExternalItem[T]{}
	b.comparator = comparator
	b.val = val
	b.old = *val
	return b
}

func baseBindExternalComparable[T bool | float64 | int | rune | string](val *T) *baseExternalItem[T] {
	if val == nil {
		val = new(T) // never allow a nil value pointer
	}
	b := &baseExternalItem[T]{}
	b.comparator = func(a, b T) bool { return a == b }
	b.val = val
	b.old = *val
	return b
}

type baseExternalItem[T any] struct {
	baseItem[T]

	old T
}

func (b *baseExternalItem[T]) Set(val T) error {
	b.lock.Lock()
	if b.comparator(b.old, val) {
		b.lock.Unlock()
		return nil
	}
	*b.val = val
	b.old = val
	b.lock.Unlock()

	b.trigger()
	return nil
}

func (b *baseExternalItem[T]) Reload() error {
	return b.Set(*b.val)
}

type prefBoundBase[T bool | float64 | int | string] struct {
	base
	key   string
	get   func(string) T
	set   func(string, T)
	cache atomic.Pointer[T]
}

func (b *prefBoundBase[T]) Get() (T, error) {
	cache := b.get(b.key)
	b.cache.Store(&cache)
	return cache, nil
}

func (b *prefBoundBase[T]) Set(v T) error {
	b.set(b.key, v)

	b.lock.RLock()
	defer b.lock.RUnlock()
	b.trigger()
	return nil
}

func (b *prefBoundBase[T]) setKey(key string) {
	b.key = key
}

func (b *prefBoundBase[T]) checkForChange() {
	val := b.cache.Load()
	if val != nil && b.get(b.key) == *val {
		return
	}
	b.trigger()
}

type genericItem[T any] interface {
	DataItem
	Get() (T, error)
	Set(T) error
}

func lookupExistingBinding[T any](key string, p fyne.Preferences) (genericItem[T], bool) {
	binds := prefBinds.getBindings(p)
	if binds == nil {
		return nil, false
	}

	if listen, ok := binds.Load(key); listen != nil && ok {
		if l, ok := listen.(genericItem[T]); ok {
			return l, ok
		}
		fyne.LogError(keyTypeMismatchError+key, nil)
	}

	return nil, false
}

func newList[T any](comparator func(T, T) bool) *boundList[T] {
	return &boundList[T]{val: new([]T), comparator: comparator}
}

func newListComparable[T bool | float64 | int | rune | string]() *boundList[T] {
	return newList(func(t1, t2 T) bool { return t1 == t2 })
}

func newExternalList[T any](v *[]T, comparator func(T, T) bool) *boundList[T] {
	return &boundList[T]{val: v, comparator: comparator, updateExternal: true}
}

func bindList[T any](v *[]T, comparator func(T, T) bool) *boundList[T] {
	if v == nil {
		return newList(comparator)
	}

	l := newExternalList(v, comparator)
	for i := range *v {
		l.appendItem(bindListItem(v, i, l.updateExternal, comparator))
	}

	return l
}

func bindListComparable[T bool | float64 | int | rune | string](v *[]T) *boundList[T] {
	return bindList(v, func(t1, t2 T) bool { return t1 == t2 })
}

type boundList[T any] struct {
	listBase

	comparator     func(T, T) bool
	updateExternal bool
	val            *[]T
}

func (l *boundList[T]) Append(val T) error {
	l.lock.Lock()
	*l.val = append(*l.val, val)

	trigger, err := l.doReload()
	l.lock.Unlock()

	if trigger {
		l.trigger()
	}

	return err
}

func (l *boundList[T]) Get() ([]T, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	return *l.val, nil
}

func (l *boundList[T]) GetValue(i int) (T, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	if i < 0 || i >= l.Length() {
		return *new(T), errOutOfBounds
	}

	return (*l.val)[i], nil
}

func (l *boundList[T]) Prepend(val T) error {
	l.lock.Lock()
	*l.val = append([]T{val}, *l.val...)

	trigger, err := l.doReload()
	l.lock.Unlock()

	if trigger {
		l.trigger()
	}

	return err
}

func (l *boundList[T]) Reload() error {
	l.lock.Lock()
	trigger, err := l.doReload()
	l.lock.Unlock()

	if trigger {
		l.trigger()
	}

	return err
}

func (l *boundList[T]) Remove(val T) error {
	l.lock.Lock()

	v := *l.val
	if len(v) == 0 {
		l.lock.Unlock()
		return nil
	}
	if l.comparator(v[0], val) {
		*l.val = v[1:]
	} else if l.comparator(v[len(v)-1], val) {
		*l.val = v[:len(v)-1]
	} else {
		id := -1
		for i, v := range v {
			if l.comparator(v, val) {
				id = i
				break
			}
		}

		if id == -1 {
			l.lock.Unlock()
			return nil
		}
		*l.val = append(v[:id], v[id+1:]...)
	}

	trigger, err := l.doReload()
	l.lock.Unlock()

	if trigger {
		l.trigger()
	}

	return err
}

func (l *boundList[T]) Set(v []T) error {
	l.lock.Lock()
	*l.val = v
	trigger, err := l.doReload()
	l.lock.Unlock()

	if trigger {
		l.trigger()
	}

	return err
}

func (l *boundList[T]) doReload() (trigger bool, retErr error) {
	oldLen := len(l.items)
	newLen := len(*l.val)
	if oldLen > newLen {
		for i := oldLen - 1; i >= newLen; i-- {
			l.deleteItem(i)
		}
		trigger = true
	} else if oldLen < newLen {
		for i := oldLen; i < newLen; i++ {
			l.appendItem(bindListItem(l.val, i, l.updateExternal, l.comparator))
		}
		trigger = true
	}

	for i, item := range l.items {
		if i > oldLen || i > newLen {
			break
		}

		var err error
		if l.updateExternal {
			err = item.(*boundExternalListItem[T]).setIfChanged((*l.val)[i])
		} else {
			err = item.(*boundListItem[T]).doSet((*l.val)[i])
		}
		if err != nil {
			retErr = err
		}
	}
	return
}

func (l *boundList[T]) SetValue(i int, v T) error {
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
	return item.(genericItem[T]).Set(v)
}

func bindListItem[T any](v *[]T, i int, external bool, comparator func(T, T) bool) genericItem[T] {
	if external {
		ret := &boundExternalListItem[T]{old: (*v)[i]}
		ret.val = v
		ret.index = i
		ret.comparator = comparator
		return ret
	}

	return &boundListItem[T]{val: v, index: i, comparator: comparator}
}

func bindListItemComparable[T bool | float64 | int | rune | string](v *[]T, i int, external bool) genericItem[T] {
	return bindListItem(v, i, external, func(t1, t2 T) bool { return t1 == t2 })
}

type boundListItem[T any] struct {
	base

	comparator func(T, T) bool
	val        *[]T
	index      int
}

func (b *boundListItem[T]) Get() (T, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.index < 0 || b.index >= len(*b.val) {
		return *new(T), errOutOfBounds
	}

	return (*b.val)[b.index], nil
}

func (b *boundListItem[T]) Set(val T) error {
	return b.doSet(val)
}

func (b *boundListItem[T]) doSet(val T) error {
	b.lock.Lock()
	(*b.val)[b.index] = val
	b.lock.Unlock()

	b.trigger()
	return nil
}

type boundExternalListItem[T any] struct {
	boundListItem[T]

	old T
}

func (b *boundExternalListItem[T]) setIfChanged(val T) error {
	b.lock.Lock()
	if b.comparator(val, b.old) {
		b.lock.Unlock()
		return nil
	}
	(*b.val)[b.index] = val
	b.old = val

	b.lock.Unlock()
	b.trigger()
	return nil
}

func newTree[T any](comparator func(T, T) bool) *boundTree[T] {
	t := &boundTree[T]{val: &map[string]T{}, comparator: comparator}
	t.ids = make(map[string][]string)
	t.items = make(map[string]DataItem)
	return t
}

func newTreeComparable[T bool | float64 | int | rune | string]() *boundTree[T] {
	return newTree(func(t1, t2 T) bool { return t1 == t2 })
}

func bindTree[T any](ids *map[string][]string, v *map[string]T, comparator func(T, T) bool) *boundTree[T] {
	if v == nil {
		return newTree[T](comparator)
	}

	t := &boundTree[T]{val: v, updateExternal: true, comparator: comparator}
	t.ids = make(map[string][]string)
	t.items = make(map[string]DataItem)

	for parent, children := range *ids {
		for _, leaf := range children {
			t.appendItem(bindTreeItem(v, leaf, t.updateExternal, t.comparator), leaf, parent)
		}
	}

	return t
}

func bindTreeComparable[T bool | float64 | int | rune | string](ids *map[string][]string, v *map[string]T) *boundTree[T] {
	return bindTree(ids, v, func(t1, t2 T) bool { return t1 == t2 })
}

type boundTree[T any] struct {
	treeBase

	comparator     func(T, T) bool
	val            *map[string]T
	updateExternal bool
}

func (t *boundTree[T]) Append(parent, id string, val T) error {
	t.lock.Lock()
	ids, ok := t.ids[parent]
	if !ok {
		ids = make([]string, 0)
	}

	t.ids[parent] = append(ids, id)
	v := *t.val
	v[id] = val

	trigger, err := t.doReload()
	t.lock.Unlock()

	if trigger {
		t.trigger()
	}

	return err
}

func (t *boundTree[T]) Get() (map[string][]string, map[string]T, error) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.ids, *t.val, nil
}

func (t *boundTree[T]) GetValue(id string) (T, error) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	if item, ok := (*t.val)[id]; ok {
		return item, nil
	}

	return *new(T), errOutOfBounds
}

func (t *boundTree[T]) Prepend(parent, id string, val T) error {
	t.lock.Lock()
	ids, ok := t.ids[parent]
	if !ok {
		ids = make([]string, 0)
	}

	t.ids[parent] = append([]string{id}, ids...)
	v := *t.val
	v[id] = val

	trigger, err := t.doReload()
	t.lock.Unlock()

	if trigger {
		t.trigger()
	}

	return err
}

func (t *boundTree[T]) Remove(id string) error {
	t.lock.Lock()
	t.removeChildren(id)
	delete(t.ids, id)
	v := *t.val
	delete(v, id)

	trigger, err := t.doReload()
	t.lock.Unlock()

	if trigger {
		t.trigger()
	}

	return err
}

func (t *boundTree[T]) removeChildren(id string) {
	for _, cid := range t.ids[id] {
		t.removeChildren(cid)

		delete(t.ids, cid)
		v := *t.val
		delete(v, cid)
	}
}

func (t *boundTree[T]) Reload() error {
	t.lock.Lock()
	trigger, err := t.doReload()
	t.lock.Unlock()

	if trigger {
		t.trigger()
	}

	return err
}

func (t *boundTree[T]) Set(ids map[string][]string, v map[string]T) error {
	t.lock.Lock()
	t.ids = ids
	*t.val = v

	trigger, err := t.doReload()
	t.lock.Unlock()

	if trigger {
		t.trigger()
	}

	return err
}

func (t *boundTree[T]) doReload() (fire bool, retErr error) {
	updated := []string{}
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
		t.appendItem(bindTreeItem(t.val, id, t.updateExternal, t.comparator), id, parentIDFor(id, t.ids))
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

	for id, item := range t.items {
		var err error
		if t.updateExternal {
			err = item.(*boundExternalTreeItem[T]).setIfChanged((*t.val)[id])
		} else {
			err = item.(*boundTreeItem[T]).doSet((*t.val)[id])
		}
		if err != nil {
			retErr = err
		}
	}
	return
}

func (t *boundTree[T]) SetValue(id string, v T) error {
	t.lock.Lock()
	(*t.val)[id] = v
	t.lock.Unlock()

	item, err := t.GetItem(id)
	if err != nil {
		return err
	}
	return item.(bindableItem[T]).Set(v)
}

func bindTreeItem[T any](v *map[string]T, id string, external bool, comparator func(T, T) bool) bindableItem[T] {
	if external {
		ret := &boundExternalTreeItem[T]{old: (*v)[id], comparator: comparator}
		ret.val = v
		ret.id = id
		return ret
	}

	return &boundTreeItem[T]{id: id, val: v}
}

type boundTreeItem[T any] struct {
	base

	val *map[string]T
	id  string
}

func (t *boundTreeItem[T]) Get() (T, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	v := *t.val
	if item, ok := v[t.id]; ok {
		return item, nil
	}

	return *new(T), errOutOfBounds
}

func (t *boundTreeItem[T]) Set(val T) error {
	return t.doSet(val)
}

func (t *boundTreeItem[T]) doSet(val T) error {
	t.lock.Lock()
	(*t.val)[t.id] = val
	t.lock.Unlock()

	t.trigger()
	return nil
}

type boundExternalTreeItem[T any] struct {
	boundTreeItem[T]

	comparator func(T, T) bool
	old        T
}

func (t *boundExternalTreeItem[T]) setIfChanged(val T) error {
	t.lock.Lock()
	if t.comparator(val, t.old) {
		t.lock.Unlock()
		return nil
	}
	(*t.val)[t.id] = val
	t.old = val
	t.lock.Unlock()

	t.trigger()
	return nil
}

func toString[T any](v bindableItem[T], formatter func(T) (string, error), comparator func(T, T) bool, parser func(string) (T, error)) *toStringFrom[T] {
	str := &toStringFrom[T]{from: v, formatter: formatter, comparator: comparator, parser: parser}
	v.AddListener(str)
	return str
}

func toStringComparable[T bool | float64 | int](v bindableItem[T], formatter func(T) (string, error), parser func(string) (T, error)) *toStringFrom[T] {
	return toString(v, formatter, func(t1, t2 T) bool { return t1 == t2 }, parser)
}

func toStringWithFormat[T any](v bindableItem[T], format, defaultFormat string, formatter func(T) (string, error), comparator func(T, T) bool, parser func(string) (T, error)) String {
	str := toString(v, formatter, comparator, parser)
	if format != defaultFormat { // Same as not using custom formatting.
		str.format = format
	}

	return str
}

func toStringWithFormatComparable[T bool | float64 | int](v bindableItem[T], format, defaultFormat string, formatter func(T) (string, error), parser func(string) (T, error)) String {
	return toStringWithFormat(v, format, defaultFormat, formatter, func(t1, t2 T) bool { return t1 == t2 }, parser)
}

type toStringFrom[T any] struct {
	base

	format string

	formatter  func(T) (string, error)
	comparator func(T, T) bool
	parser     func(string) (T, error)

	from bindableItem[T]
}

func (s *toStringFrom[T]) Get() (string, error) {
	val, err := s.from.Get()
	if err != nil {
		return "", err
	}

	if s.format != "" {
		return fmt.Sprintf(s.format, val), nil
	}

	return s.formatter(val)
}

func (s *toStringFrom[T]) Set(str string) error {
	var val T
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
		new, err := s.parser(str)
		if err != nil {
			return err
		}
		val = new
	}

	old, err := s.from.Get()
	if err != nil {
		return err
	}
	if s.comparator(val, old) {
		return nil
	}
	if err = s.from.Set(val); err != nil {
		return err
	}

	queueItem(s.DataChanged)
	return nil
}

func (s *toStringFrom[T]) DataChanged() {
	s.trigger()
}

type fromStringTo[T any] struct {
	base

	format    string
	formatter func(string) (T, error)
	parser    func(T) (string, error)

	from String
}

func (s *fromStringTo[T]) Get() (T, error) {
	str, err := s.from.Get()
	if str == "" || err != nil {
		return *new(T), err
	}

	if s.formatter != nil {
		return s.formatter(str)
	}

	var val T
	if s.format != "" {
		n, err := fmt.Sscanf(str, s.format+" ", &val) // " " denotes match to end of string
		if err != nil {
			return *new(T), err
		}
		if n != 1 {
			return *new(T), errParseFailed
		}
	} else {
		formatted, err := s.formatter(str)
		if err != nil {
			return *new(T), err
		}
		val = formatted
	}

	return val, nil
}

func (s *fromStringTo[T]) Set(val T) error {
	var str string
	if s.parser != nil {
		parsed, err := s.parser(val)
		if err != nil {
			return err
		}
		str = parsed
	} else {
		if s.format != "" {
			str = fmt.Sprintf(s.format, val)
		} else {
			str, _ = s.parser(val)
		}
	}

	old, err := s.from.Get()
	if str == old {
		return err
	}

	err = s.from.Set(str)
	if err != nil {
		return err
	}

	queueItem(s.DataChanged)
	return nil
}

func (s *fromStringTo[T]) DataChanged() {
	s.trigger()
}

type toInt[T float64] struct {
	base

	formatter func(int) (T, error)
	parser    func(T) (int, error)

	from bindableItem[T]
}

func (s *toInt[T]) Get() (int, error) {
	val, err := s.from.Get()
	if err != nil {
		return 0, err
	}
	return s.parser(val)
}

func (s *toInt[T]) Set(v int) error {
	val, err := s.formatter(v)
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

	queueItem(s.DataChanged)
	return nil
}

func (s *toInt[T]) DataChanged() {
	s.trigger()
}

type fromIntTo[T float64] struct {
	base

	formatter func(int) (T, error)
	parser    func(T) (int, error)
	from      bindableItem[int]
}

func (s *fromIntTo[T]) Get() (T, error) {
	val, err := s.from.Get()
	if err != nil {
		return 0.0, err
	}
	return s.formatter(val)
}

func (s *fromIntTo[T]) Set(val T) error {
	i, err := s.parser(val)
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
	err = s.from.Set(i)
	if err != nil {
		return err
	}

	queueItem(s.DataChanged)
	return nil
}

func (s *fromIntTo[T]) DataChanged() {
	s.trigger()
}
