package fyne

// NameableObject defines the standard behaviours of any nameable object
type NameableObject interface {
	// ID returns the object ID. It can be used to retrieve the object in the parent's object map.
	ID() string

	// SetId is used to set the object id, then it can be used to retrieve the object from the parent's object map.
	SetId(string) error
}
