package desktop

// Window defines the desktop specific extensions to a fyne.Window.
//
// Since: 2.5
type Window interface {
	// RunNative executes the given function, passing platform-specific context
	// information related to the window. The possible contexts that may be passed
	// are defined in the `driver` package and differ per OS.
	RunNative(func(ctx any) error) error
}
