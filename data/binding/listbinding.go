package binding

// DataList is the base interface for all bindable data lists.
//
// Since: 2.0.0
type DataList interface {
	DataItem
	GetItem(int) DataItem
	Length() int
}

type listBase struct {
	base
	val []DataItem
}

// GetItem returns the DataItem at the specified index.
func (b *listBase) GetItem(i int) DataItem {
	if i < 0 || i >= len(b.val) {
		return nil
	}

	return b.val[i]
}

// Length returns the number of items in this data list.
func (b *listBase) Length() int {
	return len(b.val)
}

func (b *listBase) appendItem(i DataItem) {
	b.val = append(b.val, i)

	b.trigger()
}

func (b *listBase) prependItem(i DataItem) {
	b.val = append([]DataItem{i}, b.val...)

	b.trigger()
}
