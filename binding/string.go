package binding

type String interface {
	DataItem
	Get() string
	Set(string)
}
