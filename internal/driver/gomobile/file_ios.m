// +build ios

#import <Foundation/Foundation.h>

#import <stdbool.h>

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
