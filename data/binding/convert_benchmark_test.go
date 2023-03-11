package binding

import (
	"testing"
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
