//go:build !ci && ios

#import <UIKit/UIKit.h>

char *documentsPath() {
    NSArray *paths = NSSearchPathForDirectoriesInDomains(NSDocumentDirectory, NSUserDomainMask, YES);
    NSString *path = paths.firstObject;
    return [path UTF8String];
}