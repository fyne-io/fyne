//go:build !ci && ios
// +build !ci,ios

#import <UIKit/UIKit.h>

void openURL(char *urlStr) {
    UIApplication *app = [UIApplication sharedApplication];
    NSURL *url = [NSURL URLWithString:[NSString stringWithUTF8String:urlStr]];
    [app openURL:url options:@{} completionHandler:nil];
}

char *documentsPath() {
    NSArray *paths = NSSearchPathForDirectoriesInDomains(NSDocumentDirectory, NSUserDomainMask, YES);
    NSString *path = paths.firstObject;
    return [path UTF8String];
}
