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

	b := BindStruct(&s)

	assert.Equal(t, 3, len(b.Keys()))
	assert.Equal(t, "bar", b.GetItem("Foo").(String).Get())

	b.GetItem("Foo").(String).Set("Content")
	assert.Equal(t, "Content", b.GetItem("Foo").(String).Get())
}

func TestBindUntypedMap(t *testing.T) {
	m := map[string]interface{}{
		"foo": "bar",
		"val": 5,
		"bas": 0.2,
	}

	b := BindUntypedMap(&m)

	assert.Equal(t, 3, len(b.Keys()))
	assert.Equal(t, "bar", b.Get("foo"))

	b.Set("Extra", "Content")
	assert.Equal(t, "Content", b.Get("Extra"))

	b.Set("foo", "new")
	assert.Equal(t, "new", b.Get("foo"))
	assert.Equal(t, "Content", b.Get("Extra"))
}

func TestUntypedMap_Delete(t *testing.T) {
	m := map[string]interface{}{
		"foo": "bar",
		"val": 5,
	}

	b := BindUntypedMap(&m)

	assert.Equal(t, 2, len(b.Keys()))
	assert.Equal(t, "bar", b.Get("foo"))
	assert.Equal(t, 5, b.Get("val"))

	b.Delete("foo")
	assert.Equal(t, 1, len(b.Keys()))
	assert.Equal(t, nil, b.Get("foo"))
	assert.Equal(t, 5, b.Get("val"))
}
