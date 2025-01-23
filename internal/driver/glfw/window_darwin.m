#import <Foundation/Foundation.h>
#import <AppKit/AppKit.h>

void setFullScreen(bool full, void *win) {
    NSWindow *window = (NSWindow*)win;

    NSUInteger masks = [window styleMask];
    bool isFull = masks & NSFullScreenWindowMask;
    if (isFull == full) {
        return;
    }

    [window toggleFullScreen:NULL];
}
