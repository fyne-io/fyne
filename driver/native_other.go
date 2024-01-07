//go:build !android

package driver

// RunNative provides a way to execute code within the platform-specific runtime context for various runtimes.
// This is mostly useful for Android where the JVM provides functionality that is not accessible directly in CGo.
// The call for most platforms will just execute passing an `UnknownContext` and returning any error reported.
//
// Since: 2.3
func RunNative(fn func(any) error) error {
	return fn(&UnknownContext{})
}
