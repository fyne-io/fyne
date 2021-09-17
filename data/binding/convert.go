// auto-generated
// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //

package binding

import (
	"fmt"

	"fyne.io/fyne/v2"
)

type stringFromBool struct {
	base

	format string

	from Bool
}

// BoolToString creates a binding that connects a Bool data item to a String.
// Changes to the Bool will be pushed to the String and setting the string will parse and set the
// Bool if the parse was successful.
//
// Since: 2.0
func BoolToString(v Bool) String {
	str := &stringFromBool{from: v}
	v.AddListener(str)
	return str
}

// BoolToStringWithFormat creates a binding that connects a Bool data item to a String and is
// presented using the specified format. Changes to the Bool will be pushed to the String and setting
// the string will parse and set the Bool if the string matches the format and its parse was successful.
//
// Since: 2.0
func BoolToStringWithFormat(v Bool, format string) String {
	if format == "%t" { // Same as not using custom formatting.
		return BoolToString(v)
	}

	str := &stringFromBool{from: v, format: format}
	v.AddListener(str)
	return str
}

func (s *stringFromBool) Get() (string, error) {
	val, err := s.from.Get()
	if err != nil {
		return "", err
	}

	if s.format != "" {
		return fmt.Sprintf(s.format, val), nil
	}

	return formatBool(val), nil
}

func (s *stringFromBool) Set(str string) error {
	var val bool
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
		new, err := parseBool(str)
		if err != nil {
			return err
		}
		val = new
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

func (s *stringFromBool) DataChanged() {
	s.lock.RLock()
	defer s.lock.RUnlock()
	s.trigger()
}

type stringFromFloat struct {
	base

	format string

	from Float
}

// FloatToString creates a binding that connects a Float data item to a String.
// Changes to the Float will be pushed to the String and setting the string will parse and set the
// Float if the parse was successful.
//
// Since: 2.0
func FloatToString(v Float) String {
	str := &stringFromFloat{from: v}
	v.AddListener(str)
	return str
}

// FloatToStringWithFormat creates a binding that connects a Float data item to a String and is
// presented using the specified format. Changes to the Float will be pushed to the String and setting
// the string will parse and set the Float if the string matches the format and its parse was successful.
//
// Since: 2.0
func FloatToStringWithFormat(v Float, format string) String {
	if format == "%f" { // Same as not using custom formatting.
		return FloatToString(v)
	}

	str := &stringFromFloat{from: v, format: format}
	v.AddListener(str)
	return str
}

func (s *stringFromFloat) Get() (string, error) {
	val, err := s.from.Get()
	if err != nil {
		return "", err
	}

	if s.format != "" {
		return fmt.Sprintf(s.format, val), nil
	}

	return formatFloat(val), nil
}

func (s *stringFromFloat) Set(str string) error {
	var val float64
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
		new, err := parseFloat(str)
		if err != nil {
			return err
		}
		val = new
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

func (s *stringFromFloat) DataChanged() {
	s.lock.RLock()
	defer s.lock.RUnlock()
	s.trigger()
}

type stringFromInt struct {
	base

	format string

	from Int
}

// IntToString creates a binding that connects a Int data item to a String.
// Changes to the Int will be pushed to the String and setting the string will parse and set the
// Int if the parse was successful.
//
// Since: 2.0
func IntToString(v Int) String {
	str := &stringFromInt{from: v}
	v.AddListener(str)
	return str
}

// IntToStringWithFormat creates a binding that connects a Int data item to a String and is
// presented using the specified format. Changes to the Int will be pushed to the String and setting
// the string will parse and set the Int if the string matches the format and its parse was successful.
//
// Since: 2.0
func IntToStringWithFormat(v Int, format string) String {
	if format == "%d" { // Same as not using custom formatting.
		return IntToString(v)
	}

	str := &stringFromInt{from: v, format: format}
	v.AddListener(str)
	return str
}

func (s *stringFromInt) Get() (string, error) {
	val, err := s.from.Get()
	if err != nil {
		return "", err
	}

	if s.format != "" {
		return fmt.Sprintf(s.format, val), nil
	}

	return formatInt(val), nil
}

func (s *stringFromInt) Set(str string) error {
	var val int
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
		new, err := parseInt(str)
		if err != nil {
			return err
		}
		val = new
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

func (s *stringFromInt) DataChanged() {
	s.lock.RLock()
	defer s.lock.RUnlock()
	s.trigger()
}

type stringFromURI struct {
	base

	from URI
}

// URIToString creates a binding that connects a URI data item to a String.
// Changes to the URI will be pushed to the String and setting the string will parse and set the
// URI if the parse was successful.
//
// Since: 2.1
func URIToString(v URI) String {
	str := &stringFromURI{from: v}
	v.AddListener(str)
	return str
}

func (s *stringFromURI) Get() (string, error) {
	val, err := s.from.Get()
	if err != nil {
		return "", err
	}

	return uriToString(val)
}

func (s *stringFromURI) Set(str string) error {
	val, err := uriFromString(str)
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

func (s *stringFromURI) DataChanged() {
	s.lock.RLock()
	defer s.lock.RUnlock()
	s.trigger()
}

type stringToBool struct {
	base

	format string

	from String
}

// StringToBool creates a binding that connects a String data item to a Bool.
// Changes to the String will be parsed and pushed to the Bool if the parse was successful, and setting
// the Bool update the String binding.
//
// Since: 2.0
func StringToBool(str String) Bool {
	v := &stringToBool{from: str}
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

	v := &stringToBool{from: str, format: format}
	str.AddListener(v)
	return v
}

func (s *stringToBool) Get() (bool, error) {
	str, err := s.from.Get()
	if str == "" || err != nil {
		return false, err
	}

	var val bool
	if s.format != "" {
		n, err := fmt.Sscanf(str, s.format+" ", &val) // " " denotes match to end of string
		if err != nil {
			return false, err
		}
		if n != 1 {
			return false, errParseFailed
		}
	} else {
		new, err := parseBool(str)
		if err != nil {
			return false, err
		}
		val = new
	}

	return val, nil
}

func (s *stringToBool) Set(val bool) error {
	var str string
	if s.format != "" {
		str = fmt.Sprintf(s.format, val)
	} else {
		str = formatBool(val)
	}

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

func (s *stringToBool) DataChanged() {
	s.lock.RLock()
	defer s.lock.RUnlock()
	s.trigger()
}

type stringToFloat struct {
	base

	format string

	from String
}

// StringToFloat creates a binding that connects a String data item to a Float.
// Changes to the String will be parsed and pushed to the Float if the parse was successful, and setting
// the Float update the String binding.
//
// Since: 2.0
func StringToFloat(str String) Float {
	v := &stringToFloat{from: str}
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

	v := &stringToFloat{from: str, format: format}
	str.AddListener(v)
	return v
}

func (s *stringToFloat) Get() (float64, error) {
	str, err := s.from.Get()
	if str == "" || err != nil {
		return 0.0, err
	}

	var val float64
	if s.format != "" {
		n, err := fmt.Sscanf(str, s.format+" ", &val) // " " denotes match to end of string
		if err != nil {
			return 0.0, err
		}
		if n != 1 {
			return 0.0, errParseFailed
		}
	} else {
		new, err := parseFloat(str)
		if err != nil {
			return 0.0, err
		}
		val = new
	}

	return val, nil
}

func (s *stringToFloat) Set(val float64) error {
	var str string
	if s.format != "" {
		str = fmt.Sprintf(s.format, val)
	} else {
		str = formatFloat(val)
	}

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

func (s *stringToFloat) DataChanged() {
	s.lock.RLock()
	defer s.lock.RUnlock()
	s.trigger()
}

type stringToInt struct {
	base

	format string

	from String
}

// StringToInt creates a binding that connects a String data item to a Int.
// Changes to the String will be parsed and pushed to the Int if the parse was successful, and setting
// the Int update the String binding.
//
// Since: 2.0
func StringToInt(str String) Int {
	v := &stringToInt{from: str}
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

	v := &stringToInt{from: str, format: format}
	str.AddListener(v)
	return v
}

func (s *stringToInt) Get() (int, error) {
	str, err := s.from.Get()
	if str == "" || err != nil {
		return 0, err
	}

	var val int
	if s.format != "" {
		n, err := fmt.Sscanf(str, s.format+" ", &val) // " " denotes match to end of string
		if err != nil {
			return 0, err
		}
		if n != 1 {
			return 0, errParseFailed
		}
	} else {
		new, err := parseInt(str)
		if err != nil {
			return 0, err
		}
		val = new
	}

	return val, nil
}

func (s *stringToInt) Set(val int) error {
	var str string
	if s.format != "" {
		str = fmt.Sprintf(s.format, val)
	} else {
		str = formatInt(val)
	}

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

func (s *stringToInt) DataChanged() {
	s.lock.RLock()
	defer s.lock.RUnlock()
	s.trigger()
}

type stringToURI struct {
	base

	from String
}

// StringToURI creates a binding that connects a String data item to a URI.
// Changes to the String will be parsed and pushed to the URI if the parse was successful, and setting
// the URI update the String binding.
//
// Since: 2.1
func StringToURI(str String) URI {
	v := &stringToURI{from: str}
	str.AddListener(v)
	return v
}

func (s *stringToURI) Get() (fyne.URI, error) {
	str, err := s.from.Get()
	if str == "" || err != nil {
		return fyne.URI(nil), err
	}

	return uriFromString(str)
}

func (s *stringToURI) Set(val fyne.URI) error {
	str, err := uriToString(val)
	if err != nil {
		return err
	}
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

func (s *stringToURI) DataChanged() {
	s.lock.RLock()
	defer s.lock.RUnlock()
	s.trigger()
}
