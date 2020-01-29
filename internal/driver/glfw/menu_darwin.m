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

NSMenu* menu;

void createDarwinMenu(const char* label) {
    menu = [[NSMenu alloc] initWithTitle:[NSString stringWithUTF8String:label]];

    // the menu is released at the end of the completeDarwinMenu() function
}

void addDarwinMenuItem(const char* label, int id) {
    FyneMenuItem* tapper = [[FyneMenuItem alloc] init]; // we cannot release this or the menu does not function...

    NSMenuItem* item = [[NSMenuItem alloc] initWithTitle:[NSString stringWithUTF8String:label]
        action:@selector(tapped:) keyEquivalent:@""];
    [item setTarget:tapper];
    [item setTag:id];

    [menu addItem:item];
    [item release];
}

void completeDarwinMenu() {
    NSMenu* main = nativeMainMenu();

    NSMenuItem* top = [[NSMenuItem alloc] initWithTitle:@"" action:nil keyEquivalent:@""];
    [main addItem:top];
    [main setSubmenu:menu forItem:top];
    [menu release]; // release the menu created in createDarwinMenu() function
    menu = Nil;
}
