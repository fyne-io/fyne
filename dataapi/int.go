package dataapi

import "fmt"

// Int is a DataItem that wraps an int primitive
type Int struct {
	*BaseDataItem
	value  int
	locked bool
}

// NewInt returns a new Int dataItem
func NewInt(i int) *Int {
	return &Int{
		BaseDataItem: NewBaseDataItem(),
		value:        i,
		locked:       false,
	}
}

// String returns the string representation of the Int dataItem
func (i *Int) String() string {
	return fmt.Sprintf("%d", i.value)
}

// Value gets the value as an int
func (i *Int) Value() int {
	return i.value
}

// SetInt - makes the Int dataItem implement 2 way binding
// pass a newInt to set the dataItem, and set mask to ignore
// the callback with the masked ID.  Pass a mask of 0 to trigger
// all callbacks.
func (i *Int) SetInt(newint int, mask int) {
	if i.locked {
		return
	}
	i.Lock()
	defer i.Unlock()
	i.locked = true
	i.value = newint
	for k, f := range i.callbacks {
		if k != mask {
			f(i)
		}
	}
	i.locked = false
}
