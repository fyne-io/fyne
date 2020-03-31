// +build !no_native_menus

#import <Foundation/Foundation.h>
#import <AppKit/AppKit.h>

extern void menuCallback(int);

@interface FyneMenuItem : NSObject {

}
@end

@implementation FyneMenuItem
- (void) tapped:(id) sender {
    menuCallback([sender tag]);
}
@end

NSMenu* nativeMainMenu() {
    NSApplication* app = [NSApplication sharedApplication];
    return [app mainMenu];
}

const void* darwinAppMenu() {
    return [[nativeMainMenu() itemAtIndex:0] submenu];
}

const void* createDarwinMenu(const char* label) {
    return (void*)[[NSMenu alloc] initWithTitle:[NSString stringWithUTF8String:label]];
}

void insertDarwinMenuItem(const void* m, const char* label, int id, int index, bool isSeparator) {
    NSMenu* menu = (NSMenu*)m;
    NSMenuItem* item;

    if (isSeparator) {
        item = [NSMenuItem separatorItem];
    } else {
        // NSMenuItem's target is a weak reference, therefore we must not release it.
        // TODO: Keep a reference in Go and release it when the menu of a window changes or the window is released.
        FyneMenuItem* tapper = [[FyneMenuItem alloc] init];
        item = [[NSMenuItem alloc]
            initWithTitle:[NSString stringWithUTF8String:label]
            action:@selector(tapped:)
            keyEquivalent:@""];
        [item setTarget:tapper];
        [item setTag:id];
    }

    if (index > -1) {
        [menu insertItem:item atIndex:index];
    } else {
        [menu addItem:item];
    }
    [item release];
}

void completeDarwinMenu(const void* m, bool prepend) {
    NSMenu* menu = (NSMenu*)m;
    NSMenu* main = nativeMainMenu();
    NSMenuItem* top = [[NSMenuItem alloc] initWithTitle:@"" action:nil keyEquivalent:@""];
    if (prepend) {
        [main insertItem:top atIndex:1];
    } else {
        [main addItem:top];
    }
    [main setSubmenu:menu forItem:top];
    [menu release]; // release the menu created in createDarwinMenu() function
    menu = Nil;
}
