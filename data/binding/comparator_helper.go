package binding

import "fyne.io/fyne/v2"

func compareURI(v1, v2 fyne.URI) bool {
	if v1 == nil && v1 == v2 {
		return true
	}
	if v1 == nil || v2 == nil {
		return false
	}
	return v1.String() == v2.String()
}
