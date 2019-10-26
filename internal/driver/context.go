package driver

// WithContext allows drivers to execute within another context.
// Mostly this helps GLFW code execute within the painter's GL context.
type WithContext interface {
	RunWithContext(f func())
	RescaleContext()
	Context() interface{}
}
