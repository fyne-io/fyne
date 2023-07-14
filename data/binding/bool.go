package binding

type not struct {
	Bool
}

var _ Bool = (*not)(nil)

// Not returns a Bool binding that invert the value of the given data binding.
// This is providing the logical Not boolean operation as a data binding.
//
// Since 2.4
func Not(data Bool) Bool {
	return &not{Bool: data}
}

func (n *not) Get() (bool, error) {
	v, err := n.Bool.Get()
	return !v, err
}

func (n *not) Set(value bool) error {
	return n.Bool.Set(!value)
}

type and struct {
	booleans
}

var _ Bool = (*and)(nil)

// And returns a Bool binding that return true when all the passed Bool binding are
// true and false otherwise. It does apply a logical and boolean operation on all passed
// Bool bindings. This binding is two way. In case of a Set, it will propagate the value
// identically to all the Bool bindings used for its construction.
//
// Since 2.4
func And(data ...Bool) Bool {
	return &and{booleans: booleans{data: data}}
}

func (a *and) Get() (bool, error) {
	for _, d := range a.data {
		v, err := d.Get()
		if err != nil {
			return false, err
		}
		if !v {
			return false, nil
		}
	}
	return true, nil
}

func (a *and) Set(value bool) error {
	for _, d := range a.data {
		err := d.Set(value)
		if err != nil {
			return err
		}
	}
	return nil
}

type or struct {
	booleans
}

var _ Bool = (*or)(nil)

// Or returns a Bool binding that return true when at least one of the passed Bool binding
// is true and false otherwise. It does apply a logical or boolean operation on all passed
// Bool bindings. This binding is two way. In case of a Set, it will propagate the value
// identically to all the Bool bindings used for its construction.
//
// Since 2.4
func Or(data ...Bool) Bool {
	return &or{booleans: booleans{data: data}}
}

func (o *or) Get() (bool, error) {
	for _, d := range o.data {
		v, err := d.Get()
		if err != nil {
			return false, err
		}
		if v {
			return true, nil
		}
	}
	return false, nil
}

func (o *or) Set(value bool) error {
	for _, d := range o.data {
		err := d.Set(value)
		if err != nil {
			return err
		}
	}
	return nil
}

type booleans struct {
	data []Bool
}

func (g *booleans) AddListener(listener DataListener) {
	for _, d := range g.data {
		d.AddListener(listener)
	}
}

func (g *booleans) RemoveListener(listener DataListener) {
	for _, d := range g.data {
		d.RemoveListener(listener)
	}
}
