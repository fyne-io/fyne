// +build ios

package gomobile

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework UIKit -framework MobileCoreServices

#include <stdlib.h>

void setClipboardContent(char *content);
char *getClipboardContent();
*/
import "C"
import "unsafe"

// Content returns the clipboard content for iOS
func (c *mobileClipboard) Content() string {
	content := C.getClipboardContent()

	return C.GoString(content)
}

// SetContent sets the clipboard content for iOS
func (c *mobileClipboard) SetContent(content string) {
	contentStr := C.CString(content)
	defer C.free(unsafe.Pointer(contentStr))

	C.setClipboardContent(contentStr)
}
