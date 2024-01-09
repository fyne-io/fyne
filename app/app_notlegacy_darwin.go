//go:build !ci && !legacy && !wasm && !test_web_driver

package app

/*
#cgo LDFLAGS: -framework Foundation -framework UserNotifications
*/
import "C"
