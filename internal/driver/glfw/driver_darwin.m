//go:build darwin

#import <AppKit/NSEvent.h>
#import <IOKit/pwr_mgt/IOPMLib.h>

IOPMAssertionID currentDisableID;

void setDisableDisplaySleep(BOOL disable) {
    if (!disable) {
        if (currentDisableID == 0) {
            return;
        }

        IOPMAssertionRelease(currentDisableID);
        currentDisableID = 0;
        return;
    }

    if (currentDisableID != 0) {
        return;
    }
    IOPMAssertionID assertionID;
    IOReturn success = IOPMAssertionCreateWithName(kIOPMAssertionTypeNoDisplaySleep,
        kIOPMAssertionLevelOn, (CFStringRef)@"App disabled screensaver", &assertionID);

    if (success == kIOReturnSuccess) {
        currentDisableID = assertionID;
    }
}

// https://developer.apple.com/documentation/appkit/nsevent/1528384-doubleclickinterval?language=objc
double doubleClickInterval() {
    return [NSEvent doubleClickInterval];
}
