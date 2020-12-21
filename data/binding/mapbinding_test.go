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
	item, ok := b.GetItem("Foo")
	assert.True(t, ok)
	v, err := item.(String).Get()
	assert.Nil(t, err)
	assert.Equal(t, "bar", v)

	err = item.(String).Set("Content")
	assert.Nil(t, err)
	v, err = item.(String).Get()
	assert.Nil(t, err)
	assert.Equal(t, "Content", v)
}

func TestBindUntypedMap(t *testing.T) {
	m := map[string]interface{}{
		"foo": "bar",
		"val": 5,
		"bas": 0.2,
	}

	b := BindUntypedMap(&m)

	assert.Equal(t, 3, len(b.Keys()))
	v, err := b.Get("foo")
	assert.Nil(t, err)
	assert.Equal(t, "bar", v)

	err = b.Set("Extra", "Content")
	assert.Nil(t, err)
	v, err = b.Get("Extra")
	assert.Nil(t, err)
	assert.Equal(t, "Content", v)

	err = b.Set("foo", "new")
	assert.Nil(t, err)
	v, err = b.Get("foo")
	assert.Nil(t, err)
	assert.Equal(t, "new", v)
	v, err = b.Get("Extra")
	assert.Nil(t, err)
	assert.Equal(t, "Content", v)
}

func TestUntypedMap_Delete(t *testing.T) {
	m := map[string]interface{}{
		"foo": "bar",
		"val": 5,
	}

	b := BindUntypedMap(&m)

	assert.Equal(t, 2, len(b.Keys()))
	v, err := b.Get("foo")
	assert.Nil(t, err)
	assert.Equal(t, "bar", v)
	v, err = b.Get("val")
	assert.Nil(t, err)
	assert.Equal(t, 5, v)

	b.Delete("foo")
	assert.Equal(t, 1, len(b.Keys()))
	v, err = b.Get("foo")
	assert.Nil(t, err)
	assert.Equal(t, nil, v)
	v, err = b.Get("val")
	assert.Nil(t, err)
	assert.Equal(t, 5, v)
}
