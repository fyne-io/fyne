package glfw

/*
#import <stdbool.h>

void setDisableDisplaySleep(bool);
*/
import "C"

func setDisableScreenBlank(disable bool) {
	C.setDisableDisplaySleep(C.bool(disable))
}
