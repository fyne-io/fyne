package driver

type WithContext interface {
	RunWithContext(f func())
	RescaleContext()
}
