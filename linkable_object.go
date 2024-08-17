package fyne

type LinkableObject interface {
	// SetParent is used to set the parent object pointer. Should be used by the object where this widget is added to.
	SetParent(object LinkableObject)

	// Parent is used to get the parent object pointer. Can be used to access the parent object and its object map.
	Parent() LinkableObject
}
