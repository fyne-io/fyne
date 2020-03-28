// +build ios

package app

import (
	"unsafe"

	"fyne.io/fyne"
)

/*
#import <stdlib.h>
#import <Foundation/Foundation.h>

int getPreferenceBool(const char* key, int fallback);
void setPreferenceBool(const char* key, int value);
float getPreferenceFloat(const char* key, float fallback);
void setPreferenceFloat(const char* key, float value);
int getPreferenceInt(const char* key, int fallback);
void setPreferenceInt(const char* key, int value);
const char* getPreferenceString(const char* key, const char* fallback);
void setPreferenceString(const char* key, const char* value);
*/
import "C"

type iOSPreferences struct {
}

// Declare conformity with Preferences interface
var _ fyne.Preferences = (*iOSPreferences)(nil)

func (p *iOSPreferences) Bool(key string) bool {
	return p.BoolWithFallback(key, false)
}

func (p *iOSPreferences) BoolWithFallback(key string, fallback bool) bool {
	cKey := C.CString(key)
	cBool := 0
	if fallback {
		cBool = 1
	}
	ret := C.getPreferenceBool(cKey, C.int(cBool))

	C.free(unsafe.Pointer(cKey))
	return int(ret) != 0
}

func (p *iOSPreferences) SetBool(key string, value bool) {
	cKey := C.CString(key)
	cBool := 0
	if value {
		cBool = 1
	}
	C.setPreferenceBool(cKey, C.int(cBool))

	C.free(unsafe.Pointer(cKey))
}

func (p *iOSPreferences) Float(key string) float64 {
	return p.FloatWithFallback(key, 0.0)
}

func (p *iOSPreferences) FloatWithFallback(key string, fallback float64) float64 {
	cKey := C.CString(key)
	ret := C.getPreferenceFloat(cKey, C.float(fallback))

	C.free(unsafe.Pointer(cKey))
	return float64(ret)
}

func (p *iOSPreferences) SetFloat(key string, value float64) {
	cKey := C.CString(key)
	C.setPreferenceFloat(cKey, C.float(value))

	C.free(unsafe.Pointer(cKey))
}

func (p *iOSPreferences) Int(key string) int {
	return p.IntWithFallback(key, 0)
}

func (p *iOSPreferences) IntWithFallback(key string, fallback int) int {
	cKey := C.CString(key)
	ret := C.getPreferenceInt(cKey, C.int(fallback))

	C.free(unsafe.Pointer(cKey))
	return int(ret)
}

func (p *iOSPreferences) SetInt(key string, value int) {
	cKey := C.CString(key)
	C.setPreferenceInt(cKey, C.int(value))

	C.free(unsafe.Pointer(cKey))
}

func (p *iOSPreferences) String(key string) string {
	return p.StringWithFallback(key, "")
}

func (p *iOSPreferences) StringWithFallback(key, fallback string) string {
	cKey := C.CString(key)
	cFallback := C.CString(fallback)
	cRet := C.getPreferenceString(cKey, cFallback)

	C.free(unsafe.Pointer(cKey))
	C.free(unsafe.Pointer(cFallback))
	return C.GoString(cRet)
}

func (p *iOSPreferences) SetString(key string, value string) {
	cKey := C.CString(key)
	cValue := C.CString(value)
	C.setPreferenceString(cKey, cValue)

	C.free(unsafe.Pointer(cKey))
	C.free(unsafe.Pointer(cValue))
}

func newPreferences() *iOSPreferences {
	return &iOSPreferences{}
}

func loadPreferences(_ string) *iOSPreferences {
	return newPreferences() // iOS stores the ID itself so load == new
}
