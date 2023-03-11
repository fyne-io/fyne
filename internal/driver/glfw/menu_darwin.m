//go:build !no_native_menus
// +build !no_native_menus

#import <Foundation/Foundation.h>
#import <AppKit/AppKit.h>

const int menuTagMin = 5000;

#if __MAC_OS_X_VERSION_MAX_ALLOWED >= 101400
NSControlStateValue STATE_ON = NSControlStateValueOn;
NSControlStateValue STATE_OFF = NSControlStateValueOff;
#else
NSCellStateValue STATE_ON = NSOnState;
NSCellStateValue STATE_OFF = NSOffState;
#endif


extern void menuCallback(int);
extern BOOL menuEnabled(int);
extern BOOL menuChecked(int);
extern void exceptionCallback(const char*);

@interface FyneMenuHandler : NSObject {
}
@end

@implementation FyneMenuHandler
+ (void) tapped:(NSMenuItem*) item {
    menuCallback([item tag]-menuTagMin);
}
+ (BOOL) validateMenuItem:(NSMenuItem*) item {
    BOOL checked = menuChecked([item tag]-menuTagMin);
    if (checked) {
        [item setState:STATE_ON];
    } else {
        [item setState:STATE_OFF];
    }

    return menuEnabled([item tag]-menuTagMin);
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
    [top setTag:menuTagMin];
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

void getTextColorRGBA(int* r, int* g, int* b, int* a) {
  CGFloat fr, fg, fb, fa;
  NSColor *c = [[NSColor selectedMenuItemTextColor] colorUsingColorSpace: [NSColorSpace sRGBColorSpace]];
  [c getRed: &fr green: &fg blue: &fb alpha: &fa];
  *r = fr*255.0;
  *g = fg*255.0;
  *b = fb*255.0;
  *a = fa*255.0;
}

void handleException(const char* m, id e) {
    exceptionCallback([[NSString stringWithFormat:@"%s failed: %@", m, e] UTF8String]);
}

const void* insertDarwinMenuItem(const void* m, const char* label, const char* keyEquivalent, unsigned int keyEquivalentModifierMask, int id, int index, bool isSeparator, const void *imageData, unsigned int imageDataLength) {
    NSMenu* menu = (NSMenu*)m;
    NSMenuItem* item;

    if (isSeparator) {
        item = [NSMenuItem separatorItem];
    } else {
        item = [[NSMenuItem alloc]
            initWithTitle:[NSString stringWithUTF8String:label]
            action:@selector(tapped:)
            keyEquivalent:[NSString stringWithUTF8String:keyEquivalent]];
        if (keyEquivalentModifierMask) {
            [item setKeyEquivalentModifierMask: keyEquivalentModifierMask];
        }
        [item setTarget:[FyneMenuHandler class]];
        [item setTag:id+menuTagMin];
        if (imageData) {
        char *x = (char *)imageData;
            NSData *data = [[NSData alloc] initWithBytes: imageData length: imageDataLength];
            NSImage *image = [[NSImage alloc] initWithData: data];
            [item setImage: image];
            [data release];
            [image release];
        }
    }

    if (index > -1) {
        [menu insertItem:item atIndex:index];
    } else {
        [menu addItem:item];
    }
    [item release]; // retained by the menu
    return item;
}

int menuFontSize() {
  return ceil([[NSFont menuFontOfSize: 0] pointSize]);
}

NSMenu* nativeMainMenu() {
    NSApplication* app = [NSApplication sharedApplication];
    return [app mainMenu];
}

void resetDarwinMenu() {
    NSMenu *root = nativeMainMenu();
    NSEnumerator *items = [[root itemArray] objectEnumerator];

    id object;
    while (object = [items nextObject]) {
        NSMenuItem *item = object;
        if ([item tag] < menuTagMin) {
            // check for inserted items (like Settings...)
            NSMenu *menu = [item submenu];
            NSEnumerator *subItems = [[menu itemArray] objectEnumerator];

            id sub;
            while (sub = [subItems nextObject]) {
                NSMenuItem *item = sub;
                if ([item tag] >= menuTagMin) {
                    [menu removeItem: item];
                }
            }

            continue;
        }
        [root removeItem: item];
    }
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

const char* test_NSMenuItem_keyEquivalent(const void *i) {
    NSMenuItem* item = (NSMenuItem*)i;
    return [[item keyEquivalent] UTF8String];
}

unsigned long test_NSMenuItem_keyEquivalentModifierMask(const void *i) {
    NSMenuItem* item = (NSMenuItem*)i;
    return [item keyEquivalentModifierMask];
}

const void* test_NSMenuItem_submenu(const void* i) {
    NSMenuItem* item = (NSMenuItem*)i;
    return [item submenu];
}

const char* test_NSMenuItem_title(const void* i) {
    NSMenuItem* item = (NSMenuItem*)i;
    return [[item title] UTF8String];
}
