//go:build !ci && !ios && !wasm && !test_web_driver && !mobile

extern void themeChanged();

#import <Foundation/Foundation.h>

void watchTheme() {
    [[NSDistributedNotificationCenter defaultCenter] addObserverForName:@"AppleInterfaceThemeChangedNotification" object:nil queue:nil
        usingBlock:^(NSNotification *note) {
        themeChanged(); // calls back into Go
    }];
}
