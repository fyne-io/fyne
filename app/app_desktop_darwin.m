//go:build !ci && !ios

extern void themeChanged();

#import <Foundation/Foundation.h>

void watchTheme() {
    [[NSDistributedNotificationCenter defaultCenter] addObserverForName:@"AppleInterfaceThemeChangedNotification" object:nil queue:nil
        usingBlock:^(NSNotification *note) {
        themeChanged(); // calls back into Go
    }];
}
