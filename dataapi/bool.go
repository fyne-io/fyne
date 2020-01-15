package dataapi

// Bool is a DataItem that wraps a bool primitive
type Bool struct {
	*BaseDataItem
	value  bool
	locked bool
}

// NewBool returns a new bool type
func NewBool(b bool) *Bool {
	return &Bool{
		BaseDataItem: NewBaseDataItem(),
		value:        b,
		locked:       false,
	}
}

// String returns true/false as a string, based on the bool value
func (b *Bool) String() string {
	if b.value {
		return "true"
	}
	return "false"
}

// Value gets the true/false value of the Bool
func (b *Bool) Value() bool {
	return b.value
}

// SetBoolWithMask - makes the Bool dataItem implement 2 way binding
// pass a newvalue to set the dataItem, and set mask to ignore
// the callback with the masked ID.  Pass a mask of 0 to trigger
// all callbacks.
func (b *Bool) SetBoolWithMask(newvalue bool, mask int) {
	if b.locked {
		return
	}
	b.Lock()
	defer b.Unlock()
	b.locked = true
	b.value = newvalue
	for k, f := range b.callbacks {
		if k != mask {
			f(b)
		}
	}
	b.locked = false
}

// SetBool updates the value with a mask of 0
func (b *Bool) SetBool(newvalue bool) {
	b.SetBoolWithMask(newvalue, 0)
}
