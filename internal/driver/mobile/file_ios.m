//go:build ios
// +build ios

#import <Foundation/Foundation.h>

#import <stdbool.h>

bool iosExistsPath(const char* path) {
    NSString *pathStr = [NSString stringWithUTF8String:path];
    return [[NSFileManager defaultManager] fileExistsAtPath:pathStr];
}

void* iosParseUrl(const char* url) {
    NSString *urlStr = [NSString stringWithUTF8String:url];
    return [NSURL URLWithString:urlStr];
}

const void* iosReadFromURL(void* urlPtr, int* len)  {
    NSURL* url = (NSURL*)urlPtr;
    NSData* data = [NSData dataWithContentsOfURL:url];

    *len = data.length;
    return data.bytes;
}

const int iosWriteToURL(void* urlPtr, const void* bytes, int len) {
    NSURL* url = (NSURL*)urlPtr;
    NSData *data = [NSData dataWithBytes:bytes length:len];
    BOOL ok = [data writeToURL:url atomically:YES];

    if (!ok) {
        return 0;
    }
    return data.length;
}
