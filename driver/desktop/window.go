package desktop

// Window defines the desktop specific extensions to a fyne.Window.
//
// Since: 2.5
type Window interface {
	RunNative(func(ctx any) error) error
}
