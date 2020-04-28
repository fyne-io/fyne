// +build !no_native_menus

#import <Foundation/Foundation.h>
#import <AppKit/AppKit.h>

extern void menuCallback(int);

@interface FyneMenuHandler : NSObject {
}
@end

@implementation FyneMenuHandler
+ (void) tapped:(NSMenuItem*) item {
    menuCallback([item tag]);
}
@end

// forward declaration … we want methods to be ordered alphabetically
NSMenu* nativeMainMenu();

void assignDarwinSubmenu(const void* i, const void* m) {
    NSMenu* menu = (NSMenu*)m; // this menu is created in the createDarwinMenu() function
    NSMenuItem *item = (NSMenuItem*)i;
    [item setSubmenu:menu]; // this retains the menu
    [menu release]; // release the menu
}

void completeDarwinMenu(const void* m, bool prepend) {
    NSMenu* main = nativeMainMenu();
    NSMenuItem* top = [[NSMenuItem alloc] initWithTitle:@"" action:nil keyEquivalent:@""];
    if (prepend) {
        [main insertItem:top atIndex:1];
    } else {
        [main addItem:top];
    }
    assignDarwinSubmenu(top, m);
}

const void* createDarwinMenu(const char* label) {
    return (void*)[[NSMenu alloc] initWithTitle:[NSString stringWithUTF8String:label]];
}

const void* darwinAppMenu() {
    return [[nativeMainMenu() itemAtIndex:0] submenu];
}

const void* insertDarwinMenuItem(const void* m, const char* label, int id, int index, bool isSeparator) {
    NSMenu* menu = (NSMenu*)m;
    NSMenuItem* item;

    if (isSeparator) {
        item = [NSMenuItem separatorItem];
    } else {
        item = [[NSMenuItem alloc]
            initWithTitle:[NSString stringWithUTF8String:label]
            action:@selector(tapped:)
            keyEquivalent:@""];
        [item setTarget:[FyneMenuHandler class]];
        [item setTag:id];
    }

    if (index > -1) {
        [menu insertItem:item atIndex:index];
    } else {
        [menu addItem:item];
    }
    [item release]; // retained by the menu
    return item;
}

NSMenu* nativeMainMenu() {
    NSApplication* app = [NSApplication sharedApplication];
    return [app mainMenu];
}
