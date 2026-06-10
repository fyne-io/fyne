#import <Foundation/Foundation.h>
#import <AppKit/AppKit.h>

void setFullScreen(bool full, void *win) {
    NSWindow *window = (NSWindow*)win;

    NSUInteger masks = [window styleMask];
    bool isFull = masks & NSWindowStyleMaskFullScreen;
    if (isFull == full) {
        return;
    }

    if (full) {
        // macOS 26 ignores the window's current screen for toggleFullScreen, so re-seat
        // the frame on that screen explicitly to force fullscreen on the right display.
        NSScreen *targetScreen = [window screen];
        if (targetScreen != nil) {
            NSRect frame = [window frame];
            NSRect screen = [targetScreen frame];
            frame.origin.x = screen.origin.x + (screen.size.width - frame.size.width) / 2;
            frame.origin.y = screen.origin.y + (screen.size.height - frame.size.height) / 2;
            [window setFrame:frame display:YES];
        }
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
        // pick a screen that is not the one the window is already on, so launching
        // on a secondary display still sends "secondary fullscreen" to the other display.
        NSScreen *currentScreen = [window screen];
        NSScreen *targetScreen = nil;
        for (NSScreen *screen in [NSScreen screens]) {
            if (screen != currentScreen) {
                targetScreen = screen;
                break;
            }
        }
        if (targetScreen != nil) {
            NSRect frame = [window frame];
            NSRect screen = [targetScreen frame];
            frame.origin.x = screen.origin.x + (screen.size.width - frame.size.width) / 2;
            frame.origin.y = screen.origin.y + (screen.size.height - frame.size.height) / 2;
            [window setFrame:frame display:YES];
        }
    }

    [window toggleFullScreen:NULL];
}
