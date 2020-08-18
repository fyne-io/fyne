// +build ios

package app

import (
	"path/filepath"
	"unsafe"

	"fyne.io/fyne"
)

/*
#import <stdbool.h>
#import <stdlib.h>
#import <Foundation/Foundation.h>

bool getPreferenceBool(const char* key, bool fallback);
void setPreferenceBool(const char* key, bool value);
float getPreferenceFloat(const char* key, float fallback);
void setPreferenceFloat(const char* key, float value);
int getPreferenceInt(const char* key, int fallback);
void setPreferenceInt(const char* key, int value);
const char* getPreferenceString(const char* key, const char* fallback);
void setPreferenceString(const char* key, const char* value);
void removePreferenceValue(const char* key);
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
	defer C.free(unsafe.Pointer(cKey))

	return bool(C.getPreferenceBool(cKey, C.bool(fallback)))
}

func (p *iOSPreferences) SetBool(key string, value bool) {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	C.setPreferenceBool(cKey, C.bool(value))
}

func (p *iOSPreferences) Float(key string) float64 {
	return p.FloatWithFallback(key, 0.0)
}

func (p *iOSPreferences) FloatWithFallback(key string, fallback float64) float64 {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	return float64(C.getPreferenceFloat(cKey, C.float(fallback)))
}

func (p *iOSPreferences) SetFloat(key string, value float64) {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	C.setPreferenceFloat(cKey, C.float(value))
}

func (p *iOSPreferences) Int(key string) int {
	return p.IntWithFallback(key, 0)
}

func (p *iOSPreferences) IntWithFallback(key string, fallback int) int {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	return int(C.getPreferenceInt(cKey, C.int(fallback)))
}

func (p *iOSPreferences) SetInt(key string, value int) {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	C.setPreferenceInt(cKey, C.int(value))
}

func (p *iOSPreferences) String(key string) string {
	return p.StringWithFallback(key, "")
}

func (p *iOSPreferences) StringWithFallback(key, fallback string) string {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	cFallback := C.CString(fallback)
	defer C.free(unsafe.Pointer(cFallback))

	return C.GoString(C.getPreferenceString(cKey, cFallback))
}

func (p *iOSPreferences) SetString(key string, value string) {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))

	C.setPreferenceString(cKey, cValue)
}

func (p *iOSPreferences) RemoveValue(key string) {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	C.removePreferenceValue(cKey)
}

func newPreferences(_ *fyneApp) *iOSPreferences {
	return &iOSPreferences{}
}

// storageRoot returns the location of the app storage
func (a *fyneApp) storageRoot() string {
	return filepath.Join(rootConfigDir(), a.uniqueID)
}
