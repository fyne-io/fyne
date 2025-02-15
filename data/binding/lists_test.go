package binding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type simpleList struct {
	listBase
}

func TestListBase_AddListener(t *testing.T) {
	data := &simpleList{}
	assert.Equal(t, 0, len(data.listeners))

	called := false
	fn := NewDataListener(func() {
		called = true
	})
	data.AddListener(fn)
	assert.Equal(t, 1, len(data.listeners))

	data.trigger()
	assert.True(t, called)
}

func TestListBase_GetItem(t *testing.T) {
	data := &simpleList{}
	f := 0.5
	data.appendItem(BindFloat(&f))
	assert.Equal(t, 1, len(data.items))

	item, err := data.GetItem(0)
	assert.Nil(t, err)
	val, err := item.(Float).Get()
	assert.Nil(t, err)
	assert.Equal(t, f, val)

	_, err = data.GetItem(5)
	assert.NotNil(t, err)
}

func TestListBase_Length(t *testing.T) {
	data := &simpleList{}
	assert.Equal(t, 0, data.Length())

	data.appendItem(NewFloat())
	assert.Equal(t, 1, data.Length())
}

func TestListBase_RemoveListener(t *testing.T) {
	called := false
	fn := NewDataListener(func() {
		called = true
	})
	data := &simpleList{}
	data.listeners = append(data.listeners, fn)

	assert.Equal(t, 1, len(data.listeners))
	data.RemoveListener(fn)
	assert.Equal(t, 0, len(data.listeners))

	data.trigger()
	assert.False(t, called)
}

func TestBindFloatList(t *testing.T) {
	l := []float64{1.0, 5.0, 2.3}
	f := BindFloatList(&l)

	assert.Equal(t, 3, f.Length())
	v, err := f.GetValue(1)
	assert.Nil(t, err)
	assert.Equal(t, 5.0, v)

	assert.NotNil(t, f.(*boundList[float64]).val)
	assert.Equal(t, 3, len(*(f.(*boundList[float64]).val)))

	_, err = f.GetValue(-1)
	assert.NotNil(t, err)
}

func TestExternalFloatList_Reload(t *testing.T) {
	l := []float64{1.0, 5.0, 2.3}
	f := BindFloatList(&l)

	assert.Equal(t, 3, f.Length())
	v, err := f.GetValue(1)
	assert.Nil(t, err)
	assert.Equal(t, 5.0, v)

	calledList, calledChild := false, false
	f.AddListener(NewDataListener(func() {
		calledList = true
	}))
	assert.True(t, calledList)

	child, err := f.GetItem(1)
	assert.Nil(t, err)
	child.AddListener(NewDataListener(func() {
		calledChild = true
	}))
	assert.True(t, calledChild)

	assert.NotNil(t, f.(*boundList[float64]).val)
	assert.Equal(t, 3, len(*(f.(*boundList[float64]).val)))

	_, err = f.GetValue(-1)
	assert.NotNil(t, err)

	calledList, calledChild = false, false
	l[1] = 4.8
	f.Reload()
	v, err = f.GetValue(1)
	assert.Nil(t, err)
	assert.Equal(t, 4.8, v)
	assert.False(t, calledList)
	assert.True(t, calledChild)

	calledList, calledChild = false, false
	l = []float64{1.0, 4.2}
	f.Reload()
	v, err = f.GetValue(1)
	assert.Nil(t, err)
	assert.Equal(t, 4.2, v)
	assert.True(t, calledList)
	assert.True(t, calledChild)

	calledList, calledChild = false, false
	l = []float64{1.0, 4.2, 5.3}
	f.Reload()
	v, err = f.GetValue(1)
	assert.Nil(t, err)
	assert.Equal(t, 4.2, v)
	assert.True(t, calledList)
	assert.False(t, calledChild)
}

func TestNewFloatList(t *testing.T) {
	f := NewFloatList()
	assert.Equal(t, 0, f.Length())

	_, err := f.GetValue(-1)
	assert.NotNil(t, err)
}

func TestFloatList_Append(t *testing.T) {
	f := NewFloatList()
	assert.Equal(t, 0, f.Length())

	f.Append(0.5)
	assert.Equal(t, 1, f.Length())
}

func TestFloatList_GetValue(t *testing.T) {
	f := NewFloatList()

	err := f.Append(1.3)
	assert.Nil(t, err)
	v, err := f.GetValue(0)
	assert.Nil(t, err)
	assert.Equal(t, 1.3, v)

	err = f.Append(0.2)
	assert.Nil(t, err)
	v, err = f.GetValue(1)
	assert.Nil(t, err)
	assert.Equal(t, 0.2, v)

	err = f.SetValue(1, 0.5)
	assert.Nil(t, err)
	v, err = f.GetValue(1)
	assert.Nil(t, err)
	assert.Equal(t, 0.5, v)
}

func TestFloatList_Remove(t *testing.T) {
	f := NewFloatList()
	f.Append(0.5)
	f.Append(0.7)
	f.Append(0.3)
	assert.Equal(t, 3, f.Length())

	f.Remove(0.5)
	assert.Equal(t, 2, f.Length())
	f.Remove(0.3)
	assert.Equal(t, 1, f.Length())
}

func TestFloatList_Set(t *testing.T) {
	l := []float64{1.0, 5.0, 2.3}
	f := BindFloatList(&l)
	i, err := f.GetItem(1)
	assert.Nil(t, err)
	data := i.(Float)

	assert.Equal(t, 3, f.Length())
	v, err := f.GetValue(1)
	assert.Nil(t, err)
	assert.Equal(t, 5.0, v)
	v, err = data.Get()
	assert.Nil(t, err)
	assert.Equal(t, 5.0, v)

	l = []float64{1.2, 5.2, 2.2, 4.2}
	err = f.Set(l)
	assert.Nil(t, err)

	assert.Equal(t, 4, f.Length())
	v, err = f.GetValue(1)
	assert.Nil(t, err)
	assert.Equal(t, 5.2, v)
	v, err = data.Get()
	assert.Nil(t, err)
	assert.Equal(t, 5.2, v)

	l = []float64{1.3, 5.3}
	err = f.Set(l)
	assert.Nil(t, err)

	assert.Equal(t, 2, f.Length())
	v, err = f.GetValue(0)
	assert.Nil(t, err)
	assert.Equal(t, 1.3, v)
	v, err = data.Get()
	assert.Nil(t, err)
	assert.Equal(t, 5.3, v)
}

func TestFloatList_NotifyOnlyOnceWhenChange(t *testing.T) {
	f := NewFloatList()
	triggered := 0
	f.AddListener(NewDataListener(func() {
		triggered++
	}))
	assert.Equal(t, 1, triggered)

	triggered = 0
	f.Set([]float64{55, 77})
	assert.Equal(t, 1, triggered)

	triggered = 0
	f.SetValue(0, 5)
	assert.Zero(t, triggered)

	triggered = 0
	f.Set([]float64{101, 98})
	assert.Zero(t, triggered)

	triggered = 0
	f.Append(88)
	assert.Equal(t, 1, triggered)

	triggered = 0
	f.Prepend(23)
	assert.Equal(t, 1, triggered)

	triggered = 0
	f.Set([]float64{32})
	assert.Equal(t, 1, triggered)
}
