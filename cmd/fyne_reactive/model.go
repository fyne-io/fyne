package main

// This whole section will be removed - its a test harness to try out databinding ideas
// Will move this to a separate repo shortly

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
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
	URL          *dataapi.String
	Image        *dataapi.String
	NumWindows   int
	FyneImage    *canvas.Image
}

type size int

const (
	SizeSmall size = iota + 1
	SizeMedium
	SizeLarge
)

const (
	//FyneAvatar = "https://upload.wikimedia.org/wikipedia/commons/thumb/a/aa/Loch_Fyne_from_Tighcladich.jpg/2880px-Loch_Fyne_from_Tighcladich.jpg"
	FyneAvatar = "https://avatars3.githubusercontent.com/u/36045855?s=200&v=4"
	FyneAvatarMd = "https://avatars3.githubusercontent.com/u/36045855?s=400&v=4"
	FyneAvatarLg = "https://avatars3.githubusercontent.com/u/36045855?s=600&v=4"
)

func NewDataModel() *dataModel {
	img,_ := fyne.LoadResourceFromURLString(FyneAvatar)
	return &dataModel{
		Name:         dataapi.NewString(""),
		Clock:        dataapi.NewClock(),
		IsAvailable:  dataapi.NewBool(false),
		Size:         dataapi.NewInt(int(SizeSmall)),
		OnSale:       dataapi.NewString("false"),
		DeliveryTime: dataapi.NewFloat(50.0),
		URL:          dataapi.NewString("http://myurl.com"),
		Image:        dataapi.NewString(FyneAvatar),
		NumWindows:   0,
		FyneImage:    canvas.NewImageFromResource(img),
	}
}
