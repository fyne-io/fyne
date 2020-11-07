package binding

// Float supports binding a float64 value in a Fyne application
type Float interface {
	DataItem
	Get() float64
	Set(float64)
}

// NewFloat returns a bindable float value that is managed internally.
func NewFloat() Float {
	blank := 0.0
	return &floatBind{val: &blank}
}

// BindFloat returns a new bindable value that controls the contents of the provided float64 variable.
func BindFloat(f *float64) Float {
	if f == nil {
		return NewFloat() // never allow a nil value pointer
	}

	return &floatBind{val: f}
}

type floatBind struct {
	base

	val *float64
}

func (f *floatBind) Get() float64 {
	if f.val == nil {
		return 0.0
	}
	return *f.val
}

func (f *floatBind) Set(val float64) {
	if *f.val == val {
		return
	}
	*f.val = val

	f.trigger(f)
}
