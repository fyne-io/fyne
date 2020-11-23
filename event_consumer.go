package fyne

// EventConsumer describes an event consumer
type EventConsumer interface {
	Tapped(*PointEvent)
	TappedSecondary(*PointEvent)
	DoubleTapped(*PointEvent)
	Dragged(d *DragEvent)
	DragEnd()
}
