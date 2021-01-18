// auto-generated
// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //

package binding

import (
	"fmt"
)

type stringFromBool struct {
	base

	format string
	from   Bool
}

// BoolToString creates a binding that connects a Bool data item to a String.
// Changes to the Bool will be pushed to the String and setting the string will parse and set the
// Bool if the parse was successful.
//
// Since: 2.0
func BoolToString(v Bool) String {
	return BoolToStringWithFormat(v, "%t")
}

// BoolToStringWithFormat creates a binding that connects a Bool data item to a String and is
// presented using the specified format. Changes to the Bool will be pushed to the String and setting
// the string will parse and set the Bool if the string matches the format and its parse was successful.
//
// Since: 2.0
func BoolToStringWithFormat(v Bool, format string) String {
	str := &stringFromBool{from: v, format: format}
	v.AddListener(str)
	return str
}

func (s *stringFromBool) Get() (string, error) {
	val, err := s.from.Get()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(s.format, val), nil
}

func (s *stringFromBool) Set(str string) error {
	var val bool
	n, err := fmt.Sscanf(str, s.format+" ", &val) // " " denotes match to end of string
	if err != nil {
		return err
	}
	if n != 1 {
		return errParseFailed
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
	from   Float
}

// FloatToString creates a binding that connects a Float data item to a String.
// Changes to the Float will be pushed to the String and setting the string will parse and set the
// Float if the parse was successful.
//
// Since: 2.0
func FloatToString(v Float) String {
	return FloatToStringWithFormat(v, "%f")
}

// FloatToStringWithFormat creates a binding that connects a Float data item to a String and is
// presented using the specified format. Changes to the Float will be pushed to the String and setting
// the string will parse and set the Float if the string matches the format and its parse was successful.
//
// Since: 2.0
func FloatToStringWithFormat(v Float, format string) String {
	str := &stringFromFloat{from: v, format: format}
	v.AddListener(str)
	return str
}

func (s *stringFromFloat) Get() (string, error) {
	val, err := s.from.Get()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(s.format, val), nil
}

func (s *stringFromFloat) Set(str string) error {
	var val float64
	n, err := fmt.Sscanf(str, s.format+" ", &val) // " " denotes match to end of string
	if err != nil {
		return err
	}
	if n != 1 {
		return errParseFailed
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
	from   Int
}

// IntToString creates a binding that connects a Int data item to a String.
// Changes to the Int will be pushed to the String and setting the string will parse and set the
// Int if the parse was successful.
//
// Since: 2.0
func IntToString(v Int) String {
	return IntToStringWithFormat(v, "%d")
}

// IntToStringWithFormat creates a binding that connects a Int data item to a String and is
// presented using the specified format. Changes to the Int will be pushed to the String and setting
// the string will parse and set the Int if the string matches the format and its parse was successful.
//
// Since: 2.0
func IntToStringWithFormat(v Int, format string) String {
	str := &stringFromInt{from: v, format: format}
	v.AddListener(str)
	return str
}

func (s *stringFromInt) Get() (string, error) {
	val, err := s.from.Get()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(s.format, val), nil
}

func (s *stringFromInt) Set(str string) error {
	var val int
	n, err := fmt.Sscanf(str, s.format+" ", &val) // " " denotes match to end of string
	if err != nil {
		return err
	}
	if n != 1 {
		return errParseFailed
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

type stringToBool struct {
	base

	format string
	from   String
}

// StringToBool creates a binding that connects a String data item to a Bool.
// Changes to the String will be parsed and pushed to the Bool if the parse was successful, and setting
// the Bool update the String binding.
//
// Since: 2.0
func StringToBool(str String) Bool {
	return StringToBoolWithFormat(str, "%t")
}

// StringToBoolWithFormat creates a binding that connects a String data item to a Bool and is
// presented using the specified format. Changes to the Bool will be parsed and if the format matches and
// the parse is successful it will be pushed to the String. Setting the Bool will push a formatted value
// into the String.
//
// Since: 2.0
func StringToBoolWithFormat(str String, format string) Bool {
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
	n, err := fmt.Sscanf(str, s.format+" ", &val) // " " denotes match to end of string
	if err != nil {
		return false, err
	}
	if n != 1 {
		return false, errParseFailed
	}

	return val, nil
}

func (s *stringToBool) Set(val bool) error {
	str := fmt.Sprintf(s.format, val)
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
	from   String
}

// StringToFloat creates a binding that connects a String data item to a Float.
// Changes to the String will be parsed and pushed to the Float if the parse was successful, and setting
// the Float update the String binding.
//
// Since: 2.0
func StringToFloat(str String) Float {
	return StringToFloatWithFormat(str, "%f")
}

// StringToFloatWithFormat creates a binding that connects a String data item to a Float and is
// presented using the specified format. Changes to the Float will be parsed and if the format matches and
// the parse is successful it will be pushed to the String. Setting the Float will push a formatted value
// into the String.
//
// Since: 2.0
func StringToFloatWithFormat(str String, format string) Float {
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
	n, err := fmt.Sscanf(str, s.format+" ", &val) // " " denotes match to end of string
	if err != nil {
		return 0.0, err
	}
	if n != 1 {
		return 0.0, errParseFailed
	}

	return val, nil
}

func (s *stringToFloat) Set(val float64) error {
	str := fmt.Sprintf(s.format, val)
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
	from   String
}

// StringToInt creates a binding that connects a String data item to a Int.
// Changes to the String will be parsed and pushed to the Int if the parse was successful, and setting
// the Int update the String binding.
//
// Since: 2.0
func StringToInt(str String) Int {
	return StringToIntWithFormat(str, "%d")
}

// StringToIntWithFormat creates a binding that connects a String data item to a Int and is
// presented using the specified format. Changes to the Int will be parsed and if the format matches and
// the parse is successful it will be pushed to the String. Setting the Int will push a formatted value
// into the String.
//
// Since: 2.0
func StringToIntWithFormat(str String, format string) Int {
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
	n, err := fmt.Sscanf(str, s.format+" ", &val) // " " denotes match to end of string
	if err != nil {
		return 0, err
	}
	if n != 1 {
		return 0, errParseFailed
	}

	return val, nil
}

func (s *stringToInt) Set(val int) error {
	str := fmt.Sprintf(s.format, val)
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
