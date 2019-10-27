package dataapi

// A collection of interfaces to test DataItems against
// This is used in the DataAPI widget wrappers to test if
// the given dataItem supports writing to a given type.
//
// Depending on the widget, if the dataItem implements the
// settable, then it is used to perform 2-way data binding
//
// eg:
// Lets say we have a checkbox bound to a dataItem
// When the source data is changed, the checkbox is updated.
//
// When the checkbox is toggled, AND the dataItem is of type
// SettableBool, then it will update the dataItem

// Settable - if the data item is settable using a string
type Settable interface {
	Set(string, int)
}

// SettableBool - if the data item is settable using a bool
type SettableBool interface {
	SetBool(bool, int)
}

// SettableInt - if the data item is settable using an int
type SettableInt interface {
	SetInt(int, int)
}

// SettableFloat - if the data item is settable using a float
type SettableFloat interface {
	SetFloat(float64, int)
}
