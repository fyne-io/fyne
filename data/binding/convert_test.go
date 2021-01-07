package binding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoolToString(t *testing.T) {
	b := NewBool()
	s := BoolToString(b)
	v, err := s.Get()
	assert.Nil(t, err)
	assert.Equal(t, "false", v)

	err = b.Set(true)
	assert.Nil(t, err)
	v, err = s.Get()
	assert.Nil(t, err)
	assert.Equal(t, "true", v)

	err = s.Set("trap") // bug in fmt.SScanf means "wrong" parses as "false"
	assert.NotNil(t, err)
	_, err = b.Get()
	assert.Nil(t, err)

	err = s.Set("false")
	assert.Nil(t, err)
	v2, err := b.Get()
	assert.Nil(t, err)
	assert.Equal(t, false, v2)
}

func TestBoolToStringWithFormat(t *testing.T) {
	b := NewBool()
	s := BoolToStringWithFormat(b, "%tly")
	v, err := s.Get()
	assert.Nil(t, err)
	assert.Equal(t, "falsely", v)

	err = b.Set(true)
	assert.Nil(t, err)
	v, err = s.Get()
	assert.Nil(t, err)
	assert.Equal(t, "truely", v)

	err = s.Set("true") // valid bool but not valid format
	assert.NotNil(t, err)
	_, err = b.Get()
	assert.Nil(t, err)

	err = s.Set("falsely")
	assert.Nil(t, err)
	v2, err := b.Get()
	assert.Nil(t, err)
	assert.Equal(t, false, v2)
}

func TestFloatToString(t *testing.T) {
	f := NewFloat()
	s := FloatToString(f)
	v, err := s.Get()
	assert.Nil(t, err)
	assert.Equal(t, "0.000000", v)

	err = f.Set(0.3)
	assert.Nil(t, err)
	v, err = s.Get()
	assert.Nil(t, err)
	assert.Equal(t, "0.300000", v)

	err = s.Set("wrong")
	assert.NotNil(t, err)
	_, err = f.Get()
	assert.Nil(t, err)

	err = s.Set("5.00")
	assert.Nil(t, err)
	v2, err := f.Get()
	assert.Nil(t, err)
	assert.Equal(t, 5.0, v2)
}

func TestFloatToStringWithFormat(t *testing.T) {
	f := NewFloat()
	s := FloatToStringWithFormat(f, "%f%%")
	v, err := s.Get()
	assert.Nil(t, err)
	assert.Equal(t, "0.000000%", v)

	err = f.Set(0.3)
	assert.Nil(t, err)
	v, err = s.Get()
	assert.Nil(t, err)
	assert.Equal(t, "0.300000%", v)

	err = s.Set("4.3") // valid float64 but not valid format
	assert.NotNil(t, err)
	_, err = f.Get()
	assert.Nil(t, err)

	err = s.Set("5.00%")
	assert.Nil(t, err)
	v2, err := f.Get()
	assert.Nil(t, err)
	assert.Equal(t, 5.0, v2)
}

func TestIntToString(t *testing.T) {
	i := NewInt()
	s := IntToString(i)
	v, err := s.Get()
	assert.Nil(t, err)
	assert.Equal(t, "0", v)

	err = i.Set(3)
	assert.Nil(t, err)
	v, err = s.Get()
	assert.Nil(t, err)
	assert.Equal(t, "3", v)

	err = s.Set("wrong")
	assert.NotNil(t, err)
	_, err = i.Get()
	assert.Nil(t, err)

	err = s.Set("5")
	assert.Nil(t, err)
	v2, err := i.Get()
	assert.Nil(t, err)
	assert.Equal(t, 5, v2)
}

func TestIntToStringWithFormat(t *testing.T) {
	i := NewInt()
	s := IntToStringWithFormat(i, "num%d")
	v, err := s.Get()
	assert.Nil(t, err)
	assert.Equal(t, "num0", v)

	err = i.Set(3)
	assert.Nil(t, err)
	v, err = s.Get()
	assert.Nil(t, err)
	assert.Equal(t, "num3", v)

	err = s.Set("4") // valid int but not valid format
	assert.NotNil(t, err)
	_, err = i.Get()
	assert.Nil(t, err)

	err = s.Set("num5")
	assert.Nil(t, err)
	v2, err := i.Get()
	assert.Nil(t, err)
	assert.Equal(t, 5, v2)
}

func TestStringToBool(t *testing.T) {
	s := NewString()
	b := StringToBool(s)
	v, err := b.Get()
	assert.Nil(t, err)
	assert.Equal(t, false, v)

	err = s.Set("true")
	assert.Nil(t, err)
	v, err = b.Get()
	assert.Nil(t, err)
	assert.Equal(t, true, v)

	err = s.Set("trap") // bug in fmt.SScanf means "wrong" parses as "false"
	assert.Nil(t, err)
	_, err = b.Get()
	assert.NotNil(t, err)

	err = b.Set(false)
	assert.Nil(t, err)
	v2, err := s.Get()
	assert.Nil(t, err)
	assert.Equal(t, "false", v2)
}

func TestStringToBoolWithFormat(t *testing.T) {
	start := "falsely"
	s := BindString(&start)
	b := StringToBoolWithFormat(s, "%tly")
	v, err := b.Get()
	assert.Nil(t, err)
	assert.Equal(t, false, v)

	err = s.Set("truely")
	assert.Nil(t, err)
	v, err = b.Get()
	assert.Nil(t, err)
	assert.Equal(t, true, v)

	err = s.Set("true") // valid bool but not valid format
	assert.Nil(t, err)
	_, err = b.Get()
	assert.NotNil(t, err)

	err = b.Set(false)
	assert.Nil(t, err)
	v2, err := s.Get()
	assert.Nil(t, err)
	assert.Equal(t, "falsely", v2)
}

func TestStringToFloat(t *testing.T) {
	s := NewString()
	f := StringToFloat(s)
	v, err := f.Get()
	assert.Nil(t, err)
	assert.Equal(t, 0.0, v)

	err = s.Set("3")
	assert.Nil(t, err)
	v, err = f.Get()
	assert.Nil(t, err)
	assert.Equal(t, 3.0, v)

	err = s.Set("wrong")
	assert.Nil(t, err)
	_, err = f.Get()
	assert.NotNil(t, err)

	err = f.Set(5)
	assert.Nil(t, err)
	v2, err := s.Get()
	assert.Nil(t, err)
	assert.Equal(t, "5.000000", v2)
}

func TestStringToFloatWithFormat(t *testing.T) {
	start := "0.0%"
	s := BindString(&start)
	f := StringToFloatWithFormat(s, "%f%%")
	v, err := f.Get()
	assert.Nil(t, err)
	assert.Equal(t, 0.0, v)

	err = s.Set("3.000000%")
	assert.Nil(t, err)
	v, err = f.Get()
	assert.Nil(t, err)
	assert.Equal(t, 3.0, v)

	err = s.Set("4.3") // valid float64 but not valid format
	assert.Nil(t, err)
	_, err = f.Get()
	assert.NotNil(t, err)

	err = f.Set(5)
	assert.Nil(t, err)
	v2, err := s.Get()
	assert.Nil(t, err)
	assert.Equal(t, "5.000000%", v2)
}

func TestStringToInt(t *testing.T) {
	s := NewString()
	i := StringToInt(s)
	v, err := i.Get()
	assert.Nil(t, err)
	assert.Equal(t, 0, v)

	err = s.Set("3")
	assert.Nil(t, err)
	v, err = i.Get()
	assert.Nil(t, err)
	assert.Equal(t, 3, v)

	err = s.Set("wrong")
	assert.Nil(t, err)
	_, err = i.Get()
	assert.NotNil(t, err)

	err = i.Set(5)
	assert.Nil(t, err)
	v2, err := s.Get()
	assert.Nil(t, err)
	assert.Equal(t, "5", v2)
}

func TestStringToIntWithFormat(t *testing.T) {
	start := "num0"
	s := BindString(&start)
	i := StringToIntWithFormat(s, "num%d")
	v, err := i.Get()
	assert.Nil(t, err)
	assert.Equal(t, 0, v)

	err = s.Set("num3")
	assert.Nil(t, err)
	v, err = i.Get()
	assert.Nil(t, err)
	assert.Equal(t, 3, v)

	err = s.Set("4") // valid int but not valid format
	assert.Nil(t, err)
	_, err = i.Get()
	assert.NotNil(t, err)

	err = i.Set(5)
	assert.Nil(t, err)
	v2, err := s.Get()
	assert.Nil(t, err)
	assert.Equal(t, "num5", v2)
}
