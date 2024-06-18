//go:build !ios && !wasm && !test_web_driver

#import <Foundation/Foundation.h>

bool isDarkMode() {
    NSString *style = [[NSUserDefaults standardUserDefaults] stringForKey:@"AppleInterfaceStyle"];
    return [@"Dark" isEqualToString:style];
}
