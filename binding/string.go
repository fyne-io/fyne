package binding

// String supports binding a string value in a Fyne application
type String interface {
	DataItem
	Get() string
	Set(string)
}
