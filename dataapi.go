package fyne

// DataAPI interfaces - its up to the app author to implement these

type ListenerHandle int64

// DataItem
// The DataItem interface that embeds the fmt.Stringer interface should allow to handle complex
// types and at the same time allow to use the String method to handle labels. It also provides the
// opportunity to hook in to be informed of change events so that widgets can update accordingly.
type DataItem interface {
  String() string
  AddListener(DataItemFunc) ListenerHandle
  DeleteListener(ListenerHandle)
}

type DataItemFunc func(DataItem)

// DataMap
// The DataMap interface is like a DataItem except that it has many items each with a name.
// It extends DataItem so that anything returning an item can also return a map.
// The change listener is called when an item (or multiple) within the map is changed.
// AddMapListener fires when the number of map elements change.
type DataMap interface {
  Get(string) (DataItem,bool)
  AddMapListener(DataMapFunc) ListenerHandle
  DeleteMapListener(ListenerHandle)
}

type DataMapFunc func(DataMap)

// DataSource
// The DataSource interface defines an interface that returns multiple DataItems
// you can consider it like []DataItem except that it can support lazy loading and
// advanced features like paging. The change listener is notified if the number if
// items in the source changes - an addition or deletion - but not if items within it change.
type DataSource interface {
  Count() int
  Get(int) DataItem
  AddListener(DataSourceFunc) ListenerHandle
  DeleteListener(ListenerHandle)
}

type DataSourceFunc func(DataSource)
