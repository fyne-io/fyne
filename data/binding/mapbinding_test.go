package binding

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBindReflectInt(t *testing.T) {
	i := 5
	b := bindReflectInt(reflect.ValueOf(&i).Elem())
	v, err := b.(Int).Get()
	assert.Nil(t, err)
	assert.Equal(t, 5, v)

	err = b.(Int).Set(4)
	assert.Nil(t, err)
	assert.Equal(t, 4, i)

	s := "hi"
	b = bindReflectInt(reflect.ValueOf(&s).Elem())
	_, err = b.(Int).Get()
	assert.NotNil(t, err) // don't crash
}

func TestBindReflectString(t *testing.T) {
	s := "Hi"
	b := bindReflectString(reflect.ValueOf(&s).Elem())
	v, err := b.(String).Get()
	assert.Nil(t, err)
	assert.Equal(t, "Hi", v)

	err = b.(String).Set("New")
	assert.Nil(t, err)
	assert.Equal(t, "New", s)
}

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
	item, err := b.GetItem("Foo")
	assert.Nil(t, err)
	v, err := item.(String).Get()
	assert.Nil(t, err)
	assert.Equal(t, "bar", v)

	err = item.(String).Set("Content")
	assert.Nil(t, err)
	v, err = item.(String).Get()
	assert.Nil(t, err)
	assert.Equal(t, "Content", v)

	_, err = b.GetItem("Missing")
	assert.NotNil(t, err)
}

func TestBindUntypedMap(t *testing.T) {
	m := map[string]interface{}{
		"foo": "bar",
		"val": 5,
		"bas": 0.2,
	}

	b := BindUntypedMap(&m)

	assert.Equal(t, 3, len(b.Keys()))
	v, err := b.GetValue("foo")
	assert.Nil(t, err)
	assert.Equal(t, "bar", v)

	err = b.SetValue("Extra", "Content")
	assert.Nil(t, err)
	v, err = b.GetValue("Extra")
	assert.Nil(t, err)
	assert.Equal(t, "Content", v)

	err = b.SetValue("foo", "new")
	assert.Nil(t, err)
	v, err = b.GetValue("foo")
	assert.Nil(t, err)
	assert.Equal(t, "new", v)
	v, err = b.GetValue("Extra")
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
	v, err := b.GetValue("foo")
	assert.Nil(t, err)
	assert.Equal(t, "bar", v)
	v, err = b.GetValue("val")
	assert.Nil(t, err)
	assert.Equal(t, 5, v)

	b.Delete("foo")
	assert.Equal(t, 1, len(b.Keys()))
	v, err = b.GetValue("foo")
	assert.NotNil(t, err)
	assert.Equal(t, nil, v)
	v, err = b.GetValue("val")
	assert.Nil(t, err)
	assert.Equal(t, 5, v)
}

func TestUntypedMap_Set(t *testing.T) {
	m := map[string]interface{}{
		"foo": "bar",
		"val": 5,
	}

	b := BindUntypedMap(&m)
	i, err := b.GetItem("val")
	assert.Nil(t, err)
	data := i.(Untyped)

	assert.Equal(t, 2, len(b.Keys()))
	v, err := b.GetValue("foo")
	assert.Nil(t, err)
	assert.Equal(t, "bar", v)
	v, err = data.get()
	assert.Nil(t, err)
	assert.Equal(t, 5, v)

	m = map[string]interface{}{
		"foo": "new",
		"bas": "another",
		"val": 7,
	}
	err = b.Set(m)
	assert.Nil(t, err)

	assert.Equal(t, 3, len(b.Keys()))
	v, err = b.GetValue("foo")
	assert.Nil(t, err)
	assert.Equal(t, "new", v)
	v, err = data.get()
	assert.Nil(t, err)
	assert.Equal(t, 7, v)

	m = map[string]interface{}{
		"val": 9,
	}
	err = b.Set(m)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(b.Keys()))
	v, err = b.GetValue("val")
	assert.Nil(t, err)
	assert.Equal(t, 9, v)
	v, err = data.get()
	assert.Nil(t, err)
	assert.Equal(t, 9, v)
}
