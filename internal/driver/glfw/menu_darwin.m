// +build !no_native_menus

#import <Foundation/Foundation.h>
#import <AppKit/AppKit.h>

extern void menu_callback(int);

@interface FyneMenuItem : NSObject {

}
@end

@implementation FyneMenuItem
- (void) tapped:(id) sender {
    menu_callback([sender tag]);
}
@end

NSMenu* nativeMainMenu() {
    NSApplication* app = [NSApplication sharedApplication];
    return [app mainMenu];
}

const void* createDarwinMenu(const char* label) {
    return (void*)[[NSMenu alloc] initWithTitle:[NSString stringWithUTF8String:label]];
}

void insertDarwinMenuItem(const void* m, const char* label, int id) {
    NSMenu* menu = (NSMenu*)m;
    // NSMenuItem's target is a weak reference, therefore we must not release it.
    // TODO: Keep a reference in Go and release it when the menu of a window changes or the window is released.
    FyneMenuItem* tapper = [[FyneMenuItem alloc] init];
    NSMenuItem* item = [[NSMenuItem alloc]
        initWithTitle:[NSString stringWithUTF8String:label]
        action:@selector(tapped:)
        keyEquivalent:@""];
    [item setTarget:tapper];
    [item setTag:id];

    [menu addItem:item];
    [item release];
}

void completeDarwinMenu(const void* m) {
    NSMenu* menu = (NSMenu*)m;
    NSMenu* main = nativeMainMenu();
    NSMenuItem* top = [[NSMenuItem alloc] initWithTitle:@"" action:nil keyEquivalent:@""];
    [main addItem:top];
    [main setSubmenu:menu forItem:top];
    [menu release]; // release the menu created in createDarwinMenu() function
    menu = Nil;
}
