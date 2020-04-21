// auto-generated
// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //

package binding

import (
	"fmt"
	"net/url"
	"runtime"
	"sync"

	"fyne.io/fyne"
)

// Bool defines a data binding for a bool.
type Bool interface {
	Binding
	Get() bool
	GetRef() *bool
	Set(bool)
	SetRef(*bool)
	Listen() <-chan bool
}

// baseBool implements a data binding for a bool.
type baseBool struct {
	sync.Mutex
	reference *bool
	channels  []chan bool
	traces    []string
}

// EmptyBool creates a new binding with the empty value.
func EmptyBool() Bool {
	return NewBool(false)
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

// Get returns the bound reference.
func (b *baseBool) GetRef() *bool {
	return b.reference
}

// Set updates the value of the bound reference.
func (b *baseBool) Set(value bool) {
	if *b.reference == value {
		return
	}
	*b.reference = value
	b.Update()
}

// Set updates the bound reference.
func (b *baseBool) SetRef(reference *bool) {
	if b.reference == reference {
		return
	}
	b.reference = reference
	b.Update()
}

// Listen returns a channel through which updates will be published.
func (b *baseBool) Listen() <-chan bool {
	b.Lock()
	defer b.Unlock()
	c := make(chan bool, 16)
	c <- b.Get()
	b.channels = append(b.channels, c)
	_, file, line, ok := runtime.Caller(1)
	if ok {
		b.traces = append(b.traces, fmt.Sprintf("%s#%d", file, line))
	}
	return c
}

// Update notifies all listeners after a change.
func (b *baseBool) Update() {
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

// BoolList defines a data binding for a list of bool.
type BoolList interface {
	List
	GetBinding(int) Bool
	GetBool(int) bool
	GetRef(int) *bool
	SetBinding(int, Bool)
	SetBool(int, bool)
	SetRef(int, *bool)
	AddBinding(Bool)
	AddBool(bool)
	AddRef(*bool)
}

// baseBoolList implements a data binding for a list of bool.
type baseBoolList struct {
	sync.Mutex
	references *[]*bool
	bindings   map[*bool]Bool
	channels   []chan int
	traces     []string
}

// NewBoolList creates a new list binding with the given values.
func NewBoolList(values ...bool) BoolList {
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
	return b.GetBinding(index)
}

// GetBinding returns the Bool at the given index.
func (b *baseBoolList) GetBinding(index int) Bool {
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

// SetBinding updates the Bool at the given index.
func (b *baseBoolList) SetBinding(index int, binding Bool) {
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

// AddBinding appends the Bool to the list.
func (b *baseBoolList) AddBinding(binding Bool) {
	index := b.Length()
	*b.references = append(*b.references, binding.GetRef())
	b.SetBinding(index, binding)
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

// Listen returns a channel through which updates will be published.
func (b *baseBoolList) Listen() <-chan int {
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
func (b *baseBoolList) Update() {
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

// Float64 defines a data binding for a float64.
type Float64 interface {
	Binding
	Get() float64
	GetRef() *float64
	Set(float64)
	SetRef(*float64)
	Listen() <-chan float64
}

// baseFloat64 implements a data binding for a float64.
type baseFloat64 struct {
	sync.Mutex
	reference *float64
	channels  []chan float64
	traces    []string
}

// EmptyFloat64 creates a new binding with the empty value.
func EmptyFloat64() Float64 {
	return NewFloat64(float64(0.0))
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

// Get returns the bound reference.
func (b *baseFloat64) GetRef() *float64 {
	return b.reference
}

// Set updates the value of the bound reference.
func (b *baseFloat64) Set(value float64) {
	if *b.reference == value {
		return
	}
	*b.reference = value
	b.Update()
}

// Set updates the bound reference.
func (b *baseFloat64) SetRef(reference *float64) {
	if b.reference == reference {
		return
	}
	b.reference = reference
	b.Update()
}

// Listen returns a channel through which updates will be published.
func (b *baseFloat64) Listen() <-chan float64 {
	b.Lock()
	defer b.Unlock()
	c := make(chan float64, 16)
	c <- b.Get()
	b.channels = append(b.channels, c)
	_, file, line, ok := runtime.Caller(1)
	if ok {
		b.traces = append(b.traces, fmt.Sprintf("%s#%d", file, line))
	}
	return c
}

// Update notifies all listeners after a change.
func (b *baseFloat64) Update() {
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

// Float64List defines a data binding for a list of float64.
type Float64List interface {
	List
	GetBinding(int) Float64
	GetFloat64(int) float64
	GetRef(int) *float64
	SetBinding(int, Float64)
	SetFloat64(int, float64)
	SetRef(int, *float64)
	AddBinding(Float64)
	AddFloat64(float64)
	AddRef(*float64)
}

// baseFloat64List implements a data binding for a list of float64.
type baseFloat64List struct {
	sync.Mutex
	references *[]*float64
	bindings   map[*float64]Float64
	channels   []chan int
	traces     []string
}

// NewFloat64List creates a new list binding with the given values.
func NewFloat64List(values ...float64) Float64List {
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
	return b.GetBinding(index)
}

// GetBinding returns the Float64 at the given index.
func (b *baseFloat64List) GetBinding(index int) Float64 {
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
		return float64(0.0)
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

// SetBinding updates the Float64 at the given index.
func (b *baseFloat64List) SetBinding(index int, binding Float64) {
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

// AddBinding appends the Float64 to the list.
func (b *baseFloat64List) AddBinding(binding Float64) {
	index := b.Length()
	*b.references = append(*b.references, binding.GetRef())
	b.SetBinding(index, binding)
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

// Listen returns a channel through which updates will be published.
func (b *baseFloat64List) Listen() <-chan int {
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
func (b *baseFloat64List) Update() {
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

// Int defines a data binding for a int.
type Int interface {
	Binding
	Get() int
	GetRef() *int
	Set(int)
	SetRef(*int)
	Listen() <-chan int
}

// baseInt implements a data binding for a int.
type baseInt struct {
	sync.Mutex
	reference *int
	channels  []chan int
	traces    []string
}

// EmptyInt creates a new binding with the empty value.
func EmptyInt() Int {
	return NewInt(int(0))
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

// Get returns the bound reference.
func (b *baseInt) GetRef() *int {
	return b.reference
}

// Set updates the value of the bound reference.
func (b *baseInt) Set(value int) {
	if *b.reference == value {
		return
	}
	*b.reference = value
	b.Update()
}

// Set updates the bound reference.
func (b *baseInt) SetRef(reference *int) {
	if b.reference == reference {
		return
	}
	b.reference = reference
	b.Update()
}

// Listen returns a channel through which updates will be published.
func (b *baseInt) Listen() <-chan int {
	b.Lock()
	defer b.Unlock()
	c := make(chan int, 16)
	c <- b.Get()
	b.channels = append(b.channels, c)
	_, file, line, ok := runtime.Caller(1)
	if ok {
		b.traces = append(b.traces, fmt.Sprintf("%s#%d", file, line))
	}
	return c
}

// Update notifies all listeners after a change.
func (b *baseInt) Update() {
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

// IntList defines a data binding for a list of int.
type IntList interface {
	List
	GetBinding(int) Int
	GetInt(int) int
	GetRef(int) *int
	SetBinding(int, Int)
	SetInt(int, int)
	SetRef(int, *int)
	AddBinding(Int)
	AddInt(int)
	AddRef(*int)
}

// baseIntList implements a data binding for a list of int.
type baseIntList struct {
	sync.Mutex
	references *[]*int
	bindings   map[*int]Int
	channels   []chan int
	traces     []string
}

// NewIntList creates a new list binding with the given values.
func NewIntList(values ...int) IntList {
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
	return b.GetBinding(index)
}

// GetBinding returns the Int at the given index.
func (b *baseIntList) GetBinding(index int) Int {
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
		return int(0)
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

// SetBinding updates the Int at the given index.
func (b *baseIntList) SetBinding(index int, binding Int) {
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

// AddBinding appends the Int to the list.
func (b *baseIntList) AddBinding(binding Int) {
	index := b.Length()
	*b.references = append(*b.references, binding.GetRef())
	b.SetBinding(index, binding)
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

// Listen returns a channel through which updates will be published.
func (b *baseIntList) Listen() <-chan int {
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
func (b *baseIntList) Update() {
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

// Int64 defines a data binding for a int64.
type Int64 interface {
	Binding
	Get() int64
	GetRef() *int64
	Set(int64)
	SetRef(*int64)
	Listen() <-chan int64
}

// baseInt64 implements a data binding for a int64.
type baseInt64 struct {
	sync.Mutex
	reference *int64
	channels  []chan int64
	traces    []string
}

// EmptyInt64 creates a new binding with the empty value.
func EmptyInt64() Int64 {
	return NewInt64(int64(0))
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

// Get returns the bound reference.
func (b *baseInt64) GetRef() *int64 {
	return b.reference
}

// Set updates the value of the bound reference.
func (b *baseInt64) Set(value int64) {
	if *b.reference == value {
		return
	}
	*b.reference = value
	b.Update()
}

// Set updates the bound reference.
func (b *baseInt64) SetRef(reference *int64) {
	if b.reference == reference {
		return
	}
	b.reference = reference
	b.Update()
}

// Listen returns a channel through which updates will be published.
func (b *baseInt64) Listen() <-chan int64 {
	b.Lock()
	defer b.Unlock()
	c := make(chan int64, 16)
	c <- b.Get()
	b.channels = append(b.channels, c)
	_, file, line, ok := runtime.Caller(1)
	if ok {
		b.traces = append(b.traces, fmt.Sprintf("%s#%d", file, line))
	}
	return c
}

// Update notifies all listeners after a change.
func (b *baseInt64) Update() {
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

// Int64List defines a data binding for a list of int64.
type Int64List interface {
	List
	GetBinding(int) Int64
	GetInt64(int) int64
	GetRef(int) *int64
	SetBinding(int, Int64)
	SetInt64(int, int64)
	SetRef(int, *int64)
	AddBinding(Int64)
	AddInt64(int64)
	AddRef(*int64)
}

// baseInt64List implements a data binding for a list of int64.
type baseInt64List struct {
	sync.Mutex
	references *[]*int64
	bindings   map[*int64]Int64
	channels   []chan int
	traces     []string
}

// NewInt64List creates a new list binding with the given values.
func NewInt64List(values ...int64) Int64List {
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
	return b.GetBinding(index)
}

// GetBinding returns the Int64 at the given index.
func (b *baseInt64List) GetBinding(index int) Int64 {
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
		return int64(0)
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

// SetBinding updates the Int64 at the given index.
func (b *baseInt64List) SetBinding(index int, binding Int64) {
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

// AddBinding appends the Int64 to the list.
func (b *baseInt64List) AddBinding(binding Int64) {
	index := b.Length()
	*b.references = append(*b.references, binding.GetRef())
	b.SetBinding(index, binding)
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

// Listen returns a channel through which updates will be published.
func (b *baseInt64List) Listen() <-chan int {
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
func (b *baseInt64List) Update() {
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

// Position defines a data binding for a fyne.Position.
type Position interface {
	Binding
	Get() fyne.Position
	GetRef() *fyne.Position
	Set(fyne.Position)
	SetRef(*fyne.Position)
	Listen() <-chan fyne.Position
}

// basePosition implements a data binding for a fyne.Position.
type basePosition struct {
	sync.Mutex
	reference *fyne.Position
	channels  []chan fyne.Position
	traces    []string
}

// EmptyPosition creates a new binding with the empty value.
func EmptyPosition() Position {
	return NewPosition(fyne.Position{})
}

// NewPosition creates a new binding with the given value.
func NewPosition(value fyne.Position) Position {
	return NewPositionRef(&value)
}

// NewPositionRef creates a new binding with the given reference.
func NewPositionRef(reference *fyne.Position) Position {
	return &basePosition{reference: reference}
}

// Get returns the value of the bound reference.
func (b *basePosition) Get() fyne.Position {
	return *b.reference
}

// Get returns the bound reference.
func (b *basePosition) GetRef() *fyne.Position {
	return b.reference
}

// Set updates the value of the bound reference.
func (b *basePosition) Set(value fyne.Position) {
	if *b.reference == value {
		return
	}
	*b.reference = value
	b.Update()
}

// Set updates the bound reference.
func (b *basePosition) SetRef(reference *fyne.Position) {
	if b.reference == reference {
		return
	}
	b.reference = reference
	b.Update()
}

// Listen returns a channel through which updates will be published.
func (b *basePosition) Listen() <-chan fyne.Position {
	b.Lock()
	defer b.Unlock()
	c := make(chan fyne.Position, 16)
	c <- b.Get()
	b.channels = append(b.channels, c)
	_, file, line, ok := runtime.Caller(1)
	if ok {
		b.traces = append(b.traces, fmt.Sprintf("%s#%d", file, line))
	}
	return c
}

// Update notifies all listeners after a change.
func (b *basePosition) Update() {
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

// PositionList defines a data binding for a list of fyne.Position.
type PositionList interface {
	List
	GetBinding(int) Position
	GetPosition(int) fyne.Position
	GetRef(int) *fyne.Position
	SetBinding(int, Position)
	SetPosition(int, fyne.Position)
	SetRef(int, *fyne.Position)
	AddBinding(Position)
	AddPosition(fyne.Position)
	AddRef(*fyne.Position)
}

// basePositionList implements a data binding for a list of fyne.Position.
type basePositionList struct {
	sync.Mutex
	references *[]*fyne.Position
	bindings   map[*fyne.Position]Position
	channels   []chan int
	traces     []string
}

// NewPositionList creates a new list binding with the given values.
func NewPositionList(values ...fyne.Position) PositionList {
	var references []*fyne.Position
	for i := 0; i < len(values); i++ {
		references = append(references, &values[i])
	}
	return NewPositionListRefs(&references)
}

// NewPositionListRefs creates a new list binding with the given references.
func NewPositionListRefs(references *[]*fyne.Position) PositionList {
	return &basePositionList{
		references: references,
		bindings:   make(map[*fyne.Position]Position),
	}
}

// Length returns the number of elements in the list.
func (b *basePositionList) Length() int {
	return len(*b.references)
}

// Get returns the binding at the given index.
func (b *basePositionList) Get(index int) Binding {
	return b.GetBinding(index)
}

// GetBinding returns the Position at the given index.
func (b *basePositionList) GetBinding(index int) Position {
	reference := b.GetRef(index)
	if reference == nil {
		return nil
	}
	b.Lock()
	defer b.Unlock()
	binding, ok := b.bindings[reference]
	if !ok {
		binding = NewPositionRef(reference)
		b.bindings[reference] = binding
	}
	return binding
}

// GetPosition returns the fyne.Position at the given index.
func (b *basePositionList) GetPosition(index int) fyne.Position {
	if index < 0 && index >= b.Length() {
		return fyne.Position{}
	}
	return *(*b.references)[index]
}

// GetRef returns the reference at the given index.
func (b *basePositionList) GetRef(index int) *fyne.Position {
	if index < 0 && index >= b.Length() {
		return nil
	}
	return (*b.references)[index]
}

// SetBinding updates the Position at the given index.
func (b *basePositionList) SetBinding(index int, binding Position) {
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

// SetPosition updates the fyne.Position at the given index.
func (b *basePositionList) SetPosition(index int, value fyne.Position) {
	if index < 0 && index >= b.Length() {
		return
	}
	if *(*b.references)[index] == value {
		return
	}
	b.Get(index).(Position).Set(value)
}

// SetRef updates the fyne.Position at the given index.
func (b *basePositionList) SetRef(index int, reference *fyne.Position) {
	if index < 0 && index >= b.Length() {
		return
	}
	if (*b.references)[index] == reference {
		return
	}
	(*b.references)[index] = reference
	b.Update()
}

// AddBinding appends the Position to the list.
func (b *basePositionList) AddBinding(binding Position) {
	index := b.Length()
	*b.references = append(*b.references, binding.GetRef())
	b.SetBinding(index, binding)
}

// AddPosition appends the fyne.Position to the list.
func (b *basePositionList) AddPosition(value fyne.Position) {
	b.AddRef(&value)
}

// AddRef appends the fyne.Position to the list.
func (b *basePositionList) AddRef(reference *fyne.Position) {
	*b.references = append(*b.references, reference)
	b.Update()
}

// Listen returns a channel through which updates will be published.
func (b *basePositionList) Listen() <-chan int {
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
func (b *basePositionList) Update() {
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

// Resource defines a data binding for a fyne.Resource.
type Resource interface {
	Binding
	Get() fyne.Resource
	GetRef() *fyne.Resource
	Set(fyne.Resource)
	SetRef(*fyne.Resource)
	Listen() <-chan fyne.Resource
}

// baseResource implements a data binding for a fyne.Resource.
type baseResource struct {
	sync.Mutex
	reference *fyne.Resource
	channels  []chan fyne.Resource
	traces    []string
}

// EmptyResource creates a new binding with the empty value.
func EmptyResource() Resource {
	return NewResource(nil)
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

// Get returns the bound reference.
func (b *baseResource) GetRef() *fyne.Resource {
	return b.reference
}

// Set updates the value of the bound reference.
func (b *baseResource) Set(value fyne.Resource) {
	if *b.reference == value {
		return
	}
	*b.reference = value
	b.Update()
}

// Set updates the bound reference.
func (b *baseResource) SetRef(reference *fyne.Resource) {
	if b.reference == reference {
		return
	}
	b.reference = reference
	b.Update()
}

// Listen returns a channel through which updates will be published.
func (b *baseResource) Listen() <-chan fyne.Resource {
	b.Lock()
	defer b.Unlock()
	c := make(chan fyne.Resource, 16)
	c <- b.Get()
	b.channels = append(b.channels, c)
	_, file, line, ok := runtime.Caller(1)
	if ok {
		b.traces = append(b.traces, fmt.Sprintf("%s#%d", file, line))
	}
	return c
}

// Update notifies all listeners after a change.
func (b *baseResource) Update() {
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

// ResourceList defines a data binding for a list of fyne.Resource.
type ResourceList interface {
	List
	GetBinding(int) Resource
	GetResource(int) fyne.Resource
	GetRef(int) *fyne.Resource
	SetBinding(int, Resource)
	SetResource(int, fyne.Resource)
	SetRef(int, *fyne.Resource)
	AddBinding(Resource)
	AddResource(fyne.Resource)
	AddRef(*fyne.Resource)
}

// baseResourceList implements a data binding for a list of fyne.Resource.
type baseResourceList struct {
	sync.Mutex
	references *[]*fyne.Resource
	bindings   map[*fyne.Resource]Resource
	channels   []chan int
	traces     []string
}

// NewResourceList creates a new list binding with the given values.
func NewResourceList(values ...fyne.Resource) ResourceList {
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
	return b.GetBinding(index)
}

// GetBinding returns the Resource at the given index.
func (b *baseResourceList) GetBinding(index int) Resource {
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

// SetBinding updates the Resource at the given index.
func (b *baseResourceList) SetBinding(index int, binding Resource) {
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

// AddBinding appends the Resource to the list.
func (b *baseResourceList) AddBinding(binding Resource) {
	index := b.Length()
	*b.references = append(*b.references, binding.GetRef())
	b.SetBinding(index, binding)
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

// Listen returns a channel through which updates will be published.
func (b *baseResourceList) Listen() <-chan int {
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
func (b *baseResourceList) Update() {
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

// Rune defines a data binding for a rune.
type Rune interface {
	Binding
	Get() rune
	GetRef() *rune
	Set(rune)
	SetRef(*rune)
	Listen() <-chan rune
}

// baseRune implements a data binding for a rune.
type baseRune struct {
	sync.Mutex
	reference *rune
	channels  []chan rune
	traces    []string
}

// EmptyRune creates a new binding with the empty value.
func EmptyRune() Rune {
	return NewRune(rune(0))
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

// Get returns the bound reference.
func (b *baseRune) GetRef() *rune {
	return b.reference
}

// Set updates the value of the bound reference.
func (b *baseRune) Set(value rune) {
	if *b.reference == value {
		return
	}
	*b.reference = value
	b.Update()
}

// Set updates the bound reference.
func (b *baseRune) SetRef(reference *rune) {
	if b.reference == reference {
		return
	}
	b.reference = reference
	b.Update()
}

// Listen returns a channel through which updates will be published.
func (b *baseRune) Listen() <-chan rune {
	b.Lock()
	defer b.Unlock()
	c := make(chan rune, 16)
	c <- b.Get()
	b.channels = append(b.channels, c)
	_, file, line, ok := runtime.Caller(1)
	if ok {
		b.traces = append(b.traces, fmt.Sprintf("%s#%d", file, line))
	}
	return c
}

// Update notifies all listeners after a change.
func (b *baseRune) Update() {
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

// RuneList defines a data binding for a list of rune.
type RuneList interface {
	List
	GetBinding(int) Rune
	GetRune(int) rune
	GetRef(int) *rune
	SetBinding(int, Rune)
	SetRune(int, rune)
	SetRef(int, *rune)
	AddBinding(Rune)
	AddRune(rune)
	AddRef(*rune)
}

// baseRuneList implements a data binding for a list of rune.
type baseRuneList struct {
	sync.Mutex
	references *[]*rune
	bindings   map[*rune]Rune
	channels   []chan int
	traces     []string
}

// NewRuneList creates a new list binding with the given values.
func NewRuneList(values ...rune) RuneList {
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
	return b.GetBinding(index)
}

// GetBinding returns the Rune at the given index.
func (b *baseRuneList) GetBinding(index int) Rune {
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
		return rune(0)
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

// SetBinding updates the Rune at the given index.
func (b *baseRuneList) SetBinding(index int, binding Rune) {
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

// AddBinding appends the Rune to the list.
func (b *baseRuneList) AddBinding(binding Rune) {
	index := b.Length()
	*b.references = append(*b.references, binding.GetRef())
	b.SetBinding(index, binding)
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

// Listen returns a channel through which updates will be published.
func (b *baseRuneList) Listen() <-chan int {
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
func (b *baseRuneList) Update() {
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

// Size defines a data binding for a fyne.Size.
type Size interface {
	Binding
	Get() fyne.Size
	GetRef() *fyne.Size
	Set(fyne.Size)
	SetRef(*fyne.Size)
	Listen() <-chan fyne.Size
}

// baseSize implements a data binding for a fyne.Size.
type baseSize struct {
	sync.Mutex
	reference *fyne.Size
	channels  []chan fyne.Size
	traces    []string
}

// EmptySize creates a new binding with the empty value.
func EmptySize() Size {
	return NewSize(fyne.Size{})
}

// NewSize creates a new binding with the given value.
func NewSize(value fyne.Size) Size {
	return NewSizeRef(&value)
}

// NewSizeRef creates a new binding with the given reference.
func NewSizeRef(reference *fyne.Size) Size {
	return &baseSize{reference: reference}
}

// Get returns the value of the bound reference.
func (b *baseSize) Get() fyne.Size {
	return *b.reference
}

// Get returns the bound reference.
func (b *baseSize) GetRef() *fyne.Size {
	return b.reference
}

// Set updates the value of the bound reference.
func (b *baseSize) Set(value fyne.Size) {
	if *b.reference == value {
		return
	}
	*b.reference = value
	b.Update()
}

// Set updates the bound reference.
func (b *baseSize) SetRef(reference *fyne.Size) {
	if b.reference == reference {
		return
	}
	b.reference = reference
	b.Update()
}

// Listen returns a channel through which updates will be published.
func (b *baseSize) Listen() <-chan fyne.Size {
	b.Lock()
	defer b.Unlock()
	c := make(chan fyne.Size, 16)
	c <- b.Get()
	b.channels = append(b.channels, c)
	_, file, line, ok := runtime.Caller(1)
	if ok {
		b.traces = append(b.traces, fmt.Sprintf("%s#%d", file, line))
	}
	return c
}

// Update notifies all listeners after a change.
func (b *baseSize) Update() {
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

// SizeList defines a data binding for a list of fyne.Size.
type SizeList interface {
	List
	GetBinding(int) Size
	GetSize(int) fyne.Size
	GetRef(int) *fyne.Size
	SetBinding(int, Size)
	SetSize(int, fyne.Size)
	SetRef(int, *fyne.Size)
	AddBinding(Size)
	AddSize(fyne.Size)
	AddRef(*fyne.Size)
}

// baseSizeList implements a data binding for a list of fyne.Size.
type baseSizeList struct {
	sync.Mutex
	references *[]*fyne.Size
	bindings   map[*fyne.Size]Size
	channels   []chan int
	traces     []string
}

// NewSizeList creates a new list binding with the given values.
func NewSizeList(values ...fyne.Size) SizeList {
	var references []*fyne.Size
	for i := 0; i < len(values); i++ {
		references = append(references, &values[i])
	}
	return NewSizeListRefs(&references)
}

// NewSizeListRefs creates a new list binding with the given references.
func NewSizeListRefs(references *[]*fyne.Size) SizeList {
	return &baseSizeList{
		references: references,
		bindings:   make(map[*fyne.Size]Size),
	}
}

// Length returns the number of elements in the list.
func (b *baseSizeList) Length() int {
	return len(*b.references)
}

// Get returns the binding at the given index.
func (b *baseSizeList) Get(index int) Binding {
	return b.GetBinding(index)
}

// GetBinding returns the Size at the given index.
func (b *baseSizeList) GetBinding(index int) Size {
	reference := b.GetRef(index)
	if reference == nil {
		return nil
	}
	b.Lock()
	defer b.Unlock()
	binding, ok := b.bindings[reference]
	if !ok {
		binding = NewSizeRef(reference)
		b.bindings[reference] = binding
	}
	return binding
}

// GetSize returns the fyne.Size at the given index.
func (b *baseSizeList) GetSize(index int) fyne.Size {
	if index < 0 && index >= b.Length() {
		return fyne.Size{}
	}
	return *(*b.references)[index]
}

// GetRef returns the reference at the given index.
func (b *baseSizeList) GetRef(index int) *fyne.Size {
	if index < 0 && index >= b.Length() {
		return nil
	}
	return (*b.references)[index]
}

// SetBinding updates the Size at the given index.
func (b *baseSizeList) SetBinding(index int, binding Size) {
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

// SetSize updates the fyne.Size at the given index.
func (b *baseSizeList) SetSize(index int, value fyne.Size) {
	if index < 0 && index >= b.Length() {
		return
	}
	if *(*b.references)[index] == value {
		return
	}
	b.Get(index).(Size).Set(value)
}

// SetRef updates the fyne.Size at the given index.
func (b *baseSizeList) SetRef(index int, reference *fyne.Size) {
	if index < 0 && index >= b.Length() {
		return
	}
	if (*b.references)[index] == reference {
		return
	}
	(*b.references)[index] = reference
	b.Update()
}

// AddBinding appends the Size to the list.
func (b *baseSizeList) AddBinding(binding Size) {
	index := b.Length()
	*b.references = append(*b.references, binding.GetRef())
	b.SetBinding(index, binding)
}

// AddSize appends the fyne.Size to the list.
func (b *baseSizeList) AddSize(value fyne.Size) {
	b.AddRef(&value)
}

// AddRef appends the fyne.Size to the list.
func (b *baseSizeList) AddRef(reference *fyne.Size) {
	*b.references = append(*b.references, reference)
	b.Update()
}

// Listen returns a channel through which updates will be published.
func (b *baseSizeList) Listen() <-chan int {
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
func (b *baseSizeList) Update() {
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

// String defines a data binding for a string.
type String interface {
	Binding
	Get() string
	GetRef() *string
	Set(string)
	SetRef(*string)
	Listen() <-chan string
}

// baseString implements a data binding for a string.
type baseString struct {
	sync.Mutex
	reference *string
	channels  []chan string
	traces    []string
}

// EmptyString creates a new binding with the empty value.
func EmptyString() String {
	return NewString("")
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

// Get returns the bound reference.
func (b *baseString) GetRef() *string {
	return b.reference
}

// Set updates the value of the bound reference.
func (b *baseString) Set(value string) {
	if *b.reference == value {
		return
	}
	*b.reference = value
	b.Update()
}

// Set updates the bound reference.
func (b *baseString) SetRef(reference *string) {
	if b.reference == reference {
		return
	}
	b.reference = reference
	b.Update()
}

// Listen returns a channel through which updates will be published.
func (b *baseString) Listen() <-chan string {
	b.Lock()
	defer b.Unlock()
	c := make(chan string, 16)
	c <- b.Get()
	b.channels = append(b.channels, c)
	_, file, line, ok := runtime.Caller(1)
	if ok {
		b.traces = append(b.traces, fmt.Sprintf("%s#%d", file, line))
	}
	return c
}

// Update notifies all listeners after a change.
func (b *baseString) Update() {
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

// StringList defines a data binding for a list of string.
type StringList interface {
	List
	GetBinding(int) String
	GetString(int) string
	GetRef(int) *string
	SetBinding(int, String)
	SetString(int, string)
	SetRef(int, *string)
	AddBinding(String)
	AddString(string)
	AddRef(*string)
}

// baseStringList implements a data binding for a list of string.
type baseStringList struct {
	sync.Mutex
	references *[]*string
	bindings   map[*string]String
	channels   []chan int
	traces     []string
}

// NewStringList creates a new list binding with the given values.
func NewStringList(values ...string) StringList {
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
	return b.GetBinding(index)
}

// GetBinding returns the String at the given index.
func (b *baseStringList) GetBinding(index int) String {
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

// SetBinding updates the String at the given index.
func (b *baseStringList) SetBinding(index int, binding String) {
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

// AddBinding appends the String to the list.
func (b *baseStringList) AddBinding(binding String) {
	index := b.Length()
	*b.references = append(*b.references, binding.GetRef())
	b.SetBinding(index, binding)
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

// Listen returns a channel through which updates will be published.
func (b *baseStringList) Listen() <-chan int {
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
func (b *baseStringList) Update() {
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

// URL defines a data binding for a *url.URL.
type URL interface {
	Binding
	Get() *url.URL
	GetRef() **url.URL
	Set(*url.URL)
	SetRef(**url.URL)
	Listen() <-chan *url.URL
}

// baseURL implements a data binding for a *url.URL.
type baseURL struct {
	sync.Mutex
	reference **url.URL
	channels  []chan *url.URL
	traces    []string
}

// EmptyURL creates a new binding with the empty value.
func EmptyURL() URL {
	return NewURL(nil)
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

// Get returns the bound reference.
func (b *baseURL) GetRef() **url.URL {
	return b.reference
}

// Set updates the value of the bound reference.
func (b *baseURL) Set(value *url.URL) {
	if *b.reference == value {
		return
	}
	*b.reference = value
	b.Update()
}

// Set updates the bound reference.
func (b *baseURL) SetRef(reference **url.URL) {
	if b.reference == reference {
		return
	}
	b.reference = reference
	b.Update()
}

// Listen returns a channel through which updates will be published.
func (b *baseURL) Listen() <-chan *url.URL {
	b.Lock()
	defer b.Unlock()
	c := make(chan *url.URL, 16)
	c <- b.Get()
	b.channels = append(b.channels, c)
	_, file, line, ok := runtime.Caller(1)
	if ok {
		b.traces = append(b.traces, fmt.Sprintf("%s#%d", file, line))
	}
	return c
}

// Update notifies all listeners after a change.
func (b *baseURL) Update() {
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

// URLList defines a data binding for a list of *url.URL.
type URLList interface {
	List
	GetBinding(int) URL
	GetURL(int) *url.URL
	GetRef(int) **url.URL
	SetBinding(int, URL)
	SetURL(int, *url.URL)
	SetRef(int, **url.URL)
	AddBinding(URL)
	AddURL(*url.URL)
	AddRef(**url.URL)
}

// baseURLList implements a data binding for a list of *url.URL.
type baseURLList struct {
	sync.Mutex
	references *[]**url.URL
	bindings   map[**url.URL]URL
	channels   []chan int
	traces     []string
}

// NewURLList creates a new list binding with the given values.
func NewURLList(values ...*url.URL) URLList {
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
	return b.GetBinding(index)
}

// GetBinding returns the URL at the given index.
func (b *baseURLList) GetBinding(index int) URL {
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

// SetBinding updates the URL at the given index.
func (b *baseURLList) SetBinding(index int, binding URL) {
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

// AddBinding appends the URL to the list.
func (b *baseURLList) AddBinding(binding URL) {
	index := b.Length()
	*b.references = append(*b.references, binding.GetRef())
	b.SetBinding(index, binding)
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

// Listen returns a channel through which updates will be published.
func (b *baseURLList) Listen() <-chan int {
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
func (b *baseURLList) Update() {
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

// Toggle flips the value of the bound reference.
func (b *baseBool) Toggle() {
	*b.reference = !*b.reference
	b.Update()
}
