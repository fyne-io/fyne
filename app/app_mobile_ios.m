//go:build !ci && ios
// +build !ci,ios

#import <UIKit/UIKit.h>
#import <UserNotifications/UserNotifications.h>

static int notifyNum = 0;
extern void fallbackSend(char *cTitle, char *cBody);

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


