//go:build ios

#import <UIKit/UIKit.h>
#import <MobileCoreServices/MobileCoreServices.h>

void setClipboardContent(char *content) {
    NSString *value = [NSString stringWithUTF8String:content];
    [[UIPasteboard generalPasteboard] setString:value];
}

char *getClipboardContent() {
    NSString *str = [[UIPasteboard generalPasteboard] string];

    return (char *) [str UTF8String];
}
