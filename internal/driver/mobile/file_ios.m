//go:build ios

#import <Foundation/Foundation.h>

#import <stdbool.h>

void iosDeletePath(const char* path) {
    NSString *pathStr = [NSString stringWithUTF8String:path];
    [[NSFileManager defaultManager] removeItemAtPath:pathStr error:nil];
}

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

const void* iosOpenFileWriter(void* urlPtr, bool truncate) {
    NSURL* url = (NSURL*)urlPtr;

    if (truncate || ![[NSFileManager defaultManager] fileExistsAtPath:url.path]) {
        [[NSFileManager defaultManager] createFileAtPath:url.path contents:nil attributes:nil];
    }

    NSError *err = nil;
    // TODO handle error
    NSFileHandle* handle = [NSFileHandle fileHandleForWritingToURL:url error:&err];

    if (!truncate) {
        [handle seekToEndOfFile];
    }

    return handle;
}

void iosCloseFileWriter(void* handlePtr) {
    NSFileHandle* handle = (NSFileHandle*)handlePtr;

    [handle synchronizeFile];
    [handle closeFile];
}


const int iosWriteToFile(void* handlePtr, const void* bytes, int len) {
    NSFileHandle* handle = (NSFileHandle*)handlePtr;
    NSData *data = [NSData dataWithBytes:bytes length:len];

    [handle writeData:data];
    return data.length;
}
