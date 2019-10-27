package main

import (
	"fyne.io/fyne/dataapi"
)

// This is the dataModel - a single instance of this is used by all the views in the app
type dataModel struct {
	Name         *dataapi.String
	Time         *dataapi.Clock
	IsAvailable  *dataapi.Bool
	Size         size
	DeliveryTime int
	NumWindows   int
}

type size int

const (
	SizeSmall size = iota + 1
	SizeMedium
	SizeLarge
)

func NewDataModel() *dataModel {
	return &dataModel{
		Name:         dataapi.NewString(""),
		Time:         dataapi.NewClock(),
		IsAvailable:  dataapi.NewBool(false),
		Size:         SizeSmall,
		DeliveryTime: 50,
		NumWindows:   0,
	}
}
