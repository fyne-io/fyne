package types

type container interface {
	fyneTypeContainer()
}

type ContainerRoot struct {
}

func (c *ContainerRoot) fyneTypeContainer() {}

func IsContainer(o any) bool { // any to avoid loop if using fyne.CanvasObject
	_, ok := o.(container)
	return ok
}
