package binding

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"fyne.io/fyne/v2/storage"
)

func BenchmarkBoolToString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bo := NewBool()
		s := BoolToString(bo)
		s.Get()

		bo.Set(true)
		s.Get()

		s.Set("trap")
		bo.Get()

		s.Set("false")
		bo.Get()
	}
}

func BenchmarkFloatToString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f := NewFloat()
		s := FloatToString(f)
		s.Get()

		f.Set(0.3)
		s.Get()

		s.Set("wrong")
		f.Get()

		s.Set("5.00")
		f.Get()
	}
}

func BenchmarkIntToString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		i := NewInt()
		s := IntToString(i)
		s.Get()

		i.Set(3)
		s.Get()

		s.Set("wrong")
		i.Get()

		s.Set("5")
		i.Get()
	}
}

func BenchmarkStringToBool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := NewString()
		b := StringToBool(s)
		b.Get()

		s.Set("true")
		b.Get()

		s.Set("trap") // bug in fmt.SScanf means "wrong" parses as "false"
		b.Get()

		b.Set(false)
		s.Get()
	}
}

func BenchmarkStringToFloat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := NewString()
		f := StringToFloat(s)
		f.Get()

		s.Set("3")
		f.Get()

		s.Set("wrong")
		f.Get()

		f.Set(5)
		s.Get()
	}
}

func BenchmarkStringToInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := NewString()
		i := StringToInt(s)
		i.Get()

		s.Set("3")
		i.Get()

		s.Set("wrong")
		i.Get()

		i.Set(5)
		s.Get()
	}
}

func TestBoolToString(t *testing.T) {
	b := NewBool()
	s := BoolToString(b)
	v, err := s.Get()
	require.NoError(t, err)
	assert.Equal(t, "false", v)

	err = b.Set(true)
	require.NoError(t, err)
	v, err = s.Get()
	require.NoError(t, err)
	assert.Equal(t, "true", v)

	err = s.Set("trap") // bug in fmt.SScanf means "wrong" parses as "false"
	require.Error(t, err)
	_, err = b.Get()
	require.NoError(t, err)

	err = s.Set("false")
	require.NoError(t, err)
	v2, err := b.Get()
	require.NoError(t, err)
	assert.False(t, v2)
}

func TestBoolToStringWithFormat(t *testing.T) {
	b := NewBool()
	s := BoolToStringWithFormat(b, "%tly")
	v, err := s.Get()
	require.NoError(t, err)
	assert.Equal(t, "falsely", v)

	err = b.Set(true)
	require.NoError(t, err)
	v, err = s.Get()
	require.NoError(t, err)
	assert.Equal(t, "truely", v)

	err = s.Set("true") // valid bool but not valid format
	require.Error(t, err)
	_, err = b.Get()
	require.NoError(t, err)

	err = s.Set("falsely")
	require.NoError(t, err)
	v2, err := b.Get()
	require.NoError(t, err)
	assert.False(t, v2)
}

func TestFloatToString(t *testing.T) {
	f := NewFloat()
	s := FloatToString(f)
	v, err := s.Get()
	require.NoError(t, err)
	assert.Equal(t, "0.000000", v)

	err = f.Set(0.3)
	require.NoError(t, err)
	v, err = s.Get()
	require.NoError(t, err)
	assert.Equal(t, "0.300000", v)

	err = s.Set("wrong")
	require.Error(t, err)
	_, err = f.Get()
	require.NoError(t, err)

	err = s.Set("5.00")
	require.NoError(t, err)
	v2, err := f.Get()
	require.NoError(t, err)
	assert.Equal(t, 5.0, v2)
}

func TestFloatToStringWithFormat(t *testing.T) {
	f := NewFloat()
	s := FloatToStringWithFormat(f, "%.2f%%")
	v, err := s.Get()
	require.NoError(t, err)
	assert.Equal(t, "0.00%", v)

	err = f.Set(0.3)
	require.NoError(t, err)
	v, err = s.Get()
	require.NoError(t, err)
	assert.Equal(t, "0.30%", v)

	err = s.Set("4.3") // valid float64 but not valid format
	require.Error(t, err)
	_, err = f.Get()
	require.NoError(t, err)

	err = s.Set("5.00%")
	require.NoError(t, err)
	v2, err := f.Get()
	require.NoError(t, err)
	assert.Equal(t, 5.0, v2)
}

func TestIntToString(t *testing.T) {
	i := NewInt()
	s := IntToString(i)
	v, err := s.Get()
	require.NoError(t, err)
	assert.Equal(t, "0", v)

	err = i.Set(3)
	require.NoError(t, err)
	v, err = s.Get()
	require.NoError(t, err)
	assert.Equal(t, "3", v)

	err = s.Set("wrong")
	require.Error(t, err)
	_, err = i.Get()
	require.NoError(t, err)

	err = s.Set("5")
	require.NoError(t, err)
	v2, err := i.Get()
	require.NoError(t, err)
	assert.Equal(t, 5, v2)
}

func TestIntToStringWithFormat(t *testing.T) {
	i := NewInt()
	s := IntToStringWithFormat(i, "num%d")
	v, err := s.Get()
	require.NoError(t, err)
	assert.Equal(t, "num0", v)

	err = i.Set(3)
	require.NoError(t, err)
	v, err = s.Get()
	require.NoError(t, err)
	assert.Equal(t, "num3", v)

	err = s.Set("4") // valid int but not valid format
	require.Error(t, err)
	_, err = i.Get()
	require.NoError(t, err)

	err = s.Set("num5")
	require.NoError(t, err)
	v2, err := i.Get()
	require.NoError(t, err)
	assert.Equal(t, 5, v2)
}

func TestStringToBool(t *testing.T) {
	s := NewString()
	b := StringToBool(s)
	v, err := b.Get()
	require.NoError(t, err)
	assert.False(t, v)

	err = s.Set("true")
	require.NoError(t, err)
	v, err = b.Get()
	require.NoError(t, err)
	assert.True(t, v)

	err = s.Set("trap") // bug in fmt.SScanf means "wrong" parses as "false"
	require.NoError(t, err)
	_, err = b.Get()
	require.Error(t, err)

	err = b.Set(false)
	require.NoError(t, err)
	v2, err := s.Get()
	require.NoError(t, err)
	assert.Equal(t, "false", v2)
}

func TestStringToBoolWithFormat(t *testing.T) {
	start := "falsely"
	s := BindString(&start)
	b := StringToBoolWithFormat(s, "%tly")
	v, err := b.Get()
	require.NoError(t, err)
	assert.False(t, v)

	err = s.Set("truely")
	require.NoError(t, err)
	v, err = b.Get()
	require.NoError(t, err)
	assert.True(t, v)

	err = s.Set("true") // valid bool but not valid format
	require.NoError(t, err)
	_, err = b.Get()
	require.Error(t, err)

	err = b.Set(false)
	require.NoError(t, err)
	v2, err := s.Get()
	require.NoError(t, err)
	assert.Equal(t, "falsely", v2)
}

func TestStringToFloat(t *testing.T) {
	s := NewString()
	f := StringToFloat(s)
	v, err := f.Get()
	require.NoError(t, err)
	assert.Equal(t, 0.0, v)

	err = s.Set("3")
	require.NoError(t, err)
	v, err = f.Get()
	require.NoError(t, err)
	assert.Equal(t, 3.0, v)

	err = s.Set("wrong")
	require.NoError(t, err)
	_, err = f.Get()
	require.Error(t, err)

	err = f.Set(5)
	require.NoError(t, err)
	v2, err := s.Get()
	require.NoError(t, err)
	assert.Equal(t, "5.000000", v2)
}

func TestStringToFloatWithFormat(t *testing.T) {
	start := "0.0%"
	s := BindString(&start)
	f := StringToFloatWithFormat(s, "%f%%")
	v, err := f.Get()
	require.NoError(t, err)
	assert.Equal(t, 0.0, v)

	err = s.Set("3.000000%")
	require.NoError(t, err)
	v, err = f.Get()
	require.NoError(t, err)
	assert.Equal(t, 3.0, v)

	err = s.Set("4.3") // valid float64 but not valid format
	require.NoError(t, err)
	_, err = f.Get()
	require.Error(t, err)

	err = f.Set(5)
	require.NoError(t, err)
	v2, err := s.Get()
	require.NoError(t, err)
	assert.Equal(t, "5.000000%", v2)
}

func TestStringToInt(t *testing.T) {
	s := NewString()
	i := StringToInt(s)
	v, err := i.Get()
	require.NoError(t, err)
	assert.Equal(t, 0, v)

	err = s.Set("3")
	require.NoError(t, err)
	v, err = i.Get()
	require.NoError(t, err)
	assert.Equal(t, 3, v)

	err = s.Set("wrong")
	require.NoError(t, err)
	_, err = i.Get()
	require.Error(t, err)

	err = i.Set(5)
	require.NoError(t, err)
	v2, err := s.Get()
	require.NoError(t, err)
	assert.Equal(t, "5", v2)
}

func TestStringToIntWithFormat(t *testing.T) {
	start := "num0"
	s := BindString(&start)
	i := StringToIntWithFormat(s, "num%d")
	v, err := i.Get()
	require.NoError(t, err)
	assert.Equal(t, 0, v)

	err = s.Set("num3")
	require.NoError(t, err)
	v, err = i.Get()
	require.NoError(t, err)
	assert.Equal(t, 3, v)

	err = s.Set("4") // valid int but not valid format
	require.NoError(t, err)
	_, err = i.Get()
	require.Error(t, err)

	err = i.Set(5)
	require.NoError(t, err)
	v2, err := s.Get()
	require.NoError(t, err)
	assert.Equal(t, "num5", v2)
}

func TestStringToURI(t *testing.T) {
	s := NewString()
	u := StringToURI(s)
	v, err := u.Get()
	require.NoError(t, err)
	assert.Nil(t, v)

	err = s.Set("file:///tmp/test.txt")
	require.NoError(t, err)
	v, err = u.Get()
	require.NoError(t, err)
	assert.Equal(t, "file:///tmp/test.txt", v.String())

	// TODO fix issue in URI parser whereby "wrong" is a valid URI
	//err = s.Set("wrong")
	//assert.Nil(t, err)
	//_, err = u.Get()
	//assert.NotNil(t, err)

	uri := storage.NewFileURI("/mydir/")
	err = u.Set(uri)
	require.NoError(t, err)
	v2, err := s.Get()
	require.NoError(t, err)
	assert.Equal(t, "file:///mydir/", v2)
}

func TestURIToString(t *testing.T) {
	u := NewURI()
	s := URIToString(u)
	v, err := s.Get()
	require.NoError(t, err)
	assert.Equal(t, "", v)

	err = u.Set(storage.NewFileURI("/tmp/test.txt"))
	require.NoError(t, err)
	v, err = s.Get()
	require.NoError(t, err)
	assert.Equal(t, "file:///tmp/test.txt", v)

	// TODO fix issue in URI parser whereby "wrong" is a valid URI
	//err = s.Set("wrong")
	//assert.NotNil(t, err)
	//_, err = u.Get()
	//assert.Nil(t, err)

	err = s.Set("file:///tmp/test.txt")
	require.NoError(t, err)
	v2, err := u.Get()
	require.NoError(t, err)
	assert.Equal(t, "file:///tmp/test.txt", v2.String())
}

func TestFloatToInt(t *testing.T) {
	f := NewFloat()
	i := FloatToInt(f)
	v, err := i.Get()
	require.NoError(t, err)
	assert.Equal(t, 0, v)

	err = f.Set(0.3)
	require.NoError(t, err)
	v, err = i.Get()
	require.NoError(t, err)
	assert.Equal(t, 0, v)

	err = i.Set(5)
	require.NoError(t, err)
	v2, err := f.Get()
	require.NoError(t, err)
	assert.Equal(t, 5.0, v2)
}

func TestIntToFloat(t *testing.T) {
	i := NewInt()
	f := IntToFloat(i)
	v, err := f.Get()
	require.NoError(t, err)
	assert.Equal(t, 0.0, v)

	err = i.Set(3)
	require.NoError(t, err)
	v, err = f.Get()
	require.NoError(t, err)
	assert.Equal(t, 3.0, v)

	err = f.Set(5)
	require.NoError(t, err)
	v2, err := i.Get()
	require.NoError(t, err)
	assert.Equal(t, 5, v2)
}
