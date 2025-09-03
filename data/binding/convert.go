package binding

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

// BoolToString creates a binding that connects a Bool data item to a String.
// Changes to the Bool will be pushed to the String and setting the string will parse and set the
// Bool if the parse was successful.
//
// Since: 2.0
func BoolToString(v Bool) String {
	return toStringComparable(v, formatBool, parseBool)
}

// BoolToStringWithFormat creates a binding that connects a Bool data item to a String and is
// presented using the specified format. Changes to the Bool will be pushed to the String and setting
// the string will parse and set the Bool if the string matches the format and its parse was successful.
//
// Since: 2.0
func BoolToStringWithFormat(v Bool, format string) String {
	return toStringWithFormatComparable[bool](v, format, "%t", formatBool, parseBool)
}

// FloatToString creates a binding that connects a Float data item to a String.
// Changes to the Float will be pushed to the String and setting the string will parse and set the
// Float if the parse was successful.
//
// Since: 2.0
func FloatToString(v Float) String {
	return toStringComparable(v, formatFloat, parseFloat)
}

// FloatToStringWithFormat creates a binding that connects a Float data item to a String and is
// presented using the specified format. Changes to the Float will be pushed to the String and setting
// the string will parse and set the Float if the string matches the format and its parse was successful.
//
// Since: 2.0
func FloatToStringWithFormat(v Float, format string) String {
	return toStringWithFormatComparable(v, format, "%f", formatFloat, parseFloat)
}

// IntToFloat creates a binding that connects an Int data item to a Float.
//
// Since: 2.5
func IntToFloat(val Int) Float {
	v := &fromIntTo[float64]{from: val, parser: internalFloatToInt, formatter: internalIntToFloat}
	val.AddListener(v)
	return v
}

// FloatToInt creates a binding that connects a Float data item to an Int.
//
// Since: 2.5
func FloatToInt(v Float) Int {
	i := &toInt[float64]{from: v, parser: internalFloatToInt, formatter: internalIntToFloat}
	v.AddListener(i)
	return i
}

// IntToString creates a binding that connects a Int data item to a String.
// Changes to the Int will be pushed to the String and setting the string will parse and set the
// Int if the parse was successful.
//
// Since: 2.0
func IntToString(v Int) String {
	return toStringComparable(v, formatInt, parseInt)
}

// IntToStringWithFormat creates a binding that connects a Int data item to a String and is
// presented using the specified format. Changes to the Int will be pushed to the String and setting
// the string will parse and set the Int if the string matches the format and its parse was successful.
//
// Since: 2.0
func IntToStringWithFormat(v Int, format string) String {
	return toStringWithFormatComparable(v, format, "%d", formatInt, parseInt)
}

// URIToString creates a binding that connects a URI data item to a String.
// Changes to the URI will be pushed to the String and setting the string will parse and set the
// URI if the parse was successful.
//
// Since: 2.1
func URIToString(v URI) String {
	return toString(v, uriToString, storage.EqualURI, uriFromString)
}

// StringToBool creates a binding that connects a String data item to a Bool.
// Changes to the String will be parsed and pushed to the Bool if the parse was successful, and setting
// the Bool update the String binding.
//
// Since: 2.0
func StringToBool(str String) Bool {
	v := &fromStringTo[bool]{from: str, formatter: parseBool, parser: formatBool}
	str.AddListener(v)
	return v
}

// StringToBoolWithFormat creates a binding that connects a String data item to a Bool and is
// presented using the specified format. Changes to the Bool will be parsed and if the format matches and
// the parse is successful it will be pushed to the String. Setting the Bool will push a formatted value
// into the String.
//
// Since: 2.0
func StringToBoolWithFormat(str String, format string) Bool {
	if format == "%t" { // Same as not using custom format.
		return StringToBool(str)
	}

	v := &fromStringTo[bool]{from: str, format: format}
	str.AddListener(v)
	return v
}

// StringToFloat creates a binding that connects a String data item to a Float.
// Changes to the String will be parsed and pushed to the Float if the parse was successful, and setting
// the Float update the String binding.
//
// Since: 2.0
func StringToFloat(str String) Float {
	v := &fromStringTo[float64]{from: str, formatter: parseFloat, parser: formatFloat}
	str.AddListener(v)
	return v
}

// StringToFloatWithFormat creates a binding that connects a String data item to a Float and is
// presented using the specified format. Changes to the Float will be parsed and if the format matches and
// the parse is successful it will be pushed to the String. Setting the Float will push a formatted value
// into the String.
//
// Since: 2.0
func StringToFloatWithFormat(str String, format string) Float {
	if format == "%f" { // Same as not using custom format.
		return StringToFloat(str)
	}

	v := &fromStringTo[float64]{from: str, format: format}
	str.AddListener(v)
	return v
}

// StringToInt creates a binding that connects a String data item to a Int.
// Changes to the String will be parsed and pushed to the Int if the parse was successful, and setting
// the Int update the String binding.
//
// Since: 2.0
func StringToInt(str String) Int {
	v := &fromStringTo[int]{from: str, parser: formatInt, formatter: parseInt}
	str.AddListener(v)
	return v
}

// StringToIntWithFormat creates a binding that connects a String data item to a Int and is
// presented using the specified format. Changes to the Int will be parsed and if the format matches and
// the parse is successful it will be pushed to the String. Setting the Int will push a formatted value
// into the String.
//
// Since: 2.0
func StringToIntWithFormat(str String, format string) Int {
	if format == "%d" { // Same as not using custom format.
		return StringToInt(str)
	}

	v := &fromStringTo[int]{from: str, format: format}
	str.AddListener(v)
	return v
}

// StringToURI creates a binding that connects a String data item to a URI.
// Changes to the String will be parsed and pushed to the URI if the parse was successful, and setting
// the URI update the String binding.
//
// Since: 2.1
func StringToURI(str String) URI {
	v := &fromStringTo[fyne.URI]{from: str, parser: uriToString, formatter: uriFromString}
	str.AddListener(v)
	return v
}

func toString[T any](v Item[T], formatter func(T) (string, error), comparator func(T, T) bool, parser func(string) (T, error)) *toStringFrom[T] {
	str := &toStringFrom[T]{from: v, formatter: formatter, comparator: comparator, parser: parser}
	v.AddListener(str)
	return str
}

func toStringComparable[T bool | float64 | int](v Item[T], formatter func(T) (string, error), parser func(string) (T, error)) *toStringFrom[T] {
	return toString(v, formatter, func(t1, t2 T) bool { return t1 == t2 }, parser)
}

func toStringWithFormat[T any](v Item[T], format, defaultFormat string, formatter func(T) (string, error), comparator func(T, T) bool, parser func(string) (T, error)) String {
	str := toString(v, formatter, comparator, parser)
	if format != defaultFormat { // Same as not using custom formatting.
		str.format = format
	}

	return str
}

func toStringWithFormatComparable[T bool | float64 | int](v Item[T], format, defaultFormat string, formatter func(T) (string, error), parser func(string) (T, error)) String {
	return toStringWithFormat(v, format, defaultFormat, formatter, func(t1, t2 T) bool { return t1 == t2 }, parser)
}

type convertBaseItem struct {
	base
}

func (s *convertBaseItem) DataChanged() {
	s.triggerFromMain()
}

type toStringFrom[T any] struct {
	convertBaseItem

	format string

	formatter  func(T) (string, error)
	comparator func(T, T) bool
	parser     func(string) (T, error)

	from Item[T]
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

	s.trigger()
	return nil
}

type fromStringTo[T any] struct {
	convertBaseItem

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
	if s.format != "" {
		str = fmt.Sprintf(s.format, val)
	} else {
		parsed, err := s.parser(val)
		if err != nil {
			return err
		}
		str = parsed
	}

	old, err := s.from.Get()
	if str == old {
		return err
	}

	err = s.from.Set(str)
	if err != nil {
		return err
	}

	s.trigger()
	return nil
}

type toInt[T float64] struct {
	convertBaseItem

	formatter func(int) (T, error)
	parser    func(T) (int, error)

	from Item[T]
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
	err = s.from.Set(val)
	if err != nil {
		return err
	}

	s.trigger()
	return nil
}

type fromIntTo[T float64] struct {
	convertBaseItem

	formatter func(int) (T, error)
	parser    func(T) (int, error)
	from      Item[int]
}

func (s *fromIntTo[T]) Get() (T, error) {
	val, err := s.from.Get()
	if err != nil {
		return *new(T), err
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

	s.trigger()
	return nil
}
