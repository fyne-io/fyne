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
func TestBindStruct_Reload(t *testing.T) {
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
	v, err := b.GetValue("Foo")
	assert.Nil(t, err)
	assert.Equal(t, "bar", v)

	item, err := b.GetItem("Foo")
	assert.Nil(t, err)
	v, err = item.(String).Get()
	assert.Nil(t, err)
	assert.Equal(t, "bar", v)

	calledMap, calledItem := false, false
	b.AddListener(NewDataListener(func() {
		calledMap = true
	}))
	waitForItems()
	assert.True(t, calledMap)

	item.AddListener(NewDataListener(func() {
		calledItem = true
	}))
	waitForItems()
	assert.True(t, calledItem)

	s = struct {
		Foo string
		Val int
		Bas float64
	}{
		"bas",
		2,
		1.2,
	}

	calledMap, calledItem = false, false
	b.Reload()
	waitForItems()
	v, err = b.GetValue("Foo")
	assert.Nil(t, err)
	assert.Equal(t, "bas", v)
	assert.False(t, calledMap)
	v, err = item.(String).Get()
	assert.Nil(t, err)
	assert.Equal(t, "bas", v)
	assert.True(t, calledItem)

	calledMap, calledItem = false, false
	b.Reload()
	waitForItems()
	v, err = b.GetValue("Foo")
	assert.Nil(t, err)
	assert.Equal(t, "bas", v)
	assert.False(t, calledMap)
	v, err = item.(String).Get()
	assert.Nil(t, err)
	assert.Equal(t, "bas", v)
	assert.False(t, calledItem)
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

func TestExternalUntypedMap_Reload(t *testing.T) {
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

	calledMap, calledChild := false, false
	b.AddListener(NewDataListener(func() {
		calledMap = true
	}))
	waitForItems()
	assert.True(t, calledMap)

	child, err := b.GetItem("foo")
	assert.Nil(t, err)
	child.AddListener(NewDataListener(func() {
		calledChild = true
	}))
	waitForItems()
	assert.True(t, calledChild)

	calledMap, calledChild = false, false
	m["foo"] = "boo"
	b.Reload()
	waitForItems()
	v, err = b.GetValue("foo")
	assert.Nil(t, err)
	assert.Equal(t, "boo", v)
	assert.False(t, calledMap)
	assert.True(t, calledChild)

	calledMap, calledChild = false, false
	m = map[string]interface{}{
		"foo": "bar",
		"val": 5,
	}
	b.Reload()
	waitForItems()
	v, err = b.GetValue("foo")
	assert.Nil(t, err)
	assert.Equal(t, "bar", v)
	assert.True(t, calledMap)
	assert.True(t, calledChild)

	calledMap, calledChild = false, false
	m = map[string]interface{}{
		"foo": "bar",
		"val": 5,
		"new": "longer",
	}
	b.Reload()
	waitForItems()
	v, err = b.GetValue("foo")
	assert.Nil(t, err)
	assert.Equal(t, "bar", v)
	assert.True(t, calledMap)
	assert.False(t, calledChild)
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
