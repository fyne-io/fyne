// +build !ci

// +build ios

#import <UIKit/UIKit.h>
#import <UserNotifications/UserNotifications.h>

void openURL(char *urlStr) {
    UIApplication *app = [UIApplication sharedApplication];
    NSURL *url = [NSURL URLWithString:[NSString stringWithUTF8String:urlStr]];
    [app openURL:url options:@{} completionHandler:nil];
}

static int notifyNum = 0;

void doSendNotification(UNUserNotificationCenter *center, NSString *title, NSString *body) {
    UNMutableNotificationContent *content = [UNMutableNotificationContent new];
    content.title = title;
    content.body = body;

    notifyNum++;
    NSString *identifier = [NSString stringWithFormat:@"fyne-notify-%d", notifyNum];
    UNNotificationRequest *request = [UNNotificationRequest requestWithIdentifier:identifier
        content:content trigger:nil];

    [center addNotificationRequest:request withCompletionHandler:^(NSError * _Nullable error) {
        if (error != nil) {
            NSLog(@"Could not send notification: %@", error);
        }
    }];
}

void sendNotification(char *title, char *body) {
    UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
    NSString *titleStr = [NSString stringWithUTF8String:title];
    NSString *bodyStr = [NSString stringWithUTF8String:body];

    UNAuthorizationOptions options = UNAuthorizationOptionAlert;
    [center requestAuthorizationWithOptions:options
        completionHandler:^(BOOL granted, NSError *_Nullable error) {
            if (!granted) {
                NSLog(@"Unable to get permission to send notifications");
            } else {
                doSendNotification(center, titleStr, bodyStr);
            }
        }];
}
