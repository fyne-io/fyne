package gl

import "unsafe"

// toLEByteOrder returns the byte representation of float32 values in little endian byte order.
func toLEByteOrder(values ...float32) []byte {
	b := make([]byte, 4*len(values))
	if len(values) > 0 {
		src := unsafe.Slice((*byte)(unsafe.Pointer(unsafe.SliceData(values))), len(b)) //gosec:disable G103 -- float bits are copied into a fresh byte slice
		copy(b, src)
	}
	return b
}
