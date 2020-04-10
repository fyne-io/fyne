package binding

// Notifiable represents and object that can be notified.
type Notifiable interface {
	Notify(Binding)
}

// NotifyFunction will call its function when notified.
type NotifyFunction struct {
	F func(Binding)
}

// Notify is called when the given binding has changed.
func (n *NotifyFunction) Notify(data Binding) {
	if n.F == nil {
		return
	}
	n.F(data)
}
