package main

import (
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/dataapi"
	"fyne.io/fyne/widget"
)

// List is a custom widget to list the contents of a set of files
type FileList struct {
	widget.Box
	widget.DataListener
}

func NewFileList() *FileList {
	return &FileList{}
}

// Bind will link the selected file to the dataItem
func (w *FileList) Source(data dataapi.DataSource) *FileList {
	w.DataListener.Source(data, w)
	return w
}

// SetFromSource updates the file list from a dataSource
func (w *FileList) SetFromSource(src dataapi.DataSource) {
	println("setting file list from source")
	dir, ok := src.(*dataapi.DirectoryDataSource)
	if !ok {
		println("data source is not a file list")
		return
	}
	w.Box.Children = []fyne.CanvasObject{}
	count := dir.Count()
	if count > 200 {
		count = 200 // in reality, lazy load with pagination
	}
	for i := 0; i < count; i++ {
		if v, ok := dir.Get(i); ok {
			if fi := v.(*dataapi.FileInfo).FileInfo(); fi != nil {
				// append all the things and do 1 refresh at the end, otherwise
				// we burn 6ms per append !
				ww := widget.NewLabel(fmt.Sprintf("%s - %d bytes", fi.Name(), fi.Size()))
				w.Box.Children = append(w.Box.Children, ww)
			}
		}
	}
	w.Refresh()
}

// Tapped for clicks
func (w *FileList) Tapped(ev *fyne.PointEvent) {
	println("tapped the file list")
}

// TappedSecondary for clicks with the other button
func (w *FileList) TappedSecondary(ev *fyne.PointEvent) {
	println("tapped 2nd")
}

// MinSize returns the minimum size for this list
func (w *FileList) MinSize() fyne.Size {
	return fyne.Size{800, 800}
}
