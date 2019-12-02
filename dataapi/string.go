package dataapi

// String is a DataItem that wraps a string primitive
type String struct {
	*BaseDataItem
	value  string
	locked bool
}

func NewString(s string) *String {
	return &String{
		BaseDataItem: NewBaseDataItem(),
		value:        s,
		locked:       false,
	}
}

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
