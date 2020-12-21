package binding

import "errors"

var outOfBounds = errors.New("index out of bounds")

// DataList is the base interface for all bindable data lists.
//
// Since: 2.0.0
type DataList interface {
	DataItem
	GetItem(int) (DataItem, bool)
	Length() int
}

type listBase struct {
	base
	items []DataItem
}

// GetItem returns the DataItem at the specified index.
func (b *listBase) GetItem(i int) (DataItem, bool) {
	if i < 0 || i >= len(b.items) {
		return nil, false
	}

	return b.items[i], true
}

// Length returns the number of items in this data list.
func (b *listBase) Length() int {
	return len(b.items)
}

func (b *listBase) appendItem(i DataItem) {
	b.items = append(b.items, i)

	b.trigger()
}

func (b *listBase) prependItem(i DataItem) {
	b.items = append([]DataItem{i}, b.items...)

	b.trigger()
}
