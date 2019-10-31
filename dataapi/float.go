package dataapi

import "fmt"

// Float is a DataItem that wraps a float primitive
type Float struct {
	*BaseDataItem
	value  float64
	locked bool
}

func NewFloat(i float64) *Float {
	return &Float{
		BaseDataItem: NewBaseDataItem(),
		value:        i,
		locked:       false,
	}
}

func (i *Float) String() string {
	return fmt.Sprintf("%v", i.value)
}

// SetFloat - makes the Float dataItem implement 2 way binding
// pass a newInt to set the dataItem, and set mask to ignore
// the callback with the masked ID.  Pass a mask of 0 to trigger
// all callbacks.
func (i *Float) SetFloat(newfloat float64, mask int) {
	if i.locked {
		return
	}
	i.Lock()
	defer i.Unlock()
	i.locked = true
	i.value = newfloat
	for k, f := range i.callbacks {
		if k != mask {
			f(i)
		}
	}
	i.locked = false
}

func (f *Float) Value() float64 {
	return f.value
}
