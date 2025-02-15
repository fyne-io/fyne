package binding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTreeBase_AddListener(t *testing.T) {
	data := newSimpleTree()
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

func TestTreeBase_GetItem(t *testing.T) {
	data := newSimpleTree()
	f := 0.5
	data.appendItem(BindFloat(&f), "f", "")
	assert.Equal(t, 1, len(data.items))

	item, err := data.GetItem("f")
	assert.Nil(t, err)
	val, err := item.(Float).Get()
	assert.Nil(t, err)
	assert.Equal(t, f, val)

	_, err = data.GetItem("g")
	assert.NotNil(t, err)
}

func TestListBase_IDs(t *testing.T) {
	data := newSimpleTree()
	assert.Equal(t, 0, len(data.ChildIDs("")))

	data.appendItem(NewFloat(), "1", "")
	assert.Equal(t, 1, len(data.ChildIDs("")))
	assert.Equal(t, "1", data.ChildIDs("")[0])
}

func TestTreeBase_RemoveListener(t *testing.T) {
	called := false
	fn := NewDataListener(func() {
		called = true
	})
	data := newSimpleTree()
	data.listeners = append(data.listeners, fn)

	assert.Equal(t, 1, len(data.listeners))
	data.RemoveListener(fn)
	assert.Equal(t, 0, len(data.listeners))

	data.trigger()
	assert.False(t, called)
}

type simpleTree struct {
	treeBase
}

func newSimpleTree() *simpleTree {
	t := &simpleTree{}
	t.ids = map[string][]string{}
	t.items = map[string]DataItem{}

	return t
}

func TestBindStringTree(t *testing.T) {
	ids := map[string][]string{DataTreeRootID: {"1", "5", "2"}}
	l := map[string]string{"1": "one", "5": "five", "2": "two and a half"}
	f := BindStringTree(&ids, &l)

	assert.Equal(t, 3, len(f.ChildIDs(DataTreeRootID)))
	v, err := f.GetValue("5")
	assert.Nil(t, err)
	assert.Equal(t, "five", v)

	assert.NotNil(t, f.(*boundTree[string]).val)
	assert.Equal(t, 3, len(*(f.(*boundTree[string]).val)))

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
	assert.True(t, calledTree)

	child, err := f.GetItem("2")
	assert.Nil(t, err)
	child.AddListener(NewDataListener(func() {
		calledChild = true
	}))
	assert.True(t, calledChild)

	assert.NotNil(t, f.(*boundTree[float64]).val)
	assert.Equal(t, 3, len(*(f.(*boundTree[float64]).val)))

	_, err = f.GetValue("-1")
	assert.NotNil(t, err)

	calledTree, calledChild = false, false
	m["2"] = 4.8
	f.Reload()
	v, err = f.GetValue("2")
	assert.Nil(t, err)
	assert.Equal(t, 4.8, v)
	assert.False(t, calledTree)
	assert.True(t, calledChild)

	calledTree, calledChild = false, false
	m = map[string]float64{"1": 1.0, "2": 4.2}
	f.Reload()
	v, err = f.GetValue("2")
	assert.Nil(t, err)
	assert.Equal(t, 4.2, v)
	assert.True(t, calledTree)
	assert.True(t, calledChild)

	calledTree, calledChild = false, false
	m = map[string]float64{"1": 1.0, "2": 4.2, "3": 5.3}
	f.Reload()
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
	assert.Equal(t, 1, triggered)

	triggered = 0
	f.Set(map[string][]string{"": {"1", "2"}}, map[string]float64{"1": 55, "2": 77})
	assert.Equal(t, 1, triggered)

	triggered = 0
	f.SetValue("1", 5)
	assert.Zero(t, triggered)

	triggered = 0
	f.Set(map[string][]string{"": {"1", "2"}}, map[string]float64{"1": 101, "2": 98})
	assert.Zero(t, triggered)

	triggered = 0
	f.Append("1", "3", 88)
	assert.Equal(t, 1, triggered)

	triggered = 0
	f.Prepend("", "4", 23)
	assert.Equal(t, 1, triggered)

	triggered = 0
	f.Set(map[string][]string{"": {"1"}}, map[string]float64{"1": 32})
	assert.Equal(t, 1, triggered)
}
