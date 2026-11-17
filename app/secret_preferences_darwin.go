//go:build !ci && !wasm && !test_web_driver && !mobile && !tinygo

package app

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework Security

#include <stdbool.h>
#include <stdlib.h>

bool isBundled();
int keychainLoad(const char *service, const char *account, void **data, int *length);
int keychainSave(const char *service, const char *account, const void *data, int length);
*/
import "C"

import (
	"fmt"
	"unsafe"
)

const keychainAccount = "preferences"

// newSecretStore returns the store used for secret preferences.
// Apps running from a bundle (all iOS apps, and packaged macOS apps) use the Keychain.
// Unbundled binaries, such as `go run` during development, fall back to the password encrypted file store
// because the Keychain ties access to the code signature and would prompt on every rebuild.
func (a *fyneApp) newSecretStore(password []byte) secretStore {
	if !bool(C.isBundled()) {
		return a.newEncryptedSecretStore(password)
	}
	return &keychainStore{service: a.UniqueID()}
}

// keychainStore keeps the secret preferences as a generic password item in the Keychain.
type keychainStore struct {
	service string
}

func (k *keychainStore) load() ([]byte, error) {
	service := C.CString(k.service)
	defer C.free(unsafe.Pointer(service))
	account := C.CString(keychainAccount)
	defer C.free(unsafe.Pointer(account))

	var data unsafe.Pointer
	var length C.int
	switch status := C.keychainLoad(service, account, &data, &length); status {
	case 0:
		defer C.free(data)
		return C.GoBytes(data, length), nil
	case 1:
		return nil, errEmptyPreferencesStore
	default:
		return nil, fmt.Errorf("keychain read failed with status %d", int(status))
	}
}

func (k *keychainStore) save(data []byte) error {
	service := C.CString(k.service)
	defer C.free(unsafe.Pointer(service))
	account := C.CString(keychainAccount)
	defer C.free(unsafe.Pointer(account))

	var ptr unsafe.Pointer
	if len(data) > 0 {
		ptr = unsafe.Pointer(&data[0])
	}
	if status := C.keychainSave(service, account, ptr, C.int(len(data))); status != 0 {
		return fmt.Errorf("keychain write failed with status %d", int(status))
	}
	return nil
}
