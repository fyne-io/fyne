package dataapi

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// TODO - add fsnotify to the DirectoryDataSource

// FileInfo is a DataItem wrapper around an os.FileInfo
type FileInfo struct {
	DataItem
	fileInfo os.FileInfo
}

func NewFileInfo(f os.FileInfo) *FileInfo {
	return &FileInfo{
		DataItem: NewBaseDataItem(),
		fileInfo: f,
	}
}

// String returns the filename
func (f *FileInfo) String() string {
	return f.fileInfo.Name()
}

// FileInfo returns the os.FileInfo component of the wrapper
func (f *FileInfo) FileInfo() os.FileInfo {
	return f.fileInfo
}

// DirectoryDataSource is an implementation of a data source based on a file system directory
type DirectoryDataSource struct {
	sync.RWMutex
	path      string
	data      []*FileInfo
	mData     sync.RWMutex
	callbacks map[int]func(item DataSource)
	id        int
}

// NewDirectoryDataSource returns a new DirectoryDataSource
func NewDirectoryDataSource(startPath string) *DirectoryDataSource {
	// get the path and fill in the data
	data := []*FileInfo{}

	walkFn := func(path string, fi os.FileInfo, err error) error {
		println("got", path, fi.Name(), fi.Size())
		if !fi.IsDir() {
			if strings.HasPrefix(fi.Name(), ".") {
				println("skip")
				return nil
			}
		}
		data = append(data, NewFileInfo(fi))
		return nil
	}

	err := filepath.Walk(startPath, walkFn)
	if err != nil {
		println("error", err.Error())
	}

	return &DirectoryDataSource{
		path:      startPath,
		data:      data,
		callbacks: make(map[int]func(DataSource)),
	}
}

// Count returns the size of the dataSource
func (b *DirectoryDataSource) Count() int {
	b.mData.RLock()
	defer b.mData.RUnlock()
	return len(b.data)
}

// String returns a string representation of the whole data source
func (b *DirectoryDataSource) String() string {
	//return ""
	value := ""
	i := 0
	for _, v := range b.data {
		if i > 0 {
			value = value + "\n"
		}
		value = value + v.String()
		i++
	}
	return value
}

// Get retuns the dataItem with the given index, and a flag denoting whether it was found
func (b *DirectoryDataSource) Get(idx int) (DataItem, bool) {
	b.mData.RLock()
	defer b.mData.RUnlock()
	if idx < 0 || idx >= len(b.data) {
		return nil, false
	}
	return b.data[idx], true
}

// DeleteItem removes an item from the map, and invokes the listeners
func (b *DirectoryDataSource) DeleteItem(idx int) {
	b.mData.Lock()
	if idx < 0 || idx >= len(b.data) {
		b.mData.Unlock()
		return
	}
	b.data = append(b.data[:idx], b.data[idx+1:]...)
	b.mData.Unlock()
	b.update()
}

// AddListener adds a new listener callback to this BaseDataItem
func (b *DirectoryDataSource) AddListener(f func(data DataSource)) int {
	b.Lock()
	defer b.Unlock()
	b.id++
	b.callbacks[b.id] = f
	return b.id
}

// DeleteListener removes the listener with the given ID
func (b *DirectoryDataSource) DeleteListener(i int) {
	b.Lock()
	defer b.Unlock()
	delete(b.callbacks, i)
}

func (b *DirectoryDataSource) update() {
	b.RLock()
	for _, f := range b.callbacks {
		f(b)
	}
	b.RUnlock()
}
