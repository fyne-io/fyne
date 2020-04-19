// +build darwin
// +build arm arm64

#import <Foundation/Foundation.h>

void log_wrap(const char *logStr) {
    NSLog(@"GoLog: %s", logStr);
}
