package repository

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestURIExtension(t *testing.T) {
	uri := NewFileURI("file")
	assert.Equal(t, "", uri.Extension())

	uri = NewFileURI("../file")
	assert.Equal(t, "", uri.Extension())

	uri = NewFileURI("file.txt")
	assert.Equal(t, ".txt", uri.Extension())

	uri = NewFileURI("file.tar.gz")
	assert.Equal(t, ".gz", uri.Extension())

	uri = NewFileURI("/path/.txt")
	assert.Equal(t, ".txt", uri.Extension())
}

func TestURIName(t *testing.T) {
	uri := NewFileURI("file")
	assert.Equal(t, "file", uri.Name())

	uri = NewFileURI("file.txt")
	assert.Equal(t, "file.txt", uri.Name())

	uri = NewFileURI("/somewhere/file.txt")
	assert.Equal(t, "file.txt", uri.Name())

	uri = NewFileURI("/path/.txt")
	assert.Equal(t, ".txt", uri.Name())

	if runtime.GOOS == "windows" {
		uri = NewFileURI("C:/somewhere/file.txt")
		assert.Equal(t, "file.txt", uri.Name())

		uri = NewFileURI("C:/somewhere")
		assert.Equal(t, "somewhere", uri.Name())

		uri = NewFileURI("C:/")
		assert.Equal(t, "C:", uri.Name())
	}
}

func TestURIAuthority(t *testing.T) {
	cases := []struct {
		input     string
		authority string
		user      string
		password  string
		host      string
	}{
		{"http://", "", "", "", ""},
		{"http://@", "", "", "", ""},
		{"http://foo", "foo", "", "", "foo"},
		{"http://foo@", "foo@", "foo", "", ""},
		{"http://foo:bar@", "foo:bar@", "foo", "bar", ""},
		{"http://foo::bar:@", "foo:%3Abar%3A@", "foo", ":bar:", ""},
		{"http://foo%3A::bar:@", "foo%3A:%3Abar%3A@", "foo:", ":bar:", ""},
		{"http://:bar@", ":bar@", "", "bar", ""},
		{"http://:bar@@", ":bar%40@", "", "bar@", ""},
		{"http://foo@bar", "foo@bar", "foo", "", "bar"},
		{"http://@bar", "bar", "", "", "bar"},
		{"http://foo:bar@baz", "foo:bar@baz", "foo", "bar", "baz"},
		{"http://foo:bar:baz@quux", "foo:bar%3Abaz@quux", "foo", "bar:baz", "quux"},
		{"http://foo:bar%3Abaz@quux", "foo:bar%3Abaz@quux", "foo", "bar:baz", "quux"},
		{"http://foo:bar@", "foo:bar@", "foo", "bar", ""},
	}

	for _, c := range cases {
		u, err := ParseURI(c.input)
		if !assert.NoError(t, err) || !assert.NotNil(t, u) {
			continue
		}

		assert.Equal(t, c.authority, u.Authority(), "check authority")

		if !assert.IsType(t, &uri{}, u) {
			continue
		}
		uri := u.(*uri)

		user, password := "", ""
		if uri.User != nil {
			user = uri.User.Username()
			password, _ = uri.User.Password()
		}

		assert.Equal(t, c.user, user, "check user: %q", c.input)
		assert.Equal(t, c.password, password, "check password: %q", c.input)
		assert.Equal(t, c.host, uri.Host, "check host: %q", c.input)
	}
}
