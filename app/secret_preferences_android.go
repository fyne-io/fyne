//go:build !ci && android

package app

/*
#cgo LDFLAGS: -landroid -llog

#include <stdbool.h>
#include <stdlib.h>

bool secretEncrypt(uintptr_t jni_env, uintptr_t ctx, const void *in, int inLen, void **out, int *outLen);
bool secretDecrypt(uintptr_t jni_env, uintptr_t ctx, const void *in, int inLen, void **out, int *outLen);
*/
import "C"

import (
	"errors"
	"sync"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/driver/mobile/app"
)

var errKeystoreUnavailable = errors.New("android keystore is not available")

// newSecretStore returns the store used for secret preferences.
// The data is encrypted with an AES key that lives in the Android Keystore (hardware backed where
// the device supports it) and the ciphertext is kept in the app's private files directory.
// If the Java bridge is missing, for example an app packaged with an older toolchain,
// we fall back to the encrypted file store keyed from the password that the app provides.
func (a *fyneApp) newSecretStore(password []byte) secretStore {
	return &keystoreStore{storage: a.newPlainSecretStore(), fallback: func() secretStore {
		return a.newEncryptedSecretStore(password)
	}}
}

type keystoreStore struct {
	storage  secretStore
	fallback func() secretStore

	once      sync.Once
	available bool
	alt       secretStore
}

// check probes the Java side once so that we consistently use one format for the lifetime of the app.
func (k *keystoreStore) check() bool {
	k.once.Do(func() {
		_, err := keystoreCrypt(keystoreEncrypt, []byte{0})
		k.available = err == nil
		if !k.available {
			fyne.LogError("Android keystore unavailable, using encrypted file store", err)
			k.alt = k.fallback()
		}
	})
	return k.available
}

func (k *keystoreStore) load() ([]byte, error) {
	if !k.check() {
		return k.alt.load()
	}
	data, err := k.storage.load()
	if err != nil {
		return nil, err
	}
	return keystoreCrypt(keystoreDecrypt, data)
}

func (k *keystoreStore) save(data []byte) error {
	if !k.check() {
		return k.alt.save(data)
	}
	sealed, err := keystoreCrypt(keystoreEncrypt, data)
	if err != nil {
		return err
	}
	return k.storage.save(sealed)
}

type keystoreFunc func(env, ctx C.uintptr_t, in unsafe.Pointer, inLen C.int, out *unsafe.Pointer, outLen *C.int) bool

// cgo does not allow C functions to be used as values, so wrap the two bridge calls
func keystoreEncrypt(env, ctx C.uintptr_t, in unsafe.Pointer, inLen C.int, out *unsafe.Pointer, outLen *C.int) bool {
	return bool(C.secretEncrypt(env, ctx, in, inLen, out, outLen))
}

func keystoreDecrypt(env, ctx C.uintptr_t, in unsafe.Pointer, inLen C.int, out *unsafe.Pointer, outLen *C.int) bool {
	return bool(C.secretDecrypt(env, ctx, in, inLen, out, outLen))
}

func keystoreCrypt(fn keystoreFunc, in []byte) ([]byte, error) {
	var inPtr unsafe.Pointer
	if len(in) > 0 {
		inPtr = unsafe.Pointer(&in[0])
	}

	var result []byte
	ok := false
	app.RunOnJVM(func(vm, env, ctx uintptr) error {
		var out unsafe.Pointer
		var outLen C.int
		if !fn(C.uintptr_t(env), C.uintptr_t(ctx), inPtr, C.int(len(in)), &out, &outLen) {
			return nil
		}
		result = C.GoBytes(out, outLen)
		C.free(out)
		ok = true
		return nil
	})

	if !ok {
		return nil, errKeystoreUnavailable
	}
	return result, nil
}
