//+build !android

package driver

func RunNative(fn func(interface{}) error) error {
	return fn(&UnknownContext{})
}
