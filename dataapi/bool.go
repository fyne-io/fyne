package dataapi

// Bool is a DataItem that wraps a bool primitive
type Bool struct {
	*BaseDataItem
	value  bool
	locked bool
}

func NewBool(b bool) *Bool {
	return &Bool{
		BaseDataItem: NewBaseDataItem(),
		value:        b,
		locked:       false,
	}
}

func (b *Bool) String() string {
	if b.value {
		return "true"
	}
	return "false"
}

// SetBool - makes the Bool dataItem implement 2 way binding
// pass a newvalue to set the dataItem, and set mask to ignore
// the callback with the masked ID.  Pass a mask of 0 to trigger
// all callbacks.
func (b *Bool) SetBool(newvalue bool, mask int) {
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
