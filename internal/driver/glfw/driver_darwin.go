//go:build darwin

package glfw

/*
#import <stdbool.h>

void setDisableDisplaySleep(bool);
double doubleClickInterval();
*/
import "C"
import "time"

func setDisableScreenBlank(disable bool) {
	C.setDisableDisplaySleep(C.bool(disable))
}

func (d *gLDriver) DoubleTapDelay() time.Duration {
	millis := int64(float64(C.doubleClickInterval()) * 1000)
	return time.Duration(millis) * time.Millisecond
}
