// +build android
package permisions

import (
	"fyne.io/fyne/v2/internal/driver/mobile"
)


func RequestPermision(permision string){
	mobile.RequestPermision(permision)
}