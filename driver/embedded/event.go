package embedded

// Event is the general type of all embedded device events.
//
// Since: 2.7
type Event interface {
	isEvent()
}
