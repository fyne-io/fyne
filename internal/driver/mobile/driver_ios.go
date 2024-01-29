//go:build ios

package mobile

/*
#cgo darwin LDFLAGS: -framework UIKit
#import <Foundation/Foundation.h>

void disableIdleTimer(BOOL disabled);
*/
import "C"

func setDisableScreenBlank(disable bool) {
	C.disableIdleTimer(C.BOOL(disable))
}
