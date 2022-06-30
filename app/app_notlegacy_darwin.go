//go:build !ci && !legacy && !js && !wasm && !test_web_driver
// +build !ci,!legacy,!js,!wasm,!test_web_driver

package app

/*
#cgo LDFLAGS: -framework Foundation -framework UserNotifications
*/
import "C"
