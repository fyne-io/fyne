package repository

import (
	"testing"

	"fyne.io/fyne/v2"
	"github.com/stretchr/testify/assert"
)

func TestGenericParent(t *testing.T) {
	cases := []struct {
		input  string
		expect string
		err    error
	}{
		{"foo://example.com:8042/over/there?name=ferret#nose", "foo://example.com:8042/over?name=ferret#nose", nil},
		{"file:///", "file:///", ErrURIRoot},
		{"file:///foo", "file:///", nil},
	}

	for i, c := range cases {
		t.Logf("case %d, input='%s', expect='%s', err='%v'\n", i, c.input, c.expect, c.err)

		u, err := ParseURI(c.input)
		assert.NoError(t, err)

		res, err := GenericParent(u)
		assert.Equal(t, c.err, err)

		// In the case where there is a non-nil error, res is defined
		// to be nil, so we cannot call res.String() without causing
		// a panic.
		if err == nil {
			assert.Equal(t, c.expect, res.String())
		}
	}
}

func TestGenericChild(t *testing.T) {
	cases := []struct {
		input     string
		component string
		expect    string
		err       error
	}{
		{"foo://example.com:8042/over/there?name=ferret#nose", "bar", "foo://example.com:8042/over/there/bar?name=ferret#nose", nil},
		{"file:///", "quux", "file:///quux", nil},
		{"file:///foo", "baz", "file:///foo/baz", nil},
	}

	for i, c := range cases {
		t.Logf("case %d, input='%s', component='%s', expect='%s', err='%v'\n", i, c.input, c.component, c.expect, c.err)

		u, err := ParseURI(c.input)
		assert.NoError(t, err)

		res, err := GenericChild(u, c.component)
		assert.Equal(t, c.err, err)

		assert.Equal(t, c.expect, res.String())
	}
}

var benchURI fyne.URI

func BenchmarkGenericParent(b *testing.B) {
	var uri fyne.URI
	input, _ := ParseURI("foo://example.com:8042/over/there?name=ferret#nose")

	b.ReportAllocs()

	for b.Loop() {
		uri, _ = GenericParent(input)
	}

	benchURI = uri
}

func BenchmarkGenericChild(b *testing.B) {
	var uri fyne.URI
	input, _ := ParseURI("foo://example.com:8042/over/there?name=ferret#nose")

	b.ReportAllocs()

	for b.Loop() {
		uri, _ = GenericChild(input, "bar")
	}

	benchURI = uri
}
