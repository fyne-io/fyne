package repository

import "testing"

var benchString string

func BenchmarkURIString(b *testing.B) {
	var str string
	input, _ := ParseURI("foo://example.com:8042/over/there?name=ferret#nose")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		str = input.String()
	}

	benchString = str
}
