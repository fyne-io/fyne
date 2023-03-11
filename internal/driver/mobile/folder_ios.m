//go:build ios
// +build ios

#import <Foundation/Foundation.h>

#import <stdbool.h>

NSArray *listForURL(const char* urlCstr) {
    NSString *urlStr = [NSString stringWithUTF8String:urlCstr];
    NSURL *url = [NSURL URLWithString:urlStr];

    return [[NSFileManager defaultManager] contentsOfDirectoryAtURL:url includingPropertiesForKeys:nil options:nil error:nil];
}

bool iosCanList(const char* url) {
    return listForURL(url) != nil;
}

bool iosCreateListable(const char* urlCstr) {
    NSString *urlStr = [NSString stringWithUTF8String:urlCstr];
    NSURL *url = [NSURL URLWithString:urlStr];

    return [[NSFileManager defaultManager] createDirectoryAtURL:url withIntermediateDirectories:YES attributes:nil error:nil];
}

char* iosList(const char* url) {
    NSArray *children = listForURL(url);

    return [[children componentsJoinedByString:@"|"] UTF8String];
}
