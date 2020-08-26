package binding

type DataItemListener interface {
	DataChanged(DataItem)
}

type DataItem interface {
	AddListener(DataItemListener)
	RemoveListener(DataItemListener)
}

type dataItemListener struct {
	callback func(DataItem)
}

func (d *dataItemListener) DataChanged(i DataItem) {
	d.callback(i)
}

func NewDataItemListener(f func(DataItem)) DataItemListener {
	return &dataItemListener{callback: f}
}

type DataListListener interface {
	DataChanged(DataList)
}

type DataList interface {
	AddListener(DataListListener)
	RemoveListener(DataListListener)
	Length() int
	Get(int) DataItem
}

type dataListListener struct {
	callback func(DataList)
}

func (d *dataListListener) DataChanged(l DataList) {
	d.callback(l)
}

func NewDataListListener(f func(DataList)) DataListListener {
	return &dataListListener{callback: f}
}

type DataMapListener interface {
	DataChanged(DataMap)
}

type DataMap interface {
	AddListener(DataMapListener)
	RemoveListener(DataMapListener)
	Get(string) DataItem
}

type dataMapListener struct {
	callback func(DataMap)
}

func (d *dataMapListener) DataChanged(m DataMap) {
	d.callback(m)
}

func NewDataMapListener(f func(DataMap)) DataMapListener {
	return &dataMapListener{callback: f}
}
