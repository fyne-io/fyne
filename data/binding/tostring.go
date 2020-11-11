// auto-generated
// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //

package binding

import (
	"fmt"

	"fyne.io/fyne"
)

type stringFromBool struct {
	base

	from Bool
}

// BoolToString creates a binding that connects a Bool data item to a String.
// Changes to the Bool will be pushed to the String and setting the string will parse and set the
// Bool if the parse was successful.
func BoolToString(v Bool) String {
	str := &stringFromBool{from: v}
	v.AddListener(str)
	return str
}

func (s *stringFromBool) Get() string {
	val := s.from.Get()

	return fmt.Sprintf("%t", val)
}

func (s *stringFromBool) Set(str string) {
	var val bool
	n, err := fmt.Sscanf(str, "%t", &val)
	if err != nil || n != 1 {
		fyne.LogError("bool parse error", err)
		return
	}
	if val == s.from.Get() {
		return
	}
	s.from.Set(val)

	s.trigger(s)
}

func (s *stringFromBool) DataChanged(_ DataItem) {
	s.trigger(s)
}

type stringFromFloat struct {
	base

	from Float
}

// FloatToString creates a binding that connects a Float data item to a String.
// Changes to the Float will be pushed to the String and setting the string will parse and set the
// Float if the parse was successful.
func FloatToString(v Float) String {
	str := &stringFromFloat{from: v}
	v.AddListener(str)
	return str
}

func (s *stringFromFloat) Get() string {
	val := s.from.Get()

	return fmt.Sprintf("%f", val)
}

func (s *stringFromFloat) Set(str string) {
	var val float64
	n, err := fmt.Sscanf(str, "%f", &val)
	if err != nil || n != 1 {
		fyne.LogError("float64 parse error", err)
		return
	}
	if val == s.from.Get() {
		return
	}
	s.from.Set(val)

	s.trigger(s)
}

func (s *stringFromFloat) DataChanged(_ DataItem) {
	s.trigger(s)
}

type stringFromInt struct {
	base

	from Int
}

// IntToString creates a binding that connects a Int data item to a String.
// Changes to the Int will be pushed to the String and setting the string will parse and set the
// Int if the parse was successful.
func IntToString(v Int) String {
	str := &stringFromInt{from: v}
	v.AddListener(str)
	return str
}

func (s *stringFromInt) Get() string {
	val := s.from.Get()

	return fmt.Sprintf("%d", val)
}

func (s *stringFromInt) Set(str string) {
	var val int
	n, err := fmt.Sscanf(str, "%d", &val)
	if err != nil || n != 1 {
		fyne.LogError("int parse error", err)
		return
	}
	if val == s.from.Get() {
		return
	}
	s.from.Set(val)

	s.trigger(s)
}

func (s *stringFromInt) DataChanged(_ DataItem) {
	s.trigger(s)
}
