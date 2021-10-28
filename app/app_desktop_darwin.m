//go:build !ci && !ios
// +build !ci,!ios

extern void themeChanged();

#import <Foundation/Foundation.h>

bool isDarkMode() {
    NSString *style = [[NSUserDefaults standardUserDefaults] stringForKey:@"AppleInterfaceStyle"];
    return [@"Dark" isEqualToString:style];
}

void watchTheme() {
    [[NSDistributedNotificationCenter defaultCenter] addObserverForName:@"AppleInterfaceThemeChangedNotification" object:nil queue:nil
        usingBlock:^(NSNotification *note) {
        themeChanged(); // calls back into Go
    }];
}
