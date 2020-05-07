// +build !ci
// +build !ios

#import <Foundation/Foundation.h>

@interface FyneUserNotificationCenterDelegate : NSObject<NSUserNotificationCenterDelegate>

- (BOOL)userNotificationCenter:(NSUserNotificationCenter*)center
    shouldPresentNotification:(NSUserNotification*)notification;

@end

@implementation FyneUserNotificationCenterDelegate

- (BOOL)userNotificationCenter:(NSUserNotificationCenter*)center
    shouldPresentNotification:(NSUserNotification*)notification
{
    return YES;
}

@end

static int notifyNum = 0;

void sendNSUserNotification(const char *, const char *);

bool isBundled() {
    return [[NSBundle mainBundle] bundleIdentifier] != nil;
}

void sendNotification(const char *title, const char *body) {
    NSUserNotificationCenter *center = [NSUserNotificationCenter defaultUserNotificationCenter];
    if (center.delegate == nil) {
        center.delegate = [[FyneUserNotificationCenterDelegate new] autorelease];
    }
    NSUserNotification *notification = [[NSUserNotification new] autorelease];
    notification.title = [NSString stringWithUTF8String:title];
    notification.informativeText = [NSString stringWithUTF8String:body];
    notification.identifier = [NSString stringWithFormat:@"%@-fyne-notify-%d", [[NSBundle mainBundle] bundleIdentifier], ++notifyNum];
    [center scheduleNotification:notification];
}
