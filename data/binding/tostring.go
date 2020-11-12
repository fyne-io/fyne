// auto-generated
// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //

package binding

import (
	"fmt"

	"fyne.io/fyne"
)

type stringFromBool struct {
	base

	format string
	from Bool
}

// BoolToString creates a binding that connects a Bool data item to a String.
// Changes to the Bool will be pushed to the String and setting the string will parse and set the
// Bool if the parse was successful.
func BoolToString(v Bool) String {
	return BoolToStringWithFormat(v, "%t")
}

// BoolToStringWithFormat creates a binding that connects a Bool data item to a String and is
// presented using the specified format. Changes to the Bool will be pushed to the String and setting
// the string will parse and set the Bool if the string matches the format and its parse was successful.
func BoolToStringWithFormat(v Bool, format string) String {
	str := &stringFromBool{from: v, format: format}
	v.AddListener(str)
	return str
}

func (s *stringFromBool) Get() string {
	val := s.from.Get()

	return fmt.Sprintf(s.format, val)
}

func (s *stringFromBool) Set(str string) {
	var val bool
	n, err := fmt.Sscanf(str, s.format, &val)
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

	format string
	from Float
}

// FloatToString creates a binding that connects a Float data item to a String.
// Changes to the Float will be pushed to the String and setting the string will parse and set the
// Float if the parse was successful.
func FloatToString(v Float) String {
	return FloatToStringWithFormat(v, "%f")
}

// FloatToStringWithFormat creates a binding that connects a Float data item to a String and is
// presented using the specified format. Changes to the Float will be pushed to the String and setting
// the string will parse and set the Float if the string matches the format and its parse was successful.
func FloatToStringWithFormat(v Float, format string) String {
	str := &stringFromFloat{from: v, format: format}
	v.AddListener(str)
	return str
}

func (s *stringFromFloat) Get() string {
	val := s.from.Get()

	return fmt.Sprintf(s.format, val)
}

func (s *stringFromFloat) Set(str string) {
	var val float64
	n, err := fmt.Sscanf(str, s.format, &val)
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

	format string
	from Int
}

// IntToString creates a binding that connects a Int data item to a String.
// Changes to the Int will be pushed to the String and setting the string will parse and set the
// Int if the parse was successful.
func IntToString(v Int) String {
	return IntToStringWithFormat(v, "%d")
}

// IntToStringWithFormat creates a binding that connects a Int data item to a String and is
// presented using the specified format. Changes to the Int will be pushed to the String and setting
// the string will parse and set the Int if the string matches the format and its parse was successful.
func IntToStringWithFormat(v Int, format string) String {
	str := &stringFromInt{from: v, format: format}
	v.AddListener(str)
	return str
}

func (s *stringFromInt) Get() string {
	val := s.from.Get()

	return fmt.Sprintf(s.format, val)
}

func (s *stringFromInt) Set(str string) {
	var val int
	n, err := fmt.Sscanf(str, s.format, &val)
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
