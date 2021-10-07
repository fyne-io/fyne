package mobile

import "fyne.io/fyne/v2"

func RequestPermision(permision string){
	if driver, ok := fyne.CurrentApp().Driver().(*mobileDriver); ok {
		if driver.app == nil { // not yet running
			fyne.LogError("Cannot show keyboard before app is running", nil)
			return
		}

		driver.app.RequestPermision(permision)
	}
}