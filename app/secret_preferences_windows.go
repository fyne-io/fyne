//go:build !ci && !wasm && !test_web_driver && !tinygo

package app

import (
	"errors"
	"syscall"
	"unsafe"
)

var (
	crypt32                = syscall.NewLazyDLL("crypt32.dll")
	kernel32               = syscall.NewLazyDLL("kernel32.dll")
	procCryptProtectData   = crypt32.NewProc("CryptProtectData")
	procCryptUnprotectData = crypt32.NewProc("CryptUnprotectData")
	procLocalFree          = kernel32.NewProc("LocalFree")
)

const cryptProtectUIForbidden = 0x1

// dataBlob mirrors the Win32 DATA_BLOB structure used by the Data Protection API.
type dataBlob struct {
	cbData uint32
	pbData *byte
}

func newDataBlob(data []byte) dataBlob {
	if len(data) == 0 {
		return dataBlob{}
	}
	return dataBlob{cbData: uint32(len(data)), pbData: &data[0]}
}

// newSecretStore returns the store used for secret preferences.
// On Windows the data is protected with DPAPI, which ties it to the current Windows user account,
// so the password is not needed.
func (a *fyneApp) newSecretStore([]byte) secretStore {
	return &dpapiStore{storage: a.newPlainSecretStore()}
}

type dpapiStore struct {
	storage secretStore
}

func (d *dpapiStore) load() ([]byte, error) {
	data, err := d.storage.load()
	if err != nil {
		return nil, err
	}
	return dpapiCall(procCryptUnprotectData, data)
}

func (d *dpapiStore) save(data []byte) error {
	protected, err := dpapiCall(procCryptProtectData, data)
	if err != nil {
		return err
	}
	return d.storage.save(protected)
}

// dpapiCall runs CryptProtectData or CryptUnprotectData, which share a signature:
// (DATA_BLOB *in, LPCWSTR desc, DATA_BLOB *entropy, void *reserved, PROMPTSTRUCT *prompt, DWORD flags, DATA_BLOB *out)
func dpapiCall(proc *syscall.LazyProc, in []byte) ([]byte, error) {
	if len(in) == 0 {
		return nil, errors.New("no data to protect")
	}
	inBlob := newDataBlob(in)
	var outBlob dataBlob

	ret, _, err := proc.Call(
		uintptr(unsafe.Pointer(&inBlob)),
		0,
		0,
		0,
		0,
		cryptProtectUIForbidden,
		uintptr(unsafe.Pointer(&outBlob)),
	)
	if ret == 0 {
		return nil, err
	}
	defer procLocalFree.Call(uintptr(unsafe.Pointer(outBlob.pbData)))

	out := make([]byte, outBlob.cbData)
	copy(out, unsafe.Slice(outBlob.pbData, outBlob.cbData))
	return out, nil
}
