//go:build ios

#import <Foundation/Foundation.h>
#import <UIKit/UIKit.h>

void disableIdleTimer(BOOL disabled) {
    @autoreleasepool {
        [[NSOperationQueue mainQueue] addOperationWithBlock:^ {
            [UIApplication sharedApplication].idleTimerDisabled = disabled;
        }];
    }
}
