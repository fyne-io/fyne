// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux
// +build linux

package unix

import "unsafe"

// Helpers for dealing with ifreq since it contains a union and thus requires a
// lot of unsafe.Pointer casts to use properly.

// newIfreq creates an ifreq with the input network interface name after
// validating the name does not exceed IFNAMSIZ-1 (trailing NULL required)
// bytes.
func newIfreq(name string) (*ifreq, error) {
	// Leave room for terminating NULL byte.
	if len(name) >= IFNAMSIZ {
		return nil, EINVAL
	}

	var ifr ifreq
	copy(ifr.Ifrn[:], name)

	return &ifr, nil
}

// An ifreqData is an ifreq but with a typed unsafe.Pointer field for data in
// the union. This is required in order to comply with the unsafe.Pointer rules
// since the "pointer-ness" of data would not be preserved if it were cast into
// the byte array of a raw ifreq.
type ifreqData struct {
	name [IFNAMSIZ]byte
	data unsafe.Pointer
	// Pad to the same size as ifreq.
	_ [len(ifreq{}.Ifru) - SizeofPtr]byte
}

// SetData produces an ifreqData with the pointer p set for ioctls which require
// arbitrary pointer data.
func (ifr ifreq) SetData(p unsafe.Pointer) ifreqData {
	return ifreqData{
		name: ifr.Ifrn,
		data: p,
	}
}
