// +build !ci,!cgo,!nacl

package app

import (
	"fyne.io/fyne"
)

// New panics. We were built without CGo and cannot currently continue.
// This means that playground vet will pass but also provides a slightly less friendly error when CGo is off.
func New() fyne.App {
	panic("This app was built without CGo support. Please re-build with CGO_ENABLED=1.")
}
