package binding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBindStringTree(t *testing.T) {
	ids := map[string][]string{DataTreeRootID: {"1", "5", "2"}}
	l := map[string]string{"1": "one", "5": "five", "2": "two and a half"}
	f := BindStringTree(&ids, &l)

	assert.Equal(t, 3, len(f.ChildIDs(DataTreeRootID)))
	v, err := f.GetValue("5")
	assert.Nil(t, err)
	assert.Equal(t, "five", v)

	assert.NotNil(t, f.(*boundStringTree).val)
	assert.Equal(t, 3, len(*(f.(*boundStringTree).val)))

	_, err = f.GetValue("nan")
	assert.NotNil(t, err)
}

func TestExternalFloatTree_Reload(t *testing.T) {
	i := map[string][]string{"": {"1", "2"}, "1": {"3"}}
	m := map[string]float64{"1": 1.0, "2": 5.0, "3": 2.3}
	f := BindFloatTree(&i, &m)

	assert.Equal(t, 2, len(f.ChildIDs("")))
	v, err := f.GetValue("2")
	assert.Nil(t, err)
	assert.Equal(t, 5.0, v)

	calledTree, calledChild := false, false
	f.AddListener(NewDataListener(func() {
		calledTree = true
	}))
	waitForItems()
	assert.True(t, calledTree)

	child, err := f.GetItem("2")
	assert.Nil(t, err)
	child.AddListener(NewDataListener(func() {
		calledChild = true
	}))
	waitForItems()
	assert.True(t, calledChild)

	assert.NotNil(t, f.(*boundFloatTree).val)
	assert.Equal(t, 3, len(*(f.(*boundFloatTree).val)))

	_, err = f.GetValue("-1")
	assert.NotNil(t, err)

	calledTree, calledChild = false, false
	m["2"] = 4.8
	f.Reload()
	waitForItems()
	v, err = f.GetValue("2")
	assert.Nil(t, err)
	assert.Equal(t, 4.8, v)
	assert.False(t, calledTree)
	assert.True(t, calledChild)

	calledTree, calledChild = false, false
	m = map[string]float64{"1": 1.0, "2": 4.2}
	f.Reload()
	waitForItems()
	v, err = f.GetValue("2")
	assert.Nil(t, err)
	assert.Equal(t, 4.2, v)
	assert.True(t, calledTree)
	assert.True(t, calledChild)

	calledTree, calledChild = false, false
	m = map[string]float64{"1": 1.0, "2": 4.2, "3": 5.3}
	f.Reload()
	waitForItems()
	v, err = f.GetValue("2")
	assert.Nil(t, err)
	assert.Equal(t, 4.2, v)
	assert.True(t, calledTree)
	assert.False(t, calledChild)
}

func TestNewStringTree(t *testing.T) {
	f := NewStringTree()
	assert.Equal(t, 0, len(f.ChildIDs(DataTreeRootID)))

	_, err := f.GetValue("NaN")
	assert.NotNil(t, err)
}

func TestStringTree_Append(t *testing.T) {
	f := NewStringTree()
	assert.Equal(t, 0, len(f.ChildIDs(DataTreeRootID)))

	f.Append(DataTreeRootID, "5", "five")
	assert.Equal(t, 1, len(f.ChildIDs(DataTreeRootID)))
}

func TestStringTree_GetValue(t *testing.T) {
	f := NewStringTree()

	err := f.Append(DataTreeRootID, "1", "1.3")
	assert.Nil(t, err)
	v, err := f.GetValue("1")
	assert.Nil(t, err)
	assert.Equal(t, "1.3", v)

	err = f.Append(DataTreeRootID, "fraction", "0.2")
	assert.Nil(t, err)
	v, err = f.GetValue("fraction")
	assert.Nil(t, err)
	assert.Equal(t, "0.2", v)

	err = f.SetValue("1", "0.5")
	assert.Nil(t, err)
	v, err = f.GetValue("1")
	assert.Nil(t, err)
	assert.Equal(t, "0.5", v)
}

func TestStringTree_Remove(t *testing.T) {
	f := NewStringTree()
	f.Append(DataTreeRootID, "5", "five")
	f.Append(DataTreeRootID, "3", "three")
	f.Append("5", "53", "fifty three")
	assert.Equal(t, 2, len(f.ChildIDs(DataTreeRootID)))
	assert.Equal(t, 1, len(f.ChildIDs("5")))

	f.Remove("5")
	assert.Equal(t, 1, len(f.ChildIDs(DataTreeRootID)))
	assert.Equal(t, 0, len(f.ChildIDs("5")))
}

func TestFloatTree_Set(t *testing.T) {
	ids := map[string][]string{"": {"1", "2"}, "1": {"3"}}
	m := map[string]float64{"1": 1.0, "2": 5.0, "3": 2.3}
	f := BindFloatTree(&ids, &m)
	i, err := f.GetItem("2")
	assert.Nil(t, err)
	data := i.(Float)

	assert.Equal(t, 2, len(f.ChildIDs("")))
	v, err := f.GetValue("2")
	assert.Nil(t, err)
	assert.Equal(t, 5.0, v)
	v, err = data.Get()
	assert.Nil(t, err)
	assert.Equal(t, 5.0, v)

	ids = map[string][]string{"": {"1", "2"}, "1": {"3", "4"}}
	m = map[string]float64{"1": 1.2, "2": 5.2, "3": 2.2, "4": 4.2}
	err = f.Set(ids, m)
	assert.Nil(t, err)

	assert.Equal(t, 2, len(f.ChildIDs("1")))
	v, err = f.GetValue("2")
	assert.Nil(t, err)
	assert.Equal(t, 5.2, v)
	v, err = data.Get()
	assert.Nil(t, err)
	assert.Equal(t, 5.2, v)

	ids = map[string][]string{"": {"1", "2"}}
	m = map[string]float64{"1": 1.3, "2": 5.3}
	err = f.Set(ids, m)
	assert.Nil(t, err)

	assert.Equal(t, 2, len(f.ChildIDs("")))
	v, err = f.GetValue("1")
	assert.Nil(t, err)
	assert.Equal(t, 1.3, v)
	v, err = data.Get()
	assert.Nil(t, err)
	assert.Equal(t, 5.3, v)
}

func TestFloatTree_NotifyOnlyOnceWhenChange(t *testing.T) {
	f := NewFloatTree()
	triggered := 0
	f.AddListener(NewDataListener(func() {
		triggered++
	}))
	waitForItems()
	assert.Equal(t, 1, triggered)

	triggered = 0
	f.Set(map[string][]string{"": {"1", "2"}}, map[string]float64{"1": 55, "2": 77})
	waitForItems()
	assert.Equal(t, 1, triggered)

	triggered = 0
	f.SetValue("1", 5)
	waitForItems()
	assert.Zero(t, triggered)

	triggered = 0
	f.Set(map[string][]string{"": {"1", "2"}}, map[string]float64{"1": 101, "2": 98})
	waitForItems()
	assert.Zero(t, triggered)

	triggered = 0
	f.Append("1", "3", 88)
	waitForItems()
	assert.Equal(t, 1, triggered)

	triggered = 0
	f.Prepend("", "4", 23)
	waitForItems()
	assert.Equal(t, 1, triggered)

	triggered = 0
	f.Set(map[string][]string{"": {"1"}}, map[string]float64{"1": 32})
	waitForItems()
	assert.Equal(t, 1, triggered)
}
