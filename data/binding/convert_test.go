package binding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoolToString(t *testing.T) {
	b := NewBool()
	s := BoolToString(b)
	assert.Equal(t, "false", s.Get())

	b.Set(true)
	assert.Equal(t, "true", s.Get())

	s.Set("false")
	assert.Equal(t, false, b.Get())
}

func TestBoolToStringWithFormat(t *testing.T) {
	b := NewBool()
	s := BoolToStringWithFormat(b, "%tly")
	assert.Equal(t, "falsely", s.Get())

	b.Set(true)
	assert.Equal(t, "truely", s.Get())

	s.Set("falsely")
	assert.Equal(t, false, b.Get())
}

func TestFloatToString(t *testing.T) {
	f := NewFloat()
	s := FloatToString(f)
	assert.Equal(t, "0.000000", s.Get())

	f.Set(0.3)
	assert.Equal(t, "0.300000", s.Get())

	s.Set("5.00")
	assert.Equal(t, 5.0, f.Get())
}

func TestFloatToStringWithFormat(t *testing.T) {
	f := NewFloat()
	s := FloatToStringWithFormat(f, "%f%%")
	assert.Equal(t, "0.000000%", s.Get())

	f.Set(0.3)
	assert.Equal(t, "0.300000%", s.Get())

	s.Set("5.00%")
	assert.Equal(t, 5.0, f.Get())
}

func TestIntToString(t *testing.T) {
	i := NewInt()
	s := IntToString(i)
	assert.Equal(t, "0", s.Get())

	i.Set(3)
	assert.Equal(t, "3", s.Get())

	s.Set("5")
	assert.Equal(t, 5, i.Get())
}

func TestIntToStringWithFormat(t *testing.T) {
	i := NewInt()
	s := IntToStringWithFormat(i, "num%d")
	assert.Equal(t, "num0", s.Get())

	i.Set(3)
	assert.Equal(t, "num3", s.Get())

	s.Set("num5")
	assert.Equal(t, 5, i.Get())
}

func TestStringToBool(t *testing.T) {
	s := NewString()
	b := StringToBool(s)
	assert.Equal(t, false, b.Get())

	s.Set("true")
	assert.Equal(t, true, b.Get())

	b.Set(false)
	assert.Equal(t, "false", s.Get())
}

func TestStringToBoolWithFormat(t *testing.T) {
	start := "falsely"
	s := BindString(&start)
	b := StringToBoolWithFormat(s, "%tly")
	assert.Equal(t, false, b.Get())

	s.Set("truely")
	assert.Equal(t, true, b.Get())

	b.Set(false)
	assert.Equal(t, "falsely", s.Get())
}

func TestStringToFloat(t *testing.T) {
	s := NewString()
	f := StringToFloat(s)
	assert.Equal(t, 0.0, f.Get())

	s.Set("3")
	assert.Equal(t, 3.0, f.Get())

	f.Set(5)
	assert.Equal(t, "5.000000", s.Get())
}

func TestStringToFloatWithFormat(t *testing.T) {
	start := "0.0%"
	s := BindString(&start)
	f := StringToFloatWithFormat(s, "%f%%")
	assert.Equal(t, 0.0, f.Get())

	s.Set("3.000000%")
	assert.Equal(t, 3.0, f.Get())

	f.Set(5)
	assert.Equal(t, "5.000000%", s.Get())
}

func TestStringToInt(t *testing.T) {
	s := NewString()
	i := StringToInt(s)
	assert.Equal(t, 0, i.Get())

	s.Set("3")
	assert.Equal(t, 3, i.Get())

	i.Set(5)
	assert.Equal(t, "5", s.Get())
}

func TestStringToIntWithFormat(t *testing.T) {
	start := "num0"
	s := BindString(&start)
	i := StringToIntWithFormat(s, "num%d")
	assert.Equal(t, 0, i.Get())

	s.Set("num3")
	assert.Equal(t, 3, i.Get())

	i.Set(5)
	assert.Equal(t, "num5", s.Get())
}
