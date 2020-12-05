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
	from   Bool
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

	s.trigger()
}

func (s *stringFromBool) DataChanged() {
	s.trigger()
}

type stringFromFloat struct {
	base

	format string
	from   Float
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

	s.trigger()
}

func (s *stringFromFloat) DataChanged() {
	s.trigger()
}

type stringFromInt struct {
	base

	format string
	from   Int
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

	s.trigger()
}

func (s *stringFromInt) DataChanged() {
	s.trigger()
}

type stringToBool struct {
	base

	format string
	from   String
}

// StringToBool creates a binding that connects a String data item to a Bool.
// Changes to the String will be parsed and pushed to the Bool if the parse was successful, and setting
// the Bool update the String binding.
func StringToBool(str String) Bool {
	return StringToBoolWithFormat(str, "%t")
}

// StringToBoolWithFormat creates a binding that connects a String data item to a Bool and is
// presented using the specified format. Changes to the Bool will be parsed and if the format matches and
// the parse is successful it will be pushed to the String. Setting the Bool will push a formatted value
// into the String.
func StringToBoolWithFormat(str String, format string) Bool {
	v := &stringToBool{from: str, format: format}
	str.AddListener(v)
	return v
}

func (s *stringToBool) Get() bool {
	str := s.from.Get()
	if str == "" {
		return false
	}

	var val bool
	n, err := fmt.Sscanf(str, s.format, &val)
	if err != nil || n != 1 {
		fyne.LogError("bool parse error", err)
		return false
	}

	return val
}

func (s *stringToBool) Set(val bool) {
	str := fmt.Sprintf(s.format, val)
	if str == s.from.Get() {
		return
	}

	s.from.Set(str)
	s.trigger()
}

func (s *stringToBool) DataChanged() {
	s.trigger()
}

type stringToFloat struct {
	base

	format string
	from   String
}

// StringToFloat creates a binding that connects a String data item to a Float.
// Changes to the String will be parsed and pushed to the Float if the parse was successful, and setting
// the Float update the String binding.
func StringToFloat(str String) Float {
	return StringToFloatWithFormat(str, "%f")
}

// StringToFloatWithFormat creates a binding that connects a String data item to a Float and is
// presented using the specified format. Changes to the Float will be parsed and if the format matches and
// the parse is successful it will be pushed to the String. Setting the Float will push a formatted value
// into the String.
func StringToFloatWithFormat(str String, format string) Float {
	v := &stringToFloat{from: str, format: format}
	str.AddListener(v)
	return v
}

func (s *stringToFloat) Get() float64 {
	str := s.from.Get()
	if str == "" {
		return 0.0
	}

	var val float64
	n, err := fmt.Sscanf(str, s.format, &val)
	if err != nil || n != 1 {
		fyne.LogError("float64 parse error", err)
		return 0.0
	}

	return val
}

func (s *stringToFloat) Set(val float64) {
	str := fmt.Sprintf(s.format, val)
	if str == s.from.Get() {
		return
	}

	s.from.Set(str)
	s.trigger()
}

func (s *stringToFloat) DataChanged() {
	s.trigger()
}

type stringToInt struct {
	base

	format string
	from   String
}

// StringToInt creates a binding that connects a String data item to a Int.
// Changes to the String will be parsed and pushed to the Int if the parse was successful, and setting
// the Int update the String binding.
func StringToInt(str String) Int {
	return StringToIntWithFormat(str, "%d")
}

// StringToIntWithFormat creates a binding that connects a String data item to a Int and is
// presented using the specified format. Changes to the Int will be parsed and if the format matches and
// the parse is successful it will be pushed to the String. Setting the Int will push a formatted value
// into the String.
func StringToIntWithFormat(str String, format string) Int {
	v := &stringToInt{from: str, format: format}
	str.AddListener(v)
	return v
}

func (s *stringToInt) Get() int {
	str := s.from.Get()
	if str == "" {
		return 0
	}

	var val int
	n, err := fmt.Sscanf(str, s.format, &val)
	if err != nil || n != 1 {
		fyne.LogError("int parse error", err)
		return 0
	}

	return val
}

func (s *stringToInt) Set(val int) {
	str := fmt.Sprintf(s.format, val)
	if str == s.from.Get() {
		return
	}

	s.from.Set(str)
	s.trigger()
}

func (s *stringToInt) DataChanged() {
	s.trigger()
}
