package binding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBindStruct(t *testing.T) {
	s := struct {
		Foo string
		Val int
		Bas float64
	}{
		"bar",
		5,
		0.2,
	}

	b := BindStruct(s)

	assert.Equal(t, 3, len(b.Keys()))
	assert.Equal(t, "bar", b.Get("Foo"))

	b.Set("Foo", "Content")
	assert.Equal(t, "Content", b.Get("Foo"))
}

func TestBindUntypedMap(t *testing.T) {
	m := map[string]interface{}{
		"foo": "bar",
		"val": 5,
		"bas": 0.2,
	}

	b := BindUntypedMap(m)

	assert.Equal(t, 3, len(b.Keys()))
	assert.Equal(t, "bar", b.Get("foo"))

	b.Set("Extra", "Content")
	assert.Equal(t, "Content", b.Get("Extra"))
}
