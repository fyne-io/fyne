package main

// This whole section will be removed - its a test harness to try out databinding ideas
// Will move this to a separate repo shortly

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dataapi"
)

// DataModel - a single instance of this is used by all the views in the app
type DataModel struct {
	Name         *dataapi.String
	Clock        *dataapi.Clock
	IsAvailable  *dataapi.Bool
	Size         *dataapi.Int
	OnSale       *dataapi.String
	DeliveryTime *dataapi.Float
	URL          *dataapi.String
	Image        *dataapi.String
	ActionMap    *dataapi.DataItemMap
	Actions      *dataapi.SliceDataSource

	NumWindows int
	FyneImage  *canvas.Image
}

// Images for the custom image widget
const (
	//FyneAvatar = "https://upload.wikimedia.org/wikipedia/commons/thumb/a/aa/Loch_Fyne_from_Tighcladich.jpg/2880px-Loch_Fyne_from_Tighcladich.jpg"
	FyneAvatar         = "https://avatars3.githubusercontent.com/u/36045855?s=200&v=4"
	FyneAvatarAvail    = "https://live.staticflickr.com/65535/49146531652_a12decc3af_n.jpg"
	FyneAvatarOnSale   = "https://live.staticflickr.com/65535/49216625968_aca175e6b7_z.jpg"
	FyneAvatarKangaroo = "https://live.staticflickr.com/65535/49145737868_41f6e87032_w.jpg"
	FyneAvatarSm       = "https://live.staticflickr.com/1726/28754161578_89dc92ce2c_n.jpg"
	FyneAvatarMd       = "https://live.staticflickr.com/4764/26212824798_3c3866eb7a_m.jpg"
	FyneAvatarLg       = "https://live.staticflickr.com/4177/34609667836_97ceefb52a_m.jpg"
)

// NewDataModel returns a new dataModel
func NewDataModel() *DataModel {
	img, _ := fyne.LoadResourceFromURLString(FyneAvatar)
	return &DataModel{
		Name:         dataapi.NewString(""),
		Clock:        dataapi.NewClock(),
		IsAvailable:  dataapi.NewBool(false),
		Size:         dataapi.NewInt(1),
		OnSale:       dataapi.NewString("false"),
		DeliveryTime: dataapi.NewFloat(50.0),
		URL:          dataapi.NewString("http://myurl.com"),
		Image:        dataapi.NewString(FyneAvatar),
		NumWindows:   0,
		FyneImage:    canvas.NewImageFromResource(img),
		ActionMap:    dataapi.NewDataItemMap(),
		Actions:      dataapi.NewSliceDataSource(),
	}
}
