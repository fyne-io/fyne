package binding

import (
	"fmt"

	"fyne.io/fyne/v2/storage"
)

type sprintfString struct {
	String

	format string
	source []DataItem
	err    error
}

// NewSprintf returns a String binding that format its content using the
// format string and the provide additional parameter that must be other
// data bindings. This data binding use fmt.Sprintf and fmt.Scanf internally
// and will have all the same limitation as those function.
//
// Since: 2.2
func NewSprintf(format string, b ...DataItem) String {
	ret := &sprintfString{
		String: NewString(),
		format: format,
		source: append(make([]DataItem, 0, len(b)), b...),
	}

	for _, value := range b {
		value.AddListener(ret)
	}

	return ret
}

func (s *sprintfString) DataChanged() {
	data := make([]interface{}, 0, len(s.source))

	s.err = nil
	for _, value := range s.source {
		switch x := value.(type) {
		case Bool:
			b, err := x.Get()
			if err != nil {
				s.err = err
				return
			}

			data = append(data, b)
		case Bytes:
			b, err := x.Get()
			if err != nil {
				s.err = err
				return
			}

			data = append(data, b)
		case Float:
			f, err := x.Get()
			if err != nil {
				s.err = err
				return
			}

			data = append(data, f)
		case Int:
			i, err := x.Get()
			if err != nil {
				s.err = err
				return
			}

			data = append(data, i)
		case Rune:
			r, err := x.Get()
			if err != nil {
				s.err = err
				return
			}

			data = append(data, r)
		case String:
			str, err := x.Get()
			if err != nil {
				s.err = err
				// Set error?
				return
			}

			data = append(data, str)
		case URI:
			u, err := x.Get()
			if err != nil {
				s.err = err
				return
			}

			data = append(data, u)
		}
	}

	r := fmt.Sprintf(s.format, data...)
	s.String.Set(r)
}

func (s *sprintfString) Get() (string, error) {
	if s.err != nil {
		return "", s.err
	}
	return s.String.Get()
}

func (s *sprintfString) Set(str string) error {
	data := make([]interface{}, 0, len(s.source))

	s.err = nil
	for _, value := range s.source {
		switch value.(type) {
		case Bool:
			data = append(data, new(bool))
		case Bytes:
			return fmt.Errorf("impossible to convert '%s' to []bytes type", str)
		case Float:
			data = append(data, new(float64))
		case Int:
			data = append(data, new(int))
		case Rune:
			data = append(data, new(rune))
		case String:
			data = append(data, new(string))
		case URI:
			data = append(data, new(string))
		}
	}

	count, err := fmt.Sscanf(str, s.format, data...)
	if err != nil {
		return err
	}

	if count != len(data) {
		return fmt.Errorf("impossible to decode more than %v parameters in '%s' with format '%s'", count, str, s.format)
	}

	for i, value := range s.source {
		switch x := value.(type) {
		case Bool:
			v := data[i].(*bool)

			err := x.Set(*v)
			if err != nil {
				return err
			}
		case Bytes:
			return fmt.Errorf("impossible to convert '%s' to []bytes type", str)
		case Float:
			v := data[i].(*float64)

			err := x.Set(*v)
			if err != nil {
				return err
			}
		case Int:
			v := data[i].(*int)

			err := x.Set(*v)
			if err != nil {
				return err
			}
		case Rune:
			v := data[i].(*rune)

			err := x.Set(*v)
			if err != nil {
				return err
			}
		case String:
			v := data[i].(*string)

			err := x.Set(*v)
			if err != nil {
				return err
			}
		case URI:
			v := data[i].(*string)

			if v == nil {
				return fmt.Errorf("URI can not be nil in '%s'", str)
			}

			uri, err := storage.ParseURI(*v)
			if err != nil {
				return err
			}

			err = x.Set(uri)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// StringToStringWithFormat creates a binding that converts a string to another string using the specified format.
// Changes to the returned String will be pushed to the passed in String and setting a new string value will parse and
// set the underlying String if it matches the format and the parse was successful.
//
// Since: 2.2
func StringToStringWithFormat(str String, format string) String {
	if format == "%s" { // Same as not using custom formatting.
		return str
	}

	return NewSprintf(format, str)
}
