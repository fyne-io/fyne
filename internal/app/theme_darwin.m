//go:build !ci && !ios

#import <Foundation/Foundation.h>

bool isDarkMode() {
    NSString *style = [[NSUserDefaults standardUserDefaults] stringForKey:@"AppleInterfaceStyle"];
    return [@"Dark" isEqualToString:style];
}
