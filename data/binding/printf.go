package binding

import (
	"fmt"
)

type sprintfString struct {
	String

	format string
	source []DataItem
	err    error
}

// NewSprintf returns a String binding that format its content using the
// format string and the provide additional parameter that must be other
// data bindings.
//
// Since: 2.2
func NewSprintf(format string, b ...DataItem) (String, error) {
	ret := &sprintfString{String: NewString(), format: format}

	for _, value := range b {
		value.AddListener(ret)

		ret.source = append(ret.source, value)
	}

	return ret, nil
}

func (s *sprintfString) DataChanged() {
	data := make([]interface{}, 0)

	s.err = nil
	for _, value := range s.source {
		switch x := value.(type) {
		case Bool:
			{
				b, err := x.Get()
				if err != nil {
					s.err = err
					return
				}

				data = append(data, b)
				break
			}
		case Bytes:
			{
				b, err := x.Get()
				if err != nil {
					s.err = err
					return
				}

				data = append(data, b)
				break
			}
		case Float:
			{
				f, err := x.Get()
				if err != nil {
					s.err = err
					return
				}

				data = append(data, f)
				break
			}
		case Int:
			{
				i, err := x.Get()
				if err != nil {
					s.err = err
					return
				}

				data = append(data, i)
				break
			}
		case Rune:
			{
				r, err := x.Get()
				if err != nil {
					s.err = err
					return
				}

				data = append(data, r)
				break
			}
		case String:
			{
				str, err := x.Get()
				if err != nil {
					s.err = err
					// Set error?
					return
				}

				data = append(data, str)
				break
			}
		case URI:
			{
				u, err := x.Get()
				if err != nil {
					s.err = err
					return
				}

				data = append(data, u)
				break
			}
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
	data := make([]interface{}, 0)

	s.err = nil
	for _, value := range s.source {
		switch value.(type) {
		case Bool:
			{
				data = append(data, new(bool))
				break
			}
		case Bytes:
			{
				return fmt.Errorf("impossible to convert '%s' to []bytes type", str)
			}
		case Float:
			{
				data = append(data, new(float64))
				break
			}
		case Int:
			{
				data = append(data, new(int))
				break
			}
		case Rune:
			{
				data = append(data, new(rune))
				break
			}
		case String:
			{
				data = append(data, new(string))
				break
			}
		case URI:
			{
				return fmt.Errorf("impossible to convert '%s' to fyne.URI type", str)
			}
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
			{
				v := data[i].(*bool)

				err := x.Set(*v)
				if err != nil {
					return err
				}

				break
			}
		case Bytes:
			{
				return fmt.Errorf("impossible to convert '%s' to []bytes type", str)
			}
		case Float:
			{
				v := data[i].(*float64)

				err := x.Set(*v)
				if err != nil {
					return err
				}

				break
			}
		case Int:
			{
				v := data[i].(*int)

				err := x.Set(*v)
				if err != nil {
					return err
				}

				break
			}
		case Rune:
			{
				v := data[i].(*rune)

				err := x.Set(*v)
				if err != nil {
					return err
				}

				break
			}
		case String:
			{
				v := data[i].(*string)

				err := x.Set(*v)
				if err != nil {
					return err
				}

				break
			}
		case URI:
			{
				return fmt.Errorf("impossible to convert '%s' to fyne.URI type", str)
			}
		}
	}

	return nil
}
