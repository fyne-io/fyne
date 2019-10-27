package main

import (
	"fyne.io/fyne/dataapi"
)

// This is the dataModel - a single instance of this is used by all the views in the app
type dataModel struct {
	Name         *dataapi.String
	Clock        *dataapi.Clock
	IsAvailable  *dataapi.Bool
	Size         *dataapi.Int
	OnSale       *dataapi.String
	DeliveryTime *dataapi.Float
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
		Clock:        dataapi.NewClock(),
		IsAvailable:  dataapi.NewBool(false),
		Size:         dataapi.NewInt(int(SizeSmall)),
		OnSale:       dataapi.NewString("false"),
		DeliveryTime: dataapi.NewFloat(50.0),
		NumWindows:   0,
	}
}
