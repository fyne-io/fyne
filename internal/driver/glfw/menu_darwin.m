// +build !no_native_menus

#import <Foundation/Foundation.h>
#import <AppKit/AppKit.h>

extern void menuCallback(int);
extern void exceptionCallback(const char*);

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

void handleException(const char* m, id e) {
    exceptionCallback([[NSString stringWithFormat:@"%s failed: %@", m, e] UTF8String]);
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

const void* test_darwinMainMenu() {
    return nativeMainMenu();
}

const void* test_NSMenu_itemAtIndex(const void* m, NSInteger i) {
    NSMenu* menu = (NSMenu*)m;
    @try {
        return [menu itemAtIndex: i];
    } @catch(NSException* e) {
        handleException("test_NSMenu_itemAtIndex", e);
        return NULL;
    }
}

NSInteger test_NSMenu_numberOfItems(const void* m) {
    NSMenu* menu = (NSMenu*)m;
    return [menu numberOfItems];
}

void test_NSMenu_performActionForItemAtIndex(const void* m, NSInteger i) {
    NSMenu* menu = (NSMenu*)m;
    @try {
        // Using performActionForItemAtIndex: would be better but sadly it crashes.
        // We simulate the relevant effect for now.
        // [menu performActionForItemAtIndex:i];
        NSMenuItem* item = [menu itemAtIndex:i];
        [[item target] performSelector:[item action] withObject:item];
    } @catch(NSException* e) {
        handleException("test_NSMenu_performActionForItemAtIndex", e);
    }
}

void test_NSMenu_removeItemAtIndex(const void* m, NSInteger i) {
    NSMenu* menu = (NSMenu*)m;
    @try {
        [menu removeItemAtIndex: i];
    } @catch(NSException* e) {
        handleException("test_NSMenu_removeItemAtIndex", e);
    }
}

const char* test_NSMenu_title(const void* m) {
    NSMenu* menu = (NSMenu*)m;
    return [[menu title] UTF8String];
}

bool test_NSMenuItem_isSeparatorItem(const void* i) {
    NSMenuItem* item = (NSMenuItem*)i;
    return [item isSeparatorItem];
}

const void* test_NSMenuItem_submenu(const void* i) {
    NSMenuItem* item = (NSMenuItem*)i;
    return [item submenu];
}

const char* test_NSMenuItem_title(const void* i) {
    NSMenuItem* item = (NSMenuItem*)i;
    return [[item title] UTF8String];
}
