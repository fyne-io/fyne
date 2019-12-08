package dataapi

// String is a DataItem that wraps a string primitive
type String struct {
	*BaseDataItem
	value  string
	locked bool
}

// NewString returns a new String object based on the given starting value
func NewString(s string) *String {
	return &String{
		BaseDataItem: NewBaseDataItem(),
		value:        s,
		locked:       false,
	}
}

// String returns the string representation of the String dataItem
func (s *String) String() string {
	return s.value
}

// Set - makes the String dataItem implement 2 way binding
// pass a newstring to set the dataItem, and set mask to ignore
// the callback with the masked ID.  Pass a mask of 0 to trigger
// all callbacks.
func (s *String) Set(newstring string, mask int) {
	if s.locked {
		return
	}
	s.Lock()
	defer s.Unlock()
	s.locked = true
	s.value = newstring
	for k, f := range s.callbacks {
		if k != mask {
			f(s)
		}
	}
	s.locked = false
}
