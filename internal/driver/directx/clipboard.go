//go:build windows && directx

package directx

import (
	"unsafe"

	"golang.org/x/sys/windows"

	"fyne.io/fyne/v2"
)

// Declare conformity with Clipboard.
var _ fyne.Clipboard = clipboard{}

// clipboard accesses the Windows clipboard through the standard
// Open/Get-or-Set/Close sequence. It holds no state, so instances are free.
type clipboard struct{}

// NewClipboard returns the system clipboard.
func NewClipboard() fyne.Clipboard {
	return clipboard{}
}

func (clipboard) Content() string {
	var content string
	runOnMain(func() {
		if r, _, _ := procOpenClipboard.Call(0); r == 0 {
			return
		}
		defer procCloseClipboard.Call()

		h, _, _ := procGetClipboardData.Call(cfUnicodeText)
		if h == 0 {
			return
		}
		ptr, _, _ := procGlobalLock.Call(h)
		if ptr == 0 {
			return
		}
		defer procGlobalUnlock.Call(h)

		content = windows.UTF16PtrToString((*uint16)(pointerFromAddr(ptr)))
	})
	return content
}

func (clipboard) SetContent(content string) {
	runOnMain(func() {
		text, err := windows.UTF16FromString(content)
		if err != nil {
			fyne.LogError("directx: encoding clipboard text", err)
			return
		}

		if r, _, _ := procOpenClipboard.Call(0); r == 0 {
			return
		}
		defer procCloseClipboard.Call()
		procEmptyClipboard.Call()

		size := uintptr(len(text) * int(unsafe.Sizeof(uint16(0))))
		h, _, _ := procGlobalAlloc.Call(gmemMoveable, size)
		if h == 0 {
			return
		}

		ptr, _, _ := procGlobalLock.Call(h)
		if ptr == 0 {
			procGlobalFree.Call(h)
			return
		}
		copy(unsafe.Slice((*uint16)(pointerFromAddr(ptr)), len(text)), text)
		procGlobalUnlock.Call(h)

		// On success the clipboard owns the handle, so it must not be freed here.
		if r, _, _ := procSetClipboardData.Call(cfUnicodeText, h); r == 0 {
			procGlobalFree.Call(h)
		}
	})
}
