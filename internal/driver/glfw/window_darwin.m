#import <Foundation/Foundation.h>
#import <AppKit/AppKit.h>

void setFullScreen(bool full, void *win) {
    NSWindow *window = (NSWindow*)win;

    NSUInteger masks = [window styleMask];
    bool isFull = masks & NSWindowStyleMaskFullScreen;
    if (isFull == full) {
        return;
    }

    [window toggleFullScreen:NULL];
}

void setFullScreenSecondary(bool full, void *win) {
    NSWindow *window = (NSWindow*)win;

    NSUInteger masks = [window styleMask];
    bool isFull = (masks & NSWindowStyleMaskFullScreen) != 0;
    if (isFull == full) {
        return;
    }

    if (full) {
        NSScreen *targetScreen = nil;
        for (NSScreen *screen in [NSScreen screens]) {
            if (screen != [NSScreen mainScreen]) {
                targetScreen = screen;
                break;
            }
        }
        if (targetScreen != nil) {
            NSRect frame = [window frame];
            NSRect screen = [targetScreen frame];
            frame.origin = NSMakePoint(frame.origin.x + screen.origin.x, frame.origin.y + screen.origin.y);
            [window setFrame:frame display:YES];
        }
    }

    [window toggleFullScreen:NULL];
}
