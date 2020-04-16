// auto-generated
// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //

package binding

import (
	"net/url"
	"sync"

	"fyne.io/fyne"
)

// Bool defines a data binding for a bool.
type Bool interface {
	Binding
	Get() bool
	Set(bool)
	AddBoolListener(func(bool)) *NotifyFunction
}

// baseBool implements a data binding for a bool.
type baseBool struct {
	Base
	reference *bool
}

// NewBool creates a new binding with the given value.
func NewBool(value bool) Bool {
	return NewBoolRef(&value)
}

// NewBoolRef creates a new binding with the given reference.
func NewBoolRef(reference *bool) Bool {
	return &baseBool{reference: reference}
}

// Get returns the value of the bound reference.
func (b *baseBool) Get() bool {
	return *b.reference
}

// Set updates the value of the bound reference.
func (b *baseBool) Set(value bool) {
	if *b.reference == value {
		return
	}
	*b.reference = value
	b.Update()
}

// AddBoolListener adds the given function as a listener to the binding.
// The function is wrapped in the returned NotifyFunction which can be passed to DeleteListener.
func (b *baseBool) AddBoolListener(listener func(bool)) *NotifyFunction {
	n := NewNotifyFunction(func(Binding) {
		listener(*b.reference)
	})
	b.AddListener(n)
	return n
}

// BoolList defines a data binding for a list of bool.
type BoolList interface {
	List
	GetBool(int) bool
	GetRef(int) *bool
	SetBool(int, bool)
	SetRef(int, *bool)
	AddBool(bool)
	AddRef(*bool)
}

// baseBoolList implements a data binding for a list of bool.
type baseBoolList struct {
	Base
	sync.Mutex
	references *[]*bool
	bindings   map[*bool]Bool
}

// NewBoolList creates a new list binding with the given values.
func NewBoolList(values []bool) BoolList {
	var references []*bool
	for i := 0; i < len(values); i++ {
		references = append(references, &values[i])
	}
	return NewBoolListRefs(&references)
}

// NewBoolListRefs creates a new list binding with the given references.
func NewBoolListRefs(references *[]*bool) BoolList {
	return &baseBoolList{
		references: references,
		bindings:   make(map[*bool]Bool),
	}
}

// Length returns the number of elements in the list.
func (b *baseBoolList) Length() int {
	return len(*b.references)
}

// Get returns the binding at the given index.
func (b *baseBoolList) Get(index int) Binding {
	reference := b.GetRef(index)
	if reference == nil {
		return nil
	}
	b.Lock()
	defer b.Unlock()
	binding, ok := b.bindings[reference]
	if !ok {
		binding = NewBoolRef(reference)
		b.bindings[reference] = binding
	}
	return binding
}

// GetBool returns the bool at the given index.
func (b *baseBoolList) GetBool(index int) bool {
	if index < 0 && index >= b.Length() {
		return false
	}
	return *(*b.references)[index]
}

// GetRef returns the reference at the given index.
func (b *baseBoolList) GetRef(index int) *bool {
	if index < 0 && index >= b.Length() {
		return nil
	}
	return (*b.references)[index]
}

// SetBool updates the bool at the given index.
func (b *baseBoolList) SetBool(index int, value bool) {
	if index < 0 && index >= b.Length() {
		return
	}
	if *(*b.references)[index] == value {
		return
	}
	b.Get(index).(Bool).Set(value)
}

// SetRef updates the bool at the given index.
func (b *baseBoolList) SetRef(index int, reference *bool) {
	if index < 0 && index >= b.Length() {
		return
	}
	if (*b.references)[index] == reference {
		return
	}
	(*b.references)[index] = reference
	b.Update()
}

// AddBool appends the bool to the list.
func (b *baseBoolList) AddBool(value bool) {
	b.AddRef(&value)
}

// AddRef appends the bool to the list.
func (b *baseBoolList) AddRef(reference *bool) {
	*b.references = append(*b.references, reference)
	b.Update()
}

// Float64 defines a data binding for a float64.
type Float64 interface {
	Binding
	Get() float64
	Set(float64)
	AddFloat64Listener(func(float64)) *NotifyFunction
}

// baseFloat64 implements a data binding for a float64.
type baseFloat64 struct {
	Base
	reference *float64
}

// NewFloat64 creates a new binding with the given value.
func NewFloat64(value float64) Float64 {
	return NewFloat64Ref(&value)
}

// NewFloat64Ref creates a new binding with the given reference.
func NewFloat64Ref(reference *float64) Float64 {
	return &baseFloat64{reference: reference}
}

// Get returns the value of the bound reference.
func (b *baseFloat64) Get() float64 {
	return *b.reference
}

// Set updates the value of the bound reference.
func (b *baseFloat64) Set(value float64) {
	if *b.reference == value {
		return
	}
	*b.reference = value
	b.Update()
}

// AddFloat64Listener adds the given function as a listener to the binding.
// The function is wrapped in the returned NotifyFunction which can be passed to DeleteListener.
func (b *baseFloat64) AddFloat64Listener(listener func(float64)) *NotifyFunction {
	n := NewNotifyFunction(func(Binding) {
		listener(*b.reference)
	})
	b.AddListener(n)
	return n
}

// Float64List defines a data binding for a list of float64.
type Float64List interface {
	List
	GetFloat64(int) float64
	GetRef(int) *float64
	SetFloat64(int, float64)
	SetRef(int, *float64)
	AddFloat64(float64)
	AddRef(*float64)
}

// baseFloat64List implements a data binding for a list of float64.
type baseFloat64List struct {
	Base
	sync.Mutex
	references *[]*float64
	bindings   map[*float64]Float64
}

// NewFloat64List creates a new list binding with the given values.
func NewFloat64List(values []float64) Float64List {
	var references []*float64
	for i := 0; i < len(values); i++ {
		references = append(references, &values[i])
	}
	return NewFloat64ListRefs(&references)
}

// NewFloat64ListRefs creates a new list binding with the given references.
func NewFloat64ListRefs(references *[]*float64) Float64List {
	return &baseFloat64List{
		references: references,
		bindings:   make(map[*float64]Float64),
	}
}

// Length returns the number of elements in the list.
func (b *baseFloat64List) Length() int {
	return len(*b.references)
}

// Get returns the binding at the given index.
func (b *baseFloat64List) Get(index int) Binding {
	reference := b.GetRef(index)
	if reference == nil {
		return nil
	}
	b.Lock()
	defer b.Unlock()
	binding, ok := b.bindings[reference]
	if !ok {
		binding = NewFloat64Ref(reference)
		b.bindings[reference] = binding
	}
	return binding
}

// GetFloat64 returns the float64 at the given index.
func (b *baseFloat64List) GetFloat64(index int) float64 {
	if index < 0 && index >= b.Length() {
		return 0.0
	}
	return *(*b.references)[index]
}

// GetRef returns the reference at the given index.
func (b *baseFloat64List) GetRef(index int) *float64 {
	if index < 0 && index >= b.Length() {
		return nil
	}
	return (*b.references)[index]
}

// SetFloat64 updates the float64 at the given index.
func (b *baseFloat64List) SetFloat64(index int, value float64) {
	if index < 0 && index >= b.Length() {
		return
	}
	if *(*b.references)[index] == value {
		return
	}
	b.Get(index).(Float64).Set(value)
}

// SetRef updates the float64 at the given index.
func (b *baseFloat64List) SetRef(index int, reference *float64) {
	if index < 0 && index >= b.Length() {
		return
	}
	if (*b.references)[index] == reference {
		return
	}
	(*b.references)[index] = reference
	b.Update()
}

// AddFloat64 appends the float64 to the list.
func (b *baseFloat64List) AddFloat64(value float64) {
	b.AddRef(&value)
}

// AddRef appends the float64 to the list.
func (b *baseFloat64List) AddRef(reference *float64) {
	*b.references = append(*b.references, reference)
	b.Update()
}

// Int defines a data binding for a int.
type Int interface {
	Binding
	Get() int
	Set(int)
	AddIntListener(func(int)) *NotifyFunction
}

// baseInt implements a data binding for a int.
type baseInt struct {
	Base
	reference *int
}

// NewInt creates a new binding with the given value.
func NewInt(value int) Int {
	return NewIntRef(&value)
}

// NewIntRef creates a new binding with the given reference.
func NewIntRef(reference *int) Int {
	return &baseInt{reference: reference}
}

// Get returns the value of the bound reference.
func (b *baseInt) Get() int {
	return *b.reference
}

// Set updates the value of the bound reference.
func (b *baseInt) Set(value int) {
	if *b.reference == value {
		return
	}
	*b.reference = value
	b.Update()
}

// AddIntListener adds the given function as a listener to the binding.
// The function is wrapped in the returned NotifyFunction which can be passed to DeleteListener.
func (b *baseInt) AddIntListener(listener func(int)) *NotifyFunction {
	n := NewNotifyFunction(func(Binding) {
		listener(*b.reference)
	})
	b.AddListener(n)
	return n
}

// IntList defines a data binding for a list of int.
type IntList interface {
	List
	GetInt(int) int
	GetRef(int) *int
	SetInt(int, int)
	SetRef(int, *int)
	AddInt(int)
	AddRef(*int)
}

// baseIntList implements a data binding for a list of int.
type baseIntList struct {
	Base
	sync.Mutex
	references *[]*int
	bindings   map[*int]Int
}

// NewIntList creates a new list binding with the given values.
func NewIntList(values []int) IntList {
	var references []*int
	for i := 0; i < len(values); i++ {
		references = append(references, &values[i])
	}
	return NewIntListRefs(&references)
}

// NewIntListRefs creates a new list binding with the given references.
func NewIntListRefs(references *[]*int) IntList {
	return &baseIntList{
		references: references,
		bindings:   make(map[*int]Int),
	}
}

// Length returns the number of elements in the list.
func (b *baseIntList) Length() int {
	return len(*b.references)
}

// Get returns the binding at the given index.
func (b *baseIntList) Get(index int) Binding {
	reference := b.GetRef(index)
	if reference == nil {
		return nil
	}
	b.Lock()
	defer b.Unlock()
	binding, ok := b.bindings[reference]
	if !ok {
		binding = NewIntRef(reference)
		b.bindings[reference] = binding
	}
	return binding
}

// GetInt returns the int at the given index.
func (b *baseIntList) GetInt(index int) int {
	if index < 0 && index >= b.Length() {
		return 0
	}
	return *(*b.references)[index]
}

// GetRef returns the reference at the given index.
func (b *baseIntList) GetRef(index int) *int {
	if index < 0 && index >= b.Length() {
		return nil
	}
	return (*b.references)[index]
}

// SetInt updates the int at the given index.
func (b *baseIntList) SetInt(index int, value int) {
	if index < 0 && index >= b.Length() {
		return
	}
	if *(*b.references)[index] == value {
		return
	}
	b.Get(index).(Int).Set(value)
}

// SetRef updates the int at the given index.
func (b *baseIntList) SetRef(index int, reference *int) {
	if index < 0 && index >= b.Length() {
		return
	}
	if (*b.references)[index] == reference {
		return
	}
	(*b.references)[index] = reference
	b.Update()
}

// AddInt appends the int to the list.
func (b *baseIntList) AddInt(value int) {
	b.AddRef(&value)
}

// AddRef appends the int to the list.
func (b *baseIntList) AddRef(reference *int) {
	*b.references = append(*b.references, reference)
	b.Update()
}

// Int64 defines a data binding for a int64.
type Int64 interface {
	Binding
	Get() int64
	Set(int64)
	AddInt64Listener(func(int64)) *NotifyFunction
}

// baseInt64 implements a data binding for a int64.
type baseInt64 struct {
	Base
	reference *int64
}

// NewInt64 creates a new binding with the given value.
func NewInt64(value int64) Int64 {
	return NewInt64Ref(&value)
}

// NewInt64Ref creates a new binding with the given reference.
func NewInt64Ref(reference *int64) Int64 {
	return &baseInt64{reference: reference}
}

// Get returns the value of the bound reference.
func (b *baseInt64) Get() int64 {
	return *b.reference
}

// Set updates the value of the bound reference.
func (b *baseInt64) Set(value int64) {
	if *b.reference == value {
		return
	}
	*b.reference = value
	b.Update()
}

// AddInt64Listener adds the given function as a listener to the binding.
// The function is wrapped in the returned NotifyFunction which can be passed to DeleteListener.
func (b *baseInt64) AddInt64Listener(listener func(int64)) *NotifyFunction {
	n := NewNotifyFunction(func(Binding) {
		listener(*b.reference)
	})
	b.AddListener(n)
	return n
}

// Int64List defines a data binding for a list of int64.
type Int64List interface {
	List
	GetInt64(int) int64
	GetRef(int) *int64
	SetInt64(int, int64)
	SetRef(int, *int64)
	AddInt64(int64)
	AddRef(*int64)
}

// baseInt64List implements a data binding for a list of int64.
type baseInt64List struct {
	Base
	sync.Mutex
	references *[]*int64
	bindings   map[*int64]Int64
}

// NewInt64List creates a new list binding with the given values.
func NewInt64List(values []int64) Int64List {
	var references []*int64
	for i := 0; i < len(values); i++ {
		references = append(references, &values[i])
	}
	return NewInt64ListRefs(&references)
}

// NewInt64ListRefs creates a new list binding with the given references.
func NewInt64ListRefs(references *[]*int64) Int64List {
	return &baseInt64List{
		references: references,
		bindings:   make(map[*int64]Int64),
	}
}

// Length returns the number of elements in the list.
func (b *baseInt64List) Length() int {
	return len(*b.references)
}

// Get returns the binding at the given index.
func (b *baseInt64List) Get(index int) Binding {
	reference := b.GetRef(index)
	if reference == nil {
		return nil
	}
	b.Lock()
	defer b.Unlock()
	binding, ok := b.bindings[reference]
	if !ok {
		binding = NewInt64Ref(reference)
		b.bindings[reference] = binding
	}
	return binding
}

// GetInt64 returns the int64 at the given index.
func (b *baseInt64List) GetInt64(index int) int64 {
	if index < 0 && index >= b.Length() {
		return 0
	}
	return *(*b.references)[index]
}

// GetRef returns the reference at the given index.
func (b *baseInt64List) GetRef(index int) *int64 {
	if index < 0 && index >= b.Length() {
		return nil
	}
	return (*b.references)[index]
}

// SetInt64 updates the int64 at the given index.
func (b *baseInt64List) SetInt64(index int, value int64) {
	if index < 0 && index >= b.Length() {
		return
	}
	if *(*b.references)[index] == value {
		return
	}
	b.Get(index).(Int64).Set(value)
}

// SetRef updates the int64 at the given index.
func (b *baseInt64List) SetRef(index int, reference *int64) {
	if index < 0 && index >= b.Length() {
		return
	}
	if (*b.references)[index] == reference {
		return
	}
	(*b.references)[index] = reference
	b.Update()
}

// AddInt64 appends the int64 to the list.
func (b *baseInt64List) AddInt64(value int64) {
	b.AddRef(&value)
}

// AddRef appends the int64 to the list.
func (b *baseInt64List) AddRef(reference *int64) {
	*b.references = append(*b.references, reference)
	b.Update()
}

// Resource defines a data binding for a fyne.Resource.
type Resource interface {
	Binding
	Get() fyne.Resource
	Set(fyne.Resource)
	AddResourceListener(func(fyne.Resource)) *NotifyFunction
}

// baseResource implements a data binding for a fyne.Resource.
type baseResource struct {
	Base
	reference *fyne.Resource
}

// NewResource creates a new binding with the given value.
func NewResource(value fyne.Resource) Resource {
	return NewResourceRef(&value)
}

// NewResourceRef creates a new binding with the given reference.
func NewResourceRef(reference *fyne.Resource) Resource {
	return &baseResource{reference: reference}
}

// Get returns the value of the bound reference.
func (b *baseResource) Get() fyne.Resource {
	return *b.reference
}

// Set updates the value of the bound reference.
func (b *baseResource) Set(value fyne.Resource) {
	if *b.reference == value {
		return
	}
	*b.reference = value
	b.Update()
}

// AddResourceListener adds the given function as a listener to the binding.
// The function is wrapped in the returned NotifyFunction which can be passed to DeleteListener.
func (b *baseResource) AddResourceListener(listener func(fyne.Resource)) *NotifyFunction {
	n := NewNotifyFunction(func(Binding) {
		listener(*b.reference)
	})
	b.AddListener(n)
	return n
}

// ResourceList defines a data binding for a list of fyne.Resource.
type ResourceList interface {
	List
	GetResource(int) fyne.Resource
	GetRef(int) *fyne.Resource
	SetResource(int, fyne.Resource)
	SetRef(int, *fyne.Resource)
	AddResource(fyne.Resource)
	AddRef(*fyne.Resource)
}

// baseResourceList implements a data binding for a list of fyne.Resource.
type baseResourceList struct {
	Base
	sync.Mutex
	references *[]*fyne.Resource
	bindings   map[*fyne.Resource]Resource
}

// NewResourceList creates a new list binding with the given values.
func NewResourceList(values []fyne.Resource) ResourceList {
	var references []*fyne.Resource
	for i := 0; i < len(values); i++ {
		references = append(references, &values[i])
	}
	return NewResourceListRefs(&references)
}

// NewResourceListRefs creates a new list binding with the given references.
func NewResourceListRefs(references *[]*fyne.Resource) ResourceList {
	return &baseResourceList{
		references: references,
		bindings:   make(map[*fyne.Resource]Resource),
	}
}

// Length returns the number of elements in the list.
func (b *baseResourceList) Length() int {
	return len(*b.references)
}

// Get returns the binding at the given index.
func (b *baseResourceList) Get(index int) Binding {
	reference := b.GetRef(index)
	if reference == nil {
		return nil
	}
	b.Lock()
	defer b.Unlock()
	binding, ok := b.bindings[reference]
	if !ok {
		binding = NewResourceRef(reference)
		b.bindings[reference] = binding
	}
	return binding
}

// GetResource returns the fyne.Resource at the given index.
func (b *baseResourceList) GetResource(index int) fyne.Resource {
	if index < 0 && index >= b.Length() {
		return nil
	}
	return *(*b.references)[index]
}

// GetRef returns the reference at the given index.
func (b *baseResourceList) GetRef(index int) *fyne.Resource {
	if index < 0 && index >= b.Length() {
		return nil
	}
	return (*b.references)[index]
}

// SetResource updates the fyne.Resource at the given index.
func (b *baseResourceList) SetResource(index int, value fyne.Resource) {
	if index < 0 && index >= b.Length() {
		return
	}
	if *(*b.references)[index] == value {
		return
	}
	b.Get(index).(Resource).Set(value)
}

// SetRef updates the fyne.Resource at the given index.
func (b *baseResourceList) SetRef(index int, reference *fyne.Resource) {
	if index < 0 && index >= b.Length() {
		return
	}
	if (*b.references)[index] == reference {
		return
	}
	(*b.references)[index] = reference
	b.Update()
}

// AddResource appends the fyne.Resource to the list.
func (b *baseResourceList) AddResource(value fyne.Resource) {
	b.AddRef(&value)
}

// AddRef appends the fyne.Resource to the list.
func (b *baseResourceList) AddRef(reference *fyne.Resource) {
	*b.references = append(*b.references, reference)
	b.Update()
}

// Rune defines a data binding for a rune.
type Rune interface {
	Binding
	Get() rune
	Set(rune)
	AddRuneListener(func(rune)) *NotifyFunction
}

// baseRune implements a data binding for a rune.
type baseRune struct {
	Base
	reference *rune
}

// NewRune creates a new binding with the given value.
func NewRune(value rune) Rune {
	return NewRuneRef(&value)
}

// NewRuneRef creates a new binding with the given reference.
func NewRuneRef(reference *rune) Rune {
	return &baseRune{reference: reference}
}

// Get returns the value of the bound reference.
func (b *baseRune) Get() rune {
	return *b.reference
}

// Set updates the value of the bound reference.
func (b *baseRune) Set(value rune) {
	if *b.reference == value {
		return
	}
	*b.reference = value
	b.Update()
}

// AddRuneListener adds the given function as a listener to the binding.
// The function is wrapped in the returned NotifyFunction which can be passed to DeleteListener.
func (b *baseRune) AddRuneListener(listener func(rune)) *NotifyFunction {
	n := NewNotifyFunction(func(Binding) {
		listener(*b.reference)
	})
	b.AddListener(n)
	return n
}

// RuneList defines a data binding for a list of rune.
type RuneList interface {
	List
	GetRune(int) rune
	GetRef(int) *rune
	SetRune(int, rune)
	SetRef(int, *rune)
	AddRune(rune)
	AddRef(*rune)
}

// baseRuneList implements a data binding for a list of rune.
type baseRuneList struct {
	Base
	sync.Mutex
	references *[]*rune
	bindings   map[*rune]Rune
}

// NewRuneList creates a new list binding with the given values.
func NewRuneList(values []rune) RuneList {
	var references []*rune
	for i := 0; i < len(values); i++ {
		references = append(references, &values[i])
	}
	return NewRuneListRefs(&references)
}

// NewRuneListRefs creates a new list binding with the given references.
func NewRuneListRefs(references *[]*rune) RuneList {
	return &baseRuneList{
		references: references,
		bindings:   make(map[*rune]Rune),
	}
}

// Length returns the number of elements in the list.
func (b *baseRuneList) Length() int {
	return len(*b.references)
}

// Get returns the binding at the given index.
func (b *baseRuneList) Get(index int) Binding {
	reference := b.GetRef(index)
	if reference == nil {
		return nil
	}
	b.Lock()
	defer b.Unlock()
	binding, ok := b.bindings[reference]
	if !ok {
		binding = NewRuneRef(reference)
		b.bindings[reference] = binding
	}
	return binding
}

// GetRune returns the rune at the given index.
func (b *baseRuneList) GetRune(index int) rune {
	if index < 0 && index >= b.Length() {
		return 0
	}
	return *(*b.references)[index]
}

// GetRef returns the reference at the given index.
func (b *baseRuneList) GetRef(index int) *rune {
	if index < 0 && index >= b.Length() {
		return nil
	}
	return (*b.references)[index]
}

// SetRune updates the rune at the given index.
func (b *baseRuneList) SetRune(index int, value rune) {
	if index < 0 && index >= b.Length() {
		return
	}
	if *(*b.references)[index] == value {
		return
	}
	b.Get(index).(Rune).Set(value)
}

// SetRef updates the rune at the given index.
func (b *baseRuneList) SetRef(index int, reference *rune) {
	if index < 0 && index >= b.Length() {
		return
	}
	if (*b.references)[index] == reference {
		return
	}
	(*b.references)[index] = reference
	b.Update()
}

// AddRune appends the rune to the list.
func (b *baseRuneList) AddRune(value rune) {
	b.AddRef(&value)
}

// AddRef appends the rune to the list.
func (b *baseRuneList) AddRef(reference *rune) {
	*b.references = append(*b.references, reference)
	b.Update()
}

// String defines a data binding for a string.
type String interface {
	Binding
	Get() string
	Set(string)
	AddStringListener(func(string)) *NotifyFunction
}

// baseString implements a data binding for a string.
type baseString struct {
	Base
	reference *string
}

// NewString creates a new binding with the given value.
func NewString(value string) String {
	return NewStringRef(&value)
}

// NewStringRef creates a new binding with the given reference.
func NewStringRef(reference *string) String {
	return &baseString{reference: reference}
}

// Get returns the value of the bound reference.
func (b *baseString) Get() string {
	return *b.reference
}

// Set updates the value of the bound reference.
func (b *baseString) Set(value string) {
	if *b.reference == value {
		return
	}
	*b.reference = value
	b.Update()
}

// AddStringListener adds the given function as a listener to the binding.
// The function is wrapped in the returned NotifyFunction which can be passed to DeleteListener.
func (b *baseString) AddStringListener(listener func(string)) *NotifyFunction {
	n := NewNotifyFunction(func(Binding) {
		listener(*b.reference)
	})
	b.AddListener(n)
	return n
}

// StringList defines a data binding for a list of string.
type StringList interface {
	List
	GetString(int) string
	GetRef(int) *string
	SetString(int, string)
	SetRef(int, *string)
	AddString(string)
	AddRef(*string)
}

// baseStringList implements a data binding for a list of string.
type baseStringList struct {
	Base
	sync.Mutex
	references *[]*string
	bindings   map[*string]String
}

// NewStringList creates a new list binding with the given values.
func NewStringList(values []string) StringList {
	var references []*string
	for i := 0; i < len(values); i++ {
		references = append(references, &values[i])
	}
	return NewStringListRefs(&references)
}

// NewStringListRefs creates a new list binding with the given references.
func NewStringListRefs(references *[]*string) StringList {
	return &baseStringList{
		references: references,
		bindings:   make(map[*string]String),
	}
}

// Length returns the number of elements in the list.
func (b *baseStringList) Length() int {
	return len(*b.references)
}

// Get returns the binding at the given index.
func (b *baseStringList) Get(index int) Binding {
	reference := b.GetRef(index)
	if reference == nil {
		return nil
	}
	b.Lock()
	defer b.Unlock()
	binding, ok := b.bindings[reference]
	if !ok {
		binding = NewStringRef(reference)
		b.bindings[reference] = binding
	}
	return binding
}

// GetString returns the string at the given index.
func (b *baseStringList) GetString(index int) string {
	if index < 0 && index >= b.Length() {
		return ""
	}
	return *(*b.references)[index]
}

// GetRef returns the reference at the given index.
func (b *baseStringList) GetRef(index int) *string {
	if index < 0 && index >= b.Length() {
		return nil
	}
	return (*b.references)[index]
}

// SetString updates the string at the given index.
func (b *baseStringList) SetString(index int, value string) {
	if index < 0 && index >= b.Length() {
		return
	}
	if *(*b.references)[index] == value {
		return
	}
	b.Get(index).(String).Set(value)
}

// SetRef updates the string at the given index.
func (b *baseStringList) SetRef(index int, reference *string) {
	if index < 0 && index >= b.Length() {
		return
	}
	if (*b.references)[index] == reference {
		return
	}
	(*b.references)[index] = reference
	b.Update()
}

// AddString appends the string to the list.
func (b *baseStringList) AddString(value string) {
	b.AddRef(&value)
}

// AddRef appends the string to the list.
func (b *baseStringList) AddRef(reference *string) {
	*b.references = append(*b.references, reference)
	b.Update()
}

// URL defines a data binding for a *url.URL.
type URL interface {
	Binding
	Get() *url.URL
	Set(*url.URL)
	AddURLListener(func(*url.URL)) *NotifyFunction
}

// baseURL implements a data binding for a *url.URL.
type baseURL struct {
	Base
	reference **url.URL
}

// NewURL creates a new binding with the given value.
func NewURL(value *url.URL) URL {
	return NewURLRef(&value)
}

// NewURLRef creates a new binding with the given reference.
func NewURLRef(reference **url.URL) URL {
	return &baseURL{reference: reference}
}

// Get returns the value of the bound reference.
func (b *baseURL) Get() *url.URL {
	return *b.reference
}

// Set updates the value of the bound reference.
func (b *baseURL) Set(value *url.URL) {
	if *b.reference == value {
		return
	}
	*b.reference = value
	b.Update()
}

// AddURLListener adds the given function as a listener to the binding.
// The function is wrapped in the returned NotifyFunction which can be passed to DeleteListener.
func (b *baseURL) AddURLListener(listener func(*url.URL)) *NotifyFunction {
	n := NewNotifyFunction(func(Binding) {
		listener(*b.reference)
	})
	b.AddListener(n)
	return n
}

// URLList defines a data binding for a list of *url.URL.
type URLList interface {
	List
	GetURL(int) *url.URL
	GetRef(int) **url.URL
	SetURL(int, *url.URL)
	SetRef(int, **url.URL)
	AddURL(*url.URL)
	AddRef(**url.URL)
}

// baseURLList implements a data binding for a list of *url.URL.
type baseURLList struct {
	Base
	sync.Mutex
	references *[]**url.URL
	bindings   map[**url.URL]URL
}

// NewURLList creates a new list binding with the given values.
func NewURLList(values []*url.URL) URLList {
	var references []**url.URL
	for i := 0; i < len(values); i++ {
		references = append(references, &values[i])
	}
	return NewURLListRefs(&references)
}

// NewURLListRefs creates a new list binding with the given references.
func NewURLListRefs(references *[]**url.URL) URLList {
	return &baseURLList{
		references: references,
		bindings:   make(map[**url.URL]URL),
	}
}

// Length returns the number of elements in the list.
func (b *baseURLList) Length() int {
	return len(*b.references)
}

// Get returns the binding at the given index.
func (b *baseURLList) Get(index int) Binding {
	reference := b.GetRef(index)
	if reference == nil {
		return nil
	}
	b.Lock()
	defer b.Unlock()
	binding, ok := b.bindings[reference]
	if !ok {
		binding = NewURLRef(reference)
		b.bindings[reference] = binding
	}
	return binding
}

// GetURL returns the *url.URL at the given index.
func (b *baseURLList) GetURL(index int) *url.URL {
	if index < 0 && index >= b.Length() {
		return nil
	}
	return *(*b.references)[index]
}

// GetRef returns the reference at the given index.
func (b *baseURLList) GetRef(index int) **url.URL {
	if index < 0 && index >= b.Length() {
		return nil
	}
	return (*b.references)[index]
}

// SetURL updates the *url.URL at the given index.
func (b *baseURLList) SetURL(index int, value *url.URL) {
	if index < 0 && index >= b.Length() {
		return
	}
	if *(*b.references)[index] == value {
		return
	}
	b.Get(index).(URL).Set(value)
}

// SetRef updates the *url.URL at the given index.
func (b *baseURLList) SetRef(index int, reference **url.URL) {
	if index < 0 && index >= b.Length() {
		return
	}
	if (*b.references)[index] == reference {
		return
	}
	(*b.references)[index] = reference
	b.Update()
}

// AddURL appends the *url.URL to the list.
func (b *baseURLList) AddURL(value *url.URL) {
	b.AddRef(&value)
}

// AddRef appends the *url.URL to the list.
func (b *baseURLList) AddRef(reference **url.URL) {
	*b.references = append(*b.references, reference)
	b.Update()
}
