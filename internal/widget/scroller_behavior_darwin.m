//go:build darwin

#import <Foundation/Foundation.h>

int getScrollerPagingBehavior() {
    return [[NSUserDefaults standardUserDefaults] boolForKey:@"AppleScrollerPagingBehavior"];
}
