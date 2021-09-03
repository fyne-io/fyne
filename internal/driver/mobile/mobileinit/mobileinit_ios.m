//go:build ios
// +build ios

#import <Foundation/Foundation.h>

void log_wrap(const char *logStr) {
    NSLog(@"Fyne: %s", logStr);
}
