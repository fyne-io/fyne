package dataapi

import "fmt"

// Float is a DataItem that wraps a float primitive
type Float struct {
	*BaseDataItem
	value  float64
	locked bool
}

// NewFloat returns a new Float dataItem
func NewFloat(i float64) *Float {
	return &Float{
		BaseDataItem: NewBaseDataItem(),
		value:        i,
		locked:       false,
	}
}

// String returns the string representation of the Float dataItem
func (f *Float) String() string {
	return fmt.Sprintf("%v", f.value)
}

// SetFloatWithMask - makes the Float dataItem implement 2 way binding
// pass a newInt to set the dataItem, and set mask to ignore
// the callback with the masked ID.  Pass a mask of 0 to trigger
// all callbacks.
func (f *Float) SetFloatWithMask(newfloat float64, mask int) {
	if f.locked {
		f.value = newfloat
		return
	}
	f.Lock()
	defer f.Unlock()
	f.locked = true
	f.value = newfloat
	for k, v := range f.callbacks {
		if k != mask {
			v(f)
		}
	}
	f.locked = false
}

// SetFloat sets the value with a mask of 0
func (f *Float) SetFloat(newfloat float64) {
	f.SetFloatWithMask(newfloat, 0)
}

// Value returns the float64 value of the Float dataItem
func (f *Float) Value() float64 {
	return f.value
}
