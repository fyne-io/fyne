package binding

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBindReflectInt(t *testing.T) {
	i := 5
	b := bindReflectInt(reflect.ValueOf(&i).Elem())
	v, err := b.(Int).Get()
	require.NoError(t, err)
	assert.Equal(t, 5, v)

	err = b.(Int).Set(4)
	require.NoError(t, err)
	assert.Equal(t, 4, i)

	s := "hi"
	b = bindReflectInt(reflect.ValueOf(&s).Elem())
	_, err = b.(Int).Get()
	require.Error(t, err) // don't crash
}

func TestBindReflectString(t *testing.T) {
	s := "Hi"
	b := bindReflectString(reflect.ValueOf(&s).Elem())
	v, err := b.(String).Get()
	require.NoError(t, err)
	assert.Equal(t, "Hi", v)

	err = b.(String).Set("New")
	require.NoError(t, err)
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

	assert.Len(t, b.Keys(), 3)
	item, err := b.GetItem("Foo")
	require.NoError(t, err)
	v, err := item.(String).Get()
	require.NoError(t, err)
	assert.Equal(t, "bar", v)

	err = item.(String).Set("Content")
	require.NoError(t, err)
	v, err = item.(String).Get()
	require.NoError(t, err)
	assert.Equal(t, "Content", v)

	_, err = b.GetItem("Missing")
	require.Error(t, err)
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

	assert.Len(t, b.Keys(), 3)
	v, err := b.GetValue("Foo")
	require.NoError(t, err)
	assert.Equal(t, "bar", v)

	item, err := b.GetItem("Foo")
	require.NoError(t, err)
	v, err = item.(String).Get()
	require.NoError(t, err)
	assert.Equal(t, "bar", v)

	calledMap, calledItem := false, false
	b.AddListener(NewDataListener(func() {
		calledMap = true
	}))
	assert.True(t, calledMap)

	item.AddListener(NewDataListener(func() {
		calledItem = true
	}))
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
	v, err = b.GetValue("Foo")
	require.NoError(t, err)
	assert.Equal(t, "bas", v)
	assert.False(t, calledMap)
	v, err = item.(String).Get()
	require.NoError(t, err)
	assert.Equal(t, "bas", v)
	assert.True(t, calledItem)

	calledMap, calledItem = false, false
	b.Reload()
	v, err = b.GetValue("Foo")
	require.NoError(t, err)
	assert.Equal(t, "bas", v)
	assert.False(t, calledMap)
	v, err = item.(String).Get()
	require.NoError(t, err)
	assert.Equal(t, "bas", v)
	assert.False(t, calledItem)
}

func TestBindUntypedMap(t *testing.T) {
	m := map[string]any{
		"foo": "bar",
		"val": 5,
		"bas": 0.2,
	}

	b := BindUntypedMap(&m)

	assert.Len(t, b.Keys(), 3)
	v, err := b.GetValue("foo")
	require.NoError(t, err)
	assert.Equal(t, "bar", v)

	err = b.SetValue("Extra", "Content")
	require.NoError(t, err)
	v, err = b.GetValue("Extra")
	require.NoError(t, err)
	assert.Equal(t, "Content", v)

	err = b.SetValue("foo", "new")
	require.NoError(t, err)
	v, err = b.GetValue("foo")
	require.NoError(t, err)
	assert.Equal(t, "new", v)
	v, err = b.GetValue("Extra")
	require.NoError(t, err)
	assert.Equal(t, "Content", v)
}

func TestExternalUntypedMap_Reload(t *testing.T) {
	m := map[string]any{
		"foo": "bar",
		"val": 5,
		"bas": 0.2,
	}

	b := BindUntypedMap(&m)

	assert.Len(t, b.Keys(), 3)
	v, err := b.GetValue("foo")
	require.NoError(t, err)
	assert.Equal(t, "bar", v)

	calledMap, calledChild := false, false
	b.AddListener(NewDataListener(func() {
		calledMap = true
	}))
	assert.True(t, calledMap)

	child, err := b.GetItem("foo")
	require.NoError(t, err)
	child.AddListener(NewDataListener(func() {
		calledChild = true
	}))
	assert.True(t, calledChild)

	calledMap, calledChild = false, false
	m["foo"] = "boo"
	b.Reload()
	v, err = b.GetValue("foo")
	require.NoError(t, err)
	assert.Equal(t, "boo", v)
	assert.False(t, calledMap)
	assert.True(t, calledChild)

	calledMap, calledChild = false, false
	m = map[string]any{
		"foo": "bar",
		"val": 5,
	}
	b.Reload()
	v, err = b.GetValue("foo")
	require.NoError(t, err)
	assert.Equal(t, "bar", v)
	assert.True(t, calledMap)
	assert.True(t, calledChild)

	calledMap, calledChild = false, false
	m = map[string]any{
		"foo": "bar",
		"val": 5,
		"new": "longer",
	}
	b.Reload()
	v, err = b.GetValue("foo")
	require.NoError(t, err)
	assert.Equal(t, "bar", v)
	assert.True(t, calledMap)
	assert.False(t, calledChild)
}

func TestUntypedMap_Delete(t *testing.T) {
	m := map[string]any{
		"foo": "bar",
		"val": 5,
	}

	b := BindUntypedMap(&m)

	assert.Len(t, b.Keys(), 2)
	v, err := b.GetValue("foo")
	require.NoError(t, err)
	assert.Equal(t, "bar", v)
	v, err = b.GetValue("val")
	require.NoError(t, err)
	assert.Equal(t, 5, v)

	b.Delete("foo")
	assert.Len(t, b.Keys(), 1)
	v, err = b.GetValue("foo")
	require.Error(t, err)
	assert.Nil(t, v)
	v, err = b.GetValue("val")
	require.NoError(t, err)
	assert.Equal(t, 5, v)
}

func TestUntypedMap_Set(t *testing.T) {
	m := map[string]any{
		"foo": "bar",
		"val": 5,
	}

	b := BindUntypedMap(&m)
	i, err := b.GetItem("val")
	require.NoError(t, err)
	data := i.(reflectUntyped)

	assert.Len(t, b.Keys(), 2)
	v, err := b.GetValue("foo")
	require.NoError(t, err)
	assert.Equal(t, "bar", v)
	v, err = data.get()
	require.NoError(t, err)
	assert.Equal(t, 5, v)

	m = map[string]any{
		"foo": "new",
		"bas": "another",
		"val": 7,
	}
	err = b.Set(m)
	require.NoError(t, err)

	assert.Len(t, b.Keys(), 3)
	v, err = b.GetValue("foo")
	require.NoError(t, err)
	assert.Equal(t, "new", v)
	v, err = data.get()
	require.NoError(t, err)
	assert.Equal(t, 7, v)

	m = map[string]any{
		"val": 9,
	}
	err = b.Set(m)
	require.NoError(t, err)

	assert.Len(t, b.Keys(), 1)
	v, err = b.GetValue("val")
	require.NoError(t, err)
	assert.Equal(t, 9, v)
	v, err = data.get()
	require.NoError(t, err)
	assert.Equal(t, 9, v)
}
