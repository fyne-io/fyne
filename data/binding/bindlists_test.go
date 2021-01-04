package binding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBindFloatList(t *testing.T) {
	l := []float64{1.0, 5.0, 2.3}
	f := BindFloatList(&l)

	assert.Equal(t, 3, f.Length())
	v, err := f.GetValue(1)
	assert.Nil(t, err)
	assert.Equal(t, 5.0, v)

	assert.NotNil(t, f.(*boundFloatList).val)
	assert.Equal(t, 3, len(*(f.(*boundFloatList).val)))

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
	waitForItems()
	assert.True(t, calledList)

	child, err := f.GetItem(1)
	assert.Nil(t, err)
	child.AddListener(NewDataListener(func() {
		calledChild = true
	}))
	waitForItems()
	assert.True(t, calledChild)

	assert.NotNil(t, f.(*boundFloatList).val)
	assert.Equal(t, 3, len(*(f.(*boundFloatList).val)))

	_, err = f.GetValue(-1)
	assert.NotNil(t, err)

	calledList, calledChild = false, false
	l[1] = 4.8
	f.Reload()
	waitForItems()
	v, err = f.GetValue(1)
	assert.Nil(t, err)
	assert.Equal(t, 4.8, v)
	assert.False(t, calledList)
	assert.True(t, calledChild)

	calledList, calledChild = false, false
	l = []float64{1.0, 4.2}
	f.Reload()
	waitForItems()
	v, err = f.GetValue(1)
	assert.Nil(t, err)
	assert.Equal(t, 4.2, v)
	assert.True(t, calledList)
	assert.True(t, calledChild)

	calledList, calledChild = false, false
	l = []float64{1.0, 4.2, 5.3}
	f.Reload()
	waitForItems()
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
