package widget

// SelectionMode represents the selection mode of a collection widget
//
// Since: 2.5
type SelectionMode int

const (
	// SelectionSingle allows only one item to be selected at a time
	//
	// Since: 2.5
	SelectionSingle SelectionMode = iota

	// SelectionMultiple allows multiple items to be selected at a time
	//
	// Since 2.5
	SelectionMultiple

	// SelectionNone disables selection
	//
	// Since 2.5
	SelectionNone
)
